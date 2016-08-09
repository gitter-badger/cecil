package main

import (
	"fmt"

	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// Create runs the create action.
func (c *CloudeventController) CreateImpl(ctx *app.CreateCloudeventContext) error {

	logger.Info("Create CloudEvent", "aws_account_id", ctx.Payload.AwsAccountID)

	// try to find the CloudAccount that has an upstream_account_id that matches param
	cloudAccount := models.CloudAccount{}
	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: ctx.Payload.AwsAccountID}).First(&cloudAccount)
	logger.Info("Found CloudAccount", "CloudAccount", fmt.Sprintf("%+v", cloudAccount))

	// Make sure we found a valid CloudAccount, otherwise abort
	if cloudAccount.ID == 0 {
		ctx.ResponseData.Service.Send(ctx.Context, 400, fmt.Sprintf("Could not find CloudAccount with upstream provider account id: %v", ctx.Payload.AwsAccountID))
		return nil
	}

	return nil

}
