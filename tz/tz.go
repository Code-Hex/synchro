package tz

import "time"

// UTC represents Universal Coordinated Time (UTC).
type UTC struct{}

// Location returns time.UTC.
func (u UTC) Location() *time.Location { return time.UTC }

// Local represents the system's local time zone.
// On Unix systems, Local consults the TZ environment
// variable to find the time zone to use. No TZ means
// use the system default /etc/localtime.
type Local struct{}

// Location returns time.Local.
func (u Local) Location() *time.Location { return time.Local }
