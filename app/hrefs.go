//************************************************************************//
// API "zerocloud": Application Resource Href Factories
//
// Generated with goagen v1.0.0, command line:
// $ goagen
// --design=github.com/tleyden/zerocloud/design
// --out=$(GOPATH)/src/github.com/tleyden/zerocloud
// --version=v1.0.0
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package app

import "fmt"

// AccountHref returns the resource href.
func AccountHref(accountID interface{}) string {
	return fmt.Sprintf("/accounts/%v", accountID)
}

// AwsHref returns the resource href.
func AwsHref(awsAccountID interface{}) string {
	return fmt.Sprintf("/aws/%v", awsAccountID)
}

// CloudaccountHref returns the resource href.
func CloudaccountHref(accountID, cloudAccountID interface{}) string {
	return fmt.Sprintf("/accounts/%v/cloudaccounts/%v", accountID, cloudAccountID)
}
