package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("leases", func() {

	Security(JWT, func() {
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("listLeasesForAccount", func() {
		Description("List all leases for account")
		Routing(GET("/accounts/:account_id/leases"))
		Params(func() {
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
			Param("terminated", Boolean)
		})
		Response(OK, "application/json")
	})

	Action("listLeasesForCloudaccount", func() {
		Description("List all leases for a Cloudaccount")
		Routing(GET("/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases"))
		Params(func() {
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
			Param("terminated", Boolean)
		})
		Response(OK, "application/json")
	})

	Action("show", func() {
		Description("Show a lease")
		Routing(
			GET("/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id"),
			GET("/accounts/:account_id/leases/:lease_id"),
		)
		Params(func() {
			Param("lease_id", Integer, "Lease ID",
				func() {
					Minimum(1)
				},
			)
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "application/json")
	})

	Action("terminate", func() {
		Description("Terminate a lease")
		Routing(
			POST("/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id/terminate"),
			POST("/accounts/:account_id/leases/:lease_id/terminate"),
		)
		Params(func() {
			Param("lease_id", Integer, "Lease ID",
				func() {
					Minimum(1)
				},
			)
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "application/json")
	})

	Action("deleteFromDB", func() {
		Description("Delete a lease from DB")
		Routing(
			POST("/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id/delete"),
			POST("/accounts/:account_id/leases/:lease_id/delete"),
		)
		Params(func() {
			Param("lease_id", Integer, "Lease ID",
				func() {
					Minimum(1)
				},
			)
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "application/json")
	})

	Action("setExpiry", func() {
		Description("Set expiry of a lease")
		Routing(
			POST("/accounts/:account_id/cloudaccounts/:cloudaccount_id/leases/:lease_id/expiry"),
			POST("/accounts/:account_id/leases/:lease_id/expiry"),
		)
		Params(func() {
			Param("lease_id", Integer, "Lease ID",
				func() {
					Minimum(1)
				},
			)
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
			Param("expires_at", DateTime, "Target expiry datetime")
			Required("expires_at")
		})
		Response(OK, "application/json")
	})

})
