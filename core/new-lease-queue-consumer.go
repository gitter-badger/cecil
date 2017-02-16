package core

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/satori/go.uuid"
)

// NewLeaseQueueConsumer consumes NewLeaseTask from NewLeaseQueue
func (s *Service) NewLeaseQueueConsumer(t interface{}) error {

	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	transmission := t.(NewLeaseTask).Transmission

	Logger.Info(
		"NewLeaseQueueConsumer called",
		"transmission", transmission,
	)
	defer Logger.Info(
		"NewLeaseQueueConsumer call finished",
		"transmission", transmission,
	)

	Logger.Info(
		"Creating AssumedConfig",
		"topicRegion", transmission.Topic.Region,
		"topicAWSID", transmission.Topic.AWSID,
		"externalID", transmission.Cloudaccount.ExternalID,
	)

	if err := transmission.CreateAssumedService(); err != nil {
		// TODO: this might reveal too much to the admin about the service; be selective and cautious
		s.sendMisconfigurationNotice(err, transmission.AdminAccount.Email)
		Logger.Warn("error while creating assumed service", "err", err)
		return err
	}

	if err := transmission.CreateAssumedEC2Service(); err != nil {
		// TODO: this might reveal too much to the admin about the service; be selective and cautious
		s.sendMisconfigurationNotice(err, transmission.AdminAccount.Email)
		Logger.Warn("error while creating ec2 service with assumed service", "err", err)
		return err
	}

	if err := transmission.DescribeInstance(); err != nil {
		// TODO: this might reveal too much to the admin about the service; be selective and cautious
		s.sendMisconfigurationNotice(err, transmission.AdminAccount.Email)
		Logger.Warn("error while describing instances", "err", err)
		return err
	}

	// check whether the instance specified in the event exists on aws
	if !transmission.InstanceExists() {
		Logger.Warn("Instance does not exist", "instanceID", transmission.Message.Detail.InstanceID)
		// remove message from queue
		err := transmission.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	}

	Logger.Info(
		"describeInstances",
		"response", transmission.describeInstancesResponse,
	)

	if err := transmission.FetchInstance(); err != nil {
		Logger.Warn("error while fetching instance description", "err", err)
		return err
	}

	if err := transmission.ComputeInstanceRegion(); err != nil {
		Logger.Warn("error while computing instance region", "err", err)
		return err
	}

	/// transmission.Message.InstanceID == transmission.Instance.InstanceID
	//TODO: might this happen?

	/// transmission.Instance.IsTerminated()
	/// transmission.Message.Delete()

	// CreateAssumedCloudformationService which will be used to check whether the instance is part of a cloudformation stack
	if err := transmission.CreateAssumedCloudformationService(); err != nil {
		Logger.Warn("error while creating assumed cloudformation service", "err", err)
		return err
	}

	err := transmission.DefineResourceType()
	if err != nil {
		Logger.Warn("error while DefineResourceType", "err", err)
		// if the error code anything other than AccessDenied, return error
		if e, ok := err.(awserr.Error); !ok || e.Code() != "AccessDenied" {
			return err
		}
		Logger.Warn("Cannot determine whether the instance is part of a cloudformation stack; treating as a normal instance")
		// otherwise (i.e. the error is that the user is "access denied" to perform DescribeStackResources), register the instance as a normal lease (not as a stack)
	}

	// if the message signal that an instance has been terminated, create a task
	// to mark the lease as terminated
	if transmission.InstanceIsTerminated() {
		Logger.Info(
			"NewLeaseQueueConsumer",
			"InstanceIsTerminated()", transmission,
		)

		leaseTerminatedTask := LeaseTerminatedTask{
			AWSID:        transmission.Cloudaccount.AWSID,
			TerminatedAt: transmission.Message.Time.UTC(),
		}
		if transmission.IsStack() {
			// This will not terminate the lease for the stack immediately, but only trigger a
			// check whether the stack is currently running, and if it is not, ONLY THEN
			// the lease for the stack will be terminated.
			leaseTerminatedTask.ResourceType = StackResourceType
			leaseTerminatedTask.AWSResourceID = transmission.StackInfo.StackID
		}
		if transmission.IsInstance() {
			leaseTerminatedTask.ResourceType = InstanceResourceType
			leaseTerminatedTask.AWSResourceID = transmission.InstanceID()
		}

		s.LeaseTerminatedQueue.TaskQueue <- leaseTerminatedTask

		// remove message from queue
		err := transmission.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	}

	// do not consider states other than pending and terminated
	if !transmission.InstanceIsPendingOrRunning() {
		Logger.Warn("The retrieved state is neither pending nor running:", "state", transmission.Instance.State.Name)
		// remove message from queue
		// remove message from queue
		err := transmission.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	}

	if transmission.IsStack() {
		Logger.Debug("DefineResourceType", "bool", transmission.IsStack())
		Logger.Debug("transmission.StackInfo", "transmission.StackInfo", transmission.StackInfo)

		stackHasAlreadyALease, err := transmission.StackHasAlreadyALease()
		if err != nil {
			Logger.Warn("StackHasAlreadyALease", "err", err)
			return err
		}
		Logger.Debug("StackHasAlreadyALease", "bool", stackHasAlreadyALease)

		// If the stack already has a registered lease, just ignore this message
		// because this and the other instances of this stack
		// will be terminated all together with the stack.
		if stackHasAlreadyALease {
			// remove message from queue
			err := transmission.DeleteMessage()
			if err != nil {
				Logger.Warn("DeleteMessage", "err", err)
			}
			return err
		}

	} else {
		if !transmission.LeaseIsNew() {
			// TODO: notify admin
			Logger.Warn("!transmission.LeaseIsNew()")
			return nil // TODO: return an error ???
		}

	}

	if !transmission.InstanceHasTagOrKeyName() || !transmission.ExternalOwnerIsWhitelisted() {
		// assign instance to admin, and send notification to admin
		// owner is not whitelisted: notify admin
		Logger.Info("Transmission doesn't have owner tag/keyname or owner is not whitelisted.")

		err := transmission.SetAdminAsOwner()
		if err != nil {
			Logger.Warn("Error while setting admin as owner", "err", err)
			return err
		}

		//transmission.leaseDuration = s.Config.Lease.ApprovalTimeoutDuration
		var expiresAt = transmission.LeaseExpiresAt()

		// these will be used to compose the urls and verify the requests
		leaseUUID := uuid.NewV4().String()
		instanceID := transmission.InstanceID()
		tokenOnce := uuid.NewV4().String() // one-time token

		newLease := Lease{
			UUID:      leaseUUID,
			TokenOnce: tokenOnce,

			OwnerID:        transmission.owner.ID,
			AccountID:      transmission.Cloudaccount.AccountID,
			CloudaccountID: transmission.Cloudaccount.ID,
			AWSAccountID:   transmission.Cloudaccount.AWSID,
			Region:         transmission.instanceRegion,

			NumTimesAllertedAboutExpiry: AllAlertsSent, // this will prevent any other notification before the expiry

			LaunchedAt: transmission.InstanceLaunchTimeUTC(),
			ExpiresAt:  expiresAt,
		}
		s.DB.Create(&newLease)
		if transmission.IsStack() {
			// create an entry in the StackResources DB table
			newStackResource := StackResource{}
			newStackResource.StackID = transmission.StackInfo.StackID
			newStackResource.StackName = transmission.StackInfo.StackName
			newStackResource.LeaseID = newLease.ID
			s.DB.Create(&newStackResource)
			// annotate the ID of the StackResource in the lease
			newLease.ResourceType = StackResourceType
			newLease.ResourceID = newStackResource.ID
		} else {
			// create an entry in the InstanceResources DB table
			newInstanceResource := InstanceResource{}
			newInstanceResource.InstanceID = transmission.InstanceID()
			newInstanceResource.AvailabilityZone = transmission.AvailabilityZone()
			newInstanceResource.InstanceType = transmission.InstanceType()
			newInstanceResource.LeaseID = newLease.ID
			s.DB.Create(&newInstanceResource)
			// annotate the ID of the InstanceResource in the lease
			newLease.ResourceType = InstanceResourceType
			newLease.ResourceID = newInstanceResource.ID
		}
		s.DB.Save(&newLease)
		Logger.Info("new lease created",
			"lease", newLease,
		)

		var newEmailBody string

		// URL to approve lease
		Logger.Info(
			"Creating lease signature",
			"lease_uuid", leaseUUID,
			"instance_id", instanceID,
			"action", "approve",
			"token_once", tokenOnce,
		)
		approveURL, err := s.EmailActionGenerateSignedURL("approve", leaseUUID, newLease.ResourceID, tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		// URL to terminate lease
		terminateURL, err := s.EmailActionGenerateSignedURL("terminate", leaseUUID, newLease.ResourceID, tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		var emailValues = map[string]interface{}{
			"owner_email":     transmission.owner.Email,
			"instance_id":     transmission.InstanceID(),
			"instance_type":   transmission.InstanceType(),
			"resource_region": transmission.instanceRegion,

			"termination_time": expiresAt.Format("2006-01-02 15:04:05 GMT"),
			"lease_duration":   transmission.LeaseExpiresAt().Sub(transmission.InstanceLaunchTimeUTC()).String(),

			"lease_terminate_url": terminateURL,
			"lease_approve_url":   approveURL,
		}

		if transmission.IsStack() {
			emailValues["stack_id"] = transmission.StackInfo.StackID
			emailValues["stack_name"] = transmission.StackInfo.StackName
		}

		switch {
		case !transmission.InstanceHasTagOrKeyName():
			newEmailBody, err = CompileEmailTemplate(
				"new-lease-no-owner-tag.txt",
				emailValues,
			)
			if err != nil {
				return err
			}
			break

		case !transmission.ExternalOwnerIsWhitelisted():
			newEmailBody, err = CompileEmailTemplate(
				"new-lease-owner-tag-not-whitelisted.txt",
				emailValues,
			)
			if err != nil {
				return err
			}
		}

		var emailSubject string

		if transmission.IsStack() {
			emailSubject = fmt.Sprintf("Stack (%v) needs attention", transmission.StackInfo.StackName)
		} else {
			emailSubject = fmt.Sprintf("Instance (%v) needs attention", transmission.InstanceID())
		}

		Logger.Info("Adding new NotifierTask")
		s.NotifierQueue.TaskQueue <- NotifierTask{
			AccountID: transmission.AdminAccount.ID, // this will also trigger send to Slack
			//To:       owner.Email,
			To:       transmission.AdminAccount.Email,
			Subject:  emailSubject,
			BodyHTML: newEmailBody,
			BodyText: newEmailBody,
			NotificationMeta: NotificationMeta{
				NotificationType: InstanceNeedsAttention,
				LeaseUUID:        leaseUUID,
				AWSResourceID:    transmission.AWSResourceID(),
				ResourceType:     newLease.ResourceType,
			},
		}

		Logger.Info("Delete SQS Message")
		err = transmission.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	}

	if err := transmission.SetExternalOwnerAsOwner(); err != nil {
		Logger.Warn("Error while setting external owner as owner", "err", err)
	}

	if transmission.LeaseNeedsApproval() {
		// register new lease in DB
		// expiry: 1h
		// send confirmation to owner: confirmation link, and termination link
		Logger.Info("Lease needs approval")

		//transmission.leaseDuration = s.Config.Lease.ApprovalTimeoutDuration
		var expiresAt = transmission.LeaseExpiresAt()

		// these will be used to compose the urls and verify the requests
		leaseUUID := uuid.NewV4().String()
		tokenOnce := uuid.NewV4().String() // one-time token

		newLease := Lease{
			UUID:      leaseUUID,
			TokenOnce: tokenOnce,

			OwnerID:        transmission.owner.ID,
			AccountID:      transmission.Cloudaccount.AccountID,
			CloudaccountID: transmission.Cloudaccount.ID,
			AWSAccountID:   transmission.Cloudaccount.AWSID,
			Region:         transmission.instanceRegion,

			NumTimesAllertedAboutExpiry: AllAlertsSent,

			LaunchedAt: transmission.InstanceLaunchTimeUTC(),
			ExpiresAt:  expiresAt,
		}
		s.DB.Create(&newLease)
		if transmission.IsStack() {
			// create an entry in the StackResources DB table
			newStackResource := StackResource{}
			newStackResource.StackID = transmission.StackInfo.StackID
			newStackResource.StackName = transmission.StackInfo.StackName
			newStackResource.LeaseID = newLease.ID
			s.DB.Create(&newStackResource)
			// annotate the ID of the StackResource in the lease
			newLease.ResourceType = StackResourceType
			newLease.ResourceID = newStackResource.ID
		} else {
			// create an entry in the InstanceResources DB table
			newInstanceResource := InstanceResource{}
			newInstanceResource.InstanceID = transmission.InstanceID()
			newInstanceResource.AvailabilityZone = transmission.AvailabilityZone()
			newInstanceResource.InstanceType = transmission.InstanceType()
			newInstanceResource.LeaseID = newLease.ID
			s.DB.Create(&newInstanceResource)
			// annotate the ID of the InstanceResource in the lease
			newLease.ResourceType = InstanceResourceType
			newLease.ResourceID = newInstanceResource.ID
		}
		s.DB.Save(&newLease)
		Logger.Info("new lease created",
			"lease", newLease,
		)

		// URL to approve lease
		approveURL, err := s.EmailActionGenerateSignedURL("approve", leaseUUID, newLease.ResourceID, tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		// URL to terminate lease
		terminateURL, err := s.EmailActionGenerateSignedURL("terminate", leaseUUID, newLease.ResourceID, tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		var emailValues = map[string]interface{}{
			"owner_email":        transmission.owner.Email,
			"n_of_active_leases": transmission.activeLeaseCount,
			"instance_id":        transmission.InstanceID(),
			"instance_type":      transmission.InstanceType(),
			"resource_region":    transmission.instanceRegion,

			"termination_time": expiresAt.Format("2006-01-02 15:04:05 GMT"),
			"lease_duration":   transmission.LeaseExpiresAt().Sub(transmission.InstanceLaunchTimeUTC()).String(),

			"lease_approve_url":   approveURL,
			"lease_terminate_url": terminateURL,
		}

		if transmission.IsStack() {
			emailValues["stack_id"] = transmission.StackInfo.StackID
			emailValues["stack_name"] = transmission.StackInfo.StackName
		}

		newEmailBody, err := CompileEmailTemplate(
			"new-lease-valid-owner-tag-needs-approval.txt",
			emailValues,
		)
		if err != nil {
			return err
		}

		var emailSubject string

		if transmission.IsStack() {
			emailSubject = fmt.Sprintf("Stack (%v) needs approval", transmission.StackInfo.StackName)
		} else {
			emailSubject = fmt.Sprintf("Instance (%v) needs approval", transmission.InstanceID())
		}
		s.NotifierQueue.TaskQueue <- NotifierTask{
			AccountID: transmission.AdminAccount.ID, // this will also trigger send to Slack
			To:        transmission.owner.Email,
			Subject:   emailSubject,
			BodyHTML:  newEmailBody,
			BodyText:  newEmailBody,
			NotificationMeta: NotificationMeta{
				NotificationType: InstanceNeedsApproval,
				LeaseUUID:        leaseUUID,
				AWSResourceID:    transmission.AWSResourceID(),
				ResourceType:     newLease.ResourceType,
			},
		}

		// remove message from queue
		err = transmission.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	} else {
		// register new lease in DB
		// set its expiration to zone.default_expiration (if > 0), or cloudaccount.default_expiration, or adminAccount.default_expiration
		Logger.Info("Lease is OK -- register new lease in DB")

		//transmission.DefineLeaseDuration()
		var expiresAt = transmission.LeaseExpiresAt()

		// these will be used to compose the urls and verify the requests
		leaseUUID := uuid.NewV4().String()
		tokenOnce := uuid.NewV4().String() // one-time token

		newLease := Lease{
			UUID:      leaseUUID,
			TokenOnce: tokenOnce,

			OwnerID:        transmission.owner.ID,
			AccountID:      transmission.Cloudaccount.AccountID,
			CloudaccountID: transmission.Cloudaccount.ID,
			AWSAccountID:   transmission.Cloudaccount.AWSID,
			Region:         transmission.instanceRegion,

			NumTimesAllertedAboutExpiry: NoAlertsSent, // the lease does not need an action response, no alert has been sent out

			LaunchedAt: transmission.InstanceLaunchTimeUTC(),
			ExpiresAt:  expiresAt,
		}
		s.DB.Create(&newLease)
		if transmission.IsStack() {
			// create an entry in the StackResources DB table
			// annotate the ID of the StackResource in the lease
			newStackResource := StackResource{}
			newStackResource.StackID = transmission.StackInfo.StackID
			newStackResource.StackName = transmission.StackInfo.StackName
			newStackResource.LeaseID = newLease.ID
			s.DB.Create(&newStackResource)
			// annotate the ID of the InstanceResource in the lease
			newLease.ResourceType = StackResourceType
			newLease.ResourceID = newStackResource.ID
		} else {
			// create an entry in the InstanceResources DB table
			newInstanceResource := InstanceResource{}
			newInstanceResource.InstanceID = transmission.InstanceID()
			newInstanceResource.AvailabilityZone = transmission.AvailabilityZone()
			newInstanceResource.InstanceType = transmission.InstanceType()
			newInstanceResource.LeaseID = newLease.ID
			s.DB.Create(&newInstanceResource)
			// annotate the ID of the InstanceResource in the lease
			newLease.ResourceType = InstanceResourceType
			newLease.ResourceID = newInstanceResource.ID
		}
		s.DB.Save(&newLease)
		Logger.Info("new lease created",
			"lease", newLease,
		)

		// URL to terminate lease
		terminateURL, err := s.EmailActionGenerateSignedURL("terminate", leaseUUID, newLease.ResourceID, tokenOnce)
		if err != nil {
			// TODO: notify admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		var emailValues = map[string]interface{}{
			"owner_email":     transmission.owner.Email,
			"instance_id":     transmission.InstanceID(),
			"instance_type":   transmission.InstanceType(),
			"resource_region": transmission.instanceRegion,

			"termination_time": expiresAt.Format("2006-01-02 15:04:05 GMT"),
			"lease_duration":   transmission.LeaseExpiresAt().Sub(transmission.InstanceLaunchTimeUTC()).String(),

			"lease_terminate_url": terminateURL,
		}

		if transmission.IsStack() {
			emailValues["stack_id"] = transmission.StackInfo.StackID
			emailValues["stack_name"] = transmission.StackInfo.StackName
		}

		newEmailBody, err := CompileEmailTemplate(
			"new-lease-valid-owner-tag-no-approval-needed.txt",
			emailValues,
		)
		if err != nil {
			return err
		}

		var emailSubject string

		if transmission.IsStack() {
			emailSubject = fmt.Sprintf("Stack (%v) created", transmission.StackInfo.StackName)
		} else {
			emailSubject = fmt.Sprintf("Instance (%v) created", transmission.InstanceID())
		}

		s.NotifierQueue.TaskQueue <- NotifierTask{
			AccountID: transmission.AdminAccount.ID, // this will also trigger send to Slack
			To:        transmission.owner.Email,
			Subject:   emailSubject,
			BodyHTML:  newEmailBody,
			BodyText:  newEmailBody,
			NotificationMeta: NotificationMeta{
				NotificationType: InstanceCreated,
				LeaseUUID:        leaseUUID,
				AWSResourceID:    transmission.AWSResourceID(),
				ResourceType:     newLease.ResourceType,
			},
		}

		// remove message from queue
		err = transmission.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	}

}
