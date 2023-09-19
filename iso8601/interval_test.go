package iso8601

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseInterval(t *testing.T) {
	tests := []struct {
		name    string
		want    Interval
		wantErr error
	}{
		{
			name: "2007-03-01T13:00:00Z/2008-05-11T15:30:00Z",
			want: Interval{
				start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				end:   time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "2007-03-01T13:00:00Z/P1Y2M10DT2H30M",
			want: Interval{
				start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
			},
		},
		{
			name: "P1Y2M10DT2H30M/2008-05-11T15:30:00Z",
			want: Interval{
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				end: time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "P1Y2M10DT2H30M",
			want: Interval{
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
			},
		},
		{
			name: "2007-03-01T13:00:00Z--2008-05-11T15:30:00Z",
			want: Interval{
				start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				end:   time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "2007-03-01T13:00:00Z--P1Y2M10DT2H30M",
			want: Interval{
				start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
			},
		},
		{
			name: "P1Y2M10DT2H30M--2008-05-11T15:30:00Z",
			want: Interval{
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				end: time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "R3/2007-03-01T13:00:00Z/2008-05-11T15:30:00Z",
			want: Interval{
				start:  time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				end:    time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
				repeat: 3,
			},
		},
		{
			name: "R12/2007-03-01T13:00:00Z/P1Y2M10DT2H30M",
			want: Interval{
				start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				repeat: 12,
			},
		},
		{
			name: "R/P1Y2M10DT2H30M/2008-05-11T15:30:00Z",
			want: Interval{
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				end:    time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
				repeat: -1,
			},
		},
		{
			name: "R1234/P1Y2M10DT2H30M",
			want: Interval{
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				repeat: 1234,
			},
		},
		{
			name: "R--2007-03-01T13:00:00Z--2008-05-11T15:30:00Z",
			want: Interval{
				start:  time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				end:    time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
				repeat: -1,
			},
		},
		{
			name: "R123--2007-03-01T13:00:00Z--P1Y2M10DT2H30M",
			want: Interval{
				start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				repeat: 123,
			},
		},
		{
			name: "R0--P1Y2M10DT2H30M--2008-05-11T15:30:00Z",
			want: Interval{
				duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				end:    time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
				repeat: 0,
			},
		},
		{
			name: "",
			wantErr: &UnexpectedTokenError{
				Value:      "",
				Token:      "",
				AfterToken: "",
				Expected:   "R or P or datetime",
			},
		},
		{
			name: "R",
			wantErr: &UnexpectedTokenError{
				Value:      "R",
				Token:      "",
				AfterToken: "R",
				Expected:   `internal designator "/" or "--"`,
			},
		},
		{
			name: "Rhello",
			wantErr: &UnexpectedTokenError{
				Value:      "Rhello",
				Token:      "hello",
				AfterToken: "R",
				Expected:   `internal designator "/" or "--"`,
			},
		},
		{
			name: "Rhello/",
			wantErr: &UnexpectedTokenError{
				Value:      "Rhello/",
				Token:      "hello",
				AfterToken: "R",
				Expected:   humanizeDigits(5),
			},
		},
		{
			name: "Rhello--",
			wantErr: &UnexpectedTokenError{
				Value:      "Rhello--",
				Token:      "hello",
				AfterToken: "R",
				Expected:   humanizeDigits(5),
			},
		},
		{
			name: "R12llo/",
			wantErr: &UnexpectedTokenError{
				Value:      "R12llo/",
				Token:      "12llo",
				AfterToken: "R",
				Expected:   humanizeDigits(5),
			},
		},
		{
			name: "1234",
			wantErr: &UnexpectedTokenError{
				Value:      "1234",
				Token:      "1234",
				AfterToken: "",
				Expected:   "P or + or -",
			},
		},
		{
			name: "R/1234",
			wantErr: &UnexpectedTokenError{
				Value:      "1234",
				Token:      "1234",
				AfterToken: "",
				Expected:   "P or + or -",
			},
		},
		{
			name: "P/",
			wantErr: &UnexpectedTokenError{
				Value:      "P",
				Token:      "",
				AfterToken: "P",
				Expected:   "PnYnMnDTnHnMnS or PnW format",
			},
		},
		{
			name: "R/P/",
			wantErr: &UnexpectedTokenError{
				Value:      "P",
				Token:      "",
				AfterToken: "P",
				Expected:   "PnYnMnDTnHnMnS or PnW format",
			},
		},
		{
			name: "R/P",
			wantErr: &UnexpectedTokenError{
				Value:      "P",
				Token:      "",
				AfterToken: "P",
				Expected:   "PnYnMnDTnHnMnS or PnW format",
			},
		},
		{
			name: "unknown/",
			wantErr: &UnexpectedTokenError{
				Value:      "unknown",
				Token:      humanizeDigits(0),
				AfterToken: "",
				Expected:   "date format",
			},
		},
		{
			name: "R/unknown/",
			wantErr: &UnexpectedTokenError{
				Value:      "unknown",
				Token:      humanizeDigits(0),
				AfterToken: "",
				Expected:   "date format",
			},
		},
		{
			name: "P1Y2M10DT2H30M/P1Y2M10DT2H30M",
			wantErr: &UnexpectedTokenError{
				Value:      "P1Y2M10DT2H30M/P1Y2M10DT2H30M",
				Token:      "P1Y2M10DT2H30M",
				AfterToken: "/",
				Expected:   "datetime format",
			},
		},
		{
			name: "P1Y2M10DT2H30M--P1Y2M10DT2H30M",
			wantErr: &UnexpectedTokenError{
				Value:      "P1Y2M10DT2H30M--P1Y2M10DT2H30M",
				Token:      "P1Y2M10DT2H30M",
				AfterToken: "--",
				Expected:   "datetime format",
			},
		},
		{
			name: "R/2007-03-01T13:00:00Z/Phello",
			wantErr: &UnexpectedTokenError{
				Value:      "Phello",
				Token:      "h",
				AfterToken: "P",
				Expected:   "PnYnMnDTnHnMnS or PnW format",
			},
		},
		{
			name: "P1Y2M10DT2H30M/hello",
			wantErr: &UnexpectedTokenError{
				Value:      "hello",
				Token:      humanizeDigits(0),
				AfterToken: "",
				Expected:   "date format",
			},
		},
		// different interval designators
		{
			name: "R--2007-03-01T13:00:00Z/2008-05-11T15:30:00Z",
			wantErr: &UnexpectedTokenError{
				Value:      "2007-03-01T13:00:00Z/2008-05-11T15:30:00Z",
				Token:      "2007-03-01T13:00:00Z/2008-05-11T15:30:00Z",
				AfterToken: "",
				Expected:   "P or + or -",
			},
		},
		{
			name: "R1/2007-03-01T13:00:00Z--2008-05-11T15:30:00Z",
			wantErr: &UnexpectedTokenError{
				Value:      "2007-03-01T13:00:00Z--2008-05-11T15:30:00Z",
				Token:      "2007-03-01T13:00:00Z--2008-05-11T15:30:00Z",
				AfterToken: "",
				Expected:   "P or + or -",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInterval([]byte(tt.name))
			if tt.wantErr != nil {
				if diff := cmp.Diff(tt.wantErr, err); diff != "" {
					t.Errorf("error: (-want, +got)\n%s", diff)
				}
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(Interval{})); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestInterval_Contains(t *testing.T) {
	tests := []struct {
		name  string
		parse string
		t     time.Time
		want  bool
	}{
		{
			name:  "2013-01-01/2013-01-10 contains 2013-01-01",
			parse: "2013-01-01/2013-01-10",
			t:     time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC),
			want:  true,
		},
		{
			name:  "2013-01-01/2013-01-10 contains 2013-01-10",
			parse: "2013-01-01/2013-01-10",
			t:     time.Date(2013, 1, 10, 0, 0, 0, 0, time.UTC),
			want:  true,
		},
		{
			name:  "2013-01-01/2013-01-10 does not contain 2013-01-11",
			parse: "2013-01-01/2013-01-10",
			t:     time.Date(2013, 1, 11, 0, 0, 0, 0, time.UTC),
			want:  false,
		},
		{
			name:  "2013-01-01/2013-01-10 does not contain 2012-12-31",
			parse: "2013-01-01/2013-01-10",
			t:     time.Date(2012, 12, 31, 0, 0, 0, 0, time.UTC),
			want:  false,
		},
		{
			name:  "2013-01-01/2013-01-10 contains 2013-01-09T23:59:59",
			parse: "2013-01-01/2013-01-10",
			t:     time.Date(2013, 1, 9, 23, 59, 59, 0, time.UTC),
			want:  true,
		},
		{
			name:  "2013-01-01/2013-01-10 does not contain 2013-01-10T00:00:01",
			parse: "2013-01-01/2013-01-10",
			t:     time.Date(2013, 1, 10, 0, 0, 1, 0, time.UTC),
			want:  false,
		},
		{
			name:  "2013-01-01/2013-01-10T23:59:59Z contains 2013-01-10T23:59:59",
			parse: "2013-01-01/2013-01-10T23:59:59Z",
			t:     time.Date(2013, 1, 10, 23, 59, 59, 0, time.UTC),
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := ParseInterval(tt.parse)
			if err != nil {
				t.Fatal(err)
			}
			if got := i.Contains(tt.t); got != tt.want {
				t.Errorf("Interval.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_Start(t *testing.T) {
	testCases := []struct {
		name     string
		interval Interval
		expected time.Time
	}{
		{
			name: "Start time is set",
			interval: Interval{
				start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "End time and duration are set",
			interval: Interval{
				end:      time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
				duration: Duration{Day: 1},
			},
			expected: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Neither start time, end time, nor duration are set",
			interval: Interval{},
			expected: time.Time{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.interval.Start(); !got.Equal(tc.expected) {
				t.Errorf("Expected start time %v, but got %v", tc.expected, got)
			}
		})
	}
}

func TestInterval_End(t *testing.T) {
	testCases := []struct {
		name     string
		interval Interval
		expected time.Time
	}{
		{
			name: "End time is set",
			interval: Interval{
				end: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Start time and duration are set",
			interval: Interval{
				start:    time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				duration: Duration{Day: 1},
			},
			expected: time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Neither start time, end time, nor duration are set",
			interval: Interval{},
			expected: time.Time{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.interval.End(); !got.Equal(tc.expected) {
				t.Errorf("Expected end time %v, but got %v", tc.expected, got)
			}
		})
	}
}

func TestInterval_Duration(t *testing.T) {
	testCases := []struct {
		name     string
		interval Interval
		expected Duration
	}{
		{
			name: "Duration is set",
			interval: Interval{
				duration: Duration{Day: 1},
			},
			expected: Duration{Day: 1},
		},
		{
			name: "Start time and end time are set",
			interval: Interval{
				start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				end:   time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
			},
			expected: Duration{Day: 1},
		},
		{
			name:     "Neither duration, start time, nor end time are set",
			interval: Interval{},
			expected: Duration{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.interval.Duration(); got != tc.expected {
				t.Errorf("Expected duration %v, but got %v", tc.expected, got)
			}
		})
	}
}
