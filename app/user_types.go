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
		Time *string `form:"time,omitempty" json:"time,omitempty" xml:"time,omitempty"`
		// CloudWatch Event Version
		Version *string `form:"version,omitempty" json:"version,omitempty" xml:"version,omitempty"`
	} `form:"Message,omitempty" json:"Message,omitempty" xml:"Message,omitempty"`
	// SQS Message ID
	MessageID *string `form:"MessageId,omitempty" json:"MessageId,omitempty" xml:"MessageId,omitempty"`
	// SQS Message Timestamp
	Timestamp *string `form:"Timestamp,omitempty" json:"Timestamp,omitempty" xml:"Timestamp,omitempty"`
	// SQS Topic ARN
	TopicArn *string `form:"TopicArn,omitempty" json:"TopicArn,omitempty" xml:"TopicArn,omitempty"`
	// SQS Message Type
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
			Time *string `form:"time,omitempty" json:"time,omitempty" xml:"time,omitempty"`
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
		Time *string `form:"time,omitempty" json:"time,omitempty" xml:"time,omitempty"`
		// CloudWatch Event Version
		Version *string `form:"version,omitempty" json:"version,omitempty" xml:"version,omitempty"`
	} `form:"Message" json:"Message" xml:"Message"`
	// SQS Message ID
	MessageID *string `form:"MessageId,omitempty" json:"MessageId,omitempty" xml:"MessageId,omitempty"`
	// SQS Message Timestamp
	Timestamp *string `form:"Timestamp,omitempty" json:"Timestamp,omitempty" xml:"Timestamp,omitempty"`
	// SQS Topic ARN
	TopicArn *string `form:"TopicArn,omitempty" json:"TopicArn,omitempty" xml:"TopicArn,omitempty"`
	// SQS Message Type
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
