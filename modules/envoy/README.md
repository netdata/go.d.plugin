<!--
title: "Envoy monitoring with Netdata"
description: "Monitor the health and performance of Envoy web servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/envoy/README.md"
sidebar_label: "Envoy"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Proxy"
-->

# Envoy collector

[Envoy](https://www.envoyproxy.io/docs/envoy/latest/intro/what_is_envoy) is an L7 proxy and communication bus designed
for large modern service oriented architectures.

This module will monitor one or more Envoy servers, depending on your configuration.

## Metrics

All metrics have "envoy." prefix.

| Metric                                                             | Scope  |                                               Dimensions                                               |     Units     |
|--------------------------------------------------------------------|:------:|:------------------------------------------------------------------------------------------------------:|:-------------:|
| server_state                                                       | global |                             live, draining, pre_initializing, initializing                             |     state     |
| server_connections_count                                           | global |                                              connections                                               |  connections  |
| server_parent_connections_count                                    | global |                                              connections                                               |  connections  |
| server_memory_allocated_size                                       | global |                                               allocated                                                |     bytes     |
| server_memory_heap_size                                            | global |                                                  heap                                                  |     bytes     |
| server_memory_physical_size                                        | global |                                                physical                                                |     bytes     |
| server_uptime                                                      | global |                                                 uptime                                                 |    seconds    |
| cluster_manager_cluster_count                                      | global |                                           active, not_active                                           |   clusters    |
| cluster_manager_cluster_changes_rate                               | global |                                        added, modified, removed                                        |  clusters/s   |
| cluster_manager_cluster_updates_rate                               | global |                                                cluster                                                 |   updates/s   |
| cluster_manager_cluster_updated_via_merge_rate                     | global |                                               via_merge                                                |   updates/s   |
| cluster_manager_update_merge_cancelled_rate                        | global |                                            merge_cancelled                                             |   updates/s   |
| cluster_manager_update_out_of_merge_window_rate                    | global |                                          out_of_merge_window                                           |   updates/s   |
| cluster_membership_endpoints_count                                 | global |                                      healthy, degraded, excluded                                       |   endpoints   |
| cluster_membership_changes_rate                                    | global |                                               membership                                               |   changes/s   |
| cluster_membership_updates_rate                                    | global |                                  success, failure, empty, no_rebuild                                   |   updates/s   |
| cluster_upstream_cx_active_count                                   | global |                                                 active                                                 |  connections  |
| cluster_upstream_cx_rate                                           | global |                                                created                                                 | connections/s |
| cluster_upstream_cx_http_rate                                      | global |                                          http1, http2, http3                                           | connections/s |
| cluster_upstream_cx_destroy_rate                                   | global |                                             local, remote                                              | connections/s |
| cluster_upstream_cx_connect_fail_rate                              | global |                                                 failed                                                 | connections/s |
| cluster_upstream_cx_connect_timeout_rate                           | global |                                                timeout                                                 | connections/s |
| cluster_upstream_cx_bytes_rate                                     | global |                                             received, sent                                             |    bytes/s    |
| cluster_upstream_cx_bytes_buffered_size                            | global |                                             received, send                                             |     bytes     |
| cluster_upstream_rq_active_count                                   | global |                                                 active                                                 |   requests    |
| cluster_upstream_rq_rate                                           | global |                                                requests                                                |  requests/s   |
| cluster_upstream_rq_failed_rate                                    | global | cancelled, maintenance_mode, timeout, max_duration_reached, per_try_timeout, reset_local, reset_remote |  requests/s   |
| cluster_upstream_rq_pending_active_count                           | global |                                             active_pending                                             |   requests    |
| cluster_upstream_rq_pending_rate                                   | global |                                                pending                                                 |  requests/s   |
| cluster_upstream_rq_pending_failed_rate                            | global |                                        overflow, failure_eject                                         |  requests/s   |
| cluster_upstream_rq_retry_rate                                     | global |                                                request                                                 |   retries/s   |
| cluster_upstream_rq_retry_success_rate                             | global |                                                success                                                 |   retries/s   |
| cluster_upstream_rq_retry_backoff_rate                             | global |                                        exponential, ratelimited                                        |   retries/s   |
| listener_manager_listeners_count                                   | global |                                       active, warming, draining                                        |   listeners   |
| listener_manager_listener_changes_rate                             | global |                                   added, modified, removed, stopped                                    |  listeners/s  |
| listener_manager_listener_object_events_rate                       | global |                            create_success, create_failure, in_place_updated                            |   objects/s   |
| listener_admin_downstream_cx_active_count                          | global |                                                 active                                                 |  connections  |
| listener_admin_downstream_cx_rate                                  | global |                                                created                                                 | connections/s |
| listener_admin_downstream_cx_destroy_rate                          | global |                                               destroyed                                                | connections/s |
| listener_admin_downstream_cx_transport_socket_connect_timeout_rate | global |                                                timeout                                                 | connections/s |
| listener_admin_downstream_cx_rejected_rate                         | global |                                  overflow, overload, global_overflow                                   | connections/s |
| listener_admin_downstream_listener_filter_remote_close_rate        | global |                                                 closed                                                 | connections/s |
| listener_admin_downstream_listener_filter_error_rate               | global |                                                  read                                                  |   errors/s    |
| listener_admin_downstream_pre_cx_active_count                      | global |                                                 active                                                 |    sockets    |
| listener_admin_downstream_pre_cx_timeout_rate                      | global |                                                timeout                                                 |   sockets/s   |
| listener_downstream_cx_active_count                                | global |                                                 active                                                 |  connections  |
| listener_downstream_cx_rate                                        | global |                                                created                                                 | connections/s |
| listener_downstream_cx_destroy_rate                                | global |                                               destroyed                                                | connections/s |
| listener_downstream_cx_transport_socket_connect_timeout_rate       | global |                                                timeout                                                 | connections/s |
| listener_downstream_cx_rejected_rate                               | global |                                  overflow, overload, global_overflow                                   | connections/s |
| listener_downstream_listener_filter_remote_close_rate              | global |                                                 closed                                                 | connections/s |
| listener_downstream_listener_filter_error_rate                     | global |                                                  read                                                  |   errors/s    |
| listener_downstream_pre_cx_active_count                            | global |                                                 active                                                 |    sockets    |
| listener_downstream_pre_cx_timeout_rate                            | global |                                                timeout                                                 |   sockets/s   |

## Configuration

Edit the `go.d/envoy.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/envoy.conf

```

Needs only `url` to server's `/stats/prometheus` endpoint. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9901/stats/prometheus

  - name: remote
    url: http://203.0.113.10:9901/stats/prometheus
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/envoy.conf).

## Troubleshooting

To troubleshoot issues with the `envoy` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m envoy
  ```
