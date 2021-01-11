<!--
title: "RabbitMQ monitoring with Netdata"
description: "Monitor the health and performance of RabbitMQ message brokers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/rabbitmq/README.md
sidebar_label: "RabbitMQ"
-->

# RabbitMQ monitoring with Netdata

[`RabbitMQ`](https://www.rabbitmq.com/) is the open source message broker.

This module monitors `RabbitMQ` performance and health metrics.

It collects data using following endpoints:

- `/api/overview`
- `/api/node/{node_name}`
- `/api/vhosts`

## Charts

It produces the following charts:

- Queued Messages in `messages`
- Messages in `messages/s`
- Global Counts in `counts`
- File Descriptors in `descriptors`
- Socket Descriptors in `descriptors`
- Erlang Processes in `processes`
- Erlang run queue in `processes`
- Memory in `MiB`
- Disk Space in `GiB`

Per vhost charts:

- Messages in `messages/s`

## Configuration

Edit the `go.d/rabbitmq.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/rabbitmq.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://localhost:15672

  - name: remote
    url: http://203.0.113.10:15672

```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/rabbitmq.conf).

## Troubleshooting

To troubleshoot issues with the `rabbitmq` collector, run the `go.d.plugin` with the debug option enabled. The output
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
./go.d.plugin -d -m rabbitmq
```
