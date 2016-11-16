//go:generate goagen bootstrap -d github.com/tleyden/cecil/design

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/goadesign/goa"
	goalog15 "github.com/goadesign/goa/logging/log15"
	"github.com/goadesign/goa/middleware"
	"github.com/tleyden/cecil/controllers"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
)

func main() {
	flag.BoolVar(&core.DropAllTables, "drop-all-tables", false, "If passed, drops all tables")
	flag.Parse()

	if core.DropAllTables {
		fmt.Println("You are about to drop all tables from DB;\nAre you sure? [N/y]")
		isSure := core.AskForConfirmation()
		if isSure {
			fmt.Println("Tables WILL BE dropped.")
		} else {
			fmt.Println("Tables will NOT be dropped.")
			core.DropAllTables = false
		}
	}

	// Create service
	service := goa.New("Cecil REST API")
	coreService := core.NewService()

	coreService.SetupAndRun()
	defer coreService.Stop(true)

	// make Goa use the core.Logger
	service.WithLogger(goalog15.New(core.Logger))

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// create the jwt middleware
	jwtMiddleware, err := coreService.NewJWTMiddleware()
	if err != nil {
		core.Logger.Error("Error while creating jwtMiddleware", "error", err)
		return
	}
	// mount the jwt middleware
	app.UseJWTMiddleware(service, jwtMiddleware)

	// Mount "root" controller
	c5 := controllers.NewRootController(service, time.Now().UTC())
	app.MountRootController(service, c5)
	// Mount "account" controller
	c := controllers.NewAccountController(service, coreService)
	app.MountAccountController(service, c)
	// Mount "cloudaccount" controller
	c2 := controllers.NewCloudaccountController(service, coreService)
	app.MountCloudaccountController(service, c2)
	// Mount "email_action" controller
	c3 := controllers.NewEmailActionController(service, coreService)
	app.MountEmailActionController(service, c3)
	// Mount "swagger" controller
	c4 := controllers.NewSwaggerController(service)
	app.MountSwaggerController(service, c4)

	// Start service
	if err := service.ListenAndServe(coreService.Config.Server.Port); err != nil {
		service.LogError("startup", "err", err)
	}
}
