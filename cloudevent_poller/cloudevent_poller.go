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
	SQSQueueTopicARN string
	ZeroCloudAPIURL  string
	AWSRegion        string
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

	for {

		// pull any items off queue

		// transform json

		// push to zerocloud rest API

	}

}
