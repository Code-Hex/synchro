package synchro

import "time"

type Quarter[T TimeZone] struct {
	year   int
	number int
	t      Time[T]
}

// Quarter gets current quarter.
func (t Time[T]) Quarter() Quarter[T] {
	return Quarter[T]{
		year:   t.Year(),
		number: numberOfQuarter(t.Month()),
		t:      t,
	}
}

// Year returns the year in which q occurs.
func (q Quarter[T]) Year() int { return q.year }

// Number returns the number of quarter.
func (q Quarter[T]) Number() int { return q.number }

// Start returns start time in the quarter.
func (q Quarter[T]) Start() Time[T] { return q.t.StartOfQuarter() }

// End returns end time in the quarter.
func (q Quarter[T]) End() Time[T] { return q.t.EndOfQuarter() }

// After reports whether the Quarter instant q is after u.
func (q Quarter[T]) After(u Quarter[T]) bool {
	if q.year > u.year {
		return true
	}
	if q.number > u.number {
		return true
	}
	return false
}

// Before reports whether the Quarter instant q is before u.
func (q Quarter[T]) Before(u Quarter[T]) bool {
	if q.year < u.year {
		return true
	}
	if q.number < u.number {
		return true
	}
	return false
}

// Compare compares the Quarter instant q with u. If q is before u, it returns -1;
// if q is after u, it returns +1; if they're the same, it returns 0.
func (q Quarter[T]) Compare(u Quarter[T]) int {
	if q.year == u.year && q.number == u.number {
		return 0
	}
	if q.Before(u) {
		return -1
	}
	return 1
}

func numberOfQuarter(month time.Month) int {
	if month >= time.October {
		return 4
	}
	if month >= time.July {
		return 3
	}
	if month >= time.April {
		return 2
	}
	return 1
}
