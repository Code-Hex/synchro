package iso8601

import (
	"fmt"
	"math"
)

/*
 *  Basic              Extended
 *  12                 N/A
 *  12.123456789       N/A
 *  12,123456789       N/A
 *  1230               12:30
 *  1230.123456789     12:30.123456789
 *  1230,123456789     12:30,123456789
 *  123045             12:30:45
 *  123045.123456789   12:30:45.123456789
 *  123045,123456789   12:30:45,123456789
 */
func ParseTime(b []byte) (Time, error) {
	n, t, err := parseTime(b)
	if err != nil {
		return Time{}, err
	}
	if len(b) != n {
		return Time{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b[n:]),
			AfterToken: string(b[:n]),
			Expected:   string(b[:n]),
		}
	}
	return t, nil
}

func parseTime(b []byte) (int, Time, error) {
	if len(b) > 2 && b[2] == ':' {
		return parseExtendedTime(b)
	}
	return parseBasicTime(b)
}

/*
 *  hh
 *  hh.fffffffff
 *  hh,fffffffff
 *  hh:mm
 *  hh:mm.fffffffff
 *  hh:mm,fffffffff
 *  hh:mm:ss
 *  hh:mm:ss.fffffffff
 *  hh:mm:ss,fffffffff
 */
func parseExtendedTime(b []byte) (int, Time, error) {
	var (
		h    int
		m    int
		s    int
		nsec int
	)

	parseFractionIfPresent := func(i int) (int, int) {
		if i < len(b) && (b[i] == '.' || b[i] == ',') {
			frac, digits := parseFraction(b[i+1:])
			i += digits + 1 // 1 == '.' or ','
			return frac, i
		}
		return 0, i
	}

	if c := countDigits(b, 0); c != 2 {
		return 0, Time{}, &UnexpectedTokenError{
			Value:    string(b),
			Token:    humanizeDigits(c),
			Expected: humanizeDigits(2),
		}
	}

	h = parseNumber(b, 0, 2)
	if len(b) < 3 || b[2] != ':' {
		nsec, n := parseFractionIfPresent(2)
		nsec *= 3600 // hour
		t, err := hmsfTime(h, m, s, nsec)
		return n, t, err
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
		nsec, n := parseFractionIfPresent(5)
		nsec *= 60 // hour
		t, err := hmsfTime(h, m, s, nsec)
		return n, t, err
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
	nsec, n := parseFractionIfPresent(8)
	t, err := hmsfTime(h, m, s, nsec)
	return n, t, err
}

/*
 *  hh
 *  hh.fffffffff
 *  hh,fffffffff
 *  hhmm
 *  hhmm.fffffffff
 *  hhmm,fffffffff
 *  hhmmss
 *  hhmmss.fffffffff
 *  hhmmss,fffffffff
 */
func parseBasicTime(b []byte) (int, Time, error) {
	var (
		h    int
		m    int
		s    int
		nsec int
	)
	n := countDigits(b, 0)
	switch n {
	case 2: // hh
		h = parseNumber(b, 0, 2)
	case 4: // hhmm
		h = parseNumber(b, 0, 2)
		m = parseNumber(b, 2, 2)
	case 6: // hhmmss
		h = parseNumber(b, 0, 2)
		m = parseNumber(b, 2, 2)
		s = parseNumber(b, 4, 2)
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

	// hh.fffffffff
	// hhmm.fffffffff
	// hhmmss.fffffffff
	if n < len(b) && (b[n] == '.' || b[n] == ',') {
		var digits int
		nsec, digits = parseFraction(b[n+1:])

		switch n {
		case 2: // hh
			nsec *= 3600
		case 4: // hhmm
			nsec *= 60
		}

		n += digits + 1 // 1 == '.' or ','
	}

	t, err := hmsfTime(h, m, s, nsec)
	return n, t, err
}

func parseFraction(b []byte) (int, int) {
	n := countDigits(b, 0)
	digits := n
	if digits > 9 {
		digits = 9 // nanoseconds
	}
	return parseNumber(b, 0, digits) * int(math.Pow10(9-digits)), n
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
