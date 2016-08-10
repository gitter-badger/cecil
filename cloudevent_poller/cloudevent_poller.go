package cloudevent_poller

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/inconshreveable/log15"
)

var logger log15.Logger

func init() {
	logger = log15.New()
}

type CloudEventPoller struct {
	SQSQueueURL     string
	ZeroCloudAPIURL string
	AWSRegion       string
}

func (p *CloudEventPoller) Run() error {
	logger.Info("Run() called.", "poller", fmt.Sprintf("%+v", p))

	// connect to SQS queue
	session, err := session.NewSession()
	if err != nil {
		return err
	}
	logger.Info("Session", "session", fmt.Sprintf("%+v", session))
	logger.Info("Session", "config", fmt.Sprintf("%+v", session.Config))
	logger.Info("Session", "config region", fmt.Sprintf("%+v", *session.Config.Region))
	logger.Info("Session", "config credentials", fmt.Sprintf("%+v", session.Config.Credentials))

	sqsService := sqs.New(session, &aws.Config{Region: aws.String(p.AWSRegion)})
	logger.Info("sqs service", "sqs", fmt.Sprintf("%+v", sqsService))

	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(p.SQSQueueURL),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	}
	resp, err := sqsService.ReceiveMessage(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil
	}
	logger.Info("resp", "resp", fmt.Sprintf("%+v", resp))

	// pull any items off queue

	// transform json

	// push to zerocloud rest API

	// }

	return nil

}
