// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricADFSADLoginConnectionFailures      = "windows_adfs_ad_login_connection_failures_total"
	metricADFSCertificateAuthentications     = "windows_adfs_certificate_authentications_total"
	metricADFSDBArtifactFailure              = "windows_adfs_db_artifact_failure_total"
	metricADFSDBArtifactQueryTimeSeconds     = "windows_adfs_db_artifact_query_time_seconds_total"
	metricADFSDBConfigFalure                 = "windows_adfs_db_config_failure_total"
	metricADFSDBQueryTimeSeconds             = "windows_adfs_db_config_query_time_seconds_total"
	metricADFSDeviceAuthentications          = "windows_adfs_device_authentications_total"
	metricADFSExternalAuthenticationsFailure = "windows_adfs_external_authentications_failure_total"
	metricADFSExternalAuthenticationsSuccess = "windows_adfs_external_authentications_success_total"
	metricADFSExtranetAccountLockouts        = "windows_adfs_extranet_account_lockouts_total"
	metricADFSFederatedAuthentications       = "windows_adfs_federated_authentications_total"
	metricADFSFederationMetadataRequests     = "windows_adfs_federation_metadata_requests_total"

	metricADFSOauthAuthorizationRequests                       = "windows_adfs_oauth_authorization_requests_total"
	metricADFSOauthClientAuthenticationFailure                 = "windows_adfs_oauth_client_authentication_failure_total"
	metricADFSOauthClientAuthenticationSuccess                 = "windows_adfs_oauth_client_authentication_success_total"
	metricADFSOauthClientCredentialsFailure                    = "windows_adfs_oauth_client_credentials_failure_total"
	metricADFSOauthClientCredentialsSuccess                    = "windows_adfs_oauth_client_credentials_success_total"
	metricADFSOauthClientPrivKeyJTWAuthenticationFailure       = "windows_adfs_oauth_client_privkey_jtw_authentication_failure_total"
	metricADFSOauthClientPrivKeyJTWAuthenticationSuccess       = "windows_adfs_oauth_client_privkey_jwt_authentications_success_total"
	metricADFSOauthClientSecretBasicAuthenticationsFailure     = "windows_adfs_oauth_client_secret_basic_authentications_failure_total"
	metricADFSADFSOauthClientSecretBasicAuthenticationsSuccess = "windows_adfs_oauth_client_secret_basic_authentications_success_total"
	metricADFSOauthClientSecretPostAuthenticationsFailure      = "windows_adfs_oauth_client_secret_post_authentications_failure_total"
	metricADFSOauthClientSecretPostAuthenticationsSuccess      = "windows_adfs_oauth_client_secret_post_authentications_success_total"
	metricADFSOauthClientWindowsAuthenticationsFailure         = "windows_adfs_oauth_client_windows_authentications_failure_total"
	metricADFSOauthClientWindowsAuthenticationsSuccess         = "windows_adfs_oauth_client_windows_authentications_success_total"
	metricADFSOauthLogonCertificateRequestsFailure             = "windows_adfs_oauth_logon_certificate_requests_failure_total"
	metricADFSOauthLogonCertificateTokenRequestsSuccess        = "windows_adfs_oauth_logon_certificate_token_requests_success_total"
	metricADFSOauthPasswordGrantRequestsFailure                = "windows_adfs_oauth_password_grant_requests_failure_total"
	metricADFSOauthPasswordGrantRequestsSuccess                = "windows_adfs_oauth_password_grant_requests_success_total"
	metricADFSOauthTokenRequestsSuccess                        = "windows_adfs_oauth_token_requests_success_total"

	metricADFSPassiveRequest                     = "windows_adfs_passive_requests_total"
	metricADFSPasswortAuthentications            = "windows_adfs_passport_authentications_total"
	metricADFSPasswordChangeFailed               = "windows_adfs_password_change_failed_total"
	metricADFSWPasswordChangeSucceeded           = "windows_adfs_password_change_succeeded_total"
	metricADFSSamlpToeknRequestsSuccess          = "windows_adfs_samlp_token_requests_success_total"
	metricADFSSSOAuthenticationsFailure          = "windows_adfs_sso_authentications_failure_total"
	metricADFSSSOAuthenticationsSuccess          = "windows_adfs_sso_authentications_success_total"
	metricADFSTokenRequests                      = "windows_adfs_token_requests_total"
	metricADFSUserPasswordAuthenticationsFailure = "windows_adfs_userpassword_authentications_failure_total"
	metricADFSUserPasswordAuthenticationsSuccess = "windows_adfs_userpassword_authentications_success_total"
	metricADFSWindowsIntegratedAuthentications   = "windows_adfs_windows_integrated_authentications_total"
	metricADFSWSFedTokenRequestsSuccess          = "windows_adfs_wsfed_token_requests_success_total"
	metricADFSWSTrustTokenRequestsSuccess        = "windows_adfs_wstrust_token_requests_success_total"
)

func (w *WMI) collectADFS(mx map[string]int64, pms prometheus.Series) {
	if !w.cache.collection[collectorADFS] {
		w.cache.collection[collectorADFS] = true
		w.addADFSCharts()
	}

	if pm := pms.FindByName(metricADFSADLoginConnectionFailures); pm.Len() > 0 {
		mx["adfs_ad_login_connection_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSCertificateAuthentications); pm.Len() > 0 {
		mx["adfs_certificate_authentications_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSDBArtifactFailure); pm.Len() > 0 {
		mx["adfs_db_artifact_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSDBArtifactQueryTimeSeconds); pm.Len() > 0 {
		mx["adfs_db_artifact_query_time_seconds_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSDBConfigFalure); pm.Len() > 0 {
		mx["adfs_db_config_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSDBQueryTimeSeconds); pm.Len() > 0 {
		mx["adfs_db_config_query_time_seconds_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSDeviceAuthentications); pm.Len() > 0 {
		mx["adfs_device_authentications_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSExternalAuthenticationsFailure); pm.Len() > 0 {
		mx["adfs_external_authentications_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSExternalAuthenticationsSuccess); pm.Len() > 0 {
		mx["adfs_external_authentications_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSExtranetAccountLockouts); pm.Len() > 0 {
		mx["adfs_federation_metadata_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSFederatedAuthentications); pm.Len() > 0 {
		mx["adfs_oauth_authorization_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSFederationMetadataRequests); pm.Len() > 0 {
		mx["adfs_federation_metadata_requests_total"] = int64(pm.Max())
	}

	if pm := pms.FindByName(metricADFSOauthAuthorizationRequests); pm.Len() > 0 {
		mx["adfs_oauth_authorization_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientAuthenticationFailure); pm.Len() > 0 {
		mx["adfs_oauth_client_authentication_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientAuthenticationSuccess); pm.Len() > 0 {
		mx["adfs_oauth_client_authentication_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientCredentialsFailure); pm.Len() > 0 {
		mx["adfs_oauth_client_credentials_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientCredentialsSuccess); pm.Len() > 0 {
		mx["adfs_oauth_client_credentials_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientPrivKeyJTWAuthenticationFailure); pm.Len() > 0 {
		mx["adfs_oauth_client_privkey_jwt_authentications_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientPrivKeyJTWAuthenticationSuccess); pm.Len() > 0 {
		mx["adfs_oauth_client_privkey_jwt_authentications_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientSecretBasicAuthenticationsFailure); pm.Len() > 0 {
		mx["adfs_oauth_client_secret_basic_authentications_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSADFSOauthClientSecretBasicAuthenticationsSuccess); pm.Len() > 0 {
		mx["adfs_oauth_client_secret_basic_authentications_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientSecretPostAuthenticationsFailure); pm.Len() > 0 {
		mx["adfs_oauth_client_secret_post_authentications_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientSecretPostAuthenticationsSuccess); pm.Len() > 0 {
		mx["adfs_oauth_client_secret_post_authentications_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientWindowsAuthenticationsFailure); pm.Len() > 0 {
		mx["adfs_oauth_client_windows_authentications_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthClientWindowsAuthenticationsSuccess); pm.Len() > 0 {
		mx["adfs_oauth_client_windows_authentications_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthLogonCertificateRequestsFailure); pm.Len() > 0 {
		mx["adfs_oauth_logon_certificate_requests_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthLogonCertificateTokenRequestsSuccess); pm.Len() > 0 {
		mx["adfs_oauth_logon_certificate_token_requests_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthPasswordGrantRequestsFailure); pm.Len() > 0 {
		mx["adfs_oauth_password_grant_requests_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthPasswordGrantRequestsSuccess); pm.Len() > 0 {
		mx["adfs_oauth_password_grant_requests_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSOauthTokenRequestsSuccess); pm.Len() > 0 {
		mx["adfs_oauth_token_requests_success_total"] = int64(pm.Max())
	}

	if pm := pms.FindByName(metricADFSPassiveRequest); pm.Len() > 0 {
		mx["adfs_passive_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSPasswortAuthentications); pm.Len() > 0 {
		mx["adfs_passport_authentications_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSPasswordChangeFailed); pm.Len() > 0 {
		mx["adfs_password_change_failed_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSWPasswordChangeSucceeded); pm.Len() > 0 {
		mx["adfs_password_change_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSSamlpToeknRequestsSuccess); pm.Len() > 0 {
		mx["adfs_samlp_token_requests_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSSSOAuthenticationsFailure); pm.Len() > 0 {
		mx["adfs_sso_authentications_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSSSOAuthenticationsSuccess); pm.Len() > 0 {
		mx["adfs_sso_authentications_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSTokenRequests); pm.Len() > 0 {
		mx["adfs_token_requests_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSUserPasswordAuthenticationsFailure); pm.Len() > 0 {
		mx["adfs_sso_authentications_failure_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSUserPasswordAuthenticationsSuccess); pm.Len() > 0 {
		mx["adfs_sso_authentications_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSWindowsIntegratedAuthentications); pm.Len() > 0 {
		mx["adfs_windows_integrated_authentications_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSWSFedTokenRequestsSuccess); pm.Len() > 0 {
		mx["adfs_wsfed_token_requests_success_total"] = int64(pm.Max())
	}
	if pm := pms.FindByName(metricADFSWSTrustTokenRequestsSuccess); pm.Len() > 0 {
		mx["adfs_wstrust_token_requests_success_total"] = int64(pm.Max())
	}
}
