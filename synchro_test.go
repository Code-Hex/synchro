package synchro_test

import (
	"testing"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

func TestParseISO(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cases := []struct {
			v    string
			want synchro.Time[tz.UTC]
		}{
			{
				v:    "2017-04-24T09:41:34.502+0100",
				want: synchro.New[tz.UTC](2017, 4, 24, 8, 41, 34, int(502*time.Millisecond)), // -1h
			},
			{
				v:    "2017-04-24T09:41+0100",
				want: synchro.New[tz.UTC](2017, 4, 24, 8, 41, 0, 0), // -1h
			},
			{
				v:    "2017-04-24T09+0100",
				want: synchro.New[tz.UTC](2017, 4, 24, 8, 0, 0, 0), // -1h
			},
			{
				v:    "2017-04-24",
				want: synchro.New[tz.UTC](2017, 4, 24, 0, 0, 0, 0),
			},
			{
				v:    "2017-04-24T09:41:34+0100",
				want: synchro.New[tz.UTC](2017, 4, 24, 8, 41, 34, 0), // -1h
			},
			{
				v:    "2017-04-24T09:41:34.502-0100",
				want: synchro.New[tz.UTC](2017, 4, 24, 10, 41, 34, int(502*time.Millisecond)), // +1h
			},
			{
				v:    "2017-04-24T09:41:34.502-01:00",
				want: synchro.New[tz.UTC](2017, 4, 24, 10, 41, 34, int(502*time.Millisecond)), // +1h
			},
			{
				v:    "2017-04-24T09:41-01:00",
				want: synchro.New[tz.UTC](2017, 4, 24, 10, 41, 0, 0), // +1h
			},
			{
				v:    "2017-04-24T09-01:00",
				want: synchro.New[tz.UTC](2017, 4, 24, 10, 0, 0, 0), // +1h
			},
			{
				v:    "2017-04-24T09:41:34-0100",
				want: synchro.New[tz.UTC](2017, 4, 24, 10, 41, 34, 0), // +1h
			},
			{
				v:    "2017-04-24T09:41:34.502Z",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, int(502*time.Millisecond)),
			},
			{
				v:    "2017-04-24T09:41:34Z",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, 0),
			},
			{
				v:    "2017-04-24T09:41Z",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 0, 0),
			},
			{
				v:    "2017-04-24T09Z",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 0, 0, 0),
			},
			{
				v:    "2017-04-24T09:41:34.089",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, int(89*time.Millisecond)),
			},
			{
				v:    "2017-04-24T09:41",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 0, 0),
			},
			{
				v:    "2017-04-24T09",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 0, 0, 0),
			},
			{
				v:    "2017-04-24T09:41:34.009",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, int(9*time.Millisecond)),
			},
			{
				v:    "2017-04-24T09:41:34.893",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, int(893*time.Millisecond)),
			},
			{
				v:    "2017-04-24T09:41:34.89312523Z",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, 89312523*10),
			},
			{
				v:    "2017-04-24T09:41:34.502-0530",
				want: synchro.New[tz.UTC](2017, 4, 24, 15, 11, 34, int(502*time.Millisecond)), // +5h30m
			},
			{
				v:    "2017-04-24T09:41:34.502+0530",
				want: synchro.New[tz.UTC](2017, 4, 24, 4, 11, 34, int(502*time.Millisecond)), // -5h30m
			},
			{
				v:    "2017-04-24T09:41:34.502+05:30",
				want: synchro.New[tz.UTC](2017, 4, 24, 4, 11, 34, int(502*time.Millisecond)), // +5h30m
			},
			{
				v:    "2017-04-24T09:41:34.502+05:45",
				want: synchro.New[tz.UTC](2017, 4, 24, 3, 56, 34, int(502*time.Millisecond)), // +5h45m
			},
			{
				v:    "2017-04-24T09:41:34.502+00",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, int(502*time.Millisecond)),
			},
			{
				v:    "2017-04-24T09:41:34.502+00",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, int(502*time.Millisecond)),
			},
			{
				v:    "2017-04-24T09:41:34.502+0000",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, int(502*time.Millisecond)),
			},
			{
				v:    "2017-04-24T09:41:34.502+00:00",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, int(502*time.Millisecond)),
			},
			{
				v:    "+2017-04-24T09:41:34.502+00:00",
				want: synchro.New[tz.UTC](2017, 4, 24, 9, 41, 34, int(502*time.Millisecond)),
			},
		}
		for _, tc := range cases {
			t.Run(tc.v, func(t *testing.T) {
				got, err := synchro.ParseISO[tz.UTC](tc.v)
				if err != nil {
					t.Fatal(err)
				}
				if !tc.want.Equal(got) {
					t.Fatalf("want %q but got %q", tc.want, got)
				}
			})
		}
	})

	t.Run("invalid", func(t *testing.T) {
		cases := []struct {
			v string
		}{
			// Invalid Parse Test Cases
			{
				v: "2017-+04-24T09:41:34.502-00:00",
			},
			// Invalid Range Test Cases
			{
				v: "2017-00-01T00:00:00.000+00:00",
				// month
			},
			{
				v: "2017-13-01T00:00:00.000+00:00",
				// month
			},
			{
				v: "2017-01-00T00:00:00.000+00:00",
				// day
			},
			{
				v: "2017-01-32T00:00:00.000+00:00",
				// day
			},
			{
				v: "2019-02-29T00:00:00.000+00:00",
				// day
			},
			{
				v: "2020-02-30T00:00:00.000+00:00", // Leap year
				// day
			},
			{
				v: "2017-01-01T25:00:00.000+00:00",
				// hour
			},
			{
				v: "2017-01-01T00:60:00.000+00:00",
				// minute
			},
			{
				v: "2017-01-01T00:00:60.000+00:00",
				// second
			},
		}
		for _, tc := range cases {
			t.Run(tc.v, func(t *testing.T) {
				tm, err := synchro.ParseISO[tz.UTC](tc.v)
				if err == nil {
					t.Fatalf("expected error")
				}
				if !tm.IsZero() {
					t.Fatalf("should be zero value")
				}
			})
		}

	})
}

func FuzzParseISO(f *testing.F) {
	f.Fuzz(func(t *testing.T, str string) {
		_, _ = synchro.ParseISO[tz.UTC](str)
	})
}
