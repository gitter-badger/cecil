package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

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
