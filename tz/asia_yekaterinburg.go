// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaYekaterinburgLocation  sync.Once
	cacheAsiaYekaterinburgLocation *time.Location
)

type AsiaYekaterinburg struct{}

func (AsiaYekaterinburg) Location() *time.Location {
	onceAsiaYekaterinburgLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Yekaterinburg")
		if err != nil {
			panic(err)
		}
		cacheAsiaYekaterinburgLocation = loc
	})
	return cacheAsiaYekaterinburgLocation
}