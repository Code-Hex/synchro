package synchro

import (
	"time"
)

type unit interface {
	cast() int
}

type (
	// Year sets the year component of the time being built.
	Year int

	// Month sets the month component of the time being built.
	Month int

	// Day sets the day component of the time being built.
	Day int

	// Hour sets the hour component of the time being built.
	Hour int

	// Minute sets the minute component of the time being built.
	Minute int

	// Second sets the second component of the time being built.
	Second int

	// Nanosecond sets the nanosecond component of the time being built.
	Nanosecond int
)

func (y Year) cast() int        { return int(y) }
func (m Month) cast() int       { return int(m) }
func (d Day) cast() int         { return int(d) }
func (h Hour) cast() int        { return int(h) }
func (m Minute) cast() int      { return int(m) }
func (s Second) cast() int      { return int(s) }
func (ns Nanosecond) cast() int { return int(ns) }

// Change modifies the time based on the provided unit values.
// u1 is a required unit, while u2... can be provided as additional optional units.
// This method returns a new Time[T] and does not modify the original.
func (t Time[T]) Change(u1 unit, u2 ...unit) Time[T] {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	nsec := t.Nanosecond()
	for _, u := range append([]unit{u1}, u2...) {
		switch v := u.(type) {
		case Year:
			year = v.cast()
		case Month:
			month = time.Month(v.cast())
		case Day:
			day = v.cast()
		case Hour:
			hour = v.cast()
		case Minute:
			min = v.cast()
		case Second:
			sec = v.cast()
		case Nanosecond:
			nsec = v.cast()
		}
	}
	return New[T](year, month, day, hour, min, sec, nsec)
}

// Advance adjusts the time based on the provided unit values, moving it forward in time.
// u1 is a required unit, while u2... can be provided as additional optional units.
// This method returns a new Time[T] and does not modify the original.
// The time is adjusted in the order the units are provided.
func (t Time[T]) Advance(u1 unit, u2 ...unit) Time[T] {
	ret := t
	years, months, days := 0, time.Month(0), 0
	for _, u := range append([]unit{u1}, u2...) {
		switch v := u.(type) {
		case Year:
			years += v.cast()
		case Month:
			months += time.Month(v.cast())
		case Day:
			days += v.cast()
		case Hour:
			ret = ret.Add(time.Hour * time.Duration(v.cast()))
		case Minute:
			ret = ret.Add(time.Minute * time.Duration(v.cast()))
		case Second:
			ret = ret.Add(time.Second * time.Duration(v.cast()))
		case Nanosecond:
			ret = ret.Add(time.Duration(v.cast()))
		}
	}
	year, month, day := ret.Date()
	hour, min, sec := ret.Clock()
	return New[T](year+years, month+months, day+days, hour, min, sec, ret.Nanosecond())
}
