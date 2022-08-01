// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"context"
)

func (p *Postgres) collectReplicationSlotFiles(mx map[string]int64) error {
	q := queryReplicationSlotFiles(p.serverVersion)

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
