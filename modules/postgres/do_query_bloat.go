package postgres

import "database/sql"

func (p *Postgres) doQueryBloat() error {
	if err := p.doDBQueryBloat(p.db); err != nil {
		p.Warning(err)
	}
	for _, conn := range p.dbConns {
		if conn.db == nil {
			continue
		}
		if err := p.doDBQueryBloat(conn.db); err != nil {
			p.Warning(err)
		}
	}
	return nil
}

func (p *Postgres) doDBQueryBloat(db *sql.DB) error {
	q := queryBloat()

	var dbname, schema, table, iname string
	var tableWasted, idxWasted int64
	return p.doDBQuery(db, q, func(column, value string, rowEnd bool) {
		switch column {
		case "db":
			dbname = value
		case "schemaname":
			schema = value
		case "tablename":
			table = value
		case "wastedbytes":
			tableWasted = parseInt(value)
		case "iname":
			iname = value
		case "wastedibytes":
			idxWasted = parseInt(value)
		}
		if !rowEnd {
			return
		}
		if p.hasTableMetrics(table, dbname, schema) {
			v := p.getTableMetrics(table, dbname, schema)
			v.bloatSize = tableWasted
			v.cachedTotalSize = v.totalSize
		}
		if iname != "?" && p.hasIndexMetrics(iname, table, dbname, schema) {
			v := p.getIndexMetrics(iname, table, dbname, schema)
			v.bloatSize = idxWasted
			v.cachedSize = v.size
		}
	})
}
