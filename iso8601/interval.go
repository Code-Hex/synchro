package iso8601

import (
	"bytes"
	"fmt"
	"time"
)

// Interval represents an ISO8601 time interval.
// It contains the start and end times, the duration of the interval,
// and the number of times the interval should be repeated.
type Interval struct {
	// Start represents the start time of the interval.
	start time.Time

	// End represents the end time of the interval.
	end time.Time

	// Duration represents the duration of the interval.
	duration Duration

	// Repeat represents the number of times the interval should be repeated. -1 indicates infinity.
	repeat int
}

// Start returns a time.Time representing the beginning of this interval.
// If this interval has an explicit End date specified, any existing relative
// Duration will be cleared.
//
// Note: if the interval doesn't include a time component, the start time will
// actually be zero value (January 1, year 1, 00:00:00 UTC.) of the following
// day (since the interval covers the entire day). Intervals include the start
// value (in contrast to the end value).
func (i Interval) Start() time.Time {
	if !i.start.IsZero() {
		return i.start
	}
	if !i.end.IsZero() {
		// Assuming the Duration struct has a method to convert it to time.Duration
		return i.end.Add(-i.duration.StdDuration())
	}
	return time.Time{}
}

// End returns a time.Time representing the end time of the interval.
// If this interval has an explicit Start date specified, any existing
// relative Duration will be cleared.
//
// Note: if the interval doesn't include a time component, the end
// time will actually be zero value (January 1, year 1, 00:00:00 UTC.) of
// the following day (since the interval covers the entire day).
func (i Interval) End() time.Time {
	if !i.end.IsZero() {
		return i.end
	}
	if !i.start.IsZero() {
		return i.start.Add(i.duration.StdDuration())
	}
	return time.Time{}
}

// Duration returns ISO 8601 duration.
func (i Interval) Duration() Duration {
	if !i.duration.IsZero() {
		return i.duration
	}
	return NewDuration(i.end.Sub(i.start))
}

// Contains returns a boolean indicating whether the provided time.Time
// is between the Start or End dates as defined by this interval.
func (i Interval) Contains(t time.Time) bool {
	return t.Compare(i.Start()) >= 0 && t.Compare(i.End()) <= 0
}

// ParseInterval parses an ISO8601 time interval from a byte slice or string.
// It returns the parsed Interval and any error encountered.
func ParseInterval[bytes []byte | ~string](b bytes) (Interval, error) {
	return parseInterval([]byte(b))
}

func parseInterval(b []byte) (Interval, error) {
	var (
		designator []byte
		start      time.Time
		end        time.Time
		duration   Duration
		repeat     int
	)
	if len(b) == 0 {
		return Interval{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b),
			AfterToken: "",
			Expected:   "R or P or datetime",
		}
	}

	if b[0] == 'R' {
		designatorIdx := -1
		// https://go.dev/play/p/42u_xsxugSW
		for _, candidate := range [][]byte{
			{'/'},
			{'-', '-'},
		} {
			n := bytes.Index(b, candidate)
			if n >= 0 && (n < designatorIdx || designatorIdx == -1) {
				designator = candidate
				designatorIdx = n
			}
		}
		if designatorIdx == -1 {
			return Interval{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[1:]),
				AfterToken: "R",
				Expected:   `internal designator "/" or "--"`,
			}
		}
		c := countDigits(b, 1)
		if want := designatorIdx - 1; c != want {
			return Interval{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[1:designatorIdx]),
				AfterToken: "R",
				Expected:   humanizeDigits(want),
			}
		}
		if c == 0 {
			repeat = -1 // infinity
		} else {
			repeat = parseNumber(b, 1, c)
		}
		b = b[designatorIdx+len(designator):]
	}

	designatorIdx := -1
	if designator != nil {
		n := bytes.Index(b, designator)
		if n >= 0 {
			designatorIdx = n
		}
	} else {
		for _, candidate := range [][]byte{
			{'/'},
			{'-', '-'},
		} {
			n := bytes.Index(b, candidate)
			if n >= 0 {
				designator = candidate
				designatorIdx = n
				break
			}
		}
	}
	// try to parse duration only
	if designatorIdx == -1 {
		d, err := parseDuration(b)
		if err != nil {
			return Interval{}, err
		}
		return Interval{
			duration: d,
			repeat:   repeat,
		}, nil
	}

	startb := b[:designatorIdx]
	isStartDurationFormat := b[0] == 'P'

	if isStartDurationFormat {
		d, err := parseDuration(startb)
		if err != nil {
			return Interval{}, err
		}
		duration = d
	} else {
		dt, err := parseDateTime(startb)
		if err != nil {
			return Interval{}, err
		}
		start = dt
	}

	endb := b[len(designator)+designatorIdx:]
	if endb[0] == 'P' {
		if !duration.IsZero() {
			return Interval{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(endb),
				AfterToken: string(designator),
				Expected:   "datetime format",
			}
		}
		d, err := parseDuration(endb)
		if err != nil {
			return Interval{}, err
		}
		duration = d
	} else if !isStartDurationFormat && len(startb) > len(endb) {
		dt, err := parseShortEndDuration(start, endb)
		if err != nil {
			return Interval{}, err
		}
		end = dt
	} else {
		dt, err := parseDateTime(endb)
		if err != nil {
			return Interval{}, err
		}
		end = dt
	}

	return Interval{
		start:    start,
		end:      end,
		duration: duration,
		repeat:   repeat,
	}, nil
}

func parseShortEndDuration(start time.Time, b []byte) (time.Time, error) {
	var (
		year, m, day         = start.Date()
		month                = int(m)
		hour, minute, second = start.Clock()
	)
	n := 0

	switch countDigits(b, 0) {
	case 4:
		year = parseNumber(b, 0, 4)
		if len(b) == 4 {
			return shortEndInterval(start, year, month, day, hour, minute, second)
		}
		if b[4] == 'T' {
			n = 4
			break
		}
		if b[4] != '-' {
			return time.Time{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[4]),
				AfterToken: string(b[:4]),
				Expected:   "-",
			}
		}
		if c := countDigits(b, 5); c != 2 {
			return time.Time{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      humanizeDigits(c),
				AfterToken: string("-"),
				Expected:   humanizeDigits(2),
			}
		}
		month = parseNumber(b, 5, 2)
		if len(b) == 7 {
			return shortEndInterval(start, year, month, day, hour, minute, second)
		}
		if b[7] != '-' {
			return time.Time{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[7]),
				AfterToken: string(b[:7]),
				Expected:   "-",
			}
		}
		if c := countDigits(b, 8); c != 2 {
			return time.Time{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      humanizeDigits(c),
				AfterToken: string("-"),
				Expected:   humanizeDigits(2),
			}
		}
		day = parseNumber(b, 8, 2)
		if len(b) == 10 {
			return shortEndInterval(start, year, month, day, hour, minute, second)
		}
		n = 10
	case 2:
		if len(b) == 2 {
			day = parseNumber(b, 0, 2)
			return shortEndInterval(start, year, month, day, hour, minute, second)
		}

		if b[2] == ':' {
			hour = parseNumber(b, 0, 2)
			if c := countDigits(b, 3); c != 2 {
				return time.Time{}, &UnexpectedTokenError{
					Value:      string(b),
					Token:      humanizeDigits(c),
					AfterToken: string(":"),
					Expected:   humanizeDigits(2),
				}
			}
			minute = parseNumber(b, 3, 2)
			if len(b) == 5 {
				return shortEndInterval(start, year, month, day, hour, minute, second)
			}
			if b[5] != ':' {
				return time.Time{}, &UnexpectedTokenError{
					Value:      string(b),
					Token:      string(b[5]),
					AfterToken: string(b[:5]),
					Expected:   ":",
				}
			}
			if c := countDigits(b, 6); c != 2 {
				return time.Time{}, &UnexpectedTokenError{
					Value:      string(b),
					Token:      humanizeDigits(c),
					AfterToken: string(":"),
					Expected:   humanizeDigits(2),
				}
			}
			second = parseNumber(b, 6, 2)
			return shortEndInterval(start, year, month, day, hour, minute, second)
		}

		if b[2] == '-' {
			month = parseNumber(b, 0, 2)
			if c := countDigits(b, 3); c != 2 {
				return time.Time{}, &UnexpectedTokenError{
					Value:      string(b),
					Token:      humanizeDigits(c),
					AfterToken: string("-"),
					Expected:   humanizeDigits(2),
				}
			}
			day = parseNumber(b, 3, 2)
			if len(b) == 5 {
				return shortEndInterval(start, year, month, day, hour, minute, second)
			}
			n = 5
		} else {
			day = parseNumber(b, 0, 2)
			n = 2
		}
	}

	if b[n] != 'T' {
		return time.Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b[n]),
			AfterToken: string(b[:n]),
			Expected:   "T",
		}
	}

	n++

	hour = parseNumber(b, n, 2)
	if c := countDigits(b, n); c != 2 {
		return time.Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      humanizeDigits(c),
			AfterToken: string("T"),
			Expected:   humanizeDigits(2),
		}
	}

	n += 2

	if len(b) == n {
		return shortEndInterval(start, year, month, day, hour, minute, second)
	}

	if b[n] != ':' {
		return time.Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b[n]),
			AfterToken: string(b[:n]),
			Expected:   ":",
		}
	}

	n++

	if c := countDigits(b, n); c != 2 {
		return time.Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      humanizeDigits(c),
			AfterToken: string(":"),
			Expected:   humanizeDigits(2),
		}
	}

	minute = parseNumber(b, n, 2)

	n += 2

	if len(b) == n {
		return shortEndInterval(start, year, month, day, hour, minute, second)
	}

	if b[n] != ':' {
		return time.Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b[n]),
			AfterToken: string(b[:n]),
			Expected:   ":",
		}
	}

	n++

	if c := countDigits(b, n); c != 2 {
		return time.Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      humanizeDigits(c),
			AfterToken: string(":"),
			Expected:   humanizeDigits(2),
		}
	}

	second = parseNumber(b, n, 2)
	return shortEndInterval(start, year, month, day, hour, minute, second)
}

func shortEndInterval(
	start time.Time,
	years,
	months,
	days,
	hours,
	minutes,
	seconds int,
) (time.Time, error) {
	var (
		year, m, day         = start.Date()
		month                = int(m)
		hour, minute, second = start.Clock()
	)
	if year > years {
		return time.Time{}, &IntervalRangeError{
			Element: "year",
			Value:   years,
			Min:     year,
			Max:     9999,
		}
	}
	if year == years && month > months {
		return time.Time{}, &IntervalRangeError{
			Element: "month",
			Value:   months,
			Min:     month,
			Max:     12,
		}
	}
	if err := validateIntervalDate(years, months, day); err != nil {
		if year == years {
			err.Min = month
		}
		return time.Time{}, err
	}
	if month == months && day > days {
		return time.Time{}, &IntervalRangeError{
			Element: "day of month",
			Value:   days,
			Min:     day,
			Max:     daysInMonth(year, month),
		}
	}
	if err := validateIntervalDate(years, months, days); err != nil {
		if year == years && month == months {
			err.Min = day
		}
		return time.Time{}, err
	}

	if hour > hours {
		return time.Time{}, &IntervalRangeError{
			Element: "hour",
			Value:   hours,
			Min:     hour,
			Max:     24,
		}
	}
	if err := validateIntervalTime(hours, minute, second); err != nil {
		err.Min = hour
		return time.Time{}, err
	}
	if hour == hours && minute > minutes {
		return time.Time{}, &IntervalRangeError{
			Element: "minute",
			Value:   minutes,
			Min:     minute,
			Max:     59,
		}
	}
	if err := validateIntervalTime(hours, minutes, second); err != nil {
		err.Min = minute
		return time.Time{}, err
	}
	if hour == hours && minute == minutes && second > seconds {
		return time.Time{}, &IntervalRangeError{
			Element: "second",
			Value:   seconds,
			Min:     second,
			Max:     59,
		}
	}
	if err := validateIntervalTime(hours, minutes, seconds); err != nil {
		err.Min = second
		return time.Time{}, err
	}
	return time.Date(years, time.Month(months), days, hours, minutes, seconds, start.Nanosecond(), time.UTC), nil
}

func validateIntervalDate(year, month, day int) *IntervalRangeError {
	d := Date{
		Year:  year,
		Month: time.Month(month),
		Day:   day,
	}
	err := d.Validate()
	if v, ok := err.(*DateLikeRangeError); ok && v != nil {
		return &IntervalRangeError{
			Element: v.Element,
			Value:   v.Value,
			Min:     v.Min,
			Max:     v.Max,
		}
	}
	return nil
}

func validateIntervalTime(hour, minute, second int) *IntervalRangeError {
	t := Time{
		Hour:   hour,
		Minute: minute,
		Second: second,
	}
	err := t.Validate()
	if v, ok := err.(*TimeRangeError); ok && v != nil {
		return &IntervalRangeError{
			Element: v.Element,
			Value:   v.Value,
			Min:     v.Min,
			Max:     v.Max,
		}
	}
	return nil
}

// IntervalRangeError indicates that a value is not in an expected range for time interval.
type IntervalRangeError struct {
	Element string
	Value   int
	Min     int
	Max     int
}

// Error implements the error interface.
func (e *IntervalRangeError) Error() string {
	return fmt.Sprintf("iso8601 time interval: %d %s is not in range %d-%d", e.Value, e.Element, e.Min, e.Max)
}
