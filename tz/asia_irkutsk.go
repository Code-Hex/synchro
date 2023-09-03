// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaIrkutskLocation  sync.Once
	cacheAsiaIrkutskLocation *time.Location
)

type AsiaIrkutsk struct{}

func (AsiaIrkutsk) Location() *time.Location {
	onceAsiaIrkutskLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Irkutsk")
		if err != nil {
			panic(err)
		}
		cacheAsiaIrkutskLocation = loc
	})
	return cacheAsiaIrkutskLocation
}