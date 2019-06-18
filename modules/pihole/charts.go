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
	set := make(map[string]bool)

	for _, v := range *pmx.forwarders {
		id := "target_" + v.Name
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

func (p *Pihole) updateTopDomainsChart(pmx *piholeMetrics) {
	if !pmx.hasTopQueries() {
		return
	}

	chart := p.charts.Get(topDomainsChart.ID)
	if chart == nil {
		chart = topDomainsChart.Copy()
		panicIf(p.charts.Add(chart))
	}

	set := make(map[string]bool)

	for _, v := range *pmx.topQueries {
		id := "top_domain_" + v.Name
		set[id] = true

		if chart.HasDim(id) {
			continue
		}

		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, set)
}

func (p *Pihole) updateTopAdvertisersChart(pmx *piholeMetrics) {
	if !pmx.hasTopAdvertisers() {
		return
	}

	chart := p.charts.Get(topAdvertisersChart.ID)
	if chart == nil {
		chart = topAdvertisersChart.Copy()
		panicIf(p.charts.Add(chart))
	}

	set := make(map[string]bool)

	for _, v := range *pmx.topAds {
		id := "top_ad_" + v.Name
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
