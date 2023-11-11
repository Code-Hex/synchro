package synchro

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Code-Hex/synchro/tz"
	"github.com/google/go-cmp/cmp"
)

func TestCivilDates(t *testing.T) {
	for _, test := range []struct {
		date     CivilDate[tz.UTC]
		loc      *time.Location
		wantStr  string
		wantTime Time[tz.UTC]
	}{
		{
			date: CivilDate[tz.UTC]{
				Year:  2014,
				Month: 7,
				Day:   29,
			},
			wantStr:  "2014-07-29",
			wantTime: New[tz.UTC](2014, time.July, 29, 0, 0, 0, 0),
		},
		{
			date:     CivilDateOf[tz.UTC](time.Date(2014, 8, 20, 15, 8, 43, 1, time.Local)),
			loc:      time.UTC,
			wantStr:  "2014-08-20",
			wantTime: New[tz.UTC](2014, 8, 20, 0, 0, 0, 0),
		},
		{
			date:     CivilDateOf[tz.UTC](time.Date(999, time.January, 26, 0, 0, 0, 0, time.Local)),
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

func TestCivilDateIsValid(t *testing.T) {
	for _, test := range []struct {
		date CivilDate[tz.UTC]
		want bool
	}{
		{CivilDate[tz.UTC]{Year: 2014, Month: 7, Day: 29}, true},
		{CivilDate[tz.UTC]{Year: 2000, Month: 2, Day: 29}, true},
		{CivilDate[tz.UTC]{Year: 10000, Month: 1, Day: 1}, false},
		{CivilDate[tz.UTC]{Year: 1, Month: 1, Day: 1}, true},
		{CivilDate[tz.UTC]{Year: 0, Month: 1, Day: 1}, true}, // year zero is OK
		{CivilDate[tz.UTC]{Year: -1, Month: 1, Day: 1}, false},
		{CivilDate[tz.UTC]{Year: 1, Month: 0, Day: 1}, false},
		{CivilDate[tz.UTC]{Year: 1, Month: 1, Day: 0}, false},
		{CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 32}, false},
		{CivilDate[tz.UTC]{Year: 2016, Month: 13, Day: 1}, false},
		{CivilDate[tz.UTC]{Year: 1, Month: -1, Day: 1}, false},
		{CivilDate[tz.UTC]{Year: 1, Month: 1, Day: -1}, false},
	} {
		got := test.date.IsValid()
		if got != test.want {
			t.Errorf("%#v: got %t, want %t", test.date, got, test.want)
		}
	}
}

func TestParseCivilDate(t *testing.T) {
	for _, test := range []struct {
		str  string
		want CivilDate[tz.UTC] // if empty, expect an error
	}{
		{"2016-01-02", CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 2}},
		{"20160102", CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 2}},
		{"2016-002", CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 2}},
		{"2015-W53-6", CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 2}},
		{"2016-12-31", CivilDate[tz.UTC]{Year: 2016, Month: 12, Day: 31}},
		{"0003-02-04", CivilDate[tz.UTC]{Year: 3, Month: 2, Day: 4}},
		{"0003-02-04", CivilDate[tz.UTC]{Year: 3, Month: 2, Day: 4}},
		{"999-01-26", CivilDate[tz.UTC]{}},
		{"", CivilDate[tz.UTC]{}},
		{"2016-01-02x", CivilDate[tz.UTC]{}},
	} {
		got, err := ParseCivilDate[tz.UTC](test.str)
		if got != test.want {
			t.Errorf("ParseCivilDate(%q) = %+v, want %+v", test.str, got, test.want)
		}
		if err != nil && test.want != (CivilDate[tz.UTC]{}) {
			t.Errorf("Unexpected error %v from ParseCivilDate(%q)", err, test.str)
		}
	}
}

func TestCivilDateArithmetic(t *testing.T) {
	for _, test := range []struct {
		desc  string
		start CivilDate[tz.UTC]
		end   CivilDate[tz.UTC]
		days  int
	}{
		{
			desc:  "zero days noop",
			start: CivilDate[tz.UTC]{Year: 2014, Month: 5, Day: 9},
			end:   CivilDate[tz.UTC]{Year: 2014, Month: 5, Day: 9},
			days:  0,
		},
		{
			desc:  "crossing a year boundary",
			start: CivilDate[tz.UTC]{Year: 2014, Month: 12, Day: 31},
			end:   CivilDate[tz.UTC]{Year: 2015, Month: 1, Day: 1},
			days:  1,
		},
		{
			desc:  "negative number of days",
			start: CivilDate[tz.UTC]{Year: 2015, Month: 1, Day: 1},
			end:   CivilDate[tz.UTC]{Year: 2014, Month: 12, Day: 31},
			days:  -1,
		},
		{
			desc:  "full leap year",
			start: CivilDate[tz.UTC]{Year: 2004, Month: 1, Day: 1},
			end:   CivilDate[tz.UTC]{Year: 2005, Month: 1, Day: 1},
			days:  366,
		},
		{
			desc:  "full non-leap year",
			start: CivilDate[tz.UTC]{Year: 2001, Month: 1, Day: 1},
			end:   CivilDate[tz.UTC]{Year: 2002, Month: 1, Day: 1},
			days:  365,
		},
		{
			desc:  "crossing a leap second",
			start: CivilDate[tz.UTC]{Year: 1972, Month: 6, Day: 30},
			end:   CivilDate[tz.UTC]{Year: 1972, Month: 7, Day: 1},
			days:  1,
		},
		{
			desc:  "dates before the unix epoch",
			start: CivilDate[tz.UTC]{Year: 101, Month: 1, Day: 1},
			end:   CivilDate[tz.UTC]{Year: 102, Month: 1, Day: 1},
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

func TestCivilDateBefore(t *testing.T) {
	for _, test := range []struct {
		d1, d2 CivilDate[tz.UTC]
		want   bool
	}{
		{CivilDate[tz.UTC]{Year: 2016, Month: 12, Day: 31}, CivilDate[tz.UTC]{Year: 2017, Month: 1, Day: 1}, true},
		{CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 1}, CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 1}, false},
		{CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 1}, CivilDate[tz.UTC]{Year: 2016, Month: 2, Day: 1}, true},
		{CivilDate[tz.UTC]{Year: 2016, Month: 12, Day: 30}, CivilDate[tz.UTC]{Year: 2016, Month: 12, Day: 31}, true},
	} {
		if got := test.d1.Before(test.d2); got != test.want {
			t.Errorf("%v.Before(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}

func TestCivilDateAfter(t *testing.T) {
	for _, test := range []struct {
		d1, d2 CivilDate[tz.UTC]
		want   bool
	}{
		{CivilDate[tz.UTC]{Year: 2016, Month: 12, Day: 31}, CivilDate[tz.UTC]{Year: 2017, Month: 1, Day: 1}, false},
		{CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 1}, CivilDate[tz.UTC]{Year: 2016, Month: 1, Day: 1}, false},
		{CivilDate[tz.UTC]{Year: 2016, Month: 12, Day: 30}, CivilDate[tz.UTC]{Year: 2016, Month: 12, Day: 31}, false},
	} {
		if got := test.d1.After(test.d2); got != test.want {
			t.Errorf("%v.After(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}

func TestCivilDateIsZero(t *testing.T) {
	for _, test := range []struct {
		date CivilDate[tz.UTC]
		want bool
	}{
		{CivilDate[tz.UTC]{Year: 2000, Month: 2, Day: 29}, false},
		{CivilDate[tz.UTC]{Year: 10000, Month: 12, Day: 31}, false},
		{CivilDate[tz.UTC]{Year: -1, Month: 0, Day: 0}, false},
		{CivilDate[tz.UTC]{Year: 0, Month: 0, Day: 0}, true},
		{CivilDate[tz.UTC]{}, true},
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
		{CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15}, `"1987-04-15"`},
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
	var d CivilDate[tz.UTC]
	for _, test := range []struct {
		data string
		ptr  any
		want any
	}{
		{`"1987-04-15"`, &d, &CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15}},
		{`"19870415"`, &d, &CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15}},
		{`"1987-04-\u0031\u0035"`, &d, &CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15}},
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
			t.Errorf("%q, CivilDate: got nil, want error", bad)
		}
	}
}

func TestCivilDateChange(t *testing.T) {
	for _, test := range []struct {
		date CivilDate[tz.UTC]
		u1   Unit
		u2   []Unit
		want CivilDate[tz.UTC]
	}{
		{
			date: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Year(1988),
			want: CivilDate[tz.UTC]{Year: 1988, Month: 4, Day: 15},
		},
		{
			date: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Month(5),
			want: CivilDate[tz.UTC]{Year: 1987, Month: 5, Day: 15},
		},
		{
			date: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Day(10),
			want: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 10},
		},
		{
			date: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Year(1988),
			u2:   []Unit{Month(5), Day(10), Second(10)}, // ignore second
			want: CivilDate[tz.UTC]{Year: 1988, Month: 5, Day: 10},
		},
	} {
		got := test.date.Change(test.u1, test.u2...)
		if got != test.want {
			t.Errorf("Change(%q) = %+v, want %+v", fmt.Sprint(test.u1, test.u2), got, test.want)
		}
	}
}

func TestCivilDateAdvance(t *testing.T) {
	for _, test := range []struct {
		date CivilDate[tz.UTC]
		u1   Unit
		u2   []Unit
		want CivilDate[tz.UTC]
	}{
		{
			date: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Year(1),
			want: CivilDate[tz.UTC]{Year: 1988, Month: 4, Day: 15},
		},
		{
			date: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Month(1),
			want: CivilDate[tz.UTC]{Year: 1987, Month: 5, Day: 15},
		},
		{
			date: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Day(-5),
			want: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 10},
		},
		{
			date: CivilDate[tz.UTC]{Year: 1987, Month: 4, Day: 15},
			u1:   Year(1),
			u2:   []Unit{Month(1), Day(-5), Second(10)}, // ignore second
			want: CivilDate[tz.UTC]{Year: 1988, Month: 5, Day: 10},
		},
	} {
		got := test.date.Advance(test.u1, test.u2...)
		if got != test.want {
			t.Errorf("Advance(%q) = %+v, want %+v", fmt.Sprint(test.u1, test.u2), got, test.want)
		}
	}
}
