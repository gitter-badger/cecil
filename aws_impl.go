package main

import (
	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// Find the CloudAccount corresponding to the AWS account ID in ctx.AwsAccountID
func (c *AwsController) ShowImpl(ctx *app.ShowAwsContext) error {

	cloudAccount := models.CloudAccount{}
	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: ctx.AwsAccountID}).First(&cloudAccount)
	return ctx.OK(cloudAccount.CloudAccountToCloudaccount())

}
