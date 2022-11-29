// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const precision = 1000

const (
	collectorAD          = "ad"
	collectorCPU         = "cpu"
	collectorMemory      = "memory"
	collectorNet         = "net"
	collectorLogicalDisk = "logical_disk"
	collectorOS          = "os"
	collectorSystem      = "system"
	collectorLogon       = "logon"
	collectorThermalZone = "thermalzone"
	collectorTCP         = "tcp"
	collectorIIS         = "iis"
	collectorMSSQL       = "mssql"
	collectorProcess     = "process"
	collectorService     = "service"
)

func (w *WMI) collect() (map[string]int64, error) {
	pms, err := w.prom.ScrapeSeries()
	if err != nil {
		return nil, err
	}

	mx := make(map[string]int64)
	w.collectMetrics(mx, pms)

	if hasKey(mx, "os_visible_memory_bytes", "memory_available_bytes") {
		mx["memory_used_bytes"] = 0 +
			mx["os_visible_memory_bytes"] -
			mx["memory_available_bytes"]
	}
	if hasKey(mx, "os_paging_limit_bytes", "os_paging_free_bytes") {
		mx["os_paging_used_bytes"] = 0 +
			mx["os_paging_limit_bytes"] -
			mx["os_paging_free_bytes"]
	}
	if hasKey(mx, "os_visible_memory_bytes", "os_physical_memory_free_bytes") {
		mx["os_visible_memory_used_bytes"] = 0 +
			mx["os_visible_memory_bytes"] -
			mx["os_physical_memory_free_bytes"]
	}
	if hasKey(mx, "memory_commit_limit", "memory_committed_bytes") {
		mx["memory_not_committed_bytes"] = 0 +
			mx["memory_commit_limit"] -
			mx["memory_committed_bytes"]
	}
	if hasKey(mx, "memory_standby_cache_reserve_bytes", "memory_standby_cache_normal_priority_bytes", "memory_standby_cache_core_bytes") {
		mx["memory_standby_cache_total"] = 0 +
			mx["memory_standby_cache_reserve_bytes"] +
			mx["memory_standby_cache_normal_priority_bytes"] +
			mx["memory_standby_cache_core_bytes"]
	}
	if hasKey(mx, "memory_standby_cache_total", "memory_modified_page_list_bytes") {
		mx["memory_cache_total"] = 0 +
			mx["memory_standby_cache_total"] +
			mx["memory_modified_page_list_bytes"]
	}

	return mx, nil
}

func (w *WMI) collectMetrics(mx map[string]int64, pms prometheus.Series) {
	w.collectCollector(mx, pms)
	for _, pm := range pms.FindByName(metricCollectorSuccess) {
		if pm.Value == 0 {
			continue
		}

		switch pm.Labels.Get("collector") {
		case collectorCPU:
			w.collectCPU(mx, pms)
		case collectorMemory:
			w.collectMemory(mx, pms)
		case collectorNet:
			w.collectNet(mx, pms)
		case collectorLogicalDisk:
			w.collectLogicalDisk(mx, pms)
		case collectorOS:
			w.collectOS(mx, pms)
		case collectorSystem:
			w.collectSystem(mx, pms)
		case collectorLogon:
			w.collectLogon(mx, pms)
		case collectorThermalZone:
			w.collectThermalzone(mx, pms)
		case collectorTCP:
			w.collectTCP(mx, pms)
		case collectorProcess:
			w.collectProcess(mx, pms)
		case collectorService:
			w.collectService(mx, pms)
		case collectorIIS:
			w.collectIIS(mx, pms)
		case collectorMSSQL:
			w.collectMSSQL(mx, pms)
		}
	}
}

func hasKey(mx map[string]int64, key string, keys ...string) bool {
	_, ok := mx[key]
	switch len(keys) {
	case 0:
		return ok
	default:
		return ok && hasKey(mx, keys[0], keys[1:]...)
	}
}
