package core

import (
	"fmt"
	"time"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"

	"github.com/go-stack/stack"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tleyden/cecil/awstools"
	"github.com/tleyden/cecil/mailers"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/slackbot"
	"github.com/tleyden/cecil/tools"
	"github.com/tleyden/cecil/transmission"

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

	// DBType is a string filled in from the CLI flag; decides on startup whether to use sqlite or postgres.
	DBType string
)

// Logger is the logger used in all Cecil;
var Logger log15.Logger

func init() {
	// Setup logger for this package; this logger will be used in many other packages (see next init func).
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

func init() {
	// Export Logger to other packages, to have only one, consistent logger everywhere.
	models.Logger = Logger.New(
		"package", "models",
	)
	transmission.Logger = Logger.New(
		"package", "transmission",
	)
	mailers.Logger = Logger.New(
		"package", "mailers",
	)
	slackbot.Logger = Logger.New(
		"package", "slackbot",
	)
}

// SetupAndRun runs all the initialization of Cecil.
func (service *Service) SetupAndRun() *Service {

	// Setup
	service.LoadConfig("config.yml")
	service.GenerateRSAKeys()
	service.SetupQueues()
	service.SetupDB("cecil.db")

	service.SlackBotService = slackbot.NewService(service)
	service.SetupSlack()

	// TODO: mailers.NewService returns error if dbService is nil
	service.CustomMailerService = mailers.NewService(service.DBService)
	err := service.SetupMailers()
	if err != nil {
		panic(err)
	}

	// @@@@@@@@@@@@@@@ Setup event log @@@@@@@@@@@@@@@

	viper.SetDefault("EventLogDir", "")
	EventLogDir, err := tools.ViperMustGetString("EventLogDir")
	if err != nil {
		panic(err)
	}

	if EventLogDir != "" {
		service.SetupEventRecording(true, EventLogDir)
	}

	// @@@@@@@@@@@@@@@ Setup external services @@@@@@@@@@@@@@@

	// Setup mailer client
	service.defaultMailer.Client = mailgun.NewMailgun(
		service.config.DefaultMailer.Domain,
		service.config.DefaultMailer.APIKey,
		service.config.DefaultMailer.PublicAPIKey,
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

	// max retries for request
	AWSConfig.MaxRetries = aws.Int(3)

	/*	// set custom logger
		AWSConfig.Logger = aws.LoggerFunc(func(args ...interface{}) {
			namedArgs := []interface{}{}
			for argIndex, arg := range args {
				namedArgs = append(namedArgs, strconv.Itoa(argIndex), arg)
			}
			Logger.Error("AWS DEBUG:", namedArgs...)
		})
		// set log level
		AWSConfig.LogLevel = aws.LogLevel(aws.LogDebugWithRequestErrors)*/

	service.AWS.Session = session.New(AWSConfig)

	// Setup sqs
	service.AWS.SQS = sqs.New(service.AWS.Session)

	// Setup sns
	service.AWS.SNS = sns.New(service.AWS.Session)

	// Setup EC2
	service.AWS.EC2 = awstools.DefaultEc2ServiceFactory

	// Setup CloudFormation
	service.AWS.CloudFormation = awstools.DefaultCloudFormationServiceFactory

	// Setup AutoScaling
	service.AWS.AutoScaling = awstools.DefaultAutoScalingServiceFactory

	// @@@@@@@@@@@@@@@ Schedule Periodic Jobs @@@@@@@@@@@@@@@

	commonLog := func(err error) {
		Logger.Error("SchedulePeriodicJob", "err", err)
	}

	tools.SchedulePeriodicJob(service.EventInjestorJob, time.Duration(time.Second*5), commonLog)
	tools.SchedulePeriodicJob(service.AlerterJob, time.Duration(time.Second*30), commonLog)
	tools.SchedulePeriodicJob(service.SentencerJob, time.Duration(time.Second*30), commonLog)

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
