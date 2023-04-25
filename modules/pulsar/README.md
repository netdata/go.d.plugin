<!--
title: "Apache Pulsar monitoring with Netdata"
description: "Monitor the health and performance of Apache Pulsar messaging systems with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/pulsar/README.md"
sidebar_label: "Pulsar"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Message brokers"
-->

# Apache Pulsar collector

[`Apache Pulsar`](http://pulsar.apache.org/) is an open-source distributed pub-sub messaging system.

This module will monitor one or more `Apache Pulsar` instances, depending on your configuration.

It collects broker statistics from
the [prometheus endpoint](https://pulsar.apache.org/docs/en/deploy-monitoring/#broker-stats).

`pulsar` module is tested on the following versions:

- v2.5.0

## Metrics

All metrics have "pulsar." prefix.

- topic_* metrics are available when `exposeTopicLevelMetricsInPrometheus` is set to true.
- subscription_* and namespace_subscription metrics are available when `exposeTopicLevelMetricsInPrometheus` si set to
  true.
- replication_* and namespace_replication_* metrics are available when replication is configured
  and `replicationMetricsEnabled` is set to true

| Metric                                             |   Scope   |                                Dimensions                                 |         Units         |
|----------------------------------------------------|:---------:|:-------------------------------------------------------------------------:|:---------------------:|
| broker_components                                  |  global   |          namespaces, topics, subscriptions, producers, consumers          |      components       |
| messages_rate                                      |  global   |                             publish, dispatch                             |      messages/s       |
| throughput_rate                                    |  global   |                             publish, dispatch                             |         KiB/s         |
| storage_size                                       |  global   |                                   used                                    |          KiB          |
| storage_operations_rate                            |  global   |                                read, write                                |   message batches/s   |
| msg_backlog                                        |  global   |                                  backlog                                  |       messages        |
| storage_write_latency                              |  global   | <=0.5ms, <=1ms, <=5ms, =10ms, <=20ms, <=50ms, <=100ms, <=200ms, <=1s, >1s |       entries/s       |
| entry_size                                         |  global   |     <=128B, <=512B, <=1KB, <=2KB, <=4KB, <=16KB, <=100KB, <=1MB, >1MB     |       entries/s       |
| subscription_delayed                               |  global   |                                  delayed                                  |    message bacthes    |
| subscription_msg_rate_redeliver                    |  global   |                                redelivered                                |      messages/s       |
| subscription_blocked_on_unacked_messages           |  global   |                                  blocked                                  |     subscriptions     |
| replication_rate                                   |  global   |                                  in, out                                  |      messages/s       |
| replication_throughput_rate                        |  global   |                                  in, out                                  |         KiB/s         |
| replication_backlog                                |  global   |                                  backlog                                  |       messages        |
| namespace_broker_components                        | namespace |                topics, subscriptions, producers, consumers                |      components       |
| namespace_messages_rate                            | namespace |                             publish, dispatch                             |      messages/s       |
| namespace_throughput_rate                          | namespace |                             publish, dispatch                             |         KiB/s         |
| namespace_storage_size                             | namespace |                                   used                                    |          KiB          |
| namespace_storage_operations_rate                  | namespace |                                read, write                                |   message batches/s   |
| namespace_msg_backlog                              | namespace |                                  backlog                                  |       messages        |
| namespace_storage_write_latency                    | namespace | <=0.5ms, <=1ms, <=5ms, =10ms, <=20ms, <=50ms, <=100ms, <=200ms, <=1s, >1s |       entries/s       |
| namespace_entry_size                               | namespace |     <=128B, <=512B, <=1KB, <=2KB, <=4KB, <=16KB, <=100KB, <=1MB, >1MB     |       entries/s       |
| namespace_subscription_delayed                     | namespace |                                  delayed                                  |    message bacthes    |
| namespace_subscription_msg_rate_redeliver          | namespace |                                redelivered                                |      messages/s       |
| namespace_subscription_blocked_on_unacked_messages | namespace |                                  blocked                                  |     subscriptions     |
| namespace_replication_rate                         | namespace |                                  in, out                                  |      messages/s       |
| namespace_replication_throughput_rate              | namespace |                                  in, out                                  |         KiB/s         |
| namespace_replication_backlog                      | namespace |                                  backlog                                  |       messages        |
| topic_producers                                    | namespace |                       <i>a dimension per topic</i>                        |       producers       |
| topic_subscriptions                                | namespace |                       <i>a dimension per topic</i>                        |     subscriptions     |
| topic_consumers                                    | namespace |                       <i>a dimension per topic</i>                        |       consumers       |
| topic_messages_rate_in                             | namespace |                       <i>a dimension per topic</i>                        |      publishes/s      |
| topic_messages_rate_out                            | namespace |                       <i>a dimension per topic</i>                        |     dispatches/s      |
| topic_throughput_rate_in                           | namespace |                       <i>a dimension per topic</i>                        |         KiB/s         |
| topic_throughput_rate_out                          | namespace |                       <i>a dimension per topic</i>                        |         KiB/s         |
| topic_storage_size                                 | namespace |                       <i>a dimension per topic</i>                        |          KiB          |
| topic_storage_read_rate                            | namespace |                       <i>a dimension per topic</i>                        |   message batches/s   |
| topic_storage_write_rate                           | namespace |                       <i>a dimension per topic</i>                        |   message batches/s   |
| topic_msg_backlog                                  | namespace |                       <i>a dimension per topic</i>                        |       messages        |
| topic_subscription_delayed                         | namespace |                       <i>a dimension per topic</i>                        |    message batches    |
| topic_subscription_msg_rate_redeliver              | namespace |                       <i>a dimension per topic</i>                        |      messages/s       |
| topic_subscription_blocked_on_unacked_messages     | namespace |                       <i>a dimension per topic</i>                        | blocked subscriptions |
| topic_replication_rate_in                          | namespace |                       <i>a dimension per topic</i>                        |      messages/s       |
| topic_replication_rate_out                         | namespace |                       <i>a dimension per topic</i>                        |      messages/s       |
| topic_replication_throughput_rate_in               | namespace |                       <i>a dimension per topic</i>                        |      messages/s       |
| topic_replication_throughput_rate_out              | namespace |                       <i>a dimension per topic</i>                        |      messages/s       |
| topic_replication_backlog                          | namespace |                       <i>a dimension per topic</i>                        |       messages        |

## Configuration

Edit the `go.d/pulsar.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/pulsar.conf
```

Needs only `url` to server's `/metrics` endpoint. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8080/metrics

  - name: remote
    url: http://203.0.113.10:8080/metrics
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/pulsar.conf).

## Topic filtering

By default, module collects data for all topics, but it supports topic filtering. Filtering doesn't exclude a topic
stats from the summary/namespace stats, it only removes the topic from the topic charts.

To check matcher syntax
see [matcher documentation](https://github.com/netdata/go.d.plugin/blob/master/pkg/matcher/README.md).

```yaml
  - name: local
    url: http://127.0.0.1:8080/metrics
    topic_filter:
      includes:
        - matcher1
        - matcher2
      excludes:
        - matcher1
        - matcher2
```

## Update every

Module default `update_every` is 60.

`Apache Pulsar` doesnt expose raw counters, it exposes rate. It counts rates every `statsUpdateFrequencyInSecs`. Default
value is 60 seconds.

Module `update_every` should be equal to `statsUpdateFrequencyInSecs`.

## Troubleshooting

To troubleshoot issues with the `pulsar` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m pulsar
  ```
