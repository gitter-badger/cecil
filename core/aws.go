package core

import (
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
