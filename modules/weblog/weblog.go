package weblog

import (
	"fmt"

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
		Path           string         `yaml:"path" validate:"required"`
		ExcludePath    string         `yaml:"exclude_path"`
		Filter         rawFilter      `yaml:"filter"`
		URLCategories  []RawCategory  `yaml:"categories"`
		UserCategories []RawCategory  `yaml:"user_categories"`
		CustomParser   map[string]int `yaml:"custom_log_format"`
		Histogram      []float64      `yaml:"histogram"`
		DetailedStatus bool           `yaml:"detailed_status"`
	}

	WebLog struct {
		module.Base
		Config `yaml:",inline"`
		charts *module.Charts
		worker *Worker

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

	w.metrics = NewMetricsData(w.Config)

	return true
}

func (w *WebLog) Check() bool {
	panic("TODO")
	//t, err := w.worker.tailFactory(w.Path)
	//
	//if err != nil {
	//	w.Errorf("error on creating tail : %s", err)
	//	return false
	//}
	//
	//w.worker.tail = t
	//w.Infof("used parser : %s", w.worker.parser.info())
	//
	//w.createCharts()
	//
	//go w.worker.parseLoop()
	//
	//return true
}

func (w *WebLog) Collect() map[string]int64 {
	panic("TODO")
	//w.worker.pause()
	//defer w.worker.unpause()
	//
	//for k, v := range w.worker.timings {
	//	if !v.active() {
	//		continue
	//	}
	//	w.worker.metrics[k+"_min"] += int64(v.min)
	//	w.worker.metrics[k+"_avg"] += int64(v.avg())
	//	w.worker.metrics[k+"_max"] += int64(v.max)
	//}
	//
	//for _, h := range w.worker.histograms {
	//	for _, v := range h {
	//		w.worker.metrics[v.id] = int64(v.count)
	//	}
	//}
	//
	//w.worker.timings.reset()
	//w.worker.uniqIPs = make(map[string]bool)
	//
	//m := make(map[string]int64)
	//
	//for k, v := range w.worker.metrics {
	//	m[k] = v
	//}
	//
	//for _, task := range w.worker.chartUpdate {
	//	chart := w.charts.Get(task.id)
	//	_ = chart.AddDim(task.dim)
	//	chart.MarkNotCreated()
	//}
	//w.worker.chartUpdate = w.worker.chartUpdate[:0]
	//
	//return m
}

func (w *WebLog) Cleanup() {
}

func (w *WebLog) initParser() error {
	panic("TODO:")
	//b, err := simpletail.ReadLastLine(w.Path)
	//
	//if err != nil {
	//	return err
	//}
	//
	//line := string(b)
	//var p parser
	//
	//if len(w.CustomParser) > 0 {
	//	p, err = newParser(line, w.CustomParser)
	//} else {
	//	p, err = newParser(line, csvDefaultPatterns...)
	//}
	//
	//if err != nil {
	//	return err
	//}
	//
	//w.worker.parser = p
	//w.gm, _ = p.parse(line)
	//
	//return nil
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
