package core

import "fmt"

const (
	X_CECIL_MESSAGETYPE        = "X-Cecil-MessageType"
	X_CECIL_LEASE_UUID         = "X-Cecil-LeaseUUID"
	X_CECIL_INSTANCE_ID        = "X-Cecil-InstanceID"
	X_CECIL_VERIFICATION_TOKEN = "X-Cecil-Verification-Token"
)

type NotificationMeta struct {
	NotificationType  NotificationType
	LeaseUuid         string
	InstanceId        string
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
