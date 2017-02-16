package controllers

import (
	"fmt"
	"strconv"

	"github.com/goadesign/goa"
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
	requestContextLogger := core.NewContextLogger(ctx)

	err := c.cs.EmailActionVerifySignatureParams(ctx.LeaseUUID.String(), strconv.Itoa(ctx.ResourceID), ctx.Action, ctx.Tok, ctx.Sig)
	if err != nil {
		requestContextLogger.Error("Signature verification error", "err", err)
		return core.ErrInvalidRequest(ctx, "corrupted action link")
	}

	var lease core.Lease
	var resp = make(core.HMI)

	switch ctx.Action {
	case "approve":
		core.Logger.Info("Approval of lease initiated", "instance_id", ctx.ResourceID)

		leaseToBeApproved, err := c.cs.LeaseByUUID(ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("lease for instance with id %v does not exist", ctx.ResourceID))
			}
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.ResourceID))
		}

		if leaseToBeApproved == nil {
			requestContextLogger.Error("leaseToBeApproved == nil", "err", err)
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.ResourceID))
		}

		if leaseToBeApproved.TokenOnce != ctx.Tok {
			requestContextLogger.Error("Wrong tokenOnce; link already used/expired")
			return core.ErrNotFound(ctx, "link expired")
		}

		if leaseToBeApproved.IsExpired() {
			return core.ErrInvalidRequest(ctx, "lease already expired")
		}

		c.cs.ExtenderQueue.TaskQueue <- core.ExtenderTask{
			Lease:     *leaseToBeApproved,
			Approving: true,
		}

		resp = core.HMI{
			"message":      "Approval request received",
			"resourceID":   leaseToBeApproved.ResourceID,
			"resourceType": leaseToBeApproved.ResourceType,
			"lease_id":     leaseToBeApproved.ID,
		}
		lease = *leaseToBeApproved

	case "extend":
		requestContextLogger.Info("Extension of lease initiated", "instance_id", ctx.ResourceID)

		leaseToBeExtended, err := c.cs.LeaseByUUID(ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("lease for instance with id %v does not exist", ctx.ResourceID))
			}
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.ResourceID))
		}

		if leaseToBeExtended == nil {
			requestContextLogger.Error("leaseToBeExtended == nil", "err", err)
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.ResourceID))
		}

		if leaseToBeExtended.TokenOnce != ctx.Tok {
			requestContextLogger.Error("Wrong tokenOnce; link already used/expired")
			return core.ErrNotFound(ctx, "link expired")
		}

		if leaseToBeExtended.IsExpired() {
			return core.ErrInvalidRequest(ctx, "lease already expired")
		}

		c.cs.ExtenderQueue.TaskQueue <- core.ExtenderTask{
			Lease:     *leaseToBeExtended,
			Approving: false,
		}

		resp = core.HMI{
			"message":      "Extension request received",
			"resourceID":   leaseToBeExtended.ResourceID,
			"resourceType": leaseToBeExtended.ResourceType,
			"lease_id":     leaseToBeExtended.ID,
		}
		lease = *leaseToBeExtended
	case "terminate":
		requestContextLogger.Info("Termination of lease/stack initiated", "instance_id", ctx.ResourceID)

		leaseToBeTerminated, err := c.cs.LeaseByUUID(ctx.LeaseUUID)
		if err != nil {
			requestContextLogger.Error("Error fetching lease", "err", err)
			if err == gorm.ErrRecordNotFound {
				return core.ErrInvalidRequest(ctx, fmt.Sprintf("lease for instance with id %v does not exist", ctx.ResourceID))
			}
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.ResourceID))
		}

		if leaseToBeTerminated == nil {
			requestContextLogger.Error("leaseToBeTerminated == nil", "err", err)
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving lease for instance with id %v. See logs for details", ctx.ResourceID))
		}

		if leaseToBeTerminated.TokenOnce != ctx.Tok {
			requestContextLogger.Error("Wrong tokenOnce; link already used/expired")
			return core.ErrNotFound(ctx, "link expired")
		}

		c.cs.TerminatorQueue.TaskQueue <- core.TerminatorTask{Lease: *leaseToBeTerminated}

		resp = core.HMI{
			"message":      "Termination request received",
			"resourceID":   leaseToBeTerminated.ResourceID,
			"resourceType": leaseToBeTerminated.ResourceType,
			"lease_id":     leaseToBeTerminated.ID,
		}
		lease = *leaseToBeTerminated
	default:
		core.ErrNotFound(ctx, "action not found")
	}

	var instance core.InstanceResource
	if lease.IsInstance() {
		raw, err := c.cs.ResourceOf(&lease)
		if err != nil {
			return err
		}
		instance = raw.(core.InstanceResource)
	}

	var stack core.StackResource
	if lease.IsStack() {
		raw, err := c.cs.ResourceOf(&lease)
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

	// TODO: return "non-allowed action"
}
