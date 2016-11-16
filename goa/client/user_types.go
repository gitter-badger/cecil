//************************************************************************//
// API "Cecil": Application User Types
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/cecil/design
// --out=$(GOPATH)src/github.com/tleyden/cecil/goa
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package client

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
	AwsID *string `form:"aws_id,omitempty" json:"aws_id,omitempty" xml:"aws_id,omitempty"`
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
	return &pub
}

// CloudAccountInputPayload user type.
type CloudAccountInputPayload struct {
	AwsID *string `form:"aws_id,omitempty" json:"aws_id,omitempty" xml:"aws_id,omitempty"`
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
