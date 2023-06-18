# VerneMQ collector

## Overview

[VerneMQ](https://vernemq.com) is a high-performance, distributed MQTT broker.

This collector monitors one or more VerneMQ instances, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                  |         Dimensions          |       Unit        |
|-----------------------------------------|:---------------------------:|:-----------------:|
| vernemq.sockets                         |            open             |      sockets      |
| vernemq.socket_operations               |         open, close         |     sockets/s     |
| vernemq.client_keepalive_expired        |           closed            |     sockets/s     |
| vernemq.socket_close_timeout            |           closed            |     sockets/s     |
| vernemq.socket_errors                   |           errors            |     errors/s      |
| vernemq.queue_processes                 |       queue_processes       |  queue processes  |
| vernemq.queue_processes_operations      |       setup, teardown       |     events/s      |
| vernemq.queue_process_init_from_storage |       queue_processes       | queue processes/s |
| vernemq.queue_messages                  |       received, sent        |    messages/s     |
| vernemq.queue_undelivered_messages      | dropped, expired, unhandled |    messages/s     |
| vernemq.router_subscriptions            |        subscriptions        |   subscriptions   |
| vernemq.router_matched_subscriptions    |        local, remote        |  subscriptions/s  |
| vernemq.router_memory                   |            used             |        KiB        |
| vernemq.average_scheduler_utilization   |         utilization         |    percentage     |
| vernemq.system_utilization_scheduler    |  a dimension per scheduler  |    percentage     |
| vernemq.system_processes                |          processes          |     processes     |
| vernemq.system_reductions               |         reductions          |       ops/s       |
| vernemq.system_context_switches         |      context_switches       |       ops/s       |
| vernemq.system_io                       |       received, sent        |    kilobits/s     |
| vernemq.system_run_queue                |            ready            |     processes     |
| vernemq.system_gc_count                 |             gc              |       ops/s       |
| vernemq.system_gc_words_reclaimed       |       words_reclaimed       |       ops/s       |
| vernemq.system_allocated_memory         |      processes, system      |        KiB        |
| vernemq.bandwidth                       |       received, sent        |    kilobits/s     |
| vernemq.retain_messages                 |          messages           |     messages      |
| vernemq.retain_memory                   |            used             |        KiB        |
| vernemq.cluster_bandwidth               |       received, sent        |    kilobits/s     |
| vernemq.cluster_dropped                 |           dropped           |    kilobits/s     |
| vernemq.netsplit_unresolved             |         unresolved          |     netsplits     |
| vernemq.netsplits                       |     resolved, detected      |    netsplits/s    |
| vernemq.mqtt_auth                       |       received, sent        |     packets/s     |
| vernemq.mqtt_auth_received_reason       |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_auth_sent_reason           |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_connect                    |      connect, connack       |     packets/s     |
| vernemq.mqtt_connack_sent_reason        |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_disconnect                 |       received, sent        |     packets/s     |
| vernemq.mqtt_disconnect_received_reason |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_disconnect_sent_reason     |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_subscribe                  |      subscribe, suback      |     packets/s     |
| vernemq.mqtt_subscribe_error            |           failed            |       ops/s       |
| vernemq.mqtt_subscribe_auth_error       |           unauth            |    attempts/s     |
| vernemq.mqtt_unsubscribe                |    unsubscribe, unsuback    |     packets/s     |
| vernemq.mqtt_unsubscribe                |   mqtt_unsubscribe_error    |       ops/s       |
| vernemq.mqtt_publish                    |       received, sent        |     packets/s     |
| vernemq.mqtt_publish_errors             |           failed            |       ops/s       |
| vernemq.mqtt_publish_auth_errors        |           unauth            |    attempts/s     |
| vernemq.mqtt_puback                     |       received, sent        |     packets/s     |
| vernemq.mqtt_puback_received_reason     |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_puback_sent_reason         |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_puback_invalid_error       |         unexpected          |    messages/s     |
| vernemq.mqtt_pubrec                     |       received, sent        |     packets/s     |
| vernemq.mqtt_pubrec_received_reason     |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_pubrec_sent_reason         |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_pubrec_invalid_error       |         unexpected          |    messages/s     |
| vernemq.mqtt_pubrel                     |       received, sent        |     packets/s     |
| vernemq.mqtt_pubrel_received_reason     |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_pubrel_sent_reason         |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_pubcom                     |       received, sent        |     packets/s     |
| vernemq.mqtt_pubcomp_received_reason    |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_pubcomp_sent_reason        |   a dimensions per reason   |     packets/s     |
| vernemq.mqtt_pubcomp_invalid_error      |         unexpected          |    messages/s     |
| vernemq.mqtt_ping                       |      pingreq, pingresp      |     packets/s     |
| vernemq.node_uptime                     |            time             |      seconds      |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/vernemq.conf`.

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
sudo ./edit-config go.d/vernemq.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |            Default            | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:-----------------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |               1               |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |               0               |          |
|         url          | Server URL.                                                                                               | http://127.0.0.1:8888/metrics |   yes    |
|       timeout        | HTTP request timeout.                                                                                     |               1               |          |
|       username       | Username for basic HTTP authentication.                                                                   |                               |          |
|       password       | Password for basic HTTP authentication.                                                                   |                               |          |
|      proxy_url       | Proxy URL.                                                                                                |                               |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                               |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                               |          |
|        method        | HTTP request method.                                                                                      |              GET              |          |
|         body         | HTTP request body.                                                                                        |                               |          |
|       headers        | HTTP request headers.                                                                                     |                               |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |              no               |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |              no               |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                               |          |
|       tls_cert       | Client TLS certificate.                                                                                   |                               |          |
|       tls_key        | Client TLS key.                                                                                           |                               |          |

</details>

#### Examples

##### Basic

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8888/metrics
```

</details>

##### HTTP authentication

Local instance with basic HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8888/metrics
    username: username
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
    url: http://127.0.0.1:8888/metrics

  - name: remote
    url: http://203.0.113.10:8888/metrics
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `vernemq` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m vernemq
  ```
