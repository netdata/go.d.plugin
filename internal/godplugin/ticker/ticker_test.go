package ticker

import (
	"math/rand"
	"testing"
	"time"
)

var allowedDelta = int(10 * time.Millisecond)

func TestTickerParallel(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(int(time.Second))))
			TestTicker(t)
		}()
	}
	time.Sleep(7 * time.Second)
}

func TestTicker(t *testing.T) {
	ticker := New(time.Second)
	defer ticker.Stop()
	prev := time.Now()
	for i := 0; i < 5; i++ {
		<-ticker.C
		now := time.Now()
		if now.Nanosecond() >= allowedDelta {
			t.Errorf("ticker is not aligned: expect delta < %d ns but was: %d (%s)", allowedDelta, now.Nanosecond(), now.Format(time.RFC3339Nano))
		}
		if i > 0 {
			dt := now.Sub(prev)
			if abs(int(dt)-int(time.Second)) >= allowedDelta {
				t.Errorf("ticker interval: expect delta < %d ns but was: %v", allowedDelta, dt)
			}
		}
		prev = now
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
