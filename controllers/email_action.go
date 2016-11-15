package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goadesign/goa"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
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

// Actions runs the actions action.
func (c *EmailActionController) Actions(ctx *app.ActionsEmailActionContext) error {

	err := c.cs.EmailActionVerifySignatureParams(ctx.LeaseUUID.String(), ctx.InstanceID, ctx.Action, ctx.Tok, ctx.Sig)
	if err != nil {
		core.Logger.Warn("Signature verification error", "error", err)
		return core.ErrInvalidRequest(ctx, "corrupted action link")
	}

	switch ctx.Action {
	case "approve":
		core.Logger.Info("Approval of lease initiated", "instance_id", ctx.InstanceID)

		var leaseToBeApproved core.Lease
		var leaseCount int64
		c.cs.DB.Table("leases").Where(&core.Lease{
			InstanceID: ctx.InstanceID,
			UUID:       ctx.LeaseUUID.String(),
			Terminated: false,
		}).Count(&leaseCount).First(&leaseToBeApproved)

		if leaseCount == 0 {
			core.Logger.Warn("No lease found for approval", "count", leaseCount)
			return core.ErrNotFound(ctx, "not found")
		}
		if leaseCount > 1 {
			core.Logger.Warn("Multiple leases found for approval", "count", leaseCount)
			return core.ErrInternal(ctx, "internal exception error")
		}

		if leaseToBeApproved.TokenOnce != ctx.Tok {
			core.Logger.Warn("leaseToBeApproved.TokenOnce != c.Query(\"t\")")
			return core.ErrNotFound(ctx, "link expired")
		}

		c.cs.ExtenderQueue.TaskQueue <- core.ExtenderTask{
			Lease:     leaseToBeApproved,
			ExtendBy:  time.Duration(c.cs.Config.Lease.Duration),
			Approving: true,
		}

		return ctx.Service.Send(ctx, 202, gin.H{
			"instanceId": ctx.InstanceID,
			"message":    "Approval request received",
		})

	case "extend":
		core.Logger.Info("Extension of lease initiated", "instance_id", ctx.InstanceID)

		var leaseToBeExtended core.Lease
		var leaseCount int64
		c.cs.DB.Table("leases").Where(&core.Lease{
			InstanceID: ctx.InstanceID,
			UUID:       ctx.LeaseUUID.String(),
			Terminated: false,
		}).Count(&leaseCount).First(&leaseToBeExtended)

		if leaseCount == 0 {
			core.Logger.Warn("No lease found for extension", "count", leaseCount)
			return core.ErrNotFound(ctx, "not found")
		}
		if leaseCount > 1 {
			core.Logger.Warn("Multiple leases found for extension", "count", leaseCount)
			return core.ErrInternal(ctx, "internal exception error")
		}

		if leaseToBeExtended.TokenOnce != ctx.Tok {
			core.Logger.Warn("leaseToBeExtended.TokenOnce != c.Query(\"t\")")
			return core.ErrNotFound(ctx, "link expired")
		}

		c.cs.ExtenderQueue.TaskQueue <- core.ExtenderTask{
			Lease:     leaseToBeExtended,
			ExtendBy:  time.Duration(c.cs.Config.Lease.Duration),
			Approving: false,
		}

		return ctx.Service.Send(ctx, 202, gin.H{
			"instanceId": ctx.InstanceID,
			"message":    "Extension initiated",
		})

	case "terminate":
		core.Logger.Info("Termination of lease initiated", "instance_id", ctx.InstanceID)

		var leaseCount int64
		var leaseToBeTerminated core.Lease
		c.cs.DB.Table("leases").Where(&core.Lease{
			InstanceID: ctx.InstanceID,
			UUID:       ctx.LeaseUUID.String(),
			Terminated: false,
		}).Count(&leaseCount).First(&leaseToBeTerminated)

		if leaseCount == 0 {
			core.Logger.Warn("No lease found for approval", "count", leaseCount)
			return core.ErrNotFound(ctx, "not found")
		}
		if leaseCount > 1 {
			core.Logger.Warn("Multiple leases found for approval", "count", leaseCount)
			return core.ErrInternal(ctx, "internal exception error")
		}

		if leaseToBeTerminated.TokenOnce != ctx.Tok {
			return core.ErrNotFound(ctx, "link expired")
		}

		c.cs.TerminatorQueue.TaskQueue <- core.TerminatorTask{Lease: leaseToBeTerminated}

		return ctx.Service.Send(ctx, 202, gin.H{
			"instanceId": ctx.InstanceID,
			"message":    "Termination initiated",
		})
	}
	// TODO: return "non-allowed action"
	return nil
}
