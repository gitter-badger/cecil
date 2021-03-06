// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "Cecil": Application Controllers
//
// Command:
// $ goagen
// --design=github.com/tleyden/cecil/design
// --out=$(GOPATH)/src/github.com/tleyden/cecil/goa
// --version=v1.0.0

package app

import (
	"context"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/cors"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// AccountController is the controller interface for the Account actions.
type AccountController interface {
	goa.Muxer
	Create(*CreateAccountContext) error
	MailerConfig(*MailerConfigAccountContext) error
	NewAPIToken(*NewAPITokenAccountContext) error
	RemoveMailer(*RemoveMailerAccountContext) error
	RemoveSlack(*RemoveSlackAccountContext) error
	Show(*ShowAccountContext) error
	SlackConfig(*SlackConfigAccountContext) error
	Verify(*VerifyAccountContext) error
}

// MountAccountController "mounts" a Account resource controller on the given service.
func MountAccountController(service *goa.Service, ctrl AccountController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateAccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CreateAccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/accounts", ctrl.MuxHandler("create", h, unmarshalCreateAccountPayload))
	service.LogInfo("mount", "ctrl", "Account", "action", "Create", "route", "POST /accounts")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewMailerConfigAccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*MailerConfigAccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.MailerConfig(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/mailer_config", ctrl.MuxHandler("mailerConfig", h, unmarshalMailerConfigAccountPayload))
	service.LogInfo("mount", "ctrl", "Account", "action", "MailerConfig", "route", "POST /accounts/:account_id/mailer_config", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewNewAPITokenAccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*NewAPITokenAccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.NewAPIToken(rctx)
	}
	service.Mux.Handle("POST", "/accounts/:account_id/new_api_token", ctrl.MuxHandler("new_api_token", h, unmarshalNewAPITokenAccountPayload))
	service.LogInfo("mount", "ctrl", "Account", "action", "NewAPIToken", "route", "POST /accounts/:account_id/new_api_token")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewRemoveMailerAccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.RemoveMailer(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("DELETE", "/accounts/:account_id/mailer_config", ctrl.MuxHandler("removeMailer", h, nil))
	service.LogInfo("mount", "ctrl", "Account", "action", "RemoveMailer", "route", "DELETE /accounts/:account_id/mailer_config", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewRemoveSlackAccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.RemoveSlack(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("DELETE", "/accounts/:account_id/slack_config", ctrl.MuxHandler("removeSlack", h, nil))
	service.LogInfo("mount", "ctrl", "Account", "action", "RemoveSlack", "route", "DELETE /accounts/:account_id/slack_config", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowAccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id", ctrl.MuxHandler("show", h, nil))
	service.LogInfo("mount", "ctrl", "Account", "action", "Show", "route", "GET /accounts/:account_id", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewSlackConfigAccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*SlackConfigAccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.SlackConfig(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/slack_config", ctrl.MuxHandler("slackConfig", h, unmarshalSlackConfigAccountPayload))
	service.LogInfo("mount", "ctrl", "Account", "action", "SlackConfig", "route", "POST /accounts/:account_id/slack_config", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewVerifyAccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*VerifyAccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Verify(rctx)
	}
	service.Mux.Handle("POST", "/accounts/:account_id/api_token", ctrl.MuxHandler("verify", h, unmarshalVerifyAccountPayload))
	service.LogInfo("mount", "ctrl", "Account", "action", "Verify", "route", "POST /accounts/:account_id/api_token")
}

// unmarshalCreateAccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateAccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &createAccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalMailerConfigAccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalMailerConfigAccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &mailerConfigAccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalNewAPITokenAccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalNewAPITokenAccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &newAPITokenAccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalSlackConfigAccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalSlackConfigAccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &slackConfigAccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalVerifyAccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalVerifyAccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &verifyAccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// CloudaccountController is the controller interface for the Cloudaccount actions.
type CloudaccountController interface {
	goa.Muxer
	Add(*AddCloudaccountContext) error
	AddWhitelistedOwner(*AddWhitelistedOwnerCloudaccountContext) error
	DeleteWhitelistedOwner(*DeleteWhitelistedOwnerCloudaccountContext) error
	DownloadInitialSetupTemplate(*DownloadInitialSetupTemplateCloudaccountContext) error
	DownloadRegionSetupTemplate(*DownloadRegionSetupTemplateCloudaccountContext) error
	ListRegions(*ListRegionsCloudaccountContext) error
	ListWhitelistedOwners(*ListWhitelistedOwnersCloudaccountContext) error
	Show(*ShowCloudaccountContext) error
	SubscribeSNSToSQS(*SubscribeSNSToSQSCloudaccountContext) error
	Update(*UpdateCloudaccountContext) error
	UpdateWhitelistedOwner(*UpdateWhitelistedOwnerCloudaccountContext) error
}

// MountCloudaccountController "mounts" a Cloudaccount resource controller on the given service.
func MountCloudaccountController(service *goa.Service, ctrl CloudaccountController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewAddCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*AddCloudaccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Add(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/cloudaccounts", ctrl.MuxHandler("add", h, unmarshalAddCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "Add", "route", "POST /accounts/:account_id/cloudaccounts", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewAddWhitelistedOwnerCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*AddWhitelistedOwnerCloudaccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.AddWhitelistedOwner(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", ctrl.MuxHandler("addWhitelistedOwner", h, unmarshalAddWhitelistedOwnerCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "AddWhitelistedOwner", "route", "POST /accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDeleteWhitelistedOwnerCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*DeleteWhitelistedOwnerCloudaccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.DeleteWhitelistedOwner(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("DELETE", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", ctrl.MuxHandler("deleteWhitelistedOwner", h, unmarshalDeleteWhitelistedOwnerCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "DeleteWhitelistedOwner", "route", "DELETE /accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDownloadInitialSetupTemplateCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.DownloadInitialSetupTemplate(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/tenant-aws-initial-setup.template", ctrl.MuxHandler("downloadInitialSetupTemplate", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "DownloadInitialSetupTemplate", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/tenant-aws-initial-setup.template", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDownloadRegionSetupTemplateCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.DownloadRegionSetupTemplate(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/tenant-aws-region-setup.template", ctrl.MuxHandler("downloadRegionSetupTemplate", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "DownloadRegionSetupTemplate", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/tenant-aws-region-setup.template", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListRegionsCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.ListRegions(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/regions", ctrl.MuxHandler("listRegions", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "ListRegions", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/regions", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListWhitelistedOwnersCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.ListWhitelistedOwners(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", ctrl.MuxHandler("listWhitelistedOwners", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "ListWhitelistedOwners", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id", ctrl.MuxHandler("show", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "Show", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewSubscribeSNSToSQSCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*SubscribeSNSToSQSCloudaccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.SubscribeSNSToSQS(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/subscribe-sns-to-sqs", ctrl.MuxHandler("subscribeSNSToSQS", h, unmarshalSubscribeSNSToSQSCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "SubscribeSNSToSQS", "route", "POST /accounts/:account_id/cloudaccounts/:cloudaccount_id/subscribe-sns-to-sqs", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewUpdateCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*UpdateCloudaccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Update(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("PATCH", "/accounts/:account_id/cloudaccounts/:cloudaccount_id", ctrl.MuxHandler("update", h, unmarshalUpdateCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "Update", "route", "PATCH /accounts/:account_id/cloudaccounts/:cloudaccount_id", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewUpdateWhitelistedOwnerCloudaccountContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*UpdateWhitelistedOwnerCloudaccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.UpdateWhitelistedOwner(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("PATCH", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", ctrl.MuxHandler("updateWhitelistedOwner", h, unmarshalUpdateWhitelistedOwnerCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "UpdateWhitelistedOwner", "route", "PATCH /accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", "security", "jwt")
}

// unmarshalAddCloudaccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalAddCloudaccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &addCloudaccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalAddWhitelistedOwnerCloudaccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalAddWhitelistedOwnerCloudaccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &addWhitelistedOwnerCloudaccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalDeleteWhitelistedOwnerCloudaccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalDeleteWhitelistedOwnerCloudaccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &deleteWhitelistedOwnerCloudaccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalSubscribeSNSToSQSCloudaccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalSubscribeSNSToSQSCloudaccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &subscribeSNSToSQSCloudaccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalUpdateCloudaccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalUpdateCloudaccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &updateCloudaccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalUpdateWhitelistedOwnerCloudaccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalUpdateWhitelistedOwnerCloudaccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &updateWhitelistedOwnerCloudaccountPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// EmailActionController is the controller interface for the EmailAction actions.
type EmailActionController interface {
	goa.Muxer
	Actions(*ActionsEmailActionContext) error
}

// MountEmailActionController "mounts" a EmailAction resource controller on the given service.
func MountEmailActionController(service *goa.Service, ctrl EmailActionController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewActionsEmailActionContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Actions(rctx)
	}
	service.Mux.Handle("GET", "/email_action/leases/:lease_uuid/:group_uid_hash/:action", ctrl.MuxHandler("actions", h, nil))
	service.LogInfo("mount", "ctrl", "EmailAction", "action", "Actions", "route", "GET /email_action/leases/:lease_uuid/:group_uid_hash/:action")
}

// LeasesController is the controller interface for the Leases actions.
type LeasesController interface {
	goa.Muxer
	DeleteFromDB(*DeleteFromDBLeasesContext) error
	ListLeasesForAccount(*ListLeasesForAccountLeasesContext) error
	ListLeasesForCloudaccount(*ListLeasesForCloudaccountLeasesContext) error
	SetExpiry(*SetExpiryLeasesContext) error
	Show(*ShowLeasesContext) error
	Terminate(*TerminateLeasesContext) error
}

// MountLeasesController "mounts" a Leases resource controller on the given service.
func MountLeasesController(service *goa.Service, ctrl LeasesController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDeleteFromDBLeasesContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.DeleteFromDB(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id/delete", ctrl.MuxHandler("deleteFromDB", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "DeleteFromDB", "route", "POST /accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id/delete", "security", "jwt")
	service.Mux.Handle("POST", "/accounts/:account_id/leases/:lease_id/delete", ctrl.MuxHandler("deleteFromDB", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "DeleteFromDB", "route", "POST /accounts/:account_id/leases/:lease_id/delete", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListLeasesForAccountLeasesContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.ListLeasesForAccount(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/leases", ctrl.MuxHandler("listLeasesForAccount", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "ListLeasesForAccount", "route", "GET /accounts/:account_id/leases", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListLeasesForCloudaccountLeasesContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.ListLeasesForCloudaccount(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases", ctrl.MuxHandler("listLeasesForCloudaccount", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "ListLeasesForCloudaccount", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/leases", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewSetExpiryLeasesContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.SetExpiry(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id/expiry", ctrl.MuxHandler("setExpiry", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "SetExpiry", "route", "POST /accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id/expiry", "security", "jwt")
	service.Mux.Handle("POST", "/accounts/:account_id/leases/:lease_id/expiry", ctrl.MuxHandler("setExpiry", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "SetExpiry", "route", "POST /accounts/:account_id/leases/:lease_id/expiry", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowLeasesContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id", ctrl.MuxHandler("show", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "Show", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id", "security", "jwt")
	service.Mux.Handle("GET", "/accounts/:account_id/leases/:lease_id", ctrl.MuxHandler("show", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "Show", "route", "GET /accounts/:account_id/leases/:lease_id", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewTerminateLeasesContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Terminate(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id/terminate", ctrl.MuxHandler("terminate", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "Terminate", "route", "POST /accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id/terminate", "security", "jwt")
	service.Mux.Handle("POST", "/accounts/:account_id/leases/:lease_id/terminate", ctrl.MuxHandler("terminate", h, nil))
	service.LogInfo("mount", "ctrl", "Leases", "action", "Terminate", "route", "POST /accounts/:account_id/leases/:lease_id/terminate", "security", "jwt")
}

// ReportController is the controller interface for the Report actions.
type ReportController interface {
	goa.Muxer
	OrderInstancesReport(*OrderInstancesReportReportContext) error
	ShowReport(*ShowReportReportContext) error
}

// MountReportController "mounts" a Report resource controller on the given service.
func MountReportController(service *goa.Service, ctrl ReportController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewOrderInstancesReportReportContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*OrderInstancesReportReportPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.OrderInstancesReport(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/reports/instances", ctrl.MuxHandler("orderInstancesReport", h, unmarshalOrderInstancesReportReportPayload))
	service.LogInfo("mount", "ctrl", "Report", "action", "OrderInstancesReport", "route", "POST /accounts/:account_id/cloudaccounts/:cloudaccount_id/reports/instances", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowReportReportContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.ShowReport(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/reports/generated/:report_uuid", ctrl.MuxHandler("showReport", h, nil))
	service.LogInfo("mount", "ctrl", "Report", "action", "ShowReport", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/reports/generated/:report_uuid", "security", "jwt")
}

// unmarshalOrderInstancesReportReportPayload unmarshals the request body into the context request data Payload field.
func unmarshalOrderInstancesReportReportPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &orderInstancesReportReportPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// RootController is the controller interface for the Root actions.
type RootController interface {
	goa.Muxer
	Show(*ShowRootContext) error
}

// MountRootController "mounts" a Root resource controller on the given service.
func MountRootController(service *goa.Service, ctrl RootController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowRootContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	service.Mux.Handle("GET", "/", ctrl.MuxHandler("show", h, nil))
	service.LogInfo("mount", "ctrl", "Root", "action", "Show", "route", "GET /")
}

// SwaggerController is the controller interface for the Swagger actions.
type SwaggerController interface {
	goa.Muxer
	goa.FileServer
}

// MountSwaggerController "mounts" a Swagger resource controller on the given service.
func MountSwaggerController(service *goa.Service, ctrl SwaggerController) {
	initService(service)
	var h goa.Handler
	service.Mux.Handle("OPTIONS", "/swagger.json", ctrl.MuxHandler("preflight", handleSwaggerOrigin(cors.HandlePreflight()), nil))

	h = ctrl.FileHandler("/swagger.json", "goa/swagger/swagger.json")
	h = handleSwaggerOrigin(h)
	service.Mux.Handle("GET", "/swagger.json", ctrl.MuxHandler("serve", h, nil))
	service.LogInfo("mount", "ctrl", "Swagger", "files", "goa/swagger/swagger.json", "route", "GET /swagger.json")
}

// handleSwaggerOrigin applies the CORS response headers corresponding to the origin.
func handleSwaggerOrigin(h goa.Handler) goa.Handler {

	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			return h(ctx, rw, req)
		}
		if cors.MatchOrigin(origin, "*") {
			ctx = goa.WithLogContext(ctx, "origin", origin)
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := req.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				rw.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			}
			return h(ctx, rw, req)
		}

		return h(ctx, rw, req)
	}
}
