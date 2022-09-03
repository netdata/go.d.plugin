// SPDX-License-Identifier: GPL-3.0-or-later

package pihole

import "github.com/netdata/go.d.plugin/modules/pihole/client"

type piholeMetrics struct {
	summary    *client.SummaryRaw           // ?summary
	queryTypes *client.QueryTypes           // ?getQueryTypes
	forwarders *[]client.ForwardDestination // ?getForwardedDestinations
	topClients *[]client.TopClient          // ?topClient
	topQueries *[]client.TopQuery           // ?topItem
	topAds     *[]client.TopAdvertisement   // ?topItem
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
