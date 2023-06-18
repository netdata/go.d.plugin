# RabbitMQ collector

## Overview

[RabbitMQ](https://www.rabbitmq.com/) is an open-source message broker.

This collector monitors one or more RabbitMQ instances, depending on your configuration.

It collects data using an HTTP-based API provided by the [management plugin](https://www.rabbitmq.com/management.html).
The following endpoints are used:

- `/api/overview`
- `/api/node/{node_name}`
- `/api/vhosts`
- `/api/queues` (disabled by default)

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                    |                                                             Dimensions                                                              |     Unit     |
|-------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------:|:------------:|
| rabbitmq.messages_count                   |                                                        ready, unacknowledged                                                        |   messages   |
| rabbitmq.messages_rate                    | ack, publish, publish_in, publish_out, confirm, deliver, deliver_no_ack, get, get_no_ack, deliver_get, redeliver, return_unroutable |  messages/s  |
| rabbitmq.objects_count                    |                                         channels, consumers, connections, queues, exchanges                                         |   messages   |
| rabbitmq.connection_churn_rate            |                                                           created, closed                                                           | operations/s |
| rabbitmq.channel_churn_rate               |                                                           created, closed                                                           | operations/s |
| rabbitmq.queue_churn_rate                 |                                                     created, deleted, declared                                                      | operations/s |
| rabbitmq.file_descriptors_count           |                                                           available, used                                                           |      fd      |
| rabbitmq.sockets_count                    |                                                           available, used                                                           |   sockets    |
| rabbitmq.erlang_processes_count           |                                                           available, used                                                           |  processes   |
| rabbitmq.erlang_run_queue_processes_count |                                                               length                                                                |  processes   |
| rabbitmq.memory_usage                     |                                                                used                                                                 |    bytes     |
| rabbitmq.disk_space_free_size             |                                                                free                                                                 |    bytes     |

### vhost

These metrics refer to the virtual host.

Labels:

| Label | Description       |
|-------|-------------------|
| vhost | virtual host name |

Metrics:

| Metric                        |                                                             Dimensions                                                              |    Unit    |
|-------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------:|:----------:|
| rabbitmq.vhost_messages_count |                                                        ready, unacknowledged                                                        |  messages  |
| rabbitmq.vhost_messages_rate  | ack, publish, publish_in, publish_out, confirm, deliver, deliver_no_ack, get, get_no_ack, deliver_get, redeliver, return_unroutable | messages/s |

### queue

These metrics refer to the virtual host queue.

Labels:

| Label | Description       |
|-------|-------------------|
| vhost | virtual host name |
| queue | queue name        |

Metrics:

| Metric                        |                                                             Dimensions                                                              |    Unit    |
|-------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------:|:----------:|
| rabbitmq.queue_messages_count |                                            ready, unacknowledged, paged_out, persistent                                             |  messages  |
| rabbitmq.queue_messages_rate  | ack, publish, publish_in, publish_out, confirm, deliver, deliver_no_ack, get, get_no_ack, deliver_get, redeliver, return_unroutable | messages/s |

## Setup

### Prerequisites

#### Enable management plugin.

The management plugin is included in the RabbitMQ distribution, but disabled.
To enable see [Management Plugin](https://www.rabbitmq.com/management.html#getting-started) documentation.

### Configuration

#### File

The configuration file name is `go.d/rabbitmq.conf`.

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
sudo ./edit-config go.d/rabbitmq.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|          Name          | Description                                                                                                                                           |        Default         | Required |
|:----------------------:|-------------------------------------------------------------------------------------------------------------------------------------------------------|:----------------------:|:--------:|
|      update_every      | Data collection frequency.                                                                                                                            |           1            |          |
|  autodetection_retry   | Re-check interval in seconds. Zero means not to schedule re-check.                                                                                    |           0            |          |
|          url           | Server URL.                                                                                                                                           | http://localhost:15672 |   yes    |
| collect_queues_metrics | Collect stats per vhost per queues. Enabling this can introduce serious overhead on both Netdata and RabbitMQ if many queues are configured and used. |           no           |          |
|        timeout         | HTTP request timeout.                                                                                                                                 |           1            |          |
|        username        | Username for basic HTTP authentication.                                                                                                               |                        |          |
|        password        | Password for basic HTTP authentication.                                                                                                               |                        |          |
|       proxy_url        | Proxy URL.                                                                                                                                            |                        |          |
|     proxy_username     | Username for proxy basic HTTP authentication.                                                                                                         |                        |          |
|     proxy_password     | Password for proxy basic HTTP authentication.                                                                                                         |                        |          |
|         method         | HTTP request method.                                                                                                                                  |          GET           |          |
|          body          | HTTP request body.                                                                                                                                    |                        |          |
|        headers         | HTTP request headers.                                                                                                                                 |                        |          |
|  not_follow_redirects  | Redirect handling policy. Controls whether the client follows redirects.                                                                              |           no           |          |
|    tls_skip_verify     | Server certificate chain and hostname validation policy. Controls whether the client performs this check.                                             |           no           |          |
|         tls_ca         | Certification authority that the client uses when verifying the server's certificates.                                                                |                        |          |
|        tls_cert        | Client TLS certificate.                                                                                                                               |                        |          |
|        tls_key         | Client TLS key.                                                                                                                                       |                        |          |

</details>

#### Examples

##### Basic

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:15672
```

</details>

##### Basic HTTP auth

Local server with basic HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:15672
    username: admin
    password: password
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
    url: http://127.0.0.1:15672

  - name: remote
    url: http://192.0.2.0:15672
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `rabbitmq` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m rabbitmq
  ```
