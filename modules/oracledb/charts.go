package oracledb

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "processes",
		Title: "Processes",
		Ctx:   "oracledb.processes",
		Units: "processes",
		Dims: Dims{
			{ID: "processes", Name: "processes"},
		},
	},
	{
		ID:    "sessions_total",
		Title: "Total Sessions",
		Ctx:   "oracledb.sessions_total",
		Units: "sessions",
		Dims: Dims{
			{ID: "sessions_total", Name: "total", Algo: module.Incremental},
		},
	},
	{
		ID:    "sessions",
		Title: "Sessions",
		Ctx:   "oracledb.sessions",
		Units: "sessions",
		Dims: Dims{
			{ID: "sessions_active", Name: "active"},
			{ID: "sessions_inactive", Name: "inactive"},
		},
	},
	{
		ID:    "activity",
		Title: "Activity",
		Ctx:   "oracledb.activity",
		Units: "activities",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "activity_parse_count_total", Name: "parse count (total)", Algo: module.Incremental},
			{ID: "activity_execute_count", Name: "execute count", Algo: module.Incremental},
			{ID: "activity_user_commits", Name: "user commits", Algo: module.Incremental},
			{ID: "activity_user_rollbacks", Name: "user rollbacks", Algo: module.Incremental},
		},
	},
	{
		ID:    "wait_time",
		Title: "Wait Time",
		Ctx:   "oracledb.wait_time",
		Units: "ms",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "wait_time_configuration", Name: "configuration", Div: 1000},
			{ID: "wait_time_administrative", Name: "administrative", Div: 1000},
			{ID: "wait_time_system_io", Name: "system I/O", Div: 1000},
			{ID: "wait_time_application", Name: "application", Div: 1000},
			{ID: "wait_time_concurrency", Name: "concurrency", Div: 1000},
			{ID: "wait_time_commit", Name: "commit", Div: 1000},
			{ID: "wait_time_network", Name: "network", Div: 1000},
			{ID: "wait_time_user_io", Name: "user I/O", Div: 1000},
			{ID: "wait_time_other", Name: "other", Div: 1000},
		},
	},
	/*
			{
				ID:    "tablespace",
				Title: "Tablespace Size",
				Ctx:   "oracledb.tablespace",
				Units: "KiB",
				Type:  module.Stacked,
				Dims: Dims{
					{ID: "tablespace_max_bytes_system", Name: "system: max bytes", Div: 1024 * 100000},
					{ID: "tablespace_max_bytes_sysaux", Name: "sysaux: max bytes", Div: 1024 * 100000},
					{ID: "tablespace_max_bytes_users", Name: "users: max bytes", Div: 1024 * 100000},
					{ID: "tablespace_max_bytes_temp", Name: "temp: max bytes", Div: 1024 * 100000},

					{ID: "tablespace_free_bytes_system", Name: "system: free bytes", Div: 1024 * 100000},
					{ID: "tablespace_free_bytes_sysaux", Name: "sysaux: free bytes", Div: 1024 * 100000},
					{ID: "tablespace_free_bytes_users", Name: "users: free bytes", Div: 1024 * 100000},
					{ID: "tablespace_free_bytes_temp", Name: "temp: free bytes", Div: 1024 * 100000},

					{ID: "tablespace_bytes_system", Name: "system: bytes", Div: 1024 * 100000},
					{ID: "tablespace_bytes_sysaux", Name: "sysaux: bytes", Div: 1024 * 100000},
					{ID: "tablespace_bytes_users", Name: "users: bytes", Div: 1024 * 100000},
					{ID: "tablespace_bytes_temp", Name: "temp: bytes", Div: 1024 * 100000},
				},
			},
			{
				ID:    "system",
				Title: "System Metrics",
				Ctx:   "oracledb.system",
				Type:  module.Stacked,
				Units: "metrics",
				Dims: Dims{
					{ID: "system_buffer_cachehit_ratio", Name: "system_buffer_cachehit_ratio", Algo: module.PercentOfAbsolute},
					{ID: "system_cursor_cachehit_ratio", Name: "system_cursor_cachehit_ratio", Algo: module.PercentOfAbsolute},
					{ID: "system_library_cachehit_ratio", Name: "system_library_cachehit_ratio", Algo: module.PercentOfAbsolute},
					{ID: "system_shared_pool_free", Name: "system_shared_pool_free"},
					{ID: "system_physical_reads", Name: "system_physical_writes", Algo: module.Incremental},
					{ID: "system_enqueue_timeouts", Name: "system_enqueue_timeouts"},
					{ID: "system_gc_cr_block_received", Name: "system_gc_cr_block_received"},
					{ID: "system_cache_blocks_corrupt", Name: "system_cache_blocks_corrupt"},
					{ID: "system_cache_blocks_lost", Name: "system_cache_blocks_lost"},
					{ID: "system_logons", Name: "system_logons", Algo: module.Incremental},
					{ID: "system_active_sessions", Name: "system_active_sessions"},
					{ID: "system_long_table_scans", Name: "system_long_table_scans", Algo: module.Incremental},
					{ID: "system_service_response_time", Name: "system_service_response_time"},
					{ID: "system_user_rollbacks", Name: "system_user_rollbacks", Algo: module.Incremental},
					{ID: "system_sorts_per_user_call", Name: "system_sorts_per_user_call"},
					{ID: "system_rows_per_sort", Name: "system_rows_per_sort"},
					{ID: "system_disk_sorts", Name: "system_disk_sorts", Algo: module.Incremental},
					{ID: "system_memory_sorts_ratio", Name: "system_memory_sorts_ratio", Algo: module.PercentOfAbsolute},
					{ID: "system_database_wait_time_ratio", Name: "system_database_wait_time_ratio", Algo: module.PercentOfAbsolute},
					{ID: "system_session_limit_usage", Name: "system_session_limit_usage"},
					{ID: "system_session_count", Name: "system_session_count", Algo: module.Incremental},
					{ID: "system_temp_space_used", Name: "system_temp_space_used"},
				},
		 },
	*/
}
