package main

import (
	"github.com/goadesign/goa"
	"github.com/tleyden/zerocloud/app"
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

	// CloudeventController_Create: end_implement
	return nil
}
