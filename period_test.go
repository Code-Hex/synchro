package synchro

import (
	"testing"
	"time"

	"github.com/Code-Hex/synchro/tz"
	"github.com/google/go-cmp/cmp"
)

func TestPeriod_Slice(t *testing.T) {
	want := []Time[tz.UTC]{
		New[tz.UTC](2014, 2, 5, 0, 0, 0, 0),
		New[tz.UTC](2014, 2, 6, 0, 0, 0, 0),
		New[tz.UTC](2014, 2, 7, 0, 0, 0, 0),
		New[tz.UTC](2014, 2, 8, 0, 0, 0, 0),
	}

	t.Run("Time[UTC] params", func(t *testing.T) {
		period, err := CreatePeriod[tz.UTC](
			New[tz.UTC](2014, 2, 5, 0, 0, 0, 0),
			New[tz.UTC](2014, 2, 8, 0, 0, 0, 0),
		)
		if err != nil {
			t.Fatal(err)
		}
		got := period.Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("time.Time params", func(t *testing.T) {
		period, err := CreatePeriod[tz.UTC](
			time.Date(2014, 2, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2014, 2, 8, 0, 0, 0, 0, time.UTC),
		)
		if err != nil {
			t.Fatal(err)
		}
		got := period.Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("string params", func(t *testing.T) {
		period, err := CreatePeriod[tz.UTC]("2014-02-05", "2014-02-08")
		if err != nil {
			t.Fatal(err)
		}
		got := period.Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("[]byte params", func(t *testing.T) {
		period, err := CreatePeriod[tz.UTC]([]byte("2014-02-05"), []byte("2014-02-08"))
		if err != nil {
			t.Fatal(err)
		}
		got := period.Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
	t.Run("complex params", func(t *testing.T) {
		period, err := CreatePeriod[tz.UTC](New[tz.UTC](2014, 2, 5, 0, 0, 0, 0), "2014-02-08")
		if err != nil {
			t.Fatal(err)
		}
		got := period.Slice()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("(-want, +got)\n%s", diff)
		}
	})
}
