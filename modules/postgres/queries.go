package postgres

import (
	"fmt"
	"strings"
)

// http://sqllint.com/

func queryServerVersion() string {
	return "SHOW server_version_num;"
}

//func queryIsSuperUser() string {
//	return "SELECT current_setting('is_superuser') = 'on' AS is_superuser;"
//}

func queryDatabasesList() string {
	return `
    SELECT
        datname          
    FROM
        pg_stat_database          
    WHERE
        has_database_privilege((SELECT
            CURRENT_USER), datname, 'connect')                  
        AND NOT datname ~* '^template\d'          
    ORDER BY
        datname;
`
}

func queryDatabasesStats(dbs []string) string {
	// docs: https://www.postgresql.org/docs/current/monitoring-stats.html#MONITORING-PG-STAT-DATABASE-VIEW
	// code: https://github.com/postgres/postgres/blob/366283961ac0ed6d89014444c6090f3fd02fce0a/src/backend/catalog/system_views.sql#L1018
	q := `
    SELECT
        datname,
        numbackends,
        xact_commit,
        xact_rollback,
        blks_read,
        blks_hit,
        tup_returned,
        tup_fetched,
        tup_inserted,
        tup_updated,
        tup_deleted,
        conflicts,
        pg_database_size(datname) AS size,
        temp_files,
        temp_bytes,
        deadlocks     
    FROM
        pg_stat_database
`

	if len(dbs) > 0 {
		q += fmt.Sprintf(`
    WHERE
        datname IN (
            '%s'                 
        ) 
`, strings.Join(dbs, "','"))
	}

	q += ";"

	return q
}

func queryDatabasesConflicts(dbs []string) string {
	// docs: https://www.postgresql.org/docs/current/monitoring-stats.html#MONITORING-PG-STAT-DATABASE-CONFLICTS-VIEW
	// code: https://github.com/postgres/postgres/blob/366283961ac0ed6d89014444c6090f3fd02fce0a/src/backend/catalog/system_views.sql#L1058
	q := `
    SELECT
        datname,
        confl_tablespace,
        confl_lock,
        confl_snapshot,
        confl_bufferpin,
        confl_deadlock     
    FROM
        pg_stat_database_conflicts
`

	if len(dbs) > 0 {
		q += fmt.Sprintf(`
    WHERE
        datname IN (
            '%s'                 
        ) 
`, strings.Join(dbs, "','"))
	}

	q += ";"

	return q
}
