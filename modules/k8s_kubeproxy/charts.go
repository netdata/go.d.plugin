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
		ID:    "kubeproxy_sync_proxy_rules_latency_microseconds",
		Title: "SyncProxyRules Latency",
		Units: "observes",
		Fam:   "sync proxy rules",
		Ctx:   "k8s_kubeproxy.kubeproxy_sync_proxy_rules_latency_microseconds",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "sync_proxy_rules_latency_microseconds_bucket_1000", Name: "le1000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_2000", Name: "le2000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_4000", Name: "le4000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_8000", Name: "le8000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_16000", Name: "le16000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_32000", Name: "le32000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_64000", Name: "le64000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_128000", Name: "le128000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_256000", Name: "le256000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_512000", Name: "le512000"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_1_024e+06", Name: "le1.024e+06"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_2_048e+06", Name: "le2.048e+06"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_4_096e+06", Name: "le4.096e+06"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_8_192e+06", Name: "le8.192e+06"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_1_6384e+07", Name: "le1.6384e+07"},
			{ID: "sync_proxy_rules_latency_microseconds_bucket_+Inf", Name: "+Inf"},
		},
	},
}
