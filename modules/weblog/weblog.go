package weblog

import (
	"fmt"

	"github.com/netdata/go.d.plugin/modules/weblog/parser"
	"github.com/netdata/go.d.plugin/pkg/logreader"
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/metrics"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("web_log", creator)
}

func New() *WebLog {
	return &WebLog{
		Config: Config{
			AggregateResponseCodes: true,
			Histogram:              metrics.DefBuckets,
			Parser:                 parser.DefaultConfig,
		},
		charts: charts.Copy(),
		chartsCache: chartsCache{
			created:  make(cache),
			vhosts:   make(cache),
			methods:  make(cache),
			codes:    make(cache),
			versions: make(cache),
		},
	}
}

type (
	Config struct {
		Parser                 parser.Config      `yaml:",inline"`
		Path                   string             `yaml:"path" validate:"required"`
		ExcludePath            string             `yaml:"exclude_path"`
		Filter                 matcher.SimpleExpr `yaml:"filter"`
		URLCategories          []RawCategory      `yaml:"categories"`
		UserCategories         []RawCategory      `yaml:"user_categories"`
		Histogram              []float64          `yaml:"histogram"`
		AggregateResponseCodes bool               `yaml:"aggregate_response_codes"`
	}

	WebLog struct {
		module.Base
		Config `yaml:",inline"`
		charts *module.Charts

		file   *logreader.Reader
		parser parser.Parser

		metrics        *MetricsData
		filter         matcher.Matcher
		urlCategories  []*Category
		userCategories []*Category

		collected struct {
			vhost      bool
			client     bool
			method     bool
			uri        bool
			version    bool
			status     bool
			reqSize    bool
			respSize   bool
			respTime   bool
			upRespTime bool
			custom     bool
		}
		chartsCache chartsCache
	}

	cache map[string]struct{}

	chartsCache struct {
		created  cache
		vhosts   cache
		methods  cache
		versions cache
		codes    cache
	}
)

func (c cache) has(v string) bool { _, ok := c[v]; return ok }

func (c cache) add(v string) { c[v] = struct{}{} }

func (c cache) addIfNotExist(v string) (exist bool) {
	if c.has(v) {
		return true
	}
	c.add(v)
	return
}

func (w *WebLog) initFilter() (err error) {
	if w.Filter.Empty() {
		w.filter = matcher.TRUE()
		return
	}

	m, err := w.Filter.Parse()
	if err != nil {
		return fmt.Errorf("error on creating filter %s: %v", w.Filter, err)
	}

	w.filter = m
	return
}

func (w *WebLog) initCategories() error {
	for _, raw := range w.URLCategories {
		cat, err := NewCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating url categories %s: %v", cat, err)
		}
		w.urlCategories = append(w.urlCategories, cat)
	}

	for _, raw := range w.UserCategories {
		cat, err := NewCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating user categories %s: %v", cat, err)
		}
		w.userCategories = append(w.userCategories, cat)
	}

	return nil
}

func (w *WebLog) Init() bool {
	if err := w.initFilter(); err != nil {
		w.Error(err)
		return false
	}

	if err := w.initCategories(); err != nil {
		w.Error(err)
		return false
	}

	w.metrics = NewMetricsData(w.Config)
	return true
}

func (w *WebLog) initLogReader() error {
	file, err := logreader.Open(w.Path, w.ExcludePath, w.Logger)
	if err != nil {
		return fmt.Errorf("error on creating logreader : %v", err)
	}

	w.file = file
	return nil
}

func (w *WebLog) initParser() error {
	lastLine, err := logreader.ReadLastLine(w.file.CurrentFilename(), 0)
	if err != nil {
		return fmt.Errorf("error on reading last line : %v", err)
	}

	w.parser, err = parser.NewParser(w.Config.Parser, w.file, lastLine)
	if err != nil {
		return fmt.Errorf("error on creating parser : %v", err)
	}

	log, err := w.parser.Parse(lastLine)
	if err != nil {
		return fmt.Errorf("error on parsing last line : %v (%s)", err, lastLine)
	}

	if err = log.Verify(); err != nil {
		return fmt.Errorf("error on verifying parsed log line : %v", err)
	}
	return nil
}

func (w *WebLog) Check() bool {
	// Note: these inits are here to make autodetection retry working
	if err := w.initLogReader(); err != nil {
		w.Warning("check failed: ", err)
		return false
	}

	if err := w.initParser(); err != nil {
		w.Warning("check failed: ", err)
		return false
	}
	return true
}

func (w *WebLog) Charts() *module.Charts {
	return w.charts
}

func (w *WebLog) Collect() map[string]int64 {
	mx, err := w.collect()
	if err != nil {
		w.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (w *WebLog) Cleanup() {
	_ = w.file.Close()
}
