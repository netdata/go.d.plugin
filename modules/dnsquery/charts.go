package dnsquery

import "github.com/netdata/go.d.plugin/agent/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "query_time",
		Title: "DNS Query Time",
		Units: "ms",
		Fam:   "query time",
		Ctx:   "dns_query_time.query_time",
	},
}
