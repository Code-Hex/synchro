// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaMontserratLocation  sync.Once
	cacheAmericaMontserratLocation *time.Location
)

type AmericaMontserrat struct{}

func (AmericaMontserrat) Location() *time.Location {
	onceAmericaMontserratLocation.Do(func() {
		loc, err := time.LoadLocation("America/Montserrat")
		if err != nil {
			panic(err)
		}
		cacheAmericaMontserratLocation = loc
	})
	return cacheAmericaMontserratLocation
}