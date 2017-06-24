// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package controllers

import (
	"fmt"
	"time"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
)

// LeasesController implements the leases resource.
type LeasesController struct {
	*goa.Controller
	cs *core.Service
}

// NewLeasesController creates a leases controller.
func NewLeasesController(service *goa.Service, cs *core.Service) *LeasesController {
	return &LeasesController{
		Controller: service.NewController("LeasesController"),
		cs:         cs,
	}
}

// ListLeasesForAccount handles the endpoint used to list all leases for an account.
func (c *LeasesController) ListLeasesForAccount(ctx *app.ListLeasesForAccountLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	// fetch leases for account
	leases, err := c.cs.LeasesForAccount(ctx.AccountID, ctx.Terminated)
	if err != nil {
		requestContextLogger.Error("Error fetching leases for account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, "no leases found")
		}
		return tools.ErrInternal(ctx, tools.ErrorInternal)
	}

	return tools.JSONResponse(ctx, 200, leases)
}

// ListLeasesForCloudaccount handles the endpoint used to list all leases for a cloudaccount.
func (c *LeasesController) ListLeasesForCloudaccount(ctx *app.ListLeasesForCloudaccountLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	// fetch leases for cloudaccount
	leases, err := c.cs.LeasesForCloudaccount(ctx.CloudaccountID, ctx.Terminated)
	if err != nil {
		requestContextLogger.Error("Error fetching leases for cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, "no leases found")
		}
		return tools.ErrInternal(ctx, tools.ErrorInternal)
	}

	return tools.JSONResponse(ctx, 200, leases)
}

// Show handles the endpoint used to show the info about a specific lease.
func (c *LeasesController) Show(ctx *app.ShowLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	account := core.ContextAccount(ctx)
	cloudaccount := core.ContextCloudaccount(ctx)

	cloudaccountIsSpecified := ctx.CloudaccountID > 0

	// fetch lease
	lease, err := c.cs.GetLeaseByID(ctx.LeaseID)
	if err != nil {
		requestContextLogger.Error("Error fetching lease", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, "lease not found")
		}
		return tools.ErrInternal(ctx, tools.ErrorInternal)
	}

	if cloudaccountIsSpecified {
		if !cloudaccount.IsOwnerOf(lease) {
			requestContextLogger.Error(fmt.Sprintf("Cloudaccount %v is not owner of lease %v", cloudaccount.ID, lease.ID))
			return tools.ErrInvalidRequest(ctx, "lease not found")
		}
	}

	if !account.IsOwnerOfLease(lease) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of lease %v", account.ID, lease.ID))
		return tools.ErrInvalidRequest(ctx, "lease not found")
	}

	return tools.JSONResponse(ctx, 200, lease)
}

// Terminate handles the endpoint used to terminate an instance.
func (c *LeasesController) Terminate(ctx *app.TerminateLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	account := core.ContextAccount(ctx)
	cloudaccount := core.ContextCloudaccount(ctx)

	cloudaccountIsSpecified := ctx.CloudaccountID > 0

	// fetch lease
	lease, err := c.cs.GetLeaseByID(ctx.LeaseID)
	if err != nil {
		requestContextLogger.Error("Error fetching lease", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, "lease not found")
		}
		return tools.ErrInternal(ctx, tools.ErrorInternal)
	}

	if cloudaccountIsSpecified {
		if !cloudaccount.IsOwnerOf(lease) {
			requestContextLogger.Error(fmt.Sprintf("Cloudaccount %v is not owner of lease %v", cloudaccount.ID, lease.ID))
			return tools.ErrInvalidRequest(ctx, "lease not found")
		}
	}

	if !account.IsOwnerOfLease(lease) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of lease %v", account.ID, lease.ID))
		return tools.ErrInvalidRequest(ctx, "lease not found")
	}

	c.cs.Queues().TerminatorQueue().PushTask(tasks.TerminatorTask{Lease: *lease})

	resp := tools.HMI{
		"message":    "Termination request received",
		"group_uid":  lease.GroupUID,
		"group_type": lease.GroupType.String(),
		"lease_id":   lease.ID,
	}

	return tools.JSONResponse(ctx, 202, resp)
}

// DeleteFromDB handles the endpoint used to remove a lease from the DB.
func (c *LeasesController) DeleteFromDB(ctx *app.DeleteFromDBLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	account := core.ContextAccount(ctx)
	cloudaccount := core.ContextCloudaccount(ctx)

	cloudaccountIsSpecified := ctx.CloudaccountID > 0

	// fetch lease
	lease, err := c.cs.GetLeaseByID(ctx.LeaseID)
	if err != nil {
		requestContextLogger.Error("Error fetching lease", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, "lease not found")
		}
		return tools.ErrInternal(ctx, tools.ErrorInternal)
	}

	if cloudaccountIsSpecified {
		if !cloudaccount.IsOwnerOf(lease) {
			requestContextLogger.Error(fmt.Sprintf("Cloudaccount %v is not owner of lease %v", cloudaccount.ID, lease.ID))
			return tools.ErrInvalidRequest(ctx, "lease not found")
		}
	}

	if !account.IsOwnerOfLease(lease) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of lease %v", account.ID, lease.ID))
		return tools.ErrInvalidRequest(ctx, "lease not found")
	}

	// if the lease is not terminated, cannot delete
	if lease.TerminatedAt == nil {
		requestContextLogger.Error("Lease is not terminated; cannot delete from DB")
		return tools.ErrInvalidRequest(ctx, "cannot delete non-terminated lease")
	}

	// delete lease
	if err := c.cs.DB.Delete(&lease).Error; err != nil {
		requestContextLogger.Error("Error deleting lease from DB", "err", err)
		return tools.ErrInternal(ctx, "error while deleting lease; please retry")
	}

	resp := tools.HMI{
		"message":    "Lease deleted from DB",
		"group_uid":  lease.GroupUID,
		"group_type": lease.GroupType.String(),
		"lease_id":   lease.ID,
	}

	return tools.JSONResponse(ctx, 200, resp)
}

// SetExpiry handles the endpoint used to set the expiry of a lease (if not already terminated).
func (c *LeasesController) SetExpiry(ctx *app.SetExpiryLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	account := core.ContextAccount(ctx)
	cloudaccount := core.ContextCloudaccount(ctx)

	cloudaccountIsSpecified := ctx.CloudaccountID > 0

	// fetch lease
	lease, err := c.cs.GetLeaseByID(ctx.LeaseID)
	if err != nil {
		requestContextLogger.Error("Error fetching lease", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, "lease not found")
		}
		return tools.ErrInternal(ctx, tools.ErrorInternal)
	}

	if cloudaccountIsSpecified {
		if !cloudaccount.IsOwnerOf(lease) {
			requestContextLogger.Error(fmt.Sprintf("Cloudaccount %v is not owner of lease %v", cloudaccount.ID, lease.ID))
			return tools.ErrInvalidRequest(ctx, "lease not found")
		}
	}

	if !account.IsOwnerOfLease(lease) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of lease %v", account.ID, lease.ID))
		return tools.ErrInvalidRequest(ctx, "lease not found")
	}

	// if the lease is terminated, the expiry already happened
	if lease.TerminatedAt != nil {
		requestContextLogger.Error("Lease is terminated; cannot set expiry on terminated lease")
		return tools.ErrInvalidRequest(ctx, "lease already expired")
	}

	// validate new expiry
	newExpiryIsInThePast := time.Now().UTC().After(ctx.ExpiresAt)
	if newExpiryIsInThePast {
		requestContextLogger.Error("Error newExpiryIsInThePast")
		return tools.ErrInvalidRequest(ctx, "cannot set expiry to the past; use UTC RFC3339; e.g. 2016-12-17T22:37:19Z")
	}

	// set new expiry
	lease.ExpiresAt = ctx.ExpiresAt

	// save lease to db
	if err := c.cs.DB.Save(&lease).Error; err != nil {
		requestContextLogger.Error("Error saving updated lease", "err", err)
		return tools.ErrInternal(ctx, tools.ErrorInternal)
	}

	resp := tools.HMI{
		"message":    "New expiry set",
		"group_uid":  lease.GroupUID,
		"group_type": lease.GroupType.String(),
		"lease_id":   lease.ID,
	}

	return tools.JSONResponse(ctx, 200, resp)
}
