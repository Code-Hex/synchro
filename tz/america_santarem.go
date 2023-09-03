// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaSantaremLocation  sync.Once
	cacheAmericaSantaremLocation *time.Location
)

type AmericaSantarem struct{}

func (AmericaSantarem) Location() *time.Location {
	onceAmericaSantaremLocation.Do(func() {
		loc, err := time.LoadLocation("America/Santarem")
		if err != nil {
			panic(err)
		}
		cacheAmericaSantaremLocation = loc
	})
	return cacheAmericaSantaremLocation
}