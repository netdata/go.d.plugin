package weblog

import (
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/modules/weblog/category"
	"github.com/netdata/go.d.plugin/modules/weblog/charts"
	"github.com/netdata/go.d.plugin/modules/weblog/filter"
	"github.com/netdata/go.d.plugin/modules/weblog/parser"

	"github.com/hpcloud/tail"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("web_log", creator)
}

func New() *WebLog {
	return &WebLog{}
}

type WebLog struct {
	modules.Base

	Path             string         `yaml:"path" validate:"required"`
	Filter           filter.Raw     `yaml:"filter"`
	URLCats          []category.Raw `yaml:"categories"`
	UserCats         []category.Raw `yaml:"user_categories"`
	CustomParser     string         `yaml:"custom_log_format"`
	Histogram        []int          `yaml:"histogram"`
	DoCodesDetailed  bool           `yaml:"detailed_response_codes"`
	DoCodesAggregate bool           `yaml:"detailed_response_codes_aggregate"`
	DoPerURLCharts   bool           `yaml:"per_category_charts"`
	DoAllTimeIPs     bool           `yaml:"all_time_clients"`

	tail *tail.Tail

	charts *modules.Charts

	parser.Parser
	filter filter.Filter

	matchedURL string

	categories struct {
		url  []category.Category
		user []category.Category
	}

	hooks struct {
		stop    chan struct{}
		collect chan struct{}
	}

	clientIPs struct {
		cur map[string]bool
		all map[string]bool
	}

	tailMetrics map[string]int64
	metrics     map[string]int64
}

func (WebLog) Cleanup() {

}

func (w *WebLog) Init() bool {
	return false
}

func (WebLog) Check() bool {
	return false
}

func (WebLog) Charts() *modules.Charts {
	return nil
}

func (w *WebLog) Collect() map[string]int64 {
	w.hooks.collect <- struct{}{}

	return w.metrics
}

func (w *WebLog) parseLoop() {
LOOP:
	for {
		select {
		case <-w.hooks.stop:
			w.cleanup()
			break LOOP
		case <-w.hooks.collect:
			w.copyMetrics()
		case line := <-w.tail.Lines:
			w.parseLine(line.Text)
		}
	}
}

func (w *WebLog) cleanup() {
	w.tail.Cleanup()
	_ = w.tail.Stop()
}

func (w *WebLog) copyMetrics() {
	for k, v := range w.tailMetrics {
		w.metrics[k] = v
	}

	// add timings
	// reset timings
}

func (w *WebLog) parseLine(line string) {
	if !w.filter.Filter(line) {
		return
	}

	gm, ok := w.Parse(line)

	if !ok {
		w.metrics["unmatched"]++
		return
	}

	w.codeFam(gm)

	w.codeStatus(gm)

	if w.DoCodesDetailed {
		w.codeDetailed(gm)
	}

	w.request(gm)

	if _, ok := gm.Lookup("user_defined"); ok && len(w.categories.user) > 0 {
		w.userCategory(gm)
	}

	if _, ok := gm.Lookup("bytes_sent"); ok {
		w.bytesSent(gm)
	}

	if _, ok := gm.Lookup("resp_length"); ok {
		w.respLength(gm)
	}

	if _, ok := gm.Lookup("address"); ok {
		w.ipProto(gm)
	}

	if w.DoPerURLCharts && w.matchedURL != "" {
		w.urlCategoryStats(gm)
	}

}

func (w *WebLog) codeFam(gm parser.GroupMap) {
	fam := gm.Get("code")[:1] + "xx"

	if _, ok := w.metrics[fam]; ok {
		w.metrics[fam]++
	} else {
		w.metrics["0xx"]++
	}
}

func (w *WebLog) codeDetailed(gm parser.GroupMap) {
	code := gm.Get("code")

	if _, ok := w.metrics[code]; ok {
		w.metrics[code]++
		return
	}

	if w.DoCodesAggregate {
		chart := w.charts.Get(charts.ResponseCodesDetailed.ID)
		_ = chart.AddDim(&modules.Dim{ID: code, Algo: modules.Incremental})
		w.metrics[code]++
		return
	}

	var v = "other"

	if code[0] <= 53 {
		v = code[:1] + "xx"
	}

	chart := w.charts.Get(charts.ResponseCodesDetailed.ID + "_" + v)
	_ = chart.AddDim(&modules.Dim{ID: code, Algo: modules.Incremental})
	w.metrics[code]++
}

func (w *WebLog) codeStatus(gm parser.GroupMap) {
	code, fam := gm.Get("code"), gm.Get("code")[:1]

	switch {
	case fam == "2", code == "304", fam == "1":
		w.metrics["successful_requests"]++
	case fam == "3":
		w.metrics["redirects"]++
	case fam == "4":
		w.metrics["bad_requests"]++
	case fam == "5":
		w.metrics["server_errors"]++
	default:
		w.metrics["other_requests"]++
	}
}

func (w *WebLog) request(gm parser.GroupMap) {
	request := gm.Get("request")

	if request != "" {
		gm, _ = w.Parse(request)
	}

	if _, ok := gm.Lookup("method"); ok {
		w.httpMethod(gm)
	}

	if _, ok := gm.Lookup("url"); ok {
		w.urlCategory(gm)
	}

	if _, ok := gm.Lookup("version"); ok {
		w.httpVersion(gm)
	}
}

func (w *WebLog) httpMethod(gm parser.GroupMap) {
	method := gm.Get("method")

	if _, ok := w.metrics[method]; !ok {
		chart := w.charts.Get(charts.RequestsPerHTTPMethod.ID)
		_ = chart.AddDim(&modules.Dim{
			ID:   method,
			Algo: modules.Incremental,
		})
	}

	w.metrics[method]++
}

func (w *WebLog) urlCategory(gm parser.GroupMap) {
	url := gm.Get("url")

	for _, v := range w.categories.url {
		if v.Match(url) {
			w.metrics[v.Name()]++
			w.matchedURL = v.Name()
			return
		}
	}
	w.matchedURL = ""
	w.metrics["url_other"]++
}

func (w *WebLog) userCategory(gm parser.GroupMap) {
	userDefined := gm.Get("user_defined")

	for _, cat := range w.categories.user {
		if cat.Match(userDefined) {
			w.metrics[cat.Name()]++
			return
		}
	}
	w.metrics["user_defined_other"]++
}

func (w *WebLog) httpVersion(gm parser.GroupMap) {
	version := gm.Get("version")

	dimID := strings.Replace(gm.Get("version"), ".", "_", 1)

	if _, ok := w.metrics[dimID]; !ok {
		chart := w.charts.Get(charts.RequestsPerHTTPVersion.ID)
		_ = chart.AddDim(&modules.Dim{
			ID:   dimID,
			Name: version,
			Algo: modules.Incremental,
		})
	}

	w.metrics[dimID]++
}

func (w *WebLog) bytesSent(gm parser.GroupMap) {
	v := gm.Get("bytes_sent")

	w.metrics["bytes_sent"] += toInt(v)
}

func (w *WebLog) respLength(gm parser.GroupMap) {
	v := gm.Get("resp_length")

	w.metrics["resp_length"] += toInt(v)
}

func (w *WebLog) respTime(gm parser.GroupMap) {

}

func (w *WebLog) respTimeUpstream(gm parser.GroupMap) {

}

func (w *WebLog) ipProto(gm parser.GroupMap) {
	var (
		address = gm.Get("address")
		proto   = "ipv4"
	)

	if strings.Contains(address, ":") {
		proto = "ipv6"
	}

	w.metrics["req_"+proto]++

	if _, ok := w.clientIPs.cur[address]; !ok {
		w.clientIPs.cur[address] = true
		w.metrics["unique_cur_"+proto]++
	}

	if !w.DoAllTimeIPs {
		return
	}

	if _, ok := w.clientIPs.all[address]; !ok {
		w.clientIPs.all[address] = true
		w.metrics["unique_all_"+proto]++
	}

}

func (w *WebLog) urlCategoryStats(gm parser.GroupMap) {
	code := gm.Get("code")
	v := w.matchedURL + "_" + code

	if _, ok := w.metrics[v]; !ok {
		chart := w.charts.Get(charts.ResponseCodesDetailed.ID + "_" + w.matchedURL)
		_ = chart.AddDim(&modules.Dim{
			ID:   v,
			Name: code,
			Algo: modules.Incremental,
		})
	}

	w.metrics[v]++

	if v, ok := gm.Lookup("bytes_sent"); ok {
		w.metrics[w.matchedURL+"_bytes_sent"] += toInt(v)
	}

	if v, ok := gm.Lookup("resp_length"); ok {
		w.metrics[w.matchedURL+"_resp_length"] += toInt(v)
	}

	//if v, ok := gm.Lookup("resp_time"); ok {
	//	w.timings.get(id).set(v)
	//}
}

func toInt(s string) int64 {
	if s == "-" {
		return 0
	}
	v, _ := strconv.Atoi(s)

	return int64(v)
}

//
//func init() {
//	f := func() modules.Module {
//		return &WebLog{
//			DoCodesDetailed:    true,
//			DoCodesAggregate: true,
//			DoPerURLCharts:    true,
//			DoAllTimeIPs:     true,
//			timings: timings{
//				keyRespTime:         &timing{},
//				keyRespTimeUpstream: &timing{},
//			},
//			gm:      make(groupMap),
//			uniqIPs: make(map[string]bool),
//			data: map[string]int64{
//				"successful_requests":    0,
//				"redirects":              0,
//				"bad_requests":           0,
//				"server_errors":          0,
//				"other_requests":         0,
//				"2xx":                    0,
//				"5xx":                    0,
//				"3xx":                    0,
//				"4xx":                    0,
//				"1xx":                    0,
//				"0xx":                    0,
//				"unmatched":              0,
//				"bytes_sent":             0,
//				"resp_length":            0,
//				"resp_time_min":          0,
//				"resp_time_max":          0,
//				"resp_time_avg":          0,
//				"resp_time_upstream_min": 0,
//				"resp_time_upstream_max": 0,
//				"resp_time_upstream_avg": 0,
//				"unique_cur_ipv4":        0,
//				"unique_cur_ipv6":        0,
//				"unique_tot_ipv4":        0,
//				"unique_tot_ipv6":        0,
//				"req_ipv4":               0,
//				"req_ipv6":               0,
//				"GET":                    0, // GET should be green on the dashboard
//			},
//		}
//	}
//	modules.Add(f)
//}
