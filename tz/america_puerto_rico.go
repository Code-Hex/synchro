// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaPuerto_RicoLocation  sync.Once
	cacheAmericaPuerto_RicoLocation *time.Location
)

type AmericaPuerto_Rico struct{}

func (AmericaPuerto_Rico) Location() *time.Location {
	onceAmericaPuerto_RicoLocation.Do(func() {
		loc, err := time.LoadLocation("America/Puerto_Rico")
		if err != nil {
			panic(err)
		}
		cacheAmericaPuerto_RicoLocation = loc
	})
	return cacheAmericaPuerto_RicoLocation
}