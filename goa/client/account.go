package client

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"
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

// MailerConfigAccountPayload is the account mailerConfig action payload.
type MailerConfigAccountPayload struct {
	APIKey       string `form:"api_key" json:"api_key" xml:"api_key"`
	Domain       string `form:"domain" json:"domain" xml:"domain"`
	FromName     string `form:"from_name" json:"from_name" xml:"from_name"`
	PublicAPIKey string `form:"public_api_key" json:"public_api_key" xml:"public_api_key"`
}

// MailerConfigAccountPath computes a request path to the mailerConfig action of account.
func MailerConfigAccountPath(accountID int) string {
	param0 := strconv.Itoa(accountID)

	return fmt.Sprintf("/accounts/%s/mailer_config", param0)
}

// Configure mailer
func (c *Client) MailerConfigAccount(ctx context.Context, path string, payload *MailerConfigAccountPayload) (*http.Response, error) {
	req, err := c.NewMailerConfigAccountRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewMailerConfigAccountRequest create the request corresponding to the mailerConfig action endpoint of the account resource.
func (c *Client) NewMailerConfigAccountRequest(ctx context.Context, path string, payload *MailerConfigAccountPayload) (*http.Request, error) {
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
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// RemoveSlackAccountPath computes a request path to the removeSlack action of account.
func RemoveSlackAccountPath(accountID int) string {
	param0 := strconv.Itoa(accountID)

	return fmt.Sprintf("/accounts/%s/slack_config", param0)
}

// Remove slack
func (c *Client) RemoveSlackAccount(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewRemoveSlackAccountRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRemoveSlackAccountRequest create the request corresponding to the removeSlack action endpoint of the account resource.
func (c *Client) NewRemoveSlackAccountRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// ShowAccountPath computes a request path to the show action of account.
func ShowAccountPath(accountID int) string {
	param0 := strconv.Itoa(accountID)

	return fmt.Sprintf("/accounts/%s", param0)
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

// SlackConfigAccountPayload is the account slackConfig action payload.
type SlackConfigAccountPayload struct {
	ChannelID string `form:"channel_id" json:"channel_id" xml:"channel_id"`
	Token     string `form:"token" json:"token" xml:"token"`
}

// SlackConfigAccountPath computes a request path to the slackConfig action of account.
func SlackConfigAccountPath(accountID int) string {
	param0 := strconv.Itoa(accountID)

	return fmt.Sprintf("/accounts/%s/slack_config", param0)
}

// Configure slack
func (c *Client) SlackConfigAccount(ctx context.Context, path string, payload *SlackConfigAccountPayload) (*http.Response, error) {
	req, err := c.NewSlackConfigAccountRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewSlackConfigAccountRequest create the request corresponding to the slackConfig action endpoint of the account resource.
func (c *Client) NewSlackConfigAccountRequest(ctx context.Context, path string, payload *SlackConfigAccountPayload) (*http.Request, error) {
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
	param0 := strconv.Itoa(accountID)

	return fmt.Sprintf("/accounts/%s/api_token", param0)
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
