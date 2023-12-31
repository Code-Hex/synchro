// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAustraliaHobartLocation  sync.Once
	cacheAustraliaHobartLocation *time.Location
)

type AustraliaHobart struct{}

func (AustraliaHobart) Location() *time.Location {
	onceAustraliaHobartLocation.Do(func() {
		loc, err := time.LoadLocation("Australia/Hobart")
		if err != nil {
			panic(err)
		}
		cacheAustraliaHobartLocation = loc
	})
	return cacheAustraliaHobartLocation
}
