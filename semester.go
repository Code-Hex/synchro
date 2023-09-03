package synchro

import "time"

type Semester[T TimeZone] struct {
	year   int
	number int
	t      Time[T]
}

// Quarter gets current semester.
func (t Time[T]) Semester() Semester[T] {
	return Semester[T]{
		year:   t.Year(),
		number: numberOfSemester(t.Month()),
		t:      t,
	}
}

// Year returns the year in which s occurs.
func (q Semester[T]) Year() int { return q.year }

// Number returns the number of semester.
func (q Semester[T]) Number() int { return q.number }

// Start returns start time in the semester.
func (q Semester[T]) Start() Time[T] { return q.t.StartOfSemester() }

// End returns end time in the semester.
func (q Semester[T]) End() Time[T] { return q.t.EndOfSemester() }

// After reports whether the Semester instant s is after u.
func (s Semester[T]) After(u Semester[T]) bool {
	if s.year > u.year {
		return true
	}
	if s.number > u.number {
		return true
	}
	return false
}

// Before reports whether the Semester instant s is before u.
func (s Semester[T]) Before(u Semester[T]) bool {
	if s.year < u.year {
		return true
	}
	if s.number < u.number {
		return true
	}
	return false
}

// Compare compares the Semester instant s with u. If s is before u, it returns -1;
// if s is after u, it returns +1; if they're the same, it returns 0.
func (s Semester[T]) Compare(u Semester[T]) int {
	if s.year == u.year && s.number == u.number {
		return 0
	}
	if s.Before(u) {
		return -1
	}
	return 1
}

func numberOfSemester(month time.Month) int {
	if month >= time.July {
		return 2
	}
	return 1
}
