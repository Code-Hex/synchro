package iso8601

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DateLike defines an interface for date-related structures.
// It provides methods for retrieving the date, validating the date,
// and checking if the date is valid.
type DateLike interface {
	// Date returns the underlying Date value.
	Date() Date

	// IsValid checks whether the date is valid.
	IsValid() bool

	// Validate checks the correctness of the date and returns an error if it's invalid.
	Validate() error
}

// Date represents a calendar date with year, month, and day components.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

var _ DateLike = Date{}

// Date returns itself as it directly represents a date.
func (d Date) Date() Date {
	return d
}

// IsValid checks if the date is valid based on its year, month, and day values.
func (d Date) IsValid() bool {
	return d.Validate() == nil
}

// Validate checks the individual components of the date (year, month, and day)
// and returns an error if any of them are out of the expected ranges.
func (d Date) Validate() error {
	if d.Year < 0 || d.Year > 9999 {
		return &DateLikeRangeError{
			Element: "year",
			Value:   d.Year,
			Year:    d.Year,
			Min:     0,
			Max:     9999,
		}
	}
	if d.Month < 1 || d.Month > 12 {
		return &DateLikeRangeError{
			Element: "month",
			Value:   int(d.Month),
			Year:    d.Year,
			Min:     1,
			Max:     12,
		}
	}
	daysInMonth := daysInMonth(d.Year, int(d.Month))
	if d.Day < 1 || d.Day > daysInMonth {
		return &DateLikeRangeError{
			Element: "day of month",
			Value:   d.Day,
			Year:    d.Year,
			Min:     1,
			Max:     daysInMonth,
		}
	}
	return nil
}

// StdTime converts the Date structure to a time.Time object, using UTC for the time.
func (d Date) StdTime() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

// QuarterDate represents a date within a specific quarter of a year.
// It includes the year, quarter (from 1 to 4), and day within that quarter.
type QuarterDate struct {
	Year    int
	Quarter int
	Day     int
}

var _ DateLike = QuarterDate{}

// Date converts a QuarterDate into the standard Date representation.
// It calculates the exact calendar date based on the year, quarter, and day within that quarter.
func (q QuarterDate) Date() Date {
	yday := q.Day // 1 ~ 366
	for i := 1; i < q.Quarter; i++ {
		yday += daysInQuarter(q.Year, i)
	}
	t := time.Date(q.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, yday-1)
	return Date{
		Year:  t.Year(),
		Month: t.Month(),
		Day:   t.Day(),
	}
}

// IsValid checks if the quarter date is valid based on its year, quarter, and day within the quarter values.
func (q QuarterDate) IsValid() bool {
	return q.Validate() == nil
}

// Validate checks the individual components of the quarter date (year, quarter, and day within the quarter)
// and returns an error if any of them are out of the expected ranges.
func (q QuarterDate) Validate() error {
	if q.Year < 0 || q.Year > 9999 {
		return &DateLikeRangeError{
			Element: "year",
			Value:   q.Year,
			Year:    q.Year,
			Min:     0,
			Max:     9999,
		}
	}
	if q.Quarter < 1 || q.Quarter > 4 {
		return &DateLikeRangeError{
			Element: "quarter",
			Value:   q.Quarter,
			Year:    q.Year,
			Min:     1,
			Max:     4,
		}
	}
	daysInQuarter := daysInQuarter(q.Year, q.Quarter)
	if q.Day < 1 || q.Day > daysInQuarter {
		return &DateLikeRangeError{
			Element: "day of quarter",
			Value:   q.Day,
			Year:    q.Year,
			Min:     1,
			Max:     daysInQuarter,
		}
	}
	return nil
}

// WeekDate represents a date within a specific week of a given year,
// following the ISO 8601 week-date system. It includes the year,
// week number (1 to 52 or 53), and day of the week (1 for Monday to 7 for Sunday).
type WeekDate struct {
	Year int
	Week int
	Day  int
}

var _ DateLike = WeekDate{}

// Date converts a WeekDate into the standard Date representation.
// It calculates the exact calendar date based on the year, week number, and day of the week.
func (w WeekDate) Date() Date {
	// Find the first Thursday of the given year. This will be in the first week of the year according to ISO 8601.
	thursday := time.Date(w.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	for thursday.Weekday() != time.Thursday {
		thursday = thursday.AddDate(0, 0, 1)
	}

	// Calculate the date of the Monday of week 1
	monday := thursday.AddDate(0, 0, -3)

	// Calculate the date corresponding to the given week and day
	t := monday.AddDate(0, 0, (w.Week-1)*7+w.Day-1)
	return Date{
		Year:  t.Year(),
		Month: t.Month(),
		Day:   t.Day(),
	}
}

// IsValid checks if the week date is valid based on its year, week number, and day of the week values.
func (w WeekDate) IsValid() bool {
	return w.Validate() == nil
}

// Validate checks the individual components of the week date (year, week number, and day of the week)
// and returns an error if any of them are out of the expected ranges.
func (w WeekDate) Validate() error {
	if w.Year < 0 || w.Year > 9999 {
		return &DateLikeRangeError{
			Element: "year",
			Value:   w.Year,
			Year:    w.Year,
			Min:     0,
			Max:     9999,
		}
	}
	if w.Day < 1 || w.Day > 7 {
		return &DateLikeRangeError{
			Element: "day of week",
			Value:   int(w.Day),
			Year:    w.Year,
			Min:     1,
			Max:     7,
		}
	}
	weeksInYear := weeksInYear(w.Year)
	if w.Week < 1 || w.Week > weeksInYear {
		return &DateLikeRangeError{
			Element: "week",
			Value:   w.Week,
			Year:    w.Year,
			Min:     1,
			Max:     weeksInYear,
		}
	}
	return nil
}

// OrdinalDate represents a date specified by its year and the day-of-year (ordinal date),
// where the day-of-year ranges from 1 through 365 (or 366 in a leap year).
type OrdinalDate struct {
	Year int
	Day  int
}

var _ DateLike = OrdinalDate{}

// Date converts an OrdinalDate into the standard Date representation.
// It calculates the exact calendar date based on the year and the day-of-year.
func (o OrdinalDate) Date() Date {
	yday := o.Day // 1 ~ 366
	t := time.Date(o.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, yday-1)
	return Date{
		Year:  o.Year,
		Month: t.Month(),
		Day:   t.Day(),
	}
}

// IsValid checks if the ordinal date is valid based on its year and day-of-year values.
func (o OrdinalDate) IsValid() bool {
	return o.Validate() == nil
}

// Validate checks the individual components of the ordinal date (year and day-of-year)
// and returns an error if any of them are out of the expected ranges.
func (o OrdinalDate) Validate() error {
	if o.Year < 0 || o.Year > 9999 {
		return &DateLikeRangeError{
			Element: "year",
			Value:   o.Year,
			Year:    o.Year,
			Min:     0,
			Max:     9999,
		}
	}
	daysInYear := daysInYear(o.Year)
	if o.Day < 1 || o.Day > daysInYear {
		return &DateLikeRangeError{
			Element: "day of year",
			Value:   o.Day,
			Year:    o.Year,
			Min:     1,
			Max:     daysInYear,
		}
	}
	return nil
}

// DateLikeRangeError indicates that a value is not in an expected range for DateLike.
type DateLikeRangeError struct {
	Element string
	Value   int
	Year    int
	Min     int
	Max     int
}

// Error implements the error interface.
func (e *DateLikeRangeError) Error() string {
	return fmt.Sprintf("iso8601: %d %s is not in range %d-%d in %d", e.Value, e.Element, e.Min, e.Max, e.Year)
}

// Time represents an ISO8601-compliant time without a date, specified by its hour, minute, second, and nanosecond.
//
// Note: The time '24:00:00' is valid and represents midnight at the end of the given day.
type Time struct {
	Hour       int
	Minute     int
	Second     int
	Nanosecond int
}

// Validate checks the individual components of the time (hour, minute, second, and nanosecond)
// against their respective valid ISO8601 ranges.
//
// Specifically, it validates:
//   - Minute values between 0 and 59 inclusive.
//   - Second values between 0 and 59 inclusive.
//   - Hour values between 0 and 23 inclusive, with an exception for the '24:00:00' time.
//
// Returns an error if any of the components are out of their expected ranges.
func (t Time) Validate() error {
	if t.Minute < 0 || t.Minute > 59 {
		return &TimeRangeError{
			Element: "minute",
			Value:   t.Minute,
			Min:     0,
			Max:     59,
		}
	}
	if t.Second < 0 || t.Second > 59 {
		return &TimeRangeError{
			Element: "second",
			Value:   t.Second,
			Min:     0,
			Max:     59,
		}
	}
	if t.Hour < 0 || t.Hour > 23 {
		if !(t.Hour == 24 && t.Minute == 0 && t.Second == 0 && t.Nanosecond == 0) {
			return &TimeRangeError{
				Element: "hour",
				Value:   t.Hour,
				Min:     0,
				Max:     24,
			}
		}
	}
	return nil
}

// TimeRangeError indicates that a value is not in an expected range for Time.
type TimeRangeError struct {
	Element string
	Value   int
	Min     int
	Max     int
}

// Error implements the error interface.
func (e *TimeRangeError) Error() string {
	return fmt.Sprintf("iso8601 time: %d %s is not in range %d-%d", e.Value, e.Element, e.Min, e.Max)
}

// Zone represents an ISO8601-compliant timezone offset, specified by its hour, minute, and second components.
// The "Negative" field indicates whether the offset is behind (true) or ahead (false) of Coordinated Universal Time (UTC).
type Zone struct {
	Hour     int
	Minute   int
	Second   int
	Negative bool
}

// Validate checks the individual components of the timezone offset (hour, minute, and second)
// against their respective valid ISO8601 ranges.
//
// Specifically, it validates:
//   - Minute values between 0 and 99 inclusive.
//   - Second values between 0 and 99 inclusive.
//   - Hour values between 0 and 99 inclusive.
//
// Returns an error if any of the components are out of their expected ranges.
func (z Zone) Validate() error {
	if z.Minute < 0 || z.Minute > 99 {
		return &TimeZoneRangeError{
			Element: "minute",
			Value:   z.Minute,
			Min:     0,
			Max:     99,
		}
	}
	if z.Second < 0 || z.Second > 99 {
		return &TimeZoneRangeError{
			Element: "second",
			Value:   z.Second,
			Min:     0,
			Max:     99,
		}
	}
	if z.Hour < 0 || z.Hour > 99 {
		return &TimeZoneRangeError{
			Element: "hour",
			Value:   z.Hour,
			Min:     0,
			Max:     99,
		}
	}
	return nil
}

// Offset computes the total time offset in seconds for the Zone.
// This value is positive if the zone is ahead of UTC, and negative if it's behind.
// It takes into account the hour, minute, second, and Negative fields.
func (z Zone) Offset() int {
	sign := 1
	if z.Negative {
		sign = -1
	}
	return sign * (z.Hour*3600 + z.Minute*60 + z.Second)
}

// TimeRangeError indicates that a value is not in an expected range for Time.
type TimeZoneRangeError struct {
	Element string
	Value   int
	Min     int
	Max     int
}

// Error implements the error interface.
func (e *TimeZoneRangeError) Error() string {
	return fmt.Sprintf("iso8601 time zone: %d %s is not in range %d-%d", e.Value, e.Element, e.Min, e.Max)
}

// Duration represents an ISO8601 duration with the maximum precision of nanoseconds.
// It includes components like years, months, weeks, days, hours, minutes, seconds,
// milliseconds, microseconds, and nanoseconds. The Negative field indicates whether
// the duration is negative.
type Duration struct {
	Year        int
	Month       time.Month
	Week        int
	Day         int
	Hour        int
	Minute      int
	Second      int
	Millisecond int
	Microsecond int
	Nanosecond  int
	Negative    bool
}

const yearInSecond = 31556952 * time.Second // 365.2425 days * 3600 * 24 seconds
const monthInSecond = 2630016 * time.Second // 30.44 days * 3600 * 24 seconds
const weekInSecond = 7 * dayInSecond
const dayInSecond = 24 * 3600 * time.Second

// StdDuration converts the ISO8601 Duration to a standard Go time.Duration.
// Note: This conversion is an approximation. The duration of some components
// like years and months are averaged based on typical values:
//   - 1 year is considered as 365.2425 days (or 31556952 seconds).
//   - 1 month is considered as 30.44 days (or 2630016 seconds).
//   - 1 week is considered as 7 days (or 604800 seconds).
func (d Duration) StdDuration() time.Duration {
	duration := time.Duration(d.Year)*yearInSecond +
		time.Duration(d.Month)*monthInSecond +
		time.Duration(d.Week)*weekInSecond +
		time.Duration(d.Day)*dayInSecond +
		time.Duration(d.Hour)*time.Hour +
		time.Duration(d.Minute)*time.Minute +
		time.Duration(d.Second)*time.Second +
		time.Duration(d.Millisecond)*time.Millisecond +
		time.Duration(d.Microsecond)*time.Microsecond +
		time.Duration(d.Nanosecond)

	if d.Negative {
		duration = -duration
	}
	return time.Duration(duration)
}

// NewDuration makes ISO8601 Duration struct from time.Duration.
func NewDuration(d time.Duration) Duration {
	negative := false
	if d < 0 {
		negative = true
		d = -d
	}

	year := int(d / yearInSecond)
	d -= time.Duration(year) * yearInSecond

	month := int(d / monthInSecond)
	d -= time.Duration(month) * monthInSecond

	week := int(d / weekInSecond)
	d -= time.Duration(week) * weekInSecond

	day := int(d / dayInSecond)
	d -= time.Duration(day) * dayInSecond

	hour := int(d / time.Hour)
	d -= time.Duration(hour) * time.Hour

	minute := int(d / time.Minute)
	d -= time.Duration(minute) * time.Minute

	second := int(d / time.Second)
	d -= time.Duration(second) * time.Second

	millisecond := int(d / time.Millisecond)
	d -= time.Duration(millisecond) * time.Millisecond

	microsecond := int(d / time.Microsecond)
	d -= time.Duration(microsecond) * time.Microsecond

	nanosecond := int(d)

	return Duration{
		Year:        year,
		Month:       time.Month(month),
		Week:        week,
		Day:         day,
		Hour:        hour,
		Minute:      minute,
		Second:      second,
		Millisecond: millisecond,
		Microsecond: microsecond,
		Nanosecond:  nanosecond,
		Negative:    negative,
	}
}

// Negate changes the sign of the duration.
func (d Duration) Negate() Duration {
	src := &d
	dst := &Duration{}
	*dst = *src
	dst.Negative = !dst.Negative
	return *dst
}

// IsZero checks duration is zero value.
func (d Duration) IsZero() bool {
	return d == Duration{} || d == Duration{Negative: true}
}

// String returns the ISO8601 string representation of the duration.
// It produces a minimal string without redundant zeros.
// For example, a duration of one year, two months, three days, four hours, five minutes,
// six seconds, and seven milliseconds would be "P1Y2M3DT4H5M6.007S".
func (d Duration) String() string {
	var b strings.Builder
	if d.Negative {
		b.WriteByte('-')
	}
	b.WriteByte('P')
	if d.IsZero() {
		b.WriteString("T0S")
		return b.String()
	}

	hasTime := false
	write := func(v int, designator byte, isTime bool) {
		if !hasTime && isTime {
			b.WriteByte('T')
			hasTime = true
		}
		b.WriteString(strconv.Itoa(v))
		b.WriteByte(designator)
	}

	writeSec := func(v int, fraction int) {
		if !hasTime && (fraction > 0 || d.Second > 0) {
			b.WriteByte('T')
			hasTime = true
		}
		if fraction > 0 {
			b.WriteString(strconv.Itoa(v))
			b.WriteByte('.')
			fmt.Fprintf(&b, "%09d", fraction)
			b.WriteByte('S')
		} else if d.Second > 0 {
			b.WriteString(strconv.Itoa(v))
			b.WriteByte('S')
		}
	}

	if d.Year != 0 {
		write(d.Year, 'Y', false)
	}
	if d.Month != 0 {
		write(int(d.Month), 'M', false)
	}
	if d.Week != 0 {
		write(d.Week, 'W', false)
	}
	if d.Day != 0 {
		write(d.Day, 'D', false)
	}
	if d.Hour != 0 {
		write(d.Hour, 'H', true)
	}
	if d.Minute != 0 {
		write(d.Minute, 'M', true)
	}

	nanosec := 0
	if d.Millisecond != 0 {
		nanosec += d.Millisecond * 1e6
	}
	if d.Microsecond != 0 {
		nanosec += d.Microsecond * 1000
	}
	if d.Nanosecond != 0 {
		nanosec += d.Nanosecond
	}

	writeSec(d.Second, nanosec)
	return b.String()
}
