// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceEuropeKirovLocation  sync.Once
	cacheEuropeKirovLocation *time.Location
)

type EuropeKirov struct{}

func (EuropeKirov) Location() *time.Location {
	onceEuropeKirovLocation.Do(func() {
		loc, err := time.LoadLocation("Europe/Kirov")
		if err != nil {
			panic(err)
		}
		cacheEuropeKirovLocation = loc
	})
	return cacheEuropeKirovLocation
}
