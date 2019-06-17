package client

import (
	"encoding/json"
)

type version struct {
	// API version
	Version int
}

// SummaryRaw represents summary statistics int raw format (no number formatting applied).
type SummaryRaw struct {
	DomainsBeingBlocked int64   `json:"domains_being_blocked"`
	DNSQueriesToday     int64   `json:"dns_queries_today"`
	AdsBlockedToday     int64   `json:"ads_blocked_today"`
	AdsPercentageToday  float64 `json:"ads_percentage_today"`
	UniqueDomains       int64   `json:"unique_domains"`
	QueriesForwarded    int64   `json:"queries_forwarded"`
	QueriesCached       int64   `json:"queries_cached"`
	ClientsEverSeen     int64   `json:"clients_ever_seen"`
	UniqueClients       int64   `json:"unique_clients"`
	DNSQueriesAllTypes  int64   `json:"dns_queries_all_types"`
	ReplyNODATA         int64   `json:"reply_NODATA"`
	ReplyNXDOMAIN       int64   `json:"reply_NXDOMAIN"`
	ReplyCNAME          int64   `json:"reply_CNAME"`
	ReplyIP             int64   `json:"reply_IP"`
	PrivacyLevel        int64   `json:"privacy_level"`
	Status              string  `json:"status"`
	GravityLastUpdated  struct {
		FileExists bool `json:"file_exists"`
		Absolute   int64
	} `json:"gravity_last_updated"`
}

type (
	queryTypes struct {
		Types QueryTypes `json:"querytypes"`
	}

	// QueryTypes represents DNS queries processing statistics.
	QueryTypes struct {
		A    float64 `json:"A (IPv4)"`
		AAAA float64 `json:"AAAA (IPv6)"`
		ANY  float64
		SRV  float64
		SOA  float64
		PTR  float64
		TXT  float64
	}
)

type (
	forwardDestinations struct {
		Destinations ForwardDestinations `json:"forward_destinations"`
	}

	// ForwardDestinations represents queries forwarding statistics.
	ForwardDestinations map[string]float64
)

type (
	topClients struct {
		Clients TopClients `json:"top_sources"`
	}
	// TopClients represents queries per client (source) statistics.
	TopClients map[string]int64
)

type (
	// Item Item.
	Item map[string]int

	// TopItems represents top domains and top advertisements statistics.
	TopItems struct {
		TopQueries Item `json:"top_queries"`
		TopAds     Item `json:"top_ads"`
	}
)

func (i *Item) UnmarshalJSON(data []byte) error {
	if isEmptyArray(data) {
		return nil
	}
	type tmp *Item
	return json.Unmarshal(data, (tmp)(i))
}
