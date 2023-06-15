# ActiveMQ collector

## Overview

[ActiveMQ](https://activemq.apache.org/) is an open source message broker written in Java together with a full Java
Message Service client.

This collector monitors queues and topics metrics using ActiveMQ Console API.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                        |     Dimensions     |    Unit    |
|-------------------------------|:------------------:|:----------:|
| activemq.messages             | enqueued, dequeued | messages/s |
| activemq.unprocessed_messages |    unprocessed     |  messages  |
| activemq.consumers            |     consumers      | consumers  |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/activemq.conf`.

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
sudo ./edit-config go.d/activemq.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                                                           |        Default        | Required |
|:--------------------:|---------------------------------------------------------------------------------------------------------------------------------------|:---------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                                            |           1           |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                                                    |           0           |          |
|         url          | Server URL.                                                                                                                           | http://localhost:8161 |   yes    |
|       webadmin       | Webadmin root path.                                                                                                                   |         admin         |   yes    |
|      max_queues      | Maximum number of concurrently collected queues.                                                                                      |          50           |          |
|      max_topics      | Maximum number of concurrently collected topics.                                                                                      |          50           |          |
|    queues_filter     | Queues filter. Syntax is [simple patterns](https://github.com/netdata/netdata/tree/master/libnetdata/simple_pattern#simple-patterns). |                       |          |
|    topics_filter     | Topics filter. Syntax is [simple patterns](https://github.com/netdata/netdata/tree/master/libnetdata/simple_pattern#simple-patterns). |                       |          |
|       username       | Username for basic HTTP authentication.                                                                                               |                       |          |
|       password       | Password for basic HTTP authentication.                                                                                               |                       |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                                                         |                       |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                                                         |                       |          |
|        method        | HTTP request method.                                                                                                                  |          GET          |          |
|       timeout        | HTTP request timeout.                                                                                                                 |           1           |          |
|         body         | HTTP request body.                                                                                                                    |                       |          |
|       headers        | HTTP request headers.                                                                                                                 |                       |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                                              |          no           |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check.                             |          no           |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                                                |                       |          |
|       tls_cert       | Client TLS certificate.                                                                                                               |                       |          |
|       tls_key        | Client TLS key.                                                                                                                       |                       |          |

</details>

#### Examples

##### Basic

A basic example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8161
    webadmin: admin
```

</details>

##### HTTP authentication

Basic HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8161
    webadmin: admin
    username: foo
    password: bar
```

</details>

##### Filters and limits

Using filters and limits for queues and topics.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8161
    webadmin: admin
    max_queues: 100
    max_topics: 100
    queues_filter: '!sandr* *'
    topics_filter: '!sandr* *'
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
    url: http://127.0.0.1:8161
    webadmin: admin

  - name: remote
    url: http://192.0.2.1:8161
    webadmin: admin
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `activemq` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m activemq
  ```
