package date

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func MaxProcs(t testing.TB, n int) {
	prev := runtime.GOMAXPROCS(n)
	t.Cleanup(func() { runtime.GOMAXPROCS(prev) })
}

func TestDateUpdate(t *testing.T) {
	duration := 15 * time.Millisecond // 15 ms * 400 iterations make it run in less than 10 sec
	numOfIterations := 400            // at least a year of updates
	MaxProcs(t, 1)                    // it has to be working on 1 core cpu
	service := &DateService{
		ctx: context.Background(),
	}
	now := service.startUpdatingLoop(duration)

	var actual = []int64{now.UnixMilli()}
	for range time.NewTicker(duration).C {
		if len(actual) == numOfIterations {
			break
		}
		actual = append(actual, service.CurrentDateAsTime().UnixMilli())
	}

	var expected []int64
	for i := 0; i < numOfIterations; i++ {
		expected = append(expected, now.Add(time.Duration(i)*duration).UnixMilli())
	}

	assert.Equal(t, expected, actual)
}

// BenchmarkDateUpdateVsTimeNow
// Results:
// cpu: AMD Ryzen 7 5700U with Radeon Graphics
// BenchmarkDateUpdateVsTimeNow
// BenchmarkDateUpdateVsTimeNow/DateUpdate
// BenchmarkDateUpdateVsTimeNow/DateUpdate-16         	43672665	        26.89 ns/op
// BenchmarkDateUpdateVsTimeNow/TimeNow
// BenchmarkDateUpdateVsTimeNow/TimeNow-16            	18294282	        64.91 ns/op
//
// ~30 ns/op vs ~65 ns/op
// well, seems reasonable only on a truly high load
func BenchmarkDateUpdateVsTimeNow(b *testing.B) {
	MaxProcs(b, 1)
	service := NewDateService(context.Background())
	b.Run("DateUpdate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			service.CurrentDate()
		}
	})
	b.Run("TimeNow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			time.Now()
		}
	})
}
