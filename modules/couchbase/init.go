package couchbase

import (
	"errors"
	"net/http"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (cb *Couchbase) initCharts() (*Charts, error) {
	var bucketCharts = module.Charts{
		dbPercentCharts.Copy(),
		opsPerSecCharts.Copy(),
		diskFetchesCharts.Copy(),
		diskUsedCharts.Copy(),
		dataUsedCharts.Copy(),
		memUsedCharts.Copy(),
		vbActiveNumNonResidentCharts.Copy(),
	}
	return bucketCharts.Copy(), nil
}

func (cb Couchbase) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(cb.Client)
}

func (cb Couchbase) validateConfig() error {
	if cb.URL == "" {
		return errors.New("URL not set")
	}
	if _, err := web.NewHTTPRequest(cb.Request); err != nil {
		return err
	}
	return nil
}