# CockroachDB monitoring with Netdata

[`CockroachDB`](https://www.cockroachlabs.com/)  is the SQL database for building global, scalable cloud services that survive disasters.

This module will monitor one or more `CockroachDB` databases, depending on your configuration.

## Charts

It produces the following charts:

#### Process Statistics

-   Combined CPU Time Percentage, Normalized 0-1 by Number of Cores in `percantage`
-   CPU Time Percentage in `percentage`
-   CPU Time in `ms`
-   Memory Usage in `KiB`
-   File Descriptors in `fd`
-   Uptime in `seconds`

#### Host Statistics

-   Host Disk Cumulative Bandwidth in `KiB`
-   Host Disk Cumulative Operations in `operations`
-   Host Disk Cumulative IOPS In Progress in `iops`
-   Host Network Cumulative Bandwidth in `kilobits`
-   Host Network Cumulative Packets in `packets`
-   Uptime in `seconds`

#### Liveness

-   Live Nodes in the Cluster in `num`
-   Node Liveness Heartbeats in `heartbeats`

#### Capacity

-   Total Storage Capacity in `KiB`
-   Storage Capacity Usability in `KiB`
-   Storage Usable Capacity in `KiB`
-   Storage Used Capacity Utilization in `percentage`

#### SQL

-   Active SQL Connections in `connections`
-   SQL Bandwidth in `KiB`
-   SQL Statements Total in `statements`
-   SQL Statements and Transaction Errors in `errors`
-   SQL Started DDL Statements in `statements`
-   SQL Executed DDL Statements in `statements`
-   SQL Started DML Statements in `statements`
-   SQL Executed DML Statements in `statements`
-   SQL Started TCL Statements in `statements`
-   SQL Executed TCL Statements in `statements`
-   Active Distributed SQL Queries in `queries`
-   Distributed SQL Flows in `flows`

#### Storage

-   Used Live Data in `KiB`
-   Logical Data in `KiB`
-   Logical Data Count in `num`

#### KV Transactions

-   KV Transactions in `transactions`
-   KV Transaction Restarts in `restarts`

#### Ranges

-   Ranges in `num`
-   Problem Ranges in `ranges`
-   Range Events in `events`
-   Range Snapshot Events in `events`

#### RocksDB

-   RocksDB Read Amplification in `reads/query`
-   RocksDB Table Operations in `operations`
-   RocksDB Block Cache Operations in `operations`
-   RocksDB Block Cache Hit Rate in `percentage`
-   RocksDB SSTables in `num`

#### Replication

-   Number of Replicas in `num`
-   Replicas Quiescence in `replicas`
-   Number of Raft Leaders in `num`
-   Number of Leaseholders in `num`
-   RocksDB SSTables in `num`

#### Queues

-   Queues Processing Failures in `failures`

#### Rebalancing

-   Rebalancing Average Queries in `queries/s`
-   Rebalancing Average Writes in `writes/s`

#### Time Series

-   Time Series Written Samples in `samples`
-   Time Series Write Errors in `errors`
-   Time Series Bytes Written in `KiB`

#### Slow Requests

-   Slow Requests in `requests`

#### Go/Cgo

-   Heap Memory Usage in `KiB`
-   Number of Goroutines in `num`
-   GC Runs in `invokes`
-   GC Pause Time in `us`
-   Cgo Calls in `calls`

## Configuration

Edit the `go.d/cockroachdb.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

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

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/cockroachdb.conf).

## Update every

Default `update_every` is 10 seconds because `CockroachDB` default sampling interval is 10 seconds and it is not user configurable.
It doesnt make sense to decrease the value.

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m cockroachdb
