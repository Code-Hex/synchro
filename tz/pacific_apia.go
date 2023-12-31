// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	oncePacificApiaLocation  sync.Once
	cachePacificApiaLocation *time.Location
)

type PacificApia struct{}

func (PacificApia) Location() *time.Location {
	oncePacificApiaLocation.Do(func() {
		loc, err := time.LoadLocation("Pacific/Apia")
		if err != nil {
			panic(err)
		}
		cachePacificApiaLocation = loc
	})
	return cachePacificApiaLocation
}
