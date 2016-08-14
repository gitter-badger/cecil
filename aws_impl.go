package main

import (
	"fmt"

	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// LookupCloudaccount runs the lookup_cloudaccount action.
func (c *AwsController) ShowImpl(ctx *app.ShowAwsContext) error {

	// This is a little confusing that there is an app.Cloudaccount and
	// a models.CloudAccount that are very similar (identical?).  I'm
	// passing the app.Cloudaccount .. hoping it will work.

	logger.Info("aws controller show impl called", "ctx", fmt.Sprintf("%+v", ctx))

	res := app.Cloudaccount{}

	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: ctx.AwsAccountID}).First(&res)

	logger.Info("Found cloudaccount", "cloudaccount", fmt.Sprintf("%+v", res))

	res2 := models.CloudAccount{}

	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: ctx.AwsAccountID}).First(&res2)

	res3 := res2.CloudAccountToCloudaccount()

	logger.Info("Found cloudaccount res2", "cloudaccount", fmt.Sprintf("%+v", res3))

	return ctx.OK(res3)
}
