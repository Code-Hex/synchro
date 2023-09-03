// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceEuropeVaticanLocation  sync.Once
	cacheEuropeVaticanLocation *time.Location
)

type EuropeVatican struct{}

func (EuropeVatican) Location() *time.Location {
	onceEuropeVaticanLocation.Do(func() {
		loc, err := time.LoadLocation("Europe/Vatican")
		if err != nil {
			panic(err)
		}
		cacheEuropeVaticanLocation = loc
	})
	return cacheEuropeVaticanLocation
}