<!--
title: "Cassandra monitoring with Netdata"
description: "Monitor the health and performance of Cassandra database servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/cassandra/README.md
sidebar_label: "Cassandra"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# Cassandra collector

[Cassandra](https://cassandra.apache.org/_/index.html) is an open-source NoSQL database management system.

This module will monitor one or more Cassandra servers, depending on your configuration.

## Requirements

- Cassandra with [Prometheus JMX Exporter](https://github.com/prometheus/jmx_exporter).

To configure Cassandra with the JMX Exporter:

> **Note**: paths can differ depends on your setup.

- Download latest [jmx_exporter](https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/) jar file
  and install it in a directory where Cassandra can access it.
- Add
  the [jmx_exporter.yaml](https://raw.githubusercontent.com/netdata/go.d.plugin/master/modules/cassandra/jmx_exporter.yaml)
  file to `/etc/cassandra`.
- Add the following line to `/etc/cassandra/cassandra-env.sh`
  ```
  JVM_OPTS="$JVM_OPTS $JVM_EXTRA_OPTS -javaagent:/opt/jmx_exporter/jmx_exporter.jar=7072:/etc/cassandra/jmx_exporter.yaml
  ```
- Restart cassandra service.

## Metrics

All metrics have "cassandra." prefix.

Labels per scope:

- global: no labels.
- thread pool: thread_pool.

| Metric                                 |    Scope    |  Dimensions   |    Units     |
|----------------------------------------|:-----------:|:-------------:|:------------:|
| client_requests_rate                   |   global    |  read, write  |  requests/s  |
| client_request_read_latency_histogram  |   global    |  read, write  |   seconds    |
| client_request_write_latency_histogram |   global    |  read, write  |   seconds    |
| client_requests_latency                |   global    |  read, write  |   seconds    |
| row_cache_hit_ratio                    |   global    |   hit_ratio   |  percentage  |
| row_cache_hit_rate                     |   global    | hits, misses  |   events/s   |
| row_cache_utilization                  |   global    |     used      |  percentage  |
| row_cache_size                         |   global    |     size      |    bytes     |
| key_cache_hit_ratio                    |   global    |   hit_ratio   |  percentage  |
| key_cache_hit_rate                     |   global    | hits, misses  |   events/s   |
| key_cache_utilization                  |   global    |     used      |  percentage  |
| key_cache_size                         |   global    |     size      |    bytes     |
| storage_live_disk_space_used           |   global    |     used      |    bytes     |
| compaction_completed_tasks_rate        |   global    |   completed   |   tasks/s    |
| compaction_pending_tasks_count         |   global    |    pending    |    tasks     |
| compaction_compacted_rate              |   global    |   compacted   |   bytes/s    |
| thread_pool_active_tasks_count         | thread pool |    active     |    tasks     |
| thread_pool_pending_tasks_count        | thread pool |    pending    |    tasks     |
| thread_pool_blocked_tasks_count        | thread pool |    blocked    |    tasks     |
| thread_pool_blocked_tasks_rate         | thread pool |    blocked    |   tasks/s    |
| jvm_memory_used                        |   global    | heap, nonheap |    bytes     |
| jvm_gc_rate                            |   global    |  parnew, cms  |     gc/s     |
| jvm_gc_time                            |   global    |  parnew, cms  |   seconds    |
| dropped_messages_rate                  |   global    |    dropped    |  messages/s  |
| client_requests_timeouts_rate          |   global    |  read, write  |  timeout/s   |
| client_requests_unavailables_rate      |   global    |  read, write  | exceptions/s |
| client_requests_failures_rate          |   global    |  read, write  |  failures/s  |
| storage_exceptions_rate                |   global    |    storage    | exceptions/s |

## Configuration

Edit the `go.d/cassandra.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/cassandra.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:7072/metrics

  - name: remote
    url: http://203.0.113.10:7072/metrics
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/cassandra.conf).

## Troubleshooting

To troubleshoot issues with the `cassandra` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m cassandra
  ```
