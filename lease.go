package main

import (
	"github.com/goadesign/goa"
	"github.com/tleyden/zerocloud/app"
)

// LeaseController implements the lease resource.
type LeaseController struct {
	*goa.Controller
}

// NewLeaseController creates a lease controller.
func NewLeaseController(service *goa.Service) *LeaseController {
	return &LeaseController{Controller: service.NewController("LeaseController")}
}

// List runs the list action.
func (c *LeaseController) List(ctx *app.ListLeaseContext) error {
	return c.ListImpl(ctx)
}
