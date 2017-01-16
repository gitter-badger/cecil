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
			}
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.InstanceID))
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

		resp := gin.H{
			"message":     "Approval request received",
			"instance_id": leaseToBeApproved.InstanceID,
			"lease_id":    leaseToBeApproved.ID,
		}

		leaseIsStack := leaseToBeApproved.StackName != ""
		if leaseIsStack {
			resp["logical_id"] = leaseToBeApproved.LogicalID
			resp["stack_id"] = leaseToBeApproved.StackID
			resp["stack_name"] = leaseToBeApproved.StackName
		}

		return core.JSONResponse(ctx, 202, resp)

	case "extend":
		requestContextLogger.Info("Extension of lease initiated", "instance_id", ctx.InstanceID)

		leaseToBeExtended, err := c.cs.LeaseByIDAndUUID(ctx.InstanceID, ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("lease for instance with id %v does not exist", ctx.InstanceID))
			}
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.InstanceID))
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

		resp := gin.H{
			"message":     "Extension initiated",
			"instance_id": leaseToBeExtended.InstanceID,
			"lease_id":    leaseToBeExtended.ID,
		}

		leaseIsStack := leaseToBeExtended.StackName != ""
		if leaseIsStack {
			resp["logical_id"] = leaseToBeExtended.LogicalID
			resp["stack_id"] = leaseToBeExtended.StackID
			resp["stack_name"] = leaseToBeExtended.StackName
		}

		return core.JSONResponse(ctx, 202, resp)

	case "terminate":
		requestContextLogger.Info("Termination of lease/stack initiated", "instance_id", ctx.InstanceID)

		leaseToBeTerminated, err := c.cs.LeaseByIDAndUUID(ctx.InstanceID, ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("lease for instance with id %v does not exist", ctx.InstanceID))
			}
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.InstanceID))
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

		resp := gin.H{
			"message":     "Termination initiated",
			"instance_id": leaseToBeTerminated.InstanceID,
			"lease_id":    leaseToBeTerminated.ID,
		}

		leaseIsStack := leaseToBeTerminated.StackName != ""
		if leaseIsStack {
			resp["logical_id"] = leaseToBeTerminated.LogicalID
			resp["stack_id"] = leaseToBeTerminated.StackID
			resp["stack_name"] = leaseToBeTerminated.StackName
		}

		return core.JSONResponse(ctx, 202, resp)
	}
	// TODO: return "non-allowed action"
	return nil
}
