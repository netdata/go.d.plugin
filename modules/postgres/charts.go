package postgres

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioDBTransactions = module.Priority + iota
	prioDBConnections
	prioDBBufferCache
	prioDBReadOperations
	prioDBWriteOperations
	prioDBConflicts
	prioDBTempFiles
	prioDBTempFilesData
	prioDBSize
)

var (
	dbChartsTmpl = module.Charts{
		dbTransactionsChartTmpl.Copy(),
		dbConnectionsChartTmpl.Copy(),
		dbBufferCacheChartTmpl.Copy(),
		dbReadOpsChartTmpl.Copy(),
		dbWriteOpsChartTmpl.Copy(),
		dbConflictsChartTmpl.Copy(),
		dbTempFilesChartTmpl.Copy(),
		dbTempFilesDataChartTmpl.Copy(),
		dbSizeChartTmpl.Copy(),
	}
	dbTransactionsChartTmpl = module.Chart{
		ID:       "db_%s_transactions",
		Title:    "Database transactions",
		Units:    "transactions/s",
		Fam:      "db transactions",
		Ctx:      "postgres.db_transactions",
		Priority: prioDBTransactions,
		Dims: module.Dims{
			{ID: "db_%s_xact_commit", Name: "committed", Algo: module.Incremental},
			{ID: "db_%s_xact_rollback", Name: "rollback", Algo: module.Incremental},
		},
	}
	dbConnectionsChartTmpl = module.Chart{
		ID:       "db_%s_connections",
		Title:    "Database connections",
		Units:    "connections",
		Fam:      "db connections",
		Ctx:      "postgres.db_connections",
		Priority: prioDBConnections,
		Dims: module.Dims{
			{ID: "db_%s_numbackends", Name: "connections"},
		},
	}
	dbBufferCacheChartTmpl = module.Chart{
		ID:       "db_%s_buffer_cache",
		Title:    "Database buffer cache",
		Units:    "blocks/s",
		Fam:      "db buffer cache",
		Ctx:      "postgres.db_buffer_cache",
		Priority: prioDBBufferCache,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "db_%s_blks_hit", Name: "hit", Algo: module.Incremental},
			{ID: "db_%s_blks_read", Name: "miss", Algo: module.Incremental},
		},
	}
	dbReadOpsChartTmpl = module.Chart{
		ID:       "db_%s_read_operations",
		Title:    "Database read operations",
		Units:    "rows/s",
		Fam:      "db operations",
		Ctx:      "postgres.db_read_operations",
		Priority: prioDBReadOperations,
		Dims: module.Dims{
			{ID: "db_%s_tup_returned", Name: "returned", Algo: module.Incremental},
			{ID: "db_%s_tup_fetched", Name: "fetched", Algo: module.Incremental},
		},
	}
	dbWriteOpsChartTmpl = module.Chart{
		ID:       "db_%s_write_operations",
		Title:    "Database write operations",
		Units:    "rows/s",
		Fam:      "db operations",
		Ctx:      "postgres.db_write_operations",
		Priority: prioDBWriteOperations,
		Dims: module.Dims{
			{ID: "db_%s_tup_inserted", Name: "inserted", Algo: module.Incremental},
			{ID: "db_%s_tup_deleted", Name: "deleted", Algo: module.Incremental},
			{ID: "db_%s_tup_updated", Name: "updated", Algo: module.Incremental},
		},
	}
	dbConflictsChartTmpl = module.Chart{
		ID:       "db_%s_conflicts",
		Title:    "Database canceled queries",
		Units:    "queries/s",
		Fam:      "db operations",
		Ctx:      "postgres.db_conflicts",
		Priority: prioDBConflicts,
		Dims: module.Dims{
			{ID: "db_%s_conflicts", Name: "conflicts", Algo: module.Incremental},
		},
	}
	dbTempFilesChartTmpl = module.Chart{
		ID:       "db_%s_temp_files",
		Title:    "Database temporary files written to disk",
		Units:    "files/s",
		Fam:      "db temp files",
		Ctx:      "postgres.db_temp_files",
		Priority: prioDBTempFiles,
		Dims: module.Dims{
			{ID: "db_%s_temp_files", Name: "written", Algo: module.Incremental},
		},
	}
	dbTempFilesDataChartTmpl = module.Chart{
		ID:       "db_%s_temp_files_data",
		Title:    "Database temporary files data written to disk",
		Units:    "B/s",
		Fam:      "db temp files",
		Ctx:      "postgres.db_temp_files_data",
		Priority: prioDBTempFilesData,
		Dims: module.Dims{
			{ID: "db_%s_temp_bytes", Name: "written", Algo: module.Incremental},
		},
	}
	dbSizeChartTmpl = module.Chart{
		ID:       "db_%s_size",
		Title:    "Database size",
		Units:    "B",
		Fam:      "db size",
		Ctx:      "postgres.db_size",
		Priority: prioDBSize,
		Dims: module.Dims{
			{ID: "db_%s_size", Name: "size"},
		},
	}
)

func newDatabaseCharts(dbname string) *module.Charts {
	charts := dbChartsTmpl.Copy()
	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, dbname)
		c.Labels = []module.Label{
			{Key: "database", Value: dbname},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, dbname)
		}
	}
	return charts
}

func (p *Postgres) addNewDatabaseCharts(dbname string) {
	charts := newDatabaseCharts(dbname)
	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}

func (p *Postgres) removeDatabaseCharts(dbname string) {
	prefix := fmt.Sprintf("db_%s_", dbname)
	for _, c := range *p.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}
