package synchro

import (
	"context"
	"time"

	"github.com/Code-Hex/synchro/iso8601"
)

var nowFunc = time.Now

// In returns timezone-aware time.
//
// If the given time.Time is the zero value, a zero value Time[T] is returned.
// A zero value Time[T] is not timezone-aware. Always set UTC.
func In[T TimeZone](tm time.Time) Time[T] {
	if tm.IsZero() {
		return Time[T]{}
	}
	var tz T
	return Time[T]{tm: tm.In(tz.Location())}
}

// Now returns the current time with timezone.
func Now[T TimeZone]() Time[T] {
	return In[T](nowFunc())
}

type nowContextKey[T TimeZone] struct{}

// NowContext returns the current time stored in the provided context.
// If the time is not found in the context, it returns zero value.
//
// NowContext and NowWithContext are useful when you want to store
// the current time within the context of the executing logic.
// For example, in scenarios where you consider the entire lifecycle
// of an HTTP request, from its initiation to the response.
//
// By capturing the timestamp when the request occurs and fixing it
// as the current time, you can achieve consistent handling of
// the current time throughout that request. This ensures uniformity
// in dealing with the current time within the scope of that specific request.
func NowContext[T TimeZone](ctx context.Context) Time[T] {
	t, ok := ctx.Value(nowContextKey[T]{}).(Time[T])
	if !ok {
		return Time[T]{}
	}
	return t
}

// NowWithContext returns a new context with the provided time with timezone stored in it.
// This allows you to create a context with a specific time value attached to it.
func NowWithContext[T TimeZone](ctx context.Context, t Time[T]) context.Context {
	return context.WithValue(ctx, nowContextKey[T]{}, t)
}

// New returns the Time corresponding to
//
//	yyyy-mm-dd hh:mm:ss + nsec nanoseconds
//
// in the appropriate zone for that time in the given timezone.
//
// The month, day, hour, min, sec, and nsec values may be outside
// their usual ranges and will be normalized during the conversion.
// For example, October 32 converts to November 1.
//
// A daylight savings time transition skips or repeats times.
// For example, in the United States, March 13, 2011 2:15am never occurred,
// while November 6, 2011 1:15am occurred twice. In such cases, the
// choice of time zone, and therefore the time, is not well-defined.
// Date returns a time that is correct in one of the two zones involved
// in the transition, but it does not guarantee which.
//
// This is a simple wrapper function for time.Date.
func New[T TimeZone](year int, month time.Month, day int, hour int, min int, sec int, nsec int) Time[T] {
	var tz T
	tm := time.Date(year, month, day, hour, min, sec, nsec, tz.Location())
	return Time[T]{tm: tm}
}

// Parse parses a formatted string and returns the time value it represents.
// See the documentation for the constant called Layout to see how to
// represent the format. The second argument must be parseable using
// the format string (layout) provided as the first argument.
//
// This is a simple wrapper function for time.ParseInLocation.
func Parse[T TimeZone](layout, value string) (Time[T], error) {
	var tz T
	tm, err := time.ParseInLocation(layout, value, tz.Location())
	if err != nil {
		return Time[T]{}, err
	}
	return In[T](tm), nil
}

// ParseISO parses an ISO8601-compliant date or datetime string and returns
// its representation as a Time. If the input string does not conform to the
// ISO8601 standard or if any other parsing error occurs, an error is returned.
// Supported formats include:
//
//	Basic                        Extended
//	20070301                     2007-03-01
//	2012W521                     2012-W52-1
//	2012Q485                     2012-Q4-85
//	20070301T1300Z               2007-03-01T13:00Z
//	20070301T1300Z               2007-03-01T13:00Z
//	20070301T1300+0100           2007-03-01T13:00+01:00
//	20070301T1300-0600           2007-03-01T13:00-06:00
//	20070301T130045Z             2007-03-01T13:00:45Z
//	20070301T130045+0100         2007-03-01T13:00:45+01:00
//	... and other combinations
func ParseISO[T TimeZone](value string) (Time[T], error) {
	var tz T
	tm, err := iso8601.ParseDateTime[string](value, iso8601.WithInLocation(tz.Location()))
	if err != nil {
		return Time[T]{}, err
	}
	return In[T](tm), nil
}

// Unix returns the local Time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC.
// It is valid to pass nsec outside the range [0, 999999999].
// Not all sec values have a corresponding time value. One such
// value is 1<<63-1 (the largest int64 value).
//
// This is a simple wrapper function for time.Unix.
func Unix[T TimeZone](sec int64, nsec int64) Time[T] {
	return In[T](time.Unix(sec, nsec))
}

// UnixMilli returns the local Time corresponding to the given Unix time,
// msec milliseconds since January 1, 1970 UTC.
//
// This is a simple wrapper function for time.UnixMilli.
func UnixMilli[T TimeZone](msec int64) Time[T] {
	return In[T](time.UnixMilli(msec))
}

// UnixMicro returns the local Time corresponding to the given Unix time,
// usec microseconds since January 1, 1970 UTC.
//
// This is a simple wrapper function for time.UnixMicro.
func UnixMicro[T TimeZone](usec int64) Time[T] {
	return In[T](time.UnixMicro(usec))
}

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
// It is equivalent to NewTimer(d).C.
// The underlying Timer is not recovered by the garbage collector
// until the timer fires. If efficiency is a concern, use NewTimer
// instead and call Timer.Stop if the timer is no longer needed.
//
// This is a simple wrapper function for time.After.
func After[T TimeZone](d time.Duration) <-chan Time[T] {
	c := make(chan Time[T], 1)
	go func() { c <- In[T](<-time.After(d)) }()
	return c
}

// ConvertTz can be used to convert a time from one time zone to another.
//
// For example to convert from UTC to Asia/Tokyo.
// If `2009-11-10 23:00:00 +0000 UTC` as input, the output will be `2009-11-11 08:00:00 +0900 Asia/Tokyo`.
func ConvertTz[U TimeZone, T TimeZone](from Time[T]) Time[U] {
	return In[U](from.tm)
}
