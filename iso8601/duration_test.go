package iso8601

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		want    Duration
		wantErr error
	}{
		{
			name: "PT0S",
			want: Duration{},
		},
		{
			name: "P0D",
			want: Duration{},
		},
		{
			name: "P1Y",
			want: Duration{
				Year: 1,
			},
		},
		{
			name: "-P1Y",
			want: Duration{
				Year:     1,
				Negative: true,
			},
		},
		{
			name: "P1234567890Y",
			want: Duration{
				Year: 1234567890,
			},
		},
		{
			name: "-P1234567890Y",
			want: Duration{
				Year:     1234567890,
				Negative: true,
			},
		},
		{
			name: "P1Y2M",
			want: Duration{
				Year:  1,
				Month: time.February,
			},
		},
		{
			name: "-P1Y2M",
			want: Duration{
				Year:     1,
				Month:    time.February,
				Negative: true,
			},
		},
		{
			name: "P2M",
			want: Duration{
				Month: time.February,
			},
		},
		{
			name: "-P2M",
			want: Duration{
				Month:    time.February,
				Negative: true,
			},
		},
		{
			name: "P1234567890M",
			want: Duration{
				Month: time.Month(1234567890),
			},
		},
		{
			name: "-P1234567890M",
			want: Duration{
				Month:    time.Month(1234567890),
				Negative: true,
			},
		},
		{
			name: "P1Y2M3W",
			want: Duration{
				Year:  1,
				Month: time.February,
				Week:  3,
			},
		},
		{
			name: "-P1Y2M3W",
			want: Duration{
				Year:     1,
				Month:    time.February,
				Week:     3,
				Negative: true,
			},
		},
		{
			name: "P3W",
			want: Duration{
				Week: 3,
			},
		},
		{
			name: "-P3W",
			want: Duration{
				Week:     3,
				Negative: true,
			},
		},
		{
			name: "P1Y3W",
			want: Duration{
				Year: 1,
				Week: 3,
			},
		},
		{
			name: "-P1Y3W",
			want: Duration{
				Year:     1,
				Week:     3,
				Negative: true,
			},
		},
		{
			name: "P2M3W",
			want: Duration{
				Month: time.February,
				Week:  3,
			},
		},
		{
			name: "-P2M3W",
			want: Duration{
				Month:    time.February,
				Week:     3,
				Negative: true,
			},
		},
		{
			name: "P1234567890W",
			want: Duration{
				Week: 1234567890,
			},
		},
		{
			name: "-P1234567890W",
			want: Duration{
				Week:     1234567890,
				Negative: true,
			},
		},
		{
			name: "P1Y2M3W4D",
			want: Duration{
				Year:  1,
				Month: time.February,
				Week:  3,
				Day:   4,
			},
		},
		{
			name: "-P1Y2M3W4D",
			want: Duration{
				Year:     1,
				Month:    time.February,
				Week:     3,
				Day:      4,
				Negative: true,
			},
		},
		{
			name: "P1234567890D",
			want: Duration{
				Day: 1234567890,
			},
		},
		{
			name: "-P1234567890D",
			want: Duration{
				Day:      1234567890,
				Negative: true,
			},
		},
		{
			name: "P4D",
			want: Duration{
				Day: 4,
			},
		},
		{
			name: "-P4D",
			want: Duration{
				Day:      4,
				Negative: true,
			},
		},
		{
			name: "P1Y4D",
			want: Duration{
				Year: 1,
				Day:  4,
			},
		},
		{
			name: "-P1Y4D",
			want: Duration{
				Year:     1,
				Day:      4,
				Negative: true,
			},
		},
		{
			name: "P2M4D",
			want: Duration{
				Month: time.February,
				Day:   4,
			},
		},
		{
			name: "-P2M4D",
			want: Duration{
				Month:    time.February,
				Day:      4,
				Negative: true,
			},
		},
		{
			name: "P3W4D",
			want: Duration{
				Week: 3,
				Day:  4,
			},
		},
		{
			name: "-P3W4D",
			want: Duration{
				Week:     3,
				Day:      4,
				Negative: true,
			},
		},
		{
			name: "P3DT1H",
			want: Duration{
				Day:  3,
				Hour: 1,
			},
		},
		{
			name: "-PT2H20M30S",
			want: Duration{
				Hour:     2,
				Minute:   20,
				Second:   30,
				Negative: true,
			},
		},
		{
			name: "PT6M",
			want: Duration{
				Minute: 6,
			},
		},
		{
			name: "-PT6M",
			want: Duration{
				Minute:   6,
				Negative: true,
			},
		},
		{
			name: "PT5H6M",
			want: Duration{
				Hour:   5,
				Minute: 6,
			},
		},
		{
			name: "-PT5H6M",
			want: Duration{
				Hour:     5,
				Minute:   6,
				Negative: true,
			},
		},
		{
			name: "P3WT6M",
			want: Duration{
				Week:   3,
				Minute: 6,
			},
		},
		{
			name: "-P3WT6M",
			want: Duration{
				Week:     3,
				Minute:   6,
				Negative: true,
			},
		},
		{
			name: "P4DT6M",
			want: Duration{
				Day:    4,
				Minute: 6,
			},
		},
		{
			name: "-P4DT6M",
			want: Duration{
				Day:      4,
				Minute:   6,
				Negative: true,
			},
		},
		{
			name: "PT7S",
			want: Duration{
				Second: 7,
			},
		},
		{
			name: "-PT7S",
			want: Duration{
				Second:   7,
				Negative: true,
			},
		},
		{
			name: "PT5H7S",
			want: Duration{
				Hour:   5,
				Second: 7,
			},
		},
		{
			name: "-PT5H7S",
			want: Duration{
				Hour:     5,
				Second:   7,
				Negative: true,
			},
		},
		{
			name: "PT6M7S",
			want: Duration{
				Minute: 6,
				Second: 7,
			},
		},
		{
			name: "-PT6M7S",
			want: Duration{
				Minute:   6,
				Second:   7,
				Negative: true,
			},
		},
		{
			name: "PT5H6M7S",
			want: Duration{
				Hour:   5,
				Minute: 6,
				Second: 7,
			},
		},
		{
			name: "-PT5H6M7S",
			want: Duration{
				Hour:     5,
				Minute:   6,
				Second:   7,
				Negative: true,
			},
		},
		{
			name: "P1YT5H6M7S",
			want: Duration{
				Year:   1,
				Hour:   5,
				Minute: 6,
				Second: 7,
			},
		},
		{
			name: "-P1YT5H6M7S",
			want: Duration{
				Year:     1,
				Hour:     5,
				Minute:   6,
				Second:   7,
				Negative: true,
			},
		},
		{
			name: "PT0.008S",
			want: Duration{
				Millisecond: 8,
			},
		},
		{
			name: "-PT0.008S",
			want: Duration{
				Millisecond: 8,
				Negative:    true,
			},
		},
		{
			name: "PT0.08S",
			want: Duration{
				Millisecond: 80,
			},
		},
		{
			name: "-PT0.08S",
			want: Duration{
				Millisecond: 80,
				Negative:    true,
			},
		},
		{
			name: "PT0.087S",
			want: Duration{
				Millisecond: 87,
			},
		},
		{
			name: "-PT0.087S",
			want: Duration{
				Millisecond: 87,
				Negative:    true,
			},
		},
		{
			name: "PT0.876S",
			want: Duration{
				Millisecond: 876,
			},
		},
		{
			name: "-PT0.876S",
			want: Duration{
				Millisecond: 876,
				Negative:    true,
			},
		},
		{
			name: "PT876.543S",
			want: Duration{
				Second:      876,
				Millisecond: 543,
			},
		},
		{
			name: "-PT876.543S",
			want: Duration{
				Second:      876,
				Millisecond: 543,
				Negative:    true,
			},
		},
		{
			name: "PT0.000009S",
			want: Duration{
				Microsecond: 9,
			},
		},
		{
			name: "-PT0.000009S",
			want: Duration{
				Microsecond: 9,
				Negative:    true,
			},
		},
		{
			name: "PT0.00009S",
			want: Duration{
				Microsecond: 90,
			},
		},
		{
			name: "-PT0.00009S",
			want: Duration{
				Microsecond: 90,
				Negative:    true,
			},
		},
		{
			name: "PT0.000098S",
			want: Duration{
				Microsecond: 98,
			},
		},
		{
			name: "-PT0.000098S",
			want: Duration{
				Microsecond: 98,
				Negative:    true,
			},
		},
		{
			name: "PT0.000987S",
			want: Duration{
				Microsecond: 987,
			},
		},
		{
			name: "-PT0.000987S",
			want: Duration{
				Microsecond: 987,
				Negative:    true,
			},
		},
		{
			name: "PT0.987654S",
			want: Duration{
				Millisecond: 987,
				Microsecond: 654,
			},
		},
		{
			name: "-PT0.987654S",
			want: Duration{
				Millisecond: 987,
				Microsecond: 654,
				Negative:    true,
			},
		},
		{
			name: "PT987.654321S",
			want: Duration{
				Second:      987,
				Millisecond: 654,
				Microsecond: 321,
			},
		},
		{
			name: "-PT987.654321S",
			want: Duration{
				Second:      987,
				Millisecond: 654,
				Microsecond: 321,
				Negative:    true,
			},
		},
		{
			name: "PT0.000000001S",
			want: Duration{
				Nanosecond: 1,
			},
		},
		{
			name: "-PT0.000000001S",
			want: Duration{
				Nanosecond: 1,
				Negative:   true,
			},
		},
		{
			name: "PT0.00000001S",
			want: Duration{
				Nanosecond: 10,
			},
		},
		{
			name: "-PT0.00000001S",
			want: Duration{
				Nanosecond: 10,
				Negative:   true,
			},
		},
		{
			name: "PT0.000000012S",
			want: Duration{
				Nanosecond: 12,
			},
		},
		{
			name: "-PT0.000000012S",
			want: Duration{
				Nanosecond: 12,
				Negative:   true,
			},
		},
		{
			name: "PT0.000000123S",
			want: Duration{
				Nanosecond: 123,
			},
		},
		{
			name: "-PT0.000000123S",
			want: Duration{
				Nanosecond: 123,
				Negative:   true,
			},
		},
		{
			name: "PT0.000123456S",
			want: Duration{
				Microsecond: 123,
				Nanosecond:  456,
			},
		},
		{
			name: "-PT0.000123456S",
			want: Duration{
				Microsecond: 123,
				Nanosecond:  456,
				Negative:    true,
			},
		},
		{
			name: "PT0.123456789S",
			want: Duration{
				Millisecond: 123,
				Microsecond: 456,
				Nanosecond:  789,
			},
		},
		{
			name: "-PT0.123456789S",
			want: Duration{
				Millisecond: 123,
				Microsecond: 456,
				Nanosecond:  789,
				Negative:    true,
			},
		},
		{
			name: "PT1.234567891S",
			want: Duration{
				Second:      1,
				Millisecond: 234,
				Microsecond: 567,
				Nanosecond:  891,
			},
		},
		{
			name: "-PT1.234567891S",
			want: Duration{
				Second:      1,
				Millisecond: 234,
				Microsecond: 567,
				Nanosecond:  891,
				Negative:    true,
			},
		},
		{
			name: "PT4.003002001S",
			want: Duration{
				Second:      4,
				Millisecond: 3,
				Microsecond: 2,
				Nanosecond:  1,
			},
		},
		{
			name: "-PT4.003002001S",
			want: Duration{
				Second:      4,
				Millisecond: 3,
				Microsecond: 2,
				Nanosecond:  1,
				Negative:    true,
			},
		},
		{
			name: "PT4.003092001S",
			want: Duration{
				Second:      4,
				Millisecond: 3,
				Microsecond: 92,
				Nanosecond:  1,
			},
		},
		{
			name: "+PT4.003092001S",
			want: Duration{
				Second:      4,
				Millisecond: 3,
				Microsecond: 92,
				Nanosecond:  1,
			},
		},
		{
			name: "-PT4.003092001S",
			want: Duration{
				Second:      4,
				Millisecond: 3,
				Microsecond: 92,
				Nanosecond:  1,
				Negative:    true,
			},
		},
		{
			name: "PT4.093082001S",
			want: Duration{
				Second:      4,
				Millisecond: 93,
				Microsecond: 82,
				Nanosecond:  1,
			},
		},
		{
			name: "-PT4.093082001S",
			want: Duration{
				Second:      4,
				Millisecond: 93,
				Microsecond: 82,
				Nanosecond:  1,
				Negative:    true,
			},
		},
		{
			name: "P1Y2M3W4DT5H6M7.008009001S",
			want: Duration{
				Year:        1,
				Month:       time.February,
				Week:        3,
				Day:         4,
				Hour:        5,
				Minute:      6,
				Second:      7,
				Millisecond: 8,
				Microsecond: 9,
				Nanosecond:  1,
			},
		},
		{
			name: "-P1Y2M3W4DT5H6M7.008009001S",
			want: Duration{
				Year:        1,
				Month:       time.February,
				Week:        3,
				Day:         4,
				Hour:        5,
				Minute:      6,
				Second:      7,
				Millisecond: 8,
				Microsecond: 9,
				Nanosecond:  1,
				Negative:    true,
			},
		},
		{
			name: "P1234Y2345M3456W4567DT5678H6789M7890.890901123S",
			want: Duration{
				Year:        1234,
				Month:       time.Month(2345),
				Week:        3456,
				Day:         4567,
				Hour:        5678,
				Minute:      6789,
				Second:      7890,
				Millisecond: 890,
				Microsecond: 901,
				Nanosecond:  123,
			},
		},
		{
			name: "-P1234Y2345M3456W4567DT5678H6789M7890.890901123S",
			want: Duration{
				Year:        1234,
				Month:       time.Month(2345),
				Week:        3456,
				Day:         4567,
				Hour:        5678,
				Minute:      6789,
				Second:      7890,
				Millisecond: 890,
				Microsecond: 901,
				Nanosecond:  123,
				Negative:    true,
			},
		},
		{
			name: "P1Y1M1W1DT1H1M1.001001001S",
			want: Duration{
				Year:        1,
				Month:       time.January,
				Week:        1,
				Day:         1,
				Hour:        1,
				Minute:      1,
				Second:      1,
				Millisecond: 1,
				Microsecond: 1,
				Nanosecond:  1,
			},
		},
		{
			name: "+P1Y1M1W1DT1H1M1.001001001S",
			want: Duration{
				Year:        1,
				Month:       time.January,
				Week:        1,
				Day:         1,
				Hour:        1,
				Minute:      1,
				Second:      1,
				Millisecond: 1,
				Microsecond: 1,
				Nanosecond:  1,
			},
		},
		{
			name: "-P1Y1M1W1DT1H1M1.001001001S",
			want: Duration{
				Year:        1,
				Month:       time.January,
				Week:        1,
				Day:         1,
				Hour:        1,
				Minute:      1,
				Second:      1,
				Millisecond: 1,
				Microsecond: 1,
				Nanosecond:  1,
				Negative:    true,
			},
		},
		{
			name: "PT0.999999999H",
			want: Duration{
				Minute:      59,
				Second:      59,
				Millisecond: 999,
				Microsecond: 996,
				Nanosecond:  400,
			},
		},
		{
			name: "PT0.000000011H",
			want: Duration{
				Microsecond: 39,
				Nanosecond:  600,
			},
		},
		{
			name: "PT0.999999999M",
			want: Duration{
				Second:      59,
				Millisecond: 999,
				Microsecond: 999,
				Nanosecond:  940,
			},
		},
		{
			name: "PT0.000000011M",
			want: Duration{
				Nanosecond: 660,
			},
		},
		{
			name: "PT0.999999999S",
			want: Duration{
				Millisecond: 999,
				Microsecond: 999,
				Nanosecond:  999,
			},
		},
		{
			name: "PT0.000000011S",
			want: Duration{
				Nanosecond: 11,
			},
		},
		{
			name: "PT1.03125H",
			want: Duration{
				Hour:        1,
				Minute:      1,
				Second:      52,
				Millisecond: 500,
			},
		},
		{
			name: "-PT1.03125H",
			want: Duration{
				Hour:        1,
				Minute:      1,
				Second:      52,
				Millisecond: 500,
				Negative:    true,
			},
		},
		{
			name: "PT46H66M71.50040904S",
			want: Duration{
				Hour:        46,
				Minute:      66,
				Second:      71,
				Millisecond: 500,
				Microsecond: 409,
				Nanosecond:  40,
			},
		},
		{
			name: "-PT46H66M71.50040904S",
			want: Duration{
				Hour:        46,
				Minute:      66,
				Second:      71,
				Millisecond: 500,
				Microsecond: 409,
				Nanosecond:  40,
				Negative:    true,
			},
		},
		{
			name: "",
			wantErr: &UnexpectedTokenError{
				Value:      "",
				Token:      "",
				AfterToken: "",
				Expected:   "P or + or -",
			},
		},
		{
			name: "X",
			wantErr: &UnexpectedTokenError{
				Value:      "X",
				Token:      "X",
				AfterToken: "",
				Expected:   "P or + or -",
			},
		},
		{
			name: "+",
			wantErr: &UnexpectedTokenError{
				Value:      "+",
				Token:      "+",
				AfterToken: "+",
				Expected:   "P",
			},
		},
		{
			name: "-",
			wantErr: &UnexpectedTokenError{
				Value:      "-",
				Token:      "-",
				AfterToken: "-",
				Expected:   "P",
			},
		},
		{
			name: "P1Y1Y",
			wantErr: &UnexpectedTokenError{
				Value:      "P1Y1Y",
				Token:      "Y",
				AfterToken: "P1Y1",
				Expected:   "the designator to be used only once",
			},
		},
		{
			name: "P1Y2M1Y",
			wantErr: &UnexpectedTokenError{
				Value:      "P1Y2M1Y",
				Token:      "Y",
				AfterToken: "P1Y2M1",
				Expected:   "the designator to be used only once",
			},
		},
		{
			name: "P1M1Y",
			wantErr: &UnexpectedTokenError{
				Value:      "P1M1Y",
				Token:      "Y",
				AfterToken: "P1M1",
				Expected:   "the 'Y' date designator should appear before 'M'",
			},
		},
		{
			name: "P1W1Y",
			wantErr: &UnexpectedTokenError{
				Value:      "P1W1Y",
				Token:      "Y",
				AfterToken: "P1W1",
				Expected:   "the 'Y' date designator should appear before 'W'",
			},
		},
		{
			name: "P1D1Y",
			wantErr: &UnexpectedTokenError{
				Value:      "P1D1Y",
				Token:      "Y",
				AfterToken: "P1D1",
				Expected:   "the 'Y' date designator should appear before 'D'",
			},
		},
		{
			name: "P1W1M",
			wantErr: &UnexpectedTokenError{
				Value:      "P1W1M",
				Token:      "M",
				AfterToken: "P1W1",
				Expected:   "the 'M' date designator should appear before 'W'",
			},
		},
		{
			name: "P1D1M",
			wantErr: &UnexpectedTokenError{
				Value:      "P1D1M",
				Token:      "M",
				AfterToken: "P1D1",
				Expected:   "the 'M' date designator should appear before 'D'",
			},
		},
		{
			name: "P1D1W",
			wantErr: &UnexpectedTokenError{
				Value:      "P1D1W",
				Token:      "W",
				AfterToken: "P1D1",
				Expected:   "the 'W' date designator should appear before 'D'",
			},
		},
		{
			name: "PT1D",
			wantErr: &UnexpectedTokenError{
				Value:      "PT1D",
				Token:      "D",
				AfterToken: "PT1",
				Expected:   "date duration should put before the 'T'",
			},
		},
		{
			name: "P1S",
			wantErr: &UnexpectedTokenError{
				Value:      "P1S",
				Token:      "S",
				AfterToken: "P1",
				Expected:   "the 'T' designator is required",
			},
		},
		{
			name: "PT1S1M",
			wantErr: &UnexpectedTokenError{
				Value:      "PT1S1M",
				Token:      "M",
				AfterToken: "PT1S1",
				Expected:   "the 'M' time designator should appear before 'S'",
			},
		},
		{
			name: "PT1S1H",
			wantErr: &UnexpectedTokenError{
				Value:      "PT1S1H",
				Token:      "H",
				AfterToken: "PT1S1",
				Expected:   "the 'H' time designator should appear before 'S'",
			},
		},
		{
			name: "PT1M1H",
			wantErr: &UnexpectedTokenError{
				Value:      "PT1M1H",
				Token:      "H",
				AfterToken: "PT1M1",
				Expected:   "the 'H' time designator should appear before 'M'",
			},
		},
		{
			name: "P1X",
			wantErr: &UnexpectedTokenError{
				Value:      "P1X",
				Token:      "X",
				AfterToken: "P1",
				Expected:   "PnYnMnDTnHnMnS or PnW format",
			},
		},
		{
			name: "PTT1H",
			wantErr: &UnexpectedTokenError{
				Value:      "PTT1H",
				Token:      "T",
				AfterToken: "PT",
				Expected:   "the 'T' designator should be once",
			},
		},
		{
			name: "PTN",
			wantErr: &UnexpectedTokenError{
				Value:      "PTN",
				Token:      "N",
				AfterToken: "PT",
				Expected:   "PnYnMnDTnHnMnS or PnW format",
			},
		},
		// unexpected fraction
		{
			name: "P1.123Y",
			wantErr: &UnexpectedTokenError{
				Value:      "P1.123Y",
				Token:      ".",
				AfterToken: "P1",
				Expected:   "only the time unit can be fractional",
			},
		},
		{
			name: "P12.123W",
			wantErr: &UnexpectedTokenError{
				Value:      "P12.123W",
				Token:      ".",
				AfterToken: "P12",
				Expected:   "only the time unit can be fractional",
			},
		},
		{
			name: "P123,123D",
			wantErr: &UnexpectedTokenError{
				Value:      "P123,123D",
				Token:      ",",
				AfterToken: "P123",
				Expected:   "only the time unit can be fractional",
			},
		},
		{
			name: "PT123,123H123.123S",
			wantErr: &UnexpectedTokenError{
				Value:      "PT123,123H123.123S",
				Token:      ".",
				AfterToken: "PT123,123H123",
				Expected:   "only the smallest time unit can be fractional",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration([]byte(tt.name))
			if tt.wantErr != nil {
				if diff := cmp.Diff(tt.wantErr, err); diff != "" {
					t.Errorf("error: (-want, +got)\n%s", diff)
				}
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
