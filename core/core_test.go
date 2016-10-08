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
	"github.com/inconshreveable/log15"
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
	service.SetupDB()
	defer service.Stop()

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

	// @@@@@@@@@@@@@@@ Setup mock external services @@@@@@@@@@@@@@@

	// setup mailer service
	mockMailGun := NewMockMailGun()
	service.Mailer.Client = mockMailGun

	// Create a mock SQS that will return a message indicating that an EC2 instance was luanched
	mockSQS := NewMockSQS()
	var messageBody string
	receiptHandle := "mockReceiptHandle"
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

	mockEc2 := NewMockEc2()
	service.EC2 = func(assumedService *session.Session, topicRegion string) ec2iface.EC2API {
		return mockEc2
	}

	// @@@@@@@@@@@@@@@ Schedule Periodic Jobs @@@@@@@@@@@@@@@

	schedulePeriodicJob(service.EventInjestorJob, time.Duration(time.Second*1))
	schedulePeriodicJob(service.AlerterJob, time.Duration(time.Second*1))
	schedulePeriodicJob(service.SentencerJob, time.Duration(time.Second*1))

	// @@@@@@@@@@@@@@@ Wait for Test actions To Finish @@@@@@@@@@@@@@@

	// Wait until the SQS message is sent back to the eventinjestor
	mockSQS.waitForReceivedMessageInput()
	mockSQS.waitForDeletedMessageInput(receiptHandle)

	// Wait until the event injestor tries to describe the instance
	mockEc2.waitForDescribeInstancesInput()

	// Wait until the Sentencer tries to terminate the instance
	mockEc2.waitForTerminateInstancesInput()

	// Wait until the Sentencer tries to notifies admin that the instance was terminated
	mailGunInvocation := <-mockMailGun.MailgunInvocations

	logger.Info("Received mailgunInvocation", "mailgunInvocation", mailGunInvocation)

	logger.Info("CoreTest finished")

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
		TopicArn: fmt.Sprintf("arn:aws:sns:%v:%v:ZeroCloudTopic", awsRegion, awsAccountID),
		Message:  string(messageSerialized),
	}
	// TODO: replace ZeroCloudTopic with service.AWS.COnfig.SNSTopicName

	// serialize to a string
	envelopeSerialized, err := json.Marshal(envelope)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling json: %v", err)) // TODO: return error
	}

	envSerializedString := string(envelopeSerialized)
	*result = envSerializedString
}
