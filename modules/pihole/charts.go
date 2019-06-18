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

var (
	charts = Charts{
		{
			ID:    "random",
			Title: "A Random Number", Units: "random", Fam: "random",
			Dims: Dims{
				{ID: "random0", Name: "random 0"},
				{ID: "random1", Name: "random 1"},
			},
		},
	}

	// authentication required
	authCharts = Charts{
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
	}

	topClientsChart = Chart{
		ID:       "top_clients",
		Title:    "Top Clients",
		Units:    "queries",
		Fam:      "top clients",
		Ctx:      "pihole.top_clients_queries",
		Priority: orchestrator.DefaultJobPriority + 10,
	}

	topDomainsChart = Chart{
		ID:       "top_domains",
		Title:    "Top Domains",
		Units:    "queries",
		Fam:      "top domains",
		Ctx:      "pihole.top_domains_queries",
		Priority: orchestrator.DefaultJobPriority + 20,
	}

	topAdvertisersChart = Chart{
		ID:       "top_advertisers",
		Title:    "Top Advertisers",
		Units:    "queries",
		Fam:      "top ads",
		Ctx:      "pihole.top_ads_queries",
		Priority: orchestrator.DefaultJobPriority + 30,
	}
)

func (p *Pihole) updateCharts(pmx *piholeMetrics) {
	p.updateForwardDestinationsCharts(pmx)
	p.updateTopClientChart(pmx)
	p.updateTopDomainsChart(pmx)
	p.updateTopAdvertisersChart(pmx)
}

func (p *Pihole) updateForwardDestinationsCharts(pmx *piholeMetrics) {
	if !pmx.hasForwardDestinations() {
		return
	}

	chart := p.charts.Get("forwarded_dns_queries_targets")

	for _, v := range *pmx.forwarders {
		if v.Name == "blocklist" || v.Name == "cache" {
			continue
		}

		id := "target_" + v.Name
		if chart.HasDim(id) {
			continue
		}

		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name, Div: 100}))
		chart.MarkNotCreated()
	}
}

func (p *Pihole) updateTopClientChart(pmx *piholeMetrics) {
	if !pmx.hasTopClients() {
		return
	}

	if len(p.collected.topClients) == 0 && !p.charts.Has("top_clients") {
		panicIf(p.charts.Add(topClientsChart.Copy()))
	}

	chart := p.charts.Get("top_clients")
	set := make(map[string]bool)

	for _, v := range *pmx.topClients {
		id := "top_client_" + v.Name
		set[id] = true

		if p.collected.topClients[id] {
			continue
		}

		p.collected.topClients[id] = true
		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, p.collected.topClients, set)
}

func (p *Pihole) updateTopDomainsChart(pmx *piholeMetrics) {
	if !pmx.hasTopQueries() {
		return
	}

	if len(p.collected.topDomains) == 0 && !p.charts.Has("top_domains") {
		panicIf(p.charts.Add(topDomainsChart.Copy()))
	}

	chart := p.charts.Get("top_domains")
	set := make(map[string]bool)

	for _, v := range *pmx.topQueries {
		id := "top_domain_" + v.Name
		set[id] = true

		if p.collected.topDomains[id] {
			continue
		}

		p.collected.topDomains[id] = true
		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, p.collected.topDomains, set)
}

func (p *Pihole) updateTopAdvertisersChart(pmx *piholeMetrics) {
	if !pmx.hasTopAdvertisers() {
		return
	}

	if len(p.collected.topAds) == 0 && !p.charts.Has("top_advertisers") {
		panicIf(p.charts.Add(topAdvertisersChart.Copy()))
	}

	chart := p.charts.Get("top_advertisers")
	set := make(map[string]bool)

	for _, v := range *pmx.topAds {
		id := "top_ad_" + v.Name
		set[id] = true

		if p.collected.topAds[id] {
			continue
		}

		p.collected.topAds[id] = true
		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, p.collected.topAds, set)
}

func removeNotUpdatedDims(chart *Chart, existed, updated map[string]bool) {
	for id := range existed {
		if updated[id] {
			continue
		}

		delete(existed, id)
		panicIf(chart.RemoveDim(id))
		chart.MarkNotCreated()
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
