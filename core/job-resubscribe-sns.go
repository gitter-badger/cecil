// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package core

import "fmt"

func (service *Service) ResubscribeSNSToSQSJob() error {

	// Regenerate SQS permissions for all cloudaccounts in DB.
	if err := service.RegenerateSQSPermissions(); err != nil {
		return fmt.Errorf("Error calling service.RegenerateSQSPermissions(): %v", err)
	}

	// Resubscribe to all SNS topics of all cloudaccounts present in DB.
	if err := service.ResubscribeToAllSNSTopics(); err != nil {
		return fmt.Errorf("Error calling service.ResubscribeToAllSNSTopics(): %v", err)
	}

	return nil

}
