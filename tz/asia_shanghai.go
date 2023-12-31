// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaShanghaiLocation  sync.Once
	cacheAsiaShanghaiLocation *time.Location
)

type AsiaShanghai struct{}

func (AsiaShanghai) Location() *time.Location {
	onceAsiaShanghaiLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			panic(err)
		}
		cacheAsiaShanghaiLocation = loc
	})
	return cacheAsiaShanghaiLocation
}
