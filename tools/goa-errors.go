// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package tools

import (
	"errors"

	"github.com/goadesign/goa"

	"golang.org/x/net/context"
)

// @@@ Helper functions and variables for Goa @@@

// ErrorInvalidRequest is an error
var ErrorInvalidRequest = errors.New("invalid request")

// ErrorInternal is an error
var ErrorInternal = errors.New("internal server error")

// ErrorNotFound is an error
var ErrorNotFound = errors.New("not found")

// ErrorUnauthorized is an error
var ErrorUnauthorized = errors.New("not authorized")

// ErrInvalidRequest is a shortcut for goa.ErrInvalidRequest, along with the right status code
func ErrInvalidRequest(ctx context.Context, message interface{}, keyvals ...interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	return responseData.Service.Send(ctx, 400, goa.ErrInvalidRequest(message, keyvals...))
}

// ErrInternal is a shortcut for goa.ErrInternal, along with the right status code
func ErrInternal(ctx context.Context, message interface{}, keyvals ...interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	return responseData.Service.Send(ctx, 500, goa.ErrInternal(message, keyvals...))
}

// ErrNotFound is a shortcut for goa.ErrNotFound, along with the right status code
func ErrNotFound(ctx context.Context, message interface{}, keyvals ...interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	return responseData.Service.Send(ctx, 404, goa.ErrNotFound(message, keyvals...))
}

// ErrUnauthorized is a shortcut for goa.ErrUnauthorized, along with the right status code
func ErrUnauthorized(ctx context.Context, message interface{}, keyvals ...interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	return responseData.Service.Send(ctx, 401, goa.ErrUnauthorized(message, keyvals...))
}
