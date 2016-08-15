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
		Payload(AccountPayload)
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
		Payload(AccountPayload)
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

// CloudAccountPayload defines the data structure used in the create CloudAccount request body.
// It is also the base type for the CloudAccount media type used to render CloudAccounts.
var AccountPayload = Type("AccountPayload", func() {
	Attribute("name", String, "Name of account", func() {
		MinLength(3)
		Example("BigDB")
	})
	Attribute("lease_expires_in_units", String, "The units for the lease_expires_in field", func() {
		Enum("seconds", "minutes", "hours", "days")
		Example("days")
		Default("days")
	})
	Attribute("lease_expires_in", Integer, "The lease will expire in this many lease_expires_in_units", func() {
		Example(3)
		Default(3) // defaults to 3 days
	})

	Required("name")
})

// Account is the account resource media type.
var Account = MediaType("application/vnd.account+json", func() {
	Description("An account")
	Reference(AccountPayload)
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
	Attribute("assume_role_arn", func() {
		MinLength(4)
		Example("arn:aws:iam::788612350743:role/ZeroCloud")
	})
	Attribute("assume_role_external_id", func() {
		MinLength(1)
		Example("bigdb")
	})
	Required("name", "cloudprovider", "upstream_account_id", "assume_role_arn", "assume_role_external_id")
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
		Attribute("assume_role_arn", String, "The Role ARN which allows ZeroCloud to use AssumeRole.  See https://github.com/tleyden/zerocloud/issues/1")
		Attribute("assume_role_external_id", String, "The customer and aws account specific External ID that needs to be passed when using AssumeRole.  See https://github.com/tleyden/zerocloud/issues/1")
		// Attributes below inherit from the base type
		Attribute("name")
		Attribute("cloudprovider")
		Attribute("upstream_account_id")

		Required("id", "href", "name", "cloudprovider", "upstream_account_id", "assume_role_arn", "assume_role_external_id")
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
		Attribute("assume_role_arn")
		Attribute("assume_role_external_id")
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
	Attribute("Message", func() { // Nested definition, defines a struct in Go
		Attribute("account", String, "AWS Account", func() {
			Example("868768768")
		})
		Attribute("detail", func() {
			Attribute("instance-id", String, "EC2 Instance ID", func() {
				Example("i-0a74797fd283b53de")
			})
			Attribute("state", String, "EC2 Instance State", func() {
				Example("running")
			})
			Required("instance-id", "state")
		})
		Attribute("detail-type", String, "CloudWatch Event Detail Type", func() {
			Example("EC2 Instance State-change Notification")
		})
		Attribute("id", String, "CloudWatch Event ID", func() {
			Example("2ecfc931-d9f2-4b25-9c00-87e6431d09f7")
		})
		Attribute("region", String, "CloudWatch Event Region", func() {
			Example("us-west-1")
		})
		// TODO: how do I model an array?
		// "resources":[
		//            "arn:aws:ec2:us-west-1:788612350743:instance/i-0a74797fd283b53de"
		//        ],
		Attribute("source", String, "CloudWatch Event Source", func() {
			Example("aws.ec2")
		})
		Attribute("time", DateTime, "CloudWatch Event Timestamp", func() {
			Example("2016-08-06T20:53:38Z")
		})
		Attribute("version", String, "CloudWatch Event Version", func() {
			Example("0")
		})
		Required("account")
	})
	Attribute("MessageId", String, "CloudWatch Event ID", func() {
		Example("fb7dad1a-ccee-5ac8-ac38-fd3a9c7dfe35")
	})
	Attribute("SQSPayloadBase64", String, "SQS Payload Base64", func() {
		Example("ewogICAgIkF0dHJpYnV0Z........5TlpRPT0iCn0=")
	})
	Attribute("Timestamp", DateTime, "CloudWatch Event Timestamp", func() {
		Example("2016-08-06T20:53:39.209Z")
	})
	Attribute("TopicArn", String, "CloudWatch Event Topic ARN", func() {
		Example("arn:aws:sns:us-west-1:788612350743:BigDBEC2Events")
	})
	Attribute("Type", String, "CloudWatch Event Type", func() {
		Example("Notification")
	})
	Required("Message")

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
		Attribute("created_at", DateTime, "Date of creation")
		Attribute("updated_at", DateTime, "Date of last update")
		// Attributes below inherit from the base type
		Attribute("aws_account_id")

		Required("id", "href", "aws_account_id")
		Required("created_at")
	})

	View("default", func() {
		Attribute("id")
		Attribute("href")
		Attribute("aws_account_id")
		Attribute("links")
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("href")
		Attribute("links")
	})

})

var _ = Resource("aws", func() {

	BasePath("aws")

	Action("show", func() {
		Routing(
			GET("/:awsAccountID"),
		)
		Description("Lookup the CloudAccount associated with this AWS acocunt ID")
		Response(OK, func() {
			Media(CloudAccount)
		})
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

})

// Lease is the lease resource media type.
var Lease = MediaType("application/vnd.lease+json", func() {
	Description("A lease")
	Attributes(func() {
		Attribute("id", Integer, "ID of lease", func() {
			Example(1)
		})
		Attribute("href", String, "API href of lease", func() {
			Example("/leases/1")
		})
		Attribute("expires", DateTime, "The datetime when this lease expires")
		Attribute("state", String, "The current state of the lease", func() {
			Example("Active")
		})
		Attribute("created_at", DateTime, "Date of creation")
		Attribute("updated_at", DateTime, "Date of last update")
		Attribute("account", Account, "Account that owns Lease")
		Attribute("account_id", Integer, "ID of Account that owns Lease")
		Attribute("cloud_event_id", Integer, "ID of CloudEvent that owns Lease")
		Attribute("cloud_account_id", Integer, "ID of CloudAccount that owns Lease")
		Required("id", "href", "expires", "state", "account_id", "cloud_event_id", "cloud_account_id")
	})

	Links(func() {
		Link("account")
	})

	View("default", func() {
		Attribute("id")
		Attribute("href")
		Attribute("expires")
		Attribute("state")
		Attribute("created_at")
		Attribute("account", func() {
			View("tiny")
		})
		Attribute("account_id")
		Attribute("cloud_account_id")
		Attribute("cloud_event_id")
		Attribute("links")
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("href")
		Attribute("links")
	})

	View("link", func() {
		Attribute("id")
		Attribute("href")
	})
})

var _ = Resource("lease", func() {

	DefaultMedia(Lease)
	BasePath("/leases")

	Action("list", func() {
		Routing(
			GET(""),
		)
		Description("Retrieve all leases.")
		Response(OK, CollectionOf(Lease))
		Response(NotFound)
	})

})
