// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package core

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tleyden/cecil/awstools"
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

		asgName := task.Lease.AwsContainerName // AwsContainerName is the name of the AutoScalingGroup
		assumedAutoScalingService := s.AWSServices().AutoScaling(assumedService, task.Region)

		// Set autoscaling group capacity to zero instances.
		{
			Logger.Info(
				"UpdateAutoScalingGroup to 0 on ASG",
				"asg_name", asgName,
				"lease_id", task.Lease.ID,
				"cloudaccount_id", task.CloudaccountID,
				"account_id", task.AccountID,
			)

			////
			params := &autoscaling.UpdateAutoScalingGroupInput{
				AutoScalingGroupName: aws.String(asgName), // Required
				DesiredCapacity:      aws.Int64(0),
				MinSize:              aws.Int64(0),
			}
			////
			resp, err := assumedAutoScalingService.UpdateAutoScalingGroup(params)
			if err != nil {
				Logger.Error(
					"error while UpdateAutoScalingGroup",
					"asg_name", asgName,
					"lease_id", task.Lease.ID,
					"cloudaccount_id", task.CloudaccountID,
					"account_id", task.AccountID,
					"err", err,
				)

				if awstools.IsErrNotFoundASG(err) {
					// we don't know when it has been terminated, so just use the current time
					s.DB.Save(task.Lease.MarkAsTerminated(nil))

					Logger.Debug(
						"TerminatorQueueConsumer TerminateInstances ",
						"err", err,
						"action_taken", "removing lease of already deleted/non-existent ASG from DB",
					)

					return nil
				}

				// TODO: cleaner way to do this?  cloudaccount.Account would be nice .. gorma provides this
				var account models.Account
				s.DB.First(&account, cloudaccount.AccountID)

				recipientEmail := account.Email

				s.sendMisconfigurationNotice(err, recipientEmail)
				return err
			}
			Logger.Info("UpdateAutoScalingGroup response", "resp", resp.GoString())
		}

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

	if len(instanceIDs) == 0 {
		Logger.Info(
			"Lease has no active instances",
			"group_uid", task.GroupUID,
			"group_type", task.GroupType.String(),
			"lease_id", task.ID,
			"active_instances", len(instanceIDs),
			"instanceIDs", instanceIDs,
		)
		return nil
	}

	Logger.Info(
		"Terminating lease",
		"group_uid", task.GroupUID,
		"group_type", task.GroupType.String(),
		"lease_id", task.ID,
		"active_instances", len(instanceIDs),
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

		if awstools.IsErrNotFoundInstance(err) {

			// we don't know when it has been terminated, so just use the current time
			s.DB.Save(task.Lease.MarkAsTerminated(nil))

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
