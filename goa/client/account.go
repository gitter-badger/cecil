package client

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// CreateAccountPayload is the account create action payload.
type CreateAccountPayload struct {
	Email   string `form:"email" json:"email" xml:"email"`
	Name    string `form:"name" json:"name" xml:"name"`
	Surname string `form:"surname" json:"surname" xml:"surname"`
}

// CreateAccountPath computes a request path to the create action of account.
func CreateAccountPath() string {
	return fmt.Sprintf("/accounts")
}

// Create new account
func (c *Client) CreateAccount(ctx context.Context, path string, payload *CreateAccountPayload) (*http.Response, error) {
	req, err := c.NewCreateAccountRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCreateAccountRequest create the request corresponding to the create action endpoint of the account resource.
func (c *Client) NewCreateAccountRequest(ctx context.Context, path string, payload *CreateAccountPayload) (*http.Request, error) {
	var body bytes.Buffer
	err := c.Encoder.Encode(payload, &body, "*/*")
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// ShowAccountPath computes a request path to the show action of account.
func ShowAccountPath(accountID int) string {
	return fmt.Sprintf("/accounts/%v", accountID)
}

// Show account
func (c *Client) ShowAccount(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewShowAccountRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewShowAccountRequest create the request corresponding to the show action endpoint of the account resource.
func (c *Client) NewShowAccountRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// VerifyAccountPayload is the account verify action payload.
type VerifyAccountPayload struct {
	VerificationToken string `form:"verification_token" json:"verification_token" xml:"verification_token"`
}

// VerifyAccountPath computes a request path to the verify action of account.
func VerifyAccountPath(accountID int) string {
	return fmt.Sprintf("/accounts/%v/api_token", accountID)
}

// Verify account and get API token
func (c *Client) VerifyAccount(ctx context.Context, path string, payload *VerifyAccountPayload) (*http.Response, error) {
	req, err := c.NewVerifyAccountRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewVerifyAccountRequest create the request corresponding to the verify action endpoint of the account resource.
func (c *Client) NewVerifyAccountRequest(ctx context.Context, path string, payload *VerifyAccountPayload) (*http.Request, error) {
	var body bytes.Buffer
	err := c.Encoder.Encode(payload, &body, "*/*")
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	return req, nil
}
