// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tleyden/cecil/notification"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
)

// ExtenderQueueConsumer consumes ExtenderTask from ExtenderQueue; approves or extends leases.
func (s *Service) ExtenderQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(tasks.ExtenderTask)
	// TODO: check whether fields are non-null and valid

	if task.Approving {
		Logger.Info(
			"Approving lease",

			"lease_id", task.Lease.ID,
			"group_type", task.Lease.GroupType.String(),
			"group_uid", task.Lease.GroupUID,
		)
	} else {
		Logger.Info(
			"Extending lease",

			"lease_id", task.Lease.ID,
			"group_type", task.Lease.GroupType.String(),
			"group_uid", task.Lease.GroupUID,
		)
	}

	if task.Lease.IsExpired() {
		// TODO: should the user be notified that the lease cannot be extended because it is already expired???
		err := errors.New("lease is already expired; cannot extend/approve")
		Logger.Error(
			"error while extendin lease",
			"lease_id", task.Lease.ID,
			"group_type", task.Lease.GroupType.String(),
			"group_uid", task.Lease.GroupUID,
			"err", err,
		)
		return err
	}

	task.Lease.TokenOnce = uuid.NewV4().String() // invalidates all other URLs to renew/terminate/approve
	task.Lease.NumTimesAllertedAboutExpiry = NoAlertsSent

	// define the lease duration
	leaseDuration, err := s.DefineLeaseDuration(task.Lease.AccountID, task.Lease.CloudaccountID)
	if err != nil {
		Logger.Error(
			"error while DefineLeaseDuration",

			"lease_id", task.Lease.ID,
			"group_type", task.Lease.GroupType.String(),
			"group_uid", task.Lease.GroupUID,
			"err", err,
		)
		return err
	}
	// TODO: remove leaseDuration
	_ = leaseDuration

	if task.Approving {
		now := time.Now().UTC()
		task.Lease.ApprovedAt = &now
		//task.Lease.ExpiresAt = task.Lease.CreatedAt.Add(leaseDuration)
	} else {
		leaseDuration := task.Lease.ExpiresAt.Sub(task.Lease.CreatedAt)
		task.Lease.ExpiresAt = task.Lease.ExpiresAt.Add(leaseDuration)
	}

	s.DB.Save(&task.Lease)

	owner, err := s.GetOwnerByID(task.Lease.OwnerID)
	if err != nil {
		return err
	}

	var newEmailBody string
	var newEmailSubject string
	var notificationType notification.NotificationType

	var AWSResourceID string
	AWSResourceID = task.Lease.GroupUID

	var emailValues = map[string]interface{}{
		"owner_email":     owner.Email,
		"resource_region": task.Lease.Region,

		"lease_duration": task.Lease.ExpiresAt.Sub(task.Lease.CreatedAt).String(),

		"expires_at": task.Lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
	}

	emailValues["lease_id"] = task.Lease.ID
	emailValues["group_type"] = task.Lease.GroupType.String()
	emailValues["group_uid"] = task.Lease.GroupUID
	if task.Lease.AwsContainerName != "" {
		emailValues["aws_container_name"] = task.Lease.AwsContainerName
	}
	{
		instances, err := s.ActiveInstancesForGroup(task.Lease.AccountID, &task.Lease.CloudaccountID, task.Lease.GroupUID)
		if err != nil {
			return err
		}
		emailValues["instances"] = instances
	}

	if task.Approving {
		notificationType = notification.LeaseApproved

		newEmailSubject = fmt.Sprintf("%v Lease %v approved", task.Lease.GroupType.EmailDisplayString(), task.Lease.ID)

		newEmailBody, err = tools.CompileEmailTemplate(
			"lease-approved.html",
			emailValues,
		)
		if err != nil {
			return err
		}
	} else {
		notificationType = notification.LeaseExtended

		newEmailSubject = fmt.Sprintf("%v Lease %v extended", task.Lease.GroupType.EmailDisplayString(), task.Lease.ID)

		newEmailBody, err = tools.CompileEmailTemplate(
			"lease-extended.html",
			emailValues,
		)
		if err != nil {
			return err
		}
	}

	return s.Queues().NotifierQueue().PushTask(tasks.NotifierTask{
		AccountID: task.Lease.AccountID, // this will also trigger send to Slack
		To:        owner.Email,
		Subject:   newEmailSubject,
		BodyHTML:  newEmailBody,
		BodyText:  newEmailBody,
		NotificationMeta: notification.NotificationMeta{
			NotificationType: notificationType,
			LeaseUUID:        task.Lease.UUID,
			AWSResourceID:    AWSResourceID,
			// ResourceType:     task.Lease.ResourceType,
		},
	})
}
