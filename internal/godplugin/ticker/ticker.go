package ticker

import "time"

type (
	// Ticker holds a channel that delivers ticks of a clock at intervals.
	// The ticks is aligned to interval boundaries.
	Ticker struct {
		C        <-chan int
		done     chan int
		loops    int
		interval time.Duration
	}
)

// New returns a new Ticker containing a channel that will send the time with a period specified by the duration argument.
// It adjusts the intervals or drops ticks to make up for slow receivers.
// The duration must be greater than zero; if not, NewTicker will panic. Stop the ticker to release associated resources.
func New(interval time.Duration) *Ticker {
	ticker := &Ticker{
		interval: interval,
	}
	ticker.start()
	return ticker
}

func (t *Ticker) start() {
	ch := make(chan int)
	t.loops = 0
	t.C = ch
	t.done = make(chan int, 1)
	go func() {
	LOOP:
		for {
			now := time.Now()
			nextRun := now.Truncate(t.interval).Add(t.interval)

			time.Sleep(nextRun.Sub(now))
			select {
			case <-t.done:
				close(ch)
				break LOOP
			case ch <- t.loops:
				t.loops++
			}
		}
	}()
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a read from the channel succeeding incorrectly.
func (t *Ticker) Stop() {
	t.done <- 1
}
