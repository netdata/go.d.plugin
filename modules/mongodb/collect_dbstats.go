// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"fmt"
)

func (m *Mongo) collectDbStats(mx map[string]int64) error {
	if m.dbSelector == nil {
		m.Debug("'database' selector not set, skip collecting database statistics")
		return nil
	}

	allDBs, err := m.conn.listDatabaseNames()
	if err != nil {
		return fmt.Errorf("cannot get database names: %v", err)
	}

	m.Debugf("all databases on the server: '%v'", allDBs)

	var dbs []string
	for _, db := range allDBs {
		if m.dbSelector.MatchString(db) {
			dbs = append(dbs, db)
		}
	}

	if len(allDBs) != len(dbs) {
		m.Debugf("databases remaining after filtering: %v", dbs)
	}

	seen := make(map[string]bool)
	for _, db := range dbs {
		s, err := m.conn.dbStats(db)
		if err != nil {
			return fmt.Errorf("dbStats command failed: %v", err)
		}

		seen[db] = true

		mx["database_"+db+"_collections"] = s.Collections
		mx["database_"+db+"_views"] = s.Views
		mx["database_"+db+"_indexes"] = s.Indexes
		mx["database_"+db+"_documents"] = s.Objects
		mx["database_"+db+"_data_size"] = s.DataSize
		mx["database_"+db+"_index_size"] = s.IndexSize
		mx["database_"+db+"_storage_size"] = s.StorageSize
	}

	for db := range seen {
		if !m.databases[db] {
			m.databases[db] = true
			m.Debugf("new database '%s': creating charts", db)
			m.addDatabaseCharts(db)
		}
	}

	for db := range m.databases {
		if !seen[db] {
			delete(m.databases, db)
			m.Debugf("stale database '%s': removing charts", db)
			m.removeDatabaseCharts(db)
		}
	}

	return nil
}
