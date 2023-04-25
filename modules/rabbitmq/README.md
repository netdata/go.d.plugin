<!--
title: "RabbitMQ monitoring with Netdata"
description: "Monitor the health and performance of RabbitMQ message brokers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/rabbitmq/README.md"
sidebar_label: "rabbitmq-go.d.plugin (Recommended)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Message brokers"
-->

# RabbitMQ collector

[RabbitMQ](https://www.rabbitmq.com/) is an open-source message broker.

This module monitors one or more RabbitMQ instances, depending on your configuration.

It collects data using an HTTP-based API provided by the [management plugin](https://www.rabbitmq.com/management.html).
The following endpoints are used:

- `/api/overview`
- `/api/node/{node_name}`
- `/api/vhosts`
- `/api/queues` (disabled by default)

## Requirements

RabbitMQ with [enabled](https://www.rabbitmq.com/management.html#getting-started) management plugin.

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/rabbitmq/metrics.csv) for a list
of metrics.

## Configuration

Edit the `go.d/rabbitmq.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/rabbitmq.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://localhost:15672
    collect_queues_metrics: no

  - name: remote
    url: http://203.0.113.10:15672
    collect_queues_metrics: no
```

This collector can also collect per-vhost per-queue metrics, which is disabled by
default (`collect_queues_metrics`). Enabling this can introduce serious overhead on both netdata and rabbitmq if many
queues are configured and used.

For all available options, see the
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/rabbitmq.conf).

## Troubleshooting

To troubleshoot issues with the `rabbitmq` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m rabbitmq
  ```
