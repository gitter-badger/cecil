package core

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type Ec2ServiceFactory func(*session.Session, string) ec2iface.EC2API

func DefaultEc2ServiceFactory(assumedService *session.Session, topicRegion string) ec2iface.EC2API {
	ec2Service := ec2.New(assumedService,
		&aws.Config{
			Region: aws.String(topicRegion),
		},
	)
	return ec2Service

}

// TODO: doesn't aws sdk provide this?
type SQSEnvelope struct {
	Type             string    `json:"Type"`
	MessageId        string    `json:"MessageId"`
	TopicArn         string    `json:"TopicArn"`
	Message          string    `json:"Message"`
	Timestamp        time.Time `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
	SubscribeURL     string    `json:"SubscribeURL"`
	UnsubscribeURL   string    `json:"UnsubscribeURL"`
}

type SQSMessageDetail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}

// TODO: doesn't aws sdk provide this?
type SQSMessage struct {
	Version    string           `json:"version"`
	ID         string           `json:"id"`
	DetailType string           `json:"detail-type"`
	Source     string           `json:"source"`
	Account    string           `json:"account"`
	Time       time.Time        `json:"time"`
	Region     string           `json:"region"`
	Resources  []string         `json:"resources"`
	Detail     SQSMessageDetail `json:"detail"`
}
