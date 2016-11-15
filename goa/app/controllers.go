//************************************************************************//
// API "Cecil": Application Controllers
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/cecil/design
// --out=$(GOPATH)/src/github.com/tleyden/cecil/goa
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package app

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/cors"
	"golang.org/x/net/context"
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
	Show(*ShowAccountContext) error
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
		rctx, err := NewCreateAccountContext(ctx, service)
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
	service.Mux.Handle("POST", "/accounts", ctrl.MuxHandler("Create", h, unmarshalCreateAccountPayload))
	service.LogInfo("mount", "ctrl", "Account", "action", "Create", "route", "POST /accounts")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowAccountContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id", ctrl.MuxHandler("Show", h, nil))
	service.LogInfo("mount", "ctrl", "Account", "action", "Show", "route", "GET /accounts/:account_id", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewVerifyAccountContext(ctx, service)
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
	service.Mux.Handle("POST", "/accounts/:account_id/api_token", ctrl.MuxHandler("Verify", h, unmarshalVerifyAccountPayload))
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
	AddEmailToWhitelist(*AddEmailToWhitelistCloudaccountContext) error
	DownloadInitialSetupTemplate(*DownloadInitialSetupTemplateCloudaccountContext) error
	DownloadRegionSetupTemplate(*DownloadRegionSetupTemplateCloudaccountContext) error
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
		rctx, err := NewAddCloudaccountContext(ctx, service)
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
	service.Mux.Handle("POST", "/accounts/:account_id/cloudaccounts", ctrl.MuxHandler("Add", h, unmarshalAddCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "Add", "route", "POST /accounts/:account_id/cloudaccounts", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewAddEmailToWhitelistCloudaccountContext(ctx, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*AddEmailToWhitelistCloudaccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.AddEmailToWhitelist(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("POST", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", ctrl.MuxHandler("AddEmailToWhitelist", h, unmarshalAddEmailToWhitelistCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "AddEmailToWhitelist", "route", "POST /accounts/:account_id/cloudaccounts/:cloudaccount_id/owners", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDownloadInitialSetupTemplateCloudaccountContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.DownloadInitialSetupTemplate(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/tenant-aws-initial-setup.template", ctrl.MuxHandler("DownloadInitialSetupTemplate", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "DownloadInitialSetupTemplate", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/tenant-aws-initial-setup.template", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDownloadRegionSetupTemplateCloudaccountContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.DownloadRegionSetupTemplate(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	service.Mux.Handle("GET", "/accounts/:account_id/cloudaccounts/:cloudaccount_id/tenant-aws-region-setup.template", ctrl.MuxHandler("DownloadRegionSetupTemplate", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "DownloadRegionSetupTemplate", "route", "GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/tenant-aws-region-setup.template", "security", "jwt")
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

// unmarshalAddEmailToWhitelistCloudaccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalAddEmailToWhitelistCloudaccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &addEmailToWhitelistCloudaccountPayload{}
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
		rctx, err := NewActionsEmailActionContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.Actions(rctx)
	}
	service.Mux.Handle("GET", "/email_action/leases/:lease_uuid/:instance_id/:action", ctrl.MuxHandler("Actions", h, nil))
	service.LogInfo("mount", "ctrl", "EmailAction", "action", "Actions", "route", "GET /email_action/leases/:lease_uuid/:instance_id/:action")
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

	h = ctrl.FileHandler("/swagger.json", "swagger/swagger.json")
	h = handleSwaggerOrigin(h)
	service.Mux.Handle("GET", "/swagger.json", ctrl.MuxHandler("serve", h, nil))
	service.LogInfo("mount", "ctrl", "Swagger", "files", "swagger/swagger.json", "route", "GET /swagger.json")
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
