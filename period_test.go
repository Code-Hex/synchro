package synchro

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Code-Hex/synchro/iso8601"
	"github.com/Code-Hex/synchro/tz"
	"github.com/google/go-cmp/cmp"
)

var _ interface {
	fmt.Stringer
} = (*Period[tz.UTC])(nil)

func TestNewPeriod(t *testing.T) {
	t.Run("want error for from", func(t *testing.T) {
		_, err := NewPeriod[tz.UTC]("unknown", "2023-09-27")
		if err == nil {
			t.Fatal("want error")
		}
	})
	t.Run("want error for to", func(t *testing.T) {
		_, err := NewPeriod[tz.UTC]("2023-09-27", "unknown")
		if err == nil {
			t.Fatal("want error")
		}
	})
}

func TestPeriod_Slice(t *testing.T) {
	want := []Time[tz.UTC]{
		New[tz.UTC](2014, 2, 5, 0, 0, 0, 0),
		New[tz.UTC](2014, 2, 6, 0, 0, 0, 0),
		New[tz.UTC](2014, 2, 7, 0, 0, 0, 0),
		New[tz.UTC](2014, 2, 8, 0, 0, 0, 0),
	}

	t.Run("Time[UTC] params", func(t *testing.T) {
		period, err := NewPeriod[tz.UTC](
			New[tz.UTC](2014, 2, 5, 0, 0, 0, 0),
			New[tz.UTC](2014, 2, 8, 0, 0, 0, 0),
		)
		if err != nil {
			t.Fatal(err)
		}
		got := period.PeriodicDuration(24 * time.Hour).Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("time.Time params", func(t *testing.T) {
		period, err := NewPeriod[tz.UTC](
			time.Date(2014, 2, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2014, 2, 8, 0, 0, 0, 0, time.UTC),
		)
		if err != nil {
			t.Fatal(err)
		}
		got := period.PeriodicDate(0, 0, 1).Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("string params", func(t *testing.T) {
		period, err := NewPeriod[tz.UTC]("2014-02-05", "2014-02-08")
		if err != nil {
			t.Fatal(err)
		}
		got := period.PeriodicAdvance(Day(1)).Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("alternative string params", func(t *testing.T) {
		type XString string
		period, err := NewPeriod[tz.UTC](XString("2014-02-05"), XString("2014-02-08"))
		if err != nil {
			t.Fatal(err)
		}
		got := period.PeriodicDuration(24 * time.Hour).Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("[]byte params", func(t *testing.T) {
		period, err := NewPeriod[tz.UTC]([]byte("2014-02-05"), []byte("2014-02-08"))
		if err != nil {
			t.Fatal(err)
		}
		got := period.PeriodicDate(0, 0, 1).Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("alternative []byte params", func(t *testing.T) {
		period, err := NewPeriod[tz.UTC]([]byte("2014-02-05"), []byte("2014-02-08"))
		if err != nil {
			t.Fatal(err)
		}
		got := period.PeriodicAdvance(Day(1)).Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("complex params", func(t *testing.T) {
		period, err := NewPeriod[tz.UTC](New[tz.UTC](2014, 2, 5, 0, 0, 0, 0), "2014-02-08")
		if err != nil {
			t.Fatal(err)
		}
		got := period.PeriodicDuration(24 * time.Hour).Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
}

func TestPeriod_PeriodicISO(t *testing.T) {
	cases := []struct {
		duration string
		from     string
		to       string
		want     []Time[tz.UTC]
	}{
		{
			duration: "P1Y",
			from:     "2018-08-16",
			to:       "2020-10-31",
			want: []Time[tz.UTC]{
				New[tz.UTC](2018, 8, 16, 0, 0, 0, 0),
				New[tz.UTC](2019, 8, 16, 0, 0, 0, 0),
				New[tz.UTC](2020, 8, 16, 0, 0, 0, 0),
			},
		},
		{
			duration: "P1M",
			from:     "2018-08-16",
			to:       "2018-10-31",
			want: []Time[tz.UTC]{
				New[tz.UTC](2018, 8, 16, 0, 0, 0, 0),
				New[tz.UTC](2018, 9, 16, 0, 0, 0, 0),
				New[tz.UTC](2018, 10, 16, 0, 0, 0, 0),
			},
		},
		{
			duration: "P3W",
			from:     "2018-08-16",
			to:       "2018-11-16",
			want: []Time[tz.UTC]{
				New[tz.UTC](2018, 8, 16, 0, 0, 0, 0),
				New[tz.UTC](2018, 9, 6, 0, 0, 0, 0),
				New[tz.UTC](2018, 9, 27, 0, 0, 0, 0),
				New[tz.UTC](2018, 10, 18, 0, 0, 0, 0),
				New[tz.UTC](2018, 11, 8, 0, 0, 0, 0),
			},
		},
		{
			duration: "P7D",
			from:     "2018-08-01",
			to:       "2018-08-31",
			want: []Time[tz.UTC]{
				New[tz.UTC](2018, 8, 1, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 8, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 15, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 22, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 29, 0, 0, 0, 0),
			},
		},
		{
			duration: "P1Y2M3W4D",
			from:     "2018-08-01",
			to:       "2020-08-31",
			want: []Time[tz.UTC]{
				New[tz.UTC](2018, 8, 1, 0, 0, 0, 0),
				New[tz.UTC](2019, 10, 1+21+4, 0, 0, 0, 0),
			},
		},
		{
			duration: "PT24H",
			from:     "2018-08-01",
			to:       "2018-08-05",
			want: []Time[tz.UTC]{
				New[tz.UTC](2018, 8, 1, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 2, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 3, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 4, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 5, 0, 0, 0, 0),
			},
		},
		{
			duration: "-P1Y",
			from:     "2020-10-31",
			to:       "2018-08-16",
			want: []Time[tz.UTC]{
				New[tz.UTC](2020, 10, 31, 0, 0, 0, 0),
				New[tz.UTC](2019, 10, 31, 0, 0, 0, 0),
				New[tz.UTC](2018, 10, 31, 0, 0, 0, 0),
			},
		},
		{
			duration: "-PT24H",
			from:     "2018-08-05",
			to:       "2018-08-01",
			want: []Time[tz.UTC]{
				New[tz.UTC](2018, 8, 5, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 4, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 3, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 2, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 1, 0, 0, 0, 0),
			},
		},
		{
			duration: "-P1DT1S",
			from:     "2018-08-05T00:00:00",
			to:       "2018-08-01T00:00:00",
			want: []Time[tz.UTC]{
				New[tz.UTC](2018, 8, 5, 0, 0, 0, 0),
				New[tz.UTC](2018, 8, 3, 23, 59, 59, 0),
				New[tz.UTC](2018, 8, 2, 23, 59, 58, 0),
				New[tz.UTC](2018, 8, 1, 23, 59, 57, 0),
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.duration, func(t *testing.T) {
			period, err := NewPeriod[tz.UTC](tc.from, tc.to)
			if err != nil {
				t.Fatal(err)
			}
			iter, err := period.PeriodicISODuration(tc.duration)
			if err != nil {
				t.Fatal(err)
			}
			got := iter.Slice()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("(-want, +got)\n%s", diff)
			}
		})
	}

	t.Run("invalid parse duration", func(t *testing.T) {
		period, err := NewPeriod[tz.UTC]("2023-08-16", "2023-08-17")
		if err != nil {
			t.Fatal(err)
		}
		_, err = period.PeriodicISODuration("unknown")
		if err == nil {
			t.Fatal("want error")
		}
		if _, ok := err.(*iso8601.UnexpectedTokenError); !ok {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid zero duration", func(t *testing.T) {
		period, err := NewPeriod[tz.UTC]("2023-08-16", "2023-08-17")
		if err != nil {
			t.Fatal(err)
		}

		t.Run("positive", func(t *testing.T) {
			duration := "P0D"
			_, err = period.PeriodicISODuration(duration)
			if err == nil {
				t.Fatal("want error")
			}
			if !strings.Contains(err.Error(), duration) {
				t.Fatalf("want contains %q in %q", duration, err)
			}
		})
		t.Run("negative", func(t *testing.T) {
			duration := "-PT0S"
			_, err = period.PeriodicISODuration(duration)
			if err == nil {
				t.Fatal("want error")
			}
			if !strings.Contains(err.Error(), duration) {
				t.Fatalf("want contains %q in %q", duration, err)
			}
		})
	})
}

// Unfortunate test to cover 100%
func TestConvertTime_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	convertTime[tz.UTC](true)
}

func TestPeriod_Contains(t *testing.T) {
	cases := []struct {
		name string
		t    Time[tz.UTC]
		from string
		to   string
		want int
	}{
		{
			name: "2018-08-16 is included between 2018-08-15 and 2018-08-17",
			t:    New[tz.UTC](2018, 8, 16, 0, 0, 0, 0),
			from: "2018-08-15",
			to:   "2018-08-17",
			want: 1,
		},
		{
			name: "2018-12-30 is included between 2018-12-28 and 2019-01-01",
			t:    New[tz.UTC](2018, 12, 30, 0, 0, 0, 0),
			from: "2018-12-28",
			to:   "2019-01-01",
			want: 1,
		},
		{
			name: "2018-08-15 is the same as from",
			t:    New[tz.UTC](2018, 8, 15, 0, 0, 0, 0),
			from: "2018-08-15",
			to:   "2018-08-17",
			want: 0,
		},
		{
			name: "2018-08-17 is the same as from",
			t:    New[tz.UTC](2018, 8, 17, 0, 0, 0, 0),
			from: "2018-08-15",
			to:   "2018-08-17",
			want: 0,
		},
		{
			name: "2018-08-18 is not included between 2018-08-15 and 2018-08-17",
			t:    New[tz.UTC](2018, 8, 18, 0, 0, 0, 0),
			from: "2018-08-15",
			to:   "2018-08-17",
			want: -1,
		},
		{
			name: "2018-08-14 is not included between 2018-08-15 and 2018-08-17",
			t:    New[tz.UTC](2018, 8, 14, 0, 0, 0, 0),
			from: "2018-08-15",
			to:   "2018-08-17",
			want: -1,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			period, err := NewPeriod[tz.UTC](tc.from, tc.to)
			if err != nil {
				t.Fatal(err)
			}
			got := period.Contains(tc.t)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestPeriod_From(t *testing.T) {
	want := New[tz.UTC](2018, 8, 14, 0, 0, 0, 0)
	period, err := NewPeriod[tz.UTC]("2018-08-14", "2018-08-16")
	if err != nil {
		t.Fatal(err)
	}
	got := period.From()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("(-want, +got)\n%s", diff)
	}
}

func TestPeriod_To(t *testing.T) {
	want := New[tz.UTC](2018, 8, 16, 0, 0, 0, 0)
	period, err := NewPeriod[tz.UTC]("2018-08-14", "2018-08-16")
	if err != nil {
		t.Fatal(err)
	}
	got := period.To()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("(-want, +got)\n%s", diff)
	}
}
