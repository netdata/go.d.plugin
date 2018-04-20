package web_log

import (
	"math"
	"sort"
	"strconv"
)

func newHistogram(n string, h []int) *histogram {
	sort.Ints(h)
	b := make([]int, len(h))
	s := make([]string, len(h))
	for idx, v := range h {
		b[idx] = v * 1000
		s[idx] = strconv.Itoa(v)
	}
	return &histogram{
		name:        n,
		bucketStr:   append(s, "inf"),
		bucketIndex: append(b, math.MaxInt64),
		buckets:     make([]int, len(h)+1),
	}
}

type histogram struct {
	name        string
	bucketStr   []string
	bucketIndex []int
	buckets     []int
}

func (h *histogram) set(v int) {
	for i := len(h.bucketIndex) - 1; i > -1; i-- {
		if v <= h.bucketIndex[i] {
			h.buckets[i]++
			continue
		}
		break
	}
}
