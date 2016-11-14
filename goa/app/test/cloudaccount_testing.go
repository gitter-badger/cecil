//************************************************************************//
// API "Cecil": cloudaccount TestHelpers
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/zerocloud/design
// --out=$(GOPATH)/src/github.com/tleyden/zerocloud/goa
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package test

import (
	"bytes"
	"fmt"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/goatest"
	"github.com/tleyden/zerocloud/goa/app"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// AddCloudaccountOK runs the method Add of the given controller with the given parameters and payload.
// It returns the response writer so it's possible to inspect the response headers.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func AddCloudaccountOK(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.CloudaccountController, accountID int, payload *app.AddCloudaccountPayload) http.ResponseWriter {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Validate payload
	err := payload.Validate()
	if err != nil {
		e, ok := err.(goa.ServiceError)
		if !ok {
			panic(err) // bug
		}
		if e.ResponseStatus() != 200 {
			t.Errorf("unexpected payload validation error: %+v", e)
		}
		return nil
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/accounts/%v/cloudaccounts", accountID),
	}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}
	prms := url.Values{}
	prms["account_id"] = []string{fmt.Sprintf("%v", accountID)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "CloudaccountTest"), rw, req, prms)
	addCtx, err := app.NewAddCloudaccountContext(goaCtx, service)
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}
	addCtx.Payload = payload

	// Perform action
	err = ctrl.Add(addCtx)

	// Validate response
	if err != nil {
		t.Fatalf("controller returned %s, logs:\n%s", err, logBuf.String())
	}
	if rw.Code != 200 {
		t.Errorf("invalid response status code: got %+v, expected 200", rw.Code)
	}

	// Return results
	return rw
}

// AddEmailToWhitelistCloudaccountOK runs the method AddEmailToWhitelist of the given controller with the given parameters and payload.
// It returns the response writer so it's possible to inspect the response headers.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func AddEmailToWhitelistCloudaccountOK(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.CloudaccountController, accountID int, cloudaccountID int, payload *app.AddEmailToWhitelistCloudaccountPayload) http.ResponseWriter {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Validate payload
	err := payload.Validate()
	if err != nil {
		e, ok := err.(goa.ServiceError)
		if !ok {
			panic(err) // bug
		}
		if e.ResponseStatus() != 200 {
			t.Errorf("unexpected payload validation error: %+v", e)
		}
		return nil
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/accounts/%v/cloudaccounts/%v/owners", accountID, cloudaccountID),
	}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}
	prms := url.Values{}
	prms["account_id"] = []string{fmt.Sprintf("%v", accountID)}
	prms["cloudaccount_id"] = []string{fmt.Sprintf("%v", cloudaccountID)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "CloudaccountTest"), rw, req, prms)
	addEmailToWhitelistCtx, err := app.NewAddEmailToWhitelistCloudaccountContext(goaCtx, service)
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}
	addEmailToWhitelistCtx.Payload = payload

	// Perform action
	err = ctrl.AddEmailToWhitelist(addEmailToWhitelistCtx)

	// Validate response
	if err != nil {
		t.Fatalf("controller returned %s, logs:\n%s", err, logBuf.String())
	}
	if rw.Code != 200 {
		t.Errorf("invalid response status code: got %+v, expected 200", rw.Code)
	}

	// Return results
	return rw
}

// DownloadInitialSetupTemplateCloudaccountOK runs the method DownloadInitialSetupTemplate of the given controller with the given parameters.
// It returns the response writer so it's possible to inspect the response headers.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func DownloadInitialSetupTemplateCloudaccountOK(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.CloudaccountController, accountID int, cloudaccountID int) http.ResponseWriter {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-initial-setup.template", accountID, cloudaccountID),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}
	prms := url.Values{}
	prms["account_id"] = []string{fmt.Sprintf("%v", accountID)}
	prms["cloudaccount_id"] = []string{fmt.Sprintf("%v", cloudaccountID)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "CloudaccountTest"), rw, req, prms)
	downloadInitialSetupTemplateCtx, err := app.NewDownloadInitialSetupTemplateCloudaccountContext(goaCtx, service)
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}

	// Perform action
	err = ctrl.DownloadInitialSetupTemplate(downloadInitialSetupTemplateCtx)

	// Validate response
	if err != nil {
		t.Fatalf("controller returned %s, logs:\n%s", err, logBuf.String())
	}
	if rw.Code != 200 {
		t.Errorf("invalid response status code: got %+v, expected 200", rw.Code)
	}

	// Return results
	return rw
}

// DownloadRegionSetupTemplateCloudaccountOK runs the method DownloadRegionSetupTemplate of the given controller with the given parameters.
// It returns the response writer so it's possible to inspect the response headers.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func DownloadRegionSetupTemplateCloudaccountOK(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.CloudaccountController, accountID int, cloudaccountID int) http.ResponseWriter {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-region-setup.template", accountID, cloudaccountID),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}
	prms := url.Values{}
	prms["account_id"] = []string{fmt.Sprintf("%v", accountID)}
	prms["cloudaccount_id"] = []string{fmt.Sprintf("%v", cloudaccountID)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "CloudaccountTest"), rw, req, prms)
	downloadRegionSetupTemplateCtx, err := app.NewDownloadRegionSetupTemplateCloudaccountContext(goaCtx, service)
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}

	// Perform action
	err = ctrl.DownloadRegionSetupTemplate(downloadRegionSetupTemplateCtx)

	// Validate response
	if err != nil {
		t.Fatalf("controller returned %s, logs:\n%s", err, logBuf.String())
	}
	if rw.Code != 200 {
		t.Errorf("invalid response status code: got %+v, expected 200", rw.Code)
	}

	// Return results
	return rw
}
