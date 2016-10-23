package core

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
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

func compileEmail(tpl string, values map[string]interface{}) string {
	var emailBody bytes.Buffer // A Buffer needs no initialization.

	// TODO: check errors ???

	t := template.New("new email template")
	t, _ = t.Parse(tpl)

	_ = t.Execute(&emailBody, values)

	return emailBody.String()
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

func (s *Service) emailActionSignURL(lease_uuid, instance_id, action, token_once string) ([]byte, error) {

	bytesToSign, err := s.concatBytesFromStrings(lease_uuid, instance_id, action, token_once)
	if err != nil {
		return []byte{}, err
	}

	signature, err := s.signBytes(bytesToSign)
	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

func (s *Service) emailActionVerifySignature(c *gin.Context) error {

	var bytesToVerify bytes.Buffer

	token_once, exists := c.GetQuery("t")
	token_once = strings.TrimSpace(token_once)
	if !exists || len(token_once) == 0 {
		return fmt.Errorf("token_once is not set or null in query")
	}
	_, err := bytesToVerify.WriteString(token_once)
	if err != nil {
		return err
	}

	action, exists := c.Params.Get("action")
	if !exists || len(action) == 0 {
		return fmt.Errorf("action is not set or null in query")
	}
	_, err = bytesToVerify.WriteString(action)
	if err != nil {
		return err
	}

	lease_uuid, exists := c.Params.Get("lease_uuid")
	if !exists || len(lease_uuid) == 0 {
		return fmt.Errorf("lease_uuid is not set or null in query")
	}
	_, err = bytesToVerify.WriteString(lease_uuid)
	if err != nil {
		return err
	}

	instance_id, exists := c.Params.Get("instance_id")
	if !exists || len(instance_id) == 0 {
		return fmt.Errorf("instance_id is not set or null in query")
	}
	_, err = bytesToVerify.WriteString(instance_id)
	if err != nil {
		return err
	}

	signature_base64, exists := c.GetQuery("s")
	signature_base64 = strings.TrimSpace(signature_base64)
	if !exists || len(signature_base64) == 0 {
		return fmt.Errorf("signature is not set or null in query")
	}

	signature, err := base64.URLEncoding.DecodeString(signature_base64)
	if err != nil {
		return err
	}

	return s.verifyBytes(bytesToVerify.Bytes(), signature)
}

func (s *Service) EmailActionGenerateSignedURL(action, lease_uuid, instance_id, token_once string) (string, error) {
	signature, err := s.emailActionSignURL(lease_uuid, instance_id, action, token_once)
	if err != nil {
		return "", fmt.Errorf("error while signing")
	}
	signedURL := fmt.Sprintf("%s/email_action/leases/%s/%s/%s?t=%s&s=%s",
		s.ZeroCloudHTTPAddress(),
		lease_uuid,
		instance_id,
		action,
		token_once,
		base64.URLEncoding.EncodeToString(signature),
	)
	return signedURL, nil
}

func (s *Service) GenerateAPITokenForAccount(accountID string) (string, error) {
	// generate claims
	// sign claims
	// return fmt.Sprintf("key-%v.%v",claims.Bytes().Base64(), signature.Base64())

	account_id, err := strconv.ParseInt(accountID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid account id value")
	}

	type Claims struct {
		AccountID int64 `json:"account_id"`
		IAT       int64 `json:"iat"`
	}

	claims := Claims{
		AccountID: account_id,
		IAT:       time.Now().UTC().Unix(),
	}

	claimsToSign, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	signature, err := s.signBytes(claimsToSign)
	if err != nil {
		return "", err
	}

	APIToken := fmt.Sprintf("key-%v.%v", base64.URLEncoding.EncodeToString(claimsToSign), base64.URLEncoding.EncodeToString(signature))

	return APIToken, nil
}

func (s *Service) VerifyAPIToken(APIToken string) (string, error) {
	// split APIToken string by - and then by .
	// convert to bytes the claims and the signature
	// verify signature

	/*
		signature, err := base64.URLEncoding.DecodeString(signature_base64)
		if err != nil {
			return err
		}
	*/
	return "", nil
}

func generateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey
	var err error

	// generate Private Key
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return &rsa.PrivateKey{}, &rsa.PublicKey{}, err
	}

	// precompute some calculations
	privateKey.Precompute()

	// validate Private Key
	if err = privateKey.Validate(); err != nil {
		return &rsa.PrivateKey{}, &rsa.PublicKey{}, err
	}

	// public key address of RSA key
	publicKey = &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func (s *Service) concatBytesFromStrings(str ...string) ([]byte, error) {
	var concatBytesBuffer bytes.Buffer
	for _, stringValue := range str {
		_, err := concatBytesBuffer.WriteString(stringValue)
		if err != nil {
			return []byte{}, err
		}
	}
	return concatBytesBuffer.Bytes(), nil
}

func (s *Service) signBytes(bytesToSign []byte) ([]byte, error) {
	if s.rsa.privateKey == nil {
		return []byte{}, fmt.Errorf("s.rsa.privateKey is nil")
	}

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := crypto.SHA256.New()
	pssh.Write(bytesToSign)
	hashed := pssh.Sum(nil)

	// compute signature
	return rsa.SignPSS(rand.Reader, s.rsa.privateKey, crypto.SHA256, hashed, &opts)
}

func (s *Service) verifyBytes(bytesToVerify []byte, signature []byte) error {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := crypto.SHA256.New()

	pssh.Write(bytesToVerify)
	hashed := pssh.Sum(nil)

	// verify signature
	return rsa.VerifyPSS(s.rsa.publicKey, crypto.SHA256, hashed, signature, &opts)
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

/*
TODO: use this:

func (a *Account) FetchCLoudAccountByID(cloudAccountID string) (*CloudAccount, error) {}

*/

func (s *Service) FetchAccountByID(accountID string) (*Account, error) {
	// parse parameters
	account_id, err := strconv.ParseUint(accountID, 10, 64)
	if err != nil {
		return &Account{}, fmt.Errorf("invalid account id")
	}

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the account exists
	var accountCount int64
	var account Account
	s.DB.First(&account, uint(account_id)).Count(&accountCount)
	if accountCount != 1 {
		return &Account{}, fmt.Errorf("account not found")
	}

	if uint(account_id) != account.ID {
		return &Account{}, fmt.Errorf("account not found")
	}

	return &account, nil
}

func (s *Service) AccountByEmailExists(accountEmail string) bool {

	var accountCount int64
	var account Account
	s.DB.Where(&Account{Email: accountEmail}).First(&account).Count(&accountCount)
	if accountCount != 1 {
		return false
	}

	if accountEmail != account.Email {
		return false
	}

	return true
}

func (s *Service) ZeroCloudHTTPAddress() string {
	// TODO check the prefix of Port; ignore port if 80 or 443 (decide looking at Scheme)
	return fmt.Sprintf("%v://%v%v",
		s.Config.Server.Scheme,
		s.Config.Server.HostName,
		s.Config.Server.Port,
	)
}

func (s *Service) FetchCloudAccountByID(cloudAccountID string) (*CloudAccount, error) {
	cloudaccount_id, err := strconv.ParseUint(cloudAccountID, 10, 64)
	if err != nil {
		return &CloudAccount{}, fmt.Errorf("invalid cloudAccount id")
	}

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the cloudaccount exists
	var cloudAccountCount int64
	var cloudAccount CloudAccount
	s.DB.First(&cloudAccount, uint(cloudaccount_id)).Count(&cloudAccountCount)
	if cloudAccountCount != 1 {
		return &CloudAccount{}, fmt.Errorf("cloudAccount not found")
	}

	if uint(cloudaccount_id) != cloudAccount.ID {
		return &CloudAccount{}, fmt.Errorf("cloudAccount not found")
	}

	return &cloudAccount, nil
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
		From:     s.Mailer.FromAddress,
		To:       emailRecipient,
		Subject:  "ZeroCloud configuration problem",
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
	}
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

	resp, err := s.AWS.SQS.SetQueueAttributes(&sqs.SetQueueAttributesInput{
		Attributes: map[string]*string{
			"Policy": aws.String(policyJSON),
		},
		QueueUrl: aws.String(s.SQSQueueURL()),
	})
	logger.Info("RegenerateSQSPermissions()",
		"response", resp)

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
