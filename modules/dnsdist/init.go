package dnsdist

import (
	"errors"
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (d DNSdist) validateConfig() error {
	if d.Config.Url == "" {
		return errors.New("'url' parameter is not set.")
	}

	if d.Config.User == "" {
		return errors.New("'user' parameter is not set.")
	}

	if d.Config.Pass == "" {
		return errors.New("'pass' parameter is not set.")
	}

	for i,cfg := range d.Config.Headers {
		if cfg.Name == "" {
			return fmt.Errorf("'headers[%d]->name' parameter not set", i+1)
		}

		if cfg.Value == "" {
			return fmt.Errorf("'headers[%d]->value' parameter not set", i+1)
		}
	}

	return nil
}

/*
func (d DNSdist) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(d.Client)
}
*/

func (d DNSdist) initCharts() (*module.Charts, error) {
	return charts.Copy(), nil
}