package iso8601

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDate_IsValid(t *testing.T) {
	tests := []struct {
		name string
		date Date
		want bool
	}{
		{
			name: "invalid zero day",
			date: Date{Year: 0, Month: 0, Day: 0},
			want: false,
		},
		{
			name: "valid date",
			date: Date{Year: 2022, Month: 1, Day: 1},
			want: true,
		},
		{
			name: "invalid month",
			date: Date{Year: 2022, Month: 13, Day: 1},
			want: false,
		},
		{
			name: "invalid day",
			date: Date{Year: 2022, Month: 2, Day: 29},
			want: false,
		},
		{
			name: "valid leap year",
			date: Date{Year: 2020, Month: 2, Day: 29},
			want: true,
		},
		{
			name: "invalid leap year",
			date: Date{Year: 2021, Month: 2, Day: 29},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.IsValid(); got != tt.want {
				t.Errorf("Date.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuarterDate_IsValid(t *testing.T) {
	tests := []struct {
		name string
		date QuarterDate
		want bool
	}{
		{
			name: "invalid zero day",
			date: QuarterDate{Year: 0, Quarter: 0, Day: 0},
			want: false,
		},
		{
			name: "valid date",
			date: QuarterDate{Year: 2022, Quarter: 1, Day: 1},
			want: true,
		},
		{
			name: "invalid quarter",
			date: QuarterDate{Year: 2022, Quarter: 5, Day: 1},
			want: false,
		},
		{
			name: "invalid day",
			date: QuarterDate{Year: 2022, Quarter: 1, Day: 91},
			want: false,
		},
		{
			name: "valid last day of quarter",
			date: QuarterDate{Year: 2022, Quarter: 1, Day: 90},
			want: true,
		},
		{
			name: "valid last day of leap year quarter",
			date: QuarterDate{Year: 2020, Quarter: 1, Day: 91},
			want: true,
		},
		{
			name: "invalid last day of non-leap year quarter",
			date: QuarterDate{Year: 2021, Quarter: 1, Day: 91},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.IsValid(); got != tt.want {
				t.Errorf("QuarterDate.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekDate_IsValid(t *testing.T) {
	tests := []struct {
		name string
		date WeekDate
		want bool
	}{
		{
			name: "invalid zero day",
			date: WeekDate{Year: 0, Week: 0, Day: 0},
			want: false,
		},
		{
			name: "valid date",
			date: WeekDate{Year: 2022, Week: 1, Day: 1},
			want: true,
		},
		{
			name: "invalid week",
			date: WeekDate{Year: 2022, Week: 53, Day: 1},
			want: false,
		},
		{
			name: "invalid day",
			date: WeekDate{Year: 2022, Week: 1, Day: 8},
			want: false,
		},
		{
			name: "valid last day of year",
			date: WeekDate{Year: 2022, Week: 52, Day: 7},
			want: true,
		},
		{
			name: "valid last day of leap year",
			date: WeekDate{Year: 2020, Week: 53, Day: 7},
			want: true,
		},
		{
			name: "invalid last day of non-leap year",
			date: WeekDate{Year: 2021, Week: 53, Day: 7},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.IsValid(); got != tt.want {
				t.Errorf("WeekDate.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdinalDate_IsValid(t *testing.T) {
	tests := []struct {
		name string
		date OrdinalDate
		want bool
	}{
		{
			name: "invalid zero day",
			date: OrdinalDate{Year: 0, Day: 0},
			want: false,
		},
		{
			name: "valid date",
			date: OrdinalDate{Year: 2022, Day: 1},
			want: true,
		},
		{
			name: "invalid day",
			date: OrdinalDate{Year: 2022, Day: 366},
			want: false,
		},
		{
			name: "valid last day of non-leap year",
			date: OrdinalDate{Year: 2021, Day: 365},
			want: true,
		},
		{
			name: "valid last day of leap year",
			date: OrdinalDate{Year: 2020, Day: 366},
			want: true,
		},
		{
			name: "invalid day of leap year",
			date: OrdinalDate{Year: 2020, Day: 367},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.IsValid(); got != tt.want {
				t.Errorf("OrdinalDate.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdinalDate_Date(t *testing.T) {
	tests := [...]struct {
		o    OrdinalDate
		want Date
	}{
		0: {
			o: OrdinalDate{
				Year: 2008,
				Day:  58,
			},
			want: Date{
				Year:  2008,
				Month: 2,
				Day:   27,
			},
		},
		1: {
			o: OrdinalDate{
				Year: 2008,
				Day:  59,
			},
			want: Date{
				Year:  2008,
				Month: 2,
				Day:   28,
			},
		},
		2: {
			o: OrdinalDate{
				Year: 2008,
				Day:  60,
			},
			want: Date{
				Year:  2008,
				Month: 2,
				Day:   29,
			},
		},
		3: {
			o: OrdinalDate{
				Year: 2008,
				Day:  61,
			},
			want: Date{
				Year:  2008,
				Month: 3,
				Day:   1,
			},
		},
		4: {
			o: OrdinalDate{
				Year: 2009,
				Day:  1,
			},
			want: Date{
				Year:  2009,
				Month: 1,
				Day:   1,
			},
		},
		5: {
			o: OrdinalDate{
				Year: 2009,
				Day:  2,
			},
			want: Date{
				Year:  2009,
				Month: 1,
				Day:   2,
			},
		},
		6: {
			o: OrdinalDate{
				Year: 2009,
				Day:  58,
			},
			want: Date{
				Year:  2009,
				Month: 2,
				Day:   27,
			},
		},
		7: {
			o: OrdinalDate{
				Year: 2009,
				Day:  59,
			},
			want: Date{
				Year:  2009,
				Month: 2,
				Day:   28,
			},
		},
		8: {
			o: OrdinalDate{
				Year: 2009,
				Day:  60,
			},
			want: Date{
				Year:  2009,
				Month: 3,
				Day:   1,
			},
		},
		9: {
			o: OrdinalDate{
				Year: 2009,
				Day:  305,
			},
			want: Date{
				Year:  2009,
				Month: 11,
				Day:   1,
			},
		},
		10: {
			o: OrdinalDate{
				Year: 2009,
				Day:  306,
			},
			want: Date{
				Year:  2009,
				Month: 11,
				Day:   2,
			},
		},
		11: {
			o: OrdinalDate{
				Year: 2009,
				Day:  334,
			},
			want: Date{
				Year:  2009,
				Month: 11,
				Day:   30,
			},
		},
		12: {
			o: OrdinalDate{
				Year: 2009,
				Day:  335,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   1,
			},
		},
		13: {
			o: OrdinalDate{
				Year: 2009,
				Day:  336,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   2,
			},
		},
		14: {
			o: OrdinalDate{
				Year: 2009,
				Day:  348,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   14,
			},
		},
		15: {
			o: OrdinalDate{
				Year: 2009,
				Day:  349,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   15,
			},
		},
		16: {
			o: OrdinalDate{
				Year: 2009,
				Day:  350,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   16,
			},
		},
		17: {
			o: OrdinalDate{
				Year: 2009,
				Day:  364,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   30,
			},
		},
		18: {
			o: OrdinalDate{
				Year: 2009,
				Day:  365,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   31,
			},
		},
		19: {
			o: OrdinalDate{
				Year: 2009,
				Day:  365,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   31,
			},
		},
		20: {
			o: OrdinalDate{
				Year: 2010,
				Day:  2,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   2,
			},
		},
		21: {
			o: OrdinalDate{
				Year: 2010,
				Day:  9,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   9,
			},
		},
		22: {
			o: OrdinalDate{
				Year: 2010,
				Day:  10,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   10,
			},
		},
		23: {
			o: OrdinalDate{
				Year: 2010,
				Day:  11,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   11,
			},
		},
		24: {
			o: OrdinalDate{
				Year: 2010,
				Day:  14,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   14,
			},
		},
		25: {
			o: OrdinalDate{
				Year: 2010,
				Day:  15,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   15,
			},
		},
		26: {
			o: OrdinalDate{
				Year: 2010,
				Day:  31,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   31,
			},
		},
		27: {
			o: OrdinalDate{
				Year: 2010,
				Day:  32,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   1,
			},
		},
		28: {
			o: OrdinalDate{
				Year: 2010,
				Day:  40,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   9,
			},
		},
		29: {
			o: OrdinalDate{
				Year: 2010,
				Day:  41,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   10,
			},
		},
		30: {
			o: OrdinalDate{
				Year: 2010,
				Day:  59,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   28,
			},
		},
		31: {
			o: OrdinalDate{
				Year: 2010,
				Day:  60,
			},
			want: Date{
				Year:  2010,
				Month: 3,
				Day:   1,
			},
		},
		32: {
			o: OrdinalDate{
				Year: 2010,
				Day:  68,
			},
			want: Date{
				Year:  2010,
				Month: 3,
				Day:   9,
			},
		},
		33: {
			o: OrdinalDate{
				Year: 2010,
				Day:  69,
			},
			want: Date{
				Year:  2010,
				Month: 3,
				Day:   10,
			},
		},
		34: {
			o: OrdinalDate{
				Year: 2010,
				Day:  365,
			},
			want: Date{
				Year:  2010,
				Month: 12,
				Day:   31,
			},
		},
		35: {
			o: OrdinalDate{
				Year: 2011,
				Day:  1,
			},
			want: Date{
				Year:  2011,
				Month: 1,
				Day:   1,
			},
		},
		36: {
			o: OrdinalDate{
				Year: 2011,
				Day:  9,
			},
			want: Date{
				Year:  2011,
				Month: 1,
				Day:   9,
			},
		},
		37: {
			o: OrdinalDate{
				Year: 2011,
				Day:  10,
			},
			want: Date{
				Year:  2011,
				Month: 1,
				Day:   10,
			},
		},
		38: {
			o: OrdinalDate{
				Year: 2011,
				Day:  121,
			},
			want: Date{
				Year:  2011,
				Month: 5,
				Day:   1,
			},
		},
		39: {
			o: OrdinalDate{
				Year: 2011,
				Day:  365,
			},
			want: Date{
				Year:  2011,
				Month: 12,
				Day:   31,
			},
		},
		40: {
			o: OrdinalDate{
				Year: 2012,
				Day:  1,
			},
			want: Date{
				Year:  2012,
				Month: 1,
				Day:   1,
			},
		},
		41: {
			o: OrdinalDate{
				Year: 2012,
				Day:  58,
			},
			want: Date{
				Year:  2012,
				Month: 2,
				Day:   27,
			},
		},
		42: {
			o: OrdinalDate{
				Year: 2012,
				Day:  59,
			},
			want: Date{
				Year:  2012,
				Month: 2,
				Day:   28,
			},
		},
		43: {
			o: OrdinalDate{
				Year: 2012,
				Day:  60,
			},
			want: Date{
				Year:  2012,
				Month: 2,
				Day:   29,
			},
		},
		44: {
			o: OrdinalDate{
				Year: 2014,
				Day:  58,
			},
			want: Date{
				Year:  2014,
				Month: 2,
				Day:   27,
			},
		},
		45: {
			o: OrdinalDate{
				Year: 2014,
				Day:  59,
			},
			want: Date{
				Year:  2014,
				Month: 2,
				Day:   28,
			},
		},
		46: {
			o: OrdinalDate{
				Year: 2014,
				Day:  60,
			},
			want: Date{
				Year:  2014,
				Month: 3,
				Day:   1,
			},
		},
		47: {
			o: OrdinalDate{
				Year: 2014,
				Day:  61,
			},
			want: Date{
				Year:  2014,
				Month: 3,
				Day:   2,
			},
		},
		48: {
			o: OrdinalDate{
				Year: 2016,
				Day:  59,
			},
			want: Date{
				Year:  2016,
				Month: 2,
				Day:   28,
			},
		},
		49: {
			o: OrdinalDate{
				Year: 2016,
				Day:  60,
			},
			want: Date{
				Year:  2016,
				Month: 2,
				Day:   29,
			},
		},
		50: {
			o: OrdinalDate{
				Year: 2016,
				Day:  61,
			},
			want: Date{
				Year:  2016,
				Month: 3,
				Day:   1,
			},
		},
	}
	for i, tt := range tests {
		name := fmt.Sprintf("case %d", i)
		t.Run(name, func(t *testing.T) {
			got := tt.o.Date()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestQuarterDate_Date(t *testing.T) {
	tests := [...]struct {
		q    QuarterDate
		want Date
	}{
		0: {
			q: QuarterDate{
				Year:    2012,
				Quarter: 4,
				Day:     85,
			},
			want: Date{
				Year:  2012,
				Month: time.December,
				Day:   24,
			},
		},
		1: {
			q: QuarterDate{
				Year:    2000,
				Quarter: 1,
				Day:     38,
			},
			want: Date{
				Year:  2000,
				Month: 2,
				Day:   7,
			},
		},
		2: {
			q: QuarterDate{
				Year:    2000,
				Quarter: 2,
				Day:     21,
			},
			want: Date{
				Year:  2000,
				Month: 4,
				Day:   21,
			},
		},
		3: {
			q: QuarterDate{
				Year:    2000,
				Quarter: 3,
				Day:     82,
			},
			want: Date{
				Year:  2000,
				Month: 9,
				Day:   20,
			},
		},
		4: {
			q: QuarterDate{
				Year:    2000,
				Quarter: 4,
				Day:     11,
			},
			want: Date{
				Year:  2000,
				Month: 10,
				Day:   11,
			},
		},
		5: {
			q: QuarterDate{
				Year:    2001,
				Quarter: 1,
				Day:     60,
			},
			want: Date{
				Year:  2001,
				Month: 3,
				Day:   1,
			},
		},
		6: {
			q: QuarterDate{
				Year:    2001,
				Quarter: 2,
				Day:     67,
			},
			want: Date{
				Year:  2001,
				Month: 6,
				Day:   6,
			},
		},
		7: {
			q: QuarterDate{
				Year:    2001,
				Quarter: 3,
				Day:     35,
			},
			want: Date{
				Year:  2001,
				Month: 8,
				Day:   4,
			},
		},
		8: {
			q: QuarterDate{
				Year:    2001,
				Quarter: 4,
				Day:     52,
			},
			want: Date{
				Year:  2001,
				Month: 11,
				Day:   21,
			},
		},
		9: {
			q: QuarterDate{
				Year:    2002,
				Quarter: 1,
				Day:     14,
			},
			want: Date{
				Year:  2002,
				Month: 1,
				Day:   14,
			},
		},
		10: {
			q: QuarterDate{
				Year:    2002,
				Quarter: 2,
				Day:     55,
			},
			want: Date{
				Year:  2002,
				Month: 5,
				Day:   25,
			},
		},
		11: {
			q: QuarterDate{
				Year:    2002,
				Quarter: 3,
				Day:     50,
			},
			want: Date{
				Year:  2002,
				Month: 8,
				Day:   19,
			},
		},
		12: {
			q: QuarterDate{
				Year:    2002,
				Quarter: 4,
				Day:     47,
			},
			want: Date{
				Year:  2002,
				Month: 11,
				Day:   16,
			},
		},
		13: {
			q: QuarterDate{
				Year:    2003,
				Quarter: 1,
				Day:     38,
			},
			want: Date{
				Year:  2003,
				Month: 2,
				Day:   7,
			},
		},
		14: {
			q: QuarterDate{
				Year:    2003,
				Quarter: 2,
				Day:     25,
			},
			want: Date{
				Year:  2003,
				Month: 4,
				Day:   25,
			},
		},
		15: {
			q: QuarterDate{
				Year:    2003,
				Quarter: 3,
				Day:     28,
			},
			want: Date{
				Year:  2003,
				Month: 7,
				Day:   28,
			},
		},
		16: {
			q: QuarterDate{
				Year:    2008,
				Quarter: 1,
				Day:     58,
			},
			want: Date{
				Year:  2008,
				Month: 2,
				Day:   27,
			},
		},
		17: {
			q: QuarterDate{
				Year:    2008,
				Quarter: 1,
				Day:     59,
			},
			want: Date{
				Year:  2008,
				Month: 2,
				Day:   28,
			},
		},
		18: {
			q: QuarterDate{
				Year:    2008,
				Quarter: 1,
				Day:     60,
			},
			want: Date{
				Year:  2008,
				Month: 2,
				Day:   29,
			},
		},
		19: {
			q: QuarterDate{
				Year:    2008,
				Quarter: 1,
				Day:     61,
			},
			want: Date{
				Year:  2008,
				Month: 3,
				Day:   1,
			},
		},
		20: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 1,
				Day:     1,
			},
			want: Date{
				Year:  2009,
				Month: 1,
				Day:   1,
			},
		},
		21: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 1,
				Day:     2,
			},
			want: Date{
				Year:  2009,
				Month: 1,
				Day:   2,
			},
		},
		22: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 1,
				Day:     58,
			},
			want: Date{
				Year:  2009,
				Month: 2,
				Day:   27,
			},
		},
		23: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 1,
				Day:     59,
			},
			want: Date{
				Year:  2009,
				Month: 2,
				Day:   28,
			},
		},
		24: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 1,
				Day:     60,
			},
			want: Date{
				Year:  2009,
				Month: 3,
				Day:   1,
			},
		},
		25: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     32,
			},
			want: Date{
				Year:  2009,
				Month: 11,
				Day:   1,
			},
		},
		26: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     33,
			},
			want: Date{
				Year:  2009,
				Month: 11,
				Day:   2,
			},
		},
		27: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     61,
			},
			want: Date{
				Year:  2009,
				Month: 11,
				Day:   30,
			},
		},
		28: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     62,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   1,
			},
		},
		29: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     63,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   2,
			},
		},
		30: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     75,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   14,
			},
		},
		31: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     76,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   15,
			},
		},
		32: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     77,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   16,
			},
		},
		33: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     91,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   30,
			},
		},
		34: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     92,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   31,
			},
		},
		35: {
			q: QuarterDate{
				Year:    2009,
				Quarter: 4,
				Day:     92,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   31,
			},
		},
		36: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     2,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   2,
			},
		},
		37: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     9,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   9,
			},
		},
		38: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     10,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   10,
			},
		},
		39: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     11,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   11,
			},
		},
		40: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     14,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   14,
			},
		},
		41: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     15,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   15,
			},
		},
		42: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     31,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   31,
			},
		},
		43: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     32,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   1,
			},
		},
		44: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     40,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   9,
			},
		},
		45: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     41,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   10,
			},
		},
		46: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     59,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   28,
			},
		},
		47: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     60,
			},
			want: Date{
				Year:  2010,
				Month: 3,
				Day:   1,
			},
		},
		48: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     68,
			},
			want: Date{
				Year:  2010,
				Month: 3,
				Day:   9,
			},
		},
		49: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 1,
				Day:     69,
			},
			want: Date{
				Year:  2010,
				Month: 3,
				Day:   10,
			},
		},
		50: {
			q: QuarterDate{
				Year:    2010,
				Quarter: 4,
				Day:     92,
			},
			want: Date{
				Year:  2010,
				Month: 12,
				Day:   31,
			},
		},
		51: {
			q: QuarterDate{
				Year:    2011,
				Quarter: 1,
				Day:     1,
			},
			want: Date{
				Year:  2011,
				Month: 1,
				Day:   1,
			},
		},
		52: {
			q: QuarterDate{
				Year:    2011,
				Quarter: 1,
				Day:     9,
			},
			want: Date{
				Year:  2011,
				Month: 1,
				Day:   9,
			},
		},
		53: {
			q: QuarterDate{
				Year:    2011,
				Quarter: 1,
				Day:     10,
			},
			want: Date{
				Year:  2011,
				Month: 1,
				Day:   10,
			},
		},
		54: {
			q: QuarterDate{
				Year:    2011,
				Quarter: 2,
				Day:     31,
			},
			want: Date{
				Year:  2011,
				Month: 5,
				Day:   1,
			},
		},
		55: {
			q: QuarterDate{
				Year:    2011,
				Quarter: 4,
				Day:     92,
			},
			want: Date{
				Year:  2011,
				Month: 12,
				Day:   31,
			},
		},
		56: {
			q: QuarterDate{
				Year:    2012,
				Quarter: 1,
				Day:     1,
			},
			want: Date{
				Year:  2012,
				Month: 1,
				Day:   1,
			},
		},
		57: {
			q: QuarterDate{
				Year:    2012,
				Quarter: 1,
				Day:     58,
			},
			want: Date{
				Year:  2012,
				Month: 2,
				Day:   27,
			},
		},
		58: {
			q: QuarterDate{
				Year:    2012,
				Quarter: 1,
				Day:     59,
			},
			want: Date{
				Year:  2012,
				Month: 2,
				Day:   28,
			},
		},
		59: {
			q: QuarterDate{
				Year:    2012,
				Quarter: 1,
				Day:     60,
			},
			want: Date{
				Year:  2012,
				Month: 2,
				Day:   29,
			},
		},
		60: {
			q: QuarterDate{
				Year:    2014,
				Quarter: 1,
				Day:     58,
			},
			want: Date{
				Year:  2014,
				Month: 2,
				Day:   27,
			},
		},
		61: {
			q: QuarterDate{
				Year:    2014,
				Quarter: 1,
				Day:     59,
			},
			want: Date{
				Year:  2014,
				Month: 2,
				Day:   28,
			},
		},
		62: {
			q: QuarterDate{
				Year:    2014,
				Quarter: 1,
				Day:     60,
			},
			want: Date{
				Year:  2014,
				Month: 3,
				Day:   1,
			},
		},
		63: {
			q: QuarterDate{
				Year:    2014,
				Quarter: 1,
				Day:     61,
			},
			want: Date{
				Year:  2014,
				Month: 3,
				Day:   2,
			},
		},
		64: {
			q: QuarterDate{
				Year:    2016,
				Quarter: 1,
				Day:     59,
			},
			want: Date{
				Year:  2016,
				Month: 2,
				Day:   28,
			},
		},
		65: {
			q: QuarterDate{
				Year:    2016,
				Quarter: 1,
				Day:     60,
			},
			want: Date{
				Year:  2016,
				Month: 2,
				Day:   29,
			},
		},
		66: {
			q: QuarterDate{
				Year:    2016,
				Quarter: 1,
				Day:     61,
			},
			want: Date{
				Year:  2016,
				Month: 3,
				Day:   1,
			},
		},
	}
	for i, tt := range tests {
		name := fmt.Sprintf("case %d", i)
		t.Run(name, func(t *testing.T) {
			got := tt.q.Date()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestWeekDate_Date(t *testing.T) {
	tests := [...]struct {
		w    WeekDate
		want Date
	}{
		0: {
			w: WeekDate{
				Year: 2008,
				Week: 9,
				Day:  3,
			},
			want: Date{
				Year:  2008,
				Month: 2,
				Day:   27,
			},
		},
		1: {
			w: WeekDate{
				Year: 2008,
				Week: 9,
				Day:  4,
			},
			want: Date{
				Year:  2008,
				Month: 2,
				Day:   28,
			},
		},
		2: {
			w: WeekDate{
				Year: 2008,
				Week: 9,
				Day:  5,
			},
			want: Date{
				Year:  2008,
				Month: 2,
				Day:   29,
			},
		},
		3: {
			w: WeekDate{
				Year: 2008,
				Week: 9,
				Day:  6,
			},
			want: Date{
				Year:  2008,
				Month: 3,
				Day:   1,
			},
		},
		4: {
			w: WeekDate{
				Year: 2009,
				Week: 1,
				Day:  4,
			},
			want: Date{
				Year:  2009,
				Month: 1,
				Day:   1,
			},
		},
		5: {
			w: WeekDate{
				Year: 2009,
				Week: 1,
				Day:  5,
			},
			want: Date{
				Year:  2009,
				Month: 1,
				Day:   2,
			},
		},
		6: {
			w: WeekDate{
				Year: 2009,
				Week: 9,
				Day:  5,
			},
			want: Date{
				Year:  2009,
				Month: 2,
				Day:   27,
			},
		},
		7: {
			w: WeekDate{
				Year: 2009,
				Week: 9,
				Day:  6,
			},
			want: Date{
				Year:  2009,
				Month: 2,
				Day:   28,
			},
		},
		8: {
			w: WeekDate{
				Year: 2009,
				Week: 9,
				Day:  7,
			},
			want: Date{
				Year:  2009,
				Month: 3,
				Day:   1,
			},
		},
		9: {
			w: WeekDate{
				Year: 2009,
				Week: 44,
				Day:  7,
			},
			want: Date{
				Year:  2009,
				Month: 11,
				Day:   1,
			},
		},
		10: {
			w: WeekDate{
				Year: 2009,
				Week: 45,
				Day:  1,
			},
			want: Date{
				Year:  2009,
				Month: 11,
				Day:   2,
			},
		},
		11: {
			w: WeekDate{
				Year: 2009,
				Week: 49,
				Day:  1,
			},
			want: Date{
				Year:  2009,
				Month: 11,
				Day:   30,
			},
		},
		12: {
			w: WeekDate{
				Year: 2009,
				Week: 49,
				Day:  2,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   1,
			},
		},
		13: {
			w: WeekDate{
				Year: 2009,
				Week: 49,
				Day:  3,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   2,
			},
		},
		14: {
			w: WeekDate{
				Year: 2009,
				Week: 51,
				Day:  1,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   14,
			},
		},
		15: {
			w: WeekDate{
				Year: 2009,
				Week: 51,
				Day:  2,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   15,
			},
		},
		16: {
			w: WeekDate{
				Year: 2009,
				Week: 51,
				Day:  3,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   16,
			},
		},
		17: {
			w: WeekDate{
				Year: 2009,
				Week: 53,
				Day:  3,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   30,
			},
		},
		18: {
			w: WeekDate{
				Year: 2009,
				Week: 53,
				Day:  4,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   31,
			},
		},
		19: {
			w: WeekDate{
				Year: 2009,
				Week: 53,
				Day:  4,
			},
			want: Date{
				Year:  2009,
				Month: 12,
				Day:   31,
			},
		},
		20: {
			w: WeekDate{
				Year: 2009,
				Week: 53,
				Day:  6,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   2,
			},
		},
		21: {
			w: WeekDate{
				Year: 2010,
				Week: 1,
				Day:  6,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   9,
			},
		},
		22: {
			w: WeekDate{
				Year: 2010,
				Week: 1,
				Day:  7,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   10,
			},
		},
		23: {
			w: WeekDate{
				Year: 2010,
				Week: 2,
				Day:  1,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   11,
			},
		},
		24: {
			w: WeekDate{
				Year: 2010,
				Week: 2,
				Day:  4,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   14,
			},
		},
		25: {
			w: WeekDate{
				Year: 2010,
				Week: 2,
				Day:  5,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   15,
			},
		},
		26: {
			w: WeekDate{
				Year: 2010,
				Week: 4,
				Day:  7,
			},
			want: Date{
				Year:  2010,
				Month: 1,
				Day:   31,
			},
		},
		27: {
			w: WeekDate{
				Year: 2010,
				Week: 5,
				Day:  1,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   1,
			},
		},
		28: {
			w: WeekDate{
				Year: 2010,
				Week: 6,
				Day:  2,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   9,
			},
		},
		29: {
			w: WeekDate{
				Year: 2010,
				Week: 6,
				Day:  3,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   10,
			},
		},
		30: {
			w: WeekDate{
				Year: 2010,
				Week: 8,
				Day:  7,
			},
			want: Date{
				Year:  2010,
				Month: 2,
				Day:   28,
			},
		},
		31: {
			w: WeekDate{
				Year: 2010,
				Week: 9,
				Day:  1,
			},
			want: Date{
				Year:  2010,
				Month: 3,
				Day:   1,
			},
		},
		32: {
			w: WeekDate{
				Year: 2010,
				Week: 10,
				Day:  2,
			},
			want: Date{
				Year:  2010,
				Month: 3,
				Day:   9,
			},
		},
		33: {
			w: WeekDate{
				Year: 2010,
				Week: 10,
				Day:  3,
			},
			want: Date{
				Year:  2010,
				Month: 3,
				Day:   10,
			},
		},
		34: {
			w: WeekDate{
				Year: 2010,
				Week: 52,
				Day:  5,
			},
			want: Date{
				Year:  2010,
				Month: 12,
				Day:   31,
			},
		},
		35: {
			w: WeekDate{
				Year: 2010,
				Week: 52,
				Day:  6,
			},
			want: Date{
				Year:  2011,
				Month: 1,
				Day:   1,
			},
		},
		36: {
			w: WeekDate{
				Year: 2011,
				Week: 1,
				Day:  7,
			},
			want: Date{
				Year:  2011,
				Month: 1,
				Day:   9,
			},
		},
		37: {
			w: WeekDate{
				Year: 2011,
				Week: 2,
				Day:  1,
			},
			want: Date{
				Year:  2011,
				Month: 1,
				Day:   10,
			},
		},
		38: {
			w: WeekDate{
				Year: 2011,
				Week: 17,
				Day:  7,
			},
			want: Date{
				Year:  2011,
				Month: 5,
				Day:   1,
			},
		},
		39: {
			w: WeekDate{
				Year: 2011,
				Week: 52,
				Day:  6,
			},
			want: Date{
				Year:  2011,
				Month: 12,
				Day:   31,
			},
		},
		40: {
			w: WeekDate{
				Year: 2011,
				Week: 52,
				Day:  7,
			},
			want: Date{
				Year:  2012,
				Month: 1,
				Day:   1,
			},
		},
		41: {
			w: WeekDate{
				Year: 2012,
				Week: 9,
				Day:  1,
			},
			want: Date{
				Year:  2012,
				Month: 2,
				Day:   27,
			},
		},
		42: {
			w: WeekDate{
				Year: 2012,
				Week: 9,
				Day:  2,
			},
			want: Date{
				Year:  2012,
				Month: 2,
				Day:   28,
			},
		},
		43: {
			w: WeekDate{
				Year: 2012,
				Week: 9,
				Day:  3,
			},
			want: Date{
				Year:  2012,
				Month: 2,
				Day:   29,
			},
		},
		44: {
			w: WeekDate{
				Year: 2014,
				Week: 9,
				Day:  4,
			},
			want: Date{
				Year:  2014,
				Month: 2,
				Day:   27,
			},
		},
		45: {
			w: WeekDate{
				Year: 2014,
				Week: 9,
				Day:  5,
			},
			want: Date{
				Year:  2014,
				Month: 2,
				Day:   28,
			},
		},
		46: {
			w: WeekDate{
				Year: 2014,
				Week: 9,
				Day:  6,
			},
			want: Date{
				Year:  2014,
				Month: 3,
				Day:   1,
			},
		},
		47: {
			w: WeekDate{
				Year: 2014,
				Week: 9,
				Day:  7,
			},
			want: Date{
				Year:  2014,
				Month: 3,
				Day:   2,
			},
		},
		48: {
			w: WeekDate{
				Year: 2016,
				Week: 8,
				Day:  7,
			},
			want: Date{
				Year:  2016,
				Month: 2,
				Day:   28,
			},
		},
		49: {
			w: WeekDate{
				Year: 2016,
				Week: 9,
				Day:  1,
			},
			want: Date{
				Year:  2016,
				Month: 2,
				Day:   29,
			},
		},
		50: {
			w: WeekDate{
				Year: 2016,
				Week: 9,
				Day:  2,
			},
			want: Date{
				Year:  2016,
				Month: 3,
				Day:   1,
			},
		},
	}
	for i, tt := range tests {
		name := fmt.Sprintf("case %d", i)
		t.Run(name, func(t *testing.T) {
			got := tt.w.Date()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				w := time.Date(tt.want.Year, tt.want.Month, tt.want.Day, 0, 0, 0, 0, time.UTC)
				g := time.Date(got.Year, got.Month, got.Day, 0, 0, 0, 0, time.UTC)
				t.Errorf("(-want, +got)\n%s- %q (%s)\n+ %q (%s)", diff, w, w.Weekday(), g, g.Weekday())
			}
		})
	}
}

func TestDateLikeRangeError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *DateLikeRangeError
		want string
	}{
		{
			name: "valid error",
			err: &DateLikeRangeError{
				Element: "month",
				Value:   13,
				Year:    2022,
				Min:     1,
				Max:     12,
			},
			want: "iso8601: 13 month is not in range 1-12 in 2022",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("DateLikeRangeError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Validate(t *testing.T) {
	tests := []struct {
		name    string
		d       Date
		wantErr error
	}{
		{
			name: "invalid year is less than 0",
			d: Date{
				Year:  -1,
				Month: 1,
				Day:   1,
			},
			wantErr: &DateLikeRangeError{
				Element: "year",
				Value:   -1,
				Year:    -1,
				Min:     0,
				Max:     9999,
			},
		},
		{
			name: "invalid year is more than 9999",
			d: Date{
				Year:  10000,
				Month: 1,
				Day:   1,
			},
			wantErr: &DateLikeRangeError{
				Element: "year",
				Value:   10000,
				Year:    10000,
				Min:     0,
				Max:     9999,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.d.Validate()
			if err == nil {
				t.Fatal("expected error")
			}
			if diff := cmp.Diff(tt.wantErr, err); diff != "" {
				t.Errorf("error: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestQuarterDate_Validate(t *testing.T) {
	tests := []struct {
		name    string
		q       QuarterDate
		wantErr error
	}{
		{
			name: "invalid year is less than 0",
			q: QuarterDate{
				Year:    -1,
				Quarter: 1,
				Day:     1,
			},
			wantErr: &DateLikeRangeError{
				Element: "year",
				Value:   -1,
				Year:    -1,
				Min:     0,
				Max:     9999,
			},
		},
		{
			name: "invalid year is more than 9999",
			q: QuarterDate{
				Year:    10000,
				Quarter: 1,
				Day:     1,
			},
			wantErr: &DateLikeRangeError{
				Element: "year",
				Value:   10000,
				Year:    10000,
				Min:     0,
				Max:     9999,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.q.Validate()
			if err == nil {
				t.Fatal("expected error")
			}
			if diff := cmp.Diff(tt.wantErr, err); diff != "" {
				t.Errorf("error: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestWeekDate_Validate(t *testing.T) {
	tests := []struct {
		name    string
		w       WeekDate
		wantErr error
	}{
		{
			name: "invalid year is less than 0",
			w: WeekDate{
				Year: -1,
				Week: 10,
				Day:  1,
			},
			wantErr: &DateLikeRangeError{
				Element: "year",
				Value:   -1,
				Year:    -1,
				Min:     0,
				Max:     9999,
			},
		},
		{
			name: "invalid year is more than 9999",
			w: WeekDate{
				Year: 10000,
				Week: 10,
				Day:  1,
			},
			wantErr: &DateLikeRangeError{
				Element: "year",
				Value:   10000,
				Year:    10000,
				Min:     0,
				Max:     9999,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.w.Validate()
			if err == nil {
				t.Fatal("expected error")
			}
			if diff := cmp.Diff(tt.wantErr, err); diff != "" {
				t.Errorf("error: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestOrdinalDate_Validate(t *testing.T) {
	tests := []struct {
		name    string
		o       OrdinalDate
		wantErr error
	}{
		{
			name: "invalid year is less than 0",
			o: OrdinalDate{
				Year: -1,
				Day:  365,
			},
			wantErr: &DateLikeRangeError{
				Element: "year",
				Value:   -1,
				Year:    -1,
				Min:     0,
				Max:     9999,
			},
		},
		{
			name: "invalid year is more than 9999",
			o: OrdinalDate{
				Year: 10000,
				Day:  365,
			},
			wantErr: &DateLikeRangeError{
				Element: "year",
				Value:   10000,
				Year:    10000,
				Min:     0,
				Max:     9999,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.Validate()
			if err == nil {
				t.Fatal("expected error")
			}
			if diff := cmp.Diff(tt.wantErr, err); diff != "" {
				t.Errorf("error: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestDate_Date(t *testing.T) {
	want := Date{
		Year:  2020,
		Month: 10,
		Day:   1,
	}
	got := want.Date()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want, +got)\n%s", diff)
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

func TestTimeZoneRangeError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *TimeZoneRangeError
		want string
	}{
		{
			name: "valid error",
			err: &TimeZoneRangeError{
				Element: "hour",
				Value:   24,
				Min:     0,
				Max:     23,
			},
			want: "iso8601 time zone: 24 hour is not in range 0-23",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("TimeZoneRangeError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZone_Offset(t *testing.T) {
	tests := []struct {
		name string
		z    Zone
		want int
	}{
		{
			name: "Asia/Kabul",
			z: Zone{
				Hour:   4,
				Minute: 30,
				Second: 0,
			},
			want: 16200,
		},
		{
			name: "Pacific/Marquesas",
			z: Zone{
				Hour:     9,
				Minute:   30,
				Second:   0,
				Negative: true,
			},
			want: -34200,
		},
		{
			name: "Asia/Tokyo",
			z: Zone{
				Hour:   9,
				Minute: 0,
				Second: 0,
			},
			want: 32400,
		},
		{
			name: "with second",
			z: Zone{
				Hour:     13,
				Minute:   10,
				Second:   30,
				Negative: true,
			},
			want: -47430,
		},
		{
			name: "+99:99:99",
			z: Zone{
				Hour:   99,
				Minute: 99,
				Second: 99,
			},
			want: 362439,
		},
		{
			name: "-99:99:99",
			z: Zone{
				Hour:     99,
				Minute:   99,
				Second:   99,
				Negative: true,
			},
			want: -362439,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.z.Offset(); got != tt.want {
				t.Errorf("Zone.Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZone_Validate(t *testing.T) {
	tests := []struct {
		name string
		z    Zone
		want error
	}{
		{
			name: "invalid hour more than 99",
			z: Zone{
				Hour: 100,
			},
			want: &TimeZoneRangeError{
				Element: "hour",
				Value:   100,
				Min:     0,
				Max:     99,
			},
		},
		{
			name: "invalid hour less than 0",
			z: Zone{
				Hour: -1,
			},
			want: &TimeZoneRangeError{
				Element: "hour",
				Value:   -1,
				Min:     0,
				Max:     99,
			},
		},
		{
			name: "invalid minute more than 99",
			z: Zone{
				Minute: 100,
			},
			want: &TimeZoneRangeError{
				Element: "minute",
				Value:   100,
				Min:     0,
				Max:     99,
			},
		},
		{
			name: "invalid minute less than 0",
			z: Zone{
				Minute: -1,
			},
			want: &TimeZoneRangeError{
				Element: "minute",
				Value:   -1,
				Min:     0,
				Max:     99,
			},
		},
		{
			name: "invalid second more than 99",
			z: Zone{
				Second: 100,
			},
			want: &TimeZoneRangeError{
				Element: "second",
				Value:   100,
				Min:     0,
				Max:     99,
			},
		},
		{
			name: "invalid second less than 0",
			z: Zone{
				Second: -1,
			},
			want: &TimeZoneRangeError{
				Element: "second",
				Value:   -1,
				Min:     0,
				Max:     99,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.z.Validate()
			if err == nil {
				t.Fatal("expected error")
			}
			if diff := cmp.Diff(tt.want, err); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestDuration_StdDuration(t *testing.T) {
	tests := []struct {
		name string
		d    Duration
		want time.Duration
	}{
		{
			name: "8765h49m12s", // 1 year
			d: Duration{
				Year: 1,
			},
			want: yearInSecond,
		},
		{
			name: "730h33m36s", // 1 month
			d: Duration{
				Month: 1,
			},
			want: monthInSecond,
		},
		{
			name: "168h0m0s", // 1 week
			d: Duration{
				Week: 1,
			},
			want: weekInSecond,
		},
		{
			name: "24h0m0s", // 1 day
			d: Duration{
				Day: 1,
			},
			want: dayInSecond,
		},
		{
			name: "1h0m0s",
			d: Duration{
				Hour: 1,
			},
			want: time.Hour,
		},
		{
			name: "1m0s",
			d: Duration{
				Minute: 1,
			},
			want: time.Minute,
		},
		{
			name: "1s",
			d: Duration{
				Second: 1,
			},
			want: time.Second,
		},
		{
			name: "1ms",
			d: Duration{
				Millisecond: 1,
			},
			want: time.Millisecond,
		},
		{
			name: "1µs",
			d: Duration{
				Microsecond: 1,
			},
			want: time.Microsecond,
		},
		{
			name: "1ns",
			d: Duration{
				Nanosecond: 1,
			},
			want: time.Nanosecond,
		},
		{
			name: "8765h49m13s", // 1 year + 1 sec
			d: Duration{
				Year:   1,
				Second: 1,
			},
			want: yearInSecond + time.Second,
		},
		{
			name: "-730h34m36s", // 1 month + 1 minute
			d: Duration{
				Month:    1,
				Minute:   1,
				Negative: true,
			},
			want: -1 * (monthInSecond + time.Minute),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.StdDuration(); got != tt.want {
				t.Errorf("Duration.StdDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want Duration
	}{
		{
			d: yearInSecond,
			want: Duration{
				Year: 1,
			},
		},
		{
			d: monthInSecond,
			want: Duration{
				Month: 1,
			},
		},
		{
			d: weekInSecond,
			want: Duration{
				Week: 1,
			},
		},
		{
			d: dayInSecond,
			want: Duration{
				Day: 1,
			},
		},
		{
			d: time.Hour,
			want: Duration{
				Hour: 1,
			},
		},
		{
			d: time.Minute,
			want: Duration{
				Minute: 1,
			},
		},
		{
			d: time.Second,
			want: Duration{
				Second: 1,
			},
		},
		{
			d: time.Millisecond,
			want: Duration{
				Millisecond: 1,
			},
		},
		{
			d: time.Microsecond,
			want: Duration{
				Microsecond: 1,
			},
		},
		{
			d: time.Nanosecond,
			want: Duration{
				Nanosecond: 1,
			},
		},
		{
			d: yearInSecond + time.Second,
			want: Duration{
				Year:   1,
				Second: 1,
			},
		},
		{
			d: -1 * (monthInSecond + time.Minute),
			want: Duration{
				Month:    1,
				Minute:   1,
				Negative: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.d.String(), func(t *testing.T) {
			got := NewDuration(tt.d)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestDuration_Negate(t *testing.T) {
	tests := []struct {
		d    Duration
		want Duration
	}{
		{
			d: Duration{
				Year:        1234,
				Month:       2345,
				Week:        3456,
				Day:         4567,
				Hour:        5678,
				Minute:      6789,
				Second:      7890,
				Millisecond: 890,
				Microsecond: 901,
				Nanosecond:  123,
			},
			want: Duration{
				Year:        1234,
				Month:       2345,
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
			d: Duration{
				Year:        1234,
				Month:       2345,
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
			want: Duration{
				Year:        1234,
				Month:       2345,
				Week:        3456,
				Day:         4567,
				Hour:        5678,
				Minute:      6789,
				Second:      7890,
				Millisecond: 890,
				Microsecond: 901,
				Nanosecond:  123,
				Negative:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.d.String(), func(t *testing.T) {
			got := tt.d.Negate()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestDuration_String(t *testing.T) {
	tests := []struct {
		d    Duration
		want string
	}{
		{
			d:    Duration{},
			want: "PT0S",
		},
		{
			d: Duration{
				Year: 1,
			},
			want: "P1Y",
		},
		{
			d: Duration{
				Month: 1,
			},
			want: "P1M",
		},
		{
			d: Duration{
				Week: 1,
			},
			want: "P1W",
		},
		{
			d: Duration{
				Day: 1,
			},
			want: "P1D",
		},
		{
			d: Duration{
				Hour: 1,
			},
			want: "PT1H",
		},
		{
			d: Duration{
				Minute: 1,
			},
			want: "PT1M",
		},
		{
			d: Duration{
				Second: 1,
			},
			want: "PT1S",
		},
		{
			d: Duration{
				Year:  1,
				Month: 2,
				Week:  3,
				Day:   4,
			},
			want: "P1Y2M3W4D",
		},
		{
			d: Duration{
				Year:     1,
				Month:    2,
				Week:     3,
				Day:      4,
				Negative: true,
			},
			want: "-P1Y2M3W4D",
		},
		{
			d: Duration{
				Hour:   2,
				Minute: 20,
				Second: 30,
			},
			want: "PT2H20M30S",
		},
		{
			d: Duration{
				Hour:     2,
				Minute:   20,
				Second:   30,
				Negative: true,
			},
			want: "-PT2H20M30S",
		},
		{
			d: Duration{
				Year:        1234,
				Month:       2345,
				Week:        3456,
				Day:         4567,
				Hour:        5678,
				Minute:      6789,
				Second:      7890,
				Millisecond: 890,
				Microsecond: 901,
				Nanosecond:  123,
			},
			want: "P1234Y2345M3456W4567DT5678H6789M7890.890901123S",
		},
		{
			d: Duration{
				Year:        1234,
				Month:       2345,
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
			want: "-P1234Y2345M3456W4567DT5678H6789M7890.890901123S",
		},
		{
			d: Duration{
				Year:        1,
				Month:       1,
				Week:        1,
				Day:         1,
				Hour:        1,
				Minute:      1,
				Second:      1,
				Millisecond: 1,
				Microsecond: 1,
				Nanosecond:  1,
			},
			want: "P1Y1M1W1DT1H1M1.001001001S",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("Duration.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
