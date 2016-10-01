package core

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gagliardetto/simpleQueue"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var (
	TestAWSAccountID       string = "788612350743"
	TestAWSAccountRegion   string = "us-east-1"
	TestAWSAccessKeyID     string = "WwXqFLDLbHDEIaS"               // this is a random value
	TestAWSSecretAccessKey string = "jkaeLYvjHVOmGeTYLazCgjtDqznwZ" // this is a random value
)

func TestEndToEnd(t *testing.T) {

	logger = log15.New()

	viper.SetConfigFile("temporary/config.yml") // config file path
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(err)
	}

	var service Service = Service{}

	// @@@@@@@@@@@@@@@ Setup queues @@@@@@@@@@@@@@@

	service.NewLeaseQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NewLeaseQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.NewLeaseQueueConsumer error:", "error", err)
		})
	service.NewLeaseQueue.Start()
	defer service.NewLeaseQueue.Stop()

	service.TerminatorQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.TerminatorQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.TerminatorQueueConsumer error:", "error", err)
		})
	service.TerminatorQueue.Start()
	defer service.TerminatorQueue.Stop()

	service.LeaseTerminatedQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.LeaseTerminatedQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.LeaseTerminatedQueueConsumer error:", "error", err)
		})
	service.LeaseTerminatedQueue.Start()
	defer service.LeaseTerminatedQueue.Stop()

	service.ExtenderQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.ExtenderQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.ExtenderQueueConsumer error:", "error", err)
		})
	service.ExtenderQueue.Start()
	defer service.ExtenderQueue.Stop()

	service.NotifierQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NotifierQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.NotifierQueueConsumer error:", "error", err)
		})
	service.NotifierQueue.Start()
	defer service.NotifierQueue.Stop()

	// @@@@@@@@@@@@@@@ Set defaults, parse config variables @@@@@@@@@@@@@@@

	service.AWS.Config.UseMockAWS, err = viperMustGetBool("UseMockAWS")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("AWS_REGION", TestAWSAccountRegion)
	service.AWS.Config.AWS_REGION, err = viperMustGetString("AWS_REGION")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("AWS_ACCOUNT_ID", TestAWSAccountID)
	service.AWS.Config.AWS_ACCOUNT_ID, err = viperMustGetString("AWS_ACCOUNT_ID")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("AWS_ACCESS_KEY_ID", TestAWSAccessKeyID)
	service.AWS.Config.AWS_ACCESS_KEY_ID, err = viperMustGetString("AWS_ACCESS_KEY_ID")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("AWS_SECRET_ACCESS_KEY", TestAWSSecretAccessKey)
	service.AWS.Config.AWS_SECRET_ACCESS_KEY, err = viperMustGetString("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		panic(err)
	}

	service.AWS.Config.SNSTopicName, err = viperMustGetString("SNSTopicName")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.SQSQueueName, err = viperMustGetString("SQSQueueName")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.ForeignIAMRoleName, err = viperMustGetString("ForeignIAMRoleName")
	if err != nil {
		panic(err)
	}

	// Set default values for scheme, hostname, port
	viper.SetDefault("Scheme", "http")
	service.Config.Server.Scheme, err = viperMustGetString("ServerScheme")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("HostName", "0.0.0.0")
	service.Config.Server.HostName, err = viperMustGetString("ServerHostName")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("Port", ":8080")
	service.Config.Server.Port, err = viperMustGetString("ServerPort")
	if err != nil {
		panic(err)
	}

	service.Mailer.Domain, err = viperMustGetString("MailerDomain")
	if err != nil {
		panic(err)
	}
	service.Mailer.APIKey, err = viperMustGetString("MailerAPIKey")
	if err != nil {
		panic(err)
	}
	service.Mailer.PublicAPIKey, err = viperMustGetString("MailerPublicAPIKey")
	if err != nil {
		panic(err)
	}
	service.Mailer.FromAddress = fmt.Sprintf("ZeroCloud Guardian <noreply@%v>", service.Mailer.Domain)

	// Set default values for durations
	viper.SetDefault("LeaseDuration", 3*(time.Hour*24))
	service.Config.Lease.Duration, err = viperMustGetDuration("LeaseDuration")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseApprovalTimeoutDuration", 1*time.Hour)
	service.Config.Lease.ApprovalTimeoutDuration, err = viperMustGetDuration("LeaseApprovalTimeoutDuration")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("ForewarningBeforeExpiry", 12*time.Hour)
	service.Config.Lease.ForewarningBeforeExpiry, err = viperMustGetDuration("LeaseForewarningBeforeExpiry")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseMaxPerOwner", 2)
	service.Config.Lease.MaxPerOwner, err = viperMustGetInt("LeaseMaxPerOwner")
	if err != nil {
		panic(err)
	}

	// Speed everything up
	service.Config.Lease.Duration = time.Second * 10
	service.Config.Lease.ApprovalTimeoutDuration = time.Second * 3
	service.Config.Lease.ForewarningBeforeExpiry = time.Second * 3

	// some coherency tests
	if service.Config.Lease.ForewarningBeforeExpiry >= service.Config.Lease.Duration {
		panic("service.Config.Lease.ForewarningBeforeExpiry >= service.Config.Lease.Duration")
	}
	if service.Config.Lease.ApprovalTimeoutDuration >= service.Config.Lease.Duration {
		panic("service.Config.Lease.ApprovalTimeoutDuration >= service.Config.Lease.Duration")
	}

	// @@@@@@@@@@@@@@@ Setup DB @@@@@@@@@@@@@@@

	db, err := gorm.Open("sqlite3", "zerocloud.db")
	if err != nil {
		panic(err)
	}
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
	service.DB = db

	defer service.DB.Close()

	service.DB.DropTableIfExists(
		&Account{},
		&CloudAccount{},
		&Owner{},
		&Lease{},
	)
	service.DB.AutoMigrate(
		&Account{},
		&CloudAccount{},
		&Owner{},
		&Lease{},
	)

	// <EDIT-HERE>
	firstUser := Account{
		Email: "traun.leyden@gmail.com",
		CloudAccounts: []CloudAccount{
			CloudAccount{
				Provider:   "aws",
				AWSID:      TestAWSAccountID,
				ExternalID: "bigdb_zerocloud",
			},
		},
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

	// @@@@@@@@@@@@@@@ Setup external services @@@@@@@@@@@@@@@

	// setup mailer service
	mailgunInvocations := make(chan interface{}, 100)
	mockMailGun := MockMailGun{
		MailgunInvocations: mailgunInvocations,
	}
	service.Mailer.Client = &mockMailGun

	// TODO: the mock EC2 will need to get created here
	// somehow so that a wait group can get passed in

	sqsMsgsReceivedWaitGroup := sync.WaitGroup{}
	sqsMsgsReceivedWaitGroup.Add(1)

	sqsMsgsDeletedWaitGroup := sync.WaitGroup{}
	sqsMsgsDeletedWaitGroup.Add(1)

	mockSQS := NewMockSQS(
		&sqsMsgsReceivedWaitGroup,
		&sqsMsgsDeletedWaitGroup,
	)
	var messageBody string
	receiptHandle := "todo"
	NewInstanceLaunchMessage(TestAWSAccountID, TestAWSAccountRegion, &messageBody)
	messages := []*sqs.Message{
		&sqs.Message{
			Body:          &messageBody,
			ReceiptHandle: &receiptHandle,
		},
	}
	mockSQSMessage := &sqs.ReceiveMessageOutput{
		Messages: messages,
	}
	mockSQS.Enqueue(mockSQSMessage)

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

	service.AWS.SQS = mockSQS

	ec2Invocations := make(chan interface{}, 100)
	mockEc2 := NewMockEc2(ec2Invocations)
	service.EC2 = func(assumedService *session.Session, topicRegion string) ec2iface.EC2API {
		return mockEc2
	}

	// create rsa keys
	service.rsa.privateKey, service.rsa.publicKey, err = generateRSAKeys()
	if err != nil {
		panic(err)
	}

	scheduleJob(service.EventInjestorJob, time.Duration(time.Second*1))
	scheduleJob(service.AlerterJob, time.Duration(time.Second*1))
	scheduleJob(service.SentencerJob, time.Duration(time.Second*1))

	logger.Info("Waiting for sqsMsgsReceivedWaitGroup")
	sqsMsgsReceivedWaitGroup.Wait()
	logger.Info("Done waiting for sqsMsgsReceivedWaitGroup")

	logger.Info("Waiting for sqsMsgsDeletedWaitGroup")
	sqsMsgsDeletedWaitGroup.Wait()
	logger.Info("Done waiting for sqsMsgsDeletedWaitGroup")

	logger.Info("Wait for ec2InvocationDescribeInstance")
	ec2InvocationDescribeInstance := <-ec2Invocations
	logger.Info("Received ec2InvocationDescribeInstance", "ec2InvocationDescribeInstand", ec2InvocationDescribeInstance)

	logger.Info("Wait for ec2InvocationTerminateInstance")
	ec2InvocationTerminateInstance := <-ec2Invocations
	logger.Info("Recived ec2InvocationTerminateInstance", "ec2InvocationTerminateInstance", ec2InvocationTerminateInstance)

	logger.Info("Wait for mailgunInvocation")
	mailgunInvocation := <-mailgunInvocations
	logger.Info("Received mailgunInvocation", "mailgunInvocation", mailgunInvocation)

}

func NewInstanceLaunchMessage(awsAccountID, awsRegion string, result *string) {

	// create an message
	message := SQSMessage{
		Account: awsAccountID,
		Detail: SQSMessageDetail{
			State:      ec2.InstanceStateNamePending,
			InstanceID: "i-mockinstance",
		},
	}
	messageSerialized, err := json.Marshal(message)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling json: %v", err)) // TODO: return error
	}

	// create an envelope and put the message in
	envelope := SQSEnvelope{
		TopicArn: fmt.Sprintf("todo0:todo1:todo2:%v:%v", awsRegion, awsAccountID),
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
