package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/satori/go.uuid"
	"gopkg.in/mailgun/mailgun-go.v1"

	"github.com/aws/aws-sdk-go/aws"
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

	var cloudAccount CloudAccount
	var leaseCloudOwnerCount int64
	s.DB.Model(&task.Lease).Related(&cloudAccount).Count(&leaseCloudOwnerCount)
	//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&leaseCloudOwnerCount)
	if leaseCloudOwnerCount == 0 {
		// TODO: notify admin; something fishy is going on.
		Logger.Warn("leaseCloudOwnerCount == 0")
		return fmt.Errorf("leaseCloudOwnerCount == 0")
	}

	// assume role
	assumedConfig := &aws.Config{
		Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
			Client: sts.New(s.AWS.Session, &aws.Config{Region: aws.String(task.Region)}),
			RoleARN: fmt.Sprintf(
				"arn:aws:iam::%v:role/%v",
				cloudAccount.AWSID,
				s.AWS.Config.ForeignIAMRoleName,
			),
			RoleSessionName: uuid.NewV4().String(),
			ExternalID:      aws.String(cloudAccount.ExternalID),
			ExpiryWindow:    3 * time.Minute,
		}),
	}

	assumedService := session.New(assumedConfig)

	leaseIsStack := task.Lease.StackName != ""

	if leaseIsStack {
		assumedCloudformationService := cloudformation.New(assumedService)

		DescribeStackResourcesParams := &cloudformation.DescribeStackResourcesInput{
			StackName: aws.String(task.Lease.StackName),
		}
		resp, err := assumedCloudformationService.DescribeStackResources(DescribeStackResourcesParams)
		Logger.Info("DescribeStackResources", "response", resp)
		Logger.Info("DescribeStackResources", "err", err)

		if err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				// TODO: replace this with something shorter

				var lease Lease
				var leasesFound int64
				s.DB.Table("leases").Where(&Lease{
					InstanceID:   task.InstanceID,
					AWSAccountID: task.AWSAccountID,
				}).Where("terminated_at IS NULL").First(&lease).Count(&leasesFound)

				if leasesFound == 0 {
					Logger.Warn("Lease for deletion not found", "count", leasesFound, "instanceID", task.InstanceID)
					return fmt.Errorf("Lease for deletion not found: %v=%v", "count", leasesFound)
				}
				if leasesFound > 1 {
					Logger.Warn("Found multiple leases for deletion", "count", leasesFound)
					return fmt.Errorf("Found multiple leases for deletion: %v=%v", "count", leasesFound)
				}

				lease.Terminated = true
				lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve

				// we don't know when it has been terminated, so just use the current time
				now := time.Now().UTC()
				lease.TerminatedAt = &now

				// TODO: use the ufficial time of termination, from th sqs message, because if erminated via link, the termination time is not expiresAt
				// lease.TerminatedAt = time.Now().UTC()
				s.DB.Save(&lease)

				Logger.Debug(
					"TerminatorQueueConsumer TerminateInstances ",
					"err", err,
					"action_taken", "removing lease of already deleted/non-existent stack from DB",
				)

			} else {
				// TODO: cleaner way to do this?  cloudAccount.Account would be nice .. gorma provides this
				var account Account
				s.DB.First(&account, cloudAccount.AccountID)

				recipientEmail := account.Email

				s.sendMisconfigurationNotice(err, recipientEmail)
			}
			return err
		}

		DeleteStackParams := &cloudformation.DeleteStackInput{
			StackName: aws.String(task.Lease.StackName), // Required
		}
		deleteStackResponse, err := assumedCloudformationService.DeleteStack(DeleteStackParams)
		Logger.Info("DeleteStack", "response", deleteStackResponse)
		Logger.Info("DeleteStack", "err", err)

		// TODO: handle error
		return nil
	}

	assumedEC2Service := s.EC2(assumedService, task.Region)

	terminateInstanceParams := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{ // Required
			aws.String(task.InstanceID),
		},
	}
	terminateInstanceResponse, err := assumedEC2Service.TerminateInstances(terminateInstanceParams)

	Logger.Info("TerminateInstances", "response", terminateInstanceResponse)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.

		if strings.Contains(err.Error(), "InvalidInstanceID.NotFound") {
			// TODO: replace this with something shorter

			var lease Lease
			var leasesFound int64
			s.DB.Table("leases").Where(&Lease{
				InstanceID:   task.InstanceID,
				AWSAccountID: task.AWSAccountID,
			}).Where("terminated_at IS NULL").First(&lease).Count(&leasesFound)

			if leasesFound == 0 {
				Logger.Warn("Lease for deletion not found", "count", leasesFound, "instanceID", task.InstanceID)
				return fmt.Errorf("Lease for deletion not found: %v=%v", "count", leasesFound)
			}
			if leasesFound > 1 {
				Logger.Warn("Found multiple leases for deletion", "count", leasesFound)
				return fmt.Errorf("Found multiple leases for deletion: %v=%v", "count", leasesFound)
			}

			lease.Terminated = true
			lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve

			// we don't know when it has been terminated, so just use the current time
			now := time.Now().UTC()
			lease.TerminatedAt = &now

			// TODO: use the ufficial time of termination, from th sqs message, because if erminated via link, the termination time is not expiresAt
			// lease.TerminatedAt = time.Now().UTC()
			s.DB.Save(&lease)

			Logger.Debug(
				"TerminatorQueueConsumer TerminateInstances ",
				"err", err,
				"action_taken", "removing lease of already deleted/non-existent instance from DB",
			)

		} else {
			// TODO: cleaner way to do this?  cloudAccount.Account would be nice .. gorma provides this
			var account Account
			s.DB.First(&account, cloudAccount.AccountID)

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
	Logger.Info("Marking lease as terminated",
		"InstanceID", task.InstanceID,
	)

	var lease Lease
	var leasesFound int64
	s.DB.Table("leases").Where(&Lease{
		InstanceID:   task.InstanceID,
		AWSAccountID: task.AWSID,
	}).Where("terminated_at IS NULL").First(&lease).Count(&leasesFound)

	if leasesFound == 0 {
		Logger.Warn("Lease for deletion not found", "count", leasesFound, "instanceID", task.InstanceID)
		return fmt.Errorf("Lease for deletion not found: %v=%v", "count", leasesFound)
	}
	if leasesFound > 1 {
		Logger.Warn("Found multiple leases for deletion", "count", leasesFound)
		return fmt.Errorf("Found multiple leases for deletion: %v=%v", "count", leasesFound)
	}

	lease.Terminated = true
	lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve

	// TODO: check whether this time is correct
	lease.TerminatedAt = &task.TerminatedAt

	// TODO: use the ufficial time of termination, from th sqs message, because if erminated via link, the termination time is not expiresAt
	// lease.TerminatedAt = time.Now().UTC()
	s.DB.Save(&lease)

	var owner Owner
	var ownerCount int64

	s.DB.Table("owners").Where(lease.OwnerID).First(&owner).Count(&ownerCount)

	if ownerCount != 1 {
		Logger.Warn("LeaseTerminatedQueueConsumer: ownerCount is not 1", "count", ownerCount)
		return fmt.Errorf("LeaseTerminatedQueueConsumer: ownerCount is not 1: %v=%v", "count", ownerCount)
	}

	var emailValues = map[string]interface{}{
		"owner_email":     owner.Email,
		"instance_id":     lease.InstanceID,
		"instance_type":   lease.InstanceType,
		"instance_region": lease.Region,

		"instance_duration": task.TerminatedAt.Sub(lease.CreatedAt).String(),
		"expires_at":        lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
		"terminated_at":     task.TerminatedAt.Format("2006-01-02 15:04:05 GMT"),
	}

	if lease.StackName != "" {
		emailValues["logical_id"] = lease.LogicalID
		emailValues["stack_id"] = lease.StackID
		emailValues["stack_name"] = lease.StackName
	}

	newEmailBody := CompileEmail(
		`
		{{if not .stack_name }}
		Hey {{.owner_email}}, instance with id <b>{{.instance_id}}</b>
				(of type <b>{{.instance_type}}</b>,
				on <b>{{.instance_region}}</b>, expiry on <b>{{.expires_at}}</b>) has been terminated at
				<b>{{.terminated_at}}</b> ({{.instance_duration}} after it's creation)
		{{end}}

				{{if .stack_name }}
				Hey {{.owner_email}}, the following stack (expiry on <b>{{.expires_at}}</b>) has been terminated at
						<b>{{.terminated_at}}</b> ({{.instance_duration}} after it's creation)

						<br>
						<br>

				Stack name: <b>{{.stack_name}}</b><br>
				Stack id: <b>{{.stack_id}}</b><br>
				Logical id: <b>{{.logical_id}}</b><br><br>
				{{end}}

				<br>
				<br>

				Thanks for using Cecil!
				`,

		emailValues,
	)

	var newEmailSubject string
	if lease.StackName != "" {
		newEmailSubject = fmt.Sprintf("Stack (%v) terminated", lease.StackName)
	} else {
		newEmailSubject = fmt.Sprintf("Instance (%v) terminated", lease.InstanceID)
	}

	s.NotifierQueue.TaskQueue <- NotifierTask{
		AccountID: lease.AccountID, // this will also trigger send to Slack
		To:        owner.Email,
		Subject:   newEmailSubject,
		BodyHTML:  newEmailBody,
		BodyText:  newEmailBody,
		NotificationMeta: NotificationMeta{
			NotificationType: InstanceTerminated,
			LeaseUuid:        lease.UUID,
			InstanceId:       lease.InstanceID,
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
			"InstanceID", task.InstanceID,
		)
	} else {
		Logger.Info("Extending lease",
			"InstanceID", task.InstanceID,
		)
	}

	task.Lease.TokenOnce = uuid.NewV4().String() // invalidates all other URLs to renew/terminate/approve
	task.Lease.Alerted = false

	// define the lease duration
	leaseDuration, err := s.DefineLeaseDuration(task.Lease.AccountID, task.Lease.CloudAccountID)
	if err != nil {
		Logger.Error(
			"error while DefineLeaseDuration",
			"InstanceID", task.InstanceID,
			"err", err,
		)
		return err
	}

	if task.Approving {
		task.Lease.ExpiresAt = task.Lease.CreatedAt.Add(leaseDuration)
	} else {
		task.Lease.ExpiresAt = task.Lease.ExpiresAt.Add(leaseDuration)
	}

	s.DB.Save(&task.Lease)

	var owner Owner
	var ownerCount int64

	s.DB.Table("owners").Where(task.Lease.OwnerID).First(&owner).Count(&ownerCount)

	var newEmailBody string
	var newEmailSubject string
	var notificationType NotificationType

	var emailValues = map[string]interface{}{
		"owner_email":     owner.Email,
		"instance_id":     task.Lease.InstanceID,
		"instance_type":   task.Lease.InstanceType,
		"instance_region": task.Lease.Region,

		"instance_duration": task.Lease.ExpiresAt.Sub(task.Lease.CreatedAt).String(),

		"expires_at": task.Lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
	}

	if task.Lease.StackName != "" {
		emailValues["logical_id"] = task.Lease.LogicalID
		emailValues["stack_id"] = task.Lease.StackID
		emailValues["stack_name"] = task.Lease.StackName
	}

	if task.Approving {
		notificationType = LeaseApproved

		if task.Lease.StackName != "" {
			newEmailSubject = fmt.Sprintf("Stack (%v) lease approved", task.Lease.StackName)
		} else {
			newEmailSubject = fmt.Sprintf("Instance (%v) lease approved", task.Lease.InstanceID)
		}

		newEmailBody = CompileEmail(
			`
			{{if not .stack_name }}
			Hey {{.owner_email}}, the lease of instance <b>{{.instance_id}}</b>
				(of type <b>{{.instance_type}}</b>,
				on <b>{{.instance_region}}</b>) has been approved.
			{{end}}

					{{if .stack_name }}
					Hey {{.owner_email}}, the lease of the following stack (on <b>{{.instance_region}}</b>) has been approved:

							<br>
							<br>

					Stack name: <b>{{.stack_name}}</b><br>
					Stack id: <b>{{.stack_id}}</b><br>
					Logical id: <b>{{.logical_id}}</b><br><br>
					{{end}}



				<br>
				<br>

				The current expiration is
				<b>{{.expires_at}}</b> ({{.instance_duration}} after it's creation)

				<br>
				<br>

				Thanks for using Cecil!
				`,

			emailValues,
		)
	} else {
		notificationType = LeaseExtended

		if task.Lease.StackName != "" {
			newEmailSubject = fmt.Sprintf("Stack (%v) lease extended", task.Lease.StackName)
		} else {
			newEmailSubject = fmt.Sprintf("Instance (%v) lease extended", task.Lease.InstanceID)
		}

		newEmailBody = CompileEmail(
			`
			{{if not .stack_name }}

			Hey {{.owner_email}}, the lease of instance with id <b>{{.instance_id}}</b>
				(of type <b>{{.instance_type}}</b>,
				on <b>{{.instance_region}}</b>) has been extended.

			{{end}}

					{{if .stack_name }}
					Hey {{.owner_email}}, the lease of the following stack (on <b>{{.instance_region}}</b>) has been extended:

							<br>
							<br>

					Stack name: <b>{{.stack_name}}</b><br>
					Stack id: <b>{{.stack_id}}</b><br>
					Logical id: <b>{{.logical_id}}</b><br><br>
					{{end}}



				<br>
				<br>

				The current expiration is
				<b>{{.expires_at}}</b> ({{.instance_duration}} after it's creation)

				<br>
				<br>

				Thanks for using Cecil!
				`,

			emailValues,
		)
	}

	s.NotifierQueue.TaskQueue <- NotifierTask{
		AccountID: task.Lease.AccountID, // this will also trigger send to Slack
		To:        owner.Email,
		Subject:   newEmailSubject,
		BodyHTML:  newEmailBody,
		BodyText:  newEmailBody,
		NotificationMeta: NotificationMeta{
			NotificationType: notificationType,
			LeaseUuid:        task.Lease.UUID,
			InstanceId:       task.Lease.InstanceID,
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

	// if there is a SlackInstance for the Account,
	// send in a goroutine a message to that SlackInstance.
	go func() {
		if task.AccountID == 0 {
			return
		}
		slackIns, err := s.SlackInstanceByID(task.AccountID)
		if err != nil {
			Logger.Warn("SlackInstanceByID", "err", err)
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
			Logger.Warn("MailerInstanceByID", "err", err)
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
	message.AddHeader(X_CECIL_LEASE_UUID, task.NotificationMeta.LeaseUuid)
	message.AddHeader(X_CECIL_INSTANCE_ID, task.NotificationMeta.InstanceId)
	message.AddHeader(X_CECIL_VERIFICATION_TOKEN, task.NotificationMeta.VerificationToken)

	//message.SetTracking(true)
	if task.DeliverAfter > 0 {
		message.SetDeliveryTime(time.Now().Add(task.DeliverAfter))
	}

	message.SetHtml(task.BodyHTML)

	err := retry(10, time.Second*5, func() error {
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
