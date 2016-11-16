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

	//	Security(JWT, func() {
	//		Scope("api:access")
	//	})

	Action("show", func() {
		//NoSecurity()

		Routing(GET(""))
		Description("Show info about API")
		Response(OK, "application/json")
	})
})

var _ = Resource("account", func() {
	//DefaultMedia(someOutputMedia)
	BasePath("/accounts") // Gets appended to the API base path

	Security(JWT, func() {
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("create", func() {
		NoSecurity()

		Routing(POST(""))
		Description("Create new account")
		Payload(AccountInputPayload, func() {
			Required("email", "name", "surname")
		})
		Response(OK, "application/json")
	})

	Action("verify", func() {
		NoSecurity()

		Routing(POST("/:account_id/api_token"))
		Params(func() {
			Param("account_id", Integer, "Account Id",
				func() {
					Minimum(1)
				},
			)
		})
		Payload(AccountVerificationInputPayload, func() {
			Required("verification_token")
		})
		Description("Verify account and get API token")
		Response(OK, "application/json")
	})

	Action("show", func() {
		Description("Show account")
		Routing(GET("/:account_id"))
		Params(func() {
			Param("account_id", Integer, "Account Id",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "application/json")
	})

})

var _ = Resource("cloudaccount", func() {
	BasePath("/cloudaccounts")

	Parent("account")

	Security(JWT, func() {
		Scope("api:access")
	})

	Action("add", func() {
		Routing(POST(""))
		Description("Add new cloudaccount")
		Payload(CloudAccountInputPayload, func() {
			Required("aws_id")
		})
		Response(OK, "application/json")
	})

	Action("addEmailToWhitelist", func() {
		Description("Add new email to owner tag whitelist")
		Routing(POST("/:cloudaccount_id/owners"))
		Params(func() {
			Param("cloudaccount_id", Integer, "CloudAccount Id",
				func() {
					Minimum(1)
				},
			)
		})
		Payload(OwnerInputPayload, func() {
			Required("email")
		})
		Response(OK, "application/json")
	})

	Action("downloadInitialSetupTemplate", func() {
		Description("Download AWS initial setup cloudformation template")
		Routing(GET("/:cloudaccount_id/tenant-aws-initial-setup.template"))
		Params(func() {
			Param("cloudaccount_id", Integer, "CloudAccount Id",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "text/plain")
	})

	Action("downloadRegionSetupTemplate", func() {
		Description("Download AWS region setup cloudformation template")
		Routing(GET("/:cloudaccount_id/tenant-aws-region-setup.template"))
		Params(func() {
			Param("cloudaccount_id", Integer, "CloudAccount Id",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "text/plain")
	})

})

var _ = Resource("email_action", func() {
	BasePath("/email_action")

	Action("actions", func() {
		Description("Perform an action on a lease")
		Routing(GET("/leases/:lease_uuid/:instance_id/:action"))
		Params(func() {
			Param("lease_uuid", UUID, "UUID of the lease")
			Param("instance_id", String, "ID of the lease",
				func() {
					MinLength(1)
				},
			)
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
			Required("lease_uuid", "instance_id", "action", "tok", "sig")
		})
		Response(OK, "application/json")
	})
})
