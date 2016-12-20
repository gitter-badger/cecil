//************************************************************************//
// API "Cecil": Application Resource Href Factories
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/cecil/design
// --out=$(GOPATH)/src/github.com/tleyden/cecil/goa
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package app

import (
	"fmt"
	"strings"
)

// AccountHref returns the resource href.
func AccountHref(accountID interface{}) string {
	paramaccountID := strings.TrimLeftFunc(fmt.Sprintf("%v", accountID), func(r rune) bool { return r == '/' })
	return fmt.Sprintf("/accounts/%v", paramaccountID)
}

// LeasesHref returns the resource href.
func LeasesHref(accountID, cloudaccountID, leaseID interface{}) string {
	paramaccountID := strings.TrimLeftFunc(fmt.Sprintf("%v", accountID), func(r rune) bool { return r == '/' })
	paramcloudaccountID := strings.TrimLeftFunc(fmt.Sprintf("%v", cloudaccountID), func(r rune) bool { return r == '/' })
	paramleaseID := strings.TrimLeftFunc(fmt.Sprintf("%v", leaseID), func(r rune) bool { return r == '/' })
	return fmt.Sprintf("/accounts/%v/cloudaccounts/%v/leases/%v", paramaccountID, paramcloudaccountID, paramleaseID)
}

// RootHref returns the resource href.
func RootHref() string {
	return "/"
}
