// SPDX-License-Identifier: GPL-3.0-or-later

package httpcheck

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (hc *HTTPCheck) validateConfig() error {
	if hc.URL == "" {
		return errors.New("'url' not set")
	}
	return nil
}

func (hc *HTTPCheck) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(hc.Client)
}

func (hc *HTTPCheck) initResponseMatchRegexp() (*regexp.Regexp, error) {
	if hc.ResponseMatch == "" {
		return nil, nil
	}
	return regexp.Compile(hc.ResponseMatch)
}

func (hc *HTTPCheck) initCharts() *module.Charts {
	charts := httpCheckCharts.Copy()

	for _, chart := range *charts {
		chart.Labels = []module.Label{
			{Key: "url", Value: hc.URL},
		}
	}

	return charts
}
