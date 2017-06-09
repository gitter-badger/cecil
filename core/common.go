// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/inconshreveable/log15"
	"github.com/tleyden/awsutil"
	"github.com/tleyden/cecil/awstools"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/notification"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
)

// NewContextLogger returns a new context logger which has been filled in with the request ID
func NewContextLogger(ctx context.Context) log15.Logger {
	request := goa.ContextRequest(ctx)
	return Logger.New(
		"url", request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)
}

// sendMisconfigurationNotice sends a misconfiguration notice to emailRecipient.
func (s *Service) sendMisconfigurationNotice(err error, emailRecipient string) {
	newEmailBody, err := tools.CompileEmailTemplate(
		"misconfiguration-notice.html",
		map[string]interface{}{
			"err": err,
		},
	)
	if err != nil {
		return
	}

	s.queues.NotifierQueue().PushTask(tasks.NotifierTask{
		To:               emailRecipient,
		Subject:          "Cecil configuration problem",
		BodyHTML:         newEmailBody,
		BodyText:         newEmailBody,
		NotificationMeta: notification.NotificationMeta{NotificationType: notification.Misconfiguration},
	})
}

// CecilHTTPAddress returns the complete HTTP address of cecil; e.g. https://127.0.0.1:8080
func (s *Service) CecilHTTPAddress() string {
	// TODO check the prefix of Port; ignore port if 80 or 443 (decide looking at Scheme)
	return fmt.Sprintf("%v://%v%v",
		s.config.Server.Scheme,
		s.config.Server.HostName,
		s.config.Server.Port,
	)
}

// SQSQueueURL returns the HTTP URL of the SQS queue.
func (s *Service) SQSQueueURL() string {
	return fmt.Sprintf("https://sqs.%v.amazonaws.com/%v/%v",
		s.AWS.Config.AWS_REGION,
		s.AWS.Config.AWS_ACCOUNT_ID,
		s.AWS.Config.SQSQueueName,
	)
}

// SQSQueueArn returns the AWS ARN of the SQS queue.
func (s *Service) SQSQueueArn() string {
	return fmt.Sprintf("arn:aws:sqs:%v:%v:%v",
		s.AWS.Config.AWS_REGION,
		s.AWS.Config.AWS_ACCOUNT_ID,
		s.AWS.Config.SQSQueueName,
	)
}

// SQSPolicy defines the policy of an SQS queue.
type SQSPolicy struct {
	Version   string               `json:"Version"`
	Id        string               `json:"Id"`
	Statement []SQSPolicyStatement `json:"Statement"`
}

// SQSPolicyStatement defines a single SQS queue policy statement.
type SQSPolicyStatement struct {
	Sid       string `json:"Sid"`
	Effect    string `json:"Effect"`
	Principal string `json:"Principal"`
	Action    string `json:"Action"`
	Resource  string `json:"Resource"`
	Condition struct {
		ArnEquals map[string]string `json:"ArnEquals"`
	} `json:"Condition"`
}

// NewSQSPolicy returns a new SQS policy.
func (s *Service) NewSQSPolicy() *SQSPolicy {
	return &SQSPolicy{
		Version:   "2008-10-17",
		Id:        fmt.Sprintf("%v/SQSDefaultPolicy", s.SQSQueueArn()),
		Statement: []SQSPolicyStatement{},
	}
}

// NewSQSPolicyStatement generates a new SQS queue policy statement for the given AWS account (AWSID parameter).
func (s *Service) NewSQSPolicyStatement(AWSID string) (*SQSPolicyStatement, error) {
	if AWSID == "" {
		return &SQSPolicyStatement{}, fmt.Errorf("AWSID cannot be empty")
	}

	var condition struct {
		ArnEquals map[string]string `json:"ArnEquals"`
	}
	condition.ArnEquals = make(map[string]string, 1)

	snsTopicName, err := tools.ViperMustGetString("SNSTopicName")
	if err != nil {
		panic(err)
	}

	condition.ArnEquals["aws:SourceArn"] = fmt.Sprintf("arn:aws:sns:*:%v:%v", AWSID, snsTopicName)

	return &SQSPolicyStatement{
		Sid:       fmt.Sprintf("Allow %v to send messages", AWSID),
		Effect:    "Allow",
		Principal: "*",
		Action:    "SQS:SendMessage",
		Resource:  s.SQSQueueArn(),
		Condition: condition,
	}, nil
}

// AddStatement verifies and adds a statement to an SQS policy.
func (sp *SQSPolicy) AddStatement(statement *SQSPolicyStatement) error {
	if statement.Sid == "" {
		return fmt.Errorf("Sid cannot be empty")
	}
	if statement.Effect == "" {
		return fmt.Errorf("Effect cannot be empty")
	}
	if statement.Principal == "" {
		return fmt.Errorf("Principal cannot be empty")
	}
	if statement.Action == "" {
		return fmt.Errorf("Action cannot be empty")
	}
	if statement.Resource == "" {
		return fmt.Errorf("Resource cannot be empty")
	}
	if len(statement.Condition.ArnEquals) == 0 {
		return fmt.Errorf("Condition.ArnEquals cannot be empty")
	}
	sp.Statement = append(sp.Statement, *statement)

	return nil
}

// JSON returns the string rappresentation of the JSON of the SQS policy.
func (sp *SQSPolicy) JSON() (string, error) {
	policyJSON, err := json.Marshal(*sp)
	if err != nil {
		return "", err
	}
	return string(policyJSON), nil
}

// RegenerateSQSPermissions regenerates the SQS policy adding to it every cloudaccount AWSID;
// for each Cloudaccount in the DB, allow the corresponding AWS account to send messages to the SQS queue;
// to be called after every new account is created.
func (s *Service) RegenerateSQSPermissions() error {

	var policy *SQSPolicy = s.NewSQSPolicy()

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
		if awstools.CheckArnsEqual(policyStatement.Condition.ArnEquals.SourceArn, topicArn) {
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

// DefineLeaseDuration tries to define the duration a lease should have basing the decision
// on many sources, each of which has a hierarchy
func (s *Service) DefineLeaseDuration(accountID, cloudaccountID uint) (time.Duration, error) {
	account, err := s.GetAccountByID(int(accountID))
	if err != nil {
		return 0, err
	}

	cloudaccount, err := s.GetCloudaccountByID(int(cloudaccountID))
	if err != nil {
		return 0, err
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		return 0, errors.New("!account.IsOwnerOf(cloudaccount)")
	}

	var leaseDuration time.Duration
	// Use global cecil lease duration setting
	leaseDuration = time.Duration(s.config.Lease.Duration)

	// Use lease duration setting of account
	if account.DefaultLeaseDuration > 0 {
		leaseDuration = time.Duration(account.DefaultLeaseDuration)
	}

	// Use lease duration setting of cloudaccount
	if cloudaccount.DefaultLeaseDuration > 0 {
		leaseDuration = time.Duration(cloudaccount.DefaultLeaseDuration)
	}

	return leaseDuration, nil
}
