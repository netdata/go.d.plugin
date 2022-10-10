// SPDX-License-Identifier: GPL-3.0-or-later

package cassandra

import "github.com/netdata/go.d.plugin/agent/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
	// Vars is an alias for module.Vars
	Vars = module.Vars
	// Opts is an alias for module.Dim
	Opts = module.DimOpts
)

var chartCassandraThroughput = Chart{
	ID:    "throughput_%s_%s",
	Title: "I/O requests.",
	Units: "requests/s",
	Fam:   "throughput %s %s",
	Ctx:   "cassandra.throughput",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_read", Name: "read", Algo: module.Incremental},
		{ID: "%s_%s_write", Name: "write", Algo: module.Incremental},
	},
}

var chartCassandraLatency = Chart{
	ID:    "latency_%s_%s",
	Title: "I/O latency.",
	Units: "requests/s",
	Fam:   "latency %s %s",
	Ctx:   "cassandra.latency",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_read", Name: "read", Algo: module.Incremental},
		{ID: "%s_%s_write", Name: "write", Algo: module.Incremental},
	},
}

var chartCassandraCache = Chart{
	ID:    "cache_%s_%s",
	Title: "Cache Hit",
	Units: "percentage/s",
	Fam:   "cache %s %s",
	Ctx:   "cassandra.cache",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_ratio", Name: "ratio", Algo: module.Incremental},
	},
}

var chartCassandraDiskNode = Chart{
	ID:    "node_%s_%s",
	Title: "Disk Node",
	Units: "bytes/s",
	Fam:   "node %s %s",
	Ctx:   "cassandra.node",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_node", Name: "node", Algo: module.Incremental},
	},
}

var chartCassandraDiskColumn = Chart{
	ID:    "column_%s_%s",
	Title: "Disk Column",
	Units: "bytes/s",
	Fam:   "column %s %s",
	Ctx:   "cassandra.disk_column",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_column", Name: "column", Algo: module.Incremental},
	},
}

var chartCassandraDiskCompactionCompleted = Chart{
	ID:    "compaction_completed_%s_%s",
	Title: "Completed Compaction Tasks",
	Units: "events/s",
	Fam:   "compaction %s %s",
	Ctx:   "cassandra.compaction_completed",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_compaction", Name: "compaction", Algo: module.Incremental},
	},
}

var chartCassandraDiskCompactionQueue = Chart{
	ID:    "compaction_queue_%s_%s",
	Title: "Queued Compaction Tasks",
	Units: "events/s",
	Fam:   "queue %s %s",
	Ctx:   "cassandra.compaction_queued",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_queue", Name: "queue", Algo: module.Incremental},
	},
}

var chartCassandraParNewCount = Chart{
	ID:    "gc_parnew_count_%s_%s",
	Title: "Young-generation garbage collection counter",
	Units: "garbage collection/s",
	Fam:   "par new count %s %s",
	Ctx:   "cassandra.gc_parnew_count",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_parnew", Name: "parnew", Algo: module.Incremental},
	},
}

var chartCassandraParNewTime = Chart{
	ID:    "gc_parnew_time_%s_%s",
	Title: "Young-generation garbage collection timer",
	Units: "period of time",
	Fam:   "par new time %s %s",
	Ctx:   "cassandra.gc_parnew_time",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_parnew", Name: "parnew", Algo: module.Incremental},
	},
}

var chartCassandraMarkSweepCount = Chart{
	ID:    "gc_marksweep_count_%s_%s",
	Title: "Old-generation collection",
	Units: "events/s",
	Fam:   "mark sweep %s %s",
	Ctx:   "cassandra.gc_sweep_count",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_sweep", Name: "sweep", Algo: module.Incremental},
	},
}

var chartCassandraMarkSweepTime = Chart{
	ID:    "gc_marksweep_time_%s_%s",
	Title: "Elapsed time Old-generation collection",
	Units: "period of time",
	Fam:   "mark sweep %s %s",
	Ctx:   "cassandra.gc_sweep_time",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_sweep", Name: "sweep", Algo: module.Incremental},
	},
}

var chartCassandraErrorTimeout = Chart{
	ID:    "error_timeout_%s_%s",
	Title: "Requests not unacknowledged",
	Units: "requests/s",
	Fam:   "request timeout %s %s",
	Ctx:   "cassandra.error_timeout",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_read", Name: "read", Algo: module.Incremental},
		{ID: "%s_%s_write", Name: "write", Algo: module.Incremental},
	},
}

var chartCassandraErrorUnavailable = Chart{
	ID:    "error_timeout_%s_%s",
	Title: "Request was unavailable",
	Units: "requests/s",
	Fam:   "request unavailable %s %s",
	Ctx:   "cassandra.error_unavailable",
	Type:  module.Line,
	Dims: Dims{
		{ID: "%s_%s_read", Name: "read", Algo: module.Incremental},
		{ID: "%s_%s_write", Name: "write", Algo: module.Incremental},
	},
}

func newCassandraCharts() *Charts {
	return &Charts{
		chartCassandraThroughput.Copy(),
		chartCassandraLatency.Copy(),
		chartCassandraCache.Copy(),
		chartCassandraDiskNode.Copy(),
		chartCassandraDiskColumn.Copy(),
		chartCassandraDiskCompactionCompleted.Copy(),
		chartCassandraDiskCompactionQueue.Copy(),
		chartCassandraParNewCount.Copy(),
		chartCassandraParNewTime.Copy(),
		chartCassandraMarkSweepCount.Copy(),
		chartCassandraMarkSweepTime.Copy(),
		chartCassandraErrorTimeout.Copy(),
		chartCassandraErrorUnavailable.Copy(),
	}
}

func (c *Cassandra) updateCharts(mx *metrics) {
	c.updateThrouputCharts(mx)
	c.updateLatencyCharts(mx)
	c.updateCacheCharts(mx)
	c.updateDiskCharts(mx)
}

func (c *Cassandra) updateThrouputCharts(mx *metrics) {
	if !mx.hasThrouput() {
		return
	}
}

func (c *Cassandra) updateLatencyCharts(mx *metrics) {
	if !mx.hasLatency() {
		return
	}
}

func (c *Cassandra) updateCacheCharts(mx *metrics) {
	if !mx.hasCache() {
		return
	}
}

func (c *Cassandra) updateDiskCharts(mx *metrics) {
	if !mx.hasDisk() {
		return
	}
}
