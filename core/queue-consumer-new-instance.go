package core

import (
	"errors"
	"fmt"

	"github.com/satori/go.uuid"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/notification"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
	"github.com/tleyden/cecil/transmission"
)

func ParseNewInstanceTask(t interface{}) (*transmission.Transmission, error) {
	if t == nil {
		return nil, errors.New("t is nil")
	}

	task, ok := t.(tasks.NewInstanceTask)
	if !ok {
		return nil, errors.New("t is not tasks.NewInstanceTask")
	}

	tr, ok := task.Transmission.(*transmission.Transmission)
	if !ok {
		return nil, errors.New("t.Transmission is not *transmission.Transmission")
	}

	return tr, nil
}

// NewInstanceQueueConsumer consumes NewInstanceTask from NewInstanceQueue
func (s *Service) NewInstanceQueueConsumer(t interface{}) error {

	tr, err := ParseNewInstanceTask(t)
	if err != nil {
		return err
	}

	Logger.Info(
		"NewInstanceQueueConsumer called",
		"transmission", tr,
	)
	defer Logger.Info(
		"NewInstanceQueueConsumer call finished",
		"transmission", tr,
	)

	groupUIDPtr, err := tr.DefineGroupUID()
	if err != nil {
		Logger.Error("error while DefineGroupUID", "err", err)
		return err
	}
	if groupUIDPtr == nil {
		Logger.Error("could not DefineGroupUID; groupUIDPtr is nil")
		return nil
	}

	groupUID := *groupUIDPtr
	Logger.Info("NewInstanceQueueConsumer::DefineGroupUID()", "groupUID", groupUID, "tr.GroupType", tr.GroupType.String())

	// check whether the group has a lease

	lease, err := tr.GroupHasAlreadyALease(groupUID)
	if err != nil {
		Logger.Error("error while GroupHasAlreadyALease", "err", err)
		return err
	}

	groupHasLease := lease != nil
	if groupHasLease {
		ins, err := tr.InstanceIsNew()
		if err != nil {
			Logger.Error("error while InstanceIsNew", "err", err)
			return err
		}
		instanceIsRegistered := ins != nil
		if instanceIsRegistered {
			if ins.GroupUID != groupUID {
				Logger.Warn("exirting instance is currently saved as belonging to different group", "current", ins.GroupUID, "new", groupUID)
			}
			Logger.Info("instance already registered", "existing instance in DB", ins)
			return nil
		} else {
			newInstance := models.Instance{
				LeaseID: lease.ID,

				AccountID:      tr.Cloudaccount.AccountID,
				CloudaccountID: tr.Cloudaccount.ID,
				AWSAccountID:   tr.Cloudaccount.AWSID,

				GroupUID:  groupUID,
				GroupType: tr.GroupType,

				InstanceID:       tr.InstanceID(),
				AvailabilityZone: tr.AvailabilityZone(),
				InstanceType:     tr.InstanceType(),
				Region:           tr.InstanceRegion,

				LaunchedAt: tr.InstanceLaunchTimeUTC(),
			}
			s.DB.Create(&newInstance)
			Logger.Info("saving new instance", "newInstance", newInstance)

			err = tr.DeleteMessage()
			if err != nil {
				Logger.Warn("DeleteMessage", "err", err)
			}
			return err
		}
	}

	// here we continue to figure out details about this new lease

	if !tr.InstanceHasTagOrKeyName() || !tr.ExternalOwnerIsWhitelisted() {
		// assign instance to admin, and send notification to admin
		// owner is not whitelisted: notify admin
		Logger.Info("Transmission doesn't have owner tag/keyname or owner is not whitelisted.")

		err := tr.SetAdminAsOwner()
		if err != nil {
			Logger.Warn("Error while setting admin as owner", "err", err)
			return err
		}

		//transmission.leaseDuration = s.Config.Lease.ApprovalTimeoutDuration
		var expiresAt = tr.LeaseExpiresAt()

		// these will be used to compose the urls and verify the requests
		leaseUUID := uuid.NewV4().String()
		instanceID := tr.InstanceID()
		tokenOnce := uuid.NewV4().String() // one-time token

		lease = &models.Lease{
			UUID:      leaseUUID,
			TokenOnce: tokenOnce,

			GroupUID:  groupUID,
			GroupType: tr.GroupType,

			OwnerID:        tr.Owner.ID,
			AccountID:      tr.Cloudaccount.AccountID,
			CloudaccountID: tr.Cloudaccount.ID,
			AWSAccountID:   tr.Cloudaccount.AWSID,
			Region:         tr.InstanceRegion,

			NumTimesAllertedAboutExpiry: AllAlertsSent, // this will prevent any other notification before the expiry

			LaunchedAt: tr.InstanceLaunchTimeUTC(),
			ExpiresAt:  expiresAt,
		}
		s.DB.Create(lease)
		Logger.Info("new lease created",
			"lease", lease,
		)

		newInstance := models.Instance{
			LeaseID: lease.ID,

			AccountID:      tr.Cloudaccount.AccountID,
			CloudaccountID: tr.Cloudaccount.ID,
			AWSAccountID:   tr.Cloudaccount.AWSID,

			GroupUID:  groupUID,
			GroupType: tr.GroupType,

			InstanceID:       tr.InstanceID(),
			AvailabilityZone: tr.AvailabilityZone(),
			InstanceType:     tr.InstanceType(),
			Region:           tr.InstanceRegion,

			LaunchedAt: tr.InstanceLaunchTimeUTC(),
		}
		s.DB.Create(&newInstance)
		Logger.Info("saving new instance", "newInstance", newInstance)

		var newEmailBody string

		// URL to approve lease
		Logger.Info(
			"Creating lease signature",
			"lease_uuid", leaseUUID,
			"instance_id", instanceID,
			"action", "approve",
			"token_once", tokenOnce,
		)
		approveURL, err := s.EmailActionGenerateSignedURL("approve", leaseUUID, HashString(groupUID), tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		// URL to terminate lease
		terminateURL, err := s.EmailActionGenerateSignedURL("terminate", leaseUUID, HashString(groupUID), tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		var emailValues = map[string]interface{}{
			"owner_email":     tr.Owner.Email,
			"instance_id":     tr.InstanceID(),
			"instance_type":   tr.InstanceType(),
			"resource_region": tr.InstanceRegion,

			"termination_time": expiresAt.Format("2006-01-02 15:04:05 GMT"),
			"lease_duration":   tr.LeaseExpiresAt().Sub(tr.InstanceLaunchTimeUTC()).String(),

			"lease_terminate_url": terminateURL,
			"lease_approve_url":   approveURL,
		}

		switch {
		case !tr.InstanceHasTagOrKeyName():
			newEmailBody, err = tools.CompileEmailTemplate(
				"new-lease-no-owner-tag.txt",
				emailValues,
			)
			if err != nil {
				return err
			}
			break

		case !tr.ExternalOwnerIsWhitelisted():
			newEmailBody, err = tools.CompileEmailTemplate(
				"new-lease-owner-tag-not-whitelisted.txt",
				emailValues,
			)
			if err != nil {
				return err
			}
		}

		var emailSubject string
		emailSubject = fmt.Sprintf("Instance (%v) needs attention", tr.GroupType.String())

		Logger.Info("Adding new NotifierTask")
		s.Queues().NotifierQueue().PushTask(tasks.NotifierTask{
			AccountID: tr.AdminAccount.ID, // this will also trigger send to Slack
			//To:       owner.Email,
			To:       tr.AdminAccount.Email,
			Subject:  emailSubject,
			BodyHTML: newEmailBody,
			BodyText: newEmailBody,
			NotificationMeta: notification.NotificationMeta{
				NotificationType: notification.InstanceNeedsAttention,
				LeaseUUID:        leaseUUID,
				AWSResourceID:    tr.AWSResourceID(),
				//ResourceType:     lease.ResourceType,
			},
		})

		Logger.Info("Delete SQS Message")
		err = tr.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	}

	if err := tr.SetExternalOwnerAsOwner(); err != nil {
		Logger.Warn("Error while setting external owner as owner", "err", err)
	}

	if tr.LeaseNeedsApproval() {
		// register new lease in DB
		// expiry: 1h
		// send confirmation to owner: confirmation link, and termination link
		Logger.Info("Lease needs approval")

		//transmission.leaseDuration = s.Config.Lease.ApprovalTimeoutDuration
		var expiresAt = tr.LeaseExpiresAt()

		// these will be used to compose the urls and verify the requests
		leaseUUID := uuid.NewV4().String()
		tokenOnce := uuid.NewV4().String() // one-time token

		lease = &models.Lease{
			UUID:      leaseUUID,
			TokenOnce: tokenOnce,

			GroupUID:  groupUID,
			GroupType: tr.GroupType,

			OwnerID:        tr.Owner.ID,
			AccountID:      tr.Cloudaccount.AccountID,
			CloudaccountID: tr.Cloudaccount.ID,
			AWSAccountID:   tr.Cloudaccount.AWSID,
			Region:         tr.InstanceRegion,

			NumTimesAllertedAboutExpiry: AllAlertsSent,

			LaunchedAt: tr.InstanceLaunchTimeUTC(),
			ExpiresAt:  expiresAt,
		}
		s.DB.Create(lease)
		Logger.Info("new lease created",
			"lease", lease,
		)

		newInstance := models.Instance{
			LeaseID: lease.ID,

			AccountID:      tr.Cloudaccount.AccountID,
			CloudaccountID: tr.Cloudaccount.ID,
			AWSAccountID:   tr.Cloudaccount.AWSID,

			GroupUID:  groupUID,
			GroupType: tr.GroupType,

			InstanceID:       tr.InstanceID(),
			AvailabilityZone: tr.AvailabilityZone(),
			InstanceType:     tr.InstanceType(),
			Region:           tr.InstanceRegion,

			LaunchedAt: tr.InstanceLaunchTimeUTC(),
		}
		s.DB.Create(&newInstance)
		Logger.Info("saving new instance", "newInstance", newInstance)

		// URL to approve lease
		approveURL, err := s.EmailActionGenerateSignedURL("approve", leaseUUID, HashString(groupUID), tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		// URL to terminate lease
		terminateURL, err := s.EmailActionGenerateSignedURL("terminate", leaseUUID, HashString(groupUID), tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		var emailValues = map[string]interface{}{
			"owner_email":        tr.Owner.Email,
			"n_of_active_leases": tr.ActiveLeaseCount,
			"instance_id":        tr.InstanceID(),
			"instance_type":      tr.InstanceType(),
			"resource_region":    tr.InstanceRegion,

			"termination_time": expiresAt.Format("2006-01-02 15:04:05 GMT"),
			"lease_duration":   tr.LeaseExpiresAt().Sub(tr.InstanceLaunchTimeUTC()).String(),

			"lease_approve_url":   approveURL,
			"lease_terminate_url": terminateURL,
		}

		newEmailBody, err := tools.CompileEmailTemplate(
			"new-lease-valid-owner-tag-needs-approval.txt",
			emailValues,
		)
		if err != nil {
			return err
		}

		var emailSubject string
		emailSubject = fmt.Sprintf("Instance (%v) needs approval", tr.GroupType.String())

		s.Queues().NotifierQueue().PushTask(tasks.NotifierTask{
			AccountID: tr.AdminAccount.ID, // this will also trigger send to Slack
			To:        tr.Owner.Email,
			Subject:   emailSubject,
			BodyHTML:  newEmailBody,
			BodyText:  newEmailBody,
			NotificationMeta: notification.NotificationMeta{
				NotificationType: notification.InstanceNeedsApproval,
				LeaseUUID:        leaseUUID,
				AWSResourceID:    tr.AWSResourceID(),
				//ResourceType:     lease.ResourceType,
			},
		})

		// remove message from queue
		err = tr.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	} else {
		// register new lease in DB
		// set its expiration to zone.default_expiration (if > 0), or cloudaccount.default_expiration, or adminAccount.default_expiration
		Logger.Info("Lease is OK -- register new lease in DB")

		//transmission.DefineLeaseDuration()
		var expiresAt = tr.LeaseExpiresAt()

		// these will be used to compose the urls and verify the requests
		leaseUUID := uuid.NewV4().String()
		tokenOnce := uuid.NewV4().String() // one-time token

		lease = &models.Lease{
			UUID:      leaseUUID,
			TokenOnce: tokenOnce,

			GroupUID:  groupUID,
			GroupType: tr.GroupType,

			OwnerID:        tr.Owner.ID,
			AccountID:      tr.Cloudaccount.AccountID,
			CloudaccountID: tr.Cloudaccount.ID,
			AWSAccountID:   tr.Cloudaccount.AWSID,
			Region:         tr.InstanceRegion,

			NumTimesAllertedAboutExpiry: NoAlertsSent, // the lease does not need an action response, no alert has been sent out

			LaunchedAt: tr.InstanceLaunchTimeUTC(),
			ExpiresAt:  expiresAt,
		}
		s.DB.Create(lease)
		Logger.Info("new lease created",
			"lease", lease,
		)

		newInstance := models.Instance{
			LeaseID: lease.ID,

			AccountID:      tr.Cloudaccount.AccountID,
			CloudaccountID: tr.Cloudaccount.ID,
			AWSAccountID:   tr.Cloudaccount.AWSID,

			GroupUID:  groupUID,
			GroupType: tr.GroupType,

			InstanceID:       tr.InstanceID(),
			AvailabilityZone: tr.AvailabilityZone(),
			InstanceType:     tr.InstanceType(),
			Region:           tr.InstanceRegion,

			LaunchedAt: tr.InstanceLaunchTimeUTC(),
		}
		s.DB.Create(&newInstance)
		Logger.Info("saving new instance", "newInstance", newInstance)

		// URL to terminate lease
		terminateURL, err := s.EmailActionGenerateSignedURL("terminate", leaseUUID, HashString(groupUID), tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		var emailValues = map[string]interface{}{
			"owner_email":     tr.Owner.Email,
			"instance_id":     tr.InstanceID(),
			"instance_type":   tr.InstanceType(),
			"resource_region": tr.InstanceRegion,

			"termination_time": expiresAt.Format("2006-01-02 15:04:05 GMT"),
			"lease_duration":   tr.LeaseExpiresAt().Sub(tr.InstanceLaunchTimeUTC()).String(),

			"lease_terminate_url": terminateURL,
		}

		newEmailBody, err := tools.CompileEmailTemplate(
			"new-lease-valid-owner-tag-no-approval-needed.txt",
			emailValues,
		)
		if err != nil {
			return err
		}

		var emailSubject string

		emailSubject = fmt.Sprintf("Instance (%v) created", tr.GroupType.String())

		s.Queues().NotifierQueue().PushTask(tasks.NotifierTask{
			AccountID: tr.AdminAccount.ID, // this will also trigger send to Slack
			To:        tr.Owner.Email,
			Subject:   emailSubject,
			BodyHTML:  newEmailBody,
			BodyText:  newEmailBody,
			NotificationMeta: notification.NotificationMeta{
				NotificationType: notification.InstanceCreated,
				LeaseUUID:        leaseUUID,
				AWSResourceID:    tr.AWSResourceID(),
				//ResourceType:     lease.ResourceType,
			},
		})

		// remove message from queue
		err = tr.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	}

}
