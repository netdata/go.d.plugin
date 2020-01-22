package cockroachdb

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (c *CockroachDB) collect() (map[string]int64, error) {
	scraped, err := c.prom.Scrape()
	if err != nil {
		return nil, err
	}

	mx := c.collectScraped(scraped)

	return stm.ToMap(mx), nil
}

func (c *CockroachDB) collectScraped(scraped prometheus.Metrics) *metrics {
	return nil
}
