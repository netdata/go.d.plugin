package proxysql

import (
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	Charts = module.Charts
	Chart  = module.Chart
	Dims   = module.Dims
)

var charts = Charts{
	{
		ID:    "proxysql_uptime",
		Title: "Uptime",
		Units: "seconds",
		Fam:   "global_stats",
		Ctx:   "proxysql.uptime",
		Type:  module.Line,
		Dims: Dims{
			{ID: "proxysql_uptime", Name: "uptime", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "questions",
		Title: "Questions",
		Units: "questions",
		Fam:   "global_stats",
		Ctx:   "proxysql.questions",
		Type:  module.Line,
		Dims: Dims{
			{ID: "questions", Name: "questions", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "active_transactions",
		Title: "Active transactions",
		Units: "transactions",
		Fam:   "global_stats",
		Ctx:   "proxysql.active_transactions",
		Type:  module.Line,
		Dims: Dims{
			{ID: "active_transactions", Name: "active transactions", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "slow_queries",
		Title: "Slow queries",
		Units: "queries",
		Fam:   "global_stats",
		Ctx:   "proxysql.slow_queries",
		Type:  module.Line,
		Dims: Dims{
			{ID: "slow_queries", Name: "slow queries", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "backend_lagging",
		Title: "Backend lagging during query",
		Units: "backends",
		Fam:   "global_stats",
		Ctx:   "proxysql.backend_lagging_during_query",
		Type:  module.Line,
		Dims: Dims{
			{ID: "backend_lagging_during_query", Name: "backends", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "backend_offline",
		Title: "Backend offline during query",
		Units: "backends",
		Fam:   "global_stats",
		Ctx:   "proxysql.backend_offline_during_query",
		Type:  module.Line,
		Dims: Dims{
			{ID: "backend_offline_during_query", Name: "backends", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "generated_error_packets",
		Title: "Generated error packets",
		Units: "packets",
		Fam:   "global_stats",
		Ctx:   "proxysql.generated_error_packets",
		Type:  module.Line,
		Dims: Dims{
			{ID: "generated_error_packets", Name: "packets", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "max_connect_timeouts",
		Title: "Max connection timeouts",
		Units: "connections",
		Fam:   "global_stats",
		Ctx:   "proxysql.max_connect_timeouts",
		Type:  module.Line,
		Dims: Dims{
			{ID: "max_connect_timeouts", Name: "connections", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "client_connections",
		Title: "Client connections",
		Units: "connections",
		Fam:   "global_stats",
		Ctx:   "proxysql.client_connections",
		Type:  module.Line,
		Dims: Dims{
			{ID: "client_connections_aborted", Name: "aborted", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "client_connections_connected", Name: "connected", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "client_connections_created", Name: "created", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "client_connections_hostgroup_locked", Name: "hostgroup locked", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "client_connections_non_idle", Name: "non idle", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "server_connections",
		Title: "Server connections",
		Units: "connections",
		Fam:   "global_stats",
		Ctx:   "proxysql.server_connections",
		Type:  module.Line,
		Dims: Dims{
			{ID: "server_connections_aborted", Name: "aborted", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "server_connections_connected", Name: "connected", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "server_connections_created", Name: "created", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "server_connections_delayed", Name: "delayed", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "access_denied",
		Title: "Access denied",
		Units: "number",
		Fam:   "global_stats",
		Ctx:   "proxysql.access_denied",
		Type:  module.Line,
		Dims: Dims{
			{ID: "access_denied_max_connections", Name: "max connections", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "access_denied_max_user_connections", Name: "user connections", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "access_denied_wrong_password", Name: "wrong password", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "query_time",
		Title: "Query time",
		Units: "nanoseconds",
		Fam:   "global_stats",
		Ctx:   "proxysql.query_time",
		Type:  module.Line,
		Dims: Dims{
			{ID: "backend_query_time", Name: "backend query time", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "query_processor_time_nsec", Name: "query processor time", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "commands",
		Title: "Commands",
		Units: "commands",
		Fam:   "global_stats",
		Ctx:   "proxysql.commands",
		Type:  module.Line,
		Dims: Dims{
			{ID: "com_autocommit", Name: "autocommit", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_autocommit_filtered", Name: "autocommit filterred", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_backend_change_user", Name: "backend change user", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_backend_init_db", Name: "backend init db", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_backend_set_names", Name: "backend set names", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_backend_stmt_close", Name: "backend stmt close", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_backend_stmt_execute", Name: "backend stmt execute", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_backend_stmt_prepare", Name: "backend stmt prepare", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_commit", Name: "commit", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_commit_filtered", Name: "commit filtered", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_frontend_init_db", Name: "frontend init db", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_frontend_set_names", Name: "frontend set names", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_frontend_stmt_close", Name: "frontend stmt close", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_frontend_stmt_execute", Name: "frontend stmt execute", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_frontend_stmt_prepare", Name: "frontend stmt prepare", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_frontend_use_db", Name: "frontend use db", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_rollback", Name: "rollback", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "com_rollback_filtered", Name: "rollback filtered", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "connection_pool_requests",
		Title: "Connection  pool requests",
		Units: "connections",
		Fam:   "global_stats",
		Ctx:   "proxysql.connection_pool_requests",
		Type:  module.Line,
		Dims: Dims{
			{ID: "connpool_get_conn_failure", Name: "failure", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "connpool_get_conn_immediate", Name: "immediate", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "connpool_get_conn_latency_awareness", Name: "latency awareness", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "connpool_get_conn_success", Name: "success", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "myhgm_myconnpoll_get", Name: "get", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "myhgm_myconnpoll_get_ok", Name: "get ok", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "connection_pool_connections",
		Title: "Connection  pool connections",
		Units: "connections",
		Fam:   "global_stats",
		Ctx:   "proxysql.connection_pool_connections",
		Type:  module.Line,
		Dims: Dims{
			{ID: "myhgm_myconnpoll_destroy", Name: "destroy", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "myhgm_myconnpoll_push", Name: "push", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "myhgm_myconnpoll_reset", Name: "reset", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "mysql_monitor_threads",
		Title: "Mysql monitor threads",
		Units: "threads",
		Fam:   "global_stats",
		Ctx:   "proxysql.mysql_monitor_threads",
		Type:  module.Line,
		Dims: Dims{
			{ID: "mysql_monitor_workers", Name: "workers", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_workers_aux", Name: "workers aux", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_workers_started", Name: "workers started", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_connect_check_err", Name: "connect check err", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_connect_check_ok", Name: "connect check ok", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_ping_check_err", Name: "ping check err", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_ping_check_ok", Name: "ping check ok", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_read_only_check_ERR", Name: "read only check err", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_read_only_check_ok", Name: "read only check ok", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_replication_lag_check_err", Name: "replication lag check err", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_monitor_replication_lag_check_ok", Name: "replication lag check ok", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "mysql_thread_workers",
		Title: "Mysql thread workers",
		Units: "workers",
		Fam:   "global_stats",
		Ctx:   "proxysql.mysql_thread_workers",
		Type:  module.Line,
		Dims: Dims{
			{ID: "mysql_thread_workers", Name: "workers", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "network",
		Title: "Network",
		Units: "bytes",
		Fam:   "global_stats",
		Ctx:   "proxysql.network",
		Type:  module.Line,
		Dims: Dims{
			{ID: "queries_backends_bytes_recv", Name: "query backends received", Algo: module.Absolute, Mul: 8, Div: 1000},
			{ID: "queries_backends_bytes_sent", Name: "query backends sent", Algo: module.Absolute, Mul: -8, Div: 1000},
			{ID: "queries_frontends_bytes_recv", Name: "query frontends received", Algo: module.Absolute, Mul: 8, Div: 1000},
			{ID: "queries_frontends_bytes_sent", Name: "query backends received", Algo: module.Absolute, Mul: -8, Div: 1000},
			{ID: "query_cache_bytes_in", Name: "query cache in", Algo: module.Absolute, Mul: 8, Div: 1000},
			{ID: "query_cache_bytes_out", Name: "query cache out", Algo: module.Absolute, Mul: -8, Div: 1000},
			{ID: "mysql_backend_buffers_bytes", Name: "mysql backend buffers", Algo: module.Absolute, Mul: -8, Div: 1000},
			{ID: "mysql_frontend_buffers_bytes", Name: "mysql frontend buffers", Algo: module.Absolute, Mul: -8, Div: 1000},
		},
	},
	{
		ID:    "query_cache",
		Title: "Query cache",
		Units: "number of entries",
		Fam:   "global_stats",
		Ctx:   "proxysql.query_cache",
		Type:  module.Line,
		Dims: Dims{
			{ID: "query_cache_entries", Name: "entries", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "query_cache_purged", Name: "purged", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "query_cache_count_get", Name: "count get", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "query_cache_count_get_ok", Name: "count get ok", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "query_cache_count_set", Name: "count set", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "servers_table_version",
		Title: "Servers table version",
		Units: "version",
		Fam:   "global_stats",
		Ctx:   "proxysql.servers_table_version",
		Type:  module.Line,
		Dims: Dims{
			{ID: "servers_table_version", Name: "version", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "prepared_statements",
		Title: "Prepared statements",
		Units: "prepared statements",
		Fam:   "global_stats",
		Ctx:   "proxysql.prepared_statements",
		Type:  module.Line,
		Dims: Dims{
			{ID: "stmt_cached", Name: "cached", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "stmt_client_active_total", Name: "client active total", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "stmt_client_active_unique", Name: "client active unique", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "stmt_max_stmt_id", Name: "max statement id", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "stmt_server_active_total", Name: "server active total", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "stmt_server_active_unique", Name: "server active unique", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "mysql_max_allowed_packet",
		Title: "MySQL max allowed packet",
		Units: "bytes",
		Fam:   "global_vars",
		Ctx:   "proxysql.mysql_max_allowed_packet",
		Type:  module.Line,
		Dims: Dims{
			{ID: "mysql_max_allowed_packet", Name: "mysql max allowed packets", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "memory",
		Title: "Memory",
		Units: "bytes",
		Fam:   "memory",
		Ctx:   "proxysql.memory",
		Type:  module.Line,
		Dims: Dims{
			{ID: "auth_memory", Name: "auth", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "sqlite3_memory_bytes", Name: "sqlite3", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "query_digest_memory", Name: "query digest", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "connpool_memory_bytes", Name: "connection pool", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "query_cache_memory_bytes", Name: "query cache", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_session_internal_bytes", Name: "mysql session internal", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "jemalloc_memory",
		Title: "jemalloc memory",
		Units: "bytes",
		Fam:   "memory",
		Ctx:   "proxysql.jemalloc_memory",
		Type:  module.Line,
		Dims: Dims{
			{ID: "jemalloc_active", Name: "active", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "jemalloc_allocated", Name: "allocated", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "jemalloc_mapped", Name: "mapped", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "jemalloc_metadata", Name: "metadata", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "jemalloc_resident", Name: "resident", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "jemalloc_retained", Name: "retained", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "mysql",
		Title: "MySQL",
		Units: "bytes",
		Fam:   "memory",
		Ctx:   "proxysql.mysql_memory",
		Type:  module.Line,
		Dims: Dims{
			{ID: "mysql_firewall_rules_config", Name: "firewall_rules_config", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_firewall_rules_table", Name: "firewall_rules_table", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_firewall_users_config", Name: "firewall_users_config", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_firewall_users_table", Name: "firewall_users_table", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "mysql_query_rules_memory", Name: "query_rules_memory", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
	{
		ID:    "stack_memory",
		Title: "Stack memory",
		Units: "bytes",
		Fam:   "memory",
		Ctx:   "proxysql.stack_memory",
		Type:  module.Line,
		Dims: Dims{
			{ID: "stack_memory_admin_threads", Name: "admin_threads", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "stack_memory_cluster_threads", Name: "cluster_threads", Algo: module.Absolute, Mul: 1, Div: 1},
			{ID: "stack_memory_mysql_threads", Name: "mysql_threads", Algo: module.Absolute, Mul: 1, Div: 1},
		},
	},
}

func newMysqlCommandCountersCharts(command string) module.Charts {
	command = strings.ToLower(command)
	return module.Charts{
		{
			ID:    "mysql_command_counts_" + command,
			Title: "MySQL command counts",
			Units: "commands",
			Fam:   "mysql command " + command,
			Ctx:   "proxysql.mysql_command_counts",
			Type:  module.Line,
			Dims: Dims{
				{ID: "mysql_command_" + command + "_total_cnt", Name: "total commands", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_100us", Name: "less than 100us", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_500us", Name: "less than 500us", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_1ms", Name: "less than 1ms", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_5ms", Name: "less than 5ms", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_10ms", Name: "less than 10ms", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_50ms", Name: "less than 50ms", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_100ms", Name: "less than 100ms", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_500ms", Name: "less than 500ms", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_1s", Name: "less than 1s", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_5s", Name: "less than 5s", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_10s", Name: "less than 10s", Algo: module.Incremental},
				{ID: "mysql_command_" + command + "_cnt_INFs", Name: "less than infinity", Algo: module.Incremental},
			},
		},
		{
			ID:    "mysql_command_time_" + command,
			Title: "Duration",
			Units: "microseconds",
			Fam:   "mysql command " + command,
			Ctx:   "proxysql.mysql_command_time",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "mysql_command_" + command + "_total_time_us", Name: "Total", Algo: module.Incremental},
			},
		},
	}
}

func newMysqlUsersCharts(username string) module.Charts {
	return module.Charts{
		{
			ID:    "mysql_users_" + username,
			Title: "MySQL users",
			Units: "connections",
			Fam:   "mysql user " + username,
			Ctx:   "proxysql.mysql_users",
			Type:  module.Line,
			Dims: Dims{
				{ID: "mysql_user_" + username + "_frontend_connections", Name: "frontend connections", Algo: module.Absolute},
				{ID: "mysql_user_" + username + "_frontend_max_connections", Name: "frontend max connections", Algo: module.Absolute},
			},
		},
	}
}

func (p *ProxySQL) addMysqlCommandCountersCharts(command string) {
	if err := p.Charts().Add(newMysqlCommandCountersCharts(command)...); err != nil {
		p.Warning(err)
	}
}

func (p *ProxySQL) addMysqlUsersCharts(username string) {
	if err := p.Charts().Add(newMysqlUsersCharts(username)...); err != nil {
		p.Warning(err)
	}
}
