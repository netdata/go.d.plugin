// SPDX-License-Identifier: GPL-3.0-or-later

package windows

import (
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricExchangeActiveSyncPingCmdsPending  = "windows_exchange_activesync_ping_cmds_pending"
	metricExchangeActiveSyncRequestsTotal    = "windows_exchange_activesync_requests_total"
	metricExchangeActiveSyncCMDsTotal        = "windows_exchange_activesync_sync_cmds_total"
	metricExchangeAutoDiscoverRequestsTotal  = "windows_exchange_autodiscover_requests_total"
	metricExchangeAvailServiceRequestsPerSec = "windows_exchange_avail_service_requests_per_sec"
	metricExchangeOWACurrentUniqueUsers      = "windows_exchange_owa_current_unique_users"
	metricExchangeOWARequestsTotal           = "windows_exchange_owa_requests_total"
	metricExchangeRPCActiveUserCount         = "windows_exchange_rpc_active_user_count"
	metricExchangeRPCAvgLatencySec           = "windows_exchange_rpc_avg_latency_sec"
	metricExchangeRPCConnectionCount         = "windows_exchange_rpc_connection_count"
	metricExchangeRPCOperationsTotal         = "windows_exchange_rpc_operations_total"
	metricExchangeRPCRequests                = "windows_exchange_rpc_requests"
	metricExchangeRPCUserCount               = "windows_exchange_rpc_user_count"

	metricExchangeTransportQueuesActiveMailboxDelivery         = "windows_exchange_transport_queues_active_mailbox_delivery"
	metricExchangeTransportQueuesExternanlActiveRemoteDelivery = "windows_exchange_transport_queues_external_active_remote_delivery"
	metricExchangeTransportQueuesExternalLargestDelivery       = "windows_exchange_transport_queues_external_largest_delivery"
	metricExchangeTransportQueuesInternalActiveRemoteDelivery  = "windows_exchange_transport_queues_internal_active_remote_delivery"
	metricExchangeTransportQueuesInternalLargestDelivery       = "windows_exchange_transport_queues_internal_largest_delivery"
	metricExchangeTransportQueuesPoison                        = "windows_exchange_transport_queues_poison"
	metricExchangeTransportQueuesRetryMailboxDelivery          = "windows_exchange_transport_queues_retry_mailbox_delivery"
	metricExchangeTransportQueuesUnreachable                   = "windows_exchange_transport_queues_unreachable"

	metricExchangeWorkloadActiveTasks    = "windows_exchange_workload_active_tasks"
	metricExchangeWorkloadCompletedTasks = "windows_exchange_workload_completed_tasks"
	metricExchangeWorkloadIsActive       = "windows_exchange_workload_is_active"
	metricExchangeWorkloadQueuedTasks    = "windows_exchange_workload_queued_tasks"
	metricExchangeWorkloadYeldedTasks    = "windows_exchange_workload_yielded_tasks"
)

var exchangeMetrics = []string{
	metricExchangeActiveSyncPingCmdsPending,
	metricExchangeActiveSyncRequestsTotal,
	metricExchangeActiveSyncCMDsTotal,
	metricExchangeAutoDiscoverRequestsTotal,
	metricExchangeAvailServiceRequestsPerSec,
	metricExchangeOWACurrentUniqueUsers,
	metricExchangeOWARequestsTotal,
	metricExchangeRPCActiveUserCount,
	metricExchangeRPCAvgLatencySec,
	metricExchangeRPCConnectionCount,
	metricExchangeRPCOperationsTotal,
	metricExchangeRPCRequests,
	metricExchangeRPCUserCount,
}

func (w *Windows) collectExchange(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorExchange] {
		w.cache.collection[collectorExchange] = true
		w.addExchangeCharts()
	}

	if pm := pms.FindByName(metricExchangeActiveSyncPingCmdsPending); pm.Len() > 0 {
		mx["exchange_activesync_ping_cmds_pending"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeActiveSyncRequestsTotal); pm.Len() > 0 {
		mx["exchange_activesync_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeActiveSyncCMDsTotal); pm.Len() > 0 {
		mx["exchange_activesync_sync_cmds_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeAutoDiscoverRequestsTotal); pm.Len() > 0 {
		mx["exchange_autodiscover_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeAvailServiceRequestsPerSec); pm.Len() > 0 {
		mx["exchange_avail_service_requests_per_sec"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeOWACurrentUniqueUsers); pm.Len() > 0 {
		mx["exchange_owa_current_unique_users"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeOWARequestsTotal); pm.Len() > 0 {
		mx["exchange_owa_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeRPCActiveUserCount); pm.Len() > 0 {
		mx["exchange_rpc_active_user_count"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeRPCAvgLatencySec); pm.Len() > 0 {
		mx["exchange_rpc_avg_latency_sec"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeRPCConnectionCount); pm.Len() > 0 {
		mx["exchange_rpc_connection_count"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeRPCOperationsTotal); pm.Len() > 0 {
		mx["exchange_rpc_operations_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeRPCRequests); pm.Len() > 0 {
		mx["exchange_rpc_requests"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricExchangeRPCUserCount); pm.Len() > 0 {
		mx["exchange_rpc_user_count"] = int64(pm.Max())
	}

	exchangeAddTransportQueueMetric(mx, pms)
	exchangeAddWorkloadMetric(mx, pms, w)
}

func exchangeAddTransportQueueMetric(mx map[string]int64, pms prometheus.Series) {
	pms = pms.FindByNames(
		metricExchangeTransportQueuesActiveMailboxDelivery,
		metricExchangeTransportQueuesExternanlActiveRemoteDelivery,
		metricExchangeTransportQueuesExternalLargestDelivery,
		metricExchangeTransportQueuesInternalActiveRemoteDelivery,
		metricExchangeTransportQueuesInternalLargestDelivery,
		metricExchangeTransportQueuesPoison,
		metricExchangeTransportQueuesRetryMailboxDelivery,
		metricExchangeTransportQueuesUnreachable,
	)

	for _, pm := range pms {
		if name := pm.Labels.Get("name"); name != "" && name != "total_excluding_priority_none" {
			metric := strings.TrimPrefix(pm.Name(), "windows_")
			v := pm.Value
			mx[metric+"_"+name] += int64(v)
		}
	}
}

func exchangeAddWorkloadMetric(mx map[string]int64, pms prometheus.Series, w *Windows) {
	pms = pms.FindByNames(
		metricExchangeWorkloadActiveTasks,
		metricExchangeWorkloadCompletedTasks,
		metricExchangeWorkloadIsActive,
		metricExchangeWorkloadQueuedTasks,
		metricExchangeWorkloadYeldedTasks,
	)
	seen := make(map[string]bool)

	for _, pm := range pms {
		if name := pm.Labels.Get("name"); name != "" {
			seen[name] = true
			metric := strings.TrimPrefix(pm.Name(), "windows_exchange_workload_")
			v := pm.Value
			mx["exchange_workload_"+name+"_"+metric] += int64(v)
		}
	}

	for name := range seen {
		if !w.cache.exchangeWorkload[name] {
			w.cache.exchangeWorkload[name] = true
			w.addExchangeWorkloadCharts(name)
		}
	}
	for name := range w.cache.exchangeWorkload {
		if !seen[name] {
			delete(w.cache.exchangeWorkload, name)
			w.removeCertificateTemplateCharts(name)
		}
	}
}
