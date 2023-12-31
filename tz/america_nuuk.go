// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaNuukLocation  sync.Once
	cacheAmericaNuukLocation *time.Location
)

type AmericaNuuk struct{}

func (AmericaNuuk) Location() *time.Location {
	onceAmericaNuukLocation.Do(func() {
		loc, err := time.LoadLocation("America/Nuuk")
		if err != nil {
			panic(err)
		}
		cacheAmericaNuukLocation = loc
	})
	return cacheAmericaNuukLocation
}
