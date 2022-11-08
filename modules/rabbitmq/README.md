<!--
title: "RabbitMQ monitoring with Netdata"
description: "Monitor the health and performance of RabbitMQ message brokers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/rabbitmq/README.md"
sidebar_label: "rabbitmq-go.d.plugin (Recommended)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Message brokers"
-->

# RabbitMQ monitoring with Netdata

[`RabbitMQ`](https://www.rabbitmq.com/) is the open source message broker.

This module monitors `RabbitMQ` performance and health metrics.

It collects data using following endpoints:

- `/api/overview`
- `/api/node/{node_name}`
- `/api/vhosts`

## Metrics

All metrics have "rabbitmq." prefix.

| Metric           | Scope  |                                                             Dimensions                                                              |    Units    |
|------------------|:------:|:-----------------------------------------------------------------------------------------------------------------------------------:|:-----------:|
| queued_messages  | global |                                                        ready, unacknowledged                                                        |  messages   |
| message_rates    | global | ack, publish, publish_in, publish_out, confirm, deliver, deliver_no_ack, get, get_no_ack, deliver_get, redeliver, return_unroutable | messages/s  |
| global_counts    | global |                                         channels, consumers, connections, queues, exchanges                                         |   counts    |
| file_descriptors | global |                                                                used                                                                 | descriptors |
| sockets          | global |                                                                used                                                                 | descriptors |
| processes        | global |                                                                used                                                                 |  processes  |
| erlang_run_queue | global |                                                               length                                                                |  processes  |
| memory           | global |                                                                used                                                                 |     MiB     |
| disk_space       | global |                                                                free                                                                 |     MiB     |
| disk_space       | global |                                                                free                                                                 |     GiB     |
| vhost_messages   | vhost  |                            ack, confirm, deliver, get, get_no_ack, publish, redeliver, return_unroutable                            |  messages   |

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
