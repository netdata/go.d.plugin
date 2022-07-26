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
	num, err := strconv.Atoi(v)
	if err != nil {
		return err
	}

	if p.maxConnections != 0 {
		mx["server_connections_available"] = int64(p.maxConnections - num)
		mx["server_connections_utilization"] = int64((num * 100) / p.maxConnections)
	}
	mx["server_connections_used"] = int64(num)

	return nil
}
