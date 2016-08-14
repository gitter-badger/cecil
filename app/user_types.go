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

import (
	"github.com/goadesign/goa"
	"time"
)

// accountPayload user type.
type accountPayload struct {
	// The lease will expire in this many lease_expires_in_units
	LeaseExpiresIn *int `form:"lease_expires_in,omitempty" json:"lease_expires_in,omitempty" xml:"lease_expires_in,omitempty"`
	// The units for the lease_expires_in field
	LeaseExpiresInUnits *string `form:"lease_expires_in_units,omitempty" json:"lease_expires_in_units,omitempty" xml:"lease_expires_in_units,omitempty"`
	// Name of account
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
}

// Finalize sets the default values for accountPayload type instance.
func (ut *accountPayload) Finalize() {
	var defaultLeaseExpiresIn = 3
	if ut.LeaseExpiresIn == nil {
		ut.LeaseExpiresIn = &defaultLeaseExpiresIn
	}
	var defaultLeaseExpiresInUnits = "days"
	if ut.LeaseExpiresInUnits == nil {
		ut.LeaseExpiresInUnits = &defaultLeaseExpiresInUnits
	}
}

// Validate validates the accountPayload type instance.
func (ut *accountPayload) Validate() (err error) {
	if ut.Name == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}

	if ut.LeaseExpiresInUnits != nil {
		if !(*ut.LeaseExpiresInUnits == "seconds" || *ut.LeaseExpiresInUnits == "minutes" || *ut.LeaseExpiresInUnits == "hours" || *ut.LeaseExpiresInUnits == "days") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError(`response.lease_expires_in_units`, *ut.LeaseExpiresInUnits, []interface{}{"seconds", "minutes", "hours", "days"}))
		}
	}
	if ut.Name != nil {
		if len(*ut.Name) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, *ut.Name, len(*ut.Name), 3, true))
		}
	}
	return
}

// Publicize creates AccountPayload from accountPayload
func (ut *accountPayload) Publicize() *AccountPayload {
	var pub AccountPayload
	if ut.LeaseExpiresIn != nil {
		pub.LeaseExpiresIn = *ut.LeaseExpiresIn
	}
	if ut.LeaseExpiresInUnits != nil {
		pub.LeaseExpiresInUnits = *ut.LeaseExpiresInUnits
	}
	if ut.Name != nil {
		pub.Name = *ut.Name
	}
	return &pub
}

// AccountPayload user type.
type AccountPayload struct {
	// The lease will expire in this many lease_expires_in_units
	LeaseExpiresIn int `form:"lease_expires_in" json:"lease_expires_in" xml:"lease_expires_in"`
	// The units for the lease_expires_in field
	LeaseExpiresInUnits string `form:"lease_expires_in_units" json:"lease_expires_in_units" xml:"lease_expires_in_units"`
	// Name of account
	Name string `form:"name" json:"name" xml:"name"`
}

// Validate validates the AccountPayload type instance.
func (ut *AccountPayload) Validate() (err error) {
	if ut.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}

	if !(ut.LeaseExpiresInUnits == "seconds" || ut.LeaseExpiresInUnits == "minutes" || ut.LeaseExpiresInUnits == "hours" || ut.LeaseExpiresInUnits == "days") {
		err = goa.MergeErrors(err, goa.InvalidEnumValueError(`response.lease_expires_in_units`, ut.LeaseExpiresInUnits, []interface{}{"seconds", "minutes", "hours", "days"}))
	}
	if len(ut.Name) < 3 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, ut.Name, len(ut.Name), 3, true))
	}
	return
}

// cloudAccountPayload user type.
type cloudAccountPayload struct {
	AssumeRoleArn        *string `form:"assume_role_arn,omitempty" json:"assume_role_arn,omitempty" xml:"assume_role_arn,omitempty"`
	AssumeRoleExternalID *string `form:"assume_role_external_id,omitempty" json:"assume_role_external_id,omitempty" xml:"assume_role_external_id,omitempty"`
	Cloudprovider        *string `form:"cloudprovider,omitempty" json:"cloudprovider,omitempty" xml:"cloudprovider,omitempty"`
	// Name of account
	Name              *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	UpstreamAccountID *string `form:"upstream_account_id,omitempty" json:"upstream_account_id,omitempty" xml:"upstream_account_id,omitempty"`
}

// Validate validates the cloudAccountPayload type instance.
func (ut *cloudAccountPayload) Validate() (err error) {
	if ut.Name == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if ut.Cloudprovider == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "cloudprovider"))
	}
	if ut.UpstreamAccountID == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "upstream_account_id"))
	}
	if ut.AssumeRoleArn == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "assume_role_arn"))
	}
	if ut.AssumeRoleExternalID == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "assume_role_external_id"))
	}

	if ut.AssumeRoleArn != nil {
		if len(*ut.AssumeRoleArn) < 4 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.assume_role_arn`, *ut.AssumeRoleArn, len(*ut.AssumeRoleArn), 4, true))
		}
	}
	if ut.AssumeRoleExternalID != nil {
		if len(*ut.AssumeRoleExternalID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.assume_role_external_id`, *ut.AssumeRoleExternalID, len(*ut.AssumeRoleExternalID), 1, true))
		}
	}
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
	if ut.AssumeRoleArn != nil {
		pub.AssumeRoleArn = *ut.AssumeRoleArn
	}
	if ut.AssumeRoleExternalID != nil {
		pub.AssumeRoleExternalID = *ut.AssumeRoleExternalID
	}
	if ut.Cloudprovider != nil {
		pub.Cloudprovider = *ut.Cloudprovider
	}
	if ut.Name != nil {
		pub.Name = *ut.Name
	}
	if ut.UpstreamAccountID != nil {
		pub.UpstreamAccountID = *ut.UpstreamAccountID
	}
	return &pub
}

// CloudAccountPayload user type.
type CloudAccountPayload struct {
	AssumeRoleArn        string `form:"assume_role_arn" json:"assume_role_arn" xml:"assume_role_arn"`
	AssumeRoleExternalID string `form:"assume_role_external_id" json:"assume_role_external_id" xml:"assume_role_external_id"`
	Cloudprovider        string `form:"cloudprovider" json:"cloudprovider" xml:"cloudprovider"`
	// Name of account
	Name              string `form:"name" json:"name" xml:"name"`
	UpstreamAccountID string `form:"upstream_account_id" json:"upstream_account_id" xml:"upstream_account_id"`
}

// Validate validates the CloudAccountPayload type instance.
func (ut *CloudAccountPayload) Validate() (err error) {
	if ut.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if ut.Cloudprovider == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "cloudprovider"))
	}
	if ut.UpstreamAccountID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "upstream_account_id"))
	}
	if ut.AssumeRoleArn == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "assume_role_arn"))
	}
	if ut.AssumeRoleExternalID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "assume_role_external_id"))
	}

	if len(ut.AssumeRoleArn) < 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.assume_role_arn`, ut.AssumeRoleArn, len(ut.AssumeRoleArn), 4, true))
	}
	if len(ut.AssumeRoleExternalID) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.assume_role_external_id`, ut.AssumeRoleExternalID, len(ut.AssumeRoleExternalID), 1, true))
	}
	if len(ut.Cloudprovider) < 3 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.cloudprovider`, ut.Cloudprovider, len(ut.Cloudprovider), 3, true))
	}
	if len(ut.Name) < 3 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, ut.Name, len(ut.Name), 3, true))
	}
	if len(ut.UpstreamAccountID) < 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.upstream_account_id`, ut.UpstreamAccountID, len(ut.UpstreamAccountID), 4, true))
	}
	return
}

// cloudEventPayload user type.
type cloudEventPayload struct {
	Message *struct {
		// AWS Account
		Account *string `form:"account,omitempty" json:"account,omitempty" xml:"account,omitempty"`
		Detail  *struct {
			// EC2 Instance ID
			InstanceID *string `form:"instance-id,omitempty" json:"instance-id,omitempty" xml:"instance-id,omitempty"`
			// EC2 Instance State
			State *string `form:"state,omitempty" json:"state,omitempty" xml:"state,omitempty"`
		} `form:"detail,omitempty" json:"detail,omitempty" xml:"detail,omitempty"`
		// CloudWatch Event Detail Type
		DetailType *string `form:"detail-type,omitempty" json:"detail-type,omitempty" xml:"detail-type,omitempty"`
		// CloudWatch Event ID
		ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
		// CloudWatch Event Region
		Region *string `form:"region,omitempty" json:"region,omitempty" xml:"region,omitempty"`
		// CloudWatch Event Source
		Source *string `form:"source,omitempty" json:"source,omitempty" xml:"source,omitempty"`
		// CloudWatch Event Timestamp
		Time *time.Time `form:"time,omitempty" json:"time,omitempty" xml:"time,omitempty"`
		// CloudWatch Event Version
		Version *string `form:"version,omitempty" json:"version,omitempty" xml:"version,omitempty"`
	} `form:"Message,omitempty" json:"Message,omitempty" xml:"Message,omitempty"`
	// CloudWatch Event ID
	MessageID *string `form:"MessageId,omitempty" json:"MessageId,omitempty" xml:"MessageId,omitempty"`
	// SQS Payload Base64
	SQSPayloadBase64 *string `form:"SQSPayloadBase64,omitempty" json:"SQSPayloadBase64,omitempty" xml:"SQSPayloadBase64,omitempty"`
	// CloudWatch Event Timestamp
	Timestamp *time.Time `form:"Timestamp,omitempty" json:"Timestamp,omitempty" xml:"Timestamp,omitempty"`
	// CloudWatch Event Topic ARN
	TopicArn *string `form:"TopicArn,omitempty" json:"TopicArn,omitempty" xml:"TopicArn,omitempty"`
	// CloudWatch Event Type
	Type *string `form:"Type,omitempty" json:"Type,omitempty" xml:"Type,omitempty"`
}

// Validate validates the cloudEventPayload type instance.
func (ut *cloudEventPayload) Validate() (err error) {
	if ut.Message == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "Message"))
	}

	if ut.Message != nil {
		if ut.Message.Account == nil {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response.Message`, "account"))
		}

		if ut.Message.Detail != nil {
			if ut.Message.Detail.InstanceID == nil {
				err = goa.MergeErrors(err, goa.MissingAttributeError(`response.Message.detail`, "instance-id"))
			}
			if ut.Message.Detail.State == nil {
				err = goa.MergeErrors(err, goa.MissingAttributeError(`response.Message.detail`, "state"))
			}

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
			Detail  *struct {
				// EC2 Instance ID
				InstanceID string `form:"instance-id" json:"instance-id" xml:"instance-id"`
				// EC2 Instance State
				State string `form:"state" json:"state" xml:"state"`
			} `form:"detail,omitempty" json:"detail,omitempty" xml:"detail,omitempty"`
			// CloudWatch Event Detail Type
			DetailType *string `form:"detail-type,omitempty" json:"detail-type,omitempty" xml:"detail-type,omitempty"`
			// CloudWatch Event ID
			ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
			// CloudWatch Event Region
			Region *string `form:"region,omitempty" json:"region,omitempty" xml:"region,omitempty"`
			// CloudWatch Event Source
			Source *string `form:"source,omitempty" json:"source,omitempty" xml:"source,omitempty"`
			// CloudWatch Event Timestamp
			Time *time.Time `form:"time,omitempty" json:"time,omitempty" xml:"time,omitempty"`
			// CloudWatch Event Version
			Version *string `form:"version,omitempty" json:"version,omitempty" xml:"version,omitempty"`
		}{}
		if ut.Message.Account != nil {
			pub.Message.Account = *ut.Message.Account
		}
		if ut.Message.Detail != nil {
			pub.Message.Detail = &struct {
				// EC2 Instance ID
				InstanceID string `form:"instance-id" json:"instance-id" xml:"instance-id"`
				// EC2 Instance State
				State string `form:"state" json:"state" xml:"state"`
			}{}
			if ut.Message.Detail.InstanceID != nil {
				pub.Message.Detail.InstanceID = *ut.Message.Detail.InstanceID
			}
			if ut.Message.Detail.State != nil {
				pub.Message.Detail.State = *ut.Message.Detail.State
			}
		}
		if ut.Message.DetailType != nil {
			pub.Message.DetailType = ut.Message.DetailType
		}
		if ut.Message.ID != nil {
			pub.Message.ID = ut.Message.ID
		}
		if ut.Message.Region != nil {
			pub.Message.Region = ut.Message.Region
		}
		if ut.Message.Source != nil {
			pub.Message.Source = ut.Message.Source
		}
		if ut.Message.Time != nil {
			pub.Message.Time = ut.Message.Time
		}
		if ut.Message.Version != nil {
			pub.Message.Version = ut.Message.Version
		}
	}
	if ut.MessageID != nil {
		pub.MessageID = ut.MessageID
	}
	if ut.SQSPayloadBase64 != nil {
		pub.SQSPayloadBase64 = ut.SQSPayloadBase64
	}
	if ut.Timestamp != nil {
		pub.Timestamp = ut.Timestamp
	}
	if ut.TopicArn != nil {
		pub.TopicArn = ut.TopicArn
	}
	if ut.Type != nil {
		pub.Type = ut.Type
	}
	return &pub
}

// CloudEventPayload user type.
type CloudEventPayload struct {
	Message *struct {
		// AWS Account
		Account string `form:"account" json:"account" xml:"account"`
		Detail  *struct {
			// EC2 Instance ID
			InstanceID string `form:"instance-id" json:"instance-id" xml:"instance-id"`
			// EC2 Instance State
			State string `form:"state" json:"state" xml:"state"`
		} `form:"detail,omitempty" json:"detail,omitempty" xml:"detail,omitempty"`
		// CloudWatch Event Detail Type
		DetailType *string `form:"detail-type,omitempty" json:"detail-type,omitempty" xml:"detail-type,omitempty"`
		// CloudWatch Event ID
		ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
		// CloudWatch Event Region
		Region *string `form:"region,omitempty" json:"region,omitempty" xml:"region,omitempty"`
		// CloudWatch Event Source
		Source *string `form:"source,omitempty" json:"source,omitempty" xml:"source,omitempty"`
		// CloudWatch Event Timestamp
		Time *time.Time `form:"time,omitempty" json:"time,omitempty" xml:"time,omitempty"`
		// CloudWatch Event Version
		Version *string `form:"version,omitempty" json:"version,omitempty" xml:"version,omitempty"`
	} `form:"Message" json:"Message" xml:"Message"`
	// CloudWatch Event ID
	MessageID *string `form:"MessageId,omitempty" json:"MessageId,omitempty" xml:"MessageId,omitempty"`
	// SQS Payload Base64
	SQSPayloadBase64 *string `form:"SQSPayloadBase64,omitempty" json:"SQSPayloadBase64,omitempty" xml:"SQSPayloadBase64,omitempty"`
	// CloudWatch Event Timestamp
	Timestamp *time.Time `form:"Timestamp,omitempty" json:"Timestamp,omitempty" xml:"Timestamp,omitempty"`
	// CloudWatch Event Topic ARN
	TopicArn *string `form:"TopicArn,omitempty" json:"TopicArn,omitempty" xml:"TopicArn,omitempty"`
	// CloudWatch Event Type
	Type *string `form:"Type,omitempty" json:"Type,omitempty" xml:"Type,omitempty"`
}

// Validate validates the CloudEventPayload type instance.
func (ut *CloudEventPayload) Validate() (err error) {
	if ut.Message == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "Message"))
	}

	if ut.Message != nil {
		if ut.Message.Account == "" {
			err = goa.MergeErrors(err, goa.MissingAttributeError(`response.Message`, "account"))
		}

		if ut.Message.Detail != nil {
			if ut.Message.Detail.InstanceID == "" {
				err = goa.MergeErrors(err, goa.MissingAttributeError(`response.Message.detail`, "instance-id"))
			}
			if ut.Message.Detail.State == "" {
				err = goa.MergeErrors(err, goa.MissingAttributeError(`response.Message.detail`, "state"))
			}

		}
	}
	return
}
