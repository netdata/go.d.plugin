package postgres

import (
	"fmt"
)

func (p *Postgres) queryGlobalMetrics() error {
	if err := p.doQueryConnectionsUsed(); err != nil {
		return fmt.Errorf("querying server connections used error: %v", err)
	}
	if err := p.doQueryConnectionsState(); err != nil {
		return fmt.Errorf("querying server connections state error: %v", err)
	}
	if err := p.doQueryCheckpoints(); err != nil {
		return fmt.Errorf("querying database conflicts error: %v", err)
	}
	if err := p.doQueryUptime(); err != nil {
		return fmt.Errorf("querying server uptime error: %v", err)
	}
	if err := p.doQueryTXIDWraparound(); err != nil {
		return fmt.Errorf("querying txid wraparound error: %v", err)
	}
	if err := p.doQueryWALWrites(); err != nil {
		return fmt.Errorf("querying wal writes error: %v", err)
	}
	if err := p.doQueryCatalogRelations(); err != nil {
		return fmt.Errorf("querying catalog relations error: %v", err)
	}
	if p.pgVersion >= pgVersion94 {
		if err := p.doQueryAutovacuumWorkers(); err != nil {
			return fmt.Errorf("querying autovacuum workers error: %v", err)
		}
	}

	if !p.isSuperUser() {
		return nil
	}

	if p.pgVersion >= pgVersion94 {
		if err := p.doQueryWALFiles(); err != nil {
			return fmt.Errorf("querying wal files error: %v", err)
		}
	}
	if err := p.doQueryWALArchiveFiles(); err != nil {
		return fmt.Errorf("querying wal archive files error: %v", err)
	}

	return nil
}

func (p *Postgres) doQueryConnectionsUsed() error {
	q := queryServerCurrentConnectionsUsed()

	var v string
	if err := p.doQueryRow(q).Scan(&v); err != nil {
		return err
	}

	p.metrics.connUsed = parseInt(v)

	return nil
}

func (p *Postgres) doQueryConnectionsState() error {
	q := queryServerConnectionsState()

	var state string
	return p.doQueryRows(q, func(column, value string, rowEnd bool) {
		switch column {
		case "state":
			state = value
		case "count":
			switch state {
			case "active":
				p.metrics.connStateActive = parseInt(value)
			case "idle":
				p.metrics.connStateIdle = parseInt(value)
			case "idle_in_transaction":
				p.metrics.connStateIdleInTrans = parseInt(value)
			case "idle_in_transaction (aborted)":
				p.metrics.connStateIdleInTransAborted = parseInt(value)
			case "fastpath function call":
				p.metrics.connStateFastpathFunctionCall = parseInt(value)
			case "disabled":
				p.metrics.connStateDisabled = parseInt(value)
			}
		}
	})
}

func (p *Postgres) doQueryCheckpoints() error {
	q := queryCheckpoints()

	return p.doQueryRows(q, func(column, value string, _ bool) {
		switch column {
		case "checkpoints_timed":
			p.metrics.checkpointsTimed = parseInt(value)
		case "checkpoints_req":
			p.metrics.checkpointsReq = parseInt(value)
		case "checkpoint_write_time":
			p.metrics.checkpointWriteTime = parseInt(value)
		case "checkpoint_sync_time":
			p.metrics.checkpointSyncTime = parseInt(value)
		case "buffers_checkpoint":
			p.metrics.buffersCheckpoint = parseInt(value)
		case "buffers_clean":
			p.metrics.buffersClean = parseInt(value)
		case "maxwritten_clean":
			p.metrics.maxwrittenClean = parseInt(value)
		case "buffers_backend":
			p.metrics.buffersBackend = parseInt(value)
		case "buffers_backend_fsync":
			p.metrics.buffersBackendFsync = parseInt(value)
		case "buffers_alloc":
			p.metrics.buffersAlloc = parseInt(value)
		}
	})
}

func (p *Postgres) doQueryUptime() error {
	q := queryServerUptime()

	var s string
	if err := p.doQueryRow(q).Scan(&s); err != nil {
		return err
	}

	p.metrics.uptime = parseInt(s)

	return nil
}

func (p *Postgres) doQueryTXIDWraparound() error {
	q := queryTXIDWraparound()

	return p.doQueryRows(q, func(column, value string, _ bool) {
		switch column {
		case "oldest_current_xid":
			p.metrics.oldestXID = parseInt(value)
		case "percent_towards_wraparound":
			p.metrics.percentTowardsWraparound = parseInt(value)
		case "percent_towards_emergency_autovacuum":
			p.metrics.percentTowardsEmergencyAutovacuum = parseInt(value)
		}
	})
}

func (p *Postgres) doQueryWALWrites() error {
	q := queryWALWrites(p.pgVersion)

	var v int64
	if err := p.doQueryRow(q).Scan(&v); err != nil {
		return err
	}

	p.metrics.walWrites = v

	return nil
}

func (p *Postgres) doQueryWALFiles() error {
	q := queryWALFiles(p.pgVersion)

	return p.doQueryRows(q, func(column, value string, _ bool) {
		switch column {
		case "wal_recycled_files":
			p.metrics.walRecycledFiles = parseInt(value)
		case "wal_written_files":
			p.metrics.walWrittenFiles = parseInt(value)
		}
	})
}

func (p *Postgres) doQueryWALArchiveFiles() error {
	q := queryWALArchiveFiles(p.pgVersion)

	return p.doQueryRows(q, func(column, value string, _ bool) {
		switch column {
		case "wal_archive_files_ready_count":
			p.metrics.walArchiveFilesReady = parseInt(value)
		case "wal_archive_files_done_count":
			p.metrics.walArchiveFilesDone = parseInt(value)
		}
	})
}

func (p *Postgres) doQueryCatalogRelations() error {
	q := queryCatalogRelations()

	var kind string
	var count, size int64
	return p.doQueryRows(q, func(column, value string, rowEnd bool) {
		switch column {
		case "relkind":
			kind = value
		case "count":
			count = parseInt(value)
		case "size":
			count = parseInt(value)
		}
		if !rowEnd {
			return
		}
		// https://www.postgresql.org/docs/current/catalog-pg-class.html
		switch kind {
		case "r":
			p.metrics.relkindOrdinaryTable = count
			p.metrics.relkindOrdinaryTableSize = size
		case "i":
			p.metrics.relkindIndex = count
			p.metrics.relkindIndexSize = size
		case "S":
			p.metrics.relkindSequence = count
			p.metrics.relkindSequenceSize = size
		case "t":
			p.metrics.relkindTOASTTable = count
			p.metrics.relkindTOASTTableSize = size
		case "v":
			p.metrics.relkindView = count
			p.metrics.relkindViewSize = size
		case "m":
			p.metrics.relkindMatView = count
			p.metrics.relkindMatViewSize = size
		case "c":
			p.metrics.relkindCompositeType = count
			p.metrics.relkindCompositeTypeSize = size
		case "f":
			p.metrics.relkindForeignTable = count
			p.metrics.relkindForeignTableSize = size
		case "p":
			p.metrics.relkindPartitionedTable = count
			p.metrics.relkindPartitionedTableSize = size
		case "I":
			p.metrics.relkindPartitionedIndex = count
			p.metrics.relkindPartitionedIndexSize = size
		}
	})
}

func (p *Postgres) doQueryAutovacuumWorkers() error {
	q := queryAutovacuumWorkers()

	return p.doQueryRows(q, func(column, value string, _ bool) {
		switch column {
		case "autovacuum_analyze":
			p.metrics.autovacuumWorkersAnalyze = parseInt(value)
		case "autovacuum_vacuum_analyze":
			p.metrics.autovacuumWorkersVacuumAnalyze = parseInt(value)
		case "autovacuum_vacuum":
			p.metrics.autovacuumWorkersVacuum = parseInt(value)
		case "autovacuum_vacuum_freeze":
			p.metrics.autovacuumWorkersVacuumFreeze = parseInt(value)
		case "autovacuum_brin_summarize":
			p.metrics.autovacuumWorkersBrinSummarize = parseInt(value)
		}
	})
}
