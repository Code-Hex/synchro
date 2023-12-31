// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceEuropeOsloLocation  sync.Once
	cacheEuropeOsloLocation *time.Location
)

type EuropeOslo struct{}

func (EuropeOslo) Location() *time.Location {
	onceEuropeOsloLocation.Do(func() {
		loc, err := time.LoadLocation("Europe/Oslo")
		if err != nil {
			panic(err)
		}
		cacheEuropeOsloLocation = loc
	})
	return cacheEuropeOsloLocation
}
