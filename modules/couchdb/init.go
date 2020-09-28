package couchdb

import (
	"errors"
	"net/http"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (cdb CouchDB) validateConfig() error {
	if cdb.URL == "" {
		return errors.New("URL not set")
	}
	if !(cdb.DoOverviewStats || cdb.DoSystemStats || cdb.DoDBStats) {
		return errors.New("all API calls are disabled")
	}
	if _, err := web.NewHTTPRequest(cdb.Request); err != nil {
		return err
	}
	return nil
}

func (cdb CouchDB) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(cdb.Client)
}

func (cdb CouchDB) initCharts() (*Charts, error) {
	charts := module.Charts{}
	if cdb.DoOverviewStats {
		if err := charts.Add(*overviewCharts.Copy()...); err != nil {
			return nil, err
		}
	}
	if cdb.DoSystemStats {
		if err := charts.Add(*systemCharts.Copy()...); err != nil {
			return nil, err
		}
	}
	if cdb.DoDBStats {
		if err := charts.Add(*dbCharts.Copy()...); err != nil {
			return nil, err
		}
	}
	if len(charts) == 0 {
		return nil, errors.New("zero charts")
	}
	return &charts, nil
}
