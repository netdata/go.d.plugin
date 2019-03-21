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
		Title: "Sync Proxy Rules Latency In Microseconds",
		Units: "observes per bucket",
		Fam:   "sync proxy rules",
		Ctx:   "k8s_kubeproxy.kubeproxy_sync_proxy_rules_latency_microseconds",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "sync_proxy_rules_bucket_1000", Name: "1000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_2000", Name: "2000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_4000", Name: "4000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_8000", Name: "8000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_16000", Name: "16000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_32000", Name: "32000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_64000", Name: "64000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_128000", Name: "128000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_256000", Name: "256000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_512000", Name: "512000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_1024000", Name: "1024000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_2048000", Name: "2048000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_4096000", Name: "4096000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_8192000", Name: "8192000", Algo: module.Incremental},
			{ID: "sync_proxy_rules_bucket_16384000", Name: "16384000", Algo: module.Incremental},
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
