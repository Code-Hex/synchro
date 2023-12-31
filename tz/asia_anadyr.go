// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaAnadyrLocation  sync.Once
	cacheAsiaAnadyrLocation *time.Location
)

type AsiaAnadyr struct{}

func (AsiaAnadyr) Location() *time.Location {
	onceAsiaAnadyrLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Anadyr")
		if err != nil {
			panic(err)
		}
		cacheAsiaAnadyrLocation = loc
	})
	return cacheAsiaAnadyrLocation
}
