package core

import (
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

	"github.com/tleyden/zerocloud/mocks/aws"
)

// declare task structs
// setup queues' consumer functions
// setup queues

// setup jobs
// setup services that will be used by multiple workers at the same time

// db
// sqs
// ec2
// ses

// run everything

// EventInjestorJob
// AlerterJob
// SentencerJob

// NewLeasesQueue
// TerminatorQueue
// LeaseTerminatedQueue
// RenewerQueue
// NotifiesQueue

const (
	TerminatorActionTerminate = "terminate"
	TerminatorActionShutdown  = "shutdown"

	ZCMaxLeasesPerOwner                   = 2
	ZCDefaultLeaseDuration                = time.Minute * 1
	ZCDefaultLeaseApprovalTimeoutDuration = time.Minute * 1

	// TODO: move these config values to config.yml
	maxWorkers   = 10
	maxQueueSize = 1000
)

var (
	ZCMailerFromAddress string
)

type Service struct {
	NewLeaseQueue        *simpleQueue.Queue
	TerminatorQueue      *simpleQueue.Queue
	LeaseTerminatedQueue *simpleQueue.Queue
	RenewerQueue         *simpleQueue.Queue
	NotifierQueue        *simpleQueue.Queue

	EC2    Ec2ServiceFactory
	DB     *gorm.DB
	Mailer mailgun.Mailgun
	AWS    struct {
		Session *session.Session
		SQS     sqsiface.SQSAPI
	}
}

var logger log15.Logger

func viperIsSet(key string) bool {
	if !viper.IsSet(key) {
		logger.Crit("Config parameter not set",
			key, viper.Get(key),
		)
		return false
	}
	return true
}

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

	// for more options, see https://godoc.org/github.com/spf13/viper

	// viper.SetDefault("LayoutDir", "layouts")
	// viper.GetString("logfile")
	// viper.GetBool("verbose")

	var service Service = Service{}

	// @@@@@@@@@@@@@@@ Setup queues @@@@@@@@@@@@@@@

	service.NewLeaseQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NewLeaseQueueConsumer)
	service.NewLeaseQueue.Start()
	defer service.NewLeaseQueue.Stop()

	service.TerminatorQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.TerminatorQueueConsumer)
	service.TerminatorQueue.Start()
	defer service.TerminatorQueue.Stop()

	service.LeaseTerminatedQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.LeaseTerminatedQueueConsumer)
	service.LeaseTerminatedQueue.Start()
	defer service.LeaseTerminatedQueue.Stop()

	service.RenewerQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.RenewerQueueConsumer)
	service.RenewerQueue.Start()
	defer service.RenewerQueue.Stop()

	service.NotifierQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NotifierQueueConsumer)
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
		&Lease{},
		&Region{},
		&Owner{},
	)
	service.DB.AutoMigrate(
		&Account{},
		&CloudAccount{},
		&Lease{},
		&Region{},
		&Owner{},
	)

	// <EDIT-HERE>
	if viper.IsSet("demo") {
		demo := viper.GetStringMapString("demo")

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
							Region: demo["Region"],
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
	} else {
		panic("no account")
	}
	// </EDIT-HERE>

	/*
		// <debug>
		expiredLease := Lease{
			ExpiresAt: time.Now().Add(-time.Minute),
			Region:    "some random region",
			OwnerID:   firstOwner.ID,
		}
		service.DB.Create(&expiredLease)
		// </debug>
	*/
	// @@@@@@@@@@@@@@@ Setup external services @@@@@@@@@@@@@@@

	// setup mailer service
	ZCMailerDomain := viper.GetString("ZCMailerDomain")
	ZCMailerAPIKey := viper.GetString("ZCMailerAPIKey")
	ZCMailerPublicAPIKey := viper.GetString("ZCMailerPublicAPIKey")
	service.Mailer = mailgun.NewMailgun(ZCMailerDomain, ZCMailerAPIKey, ZCMailerPublicAPIKey)
	ZCMailerFromAddress = fmt.Sprintf("ZeroCloud Guardian <postmaster@%v>", ZCMailerDomain)

	switch viper.GetBool("UseMockAWS") {
	case true:
		service.AWS.SQS = &mockaws.MockSQS{}
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

	go runForever(service.EventInjestorJob, time.Duration(time.Second*5))
	go runForever(service.AlerterJob, time.Duration(time.Second*60))
	go runForever(service.SentencerJob, time.Duration(time.Second*60))

	r := gin.Default()

	r.GET("/leases/:leaseID/terminate", service.TerminatorHandle)
	r.GET("/leases/:leaseID/renew", service.RenewerHandle)
	r.Run() // listen and server on 0.0.0.0:8080
}