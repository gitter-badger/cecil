package core

import (
	"errors"

	"github.com/goadesign/goa"

	"golang.org/x/net/context"
)

func ErrInvalidRequest(ctx context.Context, message interface{}, keyvals ...interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	return responseData.Service.Send(ctx, 400, goa.ErrInvalidRequest(message, keyvals...))
}
func ErrInternal(ctx context.Context, message interface{}, keyvals ...interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	return responseData.Service.Send(ctx, 500, goa.ErrInternal(message, keyvals...))
}
func ErrNotFound(ctx context.Context, message interface{}, keyvals ...interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	return responseData.Service.Send(ctx, 404, goa.ErrNotFound(message, keyvals...))
}
func ErrUnauthorized(ctx context.Context, message interface{}, keyvals ...interface{}) error {
	responseData := goa.ContextResponse(ctx)
	if responseData == nil {
		return errors.New("cannot extract responseData")
	}
	return responseData.Service.Send(ctx, 401, goa.ErrUnauthorized(message, keyvals...))
}