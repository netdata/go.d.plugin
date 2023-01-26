<!--
title: "Windows machine monitoring with Netdata"
description: "Monitor the health and performance of Windows machines with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/wmi/README.md"
sidebar_label: "Windows machines"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/System metrics"
-->

# Windows machine monitoring with Netdata

This module will monitor one or more Windows machines, using
the [windows_exporter](https://github.com/prometheus-community/windows_exporter).

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

## Requirements

Netdata monitors Windows hosts by utilizing the
[Prometheus exporter for Windows machines](https://github.com/prometheus-community/windows_exporter), a native Windows
agent running on each host.

To quickly test Netdata directly on a Windows machine, you can use
the [Netdata MSI installer](https://github.com/netdata/msi-installer#instructions). The installer runs Netdata in a
custom WSL deployment. WSL was not designed for production environments, so we do not recommend using the installer in
production.

For production use, you need to install Netdata on one or more nodes running Linux:

- Install the
  latest [Prometheus exporter for Windows](https://github.com/prometheus-community/windows_exporter/releases)
  on every Windows host you want to monitor.
- Get the installation commands from [Netdata Cloud](https://app.netdata.cloud) and install Netdata on one or more Linux
  nodes.
- Configure each Netdata instance to collect data remotely, from several Windows hosts. Just add one job
  for each host to  `wmi.conf`, as shown in the [configuration section](#configuration).
- [Optional] [Disable all plugins](https://learn.netdata.cloud/docs/configure/common-changes#disable-a-collector-or-plugin)
  except for go.d in `netdata.conf`, so that you only see Windows metrics.
- [Optional] Set up [replication](https://learn.netdata.cloud/docs/agent/streaming), for high availability.

Automated charts and alerts for your entire Windows infrastructure will be automatically generated.
Each Windows host (data collection job) will be identifiable as an "instance" in the Netdata Cloud charts.

## Metrics

All metrics have prefix.

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

| Metric                                                               |        Scope         |                                                                                     Dimensions                                                                                     |       Units       |
|----------------------------------------------------------------------|:--------------------:|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-----------------:|
| cpu_utilization_total                                                |        global        |                                                                          dpc, user, privileged, interrupt                                                                          |    percentage     |
| cpu_core_utilization                                                 |       cpu core       |                                                                          dpc, user, privileged, interrupt                                                                          |    percentage     |
| cpu_core_interrupts                                                  |       cpu core       |                                                                                     interrupts                                                                                     |   interrupts/s    |
| cpu_core_dpcs                                                        |       cpu core       |                                                                                        dpcs                                                                                        |      dpcs/s       |
| cpu_core_cstate                                                      |       cpu core       |                                                                                     c1, c2, c3                                                                                     |    percentage     |
| memory_utilization                                                   |        global        |                                                                                  available, used                                                                                   |       bytes       |
| memory_utilization                                                   |        global        |                                                                                  available, used                                                                                   |        KiB        |
| memory_page_faults                                                   |        global        |                                                                                    page_faults                                                                                     |     events/s      |
| memory_swap_utilization                                              |        global        |                                                                                  available, used                                                                                   |       bytes       |
| memory_swap_operations                                               |        global        |                                                                                    read, write                                                                                     |   operations/s    |
| memory_swap_pages                                                    |        global        |                                                                                   read, written                                                                                    |      pages/s      |
| memory_cached                                                        |        global        |                                                                                       cached                                                                                       |        KiB        |
| memory_cache_faults                                                  |        global        |                                                                                    cache_faults                                                                                    |     events/s      |
| memory_system_pool                                                   |        global        |                                                                                  paged, non-paged                                                                                  |       bytes       |
| logical_disk_utilization                                             |     logical disk     |                                                                                     free, used                                                                                     |       bytes       |
| logical_disk_bandwidth                                               |     logical disk     |                                                                                    read, write                                                                                     |      bytes/s      |
| logical_disk_operations                                              |     logical disk     |                                                                                   reads, writes                                                                                    |   operations/s    |
| logical_disk_latency                                                 |     logical disk     |                                                                                    read, write                                                                                     |      seconds      |
| net_nic_bandwidth                                                    |    network device    |                                                                                   received, sent                                                                                   |    kilobits/s     |
| net_nic_packets                                                      |    network device    |                                                                                   received, sent                                                                                   |     packets/s     |
| net_nic_errors                                                       |    network device    |                                                                                 inbound, outbound                                                                                  |     errors/s      |
| net_nic_discarded                                                    |    network device    |                                                                                 inbound, outbound                                                                                  |    discards/s     |
| tcp_conns_established                                                |        global        |                                                                                     ipv4, ipv6                                                                                     |    connections    |
| tcp_conns_active                                                     |        global        |                                                                                     ipv4, ipv6                                                                                     |   connections/s   |
| tcp_conns_passive                                                    |        global        |                                                                                     ipv4, ipv6                                                                                     |   connections/s   |
| tcp_conns_failures                                                   |        global        |                                                                                     ipv4, ipv6                                                                                     |    failures/s     |
| tcp_conns_resets                                                     |        global        |                                                                                     ipv4, ipv6                                                                                     |     resets/s      |
| tcp_segments_received                                                |        global        |                                                                                     ipv4, ipv6                                                                                     |    segments/s     |
| tcp_segments_sent                                                    |        global        |                                                                                     ipv4, ipv6                                                                                     |    segments/s     |
| tcp_segments_retransmitted                                           |        global        |                                                                                     ipv4, ipv6                                                                                     |    segments/s     |
| os_processes                                                         |        global        |                                                                                     processes                                                                                      |      number       |
| os_users                                                             |        global        |                                                                                       users                                                                                        |       users       |
| os_visible_memory_usage                                              |        global        |                                                                                     free, used                                                                                     |       bytes       |
| os_paging_files_usage                                                |        global        |                                                                                     free, used                                                                                     |       bytes       |
| system_threads                                                       |        global        |                                                                                      threads                                                                                       |      number       |
| system_uptime                                                        |        global        |                                                                                        time                                                                                        |      seconds      |
| logon_type_sessions                                                  |        global        | system, interactive, network, batch, service, proxy, unlock, network_clear_text, new_credentials, remote_interactive, cached_interactive, cached_remote_interactive, cached_unlock |      seconds      |
| thermalzone_temperature                                              |        global        |                                                                         <i>a dimension per thermalzone</i>                                                                         |      celsius      |
| processes_cpu_utilization                                            |        global        |                                                                           <i>a dimension per process</i>                                                                           |    percentage     |
| processes_handles                                                    |        global        |                                                                           <i>a dimension per process</i>                                                                           |      handles      |
| processes_io_bytes                                                   |        global        |                                                                           <i>a dimension per process</i>                                                                           |      bytes/s      |
| processes_io_operations                                              |        global        |                                                                           <i>a dimension per process</i>                                                                           |   operations/s    |
| processes_page_faults                                                |        global        |                                                                           <i>a dimension per process</i>                                                                           |    pgfaults/s     |
| processes_page_file_bytes                                            |        global        |                                                                           <i>a dimension per process</i>                                                                           |       bytes       |
| processes_pool_bytes                                                 |        global        |                                                                           <i>a dimension per process</i>                                                                           |       bytes       |
| processes_threads                                                    |        global        |                                                                           <i>a dimension per process</i>                                                                           |      threads      |
| service_state                                                        |       service        |                                          running, stopped, start_pending, stop_pending, continue_pending, pause_pending, paused, unknown                                           |       state       |
| service_status                                                       |       service        |                                 ok, error, unknown, degraded, pred_fail, starting, stopping, service, stressed, nonrecover, no_contact, lost_comm                                  |      status       |
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
| mssql.instance_memmgr_connection_memory_bytes                        |    mssql instance    |                                                                                       memory                                                                                       |       bytes       |
| mssql.instance_memmgr_external_benefit_of_memory                     |    mssql instance    |                                                                                      benefit                                                                                       |       bytes       |
| mssql.instance_memmgr_pending_memory_grants                          |    mssql instance    |                                                                                      pending                                                                                       |     processes     |
| mssql.instance_memmgr_server_memory                                  |    mssql instance    |                                                                                       memory                                                                                       |       bytes       |
| mssql.instance_sql_errors                                            |    mssql instance    |                                                                      db_offline, info, kill_connection, user                                                                       |      errors       |
| mssql.instance_sqlstats_auto_parameterization_attempts               |    mssql instance    |                                                                                       failed                                                                                       |    attempts/s     |
| mssql.instance_sqlstats_batch_requests                               |    mssql instance    |                                                                                       batch                                                                                        |    requests/s     |
| mssql.instance_sqlstats_safe_auto_parameterization_attempts          |    mssql instance    |                                                                                        safe                                                                                        |    attempts/s     |
| mssql_instance_sqlstats_sql_compilations                             |    mssql instance    |                                                                                    compilations                                                                                    |  compilations/s   |
| mssql.instance_sqlstats_sql_recompilations                           |    mssql instance    |                                                                                     recompiles                                                                                     |   recompiles/s    |
| mssql.database_active_transactions                                   |       database       |                                                                                       active                                                                                       |   transactions    |
| mssql.database_backup_restore_operations                             |       database       |                                                                                       backup                                                                                       |   operations/s    |
| mssql.database_data_files_size                                       |       database       |                                                                                        size                                                                                        |       bytes       |
| mssql.database_log_flushed                                           |       database       |                                                                                      flushed                                                                                       |      bytes/s      |
| mssql.database_log_flushes                                           |       database       |                                                                                        log                                                                                         |     flushes/s     |
| mssql.database_transactions                                          |       database       |                                                                                    transactions                                                                                    |  transactions/s   |
| mssql.instance_write_transactions                                    |       database       |                                                                                       write                                                                                        |  transactions/s   |
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
| adfs_ad_login_connection_failures                                    |        global        |                                                                                     connection                                                                                     |    failures/s     |
| adfs_certificate_authentications                                     |        global        |                                                                                  authentications                                                                                   | authentications/s |
| adfs_db_artifact_failures                                            |        global        |                                                                                     connection                                                                                     |    failures/s     |
| adfs_db_artifact_query_time_seconds                                  |        global        |                                                                                     query_time                                                                                     |     seconds/s     |
| adfs_db_config_failures                                              |        global        |                                                                                     connection                                                                                     |    failures/s     |
| adfs_db_config_query_time_seconds                                    |        global        |                                                                                     query_time                                                                                     |     seconds/s     |
| adfs_device_authentications                                          |        global        |                                                                                  authentications                                                                                   | authentications/s |
| adfs_external_authentications                                        |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs_federated_authentications                                       |        global        |                                                                                  authentications                                                                                   | authentications/s |
| adfs_federation_metadata_requests                                    |        global        |                                                                                      requests                                                                                      |    requests/s     |
| adfs_oauth_authorization_requests                                    |        global        |                                                                                      requests                                                                                      |    requests/s     |
| adfs_oauth_client_authentications                                    |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs_oauth_client_credentials_requests                               |        global        |                                                                                  success, failure                                                                                  |    requests/s     |
| adfs_oauth_client_privkey_jwt_authentications                        |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs_oauth_client_secret_basic_authentications                       |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs_oauth_client_secret_post_authentications                        |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs_oauth_client_windows_authentications                            |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs_oauth_logon_certificate_requests                                |        global        |                                                                                  success, failure                                                                                  |    requests/s     |
| adfs_oauth_password_grant_requests                                   |        global        |                                                                                  success, failure                                                                                  |    requests/s     |
| adfs_oauth_token_requests_success                                    |        global        |                                                                                      success                                                                                       |    requests/s     |
| adfs_passive_requests                                                |        global        |                                                                                      passive                                                                                       |    requests/s     |
| adfs_passport_authentications                                        |        global        |                                                                                      passport                                                                                      | authentications/s |
| adfs_password_change_requests                                        |        global        |                                                                                  success, failure                                                                                  |    requests/s     |
| adfs_samlp_token_requests_success                                    |        global        |                                                                                      success                                                                                       |    requests/s     |
| adfs_sso_authentications                                             |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs_token_requests                                                  |        global        |                                                                                      requests                                                                                      |    requests/s     |
| adfs_userpassword_authentications                                    |        global        |                                                                                  success, failure                                                                                  | authentications/s |
| adfs_windows_integrated_authentications                              |        global        |                                                                                  authentications                                                                                   | authentications/s |
| adfs_wsfed_token_requests_success                                    |        global        |                                                                                      success                                                                                       |    requests/s     |
| adfs_wstrust_token_requests_success                                  |        global        |                                                                                      success                                                                                       |    requests/s     |

## Configuration

Edit the `go.d/wmi.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/wmi.conf
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
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/wmi.conf).

## Troubleshooting

To troubleshoot issues with the `wmi` collector, run the `go.d.plugin` with the debug option enabled. The output should
give you clues as to why the collector isn't working.

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
  ./go.d.plugin -d -m wmi
  ```
