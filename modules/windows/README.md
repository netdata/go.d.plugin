# Windows machine collector

This module will monitor one or more Windows machines, using
the [windows_exporter](https://github.com/prometheus-community/windows_exporter).

## Requirements

Netdata monitors Windows hosts by utilizing the
[Prometheus exporter for Windows machines](https://github.com/prometheus-community/windows_exporter), a native Windows
agent running on each host.

To quickly test Netdata directly on a Windows machine, you can use
the [Netdata MSI installer](https://github.com/netdata/msi-installer#instructions). The installer runs Netdata in a
custom WSL deployment. WSL was not designed for production environments, so **we do not recommend using the MSI installer in
production**.

For production use, you need to install Netdata on one or more nodes running Linux:

![windows](https://user-images.githubusercontent.com/43294513/232522572-1fe51228-953b-43d2-81c4-dcb0a8974db5.jpg)

- Install the
  latest [Prometheus exporter for Windows](https://github.com/prometheus-community/windows_exporter/releases)
  on every Windows host you want to monitor, enabling the collectors listed here:
  ```msiexec /i "[PATH_TO_MSI]" ENABLED_COLLECTORS=process,ad,adcs,adfs,cpu,dns,memory,mssql,net,os,tcp,logical_disk```
- Get the installation commands from [Netdata Cloud](https://app.netdata.cloud) and install Netdata on one or more Linux
  nodes.
- Configure each Netdata instance to collect data remotely, from several Windows hosts. Just add one job
  for each host to  `windows.conf`, as shown in the [configuration section](#configuration).

Automated charts and alerts for your entire Windows infrastructure will be automatically generated.
Each Windows host (data collection job) will be identifiable as an "instance" in the Netdata Cloud charts.

## Metrics

The module collects metrics from the following collectors:

- [cpu](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.cpu.md)
- [iis](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.iis.md)
- [memory](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.memory.md)
- [net](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.net.md)
- [logical_disk](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.logical_disk.md)
- [os](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.os.md)
- [system](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.system.md)
- [logon](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.logon.md)
- [tcp](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.tcp.md)
- [thermalzone](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.thermalzone.md)
- [process](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.process.md)
- [service](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.service.md)
- [mssql](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.mssql.md)
- [ad](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.ad.md)
- [adcs](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.adcs.md)
- [adfs](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.adfs.md)
- [netframework_clrexceptions](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.netframework_clrexceptions.md)
- [netframework_clrinterop](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.netframework_clrinterop.md)
- [netframework_clrjit](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.netframework_clrjit.md)
- [netframework_clrloading](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.netframework_clrloading.md)
- [netframework_clrlocksandthreads](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.netframework_clrlocksandthreads.md)
- [netframework_clrmemory](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.netframework_clrmemory.md)
- [netframework_clrremoting](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.netframework_clrremoting.md)
- [exchange](https://github.com/prometheus-community/windows_exporter/blob/master/docs/collector.exchange.md)

All metrics have a prefix.

Labels per scope:

- global: no labels.
- logical disk: disk.
- network device: nic.
- thermalzone: thermalzone.
- website: website.
- mssql instance: mssql_instance.
- database: mssql_instance, database.
- certificate template: cert_template.
- service: service.
- process: process.

| Metric                                                               |        Scope         |                                                                                     Dimensions                                                                                     |       Units       |
|----------------------------------------------------------------------|:--------------------:|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-----------------:|
| windows.cpu_utilization_total                                        |        global        |                                                                          dpc, user, privileged, interrupt                                                                          |    percentage     |
| windows.cpu_core_utilization                                         |       cpu core       |                                                                          dpc, user, privileged, interrupt                                                                          |    percentage     |
| windows.cpu_core_interrupts                                          |       cpu core       |                                                                                     interrupts                                                                                     |   interrupts/s    |
| windows.cpu_core_dpcs                                                |       cpu core       |                                                                                        dpcs                                                                                        |      dpcs/s       |
| windows.cpu_core_cstate                                              |       cpu core       |                                                                                     c1, c2, c3                                                                                     |    percentage     |
| windows.memory_utilization                                           |        global        |                                                                                  available, used                                                                                   |       bytes       |
| windows.memory_utilization                                           |        global        |                                                                                  available, used                                                                                   |        KiB        |
| windows.memory_page_faults                                           |        global        |                                                                                    page_faults                                                                                     |     events/s      |
| windows.memory_swap_utilization                                      |        global        |                                                                                  available, used                                                                                   |       bytes       |
| windows.memory_swap_operations                                       |        global        |                                                                                    read, write                                                                                     |   operations/s    |
| windows.memory_swap_pages                                            |        global        |                                                                                   read, written                                                                                    |      pages/s      |
| windows.memory_cached                                                |        global        |                                                                                       cached                                                                                       |        KiB        |
| windows.memory_cache_faults                                          |        global        |                                                                                    cache_faults                                                                                    |     events/s      |
| windows.memory_system_pool                                           |        global        |                                                                                  paged, non-paged                                                                                  |       bytes       |
| windows.logical_disk_utilization                                     |     logical disk     |                                                                                     free, used                                                                                     |       bytes       |
| windows.logical_disk_bandwidth                                       |     logical disk     |                                                                                    read, write                                                                                     |      bytes/s      |
| windows.logical_disk_operations                                      |     logical disk     |                                                                                   reads, writes                                                                                    |   operations/s    |
| windows.logical_disk_latency                                         |     logical disk     |                                                                                    read, write                                                                                     |      seconds      |
| windows.net_nic_bandwidth                                            |    network device    |                                                                                   received, sent                                                                                   |    kilobits/s     |
| windows.net_nic_packets                                              |    network device    |                                                                                   received, sent                                                                                   |     packets/s     |
| windows.net_nic_errors                                               |    network device    |                                                                                 inbound, outbound                                                                                  |     errors/s      |
| windows.net_nic_discarded                                            |    network device    |                                                                                 inbound, outbound                                                                                  |    discards/s     |
| windows.tcp_conns_established                                        |        global        |                                                                                     ipv4, ipv6                                                                                     |    connections    |
| windows.tcp_conns_active                                             |        global        |                                                                                     ipv4, ipv6                                                                                     |   connections/s   |
| windows.tcp_conns_passive                                            |        global        |                                                                                     ipv4, ipv6                                                                                     |   connections/s   |
| windows.tcp_conns_failures                                           |        global        |                                                                                     ipv4, ipv6                                                                                     |    failures/s     |
| windows.tcp_conns_resets                                             |        global        |                                                                                     ipv4, ipv6                                                                                     |     resets/s      |
| windows.tcp_segments_received                                        |        global        |                                                                                     ipv4, ipv6                                                                                     |    segments/s     |
| windows.tcp_segments_sent                                            |        global        |                                                                                     ipv4, ipv6                                                                                     |    segments/s     |
| windows.tcp_segments_retransmitted                                   |        global        |                                                                                     ipv4, ipv6                                                                                     |    segments/s     |
| windows.os_processes                                                 |        global        |                                                                                     processes                                                                                      |      number       |
| windows.os_users                                                     |        global        |                                                                                       users                                                                                        |       users       |
| windows.os_visible_memory_usage                                      |        global        |                                                                                     free, used                                                                                     |       bytes       |
| windows.os_paging_files_usage                                        |        global        |                                                                                     free, used                                                                                     |       bytes       |
| windows.system_threads                                               |        global        |                                                                                      threads                                                                                       |      number       |
| windows.system_uptime                                                |        global        |                                                                                        time                                                                                        |      seconds      |
| windows.logon_type_sessions                                          |        global        | system, interactive, network, batch, service, proxy, unlock, network_clear_text, new_credentials, remote_interactive, cached_interactive, cached_remote_interactive, cached_unlock |      seconds      |
| windows.thermalzone_temperature                                      |        global        |                                                                         <i>a dimension per thermalzone</i>                                                                         |      celsius      |
| windows.processes_cpu_utilization                                    |        global        |                                                                           <i>a dimension per process</i>                                                                           |    percentage     |
| windows.processes_handles                                            |        global        |                                                                           <i>a dimension per process</i>                                                                           |      handles      |
| windows.processes_io_bytes                                           |        global        |                                                                           <i>a dimension per process</i>                                                                           |      bytes/s      |
| windows.processes_io_operations                                      |        global        |                                                                           <i>a dimension per process</i>                                                                           |   operations/s    |
| windows.processes_page_faults                                        |        global        |                                                                           <i>a dimension per process</i>                                                                           |    pgfaults/s     |
| windows.processes_page_file_bytes                                    |        global        |                                                                           <i>a dimension per process</i>                                                                           |       bytes       |
| windows.processes_pool_bytes                                         |        global        |                                                                           <i>a dimension per process</i>                                                                           |       bytes       |
| windows.processes_threads                                            |        global        |                                                                           <i>a dimension per process</i>                                                                           |      threads      |
| windows.service_state                                                |       service        |                                          running, stopped, start_pending, stop_pending, continue_pending, pause_pending, paused, unknown                                           |       state       |
| windows.service_status                                               |       service        |                                 ok, error, unknown, degraded, pred_fail, starting, stopping, service, stressed, nonrecover, no_contact, lost_comm                                  |      status       |
| iis.website_traffic                                                  |       website        |                                                                                   received, sent                                                                                   |      bytes/s      |
| iis.website_requests_rate                                            |       website        |                                                                                      requests                                                                                      |    requests/s     |
| iis.website_active_connections_count                                 |       website        |                                                                                       active                                                                                       |    connections    |
| iis.website_users_count                                              |       website        |                                                                              anonymous, non_anonymous                                                                              |       users       |
| iis.website_connection_attempts_rate                                 |       website        |                                                                                     connection                                                                                     |    attempts/s     |
| iis.website_isapi_extension_requests_count                           |       website        |                                                                                       isapi                                                                                        |     requests      |
| iis.website_isapi_extension_requests_rate                            |       website        |                                                                                       isapi                                                                                        |    requests/s     |
| iis.website_ftp_file_transfer_rate                                   |       website        |                                                                                   received, sent                                                                                   |      files/s      |
| iis.website_logon_attempts_rate                                      |       website        |                                                                                       logon                                                                                        |    attempts/s     |
| iis.website_errors_rate                                              |       website        |                                                                        document_locked, document_not_found                                                                         |     errors/s      |
| iis.website_uptime                                                   |       website        |                                                                        document_locked, document_not_found                                                                         |      seconds      |
| mssql.instance_accessmethods_page_splits                             |    mssql instance    |                                                                                        page                                                                                        |     splits/s      |
| mssql.instance_cache_hit_ratio                                       |    mssql instance    |                                                                                     hit_ratio                                                                                      |    percentage     |
| mssql.instance_bufman_checkpoint_pages                               |    mssql instance    |                                                                                      flushed                                                                                       |      pages/s      |
| mssql.instance_bufman_page_life_expectancy                           |    mssql instance    |                                                                                  life_expectancy                                                                                   |      seconds      |
| mssql.instance_bufman_iops                                           |    mssql instance    |                                                                                   read, written                                                                                    |       iops        |
| mssql.instance_blocked_processes                                     |    mssql instance    |                                                                                      blocked                                                                                       |     processes     |
| mssql.instance_user_connection                                       |    mssql instance    |                                                                                        user                                                                                        |    connections    |
| mssql.instance_locks_lock_wait                                       |    mssql instance    |                                   alloc_unit, application, database, extent, file, hobt, key, metadata, oib, object, page, rid, row_group, xact                                    |      locks/s      |
| mssql.instance_locks_deadlocks                                       |    mssql instance    |                                   alloc_unit, application, database, extent, file, hobt, key, metadata, oib, object, page, rid, row_group, xact                                    |      locks/s      |
| mssql.instance_memmgr_connection_memory_bytes                        |    mssql instance    |                                                                                       memory                                                                                       |       bytes       |
| mssql.instance_memmgr_external_benefit_of_memory                     |    mssql instance    |                                                                                      benefit                                                                                       |       bytes       |
| mssql.instance_memmgr_pending_memory_grants                          |    mssql instance    |                                                                                      pending                                                                                       |     processes     |
| mssql.instance_memmgr_server_memory                                  |    mssql instance    |                                                                                       memory                                                                                       |       bytes       |
| mssql.instance_sql_errors                                            |    mssql instance    |                                                                      db_offline, info, kill_connection, user                                                                       |      errors       |
| mssql.instance_sqlstats_auto_parameterization_attempts               |    mssql instance    |                                                                                       failed                                                                                       |    attempts/s     |
| mssql.instance_sqlstats_batch_requests                               |    mssql instance    |                                                                                       batch                                                                                        |    requests/s     |
| mssql.instance_sqlstats_safe_auto_parameterization_attempts          |    mssql instance    |                                                                                        safe                                                                                        |    attempts/s     |
| mssql.instance_sqlstats_sql_compilations                             |    mssql instance    |                                                                                    compilations                                                                                    |  compilations/s   |
| mssql.instance_sqlstats_sql_recompilations                           |    mssql instance    |                                                                                     recompiles                                                                                     |   recompiles/s    |
| mssql.database_active_transactions                                   |       database       |                                                                                       active                                                                                       |   transactions    |
| mssql.database_backup_restore_operations                             |       database       |                                                                                       backup                                                                                       |   operations/s    |
| mssql.database_data_files_size                                       |       database       |                                                                                        size                                                                                        |       bytes       |
| mssql.database_log_flushed                                           |       database       |                                                                                      flushed                                                                                       |      bytes/s      |
| mssql.database_log_flushes                                           |       database       |                                                                                        log                                                                                         |     flushes/s     |
| mssql.database_transactions                                          |       database       |                                                                                    transactions                                                                                    |  transactions/s   |
| mssql.instance_write_transactions                                    |       database       |                                                                                       write                                                                                        |  transactions/s   |
| ad.database_operations                                               |        global        |                                                                            add, delete, modify, recycle                                                                            |   operations/s    |
| ad.directory_operations                                              |        global        |                                                                                read, write, search                                                                                 |   operations/s    |
| ad.name_cache_lookups                                                |        global        |                                                                                      lookups                                                                                       |     lookups/s     |
| ad.name_cache_hits                                                   |        global        |                                                                                        hits                                                                                        |      hits/s       |
| ad.atq_average_request_latency                                       |        global        |                                                                                        time                                                                                        |      seconds      |
| ad.atq_outstanding_requests                                          |        global        |                                                                                    outstanding                                                                                     |     requests      |
| ad.dra_replication_intersite_compressed_traffic                      |        global        |                                                                                 inbound, outbound                                                                                  |      bytes/s      |
| ad.dra_replication_intrasite_compressed_traffic                      |        global        |                                                                                 inbound, outbound                                                                                  |      bytes/s      |
| ad.dra_replication_sync_objects_remaining                            |        global        |                                                                                 inbound, outbound                                                                                  |      objects      |
| ad.dra_replication_objects_filtered                                  |        global        |                                                                                 inbound, outbound                                                                                  |     objects/s     |
| ad.dra_replication_properties_updated                                |        global        |                                                                                 inbound, outbound                                                                                  |   properties/s    |
| ad.dra_replication_properties_filtered                               |        global        |                                                                                 inbound, outbound                                                                                  |   properties/s    |
| ad.dra_replication_pending_syncs                                     |        global        |                                                                                      pending                                                                                       |       syncs       |
| ad.dra_replication_sync_requests                                     |        global        |                                                                                      requests                                                                                      |    requests/s     |
| ad.ds_threads                                                        |        global        |                                                                                       in_use                                                                                       |      threads      |
| ad.ldap_last_bind_time                                               |        global        |                                                                                     last_bind                                                                                      |      seconds      |
| ad.binds                                                             |        global        |                                                                                       binds                                                                                        |      binds/s      |
| ad.ldap_searches                                                     |        global        |                                                                                      searches                                                                                      |    searches/s     |
| adcs.cert_template_requests                                          | certificate template |                                                                                      requests                                                                                      |    requests/s     |
| adcs.cert_template_failed_requests                                   | certificate template |                                                                                       failed                                                                                       |    requests/s     |
| adcs.cert_template_issued_requests                                   | certificate template |                                                                                       issued                                                                                       |    requests/s     |
| adcs.cert_template_pending_requests                                  | certificate template |                                                                                      pending                                                                                       |    requests/s     |
| adcs.cert_template_request_processing_time                           | certificate template |                                                                                  processing_time                                                                                   |      seconds      |
| adcs.cert_template_retrievals                                        | certificate template |                                                                                     retrievals                                                                                     |   retrievals/s    |
| adcs.cert_template_retrieval_processing_time                         | certificate template |                                                                                  processing_time                                                                                   |      seconds      |
| adcs.cert_template_request_cryptographic_signing_time                | certificate template |                                                                                    singing_time                                                                                    |      seconds      |
| adcs.cert_template_request_policy_module_processing                  | certificate template |                                                                                  processing_time                                                                                   |      seconds      |
| adcs.cert_template_challenge_responses                               | certificate template |                                                                                     challenge                                                                                      |    responses/s    |
| adcs.cert_template_challenge_response_processing_time                | certificate template |                                                                                  processing_time                                                                                   |      seconds      |
| adcs.cert_template_signed_certificate_timestamp_lists                | certificate template |                                                                                     processed                                                                                      |      lists/s      |
| adcs.cert_template_signed_certificate_timestamp_list_processing_time | certificate template |                                                                                  processing_time                                                                                   |      seconds      |
| adfs.ad_login_connection_failures                                    |        global        |                                                                                     connection                                                                                     |    failures/s     |
| adfs.certificate_authentications                                     |        global        |                                                                                  authentications                                                                                   | authentications/s |
| adfs.db_artifact_failures                                            |        global        |                                                                                     connection                                                                                     |    failures/s     |
| adfs.db_artifact_query_time_seconds                                  |        global        |                                                                                     query_time                                                                                     |     seconds/s     |
| adfs.db_config_failures                                              |        global        |                                                                                     connection                                                                                     |    failures/s     |
| adfs.db_config_query_time_seconds                                    |        global        |                                                                                     query_time                                                                                     |     seconds/s     |
| adfs.device_authentications                                          |        global        |                                                                                  authentications                                                                                   | authentications/s |
| adfs.external_authentications                                        |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs.federated_authentications                                       |        global        |                                                                                  authentications                                                                                   | authentications/s |
| adfs.federation_metadata_requests                                    |        global        |                                                                                      requests                                                                                      |    requests/s     |
| adfs.oauth_authorization_requests                                    |        global        |                                                                                      requests                                                                                      |    requests/s     |
| adfs.oauth_client_authentications                                    |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs.oauth_client_credentials_requests                               |        global        |                                                                                  success, failure                                                                                  |    requests/s     |
| adfs.oauth_client_privkey_jwt_authentications                        |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs.oauth_client_secret_basic_authentications                       |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs.oauth_client_secret_post_authentications                        |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs.oauth_client_windows_authentications                            |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs.oauth_logon_certificate_requests                                |        global        |                                                                                  success, failure                                                                                  |    requests/s     |
| adfs.oauth_password_grant_requests                                   |        global        |                                                                                  success, failure                                                                                  |    requests/s     |
| adfs.oauth_token_requests_success                                    |        global        |                                                                                      success                                                                                       |    requests/s     |
| adfs.passive_requests                                                |        global        |                                                                                      passive                                                                                       |    requests/s     |
| adfs.passport_authentications                                        |        global        |                                                                                      passport                                                                                      | authentications/s |
| adfs.password_change_requests                                        |        global        |                                                                                  success, failure                                                                                  |    requests/s     |
| adfs.samlp_token_requests_success                                    |        global        |                                                                                      success                                                                                       |    requests/s     |
| adfs.sso_authentications                                             |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs.token_requests                                                  |        global        |                                                                                      requests                                                                                      |    requests/s     |
| adfs.userpassword_authentications                                    |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs.windows_integrated_authentications                              |        global        |                                                                                  authentications                                                                                   | authentications/s |
| adfs.wsfed_token_requests_success                                    |        global        |                                                                                      success                                                                                       |    requests/s     |
| adfs.wstrust_token_requests_success                                  |        global        |                                                                                      success                                                                                       |    requests/s     |
| netframework.clrexception_thrown                                     |       process        |                                                                                     exceptions                                                                                     |   exceptions/s    |
| netframework.clrexception_filters                                    |       process        |                                                                                      filters                                                                                       |     filters/s     |
| netframework.clrexception_finallys                                   |       process        |                                                                                      finallys                                                                                      |    finallys/s     |
| netframework.clrexception_throw_to_catch_depth                       |       process        |                                                                                     traversed                                                                                      |  stack_frames/s   |
| netframework.clrinterop_com_callable_wrappers                        |       process        |                                                                               com_callable_wrappers                                                                                |       ccw/s       |
| netframework.clrinterop_interop_marshallings                         |       process        |                                                                                    marshallings                                                                                    |  marshallings/s   |
| netframework.clrinterop_interop_stubs_created                        |       process        |                                                                                      created                                                                                       |      stubs/s      |
| netframework.clrjit_methods                                          |       process        |                                                                                    jit-compiled                                                                                    |     methods/s     |
| netframework.clrjit_time                                             |       process        |                                                                                        time                                                                                        |    percentage     |
| netframework.clrjit_standard_failures                                |       process        |                                                                                      failures                                                                                      |    failures/s     |
| netframework.clrjit_il_bytes                                         |       process        |                                                                                   compiled_msil                                                                                    |      bytes/s      |
| netframework.clrloading_loader_heap_size                             |       process        |                                                                                     committed                                                                                      |       bytes       |
| netframework.clrloading_appdomains_loaded                            |       process        |                                                                                       loaded                                                                                       |     domain/s      |
| netframework.clrloading_appdomains_unloaded                          |       process        |                                                                                      unloaded                                                                                      |     domain/s      |
| netframework.clrloading_assemblies_loaded                            |       process        |                                                                                       loaded                                                                                       |   assemblies/s    |
| netframework.clrloading_classes_loaded                               |       process        |                                                                                       loaded                                                                                       |     classes/s     |
| netframework.clrloading_class_load_failures                          |       process        |                                                                                     class_load                                                                                     |    failures/s     |
| netframework.clrlocksandthreads_queue_length                         |       process        |                                                                                      threads                                                                                       |     threads/s     |
| netframework.clrlocksandthreads_current_logical_threads              |       process        |                                                                                      logical                                                                                       |      threads      |
| netframework.clrlocksandthreads_current_physical_threads             |       process        |                                                                                      physical                                                                                      |      threads      |
| netframework.clrlocksandthreads_recognized_threads                   |       process        |                                                                                      threads                                                                                       |     threads/s     |
| netframework.clrlocksandthreads_contentions                          |       process        |                                                                                    contentions                                                                                     |   contentions/s   |
| netframework.clrmemory_allocated_bytes                               |       process        |                                                                                     allocated                                                                                      |      bytes/s      |
| netframework.clrmemory_finalization_survivors                        |       process        |                                                                                      survived                                                                                      |      objects      |
| netframework.clrmemory_heap_size                                     |       process        |                                                                                        heap                                                                                        |       bytes       |
| netframework.clrmemory_promoted                                      |       process        |                                                                                      promoted                                                                                      |       bytes       |
| netframework.clrmemory_number_gc_handles                             |       process        |                                                                                        used                                                                                        |      handles      |
| netframework.clrmemory_collections                                   |       process        |                                                                                         gc                                                                                         |       gc/s        |
| netframework.clrmemory_induced_gc                                    |       process        |                                                                                         gc                                                                                         |       gc/s        |
| netframework.clrmemory_number_pinned_objects                         |       process        |                                                                                       pinned                                                                                       |      objects      |
| netframework.clrmemory_number_sink_blocks_in_use                     |       process        |                                                                                        used                                                                                        |      blocks       |
| netframework.clrmemory_committed                                     |       process        |                                                                                     committed                                                                                      |       bytes       |
| netframework.clrmemory_reserved                                      |       process        |                                                                                      reserved                                                                                      |       bytes       |
| netframework.clrmemory_gc_time                                       |       process        |                                                                                        time                                                                                        |    percentage     |
| netframework.clrremoting_channels                                    |       process        |                                                                                     registered                                                                                     |    channels/s     |
| netframework.clrremoting_context_bound_classes_loaded                |       process        |                                                                                       loaded                                                                                       |      classes      |
| netframework.clrremoting_context_bound_objects                       |       process        |                                                                                     allocated                                                                                      |     objects/s     |
| netframework.clrremoting_context_proxies                             |       process        |                                                                                      objects                                                                                       |     objects/s     |
| netframework.clrremoting_contexts                                    |       process        |                                                                                      contexts                                                                                      |     contexts      |
| netframework.clrremoting_remote_calls                                |       process        |                                                                                        rpc                                                                                         |      calls/s      |
| netframework.clrsecurity_link_time_checks                            |       process        |                                                                                      linktime                                                                                      |     checks/s      |
| netframework.clrsecurity_checks_time                                 |       process        |                                                                                        time                                                                                        |    percentage     |
| netframework.clrsecurity_stack_walk_depth                            |       process        |                                                                                       stack                                                                                        |       depth       |
| netframework.clrsecurity_runtime_checks                              |       process        |                                                                                      runtime                                                                                       |     checks/s      |
| exchange.activesync_ping_cmds_pending                                |        global        |                                                                                      pending                                                                                       |     commands      |
| exchange.activesync_requests                                         |        global        |                                                                                      received                                                                                      |    requests/s     |
| exchange.activesync_sync_cmds                                        |        global        |                                                                                     processed                                                                                      |    commands/s     |
| exchange.autodiscover_requests                                       |        global        |                                                                                     processed                                                                                      |    requests/s     |
| exchange.avail_service_requests                                      |        global        |                                                                                      serviced                                                                                      |    requests/s     |
| exchange.owa_current_unique_users                                    |        global        |                                                                                     logged-in                                                                                      |       users       |
| exchange.owa_requests_total                                          |        global        |                                                                                      handled                                                                                       |    requests/s     |
| exchange.rpc_active_user_count                                       |        global        |                                                                                       active                                                                                       |       users       |
| exchange.rpc_avg_latency                                             |        global        |                                                                                      latency                                                                                       |      seconds      |
| exchange.rpc_connection_count                                        |        global        |                                                                                    connections                                                                                     |    connections    |
| exchange.rpc_operations                                              |        global        |                                                                                     operations                                                                                     |   operations/s    |
| exchange.rpc_requests                                                |        global        |                                                                                     processed                                                                                      |     requests      |
| exchange.rpc_user_count                                              |        global        |                                                                                       users                                                                                        |       users       |
| exchange.transport_queues_active_mail_box_delivery                   |        global        |                                                                               low, high,none,normal                                                                                |    messages/s     |
| exchange.transport_queues_external_active_remote_delivery            |        global        |                                                                               low, high,none,normal                                                                                |    messages/s     |
| exchange.transport_queues_external_largest_delivery                  |        global        |                                                                               low, high,none,normal                                                                                |    messages/s     |
| exchange.transport_queues_internal_active_remote_delivery            |        global        |                                                                               low, high,none,normal                                                                                |    messages/s     |
| exchange.transport_queues_internal_largest_delivery                  |        global        |                                                                               low, high,none,normal                                                                                |    messages/s     |
| exchange.transport_queues_retry_mailbox_delivery                     |        global        |                                                                               low, high,none,normal                                                                                |    messages/s     |
| exchange.transport_queues_poison                                     |        global        |                                                                               low, high,none,normal                                                                                |    messages/s     |
| exchange.workload_active_tasks                                       |  exchange workload   |                                                                                       active                                                                                       |       tasks       |
| exchange.workload_completed_tasks                                    |  exchange workload   |                                                                                     completed                                                                                      |      tasks/s      |
| exchange.workload_queued_tasks                                       |  exchange workload   |                                                                                       queued                                                                                       |      tasks/s      |
| exchange.workload_yielded_tasks                                      |  exchange workload   |                                                                                      yielded                                                                                       |      tasks/s      |
| exchange.workload_activity_status                                    |  exchange workload   |                                                                                   active, paused                                                                                   |      status       |
| exchange.ldap_long_running_ops_per_sec                               |     ldap process     |                                                                                    long-running                                                                                    |   operations/s    |
| exchange.ldap_read_time                                              |     ldap process     |                                                                                        read                                                                                        |      seconds      |
| exchange.ldap_search_time                                            |     ldap process     |                                                                                       search                                                                                       |      seconds      |
| exchange.ldap_write_time                                             |     ldap process     |                                                                                       write                                                                                        |      seconds      |
| exchange.ldap_timeout_errors                                         |     ldap process     |                                                                                      timeout                                                                                       |     errors/s      |
| exchange.http_proxy_avg_auth_latency                                 |      http proxy      |                                                                                      latency                                                                                       |      seconds      |
| exchange.http_proxy_avg_cas_processing_latency_sec                   |      http proxy      |                                                                                      latency                                                                                       |      seconds      |
| exchange.http_proxy_mailbox_proxy_failure_rate                       |      http proxy      |                                                                                      failures                                                                                      |    percentage     |
| exchange.http_proxy_mailbox_server_locator_avg_latency_sec           |      http proxy      |                                                                                      latency                                                                                       |      seconds      |
| exchange.http_proxy_outstanding_proxy_requests                       |      http proxy      |                                                                                    outstanding                                                                                     |     requests      |
| exchange.http_proxy_requests_total                                   |      http proxy      |                                                                                     processed                                                                                      |    requests/s     |

## Configuration

Edit the `go.d/windows.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md#the-netdata-config-directory), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/windows.conf
```

Needs only `url` to `windows_exporter` metrics endpoint. Here is an example for 2 instances:

```yaml
jobs:
  - name: win_server1
    url: http://203.0.113.10:9182/metrics

  - name: win_server2
    url: http://203.0.113.11:9182/metrics
```

Hint: Use friendly server names for job names, as these will appear as "instances" in Netdata Cloud charts
and on the right side menu of the agent UI charts.

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/windows.conf).

### Virtual Nodes

Netdata’s new virtual nodes functionality allows you to define nodes in configuration files and have them be treated as regular nodes in all of the UI, dashboards, tabs, filters etc. For example, you can create a virtual node each for all your Windows machines and monitor them as discrete entities. Virtual nodes can help you simplify your infrastructure monitoring and focus on the individual node that matters.

To define your windows server a virtual node you need to:

  * Define virtual nodes in `/etc/netdata/vnodes/vnodes.conf`

    ```yaml
    - hostname: win_server1
      guid: <value>
    ```
    Just remember to use a valid guid (On Linux you can use `uuidgen` command to generate one, on Windows just use the `[guid]::NewGuid()` command in PowerShell)
    
  * Add the vnode config to the windows monitoring job we created earlier, see higlighted line below:
    ```yaml
      jobs:
        - name: win_server1
          vnode: win_server1
          url: http://203.0.113.10:9182/metrics
    ```

## Troubleshooting

To troubleshoot issues with the `windows` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m windows
  ```

