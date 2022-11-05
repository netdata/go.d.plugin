// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricIISCurrentAnonReq    = "windows_iis_current_anonymous_users"
	metricIISCurrentNonAnonReq = "windows_iis_current_non_anonymous_users"
	metricIISActiveConn        = "windows_iis_current_connections"
	metricIISISAPIExtConn      = "windows_iis_current_connections"
	metricIISUptime            = "windows_iis_service_uptime"
	metricIISBandwidthRecv     = "windows_iis_received_bytes_total"
	metricIISBandwidthSent     = "windows_iis_sent_bytes_total"
	metricIISTotalReqAnon      = "windows_iis_anonymous_users_total"
	metricIISTotalReq          = "windows_iis_requests_total"
	metricIISConnAttemp        = "windows_iis_connection_attempts_all_instances_total"
	metricIISReqTotal          = "windows_iis_requests_total"
	metricIISFileRecv          = "windows_iis_files_received_total"
	metricIISFileSent          = "windows_iis_files_sent_total"
	metricIISExtensionReq      = "windows_iis_ipapi_extension_requests_total"
	metricIISLogon             = "windows_iis_logon_attempts_total"
	metricIISError423          = "windows_iis_locked_errors_total"
	metricIISError404          = "windows_iis_not_found_errors_total"
)

func (w *WMI) collectIIS(mx map[string]int64, pms prometheus.Metrics) {
	if !w.cache.collection[collectorIIS] {
		w.cache.collection[collectorIIS] = true
		w.addIISCharts()
	}

	seen := make(map[string]bool)
	ix := "iis_"
	for _, pm := range pms.FindByName(metricIISCurrentAnonReq) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_active_request_anon"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISCurrentNonAnonReq) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_active_request_non_anon"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISActiveConn) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_active_conn"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISISAPIExtConn) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_isapi_ext"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISUptime) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_uptime"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISBandwidthRecv) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_bandwidth_recv"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISBandwidthSent) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_bandwidth_sent"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISTotalReqAnon) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_total_req_anon"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISTotalReq) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_total_req"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISConnAttemp) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_conns_atemp"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISReqTotal) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_req_total"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISFileRecv) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_file_transfer_recv"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISFileSent) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_file_transfer_sent"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISExtensionReq) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_extension_req"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISLogon) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_logon"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISError423) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_error_423"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricIISError404) {
		if name := cleanSiteName(pm.Labels.Get("site")); name != "" {
			seen[name] = true
			mx[ix+name+"_iis_error_404"] += int64(pm.Value)
		}
	}

	for site := range seen {
		if !w.cache.processes[site] {
			w.cache.processes[site] = true
			w.addIISToCharts(site)
		}
	}
}

func cleanSiteName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}
