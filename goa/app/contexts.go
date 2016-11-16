//************************************************************************//
// API "Cecil REST API": Application Contexts
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/cecil/design
// --out=$(GOPATH)src/github.com/tleyden/cecil/goa
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package app

import (
	"github.com/goadesign/goa"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"strconv"
	"unicode/utf8"
)

// CreateAccountContext provides the account create action context.
type CreateAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Payload *CreateAccountPayload
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

// createAccountPayload is the account create action payload.
type createAccountPayload struct {
	Email   *string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
	Name    *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	Surname *string `form:"surname,omitempty" json:"surname,omitempty" xml:"surname,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *createAccountPayload) Validate() (err error) {
	if payload.Email == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "email"))
	}
	if payload.Name == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "name"))
	}
	if payload.Surname == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "surname"))
	}

	if payload.Email != nil {
		if err2 := goa.ValidateFormat(goa.FormatEmail, *payload.Email); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`raw.email`, *payload.Email, goa.FormatEmail, err2))
		}
	}
	if payload.Name != nil {
		if utf8.RuneCountInString(*payload.Name) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.name`, *payload.Name, utf8.RuneCountInString(*payload.Name), 1, true))
		}
	}
	if payload.Name != nil {
		if utf8.RuneCountInString(*payload.Name) > 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.name`, *payload.Name, utf8.RuneCountInString(*payload.Name), 30, false))
		}
	}
	if payload.Surname != nil {
		if utf8.RuneCountInString(*payload.Surname) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.surname`, *payload.Surname, utf8.RuneCountInString(*payload.Surname), 1, true))
		}
	}
	if payload.Surname != nil {
		if utf8.RuneCountInString(*payload.Surname) > 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.surname`, *payload.Surname, utf8.RuneCountInString(*payload.Surname), 30, false))
		}
	}
	return
}

// Publicize creates CreateAccountPayload from createAccountPayload
func (payload *createAccountPayload) Publicize() *CreateAccountPayload {
	var pub CreateAccountPayload
	if payload.Email != nil {
		pub.Email = *payload.Email
	}
	if payload.Name != nil {
		pub.Name = *payload.Name
	}
	if payload.Surname != nil {
		pub.Surname = *payload.Surname
	}
	return &pub
}

// CreateAccountPayload is the account create action payload.
type CreateAccountPayload struct {
	Email   string `form:"email" json:"email" xml:"email"`
	Name    string `form:"name" json:"name" xml:"name"`
	Surname string `form:"surname" json:"surname" xml:"surname"`
}

// Validate runs the validation rules defined in the design.
func (payload *CreateAccountPayload) Validate() (err error) {
	if payload.Email == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "email"))
	}
	if payload.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "name"))
	}
	if payload.Surname == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "surname"))
	}

	if err2 := goa.ValidateFormat(goa.FormatEmail, payload.Email); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`raw.email`, payload.Email, goa.FormatEmail, err2))
	}
	if utf8.RuneCountInString(payload.Name) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.name`, payload.Name, utf8.RuneCountInString(payload.Name), 1, true))
	}
	if utf8.RuneCountInString(payload.Name) > 30 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.name`, payload.Name, utf8.RuneCountInString(payload.Name), 30, false))
	}
	if utf8.RuneCountInString(payload.Surname) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.surname`, payload.Surname, utf8.RuneCountInString(payload.Surname), 1, true))
	}
	if utf8.RuneCountInString(payload.Surname) > 30 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.surname`, payload.Surname, utf8.RuneCountInString(payload.Surname), 30, false))
	}
	return
}

// OK sends a HTTP response with status code 200.
func (ctx *CreateAccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
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
	paramAccountID := req.Params["account_id"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("account_id", rawAccountID, "integer"))
		}
		if rctx.AccountID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`account_id`, rctx.AccountID, 1, true))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowAccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// VerifyAccountContext provides the account verify action context.
type VerifyAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
	Payload   *VerifyAccountPayload
}

// NewVerifyAccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the account controller verify action.
func NewVerifyAccountContext(ctx context.Context, service *goa.Service) (*VerifyAccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := VerifyAccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["account_id"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("account_id", rawAccountID, "integer"))
		}
		if rctx.AccountID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`account_id`, rctx.AccountID, 1, true))
		}
	}
	return &rctx, err
}

// verifyAccountPayload is the account verify action payload.
type verifyAccountPayload struct {
	VerificationToken *string `form:"verification_token,omitempty" json:"verification_token,omitempty" xml:"verification_token,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *verifyAccountPayload) Validate() (err error) {
	if payload.VerificationToken == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "verification_token"))
	}

	if payload.VerificationToken != nil {
		if utf8.RuneCountInString(*payload.VerificationToken) < 108 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.verification_token`, *payload.VerificationToken, utf8.RuneCountInString(*payload.VerificationToken), 108, true))
		}
	}
	return
}

// Publicize creates VerifyAccountPayload from verifyAccountPayload
func (payload *verifyAccountPayload) Publicize() *VerifyAccountPayload {
	var pub VerifyAccountPayload
	if payload.VerificationToken != nil {
		pub.VerificationToken = *payload.VerificationToken
	}
	return &pub
}

// VerifyAccountPayload is the account verify action payload.
type VerifyAccountPayload struct {
	VerificationToken string `form:"verification_token" json:"verification_token" xml:"verification_token"`
}

// Validate runs the validation rules defined in the design.
func (payload *VerifyAccountPayload) Validate() (err error) {
	if payload.VerificationToken == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "verification_token"))
	}

	if utf8.RuneCountInString(payload.VerificationToken) < 108 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.verification_token`, payload.VerificationToken, utf8.RuneCountInString(payload.VerificationToken), 108, true))
	}
	return
}

// OK sends a HTTP response with status code 200.
func (ctx *VerifyAccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// AddCloudaccountContext provides the cloudaccount add action context.
type AddCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
	Payload   *AddCloudaccountPayload
}

// NewAddCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller add action.
func NewAddCloudaccountContext(ctx context.Context, service *goa.Service) (*AddCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := AddCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["account_id"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("account_id", rawAccountID, "integer"))
		}
		if rctx.AccountID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`account_id`, rctx.AccountID, 1, true))
		}
	}
	return &rctx, err
}

// addCloudaccountPayload is the cloudaccount add action payload.
type addCloudaccountPayload struct {
	AwsID *string `form:"aws_id,omitempty" json:"aws_id,omitempty" xml:"aws_id,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *addCloudaccountPayload) Validate() (err error) {
	if payload.AwsID == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "aws_id"))
	}

	if payload.AwsID != nil {
		if utf8.RuneCountInString(*payload.AwsID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.aws_id`, *payload.AwsID, utf8.RuneCountInString(*payload.AwsID), 1, true))
		}
	}
	return
}

// Publicize creates AddCloudaccountPayload from addCloudaccountPayload
func (payload *addCloudaccountPayload) Publicize() *AddCloudaccountPayload {
	var pub AddCloudaccountPayload
	if payload.AwsID != nil {
		pub.AwsID = *payload.AwsID
	}
	return &pub
}

// AddCloudaccountPayload is the cloudaccount add action payload.
type AddCloudaccountPayload struct {
	AwsID string `form:"aws_id" json:"aws_id" xml:"aws_id"`
}

// Validate runs the validation rules defined in the design.
func (payload *AddCloudaccountPayload) Validate() (err error) {
	if payload.AwsID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "aws_id"))
	}

	if utf8.RuneCountInString(payload.AwsID) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.aws_id`, payload.AwsID, utf8.RuneCountInString(payload.AwsID), 1, true))
	}
	return
}

// OK sends a HTTP response with status code 200.
func (ctx *AddCloudaccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// AddEmailToWhitelistCloudaccountContext provides the cloudaccount addEmailToWhitelist action context.
type AddEmailToWhitelistCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
	Payload        *AddEmailToWhitelistCloudaccountPayload
}

// NewAddEmailToWhitelistCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller addEmailToWhitelist action.
func NewAddEmailToWhitelistCloudaccountContext(ctx context.Context, service *goa.Service) (*AddEmailToWhitelistCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := AddEmailToWhitelistCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["account_id"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("account_id", rawAccountID, "integer"))
		}
		if rctx.AccountID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`account_id`, rctx.AccountID, 1, true))
		}
	}
	paramCloudaccountID := req.Params["cloudaccount_id"]
	if len(paramCloudaccountID) > 0 {
		rawCloudaccountID := paramCloudaccountID[0]
		if cloudaccountID, err2 := strconv.Atoi(rawCloudaccountID); err2 == nil {
			rctx.CloudaccountID = cloudaccountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("cloudaccount_id", rawCloudaccountID, "integer"))
		}
		if rctx.CloudaccountID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`cloudaccount_id`, rctx.CloudaccountID, 1, true))
		}
	}
	return &rctx, err
}

// addEmailToWhitelistCloudaccountPayload is the cloudaccount addEmailToWhitelist action payload.
type addEmailToWhitelistCloudaccountPayload struct {
	Email *string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *addEmailToWhitelistCloudaccountPayload) Validate() (err error) {
	if payload.Email == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "email"))
	}

	if payload.Email != nil {
		if err2 := goa.ValidateFormat(goa.FormatEmail, *payload.Email); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`raw.email`, *payload.Email, goa.FormatEmail, err2))
		}
	}
	return
}

// Publicize creates AddEmailToWhitelistCloudaccountPayload from addEmailToWhitelistCloudaccountPayload
func (payload *addEmailToWhitelistCloudaccountPayload) Publicize() *AddEmailToWhitelistCloudaccountPayload {
	var pub AddEmailToWhitelistCloudaccountPayload
	if payload.Email != nil {
		pub.Email = *payload.Email
	}
	return &pub
}

// AddEmailToWhitelistCloudaccountPayload is the cloudaccount addEmailToWhitelist action payload.
type AddEmailToWhitelistCloudaccountPayload struct {
	Email string `form:"email" json:"email" xml:"email"`
}

// Validate runs the validation rules defined in the design.
func (payload *AddEmailToWhitelistCloudaccountPayload) Validate() (err error) {
	if payload.Email == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "email"))
	}

	if err2 := goa.ValidateFormat(goa.FormatEmail, payload.Email); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`raw.email`, payload.Email, goa.FormatEmail, err2))
	}
	return
}

// OK sends a HTTP response with status code 200.
func (ctx *AddEmailToWhitelistCloudaccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// DownloadInitialSetupTemplateCloudaccountContext provides the cloudaccount downloadInitialSetupTemplate action context.
type DownloadInitialSetupTemplateCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
}

// NewDownloadInitialSetupTemplateCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller downloadInitialSetupTemplate action.
func NewDownloadInitialSetupTemplateCloudaccountContext(ctx context.Context, service *goa.Service) (*DownloadInitialSetupTemplateCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := DownloadInitialSetupTemplateCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["account_id"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("account_id", rawAccountID, "integer"))
		}
		if rctx.AccountID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`account_id`, rctx.AccountID, 1, true))
		}
	}
	paramCloudaccountID := req.Params["cloudaccount_id"]
	if len(paramCloudaccountID) > 0 {
		rawCloudaccountID := paramCloudaccountID[0]
		if cloudaccountID, err2 := strconv.Atoi(rawCloudaccountID); err2 == nil {
			rctx.CloudaccountID = cloudaccountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("cloudaccount_id", rawCloudaccountID, "integer"))
		}
		if rctx.CloudaccountID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`cloudaccount_id`, rctx.CloudaccountID, 1, true))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *DownloadInitialSetupTemplateCloudaccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// DownloadRegionSetupTemplateCloudaccountContext provides the cloudaccount downloadRegionSetupTemplate action context.
type DownloadRegionSetupTemplateCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
}

// NewDownloadRegionSetupTemplateCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller downloadRegionSetupTemplate action.
func NewDownloadRegionSetupTemplateCloudaccountContext(ctx context.Context, service *goa.Service) (*DownloadRegionSetupTemplateCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := DownloadRegionSetupTemplateCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAccountID := req.Params["account_id"]
	if len(paramAccountID) > 0 {
		rawAccountID := paramAccountID[0]
		if accountID, err2 := strconv.Atoi(rawAccountID); err2 == nil {
			rctx.AccountID = accountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("account_id", rawAccountID, "integer"))
		}
		if rctx.AccountID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`account_id`, rctx.AccountID, 1, true))
		}
	}
	paramCloudaccountID := req.Params["cloudaccount_id"]
	if len(paramCloudaccountID) > 0 {
		rawCloudaccountID := paramCloudaccountID[0]
		if cloudaccountID, err2 := strconv.Atoi(rawCloudaccountID); err2 == nil {
			rctx.CloudaccountID = cloudaccountID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("cloudaccount_id", rawCloudaccountID, "integer"))
		}
		if rctx.CloudaccountID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`cloudaccount_id`, rctx.CloudaccountID, 1, true))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *DownloadRegionSetupTemplateCloudaccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// ActionsEmailActionContext provides the email_action actions action context.
type ActionsEmailActionContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Action     string
	InstanceID string
	LeaseUUID  uuid.UUID
	Sig        string
	Tok        string
}

// NewActionsEmailActionContext parses the incoming request URL and body, performs validations and creates the
// context used by the email_action controller actions action.
func NewActionsEmailActionContext(ctx context.Context, service *goa.Service) (*ActionsEmailActionContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ActionsEmailActionContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAction := req.Params["action"]
	if len(paramAction) > 0 {
		rawAction := paramAction[0]
		rctx.Action = rawAction
		if !(rctx.Action == "approve" || rctx.Action == "terminate" || rctx.Action == "extend") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError(`action`, rctx.Action, []interface{}{"approve", "terminate", "extend"}))
		}
	}
	paramInstanceID := req.Params["instance_id"]
	if len(paramInstanceID) > 0 {
		rawInstanceID := paramInstanceID[0]
		rctx.InstanceID = rawInstanceID
		if utf8.RuneCountInString(rctx.InstanceID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`instance_id`, rctx.InstanceID, utf8.RuneCountInString(rctx.InstanceID), 1, true))
		}
	}
	paramLeaseUUID := req.Params["lease_uuid"]
	if len(paramLeaseUUID) > 0 {
		rawLeaseUUID := paramLeaseUUID[0]
		if leaseUUID, err2 := uuid.FromString(rawLeaseUUID); err2 == nil {
			rctx.LeaseUUID = leaseUUID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("lease_uuid", rawLeaseUUID, "uuid"))
		}
	}
	paramSig := req.Params["sig"]
	if len(paramSig) == 0 {
		err = goa.MergeErrors(err, goa.MissingParamError("sig"))
	} else {
		rawSig := paramSig[0]
		rctx.Sig = rawSig
		if utf8.RuneCountInString(rctx.Sig) < 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`sig`, rctx.Sig, utf8.RuneCountInString(rctx.Sig), 30, true))
		}
	}
	paramTok := req.Params["tok"]
	if len(paramTok) == 0 {
		err = goa.MergeErrors(err, goa.MissingParamError("tok"))
	} else {
		rawTok := paramTok[0]
		rctx.Tok = rawTok
		if utf8.RuneCountInString(rctx.Tok) < 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`tok`, rctx.Tok, utf8.RuneCountInString(rctx.Tok), 30, true))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ActionsEmailActionContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// ShowRootContext provides the root show action context.
type ShowRootContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
}

// NewShowRootContext parses the incoming request URL and body, performs validations and creates the
// context used by the root controller show action.
func NewShowRootContext(ctx context.Context, service *goa.Service) (*ShowRootContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ShowRootContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowRootContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}
