package couchbase

import (
	"errors"
	"net/http"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (cb *Couchbase) initCharts() (*Charts, error) {
	charts := module.Charts{}
	if err := charts.Add(dbPercentCharts.Copy()); err != nil {
		return nil, err
	}

	if err := charts.Add(opsPerSecCharts.Copy()); err != nil {
		return nil, err
	}

	if err := charts.Add(diskFetchesCharts.Copy()); err != nil {
		return nil, err
	}

	if err := charts.Add(itemCountCharts.Copy()); err != nil {
		return nil, err
	}

	if err := charts.Add(diskUsedCharts.Copy()); err != nil {
		return nil, err
	}

	if err := charts.Add(dataUsedCharts.Copy()); err != nil {
		return nil, err
	}

	if err := charts.Add(memUsedCharts.Copy()); err != nil {
		return nil, err
	}

	if err := charts.Add(vbActiveNumNonResidentCharts.Copy()); err != nil {
		return nil, err
	}
	return &charts, nil

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
