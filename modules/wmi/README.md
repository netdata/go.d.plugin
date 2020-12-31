<!--
title: "Windows machine monitoring with Netdata"
description: "Monitor the health and performance of Windows machines with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/wmi/README.md
sidebar_label: "Windows machines"
-->

# Windows machine monitoring with Netdata

This module will monitor one or more Windows machines, using
the [windows_exporter](https://github.com/prometheus-community/windows_exporter).

Module collects metrics from the following collectors:

- cpu
- memory
- net
- logical_disk
- os
- system
- logon

Run `windowns_exporter` with these collectors:

> windows_exporter-0.13.0-amd64.exe --collectors.enabled="cpu,memory,net,logical_disk,os,system,logon"


Installation: please follow the [official guide](https://github.com/prometheus-community/windows_exporter#installation).

## Requirements

- `windows_exporter` version v0.13.0+

## Charts

### cpu

- Total CPU Utilization (all cores) in `percentage`
- Received and Serviced Deferred Procedure Calls (DPC) in `dpc/c`
- Received and Serviced Hardware Interrupts in `interrupts/s`
- CPU Utilization Per Core in `percentage`
- Time Spent in Low-Power Idle State Per Core in `percentage`

### memory

- Memory Utilization in `KiB`
- Memory Page Faults in `events/s`
- Swap Utilization in `KiB`
- Swap Operations in `operations/s`
- Swap Pages in `pages/s`
- Cached Data in `KiB`
- Cache Faults in `events/s`
- System Memory Pool in `KiB`

### network

- Bandwidth Per NIC in `kilobits/s`
- Packets Per NIC in `packets/s`
- Errors Per NIC in `errors/s`
- Discards Per NIC in `discards/s`

### disk

- Utilization Per Disk in `KiB`
- Bandwidth Per Disk in `KiB/s`
- Operations Per Disk in `operations/s`
- Average Read/Write Latency Disk in `milliseconds`

### system

- Processes in `number`
- Threads in `number`
- Uptime in `seconds`

### logon

- Active User Logon Sessions By Type in `sessions`

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

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m wmi
```
