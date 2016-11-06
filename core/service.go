package core

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/gagliardetto/simpleQueue"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
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
			HostName string // hostname for links back to REST API from emails, etc
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

func NewService() *Service {
	service := &Service{}
	return service
}

func (service *Service) GenerateRSAKeys() {

	var err error

	// create rsa keys
	service.rsa.privateKey, service.rsa.publicKey, err = generateRSAKeys()
	if err != nil {
		panic(err)
	}

}

func (service *Service) SetupQueues() {

	service.NewLeaseQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NewLeaseQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.NewLeaseQueueConsumer", "error", err)
		})
	service.NewLeaseQueue.Start()

	service.TerminatorQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.TerminatorQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.TerminatorQueueConsumer", "error", err)
		})
	service.TerminatorQueue.Start()

	service.LeaseTerminatedQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.LeaseTerminatedQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.LeaseTerminatedQueueConsumer", "error", err)
		})
	service.LeaseTerminatedQueue.Start()

	service.ExtenderQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.ExtenderQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.ExtenderQueueConsumer", "error", err)
		})
	service.ExtenderQueue.Start()

	service.NotifierQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NotifierQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.NotifierQueueConsumer", "error", err)
		})
	service.NotifierQueue.Start()

}

func (service *Service) LoadConfig(configFilepath string) {

	var err error

	viper.SetConfigFile(configFilepath) // config file path
	viper.AutomaticEnv()
	err = viper.ReadInConfig() // Find and read the config file
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
	viper.SetDefault("Scheme", "http") // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.Config.Server.Scheme, err = viperMustGetString("ServerScheme")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("HostName", "0.0.0.0") // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.Config.Server.HostName, err = viperMustGetString("ServerHostName")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("Port", ":8080") // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
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
	service.Mailer.FromAddress = fmt.Sprintf("Cecil <noreply@%v>", service.Mailer.Domain)

	// Set default values for durations
	viper.SetDefault("LeaseDuration", 3*(time.Hour*24)) // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.Config.Lease.Duration, err = viperMustGetDuration("LeaseDuration")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseApprovalTimeoutDuration", 1*time.Hour) // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.Config.Lease.ApprovalTimeoutDuration, err = viperMustGetDuration("LeaseApprovalTimeoutDuration")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("ForewarningBeforeExpiry", 12*time.Hour) // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.Config.Lease.ForewarningBeforeExpiry, err = viperMustGetDuration("LeaseForewarningBeforeExpiry")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseMaxPerOwner", 2) // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
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

}

func (service *Service) SetupDB(dbname string) {

	db, err := gorm.Open("sqlite3", dbname)
	if err != nil {
		panic(err)
	}
	service.DB = db

	if dropAllTables {
		service.DB.DropTableIfExists(
			&Account{},
			&CloudAccount{},
			&Owner{},
			&Lease{},
		)
	}

	service.DB.AutoMigrate(
		&Account{},
		&CloudAccount{},
		&Owner{},
		&Lease{},
	)

}

func (service *Service) Stop(shouldCloseDb bool) {

	logger.Info("Service Stop", "service", service)

	// Stop queues
	service.NewLeaseQueue.Stop()
	service.TerminatorQueue.Stop()
	service.LeaseTerminatedQueue.Stop()
	service.ExtenderQueue.Stop()
	service.NotifierQueue.Stop()

	// Close DB
	//
	// Disabled when running tests, since it's causing "sql: database is closed" errors
	// even if different .db files are used in each test
	if shouldCloseDb {
		service.DB.Close()
	}
}
