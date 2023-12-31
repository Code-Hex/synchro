// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAustraliaPerthLocation  sync.Once
	cacheAustraliaPerthLocation *time.Location
)

type AustraliaPerth struct{}

func (AustraliaPerth) Location() *time.Location {
	onceAustraliaPerthLocation.Do(func() {
		loc, err := time.LoadLocation("Australia/Perth")
		if err != nil {
			panic(err)
		}
		cacheAustraliaPerthLocation = loc
	})
	return cacheAustraliaPerthLocation
}
