package client

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// ShowAwsPath computes a request path to the show action of aws.
func ShowAwsPath(awsAccountID string) string {
	return fmt.Sprintf("/aws/%v", awsAccountID)
}

// Lookup the CloudAccount associated with this AWS acocunt ID
func (c *Client) ShowAws(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewShowAwsRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewShowAwsRequest create the request corresponding to the show action endpoint of the aws resource.
func (c *Client) NewShowAwsRequest(ctx context.Context, path string) (*http.Request, error) {
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
