package core

import (
	"fmt"
	"time"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"

	"github.com/go-stack/stack"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	// TODO: move these config values to config.yml
	maxWorkers   = 10
	maxQueueSize = 1000
)

var (
	// DropAllTables is a bool filled in from the CLI flag; decides on startup whether to drop all tables from DB.
	DropAllTables bool
)

// Logger is the logger used in all Cecil;
var Logger log15.Logger

func init() {

	// Setup logger
	Logger = log15.New()
	//log15.Root().SetHandler(log15.CallerStackHandler("%v   %[1]n()", log15.StdoutHandler))
	log15.Root().SetHandler(func(format string, h log15.Handler) log15.Handler {
		return log15.FuncHandler(func(r *log15.Record) error {
			switch r.Lvl {
			case log15.LvlCrit, log15.LvlError, log15.LvlWarn:
				s := stack.Trace().TrimBelow(r.Call).TrimRuntime()
				if len(s) > 0 {
					r.Ctx = append(r.Ctx, "stack", fmt.Sprintf(format, s))
				}
			}
			return h.Log(r)
		})
	}("%v   %[1]n()", log15.StdoutHandler))

	// Setup gorm NowFunc callback.  This is here because it caused race condition
	// issues when it was in SetupDB() which was called from multiple tests
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
}

// SetupAndRun runs all the initialization of Cecil.
func (service *Service) SetupAndRun() *Service {

	// Setup
	service.LoadConfig("config.yml")
	service.GenerateRSAKeys()
	service.SetupQueues()
	service.SetupDB("cecil.db")
	service.SetupSlack()
	service.SetupMailers()

	// @@@@@@@@@@@@@@@ Setup event log @@@@@@@@@@@@@@@

	viper.SetDefault("EventLogDir", "")
	EventLogDir, err := viperMustGetString("EventLogDir")
	if err != nil {
		panic(err)
	}

	if EventLogDir != "" {
		service.SetupEventRecording(true, EventLogDir)
	}

	// @@@@@@@@@@@@@@@ Setup external services @@@@@@@@@@@@@@@

	// Setup mailer client
	service.DefaultMailer.Client = mailgun.NewMailgun(
		service.Config.DefaultMailer.Domain,
		service.Config.DefaultMailer.APIKey,
		service.Config.DefaultMailer.PublicAPIKey,
	)

	// Setup aws session
	AWSCreds := credentials.NewStaticCredentials(
		service.AWS.Config.AWS_ACCESS_KEY_ID,
		service.AWS.Config.AWS_SECRET_ACCESS_KEY,
		"",
	)
	AWSConfig := &aws.Config{
		Credentials: AWSCreds,
	}
	service.AWS.Session = session.New(AWSConfig)

	// Setup sqs
	service.AWS.SQS = sqs.New(service.AWS.Session)

	// Setup sns
	service.AWS.SNS = sns.New(service.AWS.Session)

	// Setup EC2
	service.EC2 = DefaultEc2ServiceFactory

	// Setup CloudFormation
	service.CloudFormation = DefaultCloudFormationServiceFactory

	// @@@@@@@@@@@@@@@ Schedule Periodic Jobs @@@@@@@@@@@@@@@

	SchedulePeriodicJob(service.EventInjestorJob, time.Duration(time.Second*5))
	SchedulePeriodicJob(service.AlerterJob, time.Duration(time.Second*30))
	SchedulePeriodicJob(service.SentencerJob, time.Duration(time.Second*30))

	// @@@@@@@@@@@@@@@ Update external services @@@@@@@@@@@@@@@

	// Regenerate SQS permissions for all cloudaccounts in DB.
	if err := service.RegenerateSQSPermissions(); err != nil {
		Logger.Info(
			"initial RegenerateSQSPermissions:",
			"err", err,
		)
	}

	// Resubscribe to all SNS topics of all cloudaccounts present in DB.
	if err := service.ResubscribeToAllSNSTopics(); err != nil {
		Logger.Info(
			"initial ResubscribeToAllSNSTopics:",
			"err", err,
		)
	}
	return service
}
