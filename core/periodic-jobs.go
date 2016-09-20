package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// @@@@@@@@@@@@@@@ Periodic Jobs @@@@@@@@@@@@@@@

func (s *Service) EventInjestorJob() error {
	// TODO: verify event origin (must be aws, not someone else)

	queueURL := fmt.Sprintf("https://sqs.%v.amazonaws.com/%v/%v",
		viper.GetString("AWS_REGION"),
		viper.GetString("AWS_ACCOUNT_ID"),
		viper.GetString("SQSQueueName"),
	)

	receiveMessageParams := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL), // Required
		//MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout: aws.Int64(3), // should be higher, like 10 (seconds), the time to finish doing everything
		WaitTimeSeconds:   aws.Int64(3),
	}

	logger.Info("EventInjestorJob(): Polling SQS", "queue", queueURL)
	receiveMessageResponse, err := s.AWS.SQS.ReceiveMessage(receiveMessageParams)

	if err != nil {
		return fmt.Errorf("EventInjestorJob() error: %v", err)
	}

	logger.Info("SQSmessages",
		"count", len(receiveMessageResponse.Messages),
	)

	for messageIndex := range receiveMessageResponse.Messages {

		transmission, err := s.parseSQSTransmission(receiveMessageResponse.Messages[messageIndex], queueURL)
		if err != nil {
			logger.Warn("Error parsing transmission", "error", err)
			continue
		}

		logger.Info("Parsed sqs message", "message", transmission.Message)

		if !transmission.TopicAndInstanceHaveSameOwner() {
			// the originating SNS topic and the instance have different owners (different AWS accounts)
			// TODO: notify zerocloud admin
			logger.Warn("topicAWSID != instanceOriginatorID", "topicAWSID", transmission.Topic.AWSID, "instanceOriginatorID", transmission.Message.Account)
			continue
		}

		// consider only pending and terminated status messages; ignore the rest
		if !transmission.MessageIsRelevant() {
			logger.Warn("Ignoring and removing message", "message.Detail.State", transmission.Message.Detail.State)
			err := transmission.DeleteMessage()
			if err != nil {
				logger.Warn("DeleteMessage", "error", err)
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

func (s *Service) AlerterJob() error {
	// find lease that expire in 24 hours
	// find owner
	// create links to extend and terminate lease
	// mark as alerted = true
	// registed new lease's token_once
	// compose email with link to extend and terminate lease
	// send email

	var expiringLeases []Lease
	var expiringLeasesCount int64

	s.DB.Table("leases").
		Where("expires_at < ?",
			time.Now().UTC().Add(ZCDefaultForewarningBeforeExpiry),
		).
		Not("terminated", true).
		Not("alerted", true).
		Find(&expiringLeases).
		Count(&expiringLeasesCount)

	logger.Info("AlerterJob(): Expiring leases", "count", expiringLeasesCount)

	// TODO: create ExpiringLeaseQueue and pass to it this task

	for _, expiringLease := range expiringLeases {

		logger.Info("Expiring lease",
			"instanceID", expiringLease.InstanceID,
			"leaseID", expiringLease.ID,
		)

		var owner Owner
		var ownerCount int64

		s.DB.Table("owners").Where(expiringLease.OwnerID).First(&owner).Count(&ownerCount)

		if ownerCount != 1 {
			logger.Warn("AlerterJob: ownerCount is not 1", "count", ownerCount)
			continue
		}

		// these will be used to compose the urls and verify the requests
		token_once := uuid.NewV4().String() // one-time token

		expiringLease.TokenOnce = token_once
		expiringLease.Alerted = true

		s.DB.Save(&expiringLease)

		// URL to extend lease
		action := "extend"
		signature, err := s.sign(expiringLease.UUID, expiringLease.InstanceID, action, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while signing: %v", err)
		}
		extend_url := fmt.Sprintf("http://0.0.0.0:8080/email_action/leases/%s/%s/%s?t=%s&s=%s",
			expiringLease.UUID,
			expiringLease.InstanceID,
			action,
			token_once,
			base64.URLEncoding.EncodeToString(signature),
		)

		// URL to terminate lease
		action = "terminate"
		signature, err = s.sign(expiringLease.UUID, expiringLease.InstanceID, action, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while signing")
		}
		terminate_url := fmt.Sprintf("http://0.0.0.0:8080/email_action/leases/%s/%s/%s?t=%s&s=%s",
			expiringLease.UUID,
			expiringLease.InstanceID,
			action,
			token_once,
			base64.URLEncoding.EncodeToString(signature),
		)

		newEmailBody := compileEmail(
			`Hey {{.owner_email}}, instance <b>{{.instance_id}}</b>
				(of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>) is expiring.

				<br>
				<br>

				The instance will expire on <b>{{.termination_time}}</b> ({{.instance_duration}} after it's creation).

				<br>
				<br>

				The instance was created on <b>{{.instance_created_at}}</b>.
				
				<br>
				<br>
				
				Terminate immediately:
				<br>
				<br>
				<a href="{{.instance_terminate_url}}" target="_blank">Click here to terminate</a>

				<br>
				<br>

				Extend by <b>{{.extend_by}}</b>:
				<br>
				<br>
				<a href="{{.instance_extend_url}}" target="_blank">Click here to extend</a>

				<br>
				<br>
				Thanks for using ZeroCloud!
				`,

			map[string]interface{}{
				"owner_email":     owner.Email,
				"instance_id":     expiringLease.InstanceID,
				"instance_type":   expiringLease.InstanceType,
				"instance_region": expiringLease.Region,

				"instance_created_at": expiringLease.CreatedAt.Format("2006-01-02 15:04:05 GMT"),
				"extend_by":           ZCDefaultLeaseDuration.String(),

				"termination_time":  expiringLease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
				"instance_duration": expiringLease.ExpiresAt.Sub(expiringLease.CreatedAt).String(),

				"instance_terminate_url": terminate_url,
				"instance_extend_url":    extend_url,
			},
		)

		s.NotifierQueue.TaskQueue <- NotifierTask{
			From:     ZCMailerFromAddress,
			To:       owner.Email,
			Subject:  fmt.Sprintf("Instance (%v) will expire soon", expiringLease.InstanceID),
			BodyHTML: newEmailBody,
			BodyText: newEmailBody,
		}
	}

	return nil
}

func (s *Service) SentencerJob() error {

	var expiredLeases []Lease
	var expiredLeasesCount int64

	s.DB.Table("leases").Where("expires_at < ?", time.Now().UTC()).Not("terminated", true).Find(&expiredLeases).Count(&expiredLeasesCount)

	logger.Info("SentencerJob(): Expired leases", "count", expiredLeasesCount)

	for _, expiredLease := range expiredLeases {
		logger.Info("expired lease",
			"instanceID", expiredLease.InstanceID,
			"leaseID", expiredLease.ID,
		)
		s.TerminatorQueue.TaskQueue <- TerminatorTask{Lease: expiredLease}
	}

	return nil
}

func (s *Service) sendMisconfigurationNotice(err error, emailRecipient string) {
	newEmailBody := compileEmail(
		`Hey it appears that ZeroCloud is mis-configured.
		<br>
		<br>
		Error:
		<br>
		{{.err}}`,
		map[string]interface{}{
			"err": err,
		},
	)

	s.NotifierQueue.TaskQueue <- NotifierTask{
		From:     ZCMailerFromAddress,
		To:       emailRecipient,
		Subject:  "ZeroCloud configuration problem",
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
	}

}

type Transmission struct {
	s       *Service
	Message SQSMessage

	Topic struct {
		Region string
		AWSID  string
	}
	deleteMessageFromQueueParams *sqs.DeleteMessageInput

	CloudAccount CloudAccount
	AdminAccount Account

	assumedService            *session.Session
	ec2Service                ec2iface.EC2API
	describeInstancesResponse *ec2.DescribeInstancesOutput

	Instance           ec2.Instance
	instanceRegion     string
	externalOwnerEmail string
	owner              Owner
	leaseDuration      time.Duration
	activeLeaseCount   int64
}

func (s *Service) parseSQSTransmission(rawMessage *sqs.Message, queueURL string) (*Transmission, error) {
	var newTransmission Transmission = Transmission{}
	newTransmission.s = s

	// parse the envelope
	var envelope SQSEnvelope
	err := json.Unmarshal([]byte(*rawMessage.Body), &envelope)
	if err != nil {
		return &Transmission{}, err
	}

	newTransmission.deleteMessageFromQueueParams = &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),                  // Required
		ReceiptHandle: aws.String(*rawMessage.ReceiptHandle), // Required
	}

	// parse the message
	err = json.Unmarshal([]byte(envelope.Message), &newTransmission.Message)
	if err != nil {
		return &Transmission{}, err
	}

	// extract some values
	// TODO: check whether these values are not empty
	/// transmission.Topic.Region
	/// transmission.Topic.AWSID
	topicArn := strings.Split(envelope.TopicArn, ":")
	newTransmission.Topic.Region = topicArn[3]
	newTransmission.Topic.AWSID = topicArn[4]

	return &newTransmission, nil
}

// the originating SNS topic and the instance have different owners (different AWS accounts)
func (t *Transmission) TopicAndInstanceHaveSameOwner() bool {
	instanceOriginatorID := t.Message.Account

	return t.Topic.AWSID == instanceOriginatorID
}

// consider only pending and terminated status messages; ignore the rest
func (t *Transmission) MessageIsRelevant() bool {
	return t.Message.Detail.State == ec2.InstanceStateNamePending ||
		t.Message.Detail.State == ec2.InstanceStateNameTerminated
}

func (t *Transmission) DeleteMessage() error {
	// remove message from queue
	return retry(5, time.Duration(3*time.Second), func() error {
		var err error
		// TODO:
		_, err = t.s.AWS.SQS.DeleteMessage(t.deleteMessageFromQueueParams)
		return err
	})
}

//check whether someone with this aws adminAccount id is registered at zerocloud
func (t *Transmission) FetchCloudAccount() error {
	var cloudOwnerCount int64
	t.s.DB.Where(&CloudAccount{AWSID: t.Topic.AWSID}).
		First(&t.CloudAccount).
		Count(&cloudOwnerCount)
	if cloudOwnerCount == 0 {
		return fmt.Errorf("No cloud account for AWSID %v", t.Topic.AWSID)
	}
	if cloudOwnerCount > 1 {
		return fmt.Errorf("Too many (%v) CloudAccounts for AWSID %v", cloudOwnerCount, t.Topic.AWSID)
	}
	return nil
}

// check whether the cloud account has an admin account
func (t *Transmission) FetchAdminAccount() error {
	var cloudAccountAdminCount int64
	t.s.DB.Model(&t.CloudAccount).Related(&t.AdminAccount).Count(&cloudAccountAdminCount)
	//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&cloudAccountAdminCount)
	if cloudAccountAdminCount == 0 {
		return fmt.Errorf("No admin for CloudAccount.  CloudAccount.ID %v", t.CloudAccount.ID)
	}
	if cloudAccountAdminCount > 1 {
		return fmt.Errorf("Too many (%v) admins for CloudAccount %v", cloudAccountAdminCount, t.CloudAccount.ID)
	}
	return nil
}

func (t *Transmission) CreateAssumedService() error {
	assumedConfig := &aws.Config{
		Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
			Client: sts.New(t.s.AWS.Session, &aws.Config{Region: aws.String(t.Topic.Region)}),
			RoleARN: fmt.Sprintf(
				"arn:aws:iam::%v:role/%v",
				t.Topic.AWSID,
				viper.GetString("ForeignRoleName"),
			),
			RoleSessionName: uuid.NewV4().String(),
			ExternalID:      aws.String(t.CloudAccount.ExternalID),
			ExpiryWindow:    60 * time.Second,
		}),
	}

	t.assumedService = session.New(assumedConfig)

	return nil
}

func (t *Transmission) CreateAssumedEC2Service() error {
	t.ec2Service = t.s.EC2(t.assumedService, t.Topic.Region)
	return nil
}

func (t *Transmission) DescribeInstance() error {
	paramsDescribeInstance := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(t.Message.Detail.InstanceID),
		},
	}
	var err error
	t.describeInstancesResponse, err = t.ec2Service.DescribeInstances(paramsDescribeInstance)
	return err
}

// check whether the instance specified in the event message exists on aws
func (t *Transmission) InstanceExists() bool {
	if len(t.describeInstancesResponse.Reservations) == 0 {
		// logger.Warn("len(describeInstancesResponse.Reservations) == 0")
		return false
	}
	if len(t.describeInstancesResponse.Reservations[0].Instances) == 0 {
		// logger.Warn("len(describeInstancesResponse.Reservations[0].Instances) == 0")
		return false
	}
	return true
}

func (t *Transmission) FetchInstance() error {
	// TODO: merge the preceding check operations here
	t.Instance = *t.describeInstancesResponse.Reservations[0].Instances[0]

	if *t.Instance.InstanceId != t.Message.Detail.InstanceID {
		return fmt.Errorf("instance.InstanceId != message.Detail.InstanceID")
	}
	return nil
}

func (t *Transmission) ComputeInstanceRegion() error {

	// TODO: is this always valid?
	// TODO: use pointer or by value?
	az := *t.Instance.Placement.AvailabilityZone
	t.instanceRegion = az[:len(az)-1]

	return nil
}

func (t *Transmission) InstanceIsTerminated() bool {
	return *t.Instance.State.Name == ec2.InstanceStateNameTerminated
}

func (t *Transmission) InstanceIsPendingOrRunning() bool {
	return *t.Instance.State.Name == ec2.InstanceStateNamePending ||
		*t.Instance.State.Name == ec2.InstanceStateNameRunning
}

// IsNew: check whether a lease with the same instanceID exists
func (t *Transmission) LeaseIsNew() bool {
	var leaseCount int64
	t.s.DB.Table("leases").Where(&Lease{InstanceID: t.Message.Detail.InstanceID}).Count(&leaseCount)

	return leaseCount == 0
}

func (t *Transmission) InstanceHasGoodOwnerTag() bool {
	// InstanceHasTags: check whether instance has tags
	if len(t.Instance.Tags) == 0 {
		return false
	}
	//logger.Warn("len(instance.Tags) == 0")

	// InstanceHasOwnerTag: check whether the instance has an zerocloudowner tag
	for _, tag := range t.Instance.Tags {
		if strings.ToLower(*tag.Key) != "zerocloudowner" {
			continue
		}

		// OwnerTagValueIsValid: check whether the zerocloudowner tag is a valid email
		ownerTag, err := t.s.Mailer.ValidateEmail(*tag.Value)
		if err != nil {
			logger.Warn("ValidateEmail", "error", err)
			// TODO: send notification to admin
			return false
		}
		if !ownerTag.IsValid {
			return false
			// TODO: notify admin: "Warning: zerocloudowner tag email not valid" (DO NOT INCLUDE IT IN THE EMAIL, OR HTML-ESCAPE IT)
		}
		// fmt.Printf("Parts local_part=%s domain=%s display_name=%s", ownerTag.Parts.LocalPart, ownerTag.Parts.Domain, ownerTag.Parts.DisplayName)
		t.externalOwnerEmail = ownerTag.Address
	}
	return true
}

// ExternalOwnerIsWhitelisted: check whether the owner email in the tag is a whitelisted owner email
func (t *Transmission) ExternalOwnerIsWhitelisted() bool {
	// TODO: select Owner by email, cloudaccountid, and region?
	if t.externalOwnerEmail == "" {
		return false
	}
	var ownerCount int64
	t.s.DB.Table("owners").Where(&Owner{Email: t.externalOwnerEmail, CloudAccountID: t.CloudAccount.ID}).Count(&ownerCount)
	if ownerCount != 1 {
		return false
	}
	return true
}

func (t *Transmission) SetExternalOwnerAsOwner() error {
	var owners []Owner
	var ownerCount int64
	t.s.DB.Table("owners").Where(&Owner{Email: t.externalOwnerEmail, CloudAccountID: t.CloudAccount.ID}).Find(&owners).Count(&ownerCount)
	if ownerCount != 1 {
		return fmt.Errorf("Too many external owners with externalOwnerEmail %v", t.externalOwnerEmail)
	}
	t.owner = owners[0]
	return nil
}
func (t *Transmission) SetAdminAsOwner() error {
	var owners []Owner
	var ownerCount int64
	t.s.DB.Table("owners").Where(&Owner{Email: t.AdminAccount.Email, CloudAccountID: t.CloudAccount.ID}).Find(&owners).Count(&ownerCount)
	if ownerCount != 1 {
		return fmt.Errorf("Too many admin owners with email %v", t.AdminAccount.Email)
	}
	t.owner = owners[0]
	return nil
}

func (t *Transmission) DefineLeaseDuration() {
	t.leaseDuration = time.Duration(ZCDefaultLeaseDuration)

	if t.AdminAccount.DefaultLeaseDuration > 0 {
		t.leaseDuration = time.Duration(t.AdminAccount.DefaultLeaseDuration)
	}
	if t.CloudAccount.DefaultLeaseDuration > 0 {
		t.leaseDuration = time.Duration(t.CloudAccount.DefaultLeaseDuration)
	}
}

func (t *Transmission) LeaseNeedsApproval() bool {
	var leases []Lease

	t.s.DB.Table("leases").Where(&Lease{
		OwnerID:        t.owner.ID,
		CloudAccountID: t.CloudAccount.ID,
		Terminated:     false,
	}).Find(&leases).Count(&t.activeLeaseCount)
	//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&activeLeaseCount)

	return t.activeLeaseCount >= ZCDefaultMaxLeasesPerOwner
}
