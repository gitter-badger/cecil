package core

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/satori/go.uuid"
	"gopkg.in/mailgun/mailgun-go.v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
)

// @@@@@@@@@@@@@@@ Task consumers @@@@@@@@@@@@@@@

// TerminatorQueueConsumer consumes TerminatorTask from TerminatorQueue;
// sends instance termination request to AWS ec2.
func (s *Service) TerminatorQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(TerminatorTask)
	// TODO: check whether fields are non-null and valid
	Logger.Info("TerminatorQueueConsumer",
		"task", task,
	)

	var cloudaccount Cloudaccount
	err := s.DB.Model(&task.Lease).Related(&cloudaccount).Error
	//s.DB.Table("accounts").Where([]uint{cloudaccount.AccountID}).First(&cloudaccount).Count(&leaseCloudOwnerCount)
	if err != nil {
		// TODO: notify admin; something fishy is going on.
		Logger.Warn("here", err.Error())
		return err
	}

	// assume role
	assumedConfig := &aws.Config{
		Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
			Client: sts.New(s.AWS.Session, &aws.Config{Region: aws.String(task.Region)}),
			RoleARN: fmt.Sprintf(
				"arn:aws:iam::%v:role/%v",
				cloudaccount.AWSID,
				s.AWS.Config.ForeignIAMRoleName,
			),
			RoleSessionName: uuid.NewV4().String(),
			ExternalID:      aws.String(cloudaccount.ExternalID),
			ExpiryWindow:    3 * time.Minute,
		}),
	}

	assumedService := session.New(assumedConfig)

	if task.Lease.IsStack() {

		var stack StackResource
		raw, err := s.ResourceOf(&task.Lease)
		if err != nil {
			return err
		}
		stack = raw.(StackResource)

		assumedCloudformationService := cloudformation.New(assumedService)

		DescribeStackResourcesParams := &cloudformation.DescribeStackResourcesInput{
			StackName: aws.String(stack.StackName),
		}
		resp, err := assumedCloudformationService.DescribeStackResources(DescribeStackResourcesParams)
		Logger.Info("DescribeStackResources", "response", resp)
		Logger.Info("DescribeStackResources", "err", err)

		if err != nil {
			// TODO: check whether this effectively is a way to catch a "not found"
			//if e, ok := err.(awserr.RequestFailure); ok && e.StatusCode() == 404 {
			if strings.Contains(err.Error(), "does not exist") {
				e, ok := err.(awserr.RequestFailure)
				if ok {
					Logger.Info("DescribeStackResources", "e", e)
				} else {
					Logger.Info("DescribeStackResources", "err", err)
				}

				task.Lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve

				// we don't know when it has been terminated, so just use the current time
				now := time.Now().UTC()
				task.Lease.TerminatedAt = &now

				// TODO: use the ufficial time of termination, from th sqs message, because if erminated via link, the termination time is not expiresAt
				// task.Lease.TerminatedAt = time.Now().UTC()
				s.DB.Save(&task.Lease)

				Logger.Debug(
					"TerminatorQueueConsumer",
					"err", err,
					"action_taken", "removing lease of already deleted/non-existent stack from DB",
				)

			} else {
				// TODO: cleaner way to do this?  cloudaccount.Account would be nice .. gorma provides this
				var account Account
				s.DB.First(&account, cloudaccount.AccountID)

				recipientEmail := account.Email

				s.sendMisconfigurationNotice(err, recipientEmail)
			}
			return err
		}

		DeleteStackParams := &cloudformation.DeleteStackInput{
			StackName: aws.String(stack.StackName), // Required
		}
		deleteStackResponse, err := assumedCloudformationService.DeleteStack(DeleteStackParams)
		Logger.Info("DeleteStack", "response", deleteStackResponse)
		Logger.Info("DeleteStack", "err", err)

		// TODO: handle error
		return nil
	}

	var instance InstanceResource
	raw, err := s.ResourceOf(&task.Lease)
	if err != nil {
		return err
	}
	instance = raw.(InstanceResource)

	assumedEC2Service := s.EC2(assumedService, task.Region)

	terminateInstanceParams := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{ // Required
			aws.String(instance.InstanceID),
		},
	}
	terminateInstanceResponse, err := assumedEC2Service.TerminateInstances(terminateInstanceParams)

	Logger.Info("TerminateInstances", "response", terminateInstanceResponse)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.

		if strings.Contains(err.Error(), "InvalidInstanceID.NotFound") {
			// TODO: replace this with something shorter

			task.Lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve

			// we don't know when it has been terminated, so just use the current time
			now := time.Now().UTC()
			task.Lease.TerminatedAt = &now

			// TODO: use the ufficial time of termination, from th sqs message, because if erminated via link, the termination time is not expiresAt
			// lease.TerminatedAt = time.Now().UTC()
			s.DB.Save(&task.Lease)

			Logger.Debug(
				"TerminatorQueueConsumer TerminateInstances ",
				"err", err,
				"action_taken", "removing lease of already deleted/non-existent instance from DB",
			)

		} else {
			// TODO: cleaner way to do this?  cloudaccount.Account would be nice .. gorma provides this
			var account Account
			s.DB.First(&account, cloudaccount.AccountID)

			recipientEmail := account.Email

			s.sendMisconfigurationNotice(err, recipientEmail)
		}

		return err
	}

	return nil
}

// LeaseTerminatedQueueConsumer consumes LeaseTerminatedTask from LeaseTerminatedQueue;
// marks leases as terminated and notifes the owner.
func (s *Service) LeaseTerminatedQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(LeaseTerminatedTask)
	Logger.Info(
		"Marking lease as terminated",
		"AWSResourceID", task.AWSResourceID,
		"resourceType", task.ResourceType,
		"task", task,
	)

	var err error
	var resourceID uint

	var stack StackResource
	if task.ResourceType == StackResourceType {
		stack, err = s.StackByStackID(task.AWSResourceID)
		if err != nil {
			return fmt.Errorf("Error while fetching stack: %v", err)
		}
		resourceID = stack.ID
	}

	var instance InstanceResource
	if task.ResourceType == InstanceResourceType {
		instance, err = s.InstanceByInstanceID(task.AWSResourceID)
		if err != nil {
			return fmt.Errorf("Error while fetching instance: %v", err)
		}
		resourceID = instance.ID
	}

	var lease Lease
	err = s.DB.Table("leases").
		Where(&Lease{
			ResourceID:   resourceID,
			ResourceType: task.ResourceType,
			AWSAccountID: task.AWSID,
		}).
		Where("terminated_at IS NULL").
		Find(&lease).
		Error

	if err != nil {
		Logger.Warn("Lease for deletion not found", "err", err)
		return fmt.Errorf("Lease for deletion not found: %v", err)
	}

	lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve

	// TODO: check whether this time is correct
	lease.TerminatedAt = &task.TerminatedAt

	// TODO: use the ufficial time of termination, from th sqs message, because if erminated via link, the termination time is not expiresAt
	// lease.TerminatedAt = time.Now().UTC()
	s.DB.Save(&lease)

	Logger.Info(
		"Marking lease as terminated",
		"lease", lease,
	)

	var owner Owner

	err = s.DB.Table("owners").Where(lease.OwnerID).First(&owner).Error

	if err != nil {
		Logger.Warn("LeaseTerminatedQueueConsumer: error fetching owner", "err", err)
		return fmt.Errorf("LeaseTerminatedQueueConsumer: error fetching owner: %v", err)
	}

	var emailValues = map[string]interface{}{
		"owner_email":     owner.Email,
		"resource_region": lease.Region,

		"lease_duration": task.TerminatedAt.Sub(lease.CreatedAt).String(),
		"expires_at":     lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
		"terminated_at":  task.TerminatedAt.Format("2006-01-02 15:04:05 GMT"),
	}

	if lease.IsInstance() {
		emailValues["instance_id"] = instance.InstanceID
		emailValues["instance_type"] = instance.InstanceType
	}
	if lease.IsStack() {
		emailValues["stack_id"] = stack.StackID
		emailValues["stack_name"] = stack.StackName
	}

	newEmailBody, err := CompileEmailTemplate(
		"lease-resource-terminated.txt",
		emailValues,
	)
	if err != nil {
		return err
	}

	var newEmailSubject string
	if lease.IsStack() {
		newEmailSubject = fmt.Sprintf("Stack (%v) terminated", stack.StackName)
	} else {
		newEmailSubject = fmt.Sprintf("Instance (%v) terminated", instance.InstanceID)
	}

	s.NotifierQueue.TaskQueue <- NotifierTask{
		AccountID: lease.AccountID, // this will also trigger send to Slack
		To:        owner.Email,
		Subject:   newEmailSubject,
		BodyHTML:  newEmailBody,
		BodyText:  newEmailBody,
		NotificationMeta: NotificationMeta{
			NotificationType: InstanceTerminated,
			LeaseUUID:        lease.UUID,
			AWSResourceID:    task.AWSResourceID,
			ResourceType:     lease.ResourceType,
		},
	}

	return nil
}

// ExtenderQueueConsumer consumes ExtenderTask from ExtenderQueue; approves or extends leases.
func (s *Service) ExtenderQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(ExtenderTask)
	// TODO: check whether fields are non-null and valid

	if task.Approving {
		Logger.Info("Approving lease",
			"resourceType", task.Lease.ResourceType,
			"resourceID", task.Lease.ResourceID,
		)
	} else {
		Logger.Info("Extending lease",
			"resourceType", task.Lease.ResourceType,
			"resourceID", task.Lease.ResourceID,
		)
	}

	if task.Lease.IsExpired() {
		// TODO: should the user be notified that the lease cannot be extended because it is already expired???
		err := errors.New("lease is already expired; cannot extend/approve")
		Logger.Error(
			"error while extendin lease",
			"resourceType", task.Lease.ResourceType,
			"resourceID", task.Lease.ResourceID,
			"err", err,
		)
		return err
	}

	task.Lease.TokenOnce = uuid.NewV4().String() // invalidates all other URLs to renew/terminate/approve
	task.Lease.NumTimesAllertedAboutExpiry = NoAlertsSent

	// define the lease duration
	leaseDuration, err := s.DefineLeaseDuration(task.Lease.AccountID, task.Lease.CloudaccountID)
	if err != nil {
		Logger.Error(
			"error while DefineLeaseDuration",
			"resourceType", task.Lease.ResourceType,
			"resourceID", task.Lease.ResourceID,
			"err", err,
		)
		return err
	}
	// TODO: remove leaseDuration
	_ = leaseDuration

	if task.Approving {
		now := time.Now().UTC()
		task.Lease.ApprovedAt = &now
		//task.Lease.ExpiresAt = task.Lease.CreatedAt.Add(leaseDuration)
	} else {
		leaseDuration := task.Lease.ExpiresAt.Sub(task.Lease.CreatedAt)
		task.Lease.ExpiresAt = task.Lease.ExpiresAt.Add(leaseDuration)
	}

	s.DB.Save(&task.Lease)

	var owner Owner
	var ownerCount int64

	s.DB.Table("owners").Where(task.Lease.OwnerID).First(&owner).Count(&ownerCount)

	var newEmailBody string
	var newEmailSubject string
	var notificationType NotificationType

	var AWSResourceID string

	var instance InstanceResource
	if task.Lease.IsInstance() {
		raw, err := s.ResourceOf(&task.Lease)
		if err != nil {
			return err
		}
		instance = raw.(InstanceResource)
		AWSResourceID = instance.InstanceID
	}

	var stack StackResource
	if task.Lease.IsStack() {
		raw, err := s.ResourceOf(&task.Lease)
		if err != nil {
			return err
		}
		stack = raw.(StackResource)
		AWSResourceID = stack.StackID
	}

	var emailValues = map[string]interface{}{
		"owner_email":     owner.Email,
		"resource_region": task.Lease.Region,

		"lease_duration": task.Lease.ExpiresAt.Sub(task.Lease.CreatedAt).String(),

		"expires_at": task.Lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
	}

	if task.Lease.IsInstance() {
		emailValues["instance_id"] = instance.InstanceID
		emailValues["instance_type"] = instance.InstanceType
	}
	if task.Lease.IsStack() {
		emailValues["stack_id"] = stack.StackID
		emailValues["stack_name"] = stack.StackName
	}

	if task.Approving {
		notificationType = LeaseApproved

		if task.Lease.IsStack() {
			newEmailSubject = fmt.Sprintf("Stack (%v) lease approved", stack.StackName)
		} else {
			newEmailSubject = fmt.Sprintf("Instance (%v) lease approved", instance.InstanceID)
		}

		newEmailBody, err = CompileEmailTemplate(
			"lease-approved.txt",
			emailValues,
		)
		if err != nil {
			return err
		}
	} else {
		notificationType = LeaseExtended

		if task.Lease.IsStack() {
			newEmailSubject = fmt.Sprintf("Stack (%v) lease extended", stack.StackName)
		} else {
			newEmailSubject = fmt.Sprintf("Instance (%v) lease extended", instance.InstanceID)
		}

		newEmailBody, err = CompileEmailTemplate(
			"lease-extended.txt",
			emailValues,
		)
		if err != nil {
			return err
		}
	}

	s.NotifierQueue.TaskQueue <- NotifierTask{
		AccountID: task.Lease.AccountID, // this will also trigger send to Slack
		To:        owner.Email,
		Subject:   newEmailSubject,
		BodyHTML:  newEmailBody,
		BodyText:  newEmailBody,
		NotificationMeta: NotificationMeta{
			NotificationType: notificationType,
			LeaseUUID:        task.Lease.UUID,
			AWSResourceID:    AWSResourceID,
			ResourceType:     task.Lease.ResourceType,
		},
	}

	return nil
}

// NotifierQueueConsumer consumes NotifierTask from NotifierQueue; sends messages
func (s *Service) NotifierQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(NotifierTask)
	// TODO: check whether fields are non-null and valid
	Logger.Info("Sending EMAIL",
		"to", task.To,
	)

	// if there is a SlackBotInstance for the Account,
	// send in a goroutine a message to that SlackBotInstance.
	go func() {
		if task.AccountID == 0 {
			return
		}
		slackIns, err := s.SlackBotInstanceByID(task.AccountID)
		if err != nil {
			//Logger.Warn("SlackBotInstanceByID", "warn", err)
			return
		}
		// HACK: the message sent to Slack should have custom formatting;
		// right now it is just the HTML of the email without html tags.
		messageWithoutHTML, err := sanitize.HTMLAllowing(task.BodyHTML, []string{"a"}, []string{"href"})
		messageWithoutHTML = strings.Replace(messageWithoutHTML, `<a href="`, "", -1)
		messageWithoutHTML = strings.Replace(messageWithoutHTML, `">Click here to terminate</a>`, "", -1)
		messageWithoutHTML = strings.Replace(messageWithoutHTML, `">Click here to approve</a>`, "", -1)
		slackIns.OutgoingMessages <- messageWithoutHTML
	}()

	// define the meailer to use (DefaultMailer or a mailer defined by account)
	var mailer *MailerInstance

	if task.AccountID > 0 {
		mailerIns, err := s.MailerInstanceByID(task.AccountID)
		if err != nil {
			//Logger.Warn("MailerInstanceByID", "warn", err)
		} else {
			mailer = mailerIns
			Logger.Info("using custom mailer", "mailer", *mailer)
		}
	}

	if mailer == nil {
		mailer = &s.DefaultMailer
	}

	message := mailgun.NewMessage(
		mailer.FromAddress,
		task.Subject,
		task.BodyText,
		task.To,
	)

	message.AddHeader(X_CECIL_MESSAGETYPE, fmt.Sprintf("%s", task.NotificationMeta.NotificationType))
	message.AddHeader(X_CECIL_LEASE_UUID, task.NotificationMeta.LeaseUUID)
	message.AddHeader(X_CECIL_AWS_RESOURCE_ID, task.NotificationMeta.AWSResourceID)
	message.AddHeader(X_CECIL_VERIFICATION_TOKEN, task.NotificationMeta.VerificationToken)

	//message.SetTracking(true)
	if task.DeliverAfter > 0 {
		message.SetDeliveryTime(time.Now().Add(task.DeliverAfter))
	}

	message.SetHtml(task.BodyHTML)

	err := Retry(10, time.Second*5, func() error {
		var err error
		_, _, err = mailer.Client.Send(message)
		return err
	}, nil)
	if err != nil {
		Logger.Error("Error while sending email", "err", err)
		return err
	}

	return nil
}
