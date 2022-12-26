// SPDX-License-Identifier: GPL-3.0-or-later

package elasticsearch

import (
	"errors"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (es *Elasticsearch) validateConfig() error {
	if es.URL == "" {
		return errors.New("URL not set")
	}
	if !(es.DoNodeStats || es.DoClusterHealth || es.DoClusterStats || es.DoIndicesStats) {
		return errors.New("all API calls are disabled")
	}
	if _, err := web.NewHTTPRequest(es.Request); err != nil {
		return err
	}
	return nil
}

func (es *Elasticsearch) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(es.Client)
}

func (es *Elasticsearch) initCharts() (*module.Charts, error) {
	charts := module.Charts{}

	if es.DoNodeStats {
		if err := charts.Add(*nodeCharts.Copy()...); err != nil {
			return nil, err
		}
	}
	if es.DoClusterHealth {
		if err := charts.Add(*clusterHealthCharts.Copy()...); err != nil {
			return nil, err
		}
	}
	if es.DoClusterStats {
		if err := charts.Add(*clusterStatsCharts.Copy()...); err != nil {
			return nil, err
		}
	}

	if !es.DoIndicesStats && len(charts) == 0 {
		return nil, errors.New("zero charts")
	}

	return &charts, nil
}
