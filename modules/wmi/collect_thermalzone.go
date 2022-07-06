// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"strings"
)

const (
	collectorThermalzone = "thermalzone"

	metricThermalzoneTemperatureCelsius = "windows_thermalzone_temperature_celsius"
)

func doCollectThermalzone(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorThermalzone)
	return enabled && success
}

func collectThermalzone(pms prometheus.Metrics) *thermalZoneMetrics {
	if !doCollectThermalzone(pms) {
		return nil
	}

	tzm := &thermalZoneMetrics{}
	var tzone *thermalZone
	for _, pm := range pms.FindByName(metricThermalzoneTemperatureCelsius) {
		zoneName := cleanZoneName(pm.Labels.Get("name"))
		if zoneName == "" {
			continue
		}

		if tzone == nil || tzone.ID != zoneName {
			tzone = tzm.Zones.get(zoneName)
		}

		tzone.Temperature = pm.Value
	}
	return tzm
}

func cleanZoneName(name string) string {
	// "\\_TZ.TZ10", "\\_TZ.X570" => TZ10, X570
	i := strings.Index(name, ".")
	if i == -1 || len(name) == i+1 {
		return ""
	}
	return name[i+1:]
}
