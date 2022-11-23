// SPDX-License-Identifier: GPL-3.0-or-later

package httpcheck

import (
	"net/http"
	"regexp"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	module.Register("httpcheck", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *HTTPCheck {
	return &HTTPCheck{
		Config: Config{
			HTTP: web.HTTP{
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second},
				},
			},
			AcceptedStatuses: []int{200},
		},
		acceptedStatuses: make(map[int]bool),
	}
}

type Config struct {
	web.HTTP         `yaml:",inline"`
	AcceptedStatuses []int  `yaml:"status_accepted"`
	ResponseMatch    string `yaml:"response_match"`
}

type (
	HTTPCheck struct {
		module.Base
		Config      `yaml:",inline"`
		UpdateEvery int `yaml:"update_every"`

		charts *module.Charts

		acceptedStatuses map[int]bool
		reResponse       *regexp.Regexp
		client           client
		metrics          metrics
	}
	client interface {
		Do(*http.Request) (*http.Response, error)
	}
)

func (hc *HTTPCheck) Init() bool {
	if err := hc.validateConfig(); err != nil {
		hc.Errorf("config validation: %v", err)
		return false
	}

	hc.charts = hc.initCharts()

	httpClient, err := hc.initHTTPClient()
	if err != nil {
		hc.Errorf("init HTTP client: %v", err)
		return false
	}
	hc.client = httpClient

	re, err := hc.initResponseMatchRegexp()
	if err != nil {
		hc.Errorf("init response match regexp: %v", err)
		return false
	}
	hc.reResponse = re

	for _, v := range hc.AcceptedStatuses {
		hc.acceptedStatuses[v] = true
	}

	hc.Debugf("using URL %s", hc.URL)
	hc.Debugf("using HTTP timeout %s", hc.Timeout.Duration)
	hc.Debugf("using accepted HTTP statuses %v", hc.AcceptedStatuses)
	if hc.reResponse != nil {
		hc.Debugf("using response match regexp %s", hc.reResponse)
	}

	return true
}

func (hc *HTTPCheck) Check() bool {
	return len(hc.Collect()) > 0
}

func (hc *HTTPCheck) Charts() *module.Charts {
	return hc.charts
}

func (hc *HTTPCheck) Collect() map[string]int64 {
	mx, err := hc.collect()
	if err != nil {
		hc.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (hc *HTTPCheck) Cleanup() {}
