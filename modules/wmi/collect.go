// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"errors"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"net/url"
	"strings"
)

const precision = 1000

const (
	collectorCPU         = "cpu"
	collectorMemory      = "memory"
	collectorNet         = "net"
	collectorLogicalDisk = "logical_disk"
	collectorOS          = "os"
	collectorSystem      = "system"
	collectorLogon       = "logon"
	collectorThermalZone = "thermalzone"
	collectorTCP         = "tcp"
	collectorProcess     = "process"
	collectorService     = "service"
)

var fastCollectors = map[string]bool{
	collectorCPU:         true,
	collectorMemory:      true,
	collectorNet:         true,
	collectorLogicalDisk: true,
	collectorOS:          true,
	collectorSystem:      true,
	collectorTCP:         true,
}

var slowCollectors = map[string]bool{
	collectorLogon:       true,
	collectorThermalZone: true,
	collectorProcess:     true,
	collectorService:     true,
}

func (w *WMI) collect() (map[string]int64, error) {
	if w.doCheck {
		if err := w.checkSupportedCollectors(); err != nil {
			return nil, err
		}
		w.doCheck = false
	}

	mx := make(map[string]int64)

	if err := w.collectMetrics(mx, w.promFast); err != nil {
		if !strings.Contains(err.Error(), "unavailable collector") {
			return nil, err
		}
		w.doCheck = true
	}

	// TODO: charts with different update_every required
	if err := w.collectMetrics(mx, w.promSlow); err != nil {
		if !strings.Contains(err.Error(), "unavailable collector") {
			return nil, err
		}
		w.doCheck = true
	}

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

func (w *WMI) collectMetrics(mx map[string]int64, prom prometheus.Prometheus) error {
	if prom == nil {
		return nil
	}

	pms, err := prom.Scrape()
	if err != nil {
		return err
	}

	if pms.Len() == 0 {
		return nil
	}

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
		}
	}

	return nil
}

func (w *WMI) checkSupportedCollectors() error {
	w.promFast, w.promSlow = nil, nil

	pms, err := w.promCheck.Scrape()
	if err != nil {
		return err
	}

	if pms = pms.FindByName(metricCollectorSuccess); pms.Len() == 0 {
		return errors.New("collected metrics aren't windows_exporter metrics")
	}

	seen := make(map[string]bool)
	var fast, slow []string
	for _, pm := range pms {
		name := pm.Labels.Get("collector")
		switch {
		case name == collectorThermalZone && pm.Value == 0:
		case fastCollectors[name]:
			seen[name] = true
			fast = append(fast, name)
		case slowCollectors[name]:
			seen[name] = true
			slow = append(slow, name)
		}
	}

	if len(seen) == 0 {
		return errors.New("no supported collectors found")
	}

	for name := range seen {
		if !w.cache.collectors[name] {
			w.cache.collectors[name] = true
			w.addCollectorCharts(name)
		}
	}
	for name := range w.cache.collectors {
		if !seen[name] {
			delete(w.cache.collectors, name)
			w.removeCollectorCharts(name)
		}
	}

	req := w.Request.Copy()
	u, err := url.Parse(req.URL)
	if err != nil {
		return err
	}

	if len(fast) > 0 {
		u.RawQuery = url.Values{"collect[]": fast}.Encode()
		req.URL = u.String()
		w.promFast = prometheus.New(w.httpClient, req.Copy())
	}
	if len(slow) > 0 {
		u.RawQuery = url.Values{"collect[]": slow}.Encode()
		req.URL = u.String()
		w.promSlow = prometheus.New(w.httpClient, req.Copy())
	}

	return nil
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
