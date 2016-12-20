package controllers

import (
	"encoding/json"
	"time"

	"github.com/goadesign/goa"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
)

// RootController implements the root resource.
type RootController struct {
	*goa.Controller
	startedAt time.Time
}

// NewRootController creates a root controller.
func NewRootController(service *goa.Service, startedAt time.Time) *RootController {
	return &RootController{
		Controller: service.NewController("RootController"),
		startedAt:  startedAt,
	}
}

// Show handles the endpoint used to show info about Cecil.
func (c *RootController) Show(ctx *app.ShowRootContext) error {
	var APIInfo struct {
		Name   string `json:"name"`
		Uptime string `json:"uptime"`
		Time   string `json:"time"`
	}

	uptime := time.Now().UTC().Sub(c.startedAt)

	APIInfo.Name = c.Service.Name
	APIInfo.Uptime = uptime.String()
	APIInfo.Time = time.Now().UTC().Format(time.RFC3339)

	resp, err := json.MarshalIndent(APIInfo, "", "  ")
	if err != nil {
		core.Logger.Error("Error while marshaling APIInfo", "err", err)
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	return ctx.OK(resp)
}
