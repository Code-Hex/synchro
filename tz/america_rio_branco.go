// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaRio_BrancoLocation  sync.Once
	cacheAmericaRio_BrancoLocation *time.Location
)

type AmericaRio_Branco struct{}

func (AmericaRio_Branco) Location() *time.Location {
	onceAmericaRio_BrancoLocation.Do(func() {
		loc, err := time.LoadLocation("America/Rio_Branco")
		if err != nil {
			panic(err)
		}
		cacheAmericaRio_BrancoLocation = loc
	})
	return cacheAmericaRio_BrancoLocation
}