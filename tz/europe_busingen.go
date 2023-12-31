// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceEuropeBusingenLocation  sync.Once
	cacheEuropeBusingenLocation *time.Location
)

type EuropeBusingen struct{}

func (EuropeBusingen) Location() *time.Location {
	onceEuropeBusingenLocation.Do(func() {
		loc, err := time.LoadLocation("Europe/Busingen")
		if err != nil {
			panic(err)
		}
		cacheEuropeBusingenLocation = loc
	})
	return cacheEuropeBusingenLocation
}
