<!--
title: "Consul monitoring with Netdata"
description: "Monitor the health and performance of Consul service meshes with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/consul/README.md"
sidebar_label: "Consul"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Webapps"
-->

# Consul monitoring with Netdata

[`Consul`](https://www.consul.io/) is a service networking solution to connect and secure services across any runtime
platform and public or private cloud.

This module monitors `Consul` health checks.

## Metrics

All metrics have "consul." prefix.

| Metric                      | Scope |               Dimensions                | Units  |
|-----------------------------|:-----:|:---------------------------------------:|:------:|
| service_health_check_status | check | passing, maintenance, warning, critical | status |
| unbound_health_check_status | check | passing, maintenance, warning, critical | status |

## Configuration

Edit the `go.d/consul.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/consul.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8500

  - name: remote
    url: http://203.0.113.10:8500
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/consul.conf).

## Troubleshooting

To troubleshoot issues with the `consul` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m consul
  ```
