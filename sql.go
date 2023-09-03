package synchro

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/Code-Hex/synchro/tz"
)

var _ interface {
	sql.Scanner
	driver.Valuer
} = (*Time[tz.UTC])(nil)

// Scan implements the sql.Scanner interface.
func (t *Time[T]) Scan(src any) error {
	if src == nil {
		*t = Time[T]{} // zero value
		return nil
	}
	switch s := src.(type) {
	case time.Time:
		*t = In[T](s)
		return nil
	default:
		return fmt.Errorf("unknown type of: %T", s)
	}
}

// Value implements the driver.Valuer interface.
func (t Time[T]) Value() (driver.Value, error) {
	return t.tm, nil
}

var _ interface {
	sql.Scanner
	driver.Valuer
} = (*NullTime[tz.UTC])(nil)

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
