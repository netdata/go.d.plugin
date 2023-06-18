# ScaleIO collector

## Overview

[Dell EMC ScaleIO](https://www.dellemc.com/en-us/storage/data-storage/software-defined-storage.htm) is a
software-defined storage product from Dell EMC that creates a server-based storage area network from local application
server storage using existing customer hardware or EMC servers.

This collector monitors one or more ScaleIO (VxFlex OS) instances via VxFlex OS Gateway API.

It collects metrics for the following ScaleIO components:

- System
- Storage Pool
- Sdc

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                              |                                                  Dimensions                                                   |    Unit    |
|-----------------------------------------------------|:-------------------------------------------------------------------------------------------------------------:|:----------:|
| scaleio.system_capacity_total                       |                                                     total                                                     |    KiB     |
| scaleio.system_capacity_in_use                      |                                                    in_use                                                     |    KiB     |
| scaleio.system_capacity_usage                       |                                thick, decreased, thin, snapshot, spare, unused                                |    KiB     |
| scaleio.system_capacity_available_volume_allocation |                                                   available                                                   |    KiB     |
| scaleio.system_capacity_health_state                |                           protected, degraded, in_maintenance, failed, unavailable                            |    KiB     |
| scaleio.system_workload_primary_bandwidth_total     |                                                     total                                                     |   KiB/s    |
| scaleio.system_workload_primary_bandwidth           |                                                  read, write                                                  |   KiB/s    |
| scaleio.system_workload_primary_iops_total          |                                                     total                                                     |   iops/s   |
| scaleio.system_workload_primary_iops                |                                                  read, write                                                  |   iops/s   |
| scaleio.system_workload_primary_io_size_total       |                                                    io_size                                                    |    KiB     |
| scaleio.system_rebalance                            |                                                  read, write                                                  |   KiB/s    |
| scaleio.system_rebalance_left                       |                                                     left                                                      |    KiB     |
| scaleio.system_rebalance_time_until_finish          |                                                     time                                                      |  seconds   |
| scaleio.system_rebuild                              |                                                  read, write                                                  |   KiB/s    |
| scaleio.system_rebuild_left                         |                                                     left                                                      |    KiB     |
| scaleio.system_defined_components                   | devices, fault_sets, protection_domains, rfcache_devices, sdc, sds, snapshots, storage_pools, volumes, vtrees | components |
| scaleio.system_components_volumes_by_type           |                                                  thick, thin                                                  |  volumes   |
| scaleio.system_components_volumes_by_mapping        |                                               mapped, unmapped                                                |  volumes   |

### storage pool

These metrics refer to the storage pool.

This scope has no labels.

Metrics:

| Metric                                                    |                        Dimensions                        |    Unit    |
|-----------------------------------------------------------|:--------------------------------------------------------:|:----------:|
| scaleio.storage_pool_capacity_total                       |                          total                           |    KiB     |
| scaleio.storage_pool_capacity_in_use                      |                          in_use                          |    KiB     |
| scaleio.storage_pool_capacity_usage                       |     thick, decreased, thin, snapshot, spare, unused      |    KiB     |
| scaleio.storage_pool_capacity_utilization                 |                           used                           | percentage |
| scaleio.storage_pool_capacity_available_volume_allocation |                        available                         |    KiB     |
| scaleio.storage_pool_capacity_health_state                | protected, degraded, in_maintenance, failed, unavailable |    KiB     |
| scaleio.storage_pool_components                           |           devices, snapshots, volumes, vtrees            | components |

### sdc

These metrics refer to the SDC (ScaleIO Data Client).

This scope has no labels.

Metrics:

| Metric                            | Dimensions  |  Unit   |
|-----------------------------------|:-----------:|:-------:|
| scaleio.sdc_mdm_connection_state  |  connected  | boolean |
| scaleio.sdc_bandwidth             | read, write |  KiB/s  |
| scaleio.sdc_iops                  | read, write | iops/s  |
| scaleio.sdc_io_size               | read, write |   KiB   |
| scaleio.sdc_num_of_mapped_volumed |   mapped    | volumes |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/scaleio.conf`.

The file format is YAML. Generally, the format is:

```yaml
update_every: 1
autodetection_retry: 0
jobs:
  - name: some_name1
  - name: some_name1
```

You can edit the configuration file using the `edit-config` script from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md#the-netdata-config-directory).

```bash
cd /etc/netdata 2>/dev/null || cd /opt/netdata/etc/netdata
sudo ./edit-config go.d/scaleio.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |       Default        | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:--------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |          5           |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |          0           |          |
|         url          | Server URL.                                                                                               | https://127.0.0.1:80 |   yes    |
|       timeout        | HTTP request timeout.                                                                                     |          1           |          |
|       username       | Username for basic HTTP authentication.                                                                   |                      |   yes    |
|       password       | Password for basic HTTP authentication.                                                                   |                      |   yes    |
|      proxy_url       | Proxy URL.                                                                                                |                      |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                      |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                      |          |
|        method        | HTTP request method.                                                                                      |         GET          |          |
|         body         | HTTP request body.                                                                                        |                      |          |
|       headers        | HTTP request headers.                                                                                     |                      |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |          no          |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |          no          |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                      |          |
|       tls_cert       | Client TLS certificate.                                                                                   |                      |          |
|       tls_key        | Client TLS key.                                                                                           |                      |          |

</details>

#### Examples

##### Basic

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: https://127.0.0.1
    username: admin
    password: password
    tls_skip_verify: yes  # self-signed certificate
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Local and remote instance.

<details>
<summary>Config</summary>

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
    tls_skip_verify: yes
```

</details>

## Troubleshooting

### Debug mode

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
