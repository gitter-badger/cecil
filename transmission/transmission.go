package transmission

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	errorwrap "github.com/pkg/errors"
	"github.com/tleyden/awsutil"
	"github.com/tleyden/cecil/awstools"
	"github.com/tleyden/cecil/interfaces"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/notification"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
)

// Logger is the logger used in this package; it is initialized by the core package (see core/core-init.go)
var Logger log15.Logger

// ErrorEnvelopeIsSubscriptionConfirmation is an error that triggers automatic
// subscription between SNS and SQS.
var ErrorEnvelopeIsSubscriptionConfirmation = errors.New("ErrEnvelopeIsSubscriptionConfirmation")

// Transmission contains the SQS message and everything else
// needed to complete the operations triggered by the message
type Transmission struct {
	s            interfaces.CoreServiceInterface
	Message      awstools.SQSMessage
	subscribeURL string

	Topic struct {
		Region string
		AWSID  string
	}
	deleteMessageFromQueueParams *sqs.DeleteMessageInput

	Cloudaccount models.Cloudaccount
	AdminAccount models.Account

	assumedService               *session.Session
	assumedEC2Service            ec2iface.EC2API
	assumedCloudformationService cloudformationiface.CloudFormationAPI
	assumedAutoScalingService    autoscalingiface.AutoScalingAPI

	DescribeInstancesResponse *ec2.DescribeInstancesOutput

	Instance           ec2.Instance
	InstanceRegion     string
	externalOwnerEmail string
	Owner              models.Owner
	LeaseDuration      time.Duration
	ActiveLeaseCount   int64

	GroupType models.GroupType
}

// GenerateSQSTransmission parses a raw SQS message into a Transmission
func GenerateSQSTransmission(s interfaces.CoreServiceInterface, rawMessage *sqs.Message, queueURL string) (*Transmission, error) {

	// Record the event
	if err := s.EventRecorder().StoreSQSMessage(rawMessage); err != nil {
		Logger.Warn("Error storing SQS message", "err", err, "msg", rawMessage)
	}

	newTransmission := Transmission{}
	newTransmission.s = s

	newTransmission.deleteMessageFromQueueParams = &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),                  // Required
		ReceiptHandle: aws.String(*rawMessage.ReceiptHandle), // Required
	}

	// parse the envelope
	var envelope awstools.SQSEnvelope
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

	if topicName != s.AWSRes().Config.SNSTopicName {
		// NOTE: here we return &newTransmission (and not an empty &Transmission{}) because
		// &newTransmission contains values that will be used
		return &newTransmission, fmt.Errorf("the SNS topic name is not %v: %v", s.AWSRes().Config.SNSTopicName, topicName)
	}

	if newTransmission.Topic.Region == "" {
		return &newTransmission, fmt.Errorf("newTransmission.Topic.Region is empty")
	}
	if newTransmission.Topic.AWSID == "" {
		return &newTransmission, fmt.Errorf("newTransmission.Topic.AWSID is empty")
	}

	//check whether someone with this aws adminAccount id is registered at cecil
	err = newTransmission.FetchCloudaccount()
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
	err = tools.Retry(5, time.Duration(3*time.Second), func() error {
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
	newEmailBody, err := tools.CompileEmailTemplate(
		"region-successfully-setup.html",
		map[string]interface{}{
			"admin_email": t.AdminAccount.Email,
			"region_name": t.Topic.Region,
		},
	)
	if err != nil {
		return err
	}

	Logger.Info("ConfirmSQSSubscription", "subscribeURL", confirmationURL.String())

	return t.s.Queues().NotifierQueue().PushTask(tasks.NotifierTask{
		AccountID:        t.AdminAccount.ID, // this will also trigger send to Slack
		To:               t.AdminAccount.Email,
		Subject:          fmt.Sprintf("Region %v has been setup", t.Topic.Region),
		BodyHTML:         newEmailBody,
		BodyText:         newEmailBody,
		NotificationMeta: notification.NotificationMeta{NotificationType: notification.RegionSetup},

		DeliverAfter: time.Duration(time.Minute), // wait for the stack to be setup before emailing that the region has been setup
	})

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

// TopicAndInstanceHaveSameOwner tells whether the originating SNS topic and the instance have different owners (different AWS accounts)
func (t *Transmission) TopicAndInstanceHaveSameOwner() bool {
	instanceOriginatorID := t.Message.Account

	return t.Topic.AWSID == instanceOriginatorID
}

// MessageIsRelevant considers only pending and terminated status messages; ignores the rest
func (t *Transmission) MessageIsRelevant() bool {
	return t.Message.Detail.State == ec2.InstanceStateNamePending ||
		t.Message.Detail.State == ec2.InstanceStateNameTerminated
}

// DeleteMessage deletes an SQS message from the SQS queue.
func (t *Transmission) DeleteMessage() error {
	// remove message from queue
	return tools.Retry(5, time.Duration(3*time.Second), func() error {
		var err error
		if t.deleteMessageFromQueueParams == nil {
			return fmt.Errorf("t.deleteMessageFromQueueParams is nil")
		}
		_, err = t.s.AWSRes().SQS.DeleteMessage(t.deleteMessageFromQueueParams)
		return err
	}, nil)
}

// FetchCloudaccount checks whether someone with this aws adminAccount id is registered at cecil,
// and fetches the cloudaccount associated with that AWS account.
func (t *Transmission) FetchCloudaccount() error {
	return t.s.GormDB().
		Where(&models.Cloudaccount{AWSID: t.Topic.AWSID}).
		First(&t.Cloudaccount).
		Error
}

// FetchAdminAccount checks whether the cloud account has an admin account,
// and fetches it.
func (t *Transmission) FetchAdminAccount() error {
	return t.s.GormDB().
		Model(&t.Cloudaccount).
		Related(&t.AdminAccount).
		Error
	//s.DB.Table("accounts").Where([]uint{cloudaccount.AccountID}).First(&cloudaccount).Count(&cloudaccountAdminCount)
}

// CreateAssumedService creates an assumed service
func (t *Transmission) CreateAssumedService() error {
	assumedConfig := &aws.Config{
		Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
			Client: sts.New(t.s.AWSRes().Session, &aws.Config{Region: aws.String(t.Topic.Region)}),
			RoleARN: fmt.Sprintf(
				"arn:aws:iam::%v:role/%v",
				t.Topic.AWSID,
				t.s.AWSRes().Config.ForeignIAMRoleName,
			),
			RoleSessionName: uuid.NewV4().String(),
			ExternalID:      aws.String(t.Cloudaccount.ExternalID),
			ExpiryWindow:    60 * time.Second,
		}),
	}

	t.assumedService = session.New(assumedConfig)

	return nil
}

// CreateAssumedEC2Service is used to create an assumed ec2 service.
func (t *Transmission) CreateAssumedEC2Service() error {
	t.assumedEC2Service = t.s.AWSRes().EC2(t.assumedService, t.Topic.Region)
	return nil
}

// CreateAssumedCloudformationService is used to create an assumed cloudformation service.
func (t *Transmission) CreateAssumedCloudformationService() error {
	t.assumedCloudformationService = t.s.AWSRes().CloudFormation(t.assumedService, t.Topic.Region)
	return nil
}

// CreateAssumedAutoscalingService is used to create an assumed autoscaling service.
func (t *Transmission) CreateAssumedAutoscalingService() error {
	t.assumedAutoScalingService = t.s.AWSRes().AutoScaling(t.assumedService, t.Topic.Region)
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
	t.DescribeInstancesResponse, err = t.assumedEC2Service.DescribeInstances(paramsDescribeInstance)
	if err != nil {
		return errorwrap.Wrap(err, "error while DescribeInstances")
	}

	if err = t.copyInstanceInfo(); err != nil {
		return errorwrap.Wrap(err, "error while copyInstanceInfo")
	}

	if err = t.computeInstanceRegion(); err != nil {
		return errorwrap.Wrap(err, "error while computing instance region")
	}

	return nil
}

// InstanceExists checks whether the instance specified in the event message exists on aws
func (t *Transmission) InstanceExists() bool {
	if len(t.DescribeInstancesResponse.Reservations) == 0 {
		// logger.Warn("len(DescribeInstancesResponse.Reservations) == 0")
		return false
	}
	if len(t.DescribeInstancesResponse.Reservations[0].Instances) == 0 {
		// logger.Warn("len(DescribeInstancesResponse.Reservations[0].Instances) == 0")
		return false
	}
	return true
}

// copyInstanceInfo extracts and copies the instance info resulting from the DescribeInstance
// to t.Instance.
func (t *Transmission) copyInstanceInfo() error {
	// TODO: merge the preceding check operations here
	t.Instance = *t.DescribeInstancesResponse.Reservations[0].Instances[0]

	if *t.Instance.InstanceId != t.Message.Detail.InstanceID {
		return fmt.Errorf("instance.InstanceId != message.Detail.InstanceID")
	}
	return nil
}

// computeInstanceRegion computes the region of the ec2 instance
func (t *Transmission) computeInstanceRegion() error {

	if t.Instance.Placement == nil {
		return fmt.Errorf("EC2Instance has nil Placement field")
	}

	if t.Instance.Placement.AvailabilityZone == nil {
		return fmt.Errorf("EC2Instance has nil Placement.AvailabilityZone field")
	}

	// TODO: is this always valid?
	az := *t.Instance.Placement.AvailabilityZone
	t.InstanceRegion = az[:len(az)-1]

	return nil
}

// InstanceIsTerminated checks whether an instance is terminated (given the info Transmission already has).
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

// InstanceStateName returns t.Instance.State.Name
func (t *Transmission) InstanceStateName() string {

	if t.Instance.State == nil {
		Logger.Warn("t.Instance.State == nil")
		return ""
	}
	if t.Instance.State.Name == nil {
		Logger.Warn("t.Instance.State.Name == nil")
		return ""
	}

	return *t.Instance.State.Name
}

// InstanceIsPendingOrRunning returns true in case the instance
// is in "pending" or "running" status.
func (t *Transmission) InstanceIsPendingOrRunning() bool {
	return *t.Instance.State.Name == ec2.InstanceStateNamePending ||
		*t.Instance.State.Name == ec2.InstanceStateNameRunning
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
		ownerTag, err := t.s.DefaultMailer().Client.ValidateEmail(*tag.Value)
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

// InstanceHasWhitelistedKeyName returns true in case there is a whitelisted owner that
// is linked to the keyname that created the instance.
func (t *Transmission) InstanceHasWhitelistedKeyName() bool {
	// TODO: make sure this works for instances created as part of a stack, and not only individual
	// instances created directly by a person (with a key).

	if t.Instance.KeyName == nil {
		return false
	}

	keyName := *t.Instance.KeyName
	var externalOwner models.Owner
	err := t.s.GormDB().Table("owners").Where(&models.Owner{KeyName: keyName, CloudaccountID: t.Cloudaccount.ID}).Find(&externalOwner).Error
	if err != nil {
		return false
	}
	t.externalOwnerEmail = externalOwner.Email

	return true
}

// InstanceHasTagOrKeyName returns true in case the instance has an owner tag that is a valid email address,
// or a KeyName that is registered in Cecil's Owners table.
func (t *Transmission) InstanceHasTagOrKeyName() bool {
	hasTag := t.InstanceHasGoodOwnerTag()
	hasKeyName := t.InstanceHasWhitelistedKeyName()
	return hasTag || hasKeyName
}

// ExternalOwnerIsWhitelisted checks whether the owner email in the tag is a whitelisted owner email
func (t *Transmission) ExternalOwnerIsWhitelisted() bool {
	// TODO: select Owner by email, cloudaccountid, and region?
	if t.externalOwnerEmail == "" {
		return false
	}
	// TODO: use Retry
	err := t.s.GormDB().
		Table("owners").
		Where(&models.Owner{Email: t.externalOwnerEmail, CloudaccountID: t.Cloudaccount.ID}).
		Error
	if err != nil {
		return false
	}
	return true
}

// SetExternalOwnerAsOwner sets the externalOwnerEmail as owner of the lease.
func (t *Transmission) SetExternalOwnerAsOwner() error {
	var owner models.Owner
	err := t.s.GormDB().
		Table("owners").
		Where(&models.Owner{Email: t.externalOwnerEmail, CloudaccountID: t.Cloudaccount.ID}).
		First(&owner).
		Error
	if err != nil {
		return err
	}
	t.Owner = owner
	return nil
}

// SetAdminAsOwner sets the admin as owner of the lease.
func (t *Transmission) SetAdminAsOwner() error {
	var owner models.Owner
	err := t.s.GormDB().
		Table("owners").
		Where(&models.Owner{Email: t.AdminAccount.Email, CloudaccountID: t.Cloudaccount.ID}).
		First(&owner).
		Error
	if err != nil {
		return err
	}
	t.Owner = owner
	return nil
}

// LeaseExpiresAt defines the duration of the lease.
func (t *Transmission) LeaseExpiresAt() time.Time {
	var expiresAt time.Time
	// Use global cecil lease duration setting
	durationFromDefault := time.Duration(t.s.Config().Lease.Duration)
	expiresAt = time.Now().UTC().Add(durationFromDefault)

	// Use lease duration setting of account
	if t.AdminAccount.DefaultLeaseDuration > 0 {
		durationFromAdminAccount := time.Duration(t.AdminAccount.DefaultLeaseDuration)
		expiresAt = time.Now().UTC().Add(durationFromAdminAccount)
	}

	// Use lease duration setting of cloudaccount
	if t.Cloudaccount.DefaultLeaseDuration > 0 {
		durationFromCloudaccount := time.Duration(t.Cloudaccount.DefaultLeaseDuration)
		expiresAt = time.Now().UTC().Add(durationFromCloudaccount)
	}

	if expiresIn := t.InstanceHasTagForExpiresIn(); expiresIn != nil {
		durationFromExpiresInTag := *expiresIn
		// Using launch_time + expires_in
		expiresAt = t.InstanceLaunchTimeUTC().Add(durationFromExpiresInTag)
	}

	if expiresOn := t.InstanceHasTagForExpiresOn(); expiresOn != nil && expiresOn.After(time.Now().UTC()) {
		expiresAt = *expiresOn
	}

	return expiresAt
}

// LeaseNeedsApproval returns true in case this lease needs to be approved by admin.
func (t *Transmission) LeaseNeedsApproval() bool {
	var leases []models.Lease

	t.s.GormDB().
		Table("leases").
		Where(&models.Lease{
			OwnerID:        t.Owner.ID,
			CloudaccountID: t.Cloudaccount.ID,
		}).
		Where("terminated_at IS NULL").
		Find(&leases).
		Count(&t.ActiveLeaseCount)
	//s.DB.Table("accounts").Where([]uint{cloudaccount.AccountID}).First(&cloudaccount).Count(&ActiveLeaseCount)

	return t.ActiveLeaseCount >= int64(t.s.Config().Lease.MaxPerOwner) && t.s.Config().Lease.MaxPerOwner >= 0
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

// DefineGroupUID defines the groupUID of the instance
func (t *Transmission) DefineGroupUID() (*string, *string, error) {

	// check if group is CecilGroupTag
	cecilGroupTag := t.GroupIsCecilGroupTag()
	if cecilGroupTag != nil {
		return cecilGroupTag, nil, nil
	}

	// check whether group is ASG
	autoScalingGroupARN, groupName, err := t.GroupIsASG()
	if err != nil {
		Logger.Error("Error while checking if part of an AutoScalingGroup", "err", err, "transmission", t)
	}
	if autoScalingGroupARN != nil {
		return autoScalingGroupARN, groupName, nil
	}

	// check whether group is CF
	cloudFormationARN, cloudFormationName, err := t.GroupIsCF()
	if err != nil {
		Logger.Error("Error while checking if part of a CloudFormation", "err", err, "transmission", t)
	}
	if cloudFormationARN != nil {
		return cloudFormationARN, cloudFormationName, nil
	}

	// TODO: figure out when group is GroupTime

	// use instanceID as groupUID
	groupUID := t.InstanceID()
	t.GroupType = models.GroupSingle

	return &groupUID, nil, nil
}

func (t *Transmission) GroupIsCecilGroupTag() *string {
	// InstanceHasTags: check whether instance has tags
	if len(t.Instance.Tags) == 0 {
		return nil
	}

	// GroupIsCecilGroupTag: check whether the instance has a CecilInstanceGroup tag
	for _, tag := range t.Instance.Tags {
		if tag == nil {
			continue
		}
		if strings.ToLower(*tag.Key) != strings.ToLower("CecilInstanceGroup") {
			continue
		}

		cecilGroupTag := *tag.Value

		t.GroupType = models.GroupCecilGroupTag

		return &cecilGroupTag
	}
	return nil
}

func (t *Transmission) GroupIsASG() (*string, *string, error) {
	asUtil, err := awsutil.NewAutoScalingUtil(t.assumedAutoScalingService, t.assumedEC2Service)
	if err != nil {
		return nil, nil, err
	}

	in, asgName, err := asUtil.InAutoScaling(*t.Instance.InstanceId)
	if err != nil {
		return nil, nil, err
	}

	if !in {
		return nil, nil, nil
	}

	asgARN, err := asUtil.ASG_ARN_fromName(asgName)
	if err != nil {
		return nil, nil, err
	}

	// use the ARN as groupID
	groupName := asgARN

	t.GroupType = models.GroupASG

	return &groupName, &asgName, nil
}

func (t *Transmission) GroupIsCF() (*string, *string, error) {
	cfnUtil, err := awsutil.NewCloudformationUtil(t.assumedCloudformationService, t.assumedEC2Service)
	if err != nil {
		return nil, nil, err
	}

	in, stackID, stackName, err := cfnUtil.InCloudformation(*t.Instance.InstanceId)
	if err != nil {
		return nil, nil, err
	}

	if !in {
		return nil, nil, nil
	}

	// stackID is the ARN of the CF stack
	groupName := stackID
	t.GroupType = models.GroupCF

	return &groupName, &stackName, nil
}

// GroupHasAlreadyALease tells whether the group to which
// the instance belongs is already registered as a lease
func (t *Transmission) GroupHasAlreadyALease(groupUID string) (*models.Lease, error) {

	lease, err := t.s.LeaseByGroupUID(t.AdminAccount.ID, &t.Cloudaccount.ID, groupUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return lease, nil
}

func (t *Transmission) LeaseByInstanceID() (*models.Lease, error) {
	Logger.Info("Calling t.s.LeaseByInstanceID",
		"t.AdminAccount.ID", t.AdminAccount.ID,
		"t.Cloudaccount.ID", t.Cloudaccount.ID,
		"t.InstanceID()", t.InstanceID(),
	)
	return t.s.LeaseByInstanceID(t.AdminAccount.ID, &t.Cloudaccount.ID, t.InstanceID())
}

// InstanceIsNew tells whether the instance is
// alredy registered with a group
func (t *Transmission) InstanceIsNew() (*models.Instance, error) {
	instance, err := t.s.GetInstanceByAWSInstanceID(t.AdminAccount.ID, &t.Cloudaccount.ID, t.InstanceID())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return instance, nil
}

// ActiveInstancesForGroup fetches from DB all non-terminated instances for a group
func (t *Transmission) ActiveInstancesForGroup(groupUID string) ([]*models.Instance, error) {
	return t.s.ActiveInstancesForGroup(t.AdminAccount.ID, &t.Cloudaccount.ID, groupUID)
}

// AWSResourceID returns the AWS StackID if the resource is a cloudformation stack, or the AWS Instance ID if the resource is an EC2 instance
func (t *Transmission) AWSResourceID() string {
	// TODO: use groupUID
	return ""
}

// InstanceHasTagForExpiresIn returns a duration if the CecilLeaseExpiresIn tag is set and is a valid duration
func (t *Transmission) InstanceHasTagForExpiresIn() *time.Duration {
	// InstanceHasTags: check whether instance has tags
	if len(t.Instance.Tags) == 0 {
		return nil
	}

	// InstanceHasOwnerTag: check whether the instance has an cecilowner tag
	for _, tag := range t.Instance.Tags {
		if strings.ToLower(*tag.Key) != strings.ToLower("CecilLeaseExpiresIN") {
			continue
		}

		expiresInStr := *tag.Value
		expiresIn, err := time.ParseDuration(expiresInStr)
		if err != nil {
			// TODO: do something with errors
			return nil
		}
		return &expiresIn
	}
	return nil
}

// InstanceHasTagForExpiresOn returns a duration if the CecilLeaseExpiresOn tag is set and is a valid time/date
func (t *Transmission) InstanceHasTagForExpiresOn() *time.Time {
	// InstanceHasTags: check whether instance has tags
	if len(t.Instance.Tags) == 0 {
		return nil
	}

	// InstanceHasOwnerTag: check whether the instance has an cecilowner tag
	for _, tag := range t.Instance.Tags {
		if strings.ToLower(*tag.Key) != strings.ToLower("CecilLeaseExpiresON") {
			continue
		}

		expiresOnStr := *tag.Value

		timeParser := when.New(nil)
		timeParser.Add(en.All...)
		timeParser.Add(common.All...)

		result, err := timeParser.Parse(expiresOnStr, time.Now().UTC())
		if err != nil {
			// TODO: do something with errors
			return nil
		}
		if result == nil {
			// TODO: do something with errors
			return nil
		}

		expiresOn := result.Time.UTC()
		return &expiresOn
	}
	return nil
}
