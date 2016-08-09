package main

import (
	"fmt"

	"github.com/goadesign/goa"
	"github.com/tleyden/zerocloud/app"
)

// CloudaccountController implements the cloudaccount resource.
type CloudaccountController struct {
	*goa.Controller
}

// NewCloudaccountController creates a cloudaccount controller.
func NewCloudaccountController(service *goa.Service) *CloudaccountController {
	return &CloudaccountController{Controller: service.NewController("CloudaccountController")}
}

// Create runs the create action.
func (c *CloudaccountController) Create(ctx *app.CreateCloudaccountContext) error {
	// CloudaccountController_Delete: start_implement

	// Put your logic here
	return c.CreateImpl(ctx)

	// CloudaccountController_Delete: end_implement

}

// Delete runs the delete action.
func (c *CloudaccountController) Delete(ctx *app.DeleteCloudaccountContext) error {
	// CloudaccountController_Delete: start_implement

	// Put your logic here

	// CloudaccountController_Delete: end_implement
	return nil
}

// List runs the list action.
func (c *CloudaccountController) List(ctx *app.ListCloudaccountContext) error {
	// CloudaccountController_List: start_implement

	// Put your logic here
	fmt.Printf("cloudaccount list called")

	// CloudaccountController_List: end_implement
	res := app.CloudaccountCollection{}
	return ctx.OK(res)
}

// Show runs the show action.
func (c *CloudaccountController) Show(ctx *app.ShowCloudaccountContext) error {
	// CloudaccountController_Show: start_implement

	// Put your logic here

	// CloudaccountController_Show: end_implement
	res := &app.Cloudaccount{}
	return ctx.OK(res)
}

// Update runs the update action.
func (c *CloudaccountController) Update(ctx *app.UpdateCloudaccountContext) error {
	// CloudaccountController_Update: start_implement

	// Put your logic here

	// CloudaccountController_Update: end_implement
	return nil
}
