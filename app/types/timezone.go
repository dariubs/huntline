package types

import (
	"time"
)

// SanFranciscoLocation returns the timezone location for San Francisco (Pacific Time)
// This is a convenience function to ensure consistent timezone usage across the application
func SanFranciscoLocation() *time.Location {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		// Fallback to UTC if timezone loading fails (should never happen)
		return time.UTC
	}
	return loc
}

