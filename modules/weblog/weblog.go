package weblog

import (
	"github.com/netdata/go.d.plugin/pkg/logs"
	logsparse "github.com/netdata/go.d.plugin/pkg/logs/parse"
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/metrics"
	"time"

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
			Parser:                 defaultParserConfig,
			TimeMultiplier:         time.Second.Seconds(),
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
	rawCategory struct {
		Name  string `yaml:"name"`
		Match string `yaml:"match"`
	}

	Config struct {
		Parser                 logsparse.Config   `yaml:",inline"`
		TimeMultiplier         float64            `yaml:"time_multiplies"`
		Path                   string             `yaml:"path" validate:"required"`
		ExcludePath            string             `yaml:"exclude_path"`
		Filter                 matcher.SimpleExpr `yaml:"filter"`
		URLCategories          []rawCategory      `yaml:"categories"`
		UserCategories         []rawCategory      `yaml:"user_categories"`
		Histogram              []float64          `yaml:"histogram"`
		AggregateResponseCodes bool               `yaml:"aggregate_response_codes"`
	}

	WebLog struct {
		module.Base
		Config `yaml:",inline"`
		charts *module.Charts

		file   *logs.Reader
		parser logsparse.Parser
		line   *LogLine

		metrics        *MetricsData
		filter         matcher.Matcher
		urlCategories  []*category
		userCategories []*category

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
)

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

func (w *WebLog) Check() bool {
	// Note: these inits are here to make auto detection retry working
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
