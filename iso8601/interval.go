package iso8601

import (
	"bytes"
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

	if b[0] == 'P' {
		d, err := parseDuration(b[:designatorIdx])
		if err != nil {
			return Interval{}, err
		}
		duration = d
	} else {
		dt, err := parseDateTime(b[:designatorIdx])
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
