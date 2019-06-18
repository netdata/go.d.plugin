package pihole

import "github.com/netdata/go.d.plugin/modules/pihole/client"

type piholeMetrics struct {
	summary    *client.SummaryRaw
	queryTypes *client.QueryTypes
	forwarders *[]client.ForwardDestination
	topClients *[]client.TopClient
	topQueries *[]client.TopQuery
	topAds     *[]client.TopAdvertisement
}

func (p piholeMetrics) hasSummary() bool {
	return p.summary != nil
}

func (p piholeMetrics) hasQueryTypes() bool {
	return p.queryTypes != nil
}

func (p piholeMetrics) hasForwardDestinations() bool {
	return p.forwarders != nil && len(*p.forwarders) > 0
}

func (p piholeMetrics) hasTopClients() bool {
	return p.topClients != nil && len(*p.topClients) > 0
}

func (p piholeMetrics) hasTopQueries() bool {
	return p.topQueries != nil && len(*p.topQueries) > 0
}

func (p piholeMetrics) hasTopAdvertisers() bool {
	return p.topAds != nil && len(*p.topAds) > 0
}
