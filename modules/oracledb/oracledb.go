package oracledb

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/modules"

	"gopkg.in/goracle.v2"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("oracledb", creator)
}

// OracleDB oracledb module.
type OracleDB struct {
	modules.Base
	DSN string `yaml:"dsn"`
	db  *sql.DB
}

// New creates OracleDB mod.
func New() *OracleDB {
	return &OracleDB{}
}

// Cleanup performs cleanup.
func (m *OracleDB) Cleanup() {
	err := m.db.Close()
	if err != nil {
		m.Errorf("cleanup: error on closing the oracle database [%s]: %v", m.DSN, err)
	}
}

// Init makes initialization of the OracleDB mod.
func (m *OracleDB) Init() bool {
	if m.DSN == "" {
		m.Errorf("dsn is missing")
		return false
	}

	// test the connectivity here.
	if err := m.openConnection(); err != nil {
		return false
	}

	// post Init debug info.
	m.Debugf("using DSN [%s]", m.DSN)
	return true
}

func (m *OracleDB) openConnection() error {
	if m.db != nil {
		if err := m.db.Ping(); err != nil {
			m.db.Close()
			m.db = nil

			return m.openConnection()
		}

		return nil
	}

	db, err := sql.Open("goracle", m.DSN)
	if err != nil {
		m.Errorf("error on opening a connection with the oracle database [%s]: %v", m.DSN, err)
		return err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		m.Errorf("error on pinging the oracle database [%s]: %v", m.DSN, err)
		return err
	}

	m.db = db
	return nil
}

// Check makes check.
func (m *OracleDB) Check() bool {
	return len(m.Collect()) > 0
}

// Charts creates Charts.
func (m *OracleDB) Charts() *Charts {
	return charts.Copy()
}

// Collect collects health checks and metrics for OracleDB.
func (m *OracleDB) Collect() map[string]int64 {
	if err := m.openConnection(); err != nil {
		return nil
	}

	metrics := make(map[string]int64)

	err := m.collectProcesses(metrics)
	if err != nil {
		m.Errorf("error on collecting processes: %v", err)
		return nil
	}

	err = m.collectSessions(metrics)
	if err != nil {
		m.Errorf("error on collecting sessions: %v", err)
		return nil
	}

	err = m.collectActivity(metrics)
	if err != nil {
		m.Errorf("error on collecting activity: %v", err)
		return nil
	}

	err = m.collectWaitTime(metrics)
	if err != nil {
		m.Errorf("error on collecting wait time: %v", err)
		return nil
	}


	// err = m.collectTablespace(metrics)
	// if err != nil {
	// 	m.Errorf("error on collecting tablespace size: %v", err)
	// 	return nil
	// }

	// err = m.collectSystemMetrics(metrics)
	// if err != nil {
	// 	m.Errorf("error on collecting system metrics: %v", err)
	// 	return nil
	// }

	return metrics
}

// collectProcesses collects information about the currently active processes.
func (m *OracleDB) collectProcesses(metrics map[string]int64) error {
	var count int64
	err := m.db.QueryRow("SELECT COUNT(*) FROM v$process").Scan(&count)
	if err != nil {
		return err
	}

	metrics["processes"] = count
	return nil
}

// collectSessions collects all (active, inactive, total) sessions.
func (m *OracleDB) collectSessions(metrics map[string]int64) error {
	// sessions (active, inactive, total).
	rows, err := m.db.Query("SELECT status, type FROM v$session GROUP BY status, type")
	if err != nil {
		return err
	}
	defer rows.Close()

	var activeCount, inactiveCount, totalCount int64

	for rows.Next() {
		var (
			status string
			typ    string
		)

		if err := rows.Scan(&status, &typ); err != nil {
			return err
		}

		totalCount++

		if status == "ACTIVE" {
			activeCount++
		} else if status == "INACTIVE" {
			inactiveCount++
		}
	}

	metrics["total_sessions"] = totalCount
	metrics["active_sessions"] = activeCount
	metrics["inactive_sessions"] += inactiveCount

	return nil
}

func cleanActivityName(s string) string {
	s = strings.Replace(s, " ", "_", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	return s
}

// collectActivity collects activity metrics.
func (m *OracleDB) collectActivity(metrics map[string]int64) error {
	rows, err := m.db.Query("SELECT name, value FROM v$sysstat WHERE name IN ('parse count (total)', 'execute count', 'user commits', 'user rollbacks')")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name  string
			value int64
		)

		if err := rows.Scan(&name, &value); err != nil {
			return err
		}

		metrics["activity_"+cleanActivityName(name)] = value
	}

	/*
	   [activity_execute_count] = [70921]
	   [activity_user_commits] = [177]
	   [activity_user_rollbacks] = [0]
	   [activity_parse_count_total] = [34340]
	*/

	return nil
}

func cleanWaitTimeClassname(s string) string {
	s = strings.Replace(s, " ", "_", -1)
	s = strings.Replace(s, "/", "", -1)
	s = strings.ToLower(s)
	return s
}

// collectWaitTime collects wait time metrics from the v$waitclassmetric view.
func (m *OracleDB) collectWaitTime(metrics map[string]int64) error {
	rows, err := m.db.Query("SELECT n.wait_class, round(m.time_waited/m.INTSIZE_CSEC,3) AAS from v$waitclassmetric  m, v$system_wait_class n where m.wait_class_id=n.wait_class_id and n.wait_class != 'Idle'")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			name  string
			value float64 // i.e "0.034".
		)

		if err := rows.Scan(&name, &value); err != nil {
			return err
		}
		metrics["wait_time_"+cleanWaitTimeClassname(name)] = int64(value * 1000)
	}

	/*
	   [wait_time_concurrency] = [0]
	   [wait_time_system_io] = [2]
	   [wait_time_other] = [3]
	   [wait_time_application] = [0]
	   [wait_time_configuration] = [0]
	   [wait_time_administrative] = [0]
	   [wait_time_commit] = [0]
	   [wait_time_network] = [0]
	   [wait_time_user_io] = [0]
	*/

	return nil
}

//  collectTablespace collects tablespace size.
func (m *OracleDB) collectTablespace(metrics map[string]int64) error {
	rows, err := m.db.Query(`
SELECT
  Z.name,
  Z.bytes,
  Z.max_bytes,
  Z.free_bytes
FROM
(
  SELECT
    X.name                   as name,
    SUM(nvl(X.free_bytes,0)) as free_bytes,
    SUM(X.bytes)             as bytes,
    SUM(X.max_bytes)         as max_bytes
  FROM
    (
      SELECT
        ddf.tablespace_name as name,
        ddf.bytes as bytes,
        sum(coalesce(dfs.bytes, 0)) as free_bytes,
        CASE
          WHEN ddf.maxbytes = 0 THEN ddf.bytes
          ELSE ddf.maxbytes
        END as max_bytes
      FROM
        sys.dba_data_files ddf,
        sys.dba_tablespaces dt,
        sys.dba_free_space dfs
      WHERE ddf.tablespace_name = dt.tablespace_name
      AND ddf.file_id = dfs.file_id(+)
      GROUP BY
        ddf.tablespace_name,
        ddf.file_name,
        ddf.bytes,
        ddf.maxbytes
    ) X
  GROUP BY X.name
  UNION ALL
  SELECT
    Y.name                   as name,
    MAX(nvl(Y.free_bytes,0)) as free_bytes,
    SUM(Y.bytes)             as bytes,
    SUM(Y.max_bytes)         as max_bytes
  FROM
    (
      SELECT
        dtf.tablespace_name as name,
        dtf.bytes as bytes,
        (
          SELECT
            ((f.total_blocks - s.tot_used_blocks)*vp.value)
          FROM
            (SELECT tablespace_name, sum(used_blocks) tot_used_blocks FROM gv$sort_segment WHERE  tablespace_name!='DUMMY' GROUP BY tablespace_name) s,
            (SELECT tablespace_name, sum(blocks) total_blocks FROM dba_temp_files where tablespace_name !='DUMMY' GROUP BY tablespace_name) f,
            (SELECT value FROM v$parameter WHERE name = 'db_block_size') vp
          WHERE f.tablespace_name=s.tablespace_name AND f.tablespace_name = dtf.tablespace_name
        ) as free_bytes,
        CASE
          WHEN dtf.maxbytes = 0 THEN dtf.bytes
          ELSE dtf.maxbytes
        END as max_bytes
      FROM
        sys.dba_temp_files dtf
    ) Y
  GROUP BY Y.name
) Z, sys.dba_tablespaces dt
WHERE
  Z.name = dt.tablespace_name
`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name  string
			bytes int64
			maxB  int64
			freeB int64
		)

		if err := rows.Scan(&name, &bytes, &maxB, &freeB); err != nil {
			return err
		}

		name = cleanTablespaceName(name)
		metrics["tablespace_max_bytes_"+name] = maxB
		metrics["tablespace_free_bytes_"+name] = freeB
		metrics["tablespace_bytes_"+name] = bytes

		// Note that, here we can have a dynamic key as well, i.e "UNDOTBS1".
		/*
		   [tablespace_free_bytes_system] = [10420224]
		   [tablespace_max_bytes_undotbs1] = [34359721984]
		   [tablespace_bytes_temp] = [34603008]
		   [tablespace_max_bytes_users] = [34359721984]
		   [tablespace_bytes_users] = [5242880]
		   [tablespace_max_bytes_sysaux] = [34359721984]
		   [tablespace_free_bytes_sysaux] = [23658496]
		   [tablespace_bytes_sysaux] = [545259520]
		   [tablespace_free_bytes_undotbs1] = [4194304]
		   [tablespace_max_bytes_temp] = [34359721984]
		   [tablespace_free_bytes_temp] = [34603008]
		   [tablespace_free_bytes_users] = [4194304]
		   [tablespace_max_bytes_system] = [34359721984]
		   [tablespace_bytes_system] = [880803840]
		   [tablespace_bytes_undotbs1] = [335544320]
		*/
	}

	return nil
}

func cleanTablespaceName(s string) string {
	return strings.ToLower(s)
}

var systemMetricsKeys = map[string]string{
	"Buffer Cache Hit Ratio":          "system_buffer_cachehit_ratio",
	"Cursor Cache Hit Ratio":          "system_cursor_cachehit_ratio",
	"Library Cache Hit Ratio":         "system_library_cachehit_ratio",
	"Shared Pool Free %":              "system_shared_pool_free",
	"Physical Reads Per Sec":          "system_physical_reads",
	"Physical Writes Per Sec":         "system_physical_writes",
	"Enqueue Timeouts Per Sec":        "system_enqueue_timeouts",
	"GC CR Block Received Per Second": "system_gc_cr_block_received",
	"Global Cache Blocks Corrupted":   "system_cache_blocks_corrupt",
	"Global Cache Blocks Lost":        "system_cache_blocks_lost",
	"Logons Per Sec":                  "system_logons",
	"Average Active Sessions":         "system_active_sessions",
	"Long Table Scans Per Sec":        "system_long_table_scans",
	"SQL Service Response Time":       "system_service_response_time",
	"User Rollbacks Per Sec":          "system_user_rollbacks",
	"Total Sorts Per User Call":       "system_sorts_per_user_call",
	"Rows Per Sort":                   "system_rows_per_sort",
	"Disk Sort Per Sec":               "system_disk_sorts",
	"Memory Sorts Ratio":              "system_memory_sorts_ratio",
	"Database Wait Time Ratio":        "system_database_wait_time_ratio",
	"Session Limit %":                 "system_session_limit_usage",
	"Session Count":                   "system_session_count",
	"Temp Space Used":                 "system_temp_space_used",
}

func (m *OracleDB) collectSystemMetrics(metrics map[string]int64) error {
	rows, err := m.db.Query(`SELECT METRIC_NAME, VALUE FROM GV$SYSMETRIC ORDER BY BEGIN_TIME`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			metricName string
			value      goracle.Number
		)

		err = rows.Scan(&metricName, &value)
		if err != nil {
			return err
		}

		for oraName, goMetricName := range systemMetricsKeys {
			if oraName == metricName {
				metrics[goMetricName], _ = strconv.ParseInt(value.String(), 10, 64)
				break
			}
		}
	}

	/* The available:

	[Buffer Cache Hit Ratio] = [100]
	[Memory Sorts Ratio] = [100]
	[Redo Allocation Hit Ratio] = [100]
	[User Transaction Per Sec] = [0]
	[Physical Reads Per Sec] = [0]
	[Physical Reads Per Txn] = [0]
	[Physical Writes Per Sec] = [0]
	[Physical Writes Per Txn] = [0]
	[Physical Reads Direct Per Sec] = [0]
	[Physical Reads Direct Per Txn] = [0]
	[Physical Writes Direct Per Sec] = [0]
	[Physical Writes Direct Per Txn] = [0]
	[Physical Reads Direct Lobs Per Sec] = [0]
	[Physical Reads Direct Lobs Per Txn] = [0]
	[Physical Writes Direct Lobs Per Sec] = [0]
	[Physical Writes Direct Lobs  Per Txn] = [0]
	[Redo Generated Per Sec] = [13.5818908122503]
	[Redo Generated Per Txn] = [816]
	[Logons Per Sec] = [0.033288948069241]
	[Logons Per Txn] = [2]
	[Open Cursors Per Sec] = [0.599201065246338]
	[Open Cursors Per Txn] = [36]
	[User Commits Per Sec] = [0]
	[User Commits Percentage] = [0]
	[User Rollbacks Per Sec] = [0]
	[User Rollbacks Percentage] = [0]
	[User Calls Per Sec] = [0.299600532623169]
	[User Calls Per Txn] = [18]
	[Recursive Calls Per Sec] = [15.2796271637816]
	[Recursive Calls Per Txn] = [918]
	[Logical Reads Per Sec] = [92.32689747004]
	[Logical Reads Per Txn] = [5547]
	[DBWR Checkpoints Per Sec] = [0]
	[Background Checkpoints Per Sec] = [0]
	[Redo Writes Per Sec] = [0.033288948069241]
	[Redo Writes Per Txn] = [2]
	[Long Table Scans Per Sec] = [0]
	[Long Table Scans Per Txn] = [0]
	[Total Table Scans Per Sec] = [0.216378162450067]
	[Total Table Scans Per Txn] = [13]
	[Full Index Scans Per Sec] = [0]
	[Full Index Scans Per Txn] = [0]
	[Total Index Scans Per Sec] = [0.532623169107856]
	[Total Index Scans Per Txn] = [32]
	[Total Parse Count Per Sec] = [0.599201065246338]
	[Total Parse Count Per Txn] = [36]
	[Hard Parse Count Per Sec] = [0]
	[Hard Parse Count Per Txn] = [0]
	[Parse Failure Count Per Sec] = [0]
	[Parse Failure Count Per Txn] = [0]
	[Cursor Cache Hit Ratio] = [41.6666666666667]
	[Disk Sort Per Sec] = [0]
	[Disk Sort Per Txn] = [0]
	[Rows Per Sort] = [10.8166666666667]
	[Execute Without Parse Ratio] = [16.2790697674419]
	[Soft Parse Ratio] = [100]
	[User Calls Ratio] = [1.92307692307692]
	[Host CPU Utilization (%)] = [6.78085991678225]
	[Network Traffic Volume Per Sec] = [338.14913448735]
	[Enqueue Timeouts Per Sec] = [0]
	[Enqueue Timeouts Per Txn] = [0]
	[Enqueue Waits Per Sec] = [0]
	[Enqueue Waits Per Txn] = [0]
	[Enqueue Deadlocks Per Sec] = [0]
	[Enqueue Deadlocks Per Txn] = [0]
	[Enqueue Requests Per Sec] = [214.064580559254]
	[Enqueue Requests Per Txn] = [12861]
	[DB Block Gets Per Sec] = [0.449400798934754]
	[DB Block Gets Per Txn] = [27]
	[Consistent Read Gets Per Sec] = [91.8774966711052]
	[Consistent Read Gets Per Txn] = [5520]
	[DB Block Changes Per Sec] = [0.066577896138482]
	[DB Block Changes Per Txn] = [4]
	[Consistent Read Changes Per Sec] = [0]
	[Consistent Read Changes Per Txn] = [0]
	[CPU Usage Per Sec] = [0.130034953395473]
	[CPU Usage Per Txn] = [7.8125]
	[CR Blocks Created Per Sec] = [0]
	[CR Blocks Created Per Txn] = [0]
	[CR Undo Records Applied Per Sec] = [0]
	[CR Undo Records Applied Per Txn] = [0]
	[User Rollback UndoRec Applied Per Sec] = [0]
	[User Rollback Undo Records Applied Per Txn] = [0]
	[Leaf Node Splits Per Sec] = [0]
	[Leaf Node Splits Per Txn] = [0]
	[Branch Node Splits Per Sec] = [0]
	[Branch Node Splits Per Txn] = [0]
	[PX downgraded 1 to 25% Per Sec] = [0]
	[PX downgraded 25 to 50% Per Sec] = [0]
	[PX downgraded 50 to 75% Per Sec] = [0]
	[PX downgraded 75 to 99% Per Sec] = [0]
	[PX downgraded to serial Per Sec] = [0]
	[Physical Read Total IO Requests Per Sec] = [8.67177097203728]
	[Physical Read Total Bytes Per Sec] = [140442.07723036]
	[GC CR Block Received Per Second] = [0]
	[GC CR Block Received Per Txn] = [0]
	[GC Current Block Received Per Second] = [0]
	[GC Current Block Received Per Txn] = [0]
	[Global Cache Average CR Get Time] = [0]
	[Global Cache Average Current Get Time] = [0]
	[Physical Write Total IO Requests Per Sec] = [0.832223701731025]
	[Global Cache Blocks Corrupted] = [0]
	[Global Cache Blocks Lost] = [0]
	[Current Logons Count] = [63]
	[Current Open Cursors Count] = [42]
	[User Limit %] = [0.0000014668330553609]
	[SQL Service Response Time] = [0.0148824786324786]
	[Database Wait Time Ratio] = [43.9160086145011]
	[Database CPU Time Ratio] = [56.0839913854989]
	[Response Time Per Txn] = [13.93]
	[Row Cache Hit Ratio] = [100]
	[Row Cache Miss Ratio] = [0]
	[Library Cache Hit Ratio] = [100]
	[Library Cache Miss Ratio] = [0]
	[Shared Pool Free %] = [34.2378282546997]
	[PGA Cache Hit %] = [100]
	[Process Limit %] = [9.375]
	[Session Limit %] = [6.89890710382514]
	[Executions Per Txn] = [43]
	[Executions Per Sec] = [0.715712383488682]
	[Txns Per Logon] = [0]
	[Database Time Per Sec] = [0.231857523302264]
	[Physical Write Total Bytes Per Sec] = [13115.3129161119]
	[Physical Read IO Requests Per Sec] = [0]
	[Physical Read Bytes Per Sec] = [0]
	[Physical Write IO Requests Per Sec] = [0]
	[Physical Write Bytes Per Sec] = [0]
	[DB Block Changes Per User Call] = [0.222222222222222]
	[DB Block Gets Per User Call] = [1.5]
	[Executions Per User Call] = [2.38888888888889]
	[Logical Reads Per User Call] = [308.166666666667]
	[Total Sorts Per User Call] = [3.33333333333333]
	[Total Table Scans Per User Call] = [0.722222222222222]
	[Current OS Load] = [0.3544921875]
	[Streams Pool Usage Percentage] = [0]
	[PQ QC Session Count] = [0]
	[PQ Slave Session Count] = [0]
	[Queries parallelized Per Sec] = [0]
	[DML statements parallelized Per Sec] = [0]
	[DDL statements parallelized Per Sec] = [0]
	[PX operations not downgraded Per Sec] = [0]
	[Session Count] = [101]
	[Average Synchronous Single-Block Read Latency] = [0.18042226487524]
	[I/O Megabytes per Second] = [0.149800266311585]
	[I/O Requests per Second] = [9.50399467376831]
	[Average Active Sessions] = [0.00231857523302264]
	[Active Serial Sessions] = [1]
	[Active Parallel Sessions] = [0]
	[Captured user calls] = [0]
	[Replayed user calls] = [0]
	[Workload Capture and Replay status] = [0]
	[Background CPU Usage Per Sec] = [0.208055925432756]
	[Background Time Per Sec] = [0.0272239513981358]
	[Host CPU Usage Per Sec] = [81.3748335552597]
	[Cell Physical IO Interconnect Bytes] = [9225728]
	[Temp Space Used] = [0]
	[Total PGA Allocated] = [335860736]
	[Total PGA Used by SQL Workareas] = [0]
	[Run Queue Per Sec] = [0]
	[VM in bytes Per Sec] = [0]
	[VM out bytes Per Sec] = [0]
	[Buffer Cache Hit Ratio] = [100]
	[Total PGA Used by SQL Workareas] = [0]
	[User Transaction Per Sec] = [0]
	[Physical Reads Per Sec] = [0]
	[Physical Reads Per Txn] = [0]
	[Physical Writes Per Sec] = [0]
	[Physical Writes Per Txn] = [0]
	[Physical Reads Direct Per Sec] = [0]
	[Physical Reads Direct Per Txn] = [0]
	[Redo Generated Per Sec] = [0]
	[Redo Generated Per Txn] = [0]
	[Logons Per Sec] = [0]
	[Logons Per Txn] = [0]
	[User Calls Per Sec] = [0]
	[User Calls Per Txn] = [0]
	[Logical Reads Per Sec] = [0.199600798403194]
	[Logical Reads Per Txn] = [3]
	[Redo Writes Per Sec] = [0]
	[Redo Writes Per Txn] = [0]
	[Total Table Scans Per Sec] = [0]
	[Total Table Scans Per Txn] = [0]
	[Full Index Scans Per Sec] = [0]
	[Full Index Scans Per Txn] = [0]
	[Execute Without Parse Ratio] = [50]
	[Soft Parse Ratio] = [100]
	[Host CPU Utilization (%)] = [6.9472166777556]
	[DB Block Gets Per Sec] = [0]
	[DB Block Gets Per Txn] = [0]
	[Consistent Read Gets Per Sec] = [0.199600798403194]
	[Consistent Read Gets Per Txn] = [3]
	[DB Block Changes Per Sec] = [0]
	[DB Block Changes Per Txn] = [0]
	[Consistent Read Changes Per Sec] = [0]
	[Consistent Read Changes Per Txn] = [0]
	[Database CPU Time Ratio] = [0]
	[Library Cache Hit Ratio] = [100]
	[Shared Pool Free %] = [34.2204599380493]
	[Executions Per Txn] = [8]
	[Executions Per Sec] = [0.53226879574185]
	[Txns Per Logon] = [0]
	[Database Time Per Sec] = [0]
	[Average Active Sessions] = [0]
	[Host CPU Usage Per Sec] = [83.3666001330672]
	[Cell Physical IO Interconnect Bytes] = [688128]
	[Temp Space Used] = [0]
	[Total PGA Allocated] = [335860736]
	[Memory Sorts Ratio] = [100]
	*/
	return nil
}
