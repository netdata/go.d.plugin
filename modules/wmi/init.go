// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"errors"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (w *WMI) validateConfig() error {
	if w.URL == "" {
		return errors.New("'url' is not set")
	}
	return nil
}

func (w *WMI) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(w.Client)
}

func (w *WMI) initPrometheusCheckClient(client *http.Client) (prometheus.Prometheus, error) {
	sr, err := selector.Parse(metricCollectorSuccess)
	if err != nil {
		return nil, err
	}

	return prometheus.NewWithSelector(client, w.Request, sr), nil
}
