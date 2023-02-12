<!--
title: "Systemd units state monitoring with Netdata"
description: "Monitor the health and performance of Systemd units states with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/systemdunits/README.md"
sidebar_label: "Systemd units"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitoring/System metrics"
-->

# Systemd units state monitoring with Netdata

[Systemd](https://www.freedesktop.org/wiki/Software/systemd/) is a suite of basic building blocks for a Linux system.

This module monitors Systemd units state.

## Requirements

- Works only on Linux systems.
- Disabled by default. Should be explicitly enabled in the `go.d.conf`:

```yaml
# go.d.conf
modules:
  systemdunits: yes
```

## Metrics

The unit types and states description can be found in
the [official documentation](https://www.freedesktop.org/software/systemd/man/systemd.html#Concepts).

All metrics have "systemd." prefix.

Labels per scope:

- unit: unit_name.

| Metric               | Scope |                     Dimensions                     | Units |
|----------------------|:-----:|:--------------------------------------------------:|:-----:|
| service_unit_state   | unit  | active, inactive, activating, deactivating, failed | state |
| socket_unit_state    | unit  | active, inactive, activating, deactivating, failed | state |
| target_unit_state    | unit  | active, inactive, activating, deactivating, failed | state |
| path_unit_state      | unit  | active, inactive, activating, deactivating, failed | state |
| device_unit_state    | unit  | active, inactive, activating, deactivating, failed | state |
| mount_unit_state     | unit  | active, inactive, activating, deactivating, failed | state |
| automount_unit_state | unit  | active, inactive, activating, deactivating, failed | state |
| swap_unit_state      | unit  | active, inactive, activating, deactivating, failed | state |
| timer_unit_state     | unit  | active, inactive, activating, deactivating, failed | state |
| scope_unit_state     | unit  | active, inactive, activating, deactivating, failed | state |
| slice_unit_state     | unit  | active, inactive, activating, deactivating, failed | state |

## Configuration

Edit the `go.d/systemdunits.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/systemdunits.conf
```

Needs only `include` option. Syntax is the [shell file name pattern](https://golang.org/pkg/path/filepath/#Match).

Here are some examples:

```yaml
jobs:
  - name: my-specific-service-unit
    include:
      - 'my-specific.service'

  - name: service-units
    include:
      - '*.service'

  - name: socket-units
    include:
      - '*.socket'
```

For all available options, see the Systemdunits
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/systemdunits.conf).

## Troubleshooting

To troubleshoot issues with the `systemdunits` collector, run the `go.d.plugin` with the debug option enabled. The
output should give you clues as to why the collector isn't working.

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
  ./go.d.plugin -d -m systemdunits
  ```
