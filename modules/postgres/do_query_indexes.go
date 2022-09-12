package postgres

import (
	"database/sql"
)

func (p *Postgres) doQueryIndexesMetrics() error {
	if err := p.doQueryStatUserIndexes(); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) doQueryStatUserIndexes() error {
	if err := p.doDBQueryStatUserIndexes(p.db); err != nil {
		p.Warning(err)
	}
	for _, conn := range p.dbConns {
		if conn.db == nil {
			continue
		}
		if err := p.doDBQueryStatUserIndexes(conn.db); err != nil {
			p.Warning(err)
		}
	}
	return nil
}

func (p *Postgres) doDBQueryStatUserIndexes(db *sql.DB) error {
	q := queryStatUserIndexes()

	var dbname, schema, table, name string
	return p.doDBQuery(db, q, func(column, value string, _ bool) {
		switch column {
		case "datname":
			dbname = value
		case "schemaname":
			schema = value
		case "relname":
			table = value
		case "indexrelname":
			name = value
			p.getIndexMetrics(name, dbname, schema, table).updated = true
		case "size":
			p.getIndexMetrics(name, dbname, schema, table).size = parseInt(value)
		}
	})
}
