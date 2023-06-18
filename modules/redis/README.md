# Redis collector

## Overview

[Redis](https://redis.io/) is an open source (BSD licensed), in-memory data structure store, used as a database, cache
and message broker.

This collector monitors one or more Redis instances, depending on your configuration.

It collects information and statistics about the server executing the following commands:

- [`INFO ALL`](https://redis.io/commands/info)

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                |                   Dimensions                   |      Unit      |
|---------------------------------------|:----------------------------------------------:|:--------------:|
| redis.connections                     |               accepted, rejected               | connections/s  |
| redis.clients                         | connected, blocked, tracking, in_timeout_table |    clients     |
| redis.ping_latency                    |                 min, max, avg                  |    seconds     |
| redis.commands                        |                   processes                    |   commands/s   |
| redis.keyspace_lookup_hit_rate        |                lookup_hit_rate                 |   percentage   |
| redis.memory                          |  max, used, rss, peak, dataset, lua, scripts   |     bytes      |
| redis.mem_fragmentation_ratio         |               mem_fragmentation                |     ratio      |
| redis.key_eviction_events             |                    evicted                     |     keys/s     |
| redis.net                             |                 received, sent                 |   kilobits/s   |
| redis.rdb_changes                     |                    changes                     |   operations   |
| redis.bgsave_now                      |              current_bgsave_time               |    seconds     |
| redis.bgsave_health                   |                  last_bgsave                   |     status     |
| redis.bgsave_last_rdb_save_since_time |                last_bgsave_time                |    seconds     |
| redis.aof_file_size                   |                 current, base                  |     bytes      |
| redis.commands_calls                  |            a dimension per command             |     calls      |
| redis.commands_usec                   |            a dimension per command             |  microseconds  |
| redis.commands_usec_per_sec           |            a dimension per command             | microseconds/s |
| redis.key_expiration_events           |                    expired                     |     keys/s     |
| redis.database_keys                   |            a dimension per database            |      keys      |
| redis.database_expires_keys           |            a dimension per database            |      keys      |
| redis.connected_replicas              |                   connected                    |    replicas    |
| redis.master_link_status              |                    up, down                    |     status     |
| redis.master_last_io_since_time       |                      time                      |    seconds     |
| redis.master_link_down_since_time     |                      time                      |    seconds     |
| redis.uptime                          |                     uptime                     |    seconds     |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/redis.conf`.

The file format is YAML. Generally, the format is:

```yaml
update_every: 1
autodetection_retry: 0
jobs:
  - name: some_name1
  - name: some_name1
```

You can edit the configuration file using the `edit-config` script from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md#the-netdata-config-directory).

```bash
cd /etc/netdata 2>/dev/null || cd /opt/netdata/etc/netdata
sudo ./edit-config go.d/redis.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                               |         Default         | Required |
|:-------------------:|-----------------------------------------------------------------------------------------------------------|:-----------------------:|:--------:|
|    update_every     | Data collection frequency.                                                                                |            5            |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                                        |            0            |          |
|       address       | Redis server address.                                                                                     | redis://@localhost:6379 |   yes    |
|       timeout       | Dial (establishing new connections), read (socket reads) and write (socket writes) timeout in seconds.    |            1            |          |
|      username       | Username used for authentication.                                                                         |                         |          |
|      password       | Password used for authentication.                                                                         |                         |          |
|   tls_skip_verify   | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |           no            |          |
|       tls_ca        | Certificate authority that client use when verifying server certificates.                                 |                         |          |
|      tls_cert       | Client tls certificate.                                                                                   |                         |          |
|       tls_key       | Client tls key.                                                                                           |                         |          |

</details>

##### address

There are two connection types: by tcp socket and by unix socket.

- Tcp connection: `redis://<user>:<password>@<host>:<port>/<db_number>`
- Unix connection: `unix://<user>:<password>@</path/to/redis.sock>?db=<db_number>`

#### Examples

##### TCP socket

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 'redis://@127.0.0.1:6379'
```

</details>

##### Unix socket

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 'unix://@/tmp/redis.sock'
```

</details>

##### TCP socket with password

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 'redis://:password@127.0.0.1:6379'
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Local and remote instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 'redis://:password@127.0.0.1:6379'

  - name: remote
    address: 'redis://user:password@203.0.113.0:6379'
```

</details>

## Troubleshooting

### Debug mode

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
