package prometheus

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

func (p *Prometheus) collectUnknown(mx map[string]int64, pms prometheus.Metrics, meta prometheus.Metadata) {
	pm := pms[0]
	switch {
	case pm.Labels.Has("quantile"):
		p.collectSummary(mx, pms, meta)
	case pm.Labels.Has("le"):
		p.collectHistogram(mx, pms, meta)
	default:
		p.collectAny(mx, pms, meta)
	}
}
