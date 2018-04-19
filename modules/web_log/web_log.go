package web_log

import (
	"regexp"
	"strconv"
	"strings"
	"errors"
	"fmt"

	"github.com/l2isbad/go.d.plugin/charts/raw"
	"github.com/l2isbad/go.d.plugin/modules"
	"github.com/l2isbad/go.d.plugin/shared"
	"github.com/l2isbad/go.d.plugin/shared/log_helper"
)

const (
	keyAddress     = "address"
	keyCode        = "code"
	keyRequest     = "request"
	keyUserDefined = "user_defined"
	keyBytesSent   = "bytes_sent"
	keyRespTime    = "resp_time"
	keyRespTimeUp  = "resp_time_upstream"
	keyRespLen     = "resp_length"

	keyHTTPMethod = "method"
	keyURL        = "url"
	keyHTTPVer    = "http_version"
)

type regex struct {
	URLCat  categories
	UserCat categories
	parser  *regexp.Regexp
}

type WebLog struct {
	modules.Charts
	modules.Logger
	Path            string        `toml:"path, required"`
	RawFilter       rawFilter     `toml:"filter"`
	RawURLCat       rawCategories `toml:"categories"`
	RawUserCat      rawCategories `toml:"user_defined"`
	RawCustomParser string        `toml:"custom_log_format"`
	DoChartURLCat   bool          `toml:"per_category_charts"`
	DoDetailCodes   bool          `toml:"detailed_response_codes"`
	DoDetailCodesA  bool          `toml:"detailed_response_codes_aggregate"`
	DoClientsAll    bool          `toml:"clients_all_time"`

	filter
	*log_helper.FileReader
	regex   regex
	uniqIPs map[string]bool
	data    map[string]int64
}

func (w *WebLog) Check() bool {
	v, err := log_helper.NewFileReader(w.Path)
	if err != nil {
		w.Error(err)
		return false
	}
	w.FileReader = v

	for _, v := range w.RawURLCat {
		re, err := regexp.Compile(v.re)
		if err != nil {
			w.Error(err)
			return false
		}
		w.regex.URLCat.add(v.name, re)
	}
	for _, v := range w.RawUserCat {
		re, err := regexp.Compile(v.re)
		if err != nil {
			w.Error(err)
			return false
		}
		w.regex.UserCat.add(v.name, re)
	}

	if f, err := getFilter(w.RawFilter); err != nil {
		w.Error(err)
		return false
	} else {
		w.filter = f
	}

	line, err := log_helper.ReadLastLine(w.Path)
	if err != nil {
		w.Error(err)
		return false
	}

	if re, err := findParser(w.RawCustomParser, line); err != nil {
		w.Error(err)
		return false
	} else {
		w.regex.parser =re
	}

	w.addCharts()

	return true
}

func (w *WebLog) GetData() *map[string]int64 {
	v, err := w.GetRawData()

	if err != nil {
		if err == log_helper.ErrSizeNotChanged {
			return &w.data
		}
		return nil
	}

	uniqIPs := make(map[string]bool)
	tr, tu := timings{name: "resp_time", min: -1}, timings{name: "resp_time_upstream", min: -1}

	for row := range v {
		if w.filter != nil && !w.filter.match(row) {
			continue
		}

		m := w.regex.parser.FindStringSubmatch(row)
		if m == nil {
			w.data["unmatched"]++
			continue
		}

		mm := createMatchMap(w.regex.parser.SubexpNames(), m)

		var URLCat string

		if v, ok := mm[keyRequest]; ok {
			URLCat = w.dataFromRequest(v)
		}

		if v, ok := mm[keyUserDefined]; ok && w.regex.UserCat.active() {
			w.reqPerCategory(v, w.regex.UserCat)
		}

		code, codeFam := mm[keyCode], mm[keyCode][:1]

		if _, ok := w.data[codeFam+"xx"]; ok {
			w.data[codeFam+"xx"]++
		} else {
			w.data["0xx"]++
		}

		if URLCat != "" && w.DoChartURLCat {
			w.dataPerCategory(URLCat, mm)
		}

		if w.DoDetailCodes {
			w.reqPerCode(code)
		}

		w.reqPerCodeFam(code)

		if v, ok := mm[keyRespTime]; ok {
			tr.set(v)
		}

		if v, ok := mm[keyRespTimeUp]; ok {
			tu.set(v)
		}

		if v, ok := mm[keyAddress]; ok {
			w.reqPerIPProto(v, uniqIPs)
		}

		if v, ok := mm[keyBytesSent]; ok {
			w.data["bytes_sent"] += int64(strToInt(v))
		}

		if v, ok := mm[keyRespLen]; ok {
			w.data["resp_length"] += int64(strToInt(v))
		}
	}

	for _, v := range []*timings{&tr, &tu} {
		if v.active() {
			w.data[v.name+"_min"] += int64(v.min)
			w.data[v.name+"_avg"] += int64(v.sum / v.count)
			w.data[v.name+"_max"] += int64(v.max)
		}
	}

	return &w.data
}

func (w *WebLog) reqPerCategory(url string, c categories) string {
	for _, v := range c.list {
		if v.re.MatchString(url) {
			w.data[v.fullname]++
			return v.fullname
		}
	}
	w.data[c.other()]++
	return ""
}

func (w *WebLog) reqPerIPProto(address string, uniqIPs map[string]bool) {
	var proto = "ipv4"

	if strings.Contains(address, ":") {
		proto = "ipv6"
	}
	w.data["req_"+proto]++

	if _, ok := uniqIPs[address]; !ok {
		uniqIPs[address] = true
		w.data["unique_cur_"+proto]++
	}

	if !w.DoClientsAll {
		return
	}

	if _, ok := w.uniqIPs[address]; !ok {
		w.uniqIPs[address] = true
		w.data["unique_tot_"+proto]++
	}
}

func (w *WebLog) reqPerCode(code string) {
	if _, ok := w.data[code]; ok {
		w.data[code]++
		return
	}

	if w.DoDetailCodesA {
		w.GetChartByID(chartDetRespCodes).AddDim(Dimension{code, "", raw.Incremental})
		w.data[code] = 0
		return
	}
	var v = "other"
	if code[0] <= 53 {
		v = code[:1] + "xx"
	}
	w.GetChartByID(chartDetRespCodes + "_" + v).AddDim(Dimension{code, "", raw.Incremental})
	w.data[code] = 0
}

func (w *WebLog) reqPerCodeFam(code string) {
	f := code[:1]
	switch {
	case f == "2", code == "304", f == "1":
		w.data["successful_requests"]++
	case f == "3":
		w.data["redirects"]++
	case f == "4":
		w.data["bad_requests"]++
	case f == "5":
		w.data["server_errors"]++
	default:
		w.data["other_requests"]++
	}
}

func (w *WebLog) dataFromRequest(req string) (URLCat string) {
	m := reRequest.FindStringSubmatch(req)
	if m == nil {
		return
	}
	mm := createMatchMap(reRequest.SubexpNames(), m)

	if w.regex.URLCat.active() {
		if v := w.reqPerCategory(mm[keyURL], w.regex.URLCat); v != "" {
			URLCat = v
		}
	}

	if _, ok := w.data[mm[keyHTTPMethod]]; !ok {
		w.GetChartByID(chartHTTPMethod).AddDim(Dimension{mm[keyHTTPMethod], "", raw.Incremental})
	}
	w.data[mm[keyHTTPMethod]]++

	dimID := strings.Replace(mm[keyHTTPVer], ".", "_", 1)
	if _, ok := w.data[dimID]; !ok {
		w.GetChartByID(chartHTTPVer).AddDim(Dimension{dimID, mm[keyHTTPVer], raw.Incremental})
	}
	w.data[dimID]++
	return
}

func (w *WebLog) dataPerCategory(fullname string, mm map[string]string) {
	code := mm[keyCode]
	dimID := code + "_" + fullname
	if _, ok := w.data[dimID]; !ok {
		w.GetChartByID(chartDetRespCodes + "_" + fullname).AddDim(Dimension{dimID, code, raw.Incremental})
	}
	w.data[dimID]++

	if v, ok := mm[keyBytesSent]; ok {
		w.data["bytes_sent_"+fullname] += int64(strToInt(v))
	}
}

func findParser(custom string, line []byte) (*regexp.Regexp, error) {
	if custom == "" {
		for _, p := range patterns {
			if p.Match(line) {
				return p, nil
			}
		}
		return nil, errors.New("can not find appropriate regex")
	}
	r, err := regexp.Compile(custom)
	if err != nil {
		return nil, err
	}

	if !r.Match(line) {
		return nil, errors.New("custom regex match fails")
	}

	if !shared.StringSlice(r.SubexpNames()).Include(keyCode) {
		return nil, fmt.Errorf("custom regex match ok, but mandatory key '%s' is missing", keyCode)
	}

	return r, nil
}

func strToInt(s string) int {
	if s != "-" {
		if v, err := strconv.Atoi(s); err != nil {
			return v
		}
	}
	return 0
}

func createMatchMap(keys, values []string) map[string]string {
	mm := make(map[string]string)
	for idx, v := range keys[1:] {
		mm[v] = values[idx+1]
	}
	return mm
}

func init() {
	f := func() modules.Module {
		return &WebLog{
			DoDetailCodes:  true,
			DoDetailCodesA: true,
			DoChartURLCat:  true,
			uniqIPs:        make(map[string]bool),
			regex: regex{
				URLCat:  categories{prefix: "url"},
				UserCat: categories{prefix: "user_defined"},
			},
			data: map[string]int64{
				"bytes_sent":             0,
				"resp_length":            0,
				"resp_time_min":          0,
				"resp_time_max":          0,
				"resp_time_avg":          0,
				"resp_time_upstream_min": 0,
				"resp_time_upstream_max": 0,
				"resp_time_upstream_avg": 0,
				"unique_cur_ipv4":        0,
				"unique_cur_ipv6":        0,
				"unique_tot_ipv4":        0,
				"unique_tot_ipv6":        0,
				"2xx":                    0,
				"5xx":                    0,
				"3xx":                    0,
				"4xx":                    0,
				"1xx":                    0,
				"0xx":                    0,
				"unmatched":              0,
				"req_ipv4":               0,
				"req_ipv6":               0,
				"successful_requests":    0,
				"redirects":              0,
				"bad_requests":           0,
				"server_errors":          0,
				"other_requests":         0,
				"GET":                    0, // GET should be green on the dashboard
			},
		}
	}
	modules.Add(f)
}
