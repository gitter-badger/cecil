package main

import (
	"fmt"

	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// Create runs the create action.
func (c *CloudaccountController) CreateImpl(ctx *app.CreateCloudaccountContext) error {

	// I think this is already called by goa somewhere
	// if err := ctx.Payload.Validate(); err != nil {
	// 	return err
	// }

	// Put your logic here
	cloudAccount := models.CloudAccount{}
	cloudAccount.Name = ctx.Payload.Name
	cloudAccount.Cloudprovider = ctx.Payload.Cloudprovider
	cloudAccount.UpstreamAccountID = ctx.Payload.UpstreamAccountID
	cloudAccount.AccountID = ctx.AccountID
	cloudAccount.AssumeRoleArn = ctx.Payload.AssumeRoleArn
	cloudAccount.AssumeRoleExternalID = ctx.Payload.AssumeRoleExternalID

	err := cdb.Add(ctx.Context, &cloudAccount)
	if err != nil {
		return ErrDatabaseError(err)
	}
	logger.Info("Created", "cloudaccount", fmt.Sprintf("%+v", cloudAccount))

	/*
		cloudEventID := cloudEvent.ID
		logger.Info("Saved CloudEvent", "ID", fmt.Sprintf("%+v", cloudEventID))
		logger.Info("Saved CloudEvent", "CloudEvent", fmt.Sprintf("%+v", cloudEvent))
		cloudEventFromDB := models.CloudEvent{}
		edb.Db.Where(&models.CloudEvent{ID: cloudEventID}).First(&cloudEventFromDB)
		logger.Info("CloudEvent from DB", "CloudEvent", fmt.Sprintf("%+v", cloudEventFromDB))

	*/
	cloudAccountID := cloudAccount.ID
	accountID := cloudAccount.AccountID

	// cloudAccountFromDB := models.CloudAccount{}
	// cdb.Db.Where(&models.CloudAccount{ID: cloudAccountID}).First(&cloudAccountFromDB)
	// OneCloudaccount(ctx context.Context, id int, accountID int) (*app.Cloudaccount, error) {
	cloudAccountFromDB, err := cdb.OneCloudaccount(ctx.Context, cloudAccountID, accountID)
	logger.Info("Loaded", "cloudaccountFromDb", fmt.Sprintf("%+v", cloudAccountFromDB))

	ctx.ResponseData.Header().Set("Location", app.CloudaccountHref(ctx.AccountID, cloudAccount.ID))
	return ctx.Created()

}
