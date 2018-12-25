package dnsquery

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dim is an alias for modules.Dim
	Dim = modules.Dim
)

var charts = Charts{
	{
		ID:    "query_time",
		Title: "DNS Response Time",
		Units: "ms",
		Fam:   "query time",
		Ctx:   "dns_query_time.query_time",
	},
}
