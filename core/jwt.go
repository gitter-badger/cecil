package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/goa/app"
	"github.com/tleyden/cecil/tools"
)

// NewJWTMiddleware creates a middleware that checks for the presence of a JWT Authorization header,
// validates signature, and content.
func (s *Service) NewJWTMiddleware() (goa.Middleware, error) {
	// TODO: use a set of keys to allow rotation, instead of using just one key
	middleware := jwt.New(
		jwt.NewSimpleResolver([]jwt.Key{s.rsa.publicKey}),
		s.additionalSecurityValidation(),
		app.NewJWTSecurity(),
	)
	return middleware, nil
}

func (s *Service) additionalSecurityValidation() goa.Middleware {
	return func(nextHandler goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {

			requestContextLogger := NewContextLogger(ctx)

			// check whether the account ID specified in the token corresponds to the ID of the accessed account
			accountID, err := tokenCanAccessAccount(ctx)
			if err != nil {
				if err == ErrParamIDNotFound {
					return nextHandler(ctx, rw, req)
				}
				requestContextLogger.Error("Error validating token", "err", err)
				return jwt.ErrJWTError(err.Error())
			}

			// get the account from the DB
			account, err := s.GetAccountByID(int(accountID))
			if err != nil {
				requestContextLogger.Error("Error fetching account", "err", err)
				if err == gorm.ErrRecordNotFound {
					return jwt.ErrJWTError("account not found")
				}
				return jwt.ErrJWTError("internal server error")
			}
			// save account in the context
			ctx = WithAccount(ctx, account)

			// get the cloudaccount_id from the params
			cloudaccountIDParam, err := intIDFromContextParams(ctx, "cloudaccount_id")
			if err != nil {
				if err == ErrParamIDNotFound {
					return nextHandler(ctx, rw, req)
				}
				requestContextLogger.Error("Error getting cloudaccount_id param", "err", err)
				return jwt.ErrJWTError("error parsing cloudaccount ID")
			}

			// get the cloudaccount from the DB
			cloudaccount, err := s.GetCloudaccountByID(cloudaccountIDParam)
			if err != nil {
				requestContextLogger.Error("Error fetching cloudaccount", "err", err)
				if err == gorm.ErrRecordNotFound {
					return jwt.ErrJWTError("cloudaccount not found")
				}
				return jwt.ErrJWTError("internal server error")
			}

			// check whether everything is consistent
			if !account.IsOwnerOf(cloudaccount) {
				requestContextLogger.Error(fmt.Sprintf("Account %v is not owner of cloudaccount %v", account.ID, cloudaccount.ID))
				return tools.ErrNotFound(ctx, "cloud account not found")
			}
			// save cloudaccount in the context
			ctx = WithCloudaccount(ctx, cloudaccount)

			//return jwt.ErrJWTError("you are not uncle ben's")
			return nextHandler(ctx, rw, req)
		}
	}
}

// tokenCanAccessAccount validates the JWT token given the context.
func tokenCanAccessAccount(ctx context.Context) (uint, error) {

	// Retrieve the token claims
	token := jwt.ContextJWT(ctx)
	if token == nil {
		Logger.Debug("tokenCanAccessAccount", "JWT token is missing from context", "context", ctx)
		return 0, fmt.Errorf("JWT token is missing from context") // internal error
	}
	claims := token.Claims.(jwtgo.MapClaims)

	// get the sub attribute
	subClaim, ok := claims["sub"]
	if !ok {
		Logger.Debug("tokenCanAccessAccount", "'sub' claim not set in claims map", "subClaim", claims)
		return 0, errors.New("'sub' claim not set in claims map")
	}

	var accountID uint

	switch v := subClaim.(type) {
	case int:
		accountID = uint(v)
	case uint:
		accountID = v
	case float64:
		accountID = uint(v)
	default:
		Logger.Debug("tokenCanAccessAccount", "'sub' claim is not any of the expected types", fmt.Sprintf("subClaim type: %T", subClaim))

		return 0, errors.New("'sub' claim is not any of the expected types")
	}

	accountIDParam, err := intIDFromContextParams(ctx, "account_id")
	if err != nil {
		return 0, err
	}

	if accountID != uint(accountIDParam) {
		Logger.Debug("tokenCanAccessAccount", "accountID != uint(accountIDParam)", "accountID", accountID, "accountIDParam", uint(accountIDParam))
		return 0, tools.ErrorUnauthorized
	}

	return accountID, nil
}

func intIDFromContextParams(ctx context.Context, paramName string) (int, error) {
	// extract parameter from URL
	reqq := goa.ContextRequest(ctx)
	paramArray := reqq.Params[paramName]

	if len(paramArray) == 0 {
		Logger.Debug("intIDFromContextParams", paramName+" param in url not set", "reqq.Params", reqq.Params)
		return 0, ErrParamIDNotFound
	}
	rawID := paramArray[0]

	IDParam, err := strconv.Atoi(rawID)
	if err != nil {
		Logger.Debug("intIDFromContextParams", "cannot parse param "+paramName, "rawID", rawID, "err", err)
		return 0, errors.New("cannot parse param" + paramName)
	}
	return IDParam, nil
}

var ErrParamIDNotFound = errors.New("param not found")
