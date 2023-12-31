// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaAnchorageLocation  sync.Once
	cacheAmericaAnchorageLocation *time.Location
)

type AmericaAnchorage struct{}

func (AmericaAnchorage) Location() *time.Location {
	onceAmericaAnchorageLocation.Do(func() {
		loc, err := time.LoadLocation("America/Anchorage")
		if err != nil {
			panic(err)
		}
		cacheAmericaAnchorageLocation = loc
	})
	return cacheAmericaAnchorageLocation
}
