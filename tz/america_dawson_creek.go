// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaDawson_CreekLocation  sync.Once
	cacheAmericaDawson_CreekLocation *time.Location
)

type AmericaDawson_Creek struct{}

func (AmericaDawson_Creek) Location() *time.Location {
	onceAmericaDawson_CreekLocation.Do(func() {
		loc, err := time.LoadLocation("America/Dawson_Creek")
		if err != nil {
			panic(err)
		}
		cacheAmericaDawson_CreekLocation = loc
	})
	return cacheAmericaDawson_CreekLocation
}