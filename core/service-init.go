package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gagliardetto/simpleQueue"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tleyden/cecil/awstools"
	"github.com/tleyden/cecil/config"
	"github.com/tleyden/cecil/eventrecord"
	"github.com/tleyden/cecil/mailers"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/queues"
	"github.com/tleyden/cecil/tools"
)

// LoadConfig loads the configuration into Service.
func (service *Service) LoadConfig(configFilepath string) {

	var err error

	viper.SetConfigFile(configFilepath) // config file path
	viper.AutomaticEnv()
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(err)
	}

	service.AWS = awstools.AWSRes{}

	service.AWS.Config.AWS_REGION, err = tools.ViperMustGetString("AWS_REGION")
	if err != nil {
		panic(err)
	}

	service.AWS.Config.AWS_ACCOUNT_ID, err = tools.ViperMustGetString("AWS_ACCOUNT_ID")
	if err != nil {
		panic(err)
	}

	service.AWS.Config.AWS_ACCESS_KEY_ID, err = tools.ViperMustGetString("AWS_ACCESS_KEY_ID")
	if err != nil {
		panic(err)
	}

	service.AWS.Config.AWS_SECRET_ACCESS_KEY, err = tools.ViperMustGetString("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		panic(err)
	}

	service.AWS.Config.SNSTopicName, err = tools.ViperMustGetString("SNSTopicName")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.SQSQueueName, err = tools.ViperMustGetString("SQSQueueName")
	if err != nil {
		panic(err)
	}
	service.AWS.Config.ForeignIAMRoleName, err = tools.ViperMustGetString("ForeignIAMRoleName")
	if err != nil {
		panic(err)
	}

	service.config = &config.Config{}

	// Set default values for scheme, hostname, port
	viper.SetDefault("Scheme", "http") // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.config.Server.Scheme, err = tools.ViperMustGetString("ServerScheme")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("HostName", "0.0.0.0") // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.config.Server.HostName, err = tools.ViperMustGetString("ServerHostName")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("Port", ":8080") // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.config.Server.Port, err = tools.ViperMustGetString("ServerPort")
	if err != nil {
		panic(err)
	}

	service.config.DefaultMailer.Domain, err = tools.ViperMustGetString("MailerDomain")
	if err != nil {
		panic(err)
	}
	service.config.DefaultMailer.APIKey, err = tools.ViperMustGetString("MailerAPIKey")
	if err != nil {
		panic(err)
	}
	service.config.DefaultMailer.PublicAPIKey, err = tools.ViperMustGetString("MailerPublicAPIKey")
	if err != nil {
		panic(err)
	}

	service.defaultMailer = &mailers.MailerInstance{}
	service.defaultMailer.FromAddress = fmt.Sprintf("Cecil <noreply@%v>", service.config.DefaultMailer.Domain)

	// Set default values for durations
	viper.SetDefault("LeaseDuration", 3*(time.Hour*24)) // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.config.Lease.Duration, err = tools.ViperMustGetDuration("LeaseDuration")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseApprovalTimeoutDuration", 24*time.Hour) // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.config.Lease.ApprovalTimeoutDuration, err = tools.ViperMustGetDuration("LeaseApprovalTimeoutDuration")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseFirstWarningBeforeExpiry", 24*time.Hour) // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.config.Lease.FirstWarningBeforeExpiry, err = tools.ViperMustGetDuration("LeaseFirstWarningBeforeExpiry")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseSecondWarningBeforeExpiry", 3*time.Hour) // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.config.Lease.SecondWarningBeforeExpiry, err = tools.ViperMustGetDuration("LeaseSecondWarningBeforeExpiry")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("LeaseMaxPerOwner", 10) // this is the default value if no value is set on config.yml or environment; default is overrident by config.yml; config.yml value is ovverriden by environment value.
	service.config.Lease.MaxPerOwner, err = tools.ViperMustGetInt("LeaseMaxPerOwner")
	if err != nil {
		panic(err)
	}

	// some coherency tests
	if service.config.Lease.FirstWarningBeforeExpiry >= service.config.Lease.Duration {
		panic("service.config.Lease.FirstWarningBeforeExpiry >= service.config.Lease.Duration")
	}
	if service.config.Lease.SecondWarningBeforeExpiry >= service.config.Lease.FirstWarningBeforeExpiry {
		panic("service.config.Lease.SecondWarningBeforeExpiry >= service.config.Lease.FirstWarningBeforeExpiry")
	}
	if service.config.Lease.FirstWarningBeforeExpiry <= service.config.Lease.SecondWarningBeforeExpiry {
		panic("service.config.Lease.FirstWarningBeforeExpiry <= service.config.Lease.SecondWarningBeforeExpiry")
	}

	if service.config.Lease.ApprovalTimeoutDuration >= service.config.Lease.Duration {
		panic("service.config.Lease.ApprovalTimeoutDuration >= service.config.Lease.Duration")
	}

}

// GenerateRSAKeys generates and sets into Service the RSA keys.
func (service *Service) GenerateRSAKeys() {

	var err error
	var privateKey *rsa.PrivateKey

	privateKeyString, err := tools.ViperMustGetString("CECIL_RSA_PRIVATE")

	if err == nil {
		// load private key
		privateKey, err = jwtgo.ParseRSAPrivateKeyFromPEM([]byte(privateKeyString))
		if err != nil {
			panic(fmt.Errorf("jwt: failed to parse private key: %s", err))
		}
	} else {
		// generate Private Key
		if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
			panic(err)
		}

		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})
		fmt.Printf("\nCECIL_RSA_PRIVATE\n%v", string(pemBytes))
	}

	privateKey.Precompute()

	// validate Private Key
	if err = privateKey.Validate(); err != nil {
		panic(err)
	}

	service.rsa.privateKey = privateKey
	service.rsa.publicKey = &privateKey.PublicKey

}

type ErrorCallbackFunc func(error)

func createErrorCallback(errorSource string) ErrorCallbackFunc {
	return func(err error) {
		Logger.Error(errorSource, "err", err)
	}
}

// SetupQueues creates and initializes queues into Service.
func (service *Service) SetupQueues() {

	service.queues = &queues.QueuesGroup{}

	NewInstanceQueue := simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NewInstanceQueueConsumer).
		SetErrorCallback(createErrorCallback("service.NewInstanceQueueConsumer"))
	service.queues.SetNewInstanceQueue(NewInstanceQueue)
	service.queues.NewInstanceQueue().Start()

	TerminatorQueue := simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.TerminatorQueueConsumer).
		SetErrorCallback(createErrorCallback("service.TerminatorQueueConsumer"))
	service.queues.SetTerminatorQueue(TerminatorQueue)
	service.queues.TerminatorQueue().Start()

	InstanceTerminatedQueue := simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.InstanceTerminatedQueueConsumer).
		SetErrorCallback(createErrorCallback("service.InstanceTerminatedQueueConsumer"))
	service.queues.SetInstanceTerminatedQueue(InstanceTerminatedQueue)
	service.queues.InstanceTerminatedQueue().Start()

	ExtenderQueue := simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.ExtenderQueueConsumer).
		SetErrorCallback(createErrorCallback("service.ExtenderQueueConsumer"))
	service.queues.SetExtenderQueue(ExtenderQueue)
	service.queues.ExtenderQueue().Start()

	NotifierQueue := simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NotifierQueueConsumer).
		SetErrorCallback(createErrorCallback("service.NotifierQueueConsumer"))
	service.queues.SetNotifierQueue(NotifierQueue)
	service.queues.NotifierQueue().Start()

}

// SetupDB setups the DB.
func (service *Service) SetupDB(dbname string) {

	Logger.Info("Setup DB", "dbname", dbname)

	db, err := gorm.Open("sqlite3", dbname)
	if err != nil {
		panic(err)
	}
	service.DB = db
	dbService := models.NewDBService(db)
	service.DBService = dbService

	if DropAllTables {
		Logger.Warn("Dropping all DB tables")
		service.DB.DropTableIfExists(
			&models.Account{},
			&models.Cloudaccount{},
			&models.Owner{},
			&models.Lease{},
			&models.Instance{},
			&models.SlackConfig{},
			&models.MailerConfig{},
		)
	}

	service.DB.AutoMigrate(
		&models.Account{},
		&models.Cloudaccount{},
		&models.Owner{},
		&models.Lease{},
		&models.Instance{},
		&models.SlackConfig{},
		&models.MailerConfig{},
	)

}

// SetupEventRecording setups event recording
func (service *Service) SetupEventRecording(persistToDisk bool, storageDir string) {

	eventRecord, err := eventrecord.NewMossEventRecord(persistToDisk, storageDir)
	if err != nil {
		panic(fmt.Sprintf("Error setting up event recording: %v", err))
	}
	service.EventRecord = eventRecord
	Logger.Info("Setup event recording", "persisted", persistToDisk, "dir", storageDir)
}

// Stop stops the service.
func (service *Service) Stop(shouldCloseDb bool) {

	Logger.Info("Service Stop", "service", service)

	// Stop queues
	service.queues.NewInstanceQueue().Stop()
	service.queues.TerminatorQueue().Stop()
	service.queues.InstanceTerminatedQueue().Stop()
	service.queues.ExtenderQueue().Stop()
	service.queues.NotifierQueue().Stop()

	// close EventRecord
	if service.EventRecord != nil {
		if err := service.EventRecord.Close(); err != nil {
			Logger.Warn("Error closing eventRecord: %v", err)
		}
	}

	// terminate all slackBotInstances
	service.SlackBotService.StopAll()

	// terminate all mailerInstances
	service.CustomMailerService.StopAll()

	// Close DB
	//
	// Disabled when running tests, since it's causing "sql: database is closed" errors
	// even if different .db files are used in each test
	if shouldCloseDb {
		service.DB.Close()
	}
}
