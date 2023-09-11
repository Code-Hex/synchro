package iso8601

import (
	"fmt"
	"strings"
	"time"
)

var defaultParseDateTimeOptions = parseDateTimeOptions{
	timeDesignators: []byte{'T'},
	local:           time.Local,
}

type parseDateTimeOptions struct {
	timeDesignators []byte
	local           *time.Location
}

// ParseDateTimeOptions is a function type that modifies the parsing behavior
// of a datetime string. It acts as a functional option.
type ParseDateTimeOptions func(*parseDateTimeOptions)

// WithTimeDesignators is an option that modifies the set of valid
// characters which can be used as time designators when parsing a datetime string.
//
// By default, if no designators are set, the parser uses only 'T'.
func WithTimeDesignators(designators ...byte) ParseDateTimeOptions {
	return func(o *parseDateTimeOptions) {
		o.timeDesignators = append(o.timeDesignators, designators...)
	}
}

// WithInLocation is an options to interpret the time as in the given location.
//
// By default, if no location is set, the parser uses current location (Local).
func WithInLocation(loc *time.Location) ParseDateTimeOptions {
	return func(o *parseDateTimeOptions) {
		o.local = loc
	}
}

// ParseDateTime attempts to parse a given byte slice representing combined date, time,
// and optionally timezone offset in supported ISO 8601 formats. Supported formats include:
//
//	Basic                        Extended
//	20070301                     2007-03-01
//	2012W521                     2012-W52-1
//	2012Q485                     2012-Q4-85
//	20070301T1300Z               2007-03-01T13:00Z
//	20070301T1300Z               2007-03-01T13:00Z
//	20070301T1300+0100           2007-03-01T13:00+01:00
//	20070301T1300-0600           2007-03-01T13:00-06:00
//	20070301T130045Z             2007-03-01T13:00:45Z
//	20070301T130045+0100         2007-03-01T13:00:45+01:00
//	... and other combinations
//
// The function returns a time.Time struct representing the parsed date-time, adjusted
// for the parsed timezone offset if provided.
//
// In the absence of a time zone indicator, Parse returns a time in UTC.
//
// When parsing a time with a zone offset like -0700, if the offset corresponds
// to a time zone used by the current location (Local), then Parse uses that
// location and zone in the returned time. Otherwise it records the time as
// being in a fabricated location with time fixed at the given zone offset.
//
// If parsing fails, an error is returned.
func ParseDateTime[bytes []byte | ~string](b bytes, opts ...ParseDateTimeOptions) (time.Time, error) {
	return parseDateTime([]byte(b), opts...)
}

func parseDateTime(b []byte, opts ...ParseDateTimeOptions) (time.Time, error) {
	o := new(parseDateTimeOptions)
	*o = defaultParseDateTimeOptions // apply default options
	for _, opt := range opts {
		opt(o)
	}

	n, d, err := parseDate(b)
	if err != nil {
		return time.Time{}, overrideUnexpectedTokenValue(err, b)
	}
	dt := d.Date()
	if len(b) == n {
		return dt.StdTime(), nil
	}

	if len(b) > n {
		var found bool
		for _, designator := range o.timeDesignators {
			if b[n] == designator {
				found = true
				break
			}
		}
		if !found {
			var buf strings.Builder
			size := len(o.timeDesignators)
			for _, designator := range o.timeDesignators[:size-1] {
				fmt.Fprintf(&buf, "%q, ", designator)
			}
			fmt.Fprintf(&buf, "%q", o.timeDesignators[size-1])
			return time.Time{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[n]),
				AfterToken: string(b[:n]),
				Expected:   buf.String(),
			}
		}
	}
	n++

	// check byte length <date> + 'T' + the minimum length of <time>
	// 2023-09-10T22
	if len(b) < n+1 {
		return time.Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b[n:]),
			AfterToken: string(b[:n]),
			Expected:   "time format is required after the 'T' designator",
		}
	}

	nt, t, err := parseTime(b[n:])
	if err != nil {
		return time.Time{}, overrideUnexpectedTokenValue(err, b)
	}
	n += nt

	result := time.Date(dt.Year, dt.Month, dt.Day, t.Hour, t.Minute, t.Second, t.Nanosecond, time.UTC)
	// There is no offset. returns as UTC.
	if len(b) == n {
		return result.In(o.local), nil
	}
	if len(b) > n && !(b[n] == 'Z' || b[n] == '+' || b[n] == '-') {
		return time.Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b[n]),
			AfterToken: string(b[:n]),
			Expected:   "time zone format after time format",
		}
	}

	zone, err := ParseZone(b[n:])
	if err != nil {
		return time.Time{}, overrideUnexpectedTokenValue(err, b)
	}
	zoneOffset := zone.Offset()
	result = result.Add(-1 * time.Duration(zoneOffset) * time.Second)

	// Try to align this part with Go's time.Parse timezone handling as closely as possible.
	// Use local zone with the given offset if possible.
	_, offset, _, _, _ := lookup(o.local, result.Unix())
	if offset == zoneOffset {
		result = result.In(o.local)
	} else {
		result = result.In(time.FixedZone("", zoneOffset))
	}

	return result, nil
}
