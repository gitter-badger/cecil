// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package tasks

import (
	"time"

	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/notification"
)

// @@@@@@@@@@@@@@@ Notifier task @@@@@@@@@@@@@@@

type NewInstanceTask struct {
	Transmission interface{}
}

func (t *NewInstanceTask) Validate() error {
	// TODO: add validation
	return nil
}

// @@@@@@@@@@@@@@@ Terminator task @@@@@@@@@@@@@@@

type TerminatorTask struct {
	models.Lease

	Action string // default is TerminatorActionTerminate
}

func (t *TerminatorTask) Validate() error {
	// TODO: add validation
	return nil
}

// @@@@@@@@@@@@@@@ Lease terminated task @@@@@@@@@@@@@@@

type InstanceTerminatedTask struct {
	AWSID         string
	AWSResourceID string
	ResourceType  string
	Transmission  interface{}
	TerminatedAt  time.Time
}

func (t *InstanceTerminatedTask) Validate() error {
	// TODO: add validation
	return nil
}

// @@@@@@@@@@@@@@@ Extender task @@@@@@@@@@@@@@@

type ExtenderTask struct {
	models.Lease

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
	NotificationMeta notification.NotificationMeta

	DeliverAfter time.Duration
}

func (t *NotifierTask) Validate() error {
	// TODO: add validation
	return nil
}
