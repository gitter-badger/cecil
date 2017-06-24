// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package controllers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goadesign/goa"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/tools"
	"gopkg.in/mailgun/mailgun-go.v1"
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
	requestContextLogger := core.NewContextLogger(ctx)

	// validate recipients
	{
		if len(ctx.Payload.Recipients) > 0 {
			for i, email := range ctx.Payload.Recipients {
				_, err := c.validateEmail(email)
				if err != nil {
					requestContextLogger.Error(
						fmt.Sprintf("Error while verifying #%v email (%q) in the recipients' list", i, email),
						"err", err,
					)
					return tools.ErrInvalidRequest(ctx, err)
				}
			}
		}
	}

	// validate minimum_lease_age
	minimumLeaseAge, err := time.ParseDuration(ctx.Payload.MinimumLeaseAge)
	if err != nil {
		return tools.ErrInvalidRequest(ctx, err)
	}
	fmt.Println("minimumLeaseAge", minimumLeaseAge)

	// fetch report, or status of report

	var leases []models.Lease
	err = c.cs.DB.
		Table("leases").
		Where("created_at <= ?", time.Now().UTC().Add(-minimumLeaseAge)).
		Where("cloudaccount_id = ?", uint(ctx.CloudaccountID)).
		Find(&leases).
		Error

	if err != nil {
		return tools.ErrInternal(ctx, err)
	}

	report := map[string]interface{}{
		"recipient_email":  "example@example.com",
		"aws_account_id":   "241535321535",
		"cecil_account_id": 123,

		"header_tracked":  []string{},
		"entries_tracked": [][]string{},

		"header_untracked": []string{},
		"entries_untracked": [][]string{
			{},
			{},
		},
	}

	report["header_tracked"] = []string{
		"Start Date",
		"Group UID",
		"Group Type",
		"Instances",
	}

	for _, lease := range leases {
		instanceIDList := []string{}
		instances, err := c.cs.ActiveInstancesForGroup(lease.AccountID, &lease.CloudaccountID, lease.GroupUID)
		if err != nil {
			// TODO: log error
			continue
		} else {
			for _, ins := range instances {
				instanceIDList = append(instanceIDList, ins.InstanceID)
			}
		}
		entry := []string{
			tools.Stringify(lease.CreatedAt),
			tools.Stringify(lease.GroupUID),
			tools.Stringify(lease.GroupType),
			strings.Join(instanceIDList, ", "),
		}
		report["entries_tracked"] = append(report["entries_tracked"].([][]string), entry)
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

	ctx.OK([]byte(emailBody))
	fmt.Println(emailBody)

	//////////////////////////////////////////////////////////////////////////////////////
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

func (c *ReportController) validateEmail(email string) (mailgun.EmailVerification, error) {

	// validate email
	newAccountInputEmail, err := c.cs.DefaultMailer().Client.ValidateEmail(email)
	if err != nil {
		return mailgun.EmailVerification{}, errors.New("error while verifying email; please retry")
	}
	if !newAccountInputEmail.IsValid {
		return mailgun.EmailVerification{}, fmt.Errorf("invalid email %q; please fix", email)
	}

	return newAccountInputEmail, nil
}

// ShowReport runs the showReport action.
func (c *ReportController) ShowReport(ctx *app.ShowReportReportContext) error {
	// fetch report, or status of report

	report := map[string]interface{}{
		"recipient_email":  "example@example.com",
		"aws_account_id":   "241535321535",
		"cecil_account_id": 123,

		"header_tracked": []string{"uno", "dos", "tres"},
		"entries_tracked": [][]string{
			{"entry_1_field1", "entry_1_field2", "entry_1_field3"},
			{"entry_2_field1", "entry_2_field2", "entry_2_field3"},
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

	fmt.Println(emailBody)

	return tools.JSONResponse(ctx, 200, map[string]interface{}{
		"status": "ready",
		"report": report,
	})
}
