package main

import (
	"fmt"
	"log"

	"github.com/goadesign/goa"
	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// CloudeventController implements the cloudevent resource.
type CloudeventController struct {
	*goa.Controller
}

// NewCloudeventController creates a cloudevent controller.
func NewCloudeventController(service *goa.Service) *CloudeventController {
	return &CloudeventController{Controller: service.NewController("CloudeventController")}
}

// Create runs the create action.
func (c *CloudeventController) Create(ctx *app.CreateCloudeventContext) error {
	// CloudeventController_Create: start_implement

	// Put your logic here

	log.Printf("aws account id: %v", ctx.Payload.AwsAccountID)

	// try to find the CloudAccount that has an upstream_account_id that matches param
	cloudAccount := models.CloudAccount{}
	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: ctx.Payload.AwsAccountID}).First(&cloudAccount)
	log.Printf("cloudAccount: %+v", cloudAccount)
	if cloudAccount.ID == 0 {
		// ctx.BadRequest(fmt.Errorf("Could not find CloudAccount with upstream provider account id: %v", ctx.Payload.AwsAccountID))
		ctx.ResponseData.Service.Send(ctx.Context, 400, fmt.Sprintf("Could not find CloudAccount with upstream provider account id: %v", ctx.Payload.AwsAccountID))
		return nil
	}

	// CloudeventController_Create: end_implement
	return nil
}
