<!--
title: "vCenter Server Appliance monitoring with Netdata"
description: "Monitor the health and performance of vCenter appliances with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/blob/master/modules/vcsa/README.md"
sidebar_label: "vCenter Server Appliance"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Virtualized environments/Virtualize hosts"
-->

# vCenter Server Appliance monitoring with Netdata

The [`vCenter Server Appliance`](https://docs.vmware.com/en/VMware-vSphere/6.5/com.vmware.vsphere.vcsa.doc/GUID-223C2821-BD98-4C7A-936B-7DBE96291BA4.html)
using [`Health API`](https://code.vmware.com/apis/60/vcenter-server-appliance-management) is a preconfigured Linux
virtual machine, which is optimized for running VMware vCenter Server® and the associated services on Linux.

This module collects health statistics from one or more `vCenter Server Appliance` servers, depending on your
configuration.

## Requirements

- `vSphere` 6.5+

## Metrics

All metrics have "vcsa." prefix.

| Metric                  | Scope  |                   Dimensions                   | Units  |
|-------------------------|:------:|:----------------------------------------------:|:------:|
| system_health           | global |                     system                     | status |
| components_health       | global | applmgmt, database_storage, mem, storage, swap | status |
| software_updates_health | global |               software_packages                | status |

## Health statuses

Overall System Health:

| Numeric |   Text    | Description                                                                                                              |
|:-------:|:---------:|:-------------------------------------------------------------------------------------------------------------------------|
|  `-1`   | `unknown` | Module failed to decode status.                                                                                          |
|   `0`   |  `green`  | All components in the appliance are healthy.                                                                             |
|   `1`   | `yellow`  | One or more components in the appliance might become overloaded soon.                                                    |
|   `2`   | `orange`  | One or more components in the appliance might be degraded.                                                               |
|   `3`   |   `red`   | One or more components in the appliance might be in an unusable status and the appliance might become unresponsive soon. |
|   `4`   |  `gray`   | No health data is available.                                                                                             |

Components Health:

| Numeric |   Text    | Description                                                  |
|:-------:|:---------:|:-------------------------------------------------------------|
|  `-1`   | `unknown` | Module failed to decode status.                              |
|   `0`   |  `green`  | The component is healthy.                                    |
|   `1`   | `yellow`  | The component is healthy, but may have some problems.        |
|   `2`   | `orange`  | The component is degraded, and may have serious problems.    |
|   `3`   |   `red`   | The component is unavailable, or will stop functioning soon. |
|   `4`   |  `gray`   | No health data is available.                                 |

Software Updates Health:

| Numeric |   Text    | Description                                          |
|:-------:|:---------:|:-----------------------------------------------------|
|  `-1`   | `unknown` | Module failed to decode status.                      |
|   `0`   |  `green`  | No updates available.                                |
|   `2`   | `orange`  | Non-security patches might be available.             |
|   `3`   |   `red`   | Security patches might be available.                 |
|   `4`   |  `gray`   | An error retrieving information on software updates. |

## Configuration

Edit the `go.d/vcsa.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/vcsa.conf
```

Needs only `url`, `username` and `password`. Here is an example for 2 servers:

```yaml
jobs:
  - name: vcsa1
    url: https://203.0.113.0
    username: admin@vsphere.local
    password: somepassword

  - name: vcsa2
    url: https://203.0.113.10
    username: admin@vsphere.local
    password: somepassword
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/vcsa.conf).

## Troubleshooting

To troubleshoot issues with the `vcsa` collector, run the `go.d.plugin` with the debug option enabled. The output should
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
  ./go.d.plugin -d -m vcsa
  ```
