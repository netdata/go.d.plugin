package wmi

import (
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorLogDisks = "logical_disk"

	metricLDReadBytesTotal    = "windows_logical_disk_read_bytes_total"
	metricLDWriteBytesTotal   = "windows_logical_disk_write_bytes_total"
	metricLDReadsTotal        = "windows_logical_disk_reads_total"
	metricLDWritesTotal       = "windows_logical_disk_writes_total"
	metricLDSizeBytes         = "windows_logical_disk_size_bytes"
	metricLDFreeBytes         = "windows_logical_disk_free_bytes"
	metricLDReadLatencyTotal  = "windows_logical_disk_read_latency_seconds_total"
	metricLDWriteLatencyTotal = "windows_logical_disk_write_latency_seconds_total"
)

var ldMetricNames = []string{
	metricLDReadBytesTotal,
	metricLDWriteBytesTotal,
	metricLDReadsTotal,
	metricLDWritesTotal,
	metricLDSizeBytes,
	metricLDFreeBytes,
	metricLDReadLatencyTotal,
	metricLDWriteLatencyTotal,
}

func doCollectLogicalDisk(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorLogDisks)
	return enabled && success
}

func collectLogicalDisk(pms prometheus.Metrics) *logicalDiskMetrics {
	if !doCollectLogicalDisk(pms) {
		return nil
	}

	dm := &logicalDiskMetrics{}
	for _, name := range ldMetricNames {
		collectLogicalDiskMetric(dm, pms, name)
	}

	for _, v := range dm.Volumes {
		v.UsedSpace = v.TotalSpace - v.FreeSpace
	}
	return dm
}

func collectLogicalDiskMetric(dm *logicalDiskMetrics, pms prometheus.Metrics, name string) {
	var vol *volume

	for _, pm := range pms.FindByName(name) {
		volumeID := pm.Labels.Get("volume")
		if volumeID == "" || strings.HasPrefix(volumeID, "HarddiskVolume") {
			continue
		}

		if vol == nil || vol.ID != volumeID {
			vol = dm.Volumes.get(volumeID)
		}

		assignVolumeMetric(vol, name, pm.Value)
	}
}

func assignVolumeMetric(vol *volume, name string, value float64) {
	switch name {
	case metricLDReadBytesTotal:
		vol.ReadBytesTotal = value
	case metricLDWriteBytesTotal:
		vol.WriteBytesTotal = value
	case metricLDReadsTotal:
		vol.ReadsTotal = value
	case metricLDWritesTotal:
		vol.WritesTotal = value
	case metricLDSizeBytes:
		vol.TotalSpace = value
	case metricLDFreeBytes:
		vol.FreeSpace = value
	case metricLDReadLatencyTotal:
		vol.ReadLatency = value
	case metricLDWriteLatencyTotal:
		vol.WriteLatency = value
	}
}
