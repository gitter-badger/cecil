package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/tools"
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

// Add handles the endpoint used to add a cloudaccount to an account.
func (c *CloudaccountController) Add(ctx *app.AddCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrNotFound(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// TODO: validate newCloudaccountInput.AWSID

	AWSIDAlreadyRegistered, err := c.cs.CloudaccountByAWSIDExists(ctx.Payload.AwsID)
	if err != nil {
		requestContextLogger.Error("Error CloudaccountByAWSIDExists", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}
	if AWSIDAlreadyRegistered {
		requestContextLogger.Error("AWSIDAlreadyRegistered", "err", err)
		return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cannot add aws %v", ctx.Payload.AwsID))
	}

	externalID := fmt.Sprintf("%v-%v-%v", uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String())
	// TODO: make sure externalID is not null

	// add newCloudaccount to DB
	newCloudaccount := models.Cloudaccount{
		AccountID:  account.ID,
		Provider:   "aws",
		AWSID:      ctx.Payload.AwsID,
		ExternalID: externalID,
	}

	// check whether the payload contains default_lease_duration
	if ctx.Payload.DefaultLeaseDuration != nil {
		defaultLeaseDuration, err := time.ParseDuration(*ctx.Payload.DefaultLeaseDuration)
		if err != nil {
			msg := "default_lease_duration not valid"
			requestContextLogger.Error(msg, "err", err)
			return tools.ErrInvalidRequest(ctx, msg, "err", err)
		}
		// set into the new cloud account the value of the default lease duration
		newCloudaccount.DefaultLeaseDuration = defaultLeaseDuration
	}

	err = c.cs.DB.Create(&newCloudaccount).Error
	if err != nil {
		requestContextLogger.Error("Error while saving new cloudaccount", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	firstOwner := models.Owner{
		Email:          account.Email,
		CloudaccountID: newCloudaccount.ID,
	}
	err = c.cs.DB.Create(&firstOwner).Error
	if err != nil {
		requestContextLogger.Error("Error while saving new owner (first)", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	// regenerate SQS permissions
	if err := c.cs.RegenerateSQSPermissions(); err != nil {
		requestContextLogger.Error("Error RegenerateSQSPermissions", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	return tools.JSONResponse(ctx, 200, tools.HMI{
		"cloudaccount_id": newCloudaccount.ID,
		"aws_id":          newCloudaccount.AWSID,
		"initial_setup_cloudformation_url": fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-initial-setup.template", account.ID, newCloudaccount.ID),
		"region_setup_cloudformation_url":  fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-region-setup.template", account.ID, newCloudaccount.ID),
	})
}

// Update handles the endpoint used to update the configuration of the cloudaccount.
func (c *CloudaccountController) Update(ctx *app.UpdateCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrNotFound(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return tools.ErrNotFound(ctx, "cloud account not found")
	}

	// parse default_lease_duration
	defaultLeaseDuration, err := time.ParseDuration(ctx.Payload.DefaultLeaseDuration)
	if err != nil {
		requestContextLogger.Error("Error parsing ctx.Payload.DefaultLeaseDuration", "err", err)
		return tools.ErrInvalidRequest(ctx, "default_lease_duration not valid", "err", err)
	}
	// set into the cloud account the value of default_lease_duration
	cloudaccount.DefaultLeaseDuration = defaultLeaseDuration

	// save to DB the updated cloudaccount
	err = c.cs.DB.Save(&cloudaccount).Error
	if err != nil {
		requestContextLogger.Error("Error saving updated cloudaccount", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	return tools.JSONResponse(ctx, 200, tools.HMI{
		"cloudaccount_id":        cloudaccount.ID,
		"aws_id":                 cloudaccount.AWSID,
		"default_lease_duration": cloudaccount.DefaultLeaseDuration.String(),
	})
}

// ListWhitelistedOwners runs the listWhitelistedOwners action.
func (c *CloudaccountController) ListWhitelistedOwners(ctx *app.ListWhitelistedOwnersCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return tools.ErrNotFound(ctx, "cloud account not found")
	}

	// check whether this owner email already exists for this cloudaccount
	var ownerList []models.Owner
	err = c.cs.DB.Table("owners").Where(&models.Owner{CloudaccountID: cloudaccount.ID}).Find(&ownerList).Error
	if err != nil {
		requestContextLogger.Error("Error fetching list of owners", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, "no owners found")
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	return tools.JSONResponse(ctx, 200, ownerList)
}

// AddWhitelistedOwner handles the endpoint used to add an email address (plus optional keyname) to the whitelist of owners
// that can start a lease without having to get an approval from the admin (i.e. account).
func (c *CloudaccountController) AddWhitelistedOwner(ctx *app.AddWhitelistedOwnerCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return tools.ErrNotFound(ctx, "cloud account not found")
	}

	// validate email
	ownerEmail, err := c.cs.DefaultMailer().Client.ValidateEmail(ctx.Payload.Email)
	if err != nil {
		requestContextLogger.Error("Error validating new owner email", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}
	if !ownerEmail.IsValid {
		requestContextLogger.Error("New owner email is not valid")
		return tools.ErrInvalidRequest(ctx, "invalid email")
	}

	// check whether this owner email already exists for this cloudaccount
	var equalOwnerEmailCount int64
	c.cs.DB.Table("owners").Where(&models.Owner{CloudaccountID: cloudaccount.ID, Email: ownerEmail.Address}).Count(&equalOwnerEmailCount)
	if equalOwnerEmailCount != 0 {
		requestContextLogger.Error("This email address is already whitelisted")
		return tools.ErrInvalidRequest(ctx, "owner already exists in whitelist")
	}

	if ctx.Payload.KeyName != nil {
		// check whether this owner keyname already exists for this cloudaccount
		var equalOwnerKeynameCount int64
		c.cs.DB.Table("owners").Where(&models.Owner{CloudaccountID: cloudaccount.ID, KeyName: *ctx.Payload.KeyName}).Count(&equalOwnerKeynameCount)
		if equalOwnerKeynameCount != 0 {
			requestContextLogger.Error("This keyname is already whitelisted")
			return tools.ErrInvalidRequest(ctx, "owner already exists in whitelist")
		}
	}

	// instert the new owner into the db
	newOwner := models.Owner{
		CloudaccountID: cloudaccount.ID,
		Email:          ownerEmail.Address,
	}
	if ctx.Payload.KeyName != nil {
		newOwner.KeyName = *ctx.Payload.KeyName
	}
	err = c.cs.DB.Create(&newOwner).Error
	if err != nil {
		requestContextLogger.Error("Error saving new owner email to DB", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	return tools.JSONResponse(ctx, 200, tools.HMI{
		"message": "Owner added successfully to whitelist",
	})
}

// UpdateWhitelistedOwner runs the updateWhitelistedOwner action.
func (c *CloudaccountController) UpdateWhitelistedOwner(ctx *app.UpdateWhitelistedOwnerCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return tools.ErrNotFound(ctx, "cloud account not found")
	}

	// validate email
	ownerEmail, err := c.cs.DefaultMailer().Client.ValidateEmail(ctx.Payload.Email)
	if err != nil {
		requestContextLogger.Error("Error validating new owner email", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}
	if !ownerEmail.IsValid {
		requestContextLogger.Error("Owner email is not valid")
		return tools.ErrInvalidRequest(ctx, "invalid email")
	}

	// check whether this owner email already exists for this cloudaccount
	var existingOwner models.Owner
	err = c.cs.DB.Table("owners").Where(&models.Owner{CloudaccountID: cloudaccount.ID, Email: ownerEmail.Address}).First(&existingOwner).Error
	if err != nil {
		requestContextLogger.Error("Error fetching owner", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("owner with email %v does not exist", ownerEmail.Address))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	if ctx.Payload.KeyName != nil {
		// check whether this owner keyname already exists for this cloudaccount
		var existingOwnerOfKeyName models.Owner
		err = c.cs.DB.Table("owners").Where(&models.Owner{CloudaccountID: cloudaccount.ID, KeyName: *ctx.Payload.KeyName}).First(&existingOwnerOfKeyName).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			requestContextLogger.Error("Error fetching owner of keyname", "err", err)
			return tools.ErrInternal(ctx, "internal server error")
		}

		if existingOwner.ID != existingOwnerOfKeyName.ID {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("This keyname (%q) is associated with another owner email (%v)", *ctx.Payload.KeyName, existingOwnerOfKeyName.Email))
		}

		existingOwner.KeyName = *ctx.Payload.KeyName
	}

	err = c.cs.DB.Save(&existingOwner).Error
	if err != nil {
		requestContextLogger.Error("Error saving updated whitelisted owner to DB", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	return tools.JSONResponse(ctx, 200, tools.HMI{
		"message": "Successfully updated whitelisted owner",
	})
}

// DeleteWhitelistedOwner runs the deleteWhitelistedOwner action.
func (c *CloudaccountController) DeleteWhitelistedOwner(ctx *app.DeleteWhitelistedOwnerCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return tools.ErrNotFound(ctx, "cloud account not found")
	}

	// validate email
	ownerEmail, err := c.cs.DefaultMailer().Client.ValidateEmail(ctx.Payload.Email)
	if err != nil {
		requestContextLogger.Error("Error validating new owner email", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}
	if !ownerEmail.IsValid {
		requestContextLogger.Error("Owner email is not valid")
		return tools.ErrInvalidRequest(ctx, "invalid email")
	}

	// check whether this owner email already exists for this cloudaccount
	var existingOwner models.Owner
	err = c.cs.DB.Table("owners").Where(&models.Owner{CloudaccountID: cloudaccount.ID, Email: ownerEmail.Address}).First(&existingOwner).Error
	if err != nil {
		requestContextLogger.Error("Error fetching owner", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("owner with email %v does not exist", ownerEmail.Address))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	if account.Email == ownerEmail.Address {
		return tools.ErrInvalidRequest(ctx, fmt.Sprintf("Cannot delete owner associated with the account holder's email address (%v)", ownerEmail.Address))
	}

	err = c.cs.DB.Unscoped().Delete(&existingOwner).Error
	if err != nil {
		requestContextLogger.Error("Error saving updated whitelisted owner to DB", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	return tools.JSONResponse(ctx, 200, tools.HMI{
		"message": "Successfully deleted whitelisted owner",
	})
}

// DownloadInitialSetupTemplate handles the endpoint used to download the Cloudformation
// template to be used to make the initial setup of cecil on an AWS account (a.k.a. cloudaccount).
func (c *CloudaccountController) DownloadInitialSetupTemplate(ctx *app.DownloadInitialSetupTemplateCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return tools.ErrNotFound(ctx, "cloud account not found")
	}

	var compiledTemplate bytes.Buffer // A Buffer needs no initialization.

	tpl, err := template.ParseFiles("./core/go-templates/tenant-aws-initial-setup.template")
	if err != nil {
		requestContextLogger.Error("Error reading tenant-aws-initial-setup.template", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	var values = map[string]interface{}{}
	values["IAMRoleExternalID"] = cloudaccount.ExternalID
	values["CecilAWSID"] = c.cs.AWS.Config.AWS_ACCOUNT_ID

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		requestContextLogger.Error("Error compiling tenant-aws-initial-setup.template with data", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	return ctx.OK(compiledTemplate.Bytes())
}

// DownloadRegionSetupTemplate handles the endpoint used to download the Cloudformation
// template to be used to setup the stack to monitor a region on that region.
func (c *CloudaccountController) DownloadRegionSetupTemplate(ctx *app.DownloadRegionSetupTemplateCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return tools.ErrNotFound(ctx, "cloud account not found")
	}

	var compiledTemplate bytes.Buffer // A Buffer needs no initialization.

	tpl, err := template.ParseFiles("./core/go-templates/tenant-aws-region-setup.template")
	if err != nil {
		requestContextLogger.Error("Error reading tenant-aws-region-setup.template", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	var values = map[string]interface{}{}
	values["CecilAWSID"] = c.cs.AWS.Config.AWS_ACCOUNT_ID
	values["CecilAWSRegion"] = c.cs.AWS.Config.AWS_REGION
	values["SNSTopicName"] = c.cs.AWS.Config.SNSTopicName
	values["SQSQueueName"] = c.cs.AWS.Config.SQSQueueName

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		requestContextLogger.Error("Error compiling tenant-aws-region-setup.template with data", "err", err)
		return tools.ErrInternal(ctx, "internal server error")
	}

	return ctx.OK(compiledTemplate.Bytes())
}

// ListRegions handles the endpoint used to list all regions and their status for a cloudaccount.
func (c *CloudaccountController) ListRegions(ctx *app.ListRegionsCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return tools.ErrNotFound(ctx, "cloud account not found")
	}

	listSubscriptions, listSubscriptionsErrors := c.cs.StatusOfAllRegions(cloudaccount.AWSID)

	requestContextLogger.Info(
		"StatusOfAllRegions()",
		"response", listSubscriptions,
		"errors", listSubscriptionsErrors,
	)

	return tools.JSONResponse(ctx, 200, listSubscriptions)
}

// SubscribeSNSToSQS handles the endpoint used to force-try subscription of Cecil to
// all or selected regions.
func (c *CloudaccountController) SubscribeSNSToSQS(ctx *app.SubscribeSNSToSQSCloudaccountContext) error {
	requestContextLogger := core.NewContextLogger(ctx)

	_, err := core.ValidateToken(ctx)
	if err != nil {
		requestContextLogger.Error("Error validating token", "err", err)
		return tools.ErrUnauthorized(ctx, tools.ErrorUnauthorized)
	}

	account, err := c.cs.FetchAccountByID(ctx.AccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching account", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	cloudaccount, err := c.cs.FetchCloudaccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return tools.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		}
		return tools.ErrInternal(ctx, "internal server error")
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudaccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
		return tools.ErrNotFound(ctx, "cloud account not found")
	}

	// TODO: what to do with non-existing regions???
	var regionsToTrySubscription = []string{}

	// check whether the payload specifies to try subscribing to all regions
	if tools.SliceContains(ctx.Payload.Regions, "all") {
		regionsToTrySubscription = core.Regions
	} else {
		regionsToTrySubscription = ctx.Payload.Regions
	}
	createdSubscriptions, createdSubscriptionsErrors := c.cs.SubscribeToRegions(regionsToTrySubscription, cloudaccount.AWSID)

	requestContextLogger.Info(
		"SubscribeToRegions()",
		"response", createdSubscriptions,
		"errors", createdSubscriptionsErrors,
	)

	return tools.JSONResponse(ctx, 200, createdSubscriptions)
}
