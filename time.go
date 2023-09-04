package synchro

import (
	"math"
	"time"
)

type empty[T TimeZone] struct{}

type Time[T TimeZone] struct {
	tm time.Time
	_  empty[T]
}

// StdTime returns the time.Time.
func (t Time[T]) StdTime() time.Time {
	return t.tm
}

// StartOfYear returns Time for start of the year.
func (t Time[T]) StartOfYear() Time[T] {
	return New[T](t.Year(), 1, 1, 0, 0, 0, 0)
}

// EndOfYear returns Time for end of the year.
func (t Time[T]) EndOfYear() Time[T] {
	return New[T](t.Year(), 12, 31, 23, 59, 59, 999999999)
}

// StartOfMonth returns Time for start of the month.
func (t Time[T]) StartOfMonth() Time[T] {
	return New[T](t.Year(), t.Month(), 1, 0, 0, 0, 0)
}

// EndOfMonth returns Time for end of the month.
func (t Time[T]) EndOfMonth() Time[T] {
	startOfMonth := t.StartOfMonth()
	return startOfMonth.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)
}

// StartOfWeek returns Time for start of the week.
func (t Time[T]) StartOfWeek() Time[T] {
	dayOfWeek := t.Weekday()
	return t.Add(-time.Duration(dayOfWeek) * 24 * time.Hour)
}

// EndOfWeek returns Time for end of the week.
func (t Time[T]) EndOfWeek() Time[T] {
	dayOfWeek := t.Weekday()
	return t.Add(time.Duration(time.Saturday-dayOfWeek+1) * 24 * time.Hour).Add(-1 * time.Nanosecond)
}

// StartOfQuarter returns a Time for start of the quarter.
func (t Time[T]) StartOfQuarter() Time[T] {
	year, quarter, day := t.Year(), numberOfQuarter(t.Month()), 1
	return New[T](year, time.Month(3*quarter-2), day, 0, 0, 0, 0)
}

// EndOfQuarter returns a Time for end of the quarter.
func (t Time[T]) EndOfQuarter() Time[T] {
	year, quarter := t.Year(), numberOfQuarter(t.Month())
	day := 31
	switch quarter {
	case 2, 3:
		day = 30
	}
	return New[T](year, time.Month(3*quarter), day, 23, 59, 59, 999999999)
}

// StartOfSemester returns a Time for start of the semester.
func (t Time[T]) StartOfSemester() Time[T] {
	year, semester, day := t.Year(), numberOfSemester(t.Month()), 1
	month := time.January
	if semester == 2 {
		month = time.July
	}
	return New[T](year, month, day, 0, 0, 0, 0)
}

// EndOfSemester returns a Time for end of the semester.
func (t Time[T]) EndOfSemester() Time[T] {
	year, semester := t.Year(), numberOfSemester(t.Month())
	month, day := time.June, 30
	if semester == 2 {
		month, day = time.December, 31
	}
	return New[T](year, month, day, 23, 59, 59, 999999999)
}

// IsLeapYear returns true if t is leap year.
func (t Time[T]) IsLeapYear() bool {
	year := t.Year()
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// IsBetween returns true if from < t && t < to.
func (t Time[T]) IsBetween(from Time[T], to Time[T]) bool {
	return from.Before(t) && to.After(t)
}

// DiffInCalendarDays calculates the difference in calendar days between t and u. (t-u)
// Calendar days are calculated by considering only the dates, excluding the times,
// and then determining the difference in days.
func (t Time[T]) DiffInCalendarDays(u Time[T]) int {
	const day = 24 * time.Hour
	t1 := t.Truncate(day)
	u1 := u.Truncate(day)
	return int(math.Ceil(float64(t1.Sub(u1)) / float64(day)))
}
