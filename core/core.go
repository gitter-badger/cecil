package core

import (
	"fmt"
	"time"

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
	"gopkg.in/mailgun/mailgun-go.v1"
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

	ZCMaxLeasesPerOwner      = 10
	ZCDefaultLeaseExpiration = time.Minute * 1
	ZCDefaultTruceDuration   = time.Minute * 1

	// TODO: move these config values to config.yml
	maxWorkers   = 10
	maxQueueSize = 1000
)

var (
	ZCMailerFromAddress string
)

type Service struct {
	counter int64

	NewLeaseQueue        *simpleQueue.Queue
	TerminatorQueue      *simpleQueue.Queue
	LeaseTerminatedQueue *simpleQueue.Queue
	RenewerQueue         *simpleQueue.Queue
	NotifierQueue        *simpleQueue.Queue

	DB     *gorm.DB
	Mailer mailgun.Mailgun
	AWS    struct {
		Session *session.Session
		SQS     *sqs.SQS
	}
}

var logger log15.Logger

func Run() {
	// Such and other options (db address, etc.) could be stored in:
	// · environment variables
	// · flags
	// · config file (read with viper)

	logger = log15.New()

	viper.SetConfigFile("config.yml") // name of config file (without extension)
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(err)
	}
	// for more options, see https://godoc.org/github.com/spf13/viper

	// viper.SetDefault("LayoutDir", "layouts")
	// viper.GetString("logfile")
	// viper.GetBool("verbose")

	var service Service = Service{}
	service.counter = 0

	// @@@@@@@@@@@@@@@ Setup queues @@@@@@@@@@@@@@@

	service.NewLeaseQueue = simpleQueue.NewQueue()
	service.NewLeaseQueue.SetMaxSize(maxQueueSize)
	service.NewLeaseQueue.SetWorkers(maxWorkers)
	service.NewLeaseQueue.Consumer = service.NewLeaseQueueConsumer
	service.NewLeaseQueue.Start()
	defer service.NewLeaseQueue.Stop()

	service.TerminatorQueue = simpleQueue.NewQueue()
	service.TerminatorQueue.SetMaxSize(maxQueueSize)
	service.TerminatorQueue.SetWorkers(maxWorkers)
	service.TerminatorQueue.Consumer = service.TerminatorQueueConsumer
	service.TerminatorQueue.Start()
	defer service.TerminatorQueue.Stop()

	service.LeaseTerminatedQueue = simpleQueue.NewQueue()
	service.LeaseTerminatedQueue.SetMaxSize(maxQueueSize)
	service.LeaseTerminatedQueue.SetWorkers(maxWorkers)
	service.LeaseTerminatedQueue.Consumer = service.LeaseTerminatedQueueConsumer
	service.LeaseTerminatedQueue.Start()
	defer service.LeaseTerminatedQueue.Stop()

	service.RenewerQueue = simpleQueue.NewQueue()
	service.RenewerQueue.SetMaxSize(maxQueueSize)
	service.RenewerQueue.SetWorkers(maxWorkers)
	service.RenewerQueue.Consumer = service.RenewerQueueConsumer
	service.RenewerQueue.Start()
	defer service.RenewerQueue.Stop()

	service.NotifierQueue = simpleQueue.NewQueue()
	service.NotifierQueue.SetMaxSize(maxQueueSize)
	service.NotifierQueue.SetWorkers(maxWorkers)
	service.NotifierQueue.Consumer = service.NotifierQueueConsumer
	service.NotifierQueue.Start()
	defer service.NotifierQueue.Stop()

	/*
		How about:

		service.NotifierQueue = simpleQueue.NewQueue().SetMaxSize(maxQueueSize).SetWorkers(maxWorkers).SetConsumer(service.NotifierQueueConsumer)
		service.NotifierQueue.Start()
	*/

	// @@@@@@@@@@@@@@@ Setup DB @@@@@@@@@@@@@@@

	db, err := gorm.Open("sqlite3", "zerocloud.db")
	if err != nil {
		panic(err)
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

	firstUser := Account{
		Email: "slv.balsan@gmail.com",
		CloudAccounts: []CloudAccount{
			CloudAccount{
				Provider:   "aws",
				AWSID:      859795398601,
				ExternalID: "slavomir",
				Regions: []Region{
					Region{
						Region: "us-east-1",
					},
				},
			},
		},
	}
	service.DB.Create(&firstUser)

	firstOwner := Owner{
		Email:          "slv.balsan@gmail.com",
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&firstOwner)

	secondaryOwner := Owner{
		Email:          "slavomir.balsan@gmail.com",
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&secondaryOwner)

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

	// setup aws session
	AWSCreds := credentials.NewStaticCredentials(viper.GetString("AWS_ACCESS_KEY_ID"), viper.GetString("AWS_SECRET_ACCESS_KEY"), "")
	AWSConfig := &aws.Config{
		Credentials: AWSCreds,
	}
	service.AWS.Session = session.New(AWSConfig)

	// setup sqs
	service.AWS.SQS = sqs.New(service.AWS.Session)

	go runForever(service.EventInjestorJob, time.Duration(time.Second*5))
	go runForever(service.AlerterJob, time.Duration(time.Second*60))
	go runForever(service.SentencerJob, time.Duration(time.Second*60))

	r := gin.Default()

	r.GET("/leases/:leaseID/terminate", service.TerminatorHandle)
	r.GET("/leases/:leaseID/renew", service.RenewerHandle)
	r.Run() // listen and server on 0.0.0.0:8080
}
