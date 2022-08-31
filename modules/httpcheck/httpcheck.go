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
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("httpcheck", creator)
}

var (
	defaultHTTPTimeout      = time.Second
	defaultAcceptedStatuses = []int{200}
)

func New() *HTTPCheck {
	config := Config{
		HTTP: web.HTTP{
			Client: web.Client{
				Timeout: web.Duration{Duration: defaultHTTPTimeout},
			},
		},
		AcceptedStatuses: defaultAcceptedStatuses,
	}
	return &HTTPCheck{
		Config:           config,
		acceptedStatuses: make(map[int]bool),
	}
}

type Config struct {
	web.HTTP         `yaml:",inline"`
	AcceptedStatuses []int  `yaml:"status_accepted"`
	ResponseMatch    string `yaml:"response_match"`
}

type client interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPCheck struct {
	module.Base
	Config      `yaml:",inline"`
	UpdateEvery int `yaml:"update_every"`

	acceptedStatuses map[int]bool
	reResponse       *regexp.Regexp
	client           client
	metrics          metrics
}

// Init makes initialization
func (hc *HTTPCheck) Init() bool {
	if hc.URL == "" {
		hc.Error("URL not set")
		return false
	}

	client, err := web.NewHTTPClient(hc.Client)
	if err != nil {
		hc.Errorf("error on creating HTTP client : %v", err)
		return false
	}
	hc.client = client

	if hc.ResponseMatch != "" {
		re, err := regexp.Compile(hc.ResponseMatch)
		if err != nil {
			hc.Errorf("error on creating regexp %s : %s", hc.ResponseMatch, err)
			return false
		}
		hc.reResponse = re
	}

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
	return charts.Copy()
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
