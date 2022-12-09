<!--
title: "Ping monitoring with Netdata"
description: "Monitor round-trip time and packet loss to network hosts with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/ping/README.md
sidebar_label: "Ping"
-->

# Ping monitoring with Netdata

This module measures round-tripe time and packet loss by sending ping messages to network hosts.

## Requirements

There are two operational modes:

- privileged (send raw ICMP ping, default). Requires
  CAP_NET_RAW [capability](https://man7.org/linux/man-pages/man7/capabilities.7.html) or root privileges:
  > **Note**: set automatically during Netdata installation.

  ```bash
  sudo setcap CAP_NET_RAW=eip <INSTALL_PREFIX>/usr/libexec/netdata/plugins.d/go.d.plugin
  ```

- unprivileged (send UDP ping, Linux only).
  Requires configuring [ping_group_range](https://www.man7.org/linux/man-pages/man7/icmp.7.html):

  ```bash
  sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
  ```
  To persist the change add `net.ipv4.ping_group_range="0 2147483647"` to `/etc/sysctl.conf` and
  execute `sudo sysctl -p`.

The mode can be changed in the module [configuration file](#Configuration).

## Metrics

All metrics have "ping." prefix.

Labels per scope:

- host: host.

| Metric           | Scope |   Dimensions   |    Units     |
|------------------|:-----:|:--------------:|:------------:|
| host_rtt         | host  | min, max, avg  | milliseconds |
| host_std_dev_rtt | host  |    std_dev     | milliseconds |
| host_packet_loss | host  |      loss      |  percentage  |
| host_packets     | host  | received, sent |   packets    |

## Configuration

Edit the `go.d/ping.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/ping.conf
```

Here is an example configuration:

```yaml
jobs:
  - name: example
    hosts:
      - 192.0.2.0
      - 192.0.2.1
      - example.com
    packets: 5       # number of ping packets to send.
    interval: 200ms  # time to wait between sending ping packets.
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/ping.conf).

## Troubleshooting

To troubleshoot issues with the `ping` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins' directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m ping
```
