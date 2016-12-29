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

	service.Config.DefaultMailer.Domain, err = viperMustGetString("MailerDomain")
	if err != nil {
		panic(err)
	}
	service.Config.DefaultMailer.APIKey, err = viperMustGetString("MailerAPIKey")
	if err != nil {
		panic(err)
	}
	service.Config.DefaultMailer.PublicAPIKey, err = viperMustGetString("MailerPublicAPIKey")
	if err != nil {
		panic(err)
	}
	service.DefaultMailer.FromAddress = fmt.Sprintf("Cecil <noreply@%v>", service.Config.DefaultMailer.Domain)

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

// GenerateRSAKeys generates and sets into Service the RSA keys.
func (service *Service) GenerateRSAKeys() {

	var err error
	var privateKey *rsa.PrivateKey

	privateKeyString, err := viperMustGetString("CECIL_RSA_PRIVATE")

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

	service.NewLeaseQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NewLeaseQueueConsumer).
		SetErrorCallback(createErrorCallback("service.NewLeaseQueueConsumer"))
	service.NewLeaseQueue.Start()

	service.TerminatorQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.TerminatorQueueConsumer).
		SetErrorCallback(createErrorCallback("service.TerminatorQueueConsumer"))
	service.TerminatorQueue.Start()

	service.LeaseTerminatedQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.LeaseTerminatedQueueConsumer).
		SetErrorCallback(createErrorCallback("service.LeaseTerminatedQueueConsumer"))
	service.LeaseTerminatedQueue.Start()

	service.ExtenderQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.ExtenderQueueConsumer).
		SetErrorCallback(createErrorCallback("service.ExtenderQueueConsumer"))
	service.ExtenderQueue.Start()

	service.NotifierQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NotifierQueueConsumer).
		SetErrorCallback(createErrorCallback("service.NotifierQueueConsumer"))
	service.NotifierQueue.Start()

}

// SetupDB setups the DB.
func (service *Service) SetupDB(dbname string) {

	Logger.Info("Setup DB", "dbname", dbname)

	db, err := gorm.Open("sqlite3", dbname)
	if err != nil {
		panic(err)
	}
	service.DB = db

	if DropAllTables {
		Logger.Warn("Dropping all DB tables")
		service.DB.DropTableIfExists(
			&Account{},
			&CloudAccount{},
			&Owner{},
			&Lease{},
			&SlackConfig{},
			&MailerConfig{},
		)
	}

	service.DB.AutoMigrate(
		&Account{},
		&CloudAccount{},
		&Owner{},
		&Lease{},
		&SlackConfig{},
		&MailerConfig{},
	)

}

func (service *Service) SetupEventRecording(persistToDisk bool, storageDir string) {

	eventRecord, err := NewMossEventRecord(persistToDisk, storageDir)
	if err != nil {
		panic(fmt.Sprintf("Error setting up event recording: %v", err))
	}
	service.eventRecord = eventRecord
	Logger.Info("Setup event recording", "persisted", persistToDisk, "dir", storageDir)

}

// Stop stops the service.
func (service *Service) Stop(shouldCloseDb bool) {

	Logger.Info("Service Stop", "service", service)

	// Stop queues
	service.NewLeaseQueue.Stop()
	service.TerminatorQueue.Stop()
	service.LeaseTerminatedQueue.Stop()
	service.ExtenderQueue.Stop()
	service.NotifierQueue.Stop()
	if service.eventRecord != nil {
		if err := service.eventRecord.Close(); err != nil {
			Logger.Warn("Error closing eventRecord: %v", err)
		}
	}

	// terminate all slackInstances
	for accountID := range service.slackInstances {
		err := service.TerminateSlackInstance(accountID)
		if err != nil {
			Logger.Warn("Error terminating Slack instance: %v", err)
		}
	}

	// terminate all mailerInstances
	for accountID := range service.mailerInstances {
		err := service.TerminateMailerInstance(accountID)
		if err != nil {
			Logger.Warn("Error terminating mailer instance: %v", err)
		}
	}

	// Close DB
	//
	// Disabled when running tests, since it's causing "sql: database is closed" errors
	// even if different .db files are used in each test
	if shouldCloseDb {
		service.DB.Close()
	}
}
