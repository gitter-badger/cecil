// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package notification

import "fmt"

const (
	X_CECIL_MESSAGETYPE        = "X-Cecil-MessageType"
	X_CECIL_LEASE_UUID         = "X-Cecil-LeaseUUID"
	X_CECIL_AWS_RESOURCE_ID    = "X-Cecil-AWS-ResourceID"
	X_CECIL_VERIFICATION_TOKEN = "X-Cecil-Verification-Token"
)

type NotificationMeta struct {
	NotificationType  NotificationType
	LeaseUUID         string
	AWSResourceID     string
	ResourceType      string
	VerificationToken string
}

type NotificationType int

const (
	Misconfiguration NotificationType = iota
	InstanceNeedsAttention
	InstanceNeedsApproval
	InstanceCreated
	InstanceWillExpire
	InstanceTerminated
	LeaseApproved
	LeaseExtended
	RegionSetup
	VerifyingAccount
)

func (nt NotificationType) String() string {
	switch nt {
	case Misconfiguration:
		return "Misconfiguration"
	case InstanceNeedsAttention:
		return "InstanceNeedsAttention"
	case InstanceNeedsApproval:
		return "InstanceNeedsApproval"
	case InstanceCreated:
		return "InstanceCreated"
	case InstanceWillExpire:
		return "InstanceWillExpire"
	case InstanceTerminated:
		return "InstanceTerminated"
	case LeaseApproved:
		return "LeaseApproved"
	case LeaseExtended:
		return "LeaseExtended"
	case RegionSetup:
		return "RegionSetup"
	case VerifyingAccount:
		return "VerifyingAccount"

	}
	return "Error"
}

func NotificationTypeFromString(notificationType string) NotificationType {

	switch notificationType {
	case "Misconfiguration":
		return Misconfiguration
	case "InstanceNeedsAttention":
		return InstanceNeedsAttention
	case "InstanceNeedsApproval":
		return InstanceNeedsApproval
	case "InstanceCreated":
		return InstanceCreated
	case "InstanceWillExpire":
		return InstanceWillExpire
	case "InstanceTerminated":
		return InstanceTerminated
	case "LeaseApproved":
		return LeaseApproved
	case "LeaseExtended":
		return LeaseExtended
	case "RegionSetup":
		return RegionSetup
	case "VerifyingAccount":
		return VerifyingAccount

	}
	panic(fmt.Sprintf("Unknown notification type: %v", notificationType))

}
