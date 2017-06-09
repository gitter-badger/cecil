// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package tools

import (
	"encoding/json"
	"errors"

	"github.com/goadesign/goa"
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
