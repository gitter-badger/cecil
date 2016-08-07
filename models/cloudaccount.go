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

// CloudAccount Model
type CloudAccount struct {
	ID                int     `gorm:"primary_key"` // primary key
	Account           Account // has one Account
	AccountID         int     // has many CloudAccount
	Cloudprovider     string
	CreatedAt         time.Time
	DeletedAt         *time.Time
	Name              string
	UpdatedAt         time.Time
	UpstreamAccountID string
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m CloudAccount) TableName() string {
	return "cloud_accounts"

}

// CloudAccountDB is the implementation of the storage interface for
// CloudAccount.
type CloudAccountDB struct {
	Db *gorm.DB
}

// NewCloudAccountDB creates a new storage type.
func NewCloudAccountDB(db *gorm.DB) *CloudAccountDB {
	return &CloudAccountDB{Db: db}
}

// DB returns the underlying database.
func (m *CloudAccountDB) DB() interface{} {
	return m.Db
}

// CloudAccountStorage represents the storage interface.
type CloudAccountStorage interface {
	DB() interface{}
	List(ctx context.Context) ([]*CloudAccount, error)
	Get(ctx context.Context, id int) (*CloudAccount, error)
	Add(ctx context.Context, cloudaccount *CloudAccount) error
	Update(ctx context.Context, cloudaccount *CloudAccount) error
	Delete(ctx context.Context, id int) error

	ListCloudaccount(ctx context.Context) []*app.Cloudaccount
	OneCloudaccount(ctx context.Context, id int) (*app.Cloudaccount, error)

	ListCloudaccountTiny(ctx context.Context) []*app.CloudaccountTiny
	OneCloudaccountTiny(ctx context.Context, id int) (*app.CloudaccountTiny, error)

	UpdateFromCloudAccountPayload(ctx context.Context, payload *app.CloudAccountPayload, id int) error
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m *CloudAccountDB) TableName() string {
	return "cloud_accounts"

}

// CRUD Functions

// Get returns a single CloudAccount as a Database Model
// This is more for use internally, and probably not what you want in  your controllers
func (m *CloudAccountDB) Get(ctx context.Context, id int) (*CloudAccount, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudAccount", "get"}, time.Now())

	var native CloudAccount
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of CloudAccount
func (m *CloudAccountDB) List(ctx context.Context) ([]*CloudAccount, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudAccount", "list"}, time.Now())

	var objs []*CloudAccount
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *CloudAccountDB) Add(ctx context.Context, model *CloudAccount) error {
	defer goa.MeasureSince([]string{"goa", "db", "cloudAccount", "add"}, time.Now())

	err := m.Db.Create(model).Error
	if err != nil {
		goa.LogError(ctx, "error adding CloudAccount", "error", err.Error())
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *CloudAccountDB) Update(ctx context.Context, model *CloudAccount) error {
	defer goa.MeasureSince([]string{"goa", "db", "cloudAccount", "update"}, time.Now())

	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		goa.LogError(ctx, "error updating CloudAccount", "error", err.Error())
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *CloudAccountDB) Delete(ctx context.Context, id int) error {
	defer goa.MeasureSince([]string{"goa", "db", "cloudAccount", "delete"}, time.Now())

	var obj CloudAccount

	err := m.Db.Delete(&obj, id).Error

	if err != nil {
		goa.LogError(ctx, "error deleting CloudAccount", "error", err.Error())
		return err
	}

	return nil
}

// CloudAccountFromCloudAccountPayload Converts source CloudAccountPayload to target CloudAccount model
// only copying the non-nil fields from the source.
func CloudAccountFromCloudAccountPayload(payload *app.CloudAccountPayload) *CloudAccount {
	cloudaccount := &CloudAccount{}
	if payload.Cloudprovider != nil {
		cloudaccount.Cloudprovider = *payload.Cloudprovider
	}
	if payload.Name != nil {
		cloudaccount.Name = *payload.Name
	}
	if payload.UpstreamAccountID != nil {
		cloudaccount.UpstreamAccountID = *payload.UpstreamAccountID
	}

	return cloudaccount
}

// UpdateFromCloudAccountPayload applies non-nil changes from CloudAccountPayload to the model and saves it
func (m *CloudAccountDB) UpdateFromCloudAccountPayload(ctx context.Context, payload *app.CloudAccountPayload, id int) error {
	defer goa.MeasureSince([]string{"goa", "db", "cloudAccount", "updatefromcloudAccountPayload"}, time.Now())

	var obj CloudAccount
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&obj).Error
	if err != nil {
		goa.LogError(ctx, "error retrieving CloudAccount", "error", err.Error())
		return err
	}
	if payload.Cloudprovider != nil {
		obj.Cloudprovider = *payload.Cloudprovider
	}
	if payload.Name != nil {
		obj.Name = *payload.Name
	}
	if payload.UpstreamAccountID != nil {
		obj.UpstreamAccountID = *payload.UpstreamAccountID
	}

	err = m.Db.Save(&obj).Error
	return err
}
