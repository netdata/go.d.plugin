package weblog

//
//type timings map[string]*timing
//
//func (t timings) reset() {
//	for _, v := range t {
//		v.reset()
//	}
//}
//
//func (t *timings) add(n string) {
//	(*t)[n] = new(timing)
//}
//
//func (t timings) get(n string) *timing {
//	return t[n]
//}
//
//func (t timings) lookup(n string) (*timing, bool) {
//	v, ok := t[n]
//	return v, ok
//}
//
//type timing struct {
//	min   int
//	max   int
//	sum   int
//	count int
//}
//
////
//func (t *timing) set(s string) int {
//	var n int
//
//	switch {
//	case s == "0.000":
//	case strings.Contains(s, "."):
//		if v, err := strconv.ParseFloat(s, 10); err == nil {
//			n = int(v * 1e6)
//		}
//	default:
//		if v, err := strconv.Atoi(s); err == nil {
//			n = v
//		}
//	}
//
//	if t.min == -1 {
//		t.min = n
//	}
//	if n > t.max {
//		t.max = n
//	} else if n < t.min {
//		t.min = n
//	}
//	t.sum += n
//	t.count++
//	return n
//}
//
//func (t *timing) avg() int {
//	return t.sum / t.count
//}
//
//func (t *timing) reset() {
//	t.min = -1
//	t.max = 0
//	t.sum = 0
//	t.count = 0
//}
//
//func (t timing) active() bool {
//	return t.min != -1
//}
