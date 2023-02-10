// SPDX-License-Identifier: GPL-3.0-or-later

package windows

import (
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricActiveSyncPingCmdsPending  = "windows_exchange_activesync_ping_cmds_pending"
	metricActiveSyncRequestsTotal    = "windows_exchange_activesync_requests_total"
	metricActiveSyncCMDsTotal        = "windows_exchange_activesync_sync_cmds_total"
	metricAutoDiscoverRequestsTotal  = "windows_exchange_autodiscover_requests_total"
	metricAvailServiceRequestsPerSec = "windows_exchange_avail_service_requests_per_sec"
	metricOWACurrentUniqueUsers      = "windows_exchange_owa_current_unique_users"
	metricOWARequestsTotal           = "windows_exchange_owa_requests_total"
	metricRPCActiveUserCount         = "windows_exchange_rpc_active_user_count"
	metricRPCAvgLatencySec           = "windows_exchange_rpc_avg_latency_sec"
	metricRPCConnectionCount         = "windows_exchange_rpc_connection_count"
	metricRPCOperationsTotal         = "windows_exchange_rpc_operations_total"
	metricRPCRequests                = "windows_exchange_rpc_requests"
	metricRPCUserCount               = "windows_exchange_rpc_user_count"

	metricTransportQueuesActiveMailboxDelivery         = "windows_exchange_transport_queues_active_mailbox_delivery"
	metricTransportQueuesExternanlActiveRemoteDelivery = "windows_exchange_transport_queues_external_active_remote_delivery"
	metricTransportQueuesExternalLargestDelivery       = "windows_exchange_transport_queues_external_largest_delivery"
	metricTransportQueuesInternalActiveRemoteDelivery  = "windows_exchange_transport_queues_internal_active_remote_delivery"
	metricTransportQueuesInternalLargestDelivery       = "windows_exchange_transport_queues_internal_largest_delivery"
	metricTransportQueuesPoison                        = "windows_exchange_transport_queues_poison"
	metricTransportQueuesRetryMailboxDelivery          = "windows_exchange_transport_queues_retry_mailbox_delivery"
	metricTransportQueuesUnreachable                   = "windows_exchange_transport_queues_unreachable"
)

var exchangeMetrics = []string{
	metricActiveSyncPingCmdsPending,
	metricActiveSyncRequestsTotal,
	metricActiveSyncCMDsTotal,
	metricAutoDiscoverRequestsTotal,
	metricAvailServiceRequestsPerSec,
	metricOWACurrentUniqueUsers,
	metricOWARequestsTotal,
	metricRPCActiveUserCount,
	metricRPCAvgLatencySec,
	metricRPCConnectionCount,
	metricRPCOperationsTotal,
	metricRPCRequests,
	metricRPCUserCount,
}

func (w *Windows) collectExchange(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorExchange] {
		w.cache.collection[collectorExchange] = true
		w.addExchangeCharts()
	}

	if pm := pms.FindByName(metricActiveSyncPingCmdsPending); pm.Len() > 0 {
		mx["exchange_activesync_ping_cmds_pending"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricActiveSyncRequestsTotal); pm.Len() > 0 {
		mx["exchange_activesync_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricActiveSyncCMDsTotal); pm.Len() > 0 {
		mx["exchange_activesync_sync_cmds_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricAutoDiscoverRequestsTotal); pm.Len() > 0 {
		mx["exchange_autodiscover_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricAvailServiceRequestsPerSec); pm.Len() > 0 {
		mx["exchange_avail_service_requests_per_sec"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricOWACurrentUniqueUsers); pm.Len() > 0 {
		mx["exchange_owa_current_unique_users"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricOWARequestsTotal); pm.Len() > 0 {
		mx["exchange_owa_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricRPCActiveUserCount); pm.Len() > 0 {
		mx["exchange_rpc_active_user_count"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricRPCAvgLatencySec); pm.Len() > 0 {
		mx["exchange_rpc_avg_latency_sec"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricRPCConnectionCount); pm.Len() > 0 {
		mx["exchange_rpc_connection_count"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricRPCOperationsTotal); pm.Len() > 0 {
		mx["exchange_rpc_operations_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricRPCRequests); pm.Len() > 0 {
		mx["exchange_rpc_requests"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricRPCUserCount); pm.Len() > 0 {
		mx["exchange_rpc_user_count"] = int64(pm.Max())
	}

	exchangeAddTransportQueueMetric(mx, pms)
}

func exchangeAddTransportQueueMetric(mx map[string]int64, pms prometheus.Series) {
	pms = pms.FindByNames(
		metricTransportQueuesActiveMailboxDelivery,
		metricTransportQueuesExternanlActiveRemoteDelivery,
		metricTransportQueuesExternalLargestDelivery,
		metricTransportQueuesInternalActiveRemoteDelivery,
		metricTransportQueuesInternalLargestDelivery,
		metricTransportQueuesPoison,
		metricTransportQueuesRetryMailboxDelivery,
		metricTransportQueuesUnreachable,
	)

	for _, pm := range pms {
		if name := pm.Labels.Get("name"); name != "" && name != "total_excluding_priority_none" {
			metric := strings.TrimPrefix(pm.Name(), "windows_")
			v := pm.Value
			mx[metric+"_"+name] += int64(v)
		}
	}
}
