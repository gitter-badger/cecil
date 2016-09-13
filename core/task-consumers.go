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

func (s *Service) NewLeaseQueueConsumer(t interface{}) error {

	return nil
}

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
	s.DB.Table("leases").Where(&Lease{InstanceID: task.InstanceID, AWSAccountID: task.AWSID}).First(&lease).Count(&leasesFound)

	if leasesFound != 1 {
		logger.Warn("Found multiple leases for deletion", "count", leasesFound)
		return fmt.Errorf("Found multiple leases for deletion", "count", leasesFound)
	}

	lease.Terminated = true
	// TODO: use the ufficial time of termination, from th sqs message
	// lease.TerminatedAt = time.Now().UTC()
	s.DB.Save(&lease)

	var owner Owner
	var ownerCount int64

	s.DB.Table("owners").Where(lease.OwnerID).First(&owner).Count(&ownerCount)

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

func (s *Service) RenewerQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(RenewerTask)
	// TODO: check whether fields are non-null and valid

	_ = task

	return nil
}

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

	message.SetTracking(true)
	//message.SetDeliveryTime(time.Now().Add(24 * time.Hour))
	message.SetHtml(task.BodyHTML)
	_, id, err := s.Mailer.Send(message)
	if err != nil {
		logger.Error("Error while sending email", "error", err)
		return err
	}
	_ = id

	return nil
}
