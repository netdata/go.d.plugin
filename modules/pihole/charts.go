package pihole

import (
	"github.com/netdata/go-orchestrator"
	"github.com/netdata/go-orchestrator/module"
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
)

// ads_percentage_today
var (
	charts = Charts{
		// queries
		{
			ID:    "queries_total",
			Title: "Total Queries (Cached, Blocked and Forwarded)",
			Units: "queries/s",
			Fam:   "queries",
			Ctx:   "pihole.queries",
			Dims: Dims{
				{ID: "dns_queries_today", Name: "queries", Algo: module.Incremental},
			},
		},
		{
			ID:    "processed_queries_total",
			Title: "Processed Queries Total",
			Units: "queries",
			Fam:   "queries",
			Ctx:   "pihole.processed_queries",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "queries_cached", Name: "cached"},
				{ID: "ads_blocked_today", Name: "blocked"},
				{ID: "queries_forwarded", Name: "forwarded"},
			},
		},
		{
			ID:    "processed_queries",
			Title: "Processed Queries",
			Units: "queries/s",
			Fam:   "queries",
			Ctx:   "pihole.processed_queries",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "queries_cached", Name: "cached", Algo: module.Incremental},
				{ID: "ads_blocked_today", Name: "blocked", Algo: module.Incremental},
				{ID: "queries_forwarded", Name: "forwarded", Algo: module.Incremental},
			},
		},
		{
			ID:    "processed_queries_ratio",
			Title: "Processed Queries Ratio",
			Units: "percentage",
			Fam:   "queries",
			Ctx:   "pihole.processed_queries_percentage",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "queries_cached", Name: "cached", Algo: module.PercentOfAbsolute},
				{ID: "ads_blocked_today", Name: "blocked", Algo: module.PercentOfAbsolute},
				{ID: "queries_forwarded", Name: "forwarded", Algo: module.PercentOfAbsolute},
			},
		},
		{
			ID:    "blocked_queries_percentage",
			Title: "Blocked Queries Percentage",
			Units: "percentage",
			Fam:   "queries",
			Ctx:   "pihole.blocked_queries_percentage",
			Dims: Dims{
				{ID: "ads_percentage_today", Name: "blocked", Div: 100},
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
				{ID: "domains_being_blocked", Name: "on blocklist"},
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
		},
	}

	// authentication required
	authCharts = Charts{
		{
			ID:    "processed_dns_queries_types",
			Title: "Processed DNS Queries By Types",
			Units: "percentage",
			Fam:   "query types",
			Ctx:   "pihole.processed_dns_queries_types",
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
			ID:    "forwarded_dns_queries_destination",
			Title: "Forwarded DNS Queries By Destination",
			Units: "percentage",
			Fam:   "queries answered by",
			Ctx:   "pihole.forwarded_dns_queries_destinations",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "destination_cache", Name: "cache", Div: 100},
				{ID: "destination_blocklist", Name: "blocklist", Div: 100},
			},
		},
	}

	topClientsChart = Chart{
		ID:       "top_clients",
		Title:    "Top Clients Total",
		Units:    "queries",
		Fam:      "top clients",
		Ctx:      "pihole.top_clients_queries",
		Type:     module.Stacked,
		Priority: orchestrator.DefaultJobPriority + 10,
	}

	topPermittedDomainsChart = Chart{
		ID:       "top_permitted_domains",
		Title:    "Top Permitted Domains",
		Units:    "hits",
		Fam:      "top domains",
		Ctx:      "pihole.top_permitted_domains",
		Type:     module.Stacked,
		Priority: orchestrator.DefaultJobPriority + 20,
	}

	topBlockedDomainsChart = Chart{
		ID:       "top_blocked_domains",
		Title:    "Top Blocked Domains",
		Units:    "hits",
		Fam:      "top domains",
		Ctx:      "pihole.top_blocked_domains",
		Type:     module.Stacked,
		Priority: orchestrator.DefaultJobPriority + 30,
	}
)

func (p *Pihole) updateCharts(pmx *piholeMetrics) {
	p.updateForwardDestinationsCharts(pmx)
	p.updateTopClientChart(pmx)
	p.updateTopPermittedDomainsChart(pmx)
	p.updateTopBlockedDomainsChart(pmx)
}

func (p *Pihole) updateForwardDestinationsCharts(pmx *piholeMetrics) {
	if !pmx.hasForwardDestinations() {
		return
	}

	chart := p.charts.Get("forwarded_dns_queries_destination")
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
