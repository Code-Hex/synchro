// Code generated by tzgen. DO NOT EDIT.

package tz

import "time"
import "sync"

var (
	onceAsiaTashkentLocation  sync.Once
	cacheAsiaTashkentLocation *time.Location
)

type AsiaTashkent struct{}

func (AsiaTashkent) Location() *time.Location {
	onceAsiaTashkentLocation.Do(func() {
		loc, err := time.LoadLocation("Asia/Tashkent")
		if err != nil {
			panic(err)
		}
		cacheAsiaTashkentLocation = loc
	})
	return cacheAsiaTashkentLocation
}
