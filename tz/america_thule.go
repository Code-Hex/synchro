// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaThuleLocation  sync.Once
	cacheAmericaThuleLocation *time.Location
)

type AmericaThule struct{}

func (AmericaThule) Location() *time.Location {
	onceAmericaThuleLocation.Do(func() {
		loc, err := time.LoadLocation("America/Thule")
		if err != nil {
			panic(err)
		}
		cacheAmericaThuleLocation = loc
	})
	return cacheAmericaThuleLocation
}
