package main

import (
	"github.com/goadesign/goa"
	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// ErrDatabaseError is the error returned when a db query fails.
var ErrDatabaseError = goa.NewErrorClass("db_error", 500)

// Create runs the create action.
func (c *AccountController) CreateImpl(ctx *app.CreateAccountContext) error {

	a := models.Account{}
	a.Name = ctx.Payload.Name
	a.LeaseExpiresIn = ctx.Payload.LeaseExpiresIn
	a.LeaseExpiresInUnits = ctx.Payload.LeaseExpiresInUnits

	err := adb.Add(ctx.Context, &a)
	if err != nil {
		return ErrDatabaseError(err)
	}
	ctx.ResponseData.Header().Set("Location", app.AccountHref(a.ID))
	return ctx.Created()

}
