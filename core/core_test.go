package core

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var (
	TestAWSAccountID       string = "788612350743"
	TestAWSAccountRegion   string = "us-east-1"
	TestAWSAccessKeyID     string = "WwXqFLDLbHDEIaS"               // this is a random value
	TestAWSSecretAccessKey string = "jkaeLYvjHVOmGeTYLazCgjtDqznwZ" // this is a random value
	TestReceiptHandle      string = "mockReceiptHandle"
	TestMockInstanceId     string = "i-mockinstance"
)

func TestBasicEndToEnd(t *testing.T) {

	// @@@@@@@@@@@@@@@ Create Test Service @@@@@@@@@@@@@@@

	service := createTestService("test_basic_end_to_end.db")
	defer service.Stop(false)

	// @@@@@@@@@@@@@@@ Setup mock external services @@@@@@@@@@@@@@@

	// Create mock Ec2
	mockEc2 := createMockEc2(service)

	// Queue up a response in mock ec2 to return "pending" state for instance
	mockEc2.describeInstanceResponses <- describeInstanceOutput(ec2.InstanceStateNamePending, TestMockInstanceId)

	// Get a reference to the mock SQS
	mockSQS := service.AWS.SQS.(*MockSQS)

	// Get a reference to the mock mailgun
	mockMailGun := service.Mailer.Client.(*MockMailGun)

	// @@@@@@@@@@@@@@@ Mock actions @@@@@@@@@@@@@@@

	// Launch mock ec2 instance
	launchMockEc2Instance(service, TestReceiptHandle)

	// @@@@@@@@@@@@@@@ Wait for Test actions To Finish @@@@@@@@@@@@@@@

	// Wait until the SQS message is sent back to the eventinjestor
	mockSQS.waitForReceivedMessageInput()
	mockSQS.waitForDeletedMessageInput(TestReceiptHandle)

	// Wait until the event injestor tries to describe the instance
	mockEc2.waitForDescribeInstancesInput()

	// Wait until the Sentencer tries to terminate the instance
	mockEc2.waitForTerminateInstancesInput()

	// Wait until the Sentencer tries to notifies admin that the instance was terminated
	mailGunInvocation := <-mockMailGun.SentMessages
	Logger.Info("Received mailgunInvocation", "mailgunInvocation", mailGunInvocation)
	Logger.Info("TestBasicEndToEnd finished")

}

func TestLeaseRenewal(t *testing.T) {

	// @@@@@@@@@@@@@@@ Create Test Service @@@@@@@@@@@@@@@

	service := createTestService("test_lease_renewal.db")
	defer service.Stop(false)

	// @@@@@@@@@@@@@@@ Setup mock external services @@@@@@@@@@@@@@@

	// Create mock Ec2
	mockEc2 := createMockEc2(service)

	// Queue up a response in mock ec2 to return "pending" state for instance
	mockEc2.describeInstanceResponses <- describeInstanceOutput(ec2.InstanceStateNamePending, TestMockInstanceId)

	// Get a reference to the mock SQS
	mockSQS := service.AWS.SQS.(*MockSQS)

	// Get a reference to the mock mailgun
	mockMailGun := service.Mailer.Client.(*MockMailGun)

	Logger.Info("mocks", "mockec2", mockEc2, "mocksqs", mockSQS)

	// @@@@@@@@@@@@@@@ Mock actions @@@@@@@@@@@@@@@

	// Launch mock ec2 instance
	launchMockEc2Instance(service, TestReceiptHandle)

	// @@@@@@@@@@@@@@@ Wait for Test actions To Finish @@@@@@@@@@@@@@@

	// Wait until the SQS message is sent back to the eventinjestor
	mockSQS.waitForReceivedMessageInput()
	mockSQS.waitForDeletedMessageInput(TestReceiptHandle)

	// Wait until the event injestor tries to describe the instance
	mockEc2.waitForDescribeInstancesInput()

	// Wait for email about launch
	notificationMeta := mockMailGun.waitForNotification(InstanceNeedsAttention)
	Logger.Info("InstanceNeedsAttention notification", "notificationMeta", notificationMeta)

	// Approve lease
	approveLease(service, notificationMeta.LeaseUuid, notificationMeta.InstanceId)

	// Wait for email about lease approval
	notificationMeta = mockMailGun.waitForNotification(LeaseApproved)
	Logger.Info("LeaseApproval notification", "notificationMeta", notificationMeta)

	// Wait for email about pending expiry
	notificationMeta = mockMailGun.waitForNotification(InstanceWillExpire)
	Logger.Info("InstanceWillExpire notification", "notificationMeta", notificationMeta)

	// Renew lease
	extendLease(service, notificationMeta.LeaseUuid, notificationMeta.InstanceId)

	// Wait for email about lease extended
	notificationMeta = mockMailGun.waitForNotification(LeaseExtended)
	Logger.Info("LeaseExtended notification", "notificationMeta", notificationMeta)

	// Wait for email about pending expiry
	notificationMeta = mockMailGun.waitForNotification(InstanceWillExpire)
	Logger.Info("InstanceWillExpire notification", "notificationMeta", notificationMeta)

	// Wait until the Sentencer tries to terminate the instance
	mockEc2.waitForTerminateInstancesInput()

	// Queue up a response in mock ec2 to return "terminated" state for instance
	mockEc2.describeInstanceResponses <- describeInstanceOutput(ec2.InstanceStateNameTerminated, TestMockInstanceId)

	// Terminate mock ec2 instance
	Logger.Info("terminateMockEc2Instance", "terminateMockEc2Instance", "terminateMockEc2Instance")
	terminateMockEc2Instance(service, TestReceiptHandle)

	// Wait for email about instance terminated
	notificationMeta = mockMailGun.waitForNotification(InstanceTerminated)
	Logger.Info("InstanceTerminated notification", "notificationMeta", notificationMeta)

	Logger.Info("TestLeaseRenewal finished")

}

func findLease(DB *gorm.DB, leaseUuid, instanceId string) Lease {
	var leaseToBeApproved Lease
	var leaseCount int64
	DB.Table("leases").Where(&Lease{
		InstanceID: instanceId,
		UUID:       leaseUuid,
		Terminated: false,
	}).Count(&leaseCount).First(&leaseToBeApproved)
	return leaseToBeApproved
}

func approveLease(service *Service, leaseUuid, instanceId string) {
	leaseToBeApproved := findLease(service.DB, leaseUuid, instanceId)
	service.ExtenderQueue.TaskQueue <- ExtenderTask{
		Lease:     leaseToBeApproved,
		ExtendBy:  time.Duration(service.Config.Lease.Duration),
		Approving: true,
	}
}

func extendLease(service *Service, leaseUuid, instanceId string) {
	leaseToBeExtended := findLease(service.DB, leaseUuid, instanceId)
	service.ExtenderQueue.TaskQueue <- ExtenderTask{
		Lease:     leaseToBeExtended,
		ExtendBy:  time.Duration(service.Config.Lease.Duration),
		Approving: false,
	}

}

func createMockEc2(service *Service) *MockEc2 {

	mockEc2 := NewMockEc2()
	service.EC2 = func(assumedService *session.Session, topicRegion string) ec2iface.EC2API {
		return mockEc2
	}
	return mockEc2

}

func launchMockEc2Instance(service *Service, receiptHandle string) {

	var messageBody string
	NewInstanceLaunchMessage(TestAWSAccountID, TestAWSAccountRegion, &messageBody)
	mockEc2InstanceAction(service, receiptHandle, messageBody)
}

func terminateMockEc2Instance(service *Service, receiptHandle string) {

	var messageBody string
	NewInstanceTerminatedMessage(TestAWSAccountID, TestAWSAccountRegion, &messageBody)
	mockEc2InstanceAction(service, receiptHandle, messageBody)
}

func mockEc2InstanceAction(service *Service, receiptHandle, messageBody string) {
	messages := []*sqs.Message{
		&sqs.Message{
			Body:          &messageBody,
			ReceiptHandle: &receiptHandle,
		},
	}
	mockSQSMessage := &sqs.ReceiveMessageOutput{
		Messages: messages,
	}
	mockSQS := service.AWS.SQS.(*MockSQS)
	mockSQS.Enqueue(mockSQSMessage)

}

func createTestService(dbname string) *Service {

	// this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	viper.SetDefault("AWS_REGION", TestAWSAccountRegion)
	viper.SetDefault("AWS_ACCOUNT_ID", TestAWSAccountID)
	viper.SetDefault("AWS_ACCESS_KEY_ID", TestAWSAccessKeyID)
	viper.SetDefault("AWS_SECRET_ACCESS_KEY", TestAWSSecretAccessKey)

	// Create a service
	service := NewService()
	service.LoadConfig("../config.yml")
	service.GenerateRSAKeys()
	service.SetupQueues()
	service.SetupDB(dbname)

	// Speed everything up for fast test execution
	service.Config.Lease.Duration = time.Second * 10
	service.Config.Lease.ApprovalTimeoutDuration = time.Second * 3
	service.Config.Lease.ForewarningBeforeExpiry = time.Second * 3

	// @@@@@@@@@@@@@@@ Add Fake Account / Admin  @@@@@@@@@@@@@@@

	// <EDIT-HERE>
	firstUser := Account{
		Email: "traun.leyden@gmail.com",
		CloudAccounts: []CloudAccount{
			CloudAccount{
				Provider:   "aws",
				AWSID:      TestAWSAccountID,
				ExternalID: "external_id",
			},
		},
		Verified: true,
	}
	service.DB.Create(&firstUser)

	firstOwner := Owner{
		Email:          "traun.leyden@gmail.com",
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&firstOwner)

	secondaryOwner := Owner{
		Email:          "tleyden@yahoo.com",
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&secondaryOwner)
	// </EDIT-HERE>

	// @@@@@@@@@@@@@@@ Setup mock external services @@@@@@@@@@@@@@@

	// setup mailer service
	mockMailGun := NewMockMailGun()
	service.Mailer.Client = mockMailGun

	// setup aws session -- TODO: mock this out
	AWSCreds := credentials.NewStaticCredentials(
		service.AWS.Config.AWS_ACCESS_KEY_ID,
		service.AWS.Config.AWS_SECRET_ACCESS_KEY,
		"",
	)
	AWSConfig := &aws.Config{
		Credentials: AWSCreds,
	}
	service.AWS.Session = session.New(AWSConfig)

	mockSQS := NewMockSQS()
	service.AWS.SQS = mockSQS

	// @@@@@@@@@@@@@@@ Schedule Periodic Jobs @@@@@@@@@@@@@@@

	schedulePeriodicJob(service.EventInjestorJob, time.Duration(time.Second*1))
	schedulePeriodicJob(service.AlerterJob, time.Duration(time.Second*1))
	schedulePeriodicJob(service.SentencerJob, time.Duration(time.Second*1))

	return service

}

func NewInstanceLaunchMessage(awsAccountID, awsRegion string, result *string) {
	NewSQSMessage(awsAccountID, awsRegion, result, ec2.InstanceStateNamePending)
}

func NewInstanceTerminatedMessage(awsAccountID, awsRegion string, result *string) {
	NewSQSMessage(awsAccountID, awsRegion, result, ec2.InstanceStateNameTerminated)
}

func NewSQSMessage(awsAccountID, awsRegion string, result *string, state string) {

	// create an message
	message := SQSMessage{
		Account: awsAccountID,
		Detail: SQSMessageDetail{
			State:      state,
			InstanceID: TestMockInstanceId,
		},
	}
	messageSerialized, err := json.Marshal(message)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling json: %v", err)) // TODO: return error
	}

	snsTopicName, err := viperMustGetString("SNSTopicName")
	if err != nil {
		panic(err)
	}

	// create an envelope and put the message in
	envelope := SQSEnvelope{
		TopicArn: fmt.Sprintf("arn:aws:sns:%v:%v:%v", awsRegion, awsAccountID, snsTopicName),
		Message:  string(messageSerialized),
	}

	// serialize to a string
	envelopeSerialized, err := json.Marshal(envelope)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling json: %v", err)) // TODO: return error
	}

	envSerializedString := string(envelopeSerialized)
	*result = envSerializedString
}
