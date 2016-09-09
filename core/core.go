package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gagliardetto/simpleQueue"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"
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

	ZCMaxLeasesPerOwner             = 10
	ZCDefaultLeaseExpiration uint64 = 100

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

// @@@@@@@@@@@@@@@ Task structs @@@@@@@@@@@@@@@

type NewLeaseTask struct {
	AWSAccountID string // message.account
	InstanceID   string // message.detail.instance-id
	Region       string // message.region

	LaunchedAt    time.Time // get from the request for tags to ec2 api, not from event
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
	// TODO: check whether fields are non-null and valid

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
	// TODO: check whether fields are non-null and valid

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
	// TODO: check whether fields are non-null and valid

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
		QueueUrl: aws.String(queueURL), // Required
		//MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout: aws.Int64(3), // should be higher, like 30, the time to finish doing everything
		WaitTimeSeconds:   aws.Int64(3),
	}
	resp, err := s.AWS.SQS.ReceiveMessage(params)

	if err != nil {
		return fmt.Errorf("EventInjestorJob() error: %v", err)
	}

	fmt.Println(resp)
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
		err := json.Unmarshal([]byte(*resp.Messages[messageIndex].Body), &envelope)
		if err != nil {
			return err
		}

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
		err = json.Unmarshal([]byte(envelope.Message), &message)
		if err != nil {
			return err
		}

		topicArn := strings.Split(envelope.TopicArn, ":")
		topicRegion := topicArn[3]
		topicOwnerID, err := strconv.ParseUint(topicArn[4], 10, 64)
		if err != nil {
			// TODO: notify
			fmt.Println(err)
			continue
		}
		// topicName := topicArn[5]
		instanceOriginatorID, err := strconv.ParseUint(message.Account, 10, 64)
		if err != nil {
			// TODO: notify
			fmt.Println(err)
			continue
		}
		// TODO: check these values are not empty

		if topicOwnerID != instanceOriginatorID {
			// the originating SNS topic and the instance have different owners
			// TODO: notify zerocloud admin
			fmt.Println("topicOwnerID != instanceOriginatorID")
			continue
		}

		// consider only pending and terminated status messages; ignore the rest
		if message.Detail.State != ec2.InstanceStateNamePending &&
			message.Detail.State != ec2.InstanceStateNameTerminated {
			fmt.Println("removing")
			// remove message from queue
			params := &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),                                   // Required
				ReceiptHandle: aws.String(*resp.Messages[messageIndex].ReceiptHandle), // Required
			}
			_, err := s.AWS.SQS.DeleteMessage(params)
			if err != nil {
				// In case of error just leave it there, and on the next turn it will be retried
				fmt.Println(err)
			}
			continue // next message
		}

		// HasOwner: check whether someone with this aws account id is registered
		var cloudAccount CloudAccount
		var cloudOwnerCount int64
		s.DB.Where(&CloudAccount{AWSID: topicOwnerID}).First(&cloudAccount).Count(&cloudOwnerCount)
		if cloudOwnerCount == 0 {
			// TODO: notify admin; something fishy is going on.
			continue
		}

		// <debug>
		var accounts []Account
		s.DB.Table("accounts").Find(&accounts)
		fmt.Printf("accounts: %#v\n", accounts)
		fmt.Printf("%s\n", accounts[0].ID)
		fmt.Println("looking for:", cloudAccount.AccountID)
		// </debug>

		var account Account
		var cloudAccountOwnerCount int64
		s.DB.Model(&cloudAccount).Related(&account).Count(&cloudAccountOwnerCount)
		//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&cloudAccountOwnerCount)
		if cloudAccountOwnerCount == 0 {
			// TODO: notify admin; something fishy is going on.
			fmt.Println("cloudAccountOwnerCount == 0")
			continue
		}

		fmt.Printf("account: %#v\n", account)

		// IsNew: check whether a lease with the same instanceID exists
		var instanceCount int64
		s.DB.Table("leases").Where(&Lease{InstanceID: message.Detail.InstanceID}).Count(&instanceCount)
		fmt.Println("here")
		if instanceCount != 0 {
			// TODO: notify admin
			fmt.Println("instanceCount != 0")
			continue
		}

		// assume role

		assumedConfig := &aws.Config{
			Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
				Client:          sts.New(s.AWS.Session, &aws.Config{Region: aws.String(topicRegion)}),
				RoleARN:         fmt.Sprintf("arn:aws:iam::%v:role/ZeroCloudRole", topicOwnerID),
				RoleSessionName: uuid.NewV4().String(),
				ExternalID:      aws.String("slavomir"),
				ExpiryWindow:    3 * time.Minute,
			}),
		}

		assumedService := session.New(assumedConfig)

		ec2Service := ec2.New(assumedService,
			&aws.Config{
				Region: aws.String(topicRegion),
			},
		)

		paramsDescribeInstance := &ec2.DescribeInstancesInput{
			// DryRun: aws.Bool(true),
			InstanceIds: []*string{
				aws.String(message.Detail.InstanceID),
			},
		}
		resp, err := ec2Service.DescribeInstances(paramsDescribeInstance)

		if err != nil {
			// TODO: notify
			fmt.Println(err)
			continue
		}

		// ExistsOnAWS: check whether the instance specified in the event exists on aws
		if len(resp.Reservations) == 0 {
			fmt.Println("len(resp.Reservations) == 0: ")
			continue
		}
		if len(resp.Reservations[0].Instances) == 0 {
			fmt.Println("len(resp.Reservations[0].Instances) == 0: ")
			continue
		}
		fmt.Println("description: ", resp)

		var instance = resp.Reservations[0].Instances[0]

		//instance.InstanceType
		//instance.LaunchTime

		if *instance.InstanceId != message.Detail.InstanceID {
			fmt.Println("instance.InstanceId !=message.Detail.InstanceID")
			continue
		}

		if *instance.State.Name != ec2.InstanceStateNamePending &&
			*instance.State.Name != ec2.InstanceStateNameRunning {
			fmt.Println("the retried state is neither pending not running:", instance.State.Name)
			continue
		}

		var ownerIsAdmin bool = false
		var ownerEmail string

		// InstanceHasTags: check whethe instance has tags
		if len(instance.Tags) == 0 {
			fmt.Println("len(instance.Tags) == 0")
			ownerIsAdmin = true
		} else {

			// InstanceHasOwnerTag: check whether the instance has an zerocloudowner tag
			for _, tag := range instance.Tags {
				if strings.ToLower(*tag.Key) == "zerocloudowner" {

					// OwnerTagValueIsValid: check whether the zerocloudowner tag is a valid email
					ownerTag, err := s.Mailer.ValidateEmail(*tag.Value)
					if err != nil {
						fmt.Println(err)
						ownerIsAdmin = true
						break
					}
					if !ownerTag.IsValid {
						fmt.Println("email not valid")
						ownerIsAdmin = true
						// TODO: notify admin: "Warning: zerocloudowner tag email not valid" (DO NOT INCLUDE IT IN THE EMAIL, OR HTML-ESCAPE IT)
						break
					}
					fmt.Printf("Parts local_part=%s domain=%s display_name=%s", ownerTag.Parts.LocalPart, ownerTag.Parts.Domain, ownerTag.Parts.DisplayName)
					ownerEmail = ownerTag.Address
					break
				}
			}
		}

		var owners []Owner
		var ownerCount int64
		if ownerEmail != "" && !ownerIsAdmin && ownerEmail != account.Email {
			// OwnerTagIsWhitelisted: check whether the owner email in the tag is a whitelisted owner email

			// TODO: select Owner by email, cloudaccountid, and region?
			s.DB.Table("owners").Where(&Owner{Email: ownerEmail, CloudAccountID: cloudAccount.ID}).Find(&owners).Count(&ownerCount)
			if ownerCount == 0 {
				// TODO: owner is not whitelisted: notify admin: "Warning: zerocloudowner tag email not in whitelist"
				ownerIsAdmin = true
			}
			if ownerCount > 1 {
				// TODO: fatal: too many owners
				ownerIsAdmin = true
			}
		}

		if ownerIsAdmin {
			ownerEmail = account.Email

			// TODO: select Owner by email, cloudaccountid, and region?
			s.DB.Table("owners").Where(&Owner{Email: ownerEmail, CloudAccountID: cloudAccount.ID}).Find(&owners).Count(&ownerCount)
			if ownerCount == 0 {
				// TODO: fatal: admin is not in the owner table
				fmt.Println("fatal: admin is not in the owner table")
				continue
			}
		}

		var owner = owners[0]

		var leases []Lease
		var activeLeaseCount int64
		s.DB.Table("leases").Where(&Lease{
			OwnerID:        owner.ID,
			CloudAccountID: cloudAccount.ID,
			Terminated:     false,
		}).Find(&leases).Count(&activeLeaseCount)
		//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&activeLeaseCount)

		var lifetime time.Duration = time.Duration(ZCDefaultLeaseExpiration)

		if account.DefaultLeaseExpiration > 0 {
			lifetime = time.Duration(account.DefaultLeaseExpiration)
		}
		if cloudAccount.DefaultLeaseExpiration > 0 {
			lifetime = time.Duration(cloudAccount.DefaultLeaseExpiration)
		}

		leaseNeedsApproval := activeLeaseCount >= ZCMaxLeasesPerOwner
		if !leaseNeedsApproval {
			// TODO:
			// register new lease in DB
			// set its expiration to zone.default_expiration (if > 0), or cloudAccount.default_expiration, or account.default_expiration
			var terminationTime = time.Now().Add(lifetime)

			newLease := Lease{
				OwnerID:        owner.ID,
				CloudAccountID: cloudAccount.ID,
				AWSAccountID:   cloudAccount.AWSID,

				InstanceID: *instance.InstanceId,
				Region:     *instance.Placement.AvailabilityZone,

				// Terminated bool `sql:"DEFAULT:false"`
				// Deleted    bool `sql:"DEFAULT:false"`

				LaunchedAt:   *instance.LaunchTime,
				ExpiresAt:    time.Now().Add(lifetime),
				InstanceType: *instance.InstanceType,
			}
			s.DB.Create(&newLease)

			newEmailBody := compileEmail(
				`Hey {{.owner_email}}, you (or someone else) created a new instance 
				id <b>({{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). That's AWESOME!

				<br>
				<br>

				Your instance will be terminated at {{.termination_time}} ({{.instance_lifetime}} after it's creation).

				<br>
				<br>
				
				Thanks for using ZeroCloud!
				`,

				map[string]interface{}{
					"owner_email":     owner.Email,
					"instance_id":     instance.InstanceId,
					"instance_type":   instance.InstanceType,
					"instance_region": instance.Placement.AvailabilityZone,

					"termination_time":  terminationTime.Format("2006-01-02 15:04:05 CET"),
					"instance_lifetime": lifetime.String(),
				},
			)
			s.NotifierQueue.TaskQueue <- NotifierTask{
				From:     ZCMailerFromAddress,
				To:       owner.Email,
				Subject:  fmt.Sprintf("Instance (%v) Created", instance.InstanceId),
				BodyHTML: newEmailBody,
				BodyText: newEmailBody,
			}

			continue
		} else {
			// TODO:
			// register new lease in DB
			//expiry: 1h
			// send confirmation to owner: confirmation link

			lifetime = time.Duration(time.Hour)
			var terminationTime = time.Now().Add(lifetime)

			newLease := Lease{
				OwnerID:        owner.ID,
				CloudAccountID: cloudAccount.ID,
				AWSAccountID:   cloudAccount.AWSID,

				InstanceID: *instance.InstanceId,
				Region:     *instance.Placement.AvailabilityZone,

				// Terminated bool `sql:"DEFAULT:false"`
				// Deleted    bool `sql:"DEFAULT:false"`

				LaunchedAt:   *instance.LaunchTime,
				ExpiresAt:    terminationTime,
				InstanceType: *instance.InstanceType,
			}
			s.DB.Create(&newLease)

			newEmailBody := compileEmail(
				`Hey {{.owner_email}}, you (or someone else) created a new instance 
				id <b>({{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). <br><br>

				At the time of writing this email, you have {{.n_of_active_leases}} active
					leases, so we need your approval for this one. <br><br>

				Please click on "Approve" to approve this instance,
					otherwise it will be terminated at {{.termination_time}} (one hour after it's creation).

				Approve:
				<br>
				{{.instance_renew_url}}

				<br>
				<br>
				
				Terminate:
				<br>
				{{.instance_terminate_url}}
				`,

				map[string]interface{}{
					"owner_email":        owner.Email,
					"n_of_active_leases": activeLeaseCount,
					"instance_id":        instance.InstanceId,
					"instance_type":      instance.InstanceType,
					"instance_region":    instance.Placement.AvailabilityZone,

					"termination_time":       terminationTime.Format("2006-01-02 15:04:05 CET"),
					"instance_renew_url":     "",
					"instance_terminate_url": "",
				},
			)
			s.NotifierQueue.TaskQueue <- NotifierTask{
				From:     ZCMailerFromAddress,
				To:       owner.Email,
				Subject:  fmt.Sprintf("Instance (%v) Needs Approval", instance.InstanceId),
				BodyHTML: newEmailBody,
				BodyText: newEmailBody,
			}
			continue
		}

		// if message.Detail.State == ec2.InstanceStateNameTerminated
		// LeaseTerminatedQueue <- LeaseTerminatedTask{} and continue

		// get zc account who has a cloudaccount with awsID == topicOwnerID
		// if no one of our customers owns this account, error
		// fetch options config
		// roleARN := fmt.Sprintf("arn:aws:iam::%v:role/ZeroCloudRole",topicOwnerID)
		// assume role
		// fetch instance info
		// check if statuses match (this message was sent by aws.ec2)
		// message.Detail.InstanceID

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
				Provider: "aws",
				AWSID:    859795398601,
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
		Email: "slv.balsan@gmail.com",
	}
	service.DB.Create(&firstOwner)

	secondaryOwner := Owner{
		Email: "slavomir.balsan@gmail.com",
	}
	service.DB.Create(&secondaryOwner)

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

// @@@@@@@@@@@@@@@ DB models @@@@@@@@@@@@@@@

type Account struct {
	gorm.Model
	Email string `sql:"size:255;unique;index"`

	Disabled bool `sql:"DEFAULT:false"`
	Deleted  bool `sql:"DEFAULT:false"`

	DefaultLeaseExpiration uint64 `sql:"DEFAULT:0"`

	CloudAccounts []CloudAccount
}

type CloudAccount struct {
	gorm.Model
	AccountID uint

	DefaultLeaseExpiration uint64 `sql:"DEFAULT:0"`
	Provider               string // e.g. AWS
	AWSID                  uint64 `sql:"size:255;unique;index"`

	Disabled bool `sql:"DEFAULT:false"`
	Deleted  bool `sql:"DEFAULT:false"`

	Leases  []Lease
	Regions []Region
	Owners  []Owner
}

type Lease struct {
	gorm.Model
	CloudAccountID uint
	OwnerID        uint

	AWSAccountID uint64
	InstanceID   string
	Region       string

	Terminated bool `sql:"DEFAULT:false"`
	Deleted    bool `sql:"DEFAULT:false"`

	LaunchedAt   time.Time
	ExpiresAt    time.Time
	InstanceType string
}

type Region struct {
	gorm.Model
	CloudAccountID uint

	Region string

	Deleted bool `sql:"DEFAULT:false"`
}

type Owner struct {
	gorm.Model
	CloudAccountID uint

	Email    string
	Disabled bool `sql:"DEFAULT:false"`
	Leases   []Lease
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
