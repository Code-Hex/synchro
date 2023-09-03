// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceIndianMaldivesLocation  sync.Once
	cacheIndianMaldivesLocation *time.Location
)

type IndianMaldives struct{}

func (IndianMaldives) Location() *time.Location {
	onceIndianMaldivesLocation.Do(func() {
		loc, err := time.LoadLocation("Indian/Maldives")
		if err != nil {
			panic(err)
		}
		cacheIndianMaldivesLocation = loc
	})
	return cacheIndianMaldivesLocation
}