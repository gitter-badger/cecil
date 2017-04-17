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

var _ = Resource("account", func() {
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

	Action("new_api_token", func() {
		NoSecurity()

		Routing(POST("/:account_id/new_api_token"))
		Params(func() {
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
		})
		Description("Create new API token")
		Payload(NewAPITokenInputPayload, func() {
			Required("email")
		})
		Response(OK, "application/json")
	})

	Action("verify", func() {
		NoSecurity()

		Routing(POST("/:account_id/api_token"))
		Params(func() {
			Param("account_id", Integer, "Account ID",
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
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "application/json")
	})

	Action("slackConfig", func() {
		Description("Configure slack")
		Routing(POST("/:account_id/slack_config"))
		Params(func() {
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
		})
		Payload(SlackConfigInputPayload, func() {
			Required("token", "channel_id")
		})
		Response(OK, "application/json")
	})

	Action("removeSlack", func() {
		Description("Remove slack")
		Routing(DELETE("/:account_id/slack_config"))
		Params(func() {
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "application/json")
	})

	Action("mailerConfig", func() {
		Description("Configure custom mailer")
		Routing(POST("/:account_id/mailer_config"))
		Params(func() {
			Param("account_id", Integer, "Account ID",
				func() {
					Minimum(1)
				},
			)
		})
		Payload(MailerConfigInputPayload, func() {
			Required("domain", "api_key", "public_api_key", "from_name")
		})
		Response(OK, "application/json")
	})

	Action("removeMailer", func() {
		Description("Remove custom mailer")
		Routing(DELETE("/:account_id/mailer_config"))
		Params(func() {
			Param("account_id", Integer, "Account ID",
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

	Action("show", func() {
		Description("Show cloudaccount")
		Routing(GET("/:cloudaccount_id"))
		Params(func() {
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "application/json")
	})

	Action("add", func() {
		Routing(POST(""))
		Description("Add new cloudaccount")
		Payload(CloudaccountInputPayload, func() {
			Required("aws_id")
		})
		Response(OK, "application/json")
	})

	Action("update", func() {
		Description("Update a cloudaccount")
		Routing(PATCH("/:cloudaccount_id"))
		Params(func() {
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
		})
		Payload(CloudaccountInputPayload, func() {
			Required("default_lease_duration")
		})
		Response(OK, "application/json")
	})

	Action("listWhitelistedOwners", func() {
		Description("List whitelisted owners")
		Routing(GET("/:cloudaccount_id/owners"))
		Params(func() {
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "application/json")
	})

	Action("addWhitelistedOwner", func() {
		Description("Add new email (plus optional KeyName) to owner tag whitelist")
		Routing(POST("/:cloudaccount_id/owners"))
		Params(func() {
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
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

	Action("updateWhitelistedOwner", func() {
		Description("Modify a whitelisted owner")
		Routing(PATCH("/:cloudaccount_id/owners"))
		Params(func() {
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
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

	Action("deleteWhitelistedOwner", func() {
		Description("Delete a whitelisted owner")
		Routing(DELETE("/:cloudaccount_id/owners"))
		Params(func() {
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
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
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
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
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "text/plain")
	})

	Action("listRegions", func() {
		Description("List all regions and their status")
		Routing(GET("/:cloudaccount_id/regions"))
		Params(func() {
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
		})
		Response(OK, "application/json")
	})

	Action("subscribeSNSToSQS", func() {
		Description("Subscribe SNS to SQS")
		Routing(POST("/:cloudaccount_id/subscribe-sns-to-sqs"))
		Params(func() {
			Param("cloudaccount_id", Integer, "Cloudaccount ID",
				func() {
					Minimum(1)
				},
			)
		})
		Payload(SubscribeSNSToSQSInputPayload, func() {
			Required("regions")
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
