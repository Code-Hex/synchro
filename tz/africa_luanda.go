// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAfricaLuandaLocation  sync.Once
	cacheAfricaLuandaLocation *time.Location
)

type AfricaLuanda struct{}

func (AfricaLuanda) Location() *time.Location {
	onceAfricaLuandaLocation.Do(func() {
		loc, err := time.LoadLocation("Africa/Luanda")
		if err != nil {
			panic(err)
		}
		cacheAfricaLuandaLocation = loc
	})
	return cacheAfricaLuandaLocation
}