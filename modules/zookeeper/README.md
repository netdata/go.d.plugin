# ZooKeeper monitoring with Netdata

[`ZooKeeper`](https://zookeeper.apache.org/) is a centralized service for maintaining configuration information, naming, providing distributed synchronization, and providing group services. 

This module monitors one or more `ZooKeeper` servers, depending on your configuration.

## Requirements

-   `Zookeeper` with accessible client port
-   whitelisted `mntr` command

## Charts

It produces the following charts:

-   Outstanding Requests in `requests`
-   Requests Latency in `ms`
-   Alive Connections in `connections`
-   Packets in `pps`
-   Open File Descriptors in `file descriptors`
-   Number of Nodes in `nodes`
-   Number of Watches in `watches`
-   Approximate Data Tree Size in `KiB`
-   Server State in `state`

## Configuration

Edit the `go.d/zookeeper.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/zookeeper.conf
```

Needs only `address` to server's client port. Here is an example for 2 servers:

```yaml
jobs:
  - name    : local
    address : 127.0.0.1:2181
      
  - name    : remote
    address : 203.0.113.10:2182
```

For all available options, please see the module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/zookeeper.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m zookeeper
