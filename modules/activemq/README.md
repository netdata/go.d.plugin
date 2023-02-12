<!--
title: "ActiveMQ monitoring with Netdata"
description: "Monitor the health and performance of ActiveMQ message brokers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/activemq/README.md"
sidebar_label: "ActiveMQ"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitoring/Message brokers"
-->

# ActiveMQ monitoring with Netdata

[`ActiveMQ`](https://activemq.apache.org/) is an open source message broker written in Java together with a full Java
Message Service client.

This plugin collects queues and topics metrics using ActiveMQ Console API.

## Metrics

All metrics have "activemq." prefix.

| Metric               | Scope  |     Dimensions     |   Units    |
|----------------------|:------:|:------------------:|:----------:|
| messages             | global | enqueued, dequeued | messages/s |
| unprocessed_messages | global |    unprocessed     |  messages  |
| consumers            | global |     consumers      | consumers  |

## Configuration

Edit the `go.d/activemq.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/activemq.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8161
    webadmin: admin
    max_queues: 100
    max_topics: 100
    queues_filter: '!sandr* *'
    topics_filter: '!sandr* *'

  - name: remote
    url: http://203.0.113.10:8161
    webadmin: admin
```

For all available options, please see the
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/activemq.conf).

## Troubleshooting

To troubleshoot issues with the `activemq` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m activemq
  ```
