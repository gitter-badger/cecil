package tools

import (
	"time"

	"github.com/satori/go.uuid"
)

// IntPtr returns a pointer to the specified int
func IntPtr(i int) *int {
	return &i
}

// StringPtr returns a pointer to the specified string
func StringPtr(s string) *string {
	return &s
}

// TimePtr returns a pointer to the specified time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// DurationPtr returns a pointer to the specified duration
func DurationPtr(t time.Duration) *time.Duration {
	return &t
}

// UUIDPtr returns a pointer to the specified uuid.UUID
func UUIDPtr(t uuid.UUID) *uuid.UUID {
	return &t
}
