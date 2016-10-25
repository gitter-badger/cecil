package core

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// AddCloudAccountHandler accepts a request to add a new CloudAccount to an account
func (s *Service) AddCloudAccountHandler(c *gin.Context) {
	/*
	   REQUEST:

	   POST /accounts/:account_id/cloudaccounts

	   {
	       "aws_id":"012345677"
	   }

	   RESPONSE:

	   {
	       "id": 1,
	       "aws_id": "012345677",
	       "initial_setup_cloudformation_url": "/accounts/1/cloudaccounts/1/zerocloud-aws-initial-setup.template",
	       "region_setup_cloudformation_url": "/accounts/1/cloudaccounts/1/zerocloud-aws-region-setup.template"
	   }
	*/

	account, err := s.FetchAccountByID(c.Param("account_id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{
				"error": fmt.Sprintf("account with id %v does not exist", c.Param("account_id")),
			})
			return
		} else {
			c.JSON(400, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	// parse json payload
	var newCloudAccountInput struct {
		AWSID string `json:"aws_id" binding:"required"`
	}
	if err := c.BindJSON(&newCloudAccountInput); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	// TODO: validate newCloudAccountInput.AWSID

	AWSIDAlreadyRegistered, err := s.CloudAccountByAWSIDExists(newCloudAccountInput.AWSID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "internal server error",
		})
		return
	}
	if AWSIDAlreadyRegistered {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("cannot add aws %v", newCloudAccountInput.AWSID),
		})
		return
	}

	externalID := fmt.Sprintf("%v-%v-%v", uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String())
	// TODO: make sure externalID is not null

	// add newCloudAccount to DB
	newCloudAccount := CloudAccount{
		AccountID:  account.ID,
		Provider:   "aws",
		AWSID:      newCloudAccountInput.AWSID,
		ExternalID: externalID,
	}
	err = s.DB.Create(&newCloudAccount).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "internal server error",
		})
		return
	}

	firstOwner := Owner{
		Email:          account.Email,
		CloudAccountID: newCloudAccount.ID,
	}
	err = s.DB.Create(&firstOwner).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "internal server error",
		})
		return
	}

	// regenerate SQS permissions
	if err := s.RegenerateSQSPermissions(); err != nil {
		c.JSON(400, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"cloudaccount_id": newCloudAccount.ID,
		"aws_id":          newCloudAccount.AWSID,
		"initial_setup_cloudformation_url": fmt.Sprintf("/accounts/%v/cloudaccounts/%v/zerocloud-aws-initial-setup.template", account.ID, newCloudAccount.ID),
		"region_setup_cloudformation_url":  fmt.Sprintf("/accounts/%v/cloudaccounts/%v/zerocloud-aws-region-setup.template", account.ID, newCloudAccount.ID),
	})
}
