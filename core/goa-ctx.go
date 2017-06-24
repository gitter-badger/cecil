package core

import (
	"context"

	"github.com/tleyden/cecil/models"
)

type cecilContextKey int

const (
	accountIDKey cecilContextKey = iota + 1
	cloudaccountIDKey
)

// WithAccount creates a child context containing the given account
func WithAccount(ctx context.Context, account *models.Account) context.Context {
	return context.WithValue(ctx, accountIDKey, account)
}

// ContextAccount retrieves the *models.Account from a `context` that went through our security middleware.
func ContextAccount(ctx context.Context) *models.Account {
	account, ok := ctx.Value(accountIDKey).(*models.Account)
	if !ok {
		return nil
	}
	return account
}

// WithCloudaccount creates a child context containing the given account
func WithCloudaccount(ctx context.Context, cloudaccount *models.Cloudaccount) context.Context {
	return context.WithValue(ctx, cloudaccountIDKey, cloudaccount)
}

// ContextCloudaccount retrieves the *models.Cloudaccount from a `context` that went through our security middleware.
func ContextCloudaccount(ctx context.Context) *models.Cloudaccount {
	cloudaccount, ok := ctx.Value(cloudaccountIDKey).(*models.Cloudaccount)
	if !ok {
		return nil
	}
	return cloudaccount
}
