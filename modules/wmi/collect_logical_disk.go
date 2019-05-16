package wmi

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricLDReadBytesTotal  = "wmi_logical_disk_read_bytes_total"
	metricLDWriteBytesTotal = "wmi_logical_disk_write_bytes_total"
	metricLDReadsTotal      = "wmi_logical_disk_reads_total"
	metricLDWritesTotal     = "wmi_logical_disk_writes_total"
	metricLDSizeBytes       = "wmi_logical_disk_size_bytes"
	metricLDFreeBytes       = "wmi_logical_disk_free_bytes"
)

func (w *WMI) collectLogicalDisk(mx *metrics, pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorLogDisks)
	if !(enabled && success) {
		return false
	}
	mx.LogicalDisk = &logicalDisk{}

	names := []string{
		metricLDReadBytesTotal,
		metricLDWriteBytesTotal,
		metricLDReadsTotal,
		metricLDWritesTotal,
		metricLDSizeBytes,
		metricLDFreeBytes,
	}

	for _, name := range names {
		collectLogicalDiskAny(mx, pms, name)
	}

	for _, v := range mx.LogicalDisk.Volumes {
		v.UsedSpace = sum(v.TotalSpace, -v.FreeSpace)
	}

	return true
}

func collectLogicalDiskAny(mx *metrics, pms prometheus.Metrics, name string) {
	vol := newVolume("")

	for _, pm := range pms.FindByName(name) {
		var (
			volumeID = pm.Labels.Get("volume")
			value    = pm.Value
		)
		if volumeID == "" || strings.HasPrefix(volumeID, "HarddiskVolume") {
			continue
		}
		if vol.ID != volumeID {
			vol = mx.LogicalDisk.Volumes.get(volumeID, true)
		}
		switch name {
		default:
			panic(fmt.Sprintf("unknown metric name during logical disk collection : %s", name))
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
		}
	}
}
