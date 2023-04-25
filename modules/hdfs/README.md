<!--
title: "HDFS monitoring with Netdata"
description: "Monitor the health and performance of HDFS nodes with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/hdfs/README.md"
sidebar_label: "HDFS"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Storage"
-->

# HDFS collector

The [`Hadoop Distributed File System (HDFS)`](https://hadoop.apache.org/docs/r1.2.1/hdfs_design.html) is a distributed
file system designed to run on commodity hardware.

This module monitors one or more `Hadoop Distributed File System` nodes, depending on your configuration.

Netdata accesses HDFS metrics over `Java Management Extensions` (JMX) through the web interface of an HDFS daemon.

## Requirements

- `hdfs` node with accessible `/jmx` endpoint

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/hdfs/metrics.csv) for a list of
metrics.

## Configuration

Edit the `go.d/hdfs.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

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


