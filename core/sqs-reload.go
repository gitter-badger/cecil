// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/tleyden/awsutil"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/tools"
	"github.com/tleyden/cecil/awstools"
)

// NewSQSPolicy returns a new SQS policy.
func (s *Service) NewSQSPolicy() *awstools.SQSPolicy {
	return awstools.NewSQSPolicy(s.SQSQueueArn())
}

// NewSQSPolicyStatement generates a new SQS queue policy statement for the given AWS account (AWSID parameter).
func (s *Service) NewSQSPolicyStatement(AWSID string) (*awstools.SQSPolicyStatement, error) {
	snsTopicName, err := tools.ViperMustGetString("SNSTopicName")
	if err != nil {
		panic(err)
	}
	return awstools.NewSQSPolicyStatement(s.SQSQueueArn(), AWSID, snsTopicName)
}

// RegenerateSQSPermissions regenerates the SQS policy adding to it every cloudaccount AWSID;
// for each Cloudaccount in the DB, allow the corresponding AWS account to send messages to the SQS queue;
// to be called after every new account is created.
func (s *Service) RegenerateSQSPermissions() error {

	var policy = s.NewSQSPolicy()

	var cloudaccounts []models.Cloudaccount

	s.DB.Where(&models.Cloudaccount{
		Provider: "aws",
	}).Find(&cloudaccounts)

	if len(cloudaccounts) == 0 {
		// No cloud accounts configured, so nothing to do
		return nil
	}

	for _, cloudaccount := range cloudaccounts {
		AWSID := cloudaccount.AWSID

		statement, err := s.NewSQSPolicyStatement(AWSID)
		if err != nil {
			Logger.Error("NewSQSPolicyStatement", "err", err)
			continue
		}

		err = policy.AddStatement(statement)
		if err != nil {
			Logger.Error("SQSPolicy.AddStatement", "err", err)
			continue
		}
	}

	if len(policy.Statement) == 0 {
		return fmt.Errorf("policy.Statement does not contain any statement")
	}

	Logger.Info("RegenerateSQSPermissions", "aws_accounts", len(policy.Statement))

	policyJSON, err := policy.JSON()
	if err != nil {
		return err
	}

	var resp *sqs.SetQueueAttributesOutput
	err = tools.Retry(10, time.Second*5, func() error {
		var err error
		resp, err = s.AWS.SQS.SetQueueAttributes(&sqs.SetQueueAttributesInput{
			Attributes: map[string]*string{
				"Policy": aws.String(policyJSON),
			},
			QueueUrl: aws.String(s.SQSQueueURL()),
		})
		return err
	}, nil)

	Logger.Info(
		"RegenerateSQSPermissions()",
		"response", resp,
	)

	return err
}

// ResubscribeToAllSNSTopics resubscribes Cecil's SQS queue to all SNS topics of the registered users.
func (s *Service) ResubscribeToAllSNSTopics() error {

	var cloudaccounts []models.Cloudaccount

	s.DB.Where(&models.Cloudaccount{
		Disabled: false,
		Provider: "aws",
	}).Find(&cloudaccounts)

	if len(cloudaccounts) == 0 {
		// No cloud accounts configured, so nothing to do
		return nil
	}

	for _, cloudaccount := range cloudaccounts {
		AWSID := cloudaccount.AWSID

		// TODO: subscribe to topic of that account

		createdSubscriptions := struct {
			mu *sync.RWMutex
			m  map[string]*sns.SubscribeOutput
		}{
			mu: &sync.RWMutex{},
			m:  make(map[string]*sns.SubscribeOutput),
		}
		createdSubscriptionsErrors := ForeachRegion(func(regionID string) error {
			resp, err := s.SubscribeToTopic(AWSID, regionID)
			if err != nil {
				return err
			}
			if resp != nil {
				createdSubscriptions.mu.Lock()
				defer createdSubscriptions.mu.Unlock()
				createdSubscriptions.m[regionID] = resp
			}
			return nil
		})

		Logger.Info(
			"ResubscribeToAllSNSTopics()",
			"response", createdSubscriptions.m,
			"errors", createdSubscriptionsErrors,
		)

		////////////////////////////////////
		listSubscriptions, listSubscriptionsErrors := s.StatusOfAllRegions(AWSID)

		Logger.Info(
			"StatusOfAllRegions()",
			"response", listSubscriptions,
			"errors", listSubscriptionsErrors,
		)

	}

	return nil
}

// ListSubscriptionsByTopic lists all the subscriptions to a specified SNS topic.
// This is later used to check whether Cecil's SQS queue is subscribed to an SNS
// topic on a specific region owned a a user.
// E.g. Is the CecilTopic on AWS account 831592357927935 in the us-east-1 region active? (ListSubscriptionsByTopic
// lists all subscriptions on that topic, and if there is CecilQueue, that region is feeding to Cecil's SQS queue.)
func (s *Service) ListSubscriptionsByTopic(AWSID string, regionID string) ([]*sns.Subscription, error) {

	topicArn := fmt.Sprintf(
		"arn:aws:sns:%v:%v:%v",
		regionID,
		AWSID,
		s.AWS.Config.SNSTopicName,
	)

	ListSubscriptionsByTopicParams := &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicArn), // Required
	}
	subscriptions := []*sns.Subscription{}
	var errString string

	// make sure to get all subscriptions by
	// iterating on eventual pages
	for {

		resp, errListSubsByTopic := s.AWS.SNS.ListSubscriptionsByTopic(ListSubscriptionsByTopicParams)

		if errListSubsByTopic != nil {
			errString += errListSubsByTopic.Error()
			errString += "; \n"
		}
		subscriptions = append(subscriptions, resp.Subscriptions...)

		if errListSubsByTopic != nil {
			if errVerifySQSPolicy := s.VerifySQSPolicy(topicArn); errVerifySQSPolicy != nil {
				errString += errVerifySQSPolicy.Error()
				errString += "; \n"
			}
		}

		if resp.NextToken != nil {
			ListSubscriptionsByTopicParams.SetNextToken(*resp.NextToken)
			continue
		} else {
			break
		}
	}

	if len(errString) > 0 {
		return subscriptions, errors.New(errString)
	}

	return subscriptions, nil
}

// Even if there is a a subscription which tells this SNS topic to try to forward events
// to an SQS, if the SQS doesn't have a policy permission to receive events, they will all be
// ignored! (see https://github.com/tleyden/cecil/issues/142).  Verify the subscription by getting
// the SQS and checking the policy permissions
func (s *Service) VerifySQSPolicy(topicArn string) error {

	policyQueueAttribute := "Policy"

	getQueueAttributesInput := &sqs.GetQueueAttributesInput{
		AttributeNames: []*string{awsutil.StringPointer(policyQueueAttribute)},
		QueueUrl:       awsutil.StringPointer(s.SQSQueueURL()),
	}
	getQueueAttributesOutput, err := s.AWS.SQS.GetQueueAttributes(getQueueAttributesInput)
	if err != nil {
		Logger.Error("SQS.GetQueueAttributes", "err", err)
		return err
	}

	policyStrPtr := getQueueAttributesOutput.Attributes[policyQueueAttribute]
	var policy awstools.SQSPolicyAttribute
	err = json.Unmarshal([]byte(*policyStrPtr), &policy)
	if err != nil {
		Logger.Error("SQS.GetQueueAttributes unmarshal policy", "err", err)
		return err
	}

	foundExpectedPolicy := false
	for _, policyStatement := range policy.Statement {
		sourceArn := policyStatement.Condition.ArnEquals["SourceArn"]
		if awstools.CheckArnsEqual(sourceArn, topicArn) {
			foundExpectedPolicy = true
		}
	}

	if !foundExpectedPolicy {
		errString := "SQS policy missing, subscription will not work.  This should fix itself soon due to periodic service.RegenerateSQSPermissions()"
		Logger.Warn("SQSPolicyMissing", "topic_arn", topicArn, "details", errString)
		return fmt.Errorf("%s", errString)
	}

	return nil

}

// SubscribeToRegions subscribes to the specified regions of the specified AWSID.
func (s *Service) SubscribeToRegions(regions []string, AWSID string) (AccountStatus, map[string]error) {
	createdSubscriptions := struct {
		mu *sync.RWMutex
		m  AccountStatus
	}{
		mu: &sync.RWMutex{},
		m:  make(AccountStatus),
	}

	createdSubscriptionsErrors := ForeachRegion(func(regionID string) error {
		isNotARequestedRegion := !tools.SliceContains(regions, regionID)
		if isNotARequestedRegion {
			// skip this region
			return nil
		}
		resp, err := s.SubscribeToTopic(AWSID, regionID)

		createdSubscriptions.mu.Lock()
		defer createdSubscriptions.mu.Unlock()

		if err != nil {
			if strings.Contains(err.Error(), "Invalid parameter: TopicArn") {
				err = errors.New("not_exists")
				createdSubscriptions.m[regionID] = RegionStatus{
					Topic:        "not_exists",
					Subscription: "not_active",
				}
			}
			// TODO: return error in response
			return err
		}
		if resp != nil {
			if resp.SubscriptionArn != nil {
				createdSubscriptions.m[regionID] = RegionStatus{
					Topic:        "exists",
					Subscription: "active",
				}
			}
		}
		return nil
	})
	return createdSubscriptions.m, createdSubscriptionsErrors
}

// RegionStatus defines the status of a single region
type RegionStatus struct {
	Topic        string `json:"topic,omitempty"`
	Subscription string `json:"subscription,omitempty"`
}

// AccountStatus defines the status of all regions of an account.
type AccountStatus map[string]RegionStatus

// StatusOfAllRegions returns the status of all regions of an account.
func (s *Service) StatusOfAllRegions(AWSID string) (AccountStatus, map[string]error) {
	listSubscriptions := struct {
		mu *sync.RWMutex
		m  AccountStatus
	}{
		mu: &sync.RWMutex{},
		m:  make(AccountStatus),
	}
	listSubscriptionsErrors := ForeachRegion(func(regionID string) error {
		resp, err := s.ListSubscriptionsByTopic(AWSID, regionID)

		listSubscriptions.mu.Lock()
		defer listSubscriptions.mu.Unlock()

		if err != nil {
			if strings.Contains(err.Error(), "Invalid parameter: TopicArn") {
				err = errors.New("not_exists")
				listSubscriptions.m[regionID] = RegionStatus{
					Topic:        "not_exists",
					Subscription: "not_active",
				}
			} else {
				listSubscriptions.m[regionID] = RegionStatus{
					Topic:        "error",
					Subscription: fmt.Sprintf("error: %v", err),
				}
			}
			// TODO: return error in response
			return err
		}
		if resp != nil {
			for _, sub := range resp {
				if sub.Endpoint == nil {
					continue
				}
				if *sub.Endpoint == s.SQSQueueArn() {
					listSubscriptions.m[regionID] = RegionStatus{
						Topic:        "exists",
						Subscription: "active",
					}
				}
			}
		}
		return nil
	})

	return listSubscriptions.m, listSubscriptionsErrors
}

// SubscribeToTopic subscribes Cecil SQS queue to the SNS topic of
// a specific region of a specific AWS account.
func (s *Service) SubscribeToTopic(AWSID string, regionID string) (*sns.SubscribeOutput, error) {
	params := &sns.SubscribeInput{
		Protocol: aws.String("sqs"), // Required
		TopicArn: aws.String(fmt.Sprintf(
			"arn:aws:sns:%v:%v:%v",
			regionID,
			AWSID,
			s.AWS.Config.SNSTopicName,
		)), // Required
		Endpoint: aws.String(fmt.Sprintf(
			"arn:aws:sqs:%v:%v:%v",
			s.AWS.Config.AWS_REGION,
			s.AWS.Config.AWS_ACCOUNT_ID,
			s.AWS.Config.SQSQueueName,
		)),
	}
	var resp *sns.SubscribeOutput
	err := tools.Retry(2, time.Second*1, func() error {
		var err error
		resp, err = s.AWS.SNS.Subscribe(params)
		return err
	}, nil)

	return resp, err
}

// Regions holds all the known regions of AWS.
var Regions = []string{
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
	"eu-west-1",
	"eu-central-1",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-south-1",
	"sa-east-1",
}

type errMap struct {
	mu *sync.RWMutex
	m  map[string]error
}

// ProcessRegion executes a specific function on a single region.
func (em errMap) ProcessRegion(regionID string, do func(regionID string) error, wg *sync.WaitGroup) {
	err := do(regionID)

	em.mu.Lock()
	defer em.mu.Unlock()

	if err != nil {
		if strings.Contains(err.Error(), "Invalid parameter: TopicArn") {
			err = errors.New("not_exists")
		}
		em.m[regionID] = err
	}

	wg.Done()
}

// ForeachRegion executes a specified function on all known regions.
func ForeachRegion(do func(regionID string) error) map[string]error {
	var mapOfErrors errMap = errMap{
		mu: &sync.RWMutex{},
		m:  make(map[string]error),
	}
	var wg sync.WaitGroup

	for regionIDIndex := range Regions {
		regionID := Regions[regionIDIndex]
		wg.Add(1)
		go mapOfErrors.ProcessRegion(regionID, do, &wg)
	}

	wg.Wait()
	return mapOfErrors.m
}

/*
var policyTest string = `
{
  "Version": "2008-10-17",
  "Id": "arn:aws:sqs:us-east-1:665102389639:CecilQueue/SQSDefaultPolicy",
  "Statement": [
    {
      "Sid": "Allow-All SQS policy",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "SQS:SendMessage",
      "Resource": "arn:aws:sqs:us-east-1:665102389639:CecilQueue",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "arn:aws:sns:*:859795398601:CecilTopic"
        }
      }
    }
  ]
}
`
*/
