package synchro_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

func ExampleTime_StdTime() {
	want := time.Date(2023, 9, 2, 14, 0, 0, 0, time.UTC)
	got1 := synchro.In[tz.UTC](want).StdTime()
	got2 := synchro.In[tz.AsiaTokyo](want).StdTime()
	fmt.Println(want.Equal(got1))
	fmt.Println(want.Equal(got2))
	// Output:
	// true
	// true
}

func TestTime_StartOfYear(t *testing.T) {
	d := synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999)
	want := synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0)
	got := d.StartOfYear()
	if !got.Equal(want) {
		t.Errorf("want %q but got %q", want, got)
	}
}

func TestTime_EndOfYear(t *testing.T) {
	d := synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0)
	want := synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999)
	got := d.EndOfYear()
	if !got.Equal(want) {
		t.Errorf("want %q but got %q", want, got)
	}
}

func TestTime_StartOfMonth(t *testing.T) {
	cases := []struct {
		time synchro.Time[tz.UTC]
		want synchro.Time[tz.UTC]
	}{
		{
			time: synchro.New[tz.UTC](2023, 1, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 2, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 2, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 3, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 3, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 4, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 4, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 5, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 5, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 6, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 7, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 8, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 8, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 10, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 10, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 11, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 11, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
			want: synchro.New[tz.UTC](2023, 12, 1, 0, 0, 0, 0),
		},
	}
	for _, tc := range cases {
		d := tc.time.StartOfMonth()
		if !d.Equal(tc.want) {
			t.Errorf("want %q but got %q", tc.want, d)
		}
	}
}

func TestTime_EndOfMonth(t *testing.T) {
	cases := []struct {
		time synchro.Time[tz.UTC]
		want synchro.Time[tz.UTC]
	}{
		{
			time: synchro.New[tz.UTC](2023, 1, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 2, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 2, 28, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 3, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 3, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 4, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 4, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 5, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 5, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 6, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 7, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 8, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 8, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 10, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 10, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 11, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 11, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
	}
	for _, tc := range cases {
		d := tc.time.EndOfMonth()
		if !d.Equal(tc.want) {
			t.Errorf("want %q but got %q", tc.want, d)
		}
	}
}

func TestTime_StartOfWeek(t *testing.T) {
	cases := []struct {
		time synchro.Time[tz.UTC]
		want synchro.Time[tz.UTC]
	}{
		{
			time: synchro.New[tz.UTC](2023, 9, 2, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 8, 27, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 3, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 3, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 4, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 3, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 5, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 3, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 6, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 3, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 7, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 3, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 8, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 3, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 9, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 3, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 10, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 10, 0, 0, 0, 0),
		},
	}
	for _, tc := range cases {
		d := tc.time.StartOfWeek()
		if !d.Equal(tc.want) {
			t.Errorf("want %q but got %q", tc.want, d)
		}
	}
}

func TestTime_EndOfWeek(t *testing.T) {
	cases := []struct {
		time synchro.Time[tz.UTC]
		want synchro.Time[tz.UTC]
	}{
		{
			time: synchro.New[tz.UTC](2023, 9, 2, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 2, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 3, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 9, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 4, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 9, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 5, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 9, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 6, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 9, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 7, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 9, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 8, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 9, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 9, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 9, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 10, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 16, 23, 59, 59, 999999999),
		},
	}
	for _, tc := range cases {
		d := tc.time.EndOfWeek()
		if !d.Equal(tc.want) {
			t.Errorf("want %q but got %q", tc.want, d)
		}
	}
}

func TestTime_StartOfQuarter(t *testing.T) {
	cases := []struct {
		time synchro.Time[tz.UTC]
		want synchro.Time[tz.UTC]
	}{
		{
			time: synchro.New[tz.UTC](2023, 1, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 2, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 3, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 4, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 4, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 5, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 4, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 6, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 4, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 7, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 8, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 10, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 10, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 11, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 10, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
			want: synchro.New[tz.UTC](2023, 10, 1, 0, 0, 0, 0),
		},
	}
	for _, tc := range cases {
		d := tc.time.StartOfQuarter()
		if !d.Equal(tc.want) {
			t.Errorf("want %q but got %q", tc.want, d)
		}
	}
}

func TestTime_EndOfQuarter(t *testing.T) {
	cases := []struct {
		time synchro.Time[tz.UTC]
		want synchro.Time[tz.UTC]
	}{
		{
			time: synchro.New[tz.UTC](2023, 1, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 3, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 2, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 3, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 3, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 3, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 4, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 5, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 6, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 7, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 8, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 9, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 10, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 11, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
	}
	for _, tc := range cases {
		d := tc.time.EndOfQuarter()
		if !d.Equal(tc.want) {
			t.Errorf("want %q but got %q", tc.want, d)
		}
	}
}

func TestTime_StartOfSemester(t *testing.T) {
	cases := []struct {
		time synchro.Time[tz.UTC]
		want synchro.Time[tz.UTC]
	}{
		{
			time: synchro.New[tz.UTC](2023, 1, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 2, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 3, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 4, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 5, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 6, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 7, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 8, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 10, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 11, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
		{
			time: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
			want: synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0),
		},
	}
	for _, tc := range cases {
		d := tc.time.StartOfSemester()
		if !d.Equal(tc.want) {
			t.Errorf("want %q but got %q", tc.want, d)
		}
	}
}

func TestTime_EndOfSemester(t *testing.T) {
	cases := []struct {
		time synchro.Time[tz.UTC]
		want synchro.Time[tz.UTC]
	}{
		{
			time: synchro.New[tz.UTC](2023, 1, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 2, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 3, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 4, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 5, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 6, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 7, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 8, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 9, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 10, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 11, 15, 0, 0, 0, 0),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
		{
			time: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
			want: synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999),
		},
	}
	for _, tc := range cases {
		d := tc.time.EndOfSemester()
		if !d.Equal(tc.want) {
			t.Errorf("want %q but got %q", tc.want, d)
		}
	}
}

func TestTime_IsLeapYear(t *testing.T) {
	cases := []struct {
		time synchro.Time[tz.UTC]
		want bool
	}{
		{
			time: synchro.New[tz.UTC](2000, 1, 1, 0, 0, 0, 0),
			want: true,
		},
		{
			time: synchro.New[tz.UTC](2004, 1, 1, 0, 0, 0, 0),
			want: true,
		},
		{
			time: synchro.New[tz.UTC](1900, 1, 1, 0, 0, 0, 0),
			want: false,
		},
		{
			time: synchro.New[tz.UTC](2001, 1, 1, 0, 0, 0, 0),
			want: false,
		},
	}
	for _, tc := range cases {
		got := tc.time.IsLeapYear()
		if got != tc.want {
			t.Errorf("for %q, want %v but got %v", tc.time, tc.want, got)
		}
	}
}

func TestTime_IsBetween(t *testing.T) {
	from := synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0)
	to := synchro.New[tz.UTC](2023, 12, 31, 23, 59, 59, 999999999)

	cases := []struct {
		time synchro.Time[tz.UTC]
		want bool
	}{
		{
			time: synchro.New[tz.UTC](2022, 12, 31, 23, 59, 59, 999999999),
			want: false,
		},
		{
			time: from,
			want: false,
		},
		{
			time: synchro.New[tz.UTC](2023, 6, 15, 0, 0, 0, 0),
			want: true,
		},
		{
			time: to,
			want: false,
		},
		{
			time: synchro.New[tz.UTC](2024, 1, 1, 0, 0, 0, 0),
			want: false,
		},
	}
	for _, tc := range cases {
		got := tc.time.IsBetween(from, to)
		if got != tc.want {
			t.Errorf("for %q, want %v but got %v", tc.time, tc.want, got)
		}
	}
}
