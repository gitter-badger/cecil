package core

import (
	"crypto/rsa"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/gagliardetto/simpleQueue"
	"github.com/jinzhu/gorm"
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
			HostName string // e.g. zerocloud.co
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
			UseMockAWS            bool
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

func (service *Service) SetupQueues() {

	service.NewLeaseQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NewLeaseQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.NewLeaseQueueConsumer error:", "error", err)
		})
	service.NewLeaseQueue.Start()

	service.TerminatorQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.TerminatorQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.TerminatorQueueConsumer error:", "error", err)
		})
	service.TerminatorQueue.Start()

	service.LeaseTerminatedQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.LeaseTerminatedQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.LeaseTerminatedQueueConsumer error:", "error", err)
		})
	service.LeaseTerminatedQueue.Start()

	service.ExtenderQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.ExtenderQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.ExtenderQueueConsumer error:", "error", err)
		})
	service.ExtenderQueue.Start()

	service.NotifierQueue = simpleQueue.NewQueue().
		SetMaxSize(maxQueueSize).
		SetWorkers(maxWorkers).
		SetConsumer(service.NotifierQueueConsumer).
		SetErrorCallback(func(err error) {
			logger.Error("service.NotifierQueueConsumer error:", "error", err)
		})
	service.NotifierQueue.Start()

}

func (service *Service) Stop() {
	service.NewLeaseQueue.Stop()
	service.TerminatorQueue.Stop()
	service.LeaseTerminatedQueue.Stop()
	service.ExtenderQueue.Stop()
	service.NotifierQueue.Stop()
}
