// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaBeirutLocation  sync.Once
	cacheAsiaBeirutLocation *time.Location
)

type AsiaBeirut struct{}

func (AsiaBeirut) Location() *time.Location {
	onceAsiaBeirutLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Beirut")
		if err != nil {
			panic(err)
		}
		cacheAsiaBeirutLocation = loc
	})
	return cacheAsiaBeirutLocation
}
