package iso8601

import (
	"encoding"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Code-Hex/synchro/internal/constraints"
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

// ParseDate attempts to parse a given byte slice representing a date in
// various supported ISO 8601 formats. Supported formats include:
//
//	Basic           Extended
//	20121224        2012-12-24    Calendar date   (ISO 8601)
//	2012359         2012-359      Ordinal date    (ISO 8601)
//	2012W521        2012-W52-1    Week date       (ISO 8601)
//	2012Q485        2012-Q4-85    Quarter date
//
// The function returns an implementation of DateLike or an error if the parsing fails.
func ParseDate[bytes constraints.Bytes](b bytes) (DateLike, error) {
	n, d, err := parseDate([]byte(b))
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

// DateLike defines an interface for date-related structures.
// It provides methods for retrieving the date, validating the date,
// and checking if the date is valid.
type DateLike interface {
	// Date returns the underlying Date value.
	Date() Date

	// IsValid checks whether the date is valid.
	IsValid() bool

	// Validate checks the correctness of the date and returns an error if it's invalid.
	Validate() error
}

// Date represents a calendar date with year, month, and day components.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

var _ interface {
	DateLike
	fmt.Stringer
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*Date)(nil)

// DateOf returns the iso8601 Date in which a time occurs in that time's location.
func DateOf(t time.Time) Date {
	var d Date
	d.Year, d.Month, d.Day = t.Date()
	return d
}

// String returns the ISO8601 string representation of the format "YYYY-MM-DD".
// For example: "2012-12-01".
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// Date returns itself as it directly represents a date.
func (d Date) Date() Date {
	return d
}

// IsValid checks if the date is valid based on its year, month, and day values.
func (d Date) IsValid() bool {
	return d.Validate() == nil
}

// Validate checks the individual components of the date (year, month, and day)
// and returns an error if any of them are out of the expected ranges.
func (d Date) Validate() error {
	if d.Year < 0 || d.Year > 9999 {
		return &DateLikeRangeError{
			Element: "year",
			Value:   d.Year,
			Year:    d.Year,
			Min:     0,
			Max:     9999,
		}
	}
	if d.Month < 1 || d.Month > 12 {
		return &DateLikeRangeError{
			Element: "month",
			Value:   int(d.Month),
			Year:    d.Year,
			Min:     1,
			Max:     12,
		}
	}
	daysInMonth := daysInMonth(d.Year, int(d.Month))
	if d.Day < 1 || d.Day > daysInMonth {
		return &DateLikeRangeError{
			Element: "day of month",
			Value:   d.Day,
			Year:    d.Year,
			Min:     1,
			Max:     daysInMonth,
		}
	}
	return nil
}

// StdTime converts the Date structure to a time.Time object, using UTC for the time.
func (d Date) StdTime() time.Time {
	return d.In(time.UTC)
}

// In returns the time corresponding to time 00:00:00 of the date in the location.
//
// In is always consistent with time.Date, even when time.Date returns a time
// on a different day. For example, if loc is America/Indiana/Vincennes, then both
//
//	time.Date(1955, time.May, 1, 0, 0, 0, 0, loc)
//
// and
//
//	iso8601.Date{Year: 1955, Month: time.May, Day: 1}.In(loc)
//
// return 23:00:00 on April 30, 1955.
//
// In panics if loc is nil.
func (d Date) In(loc *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, loc)
}

// AddDays returns the date that is n days in the future.
// n can also be negative to go into the past.
func (d Date) AddDays(n int) Date {
	return DateOf(d.StdTime().AddDate(0, 0, n))
}

// DaysSince returns the signed number of days between the date and s, not including the end day.
// This is the inverse operation to AddDays.
func (d Date) DaysSince(s Date) (days int) {
	// We convert to Unix time so we do not have to worry about leap seconds:
	// Unix time increases by exactly 86400 seconds per day.
	deltaUnix := d.StdTime().Unix() - s.StdTime().Unix()
	return int(deltaUnix / 86400)
}

// Before reports whether d occurs before d2.
func (d Date) Before(d2 Date) bool {
	if d.Year != d2.Year {
		return d.Year < d2.Year
	}
	if d.Month != d2.Month {
		return d.Month < d2.Month
	}
	return d.Day < d2.Day
}

// After reports whether d occurs after d2.
func (d Date) After(d2 Date) bool {
	return d2.Before(d)
}

// IsZero reports whether date fields are set to their default value.
func (d Date) IsZero() bool {
	return (d.Year == 0) && (int(d.Month) == 0) && (d.Day == 0)
}

// QuarterDate converts a Date to a QuarterDate.
// It calculates the quarter of the year and the day within that quarter for the given date.
func (d Date) QuarterDate() QuarterDate {
	t := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
	yday := t.YearDay() // 1 ~ 366

	quarter := 1
	for yday > daysInQuarter(d.Year, quarter) {
		yday -= daysInQuarter(d.Year, quarter)
		quarter++
	}

	return QuarterDate{
		Year:    d.Year,
		Quarter: quarter,
		Day:     yday,
	}
}

// WeekDate converts a Date to a WeekDate.
// It calculates the week of the year and the day within that week for the given date.
func (d Date) WeekDate() WeekDate {
	t := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
	year, week := t.ISOWeek()
	day := int(t.Weekday())
	if day == 0 {
		day = 7
	}
	return WeekDate{
		Year: year,
		Week: week,
		Day:  day,
	}
}

// OrdinalDate converts a Date to an OrdinalDate.
// It calculates the day of the year for the given date.
func (d Date) OrdinalDate() OrdinalDate {
	t := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
	yday := t.YearDay() // 1 ~ 366
	return OrdinalDate{
		Year: d.Year,
		Day:  yday,
	}
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of t.String().
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be a string in a format accepted by ParseDate.
func (d *Date) UnmarshalText(data []byte) error {
	var err error
	dt, err := ParseDate(data)
	if err != nil {
		return err
	}
	*d = dt.Date()
	return nil
}

// QuarterDate represents a date within a specific quarter of a year.
// It includes the year, quarter (from 1 to 4), and day within that quarter.
type QuarterDate struct {
	Year    int
	Quarter int
	Day     int
}

var _ interface {
	DateLike
	fmt.Stringer
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*QuarterDate)(nil)

// String returns the ISO8601 string representation of the format "YYYY-QX-DD".
// For example: "2012-Q4-85".
func (q QuarterDate) String() string {
	return fmt.Sprintf("%04d-Q%d-%02d", q.Year, q.Quarter, q.Day)
}

// Date converts a QuarterDate into the standard Date representation.
// It calculates the exact calendar date based on the year, quarter, and day within that quarter.
func (q QuarterDate) Date() Date {
	yday := q.Day // 1 ~ 366
	for i := 1; i < q.Quarter; i++ {
		yday += daysInQuarter(q.Year, i)
	}
	t := time.Date(q.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, yday-1)
	return Date{
		Year:  t.Year(),
		Month: t.Month(),
		Day:   t.Day(),
	}
}

// IsValid checks if the quarter date is valid based on its year, quarter, and day within the quarter values.
func (q QuarterDate) IsValid() bool {
	return q.Validate() == nil
}

// Validate checks the individual components of the quarter date (year, quarter, and day within the quarter)
// and returns an error if any of them are out of the expected ranges.
func (q QuarterDate) Validate() error {
	if q.Year < 0 || q.Year > 9999 {
		return &DateLikeRangeError{
			Element: "year",
			Value:   q.Year,
			Year:    q.Year,
			Min:     0,
			Max:     9999,
		}
	}
	if q.Quarter < 1 || q.Quarter > 4 {
		return &DateLikeRangeError{
			Element: "quarter",
			Value:   q.Quarter,
			Year:    q.Year,
			Min:     1,
			Max:     4,
		}
	}
	daysInQuarter := daysInQuarter(q.Year, q.Quarter)
	if q.Day < 1 || q.Day > daysInQuarter {
		return &DateLikeRangeError{
			Element: "day of quarter",
			Value:   q.Day,
			Year:    q.Year,
			Min:     1,
			Max:     daysInQuarter,
		}
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of t.String().
func (q QuarterDate) MarshalText() ([]byte, error) {
	return []byte(q.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The quarter date is expected to be a string in a format accepted by ParseDate.
func (q *QuarterDate) UnmarshalText(data []byte) error {
	var err error
	dt, err := ParseDate(data)
	if err != nil {
		return err
	}
	qd, ok := dt.(QuarterDate)
	if !ok {
		*q = dt.Date().QuarterDate()
		return nil
	}
	*q = qd
	return nil
}

// WeekDate represents a date within a specific week of a given year,
// following the ISO 8601 week-date system. It includes the year,
// week number (1 to 52 or 53), and day of the week (1 for Monday to 7 for Sunday).
type WeekDate struct {
	Year int
	Week int
	Day  int
}

var _ interface {
	DateLike
	fmt.Stringer
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*WeekDate)(nil)

// String returns the ISO8601 string representation of the format "YYYY-WX-DD".
// For example: "2012-W52-1".
func (w WeekDate) String() string {
	return fmt.Sprintf("%04d-W%02d-%d", w.Year, w.Week, w.Day)
}

// Date converts a WeekDate into the standard Date representation.
// It calculates the exact calendar date based on the year, week number, and day of the week.
func (w WeekDate) Date() Date {
	// Find the first Thursday of the given year. This will be in the first week of the year according to ISO 8601.
	thursday := time.Date(w.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	for thursday.Weekday() != time.Thursday {
		thursday = thursday.AddDate(0, 0, 1)
	}

	// Calculate the date of the Monday of week 1
	monday := thursday.AddDate(0, 0, -3)

	// Calculate the date corresponding to the given week and day
	t := monday.AddDate(0, 0, (w.Week-1)*7+w.Day-1)
	return Date{
		Year:  t.Year(),
		Month: t.Month(),
		Day:   t.Day(),
	}
}

// IsValid checks if the week date is valid based on its year, week number, and day of the week values.
func (w WeekDate) IsValid() bool {
	return w.Validate() == nil
}

// Validate checks the individual components of the week date (year, week number, and day of the week)
// and returns an error if any of them are out of the expected ranges.
func (w WeekDate) Validate() error {
	if w.Year < 0 || w.Year > 9999 {
		return &DateLikeRangeError{
			Element: "year",
			Value:   w.Year,
			Year:    w.Year,
			Min:     0,
			Max:     9999,
		}
	}
	if w.Day < 1 || w.Day > 7 {
		return &DateLikeRangeError{
			Element: "day of week",
			Value:   int(w.Day),
			Year:    w.Year,
			Min:     1,
			Max:     7,
		}
	}
	weeksInYear := weeksInYear(w.Year)
	if w.Week < 1 || w.Week > weeksInYear {
		return &DateLikeRangeError{
			Element: "week",
			Value:   w.Week,
			Year:    w.Year,
			Min:     1,
			Max:     weeksInYear,
		}
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of t.String().
func (w WeekDate) MarshalText() ([]byte, error) {
	return []byte(w.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The week date is expected to be a string in a format accepted by ParseDate.
func (w *WeekDate) UnmarshalText(data []byte) error {
	var err error
	dt, err := ParseDate(data)
	if err != nil {
		return err
	}
	wd, ok := dt.(WeekDate)
	if !ok {
		*w = dt.Date().WeekDate()
		return nil
	}
	*w = wd
	return nil
}

// OrdinalDate represents a date specified by its year and the day-of-year (ordinal date),
// where the day-of-year ranges from 1 through 365 (or 366 in a leap year).
type OrdinalDate struct {
	Year int
	Day  int
}

var _ interface {
	DateLike
	fmt.Stringer
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*OrdinalDate)(nil)

// String returns the ISO8601 string representation of the format "YYYY-DDD".
// For example: "2012-359".
func (o OrdinalDate) String() string {
	return fmt.Sprintf("%04d-%03d", o.Year, o.Day)
}

// Date converts an OrdinalDate into the standard Date representation.
// It calculates the exact calendar date based on the year and the day-of-year.
func (o OrdinalDate) Date() Date {
	yday := o.Day // 1 ~ 366
	t := time.Date(o.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, yday-1)
	return Date{
		Year:  o.Year,
		Month: t.Month(),
		Day:   t.Day(),
	}
}

// IsValid checks if the ordinal date is valid based on its year and day-of-year values.
func (o OrdinalDate) IsValid() bool {
	return o.Validate() == nil
}

// Validate checks the individual components of the ordinal date (year and day-of-year)
// and returns an error if any of them are out of the expected ranges.
func (o OrdinalDate) Validate() error {
	if o.Year < 0 || o.Year > 9999 {
		return &DateLikeRangeError{
			Element: "year",
			Value:   o.Year,
			Year:    o.Year,
			Min:     0,
			Max:     9999,
		}
	}
	daysInYear := daysInYear(o.Year)
	if o.Day < 1 || o.Day > daysInYear {
		return &DateLikeRangeError{
			Element: "day of year",
			Value:   o.Day,
			Year:    o.Year,
			Min:     1,
			Max:     daysInYear,
		}
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of t.String().
func (o OrdinalDate) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The week date is expected to be a string in a format accepted by ParseDate.
func (o *OrdinalDate) UnmarshalText(data []byte) error {
	var err error
	dt, err := ParseDate(data)
	if err != nil {
		return err
	}
	od, ok := dt.(OrdinalDate)
	if !ok {
		*o = dt.Date().OrdinalDate()
		return nil
	}
	*o = od
	return nil
}

// DateLikeRangeError indicates that a value is not in an expected range for DateLike.
type DateLikeRangeError struct {
	Element string
	Value   int
	Year    int
	Min     int
	Max     int
}

// Error implements the error interface.
func (e *DateLikeRangeError) Error() string {
	return fmt.Sprintf("iso8601: %d %s is not in range %d-%d in %d", e.Value, e.Element, e.Min, e.Max, e.Year)
}
