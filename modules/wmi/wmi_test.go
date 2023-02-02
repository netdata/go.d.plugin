// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v0200Metrics, _ = os.ReadFile("testdata/v0.20.0/metrics.txt")
)

func Test_TestData(t *testing.T) {
	for name, data := range map[string][]byte{
		"v0200Metrics": v0200Metrics,
	} {
		assert.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.IsType(t, (*WMI)(nil), New())
}

func TestWMI_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"success if 'url' is set": {
			config: Config{
				HTTP: web.HTTP{Request: web.Request{URL: "http://127.0.0.1:9182/metrics"}}},
		},
		"fails on default config": {
			wantFail: true,
			config:   New().Config,
		},
		"fails if 'url' is unset": {
			wantFail: true,
			config:   Config{HTTP: web.HTTP{Request: web.Request{URL: ""}}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wmi := New()
			wmi.Config = test.config

			if test.wantFail {
				assert.False(t, wmi.Init())
			} else {
				assert.True(t, wmi.Init())
			}
		})
	}
}

func TestWMI_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() (wmi *WMI, cleanup func())
		wantFail bool
	}{
		"success on valid response v0.20.0": {
			prepare: prepareWMIv0200,
		},
		"fails if endpoint returns invalid data": {
			wantFail: true,
			prepare:  prepareWMIReturnsInvalidData,
		},
		"fails on connection refused": {
			wantFail: true,
			prepare:  prepareWMIConnectionRefused,
		},
		"fails on 404 response": {
			wantFail: true,
			prepare:  prepareWMIResponse404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wmi, cleanup := test.prepare()
			defer cleanup()

			require.True(t, wmi.Init())

			if test.wantFail {
				assert.False(t, wmi.Check())
			} else {
				assert.True(t, wmi.Check())
			}
		})
	}
}

func TestWMI_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestWMI_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestWMI_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() (wmi *WMI, cleanup func())
		wantCollected map[string]int64
	}{
		"success on valid response v0.20.0": {
			prepare: prepareWMIv0200,
			wantCollected: map[string]int64{
				"ad_binds_total":                                                                                184,
				"ad_directory_service_threads":                                                                  0,
				"ad_ldap_last_bind_time_seconds":                                                                0,
				"ad_ldap_searches_total":                                                                        1382,
				"ad_replication_data_intersite_bytes_total_inbound":                                             0,
				"ad_replication_data_intersite_bytes_total_outbound":                                            0,
				"ad_replication_data_intrasite_bytes_total_inbound":                                             0,
				"ad_replication_data_intrasite_bytes_total_outbound":                                            0,
				"ad_replication_inbound_objects_filtered_total":                                                 0,
				"ad_replication_inbound_properties_filtered_total":                                              0,
				"ad_replication_inbound_properties_updated_total":                                               0,
				"ad_replication_inbound_sync_objects_remaining":                                                 0,
				"ad_replication_pending_synchronizations":                                                       0,
				"ad_replication_sync_requests_total":                                                            0,
				"adcs_cert_template_Administrator_challenge_response_processing_time_seconds":                   0,
				"adcs_cert_template_Administrator_challenge_responses_total":                                    0,
				"adcs_cert_template_Administrator_failed_requests_total":                                        0,
				"adcs_cert_template_Administrator_issued_requests_total":                                        0,
				"adcs_cert_template_Administrator_pending_requests_total":                                       0,
				"adcs_cert_template_Administrator_request_cryptographic_signing_time_seconds":                   0,
				"adcs_cert_template_Administrator_request_policy_module_processing_time_seconds":                0,
				"adcs_cert_template_Administrator_request_processing_time_seconds":                              0,
				"adcs_cert_template_Administrator_requests_total":                                               0,
				"adcs_cert_template_Administrator_retrievals_processing_time_seconds":                           0,
				"adcs_cert_template_Administrator_retrievals_total":                                             0,
				"adcs_cert_template_Administrator_signed_certificate_timestamp_list_processing_time_seconds":    0,
				"adcs_cert_template_Administrator_signed_certificate_timestamp_lists_total":                     0,
				"adcs_cert_template_DomainController_challenge_response_processing_time_seconds":                0,
				"adcs_cert_template_DomainController_challenge_responses_total":                                 0,
				"adcs_cert_template_DomainController_failed_requests_total":                                     0,
				"adcs_cert_template_DomainController_issued_requests_total":                                     1,
				"adcs_cert_template_DomainController_pending_requests_total":                                    0,
				"adcs_cert_template_DomainController_request_cryptographic_signing_time_seconds":                0,
				"adcs_cert_template_DomainController_request_policy_module_processing_time_seconds":             16,
				"adcs_cert_template_DomainController_request_processing_time_seconds":                           63,
				"adcs_cert_template_DomainController_requests_total":                                            1,
				"adcs_cert_template_DomainController_retrievals_processing_time_seconds":                        0,
				"adcs_cert_template_DomainController_retrievals_total":                                          0,
				"adcs_cert_template_DomainController_signed_certificate_timestamp_list_processing_time_seconds": 0,
				"adcs_cert_template_DomainController_signed_certificate_timestamp_lists_total":                  0,
				"adfs_ad_login_connection_failures_total":                                                       0,
				"adfs_certificate_authentications_total":                                                        0,
				"adfs_db_artifact_failure_total":                                                                0,
				"adfs_db_artifact_query_time_seconds_total":                                                     0,
				"adfs_db_config_failure_total":                                                                  0,
				"adfs_db_config_query_time_seconds_total":                                                       101,
				"adfs_device_authentications_total":                                                             0,
				"adfs_external_authentications_failure_total":                                                   0,
				"adfs_external_authentications_success_total":                                                   0,
				"adfs_extranet_account_lockouts_total":                                                          0,
				"adfs_federated_authentications_total":                                                          0,
				"adfs_federation_metadata_requests_total":                                                       1,
				"adfs_oauth_authorization_requests_total":                                                       0,
				"adfs_oauth_client_authentication_failure_total":                                                0,
				"adfs_oauth_client_authentication_success_total":                                                0,
				"adfs_oauth_client_credentials_failure_total":                                                   0,
				"adfs_oauth_client_credentials_success_total":                                                   0,
				"adfs_oauth_client_privkey_jtw_authentication_failure_total":                                    0,
				"adfs_oauth_client_privkey_jwt_authentications_success_total":                                   0,
				"adfs_oauth_client_secret_basic_authentications_failure_total":                                  0,
				"adfs_oauth_client_secret_basic_authentications_success_total":                                  0,
				"adfs_oauth_client_secret_post_authentications_failure_total":                                   0,
				"adfs_oauth_client_secret_post_authentications_success_total":                                   0,
				"adfs_oauth_client_windows_authentications_failure_total":                                       0,
				"adfs_oauth_client_windows_authentications_success_total":                                       0,
				"adfs_oauth_logon_certificate_requests_failure_total":                                           0,
				"adfs_oauth_logon_certificate_token_requests_success_total":                                     0,
				"adfs_oauth_password_grant_requests_failure_total":                                              0,
				"adfs_oauth_password_grant_requests_success_total":                                              0,
				"adfs_oauth_token_requests_success_total":                                                       0,
				"adfs_passive_requests_total":                                                                   0,
				"adfs_passport_authentications_total":                                                           0,
				"adfs_password_change_failed_total":                                                             0,
				"adfs_password_change_succeeded_total":                                                          0,
				"adfs_samlp_token_requests_success_total":                                                       0,
				"adfs_sso_authentications_failure_total":                                                        0,
				"adfs_sso_authentications_success_total":                                                        0,
				"adfs_token_requests_total":                                                                     0,
				"adfs_userpassword_authentications_failure_total":                                               0,
				"adfs_userpassword_authentications_success_total":                                               0,
				"adfs_windows_integrated_authentications_total":                                                 0,
				"adfs_wsfed_token_requests_success_total":                                                       0,
				"adfs_wstrust_token_requests_success_total":                                                     0,
				"collector_ad_duration":                                                                         769,
				"collector_ad_status_fail":                                                                      0,
				"collector_ad_status_success":                                                                   1,
				"collector_adcs_duration":                                                                       0,
				"collector_adcs_status_fail":                                                                    0,
				"collector_adcs_status_success":                                                                 1,
				"collector_adfs_duration":                                                                       3,
				"collector_adfs_status_fail":                                                                    0,
				"collector_adfs_status_success":                                                                 1,
				"collector_cpu_duration":                                                                        0,
				"collector_cpu_status_fail":                                                                     0,
				"collector_cpu_status_success":                                                                  1,
				"collector_iis_duration":                                                                        0,
				"collector_iis_status_fail":                                                                     0,
				"collector_iis_status_success":                                                                  1,
				"collector_logical_disk_duration":                                                               0,
				"collector_logical_disk_status_fail":                                                            0,
				"collector_logical_disk_status_success":                                                         1,
				"collector_logon_duration":                                                                      113,
				"collector_logon_status_fail":                                                                   0,
				"collector_logon_status_success":                                                                1,
				"collector_memory_duration":                                                                     0,
				"collector_memory_status_fail":                                                                  0,
				"collector_memory_status_success":                                                               1,
				"collector_mssql_duration":                                                                      3,
				"collector_mssql_status_fail":                                                                   0,
				"collector_mssql_status_success":                                                                1,
				"collector_net_duration":                                                                        0,
				"collector_net_status_fail":                                                                     0,
				"collector_net_status_success":                                                                  1,
				"collector_netframework_clrexceptions_duration":                                                 1437,
				"collector_netframework_clrexceptions_status_fail":                                              0,
				"collector_netframework_clrexceptions_status_success":                                           1,
				"collector_netframework_clrinterop_duration":                                                    1491,
				"collector_netframework_clrinterop_status_fail":                                                 0,
				"collector_netframework_clrinterop_status_success":                                              1,
				"collector_netframework_clrjit_duration":                                                        1278,
				"collector_netframework_clrjit_status_fail":                                                     0,
				"collector_netframework_clrjit_status_success":                                                  1,
				"collector_netframework_clrloading_duration":                                                    1323,
				"collector_netframework_clrloading_status_fail":                                                 0,
				"collector_netframework_clrloading_status_success":                                              1,
				"collector_netframework_clrlocksandthreads_duration":                                            1357,
				"collector_netframework_clrlocksandthreads_status_fail":                                         0,
				"collector_netframework_clrlocksandthreads_status_success":                                      1,
				"collector_netframework_clrmemory_duration":                                                     1406,
				"collector_netframework_clrmemory_status_fail":                                                  0,
				"collector_netframework_clrmemory_status_success":                                               1,
				"collector_netframework_clrremoting_duration":                                                   1519,
				"collector_netframework_clrremoting_status_fail":                                                0,
				"collector_netframework_clrremoting_status_success":                                             1,
				"collector_netframework_clrsecurity_duration":                                                   1467,
				"collector_netframework_clrsecurity_status_fail":                                                0,
				"collector_netframework_clrsecurity_status_success":                                             1,
				"collector_os_duration":                                                                         2,
				"collector_os_status_fail":                                                                      0,
				"collector_os_status_success":                                                                   1,
				"collector_process_duration":                                                                    115,
				"collector_process_status_fail":                                                                 0,
				"collector_process_status_success":                                                              1,
				"collector_service_duration":                                                                    101,
				"collector_service_status_fail":                                                                 0,
				"collector_service_status_success":                                                              1,
				"collector_system_duration":                                                                     0,
				"collector_system_status_fail":                                                                  0,
				"collector_system_status_success":                                                               1,
				"collector_tcp_duration":                                                                        0,
				"collector_tcp_status_fail":                                                                     0,
				"collector_tcp_status_success":                                                                  1,
				"cpu_core_0,0_cstate_c1":                                                                        160233427,
				"cpu_core_0,0_cstate_c2":                                                                        0,
				"cpu_core_0,0_cstate_c3":                                                                        0,
				"cpu_core_0,0_dpc_time":                                                                         67109,
				"cpu_core_0,0_dpcs":                                                                             4871900,
				"cpu_core_0,0_idle_time":                                                                        162455593,
				"cpu_core_0,0_interrupt_time":                                                                   77281,
				"cpu_core_0,0_interrupts":                                                                       155194331,
				"cpu_core_0,0_privileged_time":                                                                  1182109,
				"cpu_core_0,0_user_time":                                                                        1073671,
				"cpu_core_0,1_cstate_c1":                                                                        159528054,
				"cpu_core_0,1_cstate_c2":                                                                        0,
				"cpu_core_0,1_cstate_c3":                                                                        0,
				"cpu_core_0,1_dpc_time":                                                                         11093,
				"cpu_core_0,1_dpcs":                                                                             1650552,
				"cpu_core_0,1_idle_time":                                                                        159478125,
				"cpu_core_0,1_interrupt_time":                                                                   58093,
				"cpu_core_0,1_interrupts":                                                                       79325847,
				"cpu_core_0,1_privileged_time":                                                                  1801234,
				"cpu_core_0,1_user_time":                                                                        3432000,
				"cpu_core_0,2_cstate_c1":                                                                        159891723,
				"cpu_core_0,2_cstate_c2":                                                                        0,
				"cpu_core_0,2_cstate_c3":                                                                        0,
				"cpu_core_0,2_dpc_time":                                                                         16062,
				"cpu_core_0,2_dpcs":                                                                             2236469,
				"cpu_core_0,2_idle_time":                                                                        159848437,
				"cpu_core_0,2_interrupt_time":                                                                   53515,
				"cpu_core_0,2_interrupts":                                                                       67305419,
				"cpu_core_0,2_privileged_time":                                                                  1812546,
				"cpu_core_0,2_user_time":                                                                        3050250,
				"cpu_core_0,3_cstate_c1":                                                                        159544117,
				"cpu_core_0,3_cstate_c2":                                                                        0,
				"cpu_core_0,3_cstate_c3":                                                                        0,
				"cpu_core_0,3_dpc_time":                                                                         8140,
				"cpu_core_0,3_dpcs":                                                                             1185046,
				"cpu_core_0,3_idle_time":                                                                        159527546,
				"cpu_core_0,3_interrupt_time":                                                                   44484,
				"cpu_core_0,3_interrupts":                                                                       60766938,
				"cpu_core_0,3_privileged_time":                                                                  1760828,
				"cpu_core_0,3_user_time":                                                                        3422875,
				"cpu_dpc_time":                                                                                  102404,
				"cpu_idle_time":                                                                                 641309701,
				"cpu_interrupt_time":                                                                            233373,
				"cpu_privileged_time":                                                                           6556717,
				"cpu_user_time":                                                                                 10978796,
				"iis_website_Default_Web_Site_connection_attempts_all_instances_total":                          1,
				"iis_website_Default_Web_Site_current_anonymous_users":                                          0,
				"iis_website_Default_Web_Site_current_connections":                                              0,
				"iis_website_Default_Web_Site_current_isapi_extension_requests":                                 0,
				"iis_website_Default_Web_Site_current_non_anonymous_users":                                      0,
				"iis_website_Default_Web_Site_files_received_total":                                             0,
				"iis_website_Default_Web_Site_files_sent_total":                                                 2,
				"iis_website_Default_Web_Site_isapi_extension_requests_total":                                   0,
				"iis_website_Default_Web_Site_locked_errors_total":                                              0,
				"iis_website_Default_Web_Site_logon_attempts_total":                                             4,
				"iis_website_Default_Web_Site_not_found_errors_total":                                           1,
				"iis_website_Default_Web_Site_received_bytes_total":                                             10289,
				"iis_website_Default_Web_Site_requests_total":                                                   3,
				"iis_website_Default_Web_Site_sent_bytes_total":                                                 105882,
				"iis_website_Default_Web_Site_service_uptime":                                                   258633,
				"logical_disk_C:_free_space":                                                                    43636490240,
				"logical_disk_C:_read_bytes_total":                                                              17676328448,
				"logical_disk_C:_read_latency":                                                                  97420,
				"logical_disk_C:_reads_total":                                                                   350593,
				"logical_disk_C:_total_space":                                                                   67938287616,
				"logical_disk_C:_used_space":                                                                    24301797376,
				"logical_disk_C:_write_bytes_total":                                                             9135282688,
				"logical_disk_C:_write_latency":                                                                 123912,
				"logical_disk_C:_writes_total":                                                                  450705,
				"logon_type_batch_sessions":                                                                     0,
				"logon_type_cached_interactive_sessions":                                                        0,
				"logon_type_cached_remote_interactive_sessions":                                                 0,
				"logon_type_cached_unlock_sessions":                                                             0,
				"logon_type_interactive_sessions":                                                               2,
				"logon_type_network_clear_text_sessions":                                                        0,
				"logon_type_network_sessions":                                                                   0,
				"logon_type_new_credentials_sessions":                                                           0,
				"logon_type_proxy_sessions":                                                                     0,
				"logon_type_remote_interactive_sessions":                                                        0,
				"logon_type_service_sessions":                                                                   0,
				"logon_type_system_sessions":                                                                    0,
				"logon_type_unlock_sessions":                                                                    0,
				"memory_available_bytes":                                                                        1379942400,
				"memory_cache_faults_total":                                                                     8009603,
				"memory_cache_total":                                                                            1392185344,
				"memory_commit_limit":                                                                           5733113856,
				"memory_committed_bytes":                                                                        3447439360,
				"memory_modified_page_list_bytes":                                                               32653312,
				"memory_not_committed_bytes":                                                                    2285674496,
				"memory_page_faults_total":                                                                      119093924,
				"memory_pool_nonpaged_bytes_total":                                                              126865408,
				"memory_pool_paged_bytes":                                                                       303906816,
				"memory_standby_cache_core_bytes":                                                               107376640,
				"memory_standby_cache_normal_priority_bytes":                                                    1019121664,
				"memory_standby_cache_reserve_bytes":                                                            233033728,
				"memory_standby_cache_total":                                                                    1359532032,
				"memory_swap_page_reads_total":                                                                  402087,
				"memory_swap_page_writes_total":                                                                 7012,
				"memory_swap_pages_read_total":                                                                  4643279,
				"memory_swap_pages_written_total":                                                               312896,
				"memory_used_bytes":                                                                             2876776448,
				"mssql_db_master_instance_SQLEXPRESS_active_transactions":                                       0,
				"mssql_db_master_instance_SQLEXPRESS_backup_restore_operations":                                 0,
				"mssql_db_master_instance_SQLEXPRESS_data_files_size_bytes":                                     4653056,
				"mssql_db_master_instance_SQLEXPRESS_log_flushed_bytes":                                         3702784,
				"mssql_db_master_instance_SQLEXPRESS_log_flushes":                                               252,
				"mssql_db_master_instance_SQLEXPRESS_transactions":                                              2183,
				"mssql_db_master_instance_SQLEXPRESS_write_transactions":                                        236,
				"mssql_db_model_instance_SQLEXPRESS_active_transactions":                                        0,
				"mssql_db_model_instance_SQLEXPRESS_backup_restore_operations":                                  0,
				"mssql_db_model_instance_SQLEXPRESS_data_files_size_bytes":                                      8388608,
				"mssql_db_model_instance_SQLEXPRESS_log_flushed_bytes":                                          12288,
				"mssql_db_model_instance_SQLEXPRESS_log_flushes":                                                3,
				"mssql_db_model_instance_SQLEXPRESS_transactions":                                               4467,
				"mssql_db_model_instance_SQLEXPRESS_write_transactions":                                         0,
				"mssql_db_msdb_instance_SQLEXPRESS_active_transactions":                                         0,
				"mssql_db_msdb_instance_SQLEXPRESS_backup_restore_operations":                                   0,
				"mssql_db_msdb_instance_SQLEXPRESS_data_files_size_bytes":                                       15466496,
				"mssql_db_msdb_instance_SQLEXPRESS_log_flushed_bytes":                                           0,
				"mssql_db_msdb_instance_SQLEXPRESS_log_flushes":                                                 0,
				"mssql_db_msdb_instance_SQLEXPRESS_transactions":                                                4582,
				"mssql_db_msdb_instance_SQLEXPRESS_write_transactions":                                          0,
				"mssql_db_mssqlsystemresource_instance_SQLEXPRESS_active_transactions":                          0,
				"mssql_db_mssqlsystemresource_instance_SQLEXPRESS_backup_restore_operations":                    0,
				"mssql_db_mssqlsystemresource_instance_SQLEXPRESS_data_files_size_bytes":                        41943040,
				"mssql_db_mssqlsystemresource_instance_SQLEXPRESS_log_flushed_bytes":                            0,
				"mssql_db_mssqlsystemresource_instance_SQLEXPRESS_log_flushes":                                  0,
				"mssql_db_mssqlsystemresource_instance_SQLEXPRESS_transactions":                                 2,
				"mssql_db_mssqlsystemresource_instance_SQLEXPRESS_write_transactions":                           0,
				"mssql_db_tempdb_instance_SQLEXPRESS_active_transactions":                                       0,
				"mssql_db_tempdb_instance_SQLEXPRESS_backup_restore_operations":                                 0,
				"mssql_db_tempdb_instance_SQLEXPRESS_data_files_size_bytes":                                     8388608,
				"mssql_db_tempdb_instance_SQLEXPRESS_log_flushed_bytes":                                         118784,
				"mssql_db_tempdb_instance_SQLEXPRESS_log_flushes":                                               2,
				"mssql_db_tempdb_instance_SQLEXPRESS_transactions":                                              1558,
				"mssql_db_tempdb_instance_SQLEXPRESS_write_transactions":                                        29,
				"mssql_instance_SQLEXPRESS_accessmethods_page_splits":                                           429,
				"mssql_instance_SQLEXPRESS_bufman_buffer_cache_hits":                                            86,
				"mssql_instance_SQLEXPRESS_bufman_checkpoint_pages":                                             82,
				"mssql_instance_SQLEXPRESS_bufman_page_life_expectancy_seconds":                                 191350,
				"mssql_instance_SQLEXPRESS_bufman_page_reads":                                                   797,
				"mssql_instance_SQLEXPRESS_bufman_page_writes":                                                  92,
				"mssql_instance_SQLEXPRESS_cache_hit_ratio":                                                     100,
				"mssql_instance_SQLEXPRESS_genstats_blocked_processes":                                          0,
				"mssql_instance_SQLEXPRESS_genstats_user_connections":                                           1,
				"mssql_instance_SQLEXPRESS_memmgr_connection_memory_bytes":                                      1015808,
				"mssql_instance_SQLEXPRESS_memmgr_external_benefit_of_memory":                                   0,
				"mssql_instance_SQLEXPRESS_memmgr_pending_memory_grants":                                        0,
				"mssql_instance_SQLEXPRESS_memmgr_total_server_memory_bytes":                                    198836224,
				"mssql_instance_SQLEXPRESS_resource_AllocUnit_locks_deadlocks":                                  0,
				"mssql_instance_SQLEXPRESS_resource_AllocUnit_locks_lock_wait_seconds":                          0,
				"mssql_instance_SQLEXPRESS_resource_Application_locks_deadlocks":                                0,
				"mssql_instance_SQLEXPRESS_resource_Application_locks_lock_wait_seconds":                        0,
				"mssql_instance_SQLEXPRESS_resource_Database_locks_deadlocks":                                   0,
				"mssql_instance_SQLEXPRESS_resource_Database_locks_lock_wait_seconds":                           0,
				"mssql_instance_SQLEXPRESS_resource_Extent_locks_deadlocks":                                     0,
				"mssql_instance_SQLEXPRESS_resource_Extent_locks_lock_wait_seconds":                             0,
				"mssql_instance_SQLEXPRESS_resource_File_locks_deadlocks":                                       0,
				"mssql_instance_SQLEXPRESS_resource_File_locks_lock_wait_seconds":                               0,
				"mssql_instance_SQLEXPRESS_resource_HoBT_locks_deadlocks":                                       0,
				"mssql_instance_SQLEXPRESS_resource_HoBT_locks_lock_wait_seconds":                               0,
				"mssql_instance_SQLEXPRESS_resource_Key_locks_deadlocks":                                        0,
				"mssql_instance_SQLEXPRESS_resource_Key_locks_lock_wait_seconds":                                0,
				"mssql_instance_SQLEXPRESS_resource_Metadata_locks_deadlocks":                                   0,
				"mssql_instance_SQLEXPRESS_resource_Metadata_locks_lock_wait_seconds":                           0,
				"mssql_instance_SQLEXPRESS_resource_OIB_locks_deadlocks":                                        0,
				"mssql_instance_SQLEXPRESS_resource_OIB_locks_lock_wait_seconds":                                0,
				"mssql_instance_SQLEXPRESS_resource_Object_locks_deadlocks":                                     0,
				"mssql_instance_SQLEXPRESS_resource_Object_locks_lock_wait_seconds":                             0,
				"mssql_instance_SQLEXPRESS_resource_Page_locks_deadlocks":                                       0,
				"mssql_instance_SQLEXPRESS_resource_Page_locks_lock_wait_seconds":                               0,
				"mssql_instance_SQLEXPRESS_resource_RID_locks_deadlocks":                                        0,
				"mssql_instance_SQLEXPRESS_resource_RID_locks_lock_wait_seconds":                                0,
				"mssql_instance_SQLEXPRESS_resource_RowGroup_locks_deadlocks":                                   0,
				"mssql_instance_SQLEXPRESS_resource_RowGroup_locks_lock_wait_seconds":                           0,
				"mssql_instance_SQLEXPRESS_resource_Xact_locks_deadlocks":                                       0,
				"mssql_instance_SQLEXPRESS_resource_Xact_locks_lock_wait_seconds":                               0,
				"mssql_instance_SQLEXPRESS_sql_errors_total_db_offline_errors":                                  0,
				"mssql_instance_SQLEXPRESS_sql_errors_total_info_errors":                                        766,
				"mssql_instance_SQLEXPRESS_sql_errors_total_kill_connection_errors":                             0,
				"mssql_instance_SQLEXPRESS_sql_errors_total_user_errors":                                        29,
				"mssql_instance_SQLEXPRESS_sqlstats_auto_parameterization_attempts":                             37,
				"mssql_instance_SQLEXPRESS_sqlstats_batch_requests":                                             2972,
				"mssql_instance_SQLEXPRESS_sqlstats_safe_auto_parameterization_attempts":                        2,
				"mssql_instance_SQLEXPRESS_sqlstats_sql_compilations":                                           376,
				"mssql_instance_SQLEXPRESS_sqlstats_sql_recompilations":                                         8,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_bytes_received":                                 38290755856,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_bytes_sent":                                     8211165504,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_outbound_discarded":                     0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_outbound_errors":                        0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_discarded":                     0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_errors":                        0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_total":                         4120869,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_sent_total":                             1332466,
				"netframework_clrexceptions_WMSvc_filters":                                                      0,
				"netframework_clrexceptions_WMSvc_finallys":                                                     0,
				"netframework_clrexceptions_WMSvc_throw_catch_depth":                                            0,
				"netframework_clrexceptions_WMSvc_thrown":                                                       0,
				"netframework_clrexceptions_powershell_filters":                                                 0,
				"netframework_clrexceptions_powershell_finallys":                                                56,
				"netframework_clrexceptions_powershell_throw_catch_depth":                                       140,
				"netframework_clrexceptions_powershell_thrown":                                                  37,
				"netframework_clrinterop_WMSvc_com_callable_wrappers":                                           2,
				"netframework_clrinterop_WMSvc_marshalling":                                                     0,
				"netframework_clrinterop_WMSvc_stubs_created":                                                   29,
				"netframework_clrinterop_powershell_com_callable_wrappers":                                      5,
				"netframework_clrinterop_powershell_marshalling":                                                0,
				"netframework_clrinterop_powershell_stubs_created":                                              345,
				"netframework_clrjit_WMSvc_il_bytes":                                                            4007,
				"netframework_clrjit_WMSvc_methods":                                                             27,
				"netframework_clrjit_WMSvc_standard_failure":                                                    0,
				"netframework_clrjit_WMSvc_time":                                                                0,
				"netframework_clrjit_powershell_il_bytes":                                                       47021,
				"netframework_clrjit_powershell_methods":                                                        344,
				"netframework_clrjit_powershell_standard_failure":                                               0,
				"netframework_clrjit_powershell_time":                                                           0,
				"netframework_clrloading_WMSvc_app_domains_loaded":                                              1,
				"netframework_clrloading_WMSvc_app_domains_unloaded":                                            0,
				"netframework_clrloading_WMSvc_assemblies_loaded":                                               5,
				"netframework_clrloading_WMSvc_class_load_failure":                                              0,
				"netframework_clrloading_WMSvc_classes_loaded":                                                  18,
				"netframework_clrloading_WMSvc_loader_heap_size":                                                270336,
				"netframework_clrloading_powershell_app_domains_loaded":                                         1,
				"netframework_clrloading_powershell_app_domains_unloaded":                                       0,
				"netframework_clrloading_powershell_assemblies_loaded":                                          20,
				"netframework_clrloading_powershell_class_load_failure":                                         1,
				"netframework_clrloading_powershell_classes_loaded":                                             477,
				"netframework_clrloading_powershell_loader_heap_size":                                           2285568,
				"netframework_clrlockandthreads_WMSvc_contentions":                                              0,
				"netframework_clrlockandthreads_WMSvc_current_logical_threads":                                  2,
				"netframework_clrlockandthreads_WMSvc_current_physical_threads":                                 1,
				"netframework_clrlockandthreads_WMSvc_queue_length":                                             0,
				"netframework_clrlockandthreads_WMSvc_recognized_threads":                                       1,
				"netframework_clrlockandthreads_powershell_contentions":                                         10,
				"netframework_clrlockandthreads_powershell_current_logical_threads":                             16,
				"netframework_clrlockandthreads_powershell_current_physical_threads":                            13,
				"netframework_clrlockandthreads_powershell_queue_length":                                        3,
				"netframework_clrlockandthreads_powershell_recognized_threads":                                  6,
				"netframework_clrmemory_WMSvc_allocated_bytes_total":                                            227792,
				"netframework_clrmemory_WMSvc_collection":                                                       2,
				"netframework_clrmemory_WMSvc_committed":                                                        270336,
				"netframework_clrmemory_WMSvc_finalization_survivors":                                           7,
				"netframework_clrmemory_WMSvc_gc_handles":                                                       24,
				"netframework_clrmemory_WMSvc_gc_time":                                                          0,
				"netframework_clrmemory_WMSvc_heap_size":                                                        4312696,
				"netframework_clrmemory_WMSvc_induced_gc_total":                                                 0,
				"netframework_clrmemory_WMSvc_pinned_objects":                                                   1,
				"netframework_clrmemory_WMSvc_promoted":                                                         49720,
				"netframework_clrmemory_WMSvc_reserved":                                                         402644992,
				"netframework_clrmemory_WMSvc_sink_block_in_use":                                                1,
				"netframework_clrmemory_powershell_allocated_bytes_total":                                       46333800,
				"netframework_clrmemory_powershell_collection":                                                  11,
				"netframework_clrmemory_powershell_committed":                                                   20475904,
				"netframework_clrmemory_powershell_finalization_survivors":                                      244,
				"netframework_clrmemory_powershell_gc_handles":                                                  834,
				"netframework_clrmemory_powershell_gc_time":                                                     0,
				"netframework_clrmemory_powershell_heap_size":                                                   34711872,
				"netframework_clrmemory_powershell_induced_gc_total":                                            0,
				"netframework_clrmemory_powershell_pinned_objects":                                              0,
				"netframework_clrmemory_powershell_promoted":                                                    107352,
				"netframework_clrmemory_powershell_reserved":                                                    402644992,
				"netframework_clrmemory_powershell_sink_block_in_use":                                           42,
				"netframework_clrremoting_WMSvc_calls":                                                          0,
				"netframework_clrremoting_WMSvc_channels":                                                       0,
				"netframework_clrremoting_WMSvc_context_bound_classes_loaded":                                   0,
				"netframework_clrremoting_WMSvc_context_bound_objects":                                          0,
				"netframework_clrremoting_WMSvc_context_proxies":                                                0,
				"netframework_clrremoting_WMSvc_contexts":                                                       1,
				"netframework_clrremoting_powershell_calls":                                                     0,
				"netframework_clrremoting_powershell_channels":                                                  0,
				"netframework_clrremoting_powershell_context_bound_classes_loaded":                              0,
				"netframework_clrremoting_powershell_context_bound_objects":                                     0,
				"netframework_clrremoting_powershell_context_proxies":                                           0,
				"netframework_clrremoting_powershell_contexts":                                                  1,
				"netframework_clrsecurity_WMSvc_checks_time":                                                    0,
				"netframework_clrsecurity_WMSvc_link_time_checks":                                               0,
				"netframework_clrsecurity_WMSvc_runtime_checks":                                                 3,
				"netframework_clrsecurity_WMSvc_stack_walk_depth":                                               1,
				"netframework_clrsecurity_powershell_checks_time":                                               0,
				"netframework_clrsecurity_powershell_link_time_checks":                                          0,
				"netframework_clrsecurity_powershell_runtime_checks":                                            4386,
				"netframework_clrsecurity_powershell_stack_walk_depth":                                          1,
				"os_paging_free_bytes":                                                                          1414107136,
				"os_paging_limit_bytes":                                                                         1476395008,
				"os_paging_used_bytes":                                                                          62287872,
				"os_physical_memory_free_bytes":                                                                 1379946496,
				"os_processes":                                                                                  152,
				"os_processes_limit":                                                                            4294967295,
				"os_users":                                                                                      2,
				"os_visible_memory_bytes":                                                                       4256718848,
				"os_visible_memory_used_bytes":                                                                  2876772352,
				"process_msedge_cpu_time":                                                                       1919893,
				"process_msedge_handles":                                                                        5779,
				"process_msedge_io_bytes":                                                                       3978227378,
				"process_msedge_io_operations":                                                                  16738642,
				"process_msedge_page_faults":                                                                    5355941,
				"process_msedge_page_file_bytes":                                                                681603072,
				"process_msedge_threads":                                                                        213,
				"process_msedge_working_set_private_bytes":                                                      461344768,
				"service_dhcp_state_continue_pending":                                                           0,
				"service_dhcp_state_pause_pending":                                                              0,
				"service_dhcp_state_paused":                                                                     0,
				"service_dhcp_state_running":                                                                    1,
				"service_dhcp_state_start_pending":                                                              0,
				"service_dhcp_state_stop_pending":                                                               0,
				"service_dhcp_state_stopped":                                                                    0,
				"service_dhcp_state_unknown":                                                                    0,
				"service_dhcp_status_degraded":                                                                  0,
				"service_dhcp_status_error":                                                                     0,
				"service_dhcp_status_lost_comm":                                                                 0,
				"service_dhcp_status_no_contact":                                                                0,
				"service_dhcp_status_nonrecover":                                                                0,
				"service_dhcp_status_ok":                                                                        1,
				"service_dhcp_status_pred_fail":                                                                 0,
				"service_dhcp_status_service":                                                                   0,
				"service_dhcp_status_starting":                                                                  0,
				"service_dhcp_status_stopping":                                                                  0,
				"service_dhcp_status_stressed":                                                                  0,
				"service_dhcp_status_unknown":                                                                   0,
				"system_threads":                                                                                1559,
				"system_up_time":                                                                                2890557,
				"tcp_ipv4_conns_active":                                                                         4301,
				"tcp_ipv4_conns_established":                                                                    7,
				"tcp_ipv4_conns_failures":                                                                       137,
				"tcp_ipv4_conns_passive":                                                                        501,
				"tcp_ipv4_conns_resets":                                                                         1282,
				"tcp_ipv4_segments_received":                                                                    676388,
				"tcp_ipv4_segments_retransmitted":                                                               2120,
				"tcp_ipv4_segments_sent":                                                                        871379,
				"tcp_ipv6_conns_active":                                                                         214,
				"tcp_ipv6_conns_established":                                                                    0,
				"tcp_ipv6_conns_failures":                                                                       214,
				"tcp_ipv6_conns_passive":                                                                        0,
				"tcp_ipv6_conns_resets":                                                                         0,
				"tcp_ipv6_segments_received":                                                                    1284,
				"tcp_ipv6_segments_retransmitted":                                                               428,
				"tcp_ipv6_segments_sent":                                                                        856,
			},
		},
		"fails if endpoint returns invalid data": {
			prepare: prepareWMIReturnsInvalidData,
		},
		"fails on connection refused": {
			prepare: prepareWMIConnectionRefused,
		},
		"fails on 404 response": {
			prepare: prepareWMIResponse404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wmi, cleanup := test.prepare()
			defer cleanup()

			require.True(t, wmi.Init())

			mx := wmi.Collect()

			if mx != nil && test.wantCollected != nil {
				mx["system_up_time"] = test.wantCollected["system_up_time"]
			}

			assert.Equal(t, test.wantCollected, mx)
			if len(test.wantCollected) > 0 {
				testCharts(t, wmi, mx)
			}
		})
	}
}

func testCharts(t *testing.T, wmi *WMI, mx map[string]int64) {
	ensureChartsDimsCreated(t, wmi)
	ensureCollectedHasAllChartsDimsVarsIDs(t, wmi, mx)
}

func ensureChartsDimsCreated(t *testing.T, w *WMI) {
	for _, chart := range cpuCharts {
		if w.cache.collection[collectorCPU] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range memCharts {
		if w.cache.collection[collectorMemory] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range tcpCharts {
		if w.cache.collection[collectorTCP] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range osCharts {
		if w.cache.collection[collectorOS] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range systemCharts {
		if w.cache.collection[collectorSystem] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range logonCharts {
		if w.cache.collection[collectorLogon] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range processesCharts {
		if w.cache.collection[collectorProcess] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range netFrameworkCLRExceptionsChartsTmpl {
		if w.cache.collection[collectorNetFrameworkCLRExceptions] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range netFrameworkCLRInteropChartsTmpl {
		if w.cache.collection[collectorNetFrameworkCLRInterop] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range netFrameworkCLRJITChartsTmpl {
		if w.cache.collection[collectorNetFrameworkCLRJIT] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range netFrameworkCLRLoadingChartsTmpl {
		if w.cache.collection[collectorNetFrameworkCLRLoading] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range netFrameworkCLRLocksandthreadsChartsTmpl {
		if w.cache.collection[collectorNetFrameworkCLRLocksandthreads] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range netFrameworkCLRMemoryChartsTmpl {
		if w.cache.collection[collectorNetFrameworkCLRMemory] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range netFrameworkCLRRemotingChartsTmpl {
		if w.cache.collection[collectorNetFrameworkCLRRemoting] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range netFrameworkCLRSecurityChartsTmpl {
		if w.cache.collection[collectorNetFrameworkCLRSecurity] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}

	for core := range w.cache.cores {
		for _, chart := range cpuCoreChartsTmpl {
			id := fmt.Sprintf(chart.ID, core)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' core", id, core)
		}
	}
	for disk := range w.cache.volumes {
		for _, chart := range diskChartsTmpl {
			id := fmt.Sprintf(chart.ID, disk)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' disk", id, disk)
		}
	}
	for nic := range w.cache.nics {
		for _, chart := range nicChartsTmpl {
			id := fmt.Sprintf(chart.ID, nic)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' nic", id, nic)
		}
	}
	for zone := range w.cache.thermalZones {
		for _, chart := range thermalzoneChartsTmpl {
			id := fmt.Sprintf(chart.ID, zone)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' thermalzone", id, zone)
		}
	}
	for svc := range w.cache.services {
		for _, chart := range serviceChartsTmpl {
			id := fmt.Sprintf(chart.ID, svc)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' service", id, svc)
		}
	}
	for website := range w.cache.iis {
		for _, chart := range iisWebsiteChartsTmpl {
			id := fmt.Sprintf(chart.ID, website)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' website", id, website)
		}
	}
	for instance := range w.cache.mssqlInstances {
		for _, chart := range mssqlInstanceChartsTmpl {
			id := fmt.Sprintf(chart.ID, instance)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' instance", id, instance)
		}
	}
	for instanceDB := range w.cache.mssqlDBs {
		s := strings.Split(instanceDB, ":")
		if assert.Lenf(t, s, 2, "can not extract intance/database from cache.mssqlDBs") {
			instance, db := s[0], s[1]
			for _, chart := range mssqlDatabaseChartsTmpl {
				id := fmt.Sprintf(chart.ID, db, instance)
				assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' instance", id, instance)
			}
		}
	}
	for _, chart := range adCharts {
		if w.cache.collection[collectorAD] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for template := range w.cache.adcs {
		for _, chart := range adcsCertTemplateChartsTmpl {
			id := fmt.Sprintf(chart.ID, template)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' template certificate", id, template)
		}
	}
	for name := range w.cache.collectors {
		for _, chart := range collectorChartsTmpl {
			id := fmt.Sprintf(chart.ID, name)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' collector", id, name)
		}
	}

	for _, chart := range processesCharts {
		if chart = w.Charts().Get(chart.ID); chart == nil {
			continue
		}
		for proc := range w.cache.processes {
			var found bool
			for _, dim := range chart.Dims {
				if found = strings.HasPrefix(dim.ID, "process_"+proc); found {
					break
				}
			}
			assert.Truef(t, found, "chart '%s' has not dim for '%s' process", chart.ID, proc)
		}
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, w *WMI, mx map[string]int64) {
	for _, chart := range *w.Charts() {
		for _, dim := range chart.Dims {
			_, ok := mx[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := mx[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

func prepareWMIv0200() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(v0200Metrics)
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}

func prepareWMIReturnsInvalidData() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}

func prepareWMIConnectionRefused() (wmi *WMI, cleanup func()) {
	wmi = New()
	wmi.URL = "http://127.0.0.1:38001"
	return wmi, func() {}
}

func prepareWMIResponse404() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}
