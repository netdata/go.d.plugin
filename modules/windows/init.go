// SPDX-License-Identifier: GPL-3.0-or-later

package windows

import (
	"errors"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
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

func (w *WMI) initPrometheusClient(client *http.Client) (prometheus.Prometheus, error) {
	return prometheus.New(client, w.Request), nil
}
