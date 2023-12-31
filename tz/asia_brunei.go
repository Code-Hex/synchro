// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaBruneiLocation  sync.Once
	cacheAsiaBruneiLocation *time.Location
)

type AsiaBrunei struct{}

func (AsiaBrunei) Location() *time.Location {
	onceAsiaBruneiLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Brunei")
		if err != nil {
			panic(err)
		}
		cacheAsiaBruneiLocation = loc
	})
	return cacheAsiaBruneiLocation
}
