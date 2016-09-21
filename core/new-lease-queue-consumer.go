package core

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

// NewLeaseQueueConsumer consumes NewLeaseTask from NewLeaseQueue
func (s *Service) NewLeaseQueueConsumer(t interface{}) error {

	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	transmission := t.(NewLeaseTask).Transmission

	logger.Info("NewLeaseQueueConsumer called", "transmission", transmission)
	defer logger.Info("NewLeaseQueueConsumer call finished", "transmission", transmission)

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
			InstanceID: transmission.InstanceId(),
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
		logger.Info("Transmission doesn't have owner tag or owner is not whitelisted.")

		err := transmission.SetAdminAsOwner()
		if err != nil {
			logger.Warn("Error while setting admin as owner", "error", err)
		}

		transmission.leaseDuration = time.Duration(ZCDefaultLeaseApprovalTimeoutDuration)
		var expiresAt = time.Now().UTC().Add(transmission.leaseDuration)

		// these will be used to compose the urls and verify the requests
		lease_uuid := uuid.NewV4().String()
		instance_id := transmission.InstanceId()
		token_once := uuid.NewV4().String() // one-time token

		newLease := Lease{
			UUID:      lease_uuid,
			TokenOnce: token_once,

			OwnerID:        transmission.owner.ID,
			CloudAccountID: transmission.CloudAccount.ID,
			AWSAccountID:   transmission.CloudAccount.AWSID,

			InstanceID:       transmission.InstanceId(),
			Region:           transmission.instanceRegion,
			AvailabilityZone: transmission.AvailabilityZone(),
			InstanceType:     transmission.InstanceType(),

			// Terminated bool `sql:"DEFAULT:false"`
			// Deleted    bool `sql:"DEFAULT:false"`
			Alerted: true,

			LaunchedAt: transmission.InstanceLaunchTimeUTC(),
			ExpiresAt:  expiresAt,
		}
		s.DB.Create(&newLease)
		logger.Info("new lease created",
			"lease", newLease,
		)

		var newEmailBody string

		// URL to approve lease
		logger.Info(
			"Creating lease signature",
			"lease_uuid", "instance_id", "action", "token_once",
			lease_uuid, instance_id, "approve", token_once,
		)
		approve_url, err := s.generateSignedEmailActionURL("approve", lease_uuid, instance_id, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		// URL to terminate lease
		terminate_url, err := s.generateSignedEmailActionURL("terminate", lease_uuid, instance_id, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		switch {
		case !transmission.InstanceHasGoodOwnerTag():
			newEmailBody = compileEmail(
				`Hey {{.owner_email}}, someone created a new instance 
				(id <b>{{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). <br><br>

				It does not have a valid ZeroCloudOwner tag, so we assigned it to you (the admin).
				
				<br>
				<br>
				
				If not approved, it will be terminated at <b>{{.termination_time}}</b> ({{.instance_duration}} after it's creation).

				<br>
				<br>
				
				Terminate immediately:
				<br>
				<br>
				<a href="{{.instance_terminate_url}}" target="_blank">Click here to <b>terminate</b></a>

				<br>
				<br>

				Approve (you will be the owner):
				<br>
				<br>
				<a href="{{.instance_approve_url}}" target="_blank">Click here to <b>approve</b></a>

				<br>
				<br>
				Thanks for using ZeroCloud!
				`,

				map[string]interface{}{
					"owner_email":     transmission.owner.Email,
					"instance_id":     transmission.InstanceId(),
					"instance_type":   transmission.InstanceType(),
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

				The ZeroCloudOwner tag of this instance is not in the whitelist, so we assigned it to you (the admin).
				
				<br>
				<br>
				
				If not approved, it will be terminated at <b>{{.termination_time}}</b> ({{.instance_duration}} after it's creation).

				<br>
				<br>
				
				Terminate immediately:
				<br>
				<br>
				<a href="{{.instance_terminate_url}}" target="_blank">Click here to <b>terminate</b></a>

				<br>
				<br>

				Approve (you will be the owner):
				<br>
				<br>
				<a href="{{.instance_approve_url}}" target="_blank">Click here to <b>approve</b></a>

				<br>
				<br>
				Thanks for using ZeroCloud!
				`,

				map[string]interface{}{
					"owner_email":     transmission.owner.Email,
					"instance_id":     transmission.InstanceId(),
					"instance_type":   transmission.InstanceType(),
					"instance_region": transmission.instanceRegion,

					"termination_time":  expiresAt.Format("2006-01-02 15:04:05 GMT"),
					"instance_duration": transmission.leaseDuration.String(),

					"instance_terminate_url": terminate_url,
					"instance_approve_url":   approve_url,
				},
			)
		}

		logger.Info("Adding new NotifierTask")
		s.NotifierQueue.TaskQueue <- NotifierTask{
			//To:       owner.Email,
			From:     ZCMailerFromAddress,
			To:       transmission.AdminAccount.Email,
			Subject:  fmt.Sprintf("Instance (%v) needs attention", transmission.InstanceId()),
			BodyHTML: newEmailBody,
			BodyText: newEmailBody,
		}

		logger.Info("Delete SQS Message")
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
		logger.Info("Lease needs approval")

		transmission.leaseDuration = time.Duration(ZCDefaultLeaseApprovalTimeoutDuration)
		var expiresAt = time.Now().UTC().Add(transmission.leaseDuration)

		// these will be used to compose the urls and verify the requests
		lease_uuid := uuid.NewV4().String()
		instance_id := transmission.InstanceId()
		token_once := uuid.NewV4().String() // one-time token

		newLease := Lease{
			UUID:      lease_uuid,
			TokenOnce: token_once,

			OwnerID:        transmission.owner.ID,
			CloudAccountID: transmission.CloudAccount.ID,
			AWSAccountID:   transmission.CloudAccount.AWSID,

			InstanceID:       transmission.InstanceId(),
			Region:           transmission.instanceRegion,
			AvailabilityZone: transmission.AvailabilityZone(),

			// Terminated bool `sql:"DEFAULT:false"`
			// Deleted    bool `sql:"DEFAULT:false"`
			Alerted: true,

			LaunchedAt:   transmission.InstanceLaunchTimeUTC(),
			ExpiresAt:    expiresAt,
			InstanceType: transmission.InstanceType(),
		}
		s.DB.Create(&newLease)
		logger.Info("new lease created",
			"lease", newLease,
		)

		// URL to approve lease
		approve_url, err := s.generateSignedEmailActionURL("approve", lease_uuid, instance_id, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

		// URL to terminate lease
		terminate_url, err := s.generateSignedEmailActionURL("terminate", lease_uuid, instance_id, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

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
				<a href="{{.instance_approve_url}}" target="_blank">Click here to <b>approve</b></a>

				<br>
				<br>
				
				Terminate immediately:
				<br>
				<br>
				<a href="{{.instance_terminate_url}}" target="_blank">Click here to <b>terminate</b></a>
				
				<br>
				<br>
				Thanks for using ZeroCloud!
				`,

			map[string]interface{}{
				"owner_email":        transmission.owner.Email,
				"n_of_active_leases": transmission.activeLeaseCount,
				"instance_id":        transmission.InstanceId(),
				"instance_type":      transmission.InstanceType(),
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
			Subject:  fmt.Sprintf("Instance (%v) needs approval", transmission.InstanceId()),
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
		logger.Info("Lease is OK -- register new lease in DB")

		transmission.DefineLeaseDuration()
		var expiresAt = time.Now().UTC().Add(transmission.leaseDuration)

		// these will be used to compose the urls and verify the requests
		lease_uuid := uuid.NewV4().String()
		instance_id := transmission.InstanceId()
		token_once := uuid.NewV4().String() // one-time token

		newLease := Lease{
			UUID:      lease_uuid,
			TokenOnce: token_once,

			OwnerID:        transmission.owner.ID,
			CloudAccountID: transmission.CloudAccount.ID,
			AWSAccountID:   transmission.CloudAccount.AWSID,

			InstanceID:       transmission.InstanceId(),
			Region:           transmission.instanceRegion,
			AvailabilityZone: transmission.AvailabilityZone(),

			// Terminated bool `sql:"DEFAULT:false"`
			// Deleted    bool `sql:"DEFAULT:false"`
			Alerted: false, // the lease does not need an action response, no alert has been sent out

			LaunchedAt:   transmission.InstanceLaunchTimeUTC(),
			ExpiresAt:    expiresAt,
			InstanceType: transmission.InstanceType(),
		}
		s.DB.Create(&newLease)
		logger.Info("new lease created",
			"lease", newLease,
		)

		// URL to terminate lease
		terminate_url, err := s.generateSignedEmailActionURL("terminate", lease_uuid, instance_id, token_once)
		if err != nil {
			// TODO: notify ZC admins
			return fmt.Errorf("error while generating signed URL: %v", err)
		}

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
				<a href="{{.instance_terminate_url}}" target="_blank">Click here to <b>terminate</b></a>

				<br>
				<br>
				
				Thanks for using ZeroCloud!
				`,

			map[string]interface{}{
				"owner_email":     transmission.owner.Email,
				"instance_id":     transmission.InstanceId(),
				"instance_type":   transmission.InstanceType(),
				"instance_region": transmission.instanceRegion,

				"termination_time":  expiresAt.Format("2006-01-02 15:04:05 GMT"),
				"instance_duration": transmission.leaseDuration.String(),

				"instance_terminate_url": terminate_url,
			},
		)
		s.NotifierQueue.TaskQueue <- NotifierTask{
			From:     ZCMailerFromAddress,
			To:       transmission.owner.Email,
			Subject:  fmt.Sprintf("Instance (%v) created", transmission.InstanceId()),
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

}
