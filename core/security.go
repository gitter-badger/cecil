package core

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Claims struct
type APITokenClaims struct {
	AccountID uint `json:"account_id"`
	jwt.StandardClaims
}

func (s *Service) GenerateAPITokenForAccount(accountID uint) (string, error) {
	// Declare the Claims
	claims := APITokenClaims{
		accountID,
		jwt.StandardClaims{
			IssuedAt: time.Now().UTC().Unix(),
		},
	}

	// generate token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	// sign token with private key
	signedAPIToken, err := token.SignedString(s.rsa.privateKey)
	if err != nil {
		return "", err
	}

	return signedAPIToken, nil
}

func (s *Service) ParseAndVerifyAPIToken(accountID uint, APIToken string) (*APITokenClaims, error) {
	token, err := jwt.ParseWithClaims(APIToken, &APITokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
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
