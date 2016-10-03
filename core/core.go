package core

import (
	"fmt"
	"time"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"

	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	TerminatorActionTerminate = "terminate"
	TerminatorActionShutdown  = "shutdown"

	// TODO: move these config values to config.yml
	maxWorkers   = 10
	maxQueueSize = 1000
)

var logger log15.Logger

func Run() {
	// initialize global logger
	logger = log15.New()

	// @@@@@@@@@@@@@@@ Load config files @@@@@@@@@@@@@@@

	viper.SetConfigFile("config.yml") // config file path
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(err)
	}

	// create a service
	service := NewService()
	service.SetupQueues()
	defer service.Stop()

	// @@@@@@@@@@@@@@@ Set defaults, parse config variables @@@@@@@@@@@@@@@

	service.AWS.Config.UseMockAWS, err = viperMustGetBool("UseMockAWS")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.AWS_REGION, err = viperMustGetString("AWS_REGION")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.AWS_ACCOUNT_ID, err = viperMustGetString("AWS_ACCOUNT_ID")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.AWS_ACCESS_KEY_ID, err = viperMustGetString("AWS_ACCESS_KEY_ID")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.AWS_SECRET_ACCESS_KEY, err = viperMustGetString("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		panic(err)
	}

	service.AWS.Config.SNSTopicName, err = viperMustGetString("SNSTopicName")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.SQSQueueName, err = viperMustGetString("SQSQueueName")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.ForeignIAMRoleName, err = viperMustGetString("ForeignIAMRoleName")
	if err != nil {
		panic(err)
	}

	// Set default values for scheme, hostname, port
	viper.SetDefault("Scheme", "http")
	service.Config.Server.Scheme, err = viperMustGetString("ServerScheme")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("HostName", "0.0.0.0")
	service.Config.Server.HostName, err = viperMustGetString("ServerHostName")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("Port", ":8080")
	service.Config.Server.Port, err = viperMustGetString("ServerPort")
	if err != nil {
		panic(err)
	}

	service.Mailer.Domain, err = viperMustGetString("MailerDomain")
	if err != nil {
		panic(err)
	}
	service.Mailer.APIKey, err = viperMustGetString("MailerAPIKey")
	if err != nil {
		panic(err)
	}
	service.Mailer.PublicAPIKey, err = viperMustGetString("MailerPublicAPIKey")
	if err != nil {
		panic(err)
	}
	service.Mailer.FromAddress = fmt.Sprintf("ZeroCloud Guardian <noreply@%v>", service.Mailer.Domain)

	// Set default values for durations
	viper.SetDefault("LeaseDuration", 3*(time.Hour*24))
	service.Config.Lease.Duration, err = viperMustGetDuration("LeaseDuration")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseApprovalTimeoutDuration", 1*time.Hour)
	service.Config.Lease.ApprovalTimeoutDuration, err = viperMustGetDuration("LeaseApprovalTimeoutDuration")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("ForewarningBeforeExpiry", 12*time.Hour)
	service.Config.Lease.ForewarningBeforeExpiry, err = viperMustGetDuration("LeaseForewarningBeforeExpiry")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseMaxPerOwner", 2)
	service.Config.Lease.MaxPerOwner, err = viperMustGetInt("LeaseMaxPerOwner")
	if err != nil {
		panic(err)
	}

	// some coherency tests
	if service.Config.Lease.ForewarningBeforeExpiry >= service.Config.Lease.Duration {
		panic("service.Config.Lease.ForewarningBeforeExpiry >= service.Config.Lease.Duration")
	}
	if service.Config.Lease.ApprovalTimeoutDuration >= service.Config.Lease.Duration {
		panic("service.Config.Lease.ApprovalTimeoutDuration >= service.Config.Lease.Duration")
	}

	// @@@@@@@@@@@@@@@ Setup DB @@@@@@@@@@@@@@@

	db, err := gorm.Open("sqlite3", "zerocloud.db")
	if err != nil {
		panic(err)
	}
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
	service.DB = db

	defer service.DB.Close()

	service.DB.DropTableIfExists(
		&Account{},
		&CloudAccount{},
		&Owner{},
		&Lease{},
	)
	service.DB.AutoMigrate(
		&Account{},
		&CloudAccount{},
		&Owner{},
		&Lease{},
	)

	// setup mailer client
	service.Mailer.Client = mailgun.NewMailgun(
		service.Mailer.Domain,
		service.Mailer.APIKey,
		service.Mailer.PublicAPIKey,
	)

	// <EDIT-HERE>
	demo, err := viperMustGetStringMapString("demo")
	if err != nil {
		panic("no demo account set")
	}

	logger.Info("adding demo account",
		"email", demo["Email"],
	)

	firstUser := Account{
		Email: demo["Email"],
		CloudAccounts: []CloudAccount{
			CloudAccount{
				Provider:   demo["Provider"],
				AWSID:      demo["AWSID"],
				ExternalID: demo["ExternalID"],
			},
		},
	}
	service.DB.Create(&firstUser)

	firstOwner := Owner{
		Email:          demo["Email"],
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&firstOwner)

	secondaryOwner := Owner{
		Email:          demo["SecondaryEmail"],
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&secondaryOwner)
	// </EDIT-HERE>

	// @@@@@@@@@@@@@@@ Setup external services @@@@@@@@@@@@@@@

	switch service.AWS.Config.UseMockAWS {
	case true:
		service.AWS.SQS = &MockSQS{}
	default:
		// setup aws session
		AWSCreds := credentials.NewStaticCredentials(
			service.AWS.Config.AWS_ACCESS_KEY_ID,
			service.AWS.Config.AWS_SECRET_ACCESS_KEY,
			"",
		)
		AWSConfig := &aws.Config{
			Credentials: AWSCreds,
		}
		service.AWS.Session = session.New(AWSConfig)

		// setup sqs
		service.AWS.SQS = sqs.New(service.AWS.Session)
	}

	service.EC2 = DefaultEc2ServiceFactory

	// create rsa keys
	service.rsa.privateKey, service.rsa.publicKey, err = generateRSAKeys()
	if err != nil {
		panic(err)
	}

	scheduleJob(service.EventInjestorJob, time.Duration(time.Second*5))
	scheduleJob(service.AlerterJob, time.Duration(time.Second*30))
	scheduleJob(service.SentencerJob, time.Duration(time.Second*30))

	// run this because the demo account has been added
	if err := service.RegenerateSQSPermissions(); err != nil {
		panic(err)
	}

	router := gin.Default()

	router.GET("/email_action/leases/:lease_uuid/:instance_id/:action", service.EmailActionHandler)

	router.POST("/accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", service.AddOwnerHandler)

	router.Run(service.Config.Server.Port)
}
