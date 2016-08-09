package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("zerocloud", func() {
	Title("The ZeroCloud API")
	Description("An API definition for ZeroCloud in Goa")
	Host("localhost:8080")
	Scheme("http")

	ResponseTemplate(Created, func(pattern string) {
		Description("Resource created")
		Status(201)
		Headers(func() {
			Header("Location", String, "href to created resource", func() {
				Pattern(pattern)
			})
		})
	})
})

var _ = Resource("account", func() {

	DefaultMedia(Account)
	BasePath("/accounts")

	Action("list", func() {
		Routing(
			GET(""),
		)
		Description("Retrieve all accounts.")
		Response(OK, CollectionOf(Account))
		Response(NotFound)
	})

	Action("show", func() {
		Routing(
			GET("/:accountID"),
		)
		Description("Retrieve account with given id. IDs 1 and 2 pre-exist in the system.")
		Params(func() {
			Param("accountID", Integer, "Account ID")
		})
		Response(OK)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("create", func() {
		Routing(
			POST(""),
		)
		Description("Create new account")
		Payload(func() {
			Member("name")
			Required("name")
		})
		Response(Created, "/accounts/[0-9]+")
		Response(BadRequest, ErrorMedia)
	})

	Action("update", func() {
		Routing(
			PUT("/:accountID"),
		)
		Description("Change account name")
		Params(func() {
			Param("accountID", Integer, "Account ID")
		})
		Payload(func() {
			Member("name")
			Required("name")
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("delete", func() {
		Routing(
			DELETE("/:accountID"),
		)
		Params(func() {
			Param("accountID", Integer, "Account ID")
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

// Account is the account resource media type.
var Account = MediaType("application/vnd.account+json", func() {
	Description("A tenant account")
	Attributes(func() {
		Attribute("id", Integer, "ID of account", func() {
			Example(1)
		})
		Attribute("href", String, "API href of account", func() {
			Example("/accounts/1")
		})
		Attribute("name", String, "Name of account", func() {
			Example("test")
		})
		Attribute("created_at", DateTime, "Date of creation")
		Attribute("created_by", String, "Email of account owner", func() {
			Format("email")
			Example("admin@bigdb.co")
		})

		Required("id", "href", "name", "created_at", "created_by")
	})

	View("default", func() {
		Attribute("id")
		Attribute("href")
		Attribute("name")
		Attribute("created_at")
		Attribute("created_by")
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("href")
		Attribute("name")
	})
	View("link", func() {
		Attribute("id")
		Attribute("href")
	})
})

var _ = Resource("cloudaccount", func() {

	DefaultMedia(CloudAccount)
	BasePath("cloudaccounts")
	Parent("account")

	Action("list", func() {
		Routing(
			GET(""),
		)
		Description("List all cloud accounts")
		Response(OK, func() {
			Media(CollectionOf(CloudAccount, func() {
				View("default")
				View("tiny")
			}))
		})
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("show", func() {
		Routing(
			GET("/:cloudAccountID"),
		)
		Description("Retrieve cloud account with given id")
		Params(func() {
			Param("cloudAccountID", Integer)
		})
		Response(OK)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("create", func() {
		Routing(
			POST(""),
		)
		Description("Record new cloud account")
		Payload(CloudAccountPayload)
		Response(Created, "^/accounts/[0-9]+/cloudaccounts/[0-9]+$")
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("update", func() {
		Routing(
			PATCH("/:cloudAccountID"),
		)
		Params(func() {
			Param("cloudAccountID", Integer)
		})
		Payload(CloudAccountPayload)
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("delete", func() {
		Routing(
			DELETE("/:cloudAccountID"),
		)
		Params(func() {
			Param("cloudAccountID", Integer)
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

// CloudAccountPayload defines the data structure used in the create CloudAccount request body.
// It is also the base type for the CloudAccount media type used to render CloudAccounts.
var CloudAccountPayload = Type("CloudAccountPayload", func() {
	Attribute("name", String, "Name of account", func() {
		MinLength(3)
		Example("BigDB.co's perf testing AWS account")
	})
	Attribute("cloudprovider", func() {
		MinLength(3)
		Example("AWS")
	})
	Attribute("upstream_account_id", func() {
		MinLength(4)
		Example("98798079879")
	})
})

// CloudAccount is the CloudAccount resource media type.
var CloudAccount = MediaType("application/vnd.cloudaccount+json", func() {
	Description("A CloudAccount")
	Reference(CloudAccountPayload)
	Attributes(func() {
		Attribute("id", Integer, "ID of cloud account", func() {
			Example(1)
		})
		Attribute("href", String, "API href of cloud account", func() {
			Example("/accounts/1/cloudaccounts/1")
		})
		Attribute("account", Account, "Account that owns CloudAccount")
		Attribute("created_at", DateTime, "Date of creation")
		Attribute("updated_at", DateTime, "Date of last update")
		// Attributes below inherit from the base type
		Attribute("name")
		Attribute("cloudprovider")
		Attribute("upstream_account_id")

		Required("id", "href", "name", "cloudprovider", "upstream_account_id")
		Required("created_at")
	})

	Links(func() {
		Link("account")
	})

	View("default", func() {
		Attribute("id")
		Attribute("href")
		Attribute("name")
		Attribute("cloudprovider")
		Attribute("upstream_account_id")
		Attribute("account", func() {
			View("tiny")
		})
		Attribute("links")
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("href")
		Attribute("name")
		Attribute("links")
	})

	View("link", func() {
		Attribute("id")
		Attribute("href")
	})

})

// Injest AWS CloudWatch Events -- these will be pushed from customer AWS
// account into an SQS queue owned by the ZeroCloud AWS account, and there
// will be a separate process which pulls from SQS, enhnances with instance
// tags and possibly other metadata, and then calls this endpoint
var _ = Resource("cloudevent", func() {

	DefaultMedia(CloudEvent)
	BasePath("/cloudevent")

	Action("create", func() {
		Routing(
			POST(""),
		)
		Description("Save a new AWS CloudWatch event")
		Payload(CloudEventPayload)
		Response(Created, "") // What should arg to "Created" be??
		Response(BadRequest, ErrorMedia)
	})

})

// CloudEventPayload defines the data structure used in the create CloudEvent request body.
// It is also the base type for the CloudEvent media type used to render CloudEvents.
var CloudEventPayload = Type("CloudEventPayload", func() {
	Attribute("aws_account_id", func() {
		MinLength(4)
		Example("98798079879")
	})
})

// CloudEvent is the CloudEvent resource media type.
var CloudEvent = MediaType("application/vnd.cloudevent+json", func() {
	Description("A CloudEvent -- AWS CloudWatch Event")
	Reference(CloudEventPayload)
	Attributes(func() {
		Attribute("id", Integer, "ID of cloud event", func() {
			Example(1)
		})
		Attribute("href", String, "API href of cloud event", func() {
			Example("/cloudevents/1")
		})
		// Attribute("cloudaccount", CloudAccount, "CloudAccount that owns CloudEvent")
		Attribute("account", Account, "Account that owns CloudEvent")
		Attribute("created_at", DateTime, "Date of creation")
		Attribute("updated_at", DateTime, "Date of last update")
		// Attributes below inherit from the base type
		Attribute("aws_account_id")

		Required("id", "href", "aws_account_id")
		Required("created_at")
	})

	Links(func() {
		Link("account")
		// 	Link("cloudaccount")
	})

	View("default", func() {
		Attribute("id")
		Attribute("href")
		Attribute("aws_account_id")
		Attribute("account", func() {
			View("tiny")
		})
		// Attribute("cloudaccount", func() {
		// 	View("tiny")
		// })
		Attribute("links")
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("href")
		Attribute("links")
	})

})
