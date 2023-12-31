// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaPorto_VelhoLocation  sync.Once
	cacheAmericaPorto_VelhoLocation *time.Location
)

type AmericaPorto_Velho struct{}

func (AmericaPorto_Velho) Location() *time.Location {
	onceAmericaPorto_VelhoLocation.Do(func() {
		loc, err := time.LoadLocation("America/Porto_Velho")
		if err != nil {
			panic(err)
		}
		cacheAmericaPorto_VelhoLocation = loc
	})
	return cacheAmericaPorto_VelhoLocation
}
