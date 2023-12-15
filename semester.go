package synchro

import "time"

type Semester[T TimeZone] struct {
	year   int
	number int
}

// Semester gets current semester.
func (t Time[T]) Semester() Semester[T] {
	return Semester[T]{
		year:   t.Year(),
		number: numberOfSemester(t.Month()),
	}
}

// Year returns the year in which s occurs.
func (s Semester[T]) Year() int { return s.year }

// Number returns the number of semester.
func (s Semester[T]) Number() int { return s.number }

// Start returns start time in the semester.
func (s Semester[T]) Start() Time[T] {
	year, semester, day := s.year, s.number, 1
	month := time.January
	if semester == 2 {
		month = time.July
	}
	return New[T](year, month, day, 0, 0, 0, 0)
}

// End returns end time in the semester.
func (s Semester[T]) End() Time[T] {
	year, semester := s.year, s.number
	month, day := time.June, 30
	if semester == 2 {
		month, day = time.December, 31
	}
	return New[T](year, month, day, 23, 59, 59, 999999999)
}

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
