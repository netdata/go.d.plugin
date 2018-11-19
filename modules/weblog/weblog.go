package weblog

//
//import (
//	"regexp"
//	"strconv"
//	"strings"
//
//	"gopkg.in/yaml.v2"
//
//	"bufio"
//	"github.com/l2isbad/go.d.plugin/internal/modules"
//	"github.com/l2isbad/go.d.plugin/modules/pkg/tail"
//	"github.com/l2isbad/go.d.plugin/pkg/charts"
//)
//
//const (
//	keyAddress          = "address"
//	keyCode             = "code"
//	keyRequest          = "request"
//	keyUserDefined      = "user_defined"
//	keyBytesSent        = "bytes_sent"
//	keyRespTime         = "resp_time"
//	keyRespTimeUpstream = "resp_time_upstream"
//	keyRespLen          = "resp_length"
//
//	keyRespTimeHist         = "resp_time_hist"
//	keyRespTimeUpstreamHist = "resp_time_hist_upstream"
//)
//
//type WebLog struct {
//	*charts.Charts
//	modules.ModuleBase
//
//	Path             string        `yaml:"path" validate:"required"`
//	RawFilter        rawFilter     `yaml:"filter"`
//	RawURLCat        yaml.MapSlice `yaml:"categories"`
//	RawUserCat       yaml.MapSlice `yaml:"user_defined"`
//	RawCustomParser  string        `yaml:"custom_log_format"`
//	RawHistogram     []int         `yaml:"histogram"`
//	DoCodesDetail    bool          `yaml:"detailed_response_codes"`
//	DoCodesAggregate bool          `yaml:"detailed_response_codes_aggregate"`
//	DoChartURLCat    bool          `yaml:"per_category_charts"`
//	DoClientsAll     bool          `yaml:"clients_all_time"`
//
//	tail   *tail.Tail
//	parser *regexp.Regexp
//
//	fil        filter
//	urlCat     categories
//	userCat    categories
//	timings    timings
//	histograms histograms
//	uniqIPs    map[string]bool
//
//	gm   groupMap
//	data map[string]int64
//}
//
//func (WebLog) Init() {}
//
//func (w *WebLog) Check() bool {
//
//	w.tail = tail.New(w.Path)
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
//	// get parser: custom or one of predefined in patterns.go
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
//	if w.DoChartURLCat {
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
//	f, err := getFilter(w.RawFilter)
//	if err != nil {
//		w.Error(err)
//		return false
//	}
//	w.fil = f
//
//	if len(w.RawHistogram) != 0 {
//		w.histograms = getHistograms(w.RawHistogram)
//	}
//
//	w.Info("collected data:", w.parser.SubexpNames()[1:])
//	return true
//}
//
//func (w *WebLog) GetData() map[string]int64 {
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
//		if w.fil.exist() && !w.fil.filter(row) {
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
//		if w.DoCodesDetail {
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
//		if matchedURL != "" && w.DoChartURLCat {
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
//
//// Per URL and per USER_DEFINED
//func (w *WebLog) reqPerCategory(url string, c categories) string {
//	for _, v := range c.items {
//		if v.Match(url) {
//			w.data[v.id]++
//			return v.id
//		}
//	}
//	w.data[c.other]++
//	return ""
//}
//
//func (w *WebLog) reqPerIPProto(address string, uniqIPs map[string]bool) {
//	var proto = "ipv4"
//
//	if strings.Contains(address, ":") {
//		proto = "ipv6"
//	}
//	w.data["req_"+proto]++
//
//	if _, ok := uniqIPs[address]; !ok {
//		uniqIPs[address] = true
//		w.data["unique_cur_"+proto]++
//	}
//
//	if !w.DoClientsAll {
//		return
//	}
//
//	if _, ok := w.uniqIPs[address]; !ok {
//		w.uniqIPs[address] = true
//		w.data["unique_all_"+proto]++
//	}
//}
//
//func (w *WebLog) reqPerCodeDetail(code string) {
//	if _, ok := w.data[code]; ok {
//		w.data[code]++
//		return
//	}
//
//	if w.DoCodesAggregate {
//		w.Get(chartRespCodesDetailed.ID).AddDim(&Dim{ID: code, Algo: charts.Incremental})
//		w.data[code]++
//		return
//	}
//	var v = "other"
//	if code[0] <= 53 {
//		v = code[:1] + "xx"
//	}
//	w.Get(chartRespCodesDetailed.ID + "_" + v).AddDim(&Dim{ID: code, Algo: charts.Incremental})
//	w.data[code]++
//}
//
//func (w *WebLog) reqPerCodeFamily(code string) {
//	f := code[:1]
//	switch {
//	case f == "2", code == "304", f == "1":
//		w.data["successful_requests"]++
//	case f == "3":
//		w.data["redirects"]++
//	case f == "4":
//		w.data["bad_requests"]++
//	case f == "5":
//		w.data["server_errors"]++
//	default:
//		w.data["other_requests"]++
//	}
//}
//
//func (w *WebLog) reqPerHTTPMethod(method string) {
//	if _, ok := w.data[method]; !ok {
//		w.Get(chartReqPerHTTPMethod.ID).AddDim(&Dim{ID: method, Algo: charts.Incremental})
//	}
//	w.data[method]++
//}
//
//func (w *WebLog) reqPerHTTPVersion(version string) {
//	dimID := strings.Replace(version, ".", "_", 1)
//
//	if _, ok := w.data[dimID]; !ok {
//		w.Get(chartReqPerHTTPVer.ID).AddDim(&Dim{ID: dimID, Name: version, Algo: charts.Incremental})
//	}
//	w.data[dimID]++
//}
//
//func (w *WebLog) parseRequest() (matchedURL string) {
//	req := w.gm.get(keyRequest)
//	if req == "-" {
//		return
//	}
//
//	v := strings.Fields(req)
//	if len(v) != 3 {
//		return
//	}
//
//	// FIXME: assumed that 'version' part is always prefixed with 'HTTP/'
//	method, url, version := v[0], v[1], v[2][5:]
//	if w.urlCat.exist() {
//		if v := w.reqPerCategory(url, w.urlCat); v != "" {
//			matchedURL = v
//		}
//	}
//	w.reqPerHTTPMethod(method)
//	w.reqPerHTTPVersion(version)
//	return
//}
//
//func (w *WebLog) perCategoryStats(id string) {
//	code := w.gm.get(keyCode)
//	v := id + "_" + code
//	if _, ok := w.data[v]; !ok {
//		w.Get(chartRespCodesDetailed.ID + "_" + id).AddDim(&Dim{ID: v, Name: code, Algo: charts.Incremental})
//	}
//	w.data[v]++
//
//	if v, ok := w.gm.lookup(keyBytesSent); ok {
//		w.data[id+"_bytes_sent"] += toInt(v)
//	}
//
//	if v, ok := w.gm.lookup(keyRespLen); ok {
//		w.data[id+"_resp_length"] += toInt(v)
//	}
//
//	if v, ok := w.gm.lookup(keyRespTime); ok {
//		w.timings.get(id).set(v)
//	}
//}
//
//func toInt(s string) int64 {
//	if s == "-" {
//		return 0
//	}
//	v, _ := strconv.Atoi(s)
//	return int64(v)
//}
//
//func init() {
//	f := func() modules.Module {
//		return &WebLog{
//			DoCodesDetail:    true,
//			DoCodesAggregate: true,
//			DoChartURLCat:    true,
//			DoClientsAll:     true,
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
