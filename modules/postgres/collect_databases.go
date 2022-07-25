package postgres

import (
	"context"
	"database/sql"
	"strconv"
)

func (p *Postgres) queryDatabaseList() ([]string, error) {
	q := queryDatabasesList()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var dbs []string
	var db string
	for rows.Next() {
		if err := rows.Scan(&db); err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (p *Postgres) collectDatabaseList(dbs []string) {
	if len(dbs) == 0 {
		return
	}

	collected := make(map[string]bool)
	for _, db := range p.databases {
		collected[db] = true
	}
	p.databases = dbs

	seen := make(map[string]bool)
	for _, db := range dbs {
		seen[db] = true
		if !collected[db] {
			collected[db] = true
			p.addNewDatabaseCharts(db)
		}
	}
	for db := range collected {
		if !seen[db] {
			p.removeDatabaseCharts(db)
		}
	}
}

func (p *Postgres) collectDatabaseStats(mx map[string]int64) error {
	q := queryDatabasesStats(p.databases)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return collectDatabaseRows(mx, rows)
}

func (p *Postgres) collectDatabaseConflicts(mx map[string]int64) error {
	q := queryDatabasesConflicts(p.databases)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return collectDatabaseRows(mx, rows)
}

func collectDatabaseRows(mx map[string]int64, rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	vs := makeNullStrings(len(columns))

	for rows.Next() {
		if err := rows.Scan(vs...); err != nil {
			return err
		}

		var db string
		for i, name := range columns {
			s := valueToString(vs[i])
			switch name {
			case "datname":
				db = s
			default:
				if v, err := strconv.ParseInt(s, 10, 64); err == nil {
					mx["db_"+db+"_"+name] = v
				}
			}
		}
	}
	return nil
}
