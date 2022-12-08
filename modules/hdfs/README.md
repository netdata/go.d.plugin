<!--
title: "HDFS monitoring with Netdata"
description: "Monitor the health and performance of HDFS nodes with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/hdfs/README.md"
sidebar_label: "HDFS"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Storage"
-->

# HDFS monitoring with Netdata

The [`Hadoop Distributed File System (HDFS)`](https://hadoop.apache.org/docs/r1.2.1/hdfs_design.html) is a distributed
file system designed to run on commodity hardware.

This module monitors one or more `Hadoop Distributed File System` nodes, depending on your configuration.

Netdata accesses HDFS metrics over `Java Management Extensions` (JMX) through the web interface of an HDFS daemon.

## Requirements

- `hdfs` node with accessible `/jmx` endpoint

## Metrics

All metrics have "hdfs." prefix.

| Metric                  | Scope  |                         Dimensions                         |    Units    |
|-------------------------|:------:|:----------------------------------------------------------:|:-----------:|
| heap_memory             | global |                      committed, used                       |     MiB     |
| gc_count_total          | global |                             gc                             |  events/s   |
| gc_time_total           | global |                             ms                             |     ms      |
| gc_threshold            | global |                         info, warn                         |  events/s   |
| threads                 | global | new, runnable, blocked, waiting, timed_waiting, terminated |     num     |
| logs_total              | global |                  info, error, warn, fatal                  |   logs/s    |
| rpc_bandwidth           | global |                       received, sent                       | kilobits/s  |
| rpc_calls               | global |                           calls                            |   calls/s   |
| open_connections        | global |                            open                            | connections |
| call_queue_length       | global |                           length                           |     num     |
| avg_queue_time          | global |                            time                            |     ms      |
| avg_processing_time     | global |                            time                            |     ms      |
| capacity                | global |                      remaining, used                       |     KiB     |
| used_capacity           | global |                        dfs, non_dfs                        |     KiB     |
| load                    | global |                            load                            |    load     |
| volume_failures_total   | global |                          failures                          |  events/s   |
| files_total             | global |                           files                            |     num     |
| blocks_total            | global |                           blocks                           |     num     |
| blocks                  | global |             corrupt, missing, under_replicated             |     num     |
| data_nodes              | global |                     live, dead, stale                      |     num     |
| datanode_capacity       | global |                      remaining, used                       |     KiB     |
| datanode_used_capacity  | global |                        dfs, non_dfs                        |     KiB     |
| datanode_failed_volumes | global |                       failed volumes                       |     num     |
| datanode_bandwidth      | global |                       reads, writes                        |    KiB/s    |

## Configuration

Edit the `go.d/hdfs.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/hdfs.conf
```

Needs only `url` to server's `/jmx` endpoint. Here is an example for 2 servers:

```yaml
jobs:
  - name: namenode
    url: http://127.0.0.1:9870/jmx

  - name: datanode
    url: http://127.0.0.1:9864/jmx
```

For all available options, please see the
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/hdfs.conf).

## Troubleshooting

To troubleshoot issues with the `hdfs` collector, run the `go.d.plugin` with the debug option enabled. The output should
give you clues as to why the collector isn't working.

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
  ./go.d.plugin -d -m hdfs
  ```


