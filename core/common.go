package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func schedulePeriodicJob(job func() error, runEvery time.Duration) {
	go func() {
		for {
			err := job()
			if err != nil {
				Logger.Error("schedulePeriodicJob", "error", err)
			}
			time.Sleep(runEvery)
		}
	}()
}

func retry(attempts int, sleep time.Duration, callback func() error) (err error) {
	for i := 1; i <= attempts; i++ {

		err = callback()
		if err == nil {
			return nil
		}
		time.Sleep(sleep)

		fmt.Println("Retry error: ", err)
	}
	return fmt.Errorf("Abandoned after %d attempts, last error: %s", attempts, err)
}

func CompileEmail(tpl string, values map[string]interface{}) string {
	var emailBody bytes.Buffer // A Buffer needs no initialization.

	// TODO: check errors ???

	t := template.New("new email template")
	t, _ = t.Parse(tpl)

	_ = t.Execute(&emailBody, values)

	return emailBody.String()
}

func viperIsSet(key string) bool {
	if !viper.IsSet(key) {
		Logger.Crit("Config parameter not set",
			key, viper.Get(key),
		)
		return false
	}
	return true
}

func viperMustGetString(key string) (string, error) {
	if !viper.IsSet(key) {
		return "", fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetString(key), nil
}

func viperMustGetInt(key string) (int, error) {
	if !viper.IsSet(key) {
		return 0, fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetInt(key), nil
}

func viperMustGetInt64(key string) (int64, error) {
	if !viper.IsSet(key) {
		return 0, fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetInt64(key), nil
}

func viperMustGetBool(key string) (bool, error) {
	if !viper.IsSet(key) {
		return false, fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetBool(key), nil
}

func viperMustGetStringMapString(key string) (map[string]string, error) {
	if !viper.IsSet(key) {
		return map[string]string{}, fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetStringMapString(key), nil
}

func viperMustGetDuration(key string) (time.Duration, error) {
	if !viper.IsSet(key) {
		return time.Duration(0), fmt.Errorf("viper config param not set: %v", key)
	}
	return viper.GetDuration(key), nil
}

func AskForConfirmation() bool {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("fatal: ", err)
	}
	positive := []string{"y", "Y", "yes", "Yes", "YES"}
	negative := []string{"n", "N", "no", "No", "NO"}
	if SliceContains(positive, input) {
		return true
	} else if SliceContains(negative, input) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter.")
		return AskForConfirmation()
	}
}

func SliceContains(slice []string, element string) bool {
	for _, elem := range slice {
		if strings.EqualFold(element, elem) {
			return true
		}
	}
	return false
}

func (s *Service) FetchAccountByID(accountID int) (*Account, error) {

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the account exists
	var accountCount int64
	var account Account
	err := s.DB.Table("accounts").Where("id = ?", uint(accountID)).Count(&accountCount).Find(&account).Error
	if err == gorm.ErrRecordNotFound {
		return &Account{}, err
	}

	if uint(accountID) != account.ID {
		return &Account{}, gorm.ErrRecordNotFound
	}

	return &account, err
}

func (s *Service) AccountByEmailExists(accountEmail string) (bool, error) {
	var account Account
	err := s.DB.Where(&Account{Email: accountEmail}).Find(&account).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

func (s *Service) FetchCloudAccountByID(cloudAccountID int) (*CloudAccount, error) {

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the cloudAccount exists
	var cloudAccountCount int64
	var cloudAccount CloudAccount
	err := s.DB.Find(&cloudAccount, uint(cloudAccountID)).Count(&cloudAccountCount).Error
	if err == gorm.ErrRecordNotFound {
		return &CloudAccount{}, err
	}

	if uint(cloudAccountID) != cloudAccount.ID {
		return &CloudAccount{}, gorm.ErrRecordNotFound
	}

	return &cloudAccount, err
}

func (s *Service) CloudAccountByAWSIDExists(AWSID string) (bool, error) {
	var cloudAccount CloudAccount
	err := s.DB.Where(&CloudAccount{AWSID: AWSID}).Find(&cloudAccount).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

func (a *Account) IsOwnerOf(cloudAccount *CloudAccount) bool {
	return a.ID == cloudAccount.AccountID
}

func (s *Service) sendMisconfigurationNotice(err error, emailRecipient string) {
	newEmailBody := CompileEmail(
		`Hey it appears that Cecil is mis-configured.
		<br>
		<br>
		Error:
		<br>
		{{.err}}`,
		map[string]interface{}{
			"err": err,
		},
	)

	s.NotifierQueue.TaskQueue <- NotifierTask{
		From:             s.Mailer.FromAddress,
		To:               emailRecipient,
		Subject:          "Cecil configuration problem",
		BodyHTML:         newEmailBody,
		BodyText:         newEmailBody,
		NotificationMeta: NotificationMeta{NotificationType: Misconfiguration},
	}
}

func (s *Service) CecilHTTPAddress() string {
	// TODO check the prefix of Port; ignore port if 80 or 443 (decide looking at Scheme)
	return fmt.Sprintf("%v://%v%v",
		s.Config.Server.Scheme,
		s.Config.Server.HostName,
		s.Config.Server.Port,
	)
}

func (s *Service) SQSQueueURL() string {
	return fmt.Sprintf("https://sqs.%v.amazonaws.com/%v/%v",
		s.AWS.Config.AWS_REGION,
		s.AWS.Config.AWS_ACCOUNT_ID,
		s.AWS.Config.SQSQueueName,
	)
}

func (s *Service) SQSQueueArn() string {
	return fmt.Sprintf("arn:aws:sqs:%v:%v:%v",
		s.AWS.Config.AWS_REGION,
		s.AWS.Config.AWS_ACCOUNT_ID,
		s.AWS.Config.SQSQueueName,
	)
}

type SQSPolicy struct {
	Version   string               `json:"Version"`
	Id        string               `json:"Id"`
	Statement []SQSPolicyStatement `json:"Statement"`
}

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

func (s *Service) NewSQSPolicy() *SQSPolicy {
	return &SQSPolicy{
		Version:   "2008-10-17",
		Id:        fmt.Sprintf("%v/SQSDefaultPolicy", s.SQSQueueArn()),
		Statement: []SQSPolicyStatement{},
	}
}

func (s *Service) NewSQSPolicyStatement(AWSID string) (*SQSPolicyStatement, error) {
	if AWSID == "" {
		return &SQSPolicyStatement{}, fmt.Errorf("AWSID cannot be empty")
	}

	var condition struct {
		ArnEquals map[string]string `json:"ArnEquals"`
	}
	condition.ArnEquals = make(map[string]string, 1)

	snsTopicName, err := viperMustGetString("SNSTopicName")
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

func (sp *SQSPolicy) JSON() (string, error) {
	policyJSON, err := json.Marshal(*sp)
	if err != nil {
		return "", err
	}
	return string(policyJSON), nil
}

// RegenerateSQSPermissions regenerates the SQS policy adding to it every cloudAccount AWSID;
// for each CloudAccount in the DB, allow the corresponding AWS account to send messages to the SQS queue;
// to be called after every new account is created.
func (s *Service) RegenerateSQSPermissions() error {

	var policy *SQSPolicy = s.NewSQSPolicy()

	var cloudAccounts []CloudAccount

	s.DB.Where(&CloudAccount{
		Disabled: false,
		Provider: "aws",
	}).Find(&cloudAccounts)

	for _, cloudAccount := range cloudAccounts {
		AWSID := cloudAccount.AWSID

		statement, err := s.NewSQSPolicyStatement(AWSID)
		if err != nil {
			// TODO: notify ZC admins
			continue
		}

		err = policy.AddStatement(statement)
		if err != nil {
			// TODO: notify ZC admins
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
	err = retry(10, time.Second*5, func() error {
		var err error
		resp, err = s.AWS.SQS.SetQueueAttributes(&sqs.SetQueueAttributesInput{
			Attributes: map[string]*string{
				"Policy": aws.String(policyJSON),
			},
			QueueUrl: aws.String(s.SQSQueueURL()),
		})
		return err
	})

	Logger.Info(
		"RegenerateSQSPermissions()",
		"response", resp,
	)

	return err
}

func (s *Service) ResubscribeToAllSNSTopics() error {

	var cloudAccounts []CloudAccount

	s.DB.Where(&CloudAccount{
		Disabled: false,
		Provider: "aws",
	}).Find(&cloudAccounts)

	for _, cloudAccount := range cloudAccounts {
		AWSID := cloudAccount.AWSID

		// TODO: subscribe to topic of that account

		createdSubscriptions := make(map[string]*sns.SubscribeOutput)
		// not using a mutex; each key of map is handled by just one goroutine
		createdSubscriptionsErrors := ForeachRegion(func(regionID string) error {
			resp, err := s.SubscribeToTopic(AWSID, regionID)
			if err != nil {
				return err
			}
			createdSubscriptions[regionID] = resp
			return nil
		})

		Logger.Info(
			"ResubscribeToAllSNSTopics()",
			"response", createdSubscriptions,
			"errors", createdSubscriptionsErrors,
		)

		////////////////////////////////////
		listSubscriptions := make(map[string][]*sns.Subscription)
		// not using a mutex; each key of map is handled by just one goroutine
		listSubscriptionsErrors := ForeachRegion(func(regionID string) error {
			resp, err := s.ListSubscriptionsByTopic(AWSID, regionID)
			if err != nil {
				return err
			}
			if resp != nil {
				listSubscriptions[regionID] = resp
			}
			return nil
		})
		Logger.Info(
			"ListSubscriptionsByTopic()",
			"response", listSubscriptions,
			"errors", listSubscriptionsErrors,
		)

	}

	return nil
}

func (s *Service) ListSubscriptionsByTopic(AWSID string, regionID string) ([]*sns.Subscription, error) {
	ListSubscriptionsByTopicParams := &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(fmt.Sprintf(
			"arn:aws:sns:%v:%v:%v",
			regionID,
			AWSID,
			s.AWS.Config.SNSTopicName,
		)), // Required
	}
	subscriptions := make([]*sns.Subscription, 0, 10) // TODO: what length set the array to???
	var errString string

	for {

		resp, err := s.AWS.SNS.ListSubscriptionsByTopic(ListSubscriptionsByTopicParams)

		if err != nil {
			errString += err.Error()
			errString += "; \n"
		}
		subscriptions = append(subscriptions, resp.Subscriptions...)

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
	err := retry(2, time.Second*1, func() error {
		var err error
		resp, err = s.AWS.SNS.Subscribe(params)
		return err
	})

	return resp, err
}

/*
	us-east-1 // US East (N. Virginia)
	us-east-2 // US East (Ohio)
	us-west-1 // US West (N. California)
	us-west-2 // US West (Oregon)
	eu-west-1 // EU (Ireland)
	eu-central-1 // EU (Frankfurt)
	ap-northeast-1 // Asia Pacific (Tokyo)
	ap-northeast-2 // Asia Pacific (Seoul)
	ap-southeast-1 // Asia Pacific (Singapore)
	ap-southeast-2 // Asia Pacific (Sydney)
	ap-south-1 // Asia Pacific (Mumbai)
	sa-east-1 // South America (São Paulo)
*/

var regions []string = []string{
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

type errMap map[string]error

func (em errMap) ProcessRegion(regionID string, do func(regionID string) error, wg *sync.WaitGroup) {
	err := do(regionID)

	if err != nil {
		if strings.Contains(err.Error(), "Invalid parameter: TopicArn") {
			err = errors.New("not_exists")
		}
		em[regionID] = err
	}

	wg.Done()
}
func ForeachRegion(do func(regionID string) error) map[string]error {
	var mapOfErrors errMap = make(errMap)
	var wg sync.WaitGroup

	for _, regionID := range regions {
		wg.Add(1)
		go mapOfErrors.ProcessRegion(regionID, do, &wg)
	}

	wg.Wait()
	return mapOfErrors
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
