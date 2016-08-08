package main

import (
	"log"

	"github.com/goadesign/goa"
	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
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

	log.Printf("aws account id: %v", ctx.Payload.AwsAccountID)

	// try to find the CloudAccount that has an upstream_account_id that matches param
	//  cdb.Db.Model()
	// rows, err := cdb.Db.Model(&models.CloudAccount{}).Where("upstream_account_id = ?", ctx.Payload.AwsAccountID).Select("id").Rows()
	cloudAccount := models.CloudAccount{}
	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: ctx.Payload.AwsAccountID}).First(&cloudAccount)
	log.Printf("cloudAccount: %+v", cloudAccount)

	/*
		if err != nil {
			log.Printf("Got err: %s", err)
			return err
		}
		log.Printf("iterating over rows")
		for rows.Next() {
			log.Printf("row ..")
			var id string
			if err := rows.Scan(&id); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("id is %s\n", id)
		}
		log.Printf("/iterating over rows")
	*/

	// CloudeventController_Create: end_implement
	return nil
}
