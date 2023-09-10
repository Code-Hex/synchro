package iso8601

import (
	"fmt"
	"strings"
	"time"
)

var defaultParseDateTimeOptions = parseDateTimeOptions{
	timeDesignators: []byte{'T'},
}

type parseDateTimeOptions struct {
	timeDesignators []byte
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
// The function returns a time.Time structure representing the parsed date-time, adjusted
// for the parsed timezone offset if provided. If no timezone is specified, the time is
// returned in UTC.
//
// In the absence of a time zone indicator, Parse returns a time in UTC.
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
	if len(b) == n {
		return result, nil
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

	// Try to align this part with Go's time.Parse timezone handling as closely as possible.
	offset := zone.Offset()
	result = result.Add(-1 * time.Duration(offset) * time.Second)

	return result.In(time.FixedZone("", offset)), nil
}
