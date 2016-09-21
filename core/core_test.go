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
	TestAWSAccountID     string = "788612350743"
	TestAWSAccountRegion string = "us-east-1"
)

func TestEndToEnd(t *testing.T) {

	logger = log15.New()

	// Speed everything up
	ZCDefaultLeaseDuration = time.Second * 10
	ZCDefaultLeaseApprovalTimeoutDuration = time.Second * 3
	ZCDefaultForewarningBeforeExpiry = time.Second * 3

	viper.SetConfigFile("temporary/config.yml") // name of config file (without extension)
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
		&Lease{},
		&Region{},
		&Owner{},
	)
	service.DB.AutoMigrate(
		&Account{},
		&CloudAccount{},
		&Lease{},
		&Region{},
		&Owner{},
	)

	// <EDIT-HERE>
	firstUser := Account{
		Email: "traun.leyden@gmail.com",
		CloudAccounts: []CloudAccount{
			CloudAccount{
				Provider:   "aws",
				AWSID:      TestAWSAccountID,
				ExternalID: "bigdb_zerocloud",
				Regions: []Region{
					Region{
						Region: TestAWSAccountRegion,
					},
				},
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
	service.Mailer = &mockMailGun

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
		viper.GetString("AWS_ACCESS_KEY_ID"),
		viper.GetString("AWS_SECRET_ACCESS_KEY"),
		"",
	)
	AWSConfig := &aws.Config{
		Credentials: AWSCreds,
	}
	service.AWS.Session = session.New(AWSConfig)

	service.AWS.SQS = mockSQS

	ec2WaitGroup := sync.WaitGroup{}
	ec2WaitGroup.Add(2)
	ec2Invocations := make(chan interface{}, 100)
	mockEc2 := NewMockEc2(&ec2WaitGroup, ec2Invocations)
	service.EC2 = func(assumedService *session.Session, topicRegion string) ec2iface.EC2API {
		return mockEc2
	}

	// create rsa keys
	service.rsa.privateKey, service.rsa.publicKey, err = generateRSAKeys()
	if err != nil {
		panic(err)
	}

	go scheduleJob(service.EventInjestorJob, time.Duration(time.Second*1))
	go scheduleJob(service.AlerterJob, time.Duration(time.Second*1))
	go scheduleJob(service.SentencerJob, time.Duration(time.Second*1))

	logger.Info("Waiting for sqsMsgsReceivedWaitGroup")
	sqsMsgsReceivedWaitGroup.Wait()
	logger.Info("Done waiting for sqsMsgsReceivedWaitGroup")

	logger.Info("Waiting for sqsMsgsDeletedWaitGroup")
	sqsMsgsDeletedWaitGroup.Wait()
	logger.Info("Done waiting for sqsMsgsDeletedWaitGroup")

	ec2Invocation := <-ec2Invocations
	logger.Info("ec2Invocation", "ec2Invocation", ec2Invocation)

	mailgunInvocation := <-mailgunInvocations

	logger.Info("mailgunInvocation", "mailgunInvocation", mailgunInvocation)

	logger.Info("Waiting for ec2 wait group")
	ec2WaitGroup.Wait()
	logger.Info("Done waiting for ec2 wait group")

	// TODO: get the calls to mockec2 made by zerocloud and
	// make sure they are what is expected

	// TODO: ditto for mailgun mock

	// logger.Info("Waiting for timer")
	// time.Sleep(50 * time.Second)
	// logger.Info("Done waiting for timer")

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
