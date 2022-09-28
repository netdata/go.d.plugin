<!--
title: "DNS query monitoring with Netdata"
description: "Monitor the health and performance of DNS query round-trip time with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/dnsquery/README.md
sidebar_label: "DNS queries"
-->

# DNS query monitoring with Netdata

This module provides DNS query RTT in milliseconds.

## Metrics

All metrics have "dnsquery." prefix.

| Metric     | Scope  |             Dimensions             | Units |
|------------|:------:|:----------------------------------:|:-----:|
| query_time | global | <i>a dimension per name server</i> |  ms   |

## Configuration

Edit the `go.d/dns_query.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/dns_query.conf
```

Here is an example:

```yaml
jobs:
  - name: job1
    domains:
      - google.com
      - github.com
      - reddit.com
    servers:
      - 8.8.8.8
      - 8.8.4.4
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dns_query.conf).

## Troubleshooting

To troubleshoot issues with the `dns_query` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

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
  ./go.d.plugin -d -m dns_query
  ```
