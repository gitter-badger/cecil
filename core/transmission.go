package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// ErrorEnvelopeIsSubscriptionConfirmation is an error that triggers automatic
// subscription between SNS and SQS.
var ErrorEnvelopeIsSubscriptionConfirmation = errors.New("ErrEnvelopeIsSubscriptionConfirmation")

type StackInfo struct {
	LogicalID string
	StackID   string
	StackName string
}

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

	assumedService               *session.Session
	assumedEC2Service            ec2iface.EC2API
	assumedCloudformationService cloudformationiface.CloudFormationAPI
	describeInstancesResponse    *ec2.DescribeInstancesOutput

	Instance           ec2.Instance
	StackResources     []*cloudformation.StackResource
	instanceRegion     string
	externalOwnerEmail string
	owner              Owner
	leaseDuration      time.Duration
	activeLeaseCount   int64
	StackInfo          *StackInfo
}

// parseSQSTransmission parses a raw SQS message into a Transmission
func (s *Service) parseSQSTransmission(rawMessage *sqs.Message, queueURL string) (*Transmission, error) {

	// Record the eent
	if err := s.EventRecord.StoreSQSMessage(rawMessage); err != nil {
		Logger.Warn("Error storing SQS message", "err", err, "msg", rawMessage)
	}

	newTransmission := Transmission{}
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
	if len(topicArn) < 6 {
		return &Transmission{}, fmt.Errorf("cannot parse topic Arn: %v", envelope.TopicArn)
	}
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

	//check whether someone with this aws adminAccount id is registered at cecil
	err = newTransmission.FetchCloudAccount()
	if err != nil {
		// TODO: notify admin; something fishy is going on.
		Logger.Warn("originator is not registered", "AWSID", newTransmission.Topic.AWSID)
		return &newTransmission, err
	}

	// check whether the cloud account has an admin account
	err = newTransmission.FetchAdminAccount()
	if err != nil {
		// TODO: notify admin; something fishy is going on.
		Logger.Warn("transmission: Error while retrieving admin account", "err", err)
		return &newTransmission, err
	}

	Logger.Info("adminAccount",
		"adminAccount", newTransmission.AdminAccount,
	)

	// TODO: check if the user is signed up before confirming the subscription

	if envelope.Type == "SubscriptionConfirmation" {
		newTransmission.subscribeURL = envelope.SubscribeURL
		return &newTransmission, ErrorEnvelopeIsSubscriptionConfirmation
	}

	// parse the message
	err = json.Unmarshal([]byte(envelope.Message), &newTransmission.Message)
	if err != nil {
		return &newTransmission, err
	}

	return &newTransmission, nil
}

// ConfirmSQSSubscription is used to confirm the subscription of SQS to SNS.
func (t *Transmission) ConfirmSQSSubscription() error {

	confirmationURL, err := url.Parse(t.subscribeURL)
	if err != nil {
		return err
	}

	if confirmationURL.Host[len(confirmationURL.Host)-13:] != "amazonaws.com" {
		return fmt.Errorf("subscribeURL host is NOT amazonaws.com: %v", confirmationURL.Host)
	}

	var resp *http.Response
	err = retry(5, time.Duration(3*time.Second), func() error {
		var err error
		resp, err = http.Get(confirmationURL.String())
		return err
	}, nil)

	if err != nil {
		return err
	}

	// TODO: parse the response body
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("response statusCode is not 200: %v", resp.StatusCode)
	}

	// region added successfully; send confirmation email
	newEmailBody := CompileEmail(
		`Hey {{.admin_email}}, the region <b>{{.region_name}}</b> has been successfully setup!
		<br>
		<br>
		From now on, all instances on this region will be monitored by the guardian.
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
		AccountID:        t.AdminAccount.ID, // this will also trigger send to Slack
		To:               t.AdminAccount.Email,
		Subject:          fmt.Sprintf("Region %v has been setup", t.Topic.Region),
		BodyHTML:         newEmailBody,
		BodyText:         newEmailBody,
		NotificationMeta: NotificationMeta{NotificationType: RegionSetup},

		DeliverAfter: time.Duration(time.Minute), // wait for the stack to be setup before emailing that the region has been setup
	}

	Logger.Info("ConfirmSQSSubscription", "subscribeURL", confirmationURL.String())
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
					arn:aws:sns:ap-northeast-1:012345678910:CecilTopic:c1e03965-deec-4f18-aa2b-76fbb4451a04
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

// DeleteMessage deletes an SQS message from the SQS queue.
func (t *Transmission) DeleteMessage() error {
	// remove message from queue
	return retry(5, time.Duration(3*time.Second), func() error {
		var err error
		if t.deleteMessageFromQueueParams == nil {
			return fmt.Errorf("t.deleteMessageFromQueueParams is nil")
		}
		_, err = t.s.AWS.SQS.DeleteMessage(t.deleteMessageFromQueueParams)
		return err
	}, nil)
}

//check whether someone with this aws adminAccount id is registered at cecil
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

// CreateAssumedService creates an assumed service
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

// CreateAssumedEC2Service is used to create an assumed ec2 service.
func (t *Transmission) CreateAssumedEC2Service() error {
	t.assumedEC2Service = t.s.EC2(t.assumedService, t.Topic.Region)
	return nil
}

// DescribeInstance fetches the description of an instance.
func (t *Transmission) DescribeInstance() error {
	paramsDescribeInstance := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(t.Message.Detail.InstanceID),
		},
	}
	var err error
	t.describeInstancesResponse, err = t.assumedEC2Service.DescribeInstances(paramsDescribeInstance)
	return err
}

// InstanceExists checks whether the instance specified in the event message exists on aws
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

// FetchInstance extracts and copies the instance info resulting from the DescribeInstance
// to t.Instance.
func (t *Transmission) FetchInstance() error {
	// TODO: merge the preceding check operations here
	t.Instance = *t.describeInstancesResponse.Reservations[0].Instances[0]

	if *t.Instance.InstanceId != t.Message.Detail.InstanceID {
		return fmt.Errorf("instance.InstanceId != message.Detail.InstanceID")
	}
	return nil
}

// ComputeInstanceRegion computes the region of the ec2 instance
func (t *Transmission) ComputeInstanceRegion() error {

	if t.Instance.Placement == nil {
		return fmt.Errorf("EC2Instance has nil Placement field")
	}

	if t.Instance.Placement.AvailabilityZone == nil {
		return fmt.Errorf("EC2Instance has nil Placement.AvailabilityZone field")
	}

	// TODO: is this always valid?
	az := *t.Instance.Placement.AvailabilityZone
	t.instanceRegion = az[:len(az)-1]

	return nil
}

// InstanceIsTerminated checks whether an instance is terminated (give the info Transmission already has).
func (t *Transmission) InstanceIsTerminated() bool {

	if t.Instance.State == nil {
		Logger.Warn("t.Instance.State == nil")
		return false
	}
	if t.Instance.State.Name == nil {
		Logger.Warn("t.Instance.State.Name == nil")
		return false
	}

	return *t.Instance.State.Name == ec2.InstanceStateNameTerminated
}

// InstanceIsPendingOrRunning returns true in case the instance
// is in "pending" or "running" status.
func (t *Transmission) InstanceIsPendingOrRunning() bool {
	return *t.Instance.State.Name == ec2.InstanceStateNamePending ||
		*t.Instance.State.Name == ec2.InstanceStateNameRunning
}

// LeaseIsNew checks whether a lease with the same instanceID exists
func (t *Transmission) LeaseIsNew() bool {
	var leaseCount int64
	t.s.DB.Table("leases").Where(&Lease{InstanceID: t.Message.Detail.InstanceID}).Count(&leaseCount)

	return leaseCount == 0
}

// InstanceHasGoodOwnerTag checks whether the instance has a "good"
// owner tag: i.e. a "cecilowner" tag is set, and is a valid email
// address.
func (t *Transmission) InstanceHasGoodOwnerTag() bool {
	// InstanceHasTags: check whether instance has tags
	if len(t.Instance.Tags) == 0 {
		return false
	}
	//logger.Warn("len(instance.Tags) == 0")

	// InstanceHasOwnerTag: check whether the instance has an cecilowner tag
	for _, tag := range t.Instance.Tags {
		if strings.ToLower(*tag.Key) != "cecilowner" {
			continue
		}

		// OwnerTagValueIsValid: check whether the cecilowner tag is a valid email
		ownerTag, err := t.s.DefaultMailer.Client.ValidateEmail(*tag.Value)
		if err != nil {
			Logger.Warn("ValidateEmail", "err", err)
			// TODO: send notification to admin
			return false
		}
		if !ownerTag.IsValid {
			return false
			// TODO: notify admin: "Warning: cecilowner tag email not valid" (DO NOT INCLUDE IT IN THE EMAIL, OR HTML-ESCAPE IT)
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

// SetExternalOwnerAsOwner sets the externalOwnerEmail as owner of the lease.
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

// SetAdminAsOwner sets the admin as owner of the lease.
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

// DefineLeaseDuration defines the duration of the lease.
func (t *Transmission) DefineLeaseDuration() {
	// Use global cecil lease duration setting
	t.leaseDuration = time.Duration(t.s.Config.Lease.Duration)

	// Use lease duration setting of account
	if t.AdminAccount.DefaultLeaseDuration > 0 {
		t.leaseDuration = time.Duration(t.AdminAccount.DefaultLeaseDuration)
		Logger.Info("using t.AdminAccount.DefaultLeaseDuration")
	}

	// Use lease duration setting of cloudaccount
	if t.CloudAccount.DefaultLeaseDuration > 0 {
		t.leaseDuration = time.Duration(t.CloudAccount.DefaultLeaseDuration)
		Logger.Info("using t.CloudAccount.DefaultLeaseDuration")
	}
}

// LeaseNeedsApproval returns true in case this lease needs to be approved by admin.
func (t *Transmission) LeaseNeedsApproval() bool {
	var leases []Lease

	t.s.DB.Table("leases").Where(&Lease{
		OwnerID:        t.owner.ID,
		CloudAccountID: t.CloudAccount.ID,
	}).Where("terminated_at IS NULL").Find(&leases).Count(&t.activeLeaseCount)
	//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&activeLeaseCount)

	return t.activeLeaseCount >= int64(t.s.Config.Lease.MaxPerOwner) && t.s.Config.Lease.MaxPerOwner >= 0
}

// InstanceLaunchTimeUTC is a shortcut to t.Instance.LaunchTime
func (t *Transmission) InstanceLaunchTimeUTC() time.Time {
	if t.Instance.LaunchTime == nil {
		Logger.Warn("t.Instance.LaunchTime == nil")
		return time.Now().UTC()
	}
	return t.Instance.LaunchTime.UTC()
}

// InstanceType is a shortcut to *t.Instance.InstanceType
func (t *Transmission) InstanceType() string {
	if t.Instance.InstanceType == nil {
		Logger.Warn("t.Instance.InstanceType == nil")
		return "unknown"
	}
	return *t.Instance.InstanceType
}

// InstanceID is a shortuct to *t.Instance.InstanceID
func (t *Transmission) InstanceID() string {
	if t.Instance.InstanceId == nil {
		Logger.Warn("t.Instance.InstanceId == nil")
		return "i-unknown"
	}
	return *t.Instance.InstanceId
}

// AvailabilityZone is a shortcut to *t.Instance.Placement.AvailabilityZone
func (t *Transmission) AvailabilityZone() string {
	if t.Instance.Placement.AvailabilityZone == nil {
		Logger.Warn("t.Instance.Placement.AvailabilityZone == nil")
		return "somewhere-unknown"
	}
	return *t.Instance.Placement.AvailabilityZone
}

// CreateAssumedCloudformationService is used to create an assumed cloudformation service.
func (t *Transmission) CreateAssumedCloudformationService() error {
	t.assumedCloudformationService = t.s.CloudFormation(t.assumedService, t.Topic.Region)
	return nil
}

// InstanceIsPartOfStack tells whether the instance is part of
// a cloudformation stack
func (t *Transmission) InstanceIsPartOfStack() (bool, error) {
	var instanceStackInfo = StackInfo{}

	params := &cloudformation.DescribeStackResourcesInput{
		//LogicalResourceId:  aws.String("LogicalResourceId"),
		PhysicalResourceId: t.Instance.InstanceId,
		//StackName:          aws.String("StackName"),
	}
	resp, err := t.assumedCloudformationService.DescribeStackResources(params)

	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return false, nil
		}
		return false, err
	}

	Logger.Info("DescribeStackResources", "DescribeStackResources", resp)

	if len(resp.StackResources) == 0 {
		return false, nil
	}

	var instanceStackResource *cloudformation.StackResource
	for _, stackResource := range resp.StackResources {
		Logger.Info(
			"Determining instanceStackResource",
			"stackResource.PhysicalResourceId", *stackResource.PhysicalResourceId,
			"t.InstanceID()", t.InstanceID(),
		)
		if stackResource.PhysicalResourceId != nil && *stackResource.PhysicalResourceId == t.InstanceID() {
			instanceStackResource = stackResource
		}
	}
	if instanceStackResource == nil {
		return false, errors.New("instanceStackResource is nil")
	}

	if instanceStackResource.LogicalResourceId == nil {
		return false, nil
	}
	instanceStackInfo.LogicalID = *instanceStackResource.LogicalResourceId

	if instanceStackResource.StackId == nil {
		return false, nil
	}
	instanceStackInfo.StackID = *instanceStackResource.StackId

	if instanceStackResource.StackName == nil {
		return false, nil
	}
	instanceStackInfo.StackName = *instanceStackResource.StackName

	t.StackInfo = &instanceStackInfo
	return true, nil
}

// StackHasAlreadyALease tells whether the stack to which
// the instance belongs is already registered as a lease
func (t *Transmission) StackHasAlreadyALease() (bool, error) {
	if t.StackInfo == nil {
		return false, errors.New("InstanceStack is nil")
	}
	_, err := t.s.CloudformationHasLease(int(t.AdminAccount.ID), t.StackInfo.StackID, t.StackInfo.StackName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// DescribeStack fetches the description of the stack resources
func (t *Transmission) DescribeStack() error {
	params := &cloudformation.DescribeStackResourcesInput{
		StackName: aws.String(t.StackInfo.StackName),
	}
	resp, err := t.assumedCloudformationService.DescribeStackResources(params)

	if err != nil {
		return err
	}

	Logger.Info("describeStack", "describeStack", resp)

	if resp.StackResources != nil {
		t.StackResources = resp.StackResources
	}

	return err
}
