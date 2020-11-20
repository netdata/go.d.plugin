package dnsdist

import (
	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathLocalStatistics = "/jsonstat?command=stats"
)

func (d *DNSdist) collect(map[string]int64, error) {
	statistics, err := d.scrapeStatistics()
	if err != nil {
		return nil, err;
	}

	collected := make(map[string]int64)
	return collected, nil
}

func (d *DNSdist) scrapeStatistics() ([]statisticMetric, error) {
	req, err := web.NewHTTPRequest(d.Request)
	if err != nil {
		return nil, err
	}
	req.URL.Path = urlPathLocalStatistics

	var statistics statMetrics
	if err := d.doOKDecode(req, &statistics); err != nil {
		return nil, err
	}

	return statistics, nil
}