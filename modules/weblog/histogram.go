package weblog

import (
	"math"
	"sort"
	"strconv"
)

type (
	histogram []*histVal
	histVal   struct {
		id    string
		name  string
		value int
		count int
	}
)

func (h histogram) set(v int) {
	for i := len(h) - 1; i > -1; i-- {
		if v <= h[i].value {
			h[i].count++
			continue
		}
		break
	}
}

func newHistogram(prefix string, r []int) histogram {
	var h histogram

	sort.Ints(r)
	for _, v := range r {
		n := strconv.Itoa(v)
		v := &histVal{id: prefix + "_" + n, name: n, value: v * 1000}
		h = append(h, v)
	}

	h = append(h, &histVal{id: prefix + "_inf", name: "inf", value: math.MaxInt64})

	return h
}
