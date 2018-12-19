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

		worker: newWorker(),
	}
}

type WebLog struct {
	modules.Base

	Path             string        `yaml:"path" validate:"required"`
	Filter           rawfilter     `yaml:"filter"`
	URLCats          []rawcategory `yaml:"categories"`
	UserCats         []rawcategory `yaml:"user_categories"`
	CustomParser     csvPattern    `yaml:"custom_log_format"`
	Histogram        []int         `yaml:"histogram"`
	DoCodesDetailed  bool          `yaml:"detailed_response_codes"`
	DoCodesAggregate bool          `yaml:"detailed_response_codes_aggregate"`
	DoPerURLCharts   bool          `yaml:"per_category_charts"`
	DoAllTimeIPs     bool          `yaml:"all_time_clients"`

	worker *worker
	charts *modules.Charts
	gm     groupMap
}

func (w *WebLog) Cleanup() {
	w.worker.stop()
}

func (w *WebLog) initFilter() error {
	f, err := newFilter(w.Filter)

	if err != nil {
		return fmt.Errorf("error on creating filter %s: %s", w.Filter, err)
	}

	w.worker.filter = f

	return nil
}

func (w *WebLog) initCategories() error {
	for _, raw := range w.URLCats {
		cat, err := newCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating category %s : %s", raw, err)
		}
		w.worker.urlCats = append(w.worker.urlCats, cat)
		w.worker.metrics[cat.name] = 0
	}

	if w.DoPerURLCharts {
		for _, cat := range w.worker.urlCats {
			w.worker.timings.add(cat.name)
		}
	}
	w.worker.timings.reset()

	for _, raw := range w.UserCats {
		cat, err := newCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating category %s : %s", raw, err)
		}
		w.worker.userCats = append(w.worker.userCats, cat)
		w.worker.metrics[cat.name] = 0
	}

	return nil
}

func (w *WebLog) initHistograms() (err error) {
	if len(w.Histogram) == 0 {
		return nil
	}

	var h histogram

	if h, err = newHistogram(keyRespTimeHistogram, w.Histogram); err != nil {
		return fmt.Errorf("error on creating histogram %v : %s", w.Histogram, err)
	}

	w.worker.histograms[keyRespTimeHistogram] = h

	if h, err = newHistogram(keyRespTimeUpstreamHistogram, w.Histogram); err != nil {
		return fmt.Errorf("error on creating histogram %v : %s", w.Histogram, err)
	}

	w.worker.histograms[keyRespTimeUpstreamHistogram] = h

	for _, h := range w.worker.histograms {
		for _, v := range h {
			w.worker.metrics[v.id] = 0
		}
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

	w.worker.parser = p
	w.gm, _ = p.parse(line)

	return nil
}

func (w *WebLog) Init() bool {
	w.worker.doCodesDetailed = w.DoCodesDetailed
	w.worker.doCodesAggregate = w.DoCodesAggregate
	w.worker.doAllTimeIPs = w.DoAllTimeIPs
	w.worker.doPerURLCharts = w.DoPerURLCharts

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

	if err := w.initHistograms(); err != nil {
		w.Error(err)
		return false
	}

	return true
}

func (w *WebLog) Check() bool {
	t, err := w.worker.tailFactory(w.Path)

	if err != nil {
		w.Errorf("error on creating tail : %s", err)
		return false
	}

	w.worker.tail = t
	w.Infof("used parser : %s", w.worker.parser.info())

	w.createCharts()

	go w.worker.parseLoop()

	return true
}

func (w *WebLog) Collect() map[string]int64 {
	w.worker.pause()
	defer w.worker.unpause()

	for k, v := range w.worker.timings {
		if !v.active() {
			continue
		}
		w.worker.metrics[k+"_min"] += int64(v.min)
		w.worker.metrics[k+"_avg"] += int64(v.avg())
		w.worker.metrics[k+"_max"] += int64(v.max)
	}

	for _, h := range w.worker.histograms {
		for _, v := range h {
			w.worker.metrics[v.id] = int64(v.count)
		}
	}

	w.worker.timings.reset()
	w.worker.uniqIPs = make(map[string]bool)

	// NOTE: don't copy if nothing has changed?
	m := make(map[string]int64)

	for k, v := range w.worker.metrics {
		m[k] = v
	}

	for _, task := range w.worker.chartUpdate {
		_ = w.charts.Get(task.id).AddDim(task.dim)
	}
	w.worker.chartUpdate = w.worker.chartUpdate[:0]

	return m
}
