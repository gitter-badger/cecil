package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/tleyden/cecil/mailers"
	"github.com/tleyden/cecil/notification"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
	"gopkg.in/mailgun/mailgun-go.v1"
)

// NotifierQueueConsumer consumes NotifierTask from NotifierQueue; sends messages
func (s *Service) NotifierQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(tasks.NotifierTask)
	// TODO: check whether fields are non-null and valid
	Logger.Info("Sending EMAIL",
		"to", task.To,
	)

	// if there is a SlackBotInstance for the Account,
	// send in a goroutine a message to that SlackBotInstance.
	go func() {
		if task.AccountID == 0 {
			return
		}
		slackIns, err := s.SlackBotInstanceByID(task.AccountID)
		if err != nil {
			//Logger.Warn("SlackBotInstanceByID", "warn", err)
			return
		}
		// HACK: the message sent to Slack should have custom formatting;
		// right now it is just the HTML of the email without html tags.
		messageWithoutHTML, err := sanitize.HTMLAllowing(task.BodyHTML, []string{"a"}, []string{"href"})
		messageWithoutHTML = strings.Replace(messageWithoutHTML, `<a href="`, "", -1)
		messageWithoutHTML = strings.Replace(messageWithoutHTML, `">Click here to terminate</a>`, "", -1)
		messageWithoutHTML = strings.Replace(messageWithoutHTML, `">Click here to approve</a>`, "", -1)
		slackIns.OutgoingMessages <- messageWithoutHTML
	}()

	// define the meailer to use (DefaultMailer or a mailer defined by account)
	var mailer *mailers.MailerInstance

	if task.AccountID > 0 {
		mailerIns, err := s.MailerInstanceByID(task.AccountID)
		if err != nil {
			//Logger.Warn("MailerInstanceByID", "warn", err)
		} else {
			mailer = mailerIns
			Logger.Info("using custom mailer", "mailer", *mailer)
		}
	}

	if mailer == nil {
		mailer = s.defaultMailer
	}

	message := mailgun.NewMessage(
		mailer.FromAddress,
		task.Subject,
		task.BodyText,
		task.To,
	)

	message.AddHeader(notification.X_CECIL_MESSAGETYPE, fmt.Sprintf("%s", task.NotificationMeta.NotificationType))
	message.AddHeader(notification.X_CECIL_LEASE_UUID, task.NotificationMeta.LeaseUUID)
	message.AddHeader(notification.X_CECIL_AWS_RESOURCE_ID, task.NotificationMeta.AWSResourceID)
	message.AddHeader(notification.X_CECIL_VERIFICATION_TOKEN, task.NotificationMeta.VerificationToken)

	//message.SetTracking(true)
	if task.DeliverAfter > 0 {
		message.SetDeliveryTime(time.Now().Add(task.DeliverAfter))
	}

	message.SetHtml(task.BodyHTML)

	err := tools.Retry(10, time.Second*5, func() error {
		var err error
		_, _, err = mailer.Client.Send(message)
		return err
	}, nil)
	if err != nil {
		Logger.Error("Error while sending email", "err", err)
		return err
	}

	return nil
}
