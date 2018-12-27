package weblog

import (
	"errors"
	"math"
	"sort"
	"strconv"
)

type (
	histogram []*histVal
	histVal   struct {
		id    string
		name  string
		value int64
		count int
	}
)

func (h histogram) set(v int) {
	for i := len(h) - 1; i > -1; i-- {
		if int64(v) <= h[i].value {
			h[i].count++
			continue
		}
		break
	}
}

func newHistogram(prefix string, r []int) (histogram, error) {
	var h histogram

	if !sort.IntsAreSorted(r) {
		return nil, errors.New("not sorted histogram")
	}

	sort.Ints(r)
	for _, v := range r {
		if v < 0 {
			return nil, errors.New("histogram contains negative value")
		}

		n := strconv.Itoa(v)
		v := &histVal{id: prefix + "_" + n, name: n, value: int64(v * 1000)}
		h = append(h, v)
	}

	h = append(h, &histVal{id: prefix + "_inf", name: "inf", value: math.MaxInt64})

	return h, nil
}
