# Pika collector

## Overview

[Pika](https://github.com/Qihoo360/pika#introduction%E4%B8%AD%E6%96%87) is a persistent huge storage service,
compatible with the vast majority of redis interfaces (details), including string, hash, list, zset, set and management
interfaces.

This collector monitors one or more Pika instances, depending on your configuration.

It collects information and statistics about the server executing the following commands:

- [`INFO ALL`](https://github.com/Qihoo360/pika/wiki/pika-info%E4%BF%A1%E6%81%AF%E8%AF%B4%E6%98%8E)

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                             |        Dimensions        |    Unit     |
|------------------------------------|:------------------------:|:-----------:|
| pika.connections                   |         accepted         | connections |
| pika.clients                       |        connected         |   clients   |
| pika.memory                        |           used           |    bytes    |
| pika.connected_replicas            |        connected         |  replicas   |
| pika.commands                      |        processed         | commands/s  |
| pika.commands_calls                | a dimension per command  |   calls/s   |
| pika.database_strings_keys         | a dimension per database |    keys     |
| pika.database_strings_expires_keys | a dimension per database |    keys     |
| pika.database_strings_invalid_keys | a dimension per database |    keys     |
| pika.database_hashes_keys          | a dimension per database |    keys     |
| pika.database_hashes_expires_keys  | a dimension per database |    keys     |
| pika.database_hashes_invalid_keys  | a dimension per database |    keys     |
| pika.database_lists_keys           | a dimension per database |    keys     |
| pika.database_lists_expires_keys   | a dimension per database |    keys     |
| pika.database_lists_invalid_keys   | a dimension per database |    keys     |
| pika.database_zsets_keys           | a dimension per database |    keys     |
| pika.database_zsets_expires_keys   | a dimension per database |    keys     |
| pika.database_zsets_invalid_keys   | a dimension per database |    keys     |
| pika.database_sets_keys            | a dimension per database |    keys     |
| pika.database_sets_expires_keys    | a dimension per database |    keys     |
| pika.database_sets_invalid_keys    | a dimension per database |    keys     |
| pika.uptime                        |          uptime          |   seconds   |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/pika.conf`.

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
sudo ./edit-config go.d/pika.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                               |         Default         | Required |
|:-------------------:|-----------------------------------------------------------------------------------------------------------|:-----------------------:|:--------:|
|    update_every     | Data collection frequency.                                                                                |            5            |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                                        |            0            |          |
|       address       | Pika server address.                                                                                      | redis://@localhost:9221 |   yes    |
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
    address: 'redis://@localhost:9221'
```

</details>

##### TCP socket with password

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 'redis://:password@127.0.0.1:9221'
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
    address: 'redis://:password@127.0.0.1:9221'

  - name: remote
    address: 'redis://user:password@203.0.113.0:9221'
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `pika` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m pika
  ```
