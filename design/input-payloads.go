package design

// input-payloads.go contains structures that are received FROM the user
// by the API.

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var AccountInputPayload = Type("AccountInputPayload", func() {
	Attribute("email", func() {
		Format("email")
	})
	Attribute("name", func() {
		MinLength(1)
	})
	Attribute("surname", func() {
		MinLength(1)
	})
})

var AccountVerificationInputPayload = Type("AccountVerificationInputPayload", func() {
	Attribute("verification_token", func() {
		MinLength(1)
	})
})

var CloudAccountInputPayload = Type("CloudAccountInputPayload", func() {
	Attribute("aws_id", func() {
		MinLength(1)
	})
})

var OwnerInputPayload = Type("OwnerInputPayload", func() {
	Attribute("email", func() {
		Format("email")
	})
})
