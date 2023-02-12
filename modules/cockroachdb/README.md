<!--
title: "CockroachDB monitoring with Netdata"
description: "Monitor the health and performance of CockroachDB databases with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/cockroachdb/README.md"
sidebar_label: "CockroachDB"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# CockroachDB monitoring with Netdata

[`CockroachDB`](https://www.cockroachlabs.com/)  is the SQL database for building global, scalable cloud services that
survive disasters.

This module will monitor one or more `CockroachDB` databases, depending on your configuration.

## Metrics

All metrics have "cockroachdb." prefix.

| Metric                               | Scope  |                                                                              Dimensions                                                                               |    Units     |
|--------------------------------------|:------:|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:------------:|
| process_cpu_time_combined_percentage | global |                                                                                 used                                                                                  |  percentage  |
| process_cpu_time_percentage          | global |                                                                               user, sys                                                                               |  percentage  |
| process_cpu_time                     | global |                                                                               user, sys                                                                               |      ms      |
| process_memory                       | global |                                                                                  rss                                                                                  |     KiB      |
| process_file_descriptors             | global |                                                                                 open                                                                                  |      fd      |
| process_uptime                       | global |                                                                                uptime                                                                                 |   seconds    |
| host_disk_bandwidth                  | global |                                                                              read, write                                                                              |     KiB      |
| host_disk_operations                 | global |                                                                             reads, writes                                                                             |  operations  |
| host_disk_iops_in_progress           | global |                                                                              in_progress                                                                              |     iops     |
| host_network_bandwidth               | global |                                                                            received, sent                                                                             |   kilobits   |
| host_network_packets                 | global |                                                                            received, sent                                                                             |   packets    |
| live_nodes                           | global |                                                                              live_nodes                                                                               |    nodes     |
| node_liveness_heartbeats             | global |                                                                          successful, failed                                                                           |  heartbeats  |
| total_storage_capacity               | global |                                                                                 total                                                                                 |     KiB      |
| storage_capacity_usability           | global |                                                                           usable, unusable                                                                            |     KiB      |
| storage_usable_capacity              | global |                                                                            available, used                                                                            |     KiB      |
| storage_used_capacity_percentage     | global |                                                                             total, usable                                                                             |  percentage  |
| sql_connections                      | global |                                                                                active                                                                                 | connections  |
| sql_bandwidth                        | global |                                                                            received, sent                                                                             |     KiB      |
| sql_statements_total                 | global |                                                                           started, executed                                                                           |  statements  |
| sql_errors                           | global |                                                                        statement, transaction                                                                         |    errors    |
| sql_started_ddl_statements           | global |                                                                                  ddl                                                                                  |  statements  |
| sql_executed_ddl_statements          | global |                                                                                  ddl                                                                                  |  statements  |
| sql_started_dml_statements           | global |                                                                    select, update, delete, insert                                                                     |  statements  |
| sql_executed_dml_statements          | global |                                                                    select, update, delete, insert                                                                     |  statements  |
| sql_started_tcl_statements           | global |             begin, commit, rollback, savepoint, savepoint_cockroach_restart, release_savepoint_cockroach_restart, rollback_to_savepoint_cockroach_restart             |  statements  |
| sql_executed_tcl_statements          | global |             begin, commit, rollback, savepoint, savepoint_cockroach_restart, release_savepoint_cockroach_restart, rollback_to_savepoint_cockroach_restart             |  statements  |
| sql_active_distributed_queries       | global |                                                                                active                                                                                 |   queries    |
| sql_distributed_flows                | global |                                                                            active, queued                                                                             |    flows     |
| live_bytes                           | global |                                                                         applications, system                                                                          |     KiB      |
| logical_data                         | global |                                                                             keys, values                                                                              |     KiB      |
| logical_data_count                   | global |                                                                             keys, values                                                                              |     num      |
| kv_transactions                      | global |                                                                committed, fast-path_committed, aborted                                                                | transactions |
| kv_transaction_restarts              | global | write_too_old, write_too_old_multiple, forwarded_timestamp, possible_reply, async_consensus_failure, read_within_uncertainty_interval, aborted, push_failure, unknown |   restarts   |
| ranges                               | global |                                                                                ranges                                                                                 |    ranges    |
| ranges_replication_problem           | global |                                                            unavailable, under_replicated, over_replicated                                                             |    ranges    |
| range_events                         | global |                                                                       split, add, remove, merge                                                                       |    events    |
| range_snapshot_events                | global |                                                generated, applied_raft_initiated, applied_learner, applied_preemptive                                                 |    events    |
| rocksdb_read_amplification           | global |                                                                                 reads                                                                                 | reads/query  |
| rocksdb_table_operations             | global |                                                                         compactions, flushes                                                                          |  operations  |
| rocksdb_cache_usage                  | global |                                                                                 used                                                                                  |     KiB      |
| rocksdb_cache_operations             | global |                                                                             hits, misses                                                                              |  operations  |
| rocksdb_cache_hit_rate               | global |                                                                               hit_rate                                                                                |  percentage  |
| rocksdb_sstables                     | global |                                                                               sstables                                                                                |   sstables   |
| replicas                             | global |                                                                               replicas                                                                                |   replicas   |
| replicas_quiescence                  | global |                                                                           quiescent, active                                                                           |   replicas   |
| replicas_leaders                     | global |                                                                       leaders, not_leaseholders                                                                       |   replicas   |
| replicas_leaseholders                | global |                                                                             leaseholders                                                                              | leaseholders |
| queue_processing_failures            | global |                                   gc, replica_gc, replication, split, consistency, raft_log, raft_snapshot, time_series_maintenance                                   |   failures   |
| rebalancing_queries                  | global |                                                                                  avg                                                                                  |  queries/s   |
| rebalancing_writes                   | global |                                                                                  avg                                                                                  |   writes/s   |
| timeseries_samples                   | global |                                                                                written                                                                                |   samples    |
| timeseries_write_errors              | global |                                                                                 write                                                                                 |    errors    |
| timeseries_write_bytes               | global |                                                                                written                                                                                |     KiB      |
| slow_requests                        | global |                                                              acquiring_latches, acquiring_lease, in_raft                                                              |   requests   |
| code_heap_memory_usage               | global |                                                                                go, cgo                                                                                |     KiB      |
| goroutines                           | global |                                                                              goroutines                                                                               |  goroutines  |
| gc_count                             | global |                                                                                  gc                                                                                   |   invokes    |
| gc_pause                             | global |                                                                                 pause                                                                                 |      us      |
| cgo_calls                            | global |                                                                                  cgo                                                                                  |    calls     |

## Configuration

Edit the `go.d/cockroachdb.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/cockroachdb.conf
```

Needs only `url` to server's `_status/vars`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8080/_status/vars

  - name: remote
    url: http://203.0.113.10:8080/_status/vars
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/cockroachdb.conf).

## Update every

Default `update_every` is 10 seconds because `CockroachDB` default sampling interval is 10 seconds, and it is not user
configurable. It doesn't make sense to decrease the value.

## Troubleshooting

To troubleshoot issues with the `cockroachdb` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m cockroachdb
  ```

