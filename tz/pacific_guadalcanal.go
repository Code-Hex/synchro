// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	oncePacificGuadalcanalLocation  sync.Once
	cachePacificGuadalcanalLocation *time.Location
)

type PacificGuadalcanal struct{}

func (PacificGuadalcanal) Location() *time.Location {
	oncePacificGuadalcanalLocation.Do(func() {
		loc, err := time.LoadLocation("Pacific/Guadalcanal")
		if err != nil {
			panic(err)
		}
		cachePacificGuadalcanalLocation = loc
	})
	return cachePacificGuadalcanalLocation
}
