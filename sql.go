package synchro

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/Code-Hex/synchro/iso8601"
)

// Scan implements the sql.Scanner interface.
func (t *Time[T]) Scan(src any) error {
	if src == nil {
		*t = Time[T]{} // zero value
		return nil
	}
	var tz T
	switch s := src.(type) {
	case time.Time:
		*t = In[T](s)
		return nil
	case string:
		parsed, err := iso8601.ParseDateTime[string](
			s,
			iso8601.WithTimeDesignators(' '),
			iso8601.WithInLocation(tz.Location()),
		)
		if err != nil {
			return err
		}
		*t = In[T](parsed)
		return nil
	case []byte:
		parsed, err := iso8601.ParseDateTime[[]byte](
			s,
			iso8601.WithTimeDesignators(' '),
			iso8601.WithInLocation(tz.Location()),
		)
		if err != nil {
			return err
		}
		*t = In[T](parsed)
		return nil
	default:
		return fmt.Errorf("unknown type of: %T", s)
	}
}

// Value implements the driver.Valuer interface.
func (t Time[T]) Value() (driver.Value, error) {
	return t.tm, nil
}

// NullTime represents a Time[T] that may be null.
// NullTime implements the sql.Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.
type NullTime[T TimeZone] struct {
	Time  Time[T]
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the sql.Scanner interface.
func (n *NullTime[T]) Scan(src any) error {
	if src == nil {
		n.Time, n.Valid = Time[T]{}, false
		return nil
	}
	n.Valid = true // almost the same behavior as sql.NullTime
	return n.Time.Scan(src)
}

// Value implements the driver.Valuer interface.
func (n NullTime[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time.tm, nil
}
