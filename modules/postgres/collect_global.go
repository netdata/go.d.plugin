package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

func (p *Postgres) collectGlobalMetrics(mx map[string]int64) error {
	if err := p.collectConnections(mx); err != nil {
		return fmt.Errorf("querying server connections error: %v", err)
	}
	if err := p.collectConnectionsState(mx); err != nil {
		return fmt.Errorf("querying server connections state error: %v", err)
	}

	if err := p.collectCheckpoints(mx); err != nil {
		return fmt.Errorf("querying database conflicts error: %v", err)
	}

	if err := p.collectUptime(mx); err != nil {
		return fmt.Errorf("querying server uptime error: %v", err)
	}

	if err := p.collectTXIDWraparound(mx); err != nil {
		return fmt.Errorf("querying txid wraparound error: %v", err)
	}

	if err := p.collectWALWrites(mx); err != nil {
		return fmt.Errorf("querying wal writes error: %v", err)
	}

	if err := p.collectCatalogRelations(mx); err != nil {
		return fmt.Errorf("querying catalog relations error: %v", err)
	}

	if p.pgVersion >= pgVersion94 {
		if err := p.collectAutovacuumWorkers(mx); err != nil {
			return fmt.Errorf("querying autovacuum workers error: %v", err)
		}
	}

	if !p.isSuperUser() {
		return nil
	}

	if p.pgVersion >= pgVersion94 {
		if err := p.collectWALFiles(mx); err != nil {
			return fmt.Errorf("querying wal files error: %v", err)
		}
	}
	if err := p.collectWALArchiveFiles(mx); err != nil {
		return fmt.Errorf("querying wal archive files error: %v", err)
	}

	return nil
}

func (p *Postgres) collectConnections(mx map[string]int64) error {
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
	mx["server_connections_used"] = num

	return nil
}

func (p *Postgres) collectConnectionsState(mx map[string]int64) error {
	q := queryServerConnectionsState()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	mx["server_connections_state_active"] = 0
	mx["server_connections_state_idle"] = 0
	mx["server_connections_state_idle_in_transaction"] = 0
	mx["server_connections_state_idle_in_transaction_aborted"] = 0
	mx["server_connections_state_fastpath_function_call"] = 0
	mx["server_connections_state_disabled"] = 0

	r := strings.NewReplacer(" ", "_", "(", "", ")", "")
	var s string
	return collectRows(rows, func(column, value string) {
		switch column {
		case "state":
			s = r.Replace(value)
		case "count":
			mx["server_connections_state_"+s] = safeParseInt(value)
		}
	})
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

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	var s string
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
	q := queryWALWrites(p.pgVersion)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	var v int64
	if err := p.db.QueryRowContext(ctx, q).Scan(&v); err != nil {
		return err
	}

	mx["wal_writes"] = v
	return nil
}

func (p *Postgres) collectWALFiles(mx map[string]int64) error {
	q := queryWALFiles(p.pgVersion)

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
	q := queryWALArchiveFiles(p.pgVersion)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return collectRows(rows, func(column, value string) { mx[column] = safeParseInt(value) })
}

func (p *Postgres) collectCatalogRelations(mx map[string]int64) error {
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
