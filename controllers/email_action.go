// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package controllers

import (
	"fmt"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/tasks"
	"github.com/tleyden/cecil/tools"
)

// EmailActionController implements the email_action resource.
type EmailActionController struct {
	*goa.Controller
	cs *core.Service
}

// NewEmailActionController creates a email_action controller.
func NewEmailActionController(service *goa.Service, cs *core.Service) *EmailActionController {
	return &EmailActionController{
		Controller: service.NewController("EmailActionController"),
		cs:         cs,
	}
}

// Actions handles the endpoint used to receive email_actions (i.e. links sent in emails that make perform specfic actions on leases).
func (c *EmailActionController) Actions(ctx *app.ActionsEmailActionContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	err := c.cs.EmailActionVerifySignatureParams(ctx.LeaseUUID.String(), ctx.GroupUIDHash, ctx.Action, ctx.Tok, ctx.Sig)
	if err != nil {
		requestContextLogger.Error("Signature verification error", "err", err)
		return tools.ErrInvalidRequest(ctx, "corrupted action link")
	}

	var lease models.Lease
	var resp = make(tools.HMI)

	switch ctx.Action {
	case "approve":
		core.Logger.Info("Approval of lease initiated", "GroupUID", lease.GroupUID)

		leaseToBeApproved, err := c.cs.GetLeaseByUUID(ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return tools.ErrInvalidRequest(ctx, fmt.Sprintf("lease for group with id %v does not exist", lease.GroupUID))
			}
			return tools.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for group with id %v. See logs for details", lease.GroupUID))
		}

		if leaseToBeApproved == nil {
			requestContextLogger.Error("leaseToBeApproved == nil", "err", err)
			return tools.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for group with id %v. See logs for details", lease.GroupUID))
		}

		if leaseToBeApproved.TokenOnce != ctx.Tok {
			requestContextLogger.Error("Wrong tokenOnce; link already used/expired")
			return tools.ErrNotFound(ctx, "link expired")
		}

		if leaseToBeApproved.IsExpired() {
			return tools.ErrInvalidRequest(ctx, "lease already expired")
		}

		c.cs.Queues().ExtenderQueue().PushTask(tasks.ExtenderTask{
			Lease:     *leaseToBeApproved,
			Approving: true,
		})

		resp = tools.HMI{
			"message":    "Approval request received",
			"group_type": leaseToBeApproved.GroupType.String(),
			"group_uid":  leaseToBeApproved.GroupUID,
			"lease_id":   leaseToBeApproved.ID,
		}
		lease = *leaseToBeApproved

	case "extend":
		requestContextLogger.Info("Extension of lease initiated", "instance_id", lease.GroupUID)

		leaseToBeExtended, err := c.cs.GetLeaseByUUID(ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return tools.ErrInvalidRequest(ctx, fmt.Sprintf("lease for group with id %v does not exist", lease.GroupUID))
			}
			return tools.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for group with id %v. See logs for details", lease.GroupUID))
		}

		if leaseToBeExtended == nil {
			requestContextLogger.Error("leaseToBeExtended == nil", "err", err)
			return tools.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for group with id %v. See logs for details", lease.GroupUID))
		}

		if leaseToBeExtended.TokenOnce != ctx.Tok {
			requestContextLogger.Error("Wrong tokenOnce; link already used/expired")
			return tools.ErrNotFound(ctx, "link expired")
		}

		if leaseToBeExtended.IsExpired() {
			return tools.ErrInvalidRequest(ctx, "lease already expired")
		}

		c.cs.Queues().ExtenderQueue().PushTask(tasks.ExtenderTask{
			Lease:     *leaseToBeExtended,
			Approving: false,
		})

		resp = tools.HMI{
			"message":    "Extension request received",
			"group_type": leaseToBeExtended.GroupType.String(),
			"group_uid":  leaseToBeExtended.GroupUID,
			"lease_id":   leaseToBeExtended.ID,
		}
		lease = *leaseToBeExtended
	case "terminate":
		requestContextLogger.Info("Termination of lease/stack initiated", "instance_id", lease.GroupUID)

		leaseToBeTerminated, err := c.cs.GetLeaseByUUID(ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return tools.ErrInvalidRequest(ctx, fmt.Sprintf("lease for group with id %v does not exist", lease.GroupUID))
			}
			return tools.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for group with id %v. See logs for details", lease.GroupUID))
		}

		if leaseToBeTerminated == nil {
			requestContextLogger.Error("leaseToBeTerminated == nil", "err", err)
			return tools.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for group with id %v. See logs for details", lease.GroupUID))
		}

		if leaseToBeTerminated.TokenOnce != ctx.Tok {
			requestContextLogger.Error("Wrong tokenOnce; link already used/expired")
			return tools.ErrNotFound(ctx, "link expired")
		}

		c.cs.Queues().TerminatorQueue().PushTask(tasks.TerminatorTask{Lease: *leaseToBeTerminated})

		resp = tools.HMI{
			"message":    "Termination request received",
			"group_type": leaseToBeTerminated.GroupType.String(),
			"group_uid":  leaseToBeTerminated.GroupUID,
			"lease_id":   leaseToBeTerminated.ID,
		}
		lease = *leaseToBeTerminated
	default:
		tools.ErrNotFound(ctx, "action not found")
	}

	// TODO: add info about the size of the lease, etc.

	return tools.JSONResponse(ctx, 202, resp)
}
