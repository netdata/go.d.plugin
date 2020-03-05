package pulsar

import (
	"errors"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func isValidPulsarMetrics(pms prometheus.Metrics) bool {
	// TODO:
	return pms.FindByName("metricPUBLISHError").Len() > 0 && pms.FindByName("metricRouterSubscriptions").Len() > 0
}

func (p *Pulsar) collect() (map[string]int64, error) {
	pms, err := p.prom.Scrape()
	if err != nil {
		return nil, err
	}

	if !isValidPulsarMetrics(pms) {
		return nil, errors.New("returned metrics aren't VerneMQ metrics")
	}

	mx := p.collectPulsar(pms)

	return stm.ToMap(mx), nil
}

func (p *Pulsar) collectPulsar(pms prometheus.Metrics) map[string]float64 {
	return nil
}
