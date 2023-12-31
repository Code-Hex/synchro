// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaAshgabatLocation  sync.Once
	cacheAsiaAshgabatLocation *time.Location
)

type AsiaAshgabat struct{}

func (AsiaAshgabat) Location() *time.Location {
	onceAsiaAshgabatLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Ashgabat")
		if err != nil {
			panic(err)
		}
		cacheAsiaAshgabatLocation = loc
	})
	return cacheAsiaAshgabatLocation
}
