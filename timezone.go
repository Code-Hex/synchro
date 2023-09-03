package synchro

import (
	"time"
)

// TimeZone represents the timezone.
type TimeZone interface {
	Location() *time.Location
}
