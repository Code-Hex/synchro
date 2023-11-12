package iso8601

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name    string
		want    Time
		wantErr error
	}{
		{
			name: "01",
			want: Time{
				Hour: 1,
			},
		},
		{
			name: "23",
			want: Time{
				Hour: 23,
			},
		},
		{
			name: "24",
			want: Time{
				Hour: 24,
			},
		},
		{
			name: "2301",
			want: Time{
				Hour:   23,
				Minute: 1,
			},
		},
		{
			name: "2312",
			want: Time{
				Hour:   23,
				Minute: 12,
			},
		},
		{
			name: "2300",
			want: Time{
				Hour:   23,
				Minute: 0,
			},
		},
		{
			name: "235959",
			want: Time{
				Hour:   23,
				Minute: 59,
				Second: 59,
			},
		},
		{
			name: "235900",
			want: Time{
				Hour:   23,
				Minute: 59,
				Second: 0,
			},
		},
		{
			name: "230010",
			want: Time{
				Hour:   23,
				Minute: 0,
				Second: 10,
			},
		},
		{
			name: "230010.",
			want: Time{
				Hour:   23,
				Minute: 0,
				Second: 10,
			},
		},
		{
			name: "230010.1",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: int(100 * time.Millisecond),
			},
		},
		{
			name: "230010.01",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: int(10 * time.Millisecond),
			},
		},
		{
			name: "230010.987654321",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: 987654321,
			},
		},
		{
			name: "230010.98765432101",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: 987654321,
			},
		},
		{
			name: "230010,",
			want: Time{
				Hour:   23,
				Minute: 0,
				Second: 10,
			},
		},
		{
			name: "230010,1",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: int(100 * time.Millisecond),
			},
		},
		{
			name: "230010,01",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: int(10 * time.Millisecond),
			},
		},
		{
			name: "230010,987654321",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: 987654321,
			},
		},
		{
			name: "230010,98765432101",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: 987654321,
			},
		},
		{
			name: "1430.5",
			want: Time{
				Hour:       14,
				Minute:     30,
				Second:     30,
				Nanosecond: 0,
			},
		},
		{
			name: "1430,5",
			want: Time{
				Hour:       14,
				Minute:     30,
				Second:     30,
				Nanosecond: 0,
			},
		},
		{
			name: "14.5",
			want: Time{
				Hour:       14,
				Minute:     30,
				Second:     00,
				Nanosecond: 0,
			},
		},
		{
			name: "14,5",
			want: Time{
				Hour:       14,
				Minute:     30,
				Second:     00,
				Nanosecond: 0,
			},
		},
		{
			name: "24:00",
			want: Time{
				Hour: 24,
			},
		},
		{
			name: "23:01",
			want: Time{
				Hour:   23,
				Minute: 1,
			},
		},
		{
			name: "23:59",
			want: Time{
				Hour:   23,
				Minute: 59,
			},
		},
		{
			name: "23:59:00",
			want: Time{
				Hour:   23,
				Minute: 59,
				Second: 0,
			},
		},
		{
			name: "23:59:59",
			want: Time{
				Hour:   23,
				Minute: 59,
				Second: 59,
			},
		},
		{
			name: "23:00:59",
			want: Time{
				Hour:   23,
				Minute: 0,
				Second: 59,
			},
		},
		{
			name: "23:00:10.",
			want: Time{
				Hour:   23,
				Minute: 0,
				Second: 10,
			},
		},
		{
			name: "23:00:10.1",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: int(100 * time.Millisecond),
			},
		},
		{
			name: "23:00:10.01",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: int(10 * time.Millisecond),
			},
		},
		{
			name: "23:00:10.987654321",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: 987654321,
			},
		},
		{
			name: "23:00:10.987654321012",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: 987654321,
			},
		},
		{
			name: "23:00:10,",
			want: Time{
				Hour:   23,
				Minute: 0,
				Second: 10,
			},
		},
		{
			name: "23:00:10,1",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: int(100 * time.Millisecond),
			},
		},
		{
			name: "23:00:10,01",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: int(10 * time.Millisecond),
			},
		},
		{
			name: "23:00:10,987654321",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: 987654321,
			},
		},
		{
			name: "23:00:10,987654321012",
			want: Time{
				Hour:       23,
				Minute:     0,
				Second:     10,
				Nanosecond: 987654321,
			},
		},
		{
			name: "14:30.5",
			want: Time{
				Hour:       14,
				Minute:     30,
				Second:     30,
				Nanosecond: 0,
			},
		},
		{
			name: "14:30,5",
			want: Time{
				Hour:       14,
				Minute:     30,
				Second:     30,
				Nanosecond: 0,
			},
		},
		{
			name: "0",
			wantErr: &UnexpectedTokenError{
				Value:    "0",
				Token:    humanizeDigits(1),
				Expected: humanizeDigits(2),
			},
		},
		{
			name: "011",
			wantErr: &UnexpectedTokenError{
				Value:      "011",
				Token:      humanizeDigits(3),
				AfterToken: "01",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "01234",
			wantErr: &UnexpectedTokenError{
				Value:      "01234",
				Token:      humanizeDigits(5),
				AfterToken: "0123",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "0123451",
			wantErr: &UnexpectedTokenError{
				Value:      "0123451",
				Token:      humanizeDigits(7),
				AfterToken: "0123",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "01:123",
			wantErr: &UnexpectedTokenError{
				Value:      "01:123",
				Token:      humanizeDigits(3),
				AfterToken: "01:",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "01:12:345",
			wantErr: &UnexpectedTokenError{
				Value:      "01:12:345",
				Token:      humanizeDigits(3),
				AfterToken: "01:12:",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "235959hello",
			wantErr: &UnexpectedTokenError{
				Value:      "235959hello",
				Token:      "hello",
				AfterToken: "235959",
				Expected:   "235959",
			},
		},
		{
			name: "23:59:59hello",
			wantErr: &UnexpectedTokenError{
				Value:      "23:59:59hello",
				Token:      "hello",
				AfterToken: "23:59:59",
				Expected:   "23:59:59",
			},
		},
		// invalid time range
		{
			name: "2401",
			wantErr: &TimeRangeError{
				Element: "hour",
				Value:   24,
				Min:     0,
				Max:     24,
			},
		},
		{
			name: "2360",
			wantErr: &TimeRangeError{
				Element: "minute",
				Value:   60,
				Min:     0,
				Max:     59,
			},
		},
		{
			name: "235960",
			wantErr: &TimeRangeError{
				Element: "second",
				Value:   60,
				Min:     0,
				Max:     59,
			},
		},
		{
			name: "25",
			wantErr: &TimeRangeError{
				Element: "hour",
				Value:   25,
				Min:     0,
				Max:     24,
			},
		},
		{
			name: "24:01",
			wantErr: &TimeRangeError{
				Element: "hour",
				Value:   24,
				Min:     0,
				Max:     24,
			},
		},
		{
			name: "23:60",
			wantErr: &TimeRangeError{
				Element: "minute",
				Value:   60,
				Min:     0,
				Max:     59,
			},
		},
		{
			name: "23:59:60",
			wantErr: &TimeRangeError{
				Element: "second",
				Value:   60,
				Min:     0,
				Max:     59,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTime([]byte(tt.name))
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

func Test_parseExtendedTime(t *testing.T) {
	tests := []struct {
		name    string
		want    Time
		wantErr error
	}{
		{
			name: "24",
			want: Time{Hour: 24},
		},
		{
			name: "23",
			want: Time{Hour: 23},
		},
		{
			name: "01",
			want: Time{Hour: 1},
		},
		{
			name: "00",
			want: Time{},
		},
		{
			name: "0",
			wantErr: &UnexpectedTokenError{
				Value:    "0",
				Token:    humanizeDigits(1),
				Expected: humanizeDigits(2),
			},
		},
		{
			name: "123",
			wantErr: &UnexpectedTokenError{
				Value:    "123",
				Token:    humanizeDigits(3),
				Expected: humanizeDigits(2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := parseExtendedTime([]byte(tt.name))
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

func TestTimeRangeError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *TimeRangeError
		want string
	}{
		{
			name: "valid error",
			err: &TimeRangeError{
				Element: "hour",
				Value:   25,
				Min:     0,
				Max:     24,
			},
			want: "iso8601 time: 25 hour is not in range 0-24",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("TimeRangeError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_String(t *testing.T) {
	tests := []struct {
		t    Time
		want string
	}{
		{
			t: Time{
				Hour:   12,
				Minute: 59,
				Second: 59,
			},
			want: "12:59:59",
		},
		{
			t: Time{
				Hour:       12,
				Minute:     59,
				Second:     59,
				Nanosecond: 987654321,
			},
			want: "12:59:59.987654321",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("Time.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_IsZero(t *testing.T) {
	tests := []struct {
		name string
		t    Time
		want bool
	}{
		{
			name: "zero value",
			t:    Time{},
			want: true,
		},
		{
			name: "non-zero value",
			t: Time{
				Hour: 12,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.IsZero(); got != tt.want {
				t.Errorf("Time.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_MarshalText(t *testing.T) {
	tests := []struct {
		t    Time
		want string
	}{
		{
			t: Time{
				Hour:   12,
				Minute: 59,
				Second: 59,
			},
			want: "12:59:59",
		},
		{
			t: Time{
				Hour:       12,
				Minute:     59,
				Second:     59,
				Nanosecond: 987654321,
			},
			want: "12:59:59.987654321",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got, err := tt.t.MarshalText()
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Errorf("Time.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_Validate(t *testing.T) {
	tests := []struct {
		name    string
		time    Time
		wantErr bool
	}{
		{
			name: "valid time",
			time: Time{
				Hour:       12,
				Minute:     30,
				Second:     15,
				Nanosecond: 500000000,
			},
			wantErr: false,
		},
		{
			name: "invalid hour",
			time: Time{
				Hour:       25,
				Minute:     30,
				Second:     15,
				Nanosecond: 500000000,
			},
			wantErr: true,
		},
		{
			name: "invalid minute",
			time: Time{
				Hour:       12,
				Minute:     60,
				Second:     15,
				Nanosecond: 500000000,
			},
			wantErr: true,
		},
		{
			name: "invalid second",
			time: Time{
				Hour:       12,
				Minute:     30,
				Second:     60,
				Nanosecond: 500000000,
			},
			wantErr: true,
		},
		{
			name: "invalid nanosecond",
			time: Time{
				Hour:       12,
				Minute:     30,
				Second:     15,
				Nanosecond: 1000000000,
			},
			wantErr: true,
		},
		{
			name: "valid edge case - 24:00:00",
			time: Time{
				Hour:       24,
				Minute:     0,
				Second:     0,
				Nanosecond: 0,
			},
			wantErr: false,
		},
		{
			name: "invalid edge case - 24:00:01",
			time: Time{
				Hour:       24,
				Minute:     0,
				Second:     1,
				Nanosecond: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.time.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTime_UnmarshalText(t *testing.T) {
	tests := []struct {
		s    string
		want Time
	}{
		{
			s: "12:59:59",
			want: Time{
				Hour:   12,
				Minute: 59,
				Second: 59,
			},
		},
		{
			s: "12:59:59.987654321",
			want: Time{
				Hour:       12,
				Minute:     59,
				Second:     59,
				Nanosecond: 987654321,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			tm := &Time{}
			err := tm.UnmarshalText([]byte(tt.s))
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, *tm); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}

	t.Run("invalid text", func(t *testing.T) {
		tm := &Time{}
		err := tm.UnmarshalText([]byte("invalid"))
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTimeJSON(t *testing.T) {
	tests := []struct {
		name string
		time Time
	}{
		{
			name: "valid time",
			time: Time{
				Hour:       12,
				Minute:     30,
				Second:     15,
				Nanosecond: 500000000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.time)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			var got Time
			err = json.Unmarshal(bytes, &got)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			if diff := cmp.Diff(got, tt.time); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
