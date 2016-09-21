package core

import (
	"time"

	"github.com/jinzhu/gorm"
)

// @@@@@@@@@@@@@@@ DB models @@@@@@@@@@@@@@@

type Account struct {
	gorm.Model
	Email string `sql:"size:255;unique;index"`

	Disabled bool `sql:"DEFAULT:false"`
	Deleted  bool `sql:"DEFAULT:false"`

	DefaultLeaseDuration uint64 `sql:"DEFAULT:0"`

	CloudAccounts []CloudAccount
}

type CloudAccount struct {
	gorm.Model
	AccountID uint

	DefaultLeaseDuration uint64 `sql:"DEFAULT:0"`
	Provider             string // e.g. AWS
	AWSID                string `sql:"size:255;unique;index"`
	ExternalID           string

	Disabled bool `sql:"DEFAULT:false"`
	Deleted  bool `sql:"DEFAULT:false"`

	Leases  []Lease
	Regions []Region
	Owners  []Owner
}

type Lease struct {
	gorm.Model

	UUID      string `sql:"size:255;unique;index"`
	TokenOnce string

	CloudAccountID uint
	OwnerID        uint
	AWSAccountID   string

	InstanceID       string
	Region           string
	AvailabilityZone string
	InstanceType     string

	Terminated bool `sql:"DEFAULT:false"`
	Deleted    bool `sql:"DEFAULT:false"`
	Alerted    bool `sql:"DEFAULT:false"`

	LaunchedAt   time.Time
	ExpiresAt    time.Time
	TerminatedAt time.Time
}

type Owner struct {
	gorm.Model
	CloudAccountID uint

	Email    string
	Disabled bool `sql:"DEFAULT:false"`
	Leases   []Lease
}

type Region struct {
	gorm.Model
	CloudAccountID uint

	Region string

	Deleted bool `sql:"DEFAULT:false"`
}
