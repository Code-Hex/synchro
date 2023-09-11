package iso8601

import (
	"time"

	_ "unsafe"
)

//go:linkname lookup time.(*Location).lookup
func lookup(loc *time.Location, sec int64) (name string, offset int, start, end int64, isDST bool)
