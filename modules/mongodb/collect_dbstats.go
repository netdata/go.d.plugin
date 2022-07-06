// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

// collectDbStats creates the map[string]int64 for the available dims
// it calls listDatabase and then dbstats for each database internally
func (m *Mongo) collectDbStats(ms map[string]int64) error {
	if m.databasesMatcher == nil {
		return nil
	}
	allDatabases, err := m.mongoCollector.listDatabaseNames()
	if err != nil {
		return fmt.Errorf("cannot get database names: %s", err)
	}

	// filter matching databases and exclude non-matching
	var databases []string
	for _, database := range allDatabases {
		if m.databasesMatcher.MatchString(database) {
			databases = append(databases, database)
		}
	}

	// add dims for each database
	m.updateDBStatsCharts(databases)

	for _, database := range databases {
		stats, err := m.mongoCollector.dbStats(database)
		if err != nil {
			return fmt.Errorf("dbStats command failed: %s", err)
		}
		stats.toMap(database, ms)
	}
	return nil
}

// updateDBStatsCharts adds dimensions for new databases and
// removes for dropped
func (m *Mongo) updateDBStatsCharts(databases []string) {
	// remove dims for not existing databases
	m.removeDatabasesFromDBStatsCharts(databases)

	// add dimensions for new databases
	for _, database := range sliceDiff(databases, m.discoveredDBs) {
		for _, chart := range *m.chartsDbStats {
			id := chart.ID + "_" + database
			err := chart.AddDim(&module.Dim{ID: id, Name: database, Algo: module.Absolute})
			if err != nil {
				m.Warningf("failed to add dim: %s, %v", id, err)
				continue
			}
			chart.MarkNotCreated()
		}
	}

	// update the cache
	m.discoveredDBs = databases
}

// removeDatabasesFromDBStatsCharts removes dimensions from dbstats
// charts for dropped databases
func (m *Mongo) removeDatabasesFromDBStatsCharts(newDatabases []string) {
	diff := sliceDiff(m.discoveredDBs, newDatabases)
	for _, name := range diff {
		for _, chart := range *m.chartsDbStats {
			id := chart.ID + "_" + name
			err := chart.MarkDimRemove(id, true)
			if err != nil {
				m.Warningf("failed to remove dimension %s with error: %v", id, err)
				continue
			}
			chart.MarkNotCreated()
		}
	}
}
