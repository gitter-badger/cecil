package core

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/goadesign/goa/middleware/security/jwt"
)

// Claims struct
type APITokenClaims struct {
	AccountID uint `json:"account_id"`
	jwtgo.StandardClaims
}

func (s *Service) GenerateAPITokenForAccount(accountID uint) (string, error) {
	// Declare the Claims
	claims := APITokenClaims{
		accountID,
		jwtgo.StandardClaims{
			IssuedAt: time.Now().UTC().Unix(),
		},
	}

	// generate token with claims
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS512, claims)

	// sign token with private key
	signedAPIToken, err := token.SignedString(s.rsa.privateKey)
	if err != nil {
		return "", err
	}

	return signedAPIToken, nil
}

func (s *Service) ParseAndVerifyAPIToken(accountID uint, APIToken string) (*APITokenClaims, error) {
	token, err := jwtgo.ParseWithClaims(APIToken, &APITokenClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwtgo.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the key for validation
		return s.rsa.publicKey, nil
	})

	if err != nil {
		return &APITokenClaims{}, err
	}

	claims, ok := token.Claims.(*APITokenClaims)
	tokenIsValid := ok && token.Valid

	if !tokenIsValid {
		return &APITokenClaims{}, fmt.Errorf("token not valid")
	}

	if accountID != claims.AccountID {
		return &APITokenClaims{}, fmt.Errorf("account id mismatch: %v (param) != %v (token)", accountID, claims.AccountID)
	}

	return claims, nil
}

func (s *Service) mustBeAuthorized() gin.HandlerFunc {
	return func(cc *gin.Context) {
		/*
			TODO: check if https; if NOT https, return error
		*/

		accountIDString, accountIDIsSet := cc.Params.Get("account_id")
		if !accountIDIsSet {
			logger.Error("!accountIDIsSet")
			cc.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accountID, err := strconv.ParseUint(accountIDString, 10, 64)
		if err != nil {
			logger.Error("cannot parse account id", "err", err)
			cc.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		APITokenWithPrefix := cc.Request.Header.Get("Authorization")
		APIToken := strings.TrimPrefix(APITokenWithPrefix, "Bearer ")

		claims, err := s.ParseAndVerifyAPIToken(uint(accountID), APIToken)
		if err != nil {
			logger.Error("error while parsing parsing and verifying api token", "err", err)
			cc.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// TODO: verify that the account exists

		// Store token claims in gin.Context to make them accessible to endpoints
		cc.Set("claims", claims)

		logger.Info(
			"authorized user has requested a page",
			"accountID", claims.AccountID,
			"url", cc.Request.URL,
			"method", cc.Request.Method,
		)

		cc.Next() //continue
	}
}

func generateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey
	var err error

	// generate Private Key
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return &rsa.PrivateKey{}, &rsa.PublicKey{}, err
	}

	// precompute some calculations
	privateKey.Precompute()

	// validate Private Key
	if err = privateKey.Validate(); err != nil {
		return &rsa.PrivateKey{}, &rsa.PublicKey{}, err
	}

	// public key address of RSA key
	publicKey = &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func (s *Service) concatBytesFromStrings(str ...string) ([]byte, error) {
	var concatBytesBuffer bytes.Buffer
	for _, stringValue := range str {
		_, err := concatBytesBuffer.WriteString(stringValue)
		if err != nil {
			return []byte{}, err
		}
	}
	return concatBytesBuffer.Bytes(), nil
}

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

func (s *Service) verifyBytes(bytesToVerify []byte, signature []byte) error {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := crypto.SHA256.New()

	pssh.Write(bytesToVerify)
	hashed := pssh.Sum(nil)

	// verify signature
	return rsa.VerifyPSS(s.rsa.publicKey, crypto.SHA256, hashed, signature, &opts)
}

func (s *Service) emailActionSignURL(lease_uuid, instance_id, action, token_once string) ([]byte, error) {

	bytesToSign, err := s.concatBytesFromStrings(lease_uuid, instance_id, action, token_once)
	if err != nil {
		return []byte{}, err
	}

	signature, err := s.signBytes(bytesToSign)
	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

func (s *Service) emailActionVerifySignature(c *gin.Context) error {

	var bytesToVerify bytes.Buffer

	lease_uuid, exists := c.Params.Get("lease_uuid")
	if !exists || len(lease_uuid) == 0 {
		return fmt.Errorf("lease_uuid is not set or null in query")
	}
	_, err := bytesToVerify.WriteString(lease_uuid)
	if err != nil {
		return err
	}

	instance_id, exists := c.Params.Get("instance_id")
	if !exists || len(instance_id) == 0 {
		return fmt.Errorf("instance_id is not set or null in query")
	}
	_, err = bytesToVerify.WriteString(instance_id)
	if err != nil {
		return err
	}

	action, exists := c.Params.Get("action")
	if !exists || len(action) == 0 {
		return fmt.Errorf("action is not set or null in query")
	}
	_, err = bytesToVerify.WriteString(action)
	if err != nil {
		return err
	}

	token_once, exists := c.GetQuery("t")
	token_once = strings.TrimSpace(token_once)
	if !exists || len(token_once) == 0 {
		return fmt.Errorf("token_once is not set or null in query")
	}
	_, err = bytesToVerify.WriteString(token_once)
	if err != nil {
		return err
	}

	signature_base64, exists := c.GetQuery("s")
	signature_base64 = strings.TrimSpace(signature_base64)
	if !exists || len(signature_base64) == 0 {
		return fmt.Errorf("signature is not set or null in query")
	}

	signature, err := base64.URLEncoding.DecodeString(signature_base64)
	if err != nil {
		return err
	}

	return s.verifyBytes(bytesToVerify.Bytes(), signature)
}

func (s *Service) EmailActionGenerateSignedURL(action, lease_uuid, instance_id, token_once string) (string, error) {
	signature, err := s.emailActionSignURL(lease_uuid, instance_id, action, token_once)
	if err != nil {
		return "", fmt.Errorf("error while signing")
	}
	signedURL := fmt.Sprintf("%s/email_action/leases/%s/%s/%s?t=%s&s=%s",
		s.CecilHTTPAddress(),
		lease_uuid,
		instance_id,
		action,
		token_once,
		base64.URLEncoding.EncodeToString(signature),
	)
	return signedURL, nil
}

/*
// NewJWTMiddleware creates a middleware that checks for the presence of a JWT Authorization header,
// validates signature, and content.
func (s *Service) NewJWTMiddleware() (goa.Middleware, error) {
	// TODO: use a set of keys to allow rotation, instead of using just one key
	//	keys, err := LoadJWTPublicKeys()
	//	if err != nil {
	//		return nil, err
	//	}
	middleware := jwt.New(s.rsa.publicKey, nil, app.NewJWTSecurity())
	return middleware, nil
}

USE THIS TO CREATE A TOKEN:

// declare new token
token := jwtgo.New(jwtgo.SigningMethodRS512)

sevenDays := time.Duration(24*7) * time.Hour
// decide expiry
tokenExpiresAt := time.Now().UTC().Add(sevenDays).Unix()

token.Claims = jwtgo.MapClaims{
	"iss": "cecil-api-backend",   // who creates the token and signs it
	"aud": "cecil-account",    // to whom the token is intended to be sent
	"exp": tokenExpiresAt,          // time when the token will expire (10 minutes from now)
	"jti": uuid.NewV4().String(),   // a unique identifier for the token
	"iat": time.Now().UTC().Unix(), // when the token was issued/created (now)
	"nbf": 3,                       // time before which the token is not yet valid (2 minutes ago)

	"sub":    accountID, // the subject/principal is whom the token is about
	"scopes": "api:access",        // token scope - not a standard claim
}

// sign token
bearerToken, err := c.mi.SignToken(token)
if err != nil {
	return ctx.Service.Send(ctx, 500, goa.ErrInternal("internal server error"))
}
*/

func (s *Service) SignToken(token *jwtgo.Token) (string, error) {
	return token.SignedString(s.rsa.privateKey)
}

func AccountIDFromTokenFromContext(ctx context.Context) (uint, error) {

	// Retrieve the token claims
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return 0, fmt.Errorf("JWT token is missing from context") // internal error
	}
	claims := token.Claims.(jwtgo.MapClaims)

	// get the sub attribute
	subClaim, ok := claims["sub"]
	if !ok {
		return 0, errors.New("'sub' claim not set in claims map")
	}

	// assert uint type
	accountID, ok := subClaim.(uint)
	if !ok {
		return 0, errors.New("'sub' claim is not of type uint")
	}
	return accountID, nil
}
