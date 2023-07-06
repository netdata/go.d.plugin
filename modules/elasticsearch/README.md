# Elasticsearch/OpenSearch collector

## Overview

[Elasticsearch](https://www.elastic.co/elasticsearch/) is a search engine based on the Lucene library. The original
Elasticsearch project was continued as an open-source project called [OpenSearch](https://opensearch.org/) by Amazon.

This collector monitors metrics from one or more Elasticsearch/OpenSearch servers, depending on your configuration.

Used endpoints:

- Info: `/`
- Nodes metrics: `/_nodes/stats`
- Local node metrics: `/_nodes/_local/stats`
- Local node indices' metrics: `/_cat/indices?local=true`
- Cluster health metrics: `/_cluster/health`
- Cluster metrics: `/_cluster/stats`

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### node

These metrics refer to the cluster node.

Labels:

| Label        | Description                                                                                                                                                                  |
|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| cluster_name | Name of the cluster. Based on the [Cluster name setting](https://www.elastic.co/guide/en/elasticsearch/reference/current/important-settings.html#cluster-name).              |
| node_name    | Human-readable identifier for the node. Based on the [Node name setting](https://www.elastic.co/guide/en/elasticsearch/reference/current/important-settings.html#node-name). |
| host         | Network host for the node, based on the [Network host setting](https://www.elastic.co/guide/en/elasticsearch/reference/current/important-settings.html#network.host).        |

Metrics:

| Metric                                                 |                                                                             Dimensions                                                                              |     Unit     |
|--------------------------------------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:------------:|
| elasticsearch.node_indices_indexing                    |                                                                                index                                                                                | operations/s |
| elasticsearch.node_indices_indexing_current            |                                                                                index                                                                                |  operations  |
| elasticsearch.node_indices_indexing_time               |                                                                                index                                                                                | milliseconds |
| elasticsearch.node_indices_search                      |                                                                          queries, fetches                                                                           | operations/s |
| elasticsearch.node_indices_search_current              |                                                                          queries, fetches                                                                           |  operations  |
| elasticsearch.node_indices_search_time                 |                                                                          queries, fetches                                                                           | milliseconds |
| elasticsearch.node_indices_refresh                     |                                                                               refresh                                                                               | operations/s |
| elasticsearch.node_indices_refresh_time                |                                                                               refresh                                                                               | milliseconds |
| elasticsearch.node_indices_flush                       |                                                                                flush                                                                                | operations/s |
| elasticsearch.node_indices_flush_time                  |                                                                                flush                                                                                | milliseconds |
| elasticsearch.node_indices_fielddata_memory_usage      |                                                                                used                                                                                 |    bytes     |
| elasticsearch.node_indices_fielddata_evictions         |                                                                              evictions                                                                              | operations/s |
| elasticsearch.node_indices_segments_count              |                                                                              segments                                                                               |   segments   |
| elasticsearch.node_indices_segments_memory_usage_total |                                                                                used                                                                                 |    bytes     |
| elasticsearch.node_indices_segments_memory_usage       |                               terms, stored_fields, term_vectors, norms, points, doc_values, index_writer, version_map, fixed_bit_set                               |    bytes     |
| elasticsearch.node_indices_translog_operations         |                                                                         total, uncommitted                                                                          |  operations  |
| elasticsearch.node_indices_translog_size               |                                                                         total, uncommitted                                                                          |    bytes     |
| elasticsearch.node_file_descriptors                    |                                                                                open                                                                                 |      fd      |
| elasticsearch.node_jvm_heap                            |                                                                                inuse                                                                                |  percentage  |
| elasticsearch.node_jvm_heap_bytes                      |                                                                           committed, used                                                                           |    bytes     |
| elasticsearch.node_jvm_buffer_pools_count              |                                                                           direct, mapped                                                                            |    pools     |
| elasticsearch.node_jvm_buffer_pool_direct_memory       |                                                                             total, used                                                                             |    bytes     |
| elasticsearch.node_jvm_buffer_pool_mapped_memory       |                                                                             total, used                                                                             |    bytes     |
| elasticsearch.node_jvm_gc_count                        |                                                                             young, old                                                                              |     gc/s     |
| elasticsearch.node_jvm_gc_time                         |                                                                             young, old                                                                              | milliseconds |
| elasticsearch.node_thread_pool_queued                  | generic, search, search_throttled, get, analyze, write, snapshot, warmer, refresh, listener, fetch_shard_started, fetch_shard_store, flush, force_merge, management |   threads    |
| elasticsearch.node_thread_pool_rejected                | generic, search, search_throttled, get, analyze, write, snapshot, warmer, refresh, listener, fetch_shard_started, fetch_shard_store, flush, force_merge, management |   threads    |
| elasticsearch.node_cluster_communication_packets       |                                                                           received, sent                                                                            |     pps      |
| elasticsearch.node_cluster_communication_traffic       |                                                                           received, sent                                                                            |   bytes/s    |
| elasticsearch.node_http_connections                    |                                                                                open                                                                                 | connections  |
| elasticsearch.node_breakers_trips                      |                                            requests, fielddata, in_flight_requests, model_inference, accounting, parent                                             |   trips/s    |

### cluster

These metrics refer to the cluster.

Labels:

| Label        | Description                                                                                                                                                     |
|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| cluster_name | Name of the cluster. Based on the [Cluster name setting](https://www.elastic.co/guide/en/elasticsearch/reference/current/important-settings.html#cluster-name). |

Metrics:

| Metric                                          |                                                                 Dimensions                                                                 |   Unit   |
|-------------------------------------------------|:------------------------------------------------------------------------------------------------------------------------------------------:|:--------:|
| elasticsearch.cluster_health_status             |                                                             green, yellow, red                                                             |  status  |
| elasticsearch.cluster_number_of_nodes           |                                                             nodes, data_nodes                                                              |  nodes   |
| elasticsearch.cluster_shards_count              |                              active_primary, active, relocating, initializing, unassigned, delayed_unaasigned                              |  shards  |
| elasticsearch.cluster_pending_tasks             |                                                                  pending                                                                   |  tasks   |
| elasticsearch.cluster_number_of_in_flight_fetch |                                                              in_flight_fetch                                                               | fetches  |
| elasticsearch.cluster_indices_count             |                                                                  indices                                                                   | indices  |
| elasticsearch.cluster_indices_shards_count      |                                                       total, primaries, replication                                                        |  shards  |
| elasticsearch.cluster_indices_docs_count        |                                                                    docs                                                                    |   docs   |
| elasticsearch.cluster_indices_store_size        |                                                                    size                                                                    |  bytes   |
| elasticsearch.cluster_indices_query_cache       |                                                                 hit, miss                                                                  | events/s |
| elasticsearch.cluster_nodes_by_role_count       | coordinating_only, data, data_cold, data_content, data_frozen, data_hot, data_warm, ingest, master, ml, remote_cluster_client, voting_only |  nodes   |

### index

These metrics refer to the index.

Labels:

| Label        | Description                                                                                                                                                     |
|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| cluster_name | Name of the cluster. Based on the [Cluster name setting](https://www.elastic.co/guide/en/elasticsearch/reference/current/important-settings.html#cluster-name). |
| index        | Name of the index.                                                                                                                                              |

Metrics:

| Metric                                |     Dimensions     |  Unit  |
|---------------------------------------|:------------------:|:------:|
| elasticsearch.node_index_health       | green, yellow, red | status |
| elasticsearch.node_index_shards_count |       shards       | shards |
| elasticsearch.node_index_docs_count   |        docs        |  docs  |
| elasticsearch.node_index_store_size   |     store_size     | bytes  |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/elasticsearch.conf`.

The file format is YAML. Generally, the format is:

```yaml
update_every: 1
autodetection_retry: 0
jobs:
  - name: some_name1
  - name: some_name1
```

You can edit the configuration file using the `edit-config` script from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md#the-netdata-config-directory).

```bash
cd /etc/netdata 2>/dev/null || cd /opt/netdata/etc/netdata
sudo ./edit-config go.d/elasticsearch.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|          Name          | Description                                                                                               |        Default        | Required |
|:----------------------:|-----------------------------------------------------------------------------------------------------------|:---------------------:|:--------:|
|      update_every      | Data collection frequency.                                                                                |           5           |          |
|  autodetection_retry   | Re-check interval in seconds. Zero means not to schedule re-check.                                        |           0           |          |
|          url           | Server URL.                                                                                               | http://127.0.0.1:9200 |   yes    |
|      cluster_mode      | Controls whether to collect metrics for all nodes in the cluster or only for the local node.              |         false         |          |
|   collect_node_stats   | Controls whether to collect nodes metrics.                                                                |         true          |          |
| collect_cluster_health | Controls whether to collect cluster health metrics.                                                       |         true          |          |
| collect_cluster_stats  | Controls whether to collect cluster stats metrics.                                                        |         true          |          |
| collect_indices_stats  | Controls whether to collect indices metrics.                                                              |         false         |          |
|        timeout         | HTTP request timeout.                                                                                     |           5           |          |
|        username        | Username for basic HTTP authentication.                                                                   |                       |          |
|        password        | Password for basic HTTP authentication.                                                                   |                       |          |
|       proxy_url        | Proxy URL.                                                                                                |                       |          |
|     proxy_username     | Username for proxy basic HTTP authentication.                                                             |                       |          |
|     proxy_password     | Password for proxy basic HTTP authentication.                                                             |                       |          |
|         method         | HTTP request method.                                                                                      |          GET          |          |
|          body          | HTTP request body.                                                                                        |                       |          |
|        headers         | HTTP request headers.                                                                                     |                       |          |
|  not_follow_redirects  | Redirect handling policy. Controls whether the client follows redirects.                                  |          no           |          |
|    tls_skip_verify     | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |          no           |          |
|         tls_ca         | Certification authority that the client uses when verifying the server's certificates.                    |                       |          |
|        tls_cert        | Client TLS certificate.                                                                                   |                       |          |
|        tls_key         | Client TLS key.                                                                                           |                       |          |

</details>

#### Examples

##### Basic single node mode

A basic example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9200
```

</details>

##### Cluster mode

Cluster mode example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9200
    cluster_mode: yes
```

</details>

##### HTTP authentication

Basic HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9200
    username: username
    password: password
```

</details>

##### HTTPS with self-signed certificate

Elasticsearch with enabled HTTPS and self-signed certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: https://127.0.0.1:9200
    tls_skip_verify: yes
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Collecting metrics from local and remote instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9200

  - name: remote
    url: http://192.0.2.1:9200
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `elasticsearch` collector, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

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

