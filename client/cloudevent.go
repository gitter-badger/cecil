package client

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// CreateCloudeventPath computes a request path to the create action of cloudevent.
func CreateCloudeventPath() string {
	return fmt.Sprintf("/cloudevent")
}

// Save a new AWS CloudWatch event
func (c *Client) CreateCloudevent(ctx context.Context, path string, payload *CloudEventPayload, contentType string) (*http.Response, error) {
	req, err := c.NewCreateCloudeventRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCreateCloudeventRequest create the request corresponding to the create action endpoint of the cloudevent resource.
func (c *Client) NewCreateCloudeventRequest(ctx context.Context, path string, payload *CloudEventPayload, contentType string) (*http.Request, error) {
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
