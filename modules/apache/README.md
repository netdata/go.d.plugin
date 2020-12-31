<!--
title: "Apache monitoring with Netdata"
description: "Monitor the health and performance of Apache web servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/apache/README.md
sidebar_label: "Apache"
-->

# Apache monitoring with Netdata

[`Apache`](https://httpd.apache.org/) is an open-source HTTP server for modern operating systems including UNIX and
Windows.

This module will monitor one or more `Apache` servers, depending on your configuration.

## Requirements

- `Apache` with enabled [`mod_status`](https://httpd.apache.org/docs/2.4/mod/mod_status.html)

## Charts

It produces the following charts:

- Requests in `requests/s`
- Connections in `connections`
- Async Connections in `connections`
- Scoreboard in `connections`
- Bandwidth in `kilobits/s`
- Workers in `workers`
- Lifetime Average Number Of Requests Per Second in `requests/s`
- Lifetime Average Number Of Bytes Served Per Second in `KiB/s`
- Lifetime Average Response Size in `KiB`

## Configuration

Edit the `go.d/apache.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/apache.conf
```

Needs only `url` to server's `server-status?auto`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/server-status?auto

  - name: remote
    url: http://203.0.113.10/server-status?auto
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/apache.conf).

## Troubleshooting

To troubleshoot issues with the `apache` collector, run the `go.d.plugin` with the debug option enabled. The output
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
./go.d.plugin -d -m apache
```
