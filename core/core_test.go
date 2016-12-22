package core_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/goadesign/goa"
	goaclient "github.com/goadesign/goa/client"
	goalog15 "github.com/goadesign/goa/logging/log15"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tleyden/cecil/controllers"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
	apiserverclient "github.com/tleyden/cecil/goa/client"
)

var (
	TestAWSAccountID       string = "788612350743"
	TestAWSAccountRegion   string = "us-east-1"
	TestAWSAccessKeyID     string = "WwXqFLDLbHDEIaS"               // this is a random value
	TestAWSSecretAccessKey string = "jkaeLYvjHVOmGeTYLazCgjtDqznwZ" // this is a random value
	TestReceiptHandle      string = "mockReceiptHandle"
	TestMockInstanceId     string = "i-mockinstance"
	TestMockSQSMsgCount    int64  = 0
)

func createTempDBFile(filename string) *os.File {
	tmpfile, err := ioutil.TempFile("", filename)
	if err != nil {
		log.Fatal(err)
	}
	return tmpfile
}

func TestBasicEndToEnd(t *testing.T) {

	// @@@@@@@@@@@@@@@ Create Test Service @@@@@@@@@@@@@@@

	tempDBFile := createTempDBFile("test_basic_end_to_end.db")
	defer os.Remove(tempDBFile.Name())
	service := createTestService(tempDBFile.Name())
	defer service.Stop(false)

	// @@@@@@@@@@@@@@@ Setup mock external services @@@@@@@@@@@@@@@

	// Create mock Ec2
	mockEc2 := createMockEc2(service)

	// Queue up a response in mock ec2 to return "pending" state for instance
	mockEc2.DescribeInstanceResponses <- core.DescribeInstanceOutput(ec2.InstanceStateNamePending, TestMockInstanceId)

	// Get a reference to the mock SQS
	mockSQS := service.AWS.SQS.(*core.MockSQS)

	// Get a reference to the mock mailgun
	mockMailGun := service.DefaultMailer.Client.(*core.MockMailGun)

	// @@@@@@@@@@@@@@@ Mock actions @@@@@@@@@@@@@@@

	// Launch mock ec2 instance
	launchMockEc2Instance(service, TestReceiptHandle)

	// @@@@@@@@@@@@@@@ Wait for Test actions To Finish @@@@@@@@@@@@@@@

	// Wait until the SQS message is sent back to the eventinjestor
	mockSQS.WaitForReceivedMessageInput()
	mockSQS.WaitForDeletedMessageInput(TestReceiptHandle)

	// Wait until the event injestor tries to describe the instance
	mockEc2.WaitForDescribeInstancesInput()

	// Wait until the Sentencer tries to terminate the instance
	mockEc2.WaitForTerminateInstancesInput()

	// Wait until the Sentencer tries to notifies admin that the instance was terminated
	mailGunInvocation := <-mockMailGun.SentMessages
	core.Logger.Info("Received mailgunInvocation", "mailgunInvocation", mailGunInvocation)

	// Make sure the SQS event recorder works
	storedSqsMessages, err := service.EventRecord.GetStoredSQSMessages()
	if err != nil {
		panic(fmt.Sprintf("Error getting stored sqs messages: %v", err))
	}
	if len(storedSqsMessages) == 0 {
		panic(fmt.Sprintf("Expected to record sqs messages"))
	}
	for _, sqsMessage := range storedSqsMessages {
		core.Logger.Info("Recorded sqs event", "sqsMessage", sqsMessage)
	}

	core.Logger.Info("TestBasicEndToEnd finished")

}

func TestLeaseRenewal(t *testing.T) {

	// @@@@@@@@@@@@@@@ Create Test Service @@@@@@@@@@@@@@@

	tempDBFile := createTempDBFile("test_lease_renewal.db")
	defer os.Remove(tempDBFile.Name())

	service := createTestService(tempDBFile.Name())
	defer service.Stop(false)

	// @@@@@@@@@@@@@@@ Setup mock external services @@@@@@@@@@@@@@@

	// Create mock Ec2
	mockEc2 := createMockEc2(service)

	// Queue up a response in mock ec2 to return "pending" state for instance
	mockEc2.DescribeInstanceResponses <- core.DescribeInstanceOutput(ec2.InstanceStateNamePending, TestMockInstanceId)

	// Get a reference to the mock SQS
	mockSQS := service.AWS.SQS.(*core.MockSQS)

	// Get a reference to the mock mailgun
	mockMailGun := service.DefaultMailer.Client.(*core.MockMailGun)

	core.Logger.Info("mocks", "mockec2", mockEc2, "mocksqs", mockSQS)

	// @@@@@@@@@@@@@@@ Mock actions @@@@@@@@@@@@@@@

	// Launch mock ec2 instance
	launchMockEc2Instance(service, TestReceiptHandle)

	// @@@@@@@@@@@@@@@ Wait for Test actions To Finish @@@@@@@@@@@@@@@

	// Wait until the SQS message is sent back to the eventinjestor
	mockSQS.WaitForReceivedMessageInput()
	mockSQS.WaitForDeletedMessageInput(TestReceiptHandle)

	// Wait until the event injestor tries to describe the instance
	mockEc2.WaitForDescribeInstancesInput()

	// Wait for email about launch
	notificationMeta := mockMailGun.WaitForNotification(core.InstanceNeedsAttention)
	core.Logger.Info("InstanceNeedsAttention notification", "notificationMeta", notificationMeta)

	// Approve lease
	approveLease(service, notificationMeta.LeaseUuid, notificationMeta.InstanceId)

	// Wait for email about lease approval
	notificationMeta = mockMailGun.WaitForNotification(core.LeaseApproved)
	core.Logger.Info("LeaseApproval notification", "notificationMeta", notificationMeta)

	// Wait for email about pending expiry
	notificationMeta = mockMailGun.WaitForNotification(core.InstanceWillExpire)
	core.Logger.Info("InstanceWillExpire notification", "notificationMeta", notificationMeta)

	// Renew lease
	extendLease(service, notificationMeta.LeaseUuid, notificationMeta.InstanceId)

	// Wait for email about lease extended
	notificationMeta = mockMailGun.WaitForNotification(core.LeaseExtended)
	core.Logger.Info("LeaseExtended notification", "notificationMeta", notificationMeta)

	// Wait for email about pending expiry
	notificationMeta = mockMailGun.WaitForNotification(core.InstanceWillExpire)
	core.Logger.Info("InstanceWillExpire notification", "notificationMeta", notificationMeta)

	// Wait until the Sentencer tries to terminate the instance
	mockEc2.WaitForTerminateInstancesInput()

	// Queue up a response in mock ec2 to return "terminated" state for instance
	mockEc2.DescribeInstanceResponses <- core.DescribeInstanceOutput(ec2.InstanceStateNameTerminated, TestMockInstanceId)

	// Terminate mock ec2 instance
	core.Logger.Info("terminateMockEc2Instance", "terminateMockEc2Instance", "terminateMockEc2Instance")
	terminateMockEc2Instance(service, TestReceiptHandle)

	// Wait for email about instance terminated
	notificationMeta = mockMailGun.WaitForNotification(core.InstanceTerminated)
	core.Logger.Info("InstanceTerminated notification", "notificationMeta", notificationMeta)

	core.Logger.Info("TestLeaseRenewal finished")

}

func TestAccountCreation(t *testing.T) {

	// ---------------------------- Setup ----------------------------------

	// Create goa and core service
	service := goa.New("Cecil REST API")

	tempDBFile := createTempDBFile("test_account_creation.db")
	defer os.Remove(tempDBFile.Name())
	coreService := createTestService(tempDBFile.Name())

	// Mount "account" controller
	accountController := controllers.NewAccountController(service, coreService)
	app.MountAccountController(service, accountController)

	// Http and Api Client
	httpClient := http.DefaultClient
	APIClient := apiserverclient.New(goaclient.HTTPClientDoer(httpClient))

	// Goa context
	logger := goalog15.New(core.Logger)
	ctx := goa.WithLogger(context.Background(), logger)

	// ---------------------------- Create Account --------------------------

	// Create API request to create an account
	createAccountPayload := apiserverclient.CreateAccountPayload{
		Email:   "testing@test.com",
		Name:    "Test",
		Surname: "Ing",
	}
	path := apiserverclient.CreateAccountPath()

	req, err := APIClient.NewCreateAccountRequest(ctx, path, &createAccountPayload)
	if err != nil {
		panic(fmt.Sprintf("error creating new account request: %v", err))
	}
	resp := httptest.NewRecorder()

	createAccountHandler := service.Mux.Lookup(http.MethodPost, path)
	if createAccountHandler == nil {
		t.Fatalf("createAccountHandler is nil")
	}

	// Invoke API method and get response
	createAccountHandler(resp, req, req.URL.Query())

	// Process response
	if resp.Code != http.StatusOK {
		t.Fatalf("response status code is not 200", "code", resp.Code)
	}

	// Parse CreateAccount JSON response and extract account ID
	decoder := json.NewDecoder(resp.Body)
	responseJson := gin.H{}
	if err := decoder.Decode(&responseJson); err != nil {
		t.Fatalf("Could not decode response json when creating account: %v", err)
	}
	accountId := fmt.Sprintf("%v", responseJson["account_id"])

	// ---------------------------- Verify Account --------------------------------

	// Wait for verification email
	mockMailGun := coreService.DefaultMailer.Client.(*core.MockMailGun)
	notificationMeta := mockMailGun.WaitForNotification(core.VerifyingAccount)
	core.Logger.Info("Got Verification email", "notificationMeta", notificationMeta)

	// Verify account using verification token from email
	accountIdInt, err := strconv.Atoi(accountId)
	if err != nil {
		panic("Error converting string -> int")
	}

	// this path just has the placeholder variable rather than the actual account ID
	// since otherwise the service.Mux.Lookup() call will fail.  The account id is specified
	// explicitly in the call to the handler as the
	path = "/accounts/:account_id/api_token"

	// Create the API request to verify account using verification token
	// and account ID from previous step
	verifyAccountPayload := apiserverclient.VerifyAccountPayload{
		VerificationToken: notificationMeta.VerificationToken,
	}
	ctx = context.WithValue(ctx, "account_id", accountIdInt)
	req, err = APIClient.NewVerifyAccountRequest(ctx, path, &verifyAccountPayload)
	if err != nil {
		panic(fmt.Sprintf("error creating verify account request: %v", err))
	}

	// Record the response so we can later read it
	resp = httptest.NewRecorder()

	// Lookup the verify account API handler
	verifyAccountHandler := service.Mux.Lookup(http.MethodPost, path)
	if verifyAccountHandler == nil {
		t.Fatalf("verifyAccountHandler is nil")
	}

	// Create parameters that are normally extracted from URL string
	urlValues := url.Values{}
	urlValues["account_id"] = []string{accountId}

	// Invoke verify account API and get response
	verifyAccountHandler(resp, req, urlValues)

	// Make sure the response code to the account verification endpoint is 2XX
	if resp.Code != http.StatusOK {
		t.Fatalf("Unexpected response status code: %v", resp.Code)
	}

}

func findLease(DB *gorm.DB, leaseUuid, instanceId string) core.Lease {
	var leaseToBeApproved core.Lease
	var leaseCount int64
	DB.Table("leases").Where(&core.Lease{
		InstanceID: instanceId,
		UUID:       leaseUuid,
	}).Where("terminated_at IS NULL").Count(&leaseCount).First(&leaseToBeApproved)
	return leaseToBeApproved
}

func approveLease(service *core.Service, leaseUuid, instanceId string) {
	leaseToBeApproved := findLease(service.DB, leaseUuid, instanceId)
	service.ExtenderQueue.TaskQueue <- core.ExtenderTask{
		Lease:     leaseToBeApproved,
		Approving: true,
	}
}

func extendLease(service *core.Service, leaseUuid, instanceId string) {
	leaseToBeExtended := findLease(service.DB, leaseUuid, instanceId)
	service.ExtenderQueue.TaskQueue <- core.ExtenderTask{
		Lease:     leaseToBeExtended,
		Approving: false,
	}
}

func createMockEc2(service *core.Service) *core.MockEc2 {

	mockEc2 := core.NewMockEc2()
	service.EC2 = func(assumedService *session.Session, topicRegion string) ec2iface.EC2API {
		return mockEc2
	}
	return mockEc2

}

func launchMockEc2Instance(service *core.Service, receiptHandle string) {

	var messageBody string
	messageId := NewInstanceLaunchMessage(TestAWSAccountID, TestAWSAccountRegion, &messageBody)
	mockEc2InstanceAction(service, receiptHandle, messageBody, messageId)
}

func terminateMockEc2Instance(service *core.Service, receiptHandle string) {

	var messageBody string
	messageId := NewInstanceTerminatedMessage(TestAWSAccountID, TestAWSAccountRegion, &messageBody)
	mockEc2InstanceAction(service, receiptHandle, messageBody, messageId)
}

func mockEc2InstanceAction(service *core.Service, receiptHandle, messageBody, messageId string) {
	messages := []*sqs.Message{
		&sqs.Message{
			MessageId:     &messageId,
			Body:          &messageBody,
			ReceiptHandle: &receiptHandle,
		},
	}
	mockSQSMessage := &sqs.ReceiveMessageOutput{
		Messages: messages,
	}
	mockSQS := service.AWS.SQS.(*core.MockSQS)
	mockSQS.Enqueue(mockSQSMessage)

}

func createTestService(dbname string) *core.Service {

	// this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	viper.SetDefault("AWS_REGION", TestAWSAccountRegion)
	viper.SetDefault("AWS_ACCOUNT_ID", TestAWSAccountID)
	viper.SetDefault("AWS_ACCESS_KEY_ID", TestAWSAccessKeyID)
	viper.SetDefault("AWS_SECRET_ACCESS_KEY", TestAWSSecretAccessKey)

	// Create a service
	service := core.NewService()
	service.LoadConfig("../config.yml")
	service.GenerateRSAKeys()
	service.SetupQueues()
	service.SetupDB(dbname)
	service.SetupEventRecording(false, "")

	// Speed everything up for fast test execution
	service.Config.Lease.Duration = time.Millisecond * 1000
	service.Config.Lease.ApprovalTimeoutDuration = time.Millisecond * 300
	service.Config.Lease.ForewarningBeforeExpiry = time.Millisecond * 300

	// @@@@@@@@@@@@@@@ Add Fake Account / Admin  @@@@@@@@@@@@@@@

	// <EDIT-HERE>
	firstUser := core.Account{
		Email: "firstUser@gmail.com",
		CloudAccounts: []core.CloudAccount{
			core.CloudAccount{
				Provider:   "aws",
				AWSID:      TestAWSAccountID,
				ExternalID: "external_id",
			},
		},
		Verified: true,
	}
	service.DB.Create(&firstUser)

	firstOwner := core.Owner{
		Email:          "firstUser@gmail.com",
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&firstOwner)

	secondaryOwner := core.Owner{
		Email:          "secondaryOwner@yahoo.com",
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&secondaryOwner)
	// </EDIT-HERE>

	// @@@@@@@@@@@@@@@ Setup mock external services @@@@@@@@@@@@@@@

	// setup mailer service
	mockMailGun := core.NewMockMailGun()
	service.DefaultMailer.Client = mockMailGun

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

	mockSQS := core.NewMockSQS()
	service.AWS.SQS = mockSQS

	// @@@@@@@@@@@@@@@ Schedule Periodic Jobs @@@@@@@@@@@@@@@

	core.SchedulePeriodicJob(service.EventInjestorJob, time.Duration(time.Second*1))
	core.SchedulePeriodicJob(service.AlerterJob, time.Duration(time.Second*1))
	core.SchedulePeriodicJob(service.SentencerJob, time.Duration(time.Second*1))

	return service

}

func NewInstanceLaunchMessage(awsAccountID, awsRegion string, result *string) (messageId string) {
	return NewSQSMessage(awsAccountID, awsRegion, result, ec2.InstanceStateNamePending)
}

func NewInstanceTerminatedMessage(awsAccountID, awsRegion string, result *string) (messageId string) {
	return NewSQSMessage(awsAccountID, awsRegion, result, ec2.InstanceStateNameTerminated)
}

func NewSQSMessage(awsAccountID, awsRegion string, result *string, state string) (messageId string) {

	msgCounter := atomic.AddInt64(&TestMockSQSMsgCount, 1)
	msgId := fmt.Sprintf("mock_sqs_message_%d", msgCounter)

	// create an message
	message := core.SQSMessage{
		ID:      msgId,
		Account: awsAccountID,
		Detail: core.SQSMessageDetail{
			State:      state,
			InstanceID: TestMockInstanceId,
		},
	}
	messageSerialized, err := json.Marshal(message)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling json: %v", err)) // TODO: return error
	}

	snsTopicName := viper.GetString("SNSTopicName")

	// create an envelope and put the message in
	envelope := core.SQSEnvelope{
		MessageId: msgId,
		TopicArn:  fmt.Sprintf("arn:aws:sns:%v:%v:%v", awsRegion, awsAccountID, snsTopicName),
		Message:   string(messageSerialized),
	}

	core.Logger.Debug("NewSQSMessage returning mock msg", "sqsmessage", fmt.Sprintf("%+v", envelope))

	// serialize to a string
	envelopeSerialized, err := json.Marshal(envelope)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling json: %v", err)) // TODO: return error
	}

	envSerializedString := string(envelopeSerialized)
	*result = envSerializedString

	return msgId
}
