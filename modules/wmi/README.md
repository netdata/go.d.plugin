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

On your Windows machine:
- Install the latest [Prometheus exporter for Windows](https://github.com/prometheus-community/windows_exporter/releases).
- Download the [netdata.msi](https://github.com/netdata/msi-installer/releases)
- Run netdata.msi directly, or with the options provided by Netdata Cloud. It will install WSL2 or WSL1 and run Netdata on your machine.
  If you do not want to install WSL on your Windows machines, you can configure a single Netdata instance to collect data remotely from 
  multiple Windows hosts. Just add multiple jobs to `wmi.conf`, as shown in the [configuration section](#configuration).

## Metrics

All metrics have "wmi." prefix.

| Metric                     |     Scope      |                                                                                     Dimensions                                                                                     |     Units     |
|----------------------------|:--------------:|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-------------:|
| cpu_utilization_total      |     global     |                                                                          dpc, user, privileged, interrupt                                                                          |  percentage   |
| cpu_dpcs                   |     global     |                                                                            <i>a dimension per core</i>                                                                             |    dpcs/s     |
| cpu_interrupts             |     global     |                                                                            <i>a dimension per core</i>                                                                             | interrupts/s  |
| cpu_utilization            |    cpu core    |                                                                          dpc, user, privileged, interrupt                                                                          |  percentage   |
| cpu_cstate                 |    cpu core    |                                                                                     c1, c2, c3                                                                                     |  percentage   |
| memory_utilization         |     global     |                                                                                  available, used                                                                                   |      KiB      |
| memory_page_faults         |     global     |                                                                                    page_faults                                                                                     |   events/s    |
| memory_swap_utilization    |     global     |                                                                                  available, used                                                                                   |      KiB      |
| memory_swap_operations     |     global     |                                                                                    read, write                                                                                     | operations/s  |
| memory_swap_pages          |     global     |                                                                                   read, written                                                                                    |    pages/s    |
| memory_cached              |     global     |                                                                                       cached                                                                                       |      KiB      |
| memory_cache_faults        |     global     |                                                                                    cache_faults                                                                                    |   events/s    |
| memory_system_pool         |     global     |                                                                                  paged, non-paged                                                                                  |      KiB      |
| net_bandwidth              | network device |                                                                                   received, sent                                                                                   |  kilobits/s   |
| net_packets                | network device |                                                                                   received, sent                                                                                   |   packets/s   |
| net_errors                 | network device |                                                                                 inbound, outbound                                                                                  |   errors/s    |
| net_discarded              | network device |                                                                                 inbound, outbound                                                                                  |  discards/s   |
| logical_disk_utilization   |  logical disk  |                                                                                     free, used                                                                                     |      KiB      |
| logical_disk_bandwidth     |  logical disk  |                                                                                    read, write                                                                                     |     KiB/s     |
| logical_disk_operations    |  logical disk  |                                                                                   reads, writes                                                                                    | operations/s  |
| logical_disk_latency       |  logical disk  |                                                                                    read, write                                                                                     | milliseconds  |
| os_processes               |     global     |                                                                                     processes                                                                                      |    number     |
| os_users                   |     global     |                                                                                       users                                                                                        |     users     |
| os_visible_memory_usage    |     global     |                                                                                     free, used                                                                                     |     bytes     |
| os_paging_files_usage      |     global     |                                                                                     free, used                                                                                     |     bytes     |
| system_threads             |     global     |                                                                                      threads                                                                                       |    number     |
| system_uptime              |     global     |                                                                                        time                                                                                        |    seconds    |
| logon_type_sessions        |     global     | system, interactive, network, batch, service, proxy, unlock, network_clear_text, new_credentials, remote_interactive, cached_interactive, cached_remote_interactive, cached_unlock |    seconds    |
| thermalzone_temperature    |     global     |                                                                         <i>a dimension per thermalzone</i>                                                                         |    celsius    |
| tcp_conns_established      |     global     |                                                                                     ipv4, ipv6                                                                                     |  connections  |
| tcp_conns_active           |     global     |                                                                                     ipv4, ipv6                                                                                     | connections/s |
| tcp_conns_passive          |     global     |                                                                                     ipv4, ipv6                                                                                     | connections/s |
| tcp_conns_failures         |     global     |                                                                                     ipv4, ipv6                                                                                     |  failures/s   |
| tcp_conns_resets           |     global     |                                                                                     ipv4, ipv6                                                                                     |   resets/s    |
| tcp_segments_received      |     global     |                                                                                     ipv4, ipv6                                                                                     |  segments/s   |
| tcp_segments_sent          |     global     |                                                                                     ipv4, ipv6                                                                                     |  segments/s   |
| tcp_segments_retransmitted |     global     |                                                                                     ipv4, ipv6                                                                                     |  segments/s   |
| processes_cpu_utilization  |     global     |                                                                           <i>a dimension per process</i>                                                                           |  percentage   |
| processes_handles          |     global     |                                                                           <i>a dimension per process</i>                                                                           |    handles    |
| processes_io_bytes         |     global     |                                                                           <i>a dimension per process</i>                                                                           |    bytes/s    |
| processes_io_operations    |     global     |                                                                           <i>a dimension per process</i>                                                                           | operations/s  |
| processes_page_faults      |     global     |                                                                           <i>a dimension per process</i>                                                                           |  pgfaults/s   |
| processes_page_file_bytes  |     global     |                                                                           <i>a dimension per process</i>                                                                           |     bytes     |
| processes_pool_bytes       |     global     |                                                                           <i>a dimension per process</i>                                                                           |     bytes     |
| processes_threads          |     global     |                                                                           <i>a dimension per process</i>                                                                           |    threads    |
| service_state              |    service     |                                          running, stopped, start_pending, stop_pending, continue_pending, pause_pending, paused, unknown                                           |     state     |
| service_status             |    service     |                                      ok, error, unknown, pred_fail, starting, stopping, service, stressed, nonrecover, no_contact, lost_comm                                       |    status     |

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
