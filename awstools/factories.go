// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package awstools

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Ec2ServiceFactory is a function that returns ec2iface.EC2API
type Ec2ServiceFactory func(*session.Session, string) ec2iface.EC2API

// DefaultEc2ServiceFactory is the default factory of ec2iface.EC2API
func DefaultEc2ServiceFactory(assumedService *session.Session, topicRegion string) ec2iface.EC2API {
	ec2Service := ec2.New(assumedService,
		&aws.Config{
			Region: aws.String(topicRegion),
		},
	)
	return ec2Service
}

// CloudFormationServiceFactory is a function that returns cloudformationiface.CloudFormationAPI
type CloudFormationServiceFactory func(*session.Session, string) cloudformationiface.CloudFormationAPI

// DefaultCloudFormationServiceFactory is the default factory of cloudformationiface.CloudFormationAPI
func DefaultCloudFormationServiceFactory(assumedService *session.Session, topicRegion string) cloudformationiface.CloudFormationAPI {
	cloudFormationService := cloudformation.New(assumedService,
		&aws.Config{
			Region: aws.String(topicRegion),
		},
	)
	return cloudFormationService
}

// AutoScalingServiceFactory is a function that returns cloudformationiface.AutoScalingAPI
type AutoScalingServiceFactory func(*session.Session, string) autoscalingiface.AutoScalingAPI

// DefaultAutoScalingServiceFactory is the default factory of autoscalingiface.AutoScalingAPI
func DefaultAutoScalingServiceFactory(assumedService *session.Session, topicRegion string) autoscalingiface.AutoScalingAPI {
	autoScalingService := autoscaling.New(assumedService,
		&aws.Config{
			Region: aws.String(topicRegion),
		},
	)
	return autoScalingService
}
