package core

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// EventInjestorJob polls the SQS queue, verifies the message, and pushes it to the proper queue
func (s *Service) EventInjestorJob() error {
	// TODO: verify event origin (must be aws, not someone else)

	queueURL := s.SQSQueueURL()

	receiveMessageParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL), // Required
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(90), // should be higher, like 10 (seconds), the time to finish doing everything
		WaitTimeSeconds:     aws.Int64(3),
	}

	// Make sure there is a non-nil SQS
	if s.AWS.SQS == nil {
		Logger.Warn("EventInvestorJob", "SQS == nil, skipping")
		return nil
	}

	Logger.Info("EventInjestorJob(): Polling SQS", "queue", queueURL)
	receiveMessageResponse, err := s.AWS.SQS.ReceiveMessage(receiveMessageParams)

	if err != nil {
		return fmt.Errorf("EventInjestorJob() error: %v", err)
	}

	Logger.Info("SQSmessages",
		"count", len(receiveMessageResponse.Messages),
	)

	for messageIndex := range receiveMessageResponse.Messages {

		transmission, err := s.parseSQSTransmission(receiveMessageResponse.Messages[messageIndex], queueURL)

		if err != nil {
			if err == ErrorEnvelopeIsSubscriptionConfirmation {

				if err := transmission.ConfirmSQSSubscription(); err != nil {
					Logger.Warn("ConfirmSQSSubscription", "err", err)
					continue
				}

				if err := transmission.DeleteMessage(); err != nil {
					Logger.Warn("DeleteMessage", "err", err)
				}
				continue
			} else {
				Logger.Warn("Error parsing transmission", "err", err)

				err = transmission.DeleteMessage()
				if err != nil {
					Logger.Warn("DeleteMessage", "err", err)
				}
				continue
			}
		}

		Logger.Info("Parsed sqs message", "message", transmission.Message)

		if !transmission.TopicAndInstanceHaveSameOwner() {
			// the originating SNS topic and the instance have different owners (different AWS accounts)
			// TODO: notify cecil admin
			Logger.Warn("topicAWSID != instanceOriginatorID", "topicAWSID", transmission.Topic.AWSID, "instanceOriginatorID", transmission.Message.Account)
			// TODO: delete message
			continue
		}

		// consider only pending and terminated status messages; ignore the rest
		if !transmission.MessageIsRelevant() {
			Logger.Warn("Ignoring and removing message", "message.Detail.State", transmission.Message.Detail.State)
			err := transmission.DeleteMessage()
			if err != nil {
				Logger.Warn("DeleteMessage", "err", err)
			}
			continue // next message
		}

		// send transmission to NewLeaseQueue
		s.NewLeaseQueue.TaskQueue <- NewLeaseTask{
			Transmission: transmission,
		}

	}

	return nil
}

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

	var expiringLeases []Lease
	var expiringLeasesCount int64

	s.DB.Table("leases").
		Where("expires_at < ?",
			time.Now().UTC().Add(s.Config.Lease.FirstWarningBeforeExpiry),
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
				if !expiringLease.ExpiresAt.Before(time.Now().UTC().Add(s.Config.Lease.SecondWarningBeforeExpiry)) {
					continue ExpiringLeasesIterator
				}
			}
		default:
			{
				continue ExpiringLeasesIterator
			}
		}

		Logger.Info("Expiring lease",
			"leaseID", expiringLease.ID,
			"resourceType", expiringLease.ResourceType,
			"resourceID", expiringLease.ResourceID,
		)

		var owner Owner
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
		extendURL, err := s.EmailActionGenerateSignedURL("extend", expiringLease.UUID, expiringLease.ResourceID, tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		// URL to terminate lease
		terminateURL, err := s.EmailActionGenerateSignedURL("terminate", expiringLease.UUID, expiringLease.ResourceID, tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		var AWSResourceID string

		var instance InstanceResource
		if expiringLease.IsInstance() {
			raw, err := s.ResourceOf(&expiringLease)
			if err != nil {
				return err
			}
			instance = raw.(InstanceResource)
			AWSResourceID = instance.InstanceID
		}

		var stack StackResource
		if expiringLease.IsStack() {
			raw, err := s.ResourceOf(&expiringLease)
			if err != nil {
				return err
			}
			stack = raw.(StackResource)
			AWSResourceID = stack.StackID
		}

		var emailValues = map[string]interface{}{
			"owner_email": owner.Email,

			"instance_created_at": expiringLease.CreatedAt.Format("2006-01-02 15:04:05 GMT"),
			"extend_by":           s.Config.Lease.Duration.String(),

			"termination_time": expiringLease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
			"lease_duration":   expiringLease.ExpiresAt.Sub(expiringLease.CreatedAt).String(),

			"lease_terminate_url": terminateURL,
			"lease_extend_url":    extendURL,
			"resource_region":     expiringLease.Region,
		}

		if expiringLease.IsInstance() {
			emailValues["instance_id"] = instance.InstanceID
			emailValues["instance_type"] = instance.InstanceType
		}

		if expiringLease.IsStack() {
			emailValues["stack_id"] = stack.StackID
			emailValues["stack_name"] = stack.StackName
		}

		newEmailBody, err := CompileEmailTemplate(
			"expiring-lease.txt",
			emailValues,
		)
		if err != nil {
			return err
		}

		var newEmailSubject string
		if expiringLease.IsStack() {
			newEmailSubject = fmt.Sprintf("Stack (%v) will expire soon", stack.StackName)
		} else {
			newEmailSubject = fmt.Sprintf("Instance (%v) will expire soon", instance.InstanceID)
		}

		switch expiringLease.NumTimesAllertedAboutExpiry {
		case 1:
			newEmailSubject = fmt.Sprintf("%v %v", newEmailSubject, "(1st warning)")
		case 2:
			newEmailSubject = fmt.Sprintf("%v %v", newEmailSubject, "(final warning)")
		}

		s.NotifierQueue.TaskQueue <- NotifierTask{
			AccountID: expiringLease.AccountID, // this will also trigger send to Slack
			To:        owner.Email,
			Subject:   newEmailSubject,
			BodyHTML:  newEmailBody,
			BodyText:  newEmailBody,
			NotificationMeta: NotificationMeta{
				NotificationType: InstanceWillExpire,
				LeaseUUID:        expiringLease.UUID,
				AWSResourceID:    AWSResourceID,
				ResourceType:     expiringLease.ResourceType,
			},
		}
	}

	return nil
}

// SentencerJob polls the DB for expired leases and pushes them to the TerminatorQueue
func (s *Service) SentencerJob() error {

	var expiredLeases []Lease
	var expiredLeasesCount int64

	s.DB.Table("leases").
		Where("expires_at < ? AND terminated_at IS NULL", time.Now().UTC()).
		Or("approved_at IS NULL AND launched_at < ? AND terminated_at IS NULL", time.Now().UTC().Add(-s.Config.Lease.ApprovalTimeoutDuration)).
		Find(&expiredLeases).
		Count(&expiredLeasesCount)

	Logger.Info("SentencerJob(): Expired leases", "count", expiredLeasesCount)

	for _, expiredLease := range expiredLeases {
		Logger.Info("expired lease",
			"leaseID", expiredLease.ID,
			"resourceID", expiredLease.ResourceID,
			"resourceType", expiredLease.ResourceType,
		)
		s.TerminatorQueue.TaskQueue <- TerminatorTask{Lease: expiredLease}
	}

	return nil
}
