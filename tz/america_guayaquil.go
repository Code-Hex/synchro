// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaGuayaquilLocation  sync.Once
	cacheAmericaGuayaquilLocation *time.Location
)

type AmericaGuayaquil struct{}

func (AmericaGuayaquil) Location() *time.Location {
	onceAmericaGuayaquilLocation.Do(func() {
		loc, err := time.LoadLocation("America/Guayaquil")
		if err != nil {
			panic(err)
		}
		cacheAmericaGuayaquilLocation = loc
	})
	return cacheAmericaGuayaquilLocation
}