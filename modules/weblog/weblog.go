package weblog

import (
	"strconv"
	"strings"
	"sync"

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

	//tail   *tail.Tail
	//timings    timings
	//histograms histograms

	parser.Parser
	filter filter.Filter

	matchedURL string

	urlCats  []category.Category
	userCats []category.Category

	curPollIPs map[string]bool
	allTimeIPs map[string]bool

	tail *tail.Tail

	charts *modules.Charts

	mux     *sync.Mutex
	metrics map[string]int64
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

func (WebLog) Collect() map[string]int64 {
	return nil
}

func (w *WebLog) parseLoop() {
	for {
		select {
		case line := <-w.tail.Lines:
			w.parseLine(line.Text)
		}
	}
}

func (w *WebLog) parseLine(line string) {
	//if !w.filter.Filter(line) {
	//	return
	//}
	//
	//gm, ok := w.Parse(line)
	//
	//if !ok {
	//	return
	//}
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

// method, url, http version
func (w *WebLog) perRequest(gm parser.GroupMap) {
	request := gm.Get("request")

	if request != "" {
		gm, _ = w.Parse(request)
	}

	if _, ok := gm.Lookup("method"); ok {
		w.perHTTPMethod(gm)
	}

	if _, ok := gm.Lookup("url"); ok {
		w.perURLCategory(gm)
	}

	if _, ok := gm.Lookup("version"); ok {
		w.perHTTPMethod(gm)
	}
}

func (w *WebLog) perHTTPMethod(gm parser.GroupMap) {
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

func (w *WebLog) perURLCategory(gm parser.GroupMap) {
	url := gm.Get("url")

	for _, v := range w.urlCats {
		if v.Match(url) {
			w.metrics[v.Name()]++
			w.matchedURL = v.Name()
			break
		}
	}
	w.matchedURL = ""
	w.metrics["url_other"]++
}

func (w *WebLog) perUserCategory(gm parser.GroupMap) {
	url := gm.Get("user_defined")

	for _, v := range w.userCats {
		if v.Match(url) {
			w.metrics[v.Name()]++
			break
		}
	}
	w.metrics["user_defined_other"]++
}

func (w *WebLog) perHTTPVersion(gm parser.GroupMap) {
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

func (w *WebLog) perBytesSent(gm parser.GroupMap) {

}

func (w *WebLog) perRespLength(gm parser.GroupMap) {

}

func (w *WebLog) perRespTime(gm parser.GroupMap) {

}

func (w *WebLog) perRespTimeUpstream(gm parser.GroupMap) {

}

func (w *WebLog) perIPProto(gm parser.GroupMap) {
	var (
		address = gm.Get("address")
		proto   = "ipv4"
	)

	if strings.Contains(address, ":") {
		proto = "ipv6"
	}

	w.metrics["req_"+proto]++

	if _, ok := w.curPollIPs[address]; !ok {
		w.curPollIPs[address] = true
		w.metrics["unique_cur_"+proto]++
	}

	if !w.DoAllTimeIPs {
		return
	}

	if _, ok := w.allTimeIPs[address]; !ok {
		w.allTimeIPs[address] = true
		w.metrics["unique_all_"+proto]++
	}

}

func (w *WebLog) perURLCategoryStats(gm parser.GroupMap) {
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

//func (w *WebLog) Check() bool {
//
//	w.tail = tail.newMatcher(w.Path)
//	err := w.tail.Init()
//	if err != nil {
//		w.Error(err)
//		return false
//	}
//
//	// read last line
//	line, err := tail.ReadLastLine(w.Path)
//	if err != nil {
//		w.Error(err)
//		return false
//	}
//
//	// get parser: custom or one of predefined in csv.go
//	re, err := getPattern(w.RawCustomParser, line)
//	if err != nil {
//		w.Error(err)
//		return false
//	}
//	w.parser = re
//
//	c, err := getCategories(w.RawURLCat, "url")
//	if err != nil {
//		w.Error(err)
//		return false
//	}
//	w.urlCat = c
//
//	if w.DoPerURLCharts {
//		for _, v := range w.urlCat.items {
//			w.timings.add(v.id)
//		}
//	}
//
//	c, err = getCategories(w.RawUserCat, "user")
//	if err != nil {
//		w.Error(err)
//		return false
//	}
//	w.userCat = c
//
//	f, err := getFilter(w.rawFilter)
//	if err != nil {
//		w.Error(err)
//		return false
//	}
//	w.filter = f
//
//	if len(w.RawHistogram) != 0 {
//		w.histograms = getHistograms(w.RawHistogram)
//	}
//
//	w.Info("collected data:", w.parser.SubexpNames()[1:])
//	return true
//}
//
//func (w *WebLog) GatherMetrics() map[string]int64 {
//	f, err := w.tail.Tail()
//
//	if err == tail.SizeNotChanged {
//		return w.data
//	}
//
//	if err != nil {
//		return nil
//	}
//
//	defer f.Close()
//
//	uniqIPs := make(map[string]bool)
//	w.timings.reset()
//
//	s := bufio.NewScanner(f)
//
//	for s.Scan() {
//		row := s.Text()
//		if w.filter.exist() && !w.filter.rawFilter(row) {
//			continue
//		}
//
//		m := w.parser.FindStringSubmatch(row)
//		if m == nil {
//			w.data["unmatched"]++
//			continue
//		}
//
//		w.gm.update(w.parser.SubexpNames(), m)
//
//		code, codeFam := w.gm.get(keyCode), w.gm.get(keyCode)[1:]
//
//		// ResponseCodes chart
//		if _, ok := w.data[codeFam+"xx"]; ok {
//			w.data[codeFam+"xx"]++
//		} else {
//			w.data["0xx"]++
//		}
//
//		// ResponseStatuses chart
//		w.reqPerCodeFamily(code)
//
//		// ResponseCodesDetailed chart
//		if w.DoCodesDetailed {
//			w.reqPerCodeDetail(code)
//		}
//
//		// chartBandwidth chart
//		if v, ok := w.gm.lookup(keyBytesSent); ok {
//			w.data["bytes_sent"] += toInt(v)
//		}
//
//		if v, ok := w.gm.lookup(keyRespLen); ok {
//			w.data["resp_length"] += toInt(v)
//		}
//
//		// ResponseTime and ResponseTimeHistogram charts3
//		if v, ok := w.gm.lookup(keyRespTime); ok {
//			i := w.timings.get(keyRespTime).set(v)
//			if w.histograms.exist() {
//				w.histograms.get(keyRespTimeHist).set(i)
//			}
//		}
//
//		// ResponseTimeUpstream, ResponseTimeUpstreamHistogram charts3
//		if v, ok := w.gm.lookup(keyRespTimeUpstream); ok && v != "-" {
//			i := w.timings.get(keyRespTimeUpstream).set(v)
//			if w.histograms.exist() {
//				w.histograms.get(keyRespTimeUpstreamHist).set(i)
//			}
//		}
//
//		// ReqPerUrl, reqPerHTTPMethod, chartReqPerHTTPVer charts3
//		var matchedURL string
//		if w.gm.has(keyRequest) {
//			matchedURL = w.parseRequest()
//		}
//
//		// ReqPerUserDefined chart
//		if v, ok := w.gm.lookup(keyUserDefined); ok && w.userCat.exist() {
//			w.reqPerCategory(v, w.userCat)
//		}
//
//		// chartRespCodesDetailed, chartBandwidth, chartRespTime per URL (Category) charts3
//		if matchedURL != "" && w.DoPerURLCharts {
//			w.perCategoryStats(matchedURL)
//		}
//
//		// RequestsPerIPProto, chartClientsCurr, chartClientsAll charts3
//		if v, ok := w.gm.lookup(keyAddress); ok {
//			w.reqPerIPProto(v, uniqIPs)
//		}
//
//	}
//
//	for n, v := range w.timings {
//		if !v.active() {
//			continue
//		}
//		w.data[n+"_min"] += int64(v.min)
//		w.data[n+"_avg"] += int64(v.avg())
//		w.data[n+"_max"] += int64(v.max)
//	}
//
//	for _, h := range w.histograms {
//		for _, v := range h {
//			w.data[v.id] = int64(v.count)
//		}
//	}
//
//	return w.data
//}

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
