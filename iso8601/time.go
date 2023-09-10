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
	_, t, err := parseTime(b)
	return t, err
}

func parseTime(b []byte) (int, Time, error) {
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
func parseExtendedTime(b []byte) (int, Time, error) {
	var (
		h int
		m int
		s int
		f int
	)
	if c := countDigits(b, 0); c != 2 {
		return 0, Time{}, &UnexpectedTokenError{
			Value:    string(b),
			Token:    humanizeDigits(c),
			Expected: humanizeDigits(2),
		}
	}

	h = parseNumber(b, 0, 2)
	if len(b) < 3 || b[2] != ':' {
		t, err := hmsfTime(h, m, s, f)
		return 2, t, err
	}

	if c := countDigits(b, 3); c != 2 {
		return 0, Time{}, &UnexpectedTokenError{
			Value:      string(b),
			AfterToken: string(b[:3]),
			Token:      humanizeDigits(c),
			Expected:   humanizeDigits(2),
		}
	}

	m = parseNumber(b, 3, 2)
	if len(b) < 6 || b[5] != ':' {
		t, err := hmsfTime(h, m, s, f)
		return 5, t, err
	}

	if c := countDigits(b, 6); c != 2 {
		return 0, Time{}, &UnexpectedTokenError{
			Value:      string(b),
			AfterToken: string(b[:6]),
			Token:      humanizeDigits(c),
			Expected:   humanizeDigits(2),
		}
	}

	s = parseNumber(b, 6, 2)
	n := 8

	// hh:mm:ss.fffffffff
	if 8 < len(b) && (b[8] == '.' || b[8] == ',') {
		var digits int
		f, digits = parseFraction(b[9:])
		n += digits + 1 // 1 == '.' or ','
	}

	t, err := hmsfTime(h, m, s, f)
	return n, t, err
}

/*
 *  hh
 *  hhmm
 *  hhmmss
 *  hhmmss.fffffffff
 *  hhmmss,fffffffff
 */
func parseBasicTime(b []byte) (int, Time, error) {
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
		t, err := hmsfTime(h, m, s, f)
		return 2, t, err
	case 4: // hhmm
		h = parseNumber(b, 0, 2)
		m = parseNumber(b, 2, 2)
		t, err := hmsfTime(h, m, s, f)
		return 4, t, err
	case 6: // hhmmss
		h = parseNumber(b, 0, 2)
		m = parseNumber(b, 2, 2)
		s = parseNumber(b, 4, 2)

		n := 6

		// hhmmss.fffffffff
		if 7 < len(b) && (b[6] == '.' || b[6] == ',') {
			var digits int
			f, digits = parseFraction(b[7:])
			n += digits + 1
		}
		t, err := hmsfTime(h, m, s, f)
		return n, t, err
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
		return 0, Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      humanizeDigits(n),
			AfterToken: afterToken,
			Expected:   humanizeDigits(2),
		}
	}
}

func parseFraction(b []byte) (int, int) {
	digits := countDigits(b, 0)
	if digits > 9 {
		digits = 9 // nanoseconds
	}
	return parseNumber(b, 0, digits) * int(math.Pow10(9-digits)), digits
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
