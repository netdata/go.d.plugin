# zookeeper

This module monitors one or more [`Zookeeper`](https://zookeeper.apache.org/) servers, depending on your configuration.

**Requirements:**

-   `Zookeeper` with accessible client port
-   whitelisted `mntr` command

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

### Configuration

Needs only `address` to server's client port.

Here is an example for 2 servers:

```yaml
jobs:
  - name    : local
    address : 127.0.0.1:2181
      
  - name    : remote
    address : 203.0.113.10:2182
```

For all available options, please see the module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/zookeeper.conf).

---
