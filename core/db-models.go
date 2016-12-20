package core

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// @@@@@@@@@@@@@@@ DB models @@@@@@@@@@@@@@@

// Account is a Cecil account. The owner of an account is called an "Admin"
// because he/she administers one or more cloudaccounts.
type Account struct {
	gorm.Model `json:"_"`
	Email      string `sql:"size:255;unique;index" json:"email"`

	Name    string `sql:"size:255" json:"name"`
	Surname string `sql:"size:255" json:"surname"`

	Verified          bool   `sql:"DEFAULT:false" json:"verified"`
	VerificationToken string `sql:"unique" json:"_"`
	RequestedNewToken bool   `sql:"DEFAULT:false" json:"_"`

	Disabled bool `sql:"DEFAULT:false" json:"_"`
	Deleted  bool `sql:"DEFAULT:false" json:"_"`

	DefaultLeaseDuration time.Duration `sql:"DEFAULT:0" json:"_"`

	CloudAccounts []CloudAccount
	SlackConfig   SlackConfig
}

// CloudAccount is a cloud service account; e.g. an AWS account.
type CloudAccount struct {
	gorm.Model
	AccountID uint

	DefaultLeaseDuration time.Duration `sql:"DEFAULT:0"`
	Provider             string        // e.g. AWS
	AWSID                string        `sql:"size:255;unique;index"`
	ExternalID           string

	Disabled bool `sql:"DEFAULT:false"`
	Deleted  bool `sql:"DEFAULT:false"`

	Leases []Lease
	Owners []Owner
}

// Owner is a whitelisted email address that can/owns (i.e. is responsible of) leases.
type Owner struct {
	gorm.Model
	CloudAccountID uint

	Email    string
	Disabled bool `sql:"DEFAULT:false"`
	Leases   []Lease
}

// Lease is a record of a lease of an AWS EC2 instance.
type Lease struct {
	Model // using this instead of gorm.Model because gorm.Model does not have json tags

	UUID      string `sql:"size:255;unique;index" json:"-"` // TODO: sho or hide UUID on json responses?
	TokenOnce string `json:"-"`

	AccountID      uint   `json:"account_id,omitempty"`
	CloudAccountID uint   `json:"cloud_account_id,omitempty"`
	OwnerID        uint   `json:"owner_id,omitempty"`
	AWSAccountID   string `json:"aws_account_id,omitempty"`

	InstanceID       string `json:"instance_id,omitempty"`
	Region           string `json:"region,omitempty"`
	AvailabilityZone string `json:"availability_zone,omitempty"`
	InstanceType     string `json:"instance_type,omitempty"`

	Terminated bool `json:"terminated"`
	Deleted    bool `json:"deleted,omitempty"`
	Alerted    bool `json:"-"`

	LaunchedAt   time.Time  `json:"launched_at,omitempty"`
	ExpiresAt    time.Time  `json:"expires_at,omitempty"`
	TerminatedAt *time.Time `json:"terminated_at,omitempty"`
}

type SlackConfig struct {
	gorm.Model
	AccountID uint
	Token     string
	ChannelID string
}

type MailerConfig struct {
	gorm.Model
	AccountID uint

	Domain       string
	APIKey       string
	PublicAPIKey string
	FromName     string
}
