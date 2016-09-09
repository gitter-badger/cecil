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

	DefaultLeaseExpiration uint64 `sql:"DEFAULT:0"`

	CloudAccounts []CloudAccount
}

type CloudAccount struct {
	gorm.Model
	AccountID uint

	DefaultLeaseExpiration uint64 `sql:"DEFAULT:0"`
	Provider               string // e.g. AWS
	AWSID                  uint64 `sql:"size:255;unique;index"`

	Disabled bool `sql:"DEFAULT:false"`
	Deleted  bool `sql:"DEFAULT:false"`

	Leases  []Lease
	Regions []Region
	Owners  []Owner
}

type Lease struct {
	gorm.Model
	CloudAccountID uint
	OwnerID        uint

	AWSAccountID uint64
	InstanceID   string
	Region       string

	Terminated bool `sql:"DEFAULT:false"`
	Deleted    bool `sql:"DEFAULT:false"`

	LaunchedAt   time.Time
	ExpiresAt    time.Time
	InstanceType string
}

type Region struct {
	gorm.Model
	CloudAccountID uint

	Region string

	Deleted bool `sql:"DEFAULT:false"`
}

type Owner struct {
	gorm.Model
	CloudAccountID uint

	Email    string
	Disabled bool `sql:"DEFAULT:false"`
	Leases   []Lease
}
