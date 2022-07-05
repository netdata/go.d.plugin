package proxysql

import (
	"strconv"
	"strings"
)

const queryMysqlGlobalStatus = "SELECT Variable_Name, Variable_Value FROM stats_mysql_global"

var mysqlGlobalStatusMetrics = []string{
	"Access_Denied_Max_Connections",
	"Access_Denied_Max_User_Connections",
	"Access_Denied_Wrong_Password",
	"Active_Transactions",
	"Backend_query_time_nsec",
	"Client_Connections_aborted",
	"Client_Connections_connected",
	"Client_Connections_created",
	"Client_Connections_hostgroup_locked",
	"Client_Connections_non_idle",
	"Com_autocommit",
	"Com_autocommit_filtered",
	"Com_backend_change_user",
	"Com_backend_init_db",
	"Com_backend_set_names",
	"Com_backend_stmt_close",
	"Com_backend_stmt_execute",
	"Com_backend_stmt_prepare",
	"Com_commit",
	"Com_commit_filtered",
	"Com_frontend_init_db",
	"Com_frontend_set_names",
	"Com_frontend_stmt_close",
	"Com_frontend_stmt_execute",
	"Com_frontend_stmt_prepare",
	"Com_frontend_use_db",
	"Com_rollback",
	"Com_rollback_filtered",
	"ConnPool_get_conn_failure",
	"ConnPool_get_conn_immediate",
	"ConnPool_get_conn_latency_awareness",
	"ConnPool_get_conn_success",
	"ConnPool_memory_bytes",
	"MyHGM_myconnpoll_destroy",
	"MyHGM_myconnpoll_get",
	"MyHGM_myconnpoll_get_ok",
	"MyHGM_myconnpoll_push",
	"MyHGM_myconnpoll_reset",
	"MySQL_Monitor_Workers",
	"MySQL_Monitor_Workers_Aux",
	"MySQL_Monitor_Workers_Started",
	"MySQL_Monitor_connect_check_ERR",
	"MySQL_Monitor_connect_check_OK",
	"MySQL_Monitor_ping_check_ERR",
	"MySQL_Monitor_ping_check_OK",
	"MySQL_Monitor_read_only_check_ERR",
	"MySQL_Monitor_read_only_check_OK",
	"MySQL_Monitor_replication_lag_check_ERR",
	"MySQL_Monitor_replication_lag_check_OK",
	"MySQL_Thread_Workers",
	"ProxySQL_Uptime",
	"Queries_backends_bytes_recv",
	"Queries_backends_bytes_sent",
	"Queries_frontends_bytes_recv",
	"Queries_frontends_bytes_sent",
	"Query_Cache_Entries",
	"Query_Cache_Memory_bytes",
	"Query_Cache_Purged",
	"Query_Cache_bytes_IN",
	"Query_Cache_bytes_OUT",
	"Query_Cache_count_GET",
	"Query_Cache_count_GET_OK",
	"Query_Cache_count_SET",
	"Query_Processor_time_nsec",
	"Questions",
	"SQLite3_memory_bytes",
	"Selects_for_update__autocommit0",
	"Server_Connections_aborted",
	"Server_Connections_connected",
	"Server_Connections_created",
	"Server_Connections_delayed",
	"Servers_table_version",
	"Slow_queries",
	"Stmt_Cached",
	"Stmt_Client_Active_Total",
	"Stmt_Client_Active_Unique",
	"Stmt_Max_Stmt_id",
	"Stmt_Server_Active_Total",
	"Stmt_Server_Active_Unique",
	"automatic_detected_sql_injection",
	"aws_aurora_replicas_skipped_during_query",
	"backend_lagging_during_query",
	"backend_offline_during_query",
	"generated_error_packets",
	"hostgroup_locked_queries",
	"hostgroup_locked_set_cmds",
	"max_connect_timeouts",
	"mysql_backend_buffers_bytes",
	"mysql_frontend_buffers_bytes",
	"mysql_killed_backend_connections",
	"mysql_killed_backend_queries",
	"mysql_session_internal_bytes",
	"mysql_unexpected_frontend_com_quit",
	"mysql_unexpected_frontend_packets",
	"queries_with_max_lag_ms",
	"queries_with_max_lag_ms__delayed",
	"queries_with_max_lag_ms__total_wait_time_us",
	"whitelisted_sqli_fingerprint",
}

func (p *ProxySQL) collectMysqlGlobalStatus(collected map[string]int64) error {
	// https://proxysql.com/documentation/stats-statistics/#stats_mysql_global
	p.Debugf("executing query: '%s'", queryMysqlGlobalStatus)

	rows, err := p.db.Query(queryMysqlGlobalStatus)
	if err != nil {
		return err
	}
	defer rows.Close()

	set, err := rowsAsMap(rows)
	if err != nil {
		return err
	}

	for _, name := range mysqlGlobalStatusMetrics {
		strValue, ok := set[name]
		if !ok {
			continue
		}
		value, err := parseMysqlGlobalStatusValue(strValue)
		if err != nil {
			continue
		}
		collected[strings.ToLower(name)] = value
	}
	return nil
}

func parseMysqlGlobalStatusValue(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
