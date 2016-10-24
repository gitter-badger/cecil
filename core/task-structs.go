package core

import (
	"time"
)

// @@@@@@@@@@@@@@@ Notifier task @@@@@@@@@@@@@@@

type NewLeaseTask struct {
	*Transmission
}

func (t *NewLeaseTask) Validate() error {

	return nil
}

// @@@@@@@@@@@@@@@ Terminator task @@@@@@@@@@@@@@@

type TerminatorTask struct {
	Lease

	Action string // default is TerminatorActionTerminate
}

func (t *TerminatorTask) Validate() error {

	return nil
}

// @@@@@@@@@@@@@@@ Lease terminated task @@@@@@@@@@@@@@@

type LeaseTerminatedTask struct {
	AWSID        string
	InstanceID   string
	TerminatedAt time.Time
}

func (t *LeaseTerminatedTask) Validate() error {

	return nil
}

// @@@@@@@@@@@@@@@ Extender task @@@@@@@@@@@@@@@

type ExtenderTask struct {
	Lease

	ExtendBy  time.Duration
	Approving bool
}

func (t *ExtenderTask) Validate() error {

	return nil
}

// @@@@@@@@@@@@@@@ Notifier task @@@@@@@@@@@@@@@

type NotifierTask struct {
	From     string
	To       string
	Subject  string
	BodyHTML string
	BodyText string
}

func (t *NotifierTask) Validate() error {

	return nil
}
