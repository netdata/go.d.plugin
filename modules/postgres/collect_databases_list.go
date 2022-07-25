package postgres

import "context"

func (p *Postgres) queryDatabasesList() ([]string, error) {
	q := queryDatabasesList()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (p *Postgres) collectDatabasesList(dbs []string) {
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
