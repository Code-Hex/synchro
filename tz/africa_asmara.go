// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAfricaAsmaraLocation  sync.Once
	cacheAfricaAsmaraLocation *time.Location
)

type AfricaAsmara struct{}

func (AfricaAsmara) Location() *time.Location {
	onceAfricaAsmaraLocation.Do(func() {
		loc, err := time.LoadLocation("Africa/Asmara")
		if err != nil {
			panic(err)
		}
		cacheAfricaAsmaraLocation = loc
	})
	return cacheAfricaAsmaraLocation
}
