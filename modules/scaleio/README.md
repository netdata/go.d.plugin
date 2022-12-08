<!--
title: "ScaleIO monitoring with Netdata"
description: "Monitor the health and performance of ScaleIO storage with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/scaleio/README.md"
sidebar_label: "ScaleIO"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Storage"
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

## Metrics

All metrics have "scaleio." prefix.

| Metric                                            |    Scope     |                                                  Dimensions                                                   |   Units    |
|---------------------------------------------------|:------------:|:-------------------------------------------------------------------------------------------------------------:|:----------:|
| system_capacity_total                             |    global    |                                                     total                                                     |    KiB     |
| system_capacity_in_use                            |    global    |                                                    in_use                                                     |    KiB     |
| system_capacity_usage                             |    global    |                                thick, decreased, thin, snapshot, spare, unused                                |    KiB     |
| system_capacity_available_volume_allocation       |    global    |                                                   available                                                   |    KiB     |
| system_capacity_health_state                      |    global    |                           protected, degraded, in_maintenance, failed, unavailable                            |    KiB     |
| system_workload_primary_bandwidth_total           |    global    |                                                     total                                                     |   KiB/s    |
| system_workload_primary_bandwidth                 |    global    |                                                  read, write                                                  |   KiB/s    |
| system_workload_primary_iops_total                |    global    |                                                     total                                                     |   iops/s   |
| system_workload_primary_iops                      |    global    |                                                  read, write                                                  |   iops/s   |
| system_workload_primary_io_size_total             |    global    |                                                    io_size                                                    |    KiB     |
| system_rebalance                                  |    global    |                                                  read, write                                                  |   KiB/s    |
| system_rebalance_left                             |    global    |                                                     left                                                      |    KiB     |
| system_rebalance_time_until_finish                |    global    |                                                     time                                                      |  seconds   |
| system_rebuild                                    |    global    |                                                  read, write                                                  |   KiB/s    |
| system_rebuild_left                               |    global    |                                                     left                                                      |    KiB     |
| system_defined_components                         |    global    | devices, fault_sets, protection_domains, rfcache_devices, sdc, sds, snapshots, storage_pools, volumes, vtrees | components |
| system_components_volumes_by_type                 |    global    |                                                  thick, thin                                                  |  volumes   |
| system_components_volumes_by_mapping              |    global    |                                               mapped, unmapped                                                |  volumes   |
| storage_pool_capacity_total                       | storage pool |                                                     total                                                     |    KiB     |
| storage_pool_capacity_in_use                      | storage pool |                                                    in_use                                                     |    KiB     |
| storage_pool_capacity_usage                       | storage pool |                                thick, decreased, thin, snapshot, spare, unused                                |    KiB     |
| storage_pool_capacity_utilization                 | storage pool |                                                     used                                                      | percentage |
| storage_pool_capacity_available_volume_allocation | storage pool |                                                   available                                                   |    KiB     |
| storage_pool_capacity_health_state                | storage pool |                           protected, degraded, in_maintenance, failed, unavailable                            |    KiB     |
| storage_pool_components                           | storage pool |                                      devices, snapshots, volumes, vtrees                                      | components |
| sdc_mdm_connection_state                          |     sdc      |                                                   connected                                                   |  boolean   |
| sdc_bandwidth                                     |     sdc      |                                                  read, write                                                  |   KiB/s    |
| sdc_iops                                          |     sdc      |                                                  read, write                                                  |   iops/s   |
| sdc_io_size                                       |     sdc      |                                                  read, write                                                  |    KiB     |
| sdc_num_of_mapped_volumed                         |     sdc      |                                                    mapped                                                     |  volumes   |

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

