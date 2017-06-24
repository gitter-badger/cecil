// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

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
	"github.com/tleyden/cecil/mailers"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/queues"
	"github.com/tleyden/cecil/tools"
	"io/ioutil"
)

const (
	maxQueueSize = 1000
	maxWorkers   = 100
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

	service.AWS = awstools.AWSServices{}

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

	service.config.ProductName, err = tools.ViperMustGetString("ProductName")
	if err != nil {
		panic(err)
	}

	service.defaultMailer = &mailers.MailerInstance{}
	service.defaultMailer.FromAddress = fmt.Sprintf("%s <noreply@%v>", service.config.ProductName, service.config.DefaultMailer.Domain)

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
		panic(
			fmt.Sprintf(
				"service.config.Lease.FirstWarningBeforeExpiry (%v) >= service.config.Lease.Duration (%v)",
				service.config.Lease.FirstWarningBeforeExpiry,
				service.config.Lease.Duration,
			),
		)
	}
	if service.config.Lease.SecondWarningBeforeExpiry >= service.config.Lease.FirstWarningBeforeExpiry {
		panic(
			fmt.Sprintf(
				"service.config.Lease.SecondWarningBeforeExpiry (%v) >= service.config.Lease.FirstWarningBeforeExpiry (%v)",
				service.config.Lease.SecondWarningBeforeExpiry,
				service.config.Lease.FirstWarningBeforeExpiry,
			),
		)
	}
	if service.config.Lease.FirstWarningBeforeExpiry <= service.config.Lease.SecondWarningBeforeExpiry {
		panic(
			fmt.Sprintf(
				"service.config.Lease.FirstWarningBeforeExpiry (%v) <= service.config.Lease.SecondWarningBeforeExpiry (%v)",
				service.config.Lease.FirstWarningBeforeExpiry,
				service.config.Lease.SecondWarningBeforeExpiry,
			),
		)
	}
	if service.config.Lease.ApprovalTimeoutDuration >= service.config.Lease.Duration {
		panic(
			fmt.Sprintf(
				"service.config.Lease.ApprovalTimeoutDuration (%v) >= service.config.Lease.Duration (%v)",
				service.config.Lease.ApprovalTimeoutDuration,
				service.config.Lease.Duration,
			),
		)
	}

}

// GenerateRSAKeys generates and sets into Service the RSA keys.
func (service *Service) GenerateRSAKeys() {

	var err error
	var privateKey *rsa.PrivateKey

	privateKeyFilename, err := tools.ViperMustGetString("CECIL_RSA_PRIVATE")

	if err == nil {
		// load private key from file
		fmt.Printf("\nLoading CECIL_RSA_PRIVATE from file: %v\n", privateKeyFilename)

		privateKeyRaw, err := ioutil.ReadFile(privateKeyFilename)
		if err != nil {
			panic(fmt.Errorf("jwt: failed to read private key from file: %s.  Err: %v", privateKeyFilename, err))
		}
		privateKey, err = jwtgo.ParseRSAPrivateKeyFromPEM([]byte(privateKeyRaw))
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
		fmt.Printf("\nGenerated CECIL_RSA_PRIVATE\n%v", string(pemBytes))
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
		SetWorkers(1).
		SetConsumer(service.NewInstanceQueueConsumer).
		SetErrorCallback(createErrorCallback("service.NewInstanceQueueConsumer"))
	service.queues.SetNewInstanceQueue(NewInstanceQueue)
	service.queues.NewInstanceQueue().Start()

	TerminatorQueue := simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(1).
		SetConsumer(service.TerminatorQueueConsumer).
		SetErrorCallback(createErrorCallback("service.TerminatorQueueConsumer"))
	service.queues.SetTerminatorQueue(TerminatorQueue)
	service.queues.TerminatorQueue().Start()

	InstanceTerminatedQueue := simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(1).
		SetConsumer(service.InstanceTerminatedQueueConsumer).
		SetErrorCallback(createErrorCallback("service.InstanceTerminatedQueueConsumer"))
	service.queues.SetInstanceTerminatedQueue(InstanceTerminatedQueue)
	service.queues.InstanceTerminatedQueue().Start()

	ExtenderQueue := simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(1).
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

// DBSettings contains the settings of the postgres DB
type DBSettings struct {
	host         string
	username     string
	password     string
	dbName       string
	sslmode      string
	port         int
	maxOpenConns int
	maxIdleConns int
	debug        bool
}

// ConnectPosgresqlDB loads the DB settings and connects to the postgres DB
func (service *Service) ConnectPosgresqlDB() (*gorm.DB, error) {
	Logger.Info("Connecting to postgres DB...",
		"host", viper.GetString("postgres.host"),
	)

	dbSettings, err := loadDBSettings()
	if err != nil {
		return nil, err
	}

	connectionString := fmt.Sprintf(
		"host=%v user=%v dbname=%v sslmode=%v password=%v",

		dbSettings.host,
		dbSettings.username,
		dbSettings.dbName,
		dbSettings.sslmode,
		dbSettings.password,
	)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(dbSettings.maxIdleConns)
	db.DB().SetMaxOpenConns(dbSettings.maxOpenConns)

	// TODO: fetch from config
	db.DB().SetConnMaxLifetime(time.Second * 120)

	db.LogMode(dbSettings.debug)

	Logger.Info("Connected to postgres DB",
		"host", viper.GetString("postgres.host"),
	)

	return db, nil
}

// loadDBSettings check for every setting to be set and returns a DBSettings struct
func loadDBSettings() (DBSettings, error) {
	var err error
	var dbSettings = DBSettings{}

	dbSettings.host, err = tools.ViperMustGetString("postgres.host")
	if err != nil {
		return DBSettings{}, err
	}

	dbSettings.port, err = tools.ViperMustGetInt("postgres.port")
	if err != nil {
		return DBSettings{}, err
	}

	dbSettings.sslmode, err = tools.ViperMustGetString("postgres.sslmode")
	if err != nil {
		return DBSettings{}, err
	}

	dbSettings.dbName, err = tools.ViperMustGetString("postgres.dbname")
	if err != nil {
		return DBSettings{}, err
	}

	dbSettings.password, err = tools.ViperMustGetString("postgres.password")
	if err != nil {
		return DBSettings{}, err
	}

	dbSettings.username, err = tools.ViperMustGetString("postgres.user")
	if err != nil {
		return DBSettings{}, err
	}

	dbSettings.maxIdleConns, err = tools.ViperMustGetInt("postgres.maxIdleConns")
	if err != nil {
		return DBSettings{}, err
	}

	dbSettings.maxOpenConns, err = tools.ViperMustGetInt("postgres.maxOpenConns")
	if err != nil {
		return DBSettings{}, err
	}

	dbSettings.debug, err = tools.ViperMustGetBool("postgres.debug")
	if err != nil {
		return DBSettings{}, err
	}

	return dbSettings, nil
}

// SetupDB setups the DB.
func (service *Service) SetupDB(dbname string) {

	switch DBType {
	case "sqlite", "":
		Logger.Info("Setup sqlite DB", "dbname", dbname)
		db, err := gorm.Open("sqlite3", dbname)
		if err != nil {
			panic(err)
		}
		service.DB = db
		dbService := models.NewDBService(db)
		service.DBService = dbService
	case "postgres":
		Logger.Info("Setup postgres DB")
		db, err := service.ConnectPosgresqlDB()
		if err != nil {
			panic(err)
		}
		service.DB = db
		dbService := models.NewDBService(db)
		service.DBService = dbService
	default:
		panic("Unknown DB type: " + DBType)
	}

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


// Stop stops the service.
func (service *Service) Stop(shouldCloseDb bool) {

	Logger.Info("Service Stop", "service", service)

	// Stop queues
	service.queues.NewInstanceQueue().Stop()
	service.queues.TerminatorQueue().Stop()
	service.queues.InstanceTerminatedQueue().Stop()
	service.queues.ExtenderQueue().Stop()
	service.queues.NotifierQueue().Stop()

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
