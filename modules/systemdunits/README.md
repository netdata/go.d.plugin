<!--
title: "Systemd units state monitoring with Netdata"
description: "Monitor the health and performance of Systemd units states with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/systemdunits/README.md
sidebar_label: "Systemd units"
-->

# Systemd units state monitoring with Netdata

[`Systemd`](https://www.freedesktop.org/wiki/Software/systemd/) is a suite of basic building blocks for a Linux system.

This module monitors `Systemd` units state.

- Works only on linux systems.
- Disabled by default. Should be explicitly enabled in the `go.d.conf`:

```yaml
# go.d.conf
modules:
  systemdunits: yes
```

## Charts

It produces the following charts:

- Service Unit State in `state`
- Socket Unit State in `state`
- Target Unit State in `state`
- Path Unit State in `state`
- Device Unit State in `state`
- Mount Unit State in `state`
- Automount Unit State in `state`
- Swap Unit State in `state`
- Timer Unit State in `state`
- Scope Unit State in `state`
- Slice Unit State in `state`

## Unit states

| Code  | Name         | Meaning |
| ----- | ------------ | ------- |
| 1     | `active`       | started, bound, plugged in, ..., depending on the unit type |
| 2     | `inactive`     | stopped, unbound, unplugged, ..., depending on the unit type |
| 3     | `activating`   | in the process of being activated |
| 4     | `deactivating` | in the process of being deactivated |
| 5     | `failed`       | the service failed in some way (process returned error code on exit, or crashed, an operation timed out, or after too many restarts) |

## Configuration

Edit the `go.d/systemdunits.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/systemdunits.conf
```

Needs only `include` option. Syntax is [shell file name pattern](https://golang.org/pkg/path/filepath/#Match).

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

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m systemdunits
```
