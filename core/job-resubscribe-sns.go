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
