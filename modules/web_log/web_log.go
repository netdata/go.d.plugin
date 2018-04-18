package web_log

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/l2isbad/go.d.plugin/charts/raw"
	"github.com/l2isbad/go.d.plugin/modules"
	"github.com/l2isbad/go.d.plugin/shared/log_helper"
)

type timings struct {
	name  string
	min   int
	max   int
	sum   int
	count int
}

func (t *timings) set(s string) {
	var n int
	switch {
	case s == "0.000":
		n = 0
	case strings.Contains(s, "."):
		if v, err := strconv.ParseFloat(s, 10); err != nil {
			n = int(v * 1e6)
		}
	default:
		if v, err := strconv.Atoi(s); err != nil {
			n = v
		}
	}

	if t.min == -1 {
		t.min = n
	}
	if n > t.max {
		t.max = n
	} else if n < t.min {
		t.min = n
	}
	t.sum += n
	t.count++
}

func (t *timings) active() bool {
	return t.min != -1
}

type regex struct {
	URLCat  categories
	UserCat categories
	include *regexp.Regexp
	exclude *regexp.Regexp
	parser  *regexp.Regexp
}

type filter struct {
	Include string `toml:"include"`
	Exclude string `toml:"exclude"`
}

type WebLog struct {
	modules.Charts
	modules.Logger
	Path          string        `toml:"path, required"`
	Filter        filter        `toml:"filter"`
	RawURLCat     rawCategories `toml:"categories"`
	RawUserCat    rawCategories `toml:"user_defined"`
	ChartURLCat   bool          `toml:"per_category_charts"`
	DetRespCodes  bool          `toml:"detailed_response_codes"`
	DetRespCodesA bool          `toml:"detailed_response_codes_aggregate"`

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

	if w.Filter.Include != "" {
		r, err := regexp.Compile(w.Filter.Include)
		if err != nil {
			w.Error(err)
			return false
		}
		w.regex.include = r
	}

	if w.Filter.Exclude != "" {
		r, err := regexp.Compile(w.Filter.Exclude)
		if err != nil {
			w.Error(err)
			return false
		}
		w.regex.exclude = r
	}

	line, err := log_helper.ReadLastLine(w.Path)
	if err != nil {
		w.Error(err)
		return false
	}

	var found bool

	for _, p := range patterns {
		if p.Match(line) {
			w.regex.parser = p
			found = true
			break
		}
	}

	if !found {
		w.Error("can not find appropriate regex")
		return false
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
		if w.regex.include != nil && !w.regex.include.MatchString(row) {
			continue
		}
		if w.regex.exclude != nil && w.regex.exclude.MatchString(row) {
			continue
		}

		m := w.regex.parser.FindStringSubmatch(row)
		if m == nil {
			w.data["unmatched"]++
			continue
		}

		md := make(map[string]string)
		for idx, v := range w.regex.parser.SubexpNames()[1:] {
			md[v] = m[idx+1]
		}

		var URLCat string

		if v, ok := md["request"]; ok {
			URLCat = w.getDataPerRequest(v)
		}

		if v, ok := md["user_defined"]; ok && w.regex.UserCat.active() {
			w.getDataPerCategory(v, w.regex.UserCat)
		}

		code, codeFam := md["code"], md["code"][:1]

		if _, ok := w.data[codeFam+"xx"]; ok {
			w.data[codeFam+"xx"]++
		} else {
			w.data["0xx"]++
		}

		if URLCat != "" && w.ChartURLCat {
			w.perCategoriesCharts(URLCat, md)
		}

		if w.DetRespCodes {
			w.getDataPerCode(code)
		}

		w.getDataPerCodeFam(code)

		if v, ok := md["resp_time"]; ok {
			tr.set(v)
		}

		if v, ok := md["resp_time_upstream"]; ok {
			tu.set(v)
		}

		if v, ok := md["address"]; ok {
			w.getDataPerAddress(v, uniqIPs)
		}

		if v, ok := md["bytes_sent"]; ok {
			w.data["bytes_sent"] += int64(strToInt(v))
		}

		if v, ok := md["resp_length"]; ok {
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

func (w *WebLog) perCategoriesCharts(s string, md map[string]string) {
	code := md["code"]
	if _, ok := w.data[s+"_"+code]; !ok {
		w.GetChartByID(s + "_detailed_response_code").AddDim(Dimension{s + "_" + code, code, raw.Incremental})
	}
	w.data[s+"_"+code]++

	if v, ok := md["bytes_sent"]; ok {
		w.data[s+"_"+"bytes_sent"] += int64(strToInt(v))
	}
}

func (w *WebLog) getDataPerAddress(address string, uniqIPs map[string]bool) {
	var proto = "ipv4"

	if strings.Contains(address, ":") {
		proto = "ipv6"
	}
	w.data["req_"+proto]++

	if _, ok := uniqIPs[address]; !ok {
		uniqIPs[address] = true
		w.data["unique_cur_"+proto]++
	}

	if _, ok := w.uniqIPs[address]; !ok {
		w.uniqIPs[address] = true
		w.data["unique_tot_"+proto]++
	}
}

func (w *WebLog) getDataPerRequest(req string) (URLCat string) {
	// 0: method, 1: url, 2: http version
	m := reRequest.FindStringSubmatch(req)
	if m == nil {
		return
	}

	if w.regex.URLCat.active() {
		if v := w.getDataPerCategory(m[2], w.regex.URLCat); v != "" {
			URLCat = v
		}
	}

	if _, ok := w.data[m[1]]; !ok {
		w.GetChartByID(chartHttpMethod).AddDim(Dimension{m[1], "", raw.Incremental})
	}
	w.data[m[1]]++

	dimID := strings.Replace(m[3], ".", "_", 1)
	if _, ok := w.data[dimID]; !ok {
		w.GetChartByID(chartHttpVersion).AddDim(Dimension{dimID, m[3], raw.Incremental})
	}
	w.data[dimID]++
	return
}

func (w *WebLog) getDataPerCategory(s string, c categories) string {
	for _, v := range c.list {
		if v.re.MatchString(s) {
			w.data[v.fullname]++
			return v.fullname
		}
	}
	w.data[c.other()]++
	return ""
}

func (w *WebLog) getDataPerCode(code string) {
	if _, ok := w.data[code]; ok {
		w.data[code]++
		return
	}

	if w.DetRespCodesA {
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

func (w *WebLog) getDataPerCodeFam(code string) {
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

func strToInt(s string) int {
	if s != "-" {
		if v, err := strconv.Atoi(s); err != nil {
			return v
		}
	}
	return 0
}

func init() {
	f := func() modules.Module {
		return &WebLog{
			DetRespCodes:  true,
			DetRespCodesA: true,
			ChartURLCat:   true,
			uniqIPs:       make(map[string]bool),
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
