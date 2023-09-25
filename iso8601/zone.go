package iso8601

import (
	"fmt"

	"github.com/Code-Hex/synchro/internal/constraints"
)

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

// ParseZone attempts to parse a given byte slice representing a timezone offset
// in various supported ISO 8601 formats. Supported formats include:
//
//	Basic       Extended
//	Z           N/A
//	±hh         N/A
//	±hhmm       ±hh:mm
//	±hhmmss     ±hh:mm:ss
//
// The function returns a Zone structure or an error if the parsing fails.
func ParseZone[bytes constraints.Bytes](b bytes) (Zone, error) {
	if len(b) > 3 && b[3] == ':' {
		return parseExtendedZone([]byte(b))
	}
	return parseBasicZone([]byte(b))
}

/*
 *  Z
 *  ±hh
 *  ±hh:mm
 *  ±hh:mm:ss
 */
func parseExtendedZone(b []byte) (Zone, error) {
	var (
		h        int
		m        int
		s        int
		negative bool
	)
	if len(b) == 0 {
		return Zone{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b),
			AfterToken: "",
			Expected:   "Z or + or -",
		}
	}
	switch b[0] {
	case 'Z':
		if len(b) > 1 {
			return Zone{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[1:]),
				AfterToken: "Z",
				Expected:   fmt.Sprintf("non extra token (%s)", b[1:]),
			}
		}
		return timeZone(0, 0, 0, false)
	case '+':
	case '-':
		negative = true
	default:
		return Zone{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b),
			AfterToken: "",
			Expected:   "Z or + or -",
		}
	}

	if c := countDigits(b, 1); c != 2 {
		return Zone{}, &UnexpectedTokenError{
			Value:    string(b),
			Token:    humanizeDigits(c),
			Expected: humanizeDigits(2),
		}
	}

	h = parseNumber(b, 1, 2)
	if len(b) < 4 || b[3] != ':' {
		return timeZone(h, m, s, negative)
	}

	if c := countDigits(b, 4); c != 2 {
		return Zone{}, &UnexpectedTokenError{
			Value:      string(b),
			AfterToken: string(b[:4]),
			Token:      humanizeDigits(c),
			Expected:   humanizeDigits(2),
		}
	}

	m = parseNumber(b, 4, 2)
	if len(b) < 7 || b[6] != ':' {
		return timeZone(h, m, s, negative)
	}

	if c := countDigits(b, 7); c != 2 {
		return Zone{}, &UnexpectedTokenError{
			Value:      string(b),
			AfterToken: string(b[:7]),
			Token:      humanizeDigits(c),
			Expected:   humanizeDigits(2),
		}
	}

	s = parseNumber(b, 7, 2)
	return timeZone(h, m, s, negative)
}

/*
 *  Z
 *  ±hh
 *  ±hhmm
 *  ±hhmmss
 */
func parseBasicZone(b []byte) (Zone, error) {
	var (
		h        int
		m        int
		s        int
		negative bool
	)
	if len(b) == 0 {
		return Zone{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b),
			AfterToken: "",
			Expected:   "Z or + or -",
		}
	}
	switch b[0] {
	case 'Z':
		if len(b) > 1 {
			return Zone{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[1:]),
				AfterToken: "Z",
				Expected:   fmt.Sprintf("non extra token (%s)", b[1:]),
			}
		}
		return timeZone(0, 0, 0, false)
	case '+':
	case '-':
		negative = true
	default:
		return Zone{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b),
			AfterToken: "",
			Expected:   "Z or + or -",
		}
	}

	n := countDigits(b, 1)
	switch n {
	case 2: // ±hh
		h = parseNumber(b, 1, 2)
		return timeZone(h, m, s, negative)
	case 4: // ±hhmm
		h = parseNumber(b, 1, 2)
		m = parseNumber(b, 3, 2)
		return timeZone(h, m, s, negative)
	case 6: // ±hhmmss
		h = parseNumber(b, 1, 2)
		m = parseNumber(b, 3, 2)
		s = parseNumber(b, 5, 2)
		return timeZone(h, m, s, negative)
	default:
		return Zone{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      humanizeDigits(n),
			AfterToken: string(b[0]),
			Expected: fmt.Sprintf(
				"%s or %s or %s",
				humanizeDigits(2),
				humanizeDigits(4),
				humanizeDigits(6),
			),
		}
	}
}

func timeZone(h, m, s int, minus bool) (Zone, error) {
	z := Zone{
		Hour:     h,
		Minute:   m,
		Second:   s,
		Negative: minus,
	}
	if err := z.Validate(); err != nil {
		return Zone{}, err
	}
	return z, nil
}
