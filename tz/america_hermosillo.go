// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaHermosilloLocation  sync.Once
	cacheAmericaHermosilloLocation *time.Location
)

type AmericaHermosillo struct{}

func (AmericaHermosillo) Location() *time.Location {
	onceAmericaHermosilloLocation.Do(func() {
		loc, err := time.LoadLocation("America/Hermosillo")
		if err != nil {
			panic(err)
		}
		cacheAmericaHermosilloLocation = loc
	})
	return cacheAmericaHermosilloLocation
}