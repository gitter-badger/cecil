package main

import (
	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// Create runs the create action.
func (c *CloudaccountController) CreateImpl(ctx *app.CreateCloudaccountContext) error {

	if err := ctx.Payload.Validate(); err != nil {
		return err
	}

	// Put your logic here
	a := models.CloudAccount{}
	a.Name = ctx.Payload.Name
	a.Cloudprovider = ctx.Payload.Cloudprovider
	a.UpstreamAccountID = ctx.Payload.UpstreamAccountID
	a.AccountID = ctx.AccountID
	a.AssumeRoleArn = ctx.Payload.AssumeRoleArn
	a.AssumeRoleExternalID = ctx.Payload.AssumeRoleExternalID

	err := cdb.Add(ctx.Context, &a)
	if err != nil {
		return ErrDatabaseError(err)
	}
	ctx.ResponseData.Header().Set("Location", app.CloudaccountHref(ctx.AccountID, a.ID))
	return ctx.Created()

}
