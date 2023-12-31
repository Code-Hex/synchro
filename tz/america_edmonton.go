// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaEdmontonLocation  sync.Once
	cacheAmericaEdmontonLocation *time.Location
)

type AmericaEdmonton struct{}

func (AmericaEdmonton) Location() *time.Location {
	onceAmericaEdmontonLocation.Do(func() {
		loc, err := time.LoadLocation("America/Edmonton")
		if err != nil {
			panic(err)
		}
		cacheAmericaEdmontonLocation = loc
	})
	return cacheAmericaEdmontonLocation
}
