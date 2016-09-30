package core

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func scheduleJob(f func() error, runEvery time.Duration) {
	go func() {
		for {
			err := f()
			if err != nil {
				logger.Error("scheduleJob", "error", err)
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

func (s *Service) sign(lease_uuid, instance_id, action, token_once string) ([]byte, error) {

	var bytesToSign bytes.Buffer

	if s.rsa.privateKey == nil {
		return nil, fmt.Errorf("s.rsa.privateKey is nil")
	}

	_, err := bytesToSign.WriteString(token_once)
	if err != nil {
		return []byte{}, err
	}

	_, err = bytesToSign.WriteString(action)
	if err != nil {
		return []byte{}, err
	}

	_, err = bytesToSign.WriteString(lease_uuid)
	if err != nil {
		return []byte{}, err
	}

	_, err = bytesToSign.WriteString(instance_id)
	if err != nil {
		return []byte{}, err
	}

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := crypto.SHA256.New()
	pssh.Write(bytesToSign.Bytes())
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, s.rsa.privateKey, crypto.SHA256, hashed, &opts)

	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

func (s *Service) verifySignature(c *gin.Context) error {

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

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := crypto.SHA256.New()

	pssh.Write(bytesToVerify.Bytes())
	hashed := pssh.Sum(nil)

	//Verify Signature
	return rsa.VerifyPSS(s.rsa.publicKey, crypto.SHA256, hashed, signature, &opts)
}

func (s *Service) generateSignedEmailActionURL(action, lease_uuid, instance_id, token_once string) (string, error) {
	signature, err := s.sign(lease_uuid, instance_id, action, token_once)
	if err != nil {
		return "", fmt.Errorf("error while signing")
	}
	signedURL := fmt.Sprintf("%s://%s%s/email_action/leases/%s/%s/%s?t=%s&s=%s",
		ZCDefaultScheme,
		ZCDefaultHostName,
		ZCDefaultPort,
		lease_uuid,
		instance_id,
		action,
		token_once,
		base64.URLEncoding.EncodeToString(signature),
	)
	return signedURL, nil
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

func (s *Service) AllowAWSIDTOSUbscribeToSQS(AWSID string) error {

	if AWSID == "" {
		return fmt.Errorf("AWSID is empty")
	}

	params := &sqs.AddPermissionInput{
		AWSAccountIds: []*string{ // Required
			aws.String(AWSID), // Required
			// More values...
		},
		Actions: []*string{ // Required
			aws.String("SendMessage"),
		},
		Label:    aws.String(AWSID + " SQS policy"), // Required
		QueueUrl: aws.String(SQSQueueURL()),         // Required
	}
	// TODO: add this condition:
	// ArnEquals	aws:SourceArn	arn:aws:sns:*:*:ZeroCloudTopic

	_, err := s.AWS.SQS.AddPermission(params)

	if err != nil {
		return err
	}

	return nil
}

var policy string = `
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
          "aws:SourceArn": "arn:aws:sns:*:*:ZeroCloudTopic"
        }
      }
    }
  ]
}
`

func (s *Service) SetQueueAttributes() error {
	resp, err := s.AWS.SQS.SetQueueAttributes(&sqs.SetQueueAttributesInput{
		Attributes: map[string]*string{
			"Policy": aws.String(policy),
		},
		QueueUrl: aws.String(SQSQueueURL()),
	})
	fmt.Println(resp)

	return err
}

func (s *Service) SubscribeSQS() error {

	svc := sns.New(s.AWS.Session)

	params := &sns.SubscribeInput{
		Protocol: aws.String("sqs"),                                               // Required
		TopicArn: aws.String("arn:aws:sns:us-east-1:859795398601:ZeroCloudTopic"), // Required
		Endpoint: aws.String("arn:aws:sqs:us-east-1:665102389639:ZeroCloudQueue"),
	}
	resp, err := svc.Subscribe(params)

	fmt.Println(resp)

	return err
}
