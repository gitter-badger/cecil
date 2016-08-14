package main

import (
	"github.com/goadesign/goa"
	"github.com/tleyden/zerocloud/app"
)

// AwsController implements the aws resource.
type AwsController struct {
	*goa.Controller
}

// NewAwsController creates a aws controller.
func NewAwsController(service *goa.Service) *AwsController {
	return &AwsController{Controller: service.NewController("AwsController")}
}

// Show runs the show action.
func (c *AwsController) Show(ctx *app.ShowAwsContext) error {
	return c.ShowImpl(ctx)

}
