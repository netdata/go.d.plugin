<!--
title: "VerneMQ monitoring with Netdata"
description: "Monitor the health and performance of VerneMQ MQTT brokers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/vernemq/README.md
sidebar_label: "VerneMQ"
-->

# VerneMQ monitoring with Netdata

[`VerneMQ`](https://vernemq.com/) is a scalable and open source MQTT broker that connects IoT, M2M, Mobile, and web
applications.

This module will monitor one or more `VerneMQ` instances, depending on your configuration.

`vernemq` module is tested on the following versions:

- v1.10.1

## Charts

It produces the following charts:

#### Sockets

- Open Sockets in `sockets`
- Socket Open and Close Events in `sockets/s`
- Closed Sockets due to Keepalive Time Expired in `sockets/s`
- Closed Sockets due to no CONNECT Frame On Time in `sockets/s`
- Socket Errors in `errors/s`

#### Queues

- Living Queues in an Online or an Offline State in `queue processes`
- Queue Processes Setup and Teardown Events in `events/s`
- Queue Processes Initialized from Offline Storage in `queue processes/s`
- Received and Sent PUBLISH Messages in `messages/s`
- Undelivered PUBLISH Messages in `messages`

#### Subscriptions

- Subscriptions in the Routing Table in `subscriptions`
- Matched Subscriptions `subscriptions`
- Routing Table Memory Usage in `KiB`

#### Erlang VM

- Average Scheduler Utilization in `percentage`
- Erlang Processes in `processes`
- Reductions in `ops/s`
- Context Switches in `ops/s`
- Received and Sent Traffic through Ports in `kilobits/s`
- Processes that are Ready to Run on All Run-Queues in `processes`
- GC Count in `KiB/s`
- GC Words Reclaimed in `KiB/s`
- Memory Allocated by the Erlang Processes and by the Emulator in `KiB`

#### Bandwidth

- Bandwidth in `kilobits/s`

#### Retain

- Stored Retained Messages in `messages`
- Stored Retained Messages Memory Usage in `KiB`

#### Cluster

- Communication with Other Cluster Nodes in `kilobits/s`
- Traffic Dropped During Communication with Other Cluster Nodes in `kilobits/s`
- Unresolved Netsplits in `netsplits`
- Netsplits in `netsplits/s`

#### MQTT AUTH

- v5 AUTH in `packets/s`
- v5 AUTH Received by Reason in `packets/s`
- v5 AUTH Sent by Reason in `packets/s`

#### MQTT CONNECT

- v3/v5 CONNECT and CONNACK in `packets/s`
- v3/v5 CONNACK Sent by Reason in `packets/s`

#### MQTT DISCONNECT

- v3/v5 DISCONNECT in `packets/s`
- v5 DISCONNECT Received by Reason in `packets/s`
- v5 DISCONNECT Sent by Reason in `packets/s`

#### MQTT SUBSCRIBE

- v3/v5 SUBSCRIBE and SUBACK in `packets/s`
- v3/v5 Failed SUBSCRIBE Operations due to a Netsplit in `ops/s`
- v3/v5 Unauthorized SUBSCRIBE Attempts in `attempts/s`

#### MQTT UNSUBSCRIBE

- v3/v5 UNSUBSCRIBE and UNSUBACK in `packets/s`
- v3/v5 Failed UNSUBSCRIBE Operation due to a Netsplit in `ops/s`

#### MQTT PUBLISH

- v3/v5 QoS 0,1,2 PUBLISH in `packets/s`
- v3/v5 Failed PUBLISH Operations due to a Netsplit in `ops/s`
- v3/v5 Unauthorized PUBLISH Attempts in `attempts/s`
- v3/v5 QoS 1 PUBACK in `packets/s`
- v5 PUBACK QoS 1 Received by Reason in `packets/s`
- v5 PUBACK QoS 1 Sent by Reason in `packets/s`
- v3/v5 PUBACK QoS 1 Received Unexpected Messages in `messages/s`
- v3/v5 PUBREC QoS 2 in `packets/s`
- v5 PUBREC QoS 2 Received by Reason in `packets/s`
- v5 PUBREC QoS 2 Sent by Reason in `packets/s`
- v3 PUBREC QoS 2 Received Unexpected Messages in `messages/s`
- v3/v5 PUBREL QoS 2 in `packets/s`
- v5 PUBREL QoS 2 Received by Reason in `packets/s`
- v5 PUBREL QoS 2 Sent by Reason in `packets/s`
- v3/v5 PUBCOMP QoS 2 in `packets/s`
- v5 PUBCOMP QoS 2 Received by Reason in `packets/s`
- v5 PUBCOMP QoS 2 Sent by Reason in `packets/s`
- v3/v5 PUBCOMP QoS 2 Received Unexpected Messages in `messages/s`

#### MQTT PING

- v3/v5 PING in `packets/s`

#### Uptime

- Node Uptime in `seconds`

## Configuration

Edit the `go.d/vernemq.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

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
