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

// ListCloudevent returns an array of view: default.
func (m *CloudEventDB) ListCloudevent(ctx context.Context, accountID int, cloudAccountID int) []*app.Cloudevent {
	defer goa.MeasureSince([]string{"goa", "db", "cloudevent", "listcloudevent"}, time.Now())

	var native []*CloudEvent
	var objs []*app.Cloudevent
	err := m.Db.Scopes(CloudEventFilterByAccount(accountID, m.Db), CloudEventFilterByCloudAccount(cloudAccountID, m.Db)).Table(m.TableName()).Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing CloudEvent", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.CloudEventToCloudevent())
	}

	return objs
}

// CloudEventToCloudevent loads a CloudEvent and builds the default view of media type Cloudevent.
func (m *CloudEvent) CloudEventToCloudevent() *app.Cloudevent {
	cloudevent := &app.Cloudevent{}
	cloudevent.AwsAccountID = m.AwsAccountID
	cloudevent.ID = m.ID

	return cloudevent
}

// OneCloudevent loads a CloudEvent and builds the default view of media type Cloudevent.
func (m *CloudEventDB) OneCloudevent(ctx context.Context, id int, accountID int, cloudAccountID int) (*app.Cloudevent, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudevent", "onecloudevent"}, time.Now())

	var native CloudEvent
	err := m.Db.Scopes(CloudEventFilterByAccount(accountID, m.Db), CloudEventFilterByCloudAccount(cloudAccountID, m.Db)).Table(m.TableName()).Preload("Account").Preload("CloudAccount").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting CloudEvent", "error", err.Error())
		return nil, err
	}

	view := *native.CloudEventToCloudevent()
	return &view, err
}

// MediaType Retrieval Functions

// ListCloudeventTiny returns an array of view: tiny.
func (m *CloudEventDB) ListCloudeventTiny(ctx context.Context, accountID int, cloudAccountID int) []*app.CloudeventTiny {
	defer goa.MeasureSince([]string{"goa", "db", "cloudevent", "listcloudeventtiny"}, time.Now())

	var native []*CloudEvent
	var objs []*app.CloudeventTiny
	err := m.Db.Scopes(CloudEventFilterByAccount(accountID, m.Db), CloudEventFilterByCloudAccount(cloudAccountID, m.Db)).Table(m.TableName()).Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing CloudEvent", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.CloudEventToCloudeventTiny())
	}

	return objs
}

// CloudEventToCloudeventTiny loads a CloudEvent and builds the tiny view of media type Cloudevent.
func (m *CloudEvent) CloudEventToCloudeventTiny() *app.CloudeventTiny {
	cloudevent := &app.CloudeventTiny{}
	cloudevent.ID = m.ID

	return cloudevent
}

// OneCloudeventTiny loads a CloudEvent and builds the tiny view of media type Cloudevent.
func (m *CloudEventDB) OneCloudeventTiny(ctx context.Context, id int, accountID int, cloudAccountID int) (*app.CloudeventTiny, error) {
	defer goa.MeasureSince([]string{"goa", "db", "cloudevent", "onecloudeventtiny"}, time.Now())

	var native CloudEvent
	err := m.Db.Scopes(CloudEventFilterByAccount(accountID, m.Db), CloudEventFilterByCloudAccount(cloudAccountID, m.Db)).Table(m.TableName()).Preload("Account").Preload("CloudAccount").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting CloudEvent", "error", err.Error())
		return nil, err
	}

	view := *native.CloudEventToCloudeventTiny()
	return &view, err
}
