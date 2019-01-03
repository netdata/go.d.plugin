package oracledb

import (
	"database/sql"
	"strings"

	"github.com/netdata/go.d.plugin/modules"

	_ "gopkg.in/goracle.v2"
)

const defaultDSN = "SYSTEM/Oracle12345@ORCL" // there is no such thing as default connstr/dsn, but...

func init() {
	creator := modules.Creator{
		// TODO: open a discussion about API changes with @ilyam for that creator thing and more.
		Create: func() modules.Module { return New(defaultDSN) },
	}

	modules.Register("oracledb", creator)
}

// OracleDB oracledb module.
type OracleDB struct {
	modules.Base
	DSN string `json:"dsn,omitempty" yaml:"dsn,omitempty"`
}

// New creates OracleDB mod.
func New(connString string) *OracleDB {
	return &OracleDB{
		DSN: defaultDSN,
	}
}

// Cleanup performs cleanup.
func (*OracleDB) Cleanup() {}

// Init makes initialization of the OracleDB mod.
func (m *OracleDB) Init() bool {
	// test the connectivity here, this session will be not actually used.
	// each metrics scraping will have its own session.
	db, err := sql.Open("goracle", m.DSN)
	if err != nil {
		m.Errorf("error on opening a connection with the oracle database [%s]: %v", m.DSN, err)
		return false
	}

	if err = db.Ping(); err != nil {
		db.Close()
		m.Errorf("error on pinging the oracle database [%s]: %v", m.DSN, err)
		return false
	}

	if err = db.Close(); err != nil {
		return false
	}

	// post Init debug info.
	m.Debugf("using DSN [%s]", m.DSN)
	return true
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
	metrics := make(map[string]int64)

	db, err := sql.Open("goracle", m.DSN)
	if err != nil {
		m.Errorf("error on opening a connection with the oracle database [%s]: %v", m.DSN, err)
		return metrics
	}
	defer db.Close()

	err = m.collectProcesses(db, metrics)
	if err != nil {
		m.Errorf("error on collecting processes: %v", err)
	}

	err = m.collectSessions(db, metrics)
	if err != nil {
		m.Errorf("error on collecting sessions: %v", err)
	}

	err = m.collectActivity(db, metrics)
	if err != nil {
		m.Errorf("error on collecting activity: %v", err)
	}

	err = m.collectWaitTime(db, metrics)
	if err != nil {
		m.Errorf("error on collecting wait time: %v", err)
	}

	err = m.collectTablespace(db, metrics)
	if err != nil {
		m.Errorf("error on collecting tablespace size: %v", err)
	}

	return metrics
}

// collectProcesses collects information about the currently active processes.
func (m *OracleDB) collectProcesses(db *sql.DB, metrics map[string]int64) error {
	var count int64
	err := db.QueryRow("SELECT COUNT(*) FROM v$process").Scan(&count)
	if err != nil {
		return err
	}

	metrics["processes"] = count
	return nil
}

// collectSessions collects all (active, inactive, total) sessions.
func (m *OracleDB) collectSessions(db *sql.DB, metrics map[string]int64) error {
	// sessions (active, inactive, total).
	rows, err := db.Query("SELECT status, type FROM v$session GROUP BY status, type")
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
	return strings.Replace(strings.Replace(strings.Replace(s, "(", "", -1), ")", "", -1), " ", "_", -1)
}

// collectActivity collects activity metrics.
func (m *OracleDB) collectActivity(db *sql.DB, metrics map[string]int64) error {
	rows, err := db.Query("SELECT name, value FROM v$sysstat WHERE name IN ('parse count (total)', 'execute count', 'user commits', 'user rollbacks')")
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
func (m *OracleDB) collectWaitTime(db *sql.DB, metrics map[string]int64) error {
	rows, err := db.Query("SELECT n.wait_class, round(m.time_waited/m.INTSIZE_CSEC,3) AAS from v$waitclassmetric  m, v$system_wait_class n where m.wait_class_id=n.wait_class_id and n.wait_class != 'Idle'")
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
func (m *OracleDB) collectTablespace(db *sql.DB, metrics map[string]int64) error {
	rows, err := db.Query(`
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
