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

	ZCDefaultMaxLeasesPerOwner = 2

	// TODO: move these config values to config.yml
	maxWorkers   = 10
	maxQueueSize = 1000
)

var (
	ZCMailerFromAddress string

	ZCDefaultLeaseDuration                time.Duration
	ZCDefaultLeaseApprovalTimeoutDuration time.Duration
	ZCDefaultForewarningBeforeExpiry      time.Duration

	ZCDefaultScheme   string // http, or https
	ZCDefaultHostName string // e.g. zerocloud.co
	ZCDefaultPort     string
)

type Service struct {
	NewLeaseQueue        *simpleQueue.Queue
	TerminatorQueue      *simpleQueue.Queue
	LeaseTerminatedQueue *simpleQueue.Queue
	ExtenderQueue        *simpleQueue.Queue
	NotifierQueue        *simpleQueue.Queue

	// TODO: move EC2 into AWS ???
	EC2    Ec2ServiceFactory
	DB     *gorm.DB
	Mailer mailgun.Mailgun
	AWS    struct {
		Session *session.Session
		SQS     sqsiface.SQSAPI
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
	viperIsSet("ForeignRoleName")
	viperIsSet("AWS_ACCESS_KEY_ID")
	viperIsSet("AWS_SECRET_ACCESS_KEY")
	viperIsSet("ZCMailerDomain")
	viperIsSet("ZCMailerAPIKey")
	viperIsSet("UseMockAWS")
	viperIsSet("ZCMailerPublicAPIKey")
	viperIsSet("AWS_REGION")
	viperIsSet("AWS_ACCOUNT_ID")
	viperIsSet("SQSQueueName")
	viperIsSet("demo")

	// Set default values for scheme, hostname, port
	viper.SetDefault("DefaultScheme", "http")
	viper.SetDefault("DefaultHostName", "0.0.0.0")
	viper.SetDefault("DefaultPort", ":8080")
	// parse
	ZCDefaultScheme = viper.GetString("DefaultScheme")
	ZCDefaultHostName = viper.GetString("DefaultHostName")
	ZCDefaultPort = viper.GetString("DefaultPort")

	// Set default values for durations
	viper.SetDefault("DefaultLeaseDuration", 3*(time.Hour*24))
	viper.SetDefault("DefaultLeaseApprovalTimeoutDuration", 1*time.Hour)
	viper.SetDefault("DefaultForewarningBeforeExpiry", 12*time.Hour)
	// parse durations
	ZCDefaultLeaseDuration = viper.GetDuration("DefaultLeaseDuration")
	ZCDefaultLeaseApprovalTimeoutDuration = viper.GetDuration("DefaultLeaseApprovalTimeoutDuration")
	ZCDefaultForewarningBeforeExpiry = viper.GetDuration("DefaultForewarningBeforeExpiry")

	// some tests
	if ZCDefaultForewarningBeforeExpiry >= ZCDefaultLeaseDuration {
		panic("ZCDefaultForewarningBeforeExpiry >= ZCDefaultLeaseDuration")
	}
	if ZCDefaultLeaseApprovalTimeoutDuration >= ZCDefaultLeaseDuration {
		panic("ZCDefaultLeaseApprovalTimeoutDuration >= ZCDefaultLeaseDuration")
	}

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
		&Region{},
		&Owner{},
		&Lease{},
	)
	service.DB.AutoMigrate(
		&Account{},
		&CloudAccount{},
		&Region{},
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
				Regions: []Region{
					Region{
						RegionName: demo["Region"],
					},
				},
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

	// setup mailer service
	ZCMailerDomain := viper.GetString("ZCMailerDomain")
	ZCMailerAPIKey := viper.GetString("ZCMailerAPIKey")
	ZCMailerPublicAPIKey := viper.GetString("ZCMailerPublicAPIKey")
	service.Mailer = mailgun.NewMailgun(ZCMailerDomain, ZCMailerAPIKey, ZCMailerPublicAPIKey)
	ZCMailerFromAddress = fmt.Sprintf("ZeroCloud Guardian <noreply@%v>", ZCMailerDomain)

	switch viper.GetBool("UseMockAWS") {
	case true:
		service.AWS.SQS = &MockSQS{}
	default:
		// setup aws session
		AWSCreds := credentials.NewStaticCredentials(viper.GetString("AWS_ACCESS_KEY_ID"), viper.GetString("AWS_SECRET_ACCESS_KEY"), "")
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

	go scheduleJob(service.EventInjestorJob, time.Duration(time.Second*5))
	go scheduleJob(service.AlerterJob, time.Duration(time.Second*30))
	go scheduleJob(service.SentencerJob, time.Duration(time.Second*30))

	router := gin.Default()

	router.GET("/email_action/leases/:lease_uuid/:instance_id/:action", service.EmailActionHandler)

	router.POST("/accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", service.AddOwnerHandler)

	router.GET("/accounts/:account_id/cloudaccounts/:cloudaccount_id/regions", service.ListRegionsHandler)
	router.PATCH("/accounts/:account_id/cloudaccounts/:cloudaccount_id/regions", service.ChangeRegionStatusHandler)

	router.Run(ZCDefaultPort)
}
