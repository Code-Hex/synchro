package iso8601

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// NOTE(codehex): "math.MaxInt == 9223372036854775807" has 19 digits.
// So I consider the maximum to be 18 digits, which is "999999999999999999."

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
	for i := width; i > 0; i-- {
		v += int(b[start]-'0') * int(math.Pow10(i-1))
		start++
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
	n, d, err := parseDate(b)
	if err != nil {
		return nil, err
	}
	if len(b) != n {
		return nil, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b[n:]),
			AfterToken: string(b[:n]),
			Expected:   string(b[:n]),
		}
	}
	return d, err
}

func parseDate(b []byte) (int, DateLike, error) {
	var (
		y int
		x int // month or week or quarter
		d int
	)

	// To allow leading '+' signed year components.
	signed := 0
	if len(b) > 0 && b[0] == '+' {
		b = b[1:]
		signed++
	}

	n := countDigits(b, 0)
	switch n {
	case 4: /* 2012 (year) */
		y = parseNumber(b, 0, 4)
		if len(b) < 8 {
			return 0, nil, &UnexpectedTokenError{
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
				return 0, nil, &UnexpectedTokenError{
					Value:      string(b),
					Token:      humanizeDigits(n),
					AfterToken: "Q",
					Expected:   humanizeDigits(3),
				}
			}
			x = parseNumber(b, 5, 1)
			d = parseNumber(b, 6, 2)
			dt, err := yqdISODate(y, x, d)
			return 8 + signed, dt, err
		case 'W': // 2012W521
			if n != 3 {
				return 0, nil, &UnexpectedTokenError{
					Value:      string(b),
					Token:      humanizeDigits(n),
					AfterToken: "W",
					Expected:   humanizeDigits(3),
				}
			}
			x = parseNumber(b, 5, 2)
			d = parseNumber(b, 7, 1)
			dt, err := ywdISODate(y, x, d)
			return 8 + signed, dt, err
		default:
			return 0, nil, &UnexpectedTokenError{
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
						return 0, nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      humanizeDigits(n),
							AfterToken: "Q",
							Expected:   humanizeDigits(1),
						}
					}
					x = parseNumber(b, 6, 1)
					if b[7] != '-' {
						return 0, nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      string(b[7]),
							AfterToken: fmt.Sprintf("Q%d", x),
							Expected:   "-",
						}
					}
					if c := countDigits(b, 8); c != 2 {
						return 0, nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      humanizeDigits(c),
							AfterToken: fmt.Sprintf("Q%d-", x),
							Expected:   humanizeDigits(2),
						}
					}
					d = parseNumber(b, 8, 2)
					dt, err := yqdISODate(y, x, d)
					return 10 + signed, dt, err
				case 'W': // 2012-W52-1
					if n != 2 {
						return 0, nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      humanizeDigits(n),
							AfterToken: "W",
							Expected:   humanizeDigits(2),
						}
					}
					x = parseNumber(b, 6, 2)
					if b[8] != '-' {
						return 0, nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      string(b[8]),
							AfterToken: fmt.Sprintf("W%02d", x),
							Expected:   "-",
						}
					}
					if c := countDigits(b, 9); c != 1 {
						return 0, nil, &UnexpectedTokenError{
							Value:      string(b),
							Token:      humanizeDigits(c),
							AfterToken: fmt.Sprintf("W%02d-", x),
							Expected:   humanizeDigits(1),
						}
					}
					d = parseNumber(b, 9, 1)
					dt, err := ywdISODate(y, x, d)
					return 10 + signed, dt, err
				}
			}
		case 2: // 2012-12-24
			x = parseNumber(b, 5, 2)
			if b[7] != '-' {
				return 0, nil, &UnexpectedTokenError{
					Value:      string(b),
					Token:      string(b[7]),
					AfterToken: fmt.Sprintf("-%02d", x),
					Expected:   "-",
				}
			}
			if c := countDigits(b, 8); c != 2 {
				return 0, nil, &UnexpectedTokenError{
					Value:      string(b),
					Token:      humanizeDigits(c),
					AfterToken: fmt.Sprintf("-%02d-", x),
					Expected:   humanizeDigits(2),
				}
			}
			d = parseNumber(b, 8, 2)
			dt, err := ymdISODate(y, x, d)
			return 10 + signed, dt, err
		case 3: // 2012-359
			d = parseNumber(b, 5, 3)
			dt, err := ydISODate(y, d)
			return 8 + signed, dt, err
		default:
			return 0, nil, &UnexpectedTokenError{
				Value:      string(b),
				Token:      humanizeDigits(n),
				AfterToken: fmt.Sprintf("%d-", y),
				Expected:   "like -Q4-85 or -W52-1 or -359",
			}
		}
	case 7: // 2012359 (basic ordinal date)
		y = parseNumber(b, 0, 4)
		d = parseNumber(b, 4, 3)
		dt, err := ydISODate(y, d)
		return 7 + signed, dt, err
	case 8: // 20121224 (basic calendar date)
		y = parseNumber(b, 0, 4)
		x = parseNumber(b, 4, 2)
		d = parseNumber(b, 6, 2)
		dt, err := ymdISODate(y, x, d)
		return 8 + signed, dt, err
	default:
	}
	return 0, nil, &UnexpectedTokenError{
		Value:      string(b),
		Token:      humanizeDigits(n),
		AfterToken: "",
		Expected:   "date format",
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
