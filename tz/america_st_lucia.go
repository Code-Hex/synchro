// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaSt_LuciaLocation  sync.Once
	cacheAmericaSt_LuciaLocation *time.Location
)

type AmericaSt_Lucia struct{}

func (AmericaSt_Lucia) Location() *time.Location {
	onceAmericaSt_LuciaLocation.Do(func() {
		loc, err := time.LoadLocation("America/St_Lucia")
		if err != nil {
			panic(err)
		}
		cacheAmericaSt_LuciaLocation = loc
	})
	return cacheAmericaSt_LuciaLocation
}
