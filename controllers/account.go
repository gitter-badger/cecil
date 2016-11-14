package controllers

import (
	"fmt"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/tleyden/zerocloud/core"
	"github.com/tleyden/zerocloud/goa/app"
)

// AccountController implements the account resource.
type AccountController struct {
	*goa.Controller
	cs *core.Service
}

// NewAccountController creates a account controller.
func NewAccountController(service *goa.Service, cs *core.Service) *AccountController {
	return &AccountController{
		Controller: service.NewController("AccountController"),
		cs:         cs,
	}
}

// Create runs the create action.
func (c *AccountController) Create(ctx *app.CreateAccountContext) error {
	requestID := middleware.ContextRequestID(ctx)
	requestContextLog := core.Logger.New(
		"Request type", "Request to create an account",
		"url", ctx.Request.URL.String(),
		"reqID", requestID,
	)

	// validate email
	newAccountInputEmail, err := c.cs.Mailer.Client.ValidateEmail(ctx.Payload.Email)
	if err != nil {
		requestContextLog.Error("error while verifying email", "err", err)
		return core.ErrInternal(ctx, "error while verifying email; please retry")
	}
	if !newAccountInputEmail.IsValid {
		requestContextLog.Error("invalid email")
		return core.ErrInvalidRequest(ctx, "invalid email; please retry")
	}

	// check whether an account with this email address already exists
	emailAlreadyRegistered, err := c.cs.AccountByEmailExists(newAccountInputEmail.Address)
	if err != nil {
		requestContextLog.Error("internal server error", "err", err)
		return core.ErrInternal(ctx, "nternal server error; please retry")
	}
	if emailAlreadyRegistered {
		requestContextLog.Error("account with provided email already exists")
		return core.ErrInvalidRequest(ctx, fmt.Sprintf("%v already signed up", newAccountInputEmail.Address))
	}

	// verificationToken will be used to verify the account
	verificationToken := fmt.Sprintf("%v%v%v", uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String())
	if len(verificationToken) < 108 {
		requestContextLog.Error("internal exception: len(verificationToken) < 108; SOMETHING'S WRONG WITH uuid.NewV4().String()")
		return core.ErrInternal(ctx, "internal exception; please retry")
	}
	// create the new account
	var newAccount core.Account = core.Account{}

	newAccount.Email = newAccountInputEmail.Address
	newAccount.Name = strings.TrimSpace(ctx.Payload.Name)
	newAccount.Surname = strings.TrimSpace(ctx.Payload.Surname)
	newAccount.VerificationToken = verificationToken

	if err := c.cs.DB.Create(&newAccount).Error; err != nil {
		requestContextLog.Error("internal server error", "err", err)
		return core.ErrInternal(ctx, "error while creating account; please retry")
	}

	verificationTargetURL := fmt.Sprintf("%v/accounts/%v/api_token", c.cs.CecilHTTPAddress(), newAccount.ID)

	newEmailBody := core.CompileEmail(
		`Hey {{.account_name}}, to verify your account and create an API token,
		send a POST request to <b>{{.verification_target_url}}</b> with this JSON payload:
		<br>

		<br>{"verification_token":"{{.verification_token}}"}<br>

		<br>
		CURL Example:
		<br>
		<br>

		curl \
		<br>
		-H "Content-Type: application/json" \
		<br>
		-X POST \
		<br>
		-d '{"verification_token":"{{.verification_token}}"}' \
		<br>
		{{.verification_target_url}}

		<br>
		<br>
		Thanks for using Cecil!
				`,

		map[string]interface{}{
			"account_name":            newAccount.Name,
			"verification_target_url": verificationTargetURL,
			"verification_token":      newAccount.VerificationToken,
		},
	)
	c.cs.NotifierQueue.TaskQueue <- core.NotifierTask{
		From:     c.cs.Mailer.FromAddress,
		To:       newAccountInputEmail.Address,
		Subject:  "Activate account and get API token",
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
	}

	return ctx.Service.Send(ctx, 200, gin.H{
		"response":   "An email has been sent to the specified address with a verification token and instructions.",
		"account_id": newAccount.ID,
		"email":      newAccountInputEmail.Address,
		"verified":   false,
	})

}

// Show runs the show action.
func (c *AccountController) Show(ctx *app.ShowAccountContext) error {
	_, err := core.ValidateToken(ctx)
	if err != nil {
		return core.ErrUnauthorized(ctx, "unauthorized")
	}
	// TODO: return account info

	return nil
}

// Verify runs the verify action.
func (c *AccountController) Verify(ctx *app.VerifyAccountContext) error {

	// TODO: add nonce to this url to NOT allow anyone to verify which accounts are active and which are not

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		return core.ErrInvalidRequest(ctx, "account with that id does not exist")
	}

	if account.Verified {
		return core.ErrInvalidRequest(ctx, "account already verified")
	}

	if len(strings.TrimSpace(account.VerificationToken)) < 108 {
		// TODO: notify ZC admins
		return core.ErrInternal(ctx, "internal exception error")
	}

	if ctx.Payload.VerificationToken != account.VerificationToken {
		return core.ErrInvalidRequest(ctx, "cannot verify account")
	}

	// mark account as verified
	account.Verified = true
	// remove verification token
	// account.VerificationToken = "" // WARNING: this goes against the UNIQUE db constraint

	// commit to db the account
	if err := c.cs.DB.Save(&account).Error; err != nil {
		return core.ErrInternal(ctx, "internal server error")
	}

	// declare new token
	token := jwtgo.New(jwtgo.SigningMethodRS512)

	sevenDays := time.Duration(24*7) * time.Hour
	// decide expiry
	tokenExpiresAt := time.Now().UTC().Add(sevenDays).Unix()

	token.Claims = jwtgo.MapClaims{
		"iss": "cecil-api-backend",     // who creates the token and signs it
		"aud": "cecil-account",         // to whom the token is intended to be sent
		"exp": tokenExpiresAt,          // time when the token will expire (time from now)
		"jti": uuid.NewV4().String(),   // a unique identifier for the token
		"iat": time.Now().UTC().Unix(), // when the token was issued/created (now)
		"nbf": 3,                       // time before which the token is not yet valid (2 minutes ago)

		"sub":    ctx.AccountID, // the subject/principal is whom the token is about
		"scopes": "api:access",  // token scope - not a standard claim
	}

	// sign token
	APIToken, err := c.cs.SignToken(token)
	if err != nil {
		return core.ErrInternal(ctx, "internal server error")
	}

	return ctx.Service.Send(ctx, 200, gin.H{
		"account_id": account.ID,
		"email":      account.Email,
		"verified":   account.Verified,
		"api_token":  "Bearer " + APIToken,
	})
}

// TODO: add a way to regenerate the API token
