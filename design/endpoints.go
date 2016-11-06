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

	Action("activateAccount", func() {
		NoSecurity()
		Routing(GET("/:account_id/api_token"))
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

})
