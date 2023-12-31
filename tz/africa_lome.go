// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAfricaLomeLocation  sync.Once
	cacheAfricaLomeLocation *time.Location
)

type AfricaLome struct{}

func (AfricaLome) Location() *time.Location {
	onceAfricaLomeLocation.Do(func() {
		loc, err := time.LoadLocation("Africa/Lome")
		if err != nil {
			panic(err)
		}
		cacheAfricaLomeLocation = loc
	})
	return cacheAfricaLomeLocation
}
