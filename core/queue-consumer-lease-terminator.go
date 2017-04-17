package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/tasks"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
)

// @@@@@@@@@@@@@@@ Task consumers @@@@@@@@@@@@@@@

// TerminatorQueueConsumer consumes TerminatorTask from TerminatorQueue;
// sends instance termination request to AWS ec2.
func (s *Service) TerminatorQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(tasks.TerminatorTask)
	// TODO: check whether fields are non-null and valid
	Logger.Info("TerminatorQueueConsumer",
		"task", task,
	)

	var cloudaccount models.Cloudaccount
	err := s.DB.Model(&task.Lease).Related(&cloudaccount).Error
	//s.DB.Table("accounts").Where([]uint{cloudaccount.AccountID}).First(&cloudaccount).Count(&leaseCloudOwnerCount)
	if err != nil {
		// TODO: notify admin; something fishy is going on.
		Logger.Warn("here", err.Error())
		return err
	}

	// assume role
	assumedConfig := &aws.Config{
		Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
			Client: sts.New(s.AWS.Session, &aws.Config{Region: aws.String(task.Region)}),
			RoleARN: fmt.Sprintf(
				"arn:aws:iam::%v:role/%v",
				cloudaccount.AWSID,
				s.AWS.Config.ForeignIAMRoleName,
			),
			RoleSessionName: uuid.NewV4().String(),
			ExternalID:      aws.String(cloudaccount.ExternalID),
			ExpiryWindow:    3 * time.Minute,
		}),
	}

	assumedService := session.New(assumedConfig)

	if task.Lease.GroupType == models.GroupASG {

		assumedAutoScalingService := s.AWSRes().AutoScaling(assumedService, task.Region)

		// TODO: extract ASG name from ASG ARN
		asgName := task.Lease.AwsContainerName

		req := autoscaling.DeleteAutoScalingGroupInput{}
		req.SetAutoScalingGroupName(asgName)
		req.SetForceDelete(true)

		Logger.Info(
			"Deleteting ASG",
			"asg_name", asgName,
			"lease_id", task.Lease.ID,
			"account_id", task.AccountID,
		)

		resp, err := assumedAutoScalingService.DeleteAutoScalingGroup(&req)
		if err != nil {
			Logger.Error(
				"error while deleting ASG",
				"asg_name", asgName,
				"lease_id", task.Lease.ID,
				"account_id", task.AccountID,
			)
			return err
		}
		Logger.Info("DeleteAutoScalingGroup",
			"response", resp.String(),
		)
		return nil
	}

	// terminate all active instances for the group
	instances, err := s.ActiveInstancesForGroup(task.AccountID, &task.CloudaccountID, task.GroupUID)
	if err != nil {
		return err
	}

	instanceIDs := []*string{}
	for i := range instances {
		instanceIDs = append(instanceIDs, &instances[i].InstanceID)
	}

	Logger.Info(
		"Terminating lease",
		"group_uid", task.GroupUID,
		"group_type", task.GroupType.String(),
		"lease_id", task.ID,
		"total_instances", len(instanceIDs),
		"instanceIDs", instanceIDs,
	)
	assumedEC2Service := s.AWS.EC2(assumedService, task.Region)

	terminateInstanceParams := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs, // Required
	}
	terminateInstanceResponse, err := assumedEC2Service.TerminateInstances(terminateInstanceParams)

	// TODO: check if group_type is GroupASG, then terminate the ASG instead of single instances.

	Logger.Info("TerminateInstances", "response", terminateInstanceResponse)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.

		if strings.Contains(err.Error(), "InvalidInstanceID.NotFound") {
			// TODO: replace this with something shorter

			task.Lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve

			// we don't know when it has been terminated, so just use the current time
			now := time.Now().UTC()
			task.Lease.TerminatedAt = &now

			// TODO: use the ufficial time of termination, from th sqs message, because if erminated via link, the termination time is not expiresAt
			// lease.TerminatedAt = time.Now().UTC()
			s.DB.Save(&task.Lease)

			Logger.Debug(
				"TerminatorQueueConsumer TerminateInstances ",
				"err", err,
				"action_taken", "removing lease of already deleted/non-existent instance from DB",
			)

		} else {
			// TODO: cleaner way to do this?  cloudaccount.Account would be nice .. gorma provides this
			var account models.Account
			s.DB.First(&account, cloudaccount.AccountID)

			recipientEmail := account.Email

			s.sendMisconfigurationNotice(err, recipientEmail)
		}

		return err
	}

	return nil
}
