// Package date
// contains logic to update date without extra time.Now() calls
//
// BenchmarkDateUpdateVsTimeNow results:
// cpu: AMD Ryzen 7 5700U with Radeon Graphics
// BenchmarkDateUpdateVsTimeNow
// BenchmarkDateUpdateVsTimeNow/DateUpdate
// BenchmarkDateUpdateVsTimeNow/DateUpdate-16         	43672665	        26.89 ns/op
// BenchmarkDateUpdateVsTimeNow/TimeNow
// BenchmarkDateUpdateVsTimeNow/TimeNow-16            	18294282	        64.91 ns/op
//
// ~30 ns/op vs ~65 ns/op
// well, seems reasonable only on a truly high load. so now we know!
package date

import (
	"context"
	"sync/atomic"
	"time"
)

type DateService struct {
	currentDate atomic.Int64
	ctx         context.Context
}

const day = 24 * time.Hour

// NewDateService holds current date and updates it once a day
func NewDateService(ctx context.Context) *DateService {
	service := &DateService{ctx: ctx}
	service.startUpdatingLoop(day)
	return service
}

func (d *DateService) CurrentDate() (year int, month time.Month, day int) {
	return d.CurrentDateAsTime().Date()
}

const stringDateFormat = time.DateOnly

func (d *DateService) CurrentDateAsShortString() string {
	return d.CurrentDateAsTime().Format(stringDateFormat)
}

func CurrentDateFromShortString(s string) (time.Time, error) {
	return time.Parse(stringDateFormat, s)
}

func (d *DateService) CurrentDateAsTime() time.Time {
	currentDate := d.currentDate.Load()
	return fromInt64(currentDate)
}

func (d *DateService) startUpdatingLoop(duration time.Duration) time.Time {
	now := time.Now().Add(2 * time.Hour) // FIXME: on the CLoud Run time is in UTC, but converting it to Podgorica time using In(location) is not working for some reason
	currentDate := now.Truncate(duration)
	d.currentDate.Store(toInt64(currentDate))
	go func() {
		tomorrow := currentDate.Add(duration)
		timer := time.NewTimer(tomorrow.Sub(now))
		for {
			select {
			case <-d.ctx.Done():
				return
			case now = <-timer.C:
				now = now.Add(2 * time.Hour) // FIXME: on the CLoud Run time is in UTC, but converting it to Podgorica time using In(location) is not working for some reason
				d.currentDate.Store(toInt64(tomorrow))
				tomorrow = tomorrow.Add(duration)
				timer.Reset(tomorrow.Sub(now))
			}
		}
	}()
	return currentDate
}
