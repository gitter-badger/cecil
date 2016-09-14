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

func DisabledTestEndToEnd(t *testing.T) {

	logger = log15.New()

	viper.SetConfigFile("temporary/config.yml") // name of config file (without extension)
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(err)
	}

	var service Service = Service{}

	// @@@@@@@@@@@@@@@ Setup queues @@@@@@@@@@@@@@@

	service.NewLeaseQueue = simpleQueue.NewQueue()
	service.NewLeaseQueue.SetMaxSize(maxQueueSize)
	service.NewLeaseQueue.SetWorkers(maxWorkers)
	service.NewLeaseQueue.Consumer = service.NewLeaseQueueConsumer
	service.NewLeaseQueue.Start()
	defer service.NewLeaseQueue.Stop()

	service.TerminatorQueue = simpleQueue.NewQueue()
	service.TerminatorQueue.SetMaxSize(maxQueueSize)
	service.TerminatorQueue.SetWorkers(maxWorkers)
	service.TerminatorQueue.Consumer = service.TerminatorQueueConsumer
	service.TerminatorQueue.Start()
	defer service.TerminatorQueue.Stop()

	service.LeaseTerminatedQueue = simpleQueue.NewQueue()
	service.LeaseTerminatedQueue.SetMaxSize(maxQueueSize)
	service.LeaseTerminatedQueue.SetWorkers(maxWorkers)
	service.LeaseTerminatedQueue.Consumer = service.LeaseTerminatedQueueConsumer
	service.LeaseTerminatedQueue.Start()
	defer service.LeaseTerminatedQueue.Stop()

	service.ExtenderQueue = simpleQueue.NewQueue()
	service.ExtenderQueue.SetMaxSize(maxQueueSize)
	service.ExtenderQueue.SetWorkers(maxWorkers)
	service.ExtenderQueue.Consumer = service.ExtenderQueueConsumer
	service.ExtenderQueue.Start()
	defer service.ExtenderQueue.Stop()

	service.NotifierQueue = simpleQueue.NewQueue()
	service.NotifierQueue.SetMaxSize(maxQueueSize)
	service.NotifierQueue.SetWorkers(maxWorkers)
	service.NotifierQueue.Consumer = service.NotifierQueueConsumer
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
	service.Mailer = &MockMailGun{}

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
	/*
		// A list of messages.
		Messages []*Message `locationNameList:"Message" type:"list" flattened:"true"`
	*/
	mockSQS.Enqueue(mockSQSMessage)

	// setup aws session -- TODO: mock this out
	AWSCreds := credentials.NewStaticCredentials(viper.GetString("AWS_ACCESS_KEY_ID"), viper.GetString("AWS_SECRET_ACCESS_KEY"), "")
	AWSConfig := &aws.Config{
		Credentials: AWSCreds,
	}
	service.AWS.Session = session.New(AWSConfig)

	service.AWS.SQS = mockSQS

	ec2WaitGroup := sync.WaitGroup{}
	ec2WaitGroup.Add(2)
	mockEc2 := NewMockEc2(&ec2WaitGroup)
	service.EC2 = func(assumedService *session.Session, topicRegion string) ec2iface.EC2API {
		return mockEc2
	}

	go scheduleJob(service.EventInjestorJob, time.Duration(time.Second*5))
	go scheduleJob(service.AlerterJob, time.Duration(time.Second*60))
	go scheduleJob(service.SentencerJob, time.Duration(time.Second*60))

	logger.Info("Waiting for sqsMsgsReceivedWaitGroup")
	sqsMsgsReceivedWaitGroup.Wait()
	logger.Info("Done waiting for sqsMsgsReceivedWaitGroup")

	logger.Info("Waiting for sqsMsgsDeletedWaitGroup")
	sqsMsgsDeletedWaitGroup.Wait()
	logger.Info("Done waiting for sqsMsgsDeletedWaitGroup")

	logger.Info("Waiting for ec2 wait group")
	ec2WaitGroup.Wait()
	logger.Info("Done waiting for ec2 wait group")

	// TODO: get the calls to mockec2 made by zerocloud and
	// make sure they are what is expected

	// TODO: ditto for mailgun mock

	logger.Info("Waiting for timer")
	time.Sleep(50 * time.Second)
	logger.Info("Done waiting for timer")

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
