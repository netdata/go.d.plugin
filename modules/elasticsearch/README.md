<!--
title: "Elasticsearch monitoring with Netdata"
description: "Monitor the health and performance of Elasticsearch engines with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/elasticsearch/README.md"
sidebar_label: "Elasticsearch"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Apps"
-->

# Elasticsearch collector

[`Elasticsearch`](https://www.elastic.co/elasticsearch/) is a search engine based on the Lucene library.

This module monitors one or more `Elasticsearch` instances, depending on your configuration.

Used endpoints:

- Local node metrics: `/_nodes/_local/stats`
- Local node indices' metrics: `/_cat/indices?local=true`
- Cluster health metrics: `/_cluster/health`
- Cluster metrics: `/_cluster/stats`

Each endpoint can be enabled/disabled in the module configuration file.

## Metrics

All metrics have "elasticsearch." prefix.

Labels per scope:

- global: no labels.
- index: index.

| Metric                                   | Scope  |                                                                             Dimensions                                                                              |    Units     |
|------------------------------------------|:------:|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:------------:|
| node_indices_indexing                    | global |                                                                                index                                                                                | operations/s |
| node_indices_indexing_current            | global |                                                                                index                                                                                |  operations  |
| node_indices_indexing_time               | global |                                                                                index                                                                                | milliseconds |
| node_indices_search                      | global |                                                                          queries, fetches                                                                           | operations/s |
| node_indices_search_current              | global |                                                                          queries, fetches                                                                           |  operations  |
| node_indices_search_time                 | global |                                                                          queries, fetches                                                                           | milliseconds |
| node_indices_refresh                     | global |                                                                               refresh                                                                               | operations/s |
| node_indices_refresh_time                | global |                                                                               refresh                                                                               | milliseconds |
| node_indices_flush                       | global |                                                                                flush                                                                                | operations/s |
| node_indices_flush_time                  | global |                                                                                flush                                                                                | milliseconds |
| node_indices_fielddata_memory_usage      | global |                                                                                used                                                                                 |    bytes     |
| node_indices_fielddata_evictions         | global |                                                                              evictions                                                                              | operations/s |
| node_indices_segments_count              | global |                                                                              segments                                                                               |   segments   |
| node_indices_segments_memory_usage_total | global |                                                                                used                                                                                 |    bytes     |
| node_indices_segments_memory_usage       | global |                               terms, stored_fields, term_vectors, norms, points, doc_values, index_writer, version_map, fixed_bit_set                               |    bytes     |
| node_indices_translog_operations         | global |                                                                         total, uncommitted                                                                          |  operations  |
| node_indices_translog_size               | global |                                                                         total, uncommitted                                                                          |    bytes     |
| node_file_descriptors                    | global |                                                                                open                                                                                 |      fd      |
| node_jvm_heap                            | global |                                                                                inuse                                                                                |  percentage  |
| node_jvm_heap_bytes                      | global |                                                                           committed, used                                                                           |    bytes     |
| node_jvm_buffer_pools_count              | global |                                                                           direct, mapped                                                                            |    pools     |
| node_jvm_buffer_pool_direct_memory       | global |                                                                             total, used                                                                             |    bytes     |
| node_jvm_buffer_pool_mapped_memory       | global |                                                                             total, used                                                                             |    bytes     |
| node_jvm_gc_count                        | global |                                                                             young, old                                                                              |     gc/s     |
| node_jvm_gc_time                         | global |                                                                             young, old                                                                              | milliseconds |
| node_thread_pool_queued                  | global | generic, search, search_throttled, get, analyze, write, snapshot, warmer, refresh, listener, fetch_shard_started, fetch_shard_store, flush, force_merge, management |   threads    |
| node_thread_pool_rejected                | global | generic, search, search_throttled, get, analyze, write, snapshot, warmer, refresh, listener, fetch_shard_started, fetch_shard_store, flush, force_merge, management |   threads    |
| cluster_communication_packets            | global |                                                                           received, sent                                                                            |     pps      |
| cluster_communication                    | global |                                                                           received, sent                                                                            |   bytes/s    |
| http_connections                         | global |                                                                                open                                                                                 | connections  |
| breakers_trips                           | global |                                            requests, fielddata, in_flight_requests, model_inference, accounting, parent                                             |   trips/s    |
| http_connections                         | global |                                                                                open                                                                                 | connections  |
| cluster_health_status                    | global |                                                                         green, yellow, red                                                                          |    status    |
| cluster_number_of_nodes                  | global |                                                                          nodes, data_nodes                                                                          |    nodes     |
| cluster_shards_count                     | global |                                          active_primary, active, relocating, initializing, unassigned, delayed_unaasigned                                           |    shards    |
| cluster_pending_tasks                    | global |                                                                               pending                                                                               |    tasks     |
| cluster_number_of_in_flight_fetch        | global |                                                                           in_flight_fetch                                                                           |   fetches    |
| cluster_indices_count                    | global |                                                                               indices                                                                               |   indices    |
| cluster_indices_shards_count             | global |                                                                    total, primaries, replication                                                                    |    shards    |
| cluster_indices_docs_count               | global |                                                                                docs                                                                                 |     docs     |
| cluster_indices_store_size               | global |                                                                                size                                                                                 |    bytes     |
| cluster_indices_query_cache              | global |                                                                              hit, miss                                                                              |   events/s   |
| cluster_nodes_by_role_count              | global |                                           coordinating_only, data, ingest, master, ml, remote_cluster_client, voting_only                                           |    nodes     |
| node_index_health                        | index  |                                                                         green, yellow, red                                                                          |    status    |
| node_index_shards_count                  | index  |                                                                               shards                                                                                |    shards    |
| node_index_docs_count                    | index  |                                                                                docs                                                                                 |     docs     |
| node_index_store_size                    | index  |                                                                             store_size                                                                              |    bytes     |

## Configuration

Edit the `go.d/elasticsearch.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

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
  ./go.d.plugin -d -m elasticsearch
  ```
