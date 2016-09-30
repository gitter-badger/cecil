package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

var ErrEnvelopeIsSubscriptionConfirmation error = errors.New("ErrEnvelopeIsSubscriptionConfirmation")

// Transmission contains the SQS message and everything else
// needed to complete the operations triggered by the message
type Transmission struct {
	s            *Service
	Message      SQSMessage
	subscribeURL string

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

// parseSQSTransmission parses a raw SQS message into a Transmission
func (s *Service) parseSQSTransmission(rawMessage *sqs.Message, queueURL string) (*Transmission, error) {
	var newTransmission Transmission = Transmission{}
	newTransmission.s = s

	newTransmission.deleteMessageFromQueueParams = &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),                  // Required
		ReceiptHandle: aws.String(*rawMessage.ReceiptHandle), // Required
	}

	// parse the envelope
	var envelope SQSEnvelope
	err := json.Unmarshal([]byte(*rawMessage.Body), &envelope)
	if err != nil {
		// NOTE: here we return an empty &Transmission{} because the unmarshaling failed
		return &Transmission{}, err
	}

	// TODO: move the Arn parsing and validation in another function
	topicArn := strings.Split(envelope.TopicArn, ":")
	newTransmission.Topic.Region = topicArn[3]
	newTransmission.Topic.AWSID = topicArn[4]
	topicName := topicArn[5]

	if topicName != s.AWS.Config.SNSTopicName {
		// NOTE: here we return &newTransmission (and not an empty &Transmission{}) because
		// &newTransmission contains values that will be used
		return &newTransmission, fmt.Errorf("the SNS topic name is not %v: %v", s.AWS.Config.SNSTopicName, topicName)
	}

	if newTransmission.Topic.Region == "" {
		return &newTransmission, fmt.Errorf("newTransmission.Topic.Region is empty")
	}
	if newTransmission.Topic.AWSID == "" {
		return &newTransmission, fmt.Errorf("newTransmission.Topic.AWSID is empty")
	}

	//check whether someone with this aws adminAccount id is registered at zerocloud
	err = newTransmission.FetchCloudAccount()
	if err != nil {
		// TODO: notify admin; something fishy is going on.
		logger.Warn("originator is not registered", "AWSID", newTransmission.Topic.AWSID)
		return &newTransmission, err
	}

	// check whether the cloud account has an admin account
	err = newTransmission.FetchAdminAccount()
	if err != nil {
		// TODO: notify admin; something fishy is going on.
		logger.Warn("transmission: Error while retrieving admin account", "error", err)
		return &newTransmission, err
	}

	logger.Info("adminAccount",
		"adminAccount", newTransmission.AdminAccount,
	)

	// TODO: check if the user is signed up before confirming the subscription

	if envelope.Type == "SubscriptionConfirmation" {
		newTransmission.subscribeURL = envelope.SubscribeURL
		return &newTransmission, ErrEnvelopeIsSubscriptionConfirmation
	}

	// parse the message
	err = json.Unmarshal([]byte(envelope.Message), &newTransmission.Message)
	if err != nil {
		return &newTransmission, err
	}

	return &newTransmission, nil
}

func (t *Transmission) ConfirmSQSSubscription() error {

	confirmationURL, err := url.Parse(t.subscribeURL)
	if err != nil {
		return err
	}

	if confirmationURL.Host[len(confirmationURL.Host)-13:] != "amazonaws.com" {
		return fmt.Errorf("subscribeURL host is NOT amazonaws.com: %v", confirmationURL.Host)
	}

	resp, err := http.Get(confirmationURL.String())
	if err != nil {
		return err
	}
	// TODO: parse the response body
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("response statusCode is not 200: %v", resp.StatusCode)
	}

	// region added successfully; send confirmation email
	newEmailBody := compileEmail(
		`Hey {{.admin_email}}, the region <b>{{.region_name}}</b> has been successfully setup!
		<br>
		<br>
		From now on, all instances on this region will be monitored.
		<br>
		<br>
		Thanks!
		`,
		map[string]interface{}{
			"admin_email": t.AdminAccount.Email,
			"region_name": t.Topic.Region,
		},
	)

	t.s.NotifierQueue.TaskQueue <- NotifierTask{
		From:     t.s.Mailer.FromAddress,
		To:       t.AdminAccount.Email,
		Subject:  fmt.Sprintf("Region %v has been setup", t.Topic.Region),
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
	}

	logger.Info("ConfirmSQSSubscription", "subscribeURL", confirmationURL.String())
	return nil

	/*
		These are just some of the possible responses:

		<ErrorResponse xmlns="http://sns.amazonaws.com/doc/2010-03-31/">
			<Error>
				<Type>Sender</Type>
				<Code>InvalidParameter</Code>
				<Message>Invalid parameter: Token</Message>
			</Error>
			<RequestId>76c87c52-03bf-55c2-9db6-2c3409449b1e</RequestId>
		</ErrorResponse>



		<ConfirmSubscriptionResponse xmlns="http://sns.amazonaws.com/doc/2010-03-31/">
			<ConfirmSubscriptionResult>
				<SubscriptionArn>
					arn:aws:sns:ap-northeast-1:012345678910:ZeroCloudTopic:c1e03965-deec-4f18-aa2b-76fbb4451a04
				</SubscriptionArn>
			</ConfirmSubscriptionResult>
			<ResponseMetadata>
				<RequestId>83fe317d-1c8a-5b33-8058-611b16523b90</RequestId>
			</ResponseMetadata>
		</ConfirmSubscriptionResponse>
	*/
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
		_ = t.deleteMessageFromQueueParams
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
	if cloudOwnerCount == 0 || t.CloudAccount.AWSID != t.Topic.AWSID {
		return fmt.Errorf("No cloudAccount for AWSID %v", t.Topic.AWSID)
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
				t.s.AWS.Config.ForeignIAMRoleName,
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

	if t.Instance.Placement == nil {
		return fmt.Errorf("EC2Instance has nil Placement field")
	}

	if t.Instance.Placement.AvailabilityZone == nil {
		return fmt.Errorf("EC2Instance has nil Placement.AvailabilityZone field")
	}

	// TODO: is this always valid?
	// TODO: use pointer or by value?
	az := *t.Instance.Placement.AvailabilityZone
	t.instanceRegion = az[:len(az)-1]

	return nil
}

func (t *Transmission) InstanceIsTerminated() bool {

	if t.Instance.State == nil {
		logger.Warn("t.Instance.State == nil")
		return false
	}
	if t.Instance.State.Name == nil {
		logger.Warn("t.Instance.State.Name == nil")
		return false
	}

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
		ownerTag, err := t.s.Mailer.Client.ValidateEmail(*tag.Value)
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
	t.leaseDuration = time.Duration(t.s.Config.Lease.Duration)

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

	return t.activeLeaseCount >= int64(t.s.Config.Lease.MaxPerOwner) && t.s.Config.Lease.MaxPerOwner >= 0
}

// InstanceLaunchTimeUTC is a shortcut to t.Instance.LaunchTime
func (t *Transmission) InstanceLaunchTimeUTC() time.Time {
	if t.Instance.LaunchTime == nil {
		logger.Warn("t.Instance.LaunchTime == nil")
		return time.Now().UTC()
	}
	return t.Instance.LaunchTime.UTC()
}

// InstanceType is a shortcut to *t.Instance.InstanceType
func (t *Transmission) InstanceType() string {
	if t.Instance.InstanceType == nil {
		logger.Warn("t.Instance.InstanceType == nil")
		return "unknown"
	}
	return *t.Instance.InstanceType
}

// InstanceId is a shortuct to *t.Instance.InstanceId
func (t *Transmission) InstanceId() string {
	if t.Instance.InstanceId == nil {
		logger.Warn("t.Instance.InstanceId == nil")
		return "i-unknown"
	}
	return *t.Instance.InstanceId
}

// AvailabilityZone is a shortcut to *t.Instance.Placement.AvailabilityZone
func (t *Transmission) AvailabilityZone() string {
	if t.Instance.Placement.AvailabilityZone == nil {
		logger.Warn("t.Instance.Placement.AvailabilityZone == nil")
		return "somewhere-unknown"
	}
	return *t.Instance.Placement.AvailabilityZone
}
