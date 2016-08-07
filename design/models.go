package design

import (
	"github.com/goadesign/gorma"
	. "github.com/goadesign/gorma/dsl"
)

var _ = StorageGroup("ZeroCloud", func() {
	Description("This is the global storage group")
	Store("postgres", gorma.Postgres, func() {
		Description("This is the Postgres relational store")
		Model("Account", func() {
			RendersTo(Account)
			Description("ZeroCloud Account")
			Field("id", gorma.Integer, func() {
				PrimaryKey()
			})
			Field("name", gorma.String)
			HasMany("CloudAccounts", "CloudAccount")
		})

	})
})
