package iso8601

import (
	"fmt"
	"strconv"
	"time"
)

func countDigits(b []byte, i int) int {
	start := i
	for ; i < len(b); i++ {
		c := b[i] - '0'
		if c > 9 {
			break
		}
	}
	return i - start
}

func parseNumber(b []byte, start, width int) (v int) {
	if len(b) <= start {
		return
	}
	switch width {
	case 9:
		v += int(b[start]-'0') * 100000000
		start++
		fallthrough
	case 8:
		v += int(b[start]-'0') * 10000000
		start++
		fallthrough
	case 7:
		v += int(b[start]-'0') * 1000000
		start++
		fallthrough
	case 6:
		v += int(b[start]-'0') * 100000
		start++
		fallthrough
	case 5:
		v += int(b[start]-'0') * 10000
		start++
		fallthrough
	case 4:
		v += int(b[start]-'0') * 1000
		start++
		fallthrough
	case 3:
		v += int(b[start]-'0') * 100
		start++
		fallthrough
	case 2:
		v += int(b[start]-'0') * 10
		start++
		fallthrough
	case 1:
		v += int(b[start] - '0')
		start++
		fallthrough
	default:
	}
	return
}

/*
 *  Basic      Extended
 *  20121224   2012-12-24   Calendar date   (ISO 8601)
 *  2012359    2012-359     Ordinal date    (ISO 8601)
 *  2012W521   2012-W52-1   Week date       (ISO 8601)
 *  2012Q485   2012-Q4-85   Quarter date
 */
func ParseDate(b []byte) (DateLike, error) {
	var (
		y int
		x int // month or week or quarter
		d int
	)
	n := countDigits(b, 0)
	switch n {
	case 4: /* 2012 (year) */
		y = parseNumber(b, 0, 4)
		if len(b) < 8 {
			return nil, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[4:]),
				AfterToken: strconv.Itoa(y),
				Expected:   "8 or more characters",
			}
		}

		n = countDigits(b, 5)
		switch b[4] {
		case '-': // 2012-359 | 2012-12-24 | 2012-W52-1 | 2012-Q4-85
		case 'Q': // 2012Q485
			if n != 3 {
				return nil, &UnexpectedTokenError{
					Value:      string(b),
					Token:      humanizeDigits(n),
					AfterToken: "Q",
					Expected:   humanizeDigits(3),
				}
			}
			x = parseNumber(b, 5, 1)
			d = parseNumber(b, 6, 2)
			return yqdISODate(y, x, d)
		case 'W': // 2012W521
			if n != 3 {
				return nil, &UnexpectedTokenError{
					Value:      string(b),
					Token:      humanizeDigits(n),
					AfterToken: "W",
					Expected:   humanizeDigits(3),
				}
			}
			x = parseNumber(b, 5, 2)
			d = parseNumber(b, 7, 1)
			return ywdISODate(y, x, d)
		default:
			return nil, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[4:]),
				AfterToken: strconv.Itoa(y),
				Expected:   "- or Q or W",
			}
		}

		switch n {
		case 0: // 2012-Q4-85 | 2012-W52-1
			if len(b) >= 10 {
				n = countDigits(b, 6)
				switch b[5] {
				case 'Q': // 2012-Q4-85
					if n != 1 {
						return nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      humanizeDigits(n),
							AfterToken: "Q",
							Expected:   humanizeDigits(1),
						}
					}
					x = parseNumber(b, 6, 1)
					if b[7] != '-' {
						return nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      string(b[7]),
							AfterToken: fmt.Sprintf("Q%d", x),
							Expected:   "-",
						}
					}
					if c := countDigits(b, 8); c != 2 {
						return nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      humanizeDigits(c),
							AfterToken: fmt.Sprintf("Q%d-", x),
							Expected:   humanizeDigits(2),
						}
					}
					d = parseNumber(b, 8, 2)
					return yqdISODate(y, x, d)
				case 'W': // 2012-W52-1
					if n != 2 {
						return nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      humanizeDigits(n),
							AfterToken: "W",
							Expected:   humanizeDigits(2),
						}
					}
					x = parseNumber(b, 6, 2)
					if b[8] != '-' {
						return nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      string(b[8]),
							AfterToken: fmt.Sprintf("W%02d", x),
							Expected:   "-",
						}
					}
					if c := countDigits(b, 9); c != 1 {
						return nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      humanizeDigits(c),
							AfterToken: fmt.Sprintf("W%02d-", x),
							Expected:   humanizeDigits(1),
						}
					}
					d = parseNumber(b, 9, 1)
					return ywdISODate(y, x, d)
				}
			}
		case 2: // 2012-12-24
			x = parseNumber(b, 5, 2)
			if b[7] != '-' {
				return nil, &UnexpectedTokenError{
					Value:      string(b),
					Token:      string(b[7]),
					AfterToken: fmt.Sprintf("-%02d", x),
					Expected:   "-",
				}
			}
			if c := countDigits(b, 8); c != 2 {
				return nil, &UnexpectedTokenError{
					Value:      string(b),
					Token:      humanizeDigits(c),
					AfterToken: fmt.Sprintf("-%02d-", x),
					Expected:   humanizeDigits(2),
				}
			}
			d = parseNumber(b, 8, 2)
			n = 10
			return ymdISODate(y, x, d)
		case 3: // 2012-359
			d = parseNumber(b, 5, 3)
			return ydISODate(y, d)
		default:
			return nil, &UnexpectedTokenError{
				Value:      string(b),
				Token:      humanizeDigits(n),
				AfterToken: fmt.Sprintf("%d-", y),
				Expected:   "like -Q4-85 or -W52-1 or -359",
			}
		}
	case 7: // 2012359 (basic ordinal date)
		y = parseNumber(b, 0, 4)
		d = parseNumber(b, 4, 3)
		return ydISODate(y, d)
	case 8: // 20121224 (basic calendar date)
		y = parseNumber(b, 0, 4)
		x = parseNumber(b, 4, 2)
		d = parseNumber(b, 6, 2)
		return ymdISODate(y, x, d)
	default:
	}
	return nil, &UnexpectedTokenError{
		Value:      string(b),
		Token:      humanizeDigits(n),
		AfterToken: "",
		Expected:   "",
	}
}

func humanizeDigits(n int) string {
	if n <= 1 {
		return fmt.Sprintf("%d-digit", n)
	}
	return fmt.Sprintf("%d-digits", n)
}

func ydISODate(y int, d int) (DateLike, error) {
	yd := OrdinalDate{
		Year: y,
		Day:  d,
	}
	if err := yd.Validate(); err != nil {
		return nil, err
	}
	return yd, nil
}

func ymdISODate(y int, m int, d int) (DateLike, error) {
	ymd := Date{
		Year:  y,
		Month: time.Month(m),
		Day:   d,
	}
	if err := ymd.Validate(); err != nil {
		return nil, err
	}
	return ymd, nil
}

func yqdISODate(y int, q int, d int) (DateLike, error) {
	yqd := QuarterDate{
		Year:    y,
		Quarter: q,
		Day:     d,
	}
	if err := yqd.Validate(); err != nil {
		return nil, err
	}
	return yqd, nil
}

func ywdISODate(y int, w int, d int) (DateLike, error) {
	ywd := WeekDate{
		Year: y,
		Week: w,
		Day:  d,
	}
	if err := ywd.Validate(); err != nil {
		return nil, err
	}
	return ywd, nil
}

func daysInYear(y int) int {
	if isLeapYear(y) {
		return 366
	}
	return 365
}

// daysInQuarterList is the number of days for non-leap years in each quarter
var daysInQuarterList = [...]int{0, 90, 91, 92, 92}

func daysInQuarter(y int, q int) int {
	if q == 1 && isLeapYear(y) {
		return 91
	}
	return daysInQuarterList[q]
}

// daysInMonthList is the number of days for non-leap years in each calendar month
var daysInMonthList = [...]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func daysInMonth(y int, m int) int {
	if m == 2 && isLeapYear(y) {
		return 29
	}
	return daysInMonthList[m]
}

func weeksInYear(year int) int {
	if year < 1 {
		year += 400 * (1 - year/400)
	}
	y := year - 1
	d := (y + y/4 - y/100 + y/400) % 7 // [0=Mon, 6=Sun]
	if d == 3 || (d == 2 && isLeapYear(year)) {
		return 53
	}
	return 52
}

func isLeapYear(y int) bool {
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}
