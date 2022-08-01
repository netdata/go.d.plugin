// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"context"
)

func (p *Postgres) collectReplicationStandbyAppWALDelta(mx map[string]int64) error {
	q := queryReplicationStandbyAppDelta(p.serverVersion)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	var app string
	return collectRows(rows, func(column, value string) {
		switch column {
		case "application_name":
			app = value
		default:
			mx["repl_standby_app_"+app+"_wal_"+column] += safeParseInt(value)
		}
	})
}

func (p *Postgres) collectReplicationStandbyAppWALLag(mx map[string]int64) error {
	q := queryReplicationStandbyAppLag()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	var app string
	return collectRows(rows, func(column, value string) {
		switch column {
		case "application_name":
			app = value
		default:
			mx["repl_standby_app_"+app+"_wal_"+column] += safeParseInt(value)
		}
	})
}

func (p *Postgres) queryReplicationStandbyAppList() ([]string, error) {
	q := queryReplicationStandbyAppList()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var apps []string
	seen := make(map[string]bool)
	if err := collectRows(rows, func(column, value string) {
		if column == "application_name" && !seen[value] {
			seen[value] = true
			apps = append(apps, value)
		}
	}); err != nil {
		return nil, err
	}

	return apps, nil
}

func (p *Postgres) collectReplicationStandbyAppList(apps []string) {
	if len(apps) == 0 {
		return
	}

	collected := make(map[string]bool)
	for _, db := range p.replStandbyApps {
		collected[db] = true
	}
	p.replStandbyApps = apps

	seen := make(map[string]bool)
	for _, app := range apps {
		seen[app] = true
		if !collected[app] {
			collected[app] = true
			p.addNewReplicationStandbyAppCharts(app)
		}
	}
	for app := range collected {
		if !seen[app] {
			p.removeReplicationStandbyAppCharts(app)
		}
	}
}
