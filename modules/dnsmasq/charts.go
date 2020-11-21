package dnsmasq

import "github.com/netdata/go.d.plugin/agent/module"

var cacheCharts = module.Charts{
	{
		ID:    "servers_queries",
		Title: "Queries forwarded to the upstream servers",
		Units: "queries/s",
		Fam:   "servers",
		Ctx:   "dnsmasq.servers_queries",
		Dims: module.Dims{
			{ID: "queries", Name: "success", Algo: module.Incremental},
			{ID: "failed_queries", Name: "failed", Algo: module.Incremental},
		},
	},
	{
		ID:    "cache_entries",
		Title: "Cache entries",
		Units: "entries",
		Fam:   "cache",
		Ctx:   "dnsmasq.cache_entries",
		Type:  module.Area,
		Dims: module.Dims{
			{ID: "cachesize", Name: "max"},
			{ID: "cache_entries", Name: "current"},
		},
	},
	{
		ID:    "cache_operations",
		Title: "Cache operations",
		Units: "operations/s",
		Fam:   "cache",
		Ctx:   "dnsmasq.cache_operations",
		Dims: module.Dims{
			{ID: "insertions", Algo: module.Incremental},
			{ID: "evictions", Algo: module.Incremental},
		},
	},
	{
		ID:    "cache_performance",
		Title: "Cache performance",
		Units: "events/s",
		Fam:   "cache",
		Ctx:   "dnsmasq.cache_performance",
		Dims: module.Dims{
			{ID: "hits", Algo: module.Incremental},
			{ID: "misses", Algo: module.Incremental},
		},
	},
}
