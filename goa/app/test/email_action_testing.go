// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "Cecil": email_action TestHelpers
//
// Command:
// $ goagen
// --design=github.com/tleyden/cecil/design
// --out=$(GOPATH)/src/github.com/tleyden/cecil/goa
// --version=v1.0.0

package test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/goatest"
	uuid "github.com/satori/go.uuid"
	"github.com/tleyden/cecil/goa/app"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// ActionsEmailActionOK runs the method Actions of the given controller with the given parameters.
// It returns the response writer so it's possible to inspect the response headers.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func ActionsEmailActionOK(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.EmailActionController, leaseUUID uuid.UUID, groupUIDHash string, action string, sig string, tok string) http.ResponseWriter {
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
	query := url.Values{}
	{
		sliceVal := []string{sig}
		query["sig"] = sliceVal
	}
	{
		sliceVal := []string{tok}
		query["tok"] = sliceVal
	}
	u := &url.URL{
		Path:     fmt.Sprintf("/email_action/leases/%v/%v/%v", leaseUUID, groupUIDHash, action),
		RawQuery: query.Encode(),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}
	prms := url.Values{}
	prms["lease_uuid"] = []string{fmt.Sprintf("%v", leaseUUID)}
	prms["group_uid_hash"] = []string{fmt.Sprintf("%v", groupUIDHash)}
	prms["action"] = []string{fmt.Sprintf("%v", action)}
	{
		sliceVal := []string{sig}
		prms["sig"] = sliceVal
	}
	{
		sliceVal := []string{tok}
		prms["tok"] = sliceVal
	}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "EmailActionTest"), rw, req, prms)
	actionsCtx, _err := app.NewActionsEmailActionContext(goaCtx, req, service)
	if _err != nil {
		panic("invalid test data " + _err.Error()) // bug
	}

	// Perform action
	_err = ctrl.Actions(actionsCtx)

	// Validate response
	if _err != nil {
		t.Fatalf("controller returned %+v, logs:\n%s", _err, logBuf.String())
	}
	if rw.Code != 200 {
		t.Errorf("invalid response status code: got %+v, expected 200", rw.Code)
	}

	// Return results
	return rw
}
