// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"errors"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/prometheus/prometheus/model/labels"
)

func isValidWindowsExporterMetrics(pms prometheus.Metrics) bool {
	return pms.FindByName(metricCollectorSuccess).Len() > 0
}

func (w *WMI) collect() (map[string]int64, error) {
	pms, err := w.prom.Scrape()
	if err != nil {
		return nil, err
	}

	if !isValidWindowsExporterMetrics(pms) {
		return nil, errors.New("collected metrics aren't windows_exporter metrics")
	}

	mx := collect(pms)
	w.updateCharts(mx)

	return stm.ToMap(mx), nil
}

func collect(pms prometheus.Metrics) *metrics {
	mx := metrics{
		CPU:         collectCPU(pms),
		Memory:      collectMemory(pms),
		Net:         collectNet(pms),
		LogicalDisk: collectLogicalDisk(pms),
		OS:          collectOS(pms),
		System:      collectSystem(pms),
		Logon:       collectLogon(pms),
		ThermalZone: collectThermalzone(pms),
		Collectors:  collectCollection(pms),
	}

	if mx.hasOS() && mx.hasMemory() {
		v := mx.OS.VisibleMemoryBytes - mx.Memory.AvailableBytes
		mx.Memory.UsedBytes = &v
	}
	if mx.hasOS() {
		mx.OS.PagingUsedBytes = mx.OS.PagingLimitBytes - mx.OS.PagingFreeBytes
		mx.OS.VisibleMemoryUsedBytes = mx.OS.VisibleMemoryBytes - mx.OS.PhysicalMemoryFreeBytes
	}
	return &mx
}

func checkCollector(pms prometheus.Metrics, name string) (enabled, success bool) {
	m, err := labels.NewMatcher(labels.MatchEqual, "collector", name)
	if err != nil {
		panic(err)
	}

	pms = pms.FindByName(metricCollectorSuccess)
	ms := pms.Match(m)
	return ms.Len() > 0, ms.Max() == 1
}
