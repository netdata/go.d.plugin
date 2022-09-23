<!--
title: "Redis monitoring with Netdata"
description: "Monitor the health and performance of Redis storage services with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/redis/README.md
sidebar_label: "Redis"
-->

# Redis monitoring with Netdata

[`Redis`](https://redis.io/) is an open source (BSD licensed), in-memory data structure store, used as a database, cache
and message broker.

---

This module monitors one or more `Redis` instances, depending on your configuration.

It collects information and statistics about the server executing the following commands:

- [`INFO ALL`](https://redis.io/commands/info)

## Charts

### Connections

- Accepted and rejected (maxclients limit) connections in `connections/s`
- Clients in `clients`

### Performance

- Ping latency in `seconds`
- Processed commands in `commands/s`
- Keys lookup hit rate in `percentage`

### Memory

- Memory usage in `bytes`
- Ratio between used_memory_rss and used_memory in `ratio`
- Evicted keys due to maxmemory limit in `keys/s`

### Network bandwidth

- Bandwidth in `kilobits/s`

### Replication

- Connected replicas in `replicas`

### Persistence

- Operations that produced changes since the last SAVE or BGSAVE in `operations`
- Duration of the ongoing RDB save operation if any in `seconds`
- Status of the last RDB save operation in `status`
- AOF file size in `bytes`

### Commands

- Calls per command in `calls/s`
- Total CPU time consumed by the commands in `usec`
- Average CPU consumed per command execution in `usec/s`

### Keyspace

- Expired keys in `keys/s`
- Keys per database in `keys`
- Keys with an expiration per database in `keys`

### Uptime

- Uptime in `seconds`

## Configuration

Edit the `go.d/redis.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/redis.conf
```

There are two connection types: by tcp socket and by unix socket.

> **Note**: If the Redis server is password protected via the `requirepass` option, make sure you have a colon before
> the password.

```cmd
# by tcp socket
redis://<user>:<password>@<host>:<port>

# by unix socket
unix://<user>:<password>@</path/to/redis.sock
```

Needs only `address`, here is an example with two jobs:

```yaml
jobs:
  - name: local
    address: 'redis://@127.0.0.1:6379'

  - name: local
    address: 'redis://:password@127.0.0.1:6379'

  - name: remote
    address: 'redis://user:password@203.0.113.0:6379'
```

For all available options, see the `redis`
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/redis.conf).

## Troubleshooting

To troubleshoot issues with the `redis` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m redis
  ```
