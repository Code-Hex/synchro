package iso8601

import "fmt"

func ParseZone(b []byte) (Zone, error) {
	if len(b) > 3 && b[3] == ':' {
		return parseExtendedZone(b)
	}
	return parseBasicZone(b)
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
