// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricMSSQLAccessMethodPageSplits = "windows_mssql_accessmethods_page_splits"
	metricMSSQLBufferCacheHits = "windows_mssql_bufman_buffer_cache_hits"
	metricMSSQLBufferCacheLookups = "windows_mssql_bufman_buffer_cache_lookups"
	metricMSSQLBufferCheckpointPages = "windows_mssql_bufman_checkpoint_pages" 
	metricMSSQLBufferPageLifeExpectancy = "windows_mssql_bufman_page_life_expectancy_seconds" 
	metricMSSQLBufferPageRead = "windows_mssql_bufman_page_reads" 
	metricMSSQLBufferPageWrite = "windows_mssql_bufman_page_writes" 
	metricMSSQLActiveTransactions = "windows_mssql_databases_active_transactions"
	metricMSSQLBackupRestoreOperation = "windows_mssql_databases_backup_restore_operations"
	metricMSSQLDataFileSize = "windows_mssql_databases_data_files_size_bytes"
	metricMSSQLLogFlushed = "windows_mssql_databases_log_flushed_bytes"
	metricMSSQLLogFlushes = "windows_mssql_databases_log_flushes" 
	metricMSSQLTransactions = "windows_mssql_databases_transactions"
	metricMSSQLWriteTransaction = "windows_mssql_databases_write_transactions" 
	metricMSSQLBlockedProcesses = "windows_mssql_genstats_blocked_processes"
	metricMSSQLUserConnections = "windows_mssql_genstats_user_connections"
	metricMSSQLLockWait = "windows_mssql_locks_lock_wait_seconds" 
	metricMSSQLPendingMemoryGrant = "windows_mssql_memmgr_pending_memory_grants"
	metricMSSQLTotalServerMemory = "windows_mssql_memmgr_total_server_memory_bytes"
	metricMSSQLStatsAutoParameterization = "windows_mssql_sqlstats_auto_parameterization_attempts"
	metricMSSQLStatSafeAutoParameterization = "windows_mssql_sqlstats_safe_auto_parameterization_attempts"
	metricMSSQLCompilation = "windows_mssql_sqlstats_sql_compilations"
	metricMSSQLRecompilation = "windows_mssql_sqlstats_sql_recompilations"
)

func (w *WMI) collectMSSQL(mx map[string]int64, pms prometheus.Metrics) {
	seen := make(map[string]bool)
	px := "mssql_instance_"
	for _, pm := range pms.FindByName(metricMSSQLAccessMethodPageSplits) {
		if name := cleanWebsiteName(pm.Labels.Get("mssql_instance")); name != "" {
				seen[name] = true
				mx[px+name+"_access_page_split"] = int64(pm.Value)
		}
	}
}
