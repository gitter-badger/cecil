package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

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
