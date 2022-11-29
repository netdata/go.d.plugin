package prometheus

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/textparse"
)

const (
	precision = 1000
)

func (p *Prometheus) collect() (map[string]int64, error) {
	mfs, err := p.prom.Scrape()
	if err != nil {
		return nil, err
	}

	if mfs.Len() == 0 {
		p.Warningf("endpoint '%s' returned 0 metric families", p.URL)
		return nil, nil
	}

	if p.ExpectedPrefix != "" {
		if !hasPrefix(mfs, p.ExpectedPrefix) {
			return nil, fmt.Errorf("'%s' metrics have no expected prefix (%s)", p.URL, p.ExpectedPrefix)
		}
		p.ExpectedPrefix = ""
	}

	mx := make(map[string]int64)

	for name, mf := range mfs {
		if strings.HasSuffix(name, "_info") {
			continue
		}
		if len(mf.Metrics()) > p.MaxTSPerMetric {
			continue
		}

		switch mf.Type() {
		case textparse.MetricTypeGauge:
			p.collectGauge(mx, mf)
		case textparse.MetricTypeCounter:
			p.collectCounter(mx, mf)
		case textparse.MetricTypeSummary:
			p.collectSummary(mx, mf)
		case textparse.MetricTypeHistogram:
			p.collectHistogram(mx, mf)
		case textparse.MetricTypeUnknown:
			p.collectUntyped(mx, mf)
		}
	}

	return mx, nil
}

func (p *Prometheus) collectGauge(mx map[string]int64, mf *prometheus.MetricFamily) {
	for _, m := range mf.Metrics() {
		if m.Gauge() == nil || math.IsNaN(m.Gauge().Value()) {
			continue
		}

		id := mf.Name() + p.joinLabels(m.Labels())

		if !p.cache[id] {
			p.cache[id] = true
			p.addGaugeChart(id, mf.Name(), mf.Help(), m.Labels())
		}

		mx[id] = int64(m.Gauge().Value() * precision)
	}
}

func (p *Prometheus) collectCounter(mx map[string]int64, mf *prometheus.MetricFamily) {
	for _, m := range mf.Metrics() {
		if m.Counter() == nil || math.IsNaN(m.Counter().Value()) {
			continue
		}

		id := mf.Name() + p.joinLabels(m.Labels())

		if !p.cache[id] {
			p.cache[id] = true
			p.addCounterChart(id, mf.Name(), mf.Help(), m.Labels())
		}

		mx[id] = int64(m.Counter().Value() * precision)
	}
}

func (p *Prometheus) collectSummary(mx map[string]int64, mf *prometheus.MetricFamily) {
	for _, m := range mf.Metrics() {
		if m.Summary() == nil || len(m.Summary().Quantiles()) == 0 {
			continue
		}

		id := mf.Name() + p.joinLabels(m.Labels())

		if !p.cache[id] {
			p.cache[id] = true
			p.addSummaryChart(id, mf.Name(), mf.Help(), m.Labels(), m.Summary().Quantiles())
		}

		for _, v := range m.Summary().Quantiles() {
			dimID := fmt.Sprintf("%s_quantile=%s", id, strconv.FormatFloat(v.Quantile(), 'f', -1, 64))
			mx[dimID] = int64(v.Value() * precision)
		}
	}
}

func (p *Prometheus) collectHistogram(mx map[string]int64, mf *prometheus.MetricFamily) {
	for _, m := range mf.Metrics() {
		if m.Histogram() == nil || len(m.Histogram().Buckets()) == 0 {
			continue
		}

		id := mf.Name() + p.joinLabels(m.Labels())

		if !p.cache[id] {
			p.cache[id] = true
			p.addHistogramChart(id, mf.Name(), mf.Help(), m.Labels(), m.Histogram().Buckets())
		}

		for _, v := range m.Histogram().Buckets() {
			dimID := fmt.Sprintf("%s_bucket=%s", id, strconv.FormatFloat(v.UpperBound(), 'f', -1, 64))
			mx[dimID] = int64(v.CumulativeCount())
		}
	}
}

func (p *Prometheus) collectUntyped(mx map[string]int64, mf *prometheus.MetricFamily) {

}

func (p *Prometheus) joinLabels(labels labels.Labels) string {
	p.sb.Reset()
	for _, lbl := range labels {
		name, value := lbl.Name, lbl.Value
		if name == "" || value == "" {
			continue
		}

		if strings.IndexByte(value, ' ') != -1 {
			value = spaceReplacer.Replace(value)
		}
		if strings.IndexByte(value, '\\') != -1 {
			if value = decodeLabelValue(value); strings.IndexByte(value, '\\') != -1 {
				value = backslashReplacer.Replace(value)
			}
		}

		p.sb.Write([]byte("-" + name + "=" + value))
	}
	return p.sb.String()
}

func decodeLabelValue(value string) string {
	v, err := strconv.Unquote("\"" + value + "\"")
	if err != nil {
		return value
	}
	return v
}

var (
	spaceReplacer     = strings.NewReplacer(" ", "_")
	backslashReplacer = strings.NewReplacer(`\`, "_")
)

func hasPrefix(mf map[string]*prometheus.MetricFamily, prefix string) bool {
	for name := range mf {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
