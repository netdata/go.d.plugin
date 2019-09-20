# hdfs

This module will monitor one or more [`Hadoop Distributed File System`](https://hadoop.apache.org/docs/r1.2.1/hdfs_design.html) (HDFS) nodes over `Java Management Extensions` (JMX) through the web interface of an HDFS daemon.

**Requirements:**
 * `hdfs` node with accessible `/jmx` endpoint

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
  
It produces the following charts for `datanode`:
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
  

### configuration

Needs only `url` to server's `/jmx` endpoint.

Here is an example for 2 servers:

```yaml
jobs:
  - name    : master
    url     : http://127.0.0.1:9870/jmx
      
  - name    : slave
    url     : http://127.0.0.1:9864/jmx
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/hdfs.conf).

---
