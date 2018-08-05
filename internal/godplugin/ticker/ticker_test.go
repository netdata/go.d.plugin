package ticker

import (
	"testing"
	"time"
)

var allowedDelta = 5 * time.Millisecond

func TestTickerParallel(t *testing.T) {
	for i := 0; i < 100; i++ {
		i := i
		go func() {
			time.Sleep(time.Second / 100 * time.Duration(i))
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
		diff := abs(now.Round(time.Second).Sub(now))
		if diff >= allowedDelta {
			t.Errorf("ticker is not aligned: expect delta < %v but was: %v (%s)", allowedDelta, diff, now.Format(time.RFC3339Nano))
		}
		if i > 0 {
			dt := now.Sub(prev)
			if abs(dt-time.Second) >= allowedDelta {
				t.Errorf("ticker interval: expect delta < %v ns but was: %v", allowedDelta, abs(dt-time.Second))
			}
		}
		prev = now
	}
}

func abs(a time.Duration) time.Duration {
	if a < 0 {
		return -a
	}
	return a
}
