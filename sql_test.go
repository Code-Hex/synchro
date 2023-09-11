package synchro_test

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

func TestTime_Scan(t *testing.T) {
	t.Run("UTC", func(t *testing.T) {
		tests := []struct {
			name string
			src  any
			want synchro.Time[tz.UTC]
			err  bool
		}{
			{
				name: "nil",
				src:  nil,
				want: synchro.Time[tz.UTC]{},
				err:  false,
			},
			{
				name: "time.Time",
				src:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				want: synchro.New[tz.UTC](2021, 1, 1, 0, 0, 0, 0),
				err:  false,
			},
			{
				name: "datetime as string",
				src:  "2023-09-10",
				want: synchro.New[tz.UTC](2023, 9, 10, 0, 0, 0, 0),
				err:  false,
			},
			{
				name: "datetime as bytes",
				src:  []byte("2023-09-10"),
				want: synchro.New[tz.UTC](2023, 9, 10, 0, 0, 0, 0),
				err:  false,
			},
			{
				name: "timestamp as string",
				src:  "2023-09-10 14:03:54.000115898+00",
				want: synchro.New[tz.UTC](2023, 9, 10, 14, 3, 54, 115898),
				err:  false,
			},
			{
				name: "timestamp as bytes",
				src:  []byte("2023-09-10 14:03:54.000115898+00"),
				want: synchro.New[tz.UTC](2023, 9, 10, 14, 3, 54, 115898),
				err:  false,
			},
			{
				name: "invalid format as string",
				src:  "unknown",
				want: synchro.Time[tz.UTC]{},
				err:  true,
			},
			{
				name: "invalid format as byte",
				src:  []byte("unknown"),
				want: synchro.Time[tz.UTC]{},
				err:  true,
			},
			{
				name: "unknown type",
				src:  123,
				want: synchro.Time[tz.UTC]{},
				err:  true,
			},
		}
		for _, tt := range tests {
			testingTimeScan(t, tt.src, tt.want, tt.err)
		}
	})

	t.Run("Asia/Tokyo", func(t *testing.T) {
		tests := []struct {
			name string
			src  any
			want synchro.Time[tz.AsiaTokyo]
			err  bool
		}{
			{
				name: "nil",
				src:  nil,
				want: synchro.Time[tz.AsiaTokyo]{},
				err:  false,
			},
			{
				name: "time.Time",
				src:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				want: synchro.New[tz.AsiaTokyo](2021, 1, 1, 9, 0, 0, 0),
				err:  false,
			},
			{
				name: "datetime as string",
				src:  "2023-09-10",
				want: synchro.New[tz.AsiaTokyo](2023, 9, 10, 9, 0, 0, 0),
				err:  false,
			},
			{
				name: "datetime as bytes",
				src:  []byte("2023-09-10"),
				want: synchro.New[tz.AsiaTokyo](2023, 9, 10, 9, 0, 0, 0),
				err:  false,
			},
			{
				name: "timestamp as string",
				src:  "2023-09-10 14:03:54.000115898+00",
				want: synchro.New[tz.AsiaTokyo](2023, 9, 10, 23, 3, 54, 115898),
				err:  false,
			},
			{
				name: "timestamp as bytes",
				src:  []byte("2023-09-10 14:03:54.000115898+00"),
				want: synchro.New[tz.AsiaTokyo](2023, 9, 10, 23, 3, 54, 115898),
				err:  false,
			},
			{
				name: "invalid format as string",
				src:  "unknown",
				want: synchro.Time[tz.AsiaTokyo]{},
				err:  true,
			},
			{
				name: "invalid format as byte",
				src:  []byte("unknown"),
				want: synchro.Time[tz.AsiaTokyo]{},
				err:  true,
			},
			{
				name: "unknown type",
				src:  123,
				want: synchro.Time[tz.AsiaTokyo]{},
				err:  true,
			},
		}
		for _, tt := range tests {
			testingTimeScan(t, tt.src, tt.want, tt.err)
		}
	})
}

func testingTimeScan[T synchro.TimeZone](t *testing.T, src any, want synchro.Time[T], wantErr bool) {
	var got synchro.Time[T]
	err := got.Scan(src)
	if wantErr && err == nil {
		t.Errorf("Scan(%v) should have returned an error, but did not", src)
	} else if !wantErr && err != nil {
		t.Errorf("Scan(%v) returned an unexpected error: %v", src, err)
	}
	if !got.Equal(want) {
		t.Errorf("Scan(%v) = %v, want %v", src, got, want)
	}
}

func TestTime_Value(t *testing.T) {
	now := synchro.Now[tz.EuropeAthens]()
	want := now.StdTime()
	got, err := now.Value()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatalf("want %q but got %q", want, got)
	}
}

func TestNullTime_Scan(t *testing.T) {
	t.Run("UTC", func(t *testing.T) {
		tests := []struct {
			name string
			src  any
			want synchro.NullTime[tz.UTC]
			err  bool
		}{
			{
				name: "nil",
				src:  nil,
				want: synchro.NullTime[tz.UTC]{},
				err:  false,
			},
			{
				name: "time.Time",
				src:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				want: synchro.NullTime[tz.UTC]{
					Valid: true,
					Time:  synchro.New[tz.UTC](2021, 1, 1, 0, 0, 0, 0),
				},
				err: false,
			},
			{
				name: "unknown type",
				src:  123,
				want: synchro.NullTime[tz.UTC]{
					Valid: true,
				},
				err: true,
			},
		}
		for _, tt := range tests {
			testingNullTimeScan(t, tt.src, tt.want, tt.err)
		}
	})

	t.Run("Asia/Tokyo", func(t *testing.T) {
		tests := []struct {
			name string
			src  any
			want synchro.NullTime[tz.AsiaTokyo]
			err  bool
		}{{
			name: "nil",
			src:  nil,
			want: synchro.NullTime[tz.AsiaTokyo]{},
			err:  false,
		},
			{
				name: "time.Time",
				src:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				want: synchro.NullTime[tz.AsiaTokyo]{
					Valid: true,
					Time:  synchro.New[tz.AsiaTokyo](2021, 1, 1, 9, 0, 0, 0),
				},
				err: false,
			},
			{
				name: "unknown type",
				src:  123,
				want: synchro.NullTime[tz.AsiaTokyo]{
					Valid: true,
				},
				err: true,
			},
		}
		for _, tt := range tests {
			testingNullTimeScan(t, tt.src, tt.want, tt.err)
		}
	})
}

func testingNullTimeScan[T synchro.TimeZone](t *testing.T, src any, want synchro.NullTime[T], wantErr bool) {
	var got synchro.NullTime[T]
	err := got.Scan(src)
	if wantErr && err == nil {
		t.Errorf("Scan(%v) should have returned an error, but did not", src)
	} else if !wantErr && err != nil {
		t.Errorf("Scan(%v) returned an unexpected error: %v", src, err)
	}
	if got.Valid != want.Valid {
		t.Errorf("Scan(%v).Valid = %v, want %v", src, got.Valid, want.Valid)
	}
	if !got.Time.Equal(want.Time) {
		t.Errorf("Scan(%v).Time = %v, want %v", src, got.Time, want.Time)
	}
}

func TestNullTime_Value(t *testing.T) {
	now := synchro.Now[tz.EuropeAthens]()
	t.Run("valid", func(t *testing.T) {
		nullish := synchro.NullTime[tz.EuropeAthens]{
			Valid: true,
			Time:  now,
		}
		want := now.StdTime()
		got, err := nullish.Value()
		if err != nil {
			t.Fatal(err)
		}
		if want != got {
			t.Fatalf("want %q but got %q", want, got)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		nullish := synchro.NullTime[tz.EuropeAthens]{
			Valid: false,
			Time:  now,
		}
		want := driver.Valuer(nil)
		got, err := nullish.Value()
		if err != nil {
			t.Fatal(err)
		}
		if want != got {
			t.Fatalf("want %q but got %q", want, got)
		}
	})

}
