package elasticsearch

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func (es Elasticsearch) checkConfig() error {
	if es.UserURL == "" {
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

func (es Elasticsearch) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(es.Client)
}

func (es Elasticsearch) initCharts() (*Charts, error) {
	charts := module.Charts{}
	if es.DoNodeStats {
		if err := charts.Add(*nodeCharts.Copy()...); err != nil {
			return nil, fmt.Errorf("add local node charts: %v", err)
		}
	}
	if es.DoIndicesStats {
		if err := charts.Add(*nodeIndicesStatsCharts.Copy()...); err != nil {
			return nil, fmt.Errorf("add local indices charts: %v", err)
		}
	}
	if es.DoClusterHealth {
		if err := charts.Add(*clusterHealthCharts.Copy()...); err != nil {
			return nil, fmt.Errorf("add cluster health charts: %v", err)
		}
	}
	if es.DoClusterHealth {
		if err := charts.Add(*clusterStatsCharts.Copy()...); err != nil {
			return nil, fmt.Errorf("add cluster stats charts: %v", err)
		}
	}
	if len(charts) == 0 {
		return nil, errors.New("zero charts")
	}
	return &charts, nil
}
