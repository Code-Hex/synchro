// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaArubaLocation  sync.Once
	cacheAmericaArubaLocation *time.Location
)

type AmericaAruba struct{}

func (AmericaAruba) Location() *time.Location {
	onceAmericaArubaLocation.Do(func() {
		loc, err := time.LoadLocation("America/Aruba")
		if err != nil {
			panic(err)
		}
		cacheAmericaArubaLocation = loc
	})
	return cacheAmericaArubaLocation
}
