<!--
title: "VerneMQ monitoring with Netdata"
description: "Monitor the health and performance of VerneMQ MQTT brokers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/vernemq/README.md"
sidebar_label: "VerneMQ"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Message brokers"
-->

# VerneMQ collector

[`VerneMQ`](https://vernemq.com/) is a scalable and open source MQTT broker that connects IoT, M2M, Mobile, and web
applications.

This module will monitor one or more `VerneMQ` instances, depending on your configuration.

`vernemq` module is tested on the following versions:

- v1.10.1

## Metrics

All metrics have "vernemq." prefix.

| Metric                          | Scope  |            Dimensions            |       Units       |
|---------------------------------|:------:|:--------------------------------:|:-----------------:|
| sockets                         | global |           open, close            |     events/s      |
| client_keepalive_expired        | global |              closed              |     sockets/s     |
| socket_close_timeout            | global |              closed              |     sockets/s     |
| socket_errors                   | global |              errors              |     errors/s      |
| queue_processes                 | global |         queue_processes          |  queue processes  |
| queue_processes_operations      | global |         setup, teardown          |     events/s      |
| queue_process_init_from_storage | global |         queue_processes          | queue processes/s |
| queue_messages                  | global |          received, sent          |    messages/s     |
| queue_undelivered_messages      | global |   dropped, expired, unhandled    |    messages/s     |
| router_subscriptions            | global |          subscriptions           |   subscriptions   |
| router_matched_subscriptions    | global |          local, remote           |  subscriptions/s  |
| router_memory                   | global |               used               |        KiB        |
| average_scheduler_utilization   | global |           utilization            |    percentage     |
| system_utilization_scheduler    | global | <i>a dimension per scheduler</i> |    percentage     |
| system_processes                | global |            processes             |     processes     |
| system_reductions               | global |            reductions            |       ops/s       |
| system_context_switches         | global |         context_switches         |       ops/s       |
| system_io                       | global |          received, sent          |    kilobits/s     |
| system_run_queue                | global |              ready               |     processes     |
| system_gc_count                 | global |                gc                |       ops/s       |
| system_gc_words_reclaimed       | global |         words_reclaimed          |       ops/s       |
| system_allocated_memory         | global |        processes, system         |        KiB        |
| bandwidth                       | global |          received, sent          |    kilobits/s     |
| retain_messages                 | global |             messages             |     messages      |
| retain_memory                   | global |               used               |        KiB        |
| cluster_bandwidth               | global |          received, sent          |    kilobits/s     |
| cluster_dropped                 | global |             dropped              |    kilobits/s     |
| netsplit_unresolved             | global |            unresolved            |     netsplits     |
| netsplits                       | global |        resolved, detected        |    netsplits/s    |
| mqtt_auth                       | global |          received, sent          |     packets/s     |
| mqtt_auth_received_reason       | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_auth_sent_reason           | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_connect                    | global |         connect, connack         |     packets/s     |
| mqtt_connack_sent_reason        | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_disconnect                 | global |          received, sent          |     packets/s     |
| mqtt_disconnect_received_reason | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_disconnect_sent_reason     | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_subscribe                  | global |        subscribe, suback         |     packets/s     |
| mqtt_subscribe_error            | global |              failed              |       ops/s       |
| mqtt_subscribe_auth_error       | global |              unauth              |    attempts/s     |
| mqtt_unsubscribe                | global |      unsubscribe, unsuback       |     packets/s     |
| mqtt_unsubscribe                | global |      mqtt_unsubscribe_error      |       ops/s       |
| mqtt_publish                    | global |          received, sent          |     packets/s     |
| mqtt_publish_errors             | global |              failed              |       ops/s       |
| mqtt_publish_auth_errors        | global |              unauth              |    attempts/s     |
| mqtt_puback                     | global |          received, sent          |     packets/s     |
| mqtt_puback_received_reason     | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_puback_sent_reason         | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_puback_invalid_error       | global |            unexpected            |    messages/s     |
| mqtt_pubrec                     | global |          received, sent          |     packets/s     |
| mqtt_pubrec_received_reason     | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_pubrec_sent_reason         | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_pubrec_invalid_error       | global |            unexpected            |    messages/s     |
| mqtt_pubrel                     | global |          received, sent          |     packets/s     |
| mqtt_pubrel_received_reason     | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_pubrel_sent_reason         | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_pubcom                     | global |          received, sent          |     packets/s     |
| mqtt_pubcomp_received_reason    | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_pubcomp_sent_reason        | global |  <i>a dimensions per reason</i>  |     packets/s     |
| mqtt_pubcomp_invalid_error      | global |            unexpected            |    messages/s     |
| mqtt_ping                       | global |        pingreq, pingresp         |     packets/s     |
| node_uptime                     | global |        pingreq, pingresp         |     packets/s     |

## Configuration

Edit the `go.d/vernemq.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/vernemq.conf
```

Needs only `url` to server's `/metrics` endpoint. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8888/metrics

  - name: remote
    url: http://203.0.113.10:8888/metrics
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/vernemq.conf).

## Troubleshooting

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
