// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"database/sql"
	"strings"

	"github.com/jackc/pgx/v4/stdlib"
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

		if conn.db == nil && conn.connErrors < 3 {
			db, connString, err := p.openSecondaryConnection(dbname)
			if err != nil {
				p.Warning(err)
				conn.connErrors++
				continue
			}

			conn.db, conn.connString = db, connString
		}
	}

	for dbname, conn := range p.dbConns {
		if seen[dbname] {
			continue
		}

		delete(p.dbConns, dbname)
		if conn.db != nil {
			_ = conn.db.Close()
			stdlib.UnregisterConnConfig(conn.connString)
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
			p.getTableMetrics(dbname, schema, name).updated = true
		case "seq_scan":
			p.getTableMetrics(dbname, schema, name).seqScan = parseInt(value)
		case "seq_tup_read":
			p.getTableMetrics(dbname, schema, name).seqTupRead = parseInt(value)
		case "idx_scan":
			p.getTableMetrics(dbname, schema, name).idxScan = parseInt(value)
		case "idx_tup_fetch":
			p.getTableMetrics(dbname, schema, name).idxTupFetch = parseInt(value)
		case "n_tup_ins":
			p.getTableMetrics(dbname, schema, name).nTupIns = parseInt(value)
		case "n_tup_upd":
			p.getTableMetrics(dbname, schema, name).nTupUpd = parseInt(value)
		case "n_tup_del":
			p.getTableMetrics(dbname, schema, name).nTupDel = parseInt(value)
		case "n_tup_hot_upd":
			p.getTableMetrics(dbname, schema, name).nTupHotUpd = parseInt(value)
		case "n_live_tup":
			p.getTableMetrics(dbname, schema, name).nLiveTup = parseInt(value)
		case "n_dead_tup":
			p.getTableMetrics(dbname, schema, name).nDeadTup = parseInt(value)
		case "n_mod_since_analyze":
		case "n_ins_since_vacuum":
		case "last_vacuum":
			p.getTableMetrics(dbname, schema, name).lastVacuumAgo = parseFloat(value)
		case "last_autovacuum":
			p.getTableMetrics(dbname, schema, name).lastAutoVacuumAgo = parseFloat(value)
		case "last_analyze":
			p.getTableMetrics(dbname, schema, name).lastAnalyzeAgo = parseFloat(value)
		case "last_autoanalyze":
			p.getTableMetrics(dbname, schema, name).lastAutoAnalyzeAgo = parseFloat(value)
		case "vacuum_count":
			p.getTableMetrics(dbname, schema, name).vacuumCount = parseInt(value)
		case "autovacuum_count":
			p.getTableMetrics(dbname, schema, name).autovacuumCount = parseInt(value)
		case "analyze_count":
			p.getTableMetrics(dbname, schema, name).analyzeCount = parseInt(value)
		case "autoanalyze_count":
			p.getTableMetrics(dbname, schema, name).autoAnalyzeCount = parseInt(value)
		case "total_relation_size":
			p.getTableMetrics(dbname, schema, name).totalSize = parseInt(value)
		}
	})
}
