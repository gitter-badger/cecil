package models

import "time"

// Using the Account settings of the Account associated with this lease,
// figure out the lease expiry time (eg, three days from now) and set it
func (l *Lease) SetExpiryTimeFromAccountSettings() error {

	var durationMultiplier time.Duration
	switch l.Account.LeaseExpiresInUnits {
	case "hours":
		durationMultiplier = time.Hour
	case "days":
		durationMultiplier = 24 * time.Hour
	}

	deltaDuration := time.Duration(l.Account.LeaseExpiresIn) * durationMultiplier

	// TODO: it's questionable whether it should be using time.Now() here or
	// using the time that the instance was created.  If the latter is used though,
	// probably want to avoid creating leases that are already expired.
	l.Expires = time.Now().Add(deltaDuration)

	return nil
}
