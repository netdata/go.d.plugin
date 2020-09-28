package couchdb

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	Charts = module.Charts
	Dims   = module.Dims
	Vars   = module.Vars
)

var (
	overviewCharts = Charts{
		{
			ID:    "activity",
			Title: "Overall Activity",
			Units: "requests/s",
			Fam:   "dbactivity",
			Ctx:   "couchdb.activity",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "couchdb_database_reads", Name: "DB reads", Algo: module.Incremental},
				{ID: "couchdb_database_writes", Name: "DB writes", Algo: module.Incremental},
				{ID: "couchdb_httpd_view_reads", Name: "View reads", Algo: module.Incremental},
			},
		},
	}
)

var (
	systemCharts = Charts{}
)

var (
	dbCharts = Charts{}
)
