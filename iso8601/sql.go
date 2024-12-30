package iso8601

import (
	"fmt"
	"time"

	"database/sql/driver"
)

// Scan implements the sql.Scanner interface.
func (d *Date) Scan(src any) error {
	switch s := src.(type) {
	case nil:
		*d = Date{}
	case time.Time:
		*d = DateOf(s)
	case string:
		dl, err := ParseDate[string](s)
		if err != nil {
			return err
		}
		*d = dl.Date()
		return dl.Validate()
	case []byte:
		dl, err := ParseDate[[]byte](s)
		if err != nil {
			return err
		}
		*d = dl.Date()
		return dl.Validate()
	default:
		return fmt.Errorf("unknown type of: %T", s)
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (d *Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// NullDate represents a Date that may be null.
type NullDate struct {
	Date  Date
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the sql.Scanner interface.
func (n *NullDate) Scan(src any) error {
	if src == nil {
		n.Date, n.Valid = Date{}, false
		return nil
	}
	n.Valid = true // almost the same behavior as sql.NullTime
	return n.Date.Scan(src)
}

// Value implements the driver.Valuer interface.
func (n NullDate) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Date.Value()
}
