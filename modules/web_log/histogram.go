package web_log

import (
	"math"
	"sort"
	"strconv"
)

type histValue struct {
	id    string
	name  string
	value int
	count int
}

type histogram []*histValue

func (h *histogram) set(v int) {
	for i := len(*h) - 1; i > -1; i-- {
		if v <= (*h)[i].value {
			(*h)[i].count++
			continue
		}
		break
	}
}

func newHistograms(prefix string, h []int) *histogram {
	sort.Ints(h)
	rv := make(histogram, len(h))
	for idx, v := range h {
		n := strconv.Itoa(v)
		rv[idx] = &histValue{id: prefix + "_" + n, name: n, value: v * 1000}
	}
	rv = append(rv, &histValue{id: prefix + "_inf", name: "inf", value: math.MaxInt64})
	return &rv
}
