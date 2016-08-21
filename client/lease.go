package client

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// ListLeasePath computes a request path to the list action of lease.
func ListLeasePath() string {
	return fmt.Sprintf("/leases")
}

// Retrieve all leases.
func (c *Client) ListLease(ctx context.Context, path string, state string) (*http.Response, error) {
	req, err := c.NewListLeaseRequest(ctx, path, state)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListLeaseRequest create the request corresponding to the list action endpoint of the lease resource.
func (c *Client) NewListLeaseRequest(ctx context.Context, path string, state string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("state", state)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
