package iso8601

import (
	"fmt"
	"math"
)

/*
 *  Basic              Extended
 *  12                 N/A
 *  1230               12:30
 *  123045             12:30:45
 *  123045.123456789   12:30:45.123456789
 *  123045,123456789   12:30:45,123456789
 */
func ParseTime(b []byte) (Time, error) {
	if len(b) > 2 && b[2] == ':' {
		return parseExtendedTime(b)
	}
	return parseBasicTime(b)
}

/*
 *  hh
 *  hh:mm
 *  hh:mm:ss
 *  hh:mm:ss.fffffffff
 *  hh:mm:ss,fffffffff
 */
func parseExtendedTime(b []byte) (Time, error) {
	var (
		h int
		m int
		s int
		f int
	)
	if c := countDigits(b, 0); c != 2 {
		return Time{}, &UnexpectedTokenError{
			Value:    string(b),
			Token:    humanizeDigits(c),
			Expected: humanizeDigits(2),
		}
	}

	h = parseNumber(b, 0, 2)
	if len(b) < 3 || b[2] != ':' {
		return hmsfTime(h, m, s, f)
	}

	if c := countDigits(b, 3); c != 2 {
		return Time{}, &UnexpectedTokenError{
			Value:      string(b),
			AfterToken: string(b[:3]),
			Token:      humanizeDigits(c),
			Expected:   humanizeDigits(2),
		}
	}

	m = parseNumber(b, 3, 2)
	if len(b) < 6 || b[5] != ':' {
		return hmsfTime(h, m, s, f)
	}

	if c := countDigits(b, 6); c != 2 {
		return Time{}, &UnexpectedTokenError{
			Value:      string(b),
			AfterToken: string(b[:6]),
			Token:      humanizeDigits(c),
			Expected:   humanizeDigits(2),
		}
	}

	s = parseNumber(b, 6, 2)

	// hh:mm:ss.fffffffff
	if 8 < len(b) && (b[8] == '.' || b[8] == ',') {
		f = parseFraction(b[9:])
	}

	return hmsfTime(h, m, s, f)
}

/*
 *  hh
 *  hhmm
 *  hhmmss
 *  hhmmss.fffffffff
 *  hhmmss,fffffffff
 */
func parseBasicTime(b []byte) (Time, error) {
	var (
		h int
		m int
		s int
		f int
	)
	n := countDigits(b, 0)
	switch n {
	case 2: // hh
		h = parseNumber(b, 0, 2)
		return hmsfTime(h, m, s, f)
	case 4: // hhmm
		h = parseNumber(b, 0, 2)
		m = parseNumber(b, 2, 2)
		return hmsfTime(h, m, s, f)
	case 6: // hhmmss
		h = parseNumber(b, 0, 2)
		m = parseNumber(b, 2, 2)
		s = parseNumber(b, 4, 2)

		// hhmmss.fffffffff
		if 7 < len(b) && (b[6] == '.' || b[6] == ',') {
			f = parseFraction(b[7:])
		}
		return hmsfTime(h, m, s, f)
	default:
		afterToken := ""
		if n > 2 {
			h = parseNumber(b, 0, 2)
			afterToken = fmt.Sprintf("%02d", h)
		}
		if n > 4 {
			m = parseNumber(b, 2, 2)
			afterToken = fmt.Sprintf("%02d%02d", h, m)
		}
		return Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      humanizeDigits(n),
			AfterToken: afterToken,
			Expected:   humanizeDigits(2),
		}
	}
}

func parseFraction(b []byte) int {
	digits := countDigits(b, 0)
	if digits > 9 {
		digits = 9 // nanoseconds
	}
	return parseNumber(b, 0, digits) * int(math.Pow10(9-digits))
}

func hmsfTime(h, m, s, f int) (Time, error) {
	t := Time{
		Hour:       h,
		Minute:     m,
		Second:     s,
		Nanosecond: f,
	}
	if err := t.Validate(); err != nil {
		return Time{}, err
	}
	return t, nil
}
