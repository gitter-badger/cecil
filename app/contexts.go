//************************************************************************//
// API "zerocloud": Application Contexts
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/zerocloud/design
// --out=$(GOPATH)/src/github.com/tleyden/zerocloud
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package app

import (
	"github.com/goadesign/goa"
	"golang.org/x/net/context"
	"strconv"
)

// CreateAccountContext provides the account create action context.
type CreateAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Payload *AccountPayload
}

// NewCreateAccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the account controller create action.
func NewCreateAccountContext(ctx context.Context, service *goa.Service) (*CreateAccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := CreateAccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// Created sends a HTTP response with status code 201.
func (ctx *CreateAccountContext) Created() error {
	ctx.ResponseData.WriteHeader(201)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *CreateAccountContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// DeleteAccountContext provides the account delete action context.
type DeleteAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
}

// NewDeleteAccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the account controller delete action.
func NewDeleteAccountContext(ctx context.Context, service *goa.Service) (*DeleteAccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := DeleteAccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["accountID"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("accountID", rawAccountID, "integer"))
		}
	}
	return &rctx, err
}

// NoContent sends a HTTP response with status code 204.
func (ctx *DeleteAccountContext) NoContent() error {
	ctx.ResponseData.WriteHeader(204)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *DeleteAccountContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *DeleteAccountContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// ListAccountContext provides the account list action context.
type ListAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
}

// NewListAccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the account controller list action.
func NewListAccountContext(ctx context.Context, service *goa.Service) (*ListAccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ListAccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ListAccountContext) OK(r AccountCollection) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.account+json; type=collection")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// OKLink sends a HTTP response with status code 200.
func (ctx *ListAccountContext) OKLink(r AccountLinkCollection) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.account+json; type=collection")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// OKTiny sends a HTTP response with status code 200.
func (ctx *ListAccountContext) OKTiny(r AccountTinyCollection) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.account+json; type=collection")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ListAccountContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// ShowAccountContext provides the account show action context.
type ShowAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
}

// NewShowAccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the account controller show action.
func NewShowAccountContext(ctx context.Context, service *goa.Service) (*ShowAccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ShowAccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["accountID"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("accountID", rawAccountID, "integer"))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowAccountContext) OK(r *Account) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.account+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// OKLink sends a HTTP response with status code 200.
func (ctx *ShowAccountContext) OKLink(r *AccountLink) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.account+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// OKTiny sends a HTTP response with status code 200.
func (ctx *ShowAccountContext) OKTiny(r *AccountTiny) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.account+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *ShowAccountContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ShowAccountContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// UpdateAccountContext provides the account update action context.
type UpdateAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
	Payload   *UpdateAccountPayload
}

// NewUpdateAccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the account controller update action.
func NewUpdateAccountContext(ctx context.Context, service *goa.Service) (*UpdateAccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := UpdateAccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["accountID"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("accountID", rawAccountID, "integer"))
		}
	}
	return &rctx, err
}

// updateAccountPayload is the account update action payload.
type updateAccountPayload struct {
	// Name of account
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *updateAccountPayload) Validate() (err error) {
	if payload.Name == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "name"))
	}

	if payload.Name != nil {
		if len(*payload.Name) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.name`, *payload.Name, len(*payload.Name), 3, true))
		}
	}
	return
}

// Publicize creates UpdateAccountPayload from updateAccountPayload
func (payload *updateAccountPayload) Publicize() *UpdateAccountPayload {
	var pub UpdateAccountPayload
	if payload.Name != nil {
		pub.Name = *payload.Name
	}
	return &pub
}

// UpdateAccountPayload is the account update action payload.
type UpdateAccountPayload struct {
	// Name of account
	Name string `form:"name" json:"name" xml:"name"`
}

// Validate runs the validation rules defined in the design.
func (payload *UpdateAccountPayload) Validate() (err error) {
	if payload.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "name"))
	}

	if len(payload.Name) < 3 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.name`, payload.Name, len(payload.Name), 3, true))
	}
	return
}

// NoContent sends a HTTP response with status code 204.
func (ctx *UpdateAccountContext) NoContent() error {
	ctx.ResponseData.WriteHeader(204)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *UpdateAccountContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *UpdateAccountContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// ShowAwsContext provides the aws show action context.
type ShowAwsContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AwsAccountID string
}

// NewShowAwsContext parses the incoming request URL and body, performs validations and creates the
// context used by the aws controller show action.
func NewShowAwsContext(ctx context.Context, service *goa.Service) (*ShowAwsContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ShowAwsContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAwsAccountID := req.Params["awsAccountID"]
	if len(paramAwsAccountID) > 0 {
		rawAwsAccountID := paramAwsAccountID[0]
		rctx.AwsAccountID = rawAwsAccountID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowAwsContext) OK(r *Cloudaccount) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.cloudaccount+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// OKLink sends a HTTP response with status code 200.
func (ctx *ShowAwsContext) OKLink(r *CloudaccountLink) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.cloudaccount+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// OKTiny sends a HTTP response with status code 200.
func (ctx *ShowAwsContext) OKTiny(r *CloudaccountTiny) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.cloudaccount+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *ShowAwsContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ShowAwsContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// CreateCloudaccountContext provides the cloudaccount create action context.
type CreateCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
	Payload   *CloudAccountPayload
}

// NewCreateCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller create action.
func NewCreateCloudaccountContext(ctx context.Context, service *goa.Service) (*CreateCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := CreateCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["accountID"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("accountID", rawAccountID, "integer"))
		}
	}
	return &rctx, err
}

// Created sends a HTTP response with status code 201.
func (ctx *CreateCloudaccountContext) Created() error {
	ctx.ResponseData.WriteHeader(201)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *CreateCloudaccountContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *CreateCloudaccountContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// DeleteCloudaccountContext provides the cloudaccount delete action context.
type DeleteCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudAccountID int
}

// NewDeleteCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller delete action.
func NewDeleteCloudaccountContext(ctx context.Context, service *goa.Service) (*DeleteCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := DeleteCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["accountID"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("accountID", rawAccountID, "integer"))
		}
	}
	paramCloudAccountID := req.Params["cloudAccountID"]
	if len(paramCloudAccountID) > 0 {
		rawCloudAccountID := paramCloudAccountID[0]
		if cloudAccountID, err2 := strconv.Atoi(rawCloudAccountID); err2 == nil {
			rctx.CloudAccountID = cloudAccountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("cloudAccountID", rawCloudAccountID, "integer"))
		}
	}
	return &rctx, err
}

// NoContent sends a HTTP response with status code 204.
func (ctx *DeleteCloudaccountContext) NoContent() error {
	ctx.ResponseData.WriteHeader(204)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *DeleteCloudaccountContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *DeleteCloudaccountContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// ListCloudaccountContext provides the cloudaccount list action context.
type ListCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
}

// NewListCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller list action.
func NewListCloudaccountContext(ctx context.Context, service *goa.Service) (*ListCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ListCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["accountID"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("accountID", rawAccountID, "integer"))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ListCloudaccountContext) OK(r CloudaccountCollection) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.cloudaccount+json; type=collection")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// OKTiny sends a HTTP response with status code 200.
func (ctx *ListCloudaccountContext) OKTiny(r CloudaccountTinyCollection) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.cloudaccount+json; type=collection")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *ListCloudaccountContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ListCloudaccountContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// ShowCloudaccountContext provides the cloudaccount show action context.
type ShowCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudAccountID int
}

// NewShowCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller show action.
func NewShowCloudaccountContext(ctx context.Context, service *goa.Service) (*ShowCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ShowCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["accountID"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("accountID", rawAccountID, "integer"))
		}
	}
	paramCloudAccountID := req.Params["cloudAccountID"]
	if len(paramCloudAccountID) > 0 {
		rawCloudAccountID := paramCloudAccountID[0]
		if cloudAccountID, err2 := strconv.Atoi(rawCloudAccountID); err2 == nil {
			rctx.CloudAccountID = cloudAccountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("cloudAccountID", rawCloudAccountID, "integer"))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowCloudaccountContext) OK(r *Cloudaccount) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.cloudaccount+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// OKLink sends a HTTP response with status code 200.
func (ctx *ShowCloudaccountContext) OKLink(r *CloudaccountLink) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.cloudaccount+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// OKTiny sends a HTTP response with status code 200.
func (ctx *ShowCloudaccountContext) OKTiny(r *CloudaccountTiny) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.cloudaccount+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *ShowCloudaccountContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ShowCloudaccountContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// UpdateCloudaccountContext provides the cloudaccount update action context.
type UpdateCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudAccountID int
	Payload        *CloudAccountPayload
}

// NewUpdateCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller update action.
func NewUpdateCloudaccountContext(ctx context.Context, service *goa.Service) (*UpdateCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := UpdateCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["accountID"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("accountID", rawAccountID, "integer"))
		}
	}
	paramCloudAccountID := req.Params["cloudAccountID"]
	if len(paramCloudAccountID) > 0 {
		rawCloudAccountID := paramCloudAccountID[0]
		if cloudAccountID, err2 := strconv.Atoi(rawCloudAccountID); err2 == nil {
			rctx.CloudAccountID = cloudAccountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("cloudAccountID", rawCloudAccountID, "integer"))
		}
	}
	return &rctx, err
}

// NoContent sends a HTTP response with status code 204.
func (ctx *UpdateCloudaccountContext) NoContent() error {
	ctx.ResponseData.WriteHeader(204)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *UpdateCloudaccountContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *UpdateCloudaccountContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// CreateCloudeventContext provides the cloudevent create action context.
type CreateCloudeventContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Payload *CloudEventPayload
}

// NewCreateCloudeventContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudevent controller create action.
func NewCreateCloudeventContext(ctx context.Context, service *goa.Service) (*CreateCloudeventContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := CreateCloudeventContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// Created sends a HTTP response with status code 201.
func (ctx *CreateCloudeventContext) Created() error {
	ctx.ResponseData.WriteHeader(201)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *CreateCloudeventContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}
