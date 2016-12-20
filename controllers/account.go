package controllers

import (
	"fmt"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
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

// Create handles the endpoint used to create a new account on Cecil
func (c *AccountController) Create(ctx *app.CreateAccountContext) error {
	requestContextLogger := core.Logger.New(
		"url", ctx.Request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)

	// validate email
	newAccountInputEmail, err := c.cs.DefaultMailer.Client.ValidateEmail(ctx.Payload.Email)
	if err != nil {
		requestContextLogger.Error("Error while verifying email", "err", err)
		return core.ErrInternal(ctx, "error while verifying email; please retry")
	}
	if !newAccountInputEmail.IsValid {
		requestContextLogger.Error("Invalid email")
		return core.ErrInvalidRequest(ctx, "invalid email; please retry")
	}

	var account core.Account = core.Account{}
	// check whether an account with this email address already exists
	existingAccount, emailAlreadyRegistered, err := c.cs.AccountByEmailExists(newAccountInputEmail.Address)
	if err != nil {
		requestContextLogger.Error("Internal server error", "err", err)
		return core.ErrInternal(ctx, "internal server error; please retry")
	}
	if emailAlreadyRegistered && existingAccount != nil {
		// use the existing account
		account = *existingAccount
		account.RequestedNewToken = true
	} else {
		// create the new account with the provided data
		account.Email = newAccountInputEmail.Address
		account.Name = strings.TrimSpace(ctx.Payload.Name)
		account.Surname = strings.TrimSpace(ctx.Payload.Surname)
	}

	// verificationToken will be used to verify the account
	verificationToken := fmt.Sprintf(
		"%v%v%v",
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
	)

	expectedVerificationTokenSize := 108 // @gagliardetto: why 108?
	if len(verificationToken) < expectedVerificationTokenSize {
		requestContextLogger.Error("internal exception: len(verificationToken) < 108; SOMETHING'S WRONG WITH uuid.NewV4().String()")
		return core.ErrInternal(ctx, "internal exception; please retry")
	}
	if emailAlreadyRegistered {
		core.Logger.Debug("CreateAccount; get new API token", "verification_token", fmt.Sprintf("%v", verificationToken))
	} else {
		core.Logger.Debug("CreateAccount", "verification_token", fmt.Sprintf("%v", verificationToken))
	}

	account.VerificationToken = verificationToken

	if emailAlreadyRegistered {
		// save existing account
		if err := c.cs.DB.Save(&account).Error; err != nil {
			requestContextLogger.Error("Error while saving updated account", "err", err)
			return core.ErrInternal(ctx, "error while updating account; please retry")
		}
	} else {
		// create new account
		if err := c.cs.DB.Create(&account).Error; err != nil {
			requestContextLogger.Error("Error while saving new account", "err", err)
			return core.ErrInternal(ctx, "error while creating account; please retry")
		}
	}

	verificationTargetURL := fmt.Sprintf("%v/accounts/%v/api_token", c.cs.CecilHTTPAddress(), account.ID)

	var emailSubject string
	var isVerifyingAccount bool = !emailAlreadyRegistered

	if emailAlreadyRegistered {
		emailSubject = "Get another API token"
	} else {
		emailSubject = "Activate account and get API token"
	}

	newEmailBody := core.CompileEmail(
		`Hey {{.account_name}}, to{{ if .isVerifyingAccount }} verify your account and{{ end }} create an API token,
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
			"isVerifyingAccount":      isVerifyingAccount,
			"account_name":            account.Name,
			"verification_target_url": verificationTargetURL,
			"verification_token":      account.VerificationToken,
		},
	)
	c.cs.NotifierQueue.TaskQueue <- core.NotifierTask{
		To:       newAccountInputEmail.Address,
		Subject:  emailSubject,
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
	}

	return core.JSONResponse(ctx, 200, gin.H{
		"response":   "An email has been sent to the specified address with a verification token and instructions.",
		"account_id": account.ID,
		"email":      newAccountInputEmail.Address,
		"verified":   account.Verified,
	})

}

// Show handles the endpoint to show the info about an account (only the account the user is logged in to).
func (c *AccountController) Show(ctx *app.ShowAccountContext) error {
	requestContextLogger := core.Logger.New(
		"url", ctx.Request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)

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
		} else {
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving account %v. See logs for details", ctx.AccountID))
		}
	}

	return core.JSONResponse(ctx, 200, account)
}

// Verify handles the endpoint used to verify/get new API token for an account with a verification token sent via email,
// and the token must match the one in the DB
func (c *AccountController) Verify(ctx *app.VerifyAccountContext) error {
	requestContextLogger := core.Logger.New(
		"url", ctx.Request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)
	// TODO: add nonce to this url to NOT allow anyone to verify which accounts are active and which are not

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error while fetching account", "err", err)
		return core.ErrInvalidRequest(ctx, "account with that id does not exist")
	}

	if account.Verified && !account.RequestedNewToken {
		msg := fmt.Sprintf("account %v already verified, and not requested for new token", ctx.AccountID)
		requestContextLogger.Error(msg)
		return core.ErrInvalidRequest(ctx, msg)
	}

	verificationTokenIsNOTLongEnough := len(account.VerificationToken) < 108
	if verificationTokenIsNOTLongEnough {
		// TODO: notify admins
		requestContextLogger.Error(fmt.Sprintf("Verification token (%s) not long enough. Expected 108, got %d", account.VerificationToken, len(account.VerificationToken)))
		return core.ErrInternal(ctx, "internal exception error")
	}

	verificationTokensDoNOTMatch := strings.EqualFold(ctx.Payload.VerificationToken, account.VerificationToken)
	if verificationTokensDoNOTMatch {
		requestContextLogger.Error("The verification token in DB and the one from the request do not match")
		return core.ErrInvalidRequest(ctx, "cannot verify account")
	}

	// mark account as verified
	account.Verified = true

	// mark RequestedNewToken as false because the request has been fulfilled
	account.RequestedNewToken = false
	// remove verification token
	// account.VerificationToken = "" // WARNING: this goes against the UNIQUE db constraint

	// commit to db the account
	if err := c.cs.DB.Save(&account).Error; err != nil {
		requestContextLogger.Error("Error while saving account", "err", err)
		return core.ErrInternal(ctx, core.ErrorInternal)
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

		"sub":    uint(ctx.AccountID), // the subject/principal is whom the token is about
		"scopes": "api:access",        // token scope - not a standard claim
	}

	// sign token
	APIToken, err := c.cs.SignToken(token)
	if err != nil {
		requestContextLogger.Error("Error while signing token", "err", err)
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	return core.JSONResponse(ctx, 200, gin.H{
		"account_id": account.ID,
		"email":      account.Email,
		"verified":   account.Verified,
		"api_token":  "Bearer " + APIToken,
	})
}

// SlackConfig handles the endpoint used to add slack configuration to an account.
// That slack config will be used to start a SlackInstance that will send messages
// to a channel and receive commands.
func (c *AccountController) SlackConfig(ctx *app.SlackConfigAccountContext) error {
	requestContextLogger := core.Logger.New(
		"url", ctx.Request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)

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
		} else {
			return core.ErrInternal(ctx, core.ErrorInternal)
		}
	}

	// TODO: better validate payload

	newSlackConfig := core.SlackConfig{
		AccountID: uint(ctx.AccountID),
		Token:     strings.TrimSpace(ctx.Payload.Token),
		ChannelID: strings.TrimSpace(ctx.Payload.ChannelID),
	}

	err = c.cs.DB.Create(&newSlackConfig).Error
	if err != nil {
		requestContextLogger.Error("Error saving new slack config", "err", err)
		return core.ErrInternal(ctx, core.ErrorInternal)
	}

	// stop the eventual existing slack instance
	c.cs.TerminateSlackInstance(uint(ctx.AccountID))

	// start slack
	err = c.cs.StartSlackInstance(&newSlackConfig)
	if err != nil {
		requestContextLogger.Error("Error while starting slack instance", "err", err)
		return core.ErrInternal(ctx, "Internal server error starting slack instance. See logs for details")
	}

	var success struct {
		Message string `json:"message"`
	}
	success.Message = "Slack added to account"

	return core.JSONResponse(ctx, 200, success)
}

// RemoveSlack handles the endpoint used to remove slack from an account.
// The eventual running slack instance will be stopped.
func (c *AccountController) RemoveSlack(ctx *app.RemoveSlackAccountContext) error {
	requestContextLogger := core.Logger.New(
		"url", ctx.Request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)

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
		} else {
			return core.ErrInternal(ctx, "internal server error 1")
		}
	}

	// fetch existing slack config
	slackConfig, err := c.cs.FetchSlackConfig(uint(ctx.AccountID))
	if err != nil {
		requestContextLogger.Error("Error fetching slack config from DB", "err", err)
		return core.ErrInternal(ctx, err.Error())
	}

	// delete slack config from DB
	err = c.cs.DB.Delete(&slackConfig).Error
	if err != nil {
		requestContextLogger.Error("Error deleting slack config from DB", "err", err)
		return core.ErrInternal(ctx, err.Error())
	}

	// stop the eventual existing slack instance
	c.cs.TerminateSlackInstance(uint(ctx.AccountID))

	var success struct {
		Message string `json:"message"`
	}
	success.Message = "Slack removed from account"

	return core.JSONResponse(ctx, 200, success)
}

// MailerConfig handles the endpoint used to add a custom mailgun mailer instance
// for an account; this mailer instance will be used instead of the default one.
func (c *AccountController) MailerConfig(ctx *app.MailerConfigAccountContext) error {
	requestContextLogger := core.Logger.New(
		"url", ctx.Request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)

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
		} else {
			return core.ErrInternal(ctx, fmt.Sprintf("Internal server error retrieving account %v. See logs for details", ctx.AccountID))
		}
	}

	// TODO: better validate payload

	newMailerConfig := core.MailerConfig{
		AccountID: uint(ctx.AccountID),

		Domain:       strings.TrimSpace(ctx.Payload.Domain),
		APIKey:       strings.TrimSpace(ctx.Payload.APIKey),
		PublicAPIKey: strings.TrimSpace(ctx.Payload.PublicAPIKey),
		FromName:     strings.TrimSpace(ctx.Payload.FromName),
	}

	err = c.cs.DB.Create(&newMailerConfig).Error
	if err != nil {
		requestContextLogger.Error("Error saving new mailer config", "err", err)
		return core.ErrInternal(ctx, err.Error())
	}

	// stop the eventual existing mailer instance
	c.cs.TerminateMailerInstance(uint(ctx.AccountID))

	// start new mailer instance
	err = c.cs.StartMailerInstance(&newMailerConfig)
	if err != nil {
		requestContextLogger.Error("Error starting mailer instance", "err", err)
		return core.ErrInternal(ctx, "error while starting the mailer")
	}

	var success struct {
		Message string `json:"message"`
	}
	from := fmt.Sprintf("%v <noreply@%v>", newMailerConfig.FromName, newMailerConfig.Domain)
	success.Message = fmt.Sprintf(`mailer added/modified; emails will come from '%v'`, from)

	return core.JSONResponse(ctx, 200, success)
}
