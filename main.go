//go:generate goagen bootstrap -d github.com/tleyden/zerocloud/design

package main

import (
	"github.com/goadesign/goa"
	goalog15 "github.com/goadesign/goa/logging/log15"
	"github.com/goadesign/goa/middleware"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

var adb *models.AccountDB
var cdb *models.CloudAccountDB
var edb *models.CloudEventDB

var logger log15.Logger

func main() {

	// Setup logger
	logger = log15.New()

	db, err := gorm.Open("sqlite3", "zerocloud.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.DropTable(
		&models.Account{},
		&models.CloudAccount{},
		&models.CloudEvent{},
	)
	db.AutoMigrate(
		&models.Account{},
		&models.CloudAccount{},
		&models.CloudEvent{},
	)

	adb = models.NewAccountDB(db)
	cdb = models.NewCloudAccountDB(db)
	edb = models.NewCloudEventDB(db)

	// Create service
	service := goa.New("zerocloud")
	service.WithLogger(goalog15.New(logger))

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "account" controller
	ac := NewAccountController(service)
	app.MountAccountController(service, ac)

	// Mount "cloud account" controller
	ca := NewCloudaccountController(service)
	app.MountCloudaccountController(service, ca)

	// Mount "cloud event" controller
	ce := NewCloudeventController(service)
	app.MountCloudeventController(service, ce)

	// Mount "aws" controller
	aws := NewAwsController(service)
	app.MountAwsController(service, aws)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}
