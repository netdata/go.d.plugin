// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

// Windows exporter:
// https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.ad.md
// Microsoft:
// https://learn.microsoft.com/en-us/previous-versions/ms803980(v=msdn.10)
const (
	metricADReplicationInboundObjectsFilteringTotal   = "windows_ad_replication_inbound_objects_filtered_total"
	metricADReplicationInboundPropertiesFilteredTotal = "windows_ad_replication_inbound_properties_filtered_total"
	metricADReplicationInboundPropertiesUpdatedTotal  = "windows_ad_replication_inbound_properties_updated_total"
	metricADReplicationInboundSyncObjectsRemaining    = "windows_ad_replication_inbound_sync_objects_remaining"
	metricADReplicationDataInterSiteBytesTotal        = "windows_ad_replication_data_intersite_bytes_total"
	metricADReplicationDataIntraSiteBytesTotal        = "windows_ad_replication_data_intrasite_bytes_total"
	metricADReplicationPendingSyncs                   = "windows_ad_replication_pending_synchronizations"
	metricADReplicationSyncRequestsTotal              = "windows_ad_replication_sync_requests_total"
	metricADDirectoryServiceThreads                   = "windows_ad_directory_service_threads"
	metricADLDAPLastBindTimeSecondsTotal              = "windows_ad_ldap_last_bind_time_seconds"
	metricADBindsTotal                                = "windows_ad_binds_total"
	metricADLDAPSearchesTotal                         = "windows_ad_ldap_searches_total"
)

func (w *WMI) collectAD(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorAD] {
		w.cache.collection[collectorAD] = true
		w.addADCharts()
	}

	if pm := pms.FindByName(metricADReplicationInboundObjectsFilteringTotal); pm.Len() > 0 {
		mx["ad_replication_inbound_objects_filtered_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADReplicationInboundPropertiesFilteredTotal); pm.Len() > 0 {
		mx["ad_replication_inbound_properties_filtered_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADReplicationInboundPropertiesUpdatedTotal); pm.Len() > 0 {
		mx["ad_replication_inbound_properties_updated_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADReplicationInboundSyncObjectsRemaining); pm.Len() > 0 {
		mx["ad_replication_inbound_sync_objects_remaining"] = int64(pm.Max())
	}
	for _, pm := range pms.FindByName(metricADReplicationDataInterSiteBytesTotal) {
		if name := pm.Labels.Get("direction"); name != "" {
			mx["ad_replication_data_intersite_bytes_total_"+name] = int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADReplicationDataIntraSiteBytesTotal) {
		if name := pm.Labels.Get("direction"); name != "" {
			mx["ad_replication_data_intrasite_bytes_total_"+name] = int64(pm.Value)
		}
	}
	if pm := pms.FindByName(metricADReplicationPendingSyncs); pm.Len() > 0 {
		mx["ad_replication_pending_synchronizations"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADReplicationSyncRequestsTotal); pm.Len() > 0 {
		mx["ad_replication_sync_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADDirectoryServiceThreads); pm.Len() > 0 {
		mx["ad_directory_service_threads"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADLDAPLastBindTimeSecondsTotal); pm.Len() > 0 {
		mx["ad_ldap_last_bind_time_seconds"] = int64(pm.Max())
	}
	for _, pm := range pms.FindByName(metricADBindsTotal) {
		mx["ad_binds_total"] += int64(pm.Value) // sum "bind_method"'s
	}
	if pm := pms.FindByName(metricADLDAPSearchesTotal); pm.Len() > 0 {
		mx["ad_ldap_searches_total"] = int64(pm.Max())
	}
}
