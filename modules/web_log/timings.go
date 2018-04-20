package web_log

import (
	"strconv"
	"strings"
)

func newTimings(n string) *timings {
	return &timings{name: n, min: -1}
}

type timings struct {
	name  string
	min   int
	max   int
	sum   int
	count int
}

func (t *timings) set(s string) int {
	var n int
	switch {
	case s == "0.000":
		n = 0
	case strings.Contains(s, "."):
		if v, err := strconv.ParseFloat(s, 10); err == nil {
			n = int(v * 1e6)
		}
	default:
		if v, err := strconv.Atoi(s); err == nil {
			n = v
		}
	}

	if t.min == -1 {
		t.min = n
	}
	if n > t.max {
		t.max = n
	} else if n < t.min {
		t.min = n
	}
	t.sum += n
	t.count++
	return n
}

func (t *timings) active() bool {
	return t.min != -1
}

func (t *timings) reset() {
	t.min = -1
	t.max = 0
	t.sum = 0
	t.count = 0
}
