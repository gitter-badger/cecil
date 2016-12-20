package core

import "time"

// @@@@@@@@@@@@@@@ Notifier task @@@@@@@@@@@@@@@

type NewLeaseTask struct {
	*Transmission
}

func (t *NewLeaseTask) Validate() error {
	// TODO: add validation
	return nil
}

// @@@@@@@@@@@@@@@ Terminator task @@@@@@@@@@@@@@@

type TerminatorTask struct {
	Lease

	Action string // default is TerminatorActionTerminate
}

func (t *TerminatorTask) Validate() error {
	// TODO: add validation
	return nil
}

// @@@@@@@@@@@@@@@ Lease terminated task @@@@@@@@@@@@@@@

type LeaseTerminatedTask struct {
	AWSID        string
	InstanceID   string
	TerminatedAt time.Time
}

func (t *LeaseTerminatedTask) Validate() error {
	// TODO: add validation
	return nil
}

// @@@@@@@@@@@@@@@ Extender task @@@@@@@@@@@@@@@

type ExtenderTask struct {
	Lease

	Approving bool
}

func (t *ExtenderTask) Validate() error {
	// TODO: add validation
	return nil
}

// @@@@@@@@@@@@@@@ Notifier task @@@@@@@@@@@@@@@

type NotifierTask struct {
	AccountID uint

	To               string
	Subject          string
	BodyHTML         string
	BodyText         string
	NotificationMeta NotificationMeta

	DeliverAfter time.Duration
}

func (t *NotifierTask) Validate() error {
	// TODO: add validation
	return nil
}
