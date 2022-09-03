// SPDX-License-Identifier: GPL-3.0-or-later

package httpcheck

import (
	"errors"
	"github.com/netdata/go.d.plugin/pkg/web"
	"net/http"
	"regexp"
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
