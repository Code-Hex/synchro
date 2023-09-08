package iso8601

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseZone(t *testing.T) {
	tests := []struct {
		name    string
		want    Zone
		wantErr error
	}{
		{
			name: "Z",
			want: Zone{},
		},
		{
			name: "+00",
			want: Zone{},
		},
		{
			name: "+23",
			want: Zone{
				Hour: 23,
			},
		},
		{
			name: "-00",
			want: Zone{
				Minus: true,
			},
		},
		{
			name: "-23",
			want: Zone{
				Hour:  23,
				Minus: true,
			},
		},
		{
			name: "+2300",
			want: Zone{
				Hour:   23,
				Minute: 0,
			},
		},
		{
			name: "+0059",
			want: Zone{
				Hour:   0,
				Minute: 59,
			},
		},
		{
			name: "-2300",
			want: Zone{
				Hour:   23,
				Minute: 0,
				Minus:  true,
			},
		},
		{
			name: "-0059",
			want: Zone{
				Hour:   0,
				Minute: 59,
				Minus:  true,
			},
		},
		{
			name: "-0000",
			want: Zone{
				Hour:   0,
				Minute: 0,
				Minus:  true,
			},
		},
		{
			name: "+000000",
			want: Zone{
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
		},
		{
			name: "+235900",
			want: Zone{
				Hour:   23,
				Minute: 59,
				Second: 0,
			},
		},
		{
			name: "+235959",
			want: Zone{
				Hour:   23,
				Minute: 59,
				Second: 59,
			},
		},
		{
			name: "+000059",
			want: Zone{
				Hour:   0,
				Minute: 0,
				Second: 59,
			},
		},
		{
			name: "-000000",
			want: Zone{
				Hour:   0,
				Minute: 0,
				Second: 0,
				Minus:  true,
			},
		},
		{
			name: "-235900",
			want: Zone{
				Hour:   23,
				Minute: 59,
				Second: 0,
				Minus:  true,
			},
		},
		{
			name: "-235959",
			want: Zone{
				Hour:   23,
				Minute: 59,
				Second: 59,
				Minus:  true,
			},
		},
		{
			name: "-000059",
			want: Zone{
				Hour:   0,
				Minute: 0,
				Second: 59,
				Minus:  true,
			},
		},
		{
			name: "+00:00",
			want: Zone{},
		},
		{
			name: "+00:59",
			want: Zone{
				Minute: 59,
			},
		},
		{
			name: "+23:59",
			want: Zone{
				Hour:   23,
				Minute: 59,
			},
		},
		{
			name: "+23:00",
			want: Zone{
				Hour: 23,
			},
		},
		{
			name: "-00:00",
			want: Zone{
				Minus: true,
			},
		},
		{
			name: "-00:59",
			want: Zone{
				Minute: 59,
				Minus:  true,
			},
		},
		{
			name: "-23:59",
			want: Zone{
				Hour:   23,
				Minute: 59,
				Minus:  true,
			},
		},
		{
			name: "-23:00",
			want: Zone{
				Hour:  23,
				Minus: true,
			},
		},
		{
			name: "+00:00:00",
			want: Zone{},
		},
		{
			name: "+23:00:00",
			want: Zone{
				Hour: 23,
			},
		},
		{
			name: "+23:59:00",
			want: Zone{
				Hour:   23,
				Minute: 59,
			},
		},
		{
			name: "+23:59:59",
			want: Zone{
				Hour:   23,
				Minute: 59,
				Second: 59,
			},
		},
		{
			name: "+00:59:59",
			want: Zone{
				Hour:   0,
				Minute: 59,
				Second: 59,
			},
		},
		{
			name: "+00:00:59",
			want: Zone{
				Hour:   0,
				Minute: 0,
				Second: 59,
			},
		},
		{
			name: "-00:00:00",
			want: Zone{
				Minus: true,
			},
		},
		{
			name: "-23:00:00",
			want: Zone{
				Hour:  23,
				Minus: true,
			},
		},
		{
			name: "-23:59:00",
			want: Zone{
				Hour:   23,
				Minute: 59,
				Minus:  true,
			},
		},
		{
			name: "-23:59:59",
			want: Zone{
				Hour:   23,
				Minute: 59,
				Second: 59,
				Minus:  true,
			},
		},
		{
			name: "-00:59:59",
			want: Zone{
				Hour:   0,
				Minute: 59,
				Second: 59,
				Minus:  true,
			},
		},
		{
			name: "-00:00:59",
			want: Zone{
				Hour:   0,
				Minute: 0,
				Second: 59,
				Minus:  true,
			},
		},
		{
			name: "",
			wantErr: &UnexpectedTokenError{
				Value:      "",
				Token:      "",
				AfterToken: "",
				Expected:   "Z or + or -",
			},
		},
		{
			name: "Z12",
			wantErr: &UnexpectedTokenError{
				Value:      "Z12",
				Token:      "12",
				AfterToken: "Z",
				Expected:   "non extra token (12)",
			},
		},
		{
			name: "X",
			wantErr: &UnexpectedTokenError{
				Value:      "X",
				Token:      "X",
				AfterToken: "",
				Expected:   "Z or + or -",
			},
		},
		{
			name: "+000",
			wantErr: &UnexpectedTokenError{
				Value:      "+000",
				Token:      humanizeDigits(3),
				AfterToken: "+",
				Expected: fmt.Sprintf(
					"%s or %s or %s",
					humanizeDigits(2),
					humanizeDigits(4),
					humanizeDigits(6),
				),
			},
		},
		{
			name: "-00000",
			wantErr: &UnexpectedTokenError{
				Value:      "-00000",
				Token:      humanizeDigits(5),
				AfterToken: "-",
				Expected: fmt.Sprintf(
					"%s or %s or %s",
					humanizeDigits(2),
					humanizeDigits(4),
					humanizeDigits(6),
				),
			},
		},
		{
			name: "+00:000",
			wantErr: &UnexpectedTokenError{
				Value:      "+00:000",
				Token:      humanizeDigits(3),
				AfterToken: "+00:",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "-00:00:000",
			wantErr: &UnexpectedTokenError{
				Value:      "-00:00:000",
				Token:      humanizeDigits(3),
				AfterToken: "-00:00:",
				Expected:   humanizeDigits(2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseZone([]byte(tt.name))
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

func Test_parseExtendedZone(t *testing.T) {
	tests := []struct {
		name    string
		want    Zone
		wantErr error
	}{
		{
			name: "Z",
			want: Zone{},
		},
		{
			name: "+99",
			want: Zone{
				Hour: 99,
			},
		},
		{
			name: "",
			wantErr: &UnexpectedTokenError{
				Value:      "",
				Token:      "",
				AfterToken: "",
				Expected:   "Z or + or -",
			},
		},
		{
			name: "Z12",
			wantErr: &UnexpectedTokenError{
				Value:      "Z12",
				Token:      "12",
				AfterToken: "Z",
				Expected:   "non extra token (12)",
			},
		},
		{
			name: "X",
			wantErr: &UnexpectedTokenError{
				Value:      "X",
				Token:      "X",
				AfterToken: "",
				Expected:   "Z or + or -",
			},
		},
		{
			name: "-123",
			wantErr: &UnexpectedTokenError{
				Value:    "-123",
				Token:    humanizeDigits(3),
				Expected: humanizeDigits(2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExtendedZone([]byte(tt.name))
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

func Test_timeZone(t *testing.T) {
	got, err := timeZone(1000, 0, 0, false)
	if err == nil {
		t.Fatal("unexpected err is nil")
	}
	want := Zone{}
	if want != got {
		t.Errorf("got zone %v", got)
	}
}
