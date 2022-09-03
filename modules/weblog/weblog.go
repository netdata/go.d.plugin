// SPDX-License-Identifier: GPL-3.0-or-later

package weblog

import (
	"github.com/netdata/go.d.plugin/pkg/logs"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("web_log", creator)
}

func New() *WebLog {
	cfg := logs.ParserConfig{
		LogType: typeAuto,
		CSV: logs.CSVConfig{
			FieldsPerRecord:  -1,
			Delimiter:        " ",
			TrimLeadingSpace: false,
			CheckField:       checkCSVFormatField,
		},
		LTSV: logs.LTSVConfig{
			FieldDelimiter: "\t",
			ValueDelimiter: ":",
		},
		RegExp: logs.RegExpConfig{},
		JSON:   logs.JSONConfig{},
	}
	return &WebLog{
		Config: Config{
			ExcludePath:    "*.gz",
			GroupRespCodes: true,
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
	customTimeField struct {
		Name      string    `yaml:"name"`
		Histogram []float64 `yaml:"histogram"`
	}

	Config struct {
		Parser           logs.ParserConfig `yaml:",inline"`
		Path             string            `yaml:"path"`
		ExcludePath      string            `yaml:"exclude_path"`
		URLPatterns      []userPattern     `yaml:"url_patterns"`
		CustomFields     []customField     `yaml:"custom_fields"`
		CustomTimeFields []customTimeField `yaml:"custom_time_fields"`
		Histogram        []float64         `yaml:"histogram"`
		GroupRespCodes   bool              `yaml:"group_response_codes"`
	}

	WebLog struct {
		module.Base
		Config `yaml:",inline"`

		file             *logs.Reader
		parser           logs.Parser
		line             *logLine
		urlPatterns      []*pattern
		customFields     map[string][]*pattern
		customTimeFields map[string][]float64

		mx     *metricsData
		charts *module.Charts
	}
)

func (w *WebLog) Init() bool {
	if err := w.createURLPatterns(); err != nil {
		w.Error("init failed: ", err)
		return false
	}

	if err := w.createCustomFields(); err != nil {
		w.Error("init failed: ", err)
		return false
	}

	if err := w.createCustomTimeFields(); err != nil {
		w.Error("init failed: ", err)
		return false
	}

	w.createLogLine()
	w.mx = newMetricsData(w.Config)
	return true
}

func (w *WebLog) Check() bool {
	// Note: these inits are here to make auto-detection retry working
	if err := w.createLogReader(); err != nil {
		w.Warning("check failed: ", err)
		return false
	}

	if err := w.createParser(); err != nil {
		w.Warning("check failed: ", err)
		return false
	}

	if err := w.createCharts(w.line); err != nil {
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
	if w.file != nil {
		_ = w.file.Close()
	}
}
