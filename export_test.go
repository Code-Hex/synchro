package synchro

import "time"

func SetNow(f func() time.Time) {
	nowFunc = f
}
