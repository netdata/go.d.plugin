// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

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

func (w *WMI) collectExchange(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorExchange] {
		w.cache.collection[collectorExchange] = true
		w.addExchangeCharts()
	}

	for _, pm := range pms.FindByNames(exchangeMetrics...) {
		name := strings.TrimPrefix(pm.Name(), "windows_")
		v := pm.Value
		if strings.HasSuffix(name, "_sec") {
			v *= precision
		}
		mx[name] = int64(v)
	}
}
