<!--
title: "TCP endpoint monitoring with Netdata"
description: "Monitor the health and performance of any TCP endpoint with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/portcheck/README.md"
sidebar_label: "TCP endpoints"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Remotes"
-->

# TCP endpoint monitoring with Netdata

This module monitors one or more TCP services availability and response time.

## Metrics

All metrics have "portcheck." prefix.

Labels per scope:

- global: host, port.

| Metric         | Scope  |        Dimensions        |  Units  |
|----------------|:------:|:------------------------:|:-------:|
| status         | global | success, failed, timeout | boolean |
| state_duration | global |           time           | seconds |
| latency        | global |           time           |   ms    |

## Configuration

Edit the `go.d/portcheck.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/portcheck.conf
```

Here is an example for 3 servers:

> **Note**: a literal IPv6 address must be enclosed in square brackets, as in "[::1]".

```yaml
jobs:
  - name: server1
    host: 127.0.0.1
    ports:
      - 22
      - 23

  - name: server2
    host: "[2001:DB8::1]"
    ports:
      - 80
      - 8080

  - name: server3
    host: 203.0.113.10
    ports:
      - 80
      - 81
      - 8081
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/portcheck.conf).

## Troubleshooting

To troubleshoot issues with the `portcheck` collector, run the `go.d.plugin` with the debug option enabled. The output
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
./go.d.plugin -d -m portcheck
```
