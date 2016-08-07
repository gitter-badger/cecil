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
func (m *CloudAccountDB) ListCloudaccount(ctx context.Context) []*app.Cloudaccount {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "listcloudaccount"}, time.Now())

	var native []*CloudAccount
	var objs []*app.Cloudaccount
	err := m.Db.Scopes().Table(m.TableName()).Preload("Account").Find(&native).Error

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
	cloudaccount.Account = m.Account.AccountToAccount()
	cloudaccount.Cloudprovider = m.Cloudprovider
	cloudaccount.ID = m.ID
	cloudaccount.Name = m.Name
	cloudaccount.UpstreamAccountID = m.UpstreamAccountID

	return cloudaccount
}

// OneCloudaccount loads a CloudAccount and builds the default view of media type Cloudaccount.
func (m *CloudAccountDB) OneCloudaccount(ctx context.Context, id int) (*app.Cloudaccount, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "onecloudaccount"}, time.Now())

	var native CloudAccount
	err := m.Db.Scopes().Table(m.TableName()).Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting CloudAccount", "error", err.Error())
		return nil, err
	}

	view := *native.CloudAccountToCloudaccount()
	return &view, err
}

// MediaType Retrieval Functions

// ListCloudaccountTiny returns an array of view: tiny.
func (m *CloudAccountDB) ListCloudaccountTiny(ctx context.Context) []*app.CloudaccountTiny {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "listcloudaccounttiny"}, time.Now())

	var native []*CloudAccount
	var objs []*app.CloudaccountTiny
	err := m.Db.Scopes().Table(m.TableName()).Preload("Account").Find(&native).Error

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
func (m *CloudAccountDB) OneCloudaccountTiny(ctx context.Context, id int) (*app.CloudaccountTiny, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudaccount", "onecloudaccounttiny"}, time.Now())

	var native CloudAccount
	err := m.Db.Scopes().Table(m.TableName()).Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting CloudAccount", "error", err.Error())
		return nil, err
	}

	view := *native.CloudAccountToCloudaccountTiny()
	return &view, err
}
