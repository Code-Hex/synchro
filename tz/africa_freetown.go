// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAfricaFreetownLocation  sync.Once
	cacheAfricaFreetownLocation *time.Location
)

type AfricaFreetown struct{}

func (AfricaFreetown) Location() *time.Location {
	onceAfricaFreetownLocation.Do(func() {
		loc, err := time.LoadLocation("Africa/Freetown")
		if err != nil {
			panic(err)
		}
		cacheAfricaFreetownLocation = loc
	})
	return cacheAfricaFreetownLocation
}
