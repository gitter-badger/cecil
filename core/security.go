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
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/tleyden/cecil/tools"
)

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

// SignToken returns the provided token along with its singnature, in string format.
func (s *Service) SignToken(token *jwtgo.Token) (string, error) {
	return token.SignedString(s.rsa.privateKey)
}

// Hash produces a sha512 hash of a byte array
func Hash(b []byte) string {
	hasher := sha512.New()
	hasher.Write(b)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// HashString produces a sha512 hash of a string
func HashString(s string) string {
	hasher := sha512.New()
	hasher.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
