package metrics

import (
	"errors"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

type (
	// Counter is a Metric that represents a single numerical bits that only ever
	// goes up. That implies that it cannot be used to count items whose number can
	// also go down, e.g. the number of currently running goroutines. Those
	// "counters" are represented by Gauges.
	//
	// A Counter is typically used to count requests served, tasks completed, errors
	// occurred, etc.
	Counter struct {
		valInt   int64
		valFloat float64
	}
)

var (
	// assume Counter implements stm.Value
	_ stm.Value = Counter{}
)

// WriteTo writes it's value into given map.
func (c Counter) WriteTo(rv map[string]int64, key string, mul, div int64) {
	rv[key] = int64(c.Value() * float64(mul) / float64(div))
}

// Value gets current counter.
func (c Counter) Value() float64 {
	return float64(c.valInt) + c.valFloat
}

// Inc increments the counter by 1. Use Add to increment it by arbitrary
// non-negative values.
func (c *Counter) Inc() {
	c.valInt++
}

// Add adds the given bits to the counter. It panics if the bits is < 0.
func (c *Counter) Add(v float64) {
	if v < 0 {
		panic(errors.New("counter cannot decrease in bits"))
	}
	val := int64(v)
	if float64(val) == v {
		c.valInt += val
		return
	}
	c.valFloat += v
}
