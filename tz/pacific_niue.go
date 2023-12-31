// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	oncePacificNiueLocation  sync.Once
	cachePacificNiueLocation *time.Location
)

type PacificNiue struct{}

func (PacificNiue) Location() *time.Location {
	oncePacificNiueLocation.Do(func() {
		loc, err := time.LoadLocation("Pacific/Niue")
		if err != nil {
			panic(err)
		}
		cachePacificNiueLocation = loc
	})
	return cachePacificNiueLocation
}
