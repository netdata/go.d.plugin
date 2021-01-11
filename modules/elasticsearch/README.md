<!--
title: "Elasticsearch monitoring with Netdata"
description: "Monitor the health and performance of Elasticsearch engines with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/elasticsearch/README.md
sidebar_label: "Elasticsearch"
-->

# Elasticsearch monitoring with Netdata

[`Elasticsearch`](https://www.elastic.co/elasticsearch/) is a search engine based on the Lucene library.

This module monitors one or more `Elasticsearch` instances, depending on your configuration.

Used endpoints:

- Local node metrics: `/_nodes/_local/stats`
- Local node indices' metrics: `/_cat/indices?local=true`
- Cluster health metrics: `/_cluster/health`
- Cluster metrics: `/_cluster/stats`

Each endpoint can be enabled/disabled in the module configuration file.

## Charts

Number of charts depends on enabled endpoints.

### Local Node Stats

Collected from `/_nodes/_local/stats` endpoint. Controlled by `collect_node_stats` option. Enabled by default.

- Indexing Operations in `operations/s`
- Indexing Operations Current in `operations`
- Time Spent On Indexing Operations in `milliseconds`
- Search Operations in `operations/s`
- Search Operations Current in `operations`
- Time Spent On Search Operations in `milliseconds`
- Refresh Operations in `operations/s`
- Time Spent On Refresh Operations in `milliseconds`
- Flush Operations in `operations/s`
- Time Spent On Flush Operations in `milliseconds`
- Fielddata Cache Memory Usage in `bytes`
- Fielddata Evictions in `operations/s`
- Segments Count in `segments`
- Segments Memory Usage Total in `bytes`
- Segments Memory Usage in `bytes`
- Translog Operations in `operations`
- Translog Size in `bytes`
- Process File Descriptors in `fd`
- JVM Heap Percentage Currently in Use in `percentage`
- JVM Heap Commit And Usage in `bytes`
- JVM Buffer Pools Count in `pools`
- JVM Buffer Pool Direct Memory in `bytes`
- JVM Buffer Pool Mapped Memory in `bytes`
- JVM Garbage Collections in `gc/s`
- JVM Time Spent On Garbage Collections in `milliseconds`
- Thread Pool Queued Threads Count in `threads`
- Thread Pool Rejected Threads Count in `threads`
- Cluster Communication in `pps`
- Cluster Communication Bandwidth in `bytes/s`
- HTTP Connections in `connections`
- Circuit Breaker Trips Count in `trips/s`

### Local Indices Stats

Collected from `/_cat/indices?local=true` endpoint. Controlled by `collect_indices_stats` option. Disabled by default.

- Index Health in `status`
- Index Shards Count in `shards`
- Index Docs Count in `docs`
- Index Store Size in `bytes`

### Cluster Health

Collected from `/_cluster/health` endpoint. Controlled by `collect_cluster_health` option. Enabled by default.

- Cluster Status in `status`
- Cluster Nodes Count in `nodes`
- Cluster Shards Count in `shards`
- Cluster Pending Tasks in `tasks`
- Cluster Unfinished Fetches in `fetches`

### Cluster Stats

Collected from `/_cluster/stats` endpoint. Controlled by `collect_cluster_stats` option. Enabled by default.

- Cluster Indices Count in `indices`
- Cluster Indices Shards Count in `shards`
- Cluster Indices Docs Count in `docs`
- Cluster Indices Store Size in `bytes`
- Cluster Indices Query Cache in `events/s`
- Cluster Nodes By Role Count in `nodes`

## Configuration

Edit the `go.d/elasticsearch.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/elasticsearch.conf
```

To add a new endpoint to collect metrics from, or change the URL that Netdata looks for, add or configure the `name` and
`url` values. Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

Here is an example with two endpoints:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9200

  - name: remote
    url: http://203.0.113.0:9200
```

For all available options, see the Elasticsearch
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/elasticsearch.conf).

## Troubleshooting

To troubleshoot issues with the `elasticsearch` collector, run the `go.d.plugin` with the debug option enabled. The
output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m elasticsearch
```
