package postgres

import (
	"context"
	"strconv"
)

func (p *Postgres) collectConnection(mx map[string]int64) error {
	q := queryServerCurrentConnectionsNum()

	var v string
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	if err := p.db.QueryRowContext(ctx, q).Scan(&v); err != nil {
		return err
	}
	num, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return err
	}

	if p.maxConnections != 0 {
		mx["server_connections_available"] = p.maxConnections - num
		mx["server_connections_utilization"] = calcPercentage(num, p.maxConnections)
	}
	mx["server_connections_used"] = int64(num)

	return nil
}

func (p *Postgres) collectCheckpoints(mx map[string]int64) error {
	q := queryCheckpoints()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return collectRows(rows, func(column, value string) { mx[column] = safeParseInt(value) })
}

func (p *Postgres) collectUptime(mx map[string]int64) error {
	q := queryServerUptime()

	var s string
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	if err := p.db.QueryRowContext(ctx, q).Scan(&s); err != nil {
		return err
	}

	v, _ := strconv.ParseFloat(s, 64)
	mx["server_uptime"] = int64(v)

	return nil
}

func (p *Postgres) collectTXIDWraparound(mx map[string]int64) error {
	q := queryTXIDWraparound()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return collectRows(rows, func(column, value string) { mx[column] = safeParseInt(value) })
}

func (p *Postgres) collectWALWrites(mx map[string]int64) error {
	q := queryWALWrites(p.serverVersion)

	var v int64
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	if err := p.db.QueryRowContext(ctx, q).Scan(&v); err != nil {
		return err
	}

	mx["wal_writes"] = v
	return nil
}

func (p *Postgres) collectWALFiles(mx map[string]int64) error {
	q := queryWALFiles(p.serverVersion)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return collectRows(rows, func(column, value string) { mx[column] = safeParseInt(value) })
}

func (p *Postgres) collectWALArchiveFiles(mx map[string]int64) error {
	q := queryWALArchiveFiles(p.serverVersion)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return collectRows(rows, func(column, value string) { mx[column] = safeParseInt(value) })
}

func (p *Postgres) collectCatalog(mx map[string]int64) error {
	q := queryCatalogRelations()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	// https://www.postgresql.org/docs/current/catalog-pg-class.html
	// r = ordinary table
	// i = index
	// S = sequence
	// t = TOAST table
	// v = view
	// m = materialized view
	// c = composite type
	// f = foreign table
	// p = partitioned table
	// I = partitioned index

	for _, v := range []string{"r", "i", "S", "t", "v", "m", "c", "f", "p", "I"} {
		mx["catalog_relkind_"+v+"_count"] = 0
		mx["catalog_relkind_"+v+"_size"] = 0
	}

	var kind string
	return collectRows(rows, func(column, value string) {
		switch column {
		case "relkind":
			kind = value
		default:
			mx["catalog_relkind_"+kind+"_"+column] = safeParseInt(value)
		}
	})
}

func (p *Postgres) collectAutovacuumWorkers(mx map[string]int64) error {
	q := queryAutovacuumWorkers()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return collectRows(rows, func(column, value string) { mx[column] = safeParseInt(value) })
}
