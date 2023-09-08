package iso8601

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_ParseDate(t *testing.T) {
	tests := []struct {
		name    string
		want    DateLike
		wantErr error
	}{
		{
			name: "20121224",
			want: Date{
				Year:  2012,
				Month: time.December,
				Day:   24,
			},
		},
		{
			name: "00001224",
			want: Date{
				Year:  0,
				Month: time.December,
				Day:   24,
			},
		},
		{
			name: "2012359",
			want: OrdinalDate{
				Year: 2012,
				Day:  359,
			},
		},
		{
			name: "2012W521",
			want: WeekDate{
				Year: 2012,
				Week: 52,
				Day:  1,
			},
		},
		{
			name: "2012Q485",
			want: QuarterDate{
				Year:    2012,
				Quarter: 4,
				Day:     85,
			},
		},
		{
			name: "2012-12-24",
			want: Date{
				Year:  2012,
				Month: time.December,
				Day:   24,
			},
		},
		{
			name: "0000-12-24",
			want: Date{
				Year:  0,
				Month: time.December,
				Day:   24,
			},
		},
		{
			name: "2012-359",
			want: OrdinalDate{
				Year: 2012,
				Day:  359,
			},
		},
		{
			name: "0000-366",
			want: OrdinalDate{
				Year: 0,
				Day:  366,
			},
		},
		{
			name: "2012-W52-1",
			want: WeekDate{
				Year: 2012,
				Week: 52,
				Day:  1,
			},
		},
		{
			name: "0000-W52-1",
			want: WeekDate{
				Year: 0,
				Week: 52,
				Day:  1,
			},
		},
		{
			name: "2012-Q4-85",
			want: QuarterDate{
				Year:    2012,
				Quarter: 4,
				Day:     85,
			},
		},
		{
			name: "0000-Q4-85",
			want: QuarterDate{
				Year:    0,
				Quarter: 4,
				Day:     85,
			},
		},
		{
			name: "20",
			wantErr: &UnexpectedTokenError{
				Value: "20",
				Token: humanizeDigits(2),
			},
		},
		{
			name: "2000/",
			wantErr: &UnexpectedTokenError{
				Value:      "2000/",
				Token:      "/",
				AfterToken: "2000",
				Expected:   "8 or more characters",
			},
		},
		{
			name: "2000Q1",
			wantErr: &UnexpectedTokenError{
				Value:      "2000Q1",
				Token:      "Q1",
				AfterToken: "2000",
				Expected:   "8 or more characters",
			},
		},
		{
			name: "2000Q12",
			wantErr: &UnexpectedTokenError{
				Value:      "2000Q12",
				Token:      "Q12",
				AfterToken: "2000",
				Expected:   "8 or more characters",
			},
		},
		{
			name: "2000W1",
			wantErr: &UnexpectedTokenError{
				Value:      "2000W1",
				Token:      "W1",
				AfterToken: "2000",
				Expected:   "8 or more characters",
			},
		},
		{
			name: "2000W12",
			wantErr: &UnexpectedTokenError{
				Value:      "2000W12",
				Token:      "W12",
				AfterToken: "2000",
				Expected:   "8 or more characters",
			},
		},
		{
			name: "2000X1234",
			wantErr: &UnexpectedTokenError{
				Value:      "2000X1234",
				Token:      "X1234",
				AfterToken: "2000",
				Expected:   "- or Q or W",
			},
		},
		{
			name: "2000Q1234",
			wantErr: &UnexpectedTokenError{
				Value:      "2000Q1234",
				Token:      humanizeDigits(4),
				AfterToken: "Q",
				Expected:   humanizeDigits(3),
			},
		},
		{
			name: "2000W1234",
			wantErr: &UnexpectedTokenError{
				Value:      "2000W1234",
				Token:      humanizeDigits(4),
				AfterToken: "W",
				Expected:   humanizeDigits(3),
			},
		},
		{
			name: "2000-Q12-34",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-Q12-34",
				Token:      humanizeDigits(2),
				AfterToken: "Q",
				Expected:   humanizeDigits(1),
			},
		},
		{
			name: "2000-Q1=34",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-Q1=34",
				Token:      "=",
				AfterToken: "Q1",
				Expected:   "-",
			},
		},
		{
			name: "2000-Q1-123",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-Q1-123",
				Token:      humanizeDigits(3),
				AfterToken: "Q1-",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "2000-W123-34",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-W123-34",
				Token:      humanizeDigits(3),
				AfterToken: "W",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "2000-W1-234",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-W1-234",
				Token:      humanizeDigits(1),
				AfterToken: "W",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "2000-W12=34",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-W12=34",
				Token:      "=",
				AfterToken: "W12",
				Expected:   "-",
			},
		},
		{
			name: "2000-W12-34",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-W12-34",
				Token:      humanizeDigits(2),
				AfterToken: "W12-",
				Expected:   humanizeDigits(1),
			},
		},
		{
			name: "2000-12~34",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-12~34",
				Token:      "~",
				AfterToken: "-12",
				Expected:   "-",
			},
		},
		{
			name: "2000-12-345",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-12-345",
				Token:      humanizeDigits(3),
				AfterToken: "-12-",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "2000-12345",
			wantErr: &UnexpectedTokenError{
				Value:      "2000-12345",
				Token:      humanizeDigits(5),
				AfterToken: "2000-",
				Expected:   "like -Q4-85 or -W52-1 or -359",
			},
		},
		// valid format but range is invalid
		{
			name: "20121324",
			wantErr: &DateLikeRangeError{
				Element: "month",
				Value:   13,
				Year:    2012,
				Min:     1,
				Max:     12,
			},
		},
		{
			name: "20110229",
			wantErr: &DateLikeRangeError{
				Element: "day of month",
				Value:   29,
				Year:    2011,
				Min:     1,
				Max:     28,
			},
		},
		{
			name: "20120230",
			wantErr: &DateLikeRangeError{
				Element: "day of month",
				Value:   30,
				Year:    2012,
				Min:     1,
				Max:     29,
			},
		},
		{
			name: "2012367",
			wantErr: &DateLikeRangeError{
				Element: "day of year",
				Value:   367,
				Year:    2012,
				Min:     1,
				Max:     366,
			},
		},
		{
			name: "2013366",
			wantErr: &DateLikeRangeError{
				Element: "day of year",
				Value:   366,
				Year:    2013,
				Min:     1,
				Max:     365,
			},
		},
		{
			name: "2012W018",
			wantErr: &DateLikeRangeError{
				Element: "day of week",
				Value:   8,
				Year:    2012,
				Min:     1,
				Max:     7,
			},
		},
		{
			name: "2012W010",
			wantErr: &DateLikeRangeError{
				Element: "day of week",
				Value:   0,
				Year:    2012,
				Min:     1,
				Max:     7,
			},
		},
		{
			name: "2012W532",
			wantErr: &DateLikeRangeError{
				Element: "week",
				Value:   53,
				Year:    2012,
				Min:     1,
				Max:     52,
			},
		},
		{
			name: "2012Q585",
			wantErr: &DateLikeRangeError{
				Element: "quarter",
				Value:   5,
				Year:    2012,
				Min:     1,
				Max:     4,
			},
		},
		{
			name: "2012Q192",
			wantErr: &DateLikeRangeError{
				Element: "day of quarter",
				Value:   92,
				Year:    2012,
				Min:     1,
				Max:     91,
			},
		},
		{
			name: "2013Q191",
			wantErr: &DateLikeRangeError{
				Element: "day of quarter",
				Value:   91,
				Year:    2013,
				Min:     1,
				Max:     90,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate([]byte(tt.name))
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

func Test_countDigits(t *testing.T) {
	type args struct {
		b []byte
		i int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "start from position 0",
			args: args{
				b: []byte("20121224"),
				i: 0,
			},
			want: 8,
		},
		{
			name: "start from position 4",
			args: args{
				b: []byte("20121224"),
				i: 4,
			},
			want: 4,
		},
		{
			name: "stop at 4",
			args: args{
				b: []byte("2012T1224"),
				i: 0,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countDigits(tt.args.b, tt.args.i); got != tt.want {
				t.Errorf("countDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNumber(t *testing.T) {
	type args struct {
		b     []byte
		start int
		width int
	}
	tests := []struct {
		args args
		want int
	}{
		{
			args: args{
				b:     []byte("987654321"),
				start: 0,
				width: 9,
			},
			want: 987654321,
		},
		{
			args: args{
				b:     []byte("4321"),
				start: 0,
				width: 4,
			},
			want: 4321,
		},
		{
			args: args{
				b:     []byte("4321"),
				start: 1,
				width: 3,
			},
			want: 321,
		},
		{
			args: args{
				b:     []byte("4321"),
				start: 2,
				width: 2,
			},
			want: 21,
		},
		{
			args: args{
				b:     []byte("4321"),
				start: 3,
				width: 1,
			},
			want: 1,
		},
		{
			args: args{
				b:     []byte("4321"),
				start: 4,
				width: 1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%s start: %d, width: %d", tt.args.b, tt.args.start, tt.args.width)
		t.Run(name, func(t *testing.T) {
			if got := parseNumber(tt.args.b, tt.args.start, tt.args.width); got != tt.want {
				t.Errorf("parseNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
