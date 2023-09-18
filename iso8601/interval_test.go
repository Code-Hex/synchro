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
				Start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				End:   time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "2007-03-01T13:00:00Z/P1Y2M10DT2H30M",
			want: Interval{
				Start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				Duration: Duration{
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
				Duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				End: time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "P1Y2M10DT2H30M",
			want: Interval{
				Duration: Duration{
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
				Start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				End:   time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "2007-03-01T13:00:00Z--P1Y2M10DT2H30M",
			want: Interval{
				Start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				Duration: Duration{
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
				Duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				End: time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "R3/2007-03-01T13:00:00Z/2008-05-11T15:30:00Z",
			want: Interval{
				Start:  time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				End:    time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
				Repeat: 3,
			},
		},
		{
			name: "R12/2007-03-01T13:00:00Z/P1Y2M10DT2H30M",
			want: Interval{
				Start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				Duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				Repeat: 12,
			},
		},
		{
			name: "R/P1Y2M10DT2H30M/2008-05-11T15:30:00Z",
			want: Interval{
				Duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				End:    time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
				Repeat: -1,
			},
		},
		{
			name: "R1234/P1Y2M10DT2H30M",
			want: Interval{
				Duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				Repeat: 1234,
			},
		},
		{
			name: "R--2007-03-01T13:00:00Z--2008-05-11T15:30:00Z",
			want: Interval{
				Start:  time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				End:    time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
				Repeat: -1,
			},
		},
		{
			name: "R123--2007-03-01T13:00:00Z--P1Y2M10DT2H30M",
			want: Interval{
				Start: time.Date(2007, 3, 1, 13, 0, 0, 0, time.UTC),
				Duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				Repeat: 123,
			},
		},
		{
			name: "R0--P1Y2M10DT2H30M--2008-05-11T15:30:00Z",
			want: Interval{
				Duration: Duration{
					Year:   1,
					Month:  2,
					Day:    10,
					Hour:   2,
					Minute: 30,
				},
				End:    time.Date(2008, 5, 11, 15, 30, 0, 0, time.UTC),
				Repeat: 0,
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
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
			if err != nil {
				t.Error(err)
			}
		})
	}
}
