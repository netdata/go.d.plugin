<!--
title: "Redis monitoring with Netdata"
description: "Monitor the health and performance of Redis storage services with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/redis/README.md"
sidebar_label: "Redis"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# Redis collector

[`Redis`](https://redis.io/) is an open source (BSD licensed), in-memory data structure store, used as a database, cache
and message broker.

---

This module monitors one or more `Redis` instances, depending on your configuration.

It collects information and statistics about the server executing the following commands:

- [`INFO ALL`](https://redis.io/commands/info)

## Metrics

All metrics have "redis." prefix.

| Metric                          | Scope  |                   Dimensions                   |     Units      |
|---------------------------------|:------:|:----------------------------------------------:|:--------------:|
| connections                     | global |               accepted, rejected               | connections/s  |
| clients                         | global | connected, blocked, tracking, in_timeout_table |    clients     |
| ping_latency                    | global |                 min, max, avg                  |    seconds     |
| commands                        | global |                   processes                    |   commands/s   |
| keyspace_lookup_hit_rate        | global |                lookup_hit_rate                 |   percentage   |
| memory                          | global |  max, used, rss, peak, dataset, lua, scripts   |     bytes      |
| mem_fragmentation_ratio         | global |               mem_fragmentation                |     ratio      |
| key_eviction_events             | global |                    evicted                     |     keys/s     |
| net                             | global |                 received, sent                 |   kilobits/s   |
| rdb_changes                     | global |                    changes                     |   operations   |
| bgsave_now                      | global |              current_bgsave_time               |    seconds     |
| bgsave_health                   | global |                  last_bgsave                   |     status     |
| bgsave_last_rdb_save_since_time | global |                last_bgsave_time                |    seconds     |
| aof_file_size                   | global |                 current, base                  |     bytes      |
| commands_calls                  | global |         <i>a dimension per command</i>         |     calls      |
| commands_usec                   | global |         <i>a dimension per command</i>         |  microseconds  |
| commands_usec_per_sec           | global |         <i>a dimension per command</i>         | microseconds/s |
| key_expiration_events           | global |                    expired                     |     keys/s     |
| database_keys                   | global |        <i>a dimension per database</i>         |      keys      |
| database_expires_keys           | global |        <i>a dimension per database</i>         |      keys      |
| connected_replicas              | global |                   connected                    |    replicas    |
| master_link_status              | global |                    up, down                    |     status     |
| master_last_io_since_time       | global |                      time                      |    seconds     |
| master_link_down_since_time     | global |                      time                      |    seconds     |
| uptime                          | global |                     uptime                     |    seconds     |

## Configuration

Edit the `go.d/redis.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

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

  - name: local
    address: 'redis://127.0.0.1:6379'
    username: 'user'
    password: 'password'

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
