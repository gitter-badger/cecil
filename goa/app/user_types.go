//************************************************************************//
// API "Cecil": Application User Types
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/cecil/design
// --out=$(GOPATH)/src/github.com/tleyden/cecil/goa
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package app

import (
	"github.com/goadesign/goa"
	"unicode/utf8"
)

// accountInputPayload user type.
type accountInputPayload struct {
	Email   *string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
	Name    *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	Surname *string `form:"surname,omitempty" json:"surname,omitempty" xml:"surname,omitempty"`
}

// Validate validates the accountInputPayload type instance.
func (ut *accountInputPayload) Validate() (err error) {
	if ut.Email != nil {
		if err2 := goa.ValidateFormat(goa.FormatEmail, *ut.Email); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`response.email`, *ut.Email, goa.FormatEmail, err2))
		}
	}
	if ut.Name != nil {
		if utf8.RuneCountInString(*ut.Name) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, *ut.Name, utf8.RuneCountInString(*ut.Name), 1, true))
		}
	}
	if ut.Name != nil {
		if utf8.RuneCountInString(*ut.Name) > 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, *ut.Name, utf8.RuneCountInString(*ut.Name), 30, false))
		}
	}
	if ut.Surname != nil {
		if utf8.RuneCountInString(*ut.Surname) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.surname`, *ut.Surname, utf8.RuneCountInString(*ut.Surname), 1, true))
		}
	}
	if ut.Surname != nil {
		if utf8.RuneCountInString(*ut.Surname) > 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.surname`, *ut.Surname, utf8.RuneCountInString(*ut.Surname), 30, false))
		}
	}
	return
}

// Publicize creates AccountInputPayload from accountInputPayload
func (ut *accountInputPayload) Publicize() *AccountInputPayload {
	var pub AccountInputPayload
	if ut.Email != nil {
		pub.Email = ut.Email
	}
	if ut.Name != nil {
		pub.Name = ut.Name
	}
	if ut.Surname != nil {
		pub.Surname = ut.Surname
	}
	return &pub
}

// AccountInputPayload user type.
type AccountInputPayload struct {
	Email   *string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
	Name    *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	Surname *string `form:"surname,omitempty" json:"surname,omitempty" xml:"surname,omitempty"`
}

// Validate validates the AccountInputPayload type instance.
func (ut *AccountInputPayload) Validate() (err error) {
	if ut.Email != nil {
		if err2 := goa.ValidateFormat(goa.FormatEmail, *ut.Email); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`response.email`, *ut.Email, goa.FormatEmail, err2))
		}
	}
	if ut.Name != nil {
		if utf8.RuneCountInString(*ut.Name) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, *ut.Name, utf8.RuneCountInString(*ut.Name), 1, true))
		}
	}
	if ut.Name != nil {
		if utf8.RuneCountInString(*ut.Name) > 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, *ut.Name, utf8.RuneCountInString(*ut.Name), 30, false))
		}
	}
	if ut.Surname != nil {
		if utf8.RuneCountInString(*ut.Surname) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.surname`, *ut.Surname, utf8.RuneCountInString(*ut.Surname), 1, true))
		}
	}
	if ut.Surname != nil {
		if utf8.RuneCountInString(*ut.Surname) > 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.surname`, *ut.Surname, utf8.RuneCountInString(*ut.Surname), 30, false))
		}
	}
	return
}

// accountVerificationInputPayload user type.
type accountVerificationInputPayload struct {
	VerificationToken *string `form:"verification_token,omitempty" json:"verification_token,omitempty" xml:"verification_token,omitempty"`
}

// Validate validates the accountVerificationInputPayload type instance.
func (ut *accountVerificationInputPayload) Validate() (err error) {
	if ut.VerificationToken != nil {
		if utf8.RuneCountInString(*ut.VerificationToken) < 108 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.verification_token`, *ut.VerificationToken, utf8.RuneCountInString(*ut.VerificationToken), 108, true))
		}
	}
	return
}

// Publicize creates AccountVerificationInputPayload from accountVerificationInputPayload
func (ut *accountVerificationInputPayload) Publicize() *AccountVerificationInputPayload {
	var pub AccountVerificationInputPayload
	if ut.VerificationToken != nil {
		pub.VerificationToken = ut.VerificationToken
	}
	return &pub
}

// AccountVerificationInputPayload user type.
type AccountVerificationInputPayload struct {
	VerificationToken *string `form:"verification_token,omitempty" json:"verification_token,omitempty" xml:"verification_token,omitempty"`
}

// Validate validates the AccountVerificationInputPayload type instance.
func (ut *AccountVerificationInputPayload) Validate() (err error) {
	if ut.VerificationToken != nil {
		if utf8.RuneCountInString(*ut.VerificationToken) < 108 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.verification_token`, *ut.VerificationToken, utf8.RuneCountInString(*ut.VerificationToken), 108, true))
		}
	}
	return
}

// cloudAccountInputPayload user type.
type cloudAccountInputPayload struct {
	AwsID                *string `form:"aws_id,omitempty" json:"aws_id,omitempty" xml:"aws_id,omitempty"`
	DefaultLeaseDuration *string `form:"default_lease_duration,omitempty" json:"default_lease_duration,omitempty" xml:"default_lease_duration,omitempty"`
}

// Validate validates the cloudAccountInputPayload type instance.
func (ut *cloudAccountInputPayload) Validate() (err error) {
	if ut.AwsID != nil {
		if utf8.RuneCountInString(*ut.AwsID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.aws_id`, *ut.AwsID, utf8.RuneCountInString(*ut.AwsID), 1, true))
		}
	}
	return
}

// Publicize creates CloudAccountInputPayload from cloudAccountInputPayload
func (ut *cloudAccountInputPayload) Publicize() *CloudAccountInputPayload {
	var pub CloudAccountInputPayload
	if ut.AwsID != nil {
		pub.AwsID = ut.AwsID
	}
	if ut.DefaultLeaseDuration != nil {
		pub.DefaultLeaseDuration = ut.DefaultLeaseDuration
	}
	return &pub
}

// CloudAccountInputPayload user type.
type CloudAccountInputPayload struct {
	AwsID                *string `form:"aws_id,omitempty" json:"aws_id,omitempty" xml:"aws_id,omitempty"`
	DefaultLeaseDuration *string `form:"default_lease_duration,omitempty" json:"default_lease_duration,omitempty" xml:"default_lease_duration,omitempty"`
}

// Validate validates the CloudAccountInputPayload type instance.
func (ut *CloudAccountInputPayload) Validate() (err error) {
	if ut.AwsID != nil {
		if utf8.RuneCountInString(*ut.AwsID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.aws_id`, *ut.AwsID, utf8.RuneCountInString(*ut.AwsID), 1, true))
		}
	}
	return
}

// mailerConfigInputPayload user type.
type mailerConfigInputPayload struct {
	APIKey       *string `form:"api_key,omitempty" json:"api_key,omitempty" xml:"api_key,omitempty"`
	Domain       *string `form:"domain,omitempty" json:"domain,omitempty" xml:"domain,omitempty"`
	FromName     *string `form:"from_name,omitempty" json:"from_name,omitempty" xml:"from_name,omitempty"`
	PublicAPIKey *string `form:"public_api_key,omitempty" json:"public_api_key,omitempty" xml:"public_api_key,omitempty"`
}

// Validate validates the mailerConfigInputPayload type instance.
func (ut *mailerConfigInputPayload) Validate() (err error) {
	if ut.APIKey != nil {
		if utf8.RuneCountInString(*ut.APIKey) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.api_key`, *ut.APIKey, utf8.RuneCountInString(*ut.APIKey), 1, true))
		}
	}
	if ut.Domain != nil {
		if utf8.RuneCountInString(*ut.Domain) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.domain`, *ut.Domain, utf8.RuneCountInString(*ut.Domain), 1, true))
		}
	}
	if ut.FromName != nil {
		if utf8.RuneCountInString(*ut.FromName) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.from_name`, *ut.FromName, utf8.RuneCountInString(*ut.FromName), 1, true))
		}
	}
	if ut.PublicAPIKey != nil {
		if utf8.RuneCountInString(*ut.PublicAPIKey) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.public_api_key`, *ut.PublicAPIKey, utf8.RuneCountInString(*ut.PublicAPIKey), 1, true))
		}
	}
	return
}

// Publicize creates MailerConfigInputPayload from mailerConfigInputPayload
func (ut *mailerConfigInputPayload) Publicize() *MailerConfigInputPayload {
	var pub MailerConfigInputPayload
	if ut.APIKey != nil {
		pub.APIKey = ut.APIKey
	}
	if ut.Domain != nil {
		pub.Domain = ut.Domain
	}
	if ut.FromName != nil {
		pub.FromName = ut.FromName
	}
	if ut.PublicAPIKey != nil {
		pub.PublicAPIKey = ut.PublicAPIKey
	}
	return &pub
}

// MailerConfigInputPayload user type.
type MailerConfigInputPayload struct {
	APIKey       *string `form:"api_key,omitempty" json:"api_key,omitempty" xml:"api_key,omitempty"`
	Domain       *string `form:"domain,omitempty" json:"domain,omitempty" xml:"domain,omitempty"`
	FromName     *string `form:"from_name,omitempty" json:"from_name,omitempty" xml:"from_name,omitempty"`
	PublicAPIKey *string `form:"public_api_key,omitempty" json:"public_api_key,omitempty" xml:"public_api_key,omitempty"`
}

// Validate validates the MailerConfigInputPayload type instance.
func (ut *MailerConfigInputPayload) Validate() (err error) {
	if ut.APIKey != nil {
		if utf8.RuneCountInString(*ut.APIKey) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.api_key`, *ut.APIKey, utf8.RuneCountInString(*ut.APIKey), 1, true))
		}
	}
	if ut.Domain != nil {
		if utf8.RuneCountInString(*ut.Domain) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.domain`, *ut.Domain, utf8.RuneCountInString(*ut.Domain), 1, true))
		}
	}
	if ut.FromName != nil {
		if utf8.RuneCountInString(*ut.FromName) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.from_name`, *ut.FromName, utf8.RuneCountInString(*ut.FromName), 1, true))
		}
	}
	if ut.PublicAPIKey != nil {
		if utf8.RuneCountInString(*ut.PublicAPIKey) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.public_api_key`, *ut.PublicAPIKey, utf8.RuneCountInString(*ut.PublicAPIKey), 1, true))
		}
	}
	return
}

// ownerInputPayload user type.
type ownerInputPayload struct {
	Email *string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
}

// Validate validates the ownerInputPayload type instance.
func (ut *ownerInputPayload) Validate() (err error) {
	if ut.Email != nil {
		if err2 := goa.ValidateFormat(goa.FormatEmail, *ut.Email); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`response.email`, *ut.Email, goa.FormatEmail, err2))
		}
	}
	return
}

// Publicize creates OwnerInputPayload from ownerInputPayload
func (ut *ownerInputPayload) Publicize() *OwnerInputPayload {
	var pub OwnerInputPayload
	if ut.Email != nil {
		pub.Email = ut.Email
	}
	return &pub
}

// OwnerInputPayload user type.
type OwnerInputPayload struct {
	Email *string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
}

// Validate validates the OwnerInputPayload type instance.
func (ut *OwnerInputPayload) Validate() (err error) {
	if ut.Email != nil {
		if err2 := goa.ValidateFormat(goa.FormatEmail, *ut.Email); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`response.email`, *ut.Email, goa.FormatEmail, err2))
		}
	}
	return
}

// slackConfigInputPayload user type.
type slackConfigInputPayload struct {
	ChannelID *string `form:"channel_id,omitempty" json:"channel_id,omitempty" xml:"channel_id,omitempty"`
	Token     *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// Validate validates the slackConfigInputPayload type instance.
func (ut *slackConfigInputPayload) Validate() (err error) {
	if ut.ChannelID != nil {
		if utf8.RuneCountInString(*ut.ChannelID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.channel_id`, *ut.ChannelID, utf8.RuneCountInString(*ut.ChannelID), 1, true))
		}
	}
	if ut.Token != nil {
		if utf8.RuneCountInString(*ut.Token) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.token`, *ut.Token, utf8.RuneCountInString(*ut.Token), 1, true))
		}
	}
	return
}

// Publicize creates SlackConfigInputPayload from slackConfigInputPayload
func (ut *slackConfigInputPayload) Publicize() *SlackConfigInputPayload {
	var pub SlackConfigInputPayload
	if ut.ChannelID != nil {
		pub.ChannelID = ut.ChannelID
	}
	if ut.Token != nil {
		pub.Token = ut.Token
	}
	return &pub
}

// SlackConfigInputPayload user type.
type SlackConfigInputPayload struct {
	ChannelID *string `form:"channel_id,omitempty" json:"channel_id,omitempty" xml:"channel_id,omitempty"`
	Token     *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// Validate validates the SlackConfigInputPayload type instance.
func (ut *SlackConfigInputPayload) Validate() (err error) {
	if ut.ChannelID != nil {
		if utf8.RuneCountInString(*ut.ChannelID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.channel_id`, *ut.ChannelID, utf8.RuneCountInString(*ut.ChannelID), 1, true))
		}
	}
	if ut.Token != nil {
		if utf8.RuneCountInString(*ut.Token) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.token`, *ut.Token, utf8.RuneCountInString(*ut.Token), 1, true))
		}
	}
	return
}

// subscribeSNSToSQSInputPayload user type.
type subscribeSNSToSQSInputPayload struct {
	Regions []string `form:"regions,omitempty" json:"regions,omitempty" xml:"regions,omitempty"`
}

// Publicize creates SubscribeSNSToSQSInputPayload from subscribeSNSToSQSInputPayload
func (ut *subscribeSNSToSQSInputPayload) Publicize() *SubscribeSNSToSQSInputPayload {
	var pub SubscribeSNSToSQSInputPayload
	if ut.Regions != nil {
		pub.Regions = ut.Regions
	}
	return &pub
}

// SubscribeSNSToSQSInputPayload user type.
type SubscribeSNSToSQSInputPayload struct {
	Regions []string `form:"regions,omitempty" json:"regions,omitempty" xml:"regions,omitempty"`
}
