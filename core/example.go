package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gagliardetto/simpleQueue"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"gopkg.in/mailgun/mailgun-go.v1"
)

const (
	ZeroCloudGuardianSender = "ZeroCloud Guardian <guardian@zerocloud.site>"

	TerminatorActionTerminate = "terminate"
	TerminatorActionShutdown  = "shutdown"
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

// @@@@@@@@@@@@@@@ Task structs @@@@@@@@@@@@@@@

type NewLeaseTask struct {
	AWSAccountID string // message.account
	InstanceID   string // message.detail.instance-id
	Region       string // message.region

	LaunchTime    time.Time // get from the request for tags to ec2 api, not from event
	InstanceType  string
	InstanceOwner string
	//InstanceTags []string
}

type TerminatorTask struct {
	AWSAccountID string
	InstanceID   string
	Region       string // needed? arn:aws:ec2:us-east-1:859795398601:instance/i-fd1f96cc

	Action string // default is TerminatorActionTerminate
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

// @@@@@@@@@@@@@@@ Task consumers @@@@@@@@@@@@@@@

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

// @@@@@@@@@@@@@@@ Periodic Jobs @@@@@@@@@@@@@@@

func (s *Service) EventInjestorJob() error {
	// verify event origin (must be aws, not someone else)
	fmt.Println("EventInjestorJob() run")

	queueURL := fmt.Sprintf("https://sqs.%v.amazonaws.com/%v/%v",
		viper.GetString("AWS_REGION"),
		viper.GetString("AWS_ACCOUNT_ID"),
		viper.GetString("SQSQueueName"),
	)
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL), // Required
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(5), // should be higher, like 30, the time to finish doing everything
		WaitTimeSeconds:     aws.Int64(3),
	}
	resp, err := s.AWS.SQS.ReceiveMessage(params)

	if err != nil {
		return fmt.Errorf("EventInjestorJob() error: %v", err)
	}

	//fmt.Println(resp)
	for messageIndex := range resp.Messages {
		var envelope struct {
			Type             string `json:"Type"`
			MessageId        string `json:"MessageId"`
			TopicArn         string `json:"TopicArn"`
			Message          string `json:"Message"`
			Timestamp        string `json:"Timestamp"`
			SignatureVersion string `json:"SignatureVersion"`
			Signature        string `json:"Signature"`
			SigningCertURL   string `json:"SigningCertURL"`
			UnsubscribeURL   string `json:"UnsubscribeURL"`
		}
		json.Unmarshal([]byte(*resp.Messages[messageIndex].Body), envelope)

		var message struct {
			Version    string   `json:"version"`
			ID         string   `json:"id"`
			DetailType string   `json:"detail-type"`
			Source     string   `json:"source"`
			Account    string   `json:"account"`
			Time       string   `json:"time"`
			Region     string   `json:"region"`
			Resources  []string `json:"resources"`
			Detail     struct {
				InstanceID string `json:"instance-id"`
				State      string `json:"state"`
			} `json:"detail"`
		}
		json.Unmarshal([]byte(envelope.Message), message)

		fmt.Printf("%v", message)
	}

	return nil
}

func (s *Service) AlerterJob() error {

	return nil
}

func (s *Service) SentencerJob() error {

	return nil
}

var logger log15.Logger

func main() {
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

	// TODO: move these config values to config.yml
	var (
		maxWorkers   = 10
		maxQueueSize = 1000
		domain       = ""
		apiKey       = ""
		publicApiKey = ""
	)

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

	service.DB.DropTable(
		&Account{},
	)
	service.DB.AutoMigrate(
		&Account{},
	)

	// @@@@@@@@@@@@@@@ Setup external services @@@@@@@@@@@@@@@

	// setup mailer service
	service.Mailer = mailgun.NewMailgun(domain, apiKey, publicApiKey)

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

// @@@@@@@@@@@@@@@ DB models @@@@@@@@@@@@@@@

type Account struct {
	gorm.Model
	Hello string
}

// @@@@@@@@@@@@@@@ router handles @@@@@@@@@@@@@@@

func (s *Service) TerminatorHandle(c *gin.Context) {
	s.TerminatorQueue.TaskQueue <- TerminatorTask{}

	fmt.Printf("termination of %v initiated", c.Param("leaseID"))
	// /welcome?firstname=Jane&lastname=Doe
	// lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

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

func runForever(f func() error, sleepDuration time.Duration) {
	for {
		err := f()
		if err != nil {
			logger.Error("runForever", fmt.Sprintf("err: %v", err))
		}
		time.Sleep(sleepDuration)
	}
}
