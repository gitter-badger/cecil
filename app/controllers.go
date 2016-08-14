//************************************************************************//
// API "zerocloud": Application Controllers
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/zerocloud/design
// --out=$(GOPATH)/src/github.com/tleyden/zerocloud
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package app

import (
	"github.com/goadesign/goa"
	"golang.org/x/net/context"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Encoder.Register(goa.NewGobEncoder, "application/gob", "application/x-gob")
	service.Encoder.Register(goa.NewXMLEncoder, "application/xml")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")
	service.Decoder.Register(goa.NewGobDecoder, "application/gob", "application/x-gob")
	service.Decoder.Register(goa.NewXMLDecoder, "application/xml")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// AccountController is the controller interface for the Account actions.
type AccountController interface {
	goa.Muxer
	Create(*CreateAccountContext) error
	Delete(*DeleteAccountContext) error
	List(*ListAccountContext) error
	Show(*ShowAccountContext) error
	Update(*UpdateAccountContext) error
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
			rctx.Payload = rawPayload.(*AccountPayload)
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
		rctx, err := NewDeleteAccountContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.Delete(rctx)
	}
	service.Mux.Handle("DELETE", "/accounts/:accountID", ctrl.MuxHandler("Delete", h, nil))
	service.LogInfo("mount", "ctrl", "Account", "action", "Delete", "route", "DELETE /accounts/:accountID")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListAccountContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.List(rctx)
	}
	service.Mux.Handle("GET", "/accounts", ctrl.MuxHandler("List", h, nil))
	service.LogInfo("mount", "ctrl", "Account", "action", "List", "route", "GET /accounts")

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
	service.Mux.Handle("GET", "/accounts/:accountID", ctrl.MuxHandler("Show", h, nil))
	service.LogInfo("mount", "ctrl", "Account", "action", "Show", "route", "GET /accounts/:accountID")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewUpdateAccountContext(ctx, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*UpdateAccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Update(rctx)
	}
	service.Mux.Handle("PUT", "/accounts/:accountID", ctrl.MuxHandler("Update", h, unmarshalUpdateAccountPayload))
	service.LogInfo("mount", "ctrl", "Account", "action", "Update", "route", "PUT /accounts/:accountID")
}

// unmarshalCreateAccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateAccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &accountPayload{}
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

// unmarshalUpdateAccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalUpdateAccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &updateAccountPayload{}
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

// AwsController is the controller interface for the Aws actions.
type AwsController interface {
	goa.Muxer
	Show(*ShowAwsContext) error
}

// MountAwsController "mounts" a Aws resource controller on the given service.
func MountAwsController(service *goa.Service, ctrl AwsController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowAwsContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	service.Mux.Handle("GET", "/aws/:awsAccountID", ctrl.MuxHandler("Show", h, nil))
	service.LogInfo("mount", "ctrl", "Aws", "action", "Show", "route", "GET /aws/:awsAccountID")
}

// CloudaccountController is the controller interface for the Cloudaccount actions.
type CloudaccountController interface {
	goa.Muxer
	Create(*CreateCloudaccountContext) error
	Delete(*DeleteCloudaccountContext) error
	List(*ListCloudaccountContext) error
	Show(*ShowCloudaccountContext) error
	Update(*UpdateCloudaccountContext) error
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
		rctx, err := NewCreateCloudaccountContext(ctx, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CloudAccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/accounts/:accountID/cloudaccounts", ctrl.MuxHandler("Create", h, unmarshalCreateCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "Create", "route", "POST /accounts/:accountID/cloudaccounts")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDeleteCloudaccountContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.Delete(rctx)
	}
	service.Mux.Handle("DELETE", "/accounts/:accountID/cloudaccounts/:cloudAccountID", ctrl.MuxHandler("Delete", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "Delete", "route", "DELETE /accounts/:accountID/cloudaccounts/:cloudAccountID")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListCloudaccountContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.List(rctx)
	}
	service.Mux.Handle("GET", "/accounts/:accountID/cloudaccounts", ctrl.MuxHandler("List", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "List", "route", "GET /accounts/:accountID/cloudaccounts")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowCloudaccountContext(ctx, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	service.Mux.Handle("GET", "/accounts/:accountID/cloudaccounts/:cloudAccountID", ctrl.MuxHandler("Show", h, nil))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "Show", "route", "GET /accounts/:accountID/cloudaccounts/:cloudAccountID")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewUpdateCloudaccountContext(ctx, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CloudAccountPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Update(rctx)
	}
	service.Mux.Handle("PATCH", "/accounts/:accountID/cloudaccounts/:cloudAccountID", ctrl.MuxHandler("Update", h, unmarshalUpdateCloudaccountPayload))
	service.LogInfo("mount", "ctrl", "Cloudaccount", "action", "Update", "route", "PATCH /accounts/:accountID/cloudaccounts/:cloudAccountID")
}

// unmarshalCreateCloudaccountPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateCloudaccountPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &cloudAccountPayload{}
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
	payload := &cloudAccountPayload{}
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

// CloudeventController is the controller interface for the Cloudevent actions.
type CloudeventController interface {
	goa.Muxer
	Create(*CreateCloudeventContext) error
}

// MountCloudeventController "mounts" a Cloudevent resource controller on the given service.
func MountCloudeventController(service *goa.Service, ctrl CloudeventController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateCloudeventContext(ctx, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CloudEventPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/cloudevent", ctrl.MuxHandler("Create", h, unmarshalCreateCloudeventPayload))
	service.LogInfo("mount", "ctrl", "Cloudevent", "action", "Create", "route", "POST /cloudevent")
}

// unmarshalCreateCloudeventPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateCloudeventPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &cloudEventPayload{}
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
