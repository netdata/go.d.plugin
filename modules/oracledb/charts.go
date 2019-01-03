package oracledb

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{
		ID:    "processes",
		Title: "Processes",
		Ctx:   "oracledb.processes",
		Dims: Dims{
			{ID: "processes", Name: "processes"},
		},
	},
	{
		ID:    "sessions",
		Title: "Total Sessions",
		Ctx:   "oracledb.sessions",
		Dims: Dims{
			{ID: "sessions_total", Name: "sessions_total"},
			{ID: "sessions_active", Name: "sessions_active"},
			{ID: "sessions_inactive", Name: "sessions_inactive"},
		},
	},
	{
		ID:    "Activity",
		Title: "Activity",
		Ctx:   "oracledb.activity",
		Dims: Dims{
			{ID: "activity_parse_count_total", Name: "activity_parse_count_total"},
			{ID: "activity_execute_count", Name: "activity_execute_count"},
			{ID: "activity_user_commits", Name: "activity_user_commits"},
			{ID: "activity_user_rollbacks", Name: "activity_user_rollbacks"},
		},
	},
	{
		ID:    "wait_time",
		Title: "Wait Time",
		Ctx:   "oracledb.wait_time",
		Dims: Dims{
			{ID: "wait_time_configuration", Name: "wait_time_configuration", Div: 1000},
			{ID: "wait_time_administrative", Name: "wait_time_administrative", Div: 1000},
			{ID: "wait_time_system_io", Name: "wait_time_system_io", Div: 1000},
			{ID: "wait_time_application", Name: "wait_time_application", Div: 1000},
			{ID: "wait_time_concurrency", Name: "wait_time_concurrency", Div: 1000},
			{ID: "wait_time_commit", Name: "wait_time_commit", Div: 1000},
			{ID: "wait_time_network", Name: "wait_time_network", Div: 1000},
			{ID: "wait_time_user_io", Name: "wait_time_user_io", Div: 1000},
			{ID: "wait_time_other", Name: "wait_time_other", Div: 1000},
		},
	},
	{
		ID:    "tablespace",
		Title: "Tablespace Size",
		Ctx:   "oracledb.tablespace",
		Dims: Dims{
			{ID: "tablespace_max_bytes_system", Name: "tablespace_max_bytes_system"},
			{ID: "tablespace_max_bytes_sysaux", Name: "tablespace_max_bytes_sysaux"},
			{ID: "tablespace_max_bytes_users", Name: "tablespace_max_bytes_users"},
			{ID: "tablespace_max_bytes_temp", Name: "tablespace_max_bytes_temp"},

			{ID: "tablespace_free_bytes_system", Name: "tablespace_free_bytes_system"},
			{ID: "tablespace_free_bytes_sysaux", Name: "tablespace_free_bytes_sysaux"},
			{ID: "tablespace_free_bytes_users", Name: "tablespace_free_bytes_users"},
			{ID: "tablespace_free_bytes_temp", Name: "tablespace_free_bytes_temp"},

			{ID: "tablespace_bytes_system", Name: "tablespace_bytes_system"},
			{ID: "tablespace_bytes_sysaux", Name: "tablespace_bytes_sysaux"},
			{ID: "tablespace_bytes_users", Name: "tablespace_bytes_users"},
			{ID: "tablespace_bytes_temp", Name: "tablespace_bytes_temp"},
		},
	},
	{
		ID:    "system",
		Title: "System Metrics",
		Ctx:   "oracledb.system",
		Dims: Dims{
			{ID: "system_buffer_cachehit_ratio", Name: "system_buffer_cachehit_ratio"},
			{ID: "system_cursor_cachehit_ratio", Name: "system_cursor_cachehit_ratio"},
			{ID: "system_library_cachehit_ratio", Name: "system_library_cachehit_ratio"},
			{ID: "system_shared_pool_free", Name: "system_shared_pool_free"},
			{ID: "system_physical_reads", Name: "system_physical_writes"},
			{ID: "system_enqueue_timeouts", Name: "system_enqueue_timeouts"},
			{ID: "system_gc_cr_block_received", Name: "system_gc_cr_block_received"},
			{ID: "system_cache_blocks_corrupt", Name: "system_cache_blocks_corrupt"},
			{ID: "system_cache_blocks_lost", Name: "system_cache_blocks_lost"},
			{ID: "system_logons", Name: "system_logons"},
			{ID: "system_active_sessions", Name: "system_active_sessions"},
			{ID: "system_long_table_scans", Name: "system_long_table_scans"},
			{ID: "system_service_response_time", Name: "system_service_response_time"},
			{ID: "system_user_rollbacks", Name: "system_user_rollbacks"},
			{ID: "system_sorts_per_user_call", Name: "system_sorts_per_user_call"},
			{ID: "system_rows_per_sort", Name: "system_rows_per_sort"},
			{ID: "system_disk_sorts", Name: "system_disk_sorts"},
			{ID: "system_memory_sorts_ratio", Name: "system_memory_sorts_ratio"},
			{ID: "system_database_wait_time_ratio", Name: "system_database_wait_time_ratio"},
			{ID: "system_session_limit_usage", Name: "system_session_limit_usage"},
			{ID: "system_session_count", Name: "system_session_count"},
			{ID: "system_temp_space_used", Name: "system_temp_space_used"},
		},
	},
}
