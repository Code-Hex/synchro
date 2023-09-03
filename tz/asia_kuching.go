// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaKuchingLocation  sync.Once
	cacheAsiaKuchingLocation *time.Location
)

type AsiaKuching struct{}

func (AsiaKuching) Location() *time.Location {
	onceAsiaKuchingLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Kuching")
		if err != nil {
			panic(err)
		}
		cacheAsiaKuchingLocation = loc
	})
	return cacheAsiaKuchingLocation
}