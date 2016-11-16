package main

import (
	"github.com/goadesign/goa"
	"github.com/tleyden/cecil/goa-controllers/app"
)

// CloudaccountController implements the cloudaccount resource.
type CloudaccountController struct {
	*goa.Controller
}

// NewCloudaccountController creates a cloudaccount controller.
func NewCloudaccountController(service *goa.Service) *CloudaccountController {
	return &CloudaccountController{Controller: service.NewController("CloudaccountController")}
}

// Add runs the add action.
func (c *CloudaccountController) Add(ctx *app.AddCloudaccountContext) error {
	// CloudaccountController_Add: start_implement

	// Put your logic here

	// CloudaccountController_Add: end_implement
	return nil
}

// AddEmailToWhitelist runs the addEmailToWhitelist action.
func (c *CloudaccountController) AddEmailToWhitelist(ctx *app.AddEmailToWhitelistCloudaccountContext) error {
	// CloudaccountController_AddEmailToWhitelist: start_implement

	// Put your logic here

	// CloudaccountController_AddEmailToWhitelist: end_implement
	return nil
}

// DownloadInitialSetupTemplate runs the downloadInitialSetupTemplate action.
func (c *CloudaccountController) DownloadInitialSetupTemplate(ctx *app.DownloadInitialSetupTemplateCloudaccountContext) error {
	// CloudaccountController_DownloadInitialSetupTemplate: start_implement

	// Put your logic here

	// CloudaccountController_DownloadInitialSetupTemplate: end_implement
	return nil
}

// DownloadRegionSetupTemplate runs the downloadRegionSetupTemplate action.
func (c *CloudaccountController) DownloadRegionSetupTemplate(ctx *app.DownloadRegionSetupTemplateCloudaccountContext) error {
	// CloudaccountController_DownloadRegionSetupTemplate: start_implement

	// Put your logic here

	// CloudaccountController_DownloadRegionSetupTemplate: end_implement
	return nil
}
