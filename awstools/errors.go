// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

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
