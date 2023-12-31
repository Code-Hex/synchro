// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceEuropeSaratovLocation  sync.Once
	cacheEuropeSaratovLocation *time.Location
)

type EuropeSaratov struct{}

func (EuropeSaratov) Location() *time.Location {
	onceEuropeSaratovLocation.Do(func() {
		loc, err := time.LoadLocation("Europe/Saratov")
		if err != nil {
			panic(err)
		}
		cacheEuropeSaratovLocation = loc
	})
	return cacheEuropeSaratovLocation
}
