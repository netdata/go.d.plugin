<!--
title: "HDFS monitoring with Netdata"
description: "Monitor the health and performance of HDFS nodes with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/hdfs/README.md
sidebar_label: "HDFS"
-->

# HDFS monitoring with Netdata

The [`Hadoop Distributed File System (HDFS)`](https://hadoop.apache.org/docs/r1.2.1/hdfs_design.html) is a distributed
file system designed to run on commodity hardware.

This module monitors one or more `Hadoop Distributed File System` nodes, depending on your configuration.

Netdata accesses HDFS metrics over `Java Management Extensions` (JMX) through the web interface of an HDFS daemon.

## Requirements

- `hdfs` node with accessible `/jmx` endpoint

## Charts

It produces the following charts for `namenode`:

- Heap Memory in `MiB`
- GC Events in `events/s`
- GC Time in `ms`
- Number of Times That the GC Threshold is Exceeded in `events/s`
- Number of Threads in `num`
- Number of Logs in `logs/s`
- RPC Bandwidth in `kilobits/s`
- RPC Calls in `calls/s`
- RPC Open Connections in `connections`
- RPC Call Queue Length in `num`
- RPC Avg Queue Time in `ms`
- RPC Avg Processing Time in `ms`
- Capacity Across All Datanodes in `KiB`
- Used Capacity Across All Datanodes in `KiB`
- Number of Concurrent File Accesses (read/write) Across All DataNodes in `load`
- Number of Volume Failures Across All Datanodes in `events/s`
- Number of Tracked Files in `num`
- Number of Allocated Blocks in the System in `num`
- Number of Problem Blocks (can point to an unhealthy cluster) in `num`
- Number of Data Nodes By Status in `num`

For `datanode`:

- Heap Memory in `MiB`
- GC Events in `events/s`
- GC Time in `ms`
- Number of Times That the GC Threshold is Exceeded in `events/s`
- Number of Threads in `num`
- Number of Logs in `logs/s`
- RPC Bandwidth in `kilobits/s`
- RPC Calls in `calls/s`
- RPC Open Connections in `connections`
- RPC Call Queue Length in `num`
- RPC Avg Queue Time in `ms`
- RPC Avg Processing Time in `ms`
- Capacity in `KiB`
- Used Capacity in `KiB`
- Bandwidth in `KiB/s`

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

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m hdfs
```


