package prometheus

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"

	"github.com/prometheus/prometheus/pkg/textparse"
)

const (
	precision = 1000
)

// TODO: proper cleanup stale charts
func (p *Prometheus) collect() (map[string]int64, error) {
	pms, err := p.prom.Scrape()
	if err != nil {
		return nil, err
	}

	switch {
	case len(pms) == 0:
		p.Warningf("endpoint '%s' returned 0 time series", p.URL)
		return nil, nil
	case len(pms) > p.MaxTS:
		p.Warningf("endpoint '%s' returned %d time series, limit is %d", p.URL, len(pms), p.MaxTS)
		return nil, nil
	}

	if p.ExpectedPrefix != "" {
		if !hasMetricWithPrefix(pms, p.ExpectedPrefix) {
			return nil, fmt.Errorf("endpoint '%s' returned metrics without expected prefix (%s)",
				p.URL, p.ExpectedPrefix)
		}
		p.ExpectedPrefix = ""
	}

	names, metricSet := p.buildMetricSet(pms)
	meta := p.prom.Metadata()
	mx := make(map[string]int64)

	for _, name := range names {
		metrics := metricSet[name]
		if len(metrics) == 0 || p.skipMetrics[name] {
			continue
		}

		switch meta.Type(name) {
		case textparse.MetricTypeGauge, textparse.MetricTypeCounter:
			p.collectAny(mx, metrics, meta)
		case textparse.MetricTypeSummary:
			p.collectSummary(mx, metrics, meta)
		case textparse.MetricTypeHistogram:
			p.collectHistogram(mx, metrics, meta)
		case textparse.MetricTypeUnknown:
			p.collectUnknown(mx, metrics, meta)
		}
	}
	p.Debugf("time series: %d, metrics: %d, charts: %d", len(pms), len(names), len(*p.Charts()))
	mx["series"] = int64(len(pms))
	mx["metrics"] = int64(len(names))
	mx["charts"] = int64(len(*p.Charts()))
	return mx, nil
}

// TODO: should be done by prom pkg
func (p *Prometheus) buildMetricSet(pms prometheus.Metrics) (names []string, metrics map[string]prometheus.Metrics) {
	names = make([]string, 0, len(pms))
	metrics = make(map[string]prometheus.Metrics)

	for _, pm := range pms {
		if _, ok := metrics[pm.Name()]; !ok {
			names = append(names, pm.Name())
		}
		metrics[pm.Name()] = append(metrics[pm.Name()], pm)
	}

	var i int
	for _, name := range names {
		if len(metrics[name]) > p.MaxTSPerMetric {
			delete(metrics, name)
		} else {
			names[i] = name
			i++
		}
	}
	return names[:i], metrics
}

func hasMetricWithPrefix(pms prometheus.Metrics, prefix string) bool {
	for _, pm := range pms {
		if strings.HasPrefix(pm.Name(), prefix) {
			return true
		}
	}
	return false
}
