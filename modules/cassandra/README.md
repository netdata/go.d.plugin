# Cassandra collector

## Overview

[Cassandra](https://cassandra.apache.org/_/index.html) is an open-source NoSQL database management system.

This collector gathers metrics from one or more Cassandra servers, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                           |          Dimensions           |     Unit     |
|--------------------------------------------------|:-----------------------------:|:------------:|
| cassandra.client_requests_rate                   |          read, write          |  requests/s  |
| cassandra.client_request_read_latency_histogram  | p50, p75, p95, p98, p99, p999 |   seconds    |
| cassandra.client_request_write_latency_histogram | p50, p75, p95, p98, p99, p999 |   seconds    |
| cassandra.client_requests_latency                |          read, write          |   seconds    |
| cassandra.row_cache_hit_ratio                    |           hit_ratio           |  percentage  |
| cassandra.row_cache_hit_rate                     |         hits, misses          |   events/s   |
| cassandra.row_cache_utilization                  |             used              |  percentage  |
| cassandra.row_cache_size                         |             size              |    bytes     |
| cassandra.key_cache_hit_ratio                    |           hit_ratio           |  percentage  |
| cassandra.key_cache_hit_rate                     |         hits, misses          |   events/s   |
| cassandra.key_cache_utilization                  |             used              |  percentage  |
| cassandra.key_cache_size                         |             size              |    bytes     |
| cassandra.storage_live_disk_space_used           |             used              |    bytes     |
| cassandra.compaction_completed_tasks_rate        |           completed           |   tasks/s    |
| cassandra.compaction_pending_tasks_count         |            pending            |    tasks     |
| cassandra.compaction_compacted_rate              |           compacted           |   bytes/s    |
| cassandra.jvm_memory_used                        |         heap, nonheap         |    bytes     |
| cassandra.jvm_gc_rate                            |          parnew, cms          |     gc/s     |
| cassandra.jvm_gc_time                            |          parnew, cms          |   seconds    |
| cassandra.dropped_messages_rate                  |            dropped            |  messages/s  |
| cassandra.client_requests_timeouts_rate          |          read, write          |  timeout/s   |
| cassandra.client_requests_unavailables_rate      |          read, write          | exceptions/s |
| cassandra.client_requests_failures_rate          |          read, write          |  failures/s  |
| cassandra.storage_exceptions_rate                |            storage            | exceptions/s |

### thread pool

Metrics related to Cassandra's thread pools. Each thread pool provides its own set of the following metrics.

Labels:

| Label       | Description      |
|-------------|------------------|
| thread_pool | thread pool name |

Metrics:

| Metric                                    | Dimensions |  Unit   |
|-------------------------------------------|:----------:|:-------:|
| cassandra.thread_pool_active_tasks_count  |   active   |  tasks  |
| cassandra.thread_pool_pending_tasks_count |  pending   |  tasks  |
| cassandra.thread_pool_blocked_tasks_count |  blocked   |  tasks  |
| cassandra.thread_pool_blocked_tasks_rate  |  blocked   | tasks/s |

## Setup

### Prerequisites

#### Configure Cassandra with Prometheus JMX Exporter

To configure Cassandra with the [JMX Exporter](https://github.com/prometheus/jmx_exporter):

> **Note**: paths can differ depends on your setup.

- Download latest [jmx_exporter](https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/) jar file
  and install it in a directory where Cassandra can access it.
- Add
  the [jmx_exporter.yaml](https://raw.githubusercontent.com/netdata/go.d.plugin/master/modules/cassandra/jmx_exporter.yaml)
  file to `/etc/cassandra`.
- Add the following line to `/etc/cassandra/cassandra-env.sh`
  ```
  JVM_OPTS="$JVM_OPTS $JVM_EXTRA_OPTS -javaagent:/opt/jmx_exporter/jmx_exporter.jar=7072:/etc/cassandra/jmx_exporter.yaml
  ```
- Restart cassandra service.

### Configuration

#### File

The configuration file name is `go.d/cassandra.conf`.

The file format is YAML. Generally the format is:

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
sudo ./edit-config go.d/cassandra.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |            Default            | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:-----------------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |               5               |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |               0               |          |
|         url          | Server URL.                                                                                               | http://127.0.0.1:7072/metrics |   yes    |
|       username       | Username for basic HTTP authentication.                                                                   |                               |          |
|       password       | Password for basic HTTP authentication.                                                                   |                               |          |
|      proxy_url       | Proxy URL.                                                                                                |                               |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                               |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                               |          |
|       timeout        | HTTP request timeout.                                                                                     |               2               |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |              no               |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |              no               |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                               |          |
|       tls_cert       | Client TLS certificate.                                                                                   |                               |          |
|       tls_key        | Client TLS key.                                                                                           |                               |          |

</details>

#### Examples

##### Basic

A basic example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:7072/metrics
```

</details>

##### HTTP authentication

Local server with basic HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:7072/metrics
    username: foo
    password: bar
```

</details>

##### HTTPS with self-signed certificate

Local server with enabled HTTPS and self-signed certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: https://127.0.0.1:7072/metrics
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
    url: http://127.0.0.1:7072/metrics

  - name: remote
    url: http://192.0.2.1:7072/metrics
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `cassandra` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m cassandra
  ```

