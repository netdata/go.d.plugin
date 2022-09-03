<!--
title: "ScaleIO monitoring with Netdata"
description: "Monitor the health and performance of ScaleIO storage with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/scaleio/README.md
sidebar_label: "ScaleIO"
-->

# ScaleIO monitoring with Netdata

[`Dell EMC ScaleIO`](https://www.dellemc.com/en-us/storage/data-storage/software-defined-storage.htm) is a
software-defined storage product from Dell EMC that creates a server-based storage area network from local application
server storage using existing customer hardware or EMC servers.

This module monitors one or more `ScaleIO (VxFlex OS)` instances via VxFlex OS Gateway API.

It collects metrics for following `ScaleIO` components:

- System
- Storage Pool
- Sdc

`ScaleIO` module is tested on:

- VxFlex OS v2.6.1.1_113, REST API v2.5
- VxFlex OS v3.0.0.1_134, REST API v3.0

## Charts

It produces the following charts:

#### System

- Total Capacity in `KiB`
- Capacity In Use in `KiB`
- Capacity Usage in `KiB`
- Available For Volume Allocation in `KiB`
- Capacity Health State in `KiB`
- Primary Backend Bandwidth Total (Read and Write) in `KiB/s`
- Primary Backend Bandwidth in `KiB/s`
- Primary Backend IOPS Total (Read and Write) in `iops/s`
- Primary Backend IOPS in `iops/s`
- Primary Backend I/O Size Total (Read and Write) in `KiB`
- Rebalance in `KiB/s`
- Rebalance Pending Capacity in `KiB`
- Rebalance Approximate Time Until Finish in `seconds`
- Rebuild Bandwidth Total (Forward, Backward and Normal) in `KiB/s`
- Rebuild Pending Capacity Total (Forward, Backward and Normal) in `KiB`
- Components in `number`
- Volumes By Type in `number`
- Volumes By Mapping in `number`

#### Storage Pool

- Total Capacity in `KiB`
- Capacity In Use in `KiB`
- Capacity Usage in `KiB`
- Capacity Utilization in `percentage`
- Available For Volume Allocation in `KiB`
- Capacity Health State in `KiB`
- Components in `number`

#### SDC

- MDM Connection State in `boolean`
- Bandwidth in `KiB/s`
- IOPS in `iops/s`
- I/O Size in `KiB`
- Mapped Volumes in `volumes`

## Configuration

Edit the `go.d/scaleio.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/scaleio.conf
```

Needs only `url` of VxFlex OS Gateway API, MDM `username` and `password`. Here is an example for 2 instances:

```yaml
jobs:
  - name: local
    url: https://127.0.0.1
    username: admin
    password: password
    tls_skip_verify: yes  # self-signed certificate

  - name: remote
    url: https://203.0.113.10
    username: admin
    password: password
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/scaleio.conf).

## Troubleshooting

To troubleshoot issues with the `scaleio` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m scaleio
  ```

