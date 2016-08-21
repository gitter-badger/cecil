package models

import (
	"fmt"
	"time"
)

// Using the Account settings of the Account associated with this lease,
// figure out the lease expiry time (eg, three days from now) and set it
func (l *Lease) SetExpiryTime(LeaseExpiresIn int, LeaseExpiresInUnits string) error {

	var durationMultiplier time.Duration
	switch LeaseExpiresInUnits {
	case "seconds":
		durationMultiplier = time.Second
	case "minutes":
		durationMultiplier = 60 * time.Second
	case "hours":
		durationMultiplier = time.Hour
	case "days":
		durationMultiplier = 24 * time.Hour
	default:
		return fmt.Errorf("lease.SetExpiryTime called with invalid LeaseExpiresInUnits: %s", LeaseExpiresInUnits)
	}

	deltaDuration := time.Duration(LeaseExpiresIn) * durationMultiplier

	// TODO: it's questionable whether it should be using time.Now() here or
	// using the time that the instance was created.  If the latter is used though,
	// probably want to avoid creating leases that are already expired.
	l.Expires = time.Now().Add(deltaDuration)

	return nil
}
