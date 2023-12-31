// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaMatamorosLocation  sync.Once
	cacheAmericaMatamorosLocation *time.Location
)

type AmericaMatamoros struct{}

func (AmericaMatamoros) Location() *time.Location {
	onceAmericaMatamorosLocation.Do(func() {
		loc, err := time.LoadLocation("America/Matamoros")
		if err != nil {
			panic(err)
		}
		cacheAmericaMatamorosLocation = loc
	})
	return cacheAmericaMatamorosLocation
}
