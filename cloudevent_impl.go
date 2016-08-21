package main

import (
	"fmt"

	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// Create runs the create action.
func (c *CloudeventController) CreateImpl(ctx *app.CreateCloudeventContext) error {

	awsAccountId := ctx.Payload.Message.Account
	logger.Info("Create CloudEvent", "aws_account_id", awsAccountId)

	// try to find the CloudAccount that has an upstream_account_id that matches param
	cloudAccount := models.CloudAccount{}
	// cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: awsAccountId}).First(&cloudAccount)
	// err := cdb.Db.Scopes(CloudAccountFilterByAccount(accountID, m.Db)).Table(m.TableName()).Preload("CloudEvents").Preload("Leases").Preload("Account").Where("id = ?", id).Find(&native).Error
	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: awsAccountId}).Preload("Account").First(&cloudAccount)

	logger.Info("Found CloudAccount", "CloudAccount", fmt.Sprintf("%+v", cloudAccount))

	// Make sure we found a valid CloudAccount, otherwise abort
	if cloudAccount.ID == 0 {
		ctx.ResponseData.Service.Send(
			ctx.Context,
			400,
			fmt.Sprintf("Could not find CloudAccount with upstream provider account id: %v", awsAccountId),
		)
		return nil
	}

	// Save the raw CloudEvent to the database
	cloudEvent := models.CloudEvent{}
	cloudEvent.AwsAccountID = awsAccountId
	cloudEvent.CloudAccountID = cloudAccount.ID
	cloudEvent.AccountID = cloudAccount.AccountID
	cloudEvent.SqsPayloadBase64 = *ctx.Payload.SQSPayloadBase64
	cloudEvent.CwEventSource = *ctx.Payload.Message.Source
	cloudEvent.CwEventTimestamp = *ctx.Payload.Message.Time
	cloudEvent.CwEventDetailInstanceID = ctx.Payload.Message.Detail.InstanceID
	cloudEvent.CwEventDetailState = ctx.Payload.Message.Detail.State

	err := edb.Add(ctx.Context, &cloudEvent)
	if err != nil {
		return ErrDatabaseError(err)
	}

	cloudEventID := cloudEvent.ID
	logger.Info("Saved CloudEvent", "ID", fmt.Sprintf("%+v", cloudEventID))
	logger.Info("Saved CloudEvent", "CloudEvent", fmt.Sprintf("%+v", cloudEvent))

	cloudEventFromDB, err := edb.OneCloudevent(
		ctx.Context,
		cloudEventID,
		cloudAccount.AccountID,
		cloudAccount.ID,
	)
	if err != nil {
		return ErrDatabaseError(err)
	}
	// cloudEventModelFromDB := cloudEventFromDB.CloudEventToCloudevent()

	logger.Info("CloudEvent from DB", "CloudEvent", fmt.Sprintf("%+v", cloudEventFromDB))
	// cloudEventFromDB.loadRelatedModels()
	// logger.Info("CloudEvent from DB after loading related models", "CloudEvent", fmt.Sprintf("%+v", cloudEventFromDB))

	// TODO: it should be saving the EC2InstanceTags field into the CloudEvent,
	// The EC2InstanceTags field contains a JSON array of EC2 tags:
	// [{\"Key\":\"Name\",\"ResourceId\":\"i-0e730d938c710879e\",\"ResourceType\":\"instance\",\"Value\":\"blah2\"},{\"Key\":\"foo\",\"ResourceId\":\"i-0e730d938c710879e\",\"ResourceType\":\"instance\",\"Value\":\"baz3\"}]

	// Create a Lease object that references this (immutable) CloudEvent and expires
	// based on the settings in the Account
	// TODO: or can this be an AfterCreate callback on the CloudEvent?
	// file:///Users/tleyden/DevLibraries/gorm/callbacks.html
	lease, err := createLease(ctx, cloudEvent, cloudAccount.Account)
	if err != nil {
		return ErrDatabaseError(err)
	}
	logger.Info("Created lease", "lease", fmt.Sprintf("%+v", lease))

	// TODO: should this return the path to the cloudevent .. should there even be one?
	/// ctx.ResponseData.Header().Set("Location", app.CloudeventHref(ctx.AccountID, a.ID))

	return ctx.Created()

}

func createLease(ctx *app.CreateCloudeventContext, cloudEvent models.CloudEvent, account models.Account) (models.Lease, error) {

	lease := models.Lease{}
	lease.CloudEvent = cloudEvent

	// This didn't work:
	// lease.CloudAccount = cloudEvent.CloudAccount
	// lease.Account = cloudEvent.Account
	// So I switched to using ID's
	lease.CloudAccountID = cloudEvent.CloudAccountID
	lease.AccountID = cloudEvent.AccountID

	logger.Info("Creating lease from cloudEvent", "cloudevent", fmt.Sprintf("%+v", cloudEvent))

	// Set the expiration time of this lease based on Account
	leaseExpiresIn := account.LeaseExpiresIn
	leaseExpiresInUnits := account.LeaseExpiresInUnits

	if err := lease.SetExpiryTime(leaseExpiresIn, leaseExpiresInUnits); err != nil {
		return lease, err
	}

	lease.State = "Active" // TODO - create an Enum and use that

	// Save the lease to the database
	err := ldb.Add(ctx.Context, &lease)
	if err != nil {
		return lease, err
	}

	// OneLease(ctx context.Context, id int, accountID int, cloudAccountID int, cloudEventID int) (*app.Lease, error) {
	leaseID := lease.ID
	leaseFromDb, err := ldb.OneLease(ctx.Context,
		leaseID,
		cloudEvent.AccountID,
		cloudEvent.CloudAccountID,
		cloudEvent.ID,
	)
	logger.Info("Load", "leaseFromDb", fmt.Sprintf("%+v", leaseFromDb))

	return lease, nil
}
