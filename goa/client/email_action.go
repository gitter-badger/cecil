package client

import (
	"fmt"
	uuid "github.com/goadesign/goa/uuid"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// ActionsEmailActionPath computes a request path to the actions action of email_action.
func ActionsEmailActionPath(leaseUUID uuid.UUID, instanceID string, action string) string {
	return fmt.Sprintf("/email_action/leases/%v/%v/%v", leaseUUID, instanceID, action)
}

// Perform an action on a lease
func (c *Client) ActionsEmailAction(ctx context.Context, path string, sig string, tok string) (*http.Response, error) {
	req, err := c.NewActionsEmailActionRequest(ctx, path, sig, tok)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewActionsEmailActionRequest create the request corresponding to the actions action endpoint of the email_action resource.
func (c *Client) NewActionsEmailActionRequest(ctx context.Context, path string, sig string, tok string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("sig", sig)
	values.Set("tok", tok)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
