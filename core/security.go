// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package core

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/net/context"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/tleyden/cecil/goa/app"
	"github.com/tleyden/cecil/tools"
)

// APITokenClaims is the claims struct
type APITokenClaims struct {
	AccountID uint `json:"account_id"`
	jwtgo.StandardClaims
}

// signBytes returns the signature of an array of bytes.
func (s *Service) signBytes(bytesToSign []byte) ([]byte, error) {
	if s.rsa.privateKey == nil {
		return []byte{}, fmt.Errorf("s.rsa.privateKey is nil")
	}

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := crypto.SHA256.New()
	pssh.Write(bytesToSign)
	hashed := pssh.Sum(nil)

	// compute signature
	return rsa.SignPSS(rand.Reader, s.rsa.privateKey, crypto.SHA256, hashed, &opts)
}

// verifyBytes verifies the signature of a byte array.
func (s *Service) verifyBytes(bytesToVerify []byte, signature []byte) error {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := crypto.SHA256.New()

	pssh.Write(bytesToVerify)
	hashed := pssh.Sum(nil)

	// verify signature
	return rsa.VerifyPSS(s.rsa.publicKey, crypto.SHA256, hashed, signature, &opts)
}

// emailActionSignURL returns the signature of an email_action URL components.
func (s *Service) emailActionSignURL(leaseUUID string, groupUIDHash string, action, tokenOnce string) ([]byte, error) {

	bytesToSign, err := tools.ConcatBytesFromStrings(leaseUUID, groupUIDHash, action, tokenOnce)
	if err != nil {
		return []byte{}, err
	}

	signature, err := s.signBytes(bytesToSign)
	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

// EmailActionVerifySignatureParams is used to verify the email_action parameters signature.
func (s *Service) EmailActionVerifySignatureParams(leaseUUID, groupUIDHash, action, tokenOnce, signatureBase64 string) error {

	var bytesToVerify bytes.Buffer

	if len(leaseUUID) == 0 {
		return fmt.Errorf("leaseUUID is not set or null in query")
	}
	_, err := bytesToVerify.WriteString(leaseUUID)
	if err != nil {
		return err
	}

	if len(groupUIDHash) == 0 {
		return fmt.Errorf("groupUIDHash is not set or null in query")
	}
	_, err = bytesToVerify.WriteString(groupUIDHash)
	if err != nil {
		return err
	}

	if len(action) == 0 {
		return fmt.Errorf("action is not set or null in query")
	}
	_, err = bytesToVerify.WriteString(action)
	if err != nil {
		return err
	}

	if len(tokenOnce) == 0 {
		return fmt.Errorf("tokenOnce is not set or null in query")
	}
	_, err = bytesToVerify.WriteString(tokenOnce)
	if err != nil {
		return err
	}

	if len(signatureBase64) == 0 {
		return fmt.Errorf("signature is not set or null in query")
	}

	signature, err := base64.URLEncoding.DecodeString(signatureBase64)
	if err != nil {
		return err
	}

	return s.verifyBytes(bytesToVerify.Bytes(), signature)
}

// EmailActionGenerateSignedURL generates an email_action URL with the provided parameters.
func (s *Service) EmailActionGenerateSignedURL(action, leaseUUID string, groupUIDHash string, tokenOnce string) (string, error) {
	// TODO: use AWSResourceID instead of groupUIDHash

	signature, err := s.emailActionSignURL(leaseUUID, groupUIDHash, action, tokenOnce)
	if err != nil {
		return "", fmt.Errorf("error while signing")
	}
	signedURL := fmt.Sprintf("%s/email_action/leases/%s/%v/%s?tok=%s&sig=%s",
		s.CecilHTTPAddress(),
		leaseUUID,
		groupUIDHash,
		action,
		tokenOnce,
		base64.URLEncoding.EncodeToString(signature),
	)
	return signedURL, nil
}

// NewJWTMiddleware creates a middleware that checks for the presence of a JWT Authorization header,
// validates signature, and content.
func (s *Service) NewJWTMiddleware() (goa.Middleware, error) {
	// TODO: use a set of keys to allow rotation, instead of using just one key
	middleware := jwt.New(jwt.NewSimpleResolver([]jwt.Key{s.rsa.publicKey}), nil, app.NewJWTSecurity())
	return middleware, nil
}

// SignToken returns the provided token along with its singnature, in string format.
func (s *Service) SignToken(token *jwtgo.Token) (string, error) {
	return token.SignedString(s.rsa.privateKey)
}

// ValidateToken validates the JWT token given the context.
func ValidateToken(ctx context.Context) (uint, error) {

	// Retrieve the token claims
	token := jwt.ContextJWT(ctx)
	if token == nil {
		Logger.Debug("ValidateToken", "JWT token is missing from context", "context", ctx)
		return 0, fmt.Errorf("JWT token is missing from context") // internal error
	}
	claims := token.Claims.(jwtgo.MapClaims)

	// get the sub attribute
	subClaim, ok := claims["sub"]
	if !ok {
		Logger.Debug("ValidateToken", "'sub' claim not set in claims map", "subClaim", claims)
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
		Logger.Debug("ValidateToken", "'sub' claim is not any of the expected types", fmt.Sprintf("subClaim type: %T", subClaim))

		return 0, errors.New("'sub' claim is not any of the expected types")
	}

	// extract account_id parameter from URL
	reqq := goa.ContextRequest(ctx)
	paramAccountID := reqq.Params["account_id"]

	if len(paramAccountID) == 0 {
		Logger.Debug("ValidateToken", "account_id param in url not set", "reqq.Params", reqq.Params)
		return 0, errors.New("account_id param in url not set")
	}
	rawAccountID := paramAccountID[0]

	accountIDParam, err := strconv.Atoi(rawAccountID)
	if err != nil {
		Logger.Debug("ValidateToken", "cannot parse account_id param", "rawAccountID", rawAccountID, "err", err)
		return 0, errors.New("cannot parse account_id param")
	}

	if accountID != uint(accountIDParam) {
		Logger.Debug("ValidateToken", "accountID != uint(accountIDParam)", "accountID", accountID, "accountIDParam", uint(accountIDParam))
		return 0, tools.ErrorUnauthorized
	}

	return accountID, nil
}

func Hash(b []byte) string {
	hasher := sha512.New()
	hasher.Write(b)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func HashString(s string) string {
	hasher := sha512.New()
	hasher.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
