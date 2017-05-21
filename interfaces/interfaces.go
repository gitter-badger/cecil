package interfaces

import (
	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/awstools"
	"github.com/tleyden/cecil/config"
	"github.com/tleyden/cecil/eventrecord"
	"github.com/tleyden/cecil/mailers"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/queues"
)

type CoreServiceInterface interface {
	GormDB() *gorm.DB
	EventRecorder() eventrecord.EventRecord
	DefaultMailer() *mailers.MailerInstance
	AWSRes() *awstools.AWSRes
	Config() *config.Config
	models.DBServiceInterface
	Queues() queues.QueuesGroupInterface
	//queues.QueuesGroupInterface
}