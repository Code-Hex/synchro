package iso8601

import (
	"fmt"
	"math"
	"time"
)

// fraction:     (\d+)(?:[.,](\d{1,9}))?
// durationTime: (?:${fraction}H)?(?:${fraction}M)?(?:${fraction}S)?
// durationDate: (?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?
// duration:     ^([+\u2212-])?P${durationDate}(?:T(?!$)${durationTime})?$
//
// NOTE: According to the ISO 8601-1 standard, weeks are not allowed to appear
// together with any other units, and durations can only be positive. As extensions
// to the standard, ISO 8601-2 allows a sign character at the start of the string, and
// allows combining weeks with other units. If you intend to use a string such as
// P3W1D, +P1M, or -P1M for interoperability, note that other programs may not accept it.
func ParseDuration(b []byte) (Duration, error) {
	var (
		y                 int
		m                 int
		w                 int
		d                 int
		hour              int
		minute            int
		second            int
		excessNanoseconds int
		negative          bool
	)

	if len(b) == 0 {
		return Duration{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b),
			AfterToken: "",
			Expected:   "P or + or -",
		}
	}

	i := 0
	switch b[0] {
	case '+', '-':
		if b[0] == '-' {
			negative = true
		}
		if len(b[1:]) == 0 {
			return Duration{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b),
				AfterToken: string(b[0]),
				Expected:   "P",
			}
		}
		i++
		fallthrough
	case 'P':
		i++
	default:
		return Duration{}, &UnexpectedTokenError{
			Value:      string(b),
			Token:      string(b),
			AfterToken: "",
			Expected:   "P or + or -",
		}
	}

	// Separate because the 'M' designator exists in both date and time.
	// dateSeen keys will be 'Y', 'M', 'D', 'W'
	dateSeen := make(map[byte]bool, 3)
	// timeSeen keys will be 'H', 'M', 'S'
	timeSeen := make(map[byte]bool, 3)

	// Because only the smallest unit can be fractional
	seenFranction := false

	// the 'T'
	seenT := false

	setter := func(idx int, seen map[byte]bool, setter func() error) error {
		if seen[b[idx]] {
			return &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[idx]),
				AfterToken: string(b[:idx]),
				Expected:   "the designator to be used only once",
			}
		}
		if err := setter(); err != nil {
			return err
		}
		seen[b[idx]] = true
		return nil
	}
	dateSetter := func(idx int, f func() error) error {
		return setter(idx, dateSeen, func() error {
			if seenT {
				return &UnexpectedTokenError{
					Value:      string(b),
					Token:      string(b[idx]),
					AfterToken: string(b[:idx]),
					Expected:   "date duration should put before the 'T'",
				}
			}
			return f()
		})
	}
	timeSetter := func(idx int, f func() error) error {
		return setter(idx, timeSeen, func() error {
			if !seenT {
				return &UnexpectedTokenError{
					Value:      string(b),
					Token:      string(b[idx]),
					AfterToken: string(b[:idx]),
					Expected:   "the 'T' designator is required",
				}
			}
			return f()
		})
	}

	for i < len(b) {
		n := countDigits(b[i:], 0)
		if n == 0 {
			if b[i] == 'T' {
				if seenT {
					return Duration{}, &UnexpectedTokenError{
						Value:      string(b),
						Token:      "T",
						AfterToken: string(b[:i]),
						Expected:   "the 'T' designator should be once",
					}
				}
				i++
				seenT = true
				continue
			}
			return Duration{}, &UnexpectedTokenError{
				Value:      string(b),
				Token:      string(b[i]),
				AfterToken: string(b[:i]),
				Expected:   "PnYnMnDTnHnMnS or PnW format",
			}
		}

		// Parse the number
		val := parseNumber(b, i, n)
		i += n

		// check fraction
		fraction := 0
		if i < len(b) && (b[i] == '.' || b[i] == ',') {
			if !seenT {
				return Duration{}, &UnexpectedTokenError{
					Value:      string(b),
					Token:      string(b[i]),
					AfterToken: string(b[:i]),
					Expected:   "only the time unit can be fractional",
				}
			}
			if seenFranction {
				return Duration{}, &UnexpectedTokenError{
					Value:      string(b),
					Token:      string(b[i]),
					AfterToken: string(b[:i]),
					Expected:   "only the smallest time unit can be fractional",
				}
			}
			i++
			fraction = parseFraction(b[i:])
			i += countDigits(b[i:], 0) // FIXME: count in also parseFraction...
			seenFranction = true
		}

		// Decode based on the current suffix
		if i < len(b) {
			switch b[i] {
			case 'Y':
				err := dateSetter(i, func() error {
					for _, designator := range []byte{'M', 'W', 'D'} {
						if dateSeen[designator] {
							return &UnexpectedTokenError{
								Value:      string(b),
								Token:      "Y",
								AfterToken: string(b[:i]),
								Expected:   fmt.Sprintf("the 'Y' date designator should appear before '%s'", string(designator)),
							}
						}
					}
					y = val
					return nil
				})
				if err != nil {
					return Duration{}, err
				}
			case 'M':
				if seenT {
					// minute
					err := timeSetter(i, func() error {
						if timeSeen['S'] {
							return &UnexpectedTokenError{
								Value:      string(b),
								Token:      "M",
								AfterToken: string(b[:i]),
								Expected:   "the 'M' time designator should appear before 'S'",
							}
						}
						minute = val
						return nil
					})
					if err != nil {
						return Duration{}, err
					}
					if fraction > 0 {
						excessNanoseconds = fraction * 60
					}
				} else {
					// month
					err := dateSetter(i, func() error {
						for _, designator := range []byte{'W', 'D'} {
							if dateSeen[designator] {
								return &UnexpectedTokenError{
									Value:      string(b),
									Token:      "M",
									AfterToken: string(b[:i]),
									Expected:   fmt.Sprintf("the 'M' date designator should appear before '%s'", string(designator)),
								}
							}
						}
						m = val
						return nil
					})
					if err != nil {
						return Duration{}, err
					}
				}
			case 'W':
				err := dateSetter(i, func() error {
					if dateSeen['D'] {
						return &UnexpectedTokenError{
							Value:      string(b),
							Token:      "W",
							AfterToken: string(b[:i]),
							Expected:   "the 'W' date designator should appear before 'D'",
						}
					}
					w = val
					return nil
				})
				if err != nil {
					return Duration{}, err
				}
			case 'D':
				err := dateSetter(i, func() error {
					d = val
					return nil
				})
				if err != nil {
					return Duration{}, err
				}
			case 'H':
				err := timeSetter(i, func() error {
					for _, designator := range []byte{'M', 'S'} {
						if timeSeen[designator] {
							return &UnexpectedTokenError{
								Value:      string(b),
								Token:      "H",
								AfterToken: string(b[:i]),
								Expected:   fmt.Sprintf("the 'H' time designator should appear before '%s'", string(designator)),
							}
						}
					}
					hour = val
					return nil
				})
				if err != nil {
					return Duration{}, err
				}
				if fraction > 0 {
					excessNanoseconds = fraction * 3600
				}
			case 'S':
				err := timeSetter(i, func() error {
					second = val
					return nil
				})
				if err != nil {
					return Duration{}, err
				}
				if fraction > 0 {
					excessNanoseconds = fraction
				}
			default:
				return Duration{}, &UnexpectedTokenError{
					Value:      string(b),
					Token:      string(b[i]),
					AfterToken: string(b[:i]),
					Expected:   "PnYnMnDTnHnMnS or PnW format",
				}
			}
			i++
		}
	}

	// https://github.com/tc39/proposal-temporal/blob/4a5e88e2334c8ab428590bfd20a706d1bb6c6031/polyfill/lib/ecmascript.mjs#L635
	nanosec := excessNanoseconds % 1000
	microsec := int(math.Trunc(float64(excessNanoseconds)/1000)) % 1000
	millisec := int(math.Trunc(float64(excessNanoseconds)/1e6)) % 1000
	second += int(math.Trunc(float64(excessNanoseconds)/1e9)) % 60
	minute += int(math.Trunc(float64(excessNanoseconds) / 60e9))

	return Duration{
		Year:        y,
		Month:       time.Month(m),
		Week:        w,
		Day:         d,
		Hour:        hour,
		Minute:      minute,
		Second:      second,
		Millisecond: millisec,
		Microsecond: microsec,
		Nanosecond:  nanosec,
		Negative:    negative,
	}, nil
}
