package weblog

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/simpletail"

	"github.com/hpcloud/tail"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("web_log", creator)
}

func New() *WebLog {
	return &WebLog{
		DoCodesDetailed:  true,
		DoCodesAggregate: true,
		DoAllTimeIPs:     true,
		DoPerURLCharts:   false,

		stop:           make(chan struct{}),
		collect:        make(chan struct{}),
		done:           make(chan struct{}),
		uniqIPs:        make(map[string]bool),
		uniqIPsAllTime: make(map[string]bool),
	}
}

type WebLog struct {
	modules.Base

	Path string `yaml:"path" validate:"required"`

	Filter   rawFilter     `yaml:"filter"`
	URLCats  []rawCategory `yaml:"categories"`
	UserCats []rawCategory `yaml:"user_categories"`

	CustomParser csvPattern `yaml:"custom_log_format"`
	Histogram    []int      `yaml:"histogram"`

	DoCodesDetailed  bool `yaml:"detailed_response_codes"`
	DoCodesAggregate bool `yaml:"detailed_response_codes_aggregate"`
	DoPerURLCharts   bool `yaml:"per_category_charts"`
	DoAllTimeIPs     bool `yaml:"all_time_clients"`

	tail *tail.Tail

	charts *modules.Charts

	parser
	filter matcher

	matchedURL string
	updated    bool

	urlCats  []*category
	userCats []*category

	stop    chan struct{}
	collect chan struct{}
	done    chan struct{}

	uniqIPs        map[string]bool
	uniqIPsAllTime map[string]bool

	tailMetrics map[string]int64
	metrics     map[string]int64
}

func (WebLog) Cleanup() {}

func (w *WebLog) initFilter() error {
	f, err := newFilter(w.Filter)
	if err != nil {
		return fmt.Errorf("error on creating filter : %s", err)
	}
	w.filter = f

	return nil
}

func (w *WebLog) initCategories() error {
	for _, raw := range w.URLCats {
		cat, err := newCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating category %s : %s", raw, err)
		}
		w.urlCats = append(w.urlCats, cat)
	}

	for _, raw := range w.UserCats {
		cat, err := newCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating category %s : %s", raw, err)
		}
		w.userCats = append(w.userCats, cat)
	}

	return nil
}

func (w *WebLog) initParser() error {
	line, err := simpletail.ReadLastLine(w.Path)

	if err != nil {
		return err
	}

	var p parser

	if len(w.CustomParser) > 0 {
		p, err = newParser(string(line), w.CustomParser)
	} else {
		p, err = newParser(string(line), csvDefaultPatterns...)
	}

	if err != nil {
		return err
	}

	w.parser = p

	return nil
}

func (w *WebLog) Init() bool {
	if err := w.initParser(); err != nil {
		w.Error(err)
		return false
	}

	if err := w.initFilter(); err != nil {
		w.Error(err)
		return false
	}

	if err := w.initCategories(); err != nil {
		w.Error(err)
		return false
	}

	return true
}

func (WebLog) Check() bool {
	return false
}

func (WebLog) Charts() *modules.Charts {
	return nil
}

func (w *WebLog) Collect() map[string]int64 {
	w.collect <- struct{}{}
	<-w.done

	return w.metrics
}

func (w *WebLog) parseLoop() {
LOOP:
	for {
		select {
		case <-w.stop:
			w.cleanup()
			break LOOP
		case <-w.collect:
			w.copyMetrics()
			w.updated = false
			w.done <- struct{}{}
		case line := <-w.tail.Lines:
			if !w.filter.match(line.Text) {
				continue
			}
			w.updated = true
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
	gm, ok := w.parse(line)

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

	if _, ok := gm.lookup("user_defined"); ok && len(w.userCats) > 0 {
		w.userCategory(gm)
	}

	if _, ok := gm.lookup("bytes_sent"); ok {
		w.bytesSent(gm)
	}

	if _, ok := gm.lookup("resp_length"); ok {
		w.respLength(gm)
	}

	if _, ok := gm.lookup("address"); ok {
		w.ipProto(gm)
	}

	if w.DoPerURLCharts && w.matchedURL != "" {
		w.urlCategoryStats(gm)
	}

}

func (w *WebLog) codeFam(gm groupMap) {
	fam := gm.get("code")[:1] + "xx"

	if _, ok := w.metrics[fam]; ok {
		w.metrics[fam]++
	} else {
		w.metrics["0xx"]++
	}
}

func (w *WebLog) codeDetailed(gm groupMap) {
	code := gm.get("code")

	if _, ok := w.metrics[code]; ok {
		w.metrics[code]++
		return
	}

	var chart *Chart

	if w.DoCodesAggregate {
		chart = w.charts.Get(responseCodesDetailed.ID)
	} else {
		v := "other"
		if code[0] <= 53 {
			v = code[:1] + "xx"
		}
		chart = w.charts.Get(responseCodesDetailed.ID + "_" + v)
	}

	_ = chart.AddDim(&Dim{
		ID:   code,
		Algo: modules.Incremental,
	})
	chart.MarkNotCreated()

	w.metrics[code]++
}

func (w *WebLog) codeStatus(gm groupMap) {
	code, fam := gm.get("code"), gm.get("code")[:1]

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

func (w *WebLog) request(gm groupMap) {
	request := gm.get("request")

	// FIX ME: separate parser for request field
	if request != "" {
		gm, _ = w.parse(request)
	}

	if _, ok := gm.lookup("method"); ok {
		w.httpMethod(gm)
	}

	if _, ok := gm.lookup("url"); ok {
		w.urlCategory(gm)
	}

	if _, ok := gm.lookup("version"); ok {
		w.httpVersion(gm)
	}
}

func (w *WebLog) httpMethod(gm groupMap) {
	method := gm.get("method")

	if _, ok := w.metrics[method]; !ok {
		chart := w.charts.Get(requestsPerHTTPMethod.ID)
		_ = chart.AddDim(&Dim{
			ID:   method,
			Algo: modules.Incremental,
		})
		chart.MarkNotCreated()
	}

	w.metrics[method]++
}

func (w *WebLog) urlCategory(gm groupMap) {
	url := gm.get("url")

	for _, v := range w.urlCats {
		if v.match(url) {
			w.metrics[v.name]++
			w.matchedURL = v.name
			return
		}
	}
	w.matchedURL = ""
	w.metrics["url_other"]++
}

func (w *WebLog) userCategory(gm groupMap) {
	userDefined := gm.get("user_defined")

	for _, cat := range w.userCats {
		if cat.match(userDefined) {
			w.metrics[cat.name]++
			return
		}
	}
	w.metrics["user_defined_other"]++
}

func (w *WebLog) httpVersion(gm groupMap) {
	version := gm.get("version")

	dimID := strings.Replace(gm.get("version"), ".", "_", 1)

	if _, ok := w.metrics[dimID]; !ok {
		chart := w.charts.Get(requestsPerHTTPVersion.ID)
		_ = chart.AddDim(&Dim{
			ID:   dimID,
			Name: version,
			Algo: modules.Incremental,
		})
		chart.MarkNotCreated()
	}

	w.metrics[dimID]++
}

func (w *WebLog) bytesSent(gm groupMap) {
	v := gm.get("bytes_sent")

	w.metrics["bytes_sent"] += toInt(v)
}

func (w *WebLog) respLength(gm groupMap) {
	v := gm.get("resp_length")

	w.metrics["resp_length"] += toInt(v)
}

func (w *WebLog) respTime(gm groupMap) {

}

func (w *WebLog) respTimeUpstream(gm groupMap) {

}

func (w *WebLog) ipProto(gm groupMap) {
	var (
		address = gm.get("address")
		proto   = "ipv4"
	)

	if strings.Contains(address, ":") {
		proto = "ipv6"
	}

	w.metrics["req_"+proto]++

	if _, ok := w.uniqIPs[address]; !ok {
		w.uniqIPs[address] = true
		w.metrics["unique_cur_"+proto]++
	}

	if !w.DoAllTimeIPs {
		return
	}

	if _, ok := w.uniqIPsAllTime[address]; !ok {
		w.uniqIPsAllTime[address] = true
		w.metrics["unique_all_"+proto]++
	}

}

func (w *WebLog) urlCategoryStats(gm groupMap) {
	code := gm.get("code")
	id := w.matchedURL + "_" + code

	if _, ok := w.metrics[id]; !ok {
		chart := w.charts.Get(responseCodesDetailed.ID + "_" + w.matchedURL)
		_ = chart.AddDim(&Dim{
			ID:   id,
			Name: code,
			Algo: modules.Incremental,
		})
		chart.MarkNotCreated()
	}

	w.metrics[id]++

	if v, ok := gm.lookup("bytes_sent"); ok {
		w.metrics[w.matchedURL+"_bytes_sent"] += toInt(v)
	}

	if v, ok := gm.lookup("resp_length"); ok {
		w.metrics[w.matchedURL+"_resp_length"] += toInt(v)
	}

	//if id, ok := gm.Lookup("resp_time"); ok {
	//	w.timings.get(id).set(id)
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
