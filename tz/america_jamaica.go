// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaJamaicaLocation  sync.Once
	cacheAmericaJamaicaLocation *time.Location
)

type AmericaJamaica struct{}

func (AmericaJamaica) Location() *time.Location {
	onceAmericaJamaicaLocation.Do(func() {
		loc, err := time.LoadLocation("America/Jamaica")
		if err != nil {
			panic(err)
		}
		cacheAmericaJamaicaLocation = loc
	})
	return cacheAmericaJamaicaLocation
}
