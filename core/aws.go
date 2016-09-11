package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/spf13/viper"
	"github.com/tleyden/zerocloud/mocks/aws"
)

func NewEc2Service(assumedService *session.Session, topicRegion string) ec2iface.EC2API {

	switch viper.GetBool("UseMockAWS") {
	case true:
		mockEc2 := &mockaws.MockEc2{}
		return mockEc2
	default:
		ec2Service := ec2.New(assumedService,
			&aws.Config{
				Region: aws.String(topicRegion),
			},
		)
		return ec2Service
	}

}
