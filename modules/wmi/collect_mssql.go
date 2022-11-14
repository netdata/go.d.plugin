// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"math"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricMSSQLAccessMethodPageSplits       = "windows_mssql_accessmethods_page_splits"
	metricMSSQLBufferCacheHits              = "windows_mssql_bufman_buffer_cache_hits"
	metricMSSQLBufferCacheLookups           = "windows_mssql_bufman_buffer_cache_lookups"
	metricMSSQLBufferCheckpointPages        = "windows_mssql_bufman_checkpoint_pages"
	metricMSSQLBufferPageLifeExpectancy     = "windows_mssql_bufman_page_life_expectancy_seconds"
	metricMSSQLBufferPageRead               = "windows_mssql_bufman_page_reads"
	metricMSSQLBufferPageWrite              = "windows_mssql_bufman_page_writes"
	metricMSSQLActiveTransactions           = "windows_mssql_databases_active_transactions"
	metricMSSQLBackupRestoreOperation       = "windows_mssql_databases_backup_restore_operations"
	metricMSSQLDataFileSize                 = "windows_mssql_databases_data_files_size_bytes"
	metricMSSQLLogFlushed                   = "windows_mssql_databases_log_flushed_bytes"
	metricMSSQLLogFlushes                   = "windows_mssql_databases_log_flushes"
	metricMSSQLTransactions                 = "windows_mssql_databases_transactions"
	metricMSSQLWriteTransaction             = "windows_mssql_databases_write_transactions"
	metricMSSQLBlockedProcesses             = "windows_mssql_genstats_blocked_processes"
	metricMSSQLUserConnections              = "windows_mssql_genstats_user_connections"
	metricMSSQLLockWait                     = "windows_mssql_locks_lock_wait_seconds"
	metricMSSQLPendingMemoryGrant           = "windows_mssql_memmgr_pending_memory_grants"
	metricMSSQLTotalServerMemory            = "windows_mssql_memmgr_total_server_memory_bytes"
	metricMSSQLStatsAutoParameterization    = "windows_mssql_sqlstats_auto_parameterization_attempts"
	metricMSSQLStatSafeAutoParameterization = "windows_mssql_sqlstats_safe_auto_parameterization_attempts"
	metricMSSQLCompilation                  = "windows_mssql_sqlstats_sql_compilations"
	metricMSSQLRecompilation                = "windows_mssql_sqlstats_sql_recompilations"
)

func (w *WMI) collectMSSQL(mx map[string]int64, pms prometheus.Metrics) {
	seen := make(map[string]bool)
	db := make(map[string]bool)
	px := "mssql_instance_"
	for _, pm := range pms.FindByName(metricMSSQLAccessMethodPageSplits) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_access_page_split"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLBufferCacheHits) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_cache_hit_ratio"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLBufferCacheLookups) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			var hit float64
			if pm.Value != 0 {
				hit = float64(mx[px+name+"_cache_hit_ratio"]) / pm.Value
			} else {
				hit = math.NaN()
			}
			mx[px+name+"_cache_hit_ratio"] = int64(hit)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLBufferCheckpointPages) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_buffer_checkpoint_page"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLBufferPageLifeExpectancy) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_buffer_pagelife_expectancy"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLBufferPageRead) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_page_read"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLBufferPageWrite) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_page_write"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLActiveTransactions) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			if dbname := cleanInstanceDBName(pm.Labels.Get("database")); dbname != "" {
				seen[name] = true
				db[dbname] = true
				mx[px+name+"_"+dbname+"_active_transaction"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLBackupRestoreOperation) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			if dbname := cleanInstanceDBName(pm.Labels.Get("database")); dbname != "" {
				seen[name] = true
				db[dbname] = true
				mx[px+name+"_"+dbname+"_backup_restore"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLDataFileSize) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			if dbname := cleanInstanceDBName(pm.Labels.Get("database")); dbname != "" {
				seen[name] = true
				db[dbname] = true
				mx[px+name+"_"+dbname+"_database_size"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLLogFlushed) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			if dbname := cleanInstanceDBName(pm.Labels.Get("database")); dbname != "" {
				seen[name] = true
				db[dbname] = true
				mx[px+name+"_"+dbname+"_log_flushed"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLLogFlushes) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			if dbname := cleanInstanceDBName(pm.Labels.Get("database")); dbname != "" {
				seen[name] = true
				db[dbname] = true
				mx[px+name+"_"+dbname+"_log_flushes"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLTransactions) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			if dbname := cleanInstanceDBName(pm.Labels.Get("database")); dbname != "" {
				seen[name] = true
				db[dbname] = true
				mx[px+name+"_"+dbname+"_transaction"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLWriteTransaction) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			if dbname := cleanInstanceDBName(pm.Labels.Get("database")); dbname != "" {
				seen[name] = true
				db[dbname] = true
				mx[px+name+"_"+dbname+"_write_transaction"] = int64(pm.Value)
			}
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLBlockedProcesses) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_blocked_process"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLUserConnections) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_user_connection"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLLockWait) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			resource := pm.Labels.Get("resource")
			idx := buildLockWaitIndex(px, name, resource)
			mx[idx] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLPendingMemoryGrant) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_memmgr_pending"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLTotalServerMemory) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_memmgr_total"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLStatsAutoParameterization) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_state_auto_param"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLStatSafeAutoParameterization) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_state_safe_auto"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLCompilation) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_state_compilation"] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricMSSQLRecompilation) {
		if name := cleanInstanceDBName(pm.Labels.Get("mssql_instance")); name != "" {
			seen[name] = true
			mx[px+name+"_state_recompilation"] = int64(pm.Value)
		}
	}

	for instance := range seen {
		if !w.cache.mssql[instance] {
			w.cache.mssql[instance] = true
			w.addMSSQLInstanceCharts(instance, &db)
		}
	}
	for instance := range w.cache.mssql {
		if !seen[instance] {
			delete(w.cache.mssql, instance)
			w.removeMSSQLInstanceCharts(instance)
		}
	}
}

func buildLockWaitIndex(prefix string, instance string, selector string) string {
	var sufix string
	switch selector {
	case "AllocUnit":
		sufix = "allocunit"
	}

	return prefix + "_" + instance + "_lock_wait_" + sufix
}

func cleanInstanceDBName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}
