package client

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// DeleteFromDBLeasesPath computes a request path to the deleteFromDB action of leases.
func DeleteFromDBLeasesPath(accountID int, cloudaccountID int, leaseID int) string {
	param0 := strconv.Itoa(accountID)
	param1 := strconv.Itoa(cloudaccountID)
	param2 := strconv.Itoa(leaseID)

	return fmt.Sprintf("/accounts/%s/cloudaccounts/%s/leases/%s/delete", param0, param1, param2)
}

// DeleteFromDBLeasesPath2 computes a request path to the deleteFromDB action of leases.
func DeleteFromDBLeasesPath2(accountID int, leaseID int) string {
	param0 := strconv.Itoa(accountID)
	param1 := strconv.Itoa(leaseID)

	return fmt.Sprintf("/accounts/%s/leases/%s/delete", param0, param1)
}

// Delete a lease from DB
func (c *Client) DeleteFromDBLeases(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewDeleteFromDBLeasesRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDeleteFromDBLeasesRequest create the request corresponding to the deleteFromDB action endpoint of the leases resource.
func (c *Client) NewDeleteFromDBLeasesRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// ListLeasesForAccountLeasesPath computes a request path to the listLeasesForAccount action of leases.
func ListLeasesForAccountLeasesPath(accountID int) string {
	param0 := strconv.Itoa(accountID)

	return fmt.Sprintf("/accounts/%s/leases", param0)
}

// List all leases for account
func (c *Client) ListLeasesForAccountLeases(ctx context.Context, path string, terminated *bool) (*http.Response, error) {
	req, err := c.NewListLeasesForAccountLeasesRequest(ctx, path, terminated)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListLeasesForAccountLeasesRequest create the request corresponding to the listLeasesForAccount action endpoint of the leases resource.
func (c *Client) NewListLeasesForAccountLeasesRequest(ctx context.Context, path string, terminated *bool) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if terminated != nil {
		tmp25 := strconv.FormatBool(*terminated)
		values.Set("terminated", tmp25)
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// ListLeasesForCloudaccountLeasesPath computes a request path to the listLeasesForCloudaccount action of leases.
func ListLeasesForCloudaccountLeasesPath(accountID int, cloudaccountID int) string {
	param0 := strconv.Itoa(accountID)
	param1 := strconv.Itoa(cloudaccountID)

	return fmt.Sprintf("/accounts/%s/cloudaccounts/%s/leases", param0, param1)
}

// List all leases for cloudAccount
func (c *Client) ListLeasesForCloudaccountLeases(ctx context.Context, path string, terminated *bool) (*http.Response, error) {
	req, err := c.NewListLeasesForCloudaccountLeasesRequest(ctx, path, terminated)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListLeasesForCloudaccountLeasesRequest create the request corresponding to the listLeasesForCloudaccount action endpoint of the leases resource.
func (c *Client) NewListLeasesForCloudaccountLeasesRequest(ctx context.Context, path string, terminated *bool) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if terminated != nil {
		tmp26 := strconv.FormatBool(*terminated)
		values.Set("terminated", tmp26)
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// SetExpiryLeasesPath computes a request path to the setExpiry action of leases.
func SetExpiryLeasesPath(accountID int, cloudaccountID int, leaseID int) string {
	param0 := strconv.Itoa(accountID)
	param1 := strconv.Itoa(cloudaccountID)
	param2 := strconv.Itoa(leaseID)

	return fmt.Sprintf("/accounts/%s/cloudaccounts/%s/leases/%s/expiry", param0, param1, param2)
}

// SetExpiryLeasesPath2 computes a request path to the setExpiry action of leases.
func SetExpiryLeasesPath2(accountID int, leaseID int) string {
	param0 := strconv.Itoa(accountID)
	param1 := strconv.Itoa(leaseID)

	return fmt.Sprintf("/accounts/%s/leases/%s/expiry", param0, param1)
}

// Set expiry of a lease
func (c *Client) SetExpiryLeases(ctx context.Context, path string, expiresAt time.Time) (*http.Response, error) {
	req, err := c.NewSetExpiryLeasesRequest(ctx, path, expiresAt)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewSetExpiryLeasesRequest create the request corresponding to the setExpiry action endpoint of the leases resource.
func (c *Client) NewSetExpiryLeasesRequest(ctx context.Context, path string, expiresAt time.Time) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	tmp27 := expiresAt.Format(time.RFC3339)
	values.Set("expires_at", tmp27)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// ShowLeasesPath computes a request path to the show action of leases.
func ShowLeasesPath(accountID int, cloudaccountID int, leaseID int) string {
	param0 := strconv.Itoa(accountID)
	param1 := strconv.Itoa(cloudaccountID)
	param2 := strconv.Itoa(leaseID)

	return fmt.Sprintf("/accounts/%s/cloudaccounts/%s/leases/%s", param0, param1, param2)
}

// ShowLeasesPath2 computes a request path to the show action of leases.
func ShowLeasesPath2(accountID int, leaseID int) string {
	param0 := strconv.Itoa(accountID)
	param1 := strconv.Itoa(leaseID)

	return fmt.Sprintf("/accounts/%s/leases/%s", param0, param1)
}

// Show a lease
func (c *Client) ShowLeases(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewShowLeasesRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewShowLeasesRequest create the request corresponding to the show action endpoint of the leases resource.
func (c *Client) NewShowLeasesRequest(ctx context.Context, path string) (*http.Request, error) {
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

// TerminateLeasesPath computes a request path to the terminate action of leases.
func TerminateLeasesPath(accountID int, cloudaccountID int, leaseID int) string {
	param0 := strconv.Itoa(accountID)
	param1 := strconv.Itoa(cloudaccountID)
	param2 := strconv.Itoa(leaseID)

	return fmt.Sprintf("/accounts/%s/cloudaccounts/%s/leases/%s/terminate", param0, param1, param2)
}

// TerminateLeasesPath2 computes a request path to the terminate action of leases.
func TerminateLeasesPath2(accountID int, leaseID int) string {
	param0 := strconv.Itoa(accountID)
	param1 := strconv.Itoa(leaseID)

	return fmt.Sprintf("/accounts/%s/leases/%s/terminate", param0, param1)
}

// Terminate a lease
func (c *Client) TerminateLeases(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewTerminateLeasesRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewTerminateLeasesRequest create the request corresponding to the terminate action endpoint of the leases resource.
func (c *Client) NewTerminateLeasesRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}
