package synchro

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Code-Hex/synchro/internal/constraints"
)

type Period[T TimeZone] struct {
	start    Time[T]
	end      Time[T]
	duration time.Duration
}

type makablePeriod[T TimeZone] interface {
	Time[T] | time.Time | constraints.Bytes
}

var stringType = reflect.TypeOf("")

func CreatePeriod[T TimeZone, T1 makablePeriod[T], T2 makablePeriod[T]](from T1, to T2) (Period[T], error) {
	start, err := convertTime[T](any(from))
	if err != nil {
		return Period[T]{}, fmt.Errorf("failed to parse from: %w", err)
	}
	end, err := convertTime[T](any(to))
	if err != nil {
		return Period[T]{}, fmt.Errorf("failed to parse to: %w", err)
	}
	return Period[T]{
		start:    start,
		end:      end,
		duration: 24 * time.Hour,
	}, nil
}

func (p Period[T]) Iterator() *periodIterator[T] {
	return &periodIterator[T]{
		period: p,
	}
}

func (p Period[T]) Slice() (s []Time[T]) {
	iter := p.Iterator()
	for iter.Next() {
		s = append(s, iter.Get())
	}
	return s
}

type periodIterator[T TimeZone] struct {
	current Time[T]
	period  Period[T]
}

func (iter *periodIterator[T]) Next() bool {
	if iter.current.IsZero() {
		iter.current = iter.period.start
		return true
	}
	iter.current = iter.current.Add(iter.period.duration)
	return iter.current.Compare(iter.period.end) <= 0
}

func (iter *periodIterator[T]) Get() Time[T] {
	return iter.current
}

func convertTime[T TimeZone](arg any) (Time[T], error) {
	switch v := arg.(type) {
	case Time[T]:
		return v, nil
	case time.Time:
		return In[T](v), nil
	case string:
		return ParseISO[T](v)
	case []byte:
		return ParseISO[T](string(v))
	default:
		rv := reflect.ValueOf(v)
		if rv.CanConvert(stringType) {
			return ParseISO[T](rv.Convert(stringType).String())
		}
		panic("unreachable")
	}
}
