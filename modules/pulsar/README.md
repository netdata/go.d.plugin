<!--
title: "Apache Pulsar monitoring with Netdata"
description: "Monitor the health and performance of Apache Pulsar messaging systems with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/pulsar/README.md
sidebar_label: "Pulsar"
-->

# Apache Pulsar monitoring with Netdata

[`Apache Pulsar`](http://pulsar.apache.org/) is an open-source distributed pub-sub messaging system.

This module will monitor one or more `Apache Pulsar` instances, depending on your configuration.

It collects broker statistics from
the [prometheus endpoint](https://pulsar.apache.org/docs/en/deploy-monitoring/#broker-stats).

`pulsar` module is tested on the following versions:

- v2.5.0

## Charts

It produces the following charts:

### Summary

- Broker Components in `num`
- Messages Rate in `messages/s`
- Throughput Rate in `KiB/s`
- Storage Size in `KiB`
- Messages Backlog Size in `messages`
- Storage Write Latency Histogram in `entries/s`
- Entry Size Histogram in `entries/s`
- Subscriptions Delayed for Dispatching in `message batches`

If `exposeTopicLevelMetricsInPrometheus` is set to true:

- Subscriptions Redelivered Message Rate in `messages/s`
- Subscriptions Blocked On Unacked Messages in `subscriptions`

If replication is configured and `replicationMetricsEnabled` is set to true:

- Replication Rate in `messages/s`
- Replication Throughput Rate in `messages/s`
- Replication Backlog in `messages`

### Namespace

- Broker Components in `num`
- Messages Rate in `messages/s`
- Throughput Rate in `KiB/s`
- Storage Size in `KiB`
- Storage Read/Write Operations Rate in `message batches/s`
- Messages Backlog Size in `messages`
- Storage Write Latency Histogram in `entries/s`
- Entry Size Histogram in `entries/s`
- Subscriptions Delayed for Dispatching in `message batches`

If `exposeTopicLevelMetricsInPrometheus` is set to true:

- Subscriptions Redelivered Message Rate in `messages/s`
- Subscriptions Blocked On Unacked Messages in `subscriptions`

If replication is configured and `replicationMetricsEnabled` is set to true:

- Replication Rate in `messages/s`
- Replication Throughput Rate in `messages/s`
- Replication Backlog in `messages`

### Topic

Topic charts are only available when `exposeTopicLevelMetricsInPrometheus` is set to true. In addition, you need to
set `topic_filer` configuration option. If you have a lot of topics this is highly unrecommended.

- Producers in `producers`
- Subscriptions in `producers`
- Consumers in `producers`
- Publish Messages Rate in `publishes/s`
- Dispatch Messages Rate in `dispatches/s`
- Publish Throughput Rate in `KiB/s`
- Dispatch Throughput Rate in `KiB/s`
- Storage Size in `KiB`
- Storage Read Rate in `message batches/s`
- Storage Write Rate in `message batches/s`
- Messages Backlog Size in `messages`
- Subscriptions Delayed for Dispatching in `message batches`
- Subscriptions Redelivered Message Rate in `messages/s`
- Subscriptions Blocked On Unacked Messages in `blocked subscriptions`

If replication is configured and `replicationMetricsEnabled` is set to true:

- Topic Replication Rate From Remote Cluster in `messages/s`
- Topic Replication Rate To Remote Cluster in `messages/s`
- Topic Replication Throughput Rate From Remote Cluster in `KiB/s`
- Topic Replication Throughput Rate To Remote Cluster in `KiB/s`
- Topic Replication Backlog in `KiB/s`

## Configuration

Edit the `go.d/pulsar.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

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

By default module collects data for all topics but it supports topic filtering. Filtering doesnt exclude a topic stats
from the [summary](#summary)/[namespace](#namespace) stats, it only removes the topic from the [topic](#topic) charts.

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

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m pulsar
```
