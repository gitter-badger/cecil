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
	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	"github.com/tleyden/zerocloud/app"
	"golang.org/x/net/context"
	"time"
)

// ZeroCloud Account
type Account struct {
	ID                  int            `gorm:"primary_key"` // primary key
	CloudAccounts       []CloudAccount // has many CloudAccounts
	CloudEvents         []CloudEvent   // has many CloudEvents
	CreatedAt           time.Time
	DeletedAt           *time.Time
	LeaseExpiresIn      int
	LeaseExpiresInUnits string
	Leases              []Lease // has many Leases
	Name                string
	UpdatedAt           time.Time
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m Account) TableName() string {
	return "accounts"

}

// AccountDB is the implementation of the storage interface for
// Account.
type AccountDB struct {
	Db *gorm.DB
}

// NewAccountDB creates a new storage type.
func NewAccountDB(db *gorm.DB) *AccountDB {
	return &AccountDB{Db: db}
}

// DB returns the underlying database.
func (m *AccountDB) DB() interface{} {
	return m.Db
}

// AccountStorage represents the storage interface.
type AccountStorage interface {
	DB() interface{}
	List(ctx context.Context) ([]*Account, error)
	Get(ctx context.Context, id int) (*Account, error)
	Add(ctx context.Context, account *Account) error
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id int) error

	ListAccount(ctx context.Context) []*app.Account
	OneAccount(ctx context.Context, id int) (*app.Account, error)

	ListAccountLink(ctx context.Context) []*app.AccountLink
	OneAccountLink(ctx context.Context, id int) (*app.AccountLink, error)

	ListAccountTiny(ctx context.Context) []*app.AccountTiny
	OneAccountTiny(ctx context.Context, id int) (*app.AccountTiny, error)

	UpdateFromAccountPayload(ctx context.Context, payload *app.AccountPayload, id int) error
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m *AccountDB) TableName() string {
	return "accounts"

}

// CRUD Functions

// Get returns a single Account as a Database Model
// This is more for use internally, and probably not what you want in  your controllers
func (m *AccountDB) Get(ctx context.Context, id int) (*Account, error) {
	defer goa.MeasureSince([]string{"goa", "db", "account", "get"}, time.Now())

	var native Account
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of Account
func (m *AccountDB) List(ctx context.Context) ([]*Account, error) {
	defer goa.MeasureSince([]string{"goa", "db", "account", "list"}, time.Now())

	var objs []*Account
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *AccountDB) Add(ctx context.Context, model *Account) error {
	defer goa.MeasureSince([]string{"goa", "db", "account", "add"}, time.Now())

	err := m.Db.Create(model).Error
	if err != nil {
		goa.LogError(ctx, "error adding Account", "error", err.Error())
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *AccountDB) Update(ctx context.Context, model *Account) error {
	defer goa.MeasureSince([]string{"goa", "db", "account", "update"}, time.Now())

	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		goa.LogError(ctx, "error updating Account", "error", err.Error())
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *AccountDB) Delete(ctx context.Context, id int) error {
	defer goa.MeasureSince([]string{"goa", "db", "account", "delete"}, time.Now())

	var obj Account

	err := m.Db.Delete(&obj, id).Error

	if err != nil {
		goa.LogError(ctx, "error deleting Account", "error", err.Error())
		return err
	}

	return nil
}

// AccountFromAccountPayload Converts source AccountPayload to target Account model
// only copying the non-nil fields from the source.
func AccountFromAccountPayload(payload *app.AccountPayload) *Account {
	account := &Account{}
	account.LeaseExpiresIn = payload.LeaseExpiresIn
	account.LeaseExpiresInUnits = payload.LeaseExpiresInUnits
	account.Name = payload.Name

	return account
}

// UpdateFromAccountPayload applies non-nil changes from AccountPayload to the model and saves it
func (m *AccountDB) UpdateFromAccountPayload(ctx context.Context, payload *app.AccountPayload, id int) error {
	defer goa.MeasureSince([]string{"goa", "db", "account", "updatefromaccountPayload"}, time.Now())

	var obj Account
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&obj).Error
	if err != nil {
		goa.LogError(ctx, "error retrieving Account", "error", err.Error())
		return err
	}
	obj.LeaseExpiresIn = payload.LeaseExpiresIn
	obj.LeaseExpiresInUnits = payload.LeaseExpiresInUnits
	obj.Name = payload.Name

	err = m.Db.Save(&obj).Error
	return err
}
