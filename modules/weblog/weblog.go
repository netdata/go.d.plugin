package weblog

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/simpletail"

	"github.com/netdata/go.d.plugin/pkg/logreader"

	"github.com/netdata/go.d.plugin/pkg/matcher"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		DisabledByDefault: true,
		Create:            func() module.Module { return New() },
	}

	module.Register("web_log", creator)
}

func New() *WebLog {
	return &WebLog{
		Config: Config{
			DetailedStatus: true,
		},
	}
}

type (
	Config struct {
		Path           string        `yaml:"path" validate:"required"`
		ExcludePath    string        `yaml:"exclude_path"`
		Filter         rawFilter     `yaml:"filter"`
		URLCategories  []RawCategory `yaml:"categories"`
		UserCategories []RawCategory `yaml:"user_categories"`
		LogFormat      string        `yaml:"log_format"`
		LogTimeScale   float64       `yaml:"log_time_scale"`
		Histogram      []float64     `yaml:"histogram"`
		DetailedStatus bool          `yaml:"detailed_status"`
	}

	WebLog struct {
		module.Base
		Config `yaml:",inline"`
		charts *module.Charts

		file   *logreader.Reader
		parser *csv.Reader
		format *Format

		metrics        *MetricsData
		filter         matcher.Matcher
		urlCategories  []*Category
		userCategories []*Category
	}
)

func (w *WebLog) Charts() *module.Charts {
	panic("implement me")
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

func (w *WebLog) Check() bool {
	if err := w.initLogReader(); err != nil {
		w.Warning("check failed: ", err)
		return false
	}
	lastLine, err := simpletail.ReadLastLine(w.file.CurrentFilename(), 0)
	if err != nil {
		w.Warning("check failed: ", err)
		return false
	}

	parser := NewLogParser(bytes.NewBuffer(lastLine))
	fields, err := parser.Read()
	if err != nil {
		w.Warning("check failed: ", err)
		return false
	}

	if w.LogFormat != "" {
		w.format = NewFormat(w.LogTimeScale, w.LogFormat)
		if w.format.Match(fields) != nil {
			w.Warning("check failed: ", err)
			return false
		}
	} else {
		w.format = GuessFormat(fields)
		if w.format == nil {
			w.Warning("check failed: cannot determine log format")
			return false
		}
	}

	panic("TODO")
}

func (w *WebLog) Collect() map[string]int64 {
	panic("TODO")
}

func (w *WebLog) Cleanup() {
	w.file.Close()
}

func (w *WebLog) initLogReader() error {
	file, err := logreader.Open(w.Path, w.ExcludePath, w.Logger)
	if err != nil {
		return err
	}
	w.file = file
	w.parser = NewLogParser(file)
	return nil
}

func (w *WebLog) initFilter() (err error) {
	if w.filter, err = NewFilter(w.Filter); err != nil {
		err = fmt.Errorf("error on creating filter %s: %s", w.Filter, err)
	}
	return
}

func (w *WebLog) initCategories() error {
	for _, raw := range w.URLCategories {
		cat, err := NewCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating Category %s : %s", raw, err)
		}
		w.urlCategories = append(w.urlCategories, cat)
	}

	for _, raw := range w.UserCategories {
		cat, err := NewCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating Category %s : %s", raw, err)
		}
		w.userCategories = append(w.userCategories, cat)
	}

	return nil
}
