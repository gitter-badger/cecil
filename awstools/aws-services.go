// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package awstools

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// AWSServices contains the AWS services and configuration used by Cecil
type AWSServices struct {
	Session *session.Session
	SQS     sqsiface.SQSAPI
	SNS     snsiface.SNSAPI

	EC2            Ec2ServiceFactory
	CloudFormation CloudFormationServiceFactory
	AutoScaling    AutoScalingServiceFactory

	Config struct {
		AWS_REGION            string
		AWS_ACCOUNT_ID        string
		AWS_ACCESS_KEY_ID     string
		AWS_SECRET_ACCESS_KEY string

		SNSTopicName       string
		SQSQueueName       string
		ForeignIAMRoleName string
	}
}
