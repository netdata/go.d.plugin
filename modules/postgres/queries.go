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
