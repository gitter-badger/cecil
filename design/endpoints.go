// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("swagger", func() {
	Origin("*", func() {
		Methods("GET", "OPTIONS")
	})
	Files("/swagger.json", "goa/swagger/swagger.json")
})

var _ = Resource("root", func() {
	BasePath("/")

	Action("show", func() {

		Routing(GET(""))
		Description("Show info about API")
		Response(OK, "application/json")
	})
})

var _ = Resource("report", func() {
	BasePath("/reports")

	Parent("cloudaccount")

	Security(JWT, func() {
		Scope("api:access")
	})

	Action("orderInstancesReport", func() {
		Routing(POST("/instances"))
		Description("Order the creation of a report about instances")
		Payload(InstancesReportOrderInputPayload, func() {
			Required("minimum_lease_age")
		})
		Response(OK, "text/plain")
	})

	Action("showReport", func() {
		Description("Show a single report")
		Routing(GET("/generated/:report_uuid"))
		Params(func() {
			Param("report_uuid", UUID, "Report UUID")
		})
		Response(OK, "application/json")
	})

})

var _ = Resource("email_action", func() {
	BasePath("/email_action")

	Action("actions", func() {
		Description("Perform an action on a lease")
		Routing(GET("/leases/:lease_uuid/:group_uid_hash/:action"))
		Params(func() {
			Param("lease_uuid", UUID, "UUID of the lease")
			Param("group_uid_hash", String, "Hash of group UID")
			Param("action", String, "Action to be peformed on the lease", func() {
				Enum("approve", "terminate", "extend")
			})
			Param("tok", String, "The token_once of this link",
				func() {
					MinLength(30)
				},
			)
			Param("sig", String, "The signature of this link",
				func() {
					MinLength(30)
				})
			Required("lease_uuid", "group_uid_hash", "action", "tok", "sig")
		})
		Response(OK, "application/json")
	})
})
