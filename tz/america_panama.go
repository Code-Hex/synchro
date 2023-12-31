// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaPanamaLocation  sync.Once
	cacheAmericaPanamaLocation *time.Location
)

type AmericaPanama struct{}

func (AmericaPanama) Location() *time.Location {
	onceAmericaPanamaLocation.Do(func() {
		loc, err := time.LoadLocation("America/Panama")
		if err != nil {
			panic(err)
		}
		cacheAmericaPanamaLocation = loc
	})
	return cacheAmericaPanamaLocation
}
