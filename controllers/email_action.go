package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/jinzhu/gorm"
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

// Actions handles the endpoint used to receive email_actions (i.e. links sent in emails that make perform specfic actions on leases).
func (c *EmailActionController) Actions(ctx *app.ActionsEmailActionContext) error {
	requestContextLogger := core.Logger.New(
		"url", ctx.Request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)

	err := c.cs.EmailActionVerifySignatureParams(ctx.LeaseUUID.String(), ctx.InstanceID, ctx.Action, ctx.Tok, ctx.Sig)
	if err != nil {
		requestContextLogger.Error("Signature verification error", "err", err)
		return core.ErrInvalidRequest(ctx, "corrupted action link")
	}

	switch ctx.Action {
	case "approve":
		core.Logger.Info("Approval of lease initiated", "instance_id", ctx.InstanceID)

		leaseToBeApproved, err := c.cs.LeaseByIDAndUUID(ctx.InstanceID, ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("lease for instance with id %v does not exist", ctx.InstanceID))
			} else {
				return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.InstanceID))
			}
		}

		if leaseToBeApproved == nil {
			requestContextLogger.Error("leaseToBeApproved == nil", "err", err)
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.InstanceID))
		}

		if leaseToBeApproved.TokenOnce != ctx.Tok {
			requestContextLogger.Error("Wrong tokenOnce; link already used/expired")
			return core.ErrNotFound(ctx, "link expired")
		}

		c.cs.ExtenderQueue.TaskQueue <- core.ExtenderTask{
			Lease:     *leaseToBeApproved,
			Approving: true,
		}

		return core.JSONResponse(ctx, 202, gin.H{
			"instance_id": leaseToBeApproved.InstanceID,
			"lease_id":    leaseToBeApproved.ID,
			"message":     "Approval request received",
		})

	case "extend":
		requestContextLogger.Info("Extension of lease initiated", "instance_id", ctx.InstanceID)

		leaseToBeExtended, err := c.cs.LeaseByIDAndUUID(ctx.InstanceID, ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("lease for instance with id %v does not exist", ctx.InstanceID))
			} else {
				return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.InstanceID))
			}
		}

		if leaseToBeExtended == nil {
			requestContextLogger.Error("leaseToBeExtended == nil", "err", err)
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.InstanceID))
		}

		if leaseToBeExtended.TokenOnce != ctx.Tok {
			requestContextLogger.Error("Wrong tokenOnce; link already used/expired")
			return core.ErrNotFound(ctx, "link expired")
		}

		c.cs.ExtenderQueue.TaskQueue <- core.ExtenderTask{
			Lease:     *leaseToBeExtended,
			Approving: false,
		}

		return core.JSONResponse(ctx, 202, gin.H{
			"instance_id": leaseToBeExtended.InstanceID,
			"lease_id":    leaseToBeExtended.ID,
			"message":     "Extension initiated",
		})

	case "terminate":
		requestContextLogger.Info("Termination of lease initiated", "instance_id", ctx.InstanceID)

		leaseToBeTerminated, err := c.cs.LeaseByIDAndUUID(ctx.InstanceID, ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("lease for instance with id %v does not exist", ctx.InstanceID))
			} else {
				return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.InstanceID))
			}
		}

		if leaseToBeTerminated == nil {
			requestContextLogger.Error("leaseToBeTerminated == nil", "err", err)
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.InstanceID))
		}

		if leaseToBeTerminated.TokenOnce != ctx.Tok {
			requestContextLogger.Error("Wrong tokenOnce; link already used/expired")
			return core.ErrNotFound(ctx, "link expired")
		}

		c.cs.TerminatorQueue.TaskQueue <- core.TerminatorTask{Lease: *leaseToBeTerminated}

		return core.JSONResponse(ctx, 202, gin.H{
			"instance_id": leaseToBeTerminated.InstanceID,
			"lease_id":    leaseToBeTerminated.ID,
			"message":     "Termination initiated",
		})
	}
	// TODO: return "non-allowed action"
	return nil
}
