package synchro

import (
	"testing"
	"time"

	"github.com/Code-Hex/synchro/tz"
	"github.com/google/go-cmp/cmp"
)

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
	period, err := NewPeriod[tz.UTC]("2018-08-16", "2018-10-31")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("valid", func(t *testing.T) {
		iter, err := period.PeriodicISODuration("P1M")
		if err != nil {
			t.Fatal(err)
		}
		got := iter.Slice()
		want := []Time[tz.UTC]{
			New[tz.UTC](2018, 8, 16, 0, 0, 0, 0),
			New[tz.UTC](2018, 9, 16, 0, 0, 0, 0),
			New[tz.UTC](2018, 10, 16, 0, 0, 0, 0),
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
}
