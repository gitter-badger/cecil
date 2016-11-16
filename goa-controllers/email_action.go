package main

import (
	"github.com/goadesign/goa"
	"github.com/tleyden/cecil/goa-controllers/app"
)

// EmailActionController implements the email_action resource.
type EmailActionController struct {
	*goa.Controller
}

// NewEmailActionController creates a email_action controller.
func NewEmailActionController(service *goa.Service) *EmailActionController {
	return &EmailActionController{Controller: service.NewController("EmailActionController")}
}

// Actions runs the actions action.
func (c *EmailActionController) Actions(ctx *app.ActionsEmailActionContext) error {
	// EmailActionController_Actions: start_implement

	// Put your logic here

	// EmailActionController_Actions: end_implement
	return nil
}
