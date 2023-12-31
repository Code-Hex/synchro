// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAmericaCampo_GrandeLocation  sync.Once
	cacheAmericaCampo_GrandeLocation *time.Location
)

type AmericaCampo_Grande struct{}

func (AmericaCampo_Grande) Location() *time.Location {
	onceAmericaCampo_GrandeLocation.Do(func() {
		loc, err := time.LoadLocation("America/Campo_Grande")
		if err != nil {
			panic(err)
		}
		cacheAmericaCampo_GrandeLocation = loc
	})
	return cacheAmericaCampo_GrandeLocation
}
