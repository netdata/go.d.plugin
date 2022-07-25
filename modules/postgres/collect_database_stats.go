package postgres

import (
	"context"
	"strconv"
)

func (p *Postgres) collectDatabasesStats(mx map[string]int64) error {
	// https://github.com/postgres/postgres/blob/366283961ac0ed6d89014444c6090f3fd02fce0a/src/backend/catalog/system_views.sql#L1018
	q := queryDatabasesStats(p.databases)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()

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
