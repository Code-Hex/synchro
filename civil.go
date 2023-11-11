package synchro

import (
	"encoding"
	"fmt"
	"time"

	"github.com/Code-Hex/synchro/internal/constraints"
	"github.com/Code-Hex/synchro/iso8601"
	"github.com/Code-Hex/synchro/tz"
)

// A CivilDate represents a date (year, month, day).
type CivilDate[T TimeZone] struct {
	Year  int        // Year (e.g., 2014).
	Month time.Month // Month of the year (January = 1, ...).
	Day   int        // Day of the month, starting at 1.

	_ empty[T]
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*CivilDate[tz.UTC])(nil)

// CivilDateOf returns the CivilDate in which a time occurs in that time's location.
func CivilDateOf[T TimeZone, U Time[T] | time.Time](t U) CivilDate[T] {
	var d CivilDate[T]
	st, _ := convertTime[T](t) // No returns error if Time[T] | time.Time
	d.Year, d.Month, d.Day = st.Date()
	return d
}

// ParseCivilDate parses a string in ISO 8601 full-date format and returns the date value it represents.
func ParseCivilDate[T TimeZone, U constraints.Bytes](s U) (CivilDate[T], error) {
	t, err := iso8601.ParseDate(s)
	if err != nil {
		return CivilDate[T]{}, err
	}
	return CivilDateOf[T](t.Date().StdTime()), nil
}

// String returns the date in ISO 8601 full-date format.
func (d CivilDate[T]) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// Time converts Date[T] to Time[T].
func (d CivilDate[T]) Time() Time[T] {
	return New[T](d.Year, d.Month, d.Day, 0, 0, 0, 0)
}

// IsValid reports whether the date is valid.
func (d CivilDate[T]) IsValid() bool {
	iso := iso8601.Date{
		Year:  d.Year,
		Month: d.Month,
		Day:   d.Day,
	}
	return iso.IsValid()
}

// AddDate returns the time corresponding to adding the
// given number of years, months, and days to d.
func (d CivilDate[T]) AddDate(years int, months int, days int) CivilDate[T] {
	return CivilDateOf[T](d.Time().AddDate(years, months, days))
}

// AddDays returns the date that is n days in the future.
// n can also be negative to go into the past.
func (d CivilDate[T]) AddDays(n int) CivilDate[T] {
	return d.AddDate(0, 0, n)
}

// DaysSince returns the signed number of days between the date and s, not including the end day.
// This is the inverse operation to AddDays.
func (d CivilDate[T]) DaysSince(s CivilDate[T]) (days int) {
	// We convert to Unix time so we do not have to worry about leap seconds:
	// Unix time increases by exactly 86400 seconds per day.
	deltaUnix := d.Time().Unix() - s.Time().Unix()
	return int(deltaUnix / 86400)
}

// Before reports whether d occurs before d2.
func (d CivilDate[T]) Before(d2 CivilDate[T]) bool {
	if d.Year != d2.Year {
		return d.Year < d2.Year
	}
	if d.Month != d2.Month {
		return d.Month < d2.Month
	}
	return d.Day < d2.Day
}

// After reports whether d occurs after d2.
func (d CivilDate[T]) After(d2 CivilDate[T]) bool {
	return d2.Before(d)
}

// IsZero reports whether date fields are set to their default value.
func (d CivilDate[T]) IsZero() bool {
	return (d.Year == 0) && (int(d.Month) == 0) && (d.Day == 0)
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of d.String().
func (d CivilDate[T]) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be a string in a format accepted by ParseCivilDate.
func (d *CivilDate[T]) UnmarshalText(data []byte) error {
	var err error
	*d, err = ParseCivilDate[T](string(data))
	return err
}

// Change modifies the date based on the provided unit values.
//
// Note: Units related to time are ignored.
func (t CivilDate[T]) Change(u1 Unit, u2 ...Unit) CivilDate[T] {
	return CivilDateOf[T](t.Time().Change(u1, u2...))
}

// Advance adjusts the date based on the provided unit values, moving it forward in date.
//
// Note: Units related to time are ignored.
func (t CivilDate[T]) Advance(u1 Unit, u2 ...Unit) CivilDate[T] {
	return CivilDateOf[T](t.Time().Advance(u1, u2...))
}
