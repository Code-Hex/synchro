package synchro

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Code-Hex/synchro/tz"
	"github.com/google/go-cmp/cmp"
)

func TestDates(t *testing.T) {
	for _, test := range []struct {
		date     Date[tz.UTC]
		loc      *time.Location
		wantStr  string
		wantTime Time[tz.UTC]
	}{
		{
			date: Date[tz.UTC]{
				Year:  2014,
				Month: 7,
				Day:   29,
			},
			wantStr:  "2014-07-29",
			wantTime: New[tz.UTC](2014, time.July, 29, 0, 0, 0, 0),
		},
		{
			date:     DateOf[tz.UTC](time.Date(2014, 8, 20, 15, 8, 43, 1, time.Local)),
			loc:      time.UTC,
			wantStr:  "2014-08-20",
			wantTime: New[tz.UTC](2014, 8, 20, 0, 0, 0, 0),
		},
		{
			date:     DateOf[tz.UTC](time.Date(999, time.January, 26, 0, 0, 0, 0, time.Local)),
			wantStr:  "0999-01-26",
			wantTime: New[tz.UTC](999, 1, 26, 0, 0, 0, 0),
		},
	} {
		if got := test.date.String(); got != test.wantStr {
			t.Errorf("%#v.String() = %q, want %q", test.date, got, test.wantStr)
		}
		if got := test.date.Time(); !got.Equal(test.wantTime) {
			t.Errorf("%#v.In(%v) = %v, want %v", test.date, test.loc, got, test.wantTime)
		}
	}
}

func TestDateIsValid(t *testing.T) {
	for _, test := range []struct {
		date Date[tz.UTC]
		want bool
	}{
		{Date[tz.UTC]{Year: 2014, Month: 7, Day: 29}, true},
		{Date[tz.UTC]{Year: 2000, Month: 2, Day: 29}, true},
		{Date[tz.UTC]{Year: 10000, Month: 1, Day: 1}, false},
		{Date[tz.UTC]{Year: 1, Month: 1, Day: 1}, true},
		{Date[tz.UTC]{Year: 0, Month: 1, Day: 1}, true}, // year zero is OK
		{Date[tz.UTC]{Year: -1, Month: 1, Day: 1}, false},
		{Date[tz.UTC]{Year: 1, Month: 0, Day: 1}, false},
		{Date[tz.UTC]{Year: 1, Month: 1, Day: 0}, false},
		{Date[tz.UTC]{Year: 2016, Month: 1, Day: 32}, false},
		{Date[tz.UTC]{Year: 2016, Month: 13, Day: 1}, false},
		{Date[tz.UTC]{Year: 1, Month: -1, Day: 1}, false},
		{Date[tz.UTC]{Year: 1, Month: 1, Day: -1}, false},
	} {
		got := test.date.IsValid()
		if got != test.want {
			t.Errorf("%#v: got %t, want %t", test.date, got, test.want)
		}
	}
}

func TestParseDate(t *testing.T) {
	for _, test := range []struct {
		str  string
		want Date[tz.UTC] // if empty, expect an error
	}{
		{"2016-01-02", Date[tz.UTC]{Year: 2016, Month: 1, Day: 2}},
		{"20160102", Date[tz.UTC]{Year: 2016, Month: 1, Day: 2}},
		{"2016-002", Date[tz.UTC]{Year: 2016, Month: 1, Day: 2}},
		{"2015-W53-6", Date[tz.UTC]{Year: 2016, Month: 1, Day: 2}},
		{"2016-12-31", Date[tz.UTC]{Year: 2016, Month: 12, Day: 31}},
		{"0003-02-04", Date[tz.UTC]{Year: 3, Month: 2, Day: 4}},
		{"0003-02-04", Date[tz.UTC]{Year: 3, Month: 2, Day: 4}},
		{"999-01-26", Date[tz.UTC]{}},
		{"", Date[tz.UTC]{}},
		{"2016-01-02x", Date[tz.UTC]{}},
	} {
		got, err := ParseDate[tz.UTC](test.str)
		if got != test.want {
			t.Errorf("ParseDate(%q) = %+v, want %+v", test.str, got, test.want)
		}
		if err != nil && test.want != (Date[tz.UTC]{}) {
			t.Errorf("Unexpected error %v from ParseDate(%q)", err, test.str)
		}
	}
}

func TestDateArithmetic(t *testing.T) {
	for _, test := range []struct {
		desc  string
		start Date[tz.UTC]
		end   Date[tz.UTC]
		days  int
	}{
		{
			desc:  "zero days noop",
			start: Date[tz.UTC]{Year: 2014, Month: 5, Day: 9},
			end:   Date[tz.UTC]{Year: 2014, Month: 5, Day: 9},
			days:  0,
		},
		{
			desc:  "crossing a year boundary",
			start: Date[tz.UTC]{Year: 2014, Month: 12, Day: 31},
			end:   Date[tz.UTC]{Year: 2015, Month: 1, Day: 1},
			days:  1,
		},
		{
			desc:  "negative number of days",
			start: Date[tz.UTC]{Year: 2015, Month: 1, Day: 1},
			end:   Date[tz.UTC]{Year: 2014, Month: 12, Day: 31},
			days:  -1,
		},
		{
			desc:  "full leap year",
			start: Date[tz.UTC]{Year: 2004, Month: 1, Day: 1},
			end:   Date[tz.UTC]{Year: 2005, Month: 1, Day: 1},
			days:  366,
		},
		{
			desc:  "full non-leap year",
			start: Date[tz.UTC]{Year: 2001, Month: 1, Day: 1},
			end:   Date[tz.UTC]{Year: 2002, Month: 1, Day: 1},
			days:  365,
		},
		{
			desc:  "crossing a leap second",
			start: Date[tz.UTC]{Year: 1972, Month: 6, Day: 30},
			end:   Date[tz.UTC]{Year: 1972, Month: 7, Day: 1},
			days:  1,
		},
		{
			desc:  "dates before the unix epoch",
			start: Date[tz.UTC]{Year: 101, Month: 1, Day: 1},
			end:   Date[tz.UTC]{Year: 102, Month: 1, Day: 1},
			days:  365,
		},
	} {
		if got := test.start.AddDays(test.days); got != test.end {
			t.Errorf("[%s] %#v.AddDays(%v) = %#v, want %#v", test.desc, test.start, test.days, got, test.end)
		}
		if got := test.end.DaysSince(test.start); got != test.days {
			t.Errorf("[%s] %#v.Sub(%#v) = %v, want %v", test.desc, test.end, test.start, got, test.days)
		}
	}
}

func TestDateBefore(t *testing.T) {
	for _, test := range []struct {
		d1, d2 Date[tz.UTC]
		want   bool
	}{
		{Date[tz.UTC]{Year: 2016, Month: 12, Day: 31}, Date[tz.UTC]{Year: 2017, Month: 1, Day: 1}, true},
		{Date[tz.UTC]{Year: 2016, Month: 1, Day: 1}, Date[tz.UTC]{Year: 2016, Month: 1, Day: 1}, false},
		{Date[tz.UTC]{Year: 2016, Month: 1, Day: 1}, Date[tz.UTC]{Year: 2016, Month: 2, Day: 1}, true},
		{Date[tz.UTC]{Year: 2016, Month: 12, Day: 30}, Date[tz.UTC]{Year: 2016, Month: 12, Day: 31}, true},
	} {
		if got := test.d1.Before(test.d2); got != test.want {
			t.Errorf("%v.Before(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}

func TestDateAfter(t *testing.T) {
	for _, test := range []struct {
		d1, d2 Date[tz.UTC]
		want   bool
	}{
		{Date[tz.UTC]{Year: 2016, Month: 12, Day: 31}, Date[tz.UTC]{Year: 2017, Month: 1, Day: 1}, false},
		{Date[tz.UTC]{Year: 2016, Month: 1, Day: 1}, Date[tz.UTC]{Year: 2016, Month: 1, Day: 1}, false},
		{Date[tz.UTC]{Year: 2016, Month: 12, Day: 30}, Date[tz.UTC]{Year: 2016, Month: 12, Day: 31}, false},
	} {
		if got := test.d1.After(test.d2); got != test.want {
			t.Errorf("%v.After(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}

func TestDateIsZero(t *testing.T) {
	for _, test := range []struct {
		date Date[tz.UTC]
		want bool
	}{
		{Date[tz.UTC]{Year: 2000, Month: 2, Day: 29}, false},
		{Date[tz.UTC]{Year: 10000, Month: 12, Day: 31}, false},
		{Date[tz.UTC]{Year: -1, Month: 0, Day: 0}, false},
		{Date[tz.UTC]{Year: 0, Month: 0, Day: 0}, true},
		{Date[tz.UTC]{}, true},
	} {
		got := test.date.IsZero()
		if got != test.want {
			t.Errorf("%#v: got %t, want %t", test.date, got, test.want)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	for _, test := range []struct {
		value interface{}
		want  string
	}{
		{Date[tz.UTC]{Year: 1987, Month: 4, Day: 15}, `"1987-04-15"`},
	} {
		bgot, err := json.Marshal(test.value)
		if err != nil {
			t.Fatal(err)
		}
		if got := string(bgot); got != test.want {
			t.Errorf("%#v: got %s, want %s", test.value, got, test.want)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var d Date[tz.UTC]
	for _, test := range []struct {
		data string
		ptr  any
		want any
	}{
		{`"1987-04-15"`, &d, &Date[tz.UTC]{Year: 1987, Month: 4, Day: 15}},
		{`"19870415"`, &d, &Date[tz.UTC]{Year: 1987, Month: 4, Day: 15}},
		{`"1987-04-\u0031\u0035"`, &d, &Date[tz.UTC]{Year: 1987, Month: 4, Day: 15}},
	} {
		if err := json.Unmarshal([]byte(test.data), test.ptr); err != nil {
			t.Fatalf("%s: %v", test.data, err)
		}
		if !cmp.Equal(test.ptr, test.want) {
			t.Errorf("%s: got %#v, want %#v", test.data, test.ptr, test.want)
		}
	}

	for _, bad := range []string{"", `""`, `"bad"`, `"1987-04-15x"`,
		`19870415`,     // a JSON number
		`11987-04-15x`, // not a JSON string
	} {
		if json.Unmarshal([]byte(bad), &d) == nil {
			t.Errorf("%q, Date: got nil, want error", bad)
		}
	}
}

func TestDateChange(t *testing.T) {
	for _, test := range []struct {
		date Date[tz.UTC]
		u1   unit
		u2   []unit
		want Date[tz.UTC]
	}{
		{
			date: Date[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Year(1988),
			want: Date[tz.UTC]{Year: 1988, Month: 4, Day: 15},
		},
		{
			date: Date[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Month(5),
			want: Date[tz.UTC]{Year: 1987, Month: 5, Day: 15},
		},
		{
			date: Date[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Day(10),
			want: Date[tz.UTC]{Year: 1987, Month: 4, Day: 10},
		},
		{
			date: Date[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Year(1988),
			u2:   []unit{Month(5), Day(10), Second(10)}, // ignore second
			want: Date[tz.UTC]{Year: 1988, Month: 5, Day: 10},
		},
	} {
		got := test.date.Change(test.u1, test.u2...)
		if got != test.want {
			t.Errorf("Change(%q) = %+v, want %+v", fmt.Sprint(test.u1, test.u2), got, test.want)
		}
	}
}

func TestDateAdvance(t *testing.T) {
	for _, test := range []struct {
		date Date[tz.UTC]
		u1   unit
		u2   []unit
		want Date[tz.UTC]
	}{
		{
			date: Date[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Year(1),
			want: Date[tz.UTC]{Year: 1988, Month: 4, Day: 15},
		},
		{
			date: Date[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Month(1),
			want: Date[tz.UTC]{Year: 1987, Month: 5, Day: 15},
		},
		{
			date: Date[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Day(-5),
			want: Date[tz.UTC]{Year: 1987, Month: 4, Day: 10},
		},
		{
			date: Date[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Year(1),
			u2:   []unit{Month(1), Day(-5), Second(10)}, // ignore second
			want: Date[tz.UTC]{Year: 1988, Month: 5, Day: 10},
		},
	} {
		got := test.date.Advance(test.u1, test.u2...)
		if got != test.want {
			t.Errorf("Advance(%q) = %+v, want %+v", fmt.Sprint(test.u1, test.u2), got, test.want)
		}
	}
}
