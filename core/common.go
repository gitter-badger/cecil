package core

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

func runForever(f func() error, sleepDuration time.Duration) {
	for {
		err := f()
		if err != nil {
			logger.Error("runForever", "error", err)
		}
		time.Sleep(sleepDuration)
	}
}

func compileEmail(tpl string, values map[string]interface{}) string {
	var emailBody bytes.Buffer // A Buffer needs no initialization.

	// TODO: check errors ???

	t := template.New("new email template")
	t, _ = t.Parse(tpl)

	_ = t.Execute(&emailBody, values)

	return emailBody.String()
}

func retry(attempts int, sleep time.Duration, callback func() error) (err error) {
	for i := 1; i <= attempts; i++ {

		err = callback()
		if err == nil {
			return nil
		}
		time.Sleep(sleep)

		fmt.Println("Retry error: ", err)
	}
	return fmt.Errorf("Abandoned after %d attempts, last error: %s", attempts, err)
}

func (s *Service) sign(lease_uuid, instance_id, action, token_once string) ([]byte, error) {

	var bytesToSign bytes.Buffer

	_, err := bytesToSign.WriteString(token_once)
	if err != nil {
		return []byte{}, err
	}

	_, err = bytesToSign.WriteString(action)
	if err != nil {
		return []byte{}, err
	}

	_, err = bytesToSign.WriteString(lease_uuid)
	if err != nil {
		return []byte{}, err
	}

	_, err = bytesToSign.WriteString(instance_id)
	if err != nil {
		return []byte{}, err
	}

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := crypto.SHA256.New()
	pssh.Write(bytesToSign.Bytes())
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, s.rsa.privateKey, crypto.SHA256, hashed, &opts)

	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

func (s *Service) verifySignature(c *gin.Context) error {
	// I am recipient.
	// sender's public key; hash of the encrypted message; signed message;

	var bytesToVerify bytes.Buffer

	token_once := ""
	_, err := bytesToVerify.WriteString(token_once)
	if err != nil {
		return err
	}

	action := ""
	_, err = bytesToVerify.WriteString(action)
	if err != nil {
		return err
	}

	lease_uuid := ""
	_, err = bytesToVerify.WriteString(lease_uuid)
	if err != nil {
		return err
	}

	instance_id := ""
	_, err = bytesToVerify.WriteString(instance_id)
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
	// TODO: convert signature_base64 to bytes

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := crypto.SHA256.New()

	pssh.Write(bytesToVerify.Bytes())
	hashed := pssh.Sum(nil)

	//Verify Signature
	err = rsa.VerifyPSS(s.rsa.publicKey, crypto.SHA256, hashed, signature, &opts)

	if err != nil {
		return err
	} else {
		return nil
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
	fmt.Println(privateKey)

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
