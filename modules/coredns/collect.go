package coredns

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (cd *CoreDNS) collect() (map[string]int64, error) {
	raw, err := cd.prom.Scrape()

	if err != nil {
		return nil, err
	}

	mx := metrics{}

	cd.collectPanic(raw, &mx)
	cd.collectRequest(raw, &mx)

	return stm.ToMap(mx), nil
}

func (cd CoreDNS) collectPanic(raw prometheus.Metrics, mx *metrics) {
	mx.PanicCountTotal.Set(raw.FindByName("coredns_panic_count_total").Max())
}

func (cd CoreDNS) collectRequest(raw prometheus.Metrics, mx *metrics) {
	for _, metric := range raw.FindByName("coredns_dns_request_count_total") {
		mx.RequestsCountTotal.Add(metric.Value)
	}
}
