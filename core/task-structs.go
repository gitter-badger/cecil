package core

import ()

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
}

type NotifierTask struct {
	From     string
	To       string
	Subject  string
	BodyHTML string
	BodyText string
}
