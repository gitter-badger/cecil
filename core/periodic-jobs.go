package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/tleyden/zerocloud/mocks/aws"
)

// @@@@@@@@@@@@@@@ Periodic Jobs @@@@@@@@@@@@@@@

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
		}

		logger.Info("Parsed sqs message", "message", transmission.Message)

		if !transmission.TopicAndInstanceHaveSameOwner() {
			// the originating SNS topic and the instance have different owners (different AWS accounts)
			// TODO: notify zerocloud admin
			logger.Warn("topicAWSID != instanceOriginatorID", "topicAWSID", transmission.Topic.AWSID, "instanceOriginatorID", transmission.Message.Account)
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

		/// TODO: pass transmission to NewLeaseQueue

		//check whether someone with this aws adminAccount id is registered at zerocloud
		err = transmission.FetchCloudAccount()
		if err != nil {
			// TODO: notify admin; something fishy is going on.
			logger.Warn("originator is not registered", "AWSID", transmission.Topic.AWSID)
			continue
		}

		// check whether the cloud account has an admin account
		err = transmission.FetchAdminAccount()
		if err != nil {
			// TODO: notify admin; something fishy is going on.
			logger.Warn("Error while retrieving admin account", "error", err)
			continue
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
			continue
		}

		err = transmission.CreateAssumedEC2Service()
		if err != nil {
			// TODO: this might reveal too much to the admin about zerocloud; be selective and cautious
			s.sendMisconfigurationNotice(err, transmission.AdminAccount.Email)
			logger.Warn("error while creating ec2 service with assumed service", "error", err)
			continue
		}

		err = transmission.DescribeInstance()
		if err != nil {
			// TODO: this might reveal too much to the admin about zerocloud; be selective and cautious
			s.sendMisconfigurationNotice(err, transmission.AdminAccount.Email)
			logger.Warn("error while describing instances", "error", err)
			continue
		}

		// check whether the instance specified in the event exists on aws
		if !transmission.InstanceExists() {
			logger.Warn("Instance does not exist", "instanceID", transmission.Message.Detail.InstanceID)
			// remove message from queue
			err := transmission.DeleteMessage()
			if err != nil {
				logger.Warn("DeleteMessage", "error", err)
			}
			continue
		}

		logger.Info("describeInstances", "response", transmission.describeInstancesResponse)

		err = transmission.FetchInstance()
		if err != nil {
			logger.Warn("error while fetching instance description", "error", err)
			continue
		}

		err = transmission.ComputeInstanceRegion()
		if err != nil {
			logger.Warn("error while computing instance region", "error", err)
			continue
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
			continue
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
			continue
		}

		if !transmission.LeaseIsNew() {
			// TODO: notify admin
			logger.Warn("instanceCount != 0")
			continue
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

			newLease := Lease{
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
				
				Terminate now:
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

						"termination_time":       expiresAt.Format("2006-01-02 15:04:05 GMT"),
						"instance_duration":      transmission.leaseDuration.String(),
						"instance_renew_url":     "",
						"instance_terminate_url": "",
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
				
				Terminate now:
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

						"termination_time":       expiresAt.Format("2006-01-02 15:04:05 GMT"),
						"instance_duration":      transmission.leaseDuration.String(),
						"instance_renew_url":     "",
						"instance_terminate_url": "",
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
			continue
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

			newLease := Lease{
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
				{{.instance_renew_url}}

				<br>
				<br>
				
				Terminate:
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

					"termination_time":       expiresAt.Format("2006-01-02 15:04:05 GMT"),
					"instance_renew_url":     "",
					"instance_terminate_url": "",
					"instance_duration":      transmission.leaseDuration.String(),
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
			continue
		} else {
			// register new lease in DB
			// set its expiration to zone.default_expiration (if > 0), or cloudAccount.default_expiration, or adminAccount.default_expiration

			transmission.DefineLeaseDuration()
			var expiresAt = time.Now().UTC().Add(transmission.leaseDuration)

			newLease := Lease{
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

			newEmailBody := compileEmail(
				`Hey {{.owner_email}}, you (or someone else using your ZeroCloudOwner tag) created a new instance 
				(id <b>{{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). That's AWESOME!

				<br>
				<br>

				Your instance will be terminated at <b>{{.termination_time}}</b> ({{.instance_duration}} after it's creation).

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
			continue
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
	}

	return nil
}

func (s *Service) AlerterJob() error {

	return nil
}

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

type Transmission struct {
	s       *Service
	Message mockaws.SQSMessage

	Topic struct {
		Region string
		AWSID  string
	}
	deleteMessageFromQueueParams *sqs.DeleteMessageInput

	CloudAccount CloudAccount
	AdminAccount Account

	assumedService            *session.Session
	ec2Service                ec2iface.EC2API
	describeInstancesResponse *ec2.DescribeInstancesOutput

	Instance           ec2.Instance
	instanceRegion     string
	externalOwnerEmail string
	owner              Owner
	leaseDuration      time.Duration
	activeLeaseCount   int64
}

func (s *Service) parseSQSTransmission(rawMessage *sqs.Message, queueURL string) (*Transmission, error) {
	var newTransmission Transmission = Transmission{}
	newTransmission.s = s

	// parse the envelope
	var envelope mockaws.SQSEnvelope
	err := json.Unmarshal([]byte(*rawMessage.Body), &envelope)
	if err != nil {
		return &Transmission{}, err
	}

	newTransmission.deleteMessageFromQueueParams = &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),                  // Required
		ReceiptHandle: aws.String(*rawMessage.ReceiptHandle), // Required
	}

	// parse the message
	err = json.Unmarshal([]byte(envelope.Message), &newTransmission.Message)
	if err != nil {
		return &Transmission{}, err
	}

	// extract some values
	// TODO: check whether these values are not empty
	/// transmission.Topic.Region
	/// transmission.Topic.AWSID
	topicArn := strings.Split(envelope.TopicArn, ":")
	newTransmission.Topic.Region = topicArn[3]
	newTransmission.Topic.AWSID = topicArn[4]

	return &newTransmission, nil
}

// the originating SNS topic and the instance have different owners (different AWS accounts)
func (t *Transmission) TopicAndInstanceHaveSameOwner() bool {
	instanceOriginatorID := t.Message.Account

	return t.Topic.AWSID == instanceOriginatorID
}

// consider only pending and terminated status messages; ignore the rest
func (t *Transmission) MessageIsRelevant() bool {
	return t.Message.Detail.State == ec2.InstanceStateNamePending ||
		t.Message.Detail.State == ec2.InstanceStateNameTerminated
}

func (t *Transmission) DeleteMessage() error {
	// remove message from queue
	return retry(5, time.Duration(3*time.Second), func() error {
		var err error
		// TODO:
		_, err = t.s.AWS.SQS.DeleteMessage(t.deleteMessageFromQueueParams)
		return err
	})
}

//check whether someone with this aws adminAccount id is registered at zerocloud
func (t *Transmission) FetchCloudAccount() error {
	var cloudOwnerCount int64
	t.s.DB.Where(&CloudAccount{AWSID: t.Topic.AWSID}).
		First(&t.CloudAccount).
		Count(&cloudOwnerCount)
	if cloudOwnerCount == 0 {
		return fmt.Errorf("No cloud account for AWSID %v", t.Topic.AWSID)
	}
	if cloudOwnerCount > 1 {
		return fmt.Errorf("Too many (%v) CloudAccounts for AWSID %v", cloudOwnerCount, t.Topic.AWSID)
	}
	return nil
}

// check whether the cloud account has an admin account
func (t *Transmission) FetchAdminAccount() error {
	var cloudAccountAdminCount int64
	t.s.DB.Model(&t.CloudAccount).Related(&t.AdminAccount).Count(&cloudAccountAdminCount)
	//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&cloudAccountAdminCount)
	if cloudAccountAdminCount == 0 {
		return fmt.Errorf("No admin for CloudAccount", "CloudAccount.ID", t.CloudAccount.ID)
	}
	if cloudAccountAdminCount > 1 {
		return fmt.Errorf("Too many (%v) admins for CloudAccount %v", cloudAccountAdminCount, t.CloudAccount.ID)
	}
	return nil
}

func (t *Transmission) CreateAssumedService() error {
	assumedConfig := &aws.Config{
		Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
			Client: sts.New(t.s.AWS.Session, &aws.Config{Region: aws.String(t.Topic.Region)}),
			RoleARN: fmt.Sprintf(
				"arn:aws:iam::%v:role/%v",
				t.Topic.AWSID,
				viper.GetString("ForeignRoleName"),
			),
			RoleSessionName: uuid.NewV4().String(),
			ExternalID:      aws.String(t.CloudAccount.ExternalID),
			ExpiryWindow:    60 * time.Second,
		}),
	}

	t.assumedService = session.New(assumedConfig)

	return nil
}

func (t *Transmission) CreateAssumedEC2Service() error {
	t.ec2Service = t.s.EC2(t.assumedService, t.Topic.Region)
	return nil
}

func (t *Transmission) DescribeInstance() error {
	paramsDescribeInstance := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(t.Message.Detail.InstanceID),
		},
	}
	var err error
	t.describeInstancesResponse, err = t.ec2Service.DescribeInstances(paramsDescribeInstance)
	return err
}

// check whether the instance specified in the event message exists on aws
func (t *Transmission) InstanceExists() bool {
	if len(t.describeInstancesResponse.Reservations) == 0 {
		// logger.Warn("len(describeInstancesResponse.Reservations) == 0")
		return false
	}
	if len(t.describeInstancesResponse.Reservations[0].Instances) == 0 {
		// logger.Warn("len(describeInstancesResponse.Reservations[0].Instances) == 0")
		return false
	}
	return true
}

func (t *Transmission) FetchInstance() error {
	// TODO: merge the preceding check operations here
	t.Instance = *t.describeInstancesResponse.Reservations[0].Instances[0]

	if *t.Instance.InstanceId != t.Message.Detail.InstanceID {
		return fmt.Errorf("instance.InstanceId != message.Detail.InstanceID")
	}
	return nil
}

func (t *Transmission) ComputeInstanceRegion() error {

	// TODO: is this always valid?
	// TODO: use pointer or by value?
	az := *t.Instance.Placement.AvailabilityZone
	t.instanceRegion = az[:len(az)-1]

	return nil
}

func (t *Transmission) InstanceIsTerminated() bool {
	return *t.Instance.State.Name == ec2.InstanceStateNameTerminated
}

func (t *Transmission) InstanceIsPendingOrRunning() bool {
	return *t.Instance.State.Name == ec2.InstanceStateNamePending ||
		*t.Instance.State.Name == ec2.InstanceStateNameRunning
}

// IsNew: check whether a lease with the same instanceID exists
func (t *Transmission) LeaseIsNew() bool {
	var instanceCount int64
	t.s.DB.Table("leases").Where(&Lease{InstanceID: t.Message.Detail.InstanceID}).Count(&instanceCount)

	return instanceCount == 0
}

func (t *Transmission) InstanceHasGoodOwnerTag() bool {
	// InstanceHasTags: check whether instance has tags
	if len(t.Instance.Tags) == 0 {
		return false
	}
	//logger.Warn("len(instance.Tags) == 0")

	// InstanceHasOwnerTag: check whether the instance has an zerocloudowner tag
	for _, tag := range t.Instance.Tags {
		if strings.ToLower(*tag.Key) != "zerocloudowner" {
			continue
		}

		// OwnerTagValueIsValid: check whether the zerocloudowner tag is a valid email
		ownerTag, err := t.s.Mailer.ValidateEmail(*tag.Value)
		if err != nil {
			logger.Warn("ValidateEmail", "error", err)
			// TODO: send notification to admin
			return false
		}
		if !ownerTag.IsValid {
			return false
			// TODO: notify admin: "Warning: zerocloudowner tag email not valid" (DO NOT INCLUDE IT IN THE EMAIL, OR HTML-ESCAPE IT)
		}
		// fmt.Printf("Parts local_part=%s domain=%s display_name=%s", ownerTag.Parts.LocalPart, ownerTag.Parts.Domain, ownerTag.Parts.DisplayName)
		t.externalOwnerEmail = ownerTag.Address
	}
	return true
}

// ExternalOwnerIsWhitelisted: check whether the owner email in the tag is a whitelisted owner email
func (t *Transmission) ExternalOwnerIsWhitelisted() bool {
	// TODO: select Owner by email, cloudaccountid, and region?
	if t.externalOwnerEmail == "" {
		return false
	}
	var ownerCount int64
	t.s.DB.Table("owners").Where(&Owner{Email: t.externalOwnerEmail, CloudAccountID: t.CloudAccount.ID}).Count(&ownerCount)
	if ownerCount != 1 {
		return false
	}
	return true
}

func (t *Transmission) SetExternalOwnerAsOwner() error {
	var owners []Owner
	var ownerCount int64
	t.s.DB.Table("owners").Where(&Owner{Email: t.externalOwnerEmail, CloudAccountID: t.CloudAccount.ID}).Find(&owners).Count(&ownerCount)
	if ownerCount != 1 {
		return fmt.Errorf("Too many external owners with externalOwnerEmail %v", t.externalOwnerEmail)
	}
	t.owner = owners[0]
	return nil
}
func (t *Transmission) SetAdminAsOwner() error {
	var owners []Owner
	var ownerCount int64
	t.s.DB.Table("owners").Where(&Owner{Email: t.AdminAccount.Email, CloudAccountID: t.CloudAccount.ID}).Find(&owners).Count(&ownerCount)
	if ownerCount != 1 {
		return fmt.Errorf("Too many admin owners with email %v", t.AdminAccount.Email)
	}
	t.owner = owners[0]
	return nil
}

func (t *Transmission) DefineLeaseDuration() {
	t.leaseDuration = time.Duration(ZCDefaultLeaseDuration)

	if t.AdminAccount.DefaultLeaseDuration > 0 {
		t.leaseDuration = time.Duration(t.AdminAccount.DefaultLeaseDuration)
	}
	if t.CloudAccount.DefaultLeaseDuration > 0 {
		t.leaseDuration = time.Duration(t.CloudAccount.DefaultLeaseDuration)
	}
}

func (t *Transmission) LeaseNeedsApproval() bool {
	var leases []Lease

	t.s.DB.Table("leases").Where(&Lease{
		OwnerID:        t.owner.ID,
		CloudAccountID: t.CloudAccount.ID,
		Terminated:     false,
	}).Find(&leases).Count(&t.activeLeaseCount)
	//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&activeLeaseCount)

	return t.activeLeaseCount >= ZCMaxLeasesPerOwner
}
