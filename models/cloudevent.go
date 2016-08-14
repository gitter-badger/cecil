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

// CloudEvent Model
type CloudEvent struct {
	ID                      int `gorm:"primary_key"` // primary key
	AccountID               int // Belongs To Account
	AwsAccountID            string
	CloudAccountID          int // Belongs To CloudAccount
	CreatedAt               time.Time
	CwEventDetailInstanceID string
	CwEventDetailState      string
	CwEventRegion           string
	CwEventSource           string
	DeletedAt               *time.Time
	Leases                  []Lease // has many Leases
	SqsPayloadBase64        string
	UpdatedAt               time.Time
	CwEventTimestamp        time.Time // timestamp
	Account                 Account
	CloudAccount            CloudAccount
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m CloudEvent) TableName() string {
	return "cloud_events"

}

// CloudEventDB is the implementation of the storage interface for
// CloudEvent.
type CloudEventDB struct {
	Db *gorm.DB
}

// NewCloudEventDB creates a new storage type.
func NewCloudEventDB(db *gorm.DB) *CloudEventDB {
	return &CloudEventDB{Db: db}
}

// DB returns the underlying database.
func (m *CloudEventDB) DB() interface{} {
	return m.Db
}

// CloudEventStorage represents the storage interface.
type CloudEventStorage interface {
	DB() interface{}
	List(ctx context.Context) ([]*CloudEvent, error)
	Get(ctx context.Context, id int) (*CloudEvent, error)
	Add(ctx context.Context, cloudevent *CloudEvent) error
	Update(ctx context.Context, cloudevent *CloudEvent) error
	Delete(ctx context.Context, id int) error

	ListCloudevent(ctx context.Context, accountID int, cloudAccountID int) []*app.Cloudevent
	OneCloudevent(ctx context.Context, id int, accountID int, cloudAccountID int) (*app.Cloudevent, error)

	ListCloudeventTiny(ctx context.Context, accountID int, cloudAccountID int) []*app.CloudeventTiny
	OneCloudeventTiny(ctx context.Context, id int, accountID int, cloudAccountID int) (*app.CloudeventTiny, error)
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m *CloudEventDB) TableName() string {
	return "cloud_events"

}

// Belongs To Relationships

// CloudEventFilterByAccount is a gorm filter for a Belongs To relationship.
func CloudEventFilterByAccount(accountID int, originaldb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if accountID > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("account_id = ?", accountID)

		}
	}
	return func(db *gorm.DB) *gorm.DB { return db }
}

// Belongs To Relationships

// CloudEventFilterByCloudAccount is a gorm filter for a Belongs To relationship.
func CloudEventFilterByCloudAccount(cloudAccountID int, originaldb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if cloudAccountID > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("cloud_account_id = ?", cloudAccountID)

		}
	}
	return func(db *gorm.DB) *gorm.DB { return db }
}

// CRUD Functions

// Get returns a single CloudEvent as a Database Model
// This is more for use internally, and probably not what you want in  your controllers
func (m *CloudEventDB) Get(ctx context.Context, id int) (*CloudEvent, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudEvent", "get"}, time.Now())

	var native CloudEvent
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of CloudEvent
func (m *CloudEventDB) List(ctx context.Context) ([]*CloudEvent, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudEvent", "list"}, time.Now())

	var objs []*CloudEvent
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *CloudEventDB) Add(ctx context.Context, model *CloudEvent) error {
	defer goa.MeasureSince([]string{"goa", "db", "cloudEvent", "add"}, time.Now())

	err := m.Db.Create(model).Error
	if err != nil {
		goa.LogError(ctx, "error adding CloudEvent", "error", err.Error())
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *CloudEventDB) Update(ctx context.Context, model *CloudEvent) error {
	defer goa.MeasureSince([]string{"goa", "db", "cloudEvent", "update"}, time.Now())

	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		goa.LogError(ctx, "error updating CloudEvent", "error", err.Error())
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *CloudEventDB) Delete(ctx context.Context, id int) error {
	defer goa.MeasureSince([]string{"goa", "db", "cloudEvent", "delete"}, time.Now())

	var obj CloudEvent

	err := m.Db.Delete(&obj, id).Error

	if err != nil {
		goa.LogError(ctx, "error deleting CloudEvent", "error", err.Error())
		return err
	}

	return nil
}
