// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricADCSRequestotal                       = "windows_adcs_requests_total"
	metricADCSRequestProcessingTime             = "windows_adcs_request_processing_time_seconds"
	metricADCSRetrievalTotal                    = "windows_adcs_retrievals_total"
	metricADCSRetrievalProcessing               = "windows_adcs_retrievals_processing_time_seconds"
	metricADCSFailedRequest                     = "windows_adcs_failed_requests_total"
	metricADCSIssuedRequest                     = "windows_adcs_issued_requests_total"
	metricADCSPendingRequest                    = "windows_adcs_pending_requests_total"
	metricADCSReqCryptoSigning                  = "windows_adcs_request_cryptographic_signing_time_seconds"
	metricADCSPolicyModuleProcessing            = "windows_adcs_request_policy_module_processing_time_seconds"
	metricADCSChalengeResponse                  = "windows_adcs_challenge_responses_total"
	metricADCSChalengeResponseProcessing        = "windows_adcs_challenge_response_processing_time_seconds"
	metricADCSSignedCertTimestampList           = "windows_adcs_signed_certificate_timestamp_lists_total"
	metricADCSSignedCertTimestampListProcessing = "windows_adcs_signed_certificate_timestamp_list_processing_time_seconds"
)

func (w *WMI) collectADCS(mx map[string]int64, pms prometheus.Series) {
	seen := make(map[string]bool)
	px := "adcs_"
	for _, pm := range pms.FindByName(metricADCSRequestotal) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_cert_requests_total"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSRequestProcessingTime) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_cert_req_proc_time_elapsed"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSRetrievalTotal) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_cert_retrievals_total"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSRetrievalProcessing) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_cert_retrieval_time_elapsed"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSFailedRequest) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_cert_failed_request_total"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSIssuedRequest) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_cert_issued_request_total"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSPendingRequest) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_cert_pending_request_total"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSReqCryptoSigning) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_crypto_signing_time_elapsed"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSPolicyModuleProcessing) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_policy_mod_proc_time_elapsed"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSChalengeResponse) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_challenge_cert_response_total"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSChalengeResponseProcessing) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_challenge_response_proc_time_elapsed"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSSignedCertTimestampList) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_signed_cert_timestamp_list_total"] += int64(pm.Value)
		}
	}
	for _, pm := range pms.FindByName(metricADCSSignedCertTimestampListProcessing) {
		if name := pm.Labels.Get("cert_template"); name != "" && name != "_Total" {
			seen[name] = true
			mx[px+name+"_signed_cert_timestamp_list_proc_elapsed"] += int64(pm.Value)
		}
	}

	for template := range seen {
		if !w.cache.adcs[template] {
			w.cache.adcs[template] = true
			w.addTemplateCertCharts(template)
		}
	}
	for template := range w.cache.adcs {
		if !seen[template] {
			delete(w.cache.adcs, template)
			w.removeTemplateCertCharts(template)
		}
	}
}
