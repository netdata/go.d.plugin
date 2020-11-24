package couchbase

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	Charts = module.Charts
	Chart  = module.Chart
	Dims   = module.Dims
	Dim    = module.Dim
)

var dbPercentCharts = Chart{
	ID:    "couchbase_quota_percent_used_stats",
	Title: "Quota Percent Used Per Bucket",
	Units: "%",
	Fam:   "quota percent used",
	Ctx:   "couchbase.couchbase_quota_percent_used_stats",
}

var opPerSecCharts = Chart{
	ID:    "couchbase_ops_per_sec_stats",
	Title: "Operations Per Second",
	Units: "num",
	Fam:   "ops per sec",
	Ctx:   "couchbase.couchbase_ops_per_sec_stats",
}

var diskFetchesCharts = Chart{
	ID:    "couchbase_disk_fetches_stats",
	Title: "Disk Fetches",
	Units: "num/s",
	Fam:   "disk fetches",
	Ctx:   "couchbase.couchbase_disk_fetches_stats",
}

var itemCountCharts = Chart{
	ID:    "couchbase_item_count_stats",
	Title: "Item Count",
	Units: "items",
	Fam:   "item count",
	Ctx:   "couchbase.couchbase_item_count_stats",
	Type:  module.Stacked,
}

var diskUsedCharts = Chart{
	ID:    "couchbase_disk_used_stats",
	Title: "Disk Used Per Bucket",
	Units: "KiB/s",
	Fam:   "disk used",
	Ctx:   "couchbase.couchbase_disk_used_stats",
}

var dataUsedCharts = Chart{
	ID:    "couchbase_data_used_stats",
	Title: "Data Used Per Bucket",
	Units: "KiB/s",
	Fam:   "data used",
	Ctx:   "couchbase.couchbase_data_used_stats",
}

var memUsedCharts = Chart{
	ID:    "couchbase_mem_used_stats",
	Title: "Memory Used Per Bucket",
	Units: "KiB/s",
	Fam:   "memory used",
	Ctx:   "couchbase.couchbase_mem_used_stats",
}

var vbActiveNumNonResidentCharts = Chart{
	ID:    "couchbase_vb_active_num_non_resident_stats",
	Title: "Number Of Non-Resident Items",
	Units: "num/s",
	Fam:   "vb active num non resident",
	Ctx:   "couchbase.couchbase_vb_active_num_non_resident_stats",
}
