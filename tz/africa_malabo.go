// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAfricaMalaboLocation  sync.Once
	cacheAfricaMalaboLocation *time.Location
)

type AfricaMalabo struct{}

func (AfricaMalabo) Location() *time.Location {
	onceAfricaMalaboLocation.Do(func() {
		loc, err := time.LoadLocation("Africa/Malabo")
		if err != nil {
			panic(err)
		}
		cacheAfricaMalaboLocation = loc
	})
	return cacheAfricaMalaboLocation
}
