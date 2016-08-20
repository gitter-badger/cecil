//************************************************************************//
// API "zerocloud": Models
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/zerocloud/design
// --out=$(GOPATH)/src/github.com/tleyden/zerocloud
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package models

import (
	"time"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	"github.com/tleyden/zerocloud/app"
	"golang.org/x/net/context"
)

// Lease Model
type Lease struct {
	ID             int `gorm:"primary_key"` // primary key
	AccountID      int // Belongs To Account
	CloudAccountID int // Belongs To CloudAccount
	CloudEventID   int // Belongs To CloudEvent
	CreatedAt      time.Time
	DeletedAt      *time.Time
	State          string
	UpdatedAt      time.Time
	Expires        time.Time // timestamp
	Account        Account
	CloudAccount   CloudAccount
	CloudEvent     CloudEvent
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m Lease) TableName() string {
	return "leases"

}

// LeaseDB is the implementation of the storage interface for
// Lease.
type LeaseDB struct {
	Db *gorm.DB
}

// NewLeaseDB creates a new storage type.
func NewLeaseDB(db *gorm.DB) *LeaseDB {
	return &LeaseDB{Db: db}
}

// DB returns the underlying database.
func (m *LeaseDB) DB() interface{} {
	return m.Db
}

// LeaseStorage represents the storage interface.
type LeaseStorage interface {
	DB() interface{}
	List(ctx context.Context) ([]*Lease, error)
	Get(ctx context.Context, id int) (*Lease, error)
	Add(ctx context.Context, lease *Lease) error
	Update(ctx context.Context, lease *Lease) error
	Delete(ctx context.Context, id int) error

	ListLease(ctx context.Context, accountID int, cloudAccountID int, cloudEventID int) []*app.Lease
	OneLease(ctx context.Context, id int, accountID int, cloudAccountID int, cloudEventID int) (*app.Lease, error)

	ListLeaseLink(ctx context.Context, accountID int, cloudAccountID int, cloudEventID int) []*app.LeaseLink
	OneLeaseLink(ctx context.Context, id int, accountID int, cloudAccountID int, cloudEventID int) (*app.LeaseLink, error)

	ListLeaseTiny(ctx context.Context, accountID int, cloudAccountID int, cloudEventID int) []*app.LeaseTiny
	OneLeaseTiny(ctx context.Context, id int, accountID int, cloudAccountID int, cloudEventID int) (*app.LeaseTiny, error)
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m *LeaseDB) TableName() string {
	return "leases"

}

// Belongs To Relationships

// LeaseFilterByAccount is a gorm filter for a Belongs To relationship.
func LeaseFilterByAccount(accountID int, originaldb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if accountID > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("account_id = ?", accountID)

		}
	}
	return func(db *gorm.DB) *gorm.DB { return db }
}

// Belongs To Relationships

// LeaseFilterByCloudAccount is a gorm filter for a Belongs To relationship.
func LeaseFilterByCloudAccount(cloudAccountID int, originaldb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if cloudAccountID > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("cloud_account_id = ?", cloudAccountID)

		}
	}
	return func(db *gorm.DB) *gorm.DB { return db }
}

// Belongs To Relationships

// LeaseFilterByCloudEvent is a gorm filter for a Belongs To relationship.
func LeaseFilterByCloudEvent(cloudEventID int, originaldb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if cloudEventID > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("cloud_event_id = ?", cloudEventID)

		}
	}
	return func(db *gorm.DB) *gorm.DB { return db }
}

// CRUD Functions

// Get returns a single Lease as a Database Model
// This is more for use internally, and probably not what you want in  your controllers
func (m *LeaseDB) Get(ctx context.Context, id int) (*Lease, error) {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "get"}, time.Now())

	var native Lease
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of Lease
func (m *LeaseDB) List(ctx context.Context) ([]*Lease, error) {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "list"}, time.Now())

	var objs []*Lease
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *LeaseDB) Add(ctx context.Context, model *Lease) error {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "add"}, time.Now())

	err := m.Db.Create(model).Error
	if err != nil {
		goa.LogError(ctx, "error adding Lease", "error", err.Error())
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *LeaseDB) Update(ctx context.Context, model *Lease) error {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "update"}, time.Now())

	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		goa.LogError(ctx, "error updating Lease", "error", err.Error())
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *LeaseDB) Delete(ctx context.Context, id int) error {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "delete"}, time.Now())

	var obj Lease

	err := m.Db.Delete(&obj, id).Error

	if err != nil {
		goa.LogError(ctx, "error deleting Lease", "error", err.Error())
		return err
	}

	return nil
}
