package design

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var JWT = JWTSecurity("jwt", func() {
	Header("Authorization")
	Scope("api:access", "API access") // Define "api:access" scope

	//	Scope("api:read", "Provides read access")
	//	Scope("api:write", "Provides write access")
})
