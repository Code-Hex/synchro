package iso8601

import (
	"fmt"
	"time"
)

type DateLike interface {
	Date() Date
	IsValid() bool
	Validate() error
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

var _ DateLike = Date{}

func (d Date) Date() Date {
	return d
}

func (d Date) IsValid() bool {
	return d.Validate() == nil
}

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

type QuarterDate struct {
	Year    int
	Quarter int
	Day     int
}

var _ DateLike = QuarterDate{}

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

func (q QuarterDate) IsValid() bool {
	return q.Validate() == nil
}

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

type WeekDate struct {
	Year int
	Week int
	Day  int
}

var _ DateLike = WeekDate{}

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

func (w WeekDate) IsValid() bool {
	return w.Validate() == nil
}

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

type OrdinalDate struct {
	Year int
	Day  int
}

var _ DateLike = OrdinalDate{}

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

func (o OrdinalDate) IsValid() bool {
	return o.Validate() == nil
}

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

// DateLikeRangeError indicates that a value is not in an expected range for DateLike.
type DateLikeRangeError struct {
	Element string
	Value   int
	Year    int
	Min     int
	Max     int
}

func (e *DateLikeRangeError) Error() string {
	return fmt.Sprintf("iso8601: %d %s is not in range %d-%d in %d", e.Value, e.Element, e.Min, e.Max, e.Year)
}

type Time struct {
	Hour       int
	Minute     int
	Second     int
	Nanosecond int
}

func (t Time) Validate() error {
	if t.Minute > 59 {
		return &TimeRangeError{
			Element: "minute",
			Value:   t.Minute,
			Min:     0,
			Max:     59,
		}
	}
	if t.Second > 59 {
		return &TimeRangeError{
			Element: "second",
			Value:   t.Second,
			Min:     0,
			Max:     59,
		}
	}
	if t.Hour > 23 {
		if !(t.Hour == 24 && t.Minute == 0 && t.Second == 0 && t.Nanosecond == 0) {
			return &TimeRangeError{
				Element: "hour",
				Value:   t.Hour,
				Min:     0,
				Max:     24,
			}
		}
	}
	return nil
}

// TimeRangeError indicates that a value is not in an expected range for Time.
type TimeRangeError struct {
	Element string
	Value   int
	Min     int
	Max     int
}

func (e *TimeRangeError) Error() string {
	return fmt.Sprintf("iso8601: %d %s is not in range %d-%d", e.Value, e.Element, e.Min, e.Max)
}
