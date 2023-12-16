package synchro

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/Code-Hex/synchro/internal/constraints"
	"github.com/Code-Hex/synchro/iso8601"
)

type timeish[T TimeZone] interface {
	Time[T] | time.Time | constraints.Bytes
}

// Period allows iteration over a set of dates and times,
// recurring at regular intervals, over a given period.
type Period[T TimeZone] struct {
	from Time[T]
	to   Time[T]
}

// String implements the fmt.Stringer interface.
func (p Period[T]) String() string {
	return fmt.Sprintf("from %s to %s", p.from, p.to)
}

// From returns end of period.
func (p Period[T]) From() Time[T] { return p.from }

// To returns start of period.
func (p Period[T]) To() Time[T] { return p.to }

// NewPeriod creates a new Period struct between the 'from' and 'to' values you specified.
//
// If Time[T] or time.Time is specified, it guarantees no error returns.
// When a string or []byte is passed, ParseISO function is called internally. Therefore, these
// parameters should be in a format compatible with ParseISO.
func NewPeriod[T TimeZone, T1 timeish[T], T2 timeish[T]](from T1, to T2) (Period[T], error) {
	start, err := convertTime[T, T1](unsafe.Pointer(&from))
	if err != nil {
		return Period[T]{}, fmt.Errorf("failed to parse from: %w", err)
	}
	end, err := convertTime[T, T2](unsafe.Pointer(&to))
	if err != nil {
		return Period[T]{}, fmt.Errorf("failed to parse to: %w", err)
	}
	return Period[T]{
		from: start,
		to:   end,
	}, nil
}

// Contains checks whether the specified t is included within from and to.
//
// if p.from < t && t < p.to, it returns +1; if p.from == t || t == p.to, it returns 0.
// Otherwise returns -1;
func (p Period[T]) Contains(t Time[T]) int {
	cmpFrom := p.from.Compare(t)
	cmpTo := p.to.Compare(t)
	if cmpFrom == -1 && cmpTo == 1 {
		return 1
	}
	if cmpFrom == 0 || cmpTo == 0 {
		return 0
	}
	return -1
}

type periodical[T TimeZone] <-chan Time[T]

// Slice returns the slice of Time[T].
func (p periodical[T]) Slice() (s []Time[T]) {
	for current := range p {
		s = append(s, current)
	}
	return s
}

// Periodic returns a channel that emits Time[T] values at regular intervals
// between the start and end times of the Period[T]. The interval is specified
// by the next function argument.
//
// If start < end, the process will increase from start to end. In other words,
// when the current value exceeds end, the iteration is terminated.
//
// If start > end, the process will decrease from start to end. In other words,
// when the current value falls below end, the iteration is terminated.
func (p Period[T]) Periodic(next func(Time[T]) Time[T]) periodical[T] {
	compare := isNotAfter[T]
	// p.start > p.end
	if p.from.After(p.to) {
		compare = isNotBefore[T]
	}
	ch := make(chan Time[T], 1)
	go func() {
		defer close(ch)
		for current := p.from; compare(current, p.to); current = next(current) {
			ch <- current
		}
	}()
	return ch
}

func isNotAfter[T TimeZone](t1, t2 Time[T]) bool {
	// t1.Compare(t2) <= 0
	return !t1.After(t2)
}

func isNotBefore[T TimeZone](t1, t2 Time[T]) bool {
	// t1.Compare(t2) >= 0
	return !t1.Before(t2)
}

// PeriodicDuration is a wrapper for the Periodic function.
// The interval is specified by the time.Duration argument.
func (p Period[T]) PeriodicDuration(d time.Duration) periodical[T] {
	return p.Periodic(func(t Time[T]) Time[T] {
		return t.Add(d)
	})
}

// PeriodicDuration is a wrapper for the Periodic function.
// The interval is specified by the given number of years, months, and days.
func (p Period[T]) PeriodicDate(years int, months int, days int) periodical[T] {
	return p.Periodic(func(t Time[T]) Time[T] {
		return t.AddDate(years, months, days)
	})
}

// PeriodicAdvance is a wrapper for the Periodic function.
// The interval is specified by the provided date and time unit arguments.
func (p Period[T]) PeriodicAdvance(u1 Unit, u2 ...Unit) periodical[T] {
	return p.Periodic(func(t Time[T]) Time[T] {
		return t.Advance(u1, u2...)
	})
}

// PeriodicISODuration is a wrapper for the Periodic function. It accepts a duration
// in ISO 8601 format as a parameter.
//
// Examples of valid durations include:
//
//	PnYnMnDTnHnMnS (e.g., P3Y6M4DT12H30M5S)
//	PnW (e.g., P4W)
func (p Period[T]) PeriodicISODuration(duration string) (periodical[T], error) {
	d, err := iso8601.ParseDuration(duration)
	if err != nil {
		return nil, err
	}
	if d.IsZero() {
		return nil, fmt.Errorf("empty duration is not accepted: %q", duration)
	}
	sign := 1
	if d.Negative {
		sign = -1
	}
	return p.Periodic(func(t Time[T]) Time[T] {
		var (
			years  int
			months int
			days   int
		)

		if d.Year > 0 {
			years = sign * d.Year
		}
		if d.Month > 0 {
			months = sign * int(d.Month)
		}
		if d.Week > 0 {
			days += sign * 7 * d.Week
		}
		if d.Day > 0 {
			days += sign * d.Day
		}
		if years != 0 || months != 0 || days != 0 {
			t = t.AddDate(years, months, days)
		}
		return t.Add(d.StdClockDuration())
	}), nil
}

func convertTime[T TimeZone, argType timeish[T]](argPtr unsafe.Pointer) (Time[T], error) {
	var dummy argType
	switch any(dummy).(type) {
	case Time[T]:
		return *(*Time[T])(argPtr), nil
	case time.Time:
		return In[T](*(*time.Time)(argPtr)), nil
	case []byte:
		bytes := *(*[]byte)(argPtr)
		str := unsafe.String(unsafe.SliceData(bytes), len(bytes))
		return ParseISO[T](str)
	default:
		// argType is ~string, argPtr can be safely converted to *string
		return ParseISO[T](*(*string)(argPtr))
	}
}
