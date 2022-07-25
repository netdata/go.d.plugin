// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"context"
	"database/sql"
	"strconv"
)

func (p *Postgres) collectCheckpoints(mx map[string]int64) error {
	q := queryCheckpoints()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return collectRows(mx, rows)
}

func collectRows(mx map[string]int64, rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := makeNullStrings(len(columns))

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}

		for i, name := range columns {
			s := valueToString(values[i])
			if v, err := strconv.ParseInt(s, 10, 64); err == nil {
				mx[name] = v
			}
		}
	}
	return nil
}
