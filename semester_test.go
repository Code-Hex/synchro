package synchro_test

import (
	"fmt"
	"testing"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

func TestSemester(t *testing.T) {
	d := synchro.New[tz.UTC](2023, 1, 15, 0, 0, 0, 0)
	s := d.Semester()
	if want := s.Year(); want != d.Year() {
		t.Errorf("want year %q but got %q", want, s.Year())
	}
	if want := 1; want != s.Number() {
		t.Errorf("want the number of semester %q but got %q", want, s.Number())
	}
	if want := d.StartOfSemester(); !want.Equal(s.Start()) {
		t.Errorf("want start of semester %q but got %q", want, s.Start())
	}
	if want := d.EndOfSemester(); !want.Equal(s.End()) {
		t.Errorf("want end of semester %q but got %q", want, s.End())
	}
}

func ExampleSemester_After() {
	d1 := synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0)
	d2 := synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0)
	d3 := synchro.New[tz.UTC](2024, 1, 1, 0, 0, 0, 0)

	s1 := d1.Semester()
	s2 := d2.Semester()
	s3 := d3.Semester()

	fmt.Printf("s2.After(s1) = %v\n", s2.After(s1))
	fmt.Printf("s3.After(s1) = %v\n", s3.After(s1))
	fmt.Printf("s1.After(s2) = %v\n", s1.After(s2))
	// Output:
	// s2.After(s1) = true
	// s3.After(s1) = true
	// s1.After(s2) = false
}

func ExampleSemester_Before() {
	d1 := synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0)
	d2 := synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0)
	d3 := synchro.New[tz.UTC](2024, 1, 1, 0, 0, 0, 0)

	s1 := d1.Semester()
	s2 := d2.Semester()
	s3 := d3.Semester()

	fmt.Printf("s1.Before(s2) = %v\n", s1.Before(s2))
	fmt.Printf("s1.Before(s3) = %v\n", s1.Before(s3))
	fmt.Printf("s2.Before(s1) = %v\n", s2.Before(s1))
	// Output:
	// s1.Before(s2) = true
	// s1.Before(s3) = true
	// s2.Before(s1) = false
}

func ExampleSemester_Compare() {
	d1 := synchro.New[tz.UTC](2023, 1, 1, 0, 0, 0, 0)
	d2 := synchro.New[tz.UTC](2023, 6, 30, 23, 59, 59, 999999999)
	d3 := synchro.New[tz.UTC](2023, 7, 1, 0, 0, 0, 0)
	d4 := synchro.New[tz.UTC](2024, 1, 1, 0, 0, 0, 0)

	s1 := d1.Semester()
	s2 := d2.Semester()
	s3 := d3.Semester()
	s4 := d4.Semester()

	fmt.Printf("s1.Compare(s2) = %d\n", s1.Compare(s2))
	fmt.Printf("s1.Compare(s3) = %d\n", s1.Compare(s3))
	fmt.Printf("s4.Compare(s1) = %d\n", s4.Compare(s1))

	// Output:
	// s1.Compare(s2) = 0
	// s1.Compare(s3) = -1
	// s4.Compare(s1) = 1
}
