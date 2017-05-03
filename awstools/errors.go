package awstools

import "strings"

// IsErrNotFoundASG returns true if the error tells that the ASG was not found
func IsErrNotFoundASG(err error) bool {
	return strings.Contains(err.Error(), "name not found")
}

// IsErrNotFoundInstance returns true if the error tells that the instance was not found
func IsErrNotFoundInstance(err error) bool {
	return strings.Contains(err.Error(), "InvalidInstanceID.NotFound")
}
