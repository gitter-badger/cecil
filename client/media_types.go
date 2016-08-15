//************************************************************************//
// API "zerocloud": Application Media Types
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/zerocloud/design
// --out=$(GOPATH)/src/github.com/tleyden/zerocloud
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package client

import (
	"github.com/goadesign/goa"
	"net/http"
	"time"
)

// An account (default view)
//
// Identifier: application/vnd.account+json; view=default
type Account struct {
	// Date of creation
	CreatedAt time.Time `form:"created_at" json:"created_at" xml:"created_at"`
	// Email of account owner
	CreatedBy string `form:"created_by" json:"created_by" xml:"created_by"`
	// API href of account
	Href string `form:"href" json:"href" xml:"href"`
	// ID of account
	ID int `form:"id" json:"id" xml:"id"`
	// Name of account
	Name string `form:"name" json:"name" xml:"name"`
}

// Validate validates the Account media type instance.
func (mt *Account) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}
	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.CreatedBy == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "created_by"))
	}

	if err2 := goa.ValidateFormat(goa.FormatEmail, mt.CreatedBy); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`response.created_by`, mt.CreatedBy, goa.FormatEmail, err2))
	}
	if len(mt.Name) < 3 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, mt.Name, len(mt.Name), 3, true))
	}
	return
}

// An account (link view)
//
// Identifier: application/vnd.account+json; view=link
type AccountLink struct {
	// API href of account
	Href string `form:"href" json:"href" xml:"href"`
	// ID of account
	ID int `form:"id" json:"id" xml:"id"`
}

// Validate validates the AccountLink media type instance.
func (mt *AccountLink) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}

	return
}

// An account (tiny view)
//
// Identifier: application/vnd.account+json; view=tiny
type AccountTiny struct {
	// API href of account
	Href string `form:"href" json:"href" xml:"href"`
	// ID of account
	ID int `form:"id" json:"id" xml:"id"`
	// Name of account
	Name string `form:"name" json:"name" xml:"name"`
}

// Validate validates the AccountTiny media type instance.
func (mt *AccountTiny) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}
	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}

	if len(mt.Name) < 3 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, mt.Name, len(mt.Name), 3, true))
	}
	return
}

// DecodeAccount decodes the Account instance encoded in resp body.
func (c *Client) DecodeAccount(resp *http.Response) (*Account, error) {
	var decoded Account
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// DecodeAccountLink decodes the AccountLink instance encoded in resp body.
func (c *Client) DecodeAccountLink(resp *http.Response) (*AccountLink, error) {
	var decoded AccountLink
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// DecodeAccountTiny decodes the AccountTiny instance encoded in resp body.
func (c *Client) DecodeAccountTiny(resp *http.Response) (*AccountTiny, error) {
	var decoded AccountTiny
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// AccountCollection is the media type for an array of Account (default view)
//
// Identifier: application/vnd.account+json; type=collection; view=default
type AccountCollection []*Account

// Validate validates the AccountCollection media type instance.
func (mt AccountCollection) Validate() (err error) {
	for _, e := range mt {
		if e.Href == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "href"))
		}
		if e.Name == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "name"))
		}
		if e.CreatedBy == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "created_by"))
		}

		if err2 := goa.ValidateFormat(goa.FormatEmail, e.CreatedBy); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`response[*].created_by`, e.CreatedBy, goa.FormatEmail, err2))
		}
		if len(e.Name) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response[*].name`, e.Name, len(e.Name), 3, true))
		}
	}
	return
}

// AccountCollection is the media type for an array of Account (link view)
//
// Identifier: application/vnd.account+json; type=collection; view=link
type AccountLinkCollection []*AccountLink

// Validate validates the AccountLinkCollection media type instance.
func (mt AccountLinkCollection) Validate() (err error) {
	for _, e := range mt {
		if e.Href == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "href"))
		}

	}
	return
}

// AccountCollection is the media type for an array of Account (tiny view)
//
// Identifier: application/vnd.account+json; type=collection; view=tiny
type AccountTinyCollection []*AccountTiny

// Validate validates the AccountTinyCollection media type instance.
func (mt AccountTinyCollection) Validate() (err error) {
	for _, e := range mt {
		if e.Href == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "href"))
		}
		if e.Name == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "name"))
		}

		if len(e.Name) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response[*].name`, e.Name, len(e.Name), 3, true))
		}
	}
	return
}

// DecodeAccountCollection decodes the AccountCollection instance encoded in resp body.
func (c *Client) DecodeAccountCollection(resp *http.Response) (AccountCollection, error) {
	var decoded AccountCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// DecodeAccountLinkCollection decodes the AccountLinkCollection instance encoded in resp body.
func (c *Client) DecodeAccountLinkCollection(resp *http.Response) (AccountLinkCollection, error) {
	var decoded AccountLinkCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// DecodeAccountTinyCollection decodes the AccountTinyCollection instance encoded in resp body.
func (c *Client) DecodeAccountTinyCollection(resp *http.Response) (AccountTinyCollection, error) {
	var decoded AccountTinyCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// A CloudAccount (default view)
//
// Identifier: application/vnd.cloudaccount+json; view=default
type Cloudaccount struct {
	// Account that owns CloudAccount
	Account *AccountTiny `form:"account,omitempty" json:"account,omitempty" xml:"account,omitempty"`
	// The Role ARN which allows ZeroCloud to use AssumeRole.  See https://github.com/tleyden/zerocloud/issues/1
	AssumeRoleArn string `form:"assume_role_arn" json:"assume_role_arn" xml:"assume_role_arn"`
	// The customer and aws account specific External ID that needs to be passed when using AssumeRole.  See https://github.com/tleyden/zerocloud/issues/1
	AssumeRoleExternalID string `form:"assume_role_external_id" json:"assume_role_external_id" xml:"assume_role_external_id"`
	Cloudprovider        string `form:"cloudprovider" json:"cloudprovider" xml:"cloudprovider"`
	// API href of cloud account
	Href string `form:"href" json:"href" xml:"href"`
	// ID of cloud account
	ID int `form:"id" json:"id" xml:"id"`
	// Links to related resources
	Links *CloudaccountLinks `form:"links,omitempty" json:"links,omitempty" xml:"links,omitempty"`
	// Name of account
	Name              string `form:"name" json:"name" xml:"name"`
	UpstreamAccountID string `form:"upstream_account_id" json:"upstream_account_id" xml:"upstream_account_id"`
}

// Validate validates the Cloudaccount media type instance.
func (mt *Cloudaccount) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}
	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Cloudprovider == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "cloudprovider"))
	}
	if mt.UpstreamAccountID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "upstream_account_id"))
	}
	if mt.AssumeRoleArn == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "assume_role_arn"))
	}
	if mt.AssumeRoleExternalID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "assume_role_external_id"))
	}

	if mt.Account != nil {
		if err2 := mt.Account.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if len(mt.AssumeRoleArn) < 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.assume_role_arn`, mt.AssumeRoleArn, len(mt.AssumeRoleArn), 4, true))
	}
	if len(mt.AssumeRoleExternalID) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.assume_role_external_id`, mt.AssumeRoleExternalID, len(mt.AssumeRoleExternalID), 1, true))
	}
	if len(mt.Cloudprovider) < 3 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.cloudprovider`, mt.Cloudprovider, len(mt.Cloudprovider), 3, true))
	}
	if mt.Links != nil {
		if err2 := mt.Links.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if len(mt.Name) < 3 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, mt.Name, len(mt.Name), 3, true))
	}
	if len(mt.UpstreamAccountID) < 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.upstream_account_id`, mt.UpstreamAccountID, len(mt.UpstreamAccountID), 4, true))
	}
	return
}

// A CloudAccount (link view)
//
// Identifier: application/vnd.cloudaccount+json; view=link
type CloudaccountLink struct {
	// API href of cloud account
	Href string `form:"href" json:"href" xml:"href"`
	// ID of cloud account
	ID int `form:"id" json:"id" xml:"id"`
}

// Validate validates the CloudaccountLink media type instance.
func (mt *CloudaccountLink) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}

	return
}

// A CloudAccount (tiny view)
//
// Identifier: application/vnd.cloudaccount+json; view=tiny
type CloudaccountTiny struct {
	// API href of cloud account
	Href string `form:"href" json:"href" xml:"href"`
	// ID of cloud account
	ID int `form:"id" json:"id" xml:"id"`
	// Links to related resources
	Links *CloudaccountLinks `form:"links,omitempty" json:"links,omitempty" xml:"links,omitempty"`
	// Name of account
	Name string `form:"name" json:"name" xml:"name"`
}

// Validate validates the CloudaccountTiny media type instance.
func (mt *CloudaccountTiny) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}
	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}

	if mt.Links != nil {
		if err2 := mt.Links.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if len(mt.Name) < 3 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, mt.Name, len(mt.Name), 3, true))
	}
	return
}

// CloudaccountLinks contains links to related resources of Cloudaccount.
type CloudaccountLinks struct {
	Account *AccountLink `form:"account,omitempty" json:"account,omitempty" xml:"account,omitempty"`
}

// Validate validates the CloudaccountLinks type instance.
func (ut *CloudaccountLinks) Validate() (err error) {
	if ut.Account != nil {
		if err2 := ut.Account.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// DecodeCloudaccount decodes the Cloudaccount instance encoded in resp body.
func (c *Client) DecodeCloudaccount(resp *http.Response) (*Cloudaccount, error) {
	var decoded Cloudaccount
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// DecodeCloudaccountLink decodes the CloudaccountLink instance encoded in resp body.
func (c *Client) DecodeCloudaccountLink(resp *http.Response) (*CloudaccountLink, error) {
	var decoded CloudaccountLink
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// DecodeCloudaccountTiny decodes the CloudaccountTiny instance encoded in resp body.
func (c *Client) DecodeCloudaccountTiny(resp *http.Response) (*CloudaccountTiny, error) {
	var decoded CloudaccountTiny
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// CloudaccountCollection is the media type for an array of Cloudaccount (default view)
//
// Identifier: application/vnd.cloudaccount+json; type=collection; view=default
type CloudaccountCollection []*Cloudaccount

// Validate validates the CloudaccountCollection media type instance.
func (mt CloudaccountCollection) Validate() (err error) {
	for _, e := range mt {
		if e.Href == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "href"))
		}
		if e.Name == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "name"))
		}
		if e.Cloudprovider == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "cloudprovider"))
		}
		if e.UpstreamAccountID == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "upstream_account_id"))
		}
		if e.AssumeRoleArn == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "assume_role_arn"))
		}
		if e.AssumeRoleExternalID == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "assume_role_external_id"))
		}

		if e.Account != nil {
			if err2 := e.Account.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
		if len(e.AssumeRoleArn) < 4 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response[*].assume_role_arn`, e.AssumeRoleArn, len(e.AssumeRoleArn), 4, true))
		}
		if len(e.AssumeRoleExternalID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response[*].assume_role_external_id`, e.AssumeRoleExternalID, len(e.AssumeRoleExternalID), 1, true))
		}
		if len(e.Cloudprovider) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response[*].cloudprovider`, e.Cloudprovider, len(e.Cloudprovider), 3, true))
		}
		if e.Links != nil {
			if err2 := e.Links.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
		if len(e.Name) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response[*].name`, e.Name, len(e.Name), 3, true))
		}
		if len(e.UpstreamAccountID) < 4 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response[*].upstream_account_id`, e.UpstreamAccountID, len(e.UpstreamAccountID), 4, true))
		}
	}
	return
}

// CloudaccountCollection is the media type for an array of Cloudaccount (tiny view)
//
// Identifier: application/vnd.cloudaccount+json; type=collection; view=tiny
type CloudaccountTinyCollection []*CloudaccountTiny

// Validate validates the CloudaccountTinyCollection media type instance.
func (mt CloudaccountTinyCollection) Validate() (err error) {
	for _, e := range mt {
		if e.Href == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "href"))
		}
		if e.Name == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "name"))
		}

		if e.Links != nil {
			if err2 := e.Links.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
		if len(e.Name) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response[*].name`, e.Name, len(e.Name), 3, true))
		}
	}
	return
}

// CloudaccountLinksArray contains links to related resources of CloudaccountCollection.
type CloudaccountLinksArray []*CloudaccountLinks

// Validate validates the CloudaccountLinksArray type instance.
func (ut CloudaccountLinksArray) Validate() (err error) {
	for _, e := range ut {
		if e.Account != nil {
			if err2 := e.Account.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DecodeCloudaccountCollection decodes the CloudaccountCollection instance encoded in resp body.
func (c *Client) DecodeCloudaccountCollection(resp *http.Response) (CloudaccountCollection, error) {
	var decoded CloudaccountCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// DecodeCloudaccountTinyCollection decodes the CloudaccountTinyCollection instance encoded in resp body.
func (c *Client) DecodeCloudaccountTinyCollection(resp *http.Response) (CloudaccountTinyCollection, error) {
	var decoded CloudaccountTinyCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// A CloudEvent -- AWS CloudWatch Event (default view)
//
// Identifier: application/vnd.cloudevent+json; view=default
type Cloudevent struct {
	AwsAccountID string `form:"aws_account_id" json:"aws_account_id" xml:"aws_account_id"`
	// API href of cloud event
	Href string `form:"href" json:"href" xml:"href"`
	// ID of cloud event
	ID int `form:"id" json:"id" xml:"id"`
	// Links to related resources
	Links *CloudeventLinks `form:"links,omitempty" json:"links,omitempty" xml:"links,omitempty"`
}

// Validate validates the Cloudevent media type instance.
func (mt *Cloudevent) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}
	if mt.AwsAccountID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "aws_account_id"))
	}

	return
}

// A CloudEvent -- AWS CloudWatch Event (tiny view)
//
// Identifier: application/vnd.cloudevent+json; view=tiny
type CloudeventTiny struct {
	// API href of cloud event
	Href string `form:"href" json:"href" xml:"href"`
	// ID of cloud event
	ID int `form:"id" json:"id" xml:"id"`
	// Links to related resources
	Links *CloudeventLinks `form:"links,omitempty" json:"links,omitempty" xml:"links,omitempty"`
}

// Validate validates the CloudeventTiny media type instance.
func (mt *CloudeventTiny) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}

	return
}

// CloudeventLinks contains links to related resources of Cloudevent.
type CloudeventLinks struct {
}

// DecodeCloudevent decodes the Cloudevent instance encoded in resp body.
func (c *Client) DecodeCloudevent(resp *http.Response) (*Cloudevent, error) {
	var decoded Cloudevent
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// DecodeCloudeventTiny decodes the CloudeventTiny instance encoded in resp body.
func (c *Client) DecodeCloudeventTiny(resp *http.Response) (*CloudeventTiny, error) {
	var decoded CloudeventTiny
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// DecodeErrorResponse decodes the ErrorResponse instance encoded in resp body.
func (c *Client) DecodeErrorResponse(resp *http.Response) (*goa.ErrorResponse, error) {
	var decoded goa.ErrorResponse
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// A lease (default view)
//
// Identifier: application/vnd.lease+json; view=default
type Lease struct {
	// Account that owns Lease
	Account *AccountTiny `form:"account,omitempty" json:"account,omitempty" xml:"account,omitempty"`
	// Date of creation
	CreatedAt *time.Time `form:"created_at,omitempty" json:"created_at,omitempty" xml:"created_at,omitempty"`
	// The datetime when this lease expires
	Expires time.Time `form:"expires" json:"expires" xml:"expires"`
	// API href of lease
	Href string `form:"href" json:"href" xml:"href"`
	// ID of lease
	ID int `form:"id" json:"id" xml:"id"`
	// Links to related resources
	Links *LeaseLinks `form:"links,omitempty" json:"links,omitempty" xml:"links,omitempty"`
	// The current state of the lease
	State string `form:"state" json:"state" xml:"state"`
}

// Validate validates the Lease media type instance.
func (mt *Lease) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}
	if mt.State == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "state"))
	}

	if mt.Account != nil {
		if err2 := mt.Account.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if mt.Links != nil {
		if err2 := mt.Links.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// A lease (link view)
//
// Identifier: application/vnd.lease+json; view=link
type LeaseLink struct {
	// API href of lease
	Href string `form:"href" json:"href" xml:"href"`
	// ID of lease
	ID int `form:"id" json:"id" xml:"id"`
}

// Validate validates the LeaseLink media type instance.
func (mt *LeaseLink) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}

	return
}

// A lease (tiny view)
//
// Identifier: application/vnd.lease+json; view=tiny
type LeaseTiny struct {
	// API href of lease
	Href string `form:"href" json:"href" xml:"href"`
	// ID of lease
	ID int `form:"id" json:"id" xml:"id"`
	// Links to related resources
	Links *LeaseLinks `form:"links,omitempty" json:"links,omitempty" xml:"links,omitempty"`
}

// Validate validates the LeaseTiny media type instance.
func (mt *LeaseTiny) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}

	if mt.Links != nil {
		if err2 := mt.Links.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// LeaseLinks contains links to related resources of Lease.
type LeaseLinks struct {
	Account *AccountLink `form:"account,omitempty" json:"account,omitempty" xml:"account,omitempty"`
}

// Validate validates the LeaseLinks type instance.
func (ut *LeaseLinks) Validate() (err error) {
	if ut.Account != nil {
		if err2 := ut.Account.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// DecodeLease decodes the Lease instance encoded in resp body.
func (c *Client) DecodeLease(resp *http.Response) (*Lease, error) {
	var decoded Lease
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// DecodeLeaseLink decodes the LeaseLink instance encoded in resp body.
func (c *Client) DecodeLeaseLink(resp *http.Response) (*LeaseLink, error) {
	var decoded LeaseLink
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// DecodeLeaseTiny decodes the LeaseTiny instance encoded in resp body.
func (c *Client) DecodeLeaseTiny(resp *http.Response) (*LeaseTiny, error) {
	var decoded LeaseTiny
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// LeaseCollection is the media type for an array of Lease (default view)
//
// Identifier: application/vnd.lease+json; type=collection; view=default
type LeaseCollection []*Lease

// Validate validates the LeaseCollection media type instance.
func (mt LeaseCollection) Validate() (err error) {
	for _, e := range mt {
		if e.Href == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "href"))
		}
		if e.State == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "state"))
		}

		if e.Account != nil {
			if err2 := e.Account.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
		if e.Links != nil {
			if err2 := e.Links.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// LeaseCollection is the media type for an array of Lease (link view)
//
// Identifier: application/vnd.lease+json; type=collection; view=link
type LeaseLinkCollection []*LeaseLink

// Validate validates the LeaseLinkCollection media type instance.
func (mt LeaseLinkCollection) Validate() (err error) {
	for _, e := range mt {
		if e.Href == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "href"))
		}

	}
	return
}

// LeaseCollection is the media type for an array of Lease (tiny view)
//
// Identifier: application/vnd.lease+json; type=collection; view=tiny
type LeaseTinyCollection []*LeaseTiny

// Validate validates the LeaseTinyCollection media type instance.
func (mt LeaseTinyCollection) Validate() (err error) {
	for _, e := range mt {
		if e.Href == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response[*]`, "href"))
		}

		if e.Links != nil {
			if err2 := e.Links.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// LeaseLinksArray contains links to related resources of LeaseCollection.
type LeaseLinksArray []*LeaseLinks

// Validate validates the LeaseLinksArray type instance.
func (ut LeaseLinksArray) Validate() (err error) {
	for _, e := range ut {
		if e.Account != nil {
			if err2 := e.Account.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DecodeLeaseCollection decodes the LeaseCollection instance encoded in resp body.
func (c *Client) DecodeLeaseCollection(resp *http.Response) (LeaseCollection, error) {
	var decoded LeaseCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// DecodeLeaseLinkCollection decodes the LeaseLinkCollection instance encoded in resp body.
func (c *Client) DecodeLeaseLinkCollection(resp *http.Response) (LeaseLinkCollection, error) {
	var decoded LeaseLinkCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// DecodeLeaseTinyCollection decodes the LeaseTinyCollection instance encoded in resp body.
func (c *Client) DecodeLeaseTinyCollection(resp *http.Response) (LeaseTinyCollection, error) {
	var decoded LeaseTinyCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}
