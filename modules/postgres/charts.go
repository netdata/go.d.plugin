// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioConnectionsUtilization = module.Priority + iota
	prioConnectionsUsage
	prioConnectionsStateCount
	prioTransactionsDuration
	prioQueriesDuration
	prioCheckpointsRate
	prioCheckpointsTime
	prioBGWriterHaltsRate
	prioBuffersIORate
	prioBuffersBackendFsyncRate
	prioBuffersAllocRate
	prioWALIORate
	prioWALFilesCount
	prioWALArchivingFilesCount
	prioAutovacuumWorkersCount
	prioTXIDExhaustionTowardsAutovacuumPerc
	prioTXIDExhaustionPerc
	prioTXIDExhaustionOldestTXIDNum
	prioCatalogRelationsCount
	prioCatalogRelationsSize
	prioUptime
	prioReplicationAppWALLagSize
	prioReplicationAppWALLagTime
	prioReplicationSlotFilesCount
	prioDatabasesCount
	prioDBTransactionsRatio
	prioDBTransactionsRate
	prioDBConnectionsUtilization
	prioDBConnectionsCount
	prioDBCacheIORatio
	prioDBIORate
	prioDBOpsFetchedRowsRatio
	prioDBOpsReadRowsRate
	prioDBOpsWriteRowsRate
	prioDBConflictsRate
	prioDBConflictsReasonRate
	prioDBDeadlocksRate
	prioDBLocksHeldCount
	prioDBLocksAwaitedCount
	prioDBTempFilesCreatedRate
	prioDBTempFilesIORate
	prioDBSize
	prioTableRowsDeadRatio
	prioTableRowsCount
	prioTableOpsRowsRate
	prioTableOpsRowsHOTRatio
	prioTableOpsRowsHOTRate
	prioTableCacheIORatio
	prioTableIORate
	prioTableIndexCacheIORatio
	prioTableIndexIORate
	prioTableToastCacheIORatio
	prioTableToastIORate
	prioTableToastIndexCacheIORatio
	prioTableToastIndexIORate
	prioTableScansRate
	prioTableScansRowsRate
	prioTableAutovacuumSinceTime
	prioTableVacuumSinceTime
	prioTableAutoAnalyzeSinceTime
	prioTableLastAnalyzeAgo
	prioTableNullColumns
	prioTableTotalSize
	prioTableBloatSizePerc
	prioTableBloatSize
	prioIndexSize
	prioIndexBloatSizePerc
	prioIndexBloatSize
	prioIndexUsageStatus
)

var baseCharts = module.Charts{
	serverConnectionsUtilizationChart.Copy(),
	serverConnectionsUsageChart.Copy(),
	serverConnectionsStateCount.Copy(),
	checkpointsChart.Copy(),
	checkpointWriteChart.Copy(),
	buffersIORateChart.Copy(),
	buffersAllocRateChart.Copy(),
	bgWriterHaltsRateChart.Copy(),
	buffersBackendFsyncRateChart.Copy(),
	walIORateChart.Copy(),
	walFilesCountChart.Copy(),
	walArchivingFilesCountChart.Copy(),
	autovacuumWorkersCountChart.Copy(),
	txidExhaustionTowardsAutovacuumPercChart.Copy(),
	txidExhaustionPercChart.Copy(),
	txidExhaustionOldestTXIDNumChart.Copy(),

	catalogRelationSCountChart.Copy(),
	catalogRelationsSizeChart.Copy(),
	serverUptimeChart.Copy(),
	databasesCountChart.Copy(),
}

var (
	serverConnectionsUtilizationChart = module.Chart{
		ID:       "connections_utilization",
		Title:    "Connections utilization",
		Units:    "percentage",
		Fam:      "connections",
		Ctx:      "postgres.connections_utilization",
		Priority: prioConnectionsUtilization,
		Dims: module.Dims{
			{ID: "server_connections_utilization", Name: "used"},
		},
	}
	serverConnectionsUsageChart = module.Chart{
		ID:       "connections_usage",
		Title:    "Connections usage",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "postgres.connections_usage",
		Priority: prioConnectionsUsage,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "server_connections_available", Name: "available"},
			{ID: "server_connections_used", Name: "used"},
		},
	}
	serverConnectionsStateCount = module.Chart{
		ID:       "connections_state",
		Title:    "Connections in each state",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "postgres.connections_state_count",
		Priority: prioConnectionsStateCount,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "server_connections_state_active", Name: "active"},
			{ID: "server_connections_state_idle", Name: "idle"},
			{ID: "server_connections_state_idle_in_transaction", Name: "idle_in_transaction"},
			{ID: "server_connections_state_idle_in_transaction_aborted", Name: "idle_in_transaction_aborted"},
			{ID: "server_connections_state_fastpath_function_call", Name: "fastpath_function_call"},
			{ID: "server_connections_state_disabled", Name: "disabled"},
		},
	}

	checkpointsChart = module.Chart{
		ID:       "checkpoints_rate",
		Title:    "Checkpoints",
		Units:    "checkpoints/s",
		Fam:      "checkpointer",
		Ctx:      "postgres.checkpoints_rate",
		Priority: prioCheckpointsRate,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "checkpoints_timed", Name: "scheduled", Algo: module.Incremental},
			{ID: "checkpoints_req", Name: "requested", Algo: module.Incremental},
		},
	}
	// TODO: should be seconds, also it is units/s when using incremental...
	checkpointWriteChart = module.Chart{
		ID:       "checkpoints_time",
		Title:    "Checkpoint time",
		Units:    "milliseconds",
		Fam:      "checkpointer",
		Ctx:      "postgres.checkpoints_time",
		Priority: prioCheckpointsTime,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "checkpoint_write_time", Name: "write", Algo: module.Incremental},
			{ID: "checkpoint_sync_time", Name: "sync", Algo: module.Incremental},
		},
	}
	bgWriterHaltsRateChart = module.Chart{
		ID:       "bgwriter_halts_rate",
		Title:    "Background writer scan halts",
		Units:    "halts/s",
		Fam:      "bgwriter",
		Ctx:      "postgres.bgwriter_halts_rate",
		Priority: prioBGWriterHaltsRate,
		Dims: module.Dims{
			{ID: "maxwritten_clean", Name: "maxwritten", Algo: module.Incremental},
		},
	}

	buffersIORateChart = module.Chart{
		ID:       "buffers_io_rate",
		Title:    "Buffers written rate",
		Units:    "B/s",
		Fam:      "buffers",
		Ctx:      "postgres.buffers_io_rate",
		Priority: prioBuffersIORate,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "buffers_checkpoint", Name: "checkpoint", Algo: module.Incremental},
			{ID: "buffers_backend", Name: "backend", Algo: module.Incremental},
			{ID: "buffers_clean", Name: "bgwriter", Algo: module.Incremental},
		},
	}
	buffersBackendFsyncRateChart = module.Chart{
		ID:       "buffers_backend_fsync_rate",
		Title:    "Backend fsync calls",
		Units:    "calls/s",
		Fam:      "buffers",
		Ctx:      "postgres.buffers_backend_fsync_rate",
		Priority: prioBuffersBackendFsyncRate,
		Dims: module.Dims{
			{ID: "buffers_backend_fsync", Name: "fsync", Algo: module.Incremental},
		},
	}
	buffersAllocRateChart = module.Chart{
		ID:       "buffers_alloc_rate",
		Title:    "Buffers allocated",
		Units:    "B/s",
		Fam:      "buffers",
		Ctx:      "postgres.buffers_allocated_rate",
		Priority: prioBuffersAllocRate,
		Dims: module.Dims{
			{ID: "buffers_alloc", Name: "allocated", Algo: module.Incremental},
		},
	}

	walIORateChart = module.Chart{
		ID:       "wal_io_rate",
		Title:    "Write-Ahead Log writes",
		Units:    "B/s",
		Fam:      "wal",
		Ctx:      "postgres.wal_io_rate",
		Priority: prioWALIORate,
		Dims: module.Dims{
			{ID: "wal_writes", Name: "written", Algo: module.Incremental},
		},
	}
	walFilesCountChart = module.Chart{
		ID:       "wal_files_count",
		Title:    "Write-Ahead Log files",
		Units:    "files",
		Fam:      "wal",
		Ctx:      "postgres.wal_files_count",
		Priority: prioWALFilesCount,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "wal_written_files", Name: "written"},
			{ID: "wal_recycled_files", Name: "recycled"},
		},
	}

	walArchivingFilesCountChart = module.Chart{
		ID:       "wal_archiving_files_count",
		Title:    "Write-Ahead Log archived files",
		Units:    "files/s",
		Fam:      "wal archiving",
		Ctx:      "postgres.wal_archiving_files_count",
		Priority: prioWALArchivingFilesCount,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "wal_archive_files_ready_count", Name: "ready"},
			{ID: "wal_archive_files_done_count", Name: "done"},
		},
	}

	autovacuumWorkersCountChart = module.Chart{
		ID:       "autovacuum_workers_count",
		Title:    "Autovacuum workers",
		Units:    "workers",
		Fam:      "autovacuum",
		Ctx:      "postgres.autovacuum_workers_count",
		Priority: prioAutovacuumWorkersCount,
		Dims: module.Dims{
			{ID: "autovacuum_analyze", Name: "analyze"},
			{ID: "autovacuum_vacuum_analyze", Name: "vacuum_analyze"},
			{ID: "autovacuum_vacuum", Name: "vacuum"},
			{ID: "autovacuum_vacuum_freeze", Name: "vacuum_freeze"},
			{ID: "autovacuum_brin_summarize", Name: "brin_summarize"},
		},
	}

	txidExhaustionTowardsAutovacuumPercChart = module.Chart{
		ID:       "txid_exhaustion_towards_autovacuum_perc",
		Title:    "Percent towards emergency autovacuum",
		Units:    "percentage",
		Fam:      "txid exhaustion",
		Ctx:      "postgres.txid_exhaustion_towards_autovacuum_perc",
		Priority: prioTXIDExhaustionTowardsAutovacuumPerc,
		Dims: module.Dims{
			{ID: "percent_towards_emergency_autovacuum", Name: "emergency_autovacuum"},
		},
	}
	txidExhaustionPercChart = module.Chart{
		ID:       "txid_exhaustion_perc",
		Title:    "Percent towards transaction ID wraparound",
		Units:    "percentage",
		Fam:      "txid exhaustion",
		Ctx:      "postgres.txid_exhaustion_perc",
		Priority: prioTXIDExhaustionPerc,
		Dims: module.Dims{
			{ID: "percent_towards_wraparound", Name: "txid_exhaustion"},
		},
	}
	txidExhaustionOldestTXIDNumChart = module.Chart{
		ID:       "txid_exhaustion_oldest_txid_num",
		Title:    "Oldest transaction XID",
		Units:    "xid",
		Fam:      "txid exhaustion",
		Ctx:      "postgres.txid_exhaustion_oldest_txid_num",
		Priority: prioTXIDExhaustionOldestTXIDNum,
		Dims: module.Dims{
			{ID: "oldest_current_xid", Name: "xid"},
		},
	}

	catalogRelationSCountChart = module.Chart{
		ID:       "catalog_relations_count",
		Title:    "Relation count",
		Units:    "relations",
		Fam:      "catalog",
		Ctx:      "postgres.catalog_relations_count",
		Priority: prioCatalogRelationsCount,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "catalog_relkind_r_count", Name: "ordinary_table"},
			{ID: "catalog_relkind_i_count", Name: "index"},
			{ID: "catalog_relkind_S_count", Name: "sequence"},
			{ID: "catalog_relkind_t_count", Name: "toast_table"},
			{ID: "catalog_relkind_v_count", Name: "view"},
			{ID: "catalog_relkind_m_count", Name: "materialized_view"},
			{ID: "catalog_relkind_c_count", Name: "composite_type"},
			{ID: "catalog_relkind_f_count", Name: "foreign_table"},
			{ID: "catalog_relkind_p_count", Name: "partitioned_table"},
			{ID: "catalog_relkind_I_count", Name: "partitioned_index"},
		},
	}
	catalogRelationsSizeChart = module.Chart{
		ID:       "catalog_relations_size",
		Title:    "Relation size",
		Units:    "B",
		Fam:      "catalog",
		Ctx:      "postgres.catalog_relations_size",
		Priority: prioCatalogRelationsSize,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "catalog_relkind_r_size", Name: "ordinary_table"},
			{ID: "catalog_relkind_i_size", Name: "index"},
			{ID: "catalog_relkind_S_size", Name: "sequence"},
			{ID: "catalog_relkind_t_size", Name: "toast_table"},
			{ID: "catalog_relkind_v_size", Name: "view"},
			{ID: "catalog_relkind_m_size", Name: "materialized_view"},
			{ID: "catalog_relkind_c_size", Name: "composite_type"},
			{ID: "catalog_relkind_f_size", Name: "foreign_table"},
			{ID: "catalog_relkind_p_size", Name: "partitioned_table"},
			{ID: "catalog_relkind_I_size", Name: "partitioned_index"},
		},
	}

	serverUptimeChart = module.Chart{
		ID:       "server_uptime",
		Title:    "Uptime",
		Units:    "seconds",
		Fam:      "uptime",
		Ctx:      "postgres.uptime",
		Priority: prioUptime,
		Dims: module.Dims{
			{ID: "server_uptime", Name: "uptime"},
		},
	}

	databasesCountChart = module.Chart{
		ID:       "databases_count",
		Title:    "Number of databases",
		Units:    "databases",
		Fam:      "dbs",
		Ctx:      "postgres.databases_count",
		Priority: prioDatabasesCount,
		Dims: module.Dims{
			{ID: "databases_count", Name: "databases"},
		},
	}

	transactionsDurationChartTmpl = module.Chart{
		ID:       "transactions_duration",
		Title:    "Observed transactions time",
		Units:    "transactions/s",
		Fam:      "timings",
		Ctx:      "postgres.transactions_duration",
		Priority: prioTransactionsDuration,
		Type:     module.Stacked,
	}
	queriesDurationChartTmpl = module.Chart{
		ID:       "queries_duration",
		Title:    "Observed active queries time",
		Units:    "queries/s",
		Fam:      "timings",
		Ctx:      "postgres.queries_duration",
		Priority: prioQueriesDuration,
		Type:     module.Stacked,
	}
)

func newRunningTimeHistogramChart(tmpl module.Chart, prefix string, buckets []float64) (*module.Chart, error) {
	chart := tmpl.Copy()

	for i, v := range buckets {
		dim := &module.Dim{
			ID:   fmt.Sprintf("%s_hist_bucket_%d", prefix, i+1),
			Name: fmt.Sprintf("%.3f", v),
			Algo: module.Incremental,
		}
		if err := chart.AddDim(dim); err != nil {
			return nil, err
		}
	}

	dim := &module.Dim{
		ID:   fmt.Sprintf("%s_hist_bucket_inf", prefix),
		Name: "+Inf",
		Algo: module.Incremental,
	}
	if err := chart.AddDim(dim); err != nil {
		return nil, err
	}

	return chart, nil
}

func (p *Postgres) addTransactionsRunTimeHistogramChart() {
	chart, err := newRunningTimeHistogramChart(
		transactionsDurationChartTmpl,
		"transaction_running_time",
		p.XactTimeHistogram,
	)
	if err != nil {
		p.Warning(err)
		return
	}
	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addQueriesRunTimeHistogramChart() {
	chart, err := newRunningTimeHistogramChart(
		queriesDurationChartTmpl,
		"query_running_time",
		p.QueryTimeHistogram,
	)
	if err != nil {
		p.Warning(err)
		return
	}
	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

var (
	replicationStandbyAppCharts = module.Charts{
		replicationAppWALLagSizeChartTmpl.Copy(),
		replicationAppWALLagTimeChartTmpl.Copy(),
	}
	replicationAppWALLagSizeChartTmpl = module.Chart{
		ID:       "replication_app_%s_wal_lag_size",
		Title:    "Standby application WAL lag size",
		Units:    "B",
		Fam:      "replication lag",
		Ctx:      "postgres.replication_app_wal_lag_size",
		Priority: prioReplicationAppWALLagSize,
		Dims: module.Dims{
			{ID: "repl_standby_app_%s_wal_sent_delta", Name: "sent_delta"},
			{ID: "repl_standby_app_%s_wal_write_delta", Name: "write_delta"},
			{ID: "repl_standby_app_%s_wal_flush_delta", Name: "flush_delta"},
			{ID: "repl_standby_app_%s_wal_replay_delta", Name: "replay_delta"},
		},
	}
	replicationAppWALLagTimeChartTmpl = module.Chart{
		ID:       "replication_app_%s_wal_lag_time",
		Title:    "Standby application WAL lag time",
		Units:    "seconds",
		Fam:      "replication lag",
		Ctx:      "postgres.replication_app_wal_lag_time",
		Priority: prioReplicationAppWALLagTime,
		Dims: module.Dims{
			{ID: "repl_standby_app_%s_wal_write_lag", Name: "write_lag"},
			{ID: "repl_standby_app_%s_wal_flush_lag", Name: "flush_lag"},
			{ID: "repl_standby_app_%s_wal_replay_lag", Name: "replay_lag"},
		},
	}
)

func newReplicationStandbyAppCharts(app string) *module.Charts {
	charts := replicationStandbyAppCharts.Copy()
	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, app)
		c.Labels = []module.Label{
			{Key: "application", Value: app},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, app)
		}
	}
	return charts
}

func (p *Postgres) addNewReplicationStandbyAppCharts(app string) {
	charts := newReplicationStandbyAppCharts(app)
	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) removeReplicationStandbyAppCharts(app string) {
	prefix := fmt.Sprintf("replication_standby_app_%s_", app)
	for _, c := range *p.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}

var (
	replicationSlotCharts = module.Charts{
		replicationSlotFilesCountChartTmpl.Copy(),
	}
	replicationSlotFilesCountChartTmpl = module.Chart{
		ID:       "replication_slot_%s_files_count",
		Title:    "Replication slot files",
		Units:    "files",
		Fam:      "replication slot",
		Ctx:      "postgres.replication_slot_files_count",
		Priority: prioReplicationSlotFilesCount,
		Dims: module.Dims{
			{ID: "repl_slot_%s_replslot_wal_keep", Name: "wal_keep"},
			{ID: "repl_slot_%s_replslot_files", Name: "pg_replslot_files"},
		},
	}
)

func newReplicationSlotCharts(slot string) *module.Charts {
	charts := replicationSlotCharts.Copy()
	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, slot)
		c.Labels = []module.Label{
			{Key: "slot", Value: slot},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, slot)
		}
	}
	return charts
}

func (p *Postgres) addNewReplicationSlotCharts(slot string) {
	charts := newReplicationSlotCharts(slot)
	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) removeReplicationSlotCharts(slot string) {
	prefix := fmt.Sprintf("replication_slot_%s_", slot)
	for _, c := range *p.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}

var (
	dbChartsTmpl = module.Charts{
		dbTransactionsRatioChartTmpl.Copy(),
		dbTransactionsRateChartTmpl.Copy(),
		dbConnectionsUtilizationChartTmpl.Copy(),
		dbConnectionsCountChartTmpl.Copy(),
		dbCacheIORatioChartTmpl.Copy(),
		dbIORateChartTmpl.Copy(),
		dbOpsFetchedRowsRatioChartTmpl.Copy(),
		dbOpsReadRowsRateChartTmpl.Copy(),
		dbOpsWriteRowsRateChartTmpl.Copy(),
		dbConflictsRateChartTmpl.Copy(),
		dbConflictsReasonRateChartTmpl.Copy(),
		dbDeadlocksRateChartTmpl.Copy(),
		dbLocksHeldCountChartTmpl.Copy(),
		dbLocksAwaitedCountChartTmpl.Copy(),
		dbTempFilesCreatedRateChartTmpl.Copy(),
		dbTempFilesIORateChartTmpl.Copy(),
		dbSizeChartTmpl.Copy(),
	}
	dbTransactionsRatioChartTmpl = module.Chart{
		ID:       "db_%s_transactions_ratio",
		Title:    "Database transactions ratio",
		Units:    "percentage",
		Fam:      "db transactions",
		Ctx:      "postgres.db_transactions_ratio",
		Priority: prioDBTransactionsRatio,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "db_%s_xact_commit", Name: "committed", Algo: module.PercentOfIncremental},
			{ID: "db_%s_xact_rollback", Name: "rollback", Algo: module.PercentOfIncremental},
		},
	}
	dbTransactionsRateChartTmpl = module.Chart{
		ID:       "db_%s_transactions_rate",
		Title:    "Database transactions",
		Units:    "transactions/s",
		Fam:      "db transactions",
		Ctx:      "postgres.db_transactions_rate",
		Priority: prioDBTransactionsRate,
		Dims: module.Dims{
			{ID: "db_%s_xact_commit", Name: "committed", Algo: module.Incremental},
			{ID: "db_%s_xact_rollback", Name: "rollback", Algo: module.Incremental},
		},
	}
	dbConnectionsUtilizationChartTmpl = module.Chart{
		ID:       "db_%s_connections_utilization",
		Title:    "Database connections utilization",
		Units:    "percentage",
		Fam:      "db connections",
		Ctx:      "postgres.db_connections_utilization",
		Priority: prioDBConnectionsUtilization,
		Dims: module.Dims{
			{ID: "db_%s_numbackends_utilization", Name: "used"},
		},
	}
	dbConnectionsCountChartTmpl = module.Chart{
		ID:       "db_%s_connections",
		Title:    "Database connections",
		Units:    "connections",
		Fam:      "db connections",
		Ctx:      "postgres.db_connections_count",
		Priority: prioDBConnectionsCount,
		Dims: module.Dims{
			{ID: "db_%s_numbackends", Name: "connections"},
		},
	}
	dbCacheIORatioChartTmpl = module.Chart{
		ID:       "db_%s_cache_io_ratio",
		Title:    "Database buffer cache miss ratio",
		Units:    "percentage",
		Fam:      "db io",
		Ctx:      "postgres.db_cache_io_ratio",
		Priority: prioDBCacheIORatio,
		Dims: module.Dims{
			{ID: "db_%s_blks_read_perc", Name: "miss"},
		},
	}
	dbIORateChartTmpl = module.Chart{
		ID:       "db_%s_io_rate",
		Title:    "Database reads",
		Units:    "B/s",
		Fam:      "db io",
		Ctx:      "postgres.db_io_rate",
		Priority: prioDBIORate,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "db_%s_blks_hit", Name: "memory", Algo: module.Incremental},
			{ID: "db_%s_blks_read", Name: "disk", Algo: module.Incremental},
		},
	}
	dbOpsFetchedRowsRatioChartTmpl = module.Chart{
		ID:       "db_%s_db_ops_fetched_rows_ratio",
		Title:    "Database rows fetched ratio",
		Units:    "percentage",
		Fam:      "db read throughput",
		Ctx:      "postgres.db_ops_fetched_rows_ratio",
		Priority: prioDBOpsFetchedRowsRatio,
		Dims: module.Dims{
			{ID: "db_%s_tup_fetched_perc", Name: "fetched"},
		},
	}
	dbOpsReadRowsRateChartTmpl = module.Chart{
		ID:       "db_%s_ops_read_rows_rate",
		Title:    "Database rows read",
		Units:    "rows/s",
		Fam:      "db read throughput",
		Ctx:      "postgres.db_ops_read_rows_rate",
		Priority: prioDBOpsReadRowsRate,
		Dims: module.Dims{
			{ID: "db_%s_tup_returned", Name: "returned", Algo: module.Incremental},
			{ID: "db_%s_tup_fetched", Name: "fetched", Algo: module.Incremental},
		},
	}
	dbOpsWriteRowsRateChartTmpl = module.Chart{
		ID:       "db_%s_ops_write_rows_rate",
		Title:    "Database rows written",
		Units:    "rows/s",
		Fam:      "db write throughput",
		Ctx:      "postgres.db_ops_write_rows_rate",
		Priority: prioDBOpsWriteRowsRate,
		Dims: module.Dims{
			{ID: "db_%s_tup_inserted", Name: "inserted", Algo: module.Incremental},
			{ID: "db_%s_tup_deleted", Name: "deleted", Algo: module.Incremental},
			{ID: "db_%s_tup_updated", Name: "updated", Algo: module.Incremental},
		},
	}
	dbConflictsRateChartTmpl = module.Chart{
		ID:       "db_%s_conflicts_rate",
		Title:    "Database canceled queries",
		Units:    "queries/s",
		Fam:      "db query cancels",
		Ctx:      "postgres.db_conflicts_rate",
		Priority: prioDBConflictsRate,
		Dims: module.Dims{
			{ID: "db_%s_conflicts", Name: "conflicts", Algo: module.Incremental},
		},
	}
	dbConflictsReasonRateChartTmpl = module.Chart{
		ID:       "db_%s_conflicts_reason_rate",
		Title:    "Database canceled queries by reason",
		Units:    "queries/s",
		Fam:      "db query cancels",
		Ctx:      "postgres.db_conflicts_reason_rate",
		Priority: prioDBConflictsReasonRate,
		Dims: module.Dims{
			{ID: "db_%s_confl_tablespace", Name: "tablespace", Algo: module.Incremental},
			{ID: "db_%s_confl_lock", Name: "lock", Algo: module.Incremental},
			{ID: "db_%s_confl_snapshot", Name: "snapshot", Algo: module.Incremental},
			{ID: "db_%s_confl_bufferpin", Name: "bufferpin", Algo: module.Incremental},
			{ID: "db_%s_confl_deadlock", Name: "deadlock", Algo: module.Incremental},
		},
	}
	dbDeadlocksRateChartTmpl = module.Chart{
		ID:       "db_%s_deadlocks_rate",
		Title:    "Database deadlocks",
		Units:    "deadlocks/s",
		Fam:      "db deadlocks",
		Ctx:      "postgres.db_deadlocks_rate",
		Priority: prioDBDeadlocksRate,
		Dims: module.Dims{
			{ID: "db_%s_deadlocks", Name: "deadlocks", Algo: module.Incremental},
		},
	}
	dbLocksHeldCountChartTmpl = module.Chart{
		ID:       "db_%s_locks_held",
		Title:    "Database locks held",
		Units:    "locks",
		Fam:      "db locks",
		Ctx:      "postgres.db_locks_held_count",
		Priority: prioDBLocksHeldCount,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "db_%s_lock_mode_AccessShareLock_held", Name: "access_share"},
			{ID: "db_%s_lock_mode_RowShareLock_held", Name: "row_share"},
			{ID: "db_%s_lock_mode_RowExclusiveLock_held", Name: "row_exclusive"},
			{ID: "db_%s_lock_mode_ShareUpdateExclusiveLock_held", Name: "share_update"},
			{ID: "db_%s_lock_mode_ShareLock_held", Name: "share"},
			{ID: "db_%s_lock_mode_ShareRowExclusiveLock_held", Name: "share_row_exclusive"},
			{ID: "db_%s_lock_mode_ExclusiveLock_held", Name: "exclusive"},
			{ID: "db_%s_lock_mode_AccessExclusiveLock_held", Name: "access_exclusive"},
		},
	}
	dbLocksAwaitedCountChartTmpl = module.Chart{
		ID:       "db_%s_locks_awaited_count",
		Title:    "Database locks awaited",
		Units:    "locks",
		Fam:      "db locks",
		Ctx:      "postgres.db_locks_awaited_count",
		Priority: prioDBLocksAwaitedCount,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "db_%s_lock_mode_AccessShareLock_awaited", Name: "access_share"},
			{ID: "db_%s_lock_mode_RowShareLock_awaited", Name: "row_share"},
			{ID: "db_%s_lock_mode_RowExclusiveLock_awaited", Name: "row_exclusive"},
			{ID: "db_%s_lock_mode_ShareUpdateExclusiveLock_awaited", Name: "share_update"},
			{ID: "db_%s_lock_mode_ShareLock_awaited", Name: "share"},
			{ID: "db_%s_lock_mode_ShareRowExclusiveLock_awaited", Name: "share_row_exclusive"},
			{ID: "db_%s_lock_mode_ExclusiveLock_awaited", Name: "exclusive"},
			{ID: "db_%s_lock_mode_AccessExclusiveLock_awaited", Name: "access_exclusive"},
		},
	}
	dbTempFilesCreatedRateChartTmpl = module.Chart{
		ID:       "db_%s_temp_files_files_created_rate",
		Title:    "Database created temporary files",
		Units:    "files/s",
		Fam:      "db temp files",
		Ctx:      "postgres.db_temp_files_created_rate",
		Priority: prioDBTempFilesCreatedRate,
		Dims: module.Dims{
			{ID: "db_%s_temp_files", Name: "created", Algo: module.Incremental},
		},
	}
	dbTempFilesIORateChartTmpl = module.Chart{
		ID:       "db_%s_temp_files_io_rate",
		Title:    "Database temporary files data written to disk",
		Units:    "B/s",
		Fam:      "db temp files",
		Ctx:      "postgres.db_temp_files_io_rate",
		Priority: prioDBTempFilesIORate,
		Dims: module.Dims{
			{ID: "db_%s_temp_bytes", Name: "written", Algo: module.Incremental},
		},
	}
	dbSizeChartTmpl = module.Chart{
		ID:       "db_%s_size",
		Title:    "Database size",
		Units:    "B",
		Fam:      "db size",
		Ctx:      "postgres.db_size",
		Priority: prioDBSize,
		Dims: module.Dims{
			{ID: "db_%s_size", Name: "size"},
		},
	}
)

func newDatabaseCharts(dbname string) *module.Charts {
	charts := dbChartsTmpl.Copy()
	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, dbname)
		c.Labels = []module.Label{
			{Key: "database", Value: dbname},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, dbname)
		}
	}
	return charts
}

func (p *Postgres) addNewDatabaseCharts(dbname string) {
	charts := newDatabaseCharts(dbname)
	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) removeDatabaseCharts(dbname string) {
	prefix := fmt.Sprintf("db_%s_", dbname)
	for _, c := range *p.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}

var (
	tableChartsTmpl = module.Charts{
		tableRowsCountChartTmpl.Copy(),
		tableDeadRowsDeadRatioChartTmpl.Copy(),
		tableOpsRowsRateChartTmpl.Copy(),
		tableOpsRowsHOTRatioChartTmpl.Copy(),
		tableOpsRowsHOTRateChartTmpl.Copy(),
		tableScansRateChartTmpl.Copy(),
		tableScansRowsRateChartTmpl.Copy(),
		tableNullColumnsCountChartTmpl.Copy(),
		tableTotalSizeChartTmpl.Copy(),
		tableBloatSizePercChartTmpl.Copy(),
		tableBloatSizeChartTmpl.Copy(),
	}

	tableDeadRowsDeadRatioChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_rows_dead_ratio",
		Title:    "Table dead rows",
		Units:    "%",
		Fam:      "table rows",
		Ctx:      "postgres.table_rows_dead_ratio",
		Priority: prioTableRowsDeadRatio,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_n_dead_tup_perc", Name: "dead"},
		},
	}
	tableRowsCountChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_rows_count",
		Title:    "Table total rows",
		Units:    "rows",
		Fam:      "table rows",
		Ctx:      "postgres.table_rows_count",
		Priority: prioTableRowsCount,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_n_live_tup", Name: "live"},
			{ID: "table_%s_db_%s_schema_%s_n_dead_tup", Name: "dead"},
		},
	}
	tableOpsRowsRateChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_ops_rows_rate",
		Title:    "Table throughput",
		Units:    "rows/s",
		Fam:      "table throughput",
		Ctx:      "postgres.table_ops_rows_rate",
		Priority: prioTableOpsRowsRate,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_n_tup_ins", Name: "inserted", Algo: module.Incremental},
			{ID: "table_%s_db_%s_schema_%s_n_tup_del", Name: "deleted", Algo: module.Incremental},
			{ID: "table_%s_db_%s_schema_%s_n_tup_upd", Name: "updated", Algo: module.Incremental},
		},
	}
	tableOpsRowsHOTRatioChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_ops_rows_hot_ratio",
		Title:    "Table HOT updates ratio",
		Units:    "percentage",
		Fam:      "table hot updates",
		Ctx:      "postgres.table_ops_rows_hot_ratio",
		Priority: prioTableOpsRowsHOTRatio,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_n_tup_hot_upd_perc", Name: "hot"},
		},
	}
	tableOpsRowsHOTRateChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_ops_rows_hot_rate",
		Title:    "Table HOT updates",
		Units:    "rows/s",
		Fam:      "table hot updates",
		Ctx:      "postgres.table_ops_rows_hot_rate",
		Priority: prioTableOpsRowsHOTRate,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_n_tup_hot_upd", Name: "hot", Algo: module.Incremental},
		},
	}
	tableCacheIORatioChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_cache_io_ratio",
		Title:    "Table I/O cache miss ratio",
		Units:    "percentage",
		Fam:      "table io",
		Ctx:      "postgres.table_cache_io_ratio",
		Priority: prioTableCacheIORatio,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_heap_blks_read_perc", Name: "miss"},
		},
	}
	tableIORateChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_io_rate",
		Title:    "Table I/O",
		Units:    "B/s",
		Fam:      "table io",
		Ctx:      "postgres.table_io_rate",
		Priority: prioTableIORate,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_heap_blks_hit", Name: "memory", Algo: module.Incremental},
			{ID: "table_%s_db_%s_schema_%s_heap_blks_read", Name: "disk", Algo: module.Incremental},
		},
	}
	tableIndexCacheIORatioChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_index_cache_io_ratio",
		Title:    "Table index I/O cache miss ratio",
		Units:    "percentage",
		Fam:      "table index io",
		Ctx:      "postgres.table_index_cache_io_ratio",
		Priority: prioTableIndexCacheIORatio,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_idx_blks_read_perc", Name: "miss", Algo: module.Incremental},
		},
	}
	tableIndexIORateChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_index_io_rate",
		Title:    "Table index I/O",
		Units:    "B/s",
		Fam:      "table index io",
		Ctx:      "postgres.table_index_io_rate",
		Priority: prioTableIndexIORate,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_idx_blks_hit", Name: "memory", Algo: module.Incremental},
			{ID: "table_%s_db_%s_schema_%s_idx_blks_read", Name: "disk", Algo: module.Incremental},
		},
	}
	tableTOASCacheIORatioChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_toast_cache_io_ratio",
		Title:    "Table TOAST I/O cache miss ratio",
		Units:    "percentage",
		Fam:      "table toast io",
		Ctx:      "postgres.table_toast_cache_io_ratio",
		Priority: prioTableToastCacheIORatio,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_toast_blks_read_perc", Name: "miss", Algo: module.Incremental},
		},
	}
	tableTOASTIORateChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_toast_io_rate",
		Title:    "Table TOAST I/O",
		Units:    "B/s",
		Fam:      "table toast io",
		Ctx:      "postgres.table_toast_io_rate",
		Priority: prioTableToastIORate,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_toast_blks_hit", Name: "memory", Algo: module.Incremental},
			{ID: "table_%s_db_%s_schema_%s_toast_blks_read", Name: "disk", Algo: module.Incremental},
		},
	}
	tableTOASTIndexCacheIORatioChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_toast_index_cache_io_ratio",
		Title:    "Table TOAST index I/O cache miss ratio",
		Units:    "percentage",
		Fam:      "table toast index io",
		Ctx:      "postgres.table_toast_index_cache_io_ratio",
		Priority: prioTableToastIndexCacheIORatio,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_tidx_blks_read_perc", Name: "miss", Algo: module.Incremental},
		},
	}
	tableTOASTIndexIORateChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_toast_index_io_rate",
		Title:    "Table TOAST index I/O",
		Units:    "B/s",
		Fam:      "table toast index io",
		Ctx:      "postgres.table_toast_index_io_rate",
		Priority: prioTableToastIndexIORate,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_tidx_blks_hit", Name: "memory", Algo: module.Incremental},
			{ID: "table_%s_db_%s_schema_%s_tidx_blks_read", Name: "disk", Algo: module.Incremental},
		},
	}
	tableScansRateChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_scans_rate",
		Title:    "Table scans",
		Units:    "scans/s",
		Fam:      "table scans",
		Ctx:      "postgres.table_scans_rate",
		Priority: prioTableScansRate,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_idx_scan", Name: "index", Algo: module.Incremental},
			{ID: "table_%s_db_%s_schema_%s_seq_scan", Name: "sequential", Algo: module.Incremental},
		},
	}
	tableScansRowsRateChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_scans_rows_rate",
		Title:    "Table live rows fetched by scans",
		Units:    "rows/s",
		Fam:      "table scans",
		Ctx:      "postgres.table_scans_rows_rate",
		Priority: prioTableScansRowsRate,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_idx_tup_fetch", Name: "index", Algo: module.Incremental},
			{ID: "table_%s_db_%s_schema_%s_seq_tup_read", Name: "sequential", Algo: module.Incremental},
		},
	}
	tableAutoVacuumSinceTimeChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_autovacuum_since_time",
		Title:    "Table time since last auto VACUUM",
		Units:    "seconds",
		Fam:      "table autovacuum",
		Ctx:      "postgres.table_autovacuum_since_time",
		Priority: prioTableAutovacuumSinceTime,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_last_autovacuum_ago", Name: "time"},
		},
	}
	tableVacuumSinceTimeChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_vacuum_since_time",
		Title:    "Table time since last manual VACUUM",
		Units:    "seconds",
		Fam:      "table vacuum",
		Ctx:      "postgres.table_vacuum_since_time",
		Priority: prioTableVacuumSinceTime,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_last_vacuum_ago", Name: "time"},
		},
	}
	tableAutoAnalyzeSinceTimeChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_autoanalyze_since_time",
		Title:    "Table time since last auto ANALYZE",
		Units:    "seconds",
		Fam:      "table autoanalyze",
		Ctx:      "postgres.table_autoanalyze_since_time",
		Priority: prioTableAutoAnalyzeSinceTime,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_last_autoanalyze_ago", Name: "time"},
		},
	}
	tableAnalyzeSinceTimeChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_analyze_since_time",
		Title:    "Table time since last manual ANALYZE",
		Units:    "seconds",
		Fam:      "table analyze",
		Ctx:      "postgres.table_analyze_since_time",
		Priority: prioTableLastAnalyzeAgo,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_last_analyze_ago", Name: "time"},
		},
	}
	tableNullColumnsCountChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_null_columns_count",
		Title:    "Table null columns",
		Units:    "columns",
		Fam:      "table columns",
		Ctx:      "postgres.table_null_columns_count",
		Priority: prioTableNullColumns,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_null_columns", Name: "null"},
		},
	}
	tableTotalSizeChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_total_size",
		Title:    "Table total size",
		Units:    "B",
		Fam:      "table size",
		Ctx:      "postgres.table_total_size",
		Priority: prioTableTotalSize,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_total_size", Name: "size"},
		},
	}
	tableBloatSizePercChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_bloat_size_perc",
		Title:    "Table bloat size percentage",
		Units:    "percentage",
		Fam:      "table bloat",
		Ctx:      "postgres.table_bloat_size_perc",
		Priority: prioTableBloatSizePerc,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_bloat_size_perc", Name: "bloat"},
		},
	}
	tableBloatSizeChartTmpl = module.Chart{
		ID:       "table_%s_db_%s_schema_%s_bloat_size",
		Title:    "Table bloat size",
		Units:    "B",
		Fam:      "table bloat",
		Ctx:      "postgres.table_bloat_size",
		Priority: prioTableBloatSize,
		Dims: module.Dims{
			{ID: "table_%s_db_%s_schema_%s_bloat_size", Name: "bloat"},
		},
	}
)

func newTableCharts(name, dbname, schema string) *module.Charts {
	charts := tableChartsTmpl.Copy()

	for i, chart := range *charts {
		(*charts)[i] = newTableChart(chart, name, dbname, schema)
	}

	return charts
}

func newTableChart(chart *module.Chart, name, dbname, schema string) *module.Chart {
	chart = chart.Copy()
	chart.ID = fmt.Sprintf(chart.ID, name, dbname, schema)
	chart.Labels = []module.Label{
		{Key: "database", Value: dbname},
		{Key: "schema", Value: schema},
		{Key: "table", Value: name},
	}
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, name, dbname, schema)
	}
	return chart
}

func (p *Postgres) addNewTableCharts(name, dbname, schema string) {
	charts := newTableCharts(name, dbname, schema)
	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableLastAutoVacuumAgoChart(name, dbname, schema string) {
	chart := newTableChart(tableAutoVacuumSinceTimeChartTmpl.Copy(), name, dbname, schema)

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableLastVacuumAgoChart(name, dbname, schema string) {
	chart := newTableChart(tableVacuumSinceTimeChartTmpl.Copy(), name, dbname, schema)

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableLastAutoAnalyzeAgoChart(name, dbname, schema string) {
	chart := newTableChart(tableAutoAnalyzeSinceTimeChartTmpl.Copy(), name, dbname, schema)

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableLastAnalyzeAgoChart(name, dbname, schema string) {
	chart := newTableChart(tableAnalyzeSinceTimeChartTmpl.Copy(), name, dbname, schema)

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableIOChartsCharts(name, dbname, schema string) {
	charts := module.Charts{
		newTableChart(tableCacheIORatioChartTmpl.Copy(), name, dbname, schema),
		newTableChart(tableIORateChartTmpl.Copy(), name, dbname, schema),
	}

	if err := p.Charts().Add(charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableIndexIOCharts(name, dbname, schema string) {
	charts := module.Charts{
		newTableChart(tableIndexCacheIORatioChartTmpl.Copy(), name, dbname, schema),
		newTableChart(tableIndexIORateChartTmpl.Copy(), name, dbname, schema),
	}

	if err := p.Charts().Add(charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableTOASTIOCharts(name, dbname, schema string) {
	charts := module.Charts{
		newTableChart(tableTOASCacheIORatioChartTmpl.Copy(), name, dbname, schema),
		newTableChart(tableTOASTIORateChartTmpl.Copy(), name, dbname, schema),
	}

	if err := p.Charts().Add(charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableTOASTIndexIOCharts(name, dbname, schema string) {
	charts := module.Charts{
		newTableChart(tableTOASTIndexCacheIORatioChartTmpl.Copy(), name, dbname, schema),
		newTableChart(tableTOASTIndexIORateChartTmpl.Copy(), name, dbname, schema),
	}

	if err := p.Charts().Add(charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) removeTableCharts(name, dbname, schema string) {
	prefix := fmt.Sprintf("table_%s_db_%s_schema_%s", name, dbname, schema)
	for _, c := range *p.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}

var (
	indexChartsTmpl = module.Charts{
		indexSizeChartTmpl.Copy(),
		indexBloatSizePercChartTmpl.Copy(),
		indexBloatSizeChartTmpl.Copy(),
		indexUsageStatusChartTmpl.Copy(),
	}
	indexSizeChartTmpl = module.Chart{
		ID:       "index_%s_table_%s_db_%s_schema_%s_size",
		Title:    "Index size",
		Units:    "B",
		Fam:      "index size",
		Ctx:      "postgres.index_size",
		Priority: prioIndexSize,
		Dims: module.Dims{
			{ID: "index_%s_table_%s_db_%s_schema_%s_size", Name: "size"},
		},
	}
	indexBloatSizePercChartTmpl = module.Chart{
		ID:       "index_%s_table_%s_db_%s_schema_%s_bloat_size_perc",
		Title:    "Index bloat size percentage",
		Units:    "percentage",
		Fam:      "index bloat",
		Ctx:      "postgres.index_bloat_size_perc",
		Priority: prioIndexBloatSizePerc,
		Dims: module.Dims{
			{ID: "index_%s_table_%s_db_%s_schema_%s_bloat_size_perc", Name: "bloat"},
		},
	}
	indexBloatSizeChartTmpl = module.Chart{
		ID:       "index_%s_table_%s_db_%s_schema_%s_bloat_size",
		Title:    "Index bloat size",
		Units:    "B",
		Fam:      "index bloat",
		Ctx:      "postgres.index_bloat_size",
		Priority: prioIndexBloatSize,
		Dims: module.Dims{
			{ID: "index_%s_table_%s_db_%s_schema_%s_bloat_size", Name: "bloat"},
		},
	}
	indexUsageStatusChartTmpl = module.Chart{
		ID:       "index_%s_table_%s_db_%s_schema_%s_usage_status",
		Title:    "Index usage status",
		Units:    "status",
		Fam:      "index usage status",
		Ctx:      "postgres.index_usage_status",
		Priority: prioIndexUsageStatus,
		Dims: module.Dims{
			{ID: "index_%s_table_%s_db_%s_schema_%s_usage_status_used", Name: "used"},
			{ID: "index_%s_table_%s_db_%s_schema_%s_usage_status_unused", Name: "unused"},
		},
	}
)

func (p *Postgres) addNewIndexCharts(name, table, dbname, schema string) {
	charts := indexChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, name, table, dbname, schema)
		chart.Labels = []module.Label{
			{Key: "database", Value: dbname},
			{Key: "schema", Value: schema},
			{Key: "table", Value: table},
			{Key: "index", Value: name},
		}
		for _, d := range chart.Dims {
			d.ID = fmt.Sprintf(d.ID, name, table, dbname, schema)
		}
	}

	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) removeIndexCharts(name, table, dbname, schema string) {
	prefix := fmt.Sprintf("index_%s_table_%s_db_%s_schema_%s", name, table, dbname, schema)
	for _, c := range *p.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}
