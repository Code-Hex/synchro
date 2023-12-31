// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaYangonLocation  sync.Once
	cacheAsiaYangonLocation *time.Location
)

type AsiaYangon struct{}

func (AsiaYangon) Location() *time.Location {
	onceAsiaYangonLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Yangon")
		if err != nil {
			panic(err)
		}
		cacheAsiaYangonLocation = loc
	})
	return cacheAsiaYangonLocation
}
