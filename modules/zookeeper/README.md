# Zookeeper collector

## Overview

[ZooKeeper](https://zookeeper.apache.org/) is a centralized service for maintaining configuration information, naming,
providing distributed synchronization, and providing group services.

This collector monitors one or more ZooKeeper servers, depending on your configuration. It fetches metrics from
ZooKeeper by using the [mntr](https://zookeeper.apache.org/doc/r3.4.8/zookeeperAdmin.html#sc_zkCommands) command.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                          |    Dimensions     |       Unit       |
|---------------------------------|:-----------------:|:----------------:|
| zookeeper.requests              |    outstanding    |     requests     |
| zookeeper.requests_latency      |   min, avg, max   |        ms        |
| zookeeper.connections           |       alive       |   connections    |
| zookeeper.packets               |  received, sent   |       pps        |
| zookeeper.file_descriptor       |       open        | file descriptors |
| zookeeper.nodes                 | znode, ephemerals |      nodes       |
| zookeeper.watches               |      watches      |     watches      |
| zookeeper.approximate_data_size |       size        |       KiB        |
| zookeeper.server_state          |       state       |      state       |

## Setup

### Prerequisites

#### Whitelist `mntr` command

Add `mntr` to Zookeeper's [4lw.commands.whitelist](https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_4lw).

### Configuration

#### File

The configuration file name is `go.d/zookeeper.conf`.

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
sudo ./edit-config go.d/zookeeper.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                               |    Default     | Required |
|:-------------------:|-----------------------------------------------------------------------------------------------------------|:--------------:|:--------:|
|    update_every     | Data collection frequency.                                                                                |       1        |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                                        |       0        |          |
|       address       | Server address. The format is IP:PORT.                                                                    | 127.0.0.1:2181 |   yes    |
|       timeout       | Connection/read/write/ssl handshake timeout.                                                              |       1        |          |
|       use_tls       | Whether to use TLS or not.                                                                                |       no       |          |
|   tls_skip_verify   | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |       no       |          |
|       tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                |          |
|      tls_cert       | Client TLS certificate.                                                                                   |                |          |
|       tls_key       | Client TLS key.                                                                                           |                |          |

</details>

#### Examples

##### Basic

Local server.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:2181
```

</details>

##### TLS with self-signed certificate

Zookeeper with TLS and self-signed certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:2181
    use_tls: yes
    tls_skip_verify: yes
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Collecting metrics from local and remote instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:2181

  - name: remote
    address: 192.0.2.1:2181
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `zookeeper` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m zookeeper
  ```

