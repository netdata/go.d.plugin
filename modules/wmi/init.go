// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"errors"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (w WMI) validateConfig() error {
	if w.URL == "" {
		return errors.New("'url' is not set")
	}
	return nil
}

func (w WMI) initPrometheusClient() (prometheus.Prometheus, error) {
	client, err := web.NewHTTPClient(w.Client)
	if err != nil {
		return nil, err
	}
	return prometheus.New(client, w.Request), nil
}
