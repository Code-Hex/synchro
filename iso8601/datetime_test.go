package iso8601

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseDateTime(t *testing.T) {
	const hour = 3600
	const minute = 60
	tests := []struct {
		name    string
		want    time.Time
		wantErr error
	}{
		{
			name: "2017-04-24T09:41:34.502+0100",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.FixedZone("", hour)),
		},
		{
			name: "2017-04-24T09:41+0100",
			want: time.Date(2017, 4, 24, 9, 41, 0, 0, time.FixedZone("", hour)),
		},
		{
			name: "2017-04-24T09+0100",
			want: time.Date(2017, 4, 24, 9, 0, 0, 0, time.FixedZone("", hour)),
		},
		{
			name: "2017-04-24",
			want: time.Date(2017, 4, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "2017-04-24T09:41:34+0100",
			want: time.Date(2017, 4, 24, 9, 41, 34, 0, time.FixedZone("", hour)),
		},
		{
			name: "2017-04-24T09:41:34.502-0100",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.FixedZone("", -1*hour)),
		},
		{
			name: "2017-04-24T09:41:34.502-01:00",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.FixedZone("", -1*hour)),
		},
		{
			name: "2017-04-24T09:41-01:00",
			want: time.Date(2017, 4, 24, 9, 41, 0, 0, time.FixedZone("", -1*hour)),
		},
		{
			name: "2017-04-24T09-01:00",
			want: time.Date(2017, 4, 24, 9, 0, 0, 0, time.FixedZone("", -1*hour)),
		},
		{
			name: "2017-04-24T09:41:34-0100",
			want: time.Date(2017, 4, 24, 9, 41, 34, 0, time.FixedZone("", -1*hour)),
		},
		{
			name: "2017-04-24T09:41:34.502Z",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.UTC),
		},
		{
			name: "2017-04-24T09:41:34Z",
			want: time.Date(2017, 4, 24, 9, 41, 34, 0, time.UTC),
		},
		{
			name: "2017-04-24T09:41Z",
			want: time.Date(2017, 4, 24, 9, 41, 0, 0, time.UTC),
		},
		{
			name: "2017-04-24T09Z",
			want: time.Date(2017, 4, 24, 9, 0, 0, 0, time.UTC),
		},
		{
			name: "2017-04-24T09:41:34.089",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(89*time.Millisecond), time.UTC),
		},
		{
			name: "2017-04-24T09:41",
			want: time.Date(2017, 4, 24, 9, 41, 0, 0, time.UTC),
		},
		{
			name: "2017-04-24T09",
			want: time.Date(2017, 4, 24, 9, 0, 0, 0, time.UTC),
		},
		{
			name: "2017-04-24T09:41:34.009",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(9*time.Millisecond), time.UTC),
		},
		{
			name: "2017-04-24T09:41:34.893",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(893*time.Millisecond), time.UTC),
		},
		{
			name: "2017-04-24T09:41:34,893",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(893*time.Millisecond), time.UTC),
		},
		{
			name: "2017-04-24T09:41:34.89312523Z",
			want: time.Date(2017, 4, 24, 9, 41, 34, 89312523*10, time.UTC),
		},
		{
			name: "2017-04-24T09:41:34.502-0530",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.FixedZone("", -1*(5*hour+30*minute))),
		},
		{
			name: "2017-04-24T09:41:34.502+0530",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.FixedZone("", 5*hour+30*minute)),
		},
		{
			name: "2017-04-24T09:41:34.502+05:30",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.FixedZone("", 5*hour+30*minute)),
		},
		{
			name: "2017-04-24T09:41:34.502+05:45",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.FixedZone("", 5*hour+45*minute)),
		},
		{
			name: "2017-04-24T09:41:34.502+00",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.UTC),
		},
		{
			name: "2017-04-24T09:41:34.502+00",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.UTC),
		},
		{
			name: "2017-04-24T09:41:34.502+0000",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.UTC),
		},
		{
			name: "2017-04-24T09:41:34.502+00:00",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.UTC),
		},
		{
			name: "+2017-04-24T09:41:34.502+00:00",
			want: time.Date(2017, 4, 24, 9, 41, 34, int(502*time.Millisecond), time.UTC),
		},
		{
			name: "2008-W09-4T09:41:34.502+00:00",
			want: time.Date(2008, 2, 28, 9, 41, 34, int(502*time.Millisecond), time.UTC),
		},
		{
			name: "2014-Q1-59T09:41:34.502+00:00",
			want: time.Date(2014, 2, 28, 9, 41, 34, int(502*time.Millisecond), time.UTC),
		},
		{
			name: "2016-060T09:41:34.502+00:00",
			want: time.Date(2016, 2, 29, 9, 41, 34, int(502*time.Millisecond), time.UTC),
		},
		{
			name: "",
			wantErr: &UnexpectedTokenError{
				Value:      "",
				Token:      humanizeDigits(0),
				AfterToken: "",
				Expected:   "date format",
			},
		},
		{
			name: "2017-04-24T",
			wantErr: &UnexpectedTokenError{
				Value:      "2017-04-24T",
				Token:      "",
				AfterToken: "2017-04-24T",
				Expected:   "time format is required after the 'T' designator",
			},
		},
		{
			name: "2017-04-24X",
			wantErr: &UnexpectedTokenError{
				Value:      "2017-04-24X",
				Token:      "X",
				AfterToken: "2017-04-24",
				Expected:   "'T'",
			},
		},
		{
			name: "2012-W52-1T09:41:",
			wantErr: &UnexpectedTokenError{
				Value:      "2012-W52-1T09:41:",
				Token:      humanizeDigits(0),
				AfterToken: "09:41:",
				Expected:   humanizeDigits(2),
			},
		},
		{
			name: "+2017-04-24T09:41:34.502X",
			wantErr: &UnexpectedTokenError{
				Value:      "+2017-04-24T09:41:34.502X",
				Token:      "X",
				AfterToken: "+2017-04-24T09:41:34.502",
				Expected:   "time zone format after time format",
			},
		},
		{
			name: "2012-W52-1T09:41:00Z12",
			wantErr: &UnexpectedTokenError{
				Value:      "2012-W52-1T09:41:00Z12",
				Token:      "12",
				AfterToken: "Z",
				Expected:   "non extra token (12)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateTime([]byte(tt.name))
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

	t.Run("WithTimeDesignators", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			want := time.Date(2017, 4, 24, 9, 41, 34, 89312523*10, time.UTC)
			got, err := ParseDateTime("2017-04-24 09:41:34.89312523Z", WithTimeDesignators(' '))
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})

		t.Run("invalid", func(t *testing.T) {
			wantErr := &UnexpectedTokenError{
				Value:      "2017-04-24X09:41:34",
				Token:      "X",
				AfterToken: "2017-04-24",
				Expected:   `'T', ' '`,
			}
			_, err := ParseDateTime("2017-04-24X09:41:34", WithTimeDesignators(' '))
			if err == nil {
				t.Fatal("expected error")
			}
			if diff := cmp.Diff(wantErr, err); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	})

	t.Run("WithInLocation", func(t *testing.T) {
		loc, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			t.Fatal(err)
		}

		t.Run("absence timezone want as UTC", func(t *testing.T) {
			want := time.Date(2017, 4, 24, 9, 41, 34, 89312523*10, time.UTC)
			got, err := ParseDateTime("2017-04-24T09:41:34.89312523", WithInLocation(loc))
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})

		t.Run("present timezone offset is the same as Local", func(t *testing.T) {
			backup := time.Local
			time.Local = loc
			defer func() { time.Local = backup }()

			want := time.Date(2017, 4, 24, 9, 41, 34, 89312523*10, time.Local)
			got, err := ParseDateTime("2017-04-24T09:41:34.89312523+09:00")
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})

		t.Run("present timezone offset is the same as JST", func(t *testing.T) {
			want := time.Date(2017, 4, 24, 9, 41, 34, 89312523*10, loc)
			got, err := ParseDateTime("2017-04-24T09:41:34.89312523+09:00", WithInLocation(loc))
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})

		t.Run("present timezone, but different with JST", func(t *testing.T) {
			want := time.Date(2017, 4, 24, 9, 41, 34, 89312523*10, time.FixedZone("", 2*3600))
			got, err := ParseDateTime("2017-04-24T09:41:34.89312523+02:00", WithInLocation(loc))
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	})
}
