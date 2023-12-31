// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaAqtobeLocation  sync.Once
	cacheAsiaAqtobeLocation *time.Location
)

type AsiaAqtobe struct{}

func (AsiaAqtobe) Location() *time.Location {
	onceAsiaAqtobeLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Aqtobe")
		if err != nil {
			panic(err)
		}
		cacheAsiaAqtobeLocation = loc
	})
	return cacheAsiaAqtobeLocation
}
