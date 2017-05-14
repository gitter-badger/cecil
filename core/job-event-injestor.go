package core

import (
	"fmt"

	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/transmission"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// EventInjestorJob polls the SQS queue, verifies the message, and pushes it to the proper queue
func (s *Service) EventInjestorJob() error {
	// TODO: verify event origin (must be aws, not someone else)

	queueURL := s.SQSQueueURL()

	receiveMessageParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL), // Required
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(90), // should be higher, like 10 (seconds), the time to finish doing everything
		WaitTimeSeconds:     aws.Int64(3),
	}

	// Make sure there is a non-nil SQS
	if s.AWS.SQS == nil {
		Logger.Warn("EventInvestorJob", "SQS == nil, skipping")
		return nil
	}

	Logger.Info("EventInjestorJob(): Polling SQS", "queue", queueURL)
	receiveMessageResponse, err := s.AWS.SQS.ReceiveMessage(receiveMessageParams)

	if err != nil {
		return fmt.Errorf("EventInjestorJob() error: %v", err)
	}

	Logger.Info("SQSmessages",
		"count", len(receiveMessageResponse.Messages),
	)

	for messageIndex := range receiveMessageResponse.Messages {

		tr, err := transmission.GenerateSQSTransmission(s, receiveMessageResponse.Messages[messageIndex], queueURL)
		if err != nil {
			if err == transmission.ErrorEnvelopeIsSubscriptionConfirmation {

				if err := tr.ConfirmSQSSubscription(); err != nil {
					Logger.Warn("ConfirmSQSSubscription", "err", err)
					continue
				}

				if err := tr.DeleteMessage(); err != nil {
					Logger.Warn("DeleteMessage", "err", err)
				}
				continue
			} else {
				Logger.Warn("Error parsing transmission", "err", err)

				err = tr.DeleteMessage()
				if err != nil {
					Logger.Warn("DeleteMessage", "err", err)
				}
				continue
			}
		}

		err = s.ProcessTR(tr)
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *Service) ProcessTR(tr *transmission.Transmission) error {

	Logger.Info("Parsed sqs message", "message", tr.Message)

	if !tr.TopicAndInstanceHaveSameOwner() {
		// the originating SNS topic and the instance have different owners (different AWS accounts)
		// TODO: notify cecil admin
		Logger.Warn("topicAWSID != instanceOriginatorID", "topicAWSID", tr.Topic.AWSID, "instanceOriginatorID", tr.Message.Account)
		err := tr.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return nil
	}

	// consider only pending and terminated status messages; ignore the rest
	if !tr.MessageIsRelevant() {
		Logger.Warn("Ignoring and removing message", "message.Detail.State", tr.Message.Detail.State)
		err := tr.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return nil
	}

	Logger.Info(
		"Creating AssumedConfig",
		"topicRegion", tr.Topic.Region,
		"topicAWSID", tr.Topic.AWSID,
		"externalID", tr.Cloudaccount.ExternalID,
	)

	if err := tr.CreateAssumedService(); err != nil {
		// TODO: this might reveal too much to the admin about the service; be selective and cautious
		s.sendMisconfigurationNotice(err, tr.AdminAccount.Email)
		Logger.Warn("error while creating assumed service", "err", err)
		return err
	}

	if err := tr.CreateAssumedEC2Service(); err != nil {
		// TODO: this might reveal too much to the admin about the service; be selective and cautious
		s.sendMisconfigurationNotice(err, tr.AdminAccount.Email)
		Logger.Warn("error while creating ec2 service with assumed service", "err", err)
		return err
	}

	if err := tr.CreateAssumedAutoscalingService(); err != nil {
		// TODO: this might reveal too much to the admin about the service; be selective and cautious
		s.sendMisconfigurationNotice(err, tr.AdminAccount.Email)
		Logger.Warn("error while creating autoscaling service with assumed service", "err", err)
		return err
	}

	// CreateAssumedCloudformationService which will be used to check whether the instance is part of a cloudformation stack
	if err := tr.CreateAssumedCloudformationService(); err != nil {
		Logger.Warn("error while creating assumed cloudformation service", "err", err)
		return err
	}

	if err := tr.DescribeInstance(); err != nil {
		if err == transmission.ErrInstanceDoesNotExist {
			Logger.Warn("Instance does not exist", "instanceID", tr.Message.Detail.InstanceID)
			// remove message from queue
			err := tr.DeleteMessage()
			if err != nil {
				Logger.Warn("DeleteMessage", "err", err)
			}

			// send transmission to InstanceTerminatedQueue
			s.Queues().InstanceTerminatedQueue().PushTask(tasks.InstanceTerminatedTask{
				Transmission: tr,
			})
			return err
		}
		// TODO: this might reveal too much to the admin about the service; be selective and cautious
		s.sendMisconfigurationNotice(err, tr.AdminAccount.Email)
		Logger.Warn("error while describing instances", "err", err)
		return err
	}

	Logger.Info(
		"describeInstances",
		"response", tr.DescribeInstancesResponse,
	)

	switch tr.Message.Detail.State {
	case ec2.InstanceStateNamePending:
		{
			// send transmission to NewInstanceQueue
			s.Queues().NewInstanceQueue().PushTask(tasks.NewInstanceTask{
				Transmission: tr,
			})
		}
	case ec2.InstanceStateNameTerminated:
		{
			// send transmission to InstanceTerminatedQueue
			s.Queues().InstanceTerminatedQueue().PushTask(tasks.InstanceTerminatedTask{
				Transmission: tr,
			})
		}
	default:
		return nil
	}

	return nil
}
