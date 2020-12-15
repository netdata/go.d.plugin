package couchbase

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	Charts = module.Charts
	Chart  = module.Chart
	Dim    = module.Dim
)

var bucketQuotaPercentUsedChart = Chart{
	ID:    "bucket_quota_percent_used",
	Title: "Quota Percent Used Per Bucket",
	Units: "%",
	Fam:   "buckets basic stats",
	Ctx:   "couchbase.bucket_quota_percent_used",
}

var bucketOpsPerSecChart = Chart{
	ID:    "bucket_ops_per_sec",
	Title: "Operations Per Second Per Bucket",
	Units: "ops/s",
	Fam:   "buckets basic stats",
	Ctx:   "couchbase.bucket_ops_per_sec",
}

var bucketDiskFetchesChart = Chart{
	ID:    "bucket_disk_fetches",
	Title: "Disk Fetches Per Bucket",
	Units: "fetches",
	Fam:   "buckets basic stats",
	Ctx:   "couchbase.bucket_disk_fetches",
}

var bucketItemCountChart = Chart{
	ID:    "bucket_item_count",
	Title: "Item Count Per Bucket",
	Units: "items",
	Fam:   "buckets basic stats",
	Ctx:   "couchbase.bucket_item_count",
	Type:  module.Stacked,
}

var bucketDiskUsedChart = Chart{
	ID:    "bucket_disk_used_stats",
	Title: "Disk Used Per Bucket",
	Units: "bytes",
	Fam:   "buckets basic stats",
	Ctx:   "couchbase.bucket_disk_used_stats",
}

var bucketDataUsedChart = Chart{
	ID:    "bucket_data_used",
	Title: "Data Used Per Bucket",
	Units: "bytes",
	Fam:   "buckets basic stats",
	Ctx:   "couchbase.bucket_data_used",
}

var bucketMemUsedChart = Chart{
	ID:    "bucket_mem_used",
	Title: "Memory Used Per Bucket",
	Units: "bytes",
	Fam:   "buckets basic stats",
	Ctx:   "couchbase.bucket_mem_used",
}

var bucketVBActiveNumNonResidentChart = Chart{
	ID:    "bucket_vb_active_num_non_resident_stats",
	Title: "Number Of Non-Resident Items Per Bucket",
	Units: "items",
	Fam:   "buckets basic stats",
	Ctx:   "couchbase.bucket_vb_active_num_non_resident",
}
