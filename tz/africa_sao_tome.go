// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAfricaSao_TomeLocation  sync.Once
	cacheAfricaSao_TomeLocation *time.Location
)

type AfricaSao_Tome struct{}

func (AfricaSao_Tome) Location() *time.Location {
	onceAfricaSao_TomeLocation.Do(func() {
		loc, err := time.LoadLocation("Africa/Sao_Tome")
		if err != nil {
			panic(err)
		}
		cacheAfricaSao_TomeLocation = loc
	})
	return cacheAfricaSao_TomeLocation
}
