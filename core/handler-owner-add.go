package core

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// AddOwnerHandler accepts a request to add a new owner to a cloudaccount's whitelist
func (s *Service) AddOwnerHandler(c *gin.Context) {
	// validate parameters
	// check whether account exists
	// check whether cloudaccount exists
	// validate email
	// validate owner is not already in the db
	// save to db

	// TODO: only allow adding an owner if the user logged in is account_id

	account, err := s.FetchAccountByID(c.Param("account_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "account does not exist",
		})
		return
	}

	cloudAccount, err := s.FetchCloudAccountByID(c.Param("cloudaccount_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "cloud account does not exist",
		})
		return
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		c.JSON(404, gin.H{
			"error": "error",
		})
		return
	}

	// parse json payload
	var newOwnerInput struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.BindJSON(&newOwnerInput); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	// check if email field is set
	if strings.TrimSpace(newOwnerInput.Email) == "" {
		c.JSON(400, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	// validate email
	ownerEmail, err := s.Mailer.Client.ValidateEmail(newOwnerInput.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "internal error",
		})
		return
	}
	if !ownerEmail.IsValid {
		c.JSON(400, gin.H{
			"error": "invalid email",
		})
		return
	}

	// check whether this owner already exists for this cloudaccount
	var equalOwnerCount int64
	s.DB.Table("owners").Where(&Owner{CloudAccountID: cloudAccount.ID, Email: ownerEmail.Address}).Count(&equalOwnerCount)
	if equalOwnerCount != 0 {
		c.JSON(400, gin.H{
			"error": "owner already exists",
		})
		return
	}

	// instert the new owner into the db
	newOwner := Owner{
		CloudAccountID: cloudAccount.ID,
		Email:          ownerEmail.Address,
	}
	err = s.DB.Create(&newOwner).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error": "internal error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "owner added successfully",
	})
	return
}
