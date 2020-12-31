<!--
title: "Unbound monitoring with Netdata"
description: "Monitor the health and performance of Unbound DNS resolvers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/unbound/README.md
sidebar_label: "Unbound"
-->

# Unbound monitoring with Netdata

[`Unbound`](https://nlnetlabs.nl/projects/unbound/about/) is a validating, recursive, and caching DNS resolver product
from NLnet Labs.

This module monitors one or more `Unbound` servers, depending on your configuration.

## Requirements

- `Unbound` with enabled `remote-control` interface (
  see [unbound.conf](https://nlnetlabs.nl/documentation/unbound/unbound.conf))

If using unix socket:

- socket should be readable and writeable by `netdata` user

If using ip socket and TLS is disabled:

- socket should be accessible via network

If TLS is enabled, in addition:

- `control-key-file` should be readable by `netdata` user
- `control-cert-file` should be readable by `netdata` user

For auto detection parameters from `unbound.conf`:

- `unbound.conf` should be readable by `netdata` user
- if you have several configuration files (include feature) all of them should be readable by `netdata` user

## Charts

Module produces following summary charts:

- Received Queries in `queries`
- Rate Limited Queries in `queries`
- DNSCrypt Queries in `queries`
- Cache Statistics in `events`
- Cache Statistics Percentage in `percentage`
- Cache Prefetches in `prefetches`
- Replies Served From Expired Cache in `replies`
- Replies That Needed Recursive Processing in `replies`
- Time Spent On Recursive Processing in `milliseconds`
- Request List Usage in `queries`
- Current Request List Usage in `queries`
- Request List Jostle List Events in `queries`
- TCP Handler Buffers in `buffers`
- Uptime `seconds`

If `extended-statistics` is enabled:

- Queries By Type in `queries`
- Queries By Class in `queries`
- Queries By OpCode in `queries`
- Queries By Flag in `queries`
- Replies By RCode in `replies`
- Cache Items Count in `items`
- Cache Memory in `KB`
- Module Memory in `KB`
- TCP and TLS Stream Wait Buffer Memory in `KB`

Per thread charts (only if number of threads > 1):

- Received Queries in `queries`
- Rate Limited Queries in `queries`
- DNSCrypt Queries in `queries`
- Cache Statistics in `events`
- Cache Statistics Percentage in `events`
- Cache Prefetches in `prefetches`
- Replies Served From Expired Cache in `replies`
- Replies That Needed Recursive Processing in `replies`
- Time Spent On Recursive Processing in `milliseconds`
- Request List Usage in `queries`
- Current Request List Usage in `queries`
- Request List Jostle List Events in `queries`
- TCP Handler Buffers in `buffers`

## Configuration

Edit the `go.d/unbound.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/unbound.conf
```

This Unbound collector only needs the `address` to a server's `remote-control` interface if TLS is disabled or `address`
of unix socket. Otherwise you need to set path to the `control-key-file` and `control-cert-file` files.

The module tries to auto-detect following parameters reading `unbound.conf`:

- address
- cumulative_stats
- use_tls
- tls_cert
- tls_key

Module supports both cumulative and non-cumulative modes. Default is non-cumulative. If your server has enabled
`statistics-cumulative`, but the module fails to auto-detect it (`unbound.conf` is not readable or it is a remote
server), you need to set it manually in the configuration file.

Here is an example for several servers:

```yaml
jobs:
  - name: local
    address: 127.0.0.1:8953
    use_tls: yes
    tls_skip_verify: yes
    tls_cert: /etc/unbound/unbound_control.pem
    tls_key: /etc/unbound/unbound_control.key

  - name: remote
    address: 203.0.113.10:8953
    use_tls: no

  - name: remote_cumulative
    address: 203.0.113.11:8953
    use_tls: no
    cumulative_stats: yes

  - name: socket
    address: /var/run/unbound.sock
```

For all available options, please see the
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/unbound.conf).

## Troubleshooting

Ensure that the control protocol is actually configured correctly. Run following command as `root` user:

```bash
unbound-control stats_noreset
```

It should print out a bunch of info about the internal statistics of the server. If this returns an error, you don't
have the control protocol set up correctly.

To troubleshoot issues with the `unbound` collector, run the `go.d.plugin` with the debug option enabled. The output
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
./go.d.plugin -d -m unbound
```
