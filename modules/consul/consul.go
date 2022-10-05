// SPDX-License-Identifier: GPL-3.0-or-later

package consul

import (
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("consul", creator)
}

func New() *Consul {
	return &Consul{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{URL: "http://127.0.0.1:8500"},
				Client:  web.Client{Timeout: web.Duration{Duration: time.Second * 2}},
			},
		},
		checks: make(map[string]bool),
		charts: globalCharts.Copy(),
	}
}

type Config struct {
	web.HTTP `yaml:",inline"`

	ACLToken       string `yaml:"acl_token"`
	ChecksSelector string `yaml:"checks_selector"`
}

type Consul struct {
	module.Base

	Config `yaml:",inline"`

	charts *module.Charts

	httpClient *http.Client

	checks   map[string]bool
	checksSr matcher.Matcher
}

func (c *Consul) Init() bool {
	if err := c.validateConfig(); err != nil {
		c.Errorf("config validation: %v", err)
		return false
	}

	httpClient, err := c.initHTTPClient()
	if err != nil {
		c.Errorf("init HTTP client: %v", err)
		return false
	}
	c.httpClient = httpClient

	sr, err := c.initChecksSelector()
	if err != nil {
		c.Errorf("init checks filter: %v", err)
		return false
	}
	c.checksSr = sr

	return true
}

func (c *Consul) Check() bool {
	return len(c.Collect()) > 0
}

func (c *Consul) Charts() *module.Charts {
	return c.charts
}

func (c *Consul) Collect() map[string]int64 {
	mx, err := c.collect()
	if err != nil {
		c.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (c *Consul) Cleanup() {
	if c.httpClient != nil {
		c.httpClient.CloseIdleConnections()
	}
}
