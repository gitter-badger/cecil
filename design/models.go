package design

import (
	"github.com/goadesign/gorma"
	. "github.com/goadesign/gorma/dsl"
)

var _ = StorageGroup("ZeroCloud", func() {
	Description("This is the global storage group")
	Store("postgres", gorma.Postgres, func() {
		Description("This is the Postgres relational store")
		Model("Account", func() {
			RendersTo(Account)
			Description("ZeroCloud Account")
			Field("id", gorma.Integer, func() {
				PrimaryKey()
			})
			Field("name", gorma.String)
			HasMany("CloudAccounts", "CloudAccount")
			HasMany("CloudEvents", "CloudEvent")
		})
		Model("CloudAccount", func() {
			BuildsFrom(func() {
				Payload("cloudaccount", "create")
				Payload("cloudaccount", "update")
			})
			RendersTo(CloudAccount)
			Field("id", gorma.Integer, func() {
				PrimaryKey()
			})
			Field("name", gorma.String)
			Field("cloudprovider", gorma.String)
			Field("upstream_account_id", gorma.String)
			Description("CloudAccount Model")
			BelongsTo("Account")
			HasMany("CloudEvents", "CloudEvent")
		})
		Model("CloudEvent", func() {

			// I removed the BuildsFrom because it was giving errors.
			// I think it could be re-added using the MapsFrom construct
			// (see gorma/DSL docs), but for now, I'm just manually defining
			// all of the CloudEvent fields in this model definition and
			// not trying reuse the CloudEvent Payload field definitions.
			// BuildsFrom(func() {
			// 	Payload("cloudevent", "create")
			// })

			RendersTo(CloudEvent)
			Field("id", gorma.Integer, func() {
				PrimaryKey()
			})
			Field("aws_account_id", gorma.String)
			/*Field("sqs_payload", gorma.Text)
			Field("sqs_timestamp", gorma.Timestamp)
			Field("cw_event_detail_type", gorma.String)
			Field("cw_event_source", gorma.String)
			Field("cw_event_timestamp", gorma.Timestamp)
			Field("cw_event_region", gorma.String)
			Field("cw_event_detail_instance_id", gorma.String)
			Field("cw_event_detail_state", gorma.String)*/
			Description("CloudEvent Model")
			BelongsTo("Account")
			BelongsTo("CloudAccount")
		})

	})
})
