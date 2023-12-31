// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaParamariboLocation  sync.Once
	cacheAmericaParamariboLocation *time.Location
)

type AmericaParamaribo struct{}

func (AmericaParamaribo) Location() *time.Location {
	onceAmericaParamariboLocation.Do(func() {
		loc, err := time.LoadLocation("America/Paramaribo")
		if err != nil {
			panic(err)
		}
		cacheAmericaParamariboLocation = loc
	})
	return cacheAmericaParamariboLocation
}
