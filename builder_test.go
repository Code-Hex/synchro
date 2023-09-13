package synchro_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

// shorthand way
type Y = synchro.Year
type M = synchro.Month
type D = synchro.Day
type HH = synchro.Hour
type MM = synchro.Minute
type SS = synchro.Second
type NS = synchro.Nanosecond

func ExampleTime_Change() {
	utc := synchro.New[tz.UTC](2009, time.November, 10, 23, 0, 0, 0)
	c1 := utc.Change(synchro.Year(2010))
	c2 := utc.Change(synchro.Year(2010), synchro.Month(time.December))
	c3 := utc.Change(Y(2010), M(time.December), D(1))
	c4 := c3.Change(synchro.Hour(1))
	c5 := c3.Change(HH(1), MM(1))
	c6 := c3.Change(HH(1), MM(1), SS(1))
	c7 := c3.Change(HH(1), MM(1), SS(1), NS(123456789))
	fmt.Printf("Go launched at %s\n", utc)
	fmt.Println(c1)
	fmt.Println(c2)
	fmt.Println(c3)
	fmt.Println(c4)
	fmt.Println(c5)
	fmt.Println(c6)
	fmt.Println(c7)
	// Output:
	// Go launched at 2009-11-10 23:00:00 +0000 UTC
	// 2010-11-10 23:00:00 +0000 UTC
	// 2010-12-10 23:00:00 +0000 UTC
	// 2010-12-01 23:00:00 +0000 UTC
	// 2010-12-01 01:00:00 +0000 UTC
	// 2010-12-01 01:01:00 +0000 UTC
	// 2010-12-01 01:01:01 +0000 UTC
	// 2010-12-01 01:01:01.123456789 +0000 UTC
}

func ExampleTime_Advance() {
	utc := synchro.New[tz.UTC](2009, time.November, 10, 23, 0, 0, 0)
	c1 := utc.Advance(synchro.Year(1))
	c11 := utc.Advance(Y(1), Y(1)) // +2 years

	c2 := utc.Advance(Y(1), M(1))
	c3 := utc.Advance(Y(1), M(1), D(1))
	c4 := c3.Advance(HH(1))
	c5 := c3.Advance(HH(1), MM(1))
	c6 := c3.Advance(HH(1), MM(1), SS(1))
	c7 := c3.Advance(HH(1), MM(1), SS(1), NS(123456789))

	fmt.Printf("Go launched at %s\n", utc)
	fmt.Println(c1)
	fmt.Println(c11)
	fmt.Println()
	fmt.Println(c2)
	fmt.Println(c3)
	fmt.Println(c4)
	fmt.Println(c5)
	fmt.Println(c6)
	fmt.Println(c7)
	// Output:
	// Go launched at 2009-11-10 23:00:00 +0000 UTC
	// 2010-11-10 23:00:00 +0000 UTC
	// 2011-11-10 23:00:00 +0000 UTC
	//
	// 2010-12-10 23:00:00 +0000 UTC
	// 2010-12-11 23:00:00 +0000 UTC
	// 2010-12-12 00:00:00 +0000 UTC
	// 2010-12-12 00:01:00 +0000 UTC
	// 2010-12-12 00:01:01 +0000 UTC
	// 2010-12-12 00:01:01.123456789 +0000 UTC
}

func TestAdvance(t *testing.T) {
	utc := synchro.New[tz.UTC](2009, time.November, 10, 23, 0, 0, 0)

	t.Run("month twice", func(t *testing.T) {
		got := utc.Advance(M(1), M(2)) // +3 months
		want := synchro.New[tz.UTC](2010, time.February, 10, 23, 0, 0, 0)
		if want != got {
			t.Fatalf("- %s\n+ %s", want, got)
		}
	})
	t.Run("day twice", func(t *testing.T) {
		got := utc.Advance(D(1), D(2)) // +3 days
		want := synchro.New[tz.UTC](2009, time.November, 13, 23, 0, 0, 0)
		if want != got {
			t.Fatalf("- %s\n+ %s", want, got)
		}
	})
	t.Run("hour twice", func(t *testing.T) {
		got := utc.Advance(HH(1), HH(2)) // +3 hours
		want := synchro.New[tz.UTC](2009, time.November, 11, 2, 0, 0, 0)
		if want != got {
			t.Fatalf("- %s\n+ %s", want, got)
		}
	})
	t.Run("minute twice", func(t *testing.T) {
		got := utc.Advance(MM(5), MM(60)) // +65 minutes
		want := synchro.New[tz.UTC](2009, time.November, 11, 0, 5, 0, 0)
		if want != got {
			t.Fatalf("- %s\n+ %s", want, got)
		}
	})
	t.Run("second twice", func(t *testing.T) {
		got := utc.Advance(SS(5), SS(60)) // +65 seconds
		want := synchro.New[tz.UTC](2009, time.November, 10, 23, 1, 5, 0)
		if want != got {
			t.Fatalf("- %s\n+ %s", want, got)
		}
	})
	t.Run("nanosec twice", func(t *testing.T) {
		got := utc.Advance(NS(5), NS(60)) // +65 nanosec
		want := synchro.New[tz.UTC](2009, time.November, 10, 23, 0, 0, 65)
		if want != got {
			t.Fatalf("- %s\n+ %s", want, got)
		}
	})
}
