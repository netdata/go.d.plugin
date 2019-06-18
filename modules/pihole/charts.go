package pihole

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "random",
		Title: "A Random Number", Units: "random", Fam: "random",
		Dims: Dims{
			{ID: "random0", Name: "random 0"},
			{ID: "random1", Name: "random 1"},
		},
	},
}

var authCharts = Charts{
	{
		ID:    "processed_dns_queries_types",
		Title: "Processed DNS Queries By Types",
		Units: "percentage",
		Fam:   "queries",
		Ctx:   "pihole.processed_dns_queries_types",
		Dims: Dims{
			{ID: "A", Div: 100},
			{ID: "AAAA", Div: 100},
			{ID: "ANY", Div: 100},
			{ID: "PTR", Div: 100},
			{ID: "SOA", Div: 100},
			{ID: "SRV", Div: 100},
			{ID: "TXT", Div: 100},
		},
	},
	{
		ID:    "forwarded_dns_queries_targets",
		Title: "Forwarded DNS Queries By Target",
		Units: "percentage",
		Fam:   "queries",
		Ctx:   "pihole.forwarded_dns_queries_targets",
		Dims: Dims{
			{ID: "target_blocklist", Name: "blocklist", Div: 100},
			{ID: "target_cache", Name: "cache", Div: 100},
		},
	},
	{
		ID:    "top_clients",
		Title: "Top Clients",
		Units: "queries",
		Fam:   "top clients",
		Ctx:   "pihole.top_clients_queries",
	},
}

func (p *Pihole) updateCharts(rmx *rawMetrics) {

}
