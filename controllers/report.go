package controllers

import (
	"github.com/goadesign/goa"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
	"github.com/tleyden/cecil/tools"
)

// ReportController implements the report resource.
type ReportController struct {
	*goa.Controller
	cs *core.Service
}

// NewReportController creates a report controller.
func NewReportController(service *goa.Service, cs *core.Service) *ReportController {
	return &ReportController{
		Controller: service.NewController("ReportController"),
		cs:         cs,
	}
}

// OrderInstancesReport runs the orderInstancesReport action.
func (c *ReportController) OrderInstancesReport(ctx *app.OrderInstancesReportReportContext) error {
	// validate each item in recipients array
	// parse minimum_lease_age (should be duration)

	// backgroud task:
	// get leases older than (now - minimum_lease_age)
	// if len(recipients) > 0, send email to recipients with HTML table
	// return:
	/*
	   {
	   "report":http://{{host}}:{{port}}/accounts/{{account_id}}/cloudaccounts/{{cloudaccount_id}}/reports/generated/UUID"
	   }
	*/
	return nil
}

// ShowReport runs the showReport action.
func (c *ReportController) ShowReport(ctx *app.ShowReportReportContext) error {
	// fetch report, or status of report

	report := map[string]interface{}{
		"recipient_email":  "example@example.com",
		"aws_account_id":   "241535321535",
		"cecil_account_id": 123,

		"header_tracked": []string{""},
		"entries_tracked": [][]string{
			{},
			{},
		},

		"header_untracked": []string{""},
		"entries_untracked": [][]string{
			{},
			{},
		},
	}

	emailBody, err := tools.CompileEmailTemplate(
		"report.html",
		report,
	)
	if err != nil {
		return err
	}
	_ = emailBody
	// TODO: send emails

	return tools.JSONResponse(ctx, 200, map[string]interface{}{
		"status": "ready",
		"report": report,
	})
}
