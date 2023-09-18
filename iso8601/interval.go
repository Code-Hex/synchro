package iso8601

import (
	"bytes"
	"time"
)

type Interval struct {
	Start    time.Time
	End      time.Time
	Duration Duration
	Repeat   int
}

func ParseInterval[bytes []byte | ~string](b bytes) (Interval, error) {
	return parseInterval([]byte(b))
}

func parseInterval(b []byte) (Interval, error) {
	var (
		designator []byte
		start      time.Time
		end        time.Time
		duration   Duration
		repeat     int
	)
	if len(b) == 0 {
		return Interval{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b),
			AfterToken: "",
			Expected:   "R or P or datetime",
		}
	}

	if b[0] == 'R' {
		designatorIdx := -1
		// https://go.dev/play/p/42u_xsxugSW
		for _, candidate := range [][]byte{
			{'/'},
			{'-', '-'},
		} {
			n := bytes.Index(b, candidate)
			if n >= 0 && (n < designatorIdx || designatorIdx == -1) {
				designator = candidate
				designatorIdx = n
			}
		}
		if designatorIdx == -1 {
			return Interval{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[1:]),
				AfterToken: "R",
				Expected:   `internal designator "/" or "--"`,
			}
		}
		c := countDigits(b, 1)
		if want := designatorIdx - 1; c != want {
			return Interval{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[1:designatorIdx]),
				AfterToken: "R",
				Expected:   humanizeDigits(want),
			}
		}
		if c == 0 {
			repeat = -1 // infinity
		} else {
			repeat = parseNumber(b, 1, c)
		}
		b = b[designatorIdx+len(designator):]
	}

	designatorIdx := -1
	if designator != nil {
		n := bytes.Index(b, designator)
		if n >= 0 {
			designatorIdx = n
		}
	} else {
		for _, candidate := range [][]byte{
			{'/'},
			{'-', '-'},
		} {
			n := bytes.Index(b, candidate)
			if n >= 0 {
				designator = candidate
				designatorIdx = n
				break
			}
		}
	}
	// try to parse duration only
	if designatorIdx == -1 {
		d, err := parseDuration(b)
		if err != nil {
			return Interval{}, err
		}
		return Interval{
			Duration: d,
			Repeat:   repeat,
		}, nil
	}

	if b[0] == 'P' {
		d, err := parseDuration(b[:designatorIdx])
		if err != nil {
			return Interval{}, err
		}
		duration = d
	} else {
		dt, err := parseDateTime(b[:designatorIdx])
		if err != nil {
			return Interval{}, err
		}
		start = dt
	}

	endb := b[len(designator)+designatorIdx:]
	if endb[0] == 'P' {
		zero := Duration{}
		if duration != zero {
			return Interval{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(endb),
				AfterToken: string(designator),
				Expected:   "datetime format",
			}
		}
		d, err := parseDuration(endb)
		if err != nil {
			return Interval{}, err
		}
		duration = d
	} else {
		dt, err := parseDateTime(endb)
		if err != nil {
			return Interval{}, err
		}
		end = dt
	}

	return Interval{
		Start:    start,
		End:      end,
		Duration: duration,
		Repeat:   repeat,
	}, nil
}
