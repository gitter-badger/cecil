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

package app

import (
	"github.com/goadesign/goa"
	"time"
)

// A tenant account (default view)
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
	return
}

// A tenant account (link view)
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

// A tenant account (tiny view)
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

	return
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

	}
	return
}

// A CloudAccount (default view)
//
// Identifier: application/vnd.cloudaccount+json; view=default
type Cloudaccount struct {
	// Account that owns CloudAccount
	Account       *AccountTiny `form:"account,omitempty" json:"account,omitempty" xml:"account,omitempty"`
	Cloudprovider string       `form:"cloudprovider" json:"cloudprovider" xml:"cloudprovider"`
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

	if mt.Account != nil {
		if err2 := mt.Account.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
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

		if e.Account != nil {
			if err2 := e.Account.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
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
