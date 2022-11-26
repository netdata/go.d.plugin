// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricADReplicationIncomingObjectsFiltering    = "windows_ad_replication_inbound_objects_filtered_total"
	metricADReplicationInboundPropFilteredTotal    = "windows_ad_replication_inbound_properties_filtered_total"
	metricADReplicationInboundPropUpdatedTotal     = "windows_ad_replication_inbound_properties_updated_total"
	metricADReplicationInboundObjectsUpdateTotal   = "windows_ad_replication_inbound_objects_updated_total"
	metricADReplicationInboundSyncObjectsReamining = "windows_ad_replication_inbound_sync_objects_remaining"
	metricADReplicationDataIntersiteBytesTotal     = "windows_ad_replication_data_intersite_bytes_total"
	metricADReplicationDataIntrasiteBytesTotal     = "windows_ad_replication_data_intrasite_bytes_total"
	metricADReplicationPendingSync                 = "windows_ad_replication_pending_synchronizations"
	metricADReplicationSyncRequestsTotal           = "windows_ad_replication_sync_requests_total"
	metricADDSThreads                              = "windows_ad_directory_service_threads"
	metricADLDAPLastBindTimeTotal                  = "windows_ad_ldap_last_bind_time_seconds"
	metricADBindTotal                              = "windows_ad_binds_total"
	metricADLDAPSearchTotal                        = "windows_ad_ldap_searches_total"
)

func (w *WMI) collectAD(mx map[string]int64, pms prometheus.Metrics) {
	px := "ad_dra_"
	if !w.cache.collection[collectorAD] {
		w.cache.collection[collectorAD] = true
		w.addADCharts()
	}

	if pm := pms.FindByName(metricADReplicationIncomingObjectsFiltering); pm.Len() > 0 {
		mx["ad_dra_objects_filtered"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADReplicationInboundPropFilteredTotal); pm.Len() > 0 {
		mx["ad_dra_properties_applied"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADReplicationInboundPropUpdatedTotal); pm.Len() > 0 {
		mx["ad_dra_properties_filtered"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADReplicationInboundSyncObjectsReamining); pm.Len() > 0 {
		mx["ad_dra_objects_remaining"] = int64(pm.Max())
	}
	for _, pm := range pms.FindByName(metricADReplicationDataIntersiteBytesTotal) {
		if name := pm.Labels.Get("direction"); name != "" {
			mx[px+"compressed_bandwidth_"+name] = int64(pm.Value)
		}
	}
	if pm := pms.FindByName(metricADReplicationDataIntrasiteBytesTotal); pm.Len() > 0 {
		mx["ad_dra_uncompressed_inbound"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADReplicationPendingSync); pm.Len() > 0 {
		mx["ad_dra_pending_sync"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADReplicationSyncRequestsTotal); pm.Len() > 0 {
		mx["ad_dra_sync_req_made"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADDSThreads); pm.Len() > 0 {
		mx["ad_ds_thread"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADLDAPLastBindTimeTotal); pm.Len() > 0 {
		mx["ad_ldap_bind_time"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADBindTotal); pm.Len() > 0 {
		mx["ad_ldap_bind_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADLDAPSearchTotal); pm.Len() > 0 {
		mx["ad_ldap_bind_searches"] = int64(pm.Max())
	}
}
