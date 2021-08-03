package mongo

import (
	"context"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"

	"go.mongodb.org/mongo-driver/bson"
)

func (m *Mongo) serverStatusCollect(ms map[string]int64) {
	var status map[string]interface{}
	command := bson.D{{Key: "serverStatus", Value: 1}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.Timeout)
	defer cancel()
	err := m.client.Database("admin").RunCommand(ctx, command).Decode(&status)
	if err != nil {
		m.Errorf("error get server status from mongo: %s", err)
		return
	}
	m.addOptionalCharts(status)
	iterateServerStatus(ms, status)
}

func (m *Mongo) addOptionalCharts(status map[string]interface{}) {
	m.metricExists(status, "transactions", &chartTransactionsCurrent)
	m.metricExists(status, "globalLock.activeClients", &chartGlobalLockActiveClients)
	m.metricExists(status, "catalogStats", &chartCollections)
	m.metricExists(status, "tcmalloc.generic", &chartTcmallocGeneric)
	m.metricExists(status, "tcmalloc.tcmalloc", &chartTcmalloc)
	m.metricExists(status, "globalLock.currentQueue", &chartGlobalLockCurrentQueue)
	m.metricExists(status, "metrics.commands", &chartMetricsCommands)
	m.metricExists(status, "locks.Global.acquireCount", &chartGlobalLocks)
	m.metricExists(status, "flowControl", &chartFlowControl)
	// WiredTiger charts
	m.metricExists(status, "wiredTiger.block-manager", &chartWiredTigerBlockManager)
	m.metricExists(status, "wiredTiger.cache", &chartWiredTigerCache)
	m.metricExists(status, "wiredTiger.capacity", &chartWiredTigerCapacity)
	m.metricExists(status, "wiredTiger.connection", &chartWiredTigerConnection)
	m.metricExists(status, "wiredTiger.cursor", &chartWiredTigerCursor)
	m.metricExists(status, "wiredTiger.lock", &chartWiredTigerLock)
	m.metricExists(status, "wiredTiger.lock", &chartWiredTigerLockDuration)
	m.metricExists(status, "wiredTiger.log", &chartWiredTigerLogOps)
	m.metricExists(status, "wiredTiger.log", &chartWiredTigerLogBytes)
	m.metricExists(status, "wiredTiger.transaction", &chartWiredTigerTransactions)
}

func (m *Mongo) metricExists(serverStatus map[string]interface{}, key string, chart *module.Chart) {
	keys := strings.Split(key, ".")
	for _, k := range keys {
		val, ok := serverStatus[k]
		if !ok {
			return
		}
		switch val.(type) {
		case map[string]interface{}:
			serverStatus = val.(map[string]interface{})
		default:
			return
		}
	}
	if enabled, ok := m.optionalChartsEnabled[chart.ID]; !ok || !enabled {
		err := m.charts.Add(chart.Copy())
		if err != nil {
			m.Warning(err)
		}
		m.optionalChartsEnabled[chart.ID] = true
		return
	}

}

func iterateServerStatus(ms map[string]int64, status map[string]interface{}) {
	for _, chart := range serverStatusCharts {
		for _, dim := range chart.Dims {
			addIfExists(status, dim.ID, ms)
		}
	}
}

func addIfExists(serverStatus map[string]interface{}, id string, ms map[string]int64) {
	mMap := serverStatus
	keys := strings.Split(id, ".")
	for _, key := range keys {
		k := fromID(key)
		val, ok := mMap[k]
		if !ok {
			return
		}
		switch t := val.(type) {
		case map[string]interface{}:
			mMap = val.(map[string]interface{})
		case int64:
			if _, ok := mMap[fromID(key)]; ok {
				ms[id] = t
			}
		case int32:
			if _, ok := mMap[fromID(key)]; ok {
				ms[id] = int64(t)
			}
		default:
			return
		}
	}
}
