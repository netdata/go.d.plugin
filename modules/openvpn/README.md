<!--
title: "OpenVPN monitoring with Netdata"
description: "Monitor the health and performance of OpenVPN servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/openvpn/README.md
sidebar_label: "OpenVPN"
-->

# OpenVPN monitoring with Netdata

[`OpenVPN`](https://openvpn.net/) is an open-source commercial software that implements virtual private network
techniques to create secure point-to-point or site-to-site connections in routed or bridged configurations and remote
access facilities.

This module will monitor one or more `OpenVPN` instances via Management Interface.

## Requirements

- `OpenVPN` with enabled [`management-interface`](https://openvpn.net/community-resources/management-interface/).

## Charts

It produces the following charts:

- Total Number Of Active Clients in `clients`
- Total Traffic in `kilobits/s`

Per user charts (disabled by default, see `per_user_stats` in the module config file):

- User Traffic in `kilobits/s`
- User Connection Time in `seconds`

## Configuration

This collector is disabled by default. Should be explicitly enabled
in [go.d.conf](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).

Reason:
> Currently, the OpenVPN daemon can at most support a single management client any one time.

We disabled it to not break other tools which uses `Management Interface`.

Edit the `go.d/openvpn.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/openvpn.conf
```

Needs only `address` of OpenVPN `Management Interface`. Here is an example for 2 `OpenVPN` instances:

```yaml
jobs:
  - name: local
    address: /dev/openvpn

  - name: remote
    address: 203.0.113.10:7505
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/openvpn.conf).

## Troubleshooting

To troubleshoot issues with the `openvpn` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m openvpn
```
