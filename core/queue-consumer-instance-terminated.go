package core

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/notification"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
	"github.com/tleyden/cecil/transmission"
)

func ParseInstanceTerminatedTask(t interface{}) (*transmission.Transmission, error) {
	if t == nil {
		return nil, errors.New("t is nil")
	}

	task, ok := t.(tasks.InstanceTerminatedTask)
	if !ok {
		return nil, errors.New("t is not tasks.InstanceTerminatedTask")
	}

	tr, ok := task.Transmission.(*transmission.Transmission)
	if !ok {
		return nil, errors.New("t.Transmission is not *transmission.Transmission")
	}

	return tr, nil
}

// InstanceTerminatedQueueConsumer consumes InstanceTerminatedTask from InstanceTerminatedQueue;
// marks leases as terminated and notifes the owner.
func (s *Service) InstanceTerminatedQueueConsumer(t interface{}) error {

	tr, err := ParseInstanceTerminatedTask(t)
	if err != nil {
		return err
	}

	Logger.Info(
		"InstanceTerminatedQueueConsumer called",
		"transmission", tr,
	)
	defer Logger.Info(
		"InstanceTerminatedQueueConsumer call finished",
		"transmission", tr,
	)

	// check whether the group has a lease

	lease, err := tr.LeaseByInstanceID()
	if err != nil {
		Logger.Error("error while LeaseByInstanceID", "err", err)
		if err == gorm.ErrRecordNotFound {
			Logger.Info("Delete SQS Message")
			err = tr.DeleteMessage()
			if err != nil {
				Logger.Warn("DeleteMessage", "err", err)
			}
			return err
		}
		return err
	}

	groupHasLease := lease != nil
	if !groupHasLease {
		Logger.Warn("instance terminated: group hasn't a lease on it", "transmission", tr)
		Logger.Info("Delete SQS Message")
		err = tr.DeleteMessage()
		if err != nil {
			Logger.Warn("DeleteMessage", "err", err)
		}
		return err
	}

	ins, err := tr.InstanceIsNew()
	if err != nil {
		Logger.Error("error while InstanceIsNew", "err", err)
		return err
	}
	instanceIsRegistered := ins != nil
	if !instanceIsRegistered {
		Logger.Warn("instance not registered; cannot mark as deleted", "transmission", tr)
		return nil
	}

	Logger.Info("marking instance as deleted", "instance", ins)
	// TODO: use the ufficial time of termination, from th sqs message, because if erminated via link, the termination time is not expiresAt
	s.DB.Save(ins.MarkAsTerminated(&tr.Message.Time))

	instances, err := tr.ActiveInstancesForGroup(lease.GroupUID)
	if err != nil {
		return err
	}

	thisWasTheLastInstance := len(instances) == 0
	if !thisWasTheLastInstance {
		// exit, because there are other instances for this lease
		Logger.Info(
			"This lease has other instances running; not terminating the lease.",
			"lease.GroupUID", lease.GroupUID,
			"remaining", len(instances),
		)
		return nil
	}
	// check whether the instance is registered
	// determine group
	// set terminated_at of instance
	// count existing instances of groupUID
	// if == 0, lease is terminated

	// mark lease as terminated
	s.DB.Save(lease.MarkAsTerminated(ins.TerminatedAt))

	owner, err := s.GetOwnerByID(lease.OwnerID)

	if err != nil {
		Logger.Warn("InstanceTerminatedQueueConsumer: error fetching owner", "err", err)
		return fmt.Errorf("InstanceTerminatedQueueConsumer: error fetching owner: %v", err)
	}

	var emailValues = map[string]interface{}{
		"owner_email":     owner.Email,
		"resource_region": lease.Region,

		"lease_duration": ins.TerminatedAt.Sub(lease.CreatedAt).String(),
		"expires_at":     lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
		"terminated_at":  ins.TerminatedAt.Format("2006-01-02 15:04:05 GMT"),
		"created_at":     ins.CreatedAt.Format("2006-01-02 15:04:05 GMT"),
	}

	emailValues["lease_id"] = lease.ID
	emailValues["group_type"] = lease.GroupType.String()
	emailValues["group_uid"] = lease.GroupUID
	if lease.AwsContainerName != "" {
		emailValues["aws_container_name"] = lease.AwsContainerName
	}
	{
		instances, err := s.ActiveInstancesForGroup(lease.AccountID, &lease.CloudaccountID, lease.GroupUID)
		if err != nil {
			return err
		}
		emailValues["instances"] = instances
	}

	newEmailBody, err := tools.CompileEmailTemplate(
		"lease-terminated.html",
		emailValues,
	)
	if err != nil {
		return err
	}

	newEmailSubject := fmt.Sprintf("Lease %v (type %v) terminated", lease.ID, lease.GroupType.String())

	s.Queues().NotifierQueue().PushTask(tasks.NotifierTask{
		AccountID: lease.AccountID, // this will also trigger send to Slack
		To:        owner.Email,
		Subject:   newEmailSubject,
		BodyHTML:  newEmailBody,
		BodyText:  newEmailBody,
		NotificationMeta: notification.NotificationMeta{
			NotificationType: notification.InstanceTerminated,
			LeaseUUID:        lease.UUID,
			//AWSResourceID:    task.AWSResourceID,
			//ResourceType: lease.ResourceType,
		},
	})

	return nil
}
