package core

import (
	"encoding/json"
	"errors"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/inconshreveable/log15"
	"golang.org/x/net/context"
)

// JSONResponse is a tool that takes data, converts it to json and returns it
func JSONResponse(ctx context.Context, code int, v interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	responseData.ResponseWriter.Header().Set("Content-Type", "application/json")
	responseData.ResponseWriter.WriteHeader(code)
	return json.NewEncoder(responseData.ResponseWriter).Encode(v)
}

// NewContextLogger returns a new context logger which has been filled in with the request ID
func NewContextLogger(ctx context.Context) log15.Logger {
	request := goa.ContextRequest(ctx)
	return Logger.New(
		"url", request.URL.String(),
		"reqID", middleware.ContextRequestID(ctx),
	)
}
