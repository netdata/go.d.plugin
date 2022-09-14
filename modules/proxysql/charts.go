package proxysql

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	Charts = module.Charts
	Chart  = module.Chart
	Dims   = module.Dims
)

const (
	prioSmth = module.Priority + iota
	prioClientConnectionsCount
	prioClientConnectionsRate
	prioServerConnectionsCount
	prioServerConnectionsRate
	prioQueriesRate
	prioBackendStatementsCount
	prioBackendStatementsRate
	prioFrontendStatementsCount
	prioFrontendStatementsRate
	prioCachedStatementsCount
	prioQueryCacheEntriesCount
	prioQueryCacheIO
	prioQueryCacheRequestsRate
	prioQueryCacheMemoryUsedCount
	prioJemallocMemoryUsed
	prioMemoryUsed
	prioMySQLCommandExecutionsRate
	prioMySQLCommandExecutionTime
	prioMySQLCommandExecutionDurationHistogram
	prioMySQLUserConnectionsUtilization
	prioMySQLUserConnectionsCount
	prioUptime
)

var (
	baseCharts = module.Charts{
		clientConnectionsCount.Copy(),
		clientConnectionsRate.Copy(),
		serverConnectionsCount.Copy(),
		serverConnectionsRate.Copy(),
		queriesRate.Copy(),
		backendStatementsCount.Copy(),
		backendStatementsRate.Copy(),
		clientStatementsCount.Copy(),
		clientStatementsRate.Copy(),
		cachedStatementsCount.Copy(),
		queryCacheEntriesCount.Copy(),
		queryCacheIO.Copy(),
		queryCacheRequestsRate.Copy(),
		queryCacheMemoryUsedCount.Copy(),
		jemallocMemoryUsedChart.Copy(),
		memoryUsedCountChart.Copy(),
		uptimeChart.Copy(),
	}

	clientConnectionsCount = module.Chart{
		ID:       "client_connections_count",
		Title:    "Client connections",
		Units:    "connections",
		Fam:      "client connections",
		Ctx:      "proxysql.client_connections_count",
		Priority: prioClientConnectionsCount,
		Dims: Dims{
			{ID: "Client_Connections_connected", Name: "connected"},
			{ID: "Client_Connections_non_idle", Name: "non_idle"},
			{ID: "Client_Connections_hostgroup_locked", Name: "hostgroup_locked"},
		},
	}
	clientConnectionsRate = module.Chart{
		ID:       "client_connections_rate",
		Title:    "Client connections rate",
		Units:    "connections/s",
		Fam:      "client connections",
		Ctx:      "proxysql.client_connections_rate",
		Priority: prioClientConnectionsRate,
		Dims: Dims{
			{ID: "Client_Connections_created", Name: "created", Algo: module.Incremental},
			{ID: "Client_Connections_aborted", Name: "aborted", Algo: module.Incremental},
		},
	}

	serverConnectionsCount = module.Chart{
		ID:       "server_connections_count",
		Title:    "Server connections",
		Units:    "connections",
		Fam:      "server connections",
		Ctx:      "proxysql.server_connections_count",
		Priority: prioServerConnectionsCount,
		Dims: Dims{
			{ID: "Server_Connections_connected", Name: "connected"},
		},
	}
	serverConnectionsRate = module.Chart{
		ID:       "server_connections_rate",
		Title:    "Server connections rate",
		Units:    "connections/s",
		Fam:      "server connections",
		Ctx:      "proxysql.server_connections_rate",
		Priority: prioServerConnectionsRate,
		Dims: Dims{
			{ID: "Server_Connections_created", Name: "created", Algo: module.Incremental},
			{ID: "Server_Connections_aborted", Name: "aborted", Algo: module.Incremental},
			{ID: "Server_Connections_delayed", Name: "delayed", Algo: module.Incremental},
		},
	}

	queriesRate = module.Chart{
		ID:       "queries_rate",
		Title:    "Queries rate",
		Units:    "queries/s",
		Fam:      "queries",
		Ctx:      "proxysql.queries_rate",
		Priority: prioQueriesRate,
		Type:     module.Stacked,
		Dims: Dims{
			{ID: "Com_autocommit", Name: "autocommit", Algo: module.Incremental},
			{ID: "Com_autocommit_filtered", Name: "autocommit_filtered", Algo: module.Incremental},
			{ID: "Com_commit", Name: "commit", Algo: module.Incremental},
			{ID: "Com_commit_filtered", Name: "commit_filtered", Algo: module.Incremental},
			{ID: "Com_rollback", Name: "rollback", Algo: module.Incremental},
			{ID: "Com_rollback_filtered", Name: "rollback_filtered", Algo: module.Incremental},
			{ID: "Com_backend_change_user", Name: "backend_change_user", Algo: module.Incremental},
			{ID: "Com_backend_init_db", Name: "backend_init_db", Algo: module.Incremental},
			{ID: "Com_backend_set_names", Name: "backend_set_names", Algo: module.Incremental},
			{ID: "Com_frontend_init_db", Name: "frontend_init_db", Algo: module.Incremental},
			{ID: "Com_frontend_set_names", Name: "frontend_set_names", Algo: module.Incremental},
			{ID: "Com_frontend_use_db", Name: "frontend_use_db", Algo: module.Incremental},
		},
	}
	backendStatementsCount = module.Chart{
		ID:       "backend_statements_count",
		Title:    "Statements available across all backend connections",
		Units:    "statements",
		Fam:      "statements",
		Ctx:      "proxysql.backend_statements_count",
		Priority: prioBackendStatementsCount,
		Dims: Dims{
			{ID: "Stmt_Server_Active_Total", Name: "total"},
			{ID: "Stmt_Server_Active_Unique", Name: "unique"},
		},
	}
	backendStatementsRate = module.Chart{
		ID:       "backend_statements_rate",
		Title:    "Statements executed against the backends",
		Units:    "statements/s",
		Fam:      "statements",
		Ctx:      "proxysql.backend_statements_rate",
		Priority: prioBackendStatementsRate,
		Type:     module.Stacked,
		Dims: Dims{
			{ID: "Com_backend_stmt_prepare", Name: "prepare", Algo: module.Incremental},
			{ID: "Com_backend_stmt_execute", Name: "execute", Algo: module.Incremental},
			{ID: "Com_backend_stmt_close", Name: "close", Algo: module.Incremental},
		},
	}
	clientStatementsCount = module.Chart{
		ID:       "client_statements_count",
		Title:    "Statements that are in use by clients",
		Units:    "statements",
		Fam:      "statements",
		Ctx:      "proxysql.client_statements_count",
		Priority: prioFrontendStatementsCount,
		Dims: Dims{
			{ID: "Stmt_Client_Active_Total", Name: "total"},
			{ID: "Stmt_Client_Active_Unique", Name: "unique"},
		},
	}
	clientStatementsRate = module.Chart{
		ID:       "client_statements_rate",
		Title:    "Statements executed by clients",
		Units:    "statements/s",
		Fam:      "statements",
		Ctx:      "proxysql.client_statements_rate",
		Priority: prioFrontendStatementsRate,
		Type:     module.Stacked,
		Dims: Dims{
			{ID: "Com_frontend_stmt_prepare", Name: "prepare", Algo: module.Incremental},
			{ID: "Com_frontend_stmt_execute", Name: "execute", Algo: module.Incremental},
			{ID: "Com_frontend_stmt_close", Name: "close", Algo: module.Incremental},
		},
	}
	cachedStatementsCount = module.Chart{
		ID:       "cached_statements_count",
		Title:    "Global prepared statements",
		Units:    "statements",
		Fam:      "statements",
		Ctx:      "proxysql.cached_statements_count",
		Priority: prioCachedStatementsCount,
		Dims: Dims{
			{ID: "Stmt_Cached", Name: "cached"},
		},
	}

	queryCacheEntriesCount = module.Chart{
		ID:       "query_cache_entries_count",
		Title:    "Query Cache entries",
		Units:    "entries",
		Fam:      "query cache",
		Ctx:      "proxysql.query_cache_entries_count",
		Priority: prioQueryCacheEntriesCount,
		Dims: Dims{
			{ID: "Query_Cache_Entries", Name: "entries"},
		},
	}
	queryCacheMemoryUsedCount = module.Chart{
		ID:       "query_cache_memory_used_count",
		Title:    "Query Cache memory used",
		Units:    "B",
		Fam:      "query cache",
		Ctx:      "proxysql.query_cache_memory_used_count",
		Priority: prioQueryCacheMemoryUsedCount,
		Dims: Dims{
			{ID: "Query_Cache_Memory_bytes", Name: "used"},
		},
	}
	queryCacheIO = module.Chart{
		ID:       "query_cache_io",
		Title:    "Query Cache I/)",
		Units:    "B/s",
		Fam:      "query cache",
		Ctx:      "proxysql.query_cache_io",
		Priority: prioQueryCacheIO,
		Dims: Dims{
			{ID: "Query_Cache_bytes_IN", Name: "in", Algo: module.Incremental},
			{ID: "Query_Cache_bytes_OUT", Name: "out", Algo: module.Incremental},
		},
	}
	queryCacheRequestsRate = module.Chart{
		ID:       "query_cache_requests_rate",
		Title:    "Query Cache requests",
		Units:    "requests/s",
		Fam:      "query cache",
		Ctx:      "proxysql.query_cache_requests_rate",
		Priority: prioQueryCacheRequestsRate,
		Dims: Dims{
			{ID: "Query_Cache_count_GET", Name: "read", Algo: module.Incremental},
			{ID: "Query_Cache_count_SET", Name: "write", Algo: module.Incremental},
			{ID: "Query_Cache_count_GET_OK", Name: "read_success", Algo: module.Incremental},
		},
	}

	jemallocMemoryUsedChart = module.Chart{
		ID:       "jemalloc_memory_used",
		Title:    "Jemalloc used memory",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "proxysql.jemalloc_memory_used",
		Type:     module.Stacked,
		Priority: prioJemallocMemoryUsed,
		Dims: Dims{
			{ID: "jemalloc_active", Name: "active"},
			{ID: "jemalloc_allocated", Name: "allocated"},
			{ID: "jemalloc_mapped", Name: "mapped"},
			{ID: "jemalloc_metadata", Name: "metadata"},
			{ID: "jemalloc_resident", Name: "resident"},
			{ID: "jemalloc_retained", Name: "retained"},
		},
	}
	memoryUsedCountChart = module.Chart{
		ID:       "memory_used",
		Title:    "Memory used",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "proxysql.memory_used",
		Priority: prioMemoryUsed,
		Type:     module.Stacked,
		Dims: Dims{
			{ID: "Auth_memory", Name: "auth"},
			{ID: "SQLite3_memory_bytes", Name: "sqlite3"},
			{ID: "query_digest_memory", Name: "query_digest"},
			{ID: "mysql_query_rules_memory", Name: "query_rules"},
			{ID: "mysql_firewall_users_table", Name: "firewall_users_table"},
			{ID: "mysql_firewall_users_config", Name: "firewall_users_config"},
			{ID: "mysql_firewall_rules_table", Name: "firewall_rules_table"},
			{ID: "mysql_firewall_rules_config", Name: "firewall_rules_config"},
			{ID: "stack_memory_mysql_threads", Name: "mysql_threads"},
			{ID: "stack_memory_admin_threads", Name: "admin_threads"},
			{ID: "stack_memory_cluster_threads", Name: "cluster_threads"},
		},
	}
	uptimeChart = module.Chart{
		ID:       "proxysql_uptime",
		Title:    "Uptime",
		Units:    "seconds",
		Fam:      "uptime",
		Ctx:      "proxysql.uptime",
		Priority: prioUptime,
		Dims: Dims{
			{ID: "ProxySQL_Uptime", Name: "uptime"},
		},
	}
)

var charts = Charts{
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
			{ID: "mysql_backend_buffers_bytes", Name: "mysql backend buffers", Algo: module.Absolute, Mul: -8, Div: 1000},
			{ID: "mysql_frontend_buffers_bytes", Name: "mysql frontend buffers", Algo: module.Absolute, Mul: -8, Div: 1000},
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
}

var (
	mySQLCommandChartsTmpl = module.Charts{
		mySQLCommandExecutionRateChartTmpl.Copy(),
		mySQLCommandExecutionTimeChartTmpl.Copy(),
		mySQLCommandExecutionDurationHistogramChartTmpl.Copy(),
	}

	mySQLCommandExecutionRateChartTmpl = module.Chart{
		ID:       "mysql_command_%s_execution_rate",
		Title:    "MySQL command execution",
		Units:    "commands/s",
		Fam:      "command execution",
		Ctx:      "proxysql.mysql_command_execution_rate",
		Priority: prioMySQLCommandExecutionsRate,
		Dims: Dims{
			{ID: "mysql_command_%s_total_cnt", Name: "commands", Algo: module.Incremental},
		},
	}
	mySQLCommandExecutionTimeChartTmpl = module.Chart{
		ID:       "mysql_command_%s_execution_time",
		Title:    "MySQL command execution time",
		Units:    "microseconds",
		Fam:      "command execution time",
		Ctx:      "proxysql.mysql_command_execution_time",
		Priority: prioMySQLCommandExecutionTime,
		Dims: Dims{
			{ID: "mysql_command_%s_total_time_us", Name: "time", Algo: module.Incremental},
		},
	}
	mySQLCommandExecutionDurationHistogramChartTmpl = module.Chart{
		ID:       "mysql_command_%s_execution_duration",
		Title:    "MySQL command execution duration histogram",
		Units:    "commands/s",
		Fam:      "command execution duration",
		Ctx:      "proxysql.mysql_command_execution_duration",
		Type:     module.Stacked,
		Priority: prioMySQLCommandExecutionDurationHistogram,
		Dims: Dims{
			{ID: "mysql_command_%s_cnt_100us", Name: "100us", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_500us", Name: "500us", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_1ms", Name: "1ms", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_5ms", Name: "5ms", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_10ms", Name: "10ms", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_50ms", Name: "50ms", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_100ms", Name: "100ms", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_500ms", Name: "500ms", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_1s", Name: "1s", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_5s", Name: "5s", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_10s", Name: "10s", Algo: module.Incremental},
			{ID: "mysql_command_%s_cnt_INFs", Name: "+Inf", Algo: module.Incremental},
		},
	}
)

func newMySQLCommandCountersCharts(command string) *module.Charts {
	charts := mySQLCommandChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, strings.ToLower(command))
		chart.Labels = []module.Label{{Key: "command", Value: command}}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, command)
		}
	}

	return charts
}

func (p *ProxySQL) addMySQLCommandCountersCharts(command string) {
	charts := newMySQLCommandCountersCharts(command)

	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}

var (
	mySQLUserChartsTmpl = module.Charts{
		mySQLUserConnectionsUtilizationChartTmpl.Copy(),
		mySQLUserConnectionsCountChartTmpl.Copy(),
	}

	mySQLUserConnectionsUtilizationChartTmpl = module.Chart{
		ID:       "mysql_user_%s_connections_utilization",
		Title:    "MySQL user connections utilization",
		Units:    "percentage",
		Fam:      "user conns %",
		Ctx:      "proxysql.mysql_user_connections_utilization",
		Priority: prioMySQLUserConnectionsUtilization,
		Dims: Dims{
			{ID: "mysql_user_%s_frontend_connections_utilization", Name: "used"},
		},
	}
	mySQLUserConnectionsCountChartTmpl = module.Chart{
		ID:       "mysql_user_%s_connections_count",
		Title:    "MySQL user connections used",
		Units:    "connections",
		Fam:      "user conns",
		Ctx:      "proxysql.mysql_user_connections_count",
		Priority: prioMySQLUserConnectionsCount,
		Dims: Dims{
			{ID: "mysql_user_%s_frontend_connections", Name: "used"},
		},
	}
)

func newMySQLUserCharts(username string) *module.Charts {
	charts := mySQLUserChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, username)
		chart.Labels = []module.Label{{Key: "user", Value: username}}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, username)
		}
	}

	return charts
}

func (p *ProxySQL) addMysqlUsersCharts(username string) {
	charts := newMySQLUserCharts(username)

	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}
