// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"database/sql"
	"strings"
)

func (p *Postgres) doQueryTablesMetrics() error {
	if err := p.discoverQueryableDatabases(); err != nil {
		return err
	}
	if err := p.doQueryUserTableStats(); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) discoverQueryableDatabases() error {
	q := queryQueryableDatabaseList()

	var dbs []string
	if err := p.doQuery(q, func(_, value string, _ bool) { dbs = append(dbs, value) }); err != nil {
		return err
	}

	seen := make(map[string]bool, len(dbs))

	for _, dbname := range dbs {
		seen[dbname] = true

		conn, ok := p.dbConns[dbname]
		if !ok {
			conn = &dbConn{}
			p.dbConns[dbname] = conn
		}

		if conn.db != nil || conn.connErrors >= 3 {
			continue
		}

		var err error
		if conn.db, err = p.openSecondaryConnection(dbname); err != nil {
			p.Warning(err)
			conn.connErrors++
		}
	}

	for dbname, conn := range p.dbConns {
		if seen[dbname] {
			continue
		}
		delete(p.dbConns, dbname)
		if conn.db != nil {
			_ = conn.db.Close()
		}
	}

	return nil
}

func (p *Postgres) doQueryUserTableStats() error {
	if err := p.dbQueryUserTableStats(p.db); err != nil {
		p.Warning(err)
	}
	for _, conn := range p.dbConns {
		if conn.db == nil {
			continue
		}
		if err := p.dbQueryUserTableStats(conn.db); err != nil {
			p.Warning(err)
		}
	}
	return nil
}

func (p *Postgres) dbQueryUserTableStats(db *sql.DB) error {
	q := queryUserTableStats()

	var dbname, schema, name string
	return p.doDBQuery(db, q, func(column, value string, _ bool) {
		if value == "" && strings.HasPrefix(column, "last_") {
			value = "-1"
		}
		switch column {
		case "datname":
			dbname = value
		case "schemaname":
			schema = value
		case "relname":
			name = value
			p.getTableMetrics(name, dbname, schema).updated = true
		case "seq_scan":
			p.getTableMetrics(name, dbname, schema).seqScan = parseInt(value)
		case "seq_tup_read":
			p.getTableMetrics(name, dbname, schema).seqTupRead = parseInt(value)
		case "idx_scan":
			p.getTableMetrics(name, dbname, schema).idxScan = parseInt(value)
		case "idx_tup_fetch":
			p.getTableMetrics(name, dbname, schema).idxTupFetch = parseInt(value)
		case "n_tup_ins":
			p.getTableMetrics(name, dbname, schema).nTupIns = parseInt(value)
		case "n_tup_upd":
			p.getTableMetrics(name, dbname, schema).nTupUpd = parseInt(value)
		case "n_tup_del":
			p.getTableMetrics(name, dbname, schema).nTupDel = parseInt(value)
		case "n_tup_hot_upd":
			p.getTableMetrics(name, dbname, schema).nTupHotUpd = parseInt(value)
		case "n_live_tup":
			p.getTableMetrics(name, dbname, schema).nLiveTup = parseInt(value)
		case "n_dead_tup":
			p.getTableMetrics(name, dbname, schema).nDeadTup = parseInt(value)
		case "last_vacuum":
			p.getTableMetrics(name, dbname, schema).lastVacuumAgo = parseFloat(value)
		case "last_autovacuum":
			p.getTableMetrics(name, dbname, schema).lastAutoVacuumAgo = parseFloat(value)
		case "last_analyze":
			p.getTableMetrics(name, dbname, schema).lastAnalyzeAgo = parseFloat(value)
		case "last_autoanalyze":
			p.getTableMetrics(name, dbname, schema).lastAutoAnalyzeAgo = parseFloat(value)
		case "vacuum_count":
			p.getTableMetrics(name, dbname, schema).vacuumCount = parseInt(value)
		case "autovacuum_count":
			p.getTableMetrics(name, dbname, schema).autovacuumCount = parseInt(value)
		case "analyze_count":
			p.getTableMetrics(name, dbname, schema).analyzeCount = parseInt(value)
		case "autoanalyze_count":
			p.getTableMetrics(name, dbname, schema).autoAnalyzeCount = parseInt(value)
		case "total_relation_size":
			p.getTableMetrics(name, dbname, schema).totalSize = parseInt(value)
		}
	})
}
