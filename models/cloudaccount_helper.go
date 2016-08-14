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

// ListCloudaccount returns an array of view: default.
func (m *CloudAccountDB) ListCloudaccount(ctx context.Context, accountID int) []*app.Cloudaccount {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "listcloudaccount"}, time.Now())

	var native []*CloudAccount
	var objs []*app.Cloudaccount
	err := m.Db.Scopes(CloudAccountFilterByAccount(accountID, m.Db)).Table(m.TableName()).Preload("Account").Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing CloudAccount", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.CloudAccountToCloudaccount())
	}

	return objs
}

// CloudAccountToCloudaccount loads a CloudAccount and builds the default view of media type Cloudaccount.
func (m *CloudAccount) CloudAccountToCloudaccount() *app.Cloudaccount {
	cloudaccount := &app.Cloudaccount{}
	tmp1 := m.Account.AccountToAccountLink()
	cloudaccount.Links = &app.CloudaccountLinks{Account: tmp1}
	tmp2 := &m.Account
	cloudaccount.Account = tmp2.AccountToAccountTiny() // %!s(MISSING)
	cloudaccount.AssumeRoleArn = m.AssumeRoleArn
	cloudaccount.AssumeRoleExternalID = m.AssumeRoleExternalID
	cloudaccount.Cloudprovider = m.Cloudprovider
	cloudaccount.ID = m.ID
	cloudaccount.Name = m.Name
	cloudaccount.UpstreamAccountID = m.UpstreamAccountID

	return cloudaccount
}

// OneCloudaccount loads a CloudAccount and builds the default view of media type Cloudaccount.
func (m *CloudAccountDB) OneCloudaccount(ctx context.Context, id int, accountID int) (*app.Cloudaccount, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "onecloudaccount"}, time.Now())

	var native CloudAccount
	err := m.Db.Scopes(CloudAccountFilterByAccount(accountID, m.Db)).Table(m.TableName()).Preload("CloudEvents").Preload("Leases").Preload("Account").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting CloudAccount", "error", err.Error())
		return nil, err
	}

	view := *native.CloudAccountToCloudaccount()
	return &view, err
}

// MediaType Retrieval Functions

// ListCloudaccountLink returns an array of view: link.
func (m *CloudAccountDB) ListCloudaccountLink(ctx context.Context, accountID int) []*app.CloudaccountLink {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "listcloudaccountlink"}, time.Now())

	var native []*CloudAccount
	var objs []*app.CloudaccountLink
	err := m.Db.Scopes(CloudAccountFilterByAccount(accountID, m.Db)).Table(m.TableName()).Preload("Account").Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing CloudAccount", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.CloudAccountToCloudaccountLink())
	}

	return objs
}

// CloudAccountToCloudaccountLink loads a CloudAccount and builds the link view of media type Cloudaccount.
func (m *CloudAccount) CloudAccountToCloudaccountLink() *app.CloudaccountLink {
	cloudaccount := &app.CloudaccountLink{}
	cloudaccount.ID = m.ID

	return cloudaccount
}

// OneCloudaccountLink loads a CloudAccount and builds the link view of media type Cloudaccount.
func (m *CloudAccountDB) OneCloudaccountLink(ctx context.Context, id int, accountID int) (*app.CloudaccountLink, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "onecloudaccountlink"}, time.Now())

	var native CloudAccount
	err := m.Db.Scopes(CloudAccountFilterByAccount(accountID, m.Db)).Table(m.TableName()).Preload("CloudEvents").Preload("Leases").Preload("Account").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting CloudAccount", "error", err.Error())
		return nil, err
	}

	view := *native.CloudAccountToCloudaccountLink()
	return &view, err
}

// MediaType Retrieval Functions

// ListCloudaccountTiny returns an array of view: tiny.
func (m *CloudAccountDB) ListCloudaccountTiny(ctx context.Context, accountID int) []*app.CloudaccountTiny {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "listcloudaccounttiny"}, time.Now())

	var native []*CloudAccount
	var objs []*app.CloudaccountTiny
	err := m.Db.Scopes(CloudAccountFilterByAccount(accountID, m.Db)).Table(m.TableName()).Preload("Account").Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing CloudAccount", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.CloudAccountToCloudaccountTiny())
	}

	return objs
}

// CloudAccountToCloudaccountTiny loads a CloudAccount and builds the tiny view of media type Cloudaccount.
func (m *CloudAccount) CloudAccountToCloudaccountTiny() *app.CloudaccountTiny {
	cloudaccount := &app.CloudaccountTiny{}
	tmp1 := m.Account.AccountToAccountLink()
	cloudaccount.Links = &app.CloudaccountLinks{Account: tmp1}
	cloudaccount.ID = m.ID
	cloudaccount.Name = m.Name

	return cloudaccount
}

// OneCloudaccountTiny loads a CloudAccount and builds the tiny view of media type Cloudaccount.
func (m *CloudAccountDB) OneCloudaccountTiny(ctx context.Context, id int, accountID int) (*app.CloudaccountTiny, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "onecloudaccounttiny"}, time.Now())

	var native CloudAccount
	err := m.Db.Scopes(CloudAccountFilterByAccount(accountID, m.Db)).Table(m.TableName()).Preload("CloudEvents").Preload("Leases").Preload("Account").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting CloudAccount", "error", err.Error())
		return nil, err
	}

	view := *native.CloudAccountToCloudaccountTiny()
	return &view, err
}
