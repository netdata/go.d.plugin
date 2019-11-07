package weblog

import (
	"github.com/netdata/go.d.plugin/pkg/logs"
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
	cfg := logs.ParserConfig{
		LogType: typeAuto,
		CSV: logs.CSVConfig{
			Delimiter:  ' ',
			CheckField: checkCSVFormatField,
		},
		LTSV: logs.LTSVConfig{
			FieldDelimiter: '\t',
			ValueDelimiter: ':',
		},
		RegExp: logs.RegExpConfig{},
	}

	return &WebLog{
		Config: Config{
			GroupRespCodes: false,
			Histogram:      metrics.DefBuckets,
			Parser:         cfg,
		},
	}
}

type (
	userPattern struct {
		Name  string `yaml:"name"`
		Match string `yaml:"match"`
	}
	customField struct {
		Name     string        `yaml:"name"`
		Patterns []userPattern `yaml:"patterns"`
	}

	Config struct {
		Parser         logs.ParserConfig  `yaml:",inline"`
		Path           string             `yaml:"path"`
		ExcludePath    string             `yaml:"exclude_path"`
		Filter         matcher.SimpleExpr `yaml:"filter"`
		URLPatterns    []userPattern      `yaml:"url_patterns"`
		CustomFields   []customField      `yaml:"custom_fields"`
		Histogram      []float64          `yaml:"histogram"`
		GroupRespCodes bool               `yaml:"group_response_codes"`
	}

	WebLog struct {
		module.Base
		Config `yaml:",inline"`

		file         *logs.Reader
		parser       logs.Parser
		line         *logLine
		filter       matcher.Matcher
		urlPatterns  []*pattern
		customFields map[string][]*pattern

		mx     *metricsData
		charts *module.Charts
	}
)

func (w *WebLog) Init() bool {
	if err := w.initFilter(); err != nil {
		w.Error(err)
		return false
	}

	if err := w.initURLPatterns(); err != nil {
		w.Error(err)
		return false
	}

	if err := w.initCustomFields(); err != nil {
		w.Error(err)
		return false
	}

	w.initLogLine()
	w.mx = newMetricsData(w.Config)
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
	if w.charts == nil {
		w.charts = w.createCharts(w.line)
	}
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
