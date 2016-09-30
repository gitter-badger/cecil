package core

import (
	"crypto/rsa"
	"fmt"
	"time"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"

	"github.com/gagliardetto/simpleQueue"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

const (
	TerminatorActionTerminate = "terminate"
	TerminatorActionShutdown  = "shutdown"

	// TODO: move these config values to config.yml
	maxWorkers   = 10
	maxQueueSize = 1000
)

type Service struct {
	NewLeaseQueue        *simpleQueue.Queue
	TerminatorQueue      *simpleQueue.Queue
	LeaseTerminatedQueue *simpleQueue.Queue
	ExtenderQueue        *simpleQueue.Queue
	NotifierQueue        *simpleQueue.Queue

	Config struct {
		Server struct {
			Scheme   string // http, or https
			HostName string // e.g. zerocloud.co
			Port     string
		}
		Lease struct {
			Duration                time.Duration
			ApprovalTimeoutDuration time.Duration
			ForewarningBeforeExpiry time.Duration
			MaxPerOwner             int
		}
	}
	// TODO: move EC2 into AWS ???
	EC2    Ec2ServiceFactory
	DB     *gorm.DB
	Mailer struct {
		Client       mailgun.Mailgun
		Domain       string
		APIKey       string
		PublicAPIKey string
		FromAddress  string
	}
	AWS struct {
		Session *session.Session
		SQS     sqsiface.SQSAPI
		Config  struct {
			UseMockAWS            bool
			AWS_REGION            string
			AWS_ACCOUNT_ID        string
			AWS_ACCESS_KEY_ID     string
			AWS_SECRET_ACCESS_KEY string

			SNSTopicName       string
			SQSQueueName       string
			ForeignIAMRoleName string
		}
	}
	rsa struct {
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
	}
}

var logger log15.Logger

func Run() {
	// Such and other options (db address, etc.) could be stored in:
	// · environment variables
	// · flags
	// · config file (read with viper)

	logger = log15.New()

	viper.SetConfigFile("config.yml") // config file
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(err)
	}

	// @@@@@@@@@@@@@@@ Check whether these values have been set in the config @@@@@@@@@@@@@@@

	// TODO: set these variables as global, using viperMustGet*
	// viperIsSet("ForeignRoleName")
	// viperIsSet("AWS_ACCESS_KEY_ID")
	// viperIsSet("AWS_SECRET_ACCESS_KEY")
	// viperIsSet("ZCMailerDomain")
	// viperIsSet("ZCMailerAPIKey")
	// viperIsSet("ZCMailerPublicAPIKey")
	// viperIsSet("UseMockAWS")
	//viperIsSet("AWS_REGION")
	//viperIsSet("AWS_ACCOUNT_ID")
	// viperIsSet("SQSQueueName")
	viperIsSet("demo")

	var service Service = Service{}

	// @@@@@@@@@@@@@@@ Setup queues @@@@@@@@@@@@@@@

	service.NewLeaseQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NewLeaseQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.NewLeaseQueueConsumer error:", "error", err)
		})
	service.NewLeaseQueue.Start()
	defer service.NewLeaseQueue.Stop()

	service.TerminatorQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.TerminatorQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.TerminatorQueueConsumer error:", "error", err)
		})
	service.TerminatorQueue.Start()
	defer service.TerminatorQueue.Stop()

	service.LeaseTerminatedQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.LeaseTerminatedQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.LeaseTerminatedQueueConsumer error:", "error", err)
		})
	service.LeaseTerminatedQueue.Start()
	defer service.LeaseTerminatedQueue.Stop()

	service.ExtenderQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.ExtenderQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.ExtenderQueueConsumer error:", "error", err)
		})
	service.ExtenderQueue.Start()
	defer service.ExtenderQueue.Stop()

	service.NotifierQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NotifierQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.NotifierQueueConsumer error:", "error", err)
		})
	service.NotifierQueue.Start()
	defer service.NotifierQueue.Stop()

	// @@@@@@@@@@@@@@@ Parse config variables @@@@@@@@@@@@@@@

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

	// setup mailer client
	service.Mailer.Client = mailgun.NewMailgun(service.Mailer.Domain, service.Mailer.APIKey, service.Mailer.PublicAPIKey)

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

	// TODO: add permissions to SQS

	// </EDIT-HERE>

	// @@@@@@@@@@@@@@@ Setup external services @@@@@@@@@@@@@@@

	switch service.AWS.Config.UseMockAWS {
	case true:
		service.AWS.SQS = &MockSQS{}
	default:
		// setup aws session
		AWSCreds := credentials.NewStaticCredentials(service.AWS.Config.AWS_ACCESS_KEY_ID, service.AWS.Config.AWS_SECRET_ACCESS_KEY, "")
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

	// for each cloudAccount in the DB, allow the corresponding AWS account to send messages to the SQS queue
	if err := service.RegenerateSQSPermissions(); err != nil {
		panic(err)
	}

	router := gin.Default()

	router.GET("/email_action/leases/:lease_uuid/:instance_id/:action", service.EmailActionHandler)

	router.POST("/accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", service.AddOwnerHandler)

	router.Run(service.Config.Server.Port)
}
