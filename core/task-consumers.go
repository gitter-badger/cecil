package core

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gopkg.in/mailgun/mailgun-go.v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
)

// @@@@@@@@@@@@@@@ Task consumers @@@@@@@@@@@@@@@

func (s *Service) NewLeaseQueueConsumer(t interface{}) error {

	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	transmission := t.(NewLeaseTask).Transmission

	//check whether someone with this aws adminAccount id is registered at zerocloud
	err := transmission.FetchCloudAccount()
	if err != nil {
		// TODO: notify admin; something fishy is going on.
		logger.Warn("originator is not registered", "AWSID", transmission.Topic.AWSID)
		return nil // TODO: return an error ???
	}

	// check whether the cloud account has an admin account
	err = transmission.FetchAdminAccount()
	if err != nil {
		// TODO: notify admin; something fishy is going on.
		logger.Warn("Error while retrieving admin account", "error", err)
		return nil // TODO: return an error ???
	}

	logger.Info("adminAccount",
		"adminAccount", transmission.AdminAccount,
	)

	logger.Info("Creating AssumedConfig", "topicRegion", transmission.Topic.Region, "topicAWSID", transmission.Topic.AWSID, "externalID", transmission.CloudAccount.ExternalID)

	err = transmission.CreateAssumedService()
	if err != nil {
		// TODO: this might reveal too much to the admin about zerocloud; be selective and cautious
		s.sendMisconfigurationNotice(err, transmission.AdminAccount.Email)
		logger.Warn("error while creating assumed service", "error", err)
		return nil // TODO: return an error ???
	}

	err = transmission.CreateAssumedEC2Service()
	if err != nil {
		// TODO: this might reveal too much to the admin about zerocloud; be selective and cautious
		s.sendMisconfigurationNotice(err, transmission.AdminAccount.Email)
		logger.Warn("error while creating ec2 service with assumed service", "error", err)
		return nil // TODO: return an error ???
	}

	err = transmission.DescribeInstance()
	if err != nil {
		// TODO: this might reveal too much to the admin about zerocloud; be selective and cautious
		s.sendMisconfigurationNotice(err, transmission.AdminAccount.Email)
		logger.Warn("error while describing instances", "error", err)
		return nil // TODO: return an error ???
	}

	// check whether the instance specified in the event exists on aws
	if !transmission.InstanceExists() {
		logger.Warn("Instance does not exist", "instanceID", transmission.Message.Detail.InstanceID)
		// remove message from queue
		err := transmission.DeleteMessage()
		if err != nil {
			logger.Warn("DeleteMessage", "error", err)
		}
		return nil // TODO: return an error ???
	}

	logger.Info("describeInstances", "response", transmission.describeInstancesResponse)

	err = transmission.FetchInstance()
	if err != nil {
		logger.Warn("error while fetching instance description", "error", err)
		return nil // TODO: return an error ???
	}

	err = transmission.ComputeInstanceRegion()
	if err != nil {
		logger.Warn("error while computing instance region", "error", err)
		return nil // TODO: return an error ???
	}

	/// transmission.Message.InstanceID == transmission.Instance.InstanceID
	//TODO: might this happen?

	/// transmission.Instance.IsTerminated()
	/// transmission.Message.Delete()

	// if the message signal that an instance has been terminated, create a task
	// to mark the lease as terminated
	if transmission.InstanceIsTerminated() {

		s.LeaseTerminatedQueue.TaskQueue <- LeaseTerminatedTask{
			AWSID:      transmission.CloudAccount.AWSID,
			InstanceID: *transmission.Instance.InstanceId,
		}

		// remove message from queue
		err := transmission.DeleteMessage()
		if err != nil {
			logger.Warn("DeleteMessage", "error", err)
		}
		return nil // TODO: return an error ???
	}

	// do not consider states other than pending and terminated
	if !transmission.InstanceIsPendingOrRunning() {
		logger.Warn("The retrieved state is neither pending nor running:", "state", transmission.Instance.State.Name)
		// remove message from queue
		// remove message from queue
		err := transmission.DeleteMessage()
		if err != nil {
			logger.Warn("DeleteMessage", "error", err)
		}
		return nil // TODO: return an error ???
	}

	if !transmission.LeaseIsNew() {
		// TODO: notify admin
		logger.Warn("instanceCount != 0")
		return nil // TODO: return an error ???
	}

	if !transmission.InstanceHasGoodOwnerTag() || !transmission.ExternalOwnerIsWhitelisted() {
		// assign instance to admin, and send notification to admin
		// owner is not whitelisted: notify admin: "Warning: zerocloudowner tag email not in whitelist"

		err := transmission.SetAdminAsOwner()
		if err != nil {
			logger.Warn("Error while setting admin as owner", "error", err)
		}

		transmission.leaseDuration = time.Duration(ZCDefaultLeaseApprovalTimeoutDuration)
		var expiresAt = time.Now().UTC().Add(transmission.leaseDuration)

		// these will be used to compose the urls and verify the requests
		lease_uuid := uuid.NewV4().String()
		instance_id := *transmission.Instance.InstanceId
		token_once := uuid.NewV4().String() // one-time token

		newLease := Lease{
			UUID:      lease_uuid,
			TokenOnce: token_once,

			OwnerID:        transmission.owner.ID,
			CloudAccountID: transmission.CloudAccount.ID,
			AWSAccountID:   transmission.CloudAccount.AWSID,

			InstanceID:       *transmission.Instance.InstanceId,
			Region:           transmission.instanceRegion,
			AvailabilityZone: *transmission.Instance.Placement.AvailabilityZone,
			InstanceType:     *transmission.Instance.InstanceType,

			// Terminated bool `sql:"DEFAULT:false"`
			// Deleted    bool `sql:"DEFAULT:false"`

			LaunchedAt: transmission.Instance.LaunchTime.UTC(),
			ExpiresAt:  expiresAt,
		}
		s.DB.Create(&newLease)
		logger.Info("new lease created",
			"lease", newLease,
		)

		var newEmailBody string

		// URL to approve lease
		action := "approve"
		signature, err := s.sign(lease_uuid, instance_id, action, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while signing")
		}
		approve_url := fmt.Sprintf("http://0.0.0.0:8080/cmd/leases/%s/%s/%s?t=%s&s=%s",
			lease_uuid,
			instance_id,
			action,
			token_once,
			base64.URLEncoding.EncodeToString(signature),
		)

		// URL to terminate lease
		action = "terminate"
		signature, err = s.sign(lease_uuid, instance_id, action, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while signing")
		}
		terminate_url := fmt.Sprintf("http://0.0.0.0:8080/cmd/leases/%s/%s/%s?t=%s&s=%s",
			lease_uuid,
			instance_id,
			action,
			token_once,
			base64.URLEncoding.EncodeToString(signature),
		)

		switch {
		case !transmission.InstanceHasGoodOwnerTag():
			newEmailBody = compileEmail(
				`Hey {{.owner_email}}, someone created a new instance 
				(id <b>{{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). <br><br>

				It does not have a valid ZeroCloudOwner tag, so we assigned it to you.
				
				<br>
				<br>
				
				It will be terminated at <b>{{.termination_time}}</b> ({{.instance_duration}} after it's creation).

				<br>
				<br>
				
				Terminate immediately:
				<br>
				<br>
				{{.instance_terminate_url}}

				Approve (you will be the owner):
				<br>
				<br>
				{{.instance_approve_url}}

				<br>
				<br>
				Thanks for using ZeroCloud!
				`,

				map[string]interface{}{
					"owner_email":     transmission.owner.Email,
					"instance_id":     *transmission.Instance.InstanceId,
					"instance_type":   *transmission.Instance.InstanceType,
					"instance_region": transmission.instanceRegion,

					"termination_time":  expiresAt.Format("2006-01-02 15:04:05 GMT"),
					"instance_duration": transmission.leaseDuration.String(),

					"instance_terminate_url": terminate_url,
					"instance_approve_url":   approve_url,
				},
			)
			break

		case !transmission.ExternalOwnerIsWhitelisted():
			newEmailBody = compileEmail(
				`Hey {{.owner_email}}, someone created a new instance 
				(id <b>{{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). <br><br>

				The ZeroCloudOwner tag of this instance is not in the whitelist, so we assigned it to you.
				
				<br>
				<br>
				
				It will be terminated at <b>{{.termination_time}}</b> ({{.instance_duration}} after it's creation).

				<br>
				<br>
				
				Terminate immediately:
				<br>
				<br>
				{{.instance_terminate_url}}

				Approve (you will be the owner):
				<br>
				<br>
				{{.instance_approve_url}}

				<br>
				<br>
				Thanks for using ZeroCloud!
				`,

				map[string]interface{}{
					"owner_email":     transmission.owner.Email,
					"instance_id":     *transmission.Instance.InstanceId,
					"instance_type":   *transmission.Instance.InstanceType,
					"instance_region": transmission.instanceRegion,

					"termination_time":  expiresAt.Format("2006-01-02 15:04:05 GMT"),
					"instance_duration": transmission.leaseDuration.String(),

					"instance_terminate_url": terminate_url,
					"instance_approve_url":   approve_url,
				},
			)
		}

		s.NotifierQueue.TaskQueue <- NotifierTask{
			//To:       owner.Email,
			From:     ZCMailerFromAddress,
			To:       transmission.AdminAccount.Email,
			Subject:  fmt.Sprintf("Instance (%v) needs attention", *transmission.Instance.InstanceId),
			BodyHTML: newEmailBody,
			BodyText: newEmailBody,
		}

		err = transmission.DeleteMessage()
		if err != nil {
			logger.Warn("DeleteMessage", "error", err)
		}
		return nil // TODO: return an error ???
	}

	err = transmission.SetExternalOwnerAsOwner()
	if err != nil {
		logger.Warn("Error while setting external owner as owner", "error", err)
	}

	if transmission.LeaseNeedsApproval() {
		// register new lease in DB
		// expiry: 1h
		// send confirmation to owner: confirmation link, and termination link

		transmission.leaseDuration = time.Duration(ZCDefaultLeaseApprovalTimeoutDuration)
		var expiresAt = time.Now().UTC().Add(transmission.leaseDuration)

		// these will be used to compose the urls and verify the requests
		lease_uuid := uuid.NewV4().String()
		instance_id := *transmission.Instance.InstanceId
		token_once := uuid.NewV4().String() // one-time token

		newLease := Lease{
			UUID:      lease_uuid,
			TokenOnce: token_once,

			OwnerID:        transmission.owner.ID,
			CloudAccountID: transmission.CloudAccount.ID,
			AWSAccountID:   transmission.CloudAccount.AWSID,

			InstanceID:       *transmission.Instance.InstanceId,
			Region:           transmission.instanceRegion,
			AvailabilityZone: *transmission.Instance.Placement.AvailabilityZone,

			// Terminated bool `sql:"DEFAULT:false"`
			// Deleted    bool `sql:"DEFAULT:false"`

			LaunchedAt:   transmission.Instance.LaunchTime.UTC(),
			ExpiresAt:    expiresAt,
			InstanceType: *transmission.Instance.InstanceType,
		}
		s.DB.Create(&newLease)
		logger.Info("new lease created",
			"lease", newLease,
		)

		// URL to approve lease
		action := "approve"
		signature, err := s.sign(lease_uuid, instance_id, action, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while signing")
		}
		approve_url := fmt.Sprintf("http://0.0.0.0:8080/cmd/leases/%s/%s/%s?t=%s&s=%s",
			lease_uuid,
			instance_id,
			action,
			token_once,
			base64.URLEncoding.EncodeToString(signature),
		)

		// URL to terminate lease
		action = "terminate"
		signature, err = s.sign(lease_uuid, instance_id, action, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while signing")
		}
		terminate_url := fmt.Sprintf("http://0.0.0.0:8080/cmd/leases/%s/%s/%s?t=%s&s=%s",
			lease_uuid,
			instance_id,
			action,
			token_once,
			base64.URLEncoding.EncodeToString(signature),
		)

		newEmailBody := compileEmail(
			`Hey {{.owner_email}}, you (or someone else using your ZeroCloudOwner tag) created a new instance 
				(id <b>{{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). <br><br>

				At the time of writing this email, you have {{.n_of_active_leases}} active
					leases, so we need your approval for this one. <br><br>

				Please click on "Approve" to approve this instance,
					otherwise it will be terminated at <b>{{.termination_time}}</b> ({{.instance_duration}} after it's creation).

				<br>
				<br>

				Approve:
				<br>
				<br>
				{{.instance_approve_url}}

				<br>
				<br>
				
				Terminate immediately:
				<br>
				<br>
				{{.instance_terminate_url}}
				
				<br>
				<br>
				Thanks for using ZeroCloud!
				`,

			map[string]interface{}{
				"owner_email":        transmission.owner.Email,
				"n_of_active_leases": transmission.activeLeaseCount,
				"instance_id":        *transmission.Instance.InstanceId,
				"instance_type":      *transmission.Instance.InstanceType,
				"instance_region":    transmission.instanceRegion,

				"termination_time":  expiresAt.Format("2006-01-02 15:04:05 GMT"),
				"instance_duration": transmission.leaseDuration.String(),

				"instance_approve_url":   approve_url,
				"instance_terminate_url": terminate_url,
			},
		)
		s.NotifierQueue.TaskQueue <- NotifierTask{
			From:     ZCMailerFromAddress,
			To:       transmission.owner.Email,
			Subject:  fmt.Sprintf("Instance (%v) needs approval", *transmission.Instance.InstanceId),
			BodyHTML: newEmailBody,
			BodyText: newEmailBody,
		}

		// remove message from queue
		err = transmission.DeleteMessage()
		if err != nil {
			logger.Warn("DeleteMessage", "error", err)
		}
		return nil // TODO: return an error ???
	} else {
		// register new lease in DB
		// set its expiration to zone.default_expiration (if > 0), or cloudAccount.default_expiration, or adminAccount.default_expiration

		transmission.DefineLeaseDuration()
		var expiresAt = time.Now().UTC().Add(transmission.leaseDuration)

		// these will be used to compose the urls and verify the requests
		lease_uuid := uuid.NewV4().String()
		instance_id := *transmission.Instance.InstanceId
		token_once := uuid.NewV4().String() // one-time token

		newLease := Lease{
			UUID:      lease_uuid,
			TokenOnce: token_once,

			OwnerID:        transmission.owner.ID,
			CloudAccountID: transmission.CloudAccount.ID,
			AWSAccountID:   transmission.CloudAccount.AWSID,

			InstanceID:       *transmission.Instance.InstanceId,
			Region:           transmission.instanceRegion,
			AvailabilityZone: *transmission.Instance.Placement.AvailabilityZone,

			// Terminated bool `sql:"DEFAULT:false"`
			// Deleted    bool `sql:"DEFAULT:false"`

			LaunchedAt:   transmission.Instance.LaunchTime.UTC(),
			ExpiresAt:    expiresAt,
			InstanceType: *transmission.Instance.InstanceType,
		}
		s.DB.Create(&newLease)
		logger.Info("new lease created",
			"lease", newLease,
		)

		// URL to terminate lease
		action := "terminate"
		signature, err := s.sign(lease_uuid, instance_id, action, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while signing")
		}
		terminate_url := fmt.Sprintf("http://0.0.0.0:8080/cmd/leases/%s/%s/%s?t=%s&s=%s",
			lease_uuid,
			instance_id,
			action,
			token_once,
			base64.URLEncoding.EncodeToString(signature),
		)

		newEmailBody := compileEmail(
			`Hey {{.owner_email}}, you (or someone else using your ZeroCloudOwner tag) created a new instance 
				(id <b>{{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). That's AWESOME!

				<br>
				<br>

				Your instance will be terminated at <b>{{.termination_time}}</b> ({{.instance_duration}} after it's creation).

				<br>
				<br>

				Terminate immediately:
				<br>
				<br>
				{{.instance_terminate_url}}

				<br>
				<br>
				
				Thanks for using ZeroCloud!
				`,

			map[string]interface{}{
				"owner_email":     transmission.owner.Email,
				"instance_id":     *transmission.Instance.InstanceId,
				"instance_type":   *transmission.Instance.InstanceType,
				"instance_region": transmission.instanceRegion,

				"termination_time":  expiresAt.Format("2006-01-02 15:04:05 GMT"),
				"instance_duration": transmission.leaseDuration.String(),

				"instance_terminate_url": terminate_url,
			},
		)
		s.NotifierQueue.TaskQueue <- NotifierTask{
			From:     ZCMailerFromAddress,
			To:       transmission.owner.Email,
			Subject:  fmt.Sprintf("Instance (%v) created", *transmission.Instance.InstanceId),
			BodyHTML: newEmailBody,
			BodyText: newEmailBody,
		}

		// remove message from queue
		err = transmission.DeleteMessage()
		if err != nil {
			logger.Warn("DeleteMessage", "error", err)
		}
		return nil // TODO: return an error ???
	}

	// if message.Detail.State == ec2.InstanceStateNameTerminated
	// LeaseTerminatedQueue <- LeaseTerminatedTask{} and continue

	// get zc adminAccount who has a cloudaccount with awsID == topicAWSID
	// if no one of our customers owns this adminAccount, error
	// fetch options config
	// roleARN := fmt.Sprintf("arn:aws:iam::%v:role/ZeroCloudRole",topicAWSID)
	// assume role
	// fetch instance info
	// check if statuses match (this message was sent by aws.ec2)
	// message.Detail.InstanceID

	// fmt.Printf("%v", message)
	return nil
}

func (s *Service) TerminatorQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(TerminatorTask)
	// TODO: check whether fields are non-null and valid
	logger.Info("TerminatorQueueConsumer",
		"task", task,
	)

	// need:
	// region
	// roleARN
	// external ID

	var cloudAccount CloudAccount
	var leaseCloudOwnerCount int64
	s.DB.Model(&task.Lease).Related(&cloudAccount).Count(&leaseCloudOwnerCount)
	//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&leaseCloudOwnerCount)
	if leaseCloudOwnerCount == 0 {
		// TODO: notify admin; something fishy is going on.
		logger.Warn("leaseCloudOwnerCount == 0")
		return fmt.Errorf("leaseCloudOwnerCount == 0")
	}

	// assume role
	assumedConfig := &aws.Config{
		Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
			Client: sts.New(s.AWS.Session, &aws.Config{Region: aws.String(task.Region)}),
			RoleARN: fmt.Sprintf(
				"arn:aws:iam::%v:role/%v",
				cloudAccount.AWSID,
				viper.GetString("ForeignRoleName"),
			),
			RoleSessionName: uuid.NewV4().String(),
			ExternalID:      aws.String(cloudAccount.ExternalID),
			ExpiryWindow:    3 * time.Minute,
		}),
	}

	assumedService := session.New(assumedConfig)

	ec2Service := s.EC2(assumedService, task.Region)

	terminateInstanceParams := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{ // Required
			aws.String(task.InstanceID),
		},
	}
	terminateInstanceResponse, err := ec2Service.TerminateInstances(terminateInstanceParams)
	_ = terminateInstanceResponse

	logger.Info("TerminateInstances", "response", terminateInstanceResponse)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.

		// TODO: cleaner way to do this?  cloudAccount.Account would be nice .. gorma provides this
		var account Account
		s.DB.First(&account, cloudAccount.AccountID)

		recipientEmail := account.Email

		s.sendMisconfigurationNotice(err, recipientEmail)
		return err
	}

	return nil
}

func (s *Service) LeaseTerminatedQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(LeaseTerminatedTask)
	logger.Info("Marking lease as terminated",
		"InstanceID", task.InstanceID,
	)

	var lease Lease
	var leasesFound int64
	s.DB.Table("leases").Where(&Lease{InstanceID: task.InstanceID, AWSAccountID: task.AWSID}).First(&lease).Count(&leasesFound)

	if leasesFound != 1 {
		logger.Warn("Found multiple leases for deletion", "count", leasesFound)
		return fmt.Errorf("Found multiple leases for deletion", "count", leasesFound)
	}

	lease.Terminated = true
	lease.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve

	// TODO: use the ufficial time of termination, from th sqs message
	// lease.TerminatedAt = time.Now().UTC()
	s.DB.Save(&lease)

	var owner Owner
	var ownerCount int64

	s.DB.Table("owners").Where(lease.OwnerID).First(&owner).Count(&ownerCount)

	newEmailBody := compileEmail(
		`Hey {{.owner_email}}, instance with id <b>{{.instance_id}}</b>
				(of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>) has been terminated at 
				<b>{{.terminated_at}}</b> ({{.instance_duration}} after it's creation)

				<br>
				<br>
				
				Thanks for using ZeroCloud!
				`,

		map[string]interface{}{
			"owner_email":     owner.Email,
			"instance_id":     lease.InstanceID,
			"instance_type":   lease.InstanceType,
			"instance_region": lease.Region,

			"instance_duration": lease.ExpiresAt.Sub(lease.CreatedAt).String(),

			"terminated_at": lease.ExpiresAt.Format("2006-01-02 15:04:05 GMT"),
		},
	)
	s.NotifierQueue.TaskQueue <- NotifierTask{
		From:     ZCMailerFromAddress,
		To:       owner.Email,
		Subject:  fmt.Sprintf("Instance (%v) terminated", lease.InstanceID),
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
	}

	return nil
}

func (s *Service) ExtenderQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(ExtenderTask)
	// TODO: check whether fields are non-null and valid

	_ = task

	return nil
}

func (s *Service) NotifierQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(NotifierTask)
	// TODO: check whether fields are non-null and valid
	logger.Info("Sending EMAIL",
		"to", task.To,
	)

	message := mailgun.NewMessage(
		task.From,
		task.Subject,
		task.BodyText,
		task.To,
	)

	//message.SetTracking(true)
	//message.SetDeliveryTime(time.Now().Add(24 * time.Hour))
	message.SetHtml(task.BodyHTML)
	_, id, err := s.Mailer.Send(message)
	if err != nil {
		logger.Error("Error while sending email", "error", err)
		return err
	}
	_ = id

	return nil
}
