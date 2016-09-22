package core

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gopkg.in/mailgun/mailgun-go.v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
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
	logger.Info("TerminatorQueueConsumer",
		"task", task,
	)

	// need:
	// region
	// roleARN
	// external ID

	var cloudAccount CloudAccount
	var leaseCloudOwnerCount int64
	s.DB.Model(&task.Lease).Related(&cloudAccount).Count(&leaseCloudOwnerCount)
	//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&leaseCloudOwnerCount)
	if leaseCloudOwnerCount == 0 {
		// TODO: notify admin; something fishy is going on.
		logger.Warn("leaseCloudOwnerCount == 0")
		return fmt.Errorf("leaseCloudOwnerCount == 0")
	}

	// assume role
	assumedConfig := &aws.Config{
		Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
			Client: sts.New(s.AWS.Session, &aws.Config{Region: aws.String(task.Region)}),
			RoleARN: fmt.Sprintf(
				"arn:aws:iam::%v:role/%v",
				cloudAccount.AWSID,
				viper.GetString("ForeignRoleName"),
			),
			RoleSessionName: uuid.NewV4().String(),
			ExternalID:      aws.String(cloudAccount.ExternalID),
			ExpiryWindow:    3 * time.Minute,
		}),
	}

	assumedService := session.New(assumedConfig)

	ec2Service := s.EC2(assumedService, task.Region)

	terminateInstanceParams := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{ // Required
			aws.String(task.InstanceID),
		},
	}
	terminateInstanceResponse, err := ec2Service.TerminateInstances(terminateInstanceParams)
	_ = terminateInstanceResponse

	logger.Info("TerminateInstances", "response", terminateInstanceResponse)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.

		// TODO: cleaner way to do this?  cloudAccount.Account would be nice .. gorma provides this
		var account Account
		s.DB.First(&account, cloudAccount.AccountID)

		recipientEmail := account.Email

		s.sendMisconfigurationNotice(err, recipientEmail)
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
	logger.Info("Marking lease as terminated",
		"InstanceID", task.InstanceID,
	)

	var lease Lease
	var leasesFound int64
	s.DB.Table("leases").Where(&Lease{
		InstanceID:   task.InstanceID,
		AWSAccountID: task.AWSID,
		Terminated:   false,
	}).First(&lease).Count(&leasesFound)

	if leasesFound == 0 {
		logger.Warn("Lease for deletion not found", "count", leasesFound, "instanceID", task.InstanceID)
		return fmt.Errorf("Lease for deletion not found: %v=%v", "count", leasesFound)
	}
	if leasesFound > 1 {
		logger.Warn("Found multiple leases for deletion", "count", leasesFound)
		return fmt.Errorf("Found multiple leases for deletion: %v=%v", "count", leasesFound)
	}

	lease.Terminated = true
	lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve

	// TODO: use the ufficial time of termination, from th sqs message, because if erminated via link, the termination time is not expiresAt
	// lease.TerminatedAt = time.Now().UTC()
	s.DB.Save(&lease)

	var owner Owner
	var ownerCount int64

	s.DB.Table("owners").Where(lease.OwnerID).First(&owner).Count(&ownerCount)

	if ownerCount != 1 {
		logger.Warn("LeaseTerminatedQueueConsumer: ownerCount is not 1", "count", ownerCount)
		return fmt.Errorf("LeaseTerminatedQueueConsumer: ownerCount is not 1: %v=%v", "count", ownerCount)
	}

	newEmailBody := compileEmail(
		`Hey {{.owner_email}}, instance with id <b>{{.instance_id}}</b>
				(of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>) has been terminated at 
				<b>{{.terminated_at}}</b> ({{.instance_duration}} after it's creation)

				<br>
				<br>
				
				Thanks for using ZeroCloud!
				`,

		map[string]interface{}{
			"owner_email":     owner.Email,
			"instance_id":     lease.InstanceID,
			"instance_type":   lease.InstanceType,
			"instance_region": lease.Region,

			"instance_duration": lease.ExpiresAt.Sub(lease.CreatedAt).String(),

			"terminated_at": lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
		},
	)
	s.NotifierQueue.TaskQueue <- NotifierTask{
		From:     ZCMailerFromAddress,
		To:       owner.Email,
		Subject:  fmt.Sprintf("Instance (%v) terminated", lease.InstanceID),
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
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
		logger.Info("Approving lease",
			"InstanceID", task.InstanceID,
		)
	} else {
		logger.Info("Extending lease",
			"InstanceID", task.InstanceID,
		)
	}

	task.Lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve
	task.Lease.Alerted = false

	if task.Approving {
		task.Lease.ExpiresAt = task.Lease.CreatedAt.Add(task.ExtendBy)
	} else {
		task.Lease.ExpiresAt = task.Lease.ExpiresAt.Add(task.ExtendBy)
	}

	s.DB.Save(&task.Lease)

	var owner Owner
	var ownerCount int64

	s.DB.Table("owners").Where(task.Lease.OwnerID).First(&owner).Count(&ownerCount)

	var newEmailBody string
	var newEmailSubject string
	if task.Approving {
		newEmailSubject = fmt.Sprintf("Instance (%v) lease approved", task.Lease.InstanceID)
		newEmailBody = compileEmail(
			`Hey {{.owner_email}}, the lease of instance <b>{{.instance_id}}</b>
				(of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>) has been approved.

				<br>
				<br>

				The current expiration is 
				<b>{{.expires_at}}</b> ({{.instance_duration}} after it's creation)

				<br>
				<br>
				
				Thanks for using ZeroCloud!
				`,

			map[string]interface{}{
				"owner_email":     owner.Email,
				"instance_id":     task.Lease.InstanceID,
				"instance_type":   task.Lease.InstanceType,
				"instance_region": task.Lease.Region,

				"instance_duration": task.Lease.ExpiresAt.Sub(task.Lease.CreatedAt).String(),

				"expires_at": task.Lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
			},
		)
	} else {
		newEmailSubject = fmt.Sprintf("Instance (%v) lease extended", task.Lease.InstanceID)
		newEmailBody = compileEmail(
			`Hey {{.owner_email}}, the lease of instance with id <b>{{.instance_id}}</b>
				(of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>) has been extended.

				<br>
				<br>

				The current expiration is 
				<b>{{.expires_at}}</b> ({{.instance_duration}} after it's creation)

				<br>
				<br>
				
				Thanks for using ZeroCloud!
				`,

			map[string]interface{}{
				"owner_email":     owner.Email,
				"instance_id":     task.Lease.InstanceID,
				"instance_type":   task.Lease.InstanceType,
				"instance_region": task.Lease.Region,

				"instance_duration": task.Lease.ExpiresAt.Sub(task.Lease.CreatedAt).String(),

				"expires_at": task.Lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
			},
		)
	}

	s.NotifierQueue.TaskQueue <- NotifierTask{
		From:     ZCMailerFromAddress,
		To:       owner.Email,
		Subject:  newEmailSubject,
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
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
	logger.Info("Sending EMAIL",
		"to", task.To,
	)

	message := mailgun.NewMessage(
		task.From,
		task.Subject,
		task.BodyText,
		task.To,
	)

	//message.SetTracking(true)
	//message.SetDeliveryTime(time.Now().Add(24 * time.Hour))
	message.SetHtml(task.BodyHTML)

	err := retry(10, time.Second*5, func() error {
		var err error
		_, _, err = s.Mailer.Send(message)
		return err
	})
	if err != nil {
		logger.Error("Error while sending email", "error", err)
		return err
	}

	return nil
}
