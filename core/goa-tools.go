package core

import (
	"encoding/json"
	"errors"

	"github.com/goadesign/goa"
	"golang.org/x/net/context"
)

func JSONResponse(ctx context.Context, code int, v interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	responseData.ResponseWriter.Header().Set("Content-Type", "application/json")
	responseData.ResponseWriter.WriteHeader(code)
	return json.NewEncoder(responseData.ResponseWriter).Encode(v)
}
