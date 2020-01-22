package cockroachdb

import (
	"errors"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (c *CockroachDB) collect() (map[string]int64, error) {
	scraped, err := c.prom.Scrape()
	if err != nil {
		return nil, err
	}

	if !validCockroachDBMetrics(scraped) {
		return nil, errors.New("returned metrics aren't CockroachDB metrics")
	}

	mx := collectScraped(scraped)

	return stm.ToMap(mx), nil
}

func collectScraped(scraped prometheus.Metrics) metrics {
	return metrics{
		Storage: collectStorage(scraped),
	}
}

func validCockroachDBMetrics(scraped prometheus.Metrics) bool {
	// TODO: enough?
	return scraped.FindByName("sql_restart_savepoint_count_internal").Len() > 0
}
