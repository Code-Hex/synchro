package synchro_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

func init() {
	synchro.SetNow(func() time.Time {
		return time.Date(2023, 9, 2, 14, 0, 0, 0, time.UTC)
	})
}

func ExampleIn() {
	d := time.Date(2023, 9, 2, 14, 0, 0, 0, time.UTC)
	utc := synchro.In[tz.UTC](d)
	fmt.Println(utc)

	jst := synchro.In[tz.AsiaTokyo](d)
	fmt.Println(jst)
	// Output:
	// 2023-09-02 14:00:00 +0000 UTC
	// 2023-09-02 23:00:00 +0900 JST
}

func ExampleNow() {
	// The current UTC time is fixed to `2023-09-02 14:00:00`.
	utcNow := synchro.Now[tz.UTC]()
	fmt.Println(utcNow)

	jstNow := synchro.Now[tz.AsiaTokyo]()
	fmt.Println(jstNow)
	// Output:
	// 2023-09-02 14:00:00 +0000 UTC
	// 2023-09-02 23:00:00 +0900 JST
}

func ExampleNowContext() {
	// This is the timestamp when the request occurred.
	// The current UTC time is fixed to `2023-09-02 14:00:00`.
	timestamp := synchro.Now[tz.UTC]()

	// The context of the request.
	ctx := context.Background()

	// Set the current time within the request context.Context.
	ctx = synchro.NowWithContext[tz.UTC](ctx, timestamp)
	utcNow := synchro.NowContext[tz.UTC](ctx)
	fmt.Println(utcNow)

	// A zero value is returned because the time is not stored in the same timezone.
	jstNow := synchro.NowContext[tz.AsiaTokyo](ctx)
	fmt.Println(jstNow)
	// Output:
	// 2023-09-02 14:00:00 +0000 UTC
	// 0001-01-01 00:00:00 +0000 UTC
}

func ExampleNew() {
	utc := synchro.New[tz.UTC](2009, time.November, 10, 23, 0, 0, 0)
	fmt.Printf("Go launched at %s\n", utc)
	// Output:
	// Go launched at 2009-11-10 23:00:00 +0000 UTC
}

type EuropeBerlin struct{}

func (EuropeBerlin) Location() *time.Location {
	loc, _ := time.LoadLocation("Europe/Berlin")
	return loc
}

func ExampleParse() {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	t, _ := synchro.Parse[EuropeBerlin](longForm, "Jul 9, 2012 at 5:02am (CEST)")
	fmt.Println(t)

	// Note: without explicit zone, returns time in given location.
	const shortForm = "2006-Jan-02"
	t, _ = synchro.Parse[EuropeBerlin](shortForm, "2012-Jul-09")
	fmt.Println(t)

	_, err := synchro.Parse[tz.UTC](time.RFC3339, time.RFC3339)
	fmt.Println("error", err) // Returns an error as the layout is not a valid time value

	// Output:
	// 2012-07-09 05:02:00 +0200 CEST
	// 2012-07-09 00:00:00 +0200 CEST
	// error parsing time "2006-01-02T15:04:05Z07:00": extra text: "07:00"
}

func ExampleUnix() {
	unixTime := synchro.New[tz.UTC](2009, time.November, 10, 23, 0, 0, 0)
	fmt.Println(unixTime.Unix())
	t := synchro.Unix[tz.UTC](unixTime.Unix(), 0)
	fmt.Println(t)

	// Output:
	// 1257894000
	// 2009-11-10 23:00:00 +0000 UTC
}

func ExampleUnixMicro() {
	umt := synchro.New[tz.UTC](2009, time.November, 10, 23, 0, 0, 0)
	fmt.Println(umt.UnixMicro())
	t := synchro.UnixMicro[tz.UTC](umt.UnixMicro())
	fmt.Println(t)

	// Output:
	// 1257894000000000
	// 2009-11-10 23:00:00 +0000 UTC
}

func ExampleUnixMilli() {
	umt := synchro.New[tz.UTC](2009, time.November, 10, 23, 0, 0, 0)
	fmt.Println(umt.UnixMilli())
	t := synchro.UnixMilli[tz.UTC](umt.UnixMilli())
	fmt.Println(t)

	// Output:
	// 1257894000000
	// 2009-11-10 23:00:00 +0000 UTC
}

func ExampleTime_Unix() {
	// 1 billion seconds of Unix, three ways.
	fmt.Println(synchro.Unix[tz.UTC](1e9, 0))     // 1e9 seconds
	fmt.Println(synchro.Unix[tz.UTC](0, 1e18))    // 1e18 nanoseconds
	fmt.Println(synchro.Unix[tz.UTC](2e9, -1e18)) // 2e9 seconds - 1e18 nanoseconds

	t := synchro.New[tz.UTC](2001, time.September, 9, 1, 46, 40, 0)
	fmt.Println(t.Unix())     // seconds since 1970
	fmt.Println(t.UnixNano()) // nanoseconds since 1970

	// Output:
	// 2001-09-09 01:46:40 +0000 UTC
	// 2001-09-09 01:46:40 +0000 UTC
	// 2001-09-09 01:46:40 +0000 UTC
	// 1000000000
	// 1000000000000000000
}

func ExampleTime_Round() {
	t := synchro.New[tz.UTC](0, 0, 0, 12, 15, 30, 918273645)
	round := []time.Duration{
		time.Nanosecond,
		time.Microsecond,
		time.Millisecond,
		time.Second,
		2 * time.Second,
		time.Minute,
		10 * time.Minute,
		time.Hour,
	}

	for _, d := range round {
		fmt.Printf("t.Round(%6s) = %s\n", d, t.Round(d).Format("15:04:05.999999999"))
	}
	// Output:
	// t.Round(   1ns) = 12:15:30.918273645
	// t.Round(   1µs) = 12:15:30.918274
	// t.Round(   1ms) = 12:15:30.918
	// t.Round(    1s) = 12:15:31
	// t.Round(    2s) = 12:15:30
	// t.Round(  1m0s) = 12:16:00
	// t.Round( 10m0s) = 12:20:00
	// t.Round(1h0m0s) = 12:00:00
}

func ExampleTime_Truncate() {
	t, _ := synchro.Parse[tz.UTC]("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")
	trunc := []time.Duration{
		time.Nanosecond,
		time.Microsecond,
		time.Millisecond,
		time.Second,
		2 * time.Second,
		time.Minute,
		10 * time.Minute,
	}

	for _, d := range trunc {
		fmt.Printf("t.Truncate(%5s) = %s\n", d, t.Truncate(d).Format("15:04:05.999999999"))
	}
	// To round to the last midnight in the local timezone, create a new Date.
	midnight := synchro.New[tz.UTC](t.Year(), t.Month(), t.Day(), 0, 0, 0, 0)
	_ = midnight

	// Output:
	// t.Truncate(  1ns) = 12:15:30.918273645
	// t.Truncate(  1µs) = 12:15:30.918273
	// t.Truncate(  1ms) = 12:15:30.918
	// t.Truncate(   1s) = 12:15:30
	// t.Truncate(   2s) = 12:15:30
	// t.Truncate( 1m0s) = 12:15:00
	// t.Truncate(10m0s) = 12:10:00
}

func ExampleTime_Add() {
	start := synchro.New[tz.UTC](2009, 1, 1, 12, 0, 0, 0)
	afterTenSeconds := start.Add(time.Second * 10)
	afterTenMinutes := start.Add(time.Minute * 10)
	afterTenHours := start.Add(time.Hour * 10)
	afterTenDays := start.Add(time.Hour * 24 * 10)

	fmt.Printf("start = %v\n", start)
	fmt.Printf("start.Add(time.Second * 10) = %v\n", afterTenSeconds)
	fmt.Printf("start.Add(time.Minute * 10) = %v\n", afterTenMinutes)
	fmt.Printf("start.Add(time.Hour * 10) = %v\n", afterTenHours)
	fmt.Printf("start.Add(time.Hour * 24 * 10) = %v\n", afterTenDays)

	// Output:
	// start = 2009-01-01 12:00:00 +0000 UTC
	// start.Add(time.Second * 10) = 2009-01-01 12:00:10 +0000 UTC
	// start.Add(time.Minute * 10) = 2009-01-01 12:10:00 +0000 UTC
	// start.Add(time.Hour * 10) = 2009-01-01 22:00:00 +0000 UTC
	// start.Add(time.Hour * 24 * 10) = 2009-01-11 12:00:00 +0000 UTC
}

func ExampleTime_AddDate() {
	start := synchro.New[tz.UTC](2009, 1, 1, 0, 0, 0, 0)
	oneDayLater := start.AddDate(0, 0, 1)
	oneMonthLater := start.AddDate(0, 1, 0)
	oneYearLater := start.AddDate(1, 0, 0)

	fmt.Printf("oneDayLater: start.AddDate(0, 0, 1) = %v\n", oneDayLater)
	fmt.Printf("oneMonthLater: start.AddDate(0, 1, 0) = %v\n", oneMonthLater)
	fmt.Printf("oneYearLater: start.AddDate(1, 0, 0) = %v\n", oneYearLater)

	// Output:
	// oneDayLater: start.AddDate(0, 0, 1) = 2009-01-02 00:00:00 +0000 UTC
	// oneMonthLater: start.AddDate(0, 1, 0) = 2009-02-01 00:00:00 +0000 UTC
	// oneYearLater: start.AddDate(1, 0, 0) = 2010-01-01 00:00:00 +0000 UTC
}

func ExampleTime_After() {
	year2000 := synchro.New[tz.UTC](2000, 1, 1, 0, 0, 0, 0)
	year3000 := synchro.New[tz.UTC](3000, 1, 1, 0, 0, 0, 0)

	isYear3000AfterYear2000 := year3000.After(year2000) // True
	isYear2000AfterYear3000 := year2000.After(year3000) // False

	fmt.Printf("year3000.After(year2000) = %v\n", isYear3000AfterYear2000)
	fmt.Printf("year2000.After(year3000) = %v\n", isYear2000AfterYear3000)

	// Output:
	// year3000.After(year2000) = true
	// year2000.After(year3000) = false
}

func ExampleTime_Before() {
	year2000 := synchro.New[tz.UTC](2000, 1, 1, 0, 0, 0, 0)
	year3000 := synchro.New[tz.UTC](3000, 1, 1, 0, 0, 0, 0)

	isYear2000BeforeYear3000 := year2000.Before(year3000) // True
	isYear3000BeforeYear2000 := year3000.Before(year2000) // False

	fmt.Printf("year2000.Before(year3000) = %v\n", isYear2000BeforeYear3000)
	fmt.Printf("year3000.Before(year2000) = %v\n", isYear3000BeforeYear2000)

	// Output:
	// year2000.Before(year3000) = true
	// year3000.Before(year2000) = false
}

func ExampleTime_Compare() {
	base := synchro.New[tz.UTC](2000, 1, 1, 1, 0, 0, 0)
	before := synchro.New[tz.UTC](2000, 1, 1, 0, 0, 0, 0)
	after := synchro.New[tz.UTC](3000, 1, 1, 2, 0, 0, 0)

	isSame := base.Compare(base)
	isAfter := base.Compare(after)
	isBefore := base.Compare(before)

	fmt.Printf("base.Compare(base) = %d\n", isSame)
	fmt.Printf("base.Compare(after) = %d\n", isAfter)
	fmt.Printf("base.Compare(before) = %d\n", isBefore)

	// Output:
	// base.Compare(base) = 0
	// base.Compare(after) = -1
	// base.Compare(before) = 1
}

func ExampleTime_Date() {
	d := synchro.New[tz.UTC](2000, 2, 1, 12, 30, 0, 0)
	year, month, day := d.Date()

	fmt.Printf("year = %v\n", year)
	fmt.Printf("month = %v\n", month)
	fmt.Printf("day = %v\n", day)

	// Output:
	// year = 2000
	// month = February
	// day = 1
}

func ExampleTime_Day() {
	d := synchro.New[tz.UTC](2000, 2, 1, 12, 30, 0, 0)
	day := d.Day()

	fmt.Printf("day = %v\n", day)

	// Output:
	// day = 1
}

func ExampleTime_Equal() {
	d1 := synchro.New[tz.UTC](2000, 2, 1, 12, 30, 0, 0)
	d2, _ := synchro.Parse[tz.UTC]("2006 Jan 02 15:04:05 (MST)", "2000 Feb 01 12:30:00 (JST)")

	datesEqualUsingEqualOperator := d1 == d2
	datesEqualUsingFunction := d1.Equal(d2)

	fmt.Printf("d1 = %q\n", d1)
	fmt.Printf("d2 = %q\n", d2)
	fmt.Printf("datesEqualUsingEqualOperator = %v\n", datesEqualUsingEqualOperator)
	fmt.Printf("datesEqualUsingFunction = %v\n", datesEqualUsingFunction)

	// Output:
	// d1 = "2000-02-01 12:30:00 +0000 UTC"
	// d2 = "2000-02-01 12:30:00 +0000 UTC"
	// datesEqualUsingEqualOperator = true
	// datesEqualUsingFunction = true
}

func ExampleTime_String() {
	timeWithNanoseconds := synchro.New[tz.UTC](2000, 2, 1, 12, 13, 14, 15)
	withNanoseconds := timeWithNanoseconds.String()

	timeWithoutNanoseconds := synchro.New[tz.UTC](2000, 2, 1, 12, 13, 14, 0)
	withoutNanoseconds := timeWithoutNanoseconds.String()

	fmt.Printf("withNanoseconds = %v\n", string(withNanoseconds))
	fmt.Printf("withoutNanoseconds = %v\n", string(withoutNanoseconds))

	// Output:
	// withNanoseconds = 2000-02-01 12:13:14.000000015 +0000 UTC
	// withoutNanoseconds = 2000-02-01 12:13:14 +0000 UTC
}

func ExampleTime_Sub() {
	start := synchro.New[tz.UTC](2000, 1, 1, 0, 0, 0, 0)
	end := synchro.New[tz.UTC](2000, 1, 1, 12, 0, 0, 0)

	difference := end.Sub(start)
	fmt.Printf("difference = %v\n", difference)

	// Output:
	// difference = 12h0m0s
}

func ExampleTime_AppendFormat() {
	t := synchro.New[tz.UTC](2017, time.November, 4, 11, 0, 0, 0)
	text := []byte("Time: ")

	text = t.AppendFormat(text, time.Kitchen)
	fmt.Println(string(text))

	// Output:
	// Time: 11:00AM
}

func ExampleConvertTz() {
	utc := synchro.New[tz.UTC](2009, time.November, 10, 23, 0, 0, 0)
	fmt.Printf("Go launched at %s\n", utc)

	jst := synchro.ConvertTz[tz.UTC, tz.AsiaTokyo](utc)
	fmt.Printf("Go launched at %s\n", jst)
	// Output:
	// Go launched at 2009-11-10 23:00:00 +0000 UTC
	// Go launched at 2009-11-11 08:00:00 +0900 JST
}

var c chan int

func handle(int) {}

func ExampleAfter() {
	select {
	case m := <-c:
		handle(m)
	case <-synchro.After[tz.UTC](time.Millisecond):
		fmt.Println("timed out")
	}
	// Output: timed out
}

func ExampleNewPeriod() {
	// with synchro.Time params
	p1, _ := synchro.NewPeriod[tz.UTC](
		synchro.New[tz.UTC](2009, 1, 1, 0, 0, 0, 0),
		synchro.New[tz.UTC](2009, 1, 10, 23, 59, 59, 0),
	)
	// with time.Time params
	p2, _ := synchro.NewPeriod[tz.UTC](
		time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2009, 1, 10, 23, 59, 59, 0, time.UTC),
	)
	// with ISO 8601 date and time format string or bytes
	p3, _ := synchro.NewPeriod[tz.UTC](
		"2009-01-01",
		[]byte("2009-01-10T23:59:59Z"),
	)
	fmt.Printf("p1: %s\n", p1)
	fmt.Printf("p2: %s\n", p2)
	fmt.Printf("p3: %s\n", p3)
	// Output:
	// p1: from 2009-01-01 00:00:00 +0000 UTC to 2009-01-10 23:59:59 +0000 UTC
	// p2: from 2009-01-01 00:00:00 +0000 UTC to 2009-01-10 23:59:59 +0000 UTC
	// p3: from 2009-01-01 00:00:00 +0000 UTC to 2009-01-10 23:59:59 +0000 UTC
}

func ExamplePeriod_PeriodicDuration() {
	p1, _ := synchro.NewPeriod[tz.UTC](
		"2009-01-01",
		synchro.New[tz.UTC](2009, 1, 1, 3, 59, 59, 0),
	)
	for current := range p1.PeriodicDuration(time.Hour) {
		fmt.Println(current)
	}
	// Output:
	// 2009-01-01 00:00:00 +0000 UTC
	// 2009-01-01 01:00:00 +0000 UTC
	// 2009-01-01 02:00:00 +0000 UTC
	// 2009-01-01 03:00:00 +0000 UTC
}

func ExamplePeriod_PeriodicDate() {
	p1, _ := synchro.NewPeriod[tz.UTC](
		"2009-01-01",
		"2009-01-03",
	)
	for current := range p1.PeriodicDate(0, 0, 1) {
		fmt.Println(current)
	}
	// Output:
	// 2009-01-01 00:00:00 +0000 UTC
	// 2009-01-02 00:00:00 +0000 UTC
	// 2009-01-03 00:00:00 +0000 UTC
}

func ExamplePeriod_PeriodicAdvance() {
	p1, _ := synchro.NewPeriod[tz.UTC](
		"2009-01-01",
		time.Date(2013, 1, 3, 0, 0, 0, 0, time.UTC),
	)
	for current := range p1.PeriodicAdvance(synchro.Year(1), synchro.Day(1)) {
		fmt.Println(current)
	}
	// Output:
	// 2009-01-01 00:00:00 +0000 UTC
	// 2010-01-02 00:00:00 +0000 UTC
	// 2011-01-03 00:00:00 +0000 UTC
	// 2012-01-04 00:00:00 +0000 UTC
}

func ExamplePeriod_PeriodicISODuration() {
	p1, _ := synchro.NewPeriod[tz.UTC](
		"2009-01-01",
		"2009-03-01",
	)
	iter, _ := p1.PeriodicISODuration("P1M")
	for current := range iter {
		fmt.Println(current)
	}
	// Output:
	// 2009-01-01 00:00:00 +0000 UTC
	// 2009-02-01 00:00:00 +0000 UTC
	// 2009-03-01 00:00:00 +0000 UTC
}

func ExamplePeriod_periodicalSlice() {
	p1, _ := synchro.NewPeriod[tz.UTC](
		"2009-01-01T00:00:00",
		"2009-01-01T03:00:00",
	)

	got := p1.PeriodicDuration(time.Hour).Slice()
	fmt.Println("len:", len(got))
	for _, current := range got {
		fmt.Println(current)
	}
	// Output:
	// len: 4
	// 2009-01-01 00:00:00 +0000 UTC
	// 2009-01-01 01:00:00 +0000 UTC
	// 2009-01-01 02:00:00 +0000 UTC
	// 2009-01-01 03:00:00 +0000 UTC
}
