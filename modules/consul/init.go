// SPDX-License-Identifier: GPL-3.0-or-later

package consul

import (
	"errors"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (c *Consul) validateConfig() error {
	if c.URL == "" {
		return errors.New("'url' not set")
	}
	return nil
}

func (c *Consul) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(c.Client)
}

func (c *Consul) initChecksSelector() (matcher.Matcher, error) {
	if c.ChecksSelector == "" {
		return matcher.TRUE(), nil
	}

	return matcher.NewSimplePatternsMatcher(c.ChecksSelector)
}
