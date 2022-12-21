<!--
title: "Consul monitoring with Netdata"
description: "Monitor the health and performance of Consul service meshes with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/consul/README.md"
sidebar_label: "Consul"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Webapps"
-->

# Consul monitoring with Netdata

[Consul](https://www.consul.io/) is a service networking solution to connect and secure services across any runtime
platform and public or private cloud.

This module collects the [Key Metrics](https://developer.hashicorp.com/consul/docs/agent/telemetry#key-metrics) of the
Consul Agent.

## Requirements

- Consul
  with [enabled](https://developer.hashicorp.com/consul/docs/agent/config/config-files#telemetry-prometheus_retention_time)
  Prometheus telemetry.

## Metrics

Depending on
the [mode](https://developer.hashicorp.com/consul/docs/install/glossary#agent), the collector collects a different
number of metrics.

All metrics have "consul." prefix.

Labels per scope:

- global: datacenter, node_name.
- node check: datacenter, node_name, check_name.
- service check: datacenter, node_name, check_name, service_name.

| Metric                                 |     Scope     |                Dimensions                 |     Units     | Server Leader | Server Follower | Client |
|----------------------------------------|:-------------:|:-----------------------------------------:|:-------------:|:-------------:|:---------------:|:------:|
| node_health_check_status               |  node check   |  passing, maintenance, warning, critical  |    status     |      yes      |       yes       |  yes   |
| service_health_check_status            | service check |  passing, maintenance, warning, critical  |    status     |      yes      |       yes       |  yes   |
| client_rpc_requests_rate               |    global     |                    rpc                    |  requests/s   |      yes      |       yes       |  yes   |
| client_rpc_requests_exceeded_rate      |    global     |                 exceeded                  |  requests/s   |      yes      |       yes       |  yes   |
| client_rpc_requests_failed_rate        |    global     |                  failed                   |  requests/s   |      yes      |       yes       |  yes   |
| memory_allocated                       |    global     |                 allocated                 |     bytes     |      yes      |       yes       |  yes   |
| memory_sys                             |    global     |                    sys                    |     bytes     |      yes      |       yes       |  yes   |
| gc_pause_time                          |    global     |                 gc_pause                  |    seconds    |      yes      |       yes       |  yes   |
| kvs_apply_time                         |    global     | quantile_0.5, quantile_0.9, quantile_0.99 |      ms       |      yes      |       yes       |   no   |
| kvs_apply_operations_rate              |    global     |                 kvs_apply                 |     ops/s     |      yes      |       yes       |   no   |
| txn_apply_time                         |    global     | quantile_0.5, quantile_0.9, quantile_0.99 |      ms       |      yes      |       yes       |   no   |
| txn_apply_operations_rate              |    global     |                 txn_apply                 |     ops/s     |      yes      |       yes       |   no   |
| raft_commit_time                       |    global     | quantile_0.5, quantile_0.9, quantile_0.99 |      ms       |      yes      |       no        |   no   |
| raft_commits_rate                      |    global     |                  commits                  |   commits/s   |      yes      |       no        |   no   |
| autopilot_health_status                |    global     |            healthy, unhealthy             |    status     |      yes      |       yes       |   no   |
| autopilot_failure_tolerance            |    global     |             failure_tolerance             |    servers    |      yes      |       yes       |   no   |
| autopilot_server_health_status         |    global     |            healthy, unhealthy             |    status     |      yes      |       yes       |   no   |
| autopilot_server_stable_time           |    global     |                  stable                   |    seconds    |      yes      |       yes       |   no   |
| autopilot_server_serf_status           |    global     |        active, failed, left, none         |    status     |      yes      |       yes       |   no   |
| autopilot_server_voter_status          |    global     |             voter, not_voter              |    status     |      yes      |       yes       |   no   |
| raft_leader_last_contact_time          |    global     | quantile_0.5, quantile_0.9, quantile_0.99 |      ms       |      yes      |       no        |   no   |
| raft_follower_last_contact_leader_time |    global     |            leader_last_contact            |      ms       |      no       |       yes       |   no   |
| raft_leader_elections_rate             |    global     |                  leader                   |  elections/s  |      yes      |       yes       |   no   |
| raft_leadership_transitions_rate       |    global     |                leadership                 | transitions/s |      yes      |       yes       |   no   |
| server_leadership_status               |    global     |            leader, not_leader             |    status     |      yes      |       yes       |   no   |
| raft_thread_main_saturation_perc       |    global     | quantile_0.5, quantile_0.9, quantile_0.99 |  percentage   |      yes      |       yes       |   no   |
| raft_thread_fsm_saturation_perc        |    global     | quantile_0.5, quantile_0.9, quantile_0.99 |  percentage   |      yes      |       yes       |   no   |
| raft_fsm_last_restore_duration         |    global     |           last_restore_duration           |      ms       |      yes      |       yes       |   no   |
| raft_leader_oldest_log_age             |    global     |              oldest_log_age               |    seconds    |      yes      |       no        |   no   |
| raft_rpc_install_snapshot_time         |    global     | quantile_0.5, quantile_0.9, quantile_0.99 |      ms       |      no       |       yes       |   no   |
| raft_boltdb_freelist_bytes             |    global     |                 freelist                  |     bytes     |      yes      |       yes       |   no   |
| raft_boltdb_logs_per_batch_rate        |    global     |                  written                  |    logs/s     |      yes      |       yes       |   no   |
| raft_boltdb_store_logs_time            |    global     | quantile_0.5, quantile_0.9, quantile_0.99 |      ms       |      yes      |       yes       |   no   |

## Configuration

Edit the `go.d/consul.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/consul.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8500
    acl_token: "ec15675e-2999-d789-832e-8c4794daa8d7"

  - name: remote
    url: http://203.0.113.10:8500
    acl_token: "ada7f751-f654-8872-7f93-498e799158b6"
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/consul.conf).

## Troubleshooting

To troubleshoot issues with the `consul` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m consul
  ```
