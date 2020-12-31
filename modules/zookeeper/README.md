<!--
title: "Apache ZooKeeper monitoring with Netdata"
description: "Monitor the health and performance of Zookeeper servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/zookeeper/README.md
sidebar_label: "ZooKeeper"
-->

# Apache ZooKeeper monitoring with Netdata

[`ZooKeeper`](https://zookeeper.apache.org/) is a centralized service for maintaining configuration information, naming,
providing distributed synchronization, and providing group services.

This module monitors one or more `ZooKeeper` servers, depending on your configuration.

## Requirements

- `Zookeeper` with accessible client port
- whitelisted `mntr` command

## Charts

It produces the following charts:

- Outstanding Requests in `requests`
- Requests Latency in `ms`
- Alive Connections in `connections`
- Packets in `pps`
- Open File Descriptors in `file descriptors`
- Number of Nodes in `nodes`
- Number of Watches in `watches`
- Approximate Data Tree Size in `KiB`
- Server State in `state`

## Configuration

Edit the `go.d/zookeeper.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/zookeeper.conf
```

Needs only `address` to server's client port. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    address: 127.0.0.1:2181

  - name: remote
    address: 203.0.113.10:2182
```

For all available options, please see the
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/zookeeper.conf).

## Troubleshooting

To troubleshoot issues with the `zookeeper` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m zookeeper
```
