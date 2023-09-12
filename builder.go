package synchro

import (
	"time"

	"github.com/Code-Hex/synchro/tz"
)

func toPtr[T any](x T) *T {
	return &x
}

// TimeBuilder defines an interface for constructing a Time[T] value
// with specific date and time components using method chaining.
//
// Implementations should allow for setting each component step by step,
// and then finalizing the construction with the Do() method.
type TimeBuilder[T TimeZone] interface {
	// Year sets the year component of the time being built.
	Year(int) TimeBuilder[T]

	// Month sets the month component of the time being built.
	Month(time.Month) TimeBuilder[T]

	// Day sets the day component of the time being built.
	Day(int) TimeBuilder[T]

	// Hour sets the hour component of the time being built.
	Hour(int) TimeBuilder[T]

	// Minute sets the minute component of the time being built.
	Minute(int) TimeBuilder[T]

	// Second sets the second component of the time being built.
	Second(int) TimeBuilder[T]

	// Nanosecond sets the nanosecond component of the time being built.
	Nanosecond(int) TimeBuilder[T]

	// Do finalizes and returns the constructed Time value
	// based on the components set in the chain.
	Do() Time[T]
}

type changeBuilder[T TimeZone] struct {
	year   *int
	month  *time.Month
	day    *int
	hour   *int
	minute *int
	second *int
	nsec   *int

	t Time[T]
}

var _ TimeBuilder[tz.UTC] = changeBuilder[tz.UTC]{}

// Year implements the TimeBuilder interface to change years.
func (c changeBuilder[T]) Year(year int) TimeBuilder[T] {
	c.year = toPtr(year)
	return c
}

// Month implements the TimeBuilder interface to change months.
func (c changeBuilder[T]) Month(month time.Month) TimeBuilder[T] {
	c.month = toPtr(month)
	return c
}

// Day implements the TimeBuilder interface to change days.
func (c changeBuilder[T]) Day(day int) TimeBuilder[T] {
	c.day = toPtr(day)
	return c
}

// Hour implements the TimeBuilder interface to change hours.
func (c changeBuilder[T]) Hour(hour int) TimeBuilder[T] {
	c.hour = toPtr(hour)
	return c
}

// Minute implements the TimeBuilder interface to change minutes.
func (c changeBuilder[T]) Minute(minute int) TimeBuilder[T] {
	c.minute = toPtr(minute)
	return c
}

// Second implements the TimeBuilder interface to change seconds.
func (c changeBuilder[T]) Second(second int) TimeBuilder[T] {
	c.second = toPtr(second)
	return c
}

// Nanosecond implements the TimeBuilder interface to change nanoseconds.
func (c changeBuilder[T]) Nanosecond(nsec int) TimeBuilder[T] {
	c.nsec = toPtr(nsec)
	return c
}

// Do implements the TimeBuilder interface to change any date and time components.
func (c changeBuilder[T]) Do() Time[T] {
	year, month, day := c.t.Date()
	hour, min, sec := c.t.Clock()
	nsec := c.t.Nanosecond()
	if c.year != nil {
		year = *c.year
	}
	if c.month != nil {
		month = *c.month
	}
	if c.day != nil {
		day = *c.day
	}
	if c.hour != nil {
		hour = *c.hour
	}
	if c.minute != nil {
		min = *c.minute
	}
	if c.second != nil {
		sec = *c.second
	}
	if c.nsec != nil {
		nsec = *c.nsec
	}
	return New[T](year, month, day, hour, min, sec, nsec)
}

// Change initializes and returns a TimeBuilder for the current Time value.
//
// This provides a method chain approach for specifying which parts of the time
// you want to change, allowing for the creation of a new Time[T] instance
// with the specified modifications.
func (t Time[T]) Change() TimeBuilder[T] {
	return changeBuilder[T]{t: t}
}

type advanceBuilder[T TimeZone] struct {
	year   int
	month  time.Month
	day    int
	hour   int
	minute int
	second int
	nsec   int

	t Time[T]
}

var _ TimeBuilder[tz.UTC] = advanceBuilder[tz.UTC]{}

// Advance initializes and returns a TimeBuilder for the current Time value.
//
// This provides a method chain approach for specifying which parts of the time
// you want to increment, allowing for the creation of a new Time[T] instance
// with the specified modifications.
func (t Time[T]) Advance() TimeBuilder[T] {
	return advanceBuilder[T]{t: t}
}

// Year implements the TimeBuilder interface to increment years.
func (a advanceBuilder[T]) Year(year int) TimeBuilder[T] {
	a.year += year
	return a
}

// Month implements the TimeBuilder interface to increment months.
func (a advanceBuilder[T]) Month(month time.Month) TimeBuilder[T] {
	a.month += month
	return a
}

// Day implements the TimeBuilder interface to increment days.
func (a advanceBuilder[T]) Day(day int) TimeBuilder[T] {
	a.day += day
	return a
}

// Hour implements the TimeBuilder interface to increment hours.
func (a advanceBuilder[T]) Hour(hour int) TimeBuilder[T] {
	a.hour += hour
	return a
}

// Minute implements the TimeBuilder interface to increment minutes.
func (a advanceBuilder[T]) Minute(minute int) TimeBuilder[T] {
	a.minute += minute
	return a
}

// Second implements the TimeBuilder interface to increment seconds.
func (a advanceBuilder[T]) Second(second int) TimeBuilder[T] {
	a.second += second
	return a
}

// Nanosecond implements the TimeBuilder interface to increment nanoseconds.
func (a advanceBuilder[T]) Nanosecond(nsec int) TimeBuilder[T] {
	a.nsec += nsec
	return a
}

// Do implements the TimeBuilder interface to increment any date and time components.
func (a advanceBuilder[T]) Do() Time[T] {
	t := a.t.AddDate(a.year, int(a.month), a.day)

	if a.hour != 0 {
		t = t.Add(time.Hour * time.Duration(a.hour))
	}
	if a.minute != 0 {
		t = t.Add(time.Minute * time.Duration(a.minute))
	}
	if a.second != 0 {
		t = t.Add(time.Second * time.Duration(a.second))
	}
	if a.nsec != 0 {
		t = t.Add(time.Duration(a.nsec))
	}
	return t
}
