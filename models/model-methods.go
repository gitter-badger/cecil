// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/tleyden/cecil/tools"
)

// SetGroupType is a helper method to set the groupType on a lease
func (ll *Lease) SetGroupType(groupType GroupType) {
	ll.GroupType = groupType
}

// SetGroupUID is a helper method to set the groupUID on a lease
func (ll *Lease) SetGroupUID(groupUID string) {
	ll.GroupUID = groupUID
}

// MarkAsTerminated is a helper method to mark a lease as terminated; this also returns a pointer
// to the lease, as a convenience for passing it to DB.Save(.)
func (ll *Lease) MarkAsTerminated(when *time.Time) *Lease {
	Logger.Info("marking LEASE as terminated",
		"lease", ll,
	)
	ll.TokenOnce = uuid.NewV4().String() // invalidates all url to renew/terminate/approve
	if when == nil {
		when = tools.TimePtr(time.Now().UTC())
	}
	ll.TerminatedAt = when
	return ll
}

// IsExpired returns true in case the lease has expired
func (ll *Lease) IsExpired() bool {
	return ll.ExpiresAt.Before(time.Now().UTC())
}

// MarkAsTerminated is a helper method to mark an Instance as terminated; this also returns a pointer
// to the lease, as a convenience for passing it to DB.Save(.)
func (ins *Instance) MarkAsTerminated(when *time.Time) *Instance {
	Logger.Info("marking INSTANCE as terminated",
		"lease", ins,
	)
	if when == nil {
		when = tools.TimePtr(time.Now().UTC())
	}
	ins.TerminatedAt = when
	return ins
}

// IsOwnerOf returns true if the account owns the cloudaccount.
func (aa *Account) IsOwnerOf(cloudaccount *Cloudaccount) bool {
	return aa.ID == cloudaccount.AccountID
}

// IsOwnerOfLease returns true if the account owns the lease
func (aa *Account) IsOwnerOfLease(lease *Lease) bool {
	return aa.ID == lease.AccountID
}

// IsOwnerOf returns true if the cloudaccount owns the lease.
func (ca *Cloudaccount) IsOwnerOf(lease *Lease) bool {
	return ca.ID == lease.CloudaccountID
}
