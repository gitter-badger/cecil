package controllers

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/tleyden/zerocloud/core"
	"github.com/tleyden/zerocloud/goa/app"
)

// CloudaccountController implements the cloudaccount resource.
type CloudaccountController struct {
	*goa.Controller
	cs *core.Service
}

// NewCloudaccountController creates a cloudaccount controller.
func NewCloudaccountController(service *goa.Service, cs *core.Service) *CloudaccountController {
	return &CloudaccountController{
		Controller: service.NewController("CloudaccountController"),
		cs:         cs,
	}
}

// Add runs the add action.
func (c *CloudaccountController) Add(ctx *app.AddCloudaccountContext) error {

	_, err := core.ValidateToken(ctx)
	if err != nil {
		return core.ErrUnauthorized(ctx, "unauthorized")
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.ErrNotFound(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// TODO: validate newCloudAccountInput.AWSID

	AWSIDAlreadyRegistered, err := c.cs.CloudAccountByAWSIDExists(ctx.Payload.AwsID)
	if err != nil {
		return core.ErrInternal(ctx, "internal server error")
	}
	if AWSIDAlreadyRegistered {
		return core.ErrInvalidRequest(ctx, fmt.Sprintf("cannot add aws %v", ctx.Payload.AwsID))
	}

	externalID := fmt.Sprintf("%v-%v-%v", uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String())
	// TODO: make sure externalID is not null

	// add newCloudAccount to DB
	newCloudAccount := core.CloudAccount{
		AccountID:  account.ID,
		Provider:   "aws",
		AWSID:      ctx.Payload.AwsID,
		ExternalID: externalID,
	}
	err = c.cs.DB.Create(&newCloudAccount).Error
	if err != nil {
		return core.ErrInternal(ctx, "internal server error")
	}

	firstOwner := core.Owner{
		Email:          account.Email,
		CloudAccountID: newCloudAccount.ID,
	}
	err = c.cs.DB.Create(&firstOwner).Error
	if err != nil {
		return core.ErrInternal(ctx, "internal server error")
	}

	// regenerate SQS permissions
	if err := c.cs.RegenerateSQSPermissions(); err != nil {
		return core.ErrInternal(ctx, "internal server error")
	}

	return ctx.Service.Send(ctx, 200, gin.H{
		"cloudaccount_id": newCloudAccount.ID,
		"aws_id":          newCloudAccount.AWSID,
		"initial_setup_cloudformation_url": fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-initial-setup.template", account.ID, newCloudAccount.ID),
		"region_setup_cloudformation_url":  fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-region-setup.template", account.ID, newCloudAccount.ID),
	})
}

// AddEmailToWhitelist runs the addEmailToWhitelist action.
func (c *CloudaccountController) AddEmailToWhitelist(ctx *app.AddEmailToWhitelistCloudaccountContext) error {

	_, err := core.ValidateToken(ctx)
	if err != nil {
		return core.ErrUnauthorized(ctx, "unauthorized")
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	cloudAccount, err := c.cs.FetchCloudAccountByID(ctx.CloudaccountID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	// validate email
	ownerEmail, err := c.cs.Mailer.Client.ValidateEmail(ctx.Payload.Email)
	if err != nil {
		return core.ErrInternal(ctx, "internal server error")
	}
	if !ownerEmail.IsValid {
		return core.ErrInvalidRequest(ctx, "invalid email")
	}

	// check whether this owner already exists for this cloudaccount
	var equalOwnerCount int64
	c.cs.DB.Table("owners").Where(&core.Owner{CloudAccountID: cloudAccount.ID, Email: ownerEmail.Address}).Count(&equalOwnerCount)
	if equalOwnerCount != 0 {
		return core.ErrInvalidRequest(ctx, "owner already exists in whitelist")
	}

	// instert the new owner into the db
	newOwner := core.Owner{
		CloudAccountID: cloudAccount.ID,
		Email:          ownerEmail.Address,
	}
	err = c.cs.DB.Create(&newOwner).Error

	if err != nil {
		return core.ErrInternal(ctx, "internal server error")
	}

	return ctx.Service.Send(ctx, 200, gin.H{
		"message": "owner added successfully to whitelist",
	})
	return nil
}

// DownloadInitialSetupTemplate runs the downloadInitialSetupTemplate action.
func (c *CloudaccountController) DownloadInitialSetupTemplate(ctx *app.DownloadInitialSetupTemplateCloudaccountContext) error {

	_, err := core.ValidateToken(ctx)
	if err != nil {
		return core.ErrUnauthorized(ctx, "unauthorized")
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	cloudAccount, err := c.cs.FetchCloudAccountByID(ctx.CloudaccountID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	var compiledTemplate bytes.Buffer // A Buffer needs no initialization.

	tpl, err := template.ParseFiles("./core/go-templates/tenant-aws-initial-setup.template")
	if err != nil {
		core.Logger.Error("1:", "error", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	var values map[string]interface{} = map[string]interface{}{}
	values["IAMRoleExternalID"] = cloudAccount.ExternalID
	values["CecilAWSID"] = c.cs.AWS.Config.AWS_ACCOUNT_ID

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		core.Logger.Error("2:", "error", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	return ctx.OK(compiledTemplate.Bytes())
}

// DownloadRegionSetupTemplate runs the downloadRegionSetupTemplate action.
func (c *CloudaccountController) DownloadRegionSetupTemplate(ctx *app.DownloadRegionSetupTemplateCloudaccountContext) error {

	_, err := core.ValidateToken(ctx)
	if err != nil {
		return core.ErrUnauthorized(ctx, "unauthorized")
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	cloudAccount, err := c.cs.FetchCloudAccountByID(ctx.CloudaccountID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	var compiledTemplate bytes.Buffer // A Buffer needs no initialization.

	tpl, err := template.ParseFiles("./core/go-templates/tenant-aws-region-setup.template")
	if err != nil {
		core.Logger.Error("1:", "error", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	var values map[string]interface{} = map[string]interface{}{}
	values["CecilAWSID"] = c.cs.AWS.Config.AWS_ACCOUNT_ID

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		core.Logger.Error("2:", "error", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	return ctx.OK(compiledTemplate.Bytes())
}
