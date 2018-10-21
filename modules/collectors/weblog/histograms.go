package weblog

//
//import (
//	"math"
//	"sort"
//	"strconv"
//)
//
//type (
//	histograms map[string]histogram
//	histogram  []*histVal
//	histVal    struct {
//		id    string
//		name  string
//		value int
//		count int
//	}
//)
//
//func (hs histograms) exist() bool {
//	return len(hs) != 0
//}
//
//func (hs histograms) get(n string) histogram {
//	return hs[n]
//}
//
//func (h histogram) set(v int) {
//	for i := len(h) - 1; i > -1; i-- {
//		if v <= h[i].value {
//			h[i].count++
//			continue
//		}
//		break
//	}
//}
//
//func getHistograms(r []int) histograms {
//	var h1, h2 histogram
//	h := make(histograms)
//
//	sort.Ints(r)
//	for _, v := range r {
//		n := strconv.Itoa(v)
//		h1 = append(h1, &histVal{id: keyRespTimeHist + "_" + n, name: n, value: v * 1000})
//		h2 = append(h2, &histVal{id: keyRespTimeUpstreamHist + "_" + n, name: n, value: v * 1000})
//	}
//	h1 = append(h1, &histVal{id: keyRespTimeHist + "_inf", name: "inf", value: math.MaxInt64})
//	h2 = append(h2, &histVal{id: keyRespTimeUpstreamHist + "_inf", name: "inf", value: math.MaxInt64})
//
//	h[keyRespTimeHist] = h1
//	h[keyRespTimeUpstreamHist] = h2
//	return h
//}
