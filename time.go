package synchro

import (
	"math"
	"time"

	"github.com/itchyny/timefmt-go"
)

type empty[T TimeZone] struct{}

type Time[T TimeZone] struct {
	_  empty[T]
	tm time.Time
}

// StdTime returns the time.Time.
func (t Time[T]) StdTime() time.Time {
	return t.tm
}

// StartOfYear returns Time for start of the year.
func (t Time[T]) StartOfYear() Time[T] {
	return New[T](t.Year(), 1, 1, 0, 0, 0, 0)
}

// EndOfYear returns Time for end of the year.
func (t Time[T]) EndOfYear() Time[T] {
	return New[T](t.Year(), 12, 31, 23, 59, 59, 999999999)
}

// StartOfMonth returns Time for start of the month.
func (t Time[T]) StartOfMonth() Time[T] {
	return New[T](t.Year(), t.Month(), 1, 0, 0, 0, 0)
}

// EndOfMonth returns Time for end of the month.
func (t Time[T]) EndOfMonth() Time[T] {
	startOfMonth := t.StartOfMonth()
	return startOfMonth.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)
}

// StartOfWeek returns Time for start of the week.
func (t Time[T]) StartOfWeek() Time[T] {
	dayOfWeek := t.Weekday()
	return t.Add(-time.Duration(dayOfWeek) * 24 * time.Hour)
}

// EndOfWeek returns Time for end of the week.
func (t Time[T]) EndOfWeek() Time[T] {
	dayOfWeek := t.Weekday()
	return t.Add(time.Duration(time.Saturday-dayOfWeek+1) * 24 * time.Hour).Add(-1 * time.Nanosecond)
}

// StartOfQuarter returns a Time for start of the quarter.
func (t Time[T]) StartOfQuarter() Time[T] {
	year, quarter, day := t.Year(), numberOfQuarter(t.Month()), 1
	return New[T](year, time.Month(3*quarter-2), day, 0, 0, 0, 0)
}

// EndOfQuarter returns a Time for end of the quarter.
func (t Time[T]) EndOfQuarter() Time[T] {
	year, quarter := t.Year(), numberOfQuarter(t.Month())
	day := 31
	switch quarter {
	case 2, 3:
		day = 30
	}
	return New[T](year, time.Month(3*quarter), day, 23, 59, 59, 999999999)
}

// StartOfSemester returns a Time for start of the semester.
func (t Time[T]) StartOfSemester() Time[T] {
	year, semester, day := t.Year(), numberOfSemester(t.Month()), 1
	month := time.January
	if semester == 2 {
		month = time.July
	}
	return New[T](year, month, day, 0, 0, 0, 0)
}

// EndOfSemester returns a Time for end of the semester.
func (t Time[T]) EndOfSemester() Time[T] {
	year, semester := t.Year(), numberOfSemester(t.Month())
	month, day := time.June, 30
	if semester == 2 {
		month, day = time.December, 31
	}
	return New[T](year, month, day, 23, 59, 59, 999999999)
}

// IsLeapYear returns true if t is leap year.
func (t Time[T]) IsLeapYear() bool {
	year := t.Year()
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// IsBetween returns true if from < t && t < to.
func (t Time[T]) IsBetween(from Time[T], to Time[T]) bool {
	p, _ := NewPeriod[T](from, to)
	return p.Contains(t) == 1
}

// DiffInCalendarDays calculates the difference in calendar days between t and u. (t-u)
// Calendar days are calculated by considering only the dates, excluding the times,
// and then determining the difference in days.
func (t Time[T]) DiffInCalendarDays(u Time[T]) int {
	const day = 24 * time.Hour
	t1 := t.Truncate(day)
	u1 := u.Truncate(day)
	return int(math.Ceil(float64(t1.Sub(u1)) / float64(day)))
}

// Strftime formats the time according to the given format string.
//
// This method is a wrapper for the [github.com/itchyny/timefmt-go] library.
//
// Example:
//   - %Y-%m-%d %H:%M:%S => 2023-09-02 14:09:56
//   - %a, %d %b %Y %T %z => Sat, 02 Sep 2023 14:09:56 +0900
//
// The format string should follow the format of [strftime(3)] in man pages.
// The following list shows the supported format specifiers:
//   - %a: Abbreviated weekday name (Sun)
//   - %A: Full weekday name (Sunday)
//   - %b: Abbreviated month name (Jan)
//   - %B: Full month name (January)
//   - %c: Date and time representation
//   - %C: Year divided by 100 (00-99)
//   - %d: Day of the month (01-31)
//   - %D: Short MM/DD/YY date, equivalent to %m/%d/%y
//   - %e: Day of the month, with a space preceding single digits ( 1-31)
//   - %F: Equivalent to %Y-%m-%d (the ISO 8601 date format)
//   - %g: Week-based year, last two digits (00-99)
//   - %G: Week-based year
//   - %h: Abbreviated month name (Jan)
//   - %H: Hour in 24h format (00-23)
//   - %I: Hour in 12h format (01-12)
//   - %j: Day of the year (001-366)
//   - %m: Month as a decimal number (01-12)
//   - %M: Minute (00-59)
//   - %n: New-line character
//   - %p: AM or PM designation
//   - %P: am or pm designation
//   - %r: 12-hour clock time
//   - %R: 24-hour HH:MM time, equivalent to %H:%M
//   - %S: Second (00-59)
//   - %t: Horizontal-tab character
//   - %T: 24-hour clock time, equivalent to %H:%M:%S
//   - %u: ISO 8601 weekday as number with Monday as 1 (1-7)
//   - %U: Week number with the first Sunday as the first day of week (00-53)
//   - %V: ISO 8601 week number (01-53)
//   - %w: Weekday as a decimal number with Sunday as 0 (0-6)
//   - %W: Week number with the first Monday as the first day of week (00-53)
//   - %x: Date representation
//   - %X: Time representation
//   - %y: Year, last two digits (00-99)
//   - %Y: Year
//   - %z: ISO 8601 offset from UTC in timezone (+HHMM)
//   - %Z: Timezone name or abbreviation
//   - %+: Extended date and time representation
//   - %::z: Colon-separated offset from UTC in timezone (e.g. +05:00)
//   - %:::z: Like %::z, but with optional seconds
//
// [strftime(3)]: https://linux.die.net/man/3/strftime
// [github.com/itchyny/timefmt-go]: https://github.com/itchyny/timefmt-go
func (t Time[T]) Strftime(format string) string {
	return timefmt.Format(t.StdTime(), format)
}

// Strptime parses time string with the default location.
// The location is also used to parse the time zone name (%Z).
//
// The format string should follow the format of [strptime(3)] in man pages.
// The following list shows the supported format specifiers:
//   - %a: abbreviated weekday name
//   - %A: full weekday name
//   - %b: abbreviated month name
//   - %B: full month name
//   - %c: preferred date and time representation
//   - %C: century number (00-99)
//   - %d: day of the month (01-31)
//   - %D: same as %m/%d/%y
//   - %e: day of the month (1-31)
//   - %F: same as %Y-%m-%d
//   - %g: last two digits of the year (00-99)
//   - %G: year as a 4-digit number
//   - %h: same as %b
//   - %H: hour (00-23)
//   - %I: hour (01-12)
//   - %j: day of the year (001-366)
//   - %m: month (01-12)
//   - %M: minute (00-59)
//   - %n: newline character
//   - %p: either "am" or "pm" according to the given time value
//   - %r: time in a.m. and p.m. notation
//   - %R: time in 24 hour notation
//   - %S: second (00-60)
//   - %t: tab character
//   - %T: current time, equal to %H:%M:%S
//   - %u: weekday as a number (1-7)
//   - %U: week number of the current year, starting with the first Sunday as the first day of the first week
//   - %V: week number of the current year, starting with the first week that has at least 4 days in the new year
//   - %w: day of the week as a decimal, Sunday being 0
//   - %W: week number of the current year, starting with the first Monday as the first day of the first week
//   - %x: preferred date representation without the time
//   - %X: preferred time representation without the date
//   - %y: year without a century (00-99)
//   - %Y: year with century
//   - %z: time zone offset, such as "-0700"
//   - %Z: time zone name, such as "UTC" or "GMT"
//
// This is a wrapper for the [github.com/itchyny/timefmt-go] library.
//
// [strptime(3)]: https://linux.die.net/man/3/strptime
// [github.com/itchyny/timefmt-go]: https://github.com/itchyny/timefmt-go
func Strptime[T TimeZone](source string, format string) (Time[T], error) {
	var tz T
	tm, err := timefmt.ParseInLocation(source, format, tz.Location())
	if err != nil {
		return Time[T]{}, err
	}
	return In[T](tm), nil
}
