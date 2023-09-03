// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaArgentinaLa_RiojaLocation  sync.Once
	cacheAmericaArgentinaLa_RiojaLocation *time.Location
)

type AmericaArgentinaLa_Rioja struct{}

func (AmericaArgentinaLa_Rioja) Location() *time.Location {
	onceAmericaArgentinaLa_RiojaLocation.Do(func() {
		loc, err := time.LoadLocation("America/Argentina/La_Rioja")
		if err != nil {
			panic(err)
		}
		cacheAmericaArgentinaLa_RiojaLocation = loc
	})
	return cacheAmericaArgentinaLa_RiojaLocation
}