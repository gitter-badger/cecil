// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

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

var NewAPITokenInputPayload = Type("NewAPITokenInputPayload", func() {
	Attribute("email", func() {
		Format("email")
	})
})

var AccountVerificationInputPayload = Type("AccountVerificationInputPayload", func() {
	Attribute("verification_token", func() {
		MinLength(108) // it is 3 consecutive UUIDs (each long 36 characters)
	})
})

var CloudaccountInputPayload = Type("CloudaccountInputPayload", func() {
	Attribute("aws_id", func() {
		MinLength(1)
	})
	Attribute("default_lease_duration")
})

var OwnerInputPayload = Type("OwnerInputPayload", func() {
	Attribute("email", func() {
		Format("email")
	})
	Attribute("key_name")
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

var InstancesReportOrderInputPayload = Type("InstancesReportOrderInputPayload", func() {
	Attribute("recipients", ArrayOf(String))
	Attribute("minimum_lease_age", String)
})
