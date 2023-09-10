package iso8601

import (
	"time"
)

// ParseDateTime attempts to parse a given byte slice representing combined date, time,
// and optionally timezone offset in supported ISO 8601 formats. Supported formats include:
//
//	Basic                        Extended
//	20070301                     2007-03-01
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
func ParseDateTime(b []byte) (time.Time, error) {
	n, d, err := parseDate(b)
	if err != nil {
		return time.Time{}, overrideUnexpectedTokenValue(err, b)
	}
	dt := d.Date()
	if len(b) == n {
		return dt.StdTime(), nil
	}

	if len(b) > n && b[n] != 'T' {
		return time.Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b[n]),
			AfterToken: string(b[:n]),
			Expected:   "T",
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
