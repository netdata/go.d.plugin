package web_log

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/l2isbad/go.d.plugin/charts/raw"
	"github.com/l2isbad/go.d.plugin/modules"
	"github.com/l2isbad/go.d.plugin/shared/log_helper"
)

type (
	Charts      = raw.Charts
	Order       = raw.Order
	Definitions = raw.Definitions
	Chart       = raw.Chart
	Options     = raw.Options
	Dimensions  = raw.Dimensions
	Dimension   = raw.Dimension
)

var uCharts = Charts{
	Order: Order{
		"response_statuses", "response_codes", "bandwidth", "response_time", "response_time_upstream",
		"requests_per_url", "requests_per_user_defined", "http_method", "http_version",
		"requests_per_ipproto", "clients", "clients_all",
	},
	Definitions: Definitions{
		Chart{
			ID:      "response_statuses",
			Options: Options{"Response Statuses", "requests/s", "responses", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"successful_requests", "success", raw.Incremental},
				Dimension{"server_errors", "error", raw.Incremental},
				Dimension{"redirects", "redirect", raw.Incremental},
				Dimension{"bad_requests", "bad", raw.Incremental},
				Dimension{"other_requests", "other", raw.Incremental},
			},
		},
		Chart{
			ID:      "response_codes",
			Options: Options{"Response Codes", "requests/s", "responses", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"2xx", "", raw.Incremental},
				Dimension{"5xx", "", raw.Incremental},
				Dimension{"3xx", "", raw.Incremental},
				Dimension{"4xx", "", raw.Incremental},
				Dimension{"1xx", "", raw.Incremental},
				Dimension{"0xx", "", raw.Incremental},
				Dimension{"unmatched", "", raw.Incremental},
			},
		},
		Chart{
			ID:      "bandwidth",
			Options: Options{"Bandwidth", "kilobits/s", "bandwidth", "", raw.Area},
			Dimensions: Dimensions{
				Dimension{"resp_length", "received", raw.Incremental},
				Dimension{"bytes_sent", "sent", raw.Incremental},
			},
		},
		Chart{
			ID:      "response_time",
			Options: Options{"Processing Time", "milliseconds", "timings", "", raw.Area},
			Dimensions: Dimensions{
				Dimension{"resp_time_min", "min", raw.Incremental, 1, 1000},
				Dimension{"resp_time_max", "max", raw.Incremental, 1, 1000},
				Dimension{"resp_time_avg", "avg", raw.Incremental, 1, 1000},
			},
		},
		Chart{
			ID:      "response_time_upstream",
			Options: Options{"Processing Time Upstream", "milliseconds", "timings", "", raw.Area},
			Dimensions: Dimensions{
				Dimension{"resp_time_upstream_min", "min", raw.Incremental, 1, 1000},
				Dimension{"resp_time_upstream_max", "max", raw.Incremental, 1, 1000},
				Dimension{"resp_time_upstream_avg", "avg", raw.Incremental, 1, 1000},
			},
		},
		Chart{
			ID:      "clients",
			Options: Options{"Current Poll Unique Client IPs", "unique ips", "clients", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"unique_cur_ipv4", "ipv4", raw.Incremental},
				Dimension{"unique_cur_ipv6", "ipv6", raw.Incremental},
			},
		},
		Chart{
			ID:      "clients_all",
			Options: Options{"All Time Unique Client IPs", "unique ips", "clients", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"unique_tot_ipv4", "ipv4"},
				Dimension{"unique_tot_ipv6", "ipv6"},
			},
		},
		Chart{
			ID:      "http_method",
			Options: Options{"Requests Per HTTP Method", "requests/s", "http methods", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"GET", "", raw.Incremental},
			},
		},
		Chart{
			ID:         "http_version",
			Options:    Options{"Requests Per HTTP Version", "requests/s", "http versions", "", raw.Stacked},
			Dimensions: Dimensions{},
		},
		Chart{
			ID:         "requests_per_ipproto",
			Options:    Options{"Requests Per IP Protocol", "requests/s", "ip protocols", "", raw.Stacked},
			Dimensions: Dimensions{},
		},
		Chart{
			ID:      "requests_per_url",
			Options: Options{"Requests Per Url", "requests/s", "urls", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"url_other", "other", raw.Incremental},
			},
		},
		Chart{
			ID:      "requests_per_user_defined",
			Options: Options{"Requests Per User Defined Pattern", "requests/s", "user defined", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"user_pattern_other", "other", raw.Incremental},
			},
		},
	},
}

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
	categories  map[string]*regexp.Regexp
	categoriesU map[string]*regexp.Regexp
	include     *regexp.Regexp
	exclude     *regexp.Regexp
	parser      *regexp.Regexp
}

type filter struct {
	Include string `toml:"include"`
	Exclude string `toml:"exclude"`
}

type WebLog struct {
	modules.Charts
	modules.Logger
	Path           string            `toml:"path, required"`
	Filter         filter            `toml:"filter"`
	Categories     map[string]string `toml:"categories"`
	CategoriesU    map[string]string `toml:"categories_user_defined"`
	CategoryCharts bool              `toml:"per_category_charts"`
	DetRespCodes   bool              `toml:"detailed_response_codes"`
	DetRespCodesA  bool              `toml:"detailed_response_codes_aggregate"`

	*log_helper.FileReader
	regex   regex
	uniqIPs map[string]bool
	data    map[string]int64
}

var reRequest = regexp.MustCompile(`(?P<method>[A-Z]+) (?P<url>[^ ]+) [A-Z]+/(?P<http_version>\d(?:.\d)?)`)

func (w *WebLog) Check() bool {
	v, err := log_helper.NewFileReader(w.Path)
	if err != nil {
		w.Error(err)
		return false
	}
	w.FileReader = v

	if len(w.Categories) != 0 {
		for k, v := range w.Categories {
			r, err := regexp.Compile(v)
			if err != nil {
				w.Error(err)
				return false
			}
			w.regex.categories["url_"+k] = r
		}
	}

	if len(w.CategoriesU) != 0 {
		for k, v := range w.CategoriesU {
			r, err := regexp.Compile(v)
			if err != nil {
				w.Error(err)
				return false
			}
			w.regex.categoriesU["user_defined_"+k] = r
		}
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

	w.createCharts()

	return true
}

func (w *WebLog) createCharts() {
	c := uCharts.Copy()
	if w.DetRespCodes && w.DetRespCodesA {
		n := raw.NewChart(
			"detailed_response_codes",
			Options{"Detailed Response Codes", "requests/s", "responses", "", raw.Stacked})
		c.AddChart(n, true)
	}

	if w.DetRespCodes && !w.DetRespCodesA {
		for _, v := range []string{"1xx", "2xx", "3xx", "4xx", "5xx", "other"} {
			n := raw.NewChart(
				fmt.Sprintf("detailed_response_codes_%s", v),
				Options{fmt.Sprintf("Detailed Response Codes %s", v), "requests/s", "responses", "", raw.Stacked})
			c.AddChart(n, true)
		}
	}

	if len(w.regex.categories) != 0 {
		for key := range w.regex.categories {
			c.GetChartByID("requests_per_url").AddDim(Dimension{key, key[4:], raw.Incremental})
			w.data[key] = 0
			w.data["url_other"] = 0
		}
	}

	if len(w.regex.categoriesU) != 0 {
		for key := range w.regex.categoriesU {
			c.GetChartByID("requests_per_user_defined").AddDim(Dimension{key, key[13:], raw.Incremental})
			w.data[key] = 0
			w.data["user_pattern_other"] = 0
		}
	}
	w.AddMany(c)
}

func (w *WebLog) GetData() *map[string]int64 {
	v, err := w.GetRawData()

	if err != nil {
		if err == log_helper.ErrSizeNotChanged {
			w.Warning(err)
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

		code, codeFam := md["code"], md["code"][:1]

		if _, ok := w.data[codeFam+"xx"]; ok {
			w.data[codeFam+"xx"]++
		} else {
			w.data["0xx"]++
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

		if v, ok := md["user_defined"]; ok && w.regex.categoriesU != nil {
			w.getDataPerPattern(v, "user_pattern_other", w.regex.categoriesU)
		}

		if v, ok := md["request"]; ok {
			w.getDataPerRequest(v)
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

func (w *WebLog) getDataPerAddress(s string, m map[string]bool) {
	var proto = "ipv4"

	if strings.Contains(s, ":") {
		proto = "ipv6"
	}
	w.data["req_"+proto]++

	if _, ok := m[s]; !ok {
		m[s] = true
	}
	w.data["unique_cur_"+proto]++

	if _, ok := w.uniqIPs[s]; !ok {
		w.uniqIPs[s] = true
	}
	w.data["unique_tot_"+proto]++
}

func (w *WebLog) getDataPerRequest(req string) {
	m := reRequest.FindStringSubmatch(req)
	if m == nil {
		return
	}

	if w.regex.categories != nil {
		w.getDataPerPattern(m[2], "url_other", w.regex.categories)
	}

	if _, ok := w.data[m[1]]; !ok {
		w.GetChartByID("http_method").AddDim(Dimension{m[1], "", raw.Incremental})
	}
	w.data[m[1]]++

	dimID := strings.Replace(m[3], ".", "_", 1)
	if _, ok := w.data[dimID]; !ok {
		w.GetChartByID("http_version").AddDim(Dimension{dimID, m[3], raw.Incremental})
	}
	w.data[dimID]++
}

func (w *WebLog) getDataPerPattern(r, other string, p map[string]*regexp.Regexp) string {
	for k, v := range p {
		if v.MatchString(r) {
			w.data[k]++
			return k
		}
	}
	w.data[other]++
	return ""
}

func (w *WebLog) getDataPerCode(code string) {
	if _, ok := w.data[code]; ok {
		w.data[code]++
		return
	}

	if w.DetRespCodesA {
		w.GetChartByID("detailed_response_codes").AddDim(Dimension{code, "", raw.Incremental})
		w.data[code] = 0
		return
	}
	var v = "other"
	if code[0] <= 53 {
		v = code[:1] + "xx"
	}
	w.GetChartByID("detailed_response_codes_" + v).AddDim(Dimension{code, "", raw.Incremental})
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
	var l int
	if s != "-" {
		if v, err := strconv.Atoi(s); err != nil {
			l = v
		}
	}
	return l
}

func init() {
	f := func() modules.Module {
		return &WebLog{
			DetRespCodes:  true,
			DetRespCodesA: true,
			regex: regex{
				categories:  make(map[string]*regexp.Regexp),
				categoriesU: make(map[string]*regexp.Regexp)},
			uniqIPs: make(map[string]bool),
			data: map[string]int64{
				"bytes_sent": 0, "resp_length": 0, "resp_time_min": 0, "resp_time_max": 0,
				"resp_time_avg": 0, "resp_time_upstream_min": 0, "resp_time_upstream_max": 0,
				"resp_time_upstream_avg": 0, "unique_cur_ipv4": 0, "unique_cur_ipv6": 0,
				"2xx": 0, "5xx": 0, "3xx": 0, "4xx": 0, "1xx": 0, "0xx": 0, "unmatched": 0, "req_ipv4": 0,
				"req_ipv6": 0, "unique_tot_ipv4": 0, "unique_tot_ipv6": 0, "successful_requests": 0,
				"redirects": 0, "bad_requests": 0, "server_errors": 0, "other_requests": 0, "GET": 0,
			},
		}
	}
	modules.Add(f)
}
