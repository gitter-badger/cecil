package models

import (
	"fmt"
	"log"
	"time"
)

// Using the Account settings of the Account associated with this lease,
// figure out the lease expiry time (eg, three days from now) and set it
func (l *Lease) SetExpiryTime(LeaseExpiresIn int, LeaseExpiresInUnits string) error {

	log.Printf("SetExpiry called.  LeaseExpiresIn: %v LeaseExpiresInUnits: %v", LeaseExpiresIn, LeaseExpiresInUnits)

	var durationMultiplier time.Duration
	switch LeaseExpiresInUnits {
	case "hours":
		log.Printf("SetExpiryTime: hours")
		durationMultiplier = time.Hour
	case "days":
		log.Printf("SetExpiryTime: days")
		durationMultiplier = 24 * time.Hour
	default:
		return fmt.Errorf("lease.SetExpiryTime called with invalid LeaseExpiresInUnits: %s", LeaseExpiresInUnits)
	}

	deltaDuration := time.Duration(LeaseExpiresIn) * durationMultiplier
	log.Printf("SetExpiryTime deltaDuration: %v", deltaDuration)

	// TODO: it's questionable whether it should be using time.Now() here or
	// using the time that the instance was created.  If the latter is used though,
	// probably want to avoid creating leases that are already expired.
	l.Expires = time.Now().Add(deltaDuration)
	log.Printf("SetExpiryTimeFromAccountSettings expires: %v", l.Expires)

	return nil
}
