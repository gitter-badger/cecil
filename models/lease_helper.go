//************************************************************************//
// API "zerocloud": Model Helpers
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

// MediaType Retrieval Functions

// ListLease returns an array of view: default.
func (m *LeaseDB) ListLease(ctx context.Context, accountID int, cloudAccountID int, cloudEventID int) []*app.Lease {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "listlease"}, time.Now())

	var native []*Lease
	var objs []*app.Lease
	err := m.Db.Scopes(LeaseFilterByAccount(accountID, m.Db), LeaseFilterByCloudAccount(cloudAccountID, m.Db), LeaseFilterByCloudEvent(cloudEventID, m.Db)).Table(m.TableName()).Preload("Account").Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing Lease", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.LeaseToLease())
	}

	return objs
}

// LeaseToLease loads a Lease and builds the default view of media type Lease.
func (m *Lease) LeaseToLease() *app.Lease {
	lease := &app.Lease{}
	tmp1 := m.Account.AccountToAccountLink()
	lease.Links = &app.LeaseLinks{Account: tmp1}
	tmp2 := &m.Account
	lease.Account = tmp2.AccountToAccountTiny() // %!s(MISSING)
	lease.AccountID = m.AccountID
	lease.CloudAccountID = m.CloudAccountID
	lease.CloudEventID = m.CloudEventID
	lease.CreatedAt = &m.CreatedAt
	lease.Expires = m.Expires
	lease.ID = m.ID
	lease.State = m.State

	return lease
}

// OneLease loads a Lease and builds the default view of media type Lease.
func (m *LeaseDB) OneLease(ctx context.Context, id int, accountID int, cloudAccountID int, cloudEventID int) (*app.Lease, error) {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "onelease"}, time.Now())

	var native Lease
	err := m.Db.Scopes(LeaseFilterByAccount(accountID, m.Db), LeaseFilterByCloudAccount(cloudAccountID, m.Db), LeaseFilterByCloudEvent(cloudEventID, m.Db)).Table(m.TableName()).Preload("Account").Preload("CloudAccount").Preload("CloudEvent").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting Lease", "error", err.Error())
		return nil, err
	}

	view := *native.LeaseToLease()
	return &view, err
}

// MediaType Retrieval Functions

// ListLeaseLink returns an array of view: link.
func (m *LeaseDB) ListLeaseLink(ctx context.Context, accountID int, cloudAccountID int, cloudEventID int) []*app.LeaseLink {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "listleaselink"}, time.Now())

	var native []*Lease
	var objs []*app.LeaseLink
	err := m.Db.Scopes(LeaseFilterByAccount(accountID, m.Db), LeaseFilterByCloudAccount(cloudAccountID, m.Db), LeaseFilterByCloudEvent(cloudEventID, m.Db)).Table(m.TableName()).Preload("Account").Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing Lease", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.LeaseToLeaseLink())
	}

	return objs
}

// LeaseToLeaseLink loads a Lease and builds the link view of media type Lease.
func (m *Lease) LeaseToLeaseLink() *app.LeaseLink {
	lease := &app.LeaseLink{}
	lease.ID = m.ID

	return lease
}

// OneLeaseLink loads a Lease and builds the link view of media type Lease.
func (m *LeaseDB) OneLeaseLink(ctx context.Context, id int, accountID int, cloudAccountID int, cloudEventID int) (*app.LeaseLink, error) {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "oneleaselink"}, time.Now())

	var native Lease
	err := m.Db.Scopes(LeaseFilterByAccount(accountID, m.Db), LeaseFilterByCloudAccount(cloudAccountID, m.Db), LeaseFilterByCloudEvent(cloudEventID, m.Db)).Table(m.TableName()).Preload("Account").Preload("CloudAccount").Preload("CloudEvent").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting Lease", "error", err.Error())
		return nil, err
	}

	view := *native.LeaseToLeaseLink()
	return &view, err
}

// MediaType Retrieval Functions

// ListLeaseTiny returns an array of view: tiny.
func (m *LeaseDB) ListLeaseTiny(ctx context.Context, accountID int, cloudAccountID int, cloudEventID int) []*app.LeaseTiny {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "listleasetiny"}, time.Now())

	var native []*Lease
	var objs []*app.LeaseTiny
	err := m.Db.Scopes(LeaseFilterByAccount(accountID, m.Db), LeaseFilterByCloudAccount(cloudAccountID, m.Db), LeaseFilterByCloudEvent(cloudEventID, m.Db)).Table(m.TableName()).Preload("Account").Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing Lease", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.LeaseToLeaseTiny())
	}

	return objs
}

// LeaseToLeaseTiny loads a Lease and builds the tiny view of media type Lease.
func (m *Lease) LeaseToLeaseTiny() *app.LeaseTiny {
	lease := &app.LeaseTiny{}
	tmp1 := m.Account.AccountToAccountLink()
	lease.Links = &app.LeaseLinks{Account: tmp1}
	lease.ID = m.ID

	return lease
}

// OneLeaseTiny loads a Lease and builds the tiny view of media type Lease.
func (m *LeaseDB) OneLeaseTiny(ctx context.Context, id int, accountID int, cloudAccountID int, cloudEventID int) (*app.LeaseTiny, error) {
	defer goa.MeasureSince([]string{"goa", "db", "lease", "oneleasetiny"}, time.Now())

	var native Lease
	err := m.Db.Scopes(LeaseFilterByAccount(accountID, m.Db), LeaseFilterByCloudAccount(cloudAccountID, m.Db), LeaseFilterByCloudEvent(cloudEventID, m.Db)).Table(m.TableName()).Preload("Account").Preload("CloudAccount").Preload("CloudEvent").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting Lease", "error", err.Error())
		return nil, err
	}

	view := *native.LeaseToLeaseTiny()
	return &view, err
}
