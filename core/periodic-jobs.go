package core

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// EventInjestorJob polls the SQS queue, verifies the message, and pushes it to the proper queue
func (s *Service) EventInjestorJob() error {
	// TODO: verify event origin (must be aws, not someone else)

	queueURL := fmt.Sprintf("https://sqs.%v.amazonaws.com/%v/%v",
		viper.GetString("AWS_REGION"),
		viper.GetString("AWS_ACCOUNT_ID"),
		viper.GetString("SQSQueueName"),
	)

	receiveMessageParams := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL), // Required
		//MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout: aws.Int64(3), // should be higher, like 10 (seconds), the time to finish doing everything
		WaitTimeSeconds:   aws.Int64(3),
	}

	logger.Info("EventInjestorJob(): Polling SQS", "queue", queueURL)
	receiveMessageResponse, err := s.AWS.SQS.ReceiveMessage(receiveMessageParams)

	if err != nil {
		return fmt.Errorf("EventInjestorJob() error: %v", err)
	}

	logger.Info("SQSmessages",
		"count", len(receiveMessageResponse.Messages),
	)

	for messageIndex := range receiveMessageResponse.Messages {

		transmission, err := s.parseSQSTransmission(receiveMessageResponse.Messages[messageIndex], queueURL)
		if err != nil {
			logger.Warn("Error parsing transmission", "error", err)
			// TODO: delete message
			err = transmission.DeleteMessage()
			if err != nil {
				logger.Warn("DeleteMessage", "error", err)
			}
			continue
		}

		logger.Info("Parsed sqs message", "message", transmission.Message)

		if !transmission.TopicAndInstanceHaveSameOwner() {
			// the originating SNS topic and the instance have different owners (different AWS accounts)
			// TODO: notify zerocloud admin
			logger.Warn("topicAWSID != instanceOriginatorID", "topicAWSID", transmission.Topic.AWSID, "instanceOriginatorID", transmission.Message.Account)
			// TODO: delete message
			continue
		}

		// consider only pending and terminated status messages; ignore the rest
		if !transmission.MessageIsRelevant() {
			logger.Warn("Ignoring and removing message", "message.Detail.State", transmission.Message.Detail.State)
			err := transmission.DeleteMessage()
			if err != nil {
				logger.Warn("DeleteMessage", "error", err)
			}
			continue // next message
		}

		// send transmission to NewLeaseQueue
		s.NewLeaseQueue.TaskQueue <- NewLeaseTask{
			Transmission: transmission,
		}

	}

	return nil
}

// AlerterJob polls the DB for leases that are about to expire, and notifes the owner of the imminent expiry
func (s *Service) AlerterJob() error {
	// find lease that expire in 24 hours
	// find owner
	// create links to extend and terminate lease
	// mark as alerted = true
	// registed new lease's token_once
	// compose email with link to extend and terminate lease
	// send email

	var expiringLeases []Lease
	var expiringLeasesCount int64

	s.DB.Table("leases").
		Where("expires_at < ?",
			time.Now().UTC().Add(ZCDefaultForewarningBeforeExpiry),
		).
		Not("terminated", true).
		Not("alerted", true).
		Find(&expiringLeases).
		Count(&expiringLeasesCount)

	logger.Info("AlerterJob(): Expiring leases", "count", expiringLeasesCount)

	// TODO: create ExpiringLeaseQueue and pass to it this task

	for _, expiringLease := range expiringLeases {

		logger.Info("Expiring lease",
			"instanceID", expiringLease.InstanceID,
			"leaseID", expiringLease.ID,
		)

		var owner Owner
		var ownerCount int64

		s.DB.Table("owners").Where(expiringLease.OwnerID).First(&owner).Count(&ownerCount)

		if ownerCount != 1 {
			logger.Warn("AlerterJob: ownerCount is not 1", "count", ownerCount)
			continue
		}

		// these will be used to compose the urls and verify the requests
		token_once := uuid.NewV4().String() // one-time token

		expiringLease.TokenOnce = token_once
		expiringLease.Alerted = true

		s.DB.Save(&expiringLease)

		// URL to extend lease
		extend_url, err := s.generateSignedEmailActionURL("extend", expiringLease.UUID, expiringLease.InstanceID, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		// URL to terminate lease
		terminate_url, err := s.generateSignedEmailActionURL("terminate", expiringLease.UUID, expiringLease.InstanceID, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		newEmailBody := compileEmail(
			`Hey {{.owner_email}}, instance <b>{{.instance_id}}</b>
				(of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>) is expiring.

				<br>
				<br>

				The instance will expire on <b>{{.termination_time}}</b> ({{.instance_duration}} after it's creation).

				<br>
				<br>

				The instance was created on {{.instance_created_at}}.
				
				<br>
				<br>
				
				Terminate immediately:
				<br>
				<br>
				<a href="{{.instance_terminate_url}}" target="_blank">Click here to <b>terminate</b></a>

				<br>
				<br>

				Extend lease by <b>{{.extend_by}}</b>:
				<br>
				<br>
				<a href="{{.instance_extend_url}}" target="_blank">Click here to <b>extend</b></a>

				<br>
				<br>
				Thanks for using ZeroCloud!
				`,

			map[string]interface{}{
				"owner_email":     owner.Email,
				"instance_id":     expiringLease.InstanceID,
				"instance_type":   expiringLease.InstanceType,
				"instance_region": expiringLease.Region,

				"instance_created_at": expiringLease.CreatedAt.Format("2006-01-02 15:04:05 GMT"),
				"extend_by":           ZCDefaultLeaseDuration.String(),

				"termination_time":  expiringLease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
				"instance_duration": expiringLease.ExpiresAt.Sub(expiringLease.CreatedAt).String(),

				"instance_terminate_url": terminate_url,
				"instance_extend_url":    extend_url,
			},
		)

		s.NotifierQueue.TaskQueue <- NotifierTask{
			From:     ZCMailerFromAddress,
			To:       owner.Email,
			Subject:  fmt.Sprintf("Instance (%v) will expire soon", expiringLease.InstanceID),
			BodyHTML: newEmailBody,
			BodyText: newEmailBody,
		}
	}

	return nil
}

// SentencerJobpolls the DB for expired leases and pushes them to the TerminatorQueue
func (s *Service) SentencerJob() error {

	var expiredLeases []Lease
	var expiredLeasesCount int64

	s.DB.Table("leases").Where("expires_at < ?", time.Now().UTC()).Not("terminated", true).Find(&expiredLeases).Count(&expiredLeasesCount)

	logger.Info("SentencerJob(): Expired leases", "count", expiredLeasesCount)

	for _, expiredLease := range expiredLeases {
		logger.Info("expired lease",
			"instanceID", expiredLease.InstanceID,
			"leaseID", expiredLease.ID,
		)
		s.TerminatorQueue.TaskQueue <- TerminatorTask{Lease: expiredLease}
	}

	return nil
}

func (s *Service) sendMisconfigurationNotice(err error, emailRecipient string) {
	newEmailBody := compileEmail(
		`Hey it appears that ZeroCloud is mis-configured.
		<br>
		<br>
		Error:
		<br>
		{{.err}}`,
		map[string]interface{}{
			"err": err,
		},
	)

	s.NotifierQueue.TaskQueue <- NotifierTask{
		From:     ZCMailerFromAddress,
		To:       emailRecipient,
		Subject:  "ZeroCloud configuration problem",
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
	}

}
