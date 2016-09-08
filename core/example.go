package main

import (
	"fmt"
	"sync/atomic"

	"github.com/gagliardetto/simpleQueue"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	"gopkg.in/mailgun/mailgun-go.v1"
)

const (
	ZeroCloudGuardianSender = "ZeroCloud Guardian <guardian@zerocloud.site>"
)

// NewLeaseTask{} == Lease{}
type NewLeaseTask struct {
}

type TerminatorTask struct {
}

type LeaseTerminatedTask struct {
}

type RenewerTask struct {
}

type NotifierTask struct {
	From     string
	To       string
	Subject  string
	BodyHTML string
	BodyText string
}

type Service struct {
	counter int64

	NewLeaseQueue        *simpleQueue.Queue
	TerminatorQueue      *simpleQueue.Queue
	LeaseTerminatedQueue *simpleQueue.Queue
	RenewerQueue         *simpleQueue.Queue
	NotifierQueue        *simpleQueue.Queue

	DB     *gorm.DB
	Mailer mailgun.Mailgun
}

func (s *Service) NewLeaseQueueConsumer(t interface{}) error {

	return nil
}

func (s *Service) TerminatorQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(TerminatorTask)

	_ = task

	atomic.AddInt64(&s.counter, 1)
	fmt.Println(s.counter)
	return nil
}

func (s *Service) LeaseTerminatedQueueConsumer(t interface{}) error {

	return nil
}

func (s *Service) RenewerQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(RenewerTask)

	_ = task

	atomic.AddInt64(&s.counter, 1)
	fmt.Println(s.counter)

	return nil
}

func (s *Service) NotifierQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(NotifierTask)

	message := mailgun.NewMessage(
		task.From,
		task.Subject,
		task.BodyText,
		task.To,
	)

	message.SetTracking(true)
	//message.SetDeliveryTime(time.Now().Add(24 * time.Hour))
	message.SetHtml(task.BodyHTML)
	_, id, err := s.Mailer.Send(message)
	if err != nil {
		logger.Error("error while sending email", err)
		return err
	}
	_ = id

	return nil
}

var logger log15.Logger

func main() {
	// Such and other options (db address, etc.) could be stored in:
	// · environment variables
	// · flags
	// · config file (read with viper)

	logger = log15.New()

	var (
		maxWorkers   = 10
		maxQueueSize = 1000
		domain       = ""
		apiKey       = ""
		publicApiKey = ""
	)

	var service Service = Service{}
	service.counter = 0

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

	db, err := gorm.Open("sqlite3", "zerocloud.db")
	if err != nil {
		panic(err)
	}
	service.DB = db

	defer service.DB.Close()

	service.DB.DropTable(
		&Account{},
	)
	service.DB.AutoMigrate(
		&Account{},
	)

	service.Mailer = mailgun.NewMailgun(domain, apiKey, publicApiKey)

	r := gin.Default()

	r.GET("/leases/:leaseID/terminate", service.TerminatorHandle)
	r.GET("/leases/:leaseID/renew", service.RenewerHandle)
	r.Run() // listen and server on 0.0.0.0:8080
}

type Account struct {
	gorm.Model
	Hello string
}

func (s *Service) TerminatorHandle(c *gin.Context) {
	s.TerminatorQueue.TaskQueue <- TerminatorTask{}

	fmt.Printf("termination of %v initiated", c.Param("leaseID"))

	c.JSON(200, gin.H{
		"message": s.counter,
	})
}

func (s *Service) RenewerHandle(c *gin.Context) {
	s.RenewerQueue.TaskQueue <- RenewerTask{}

	fmt.Printf("renewal of %v initiated", c.Param("leaseID"))

	c.JSON(200, gin.H{
		"message": s.counter,
	})
}
