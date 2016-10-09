package core

import (
	"time"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"

	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	_ "github.com/mattn/go-sqlite3"

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

func (service *Service) SetupAndRun() *Service {
	// Initialize global logger
	logger = log15.New()

	// Setup
	service.LoadConfig("config.yml")
	service.GenerateRSAKeys()
	service.SetupQueues()
	service.SetupDB()

	// @@@@@@@@@@@@@@@ Add Fake Account / Admin  @@@@@@@@@@@@@@@

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

	// Setup mailer client
	service.Mailer.Client = mailgun.NewMailgun(
		service.Mailer.Domain,
		service.Mailer.APIKey,
		service.Mailer.PublicAPIKey,
	)

	// Setup aws session
	AWSCreds := credentials.NewStaticCredentials(
		service.AWS.Config.AWS_ACCESS_KEY_ID,
		service.AWS.Config.AWS_SECRET_ACCESS_KEY,
		"",
	)
	AWSConfig := &aws.Config{
		Credentials: AWSCreds,
	}
	service.AWS.Session = session.New(AWSConfig)

	// Setup sqs
	service.AWS.SQS = sqs.New(service.AWS.Session)

	// Setup EC2
	service.EC2 = DefaultEc2ServiceFactory

	// @@@@@@@@@@@@@@@ Schedule Periodic Jobs @@@@@@@@@@@@@@@

	schedulePeriodicJob(service.EventInjestorJob, time.Duration(time.Second*5))
	schedulePeriodicJob(service.AlerterJob, time.Duration(time.Second*30))
	schedulePeriodicJob(service.SentencerJob, time.Duration(time.Second*30))

	// @@@@@@@@@@@@@@@ Update external services @@@@@@@@@@@@@@@

	// run this because the demo account has been added
	if err := service.RegenerateSQSPermissions(); err != nil {
		panic(err)
	}

	return service
}

func (service *Service) RunHTTPServer() error {

	router := gin.Default()

	router.GET("/email_action/leases/:lease_uuid/:instance_id/:action", service.EmailActionHandler)

	router.POST("/accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", service.AddOwnerHandler)

	router.POST("/accounts", service.CreateAccount)
	router.POST("/accounts/:account_id/api_token", service.ValidateAccount)

	return router.Run(service.Config.Server.Port)
}
