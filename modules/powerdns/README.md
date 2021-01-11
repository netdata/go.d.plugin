<!--
title: "PowerDNS Authoritative Server monitoring with Netdata"
description: "Monitor the health and performance of PowerDNS nameservers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/powerdns/README.md
sidebar_label: "PowerDNS Authoritative Server"
-->

# PowerDNS Authoritative Server monitoring with Netdata

[`PowerDNS Authoritative Server`](https://doc.powerdns.com/authoritative/) is a versatile nameserver which supports a
large number of backends.

This module monitors one or more `PowerDNS Authoritative Server` instances, depending on your configuration.

It collects metrics from [the internal webserver](https://doc.powerdns.com/authoritative/http-api/index.html#webserver).

Used endpoints:

- [`/api/v1/servers/localhost/statistics`](https://doc.powerdns.com/authoritative/http-api/statistics.html)

## Requirements

For collecting metrics via HTTP, we need:

- [enabled webserver](https://doc.powerdns.com/authoritative/http-api/index.html#webserver).
- [enabled HTTP API](https://doc.powerdns.com/authoritative/http-api/index.html#enabling-the-api).

## Charts

- Incoming questions in `questions/s`
- Outgoing questions in `questions/s`
- Cache Usage in `events/s`
- Cache Size in `entries`
- Latency in `microseconds`

## Configuration

Edit the `go.d/powerdns.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/powerdns.conf
```

To add a new endpoint to collect metrics from, or change the URL that Netdata looks for, add or configure the `name` and
`url` values. Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

Here is an example with two endpoints:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8081
    headers:
      X-API-KEY: secret  # static pre-shared authentication key for access to the REST API (api-key).

  - name: remote
    url: http://203.0.113.0:8081
    headers:
      X-API-KEY: secret  # static pre-shared authentication key for access to the REST API (api-key).
```

For all available options, see the PowerDNS Authoritative Server
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/powerdns.conf).

## Troubleshooting

To troubleshoot issues with the `powerdns` collector, run the `go.d.plugin` with the debug option enabled. The output
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
./go.d.plugin -d -m powerdns
```
