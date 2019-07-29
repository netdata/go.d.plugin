package metrics

import (
	"math"
)

type (
	// A Summary captures individual observations from an event or sample stream and
	// summarizes them in a manner similar to traditional summary statistics:
	//   sum of observations
	//   observation count
	//   observation average.
	//
	// To create summary instances, use NewSummary.
	Summary interface {
		Observer
		Reset()
	}

	summary struct {
		min   float64
		max   float64
		sum   float64
		count int64
	}
)

// NewSummary creates a new Summary.
func NewSummary() Summary {
	return &summary{
		min: math.MaxFloat64,
		max: -math.MaxFloat64,
	}
}

// WriteTo writes it's values into given map.
// It add those key-value pairs:
//   ${key}_sum        gauge, for sum of it's observed values from last Reset calls
//   ${key}_count      counter, for count of it's observed values from last Reset calls
//   ${key}_min        gauge, for min of it's observed values from last Reset calls (only exists if count > 0)
//   ${key}_max        gauge, for max of it's observed values from last Reset calls (only exists if count > 0)
//   ${key}_avg        gauge, for avg of it's observed values from last Reset calls (only exists if count > 0)
func (s *summary) WriteTo(rv map[string]int64, key string, mul, div int) {
	if s.count > 0 {
		rv[key+"_min"] = int64(s.min * float64(mul) / float64(div))
		rv[key+"_max"] = int64(s.max * float64(mul) / float64(div))
		rv[key+"_sum"] = int64(s.sum * float64(mul) / float64(div))
		rv[key+"_count"] = s.count
		rv[key+"_avg"] = int64(s.sum / float64(s.count) * float64(mul) / float64(div))
	} else {
		rv[key+"_count"] = 0
		rv[key+"_sum"] = 0
		delete(rv, key+"_min")
		delete(rv, key+"_max")
		delete(rv, key+"_avg")
	}
}

// Reset resets all of its counters.
// Call it before every scrape loop.
func (s *summary) Reset() {
	s.min = math.MaxFloat64
	s.max = -math.MaxFloat64
	s.sum = 0
	s.count = 0
}

// Observe observes a value
func (s *summary) Observe(v float64) {
	if v > s.max {
		s.max = v
	}
	if v < s.min {
		s.min = v
	}
	s.sum += v
	s.count++
}
