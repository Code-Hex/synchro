package synchro_test

import (
	"fmt"
	"testing"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

func TestQuarter(t *testing.T) {
	d := synchro.New[tz.UTC](2023, 1, 15, 0, 0, 0, 0)
	q := d.Quarter()
	if want := q.Year(); want != d.Year() {
		t.Errorf("want year %q but got %q", want, q.Year())
	}
	if want := 1; want != q.Number() {
		t.Errorf("want the number of quarter %q but got %q", want, q.Number())
	}
	if want := d.StartOfQuarter(); !want.Equal(q.Start()) {
		t.Errorf("want start of quarter %q but got %q", want, q.Start())
	}
	if want := d.EndOfQuarter(); !want.Equal(q.End()) {
		t.Errorf("want end of quarter %q but got %q", want, q.End())
	}
}

func ExampleQuarter_After() {
	d1 := synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0)
	d2 := synchro.New[tz.UTC](2023, 4, 1, 0, 0, 0, 0)
	d3 := synchro.New[tz.UTC](2024, 1, 1, 0, 0, 0, 0)

	q1 := d1.Quarter()
	q2 := d2.Quarter()
	q3 := d3.Quarter()

	fmt.Printf("q2.After(q1) = %v\n", q2.After(q1))
	fmt.Printf("q3.After(q1) = %v\n", q3.After(q1))
	fmt.Printf("q1.After(q2) = %v\n", q1.After(q2))
	// Output:
	// q2.After(q1) = true
	// q3.After(q1) = true
	// q1.After(q2) = false
}

func ExampleQuarter_Before() {
	d1 := synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0)
	d2 := synchro.New[tz.UTC](2023, 4, 1, 0, 0, 0, 0)
	d3 := synchro.New[tz.UTC](2024, 1, 1, 0, 0, 0, 0)

	q1 := d1.Quarter()
	q2 := d2.Quarter()
	q3 := d3.Quarter()

	fmt.Printf("q1.Before(q2) = %v\n", q1.Before(q2))
	fmt.Printf("q1.Before(q3) = %v\n", q1.Before(q3))
	fmt.Printf("q2.Before(q1) = %v\n", q2.Before(q1))
	// Output:
	// q1.Before(q2) = true
	// q1.Before(q3) = true
	// q2.Before(q1) = false
}

func ExampleQuarter_Compare() {
	d1 := synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0)
	d2 := synchro.New[tz.UTC](2023, 3, 31, 23, 59, 59, 999999999)
	d3 := synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0)
	d4 := synchro.New[tz.UTC](2024, 1, 1, 0, 0, 0, 0)

	q1 := d1.Quarter()
	q2 := d2.Quarter()
	q3 := d3.Quarter()
	q4 := d4.Quarter()

	fmt.Printf("q1.Compare(q2) = %d\n", q1.Compare(q2))
	fmt.Printf("q1.Compare(q3) = %d\n", q1.Compare(q3))
	fmt.Printf("q4.Compare(q1) = %d\n", q4.Compare(q1))

	// Output:
	// q1.Compare(q2) = 0
	// q1.Compare(q3) = -1
	// q4.Compare(q1) = 1
}
