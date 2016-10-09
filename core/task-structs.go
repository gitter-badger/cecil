package core

import (
	"fmt"
	"time"
)

// @@@@@@@@@@@@@@@ Task structs @@@@@@@@@@@@@@@

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

	}
	panic(fmt.Sprintf("Unknown notification type: %v", notificationType))

}

type NewLeaseTask struct {
	*Transmission
}

type TerminatorTask struct {
	Lease

	Action string // default is TerminatorActionTerminate
}

type LeaseTerminatedTask struct {
	AWSID        string
	InstanceID   string
	TerminatedAt time.Time
}

type ExtenderTask struct {
	Lease

	ExtendBy  time.Duration
	Approving bool
}

type NotifierTask struct {
	From             string
	To               string
	Subject          string
	BodyHTML         string
	BodyText         string
	NotificationType NotificationType
}
