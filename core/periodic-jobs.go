package core

import (
	"encoding/json"
	"fmt"
	"strconv"
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
)

// @@@@@@@@@@@@@@@ Periodic Jobs @@@@@@@@@@@@@@@

func (s *Service) EventInjestorJob() error {
	// TODO: verify event origin (must be aws, not someone else)
	fmt.Println("EventInjestorJob() run")

	queueURL := fmt.Sprintf("https://sqs.%v.amazonaws.com/%v/%v",
		viper.GetString("AWS_REGION"),
		viper.GetString("AWS_ACCOUNT_ID"),
		viper.GetString("SQSQueueName"),
	)
	receiveMessageParams := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL), // Required
		//MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout: aws.Int64(3), // should be higher, like 30, the time to finish doing everything
		WaitTimeSeconds:   aws.Int64(3),
	}
	receiveMessageResponse, err := s.AWS.SQS.ReceiveMessage(receiveMessageParams)

	if err != nil {
		return fmt.Errorf("EventInjestorJob() error: %v", err)
	}

	fmt.Println(receiveMessageResponse)

OnMessagesLoop:
	for messageIndex := range receiveMessageResponse.Messages {
		var envelope struct {
			Type             string `json:"Type"`
			MessageId        string `json:"MessageId"`
			TopicArn         string `json:"TopicArn"`
			Message          string `json:"Message"`
			Timestamp        string `json:"Timestamp"`
			SignatureVersion string `json:"SignatureVersion"`
			Signature        string `json:"Signature"`
			SigningCertURL   string `json:"SigningCertURL"`
			UnsubscribeURL   string `json:"UnsubscribeURL"`
		}
		err := json.Unmarshal([]byte(*receiveMessageResponse.Messages[messageIndex].Body), &envelope)
		if err != nil {
			return err
		}

		var deleteMessageFromQueueParams = &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),                                                     // Required
			ReceiptHandle: aws.String(*receiveMessageResponse.Messages[messageIndex].ReceiptHandle), // Required
		}

		var message struct {
			Version    string   `json:"version"`
			ID         string   `json:"id"`
			DetailType string   `json:"detail-type"`
			Source     string   `json:"source"`
			Account    string   `json:"account"`
			Time       string   `json:"time"`
			Region     string   `json:"region"`
			Resources  []string `json:"resources"`
			Detail     struct {
				InstanceID string `json:"instance-id"`
				State      string `json:"state"`
			} `json:"detail"`
		}
		err = json.Unmarshal([]byte(envelope.Message), &message)
		if err != nil {
			return err
		}

		topicArn := strings.Split(envelope.TopicArn, ":")
		topicRegion := topicArn[3]
		topicOwnerID, err := strconv.ParseUint(topicArn[4], 10, 64)
		if err != nil {
			// TODO: notify
			fmt.Println(err)
			continue
		}
		// topicName := topicArn[5]
		instanceOriginatorID, err := strconv.ParseUint(message.Account, 10, 64)
		if err != nil {
			// TODO: notify
			fmt.Println(err)
			continue
		}
		// TODO: check these values are not empty

		if topicOwnerID != instanceOriginatorID {
			// the originating SNS topic and the instance have different owners
			// TODO: notify zerocloud admin
			fmt.Println("topicOwnerID != instanceOriginatorID")
			continue
		}

		// consider only pending and terminated status messages; ignore the rest
		if message.Detail.State != ec2.InstanceStateNamePending &&
			message.Detail.State != ec2.InstanceStateNameTerminated {
			fmt.Println("removing")
			// remove message from queue
			err := retry(5, time.Duration(3*time.Second), func() error {
				var err error
				_, err = s.AWS.SQS.DeleteMessage(deleteMessageFromQueueParams)
				return err
			})
			if err != nil {
				fmt.Println(err)
			}
			continue // next message
		}

		// HasOwner: check whether someone with this aws account id is registered
		var cloudAccount CloudAccount
		var cloudOwnerCount int64
		s.DB.Where(&CloudAccount{AWSID: topicOwnerID}).First(&cloudAccount).Count(&cloudOwnerCount)
		if cloudOwnerCount == 0 {
			// TODO: notify admin; something fishy is going on.
			continue
		}

		// <debug>
		var accounts []Account
		s.DB.Table("accounts").Find(&accounts)
		fmt.Printf("accounts: %#v\n", accounts)
		fmt.Printf("%s\n", accounts[0].ID)
		fmt.Println("looking for:", cloudAccount.AccountID)
		// </debug>

		var account Account
		var cloudAccountOwnerCount int64
		s.DB.Model(&cloudAccount).Related(&account).Count(&cloudAccountOwnerCount)
		//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&cloudAccountOwnerCount)
		if cloudAccountOwnerCount == 0 {
			// TODO: notify admin; something fishy is going on.
			fmt.Println("cloudAccountOwnerCount == 0")
			continue
		}

		fmt.Printf("account: %#v\n", account)

		// IsNew: check whether a lease with the same instanceID exists
		var instanceCount int64
		s.DB.Table("leases").Where(&Lease{InstanceID: message.Detail.InstanceID}).Count(&instanceCount)
		fmt.Println("here")
		if instanceCount != 0 {
			// TODO: notify admin
			fmt.Println("instanceCount != 0")
			continue
		}

		// assume role

		assumedConfig := &aws.Config{
			Credentials: credentials.NewCredentials(&stscreds.AssumeRoleProvider{
				Client:          sts.New(s.AWS.Session, &aws.Config{Region: aws.String(topicRegion)}),
				RoleARN:         fmt.Sprintf("arn:aws:iam::%v:role/ZeroCloudRole", topicOwnerID),
				RoleSessionName: uuid.NewV4().String(),
				ExternalID:      aws.String("slavomir"),
				ExpiryWindow:    3 * time.Minute,
			}),
		}

		assumedService := session.New(assumedConfig)

		ec2Service := ec2.New(assumedService,
			&aws.Config{
				Region: aws.String(topicRegion),
			},
		)

		paramsDescribeInstance := &ec2.DescribeInstancesInput{
			// DryRun: aws.Bool(true),
			InstanceIds: []*string{
				aws.String(message.Detail.InstanceID),
			},
		}
		describeInstancesResponse, err := ec2Service.DescribeInstances(paramsDescribeInstance)

		if err != nil {
			// TODO: notify
			fmt.Println(err)
			continue
		}

		// ExistsOnAWS: check whether the instance specified in the event exists on aws
		if len(describeInstancesResponse.Reservations) == 0 {
			fmt.Println("len(describeInstancesResponse.Reservations) == 0: ")
			// remove message from queue
			err := retry(5, time.Duration(3*time.Second), func() error {
				var err error
				_, err = s.AWS.SQS.DeleteMessage(deleteMessageFromQueueParams)
				return err
			})
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		if len(describeInstancesResponse.Reservations[0].Instances) == 0 {
			fmt.Println("len(describeInstancesResponse.Reservations[0].Instances) == 0: ")
			// remove message from queue
			err := retry(5, time.Duration(3*time.Second), func() error {
				var err error
				_, err = s.AWS.SQS.DeleteMessage(deleteMessageFromQueueParams)
				return err
			})
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		fmt.Println("description: ", describeInstancesResponse)

		var instance = describeInstancesResponse.Reservations[0].Instances[0]

		//instance.InstanceType
		//instance.LaunchTime

		if *instance.InstanceId != message.Detail.InstanceID {
			fmt.Println("instance.InstanceId !=message.Detail.InstanceID")
			continue
		}

		if *instance.State.Name != ec2.InstanceStateNamePending &&
			*instance.State.Name != ec2.InstanceStateNameRunning {
			fmt.Println("the retried state is neither pending not running:", *instance.State.Name)
			continue
		}

		var instanceHasValidOwnerTag bool = false
		var ownerIsWhitelisted bool = false
		var ownerEmail string = account.Email

		// InstanceHasTags: check whethe instance has tags
		if len(instance.Tags) > 0 {
			fmt.Println("len(instance.Tags) == 0")

			// InstanceHasOwnerTag: check whether the instance has an zerocloudowner tag
			for _, tag := range instance.Tags {
				if strings.ToLower(*tag.Key) != "zerocloudowner" {
					continue
				}

				// OwnerTagValueIsValid: check whether the zerocloudowner tag is a valid email
				ownerTag, err := s.Mailer.ValidateEmail(*tag.Value)
				if err != nil {
					fmt.Println(err)
					break
				}
				if !ownerTag.IsValid {
					fmt.Println("email not valid")
					// TODO: notify admin: "Warning: zerocloudowner tag email not valid" (DO NOT INCLUDE IT IN THE EMAIL, OR HTML-ESCAPE IT)
					break
				}
				fmt.Printf("Parts local_part=%s domain=%s display_name=%s", ownerTag.Parts.LocalPart, ownerTag.Parts.Domain, ownerTag.Parts.DisplayName)
				ownerEmail = ownerTag.Address
				instanceHasValidOwnerTag = true
				break
			}
		}

		var owners []Owner
		var ownerCount int64

		// OwnerTagIsWhitelisted: check whether the owner email in the tag is a whitelisted owner email

		// TODO: select Owner by email, cloudaccountid, and region?
		s.DB.Table("owners").Where(&Owner{Email: ownerEmail, CloudAccountID: cloudAccount.ID}).Find(&owners).Count(&ownerCount)
		if ownerEmail != account.Email && ownerCount != 1 {
			ownerIsWhitelisted = false
			s.DB.Table("owners").Where(&Owner{Email: account.Email, CloudAccountID: cloudAccount.ID}).Find(&owners).Count(&ownerCount)
		}
		if ownerCount == 0 {
			// TODO: fatal: admin does not have an entry in the owners table
			fmt.Println("fatal: admin does not have an entry in the owners table")
			continue OnMessagesLoop
		}

		var owner = owners[0] // assuming that each admin has an entry in the owners table
		var lifetime time.Duration = time.Duration(ZCDefaultLeaseExpiration)

		if account.DefaultLeaseExpiration > 0 {
			lifetime = time.Duration(account.DefaultLeaseExpiration)
		}
		if cloudAccount.DefaultLeaseExpiration > 0 {
			lifetime = time.Duration(cloudAccount.DefaultLeaseExpiration)
		}

		if !instanceHasValidOwnerTag || !ownerIsWhitelisted {
			// assign instance to admin, and send notification to admin
			// owner is not whitelisted: notify admin: "Warning: zerocloudowner tag email not in whitelist"

			lifetime = time.Duration(ZCDefaultTruceDuration)
			var terminationTime = time.Now().Add(lifetime)

			newLease := Lease{
				OwnerID:        owner.ID,
				CloudAccountID: cloudAccount.ID,
				AWSAccountID:   cloudAccount.AWSID,

				InstanceID: *instance.InstanceId,
				Region:     *instance.Placement.AvailabilityZone,

				// Terminated bool `sql:"DEFAULT:false"`
				// Deleted    bool `sql:"DEFAULT:false"`

				LaunchedAt:   *instance.LaunchTime,
				ExpiresAt:    terminationTime,
				InstanceType: *instance.InstanceType,
			}
			s.DB.Create(&newLease)

			var newEmailBody string

			if !instanceHasValidOwnerTag {
				newEmailBody = compileEmail(
					`Hey {{.owner_email}}, someone created a new instance 
				id <b>({{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). <br><br>

				It does not have a valid ZeroCloudOwner tag, so we assigned it to you.
				
				<br>
				<br>
				
				It will be terminated at <b>{{.termination_time}}</b> ({{.instance_lifetime}} after it's creation).

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
						"owner_email":     owner.Email,
						"instance_id":     *instance.InstanceId,
						"instance_type":   *instance.InstanceType,
						"instance_region": *instance.Placement.AvailabilityZone,

						"termination_time":       terminationTime.Format("2006-01-02 15:04:05 CET"),
						"instance_lifetime":      lifetime.String(),
						"instance_renew_url":     "",
						"instance_terminate_url": "",
					},
				)
			} else {
				newEmailBody = compileEmail(
					`Hey {{.owner_email}}, someone created a new instance 
				id <b>({{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). <br><br>

				The ZeroCloudOwner tag this instance has is not in the whitelist, so we assigned it to you.
				
				<br>
				<br>
				
				It will be terminated at <b>{{.termination_time}}</b> ({{.instance_lifetime}} after it's creation).

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
						"owner_email":     owner.Email,
						"instance_id":     *instance.InstanceId,
						"instance_type":   *instance.InstanceType,
						"instance_region": *instance.Placement.AvailabilityZone,

						"termination_time":       terminationTime.Format("2006-01-02 15:04:05 CET"),
						"instance_lifetime":      lifetime.String(),
						"instance_renew_url":     "",
						"instance_terminate_url": "",
					},
				)
			}
			s.NotifierQueue.TaskQueue <- NotifierTask{
				From:     ZCMailerFromAddress,
				To:       owner.Email,
				Subject:  fmt.Sprintf("Instance (%v) Needs Attention", *instance.InstanceId),
				BodyHTML: newEmailBody,
				BodyText: newEmailBody,
			}

			// remove message from queue
			err := retry(5, time.Duration(3*time.Second), func() error {
				var err error
				_, err = s.AWS.SQS.DeleteMessage(deleteMessageFromQueueParams)
				return err
			})
			if err != nil {
				fmt.Println(err)
			}

			continue
		}

		var leases []Lease
		var activeLeaseCount int64
		s.DB.Table("leases").Where(&Lease{
			OwnerID:        owner.ID,
			CloudAccountID: cloudAccount.ID,
			Terminated:     false,
		}).Find(&leases).Count(&activeLeaseCount)
		//s.DB.Table("accounts").Where([]uint{cloudAccount.AccountID}).First(&cloudAccount).Count(&activeLeaseCount)

		leaseNeedsApproval := activeLeaseCount >= ZCMaxLeasesPerOwner

		if !leaseNeedsApproval {
			// register new lease in DB
			// set its expiration to zone.default_expiration (if > 0), or cloudAccount.default_expiration, or account.default_expiration
			var terminationTime = time.Now().Add(lifetime)

			newLease := Lease{
				OwnerID:        owner.ID,
				CloudAccountID: cloudAccount.ID,
				AWSAccountID:   cloudAccount.AWSID,

				InstanceID: *instance.InstanceId,
				Region:     *instance.Placement.AvailabilityZone,

				// Terminated bool `sql:"DEFAULT:false"`
				// Deleted    bool `sql:"DEFAULT:false"`

				LaunchedAt:   *instance.LaunchTime,
				ExpiresAt:    time.Now().Add(lifetime),
				InstanceType: *instance.InstanceType,
			}
			s.DB.Create(&newLease)

			newEmailBody := compileEmail(
				`Hey {{.owner_email}}, you (or someone else) created a new instance 
				id <b>({{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). That's AWESOME!

				<br>
				<br>

				Your instance will be terminated at <b>{{.termination_time}}</b> ({{.instance_lifetime}} after it's creation).

				<br>
				<br>
				
				Thanks for using ZeroCloud!
				`,

				map[string]interface{}{
					"owner_email":     owner.Email,
					"instance_id":     *instance.InstanceId,
					"instance_type":   *instance.InstanceType,
					"instance_region": *instance.Placement.AvailabilityZone,

					"termination_time":  terminationTime.Format("2006-01-02 15:04:05 CET"),
					"instance_lifetime": lifetime.String(),
				},
			)
			s.NotifierQueue.TaskQueue <- NotifierTask{
				From:     ZCMailerFromAddress,
				To:       owner.Email,
				Subject:  fmt.Sprintf("Instance (%v) Created", *instance.InstanceId),
				BodyHTML: newEmailBody,
				BodyText: newEmailBody,
			}

			// remove message from queue
			err := retry(5, time.Duration(3*time.Second), func() error {
				var err error
				_, err = s.AWS.SQS.DeleteMessage(deleteMessageFromQueueParams)
				return err
			})
			if err != nil {
				fmt.Println(err)
			}

			continue
		} else {
			// register new lease in DB
			// expiry: 1h
			// send confirmation to owner: confirmation link, and termination link

			lifetime = time.Duration(ZCDefaultTruceDuration)
			var terminationTime = time.Now().Add(lifetime)

			newLease := Lease{
				OwnerID:        owner.ID,
				CloudAccountID: cloudAccount.ID,
				AWSAccountID:   cloudAccount.AWSID,

				InstanceID: *instance.InstanceId,
				Region:     *instance.Placement.AvailabilityZone,

				// Terminated bool `sql:"DEFAULT:false"`
				// Deleted    bool `sql:"DEFAULT:false"`

				LaunchedAt:   *instance.LaunchTime,
				ExpiresAt:    terminationTime,
				InstanceType: *instance.InstanceType,
			}
			s.DB.Create(&newLease)

			newEmailBody := compileEmail(
				`Hey {{.owner_email}}, you (or someone else) created a new instance 
				id <b>({{.instance_id}}</b>, of type <b>{{.instance_type}}</b>, 
				on <b>{{.instance_region}}</b>). <br><br>

				At the time of writing this email, you have {{.n_of_active_leases}} active
					leases, so we need your approval for this one. <br><br>

				Please click on "Approve" to approve this instance,
					otherwise it will be terminated at <b>{{.termination_time}}</b> (one hour after it's creation).

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
					"owner_email":        owner.Email,
					"n_of_active_leases": activeLeaseCount,
					"instance_id":        *instance.InstanceId,
					"instance_type":      *instance.InstanceType,
					"instance_region":    *instance.Placement.AvailabilityZone,

					"termination_time":       terminationTime.Format("2006-01-02 15:04:05 CET"),
					"instance_renew_url":     "",
					"instance_terminate_url": "",
				},
			)
			s.NotifierQueue.TaskQueue <- NotifierTask{
				From:     ZCMailerFromAddress,
				To:       owner.Email,
				Subject:  fmt.Sprintf("Instance (%v) Needs Approval", *instance.InstanceId),
				BodyHTML: newEmailBody,
				BodyText: newEmailBody,
			}

			// remove message from queue
			err := retry(5, time.Duration(3*time.Second), func() error {
				var err error
				_, err = s.AWS.SQS.DeleteMessage(deleteMessageFromQueueParams)
				return err
			})
			if err != nil {
				fmt.Println(err)
			}

			continue
		}

		// if message.Detail.State == ec2.InstanceStateNameTerminated
		// LeaseTerminatedQueue <- LeaseTerminatedTask{} and continue

		// get zc account who has a cloudaccount with awsID == topicOwnerID
		// if no one of our customers owns this account, error
		// fetch options config
		// roleARN := fmt.Sprintf("arn:aws:iam::%v:role/ZeroCloudRole",topicOwnerID)
		// assume role
		// fetch instance info
		// check if statuses match (this message was sent by aws.ec2)
		// message.Detail.InstanceID

		fmt.Printf("%v", message)
	}

	return nil
}

func (s *Service) AlerterJob() error {

	return nil
}

func (s *Service) SentencerJob() error {
	fmt.Println("SentencerJob() run")

	var expiredLeases []Lease
	var expiredLeasesCount int64

	fmt.Println("expired leases count: ", expiredLeasesCount)

	s.DB.Table("leases").Where("expires_at < ?", time.Now()).Find(&expiredLeases).Count(&expiredLeasesCount)

	for expiredLeaseIndex := range expiredLeases {
		fmt.Println("expired lease: ", expiredLeases[expiredLeaseIndex])
	}

	return nil
}
