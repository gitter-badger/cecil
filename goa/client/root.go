package client

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// ShowRootPath computes a request path to the show action of root.
func ShowRootPath() string {

	return fmt.Sprintf("/")
}

// Show info about API
func (c *Client) ShowRoot(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewShowRootRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewShowRootRequest create the request corresponding to the show action endpoint of the root resource.
func (c *Client) NewShowRootRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
