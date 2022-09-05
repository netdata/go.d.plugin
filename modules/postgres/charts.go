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
	prioConnectionsState
	prioCheckpoints
	prioCheckpointTime
	prioBGWriterBuffersAllocated
	prioBGWriterBuffersWritten
	prioBGWriterMaxWrittenClean
	prioBGWriterBackedFsync
	prioWALWrites
	prioWALFiles
	prioWALArchive
	prioAutovacuumWorkers
	prioAutovacuumPercentTowards
	prioTXIDWraparoundPercentTowards
	prioTXIDWraparoundOldestTXID
	prioCatalogRelationCount
	prioCatalogRelationSize
	prioUptime
	prioReplicationWALDelta
	prioReplicationWALLag
	prioReplicationSlotFiles
	prioDBTransactionsRatio
	prioDBTransactions
	prioDBConnectionsUtilization
	prioDBConnections
	prioDBBufferCacheRatio
	prioDBBlockReads
	prioDBRowsReadRatio
	prioDBRowsRead
	prioDBRowsWritten
	prioDBConflicts
	prioDBConflictsStat
	prioDBDeadlocks
	prioDBLocksHeld
	prioDBLocksAwaited
	prioDBTempFiles
	prioDBTempFilesData
	prioDBSize
	prioTableRowsTotalPercentage
	prioTableRowsTotal
	prioTableRowsOperations
	prioTableHOTUpdates
	prioTableScans
	prioTableScansRows
	prioTableLastAutovacuumAgo
	prioTableLastVacuumAgo
	prioTableLastAutoAnalyzeAgo
	prioTableLastAnalyzeAgo
	prioTableTotalSize
)

var baseCharts = module.Charts{
	serverConnectionsUtilizationChart.Copy(),
	serverConnectionsUsageChart.Copy(),
	serverConnectionsState.Copy(),
	checkpointsChart.Copy(),
	checkpointWriteChart.Copy(),
	bgWriterBuffersWrittenChart.Copy(),
	bgWriterBuffersAllocChart.Copy(),
	bgWriterMaxWrittenCleanChart.Copy(),
	bgWriterBuffersBackendFsyncChart.Copy(),
	walWritesChart.Copy(),
	walFilesChart.Copy(),
	walArchiveFilesChart.Copy(),
	autovacuumWorkersChart.Copy(),
	percentTowardsEmergencyAutovacuumChart.Copy(),
	percentTowardTXIDWraparoundChart.Copy(),
	oldestTXIDChart.Copy(),

	catalogRelationCountChart.Copy(),
	catalogRelationSizeChart.Copy(),
	serverUptimeChart.Copy(),
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
	serverConnectionsState = module.Chart{
		ID:       "connections_state",
		Title:    "Connections in each state",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "postgres.connections_state",
		Priority: prioConnectionsState,
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
		ID:       "checkpoints",
		Title:    "Checkpoints",
		Units:    "checkpoints/s",
		Fam:      "checkpointer",
		Ctx:      "postgres.checkpoints",
		Priority: prioCheckpoints,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "checkpoints_timed", Name: "scheduled", Algo: module.Incremental},
			{ID: "checkpoints_req", Name: "requested", Algo: module.Incremental},
		},
	}
	// TODO: should be seconds, also it is units/s when using incremental...
	checkpointWriteChart = module.Chart{
		ID:       "checkpoint_time",
		Title:    "Checkpoint time",
		Units:    "milliseconds",
		Fam:      "checkpointer",
		Ctx:      "postgres.checkpoint_time",
		Priority: prioCheckpointTime,
		Dims: module.Dims{
			{ID: "checkpoint_write_time", Name: "write", Algo: module.Incremental},
			{ID: "checkpoint_sync_time", Name: "sync", Algo: module.Incremental},
		},
	}

	bgWriterBuffersAllocChart = module.Chart{
		ID:       "bgwriter_buffers_alloc",
		Title:    "Background writer buffers allocated",
		Units:    "B/s",
		Fam:      "background writer",
		Ctx:      "postgres.bgwriter_buffers_alloc",
		Priority: prioBGWriterBuffersAllocated,
		Dims: module.Dims{
			{ID: "buffers_alloc", Name: "allocated", Algo: module.Incremental},
		},
	}
	bgWriterBuffersWrittenChart = module.Chart{
		ID:       "bgwriter_buffers_written",
		Title:    "Background writer buffers written",
		Units:    "B/s",
		Fam:      "background writer",
		Ctx:      "postgres.bgwriter_buffers_written",
		Priority: prioBGWriterBuffersWritten,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "buffers_checkpoint", Name: "checkpoint", Algo: module.Incremental},
			{ID: "buffers_backend", Name: "backend", Algo: module.Incremental},
			{ID: "buffers_clean", Name: "clean", Algo: module.Incremental},
		},
	}
	bgWriterMaxWrittenCleanChart = module.Chart{
		ID:       "bgwriter_maxwritten_clean",
		Title:    "Background writer cleaning scan stops",
		Units:    "events/s",
		Fam:      "background writer",
		Ctx:      "postgres.bgwriter_maxwritten_clean",
		Priority: prioBGWriterMaxWrittenClean,
		Dims: module.Dims{
			{ID: "maxwritten_clean", Name: "maxwritten", Algo: module.Incremental},
		},
	}
	bgWriterBuffersBackendFsyncChart = module.Chart{
		ID:       "bgwriter_buffers_backend_fsync",
		Title:    "Backend fsync",
		Units:    "operations/s",
		Fam:      "background writer",
		Ctx:      "postgres.bgwriter_buffers_backend_fsync",
		Priority: prioBGWriterBackedFsync,
		Dims: module.Dims{
			{ID: "buffers_backend_fsync", Name: "fsync", Algo: module.Incremental},
		},
	}

	walWritesChart = module.Chart{
		ID:       "wal_writes",
		Title:    "Write-Ahead Log",
		Units:    "B/s",
		Fam:      "wal",
		Ctx:      "postgres.wal_writes",
		Priority: prioWALWrites,
		Dims: module.Dims{
			{ID: "wal_writes", Name: "writes", Algo: module.Incremental},
		},
	}
	walFilesChart = module.Chart{
		ID:       "wal_files",
		Title:    "Write-Ahead Log files",
		Units:    "files",
		Fam:      "wal",
		Ctx:      "postgres.wal_files",
		Priority: prioWALFiles,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "wal_written_files", Name: "written"},
			{ID: "wal_recycled_files", Name: "recycled"},
		},
	}

	walArchiveFilesChart = module.Chart{
		ID:       "wal_archive_files",
		Title:    "Write-Ahead Log archive files",
		Units:    "files/s",
		Fam:      "wal archive",
		Ctx:      "postgres.wal_archive_files",
		Priority: prioWALArchive,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "wal_archive_files_ready_count", Name: "ready", Algo: module.Incremental},
			{ID: "wal_archive_files_done_count", Name: "done", Algo: module.Incremental},
		},
	}

	autovacuumWorkersChart = module.Chart{
		ID:       "autovacuum_workers",
		Title:    "Autovacuum workers",
		Units:    "workers",
		Fam:      "autovacuum",
		Ctx:      "postgres.autovacuum_workers",
		Priority: prioAutovacuumWorkers,
		Dims: module.Dims{
			{ID: "autovacuum_analyze", Name: "analyze"},
			{ID: "autovacuum_vacuum_analyze", Name: "vacuum_analyze"},
			{ID: "autovacuum_vacuum", Name: "vacuum"},
			{ID: "autovacuum_vacuum_freeze", Name: "vacuum_freeze"},
			{ID: "autovacuum_brin_summarize", Name: "brin_summarize"},
		},
	}
	percentTowardsEmergencyAutovacuumChart = module.Chart{
		ID:       "percent_towards_emergency_autovacuum",
		Title:    "Percent towards emergency autovacuum",
		Units:    "percentage",
		Fam:      "autovacuum",
		Ctx:      "postgres.percent_towards_emergency_autovacuum",
		Priority: prioAutovacuumPercentTowards,
		Dims: module.Dims{
			{ID: "percent_towards_emergency_autovacuum", Name: "emergency_autovacuum"},
		},
	}

	percentTowardTXIDWraparoundChart = module.Chart{
		ID:       "percent_towards_txid_wraparound",
		Title:    "Percent towards transaction ID wraparound",
		Units:    "percentage",
		Fam:      "txid wraparound",
		Ctx:      "postgres.percent_towards_txid_wraparound",
		Priority: prioTXIDWraparoundPercentTowards,
		Dims: module.Dims{
			{ID: "percent_towards_wraparound", Name: "txid_wraparound"},
		},
	}
	oldestTXIDChart = module.Chart{
		ID:       "oldest_transaction_xid",
		Title:    "Oldest transaction XID",
		Units:    "xid",
		Fam:      "txid wraparound",
		Ctx:      "postgres.oldest_transaction_xid",
		Priority: prioTXIDWraparoundOldestTXID,
		Dims: module.Dims{
			{ID: "oldest_current_xid", Name: "xid"},
		},
	}

	catalogRelationCountChart = module.Chart{
		ID:       "catalog_relation_count",
		Title:    "Relation count",
		Units:    "relations",
		Fam:      "catalog",
		Ctx:      "postgres.catalog_relation_count",
		Priority: prioCatalogRelationCount,
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
	catalogRelationSizeChart = module.Chart{
		ID:       "catalog_relation_size",
		Title:    "Relation size",
		Units:    "B",
		Fam:      "catalog",
		Ctx:      "postgres.catalog_relation_size",
		Priority: prioCatalogRelationSize,
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
)

var (
	replicationStandbyAppCharts = module.Charts{
		replicationStandbyAppWALDeltaChartTmpl.Copy(),
		replicationStandbyAppWALLagChartTmpl.Copy(),
	}
	replicationStandbyAppWALDeltaChartTmpl = module.Chart{
		ID:       "replication_standby_app_%s_wal_delta",
		Title:    "Standby application WAL delta",
		Units:    "B",
		Fam:      "replication delta",
		Ctx:      "postgres.replication_standby_app_wal_delta",
		Priority: prioReplicationWALDelta,
		Dims: module.Dims{
			{ID: "repl_standby_app_%s_wal_sent_delta", Name: "sent_delta"},
			{ID: "repl_standby_app_%s_wal_write_delta", Name: "write_delta"},
			{ID: "repl_standby_app_%s_wal_flush_delta", Name: "flush_delta"},
			{ID: "repl_standby_app_%s_wal_replay_delta", Name: "replay_delta"},
		},
	}
	replicationStandbyAppWALLagChartTmpl = module.Chart{
		ID:       "replication_standby_app_%s_wal_lag",
		Title:    "Standby application WAL lag",
		Units:    "seconds",
		Fam:      "replication lag",
		Ctx:      "postgres.replication_standby_app_wal_lag",
		Priority: prioReplicationWALLag,
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
		replicationSlotFilesChartTmpl.Copy(),
	}
	replicationSlotFilesChartTmpl = module.Chart{
		ID:       "replication_slot_%s_files",
		Title:    "Replication slot files",
		Units:    "files",
		Fam:      "replication slot",
		Ctx:      "postgres.replication_slot_files",
		Priority: prioReplicationSlotFiles,
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
		dbTransactionsChartTmpl.Copy(),
		dbConnectionsUtilizationChartTmpl.Copy(),
		dbConnectionsChartTmpl.Copy(),
		dbBufferCacheHitRatioChartTmpl.Copy(),
		dbBlocksReadChartTmpl.Copy(),
		dbRowsReadRatioChartTmpl.Copy(),
		dbRowsReadChartTmpl.Copy(),
		dbRowsWrittenChartTmpl.Copy(),
		dbConflictsChartTmpl.Copy(),
		dbConflictsStatChartTmpl.Copy(),
		dbDeadlocksChartTmpl.Copy(),
		dbLocksHeldChartTmpl.Copy(),
		dbLocksAwaitedChartTmpl.Copy(),
		dbTempFilesChartTmpl.Copy(),
		dbTempFilesDataChartTmpl.Copy(),
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
	dbTransactionsChartTmpl = module.Chart{
		ID:       "db_%s_transactions",
		Title:    "Database transactions",
		Units:    "transactions/s",
		Fam:      "db transactions",
		Ctx:      "postgres.db_transactions",
		Priority: prioDBTransactions,
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
	dbConnectionsChartTmpl = module.Chart{
		ID:       "db_%s_connections",
		Title:    "Database connections",
		Units:    "connections",
		Fam:      "db connections",
		Ctx:      "postgres.db_connections",
		Priority: prioDBConnections,
		Dims: module.Dims{
			{ID: "db_%s_numbackends", Name: "connections"},
		},
	}
	dbBufferCacheHitRatioChartTmpl = module.Chart{
		ID:       "db_%s_buffer_cache_hit_ratio",
		Title:    "Database buffer cache hit ratio",
		Units:    "percentage",
		Fam:      "db buffer cache",
		Ctx:      "postgres.db_buffer_cache_hit_ratio",
		Priority: prioDBBufferCacheRatio,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "db_%s_blks_hit", Name: "hit", Algo: module.PercentOfIncremental},
			{ID: "db_%s_blks_read", Name: "miss", Algo: module.PercentOfIncremental},
		},
	}
	dbBlocksReadChartTmpl = module.Chart{
		ID:       "db_%s_blocks_read",
		Title:    "Database blocks read",
		Units:    "blocks/s",
		Fam:      "db buffer cache",
		Ctx:      "postgres.db_blocks_read",
		Priority: prioDBBlockReads,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "db_%s_blks_hit", Name: "memory", Algo: module.Incremental},
			{ID: "db_%s_blks_read", Name: "disk", Algo: module.Incremental},
		},
	}
	dbRowsReadRatioChartTmpl = module.Chart{
		ID:       "db_%s_rows_read_ratio",
		Title:    "Database rows read ratio",
		Units:    "percentage",
		Fam:      "db read throughput",
		Ctx:      "postgres.db_rows_read_ratio",
		Priority: prioDBRowsReadRatio,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "db_%s_tup_returned", Name: "returned", Algo: module.PercentOfIncremental},
			{ID: "db_%s_tup_fetched", Name: "fetched", Algo: module.PercentOfIncremental},
		},
	}
	dbRowsReadChartTmpl = module.Chart{
		ID:       "db_%s_rows_read",
		Title:    "Database rows read",
		Units:    "rows/s",
		Fam:      "db read throughput",
		Ctx:      "postgres.db_rows_read",
		Priority: prioDBRowsRead,
		Dims: module.Dims{
			{ID: "db_%s_tup_returned", Name: "returned", Algo: module.Incremental},
			{ID: "db_%s_tup_fetched", Name: "fetched", Algo: module.Incremental},
		},
	}
	dbRowsWrittenChartTmpl = module.Chart{
		ID:       "db_%s_rows_written",
		Title:    "Database rows written",
		Units:    "rows/s",
		Fam:      "db write throughput",
		Ctx:      "postgres.db_rows_written",
		Priority: prioDBRowsWritten,
		Dims: module.Dims{
			{ID: "db_%s_tup_inserted", Name: "inserted", Algo: module.Incremental},
			{ID: "db_%s_tup_deleted", Name: "deleted", Algo: module.Incremental},
			{ID: "db_%s_tup_updated", Name: "updated", Algo: module.Incremental},
		},
	}
	dbConflictsChartTmpl = module.Chart{
		ID:       "db_%s_conflicts",
		Title:    "Database canceled queries",
		Units:    "queries/s",
		Fam:      "db query cancels",
		Ctx:      "postgres.db_conflicts",
		Priority: prioDBConflicts,
		Dims: module.Dims{
			{ID: "db_%s_conflicts", Name: "conflicts", Algo: module.Incremental},
		},
	}
	dbConflictsStatChartTmpl = module.Chart{
		ID:       "db_%s_conflicts_stat",
		Title:    "Database canceled queries by reason",
		Units:    "queries/s",
		Fam:      "db query cancels",
		Ctx:      "postgres.db_conflicts_stat",
		Priority: prioDBConflictsStat,
		Dims: module.Dims{
			{ID: "db_%s_confl_tablespace", Name: "tablespace", Algo: module.Incremental},
			{ID: "db_%s_confl_lock", Name: "lock", Algo: module.Incremental},
			{ID: "db_%s_confl_snapshot", Name: "snapshot", Algo: module.Incremental},
			{ID: "db_%s_confl_bufferpin", Name: "bufferpin", Algo: module.Incremental},
			{ID: "db_%s_confl_deadlock", Name: "deadlock", Algo: module.Incremental},
		},
	}
	dbDeadlocksChartTmpl = module.Chart{
		ID:       "db_%s_deadlocks",
		Title:    "Database deadlocks",
		Units:    "deadlocks/s",
		Fam:      "db deadlocks",
		Ctx:      "postgres.db_deadlocks",
		Priority: prioDBDeadlocks,
		Dims: module.Dims{
			{ID: "db_%s_deadlocks", Name: "deadlocks", Algo: module.Incremental},
		},
	}
	dbLocksHeldChartTmpl = module.Chart{
		ID:       "db_%s_locks_held",
		Title:    "Database locks held",
		Units:    "locks",
		Fam:      "db locks",
		Ctx:      "postgres.db_locks_held",
		Priority: prioDBLocksHeld,
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
	dbLocksAwaitedChartTmpl = module.Chart{
		ID:       "db_%s_locks_awaited",
		Title:    "Database locks awaited",
		Units:    "locks",
		Fam:      "db locks",
		Ctx:      "postgres.db_locks_awaited",
		Priority: prioDBLocksAwaited,
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
	dbTempFilesChartTmpl = module.Chart{
		ID:       "db_%s_temp_files",
		Title:    "Database temporary files written to disk",
		Units:    "files/s",
		Fam:      "db temp files",
		Ctx:      "postgres.db_temp_files",
		Priority: prioDBTempFiles,
		Dims: module.Dims{
			{ID: "db_%s_temp_files", Name: "written", Algo: module.Incremental},
		},
	}
	dbTempFilesDataChartTmpl = module.Chart{
		ID:       "db_%s_temp_files_data",
		Title:    "Database temporary files data written to disk",
		Units:    "B/s",
		Fam:      "db temp files",
		Ctx:      "postgres.db_temp_files_data",
		Priority: prioDBTempFilesData,
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
		tableRowsChartTmpl.Copy(),
		tableRowsPercChartTmpl.Copy(),
		tableRowsOperationsChartTmpl.Copy(),
		tableHOTUpdatesChartTmpl.Copy(),
		tableScansChartTmpl.Copy(),
		tableScansRowsChartTmpl.Copy(),
		tableTotalSizeChartTmpl.Copy(),
	}

	tableRowsPercChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_rows_perc",
		Title:    "Table total rows",
		Units:    "%",
		Fam:      "table rows",
		Ctx:      "postgres.table_rows_perc",
		Priority: prioTableRowsTotalPercentage,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_n_live_tup", Name: "live", Algo: module.PercentOfAbsolute},
			{ID: "db_%s_schema_%s_table_%s_n_dead_tup", Name: "dead", Algo: module.PercentOfAbsolute},
		},
	}
	tableRowsChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_rows",
		Title:    "Table total rows",
		Units:    "rows",
		Fam:      "table rows",
		Ctx:      "postgres.table_rows",
		Priority: prioTableRowsTotal,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_n_live_tup", Name: "live"},
			{ID: "db_%s_schema_%s_table_%s_n_dead_tup", Name: "dead"},
		},
	}
	tableRowsOperationsChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_rows_operations",
		Title:    "Table throughput",
		Units:    "rows/s",
		Fam:      "table throughput",
		Ctx:      "postgres.table_rows_operations",
		Priority: prioTableRowsOperations,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_n_tup_ins", Name: "inserted", Algo: module.Incremental},
			{ID: "db_%s_schema_%s_table_%s_n_tup_del", Name: "deleted", Algo: module.Incremental},
			{ID: "db_%s_schema_%s_table_%s_n_tup_upd", Name: "updated", Algo: module.Incremental},
		},
	}
	tableHOTUpdatesChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_hot_updates",
		Title:    "Table HOT updates",
		Units:    "updates/s",
		Fam:      "table hot updates",
		Ctx:      "postgres.table_hot_updates",
		Priority: prioTableHOTUpdates,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_n_tup_hot_upd", Name: "hot", Algo: module.Incremental},
		},
	}
	tableScansChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_scans",
		Title:    "Table scans",
		Units:    "scans/s",
		Fam:      "table scans",
		Ctx:      "postgres.table_scans",
		Priority: prioTableScans,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_idx_scan", Name: "index"},
			{ID: "db_%s_schema_%s_table_%s_seq_scan", Name: "sequential"},
		},
	}
	tableScansRowsChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_scans_rows",
		Title:    "Table live rows fetched by scans",
		Units:    "rows/s",
		Fam:      "table scans",
		Ctx:      "postgres.table_scans_rows",
		Priority: prioTableScansRows,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_idx_tup_fetch", Name: "sequential"},
			{ID: "db_%s_schema_%s_table_%s_seq_tup_read", Name: "index"},
		},
	}
	tableLastAutoVacuumAgoChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_last_autovacuum_ago",
		Title:    "Table time since last auto VACUUM",
		Units:    "seconds",
		Fam:      "table autovacuum",
		Ctx:      "postgres.table_last_autovacuum_ago",
		Priority: prioTableLastAutovacuumAgo,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_last_autovacuum_ago", Name: "time"},
		},
	}
	tableLastVacuumAgoChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_last_vacuum_ago",
		Title:    "Table time since last manual VACUUM",
		Units:    "seconds",
		Fam:      "table vacuum",
		Ctx:      "postgres.table_last_vacuum_ago",
		Priority: prioTableLastVacuumAgo,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_last_vacuum_ago", Name: "time"},
		},
	}
	tableLastAutoAnalyzeAgoChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_last_autoanalyze_ago",
		Title:    "Table time since last auto ANALYZE",
		Units:    "seconds",
		Fam:      "table autoanalyze",
		Ctx:      "postgres.table_last_autoanalyze_ago",
		Priority: prioTableLastAutoAnalyzeAgo,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_last_autoanalyze_ago", Name: "time"},
		},
	}
	tableLastAnalyzeAgoChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_last_analyze_ago",
		Title:    "Table time since last manual ANALYZE",
		Units:    "seconds",
		Fam:      "table analyze",
		Ctx:      "postgres.table_last_analyze_ago",
		Priority: prioTableLastAnalyzeAgo,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_last_analyze_ago", Name: "time"},
		},
	}
	tableTotalSizeChartTmpl = module.Chart{
		ID:       "db_%s_schema_%s_table_%s_total_size",
		Title:    "Table total size",
		Units:    "B",
		Fam:      "table size",
		Ctx:      "postgres.table_total_size",
		Priority: prioTableTotalSize,
		Dims: module.Dims{
			{ID: "db_%s_schema_%s_table_%s_total_size", Name: "size"},
		},
	}
)

func newTableCharts(dbname, schema, name string) *module.Charts {
	charts := tableChartsTmpl.Copy()

	for i, chart := range *charts {
		(*charts)[i] = newTableChart(chart, dbname, schema, name)
	}

	return charts
}

func newTableChart(chart *module.Chart, dbname, schema, name string) *module.Chart {
	chart = chart.Copy()
	chart.ID = fmt.Sprintf(chart.ID, dbname, schema, name)
	chart.Labels = []module.Label{
		{Key: "database", Value: dbname},
		{Key: "schema", Value: schema},
		{Key: "table", Value: name},
	}
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, dbname, schema, name)
	}
	return chart
}

func (p *Postgres) addNewTableCharts(dbname, schema, name string) {
	charts := newTableCharts(dbname, schema, name)
	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableLastAutoVacuumAgoChart(dbname, schema, name string) {
	chart := newTableChart(tableLastAutoVacuumAgoChartTmpl.Copy(), dbname, schema, name)

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableLastVacuumAgoChart(dbname, schema, name string) {
	chart := newTableChart(tableLastVacuumAgoChartTmpl.Copy(), dbname, schema, name)

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableLastAutoAnalyzeAgoChart(dbname, schema, name string) {
	chart := newTableChart(tableLastAutoAnalyzeAgoChartTmpl.Copy(), dbname, schema, name)

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) addTableLastAnalyzeAgoChart(dbname, schema, name string) {
	chart := newTableChart(tableLastAnalyzeAgoChartTmpl.Copy(), dbname, schema, name)

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) removeTableCharts(dbname, schema, name string) {
	prefix := fmt.Sprintf("db_%s_schema_%s_table_%s", dbname, schema, name)
	for _, c := range *p.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}
