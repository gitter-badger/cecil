package client

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// AddCloudaccountPayload is the cloudaccount add action payload.
type AddCloudaccountPayload struct {
	AwsID string `form:"aws_id" json:"aws_id" xml:"aws_id"`
}

// AddCloudaccountPath computes a request path to the add action of cloudaccount.
func AddCloudaccountPath(accountID int) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts", accountID)
}

// Add new cloudaccount
func (c *Client) AddCloudaccount(ctx context.Context, path string, payload *AddCloudaccountPayload) (*http.Response, error) {
	req, err := c.NewAddCloudaccountRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewAddCloudaccountRequest create the request corresponding to the add action endpoint of the cloudaccount resource.
func (c *Client) NewAddCloudaccountRequest(ctx context.Context, path string, payload *AddCloudaccountPayload) (*http.Request, error) {
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

// AddEmailToWhitelistCloudaccountPayload is the cloudaccount addEmailToWhitelist action payload.
type AddEmailToWhitelistCloudaccountPayload struct {
	Email string `form:"email" json:"email" xml:"email"`
}

// AddEmailToWhitelistCloudaccountPath computes a request path to the addEmailToWhitelist action of cloudaccount.
func AddEmailToWhitelistCloudaccountPath(accountID int, cloudaccountID int) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts/%v/owners", accountID, cloudaccountID)
}

// Add new email to owner tag whitelist
func (c *Client) AddEmailToWhitelistCloudaccount(ctx context.Context, path string, payload *AddEmailToWhitelistCloudaccountPayload) (*http.Response, error) {
	req, err := c.NewAddEmailToWhitelistCloudaccountRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewAddEmailToWhitelistCloudaccountRequest create the request corresponding to the addEmailToWhitelist action endpoint of the cloudaccount resource.
func (c *Client) NewAddEmailToWhitelistCloudaccountRequest(ctx context.Context, path string, payload *AddEmailToWhitelistCloudaccountPayload) (*http.Request, error) {
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

// DownloadInitialSetupTemplateCloudaccountPath computes a request path to the downloadInitialSetupTemplate action of cloudaccount.
func DownloadInitialSetupTemplateCloudaccountPath(accountID int, cloudaccountID int) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-initial-setup.template", accountID, cloudaccountID)
}

// Download AWS initial setup cloudformation template
func (c *Client) DownloadInitialSetupTemplateCloudaccount(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewDownloadInitialSetupTemplateCloudaccountRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDownloadInitialSetupTemplateCloudaccountRequest create the request corresponding to the downloadInitialSetupTemplate action endpoint of the cloudaccount resource.
func (c *Client) NewDownloadInitialSetupTemplateCloudaccountRequest(ctx context.Context, path string) (*http.Request, error) {
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

// DownloadRegionSetupTemplateCloudaccountPath computes a request path to the downloadRegionSetupTemplate action of cloudaccount.
func DownloadRegionSetupTemplateCloudaccountPath(accountID int, cloudaccountID int) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-region-setup.template", accountID, cloudaccountID)
}

// Download AWS region setup cloudformation template
func (c *Client) DownloadRegionSetupTemplateCloudaccount(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewDownloadRegionSetupTemplateCloudaccountRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDownloadRegionSetupTemplateCloudaccountRequest create the request corresponding to the downloadRegionSetupTemplate action endpoint of the cloudaccount resource.
func (c *Client) NewDownloadRegionSetupTemplateCloudaccountRequest(ctx context.Context, path string) (*http.Request, error) {
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
