<!--
title: "RabbitMQ monitoring with Netdata"
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

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m rabbitmq
