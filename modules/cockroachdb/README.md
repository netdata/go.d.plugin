# CockroachDB collector

## Overview

[CockroachDB](https://www.cockroachlabs.com/)  is the SQL database for building global, scalable cloud services that
survive disasters.

This collector monitors one or more CockroachDB databases, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                           |                                                                              Dimensions                                                                               |     Unit     |
|--------------------------------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:------------:|
| cockroachdb.process_cpu_time_combined_percentage |                                                                                 used                                                                                  |  percentage  |
| cockroachdb.process_cpu_time_percentage          |                                                                               user, sys                                                                               |  percentage  |
| cockroachdb.process_cpu_time                     |                                                                               user, sys                                                                               |      ms      |
| cockroachdb.process_memory                       |                                                                                  rss                                                                                  |     KiB      |
| cockroachdb.process_file_descriptors             |                                                                                 open                                                                                  |      fd      |
| cockroachdb.process_uptime                       |                                                                                uptime                                                                                 |   seconds    |
| cockroachdb.host_disk_bandwidth                  |                                                                              read, write                                                                              |     KiB      |
| cockroachdb.host_disk_operations                 |                                                                             reads, writes                                                                             |  operations  |
| cockroachdb.host_disk_iops_in_progress           |                                                                              in_progress                                                                              |     iops     |
| cockroachdb.host_network_bandwidth               |                                                                            received, sent                                                                             |   kilobits   |
| cockroachdb.host_network_packets                 |                                                                            received, sent                                                                             |   packets    |
| cockroachdb.live_nodes                           |                                                                              live_nodes                                                                               |    nodes     |
| cockroachdb.node_liveness_heartbeats             |                                                                          successful, failed                                                                           |  heartbeats  |
| cockroachdb.total_storage_capacity               |                                                                                 total                                                                                 |     KiB      |
| cockroachdb.storage_capacity_usability           |                                                                           usable, unusable                                                                            |     KiB      |
| cockroachdb.storage_usable_capacity              |                                                                            available, used                                                                            |     KiB      |
| cockroachdb.storage_used_capacity_percentage     |                                                                             total, usable                                                                             |  percentage  |
| cockroachdb.sql_connections                      |                                                                                active                                                                                 | connections  |
| cockroachdb.sql_bandwidth                        |                                                                            received, sent                                                                             |     KiB      |
| cockroachdb.sql_statements_total                 |                                                                           started, executed                                                                           |  statements  |
| cockroachdb.sql_errors                           |                                                                        statement, transaction                                                                         |    errors    |
| cockroachdb.sql_started_ddl_statements           |                                                                                  ddl                                                                                  |  statements  |
| cockroachdb.sql_executed_ddl_statements          |                                                                                  ddl                                                                                  |  statements  |
| cockroachdb.sql_started_dml_statements           |                                                                    select, update, delete, insert                                                                     |  statements  |
| cockroachdb.sql_executed_dml_statements          |                                                                    select, update, delete, insert                                                                     |  statements  |
| cockroachdb.sql_started_tcl_statements           |             begin, commit, rollback, savepoint, savepoint_cockroach_restart, release_savepoint_cockroach_restart, rollback_to_savepoint_cockroach_restart             |  statements  |
| cockroachdb.sql_executed_tcl_statements          |             begin, commit, rollback, savepoint, savepoint_cockroach_restart, release_savepoint_cockroach_restart, rollback_to_savepoint_cockroach_restart             |  statements  |
| cockroachdb.sql_active_distributed_queries       |                                                                                active                                                                                 |   queries    |
| cockroachdb.sql_distributed_flows                |                                                                            active, queued                                                                             |    flows     |
| cockroachdb.live_bytes                           |                                                                         applications, system                                                                          |     KiB      |
| cockroachdb.logical_data                         |                                                                             keys, values                                                                              |     KiB      |
| cockroachdb.logical_data_count                   |                                                                             keys, values                                                                              |     num      |
| cockroachdb.kv_transactions                      |                                                                committed, fast-path_committed, aborted                                                                | transactions |
| cockroachdb.kv_transaction_restarts              | write_too_old, write_too_old_multiple, forwarded_timestamp, possible_reply, async_consensus_failure, read_within_uncertainty_interval, aborted, push_failure, unknown |   restarts   |
| cockroachdb.ranges                               |                                                                                ranges                                                                                 |    ranges    |
| cockroachdb.ranges_replication_problem           |                                                            unavailable, under_replicated, over_replicated                                                             |    ranges    |
| cockroachdb.range_events                         |                                                                       split, add, remove, merge                                                                       |    events    |
| cockroachdb.range_snapshot_events                |                                                generated, applied_raft_initiated, applied_learner, applied_preemptive                                                 |    events    |
| cockroachdb.rocksdb_read_amplification           |                                                                                 reads                                                                                 | reads/query  |
| cockroachdb.rocksdb_table_operations             |                                                                         compactions, flushes                                                                          |  operations  |
| cockroachdb.rocksdb_cache_usage                  |                                                                                 used                                                                                  |     KiB      |
| cockroachdb.rocksdb_cache_operations             |                                                                             hits, misses                                                                              |  operations  |
| cockroachdb.rocksdb_cache_hit_rate               |                                                                               hit_rate                                                                                |  percentage  |
| cockroachdb.rocksdb_sstables                     |                                                                               sstables                                                                                |   sstables   |
| cockroachdb.replicas                             |                                                                               replicas                                                                                |   replicas   |
| cockroachdb.replicas_quiescence                  |                                                                           quiescent, active                                                                           |   replicas   |
| cockroachdb.replicas_leaders                     |                                                                       leaders, not_leaseholders                                                                       |   replicas   |
| cockroachdb.replicas_leaseholders                |                                                                             leaseholders                                                                              | leaseholders |
| cockroachdb.queue_processing_failures            |                                   gc, replica_gc, replication, split, consistency, raft_log, raft_snapshot, time_series_maintenance                                   |   failures   |
| cockroachdb.rebalancing_queries                  |                                                                                  avg                                                                                  |  queries/s   |
| cockroachdb.rebalancing_writes                   |                                                                                  avg                                                                                  |   writes/s   |
| cockroachdb.timeseries_samples                   |                                                                                written                                                                                |   samples    |
| cockroachdb.timeseries_write_errors              |                                                                                 write                                                                                 |    errors    |
| cockroachdb.timeseries_write_bytes               |                                                                                written                                                                                |     KiB      |
| cockroachdb.slow_requests                        |                                                              acquiring_latches, acquiring_lease, in_raft                                                              |   requests   |
| cockroachdb.code_heap_memory_usage               |                                                                                go, cgo                                                                                |     KiB      |
| cockroachdb.goroutines                           |                                                                              goroutines                                                                               |  goroutines  |
| cockroachdb.gc_count                             |                                                                                  gc                                                                                   |   invokes    |
| cockroachdb.gc_pause                             |                                                                                 pause                                                                                 |      us      |
| cockroachdb.cgo_calls                            |                                                                                  cgo                                                                                  |    calls     |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/cockroachdb.conf`.

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
sudo ./edit-config go.d/cockroachdb.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |              Default               | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:----------------------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |                 10                 |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |                 0                  |          |
|         url          | Server URL.                                                                                               | http://127.0.0.1:8080/_status/vars |   yes    |
|       timeout        | HTTP request timeout.                                                                                     |                 1                  |          |
|       username       | Username for basic HTTP authentication.                                                                   |                                    |          |
|       password       | Password for basic HTTP authentication.                                                                   |                                    |          |
|      proxy_url       | Proxy URL.                                                                                                |                                    |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                                    |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                                    |          |
|        method        | HTTP request method.                                                                                      |                GET                 |          |
|         body         | HTTP request body.                                                                                        |                                    |          |
|       headers        | HTTP request header.                                                                                      |                                    |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |                 no                 |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |                 no                 |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                                    |          |
|       tls_cert       | Client TLS certificate.                                                                                   |                                    |          |
|       tls_key        | Client TLS key.                                                                                           |                                    |          |

</details>

#### Examples

##### Basic

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8080/_status/vars
```

</details>

##### HTTP authentication

Local server with basic HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8080/_status/vars
    username: username
    password: password
```

</details>

##### HTTPS with self-signed certificate

CockroachDB with enabled HTTPS and self-signed certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: https://127.0.0.1:8080/_status/vars
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
    url: http://127.0.0.1:8080/_status/vars

  - name: remote
    url: http://203.0.113.10:8080/_status/vars
```

</details>

## Troubleshooting

### Debug mode

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
