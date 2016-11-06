package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("swagger", func() {
	Origin("*", func() {
		Methods("GET", "OPTIONS")
	})
	Files("/swagger.json", "public/swagger/swagger.json")
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
		Response(OK)
	})

	Action("verify", func() {
		NoSecurity()
		Routing(POST("/:account_id/api_token"))
		Params(func() {
			Param("account_id", Integer, "Account Id")
		})
		Payload(AccountVerificationInputPayload, func() {
			Required("verification_token")
		})
		Description("Verify account and get API token")
		Response(OK)
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
		Response(OK)
	})

	Action("addEmailToWhitelist", func() {
		Routing(POST("/:cloudaccount_id/owners"))
		Description("Add new email to owner tag whitelist")
		Payload(OwnerInputPayload, func() {
			Required("email")
		})
		Response(OK)
	})

	Action("downloadInitialSetupTemplate", func() {
		Routing(GET("/:cloudaccount_id/cecil-aws-initial-setup.template"))
		Description("Download AWS initial setup cloudformation template")
		Response(OK)
	})

	Action("downloadRegionSetupTemplate", func() {
		Routing(GET("/:cloudaccount_id/cecil-aws-region-setup.template"))
		Description("Download AWS region setup cloudformation template")
		Response(OK)
	})

})
