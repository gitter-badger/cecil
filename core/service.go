package core

import (
	"crypto/rsa"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/gagliardetto/simpleQueue"
	"github.com/jinzhu/gorm"
)

// Service is fundamental element of Cecil, and holds most of what is used by Cecil.
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
		DefaultMailer struct {
			Domain       string
			APIKey       string
			PublicAPIKey string
		}
	}
	// TODO: move EC2 into AWS ???
	EC2           Ec2ServiceFactory
	DB            *gorm.DB
	DefaultMailer MailerInstance
	AWS           struct {
		Session *session.Session
		SQS     sqsiface.SQSAPI
		SNS     snsiface.SNSAPI
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

	// The eventRecorder is a KV store used to record events for later
	// analysis.  Events like all SQS messages received, etc.
	eventRecord EventRecord

	slackInstances  map[uint]*SlackInstance  // map account_id to *SlackInstance
	mailerInstances map[uint]*MailerInstance // map account_id to *MailerInstance
	mu              *sync.RWMutex
}

// NewService returns a new service
func NewService() *Service {
	service := &Service{
		eventRecord:     NoOpEventRecord{},
		mu:              &sync.RWMutex{},
		slackInstances:  make(map[uint]*SlackInstance),
		mailerInstances: make(map[uint]*MailerInstance),
	}
	return service
}
