//go:generate goagen bootstrap -d github.com/tleyden/cecil/design

package main

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/tleyden/cecil/goa-controllers/app"
)

func main() {
	// Create service
	service := goa.New("Cecil")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "account" controller
	c := NewAccountController(service)
	app.MountAccountController(service, c)
	// Mount "cloudaccount" controller
	c2 := NewCloudaccountController(service)
	app.MountCloudaccountController(service, c2)
	// Mount "email_action" controller
	c3 := NewEmailActionController(service)
	app.MountEmailActionController(service, c3)
	// Mount "swagger" controller
	c4 := NewSwaggerController(service)
	app.MountSwaggerController(service, c4)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}
