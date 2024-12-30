package iso8601

import (
	"testing"
	"time"

	"database/sql"
	"database/sql/driver"
)

var _ interface {
	sql.Scanner
	driver.Valuer
} = (*Date)(nil)

var _ interface {
	sql.Scanner
	driver.Valuer
} = (*NullDate)(nil)

func TestDate_Scan(t *testing.T) {
	testCases := []struct {
		Name     string
		Value    any
		Error    bool
		Expected Date
	}{
		{
			Name:     "Valid string date",
			Value:    "2025-11-13",
			Expected: Date{2025, 11, 13},
		},
		{
			Name:     "Valid time.Time date",
			Value:    time.Date(2056, 10, 31, 0, 0, 0, 0, time.UTC),
			Expected: Date{2056, 10, 31},
		},
		{
			Name:     "Valid byte slice date",
			Value:    []byte("2157-12-31"),
			Expected: Date{2157, 12, 31},
		},
		{
			Name:  "Invalid byte slice date",
			Value: []byte("xxx"),
			Error: true,
		},
		{
			Name:     "Nil value",
			Value:    nil,
			Expected: Date{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var d Date
			err := d.Scan(tc.Value)
			if tc.Error {
				if err == nil {
					t.Errorf("expected error for value %v, but got none", tc.Value)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for value %v: %v", tc.Value, err)
				} else if d != tc.Expected {
					t.Errorf("expected %v, but got %v", tc.Expected, d)
				}
			}
		})
	}
}

func TestNullDate_Scan(t *testing.T) {
	testCases := []struct {
		Name     string
		Value    any
		Expected NullDate
	}{
		{
			Name:     "Valid string date",
			Value:    "2056-11-13",
			Expected: NullDate{Date: Date{2056, 11, 13}, Valid: true},
		},
		{
			Name:     "Nil value",
			Value:    nil,
			Expected: NullDate{Valid: false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var nd NullDate
			err := nd.Scan(tc.Value)
			if err != nil {
				t.Errorf("unexpected error for value %v: %v", tc.Value, err)
			} else if nd != tc.Expected {
				t.Errorf("expected %v, but got %v", tc.Expected, nd)
			}
		})
	}
}

func TestDate_Value(t *testing.T) {
	testCases := []struct {
		Name     string
		Date     Date
		Expected driver.Value
	}{
		{
			Name:     "Valid date",
			Date:     Date{2056, 11, 13},
			Expected: "2056-11-13",
		},
		{
			Name:     "Zero date",
			Date:     Date{},
			Expected: "0000-00-00",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			value, err := tc.Date.Value()
			if err != nil {
				t.Errorf("unexpected error for date %v: %v", tc.Date, err)
			}
			if value != tc.Expected {
				t.Errorf("expected %v, but got %v", tc.Expected, value)
			}
		})
	}
}

func TestNullDate_Value(t *testing.T) {
	testCases := []struct {
		Name     string
		NullDate NullDate
		Expected driver.Value
	}{
		{
			Name:     "Valid NullDate",
			NullDate: NullDate{Date: Date{2056, 11, 13}, Valid: true},
			Expected: "2056-11-13",
		},
		{
			Name:     "Invalid NullDate",
			NullDate: NullDate{Valid: false},
			Expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			value, err := tc.NullDate.Value()
			if err != nil {
				t.Errorf("unexpected error for NullDate %v: %v", tc.NullDate, err)
			}
			if value != tc.Expected {
				t.Errorf("expected %v, but got %v", tc.Expected, value)
			}
		})
	}
}
