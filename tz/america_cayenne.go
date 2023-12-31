// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaCayenneLocation  sync.Once
	cacheAmericaCayenneLocation *time.Location
)

type AmericaCayenne struct{}

func (AmericaCayenne) Location() *time.Location {
	onceAmericaCayenneLocation.Do(func() {
		loc, err := time.LoadLocation("America/Cayenne")
		if err != nil {
			panic(err)
		}
		cacheAmericaCayenneLocation = loc
	})
	return cacheAmericaCayenneLocation
}
