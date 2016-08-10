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
		ctx.ResponseData.Service.Send(
			ctx.Context,
			400,
			fmt.Sprintf("Could not find CloudAccount with upstream provider account id: %v", awsAccountId),
		)
		return nil
	}

	/*
		Field("sqs_payload", gorma.Text)
		Field("sqs_timestamp", gorma.Timestamp)
		Field("cw_event_detail_type", gorma.String)
		Field("cw_event_source", gorma.String)
		Field("cw_event_timestamp", gorma.Timestamp)
		Field("cw_event_region", gorma.String)
		Field("cw_event_detail_instance_id", gorma.String)
		Field("cw_event_detail_state", gorma.String)

	*/

	// Save the raw CloudEvent to the database
	e := models.CloudEvent{}
	e.AwsAccountID = awsAccountId
	e.CloudAccountID = cloudAccount.ID
	e.AccountID = cloudAccount.AccountID
	err := edb.Add(ctx.Context, &e)
	if err != nil {
		return ErrDatabaseError(err)
	}

	// Create a Lease object that references this (immutable) CloudEvent and expires
	// based on the settings in the Account
	// TODO: or can this be an AfterCreate callback on the CloudEvent?
	// file:///Users/tleyden/DevLibraries/gorm/callbacks.html
	err = createLease(e)
	if err != nil {
		return ErrDatabaseError(err)
	}

	// TODO: should this return the path to the cloudevent .. should there even be one?
	/// ctx.ResponseData.Header().Set("Location", app.CloudeventHref(ctx.AccountID, a.ID))

	return ctx.Created()

}

func createLease(cloudEvent models.CloudEvent) error {
	// TODO
	return nil
}
