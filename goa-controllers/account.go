package main

import (
	"github.com/goadesign/goa"
	"github.com/tleyden/cecil/goa-controllers/app"
)

// AccountController implements the account resource.
type AccountController struct {
	*goa.Controller
}

// NewAccountController creates a account controller.
func NewAccountController(service *goa.Service) *AccountController {
	return &AccountController{Controller: service.NewController("AccountController")}
}

// Create runs the create action.
func (c *AccountController) Create(ctx *app.CreateAccountContext) error {
	// AccountController_Create: start_implement

	// Put your logic here

	// AccountController_Create: end_implement
	return nil
}

// Show runs the show action.
func (c *AccountController) Show(ctx *app.ShowAccountContext) error {
	// AccountController_Show: start_implement

	// Put your logic here

	// AccountController_Show: end_implement
	return nil
}

// Verify runs the verify action.
func (c *AccountController) Verify(ctx *app.VerifyAccountContext) error {
	// AccountController_Verify: start_implement

	// Put your logic here

	// AccountController_Verify: end_implement
	return nil
}
