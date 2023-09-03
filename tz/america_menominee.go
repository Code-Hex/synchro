// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaMenomineeLocation  sync.Once
	cacheAmericaMenomineeLocation *time.Location
)

type AmericaMenominee struct{}

func (AmericaMenominee) Location() *time.Location {
	onceAmericaMenomineeLocation.Do(func() {
		loc, err := time.LoadLocation("America/Menominee")
		if err != nil {
			panic(err)
		}
		cacheAmericaMenomineeLocation = loc
	})
	return cacheAmericaMenomineeLocation
}