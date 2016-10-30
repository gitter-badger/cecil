package core

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func (s *Service) CloudformationInitialSetupHandler(c *gin.Context) {

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

	cloudAccount, err := s.FetchCloudAccountByID(c.Param("cloudaccount_id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{
				"error": fmt.Sprintf("cloud account with id %v does not exist", c.Param("cloudaccount_id")),
			})
			return
		} else {
			c.JSON(400, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		c.JSON(404, gin.H{
			"error": "internal server error",
		})
		return
	}

	var compiledTemplate bytes.Buffer // A Buffer needs no initialization.

	tpl, err := template.ParseFiles("./core/go-templates/cecil-aws-initial-setup.template")
	if err != nil {
		logger.Error("1:", "error", err)
		c.JSON(404, gin.H{
			"error": "internal server error",
		})
		return
	}

	var values map[string]interface{} = map[string]interface{}{}
	values["IAMRoleExternalID"] = cloudAccount.ExternalID
	values["ZeroCloudAWSID"] = s.AWS.Config.AWS_ACCOUNT_ID

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		logger.Error("2:", "error", err)
		c.JSON(404, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.Data(200, "text/plain", compiledTemplate.Bytes())
}

func (s *Service) CloudformationRegionSetupHandler(c *gin.Context) {

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

	cloudAccount, err := s.FetchCloudAccountByID(c.Param("cloudaccount_id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{
				"error": fmt.Sprintf("cloud account with id %v does not exist", c.Param("cloudaccount_id")),
			})
			return
		} else {
			c.JSON(400, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		c.JSON(404, gin.H{
			"error": "internal server error",
		})
		return
	}

	var compiledTemplate bytes.Buffer // A Buffer needs no initialization.

	tpl, err := template.ParseFiles("./core/go-templates/zerocloud-aws-region-setup.template")
	if err != nil {
		logger.Error("1:", "error", err)
		c.JSON(404, gin.H{
			"error": "internal server error",
		})
		return
	}

	var values map[string]interface{} = map[string]interface{}{}
	values["ZeroCloudAWSID"] = s.AWS.Config.AWS_ACCOUNT_ID

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		logger.Error("2:", "error", err)
		c.JSON(404, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.Data(200, "text/plain", compiledTemplate.Bytes())

}
