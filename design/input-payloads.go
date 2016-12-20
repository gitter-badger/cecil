package design

// input-payloads.go contains structures that are received FROM the user
// by the API.

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AccountInputPayload = Type("AccountInputPayload", func() {
	Attribute("email", func() {
		Format("email")
	})
	Attribute("name", func() {
		MinLength(1)
		MaxLength(30)
	})
	Attribute("surname", func() {
		MinLength(1)
		MaxLength(30)
	})
})

var AccountVerificationInputPayload = Type("AccountVerificationInputPayload", func() {
	Attribute("verification_token", func() {
		MinLength(108) // it is 3 consecutive UUIDs (each long 36 characters)
	})
})

var CloudAccountInputPayload = Type("CloudAccountInputPayload", func() {
	Attribute("aws_id", func() {
		MinLength(1)
	})
	Attribute("default_lease_duration")
})

var OwnerInputPayload = Type("OwnerInputPayload", func() {
	Attribute("email", func() {
		Format("email")
	})
})

var SlackConfigInputPayload = Type("SlackConfigInputPayload", func() {
	Attribute("token", func() {
		MinLength(1)
	})
	Attribute("channel_id", func() {
		MinLength(1)
	})
})

var MailerConfigInputPayload = Type("MailerConfigInputPayload", func() {
	Attribute("domain", func() {
		MinLength(1)
	})
	Attribute("api_key", func() {
		MinLength(1)
	})
	Attribute("public_api_key", func() {
		MinLength(1)
	})
	Attribute("from_name", func() {
		MinLength(1)
	})
})

var SubscribeSNSToSQSInputPayload = Type("SubscribeSNSToSQSInputPayload", func() {
	Attribute("regions", ArrayOf(String))
})
