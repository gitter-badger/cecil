//go:generate goagen bootstrap -d github.com/tleyden/zerocloud/design

package main

import (
	"log"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

var adb *models.AccountDB

func main() {

	db, err := gorm.Open("sqlite3", "zerocloud.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.DropTable(&models.Account{}, &models.CloudAccount{})
	db.AutoMigrate(&models.Account{}, &models.CloudAccount{})

	adb = models.NewAccountDB(db)
	log.Printf("adb: %v", adb)

	// Create service
	service := goa.New("zerocloud")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "account" controller
	c := NewAccountController(service)
	app.MountAccountController(service, c)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}
