package k8s_kubeproxy

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "kubeproxy_sync_proxy_rules",
		Title: "Sync Proxy Rules",
		Units: "events/s",
		Fam:   "sync proxy rules",
		Ctx:   "k8s_kubeproxy.kubeproxy_sync_proxy_rules",
		Dims: Dims{
			{ID: "sync_proxy_rules_count", Name: "sync proxy rules", Algo: module.Incremental},
		},
	},
	{
		ID:    "kubeproxy_sync_proxy_rules_latency_microseconds",
		Title: "Sync Proxy Rules Latency",
		Units: "observes",
		Fam:   "sync proxy rules",
		Ctx:   "k8s_kubeproxy.kubeproxy_sync_proxy_rules_latency_microseconds",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "sync_proxy_rules_bucket_1000", Name: "le1000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_2000", Name: "le2000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_4000", Name: "le4000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_8000", Name: "le8000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_16000", Name: "le16000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_32000", Name: "le32000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_64000", Name: "le64000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_128000", Name: "le128000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_256000", Name: "le256000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_512000", Name: "le512000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_1_024e+06", Name: "le1.024e+06", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_2_048e+06", Name: "le2.048e+06", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_4_096e+06", Name: "le4.096e+06", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_8_192e+06", Name: "le8.192e+06", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_1_6384e+07", Name: "le1.6384e+07", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_+Inf", Name: "+Inf", Algo: module.Incremental},
		},
	},
	{
		ID:    "rest_client_requests_by_code",
		Title: "HTTP Requests By Status Code",
		Units: "requests/s",
		Fam:   "rest client",
		Ctx:   "k8s_kubeproxy.rest_client_requests_by_code",
	},
	{
		ID:    "rest_client_requests_by_method",
		Title: "HTTP Requests By Status Method",
		Units: "requests/s",
		Fam:   "rest client",
		Ctx:   "k8s_kubeproxy.rest_client_requests_by_method",
	},
}
