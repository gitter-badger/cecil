package controllers

import (
	"fmt"
	"time"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
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

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return core.ErrUnauthorized(ctx, core.ErrorUnauthorized)
	}

	_, err = c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	// fetch leases for account
	leases, err := c.cs.LeasesForAccount(ctx.AccountID, ctx.Terminated)
	if err != nil {
		requestContextLogger.Error("Error fetching leases for account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, "no leases found")
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	return core.JSONResponse(ctx, 200, leases)
}

// ListLeasesForCloudaccount handles the endpoint used to list all leases for a cloudaccount.
func (c *LeasesController) ListLeasesForCloudaccount(ctx *app.ListLeasesForCloudaccountLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return core.ErrUnauthorized(ctx, core.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	// fetch leases for cloudaccount
	leases, err := c.cs.LeasesForCloudaccount(ctx.CloudaccountID, ctx.Terminated)
	if err != nil {
		requestContextLogger.Error("Error fetching leases for cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, "no leases found")
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	return core.JSONResponse(ctx, 200, leases)
}

// Show handles the endpoint used to show the info about a specific lease.
func (c *LeasesController) Show(ctx *app.ShowLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return core.ErrUnauthorized(ctx, core.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	cloudaccountIsSpecified := ctx.CloudaccountID > 0
	var cloudaccount *core.Cloudaccount

	if cloudaccountIsSpecified {
		cloudaccount, err = c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
		if err != nil {
			requestContextLogger.Error("Error fetching cloudaccount", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
			}
			return core.ErrInternal(ctx, core.ErrorInternal)
		}

		// check whether everything is consistent
		if !account.IsOwnerOf(cloudaccount) {
			requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
			return core.ErrNotFound(ctx, "cloud account not found")
		}
	}

	// fetch lease
	lease, err := c.cs.FetchLeaseByID(ctx.LeaseID)
	if err != nil {
		requestContextLogger.Error("Error fetching lease", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, "lease not found")
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	if cloudaccountIsSpecified {
		if !cloudaccount.IsOwnerOf(lease) {
			requestContextLogger.Error(fmt.Sprintf("Cloudaccount %v is not owner of lease %v", cloudaccount.ID, lease.ID))
			return core.ErrInvalidRequest(ctx, "lease not found")
		}
	}

	if !account.IsOwnerOfLease(lease) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of lease %v", account.ID, lease.ID))
		return core.ErrInvalidRequest(ctx, "lease not found")
	}

	return core.JSONResponse(ctx, 200, lease)
}

// Terminate handles the endpoint used to terminate an instance.
func (c *LeasesController) Terminate(ctx *app.TerminateLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return core.ErrUnauthorized(ctx, core.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	cloudaccountIsSpecified := ctx.CloudaccountID > 0
	var cloudaccount *core.Cloudaccount

	if cloudaccountIsSpecified {
		cloudaccount, err = c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
		if err != nil {
			requestContextLogger.Error("Error fetching cloudaccount", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
			}
			return core.ErrInternal(ctx, core.ErrorInternal)
		}

		// check whether everything is consistent
		if !account.IsOwnerOf(cloudaccount) {
			requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
			return core.ErrNotFound(ctx, "cloud account not found")
		}
	}

	// fetch lease
	lease, err := c.cs.FetchLeaseByID(ctx.LeaseID)
	if err != nil {
		requestContextLogger.Error("Error fetching lease", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, "lease not found")
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	if cloudaccountIsSpecified {
		if !cloudaccount.IsOwnerOf(lease) {
			requestContextLogger.Error(fmt.Sprintf("Cloudaccount %v is not owner of lease %v", cloudaccount.ID, lease.ID))
			return core.ErrInvalidRequest(ctx, "lease not found")
		}
	}

	if !account.IsOwnerOfLease(lease) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of lease %v", account.ID, lease.ID))
		return core.ErrInvalidRequest(ctx, "lease not found")
	}

	c.cs.TerminatorQueue.TaskQueue <- core.TerminatorTask{Lease: *lease}

	resp := core.HMI{
		"message":      "Termination request received",
		"resourceID":   lease.ResourceID,
		"resourceType": lease.ResourceType,
		"lease_id":     lease.ID,
	}

	var instance core.InstanceResource
	if lease.IsInstance() {
		raw, err := c.cs.ResourceOf(lease)
		if err != nil {
			return err
		}
		instance = raw.(core.InstanceResource)
	}

	var stack core.StackResource
	if lease.IsStack() {
		raw, err := c.cs.ResourceOf(lease)
		if err != nil {
			return err
		}
		stack = raw.(core.StackResource)
	}

	if lease.IsInstance() {
		resp["instance_id"] = instance.InstanceID
		resp["instance_type"] = instance.InstanceType
	}

	if lease.IsStack() {
		resp["stack_id"] = stack.StackID
		resp["stack_name"] = stack.StackName
	}

	return core.JSONResponse(ctx, 202, resp)
}

// DeleteFromDB handles the endpoint used to remove a lease from the DB.
func (c *LeasesController) DeleteFromDB(ctx *app.DeleteFromDBLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return core.ErrUnauthorized(ctx, core.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	cloudaccountIsSpecified := ctx.CloudaccountID > 0
	var cloudaccount *core.Cloudaccount

	if cloudaccountIsSpecified {
		cloudaccount, err = c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
		if err != nil {
			requestContextLogger.Error("Error fetching cloudaccount", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
			}
			return core.ErrInternal(ctx, core.ErrorInternal)
		}

		// check whether everything is consistent
		if !account.IsOwnerOf(cloudaccount) {
			requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
			return core.ErrNotFound(ctx, "cloud account not found")
		}
	}

	// fetch lease
	lease, err := c.cs.FetchLeaseByID(ctx.LeaseID)
	if err != nil {
		requestContextLogger.Error("Error fetching lease", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, "lease not found")
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	if cloudaccountIsSpecified {
		if !cloudaccount.IsOwnerOf(lease) {
			requestContextLogger.Error(fmt.Sprintf("Cloudaccount %v is not owner of lease %v", cloudaccount.ID, lease.ID))
			return core.ErrInvalidRequest(ctx, "lease not found")
		}
	}

	if !account.IsOwnerOfLease(lease) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of lease %v", account.ID, lease.ID))
		return core.ErrInvalidRequest(ctx, "lease not found")
	}

	// if the lease is not terminated, cannot delete
	if lease.TerminatedAt == nil {
		requestContextLogger.Error("Lease is not terminated; cannot delete from DB")
		return core.ErrInvalidRequest(ctx, "cannot delete non-terminated lease")
	}

	// delete lease
	if err := c.cs.DB.Delete(&lease).Error; err != nil {
		requestContextLogger.Error("Error deleting lease from DB", "err", err)
		return core.ErrInternal(ctx, "error while deleting lease; please retry")
	}

	resp := core.HMI{
		"message":      "Lease deleted from DB",
		"resourceID":   lease.ResourceID,
		"resourceType": lease.ResourceType,
		"lease_id":     lease.ID,
	}

	var instance core.InstanceResource
	if lease.IsInstance() {
		raw, err := c.cs.ResourceOf(lease)
		if err != nil {
			return err
		}
		instance = raw.(core.InstanceResource)
	}

	var stack core.StackResource
	if lease.IsStack() {
		raw, err := c.cs.ResourceOf(lease)
		if err != nil {
			return err
		}
		stack = raw.(core.StackResource)
	}

	if lease.IsInstance() {
		resp["instance_id"] = instance.InstanceID
		resp["instance_type"] = instance.InstanceType
	}

	if lease.IsStack() {
		resp["stack_id"] = stack.StackID
		resp["stack_name"] = stack.StackName
	}

	return core.JSONResponse(ctx, 200, resp)
}

// SetExpiry handles the endpoint used to set the expiry of a lease (if not already terminated).
func (c *LeasesController) SetExpiry(ctx *app.SetExpiryLeasesContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return core.ErrUnauthorized(ctx, core.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	cloudaccountIsSpecified := ctx.CloudaccountID > 0
	var cloudaccount *core.Cloudaccount

	if cloudaccountIsSpecified {
		cloudaccount, err = c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
		if err != nil {
			requestContextLogger.Error("Error fetching cloudaccount", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
			}
			return core.ErrInternal(ctx, core.ErrorInternal)
		}

		// check whether everything is consistent
		if !account.IsOwnerOf(cloudaccount) {
			requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
			return core.ErrNotFound(ctx, "cloud account not found")
		}
	}

	// fetch lease
	lease, err := c.cs.FetchLeaseByID(ctx.LeaseID)
	if err != nil {
		requestContextLogger.Error("Error fetching lease", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, "lease not found")
		}
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	if cloudaccountIsSpecified {
		if !cloudaccount.IsOwnerOf(lease) {
			requestContextLogger.Error(fmt.Sprintf("Cloudaccount %v is not owner of lease %v", cloudaccount.ID, lease.ID))
			return core.ErrInvalidRequest(ctx, "lease not found")
		}
	}

	if !account.IsOwnerOfLease(lease) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of lease %v", account.ID, lease.ID))
		return core.ErrInvalidRequest(ctx, "lease not found")
	}

	// if the lease is terminated, the expiry already happened
	if lease.TerminatedAt != nil {
		requestContextLogger.Error("Lease is terminated; cannot set expiry on terminated lease")
		return core.ErrInvalidRequest(ctx, "lease already expired")
	}

	// validate new expiry
	newExpiryIsInThePast := time.Now().UTC().After(ctx.ExpiresAt)
	if newExpiryIsInThePast {
		requestContextLogger.Error("Error newExpiryIsInThePast")
		return core.ErrInvalidRequest(ctx, "cannot set expiry to the past; use UTC RFC3339; e.g. 2016-12-17T22:37:19Z")
	}

	// set new expiry
	lease.ExpiresAt = ctx.ExpiresAt

	// save lease to db
	if err := c.cs.DB.Save(&lease).Error; err != nil {
		requestContextLogger.Error("Error saving updated lease", "err", err)
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	resp := core.HMI{
		"message":      "New expiry set",
		"resourceID":   lease.ResourceID,
		"resourceType": lease.ResourceType,
		"lease_id":     lease.ID,
	}

	var instance core.InstanceResource
	if lease.IsInstance() {
		raw, err := c.cs.ResourceOf(lease)
		if err != nil {
			return err
		}
		instance = raw.(core.InstanceResource)
	}

	var stack core.StackResource
	if lease.IsStack() {
		raw, err := c.cs.ResourceOf(lease)
		if err != nil {
			return err
		}
		stack = raw.(core.StackResource)
	}

	if lease.IsInstance() {
		resp["instance_id"] = instance.InstanceID
		resp["instance_type"] = instance.InstanceType
	}

	if lease.IsStack() {
		resp["stack_id"] = stack.StackID
		resp["stack_name"] = stack.StackName
	}

	return core.JSONResponse(ctx, 200, resp)
}
