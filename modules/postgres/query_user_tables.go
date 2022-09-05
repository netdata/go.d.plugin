package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func (p *Postgres) doQueryUserTableStats() error {
	q := queryDatabaseList()

	var dbs []string
	if err := p.doQuery(q, func(_, value string, _ bool) { dbs = append(dbs, value) }); err != nil {
		return err
	}

	for _, dbname := range dbs {
		conn, ok := p.dbConns[dbname]
		if !ok {
			cfg, err := pgx.ParseConfig(p.DSN)
			if err != nil {
				p.Warningf("error on parsing DSN: %v", err)
				continue
			}

			cfg.Database = dbname
			connString := stdlib.RegisterConnConfig(cfg)

			db, err := sql.Open("pgx", connString)
			if err != nil {
				p.Warningf("error on opening a connection with the Postgres database [%s]: %v", p.DSN, err)
				continue
			}

			db.SetMaxOpenConns(1)
			db.SetMaxIdleConns(1)
			db.SetConnMaxLifetime(10 * time.Minute)

			func() {
				ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
				defer cancel()

				if err := db.PingContext(ctx); err != nil {
					_ = db.Close()
					p.Warningf("error on pinging the Postgres database [%s]: %v", p.DSN, err)
					return
				}
				conn = &dbConn{
					connString: connString,
					db:         db,
				}
				p.dbConns[dbname] = conn

			}()
		}
		if conn == nil || conn.db == nil {
			continue
		}

		var schema, name string
		err := p.doDBQuery(conn.db, queryUserTableStats(), func(column, value string, _ bool) {
			if value == "" {
				value = "-1"
			}
			switch column {
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
		if err != nil {
			p.Warning(err)
		}
	}
	return nil
}
