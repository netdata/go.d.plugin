package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
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
		zoneName := pm.Labels.Get("name")
		if zoneName == "" {
			continue
		}

		if zoneName = cleanZoneName(zoneName); tzone == nil || tzone.ID != zoneName {
			tzone = tzm.Zones.get(zoneName)
		}

		tzone.Temperature = pm.Value
	}
	return tzm
}

func cleanZoneName(name string) string {
	return name
}
