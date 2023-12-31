// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceEuropeLuxembourgLocation  sync.Once
	cacheEuropeLuxembourgLocation *time.Location
)

type EuropeLuxembourg struct{}

func (EuropeLuxembourg) Location() *time.Location {
	onceEuropeLuxembourgLocation.Do(func() {
		loc, err := time.LoadLocation("Europe/Luxembourg")
		if err != nil {
			panic(err)
		}
		cacheEuropeLuxembourgLocation = loc
	})
	return cacheEuropeLuxembourgLocation
}
