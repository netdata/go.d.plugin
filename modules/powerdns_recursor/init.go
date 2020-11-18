package powerdns_recursor

import (
	"net/http"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (r *Recursor) validateConfig() error {
	return nil
}

func (r *Recursor) initHTTPClient() (*http.Client, error) {
	return nil, nil
}

func (r *Recursor) initCharts() (*module.Charts, error) {
	return nil, nil
}
