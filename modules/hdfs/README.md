# HDFS monitoring with Netdata

The [`Hadoop Distributed File System (HDFS)`](https://hadoop.apache.org/docs/r1.2.1/hdfs_design.html) is a distributed file system designed to run on commodity hardware.

This module monitors one or more `Hadoop Distributed File System` nodes, depending on your configuration.

Netdata accesses HDFS metrics over `Java Management Extensions` (JMX) through the web interface of an HDFS daemon.

## Requirements

-   `hdfs` node with accessible `/jmx` endpoint

## Charts

It produces the following charts for `namenode`:

-   Heap Memory in `MiB`
-   GC Events in `events/s`
-   GC Time in `ms`
-   Number of Times That the GC Threshold is Exceeded in `events/s`
-   Number of Threads in `num`
-   Number of Logs in `logs/s`
-   RPC Bandwidth in `kilobits/s`
-   RPC Calls in `calls/s`
-   RPC Open Connections in `connections`
-   RPC Call Queue Length in `num`
-   RPC Avg Queue Time in `ms`
-   RPC Avg Processing Time in `ms`
-   Capacity Across All Datanodes in `KiB`
-   Used Capacity Across All Datanodes in `KiB`
-   Number of Concurrent File Accesses (read/write) Across All DataNodes in `load`
-   Number of Volume Failures Across All Datanodes in `events/s`
-   Number of Tracked Files in `num`
-   Number of Allocated Blocks in the System in `num`
-   Number of Problem Blocks (can point to an unhealthy cluster) in `num`
-   Number of Data Nodes By Status in `num`
  
For `datanode`:

-   Heap Memory in `MiB`
-   GC Events in `events/s`
-   GC Time in `ms`
-   Number of Times That the GC Threshold is Exceeded in `events/s`
-   Number of Threads in `num`
-   Number of Logs in `logs/s`
-   RPC Bandwidth in `kilobits/s`
-   RPC Calls in `calls/s`
-   RPC Open Connections in `connections`
-   RPC Call Queue Length in `num`
-   RPC Avg Queue Time in `ms`
-   RPC Avg Processing Time in `ms`
-   Capacity in `KiB`
-   Used Capacity in `KiB`
-   Bandwidth in `KiB/s`

## Configuration

Edit the `go.d/hdfs.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

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

For all available options, please see the module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/hdfs.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m hdfs

