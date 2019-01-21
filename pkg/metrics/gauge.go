package metrics

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

type (
	// Gauge is a Metric that represents a single numerical value that can
	// arbitrarily go up and down.
	//
	// A Gauge is typically used for measured values like temperatures or current
	// memory usage, but also "counts" that can go up and down, like the number of
	// running goroutines.
	Gauge float64
)

var (
	// assume Counter implements stm.Value
	_ stm.Value = Gauge(0)
)

// WriteTo writes it's value into given map.
func (g Gauge) WriteTo(rv map[string]int64, key string, mul, div int) {
	rv[key] = int64(float64(g) * float64(mul) / float64(div))
}

// Value gets current counter.
func (g Gauge) Value() float64 {
	return float64(g)
}

// Set sets the atomicGauge to an arbitrary bits.
func (g *Gauge) Set(v float64) {
	*g = Gauge(v)
}

// Inc increments the atomicGauge by 1. Use Add to increment it by arbitrary
// values.
func (g *Gauge) Inc() {
	*g++
}

// Dec decrements the atomicGauge by 1. Use Sub to decrement it by arbitrary
// values.
func (g *Gauge) Dec() {
	*g--
}

// Add adds the given bits to the atomicGauge. (The bits can be negative,
// resulting in a decrease of the atomicGauge.)
func (g *Gauge) Add(delta float64) {
	*g += Gauge(delta)
}

// Sub subtracts the given bits from the atomicGauge. (The bits can be
// negative, resulting in an increase of the atomicGauge.)
func (g *Gauge) Sub(delta float64) {
	*g -= Gauge(delta)
}

// SetToCurrentTime sets the atomicGauge to the current Unix time in second.
func (g *Gauge) SetToCurrentTime() {
	*g = Gauge(time.Now().UnixNano()) / 1e9
}
