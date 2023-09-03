package synchro

import (
	"encoding"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Code-Hex/synchro/tz"
)

var _ interface {
	fmt.Stringer
	fmt.GoStringer
	gob.GobEncoder
	gob.GobDecoder
	json.Marshaler
	json.Unmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
} = (*Time[tz.UTC])(nil)

// Local returns t with the location set to local time.
func (t Time[T]) Local() Time[tz.Local] {
	return In[tz.Local](t.tm)
}

// Add returns the time t+d.
//
// This is a simple wrapper method for (time.Time{}).Add.
func (t Time[T]) Add(d time.Duration) Time[T] {
	return In[T](t.tm.Add(d))
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute t-d for a duration d, use t.Add(-d).
//
// This is a simple wrapper method for (time.Time{}).Sub.
func (t Time[T]) Sub(u Time[T]) time.Duration {
	return t.tm.Sub(u.tm)
}

// AddDate returns the time corresponding to adding the
// given number of years, months, and days to t.
//
// This is a simple wrapper method for (time.Time{}).AddDate.
func (t Time[T]) AddDate(years int, months int, days int) Time[T] {
	return In[T](t.tm.AddDate(years, months, days))
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// This is a simple wrapper method for (time.Time{}).Truncate.
func (t Time[T]) Truncate(d time.Duration) Time[T] {
	return In[T](t.tm.Truncate(d))
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// This is a simple wrapper method for (time.Time{}).Round.
func (t Time[T]) Round(d time.Duration) Time[T] {
	return In[T](t.tm.Round(d))
}

// After reports whether the time instant t is after u.
//
// If you want to compare time.Time as a parameter, please use with the Time method.
func (t Time[T]) After(u Time[T]) bool {
	return t.tm.After(u.tm)
}

// Before reports whether the time instant t is before u.
//
// If you want to compare time.Time as a parameter, please use with the Time method.
func (t Time[T]) Before(u Time[T]) bool {
	return t.tm.Before(u.tm)
}

// Compare compares the time instant t with u. If t is before u, it returns -1;
// if t is after u, it returns +1; if they're the same, it returns 0.
//
// If you want to compare time.Time as a parameter, please use with the Time method.
func (t Time[T]) Compare(u Time[T]) int {
	return t.tm.Compare(u.tm)
}

// Equal reports whether t and u represent the same time instant.
//
// If you want to compare time.Time as a parameter, please use with the Time method.
func (t Time[T]) Equal(u Time[T]) bool {
	return t.tm.Equal(u.tm)
}

// String returns the time formatted using the format string
//
//	"2006-01-02 15:04:05.999999999 -0700 MST"
//
// If the time has a monotonic clock reading, the returned string
// includes a final field "m=±<value>", where value is the monotonic
// clock reading formatted as a decimal number of seconds.
//
// This is a simple wrapper method for (time.Time{}).String.
func (t Time[T]) String() string {
	return t.tm.String()
}

// GoString implements fmt.GoStringer and formats t to be printed in Go source
// code.
//
// This is a simple wrapper method for (time.Time{}).GoString.
func (t Time[T]) GoString() string {
	return t.tm.GoString()
}

// Format returns a textual representation of the time value formatted according
// to the layout defined by the argument.
//
// This is a simple wrapper method for (time.Time{}).Format.
func (t Time[T]) Format(layout string) string {
	return t.tm.Format(layout)
}

// AppendFormat is like Format but appends the textual
// representation to b and returns the extended buffer.
//
// This is a simple wrapper method for (time.Time{}).AppendFormat.
func (t Time[T]) AppendFormat(b []byte, layout string) []byte {
	return t.tm.AppendFormat(b, layout)
}

// Clock returns the hour, minute, and second within the day specified by t.
//
// This is a simple wrapper method for (time.Time{}).Clock.
func (t Time[T]) Clock() (hour, min, sec int) {
	return t.tm.Clock()
}

// Hour returns the hour within the day specified by t, in the range [0, 23].
//
// This is a simple wrapper method for (time.Time{}).Hour.
func (t Time[T]) Hour() int {
	return t.tm.Hour()
}

// Minute returns the minute offset within the hour specified by t, in the range [0, 59].
//
// This is a simple wrapper method for (time.Time{}).Minute.
func (t Time[T]) Minute() int {
	return t.tm.Minute()
}

// Second returns the second offset within the minute specified by t, in the range [0, 59].
//
// This is a simple wrapper method for (time.Time{}).Second.
func (t Time[T]) Second() int {
	return t.tm.Second()
}

// Nanosecond returns the nanosecond offset within the second specified by t,
// in the range [0, 999999999].
//
// This is a simple wrapper method for (time.Time{}).Nanosecond.
func (t Time[T]) Nanosecond() int {
	return t.tm.Nanosecond()
}

// YearDay returns the day of the year specified by t, in the range [1,365] for non-leap years,
// and [1,366] in leap years.
//
// This is a simple wrapper method for (time.Time{}).YearDay.
func (t Time[T]) YearDay() int {
	return t.tm.YearDay()
}

// Date returns the year, month, and day in which t occurs.
//
// This is a simple wrapper method for (time.Time{}).Date.
func (t Time[T]) Date() (year int, month time.Month, day int) {
	return t.tm.Date()
}

// Year returns the year in which t occurs.
//
// This is a simple wrapper method for (time.Time{}).Year.
func (t Time[T]) Year() int {
	return t.tm.Year()
}

// Month returns the month of the year specified by t.
//
// This is a simple wrapper method for (time.Time{}).Month.
func (t Time[T]) Month() time.Month {
	return t.tm.Month()
}

// Day returns the day of the month specified by t.
//
// This is a simple wrapper method for (time.Time{}).Day.
func (t Time[T]) Day() int {
	return t.tm.Day()
}

// Weekday returns the day of the week specified by t.
//
// This is a simple wrapper method for (time.Time{}).Weekday.
func (t Time[T]) Weekday() time.Weekday {
	return t.tm.Weekday()
}

// ISOWeek returns the ISO 8601 year and week number in which t occurs.
//
// This is a simple wrapper method for (time.Time{}).ISOWeek.
func (t Time[T]) ISOWeek() (year, week int) {
	return t.tm.ISOWeek()
}

// IsDST reports whether the time in the configured location is in Daylight Savings Time.
//
// This is a simple wrapper method for (time.Time{}).IsDST.
func (t Time[T]) IsDST() bool {
	return t.tm.IsDST()
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
//
// This is a simple wrapper method for (time.Time{}).IsZero.
func (t Time[T]) IsZero() bool {
	return t.tm.IsZero()
}

// Location returns the time zone information associated with t.
//
// This is a simple wrapper method for (time.Time{}).Location.
func (t Time[T]) Location() *time.Location {
	if t.tm.IsZero() {
		var tz T
		return tz.Location()
	}
	return t.tm.Location()
}

var nonZero = time.Unix(1, 0)

// Zone computes the time zone in effect at time t, returning the abbreviated
// name of the zone (such as "CET") and its offset in seconds east of UTC.
//
// This is a simple wrapper method for (time.Time{}).Zone.
func (t Time[T]) Zone() (name string, offset int) {
	// When time.Time is zero value, the location is set to UTC.
	if t.tm.IsZero() {
		return In[T](nonZero).Zone()
	}
	return t.tm.Zone()
}

// ZoneBounds returns the bounds of the time zone in effect at time t.
// The zone begins at start and the next zone begins at end.
// If the zone begins at the beginning of time, start will be returned as a zero Time.
// If the zone goes on forever, end will be returned as a zero Time.
// The Location of the returned times will be the same as t.
//
// This is a simple wrapper method for (time.Time{}).ZoneBounds.
func (t Time[T]) ZoneBounds() (start, end Time[T]) {
	tmStart, tmEnd := t.tm.ZoneBounds()
	return In[T](tmStart), In[T](tmEnd)
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
// Unix-like operating systems often record time as a 32-bit
// count of seconds, but since the method here returns a 64-bit
// value it is valid for billions of years into the past or future.
//
// This is a simple wrapper method for (time.Time{}).Unix.
func (t Time[T]) Unix() int64 {
	return t.tm.Unix()
}

// UnixMilli returns t as a Unix time, the number of milliseconds elapsed since
// January 1, 1970 UTC. The result is undefined if the Unix time in
// milliseconds cannot be represented by an int64 (a date more than 292 million
// years before or after 1970). The result does not depend on the
// location associated with t.
//
// This is a simple wrapper method for (time.Time{}).UnixMilli.
func (t Time[T]) UnixMilli() int64 {
	return t.tm.UnixMilli()
}

// UnixMicro returns t as a Unix time, the number of microseconds elapsed since
// January 1, 1970 UTC. The result is undefined if the Unix time in
// microseconds cannot be represented by an int64 (a date before year -290307 or
// after year 294246). The result does not depend on the location associated
// with t.
//
// This is a simple wrapper method for (time.Time{}).UnixMicro.
func (t Time[T]) UnixMicro() int64 {
	return t.tm.UnixMicro()
}

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC. The result is undefined if the Unix time
// in nanoseconds cannot be represented by an int64 (a date before the year
// 1678 or after 2262). Note that this means the result of calling UnixNano
// on the zero Time is undefined. The result does not depend on the
// location associated with t.
//
// This is a simple wrapper method for (time.Time{}).UnixNano.
func (t Time[T]) UnixNano() int64 {
	return t.tm.UnixNano()
}

// GobEncode implements the gob.GobEncoder interface.
//
// This is a simple wrapper method for (time.Time{}).GobEncode.
func (t Time[T]) GobEncode() ([]byte, error) {
	return t.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface.
//
// This is a simple wrapper method for (time.Time{}).GobDecode.
func (t *Time[T]) GobDecode(data []byte) error {
	return t.UnmarshalBinary(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
//
// This is a simple wrapper method for (time.Time{}).MarshalBinary.
func (t Time[T]) MarshalBinary() ([]byte, error) {
	return t.tm.MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
//
// This is a simple wrapper method for (time.Time{}).UnmarshalBinary.
func (t *Time[T]) UnmarshalBinary(data []byte) error {
	tm := time.Time{}
	if err := tm.UnmarshalBinary(data); err != nil {
		return err
	}
	*t = In[T](tm)
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in RFC 3339 format with sub-second precision.
// If the timestamp cannot be represented as valid RFC 3339
// (e.g., the year is out of range), then an error is reported.
//
// This is a simple wrapper method for (time.Time{}).MarshalText.
func (t Time[T]) MarshalText() ([]byte, error) {
	return t.tm.MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time must be in the RFC 3339 format.
//
// This is a simple wrapper method for (time.Time{}).MarshalText.
func (t *Time[T]) UnmarshalText(data []byte) error {
	tm := time.Time{}
	if err := tm.UnmarshalText(data); err != nil {
		return err
	}
	*t = In[T](tm)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in the RFC 3339 format with sub-second precision.
// If the timestamp cannot be represented as valid RFC 3339
// (e.g., the year is out of range), then an error is reported.
//
// This is a simple wrapper method for (time.Time{}).MarshalJSON.
func (t Time[T]) MarshalJSON() ([]byte, error) {
	return t.tm.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time must be a quoted string in the RFC 3339 format.
//
// This is a simple wrapper method for (time.Time{}).UnmarshalJSON.
func (t *Time[T]) UnmarshalJSON(data []byte) error {
	tm := time.Time{}
	if err := tm.UnmarshalJSON(data); err != nil {
		return err
	}
	*t = In[T](tm)
	return nil
}
