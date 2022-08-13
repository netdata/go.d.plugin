<!--
title: "Traefik monitoring with Netdata"
description: "Monitor the health and performance of Traefik with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/traefik/README.md
sidebar_label: "Traefik"
-->

# Traefik monitoring with Netdata

[`Traefik`](https://traefik.io/traefik/) is a leading modern reverse proxy and load balancer that makes deploying
microservices easy. .

This module will monitor one or more `Traefik` instances, depending on your configuration.

## Requirements

- `Traefik` with enabled [Prometheus exporter](https://doc.traefik.io/traefik/observability/metrics/prometheus/).

## Charts

Current implementation collects only [entrypoint](https://doc.traefik.io/traefik/routing/entrypoints/) metrics.

### Entrypoint

- Processed HTTP requests, partitioned by code class in `requests/s`
- Average HTTP request processing time, partitioned by code class in `milliseconds`
- Open connections, partitioned by method in `connections`

## Configuration

Edit the `go.d/traefik.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/traefik.conf
```

Needs only `url` to server's `/metrics` endpoint. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8082/metrics

  - name: remote
    url: http://203.0.113.10:8082/metrics
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/traefik.conf).

## Troubleshooting

To troubleshoot issues with the `traefik` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m traefik
  ```
