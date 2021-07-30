<!--
title: "HAProxy monitoring with Netdata"
description: "Monitor the health and performance of HAProxy with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/haproxy/README.md
sidebar_label: "HAProxy"
-->

# HAProxy monitoring with Netdata

[`HAProxy`](http://www.haproxy.org/) is a free, very fast and reliable solution offering high availability, load
balancing, and proxying for TCP and HTTP-based applications.

This module will monitor one or more `HAProxy` instances, depending on your configuration.

## Requirements

- `HAProxy` v2.0+ (or 1.9r1+ for Enterprise users) with enabled PROMEX addon. PROMEX is not built by default with
  `HAProxy`. It is provided as an extra component for
  everyone [who wants to use it](https://github.com/haproxy/haproxy/tree/master/addons/promex).

## Charts

Current implementation collects
only [backend](https://www.haproxy.com/documentation/hapee/latest/configuration/config-sections/backend/) metrics.

### Backend

- Sessions
    - Current number of active sessions in `sessions`
    - Sessions rate in `sessions/s`
- Responses
    - Average response time for last 1024 successful connections in `milliseconds`
    - HTTP responses by code class in `responses/s`
- Queue
    - Average queue time for last 1024 successful connections in `milliseconds`
    - Current number of queued requests in `requests`
- Network
    - Network traffic in `bytes/s`

## Configuration

Edit the `go.d/haproxy.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/haproxy.conf
```

Needs only `url` to server's `/metrics` endpoint. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8404/metrics

  - name: remote
    url: http://203.0.113.10:8404/metrics
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/haproxy.conf).

## Troubleshooting

To troubleshoot issues with the `haproxy` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory,
  switch to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

- You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m haproxy
```
