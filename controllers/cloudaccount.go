package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/tleyden/cecil/core"
	"github.com/tleyden/cecil/goa/app"
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
			return core.ErrNotFound(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// TODO: validate newCloudAccountInput.AWSID

	AWSIDAlreadyRegistered, err := c.cs.CloudAccountByAWSIDExists(ctx.Payload.AwsID)
	if err != nil {
		requestContextLogger.Error("Error CloudAccountByAWSIDExists", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}
	if AWSIDAlreadyRegistered {
		requestContextLogger.Error("AWSIDAlreadyRegistered", "err", err)
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

	// check whether the payload contains default_lease_duration
	if ctx.Payload.DefaultLeaseDuration != nil {
		defaultLeaseDuration, err := time.ParseDuration(*ctx.Payload.DefaultLeaseDuration)
		if err != nil {
			msg := "default_lease_duration not valid"
			requestContextLogger.Error(msg, "err", err)
			return core.ErrInvalidRequest(ctx, msg, "err", err)
		}
		// set into the new cloud account the value of the default lease duration
		newCloudAccount.DefaultLeaseDuration = defaultLeaseDuration
	}

	err = c.cs.DB.Create(&newCloudAccount).Error
	if err != nil {
		requestContextLogger.Error("Error while saving new cloudaccount", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	firstOwner := core.Owner{
		Email:          account.Email,
		CloudAccountID: newCloudAccount.ID,
	}
	err = c.cs.DB.Create(&firstOwner).Error
	if err != nil {
		requestContextLogger.Error("Error while saving new owner (first)", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	// regenerate SQS permissions
	if err := c.cs.RegenerateSQSPermissions(); err != nil {
		requestContextLogger.Error("Error RegenerateSQSPermissions", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	return core.JSONResponse(ctx, 200, gin.H{
		"cloudaccount_id": newCloudAccount.ID,
		"aws_id":          newCloudAccount.AWSID,
		"initial_setup_cloudformation_url": fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-initial-setup.template", account.ID, newCloudAccount.ID),
		"region_setup_cloudformation_url":  fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-region-setup.template", account.ID, newCloudAccount.ID),
	})
}

// Update handles the endpoint used to update the configuration of the cloudaccount.
func (c *CloudaccountController) Update(ctx *app.UpdateCloudaccountContext) error {
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
			return core.ErrNotFound(ctx, fmt.Sprintf("account with id %v does not exist", ctx.AccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	cloudAccount, err := c.cs.FetchCloudAccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudAccount.ID))
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	// parse default_lease_duration
	defaultLeaseDuration, err := time.ParseDuration(ctx.Payload.DefaultLeaseDuration)
	if err != nil {
		requestContextLogger.Error("Error parsing ctx.Payload.DefaultLeaseDuration", "err", err)
		return core.ErrInvalidRequest(ctx, "default_lease_duration not valid", "err", err)
	}
	// set into the cloud account the value of default_lease_duration
	cloudAccount.DefaultLeaseDuration = defaultLeaseDuration

	// save to DB the updated cloudAccount
	err = c.cs.DB.Save(&cloudAccount).Error
	if err != nil {
		requestContextLogger.Error("Error saving updated cloudaccount", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	return core.JSONResponse(ctx, 200, gin.H{
		"cloudaccount_id":        cloudAccount.ID,
		"aws_id":                 cloudAccount.AWSID,
		"default_lease_duration": cloudAccount.DefaultLeaseDuration.String(),
	})
}

// AddEmailToWhitelist handles the endpoint used to add an email address to the whitelist of owners
// that can start a lease without having to get an approval from the admin (i.e. account).
func (c *CloudaccountController) AddEmailToWhitelist(ctx *app.AddEmailToWhitelistCloudaccountContext) error {
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
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	cloudAccount, err := c.cs.FetchCloudAccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudAccount.ID))
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	// validate email
	ownerEmail, err := c.cs.DefaultMailer.Client.ValidateEmail(ctx.Payload.Email)
	if err != nil {
		requestContextLogger.Error("Error validating new owner email", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}
	if !ownerEmail.IsValid {
		requestContextLogger.Error("New owner email is not valid")
		return core.ErrInvalidRequest(ctx, "invalid email")
	}

	// check whether this owner already exists for this cloudaccount
	var equalOwnerCount int64
	c.cs.DB.Table("owners").Where(&core.Owner{CloudAccountID: cloudAccount.ID, Email: ownerEmail.Address}).Count(&equalOwnerCount)
	if equalOwnerCount != 0 {
		requestContextLogger.Error("New owner email is already registered")
		return core.ErrInvalidRequest(ctx, "owner already exists in whitelist")
	}

	// instert the new owner into the db
	newOwner := core.Owner{
		CloudAccountID: cloudAccount.ID,
		Email:          ownerEmail.Address,
	}
	err = c.cs.DB.Create(&newOwner).Error
	if err != nil {
		requestContextLogger.Error("Error saving new owner email to DB", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	return core.JSONResponse(ctx, 200, gin.H{
		"message": "Owner added successfully to whitelist",
	})
	return nil
}

// DownloadInitialSetupTemplate handles the endpoint used to download the Cloudformation
// template to be used to make the initial setup of cecil on an AWS account (a.k.a. cloudaccount).
func (c *CloudaccountController) DownloadInitialSetupTemplate(ctx *app.DownloadInitialSetupTemplateCloudaccountContext) error {
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
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	cloudAccount, err := c.cs.FetchCloudAccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudAccount.ID))
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	var compiledTemplate bytes.Buffer // A Buffer needs no initialization.

	tpl, err := template.ParseFiles("./core/go-templates/tenant-aws-initial-setup.template")
	if err != nil {
		requestContextLogger.Error("Error reading tenant-aws-initial-setup.template", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	var values map[string]interface{} = map[string]interface{}{}
	values["IAMRoleExternalID"] = cloudAccount.ExternalID
	values["CecilAWSID"] = c.cs.AWS.Config.AWS_ACCOUNT_ID

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		requestContextLogger.Error("Error compiling tenant-aws-initial-setup.template with data", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	return ctx.OK(compiledTemplate.Bytes())
}

// DownloadRegionSetupTemplate handles the endpoint used to download the Cloudformation
// template to be used to setup the stack to monitor a region on that region.
func (c *CloudaccountController) DownloadRegionSetupTemplate(ctx *app.DownloadRegionSetupTemplateCloudaccountContext) error {
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
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	cloudAccount, err := c.cs.FetchCloudAccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudAccount.ID))
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	var compiledTemplate bytes.Buffer // A Buffer needs no initialization.

	tpl, err := template.ParseFiles("./core/go-templates/tenant-aws-region-setup.template")
	if err != nil {
		requestContextLogger.Error("Error reading tenant-aws-region-setup.template", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	var values map[string]interface{} = map[string]interface{}{}
	values["CecilAWSID"] = c.cs.AWS.Config.AWS_ACCOUNT_ID
	values["CecilAWSRegion"] = c.cs.AWS.Config.AWS_REGION
	values["SNSTopicName"] = c.cs.AWS.Config.SNSTopicName
	values["SQSQueueName"] = c.cs.AWS.Config.SQSQueueName

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		requestContextLogger.Error("Error compiling tenant-aws-region-setup.template with data", "err", err)
		return core.ErrInternal(ctx, "internal server error")
	}

	return ctx.OK(compiledTemplate.Bytes())
}

// ListRegions handles the endpoint used to list all regions and their status for a cloudaccount.
func (c *CloudaccountController) ListRegions(ctx *app.ListRegionsCloudaccountContext) error {
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
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	cloudAccount, err := c.cs.FetchCloudAccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudAccount.ID))
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	listSubscriptions, listSubscriptionsErrors := c.cs.StatusOfAllRegions(cloudAccount.AWSID)

	requestContextLogger.Info(
		"StatusOfAllRegions()",
		"response", listSubscriptions,
		"errors", listSubscriptionsErrors,
	)

	return core.JSONResponse(ctx, 200, listSubscriptions)
}

// SubscribeSNSToSQS handles the endpoint used to force-try subscription of Cecil to
// all or selected regions.
func (c *CloudaccountController) SubscribeSNSToSQS(ctx *app.SubscribeSNSToSQSCloudaccountContext) error {
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
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	cloudAccount, err := c.cs.FetchCloudAccountByID(ctx.CloudaccountID)
	if err != nil {
		requestContextLogger.Error("Error fetching cloudaccount", "err", err)
		if err == gorm.ErrRecordNotFound {
			return core.ErrInvalidRequest(ctx, fmt.Sprintf("cloud account with id %v does not exist", ctx.CloudaccountID))
		} else {
			return core.ErrInternal(ctx, "internal server error")
		}
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudAccount.ID))
		return core.ErrNotFound(ctx, "cloud account not found")
	}

	// TODO: what to do with non-existing regions???
	var regionsToTrySubscription []string = []string{}

	// check whether the payload specifies to try subscribing to all regions
	if core.SliceContains(ctx.Payload.Regions, "all") {
		regionsToTrySubscription = core.Regions
	} else {
		regionsToTrySubscription = ctx.Payload.Regions
	}
	createdSubscriptions, createdSubscriptionsErrors := c.cs.SubscribeToRegions(regionsToTrySubscription, cloudAccount.AWSID)

	requestContextLogger.Info(
		"SubscribeToRegions()",
		"response", createdSubscriptions,
		"errors", createdSubscriptionsErrors,
	)

	return core.JSONResponse(ctx, 200, createdSubscriptions)
}
