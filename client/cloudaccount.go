package client

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// CreateCloudaccountPath computes a request path to the create action of cloudaccount.
func CreateCloudaccountPath(accountID int) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts", accountID)
}

// Record new cloud account
func (c *Client) CreateCloudaccount(ctx context.Context, path string, payload *CloudAccountPayload, contentType string) (*http.Response, error) {
	req, err := c.NewCreateCloudaccountRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCreateCloudaccountRequest create the request corresponding to the create action endpoint of the cloudaccount resource.
func (c *Client) NewCreateCloudaccountRequest(ctx context.Context, path string, payload *CloudAccountPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType != "*/*" {
		header.Set("Content-Type", contentType)
	}
	return req, nil
}

// DeleteCloudaccountPath computes a request path to the delete action of cloudaccount.
func DeleteCloudaccountPath(accountID int, cloudAccountID int) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts/%v", accountID, cloudAccountID)
}

// DeleteCloudaccount makes a request to the delete action endpoint of the cloudaccount resource
func (c *Client) DeleteCloudaccount(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewDeleteCloudaccountRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDeleteCloudaccountRequest create the request corresponding to the delete action endpoint of the cloudaccount resource.
func (c *Client) NewDeleteCloudaccountRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// ListCloudaccountPath computes a request path to the list action of cloudaccount.
func ListCloudaccountPath(accountID int) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts", accountID)
}

// List all cloud accounts
func (c *Client) ListCloudaccount(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewListCloudaccountRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListCloudaccountRequest create the request corresponding to the list action endpoint of the cloudaccount resource.
func (c *Client) NewListCloudaccountRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// ShowCloudaccountPath computes a request path to the show action of cloudaccount.
func ShowCloudaccountPath(accountID int, cloudAccountID int) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts/%v", accountID, cloudAccountID)
}

// Retrieve cloud account with given id
func (c *Client) ShowCloudaccount(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewShowCloudaccountRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewShowCloudaccountRequest create the request corresponding to the show action endpoint of the cloudaccount resource.
func (c *Client) NewShowCloudaccountRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// UpdateCloudaccountPath computes a request path to the update action of cloudaccount.
func UpdateCloudaccountPath(accountID int, cloudAccountID int) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts/%v", accountID, cloudAccountID)
}

// UpdateCloudaccount makes a request to the update action endpoint of the cloudaccount resource
func (c *Client) UpdateCloudaccount(ctx context.Context, path string, payload *CloudAccountPayload, contentType string) (*http.Response, error) {
	req, err := c.NewUpdateCloudaccountRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUpdateCloudaccountRequest create the request corresponding to the update action endpoint of the cloudaccount resource.
func (c *Client) NewUpdateCloudaccountRequest(ctx context.Context, path string, payload *CloudAccountPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("PATCH", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType != "*/*" {
		header.Set("Content-Type", contentType)
	}
	return req, nil
}
