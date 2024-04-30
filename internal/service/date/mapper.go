package date

import (
	"time"
)

func fromInt64(i int64) time.Time {
	return time.UnixMilli(i)
}

func toInt64(t time.Time) int64 {
	return t.UnixMilli()
}
