package pihole

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
	// Dim is an alias for module.Dim
	Vars = module.Vars
)

// ads_percentage_today
var (
	charts = Charts{
		// queries
		{
			ID:    "dns_queries_total",
			Title: "DNS Queries Total (Cached, Blocked and Forwarded)",
			Units: "queries",
			Fam:   "queries",
			Ctx:   "pihole.dns_queries_total",
			Dims: Dims{
				{ID: "dns_queries_today", Name: "queries"},
			},
		},
		{
			ID:    "dns_queries",
			Title: "DNS Queries",
			Units: "queries",
			Fam:   "queries",
			Ctx:   "pihole.dns_queries",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "queries_cached", Name: "cached"},
				{ID: "ads_blocked_today", Name: "blocked"},
				{ID: "queries_forwarded", Name: "forwarded"},
			},
		},
		{
			ID:    "dns_queries_percentage",
			Title: "DNS Queries Percentage",
			Units: "percentage",
			Fam:   "queries",
			Ctx:   "pihole.dns_queries_percentage",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "queries_cached", Name: "cached", Algo: module.PercentOfAbsolute},
				{ID: "ads_blocked_today", Name: "blocked", Algo: module.PercentOfAbsolute},
				{ID: "queries_forwarded", Name: "forwarded", Algo: module.PercentOfAbsolute},
			},
		},
		// clients
		{
			ID:    "unique_clients",
			Title: "Unique Clients",
			Units: "clients",
			Fam:   "clients",
			Ctx:   "pihole.unique_clients",
			Dims: Dims{
				{ID: "unique_clients", Name: "unique"},
			},
		},
		// blocklist
		{
			ID:    "domains_on_blocklist",
			Title: "Domains On Blocklist",
			Units: "domains",
			Fam:   "blocklist",
			Ctx:   "pihole.domains_on_blocklist",
			Dims: Dims{
				{ID: "domains_being_blocked", Name: "blocklist"},
			},
		},
		{
			ID:    "blocklist_last_update",
			Title: "Blocklist Last Update",
			Units: "seconds",
			Fam:   "blocklist",
			Ctx:   "pihole.blocklist_last_update",
			Dims: Dims{
				{ID: "blocklist_last_update", Name: "ago"},
			},
			Vars: Vars{
				{ID: "file_exists", Value: 1},
			},
		},
		// ads blocking
		{
			ID:    "unwanted_domains_blocking_status",
			Title: "Unwanted Domains Blocking Status",
			Units: "boolean",
			Fam:   "status",
			Ctx:   "pihole.unwanted_domains_blocking_status",
			Dims: Dims{
				{ID: "status", Name: "enabled"},
			},
		},
	}

	// authentication required
	authCharts = Charts{
		{
			ID:    "dns_queries_types",
			Title: "DNS Queries Per Type",
			Units: "percentage",
			Fam:   "query types",
			Ctx:   "pihole.dns_queries_types",
			Type:  module.Stacked,
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
			ID:    "dns_queries_forwarded_destination",
			Title: "DNS Queries Per Destination",
			Units: "percentage",
			Fam:   "queries answered by",
			Ctx:   "pihole.dns_queries_forwarded_destination",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "destination_cache", Name: "cache", Div: 100},
				{ID: "destination_blocklist", Name: "blocklist", Div: 100},
			},
		},
	}

	topClientsChart = Chart{
		ID:       "top_clients",
		Title:    "Top Clients",
		Units:    "requests",
		Fam:      "top clients",
		Ctx:      "pihole.top_clients",
		Type:     module.Stacked,
		Priority: module.Priority + 10,
	}

	topPermittedDomainsChart = Chart{
		ID:       "top_permitted_domains",
		Title:    "Top Permitted Domains",
		Units:    "hits",
		Fam:      "top domains",
		Ctx:      "pihole.top_permitted_domains",
		Type:     module.Stacked,
		Priority: module.Priority + 20,
	}

	topBlockedDomainsChart = Chart{
		ID:       "top_blocked_domains",
		Title:    "Top Blocked Domains",
		Units:    "hits",
		Fam:      "top domains",
		Ctx:      "pihole.top_blocked_domains",
		Type:     module.Stacked,
		Priority: module.Priority + 30,
	}
)

func (p *Pihole) updateCharts(pmx *piholeMetrics) {
	// auth charts
	p.updateForwardDestinationsCharts(pmx)
	p.updateTopClientChart(pmx)
	p.updateTopPermittedDomainsChart(pmx)
	p.updateTopBlockedDomainsChart(pmx)
}

func (p *Pihole) updateForwardDestinationsCharts(pmx *piholeMetrics) {
	if !pmx.hasForwardDestinations() {
		return
	}

	chart := p.charts.Get("dns_queries_forwarded_destination")
	set := make(map[string]bool)

	for _, v := range *pmx.forwarders {
		id := "destination_" + v.Name
		set[id] = true

		if chart.HasDim(id) {
			continue
		}

		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name, Div: 100}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, set)
}

func (p *Pihole) updateTopClientChart(pmx *piholeMetrics) {
	if !pmx.hasTopClients() {
		return
	}

	chart := p.charts.Get(topClientsChart.ID)
	if chart == nil {
		chart = topClientsChart.Copy()
		panicIf(p.charts.Add(chart))
	}

	set := make(map[string]bool)

	for _, v := range *pmx.topClients {
		id := "top_client_" + v.Name
		set[id] = true

		if chart.HasDim(id) {
			continue
		}

		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, set)
}

func (p *Pihole) updateTopPermittedDomainsChart(pmx *piholeMetrics) {
	if !pmx.hasTopQueries() {
		return
	}

	chart := p.charts.Get(topPermittedDomainsChart.ID)
	if chart == nil {
		chart = topPermittedDomainsChart.Copy()
		panicIf(p.charts.Add(chart))
	}

	set := make(map[string]bool)

	for _, v := range *pmx.topQueries {
		id := "top_perm_domain_" + v.Name
		set[id] = true

		if chart.HasDim(id) {
			continue
		}

		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, set)
}

func (p *Pihole) updateTopBlockedDomainsChart(pmx *piholeMetrics) {
	if !pmx.hasTopAdvertisers() {
		return
	}

	chart := p.charts.Get(topBlockedDomainsChart.ID)
	if chart == nil {
		chart = topBlockedDomainsChart.Copy()
		panicIf(p.charts.Add(chart))
	}

	set := make(map[string]bool)

	for _, v := range *pmx.topAds {
		id := "top_blocked_domain_" + v.Name
		set[id] = true

		if chart.HasDim(id) {
			continue
		}

		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, set)
}

func removeNotUpdatedDims(chart *Chart, updated map[string]bool) {
	if len(updated) == 0 {
		return
	}

	var notUpdated []string

	for _, d := range chart.Dims {
		if updated[d.ID] {
			continue
		}
		notUpdated = append(notUpdated, d.ID)
	}

	for _, v := range notUpdated {
		panicIf(chart.MarkDimRemove(v, true))
		chart.MarkNotCreated()
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
