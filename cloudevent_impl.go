package main

import (
	"fmt"
	"time"

	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// Create runs the create action.
func (c *CloudeventController) CreateImpl(ctx *app.CreateCloudeventContext) error {

	awsAccountId := ctx.Payload.Message.Account
	logger.Info("Create CloudEvent", "aws_account_id", awsAccountId)

	// try to find the CloudAccount that has an upstream_account_id that matches param
	cloudAccount := models.CloudAccount{}
	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: awsAccountId}).First(&cloudAccount)
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

	// TODO: it should be saving the EC2InstanceTags field into the CloudEvent,
	// The EC2InstanceTags field contains a JSON array of EC2 tags:
	// [{\"Key\":\"Name\",\"ResourceId\":\"i-0e730d938c710879e\",\"ResourceType\":\"instance\",\"Value\":\"blah2\"},{\"Key\":\"foo\",\"ResourceId\":\"i-0e730d938c710879e\",\"ResourceType\":\"instance\",\"Value\":\"baz3\"}]

	// Create a Lease object that references this (immutable) CloudEvent and expires
	// based on the settings in the Account
	// TODO: or can this be an AfterCreate callback on the CloudEvent?
	// file:///Users/tleyden/DevLibraries/gorm/callbacks.html
	err = createLease(ctx, cloudEvent)
	if err != nil {
		return ErrDatabaseError(err)
	}

	// TODO: should this return the path to the cloudevent .. should there even be one?
	/// ctx.ResponseData.Header().Set("Location", app.CloudeventHref(ctx.AccountID, a.ID))

	return ctx.Created()

}

func createLease(ctx *app.CreateCloudeventContext, cloudEvent models.CloudEvent) error {

	lease := models.Lease{}
	lease.CloudEvent = cloudEvent

	// This didn't work:
	// lease.CloudAccount = cloudEvent.CloudAccount
	// lease.Account = cloudEvent.Account
	// So I switched to using ID's
	lease.CloudAccountID = cloudEvent.CloudAccountID
	lease.AccountID = cloudEvent.AccountID

	logger.Info("Creating lease from cloudEvent", "cloudevent", fmt.Sprintf("%+v", cloudEvent))

	lease.Expires = time.Now().Add(time.Duration(3 * 24 * time.Hour)) // TODO calculate this based on Account settings
	lease.State = "Active"                                            // TODO - create an Enum and use that

	err := ldb.Add(ctx.Context, &lease)
	if err != nil {
		return err
	}

	return nil
}
