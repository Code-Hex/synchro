// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceIndianMayotteLocation  sync.Once
	cacheIndianMayotteLocation *time.Location
)

type IndianMayotte struct{}

func (IndianMayotte) Location() *time.Location {
	onceIndianMayotteLocation.Do(func() {
		loc, err := time.LoadLocation("Indian/Mayotte")
		if err != nil {
			panic(err)
		}
		cacheIndianMayotteLocation = loc
	})
	return cacheIndianMayotteLocation
}
