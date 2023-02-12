<!--
title: "PowerDNS Recursor monitoring with Netdata"
description: "Monitor the health and performance of PowerDNS recursor instances with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/powerdns_recursor/README.md"
sidebar_label: "PowerDNS Recursor"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Devices"
-->

# PowerDNS Recursor monitoring with Netdata

[`PowerDNS Recursor`](https://doc.powerdns.com/recursor/) is a high-performance DNS recursor with built-in scripting
capabilities.

This module monitors one or more `PowerDNS Recursor` instances, depending on your configuration.

It collects metrics
from [the internal webserver](https://doc.powerdns.com/recursor/http-api/index.html#built-in-webserver-and-http-api).

Used endpoints:

- [`/api/v1/servers/localhost/statistics`](https://doc.powerdns.com/recursor/common/api/endpoint-statistics.html)

## Requirements

For collecting metrics via HTTP, we need:

- [enabled webserver](https://doc.powerdns.com/recursor/http-api/index.html#webserver).
- [enabled HTTP API](https://doc.powerdns.com/recursor/http-api/index.html#enabling-the-api).

## Metrics

All metrics have "powerdns_recursor." prefix.

| Metric        | Scope  |                                        Dimensions                                         |    Units    |
|---------------|:------:|:-----------------------------------------------------------------------------------------:|:-----------:|
| questions_in  | global |                                     total, tcp, ipv6                                      | questions/s |
| questions_out | global |                                 udp, tcp, ipv6, throttled                                 | questions/s |
| answer_time   | global |                         0-1ms, 1-10ms, 10-100ms, 100-1000ms, slow                         |  queries/s  |
| timeouts      | global |                                     total, ipv4, ipv6                                     | timeouts/s  |
| drops         | global | over-capacity-drops, query-pipe-full-drops, too-old-drops, truncated-drops, empty-queries |   drops/s   |
| cache_size    | global |                            cache, packet-cache, negative-cache                            |   entries   |

## Configuration

Edit the `go.d/powerdns_recursor.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/powerdns_recursor.conf
```

To add a new endpoint to collect metrics from, or change the URL that Netdata looks for, add or configure the `name` and
`url` values. Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

Here is an example with two endpoints:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8081

  - name: remote
    url: http://203.0.113.0:8081
```

For all available options, see the PowerDNS Recursor
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/powerdns_recursor.conf).

## Troubleshooting

To troubleshoot issues with the `powerdns_recursor` collector, run the `go.d.plugin` with the debug option enabled. The
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
  ./go.d.plugin -d -m powerdns_recursor
  ```
