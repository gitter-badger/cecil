package core

import (
	"time"
)

// @@@@@@@@@@@@@@@ Task structs @@@@@@@@@@@@@@@

type NewLeaseTask struct {
	*Transmission
}

type TerminatorTask struct {
	Lease

	Action string // default is TerminatorActionTerminate
}

type LeaseTerminatedTask struct {
	AWSID      string
	InstanceID string
}

type ExtenderTask struct {
	Lease

	//TokenOnce  string
	//UUID       string
	//InstanceID string
	ExtendBy  time.Duration
	Approving bool
}

type NotifierTask struct {
	From     string
	To       string
	Subject  string
	BodyHTML string
	BodyText string
}
