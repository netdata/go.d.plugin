// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"context"
	"fmt"
)

func (p *Postgres) collectReplicationMetrics(mx map[string]int64) error {
	if len(p.replStandbyApps) > 0 {
		if err := p.collectReplicationStandbyAppWALDelta(mx); err != nil {
			return fmt.Errorf("querying replication standby app wal delta error: %v", err)
		}
		if p.pgVersion >= pgVersion10 {
			if err := p.collectReplicationStandbyAppWALLag(mx); err != nil {
				return fmt.Errorf("querying replication standby app wal lag error: %v", err)
			}
		}
	}

	if p.pgVersion >= pgVersion10 && len(p.replSlots) > 0 && p.isSuperUser() {
		if err := p.collectReplicationSlotFiles(mx); err != nil {
			return fmt.Errorf("querying replication slot files error: %v", err)
		}
	}
	return nil
}

func (p *Postgres) collectReplicationStandbyAppWALDelta(mx map[string]int64) error {
	q := queryReplicationStandbyAppDelta(p.pgVersion)

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

func (p *Postgres) collectReplicationSlotFiles(mx map[string]int64) error {
	q := queryReplicationSlotFiles(p.pgVersion)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	var slot string
	return collectRows(rows, func(column, value string) {
		switch column {
		case "slot_name":
			slot = value
		case "slot_type":
		default:
			mx["repl_slot_"+slot+"_"+column] += safeParseInt(value)
		}
	})
}

func (p *Postgres) queryReplicationSlotList() ([]string, error) {
	q := queryReplicationSlotList()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var slots []string
	seen := make(map[string]bool)
	if err := collectRows(rows, func(column, value string) {
		if column == "slot_name" && !seen[value] {
			seen[value] = true
			slots = append(slots, value)
		}
	}); err != nil {
		return nil, err
	}

	return slots, nil
}

func (p *Postgres) collectReplicationSlotList(slots []string) {
	if len(slots) == 0 {
		return
	}

	collected := make(map[string]bool)
	for _, db := range p.replSlots {
		collected[db] = true
	}
	p.replSlots = slots

	seen := make(map[string]bool)
	for _, app := range slots {
		seen[app] = true
		if !collected[app] {
			collected[app] = true
			p.addNewReplicationSlotCharts(app)
		}
	}
	for app := range collected {
		if !seen[app] {
			p.removeReplicationSlotCharts(app)
		}
	}
}
