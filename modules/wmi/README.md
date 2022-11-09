<!--
title: "Windows machine monitoring with Netdata"
description: "Monitor the health and performance of Windows machines with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/wmi/README.md
sidebar_label: "Windows machines"
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

## Requirements

There are two ways to monitor Windows with Netdata. Install Netdata over WSL on each host, or remotely collect
data from one or more centralized agents, running on dedicated Linux machines.

### Netdata on each Windows machine

Use the [Netdata MSI installer](https://github.com/netdata/msi-installer#instructions).
  
### Remote data collection

- Install the
  latest [Prometheus exporter for Windows](https://github.com/prometheus-community/windows_exporter/releases)
  on every Windows host you want to monitor.
- Install Netdata on one or more Linux servers.
- Configure each Netdata instance to collect data remotely, from several Windows hosts. Just add one job
  for each host to  `wmi.conf`, as shown in the [configuration section](#configuration).

## Metrics

All metrics have "wmi." prefix.

| Metric                                     |     Scope      |                                                                                     Dimensions                                                                                     |     Units     |
|--------------------------------------------|:--------------:|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-------------:|
| cpu_utilization_total                      |     global     |                                                                          dpc, user, privileged, interrupt                                                                          |  percentage   |
| cpu_core_utilization                       |    cpu core    |                                                                          dpc, user, privileged, interrupt                                                                          |  percentage   |
| cpu_core_interrupts                        |    cpu core    |                                                                                     interrupts                                                                                     | interrupts/s  |
| cpu_core_dpcs                              |    cpu core    |                                                                                        dpcs                                                                                        |    dpcs/s     |
| cpu_core_cstate                            |    cpu core    |                                                                                     c1, c2, c3                                                                                     |  percentage   |
| memory_utilization                         |     global     |                                                                                  available, used                                                                                   |     bytes     |
| memory_utilization                         |     global     |                                                                                  available, used                                                                                   |      KiB      |
| memory_page_faults                         |     global     |                                                                                    page_faults                                                                                     |   events/s    |
| memory_swap_utilization                    |     global     |                                                                                  available, used                                                                                   |     bytes     |
| memory_swap_operations                     |     global     |                                                                                    read, write                                                                                     | operations/s  |
| memory_swap_pages                          |     global     |                                                                                   read, written                                                                                    |    pages/s    |
| memory_cached                              |     global     |                                                                                       cached                                                                                       |      KiB      |
| memory_cache_faults                        |     global     |                                                                                    cache_faults                                                                                    |   events/s    |
| memory_system_pool                         |     global     |                                                                                  paged, non-paged                                                                                  |     bytes     |
| logical_disk_utilization                   |  logical disk  |                                                                                     free, used                                                                                     |     bytes     |
| logical_disk_bandwidth                     |  logical disk  |                                                                                    read, write                                                                                     |    bytes/s    |
| logical_disk_operations                    |  logical disk  |                                                                                   reads, writes                                                                                    | operations/s  |
| logical_disk_latency                       |  logical disk  |                                                                                    read, write                                                                                     |    seconds    |
| net_nic_bandwidth                          | network device |                                                                                   received, sent                                                                                   |  kilobits/s   |
| net_nic_packets                            | network device |                                                                                   received, sent                                                                                   |   packets/s   |
| net_nic_errors                             | network device |                                                                                 inbound, outbound                                                                                  |   errors/s    |
| net_nic_discarded                          | network device |                                                                                 inbound, outbound                                                                                  |  discards/s   |
| tcp_conns_established                      |     global     |                                                                                     ipv4, ipv6                                                                                     |  connections  |
| tcp_conns_active                           |     global     |                                                                                     ipv4, ipv6                                                                                     | connections/s |
| tcp_conns_passive                          |     global     |                                                                                     ipv4, ipv6                                                                                     | connections/s |
| tcp_conns_failures                         |     global     |                                                                                     ipv4, ipv6                                                                                     |  failures/s   |
| tcp_conns_resets                           |     global     |                                                                                     ipv4, ipv6                                                                                     |   resets/s    |
| tcp_segments_received                      |     global     |                                                                                     ipv4, ipv6                                                                                     |  segments/s   |
| tcp_segments_sent                          |     global     |                                                                                     ipv4, ipv6                                                                                     |  segments/s   |
| tcp_segments_retransmitted                 |     global     |                                                                                     ipv4, ipv6                                                                                     |  segments/s   |
| os_processes                               |     global     |                                                                                     processes                                                                                      |    number     |
| os_users                                   |     global     |                                                                                       users                                                                                        |     users     |
| os_visible_memory_usage                    |     global     |                                                                                     free, used                                                                                     |     bytes     |
| os_paging_files_usage                      |     global     |                                                                                     free, used                                                                                     |     bytes     |
| system_threads                             |     global     |                                                                                      threads                                                                                       |    number     |
| system_uptime                              |     global     |                                                                                        time                                                                                        |    seconds    |
| logon_type_sessions                        |     global     | system, interactive, network, batch, service, proxy, unlock, network_clear_text, new_credentials, remote_interactive, cached_interactive, cached_remote_interactive, cached_unlock |    seconds    |
| thermalzone_temperature                    |     global     |                                                                         <i>a dimension per thermalzone</i>                                                                         |    celsius    |
| processes_cpu_utilization                  |     global     |                                                                           <i>a dimension per process</i>                                                                           |  percentage   |
| processes_handles                          |     global     |                                                                           <i>a dimension per process</i>                                                                           |    handles    |
| processes_io_bytes                         |     global     |                                                                           <i>a dimension per process</i>                                                                           |    bytes/s    |
| processes_io_operations                    |     global     |                                                                           <i>a dimension per process</i>                                                                           | operations/s  |
| processes_page_faults                      |     global     |                                                                           <i>a dimension per process</i>                                                                           |  pgfaults/s   |
| processes_page_file_bytes                  |     global     |                                                                           <i>a dimension per process</i>                                                                           |     bytes     |
| processes_pool_bytes                       |     global     |                                                                           <i>a dimension per process</i>                                                                           |     bytes     |
| processes_threads                          |     global     |                                                                           <i>a dimension per process</i>                                                                           |    threads    |
| service_state                              |    service     |                                          running, stopped, start_pending, stop_pending, continue_pending, pause_pending, paused, unknown                                           |     state     |
| service_status                             |    service     |                                 ok, error, unknown, degraded, pred_fail, starting, stopping, service, stressed, nonrecover, no_contact, lost_comm                                  |    status     |
| iis_website_traffic                        |    website     |                                                                                   received, sent                                                                                   |    bytes/s    |
| iis_website_requests_rate                  |    website     |                                                                                      requests                                                                                      |  requests/s   |
| iis_website_active_connections_count       |    website     |                                                                                       active                                                                                       |  connections  |
| iis_website_users_count                    |    website     |                                                                              anonymous, non_anonymous                                                                              |     users     |
| iis_website_connection_attempts_rate       |    website     |                                                                                     connection                                                                                     |  attempts/s   |
| iis_website_isapi_extension_requests_count |    website     |                                                                                       isapi                                                                                        |   requests    |
| iis_website_isapi_extension_requests_rate  |    website     |                                                                                       isapi                                                                                        |  requests/s   |
| iis_website_ftp_file_transfer_rate         |    website     |                                                                                   received, sent                                                                                   |    files/s    |
| iis_website_logon_attempts_rate            |    website     |                                                                                       logon                                                                                        |  attempts/s   |
| iis_website_errors_rate                    |    website     |                                                                        document_locked, document_not_found                                                                         |   errors/s    |
| iis_website_uptime                         |    website     |                                                                        document_locked, document_not_found                                                                         |    seconds    |

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
