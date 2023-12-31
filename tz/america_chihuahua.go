// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaChihuahuaLocation  sync.Once
	cacheAmericaChihuahuaLocation *time.Location
)

type AmericaChihuahua struct{}

func (AmericaChihuahua) Location() *time.Location {
	onceAmericaChihuahuaLocation.Do(func() {
		loc, err := time.LoadLocation("America/Chihuahua")
		if err != nil {
			panic(err)
		}
		cacheAmericaChihuahuaLocation = loc
	})
	return cacheAmericaChihuahuaLocation
}
