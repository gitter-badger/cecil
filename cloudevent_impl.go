package main

import (
	"fmt"

	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// Create runs the create action.
func (c *CloudeventController) CreateImpl(ctx *app.CreateCloudeventContext) error {

	awsAccountId := *ctx.Payload.AwsAccountID
	logger.Info("Create CloudEvent", "aws_account_id", awsAccountId)

	// try to find the CloudAccount that has an upstream_account_id that matches param
	cloudAccount := models.CloudAccount{}
	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: awsAccountId}).First(&cloudAccount)
	logger.Info("Found CloudAccount", "CloudAccount", fmt.Sprintf("%+v", cloudAccount))

	// Make sure we found a valid CloudAccount, otherwise abort
	if cloudAccount.ID == 0 {
		ctx.ResponseData.Service.Send(ctx.Context, 400, fmt.Sprintf("Could not find CloudAccount with upstream provider account id: %v", awsAccountId))
		return nil
	}

	// Get the instance_tags (assumes they have been looked up after pulling from SQS and calling this REST api)

	// Save the raw CloudEvent to the database

	// Call CloudEvent.createLeaseMaybe()
	// Or can this be an AfterCreate callback on the CloudEvent?
	// file:///Users/tleyden/DevLibraries/gorm/callbacks.html

	// Put your logic here
	e := models.CloudEvent{}
	e.AwsAccountID = awsAccountId
	e.CloudAccountID = cloudAccount.ID
	e.AccountID = cloudAccount.AccountID

	err := edb.Add(ctx.Context, &e)
	if err != nil {
		return ErrDatabaseError(err)
	}
	/// ctx.ResponseData.Header().Set("Location", app.CloudeventHref(ctx.AccountID, a.ID))
	return ctx.Created()

}
