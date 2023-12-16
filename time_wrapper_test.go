package synchro_test

import (
	"bytes"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

var _ interface {
	fmt.Stringer
	fmt.GoStringer
	gob.GobEncoder
	gob.GobDecoder
	json.Marshaler
	json.Unmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
} = (*synchro.Time[tz.UTC])(nil)

func TestMain(m *testing.M) {
	os.Setenv("TZ", "America/Los_Angeles")
	defer os.Unsetenv("TZ")
	os.Exit(m.Run())
}

type FixedZone struct{}

func (f FixedZone) Location() *time.Location {
	return time.FixedZone("FixedZone", 9*3600)
}

func TestTime_Local(t *testing.T) {
	d := time.Date(2023, 9, 2, 0, 0, 0, 0, FixedZone{}.Location())
	got := synchro.In[tz.UTC](d).Local()
	want := synchro.In[tz.Local](d)

	if want != got {
		t.Fatalf("want %q but got %q", want, got)
	}
}

func TestTime_GoString(t *testing.T) {
	d := time.Date(2023, 9, 2, 0, 0, 0, 0, time.UTC)
	want := d.GoString()

	s := synchro.In[tz.UTC](d)
	if got := s.GoString(); want != got {
		t.Fatalf("want %q but got %q", want, got)
	}
}

func TestTime_Clock(t *testing.T) {
	d := time.Date(2023, 9, 2, 3, 2, 1, 999999999, time.UTC)
	wantHour, wantMin, wantSec := d.Clock()

	s := synchro.In[tz.UTC](d)
	hour, min, sec := s.Clock()
	if wantHour != hour {
		t.Errorf("want hour %q but got %q", wantHour, hour)
	}
	if wantMin != min {
		t.Errorf("want min %q but got %q", wantMin, min)
	}
	if wantSec != sec {
		t.Errorf("want sec %q but got %q", wantSec, sec)
	}

	t.Run("Hour", func(t *testing.T) {
		want := d.Hour()
		if got := s.Hour(); want != got {
			t.Errorf("want hour %q but got %q", want, got)
		}
	})

	t.Run("Minute", func(t *testing.T) {
		want := d.Minute()
		if got := s.Minute(); want != got {
			t.Errorf("want min %q but got %q", want, got)
		}
	})

	t.Run("Second", func(t *testing.T) {
		want := d.Second()
		if got := s.Second(); want != got {
			t.Errorf("want sec %q but got %q", want, got)
		}
	})

	t.Run("Nanosecond", func(t *testing.T) {
		want := d.Nanosecond()
		if got := s.Nanosecond(); want != got {
			t.Errorf("want nsec %q but got %q", want, got)
		}
	})
}

func TestTime_YearDay(t *testing.T) {
	d := time.Date(2023, 9, 2, 0, 0, 0, 0, time.UTC)
	want := d.YearDay()

	s := synchro.In[tz.UTC](d)
	if got := s.YearDay(); want != got {
		t.Fatalf("want %q but got %q", want, got)
	}
}

func TestTime_ISOWeek(t *testing.T) {
	d := time.Date(2023, 9, 2, 0, 0, 0, 0, time.UTC)
	want1, want2 := d.ISOWeek()

	s := synchro.In[tz.UTC](d)
	year, week := s.ISOWeek()
	if want1 != year {
		t.Fatalf("want year %q but got %q", want1, year)
	}
	if want2 != week {
		t.Fatalf("want week %q but got %q", want2, week)
	}
}

type AustraliaSydney struct{}

func (AustraliaSydney) Location() *time.Location {
	tzWithDST, _ := time.LoadLocation("Australia/Sydney")
	return tzWithDST
}

func TestTime_IsDST(t *testing.T) {
	tests := []struct {
		time synchro.Time[AustraliaSydney]
		want bool
	}{
		{
			time: synchro.New[AustraliaSydney](2009, 1, 1, 12, 0, 0, 0),
			want: true,
		},
		{
			time: synchro.New[AustraliaSydney](2009, 6, 1, 12, 0, 0, 0),
			want: false,
		},
	}
	for _, tt := range tests {
		got := tt.time.IsDST()
		if got != tt.want {
			t.Errorf("(%#v).IsDST()=%t, want %t", tt.time.Format(time.RFC3339), got, tt.want)
		}
	}
}

func TestTime_IsZero(t *testing.T) {
	tm1 := synchro.Time[tz.UTC]{}
	if !tm1.IsZero() {
		t.Errorf("want IsZero = true")
	}
	tm2 := synchro.Now[tz.UTC]()
	if tm2.IsZero() {
		t.Errorf("want IsZero = false")
	}
}

func TestTime_Location(t *testing.T) {
	local := synchro.Time[tz.Local]{}
	got1 := local.Location().String()
	want1 := time.Local.String()
	if got1 != want1 {
		t.Errorf("want Local %q but got %q", want1, got1)
	}

	utc := synchro.Time[tz.UTC]{}
	got2 := utc.Location().String()
	want2 := time.UTC.String()
	if got2 != want2 {
		t.Errorf("want UTC %q but got %q", want2, got2)
	}

	localNow := synchro.Now[tz.Local]()
	got3 := localNow.Location().String()
	want3 := time.Local.String()
	if got3 != want3 {
		t.Errorf("want Local now %q but got %q", want3, got3)
	}
}

func TestTime_Zone(t *testing.T) {
	fz := synchro.Time[FixedZone]{}
	name1, zone1 := fz.Zone()
	want1 := "FixedZone"
	want2 := 9 * 3600
	if want1 != name1 {
		t.Errorf("want name %q but got %q", want1, name1)
	}
	if want2 != zone1 {
		t.Errorf("want zone %d but got %d", want2, zone1)
	}

	fz2 := synchro.Now[FixedZone]()
	name2, zone2 := fz2.Zone()
	if want1 != name2 {
		t.Errorf("want name %q but got %q", want1, name2)
	}
	if want2 != zone2 {
		t.Errorf("want zone %d but got %d", want2, zone2)
	}
}

func TestTime_ZoneBounds(t *testing.T) {
	fz := synchro.New[tz.EuropeAmsterdam](2023, 9, 2, 14, 0, 0, 0)
	start, end := fz.ZoneBounds()

	unwrapStart := start.StdTime()
	name1, zone1 := unwrapStart.Zone()
	if want := "CEST"; want != name1 {
		t.Errorf("want name %q but got %q", want, name1)
	}
	if want := 7200; want != zone1 {
		t.Errorf("want zone %d but got %d", want, zone1)
	}

	unwrapEnd := end.StdTime()
	name2, zone2 := unwrapEnd.Zone()
	if want := "CET"; want != name2 {
		t.Errorf("want name %q but got %q", want, name2)
	}
	if want := 3600; want != zone2 {
		t.Errorf("want zone %d but got %d", want, zone2)
	}
}

func TestTime_Gob(t *testing.T) {
	gobTests := []synchro.Time[tz.UTC]{
		synchro.New[tz.UTC](0, 1, 2, 3, 4, 5, 6),
		synchro.Unix[tz.UTC](81985467080890095, 0x76543210), // time.sec: 0x0123456789ABCDEF
		{}, // nil location
	}

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	dec := gob.NewDecoder(&b)
	for _, tt := range gobTests {
		var gobtt synchro.Time[tz.UTC]
		if err := enc.Encode(&tt); err != nil {
			t.Errorf("%v gob Encode error = %q, want nil", tt, err)
		} else if err := dec.Decode(&gobtt); err != nil {
			t.Errorf("%v gob Decode error = %q, want nil", tt, err)
		} else if !equalTimeAndZone(gobtt, tt) {
			t.Errorf("Decoded time = %v, want %v", gobtt, tt)
		}
		b.Reset()
	}

	{
		tt := synchro.New[FixedZone](2023, 9, 2, 23, 59, 59, 999999999)
		var gobtt synchro.Time[FixedZone]
		if err := enc.Encode(&tt); err != nil {
			t.Errorf("%v gob Encode error = %q, want nil", tt, err)
		} else if err := dec.Decode(&gobtt); err != nil {
			t.Errorf("%v gob Decode error = %q, want nil", tt, err)
		} else if !equalTimeAndZone(gobtt, tt) {
			t.Errorf("Decoded time = %v, want %v", gobtt, tt)
		}
		b.Reset()
	}
}

func equalTimeAndZone[T synchro.TimeZone](a, b synchro.Time[T]) bool {
	aname, aoffset := a.Zone()
	bname, boffset := b.Zone()
	return a.Equal(b) && aoffset == boffset && aname == bname
}

func TestMarshalBinaryZeroTime(t *testing.T) {
	t0 := synchro.Time[tz.UTC]{}
	enc, err := t0.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	t1 := synchro.Now[tz.UTC]() // not zero
	if err := t1.UnmarshalBinary(enc); err != nil {
		t.Fatal(err)
	}
	if t1 != t0 {
		t.Errorf("t0=%#v\nt1=%#v\nwant identical structures", t0, t1)
	}
}

func TestInvalidTimeGob(t *testing.T) {
	invalidEncodingTests := []struct {
		bytes []byte
		want  string
	}{
		{[]byte{}, "Time.UnmarshalBinary: no data"},
		{[]byte{0, 2, 3}, "Time.UnmarshalBinary: unsupported version"},
		{[]byte{1, 2, 3}, "Time.UnmarshalBinary: invalid length"},
	}
	for _, tt := range invalidEncodingTests {
		var ignored synchro.Time[tz.UTC]
		err := ignored.GobDecode(tt.bytes)
		if err == nil || err.Error() != tt.want {
			t.Errorf("time.GobDecode(%#v) error = %v, want %v", tt.bytes, err, tt.want)
		}
		err = ignored.UnmarshalBinary(tt.bytes)
		if err == nil || err.Error() != tt.want {
			t.Errorf("time.UnmarshalBinary(%#v) error = %v, want %v", tt.bytes, err, tt.want)
		}
	}
}

func TestTime_JSON(t *testing.T) {
	testingJSON(
		t,
		synchro.New[tz.UTC](9999, 4, 12, 23, 20, 50, 520*1e6),
		`"9999-04-12T23:20:50.52Z"`,
	)
	testingJSON(
		t,
		synchro.New[tz.Local](1996, 12, 19, 16, 39, 57, 0),
		`"1996-12-19T16:39:57-08:00"`,
	)
	testingJSON(
		t,
		synchro.New[FixedZone](0, 1, 1, 0, 0, 0, 1),
		`"0000-01-01T00:00:00.000000001+09:00"`,
	)
}

func testingJSON[T synchro.TimeZone](t *testing.T, time synchro.Time[T], want string) {
	var jsonTime synchro.Time[T]

	if jsonBytes, err := json.Marshal(time); err != nil {
		t.Errorf("%v json.Marshal error = %v, want nil", time, err)
	} else if string(jsonBytes) != want {
		t.Errorf("%v JSON = %#q, want %#q", time, string(jsonBytes), want)
	} else if err = json.Unmarshal(jsonBytes, &jsonTime); err != nil {
		t.Errorf("%v json.Unmarshal error = %v, want nil", time, err)
	} else if !equalTimeAndZone(jsonTime, time) {
		t.Errorf("Unmarshaled time = %v, want %v", jsonTime, time)
	}
}

func TestMarshalInvalidTimes(t *testing.T) {
	tests := []struct {
		time synchro.Time[tz.UTC]
		want string
	}{
		{synchro.New[tz.UTC](10000, 1, 1, 0, 0, 0, 0), "Time.MarshalJSON: year outside of range [0,9999]"},
		{synchro.New[tz.UTC](-998, 1, 1, 0, 0, 0, 0).Add(-time.Second), "Time.MarshalJSON: year outside of range [0,9999]"},
		{synchro.New[tz.UTC](0, 1, 1, 0, 0, 0, 0).Add(-time.Nanosecond), "Time.MarshalJSON: year outside of range [0,9999]"},
	}

	for _, tt := range tests {
		want := tt.want
		b, err := tt.time.MarshalJSON()
		switch {
		case b != nil:
			t.Errorf("(%v).MarshalText() = %q, want nil", tt.time, b)
		case err == nil || err.Error() != want:
			t.Errorf("(%v).MarshalJSON() error = %v, want %v", tt.time, err, want)
		}

		want = strings.ReplaceAll(tt.want, "JSON", "Text")
		b, err = tt.time.MarshalText()
		switch {
		case b != nil:
			t.Errorf("(%v).MarshalText() = %q, want nil", tt.time, b)
		case err == nil || err.Error() != want:
			t.Errorf("(%v).MarshalText() error = %v, want %v", tt.time, err, want)
		}
	}
}

func TestUnmarshalInvalidTimes(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`{}`, "Time.UnmarshalJSON: input is not a JSON string"},
		{`[]`, "Time.UnmarshalJSON: input is not a JSON string"},
		{`"2000-01-01T1:12:34Z"`, `<nil>`},
		{`"2000-01-01T00:00:00,000Z"`, `<nil>`},
		{`"2000-01-01T00:00:00+24:00"`, `<nil>`},
		{`"2000-01-01T00:00:00+00:60"`, `<nil>`},
		{`"2000-01-01T00:00:00+123:45"`, `parsing time "2000-01-01T00:00:00+123:45" as "2006-01-02T15:04:05Z07:00": cannot parse "+123:45" as "Z07:00"`},
	}

	for _, tt := range tests {
		var ts synchro.Time[tz.UTC]

		want := tt.want
		err := json.Unmarshal([]byte(tt.in), &ts)
		if fmt.Sprint(err) != want {
			t.Errorf("Time.UnmarshalJSON(%s) = %v, want %v", tt.in, err, want)
		}

		if strings.HasPrefix(tt.in, `"`) && strings.HasSuffix(tt.in, `"`) {
			err = ts.UnmarshalText([]byte(strings.Trim(tt.in, `"`)))
			if fmt.Sprint(err) != want {
				t.Errorf("Time.UnmarshalText(%s) = %v, want %v", tt.in, err, want)
			}
		}
	}
}
