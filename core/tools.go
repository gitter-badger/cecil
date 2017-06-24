package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/inconshreveable/log15"
	"github.com/tleyden/cecil/notification"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
)

// NewContextLogger returns a new context logger which has been filled in with the request ID
func NewContextLogger(ctx context.Context) log15.Logger {
	request := goa.ContextRequest(ctx)
	if request == nil {
		return Logger.New()
	}
	return Logger.New(
		"url", request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)
}

// sendMisconfigurationNotice sends a misconfiguration notice to emailRecipient.
func (s *Service) sendMisconfigurationNotice(err error, emailRecipient string) {
	newEmailBody, err := tools.CompileEmailTemplate(
		"misconfiguration-notice.html",
		map[string]interface{}{
			"err": err,
		},
	)
	if err != nil {
		return
	}

	s.queues.NotifierQueue().PushTask(tasks.NotifierTask{
		To:               emailRecipient,
		Subject:          "Cecil configuration problem",
		BodyHTML:         newEmailBody,
		BodyText:         newEmailBody,
		NotificationMeta: notification.NotificationMeta{NotificationType: notification.Misconfiguration},
	})
}

// CecilHTTPAddress returns the complete HTTP address of cecil; e.g. https://127.0.0.1:8080
func (s *Service) CecilHTTPAddress() string {
	// TODO check the prefix of Port; ignore port if 80 or 443 (decide looking at Scheme)
	return fmt.Sprintf("%v://%v%v",
		s.config.Server.Scheme,
		s.config.Server.HostName,
		s.config.Server.Port,
	)
}

// SQSQueueURL returns the HTTP URL of the SQS queue.
func (s *Service) SQSQueueURL() string {
	return fmt.Sprintf("https://sqs.%v.amazonaws.com/%v/%v",
		s.AWS.Config.AWS_REGION,
		s.AWS.Config.AWS_ACCOUNT_ID,
		s.AWS.Config.SQSQueueName,
	)
}

// SQSQueueArn returns the AWS ARN of the SQS queue.
func (s *Service) SQSQueueArn() string {
	return fmt.Sprintf("arn:aws:sqs:%v:%v:%v",
		s.AWS.Config.AWS_REGION,
		s.AWS.Config.AWS_ACCOUNT_ID,
		s.AWS.Config.SQSQueueName,
	)
}

// DefineLeaseDuration tries to define the duration a lease should have basing the decision
// on many sources, each of which has a hierarchy
func (s *Service) DefineLeaseDuration(accountID, cloudaccountID uint) (time.Duration, error) {
	account, err := s.GetAccountByID(int(accountID))
	if err != nil {
		return 0, err
	}

	cloudaccount, err := s.GetCloudaccountByID(int(cloudaccountID))
	if err != nil {
		return 0, err
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		return 0, errors.New("!account.IsOwnerOf(cloudaccount)")
	}

	var leaseDuration time.Duration
	// Use global cecil lease duration setting
	leaseDuration = time.Duration(s.config.Lease.Duration)

	// Use lease duration setting of account
	if account.DefaultLeaseDuration > 0 {
		leaseDuration = time.Duration(account.DefaultLeaseDuration)
	}

	// Use lease duration setting of cloudaccount
	if cloudaccount.DefaultLeaseDuration > 0 {
		leaseDuration = time.Duration(cloudaccount.DefaultLeaseDuration)
	}

	return leaseDuration, nil
}
