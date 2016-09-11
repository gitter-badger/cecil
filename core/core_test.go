package core

import (
	"testing"
	"time"

	"github.com/gagliardetto/simpleQueue"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tleyden/zerocloud/mocks/aws"
	"github.com/tleyden/zerocloud/mocks/mailgun"
)

func TestEndToEnd(t *testing.T) {

	logger = log15.New()

	viper.SetConfigFile("temporary/config.yml") // name of config file (without extension)
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(err)
	}

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
	firstUser := Account{
		Email: "traun.leyden@gmail.com",
		CloudAccounts: []CloudAccount{
			CloudAccount{
				Provider:   "aws",
				AWSID:      "788612350743",
				ExternalID: "bigdb_zerocloud",
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
		Email:          "traun.leyden@gmail.com",
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&firstOwner)

	secondaryOwner := Owner{
		Email:          "tleyden@yahoo.com",
		CloudAccountID: firstUser.CloudAccounts[0].ID,
	}
	service.DB.Create(&secondaryOwner)
	// </EDIT-HERE>

	// @@@@@@@@@@@@@@@ Setup external services @@@@@@@@@@@@@@@

	// setup mailer service
	service.Mailer = &mockmailgun.MockMailGun{}

	service.AWS.SQS = &mockaws.MockSQS{}

	go runForever(service.EventInjestorJob, time.Duration(time.Second*5))
	go runForever(service.AlerterJob, time.Duration(time.Second*60))
	go runForever(service.SentencerJob, time.Duration(time.Second*60))

	logger.Info("Sleeping while stuff happens")
	time.Sleep(100 * time.Second)

}
