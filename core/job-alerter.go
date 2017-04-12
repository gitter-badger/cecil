package core

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/notification"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
)

const AllAlertsSent = 2
const NoAlertsSent = 0

// AlerterJob polls the DB for leases that are about to expire, and notifes the owner of the imminent expiry
func (s *Service) AlerterJob() error {
	// find lease that expire in 24 hours
	// find owner
	// create links to extend and terminate lease
	// mark as num_times_allerted_about_expiry =+ 1
	// registed new lease's token_once
	// compose email with link to extend and terminate lease
	// send email

	var expiringLeases []models.Lease
	var expiringLeasesCount int64

	s.DB.Table("leases").
		Where("expires_at < ?",
			time.Now().UTC().Add(s.Config().Lease.FirstWarningBeforeExpiry),
		).
		Where("num_times_allerted_about_expiry < ? AND terminated_at IS NULL", AllAlertsSent).
		Not("approved_at IS NULL").
		Find(&expiringLeases).
		Count(&expiringLeasesCount)

	Logger.Info("AlerterJob(): Expiring leases", "count", expiringLeasesCount)

	// TODO: create ExpiringLeaseQueue and pass to it this task
ExpiringLeasesIterator:
	for _, expiringLease := range expiringLeases {

		switch expiringLease.NumTimesAllertedAboutExpiry {
		case 0:
			{
			}
		case 1:
			{
				if !expiringLease.ExpiresAt.Before(time.Now().UTC().Add(s.Config().Lease.SecondWarningBeforeExpiry)) {
					continue ExpiringLeasesIterator
				}
			}
		default:
			{
				continue ExpiringLeasesIterator
			}
		}

		Logger.Info("Expiring lease",
			"lease_id", expiringLease.ID,
			"group_type", expiringLease.GroupType.String(),
			"group_uid", expiringLease.GroupUID,
		)

		var owner models.Owner
		err := s.DB.Table("owners").Where(expiringLease.OwnerID).First(&owner).Error
		if err != nil {
			Logger.Error("error while fetching owner of expiring lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return err
			}
			return err
		}

		// these will be used to compose the urls and verify the requests
		tokenOnce := uuid.NewV4().String() // one-time token

		expiringLease.TokenOnce = tokenOnce
		expiringLease.NumTimesAllertedAboutExpiry++

		s.DB.Save(&expiringLease)

		// URL to extend lease
		extendURL, err := s.EmailActionGenerateSignedURL("extend", expiringLease.UUID, HashString(expiringLease.GroupUID), tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		// URL to terminate lease
		terminateURL, err := s.EmailActionGenerateSignedURL("terminate", expiringLease.UUID, HashString(expiringLease.GroupUID), tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		var AWSResourceID string

		AWSResourceID = expiringLease.GroupUID

		var emailValues = map[string]interface{}{
			"owner_email": owner.Email,

			"instance_created_at": expiringLease.CreatedAt.Format("2006-01-02 15:04:05 GMT"),
			"extend_by":           s.Config().Lease.Duration.String(),

			"termination_time": expiringLease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
			"lease_duration":   expiringLease.ExpiresAt.Sub(expiringLease.CreatedAt).String(),

			"lease_terminate_url": terminateURL,
			"lease_extend_url":    extendURL,
			"resource_region":     expiringLease.Region,
		}

		emailValues["lease_id"] = expiringLease.ID
		emailValues["group_type"] = expiringLease.GroupType.String()
		emailValues["group_uid"] = expiringLease.GroupUID

		newEmailBody, err := tools.CompileEmailTemplate(
			"expiring-lease.txt",
			emailValues,
		)
		if err != nil {
			return err
		}

		var newEmailSubject string
		newEmailSubject = fmt.Sprintf("Stack (%v) will expire soon", expiringLease.ID) // TODO: change email subject

		switch expiringLease.NumTimesAllertedAboutExpiry {
		case 1:
			newEmailSubject = fmt.Sprintf("%v %v", newEmailSubject, "(1st warning)")
		case 2:
			newEmailSubject = fmt.Sprintf("%v %v", newEmailSubject, "(final warning)")
		}

		s.Queues().NotifierQueue().PushTask(tasks.NotifierTask{
			AccountID: expiringLease.AccountID, // this will also trigger send to Slack
			To:        owner.Email,
			Subject:   newEmailSubject,
			BodyHTML:  newEmailBody,
			BodyText:  newEmailBody,
			NotificationMeta: notification.NotificationMeta{
				NotificationType: notification.InstanceWillExpire,
				LeaseUUID:        expiringLease.UUID,
				AWSResourceID:    AWSResourceID,
				//ResourceType:     expiringLease.ResourceType,
			},
		})

	}
	return nil
}
