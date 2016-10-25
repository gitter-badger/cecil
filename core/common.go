package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func schedulePeriodicJob(job func() error, runEvery time.Duration) {
	go func() {
		for {
			err := job()
			if err != nil {
				logger.Error("schedulePeriodicJob", "error", err)
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

func compileEmail(tpl string, values map[string]interface{}) string {
	var emailBody bytes.Buffer // A Buffer needs no initialization.

	// TODO: check errors ???

	t := template.New("new email template")
	t, _ = t.Parse(tpl)

	_ = t.Execute(&emailBody, values)

	return emailBody.String()
}

func viperIsSet(key string) bool {
	if !viper.IsSet(key) {
		logger.Crit("Config parameter not set",
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

func (s *Service) FetchAccountByID(accountIDString string) (*Account, error) {
	// parse parameters
	accountID, err := strconv.ParseUint(accountIDString, 10, 64)
	if err != nil {
		return &Account{}, fmt.Errorf("invalid account id")
	}

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the account exists
	var accountCount int64
	var account Account
	err = s.DB.Table("accounts").Where("id = ?", uint(accountID)).Count(&accountCount).Find(&account).Error
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

func (s *Service) FetchCloudAccountByID(cloudAccountIDString string) (*CloudAccount, error) {
	// parse parameters
	cloudAccountID, err := strconv.ParseUint(cloudAccountIDString, 10, 64)
	if err != nil {
		return &CloudAccount{}, fmt.Errorf("invalid cloudAccount id")
	}

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the cloudAccount exists
	var cloudAccountCount int64
	var cloudAccount CloudAccount
	err = s.DB.Find(&cloudAccount, uint(cloudAccountID)).Count(&cloudAccountCount).Error
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
	newEmailBody := compileEmail(
		`Hey it appears that ZeroCloud is mis-configured.
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
		Subject:          "ZeroCloud configuration problem",
		BodyHTML:         newEmailBody,
		BodyText:         newEmailBody,
		NotificationMeta: NotificationMeta{NotificationType: Misconfiguration},
	}
}

func (s *Service) ZeroCloudHTTPAddress() string {
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
	condition.ArnEquals["aws:SourceArn"] = fmt.Sprintf("arn:aws:sns:*:%v:ZeroCloudTopic", AWSID)

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

	logger.Info("RegenerateSQSPermissions", "aws_accounts", len(policy.Statement))

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

	logger.Info(
		"RegenerateSQSPermissions()",
		"response", resp,
	)

	return err
}

/*
var policyTest string = `
{
  "Version": "2008-10-17",
  "Id": "arn:aws:sqs:us-east-1:665102389639:ZeroCloudQueue/SQSDefaultPolicy",
  "Statement": [
    {
      "Sid": "Allow-All SQS policy",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "SQS:SendMessage",
      "Resource": "arn:aws:sqs:us-east-1:665102389639:ZeroCloudQueue",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "arn:aws:sns:*:859795398601:ZeroCloudTopic"
        }
      }
    }
  ]
}
`
*/
