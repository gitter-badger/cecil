//************************************************************************//
// API "Cecil": Application Contexts
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/cecil/design
// --out=$(GOPATH)/src/github.com/tleyden/cecil/goa
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
	"time"
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

// MailerConfigAccountContext provides the account mailerConfig action context.
type MailerConfigAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
	Payload   *MailerConfigAccountPayload
}

// NewMailerConfigAccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the account controller mailerConfig action.
func NewMailerConfigAccountContext(ctx context.Context, service *goa.Service) (*MailerConfigAccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := MailerConfigAccountContext{Context: ctx, ResponseData: resp, RequestData: req}
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

// mailerConfigAccountPayload is the account mailerConfig action payload.
type mailerConfigAccountPayload struct {
	APIKey       *string `form:"api_key,omitempty" json:"api_key,omitempty" xml:"api_key,omitempty"`
	Domain       *string `form:"domain,omitempty" json:"domain,omitempty" xml:"domain,omitempty"`
	FromName     *string `form:"from_name,omitempty" json:"from_name,omitempty" xml:"from_name,omitempty"`
	PublicAPIKey *string `form:"public_api_key,omitempty" json:"public_api_key,omitempty" xml:"public_api_key,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *mailerConfigAccountPayload) Validate() (err error) {
	if payload.Domain == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "domain"))
	}
	if payload.APIKey == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "api_key"))
	}
	if payload.PublicAPIKey == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "public_api_key"))
	}
	if payload.FromName == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "from_name"))
	}
	if payload.APIKey != nil {
		if utf8.RuneCountInString(*payload.APIKey) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.api_key`, *payload.APIKey, utf8.RuneCountInString(*payload.APIKey), 1, true))
		}
	}
	if payload.Domain != nil {
		if utf8.RuneCountInString(*payload.Domain) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.domain`, *payload.Domain, utf8.RuneCountInString(*payload.Domain), 1, true))
		}
	}
	if payload.FromName != nil {
		if utf8.RuneCountInString(*payload.FromName) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.from_name`, *payload.FromName, utf8.RuneCountInString(*payload.FromName), 1, true))
		}
	}
	if payload.PublicAPIKey != nil {
		if utf8.RuneCountInString(*payload.PublicAPIKey) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.public_api_key`, *payload.PublicAPIKey, utf8.RuneCountInString(*payload.PublicAPIKey), 1, true))
		}
	}
	return
}

// Publicize creates MailerConfigAccountPayload from mailerConfigAccountPayload
func (payload *mailerConfigAccountPayload) Publicize() *MailerConfigAccountPayload {
	var pub MailerConfigAccountPayload
	if payload.APIKey != nil {
		pub.APIKey = *payload.APIKey
	}
	if payload.Domain != nil {
		pub.Domain = *payload.Domain
	}
	if payload.FromName != nil {
		pub.FromName = *payload.FromName
	}
	if payload.PublicAPIKey != nil {
		pub.PublicAPIKey = *payload.PublicAPIKey
	}
	return &pub
}

// MailerConfigAccountPayload is the account mailerConfig action payload.
type MailerConfigAccountPayload struct {
	APIKey       string `form:"api_key" json:"api_key" xml:"api_key"`
	Domain       string `form:"domain" json:"domain" xml:"domain"`
	FromName     string `form:"from_name" json:"from_name" xml:"from_name"`
	PublicAPIKey string `form:"public_api_key" json:"public_api_key" xml:"public_api_key"`
}

// Validate runs the validation rules defined in the design.
func (payload *MailerConfigAccountPayload) Validate() (err error) {
	if payload.Domain == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "domain"))
	}
	if payload.APIKey == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "api_key"))
	}
	if payload.PublicAPIKey == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "public_api_key"))
	}
	if payload.FromName == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "from_name"))
	}
	if utf8.RuneCountInString(payload.APIKey) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.api_key`, payload.APIKey, utf8.RuneCountInString(payload.APIKey), 1, true))
	}
	if utf8.RuneCountInString(payload.Domain) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.domain`, payload.Domain, utf8.RuneCountInString(payload.Domain), 1, true))
	}
	if utf8.RuneCountInString(payload.FromName) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.from_name`, payload.FromName, utf8.RuneCountInString(payload.FromName), 1, true))
	}
	if utf8.RuneCountInString(payload.PublicAPIKey) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.public_api_key`, payload.PublicAPIKey, utf8.RuneCountInString(payload.PublicAPIKey), 1, true))
	}
	return
}

// OK sends a HTTP response with status code 200.
func (ctx *MailerConfigAccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// RemoveSlackAccountContext provides the account removeSlack action context.
type RemoveSlackAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
}

// NewRemoveSlackAccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the account controller removeSlack action.
func NewRemoveSlackAccountContext(ctx context.Context, service *goa.Service) (*RemoveSlackAccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := RemoveSlackAccountContext{Context: ctx, ResponseData: resp, RequestData: req}
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
func (ctx *RemoveSlackAccountContext) OK(resp []byte) error {
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

// SlackConfigAccountContext provides the account slackConfig action context.
type SlackConfigAccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID int
	Payload   *SlackConfigAccountPayload
}

// NewSlackConfigAccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the account controller slackConfig action.
func NewSlackConfigAccountContext(ctx context.Context, service *goa.Service) (*SlackConfigAccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := SlackConfigAccountContext{Context: ctx, ResponseData: resp, RequestData: req}
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

// slackConfigAccountPayload is the account slackConfig action payload.
type slackConfigAccountPayload struct {
	ChannelID *string `form:"channel_id,omitempty" json:"channel_id,omitempty" xml:"channel_id,omitempty"`
	Token     *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *slackConfigAccountPayload) Validate() (err error) {
	if payload.Token == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "token"))
	}
	if payload.ChannelID == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "channel_id"))
	}
	if payload.ChannelID != nil {
		if utf8.RuneCountInString(*payload.ChannelID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.channel_id`, *payload.ChannelID, utf8.RuneCountInString(*payload.ChannelID), 1, true))
		}
	}
	if payload.Token != nil {
		if utf8.RuneCountInString(*payload.Token) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.token`, *payload.Token, utf8.RuneCountInString(*payload.Token), 1, true))
		}
	}
	return
}

// Publicize creates SlackConfigAccountPayload from slackConfigAccountPayload
func (payload *slackConfigAccountPayload) Publicize() *SlackConfigAccountPayload {
	var pub SlackConfigAccountPayload
	if payload.ChannelID != nil {
		pub.ChannelID = *payload.ChannelID
	}
	if payload.Token != nil {
		pub.Token = *payload.Token
	}
	return &pub
}

// SlackConfigAccountPayload is the account slackConfig action payload.
type SlackConfigAccountPayload struct {
	ChannelID string `form:"channel_id" json:"channel_id" xml:"channel_id"`
	Token     string `form:"token" json:"token" xml:"token"`
}

// Validate runs the validation rules defined in the design.
func (payload *SlackConfigAccountPayload) Validate() (err error) {
	if payload.Token == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "token"))
	}
	if payload.ChannelID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "channel_id"))
	}
	if utf8.RuneCountInString(payload.ChannelID) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.channel_id`, payload.ChannelID, utf8.RuneCountInString(payload.ChannelID), 1, true))
	}
	if utf8.RuneCountInString(payload.Token) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.token`, payload.Token, utf8.RuneCountInString(payload.Token), 1, true))
	}
	return
}

// OK sends a HTTP response with status code 200.
func (ctx *SlackConfigAccountContext) OK(resp []byte) error {
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
	AwsID                *string `form:"aws_id,omitempty" json:"aws_id,omitempty" xml:"aws_id,omitempty"`
	DefaultLeaseDuration *string `form:"default_lease_duration,omitempty" json:"default_lease_duration,omitempty" xml:"default_lease_duration,omitempty"`
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
	if payload.DefaultLeaseDuration != nil {
		pub.DefaultLeaseDuration = payload.DefaultLeaseDuration
	}
	return &pub
}

// AddCloudaccountPayload is the cloudaccount add action payload.
type AddCloudaccountPayload struct {
	AwsID                string  `form:"aws_id" json:"aws_id" xml:"aws_id"`
	DefaultLeaseDuration *string `form:"default_lease_duration,omitempty" json:"default_lease_duration,omitempty" xml:"default_lease_duration,omitempty"`
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

// ListRegionsCloudaccountContext provides the cloudaccount listRegions action context.
type ListRegionsCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
}

// NewListRegionsCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller listRegions action.
func NewListRegionsCloudaccountContext(ctx context.Context, service *goa.Service) (*ListRegionsCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ListRegionsCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
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
func (ctx *ListRegionsCloudaccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// SubscribeSNSToSQSCloudaccountContext provides the cloudaccount subscribeSNSToSQS action context.
type SubscribeSNSToSQSCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
	Payload        *SubscribeSNSToSQSCloudaccountPayload
}

// NewSubscribeSNSToSQSCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller subscribeSNSToSQS action.
func NewSubscribeSNSToSQSCloudaccountContext(ctx context.Context, service *goa.Service) (*SubscribeSNSToSQSCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := SubscribeSNSToSQSCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
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

// subscribeSNSToSQSCloudaccountPayload is the cloudaccount subscribeSNSToSQS action payload.
type subscribeSNSToSQSCloudaccountPayload struct {
	Regions []string `form:"regions,omitempty" json:"regions,omitempty" xml:"regions,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *subscribeSNSToSQSCloudaccountPayload) Validate() (err error) {
	if payload.Regions == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "regions"))
	}
	return
}

// Publicize creates SubscribeSNSToSQSCloudaccountPayload from subscribeSNSToSQSCloudaccountPayload
func (payload *subscribeSNSToSQSCloudaccountPayload) Publicize() *SubscribeSNSToSQSCloudaccountPayload {
	var pub SubscribeSNSToSQSCloudaccountPayload
	if payload.Regions != nil {
		pub.Regions = payload.Regions
	}
	return &pub
}

// SubscribeSNSToSQSCloudaccountPayload is the cloudaccount subscribeSNSToSQS action payload.
type SubscribeSNSToSQSCloudaccountPayload struct {
	Regions []string `form:"regions" json:"regions" xml:"regions"`
}

// Validate runs the validation rules defined in the design.
func (payload *SubscribeSNSToSQSCloudaccountPayload) Validate() (err error) {
	if payload.Regions == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "regions"))
	}
	return
}

// OK sends a HTTP response with status code 200.
func (ctx *SubscribeSNSToSQSCloudaccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// UpdateCloudaccountContext provides the cloudaccount update action context.
type UpdateCloudaccountContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
	Payload        *UpdateCloudaccountPayload
}

// NewUpdateCloudaccountContext parses the incoming request URL and body, performs validations and creates the
// context used by the cloudaccount controller update action.
func NewUpdateCloudaccountContext(ctx context.Context, service *goa.Service) (*UpdateCloudaccountContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := UpdateCloudaccountContext{Context: ctx, ResponseData: resp, RequestData: req}
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

// updateCloudaccountPayload is the cloudaccount update action payload.
type updateCloudaccountPayload struct {
	AwsID                *string `form:"aws_id,omitempty" json:"aws_id,omitempty" xml:"aws_id,omitempty"`
	DefaultLeaseDuration *string `form:"default_lease_duration,omitempty" json:"default_lease_duration,omitempty" xml:"default_lease_duration,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *updateCloudaccountPayload) Validate() (err error) {
	if payload.DefaultLeaseDuration == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "default_lease_duration"))
	}
	if payload.AwsID != nil {
		if utf8.RuneCountInString(*payload.AwsID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.aws_id`, *payload.AwsID, utf8.RuneCountInString(*payload.AwsID), 1, true))
		}
	}
	return
}

// Publicize creates UpdateCloudaccountPayload from updateCloudaccountPayload
func (payload *updateCloudaccountPayload) Publicize() *UpdateCloudaccountPayload {
	var pub UpdateCloudaccountPayload
	if payload.AwsID != nil {
		pub.AwsID = payload.AwsID
	}
	if payload.DefaultLeaseDuration != nil {
		pub.DefaultLeaseDuration = *payload.DefaultLeaseDuration
	}
	return &pub
}

// UpdateCloudaccountPayload is the cloudaccount update action payload.
type UpdateCloudaccountPayload struct {
	AwsID                *string `form:"aws_id,omitempty" json:"aws_id,omitempty" xml:"aws_id,omitempty"`
	DefaultLeaseDuration string  `form:"default_lease_duration" json:"default_lease_duration" xml:"default_lease_duration"`
}

// Validate runs the validation rules defined in the design.
func (payload *UpdateCloudaccountPayload) Validate() (err error) {
	if payload.DefaultLeaseDuration == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "default_lease_duration"))
	}
	if payload.AwsID != nil {
		if utf8.RuneCountInString(*payload.AwsID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`raw.aws_id`, *payload.AwsID, utf8.RuneCountInString(*payload.AwsID), 1, true))
		}
	}
	return
}

// OK sends a HTTP response with status code 200.
func (ctx *UpdateCloudaccountContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
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

// DeleteFromDBLeasesContext provides the leases deleteFromDB action context.
type DeleteFromDBLeasesContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
	LeaseID        int
}

// NewDeleteFromDBLeasesContext parses the incoming request URL and body, performs validations and creates the
// context used by the leases controller deleteFromDB action.
func NewDeleteFromDBLeasesContext(ctx context.Context, service *goa.Service) (*DeleteFromDBLeasesContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := DeleteFromDBLeasesContext{Context: ctx, ResponseData: resp, RequestData: req}
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
	paramLeaseID := req.Params["lease_id"]
	if len(paramLeaseID) > 0 {
		rawLeaseID := paramLeaseID[0]
		if leaseID, err2 := strconv.Atoi(rawLeaseID); err2 == nil {
			rctx.LeaseID = leaseID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("lease_id", rawLeaseID, "integer"))
		}
		if rctx.LeaseID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`lease_id`, rctx.LeaseID, 1, true))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *DeleteFromDBLeasesContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// ListLeasesForAccountLeasesContext provides the leases listLeasesForAccount action context.
type ListLeasesForAccountLeasesContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID  int
	Terminated *bool
}

// NewListLeasesForAccountLeasesContext parses the incoming request URL and body, performs validations and creates the
// context used by the leases controller listLeasesForAccount action.
func NewListLeasesForAccountLeasesContext(ctx context.Context, service *goa.Service) (*ListLeasesForAccountLeasesContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ListLeasesForAccountLeasesContext{Context: ctx, ResponseData: resp, RequestData: req}
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
	paramTerminated := req.Params["terminated"]
	if len(paramTerminated) > 0 {
		rawTerminated := paramTerminated[0]
		if terminated, err2 := strconv.ParseBool(rawTerminated); err2 == nil {
			tmp24 := &terminated
			rctx.Terminated = tmp24
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("terminated", rawTerminated, "boolean"))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ListLeasesForAccountLeasesContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// ListLeasesForCloudaccountLeasesContext provides the leases listLeasesForCloudaccount action context.
type ListLeasesForCloudaccountLeasesContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
	Terminated     *bool
}

// NewListLeasesForCloudaccountLeasesContext parses the incoming request URL and body, performs validations and creates the
// context used by the leases controller listLeasesForCloudaccount action.
func NewListLeasesForCloudaccountLeasesContext(ctx context.Context, service *goa.Service) (*ListLeasesForCloudaccountLeasesContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ListLeasesForCloudaccountLeasesContext{Context: ctx, ResponseData: resp, RequestData: req}
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
	paramTerminated := req.Params["terminated"]
	if len(paramTerminated) > 0 {
		rawTerminated := paramTerminated[0]
		if terminated, err2 := strconv.ParseBool(rawTerminated); err2 == nil {
			tmp27 := &terminated
			rctx.Terminated = tmp27
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("terminated", rawTerminated, "boolean"))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ListLeasesForCloudaccountLeasesContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// SetExpiryLeasesContext provides the leases setExpiry action context.
type SetExpiryLeasesContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
	ExpiresAt      time.Time
	LeaseID        int
}

// NewSetExpiryLeasesContext parses the incoming request URL and body, performs validations and creates the
// context used by the leases controller setExpiry action.
func NewSetExpiryLeasesContext(ctx context.Context, service *goa.Service) (*SetExpiryLeasesContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := SetExpiryLeasesContext{Context: ctx, ResponseData: resp, RequestData: req}
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
	paramExpiresAt := req.Params["expires_at"]
	if len(paramExpiresAt) == 0 {
		err = goa.MergeErrors(err, goa.MissingParamError("expires_at"))
	} else {
		rawExpiresAt := paramExpiresAt[0]
		if expiresAt, err2 := time.Parse(time.RFC3339, rawExpiresAt); err2 == nil {
			rctx.ExpiresAt = expiresAt
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("expires_at", rawExpiresAt, "datetime"))
		}
	}
	paramLeaseID := req.Params["lease_id"]
	if len(paramLeaseID) > 0 {
		rawLeaseID := paramLeaseID[0]
		if leaseID, err2 := strconv.Atoi(rawLeaseID); err2 == nil {
			rctx.LeaseID = leaseID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("lease_id", rawLeaseID, "integer"))
		}
		if rctx.LeaseID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`lease_id`, rctx.LeaseID, 1, true))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *SetExpiryLeasesContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// ShowLeasesContext provides the leases show action context.
type ShowLeasesContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
	LeaseID        int
}

// NewShowLeasesContext parses the incoming request URL and body, performs validations and creates the
// context used by the leases controller show action.
func NewShowLeasesContext(ctx context.Context, service *goa.Service) (*ShowLeasesContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ShowLeasesContext{Context: ctx, ResponseData: resp, RequestData: req}
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
	paramLeaseID := req.Params["lease_id"]
	if len(paramLeaseID) > 0 {
		rawLeaseID := paramLeaseID[0]
		if leaseID, err2 := strconv.Atoi(rawLeaseID); err2 == nil {
			rctx.LeaseID = leaseID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("lease_id", rawLeaseID, "integer"))
		}
		if rctx.LeaseID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`lease_id`, rctx.LeaseID, 1, true))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowLeasesContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// TerminateLeasesContext provides the leases terminate action context.
type TerminateLeasesContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AccountID      int
	CloudaccountID int
	LeaseID        int
}

// NewTerminateLeasesContext parses the incoming request URL and body, performs validations and creates the
// context used by the leases controller terminate action.
func NewTerminateLeasesContext(ctx context.Context, service *goa.Service) (*TerminateLeasesContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := TerminateLeasesContext{Context: ctx, ResponseData: resp, RequestData: req}
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
	paramLeaseID := req.Params["lease_id"]
	if len(paramLeaseID) > 0 {
		rawLeaseID := paramLeaseID[0]
		if leaseID, err2 := strconv.Atoi(rawLeaseID); err2 == nil {
			rctx.LeaseID = leaseID
		} else {
			err = goa.MergeErrors(err, goa.InvalidParamTypeError("lease_id", rawLeaseID, "integer"))
		}
		if rctx.LeaseID < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`lease_id`, rctx.LeaseID, 1, true))
		}
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *TerminateLeasesContext) OK(resp []byte) error {
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
