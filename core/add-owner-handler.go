package core

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Service) AddOwnerHandler(c *gin.Context) {
	// validate parameters (are set, their type, ...)
	// check whether cloudaccount exists
	// validate email
	// save to db

	// TODO: only allow adding an owner if the user logged in is account_id
	// parse parameters
	account_id, err := strconv.ParseUint(c.Param("account_id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	cloudaccount_id, err := strconv.ParseUint(c.Param("cloudaccount_id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	// TODO: figure out why it always finds one result, even if non are in the db
	// check whether the account exists
	var accountCount int64
	var account Account
	s.DB.First(&account, uint(account_id)).Count(&accountCount)
	if accountCount != 1 {
		c.JSON(404, gin.H{
			"error": "not found",
		})
		return
	}

	// TODO: figure out why it always finds one result, even if non are in the db
	// check whether the cloudaccount exists
	var cloudAccountCount int64
	var cloudAccount CloudAccount
	s.DB.First(&cloudAccount, uint(cloudaccount_id)).Count(&cloudAccountCount)
	if cloudAccountCount != 1 {
		c.JSON(404, gin.H{
			"error": "not found",
		})
		return
	}

	// check if everything is consistent
	if !(uint(account_id) == account.ID && uint(cloudaccount_id) == cloudAccount.ID && account.ID == cloudAccount.AccountID) {
		c.JSON(404, gin.H{
			"error": "error",
		})
		return
	}

	// parse json payload
	var newOwnerInput struct {
		Email string `json:"email"`
	}
	if err := c.BindJSON(&newOwnerInput); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	if newOwnerInput.Email == "" {
		c.JSON(400, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	// validate email
	ownerTag, err := s.Mailer.ValidateEmail(newOwnerInput.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "internal error",
		})
		return
	}
	if !ownerTag.IsValid {
		c.JSON(400, gin.H{
			"error": "invalid email",
		})
		return
	}

	// check whether this owner already exists for this cloudaccount
	var equalOwnerCount int64
	s.DB.Table("owners").Where(&Owner{CloudAccountID: cloudAccount.ID, Email: ownerTag.Address}).Count(&equalOwnerCount)
	if equalOwnerCount != 0 {
		c.JSON(400, gin.H{
			"error": "owner already exists",
		})
		return
	}

	// instert the owner into the db
	newOwner := Owner{
		CloudAccountID: cloudAccount.ID,
		Email:          ownerTag.Address,
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
