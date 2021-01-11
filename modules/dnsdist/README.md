<!--
title: "DNSdist monitoring with Netdata"
description: "Monitor the health and performance of DNSdist load balancers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/dnsdist/README.md
sidebar_label: "DNSdist"
-->

# DNSdist monitoring with Netdata

[`DNSdist`](https://dnsdist.org/) is a highly DNS-, DoS- and abuse-aware loadbalancer.

This module monitors load-balancer performance and health metrics.

It collects metrics from [the internal webserver](https://dnsdist.org/guides/webserver.html).

Used endpoints:

- [/jsonstat?command=stats](https://dnsdist.org/guides/webserver.html#get--jsonstat).

## Requirements

For collecting metrics via HTTP, we need [enabled webserver](https://dnsdist.org/guides/webserver.html).

## Charts

- Client queries received in `queries/s`
- Client queries dropped in `queries/s`
- Packets dropped in `packets/s`
- Answers statistics in `answers/s`
- Backend responses in `responses/s`
- Backend communication errors in `errors/s`
- Backend error responses in `responses/s`
- Cache performance in `answers/s`
- DNSdist server CPU utilization in `ms/s`
- DNSdist server memory utilization in `MiB`
- Query latency in `queries/s`
- Average latency for the last N queries in `microseconds`

## Configuration

Edit the `go.d/dnsdist.conf` configuration file using `edit-config` from the
Agent's [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/dnsdist.conf
```

Needs `url` and `password` or _apikey_ to access the webserver.

Here is a configuration example:

```yaml
jobs:
  - name: local
    url: 'http://127.0.0.1:8083'
    headers:
      X-API-Key: 'your-api-key' # static pre-shared authentication key for access to the REST API (api-key).

  - name: remote
    url: 'http://203.0.113.0:8083'
    headers:
      X-API-Key: 'your-api-key' # static pre-shared authentication key for access to the REST API (api-key).
```

For all available options, see the `dnsdist`
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dnsdist.conf).

## Troubleshooting

To troubleshoot issues with the `dnsdist` collector, run the `go.d.plugin` with the debug option enabled. The output
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
./go.d.plugin -d -m dnsdist
```
