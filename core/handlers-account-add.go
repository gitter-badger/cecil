package core

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func (s *Service) CreateAccountHandler(c *gin.Context) {
	requestID := uuid.NewV4().String()
	requestContextLog := logger.New(
		"Request to create an account",
		"url", c.Request.URL,
		"reqID", requestID,
	)

	// parse json payload
	var newAccountInput struct {
		Email   string `json:"email" binding:"required"`
		Name    string `json:"name" binding:"required"`
		Surname string `json:"surname"` // optional
	}
	if err := c.BindJSON(&newAccountInput); err != nil {
		requestContextLog.Error("cannot parse request payload", "err", err)
		c.JSON(400, gin.H{
			"error": "cannot parse request payload",
		})
		return
	}

	// check if email field is set
	if strings.TrimSpace(newAccountInput.Email) == "" {
		requestContextLog.Error("email not provided")
		c.JSON(400, gin.H{
			"error": "email must be provided",
		})
		return
	}

	// validate email
	newAccountInputEmail, err := s.Mailer.Client.ValidateEmail(newAccountInput.Email)
	if err != nil {
		requestContextLog.Error("error while verifying email", "err", err)
		c.JSON(500, gin.H{
			"error": "error while verifying email; please retry",
		})
		return
	}
	if !newAccountInputEmail.IsValid {
		requestContextLog.Error("invalid email")
		c.JSON(400, gin.H{
			"error": "invalid email; please retry",
		})
		return
	}

	// check max name and surname length
	if len(newAccountInput.Name) > 30 || len(newAccountInput.Surname) > 30 {
		requestContextLog.Error("name or surname too long")
		c.JSON(400, gin.H{
			"error": "name or surname too long",
		})
		return
	}

	// check whether an account with this email address already exists
	emailAlreadyRegistered, err := s.AccountByEmailExists(newAccountInputEmail.Address)
	if err != nil {
		requestContextLog.Error("internal server error", "err", err)
		c.JSON(500, gin.H{
			"error": "internal server error; please retry",
		})
		return
	}
	if emailAlreadyRegistered {
		requestContextLog.Error("account with provided email already exists")
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("%v already signed up", newAccountInputEmail.Address),
		})
		return
	}

	// verificationToken will be used to verify the account
	verificationToken := fmt.Sprintf("%v%v%v", uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String())
	if len(verificationToken) < 108 {
		requestContextLog.Error("internal exception: len(verificationToken) < 108; SOMETHING'S WRONG WITH uuid.NewV4().String()")
		c.JSON(500, gin.H{
			"error": "internal exception; please retry",
		})
		return
	}
	// create the new account
	var newAccount Account = Account{}

	newAccount.Email = newAccountInputEmail.Address
	newAccount.Name = strings.TrimSpace(newAccountInput.Name)
	newAccount.Surname = strings.TrimSpace(newAccountInput.Surname)
	newAccount.VerificationToken = verificationToken

	if err := s.DB.Create(&newAccount).Error; err != nil {
		requestContextLog.Error("internal server error", "err", err)
		c.JSON(500, gin.H{
			"error": "error while creating account; please retry",
		})
		return
	}

	verificationTargetURL := fmt.Sprintf("%v/accounts/%v/api_token", s.ZeroCloudHTTPAddress(), newAccount.ID)

	newEmailBody := compileEmail(
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
		Thanks for using ZeroCloud!
				`,

		map[string]interface{}{
			"account_name":            newAccount.Name,
			"verification_target_url": verificationTargetURL,
			"verification_token":      newAccount.VerificationToken,
		},
	)
	s.NotifierQueue.TaskQueue <- NotifierTask{
		From:     s.Mailer.FromAddress,
		To:       newAccountInputEmail.Address,
		Subject:  "Activate account and get API token",
		BodyHTML: newEmailBody,
		BodyText: newEmailBody,
	}

	c.JSON(200, gin.H{
		"response":   "An email has been sent to the specified address with a verification token and instructions.",
		"account_id": newAccount.ID,
		"email":      newAccountInputEmail.Address,
		"verified":   false,
	})
	return
	/*
				POST /accounts

		REQUEST:
		{
			"email":"example@example.com",
			"name":"Example",
			"surname":"example"
		}

		// validate email
		// check whether there is already an account with that same email address
		// create a new account in db: verified:false, verification_token:78w3t823gt32tg4gt674gt74g..., etc.
		// send email with verification token and instructions
		// return response

		RESPONSE:
		   {
				"id":1,
				"email":"example@example.com",
				"verified":false
		   }
	*/

	/*
				    Email with Verification token +
		           instructions to create API token
	*/

}

func (s *Service) ValidateAccountHandler(c *gin.Context) {
	/*
					 POST /accounts/:account_id/api_token

		REQUEST:
					 {
						"verification_token":"98wtyw4t8h3nc94t34t3gtgc643n7t347gtc396tbgb36"
					 }

		// check verification_token length
		// find in db a non-verified account with that verification_token
		// check whether they match
		// generate api_token

		RESPONSE:
					 {
						"id":1,
						"email":"example@example.com",
						"verified":true
						"api_token":"key-giowg9w9g49tgh439hy9384hy943hy934hy4u39t8439y"
					 }
	*/

	// TODO: add nonce to this url to NOT allow anyone to verify which accounts are active and which are not

	// parse json payload
	var validateAccountInput struct {
		VerificationToken string `json:"verification_token" binding:"required"`
	}

	if err := c.BindJSON(&validateAccountInput); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	validateAccountInput.VerificationToken = strings.TrimSpace(validateAccountInput.VerificationToken)

	if len(validateAccountInput.VerificationToken) < 108 {
		c.JSON(400, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	account, err := s.FetchAccountByID(c.Param("account_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "account with that id does not exist",
		})
		return
	}

	if account.Verified {
		c.JSON(400, gin.H{
			"error": "account already verified",
		})
		return
	}

	if len(strings.TrimSpace(account.VerificationToken)) < 108 {
		// TODO: notify ZC admins
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	if validateAccountInput.VerificationToken != account.VerificationToken {
		c.JSON(400, gin.H{
			"error": "cannot verify account",
		})
		return
	}

	// mark account as verified
	account.Verified = true
	// remove verification token
	// account.VerificationToken = "" // WARNING: this goes against the UNIQUE db constraint

	// commit to db the account
	if err := s.DB.Save(&account).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// generate api token
	var APIToken string

	if APIToken, err = s.GenerateAPITokenForAccount(account.ID); err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"account_id": account.ID,
		"email":      account.Email,
		"verified":   account.Verified,
		"api_token":  APIToken,
	})
}

// TODO: add a way to regenerate the API token
// TODO: add a middleware that verifies the API token
// TODO: add a way to add cloudAccounts
// TODO:
// TODO:
