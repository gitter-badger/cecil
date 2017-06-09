// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("Cecil", func() {
	Title("Cecil APIs")
	Description("")

	Version("0.1")

	License(func() {
		Name("Apache 2.0")
		URL("http://www.apache.org/licenses/LICENSE-2.0.html")
	})
	Docs(func() {
		Description("Cecil APIs docs")
		URL("")
	})
	Host("127.0.0.1:8080")

	Scheme("http", "https")

	Consumes("application/json")
	Produces("application/json")
	NoExample()
})
