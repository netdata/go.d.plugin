package weblog

import (
	"fmt"
	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/simpletail"
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
		DoPerURLCharts:   true,

		reqParser: newCSVParser(csvPattern{
			{"method", 0},
			{"url", 1},
			{"version", 2},
		}),
		stop:  make(chan struct{}),
		pause: make(chan struct{}),
		timings: timings{
			keyResponseTime:         &timing{},
			keyResponseTimeUpstream: &timing{},
		},
		uniqIPs:        make(map[string]bool),
		uniqIPsAllTime: make(map[string]bool),
		metrics: map[string]int64{
			"successful_requests":      0,
			"redirects":                0,
			"bad_requests":             0,
			"server_errors":            0,
			"other_requests":           0,
			"2xx":                      0,
			"5xx":                      0,
			"3xx":                      0,
			"4xx":                      0,
			"1xx":                      0,
			"0xx":                      0,
			"unmatched":                0,
			"bytes_sent":               0,
			"resp_length":              0,
			"resp_time_min":            0,
			"resp_time_max":            0,
			"resp_time_avg":            0,
			"resp_time_upstream_min":   0,
			"resp_time_upstream_max":   0,
			"resp_time_upstream_avg":   0,
			"unique_current_poll_ipv4": 0,
			"unique_current_poll_ipv6": 0,
			"unique_all_time_ipv4":     0,
			"unique_all_time_ipv6":     0,
			"req_ipv4":                 0,
			"req_ipv6":                 0,
			"GET":                      0, // GET should be green on the dashboard
		},
	}
}

type WebLog struct {
	modules.Base

	Path             string        `yaml:"path" validate:"required"`
	Filter           rawFilter     `yaml:"filter"`
	URLCats          []rawCategory `yaml:"categories"`
	UserCats         []rawCategory `yaml:"user_categories"`
	CustomParser     csvPattern    `yaml:"custom_log_format"`
	Histogram        []int         `yaml:"histogram"`
	DoCodesDetailed  bool          `yaml:"detailed_response_codes"`
	DoCodesAggregate bool          `yaml:"detailed_response_codes_aggregate"`
	DoPerURLCharts   bool          `yaml:"per_category_charts"`
	DoAllTimeIPs     bool          `yaml:"all_time_clients"`

	charts *modules.Charts

	tail      follower
	filter    matcher
	parser    parser
	reqParser parser

	gm         groupMap // for creating charts
	matchedURL string

	urlCats  []*category
	userCats []*category

	stop  chan struct{}
	pause chan struct{}

	timings        timings
	uniqIPs        map[string]bool
	uniqIPsAllTime map[string]bool

	metrics map[string]int64
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

	if w.DoPerURLCharts {
		for _, cat := range w.urlCats {
			w.timings.add(cat.name)
		}
	}
	w.timings.reset()

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
	b, err := simpletail.ReadLastLine(w.Path)

	if err != nil {
		return err
	}

	line := string(b)
	var p parser

	if len(w.CustomParser) > 0 {
		p, err = newParser(line, w.CustomParser)
	} else {
		p, err = newParser(line, csvDefaultPatterns...)
	}

	if err != nil {
		return err
	}

	w.parser = p
	w.gm, _ = w.parser.parse(line)

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

func (w *WebLog) Check() bool {
	t, err := newFollower(w.Path)

	if err != nil {
		w.Error(err)
		return false
	}
	w.tail = t
	go w.parseLoop()

	return true
}

func (w *WebLog) Collect() map[string]int64 {
	w.pause <- struct{}{}
	defer func() { <-w.pause }()

	for k, v := range w.timings {
		if !v.active() {
			continue
		}
		fmt.Println(v, 1111111111, k+"_min")
		w.metrics[k+"_min"] += int64(v.min)
		w.metrics[k+"_avg"] += int64(v.avg())
		w.metrics[k+"_max"] += int64(v.max)
	}

	w.timings.reset()

	// NOTE: don't copy if nothing has changed?
	m := make(map[string]int64)

	for k, v := range w.metrics {
		m[k] = v
	}

	return m
}
