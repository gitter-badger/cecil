//************************************************************************//
// API "zerocloud": Application User Types
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

import "github.com/goadesign/goa"

// cloudAccountPayload user type.
type cloudAccountPayload struct {
	Cloudprovider *string `form:"cloudprovider,omitempty" json:"cloudprovider,omitempty" xml:"cloudprovider,omitempty"`
	// Name of account
	Name              *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	UpstreamAccountID *string `form:"upstream_account_id,omitempty" json:"upstream_account_id,omitempty" xml:"upstream_account_id,omitempty"`
}

// Validate validates the cloudAccountPayload type instance.
func (ut *cloudAccountPayload) Validate() (err error) {
	if ut.Cloudprovider != nil {
		if len(*ut.Cloudprovider) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.cloudprovider`, *ut.Cloudprovider, len(*ut.Cloudprovider), 3, true))
		}
	}
	if ut.Name != nil {
		if len(*ut.Name) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, *ut.Name, len(*ut.Name), 3, true))
		}
	}
	if ut.UpstreamAccountID != nil {
		if len(*ut.UpstreamAccountID) < 4 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.upstream_account_id`, *ut.UpstreamAccountID, len(*ut.UpstreamAccountID), 4, true))
		}
	}
	return
}

// Publicize creates CloudAccountPayload from cloudAccountPayload
func (ut *cloudAccountPayload) Publicize() *CloudAccountPayload {
	var pub CloudAccountPayload
	if ut.Cloudprovider != nil {
		pub.Cloudprovider = ut.Cloudprovider
	}
	if ut.Name != nil {
		pub.Name = ut.Name
	}
	if ut.UpstreamAccountID != nil {
		pub.UpstreamAccountID = ut.UpstreamAccountID
	}
	return &pub
}

// CloudAccountPayload user type.
type CloudAccountPayload struct {
	Cloudprovider *string `form:"cloudprovider,omitempty" json:"cloudprovider,omitempty" xml:"cloudprovider,omitempty"`
	// Name of account
	Name              *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	UpstreamAccountID *string `form:"upstream_account_id,omitempty" json:"upstream_account_id,omitempty" xml:"upstream_account_id,omitempty"`
}

// Validate validates the CloudAccountPayload type instance.
func (ut *CloudAccountPayload) Validate() (err error) {
	if ut.Cloudprovider != nil {
		if len(*ut.Cloudprovider) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.cloudprovider`, *ut.Cloudprovider, len(*ut.Cloudprovider), 3, true))
		}
	}
	if ut.Name != nil {
		if len(*ut.Name) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, *ut.Name, len(*ut.Name), 3, true))
		}
	}
	if ut.UpstreamAccountID != nil {
		if len(*ut.UpstreamAccountID) < 4 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.upstream_account_id`, *ut.UpstreamAccountID, len(*ut.UpstreamAccountID), 4, true))
		}
	}
	return
}

// cloudEventPayload user type.
type cloudEventPayload struct {
	Message *struct {
		// AWS Account
		Account *string `form:"account,omitempty" json:"account,omitempty" xml:"account,omitempty"`
	} `form:"Message,omitempty" json:"Message,omitempty" xml:"Message,omitempty"`
}

// Validate validates the cloudEventPayload type instance.
func (ut *cloudEventPayload) Validate() (err error) {
	if ut.Message != nil {
		if ut.Message.Account == nil {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response.Message`, "account"))
		}

	}
	return
}

// Publicize creates CloudEventPayload from cloudEventPayload
func (ut *cloudEventPayload) Publicize() *CloudEventPayload {
	var pub CloudEventPayload
	if ut.Message != nil {
		pub.Message = &struct {
			// AWS Account
			Account string `form:"account" json:"account" xml:"account"`
		}{}
		if ut.Message.Account != nil {
			pub.Message.Account = *ut.Message.Account
		}
	}
	return &pub
}

// CloudEventPayload user type.
type CloudEventPayload struct {
	Message *struct {
		// AWS Account
		Account string `form:"account" json:"account" xml:"account"`
	} `form:"Message,omitempty" json:"Message,omitempty" xml:"Message,omitempty"`
}

// Validate validates the CloudEventPayload type instance.
func (ut *CloudEventPayload) Validate() (err error) {
	if ut.Message != nil {
		if ut.Message.Account == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response.Message`, "account"))
		}

	}
	return
}
