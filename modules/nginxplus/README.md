<!--
title: "NGINX Plus monitoring"
description: "Monitor the health and performance of NGINX Plus web servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/nginxplus/README.md
sidebar_label: "NGINX Plus"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Webapps"
-->

# NGINX Plus monitoring

[NGINX Plus](https://www.nginx.com/products/nginx/) is a software load balancer, API gateway, and reverse proxy built on
top of NGINX.

This module will monitor one or more NGINX Plus servers, depending on your configuration.

## Requirements

- NGINX Plus with
  configured [API](https://docs.nginx.com/nginx/admin-guide/monitoring/live-activity-monitoring/#configuring-the-api).

## Metrics

All metrics have "nginxplus." prefix.

Labels per scope:

- global: no labels.
- http server zone: http_server_zone.
- http location zone: http_location_zone.
- http upstream: http_upstream_name, http_upstream_zone.
- http upstream server: http_upstream_name, http_upstream_zone, http_upstream_server_address, http_upstream_server_name.
- http cache: http_cache.
- stream server zone: stream_server_zone.
- stream upstream: stream_upstream_name, stream_upstream_zone.
- stream upstream server: stream_upstream_name, stream_upstream_zone, stream_upstream_server_address,
  stream_upstream_server_name.
- resolver zone: resolver_zone.

| Metric                                             |         Scope          |                                Dimensions                                |     Units     |
|----------------------------------------------------|:----------------------:|:------------------------------------------------------------------------:|:-------------:|
| client_connections_rate                            |         global         |                            accepted, dropped                             | connections/s |
| client_connections_count                           |         global         |                               active, idle                               |  connections  |
| ssl_handshakes_rate                                |         global         |                            successful, failed                            | handshakes/s  |
| ssl_handshakes_failures_rate                       |         global         |    no_common_protocol, no_common_cipher, timeout, peer_rejected_cert     |  failures/s   |
| ssl_verification_errors_rate                       |         global         |      no_cert, expired_cert, revoked_cert, hostname_mismatch, other       |   errors/s    |
| ssl_session_reuses_rate                            |         global         |                               ssl_session                                |   reuses/s    |
| http_requests_rate                                 |         global         |                                 requests                                 |  requests/s   |
| http_requests_count                                |         global         |                                 requests                                 |   requests    |
| uptime                                             |         global         |                                  uptime                                  |    seconds    |
| http_server_zone_requests_rate                     |    http server zone    |                                 requests                                 |  requests/s   |
| http_server_zone_responses_per_code_class_rate     |    http server zone    |                         1xx, 2xx, 3xx, 4xx, 5xx                          |  responses/s  |
| http_server_zone_traffic_rate                      |    http server zone    |                              received, sent                              |    bytes/s    |
| http_server_zone_requests_processing_count         |    http server zone    |                                processing                                |   requests    |
| http_server_zone_requests_discarded_rate           |    http server zone    |                                discarded                                 |  requests/s   |
| http_location_zone_requests_rate                   |   http location zone   |                                 requests                                 |  requests/s   |
| http_location_zone_responses_per_code_class_rate   |   http location zone   |                         1xx, 2xx, 3xx, 4xx, 5xx                          |  responses/s  |
| http_location_zone_traffic_rate                    |   http location zone   |                              received, sent                              |    bytes/s    |
| http_location_zone_requests_discarded_rate         |   http location zone   |                                discarded                                 |  requests/s   |
| http_upstream_peers_count                          |     http upstream      |                                  peers                                   |     peers     |
| http_upstream_zombies_count                        |     http upstream      |                                  zombie                                  |    servers    |
| http_upstream_keepalive_count                      |     http upstream      |                                keepalive                                 |  connections  |
| http_upstream_server_requests_rate                 |  http upstream server  |                                 requests                                 |  requests/s   |
| http_upstream_server_responses_per_code_class_rate |  http upstream server  |                         1xx, 2xx, 3xx, 4xx, 5xx                          |  responses/s  |
| http_upstream_server_response_time                 |  http upstream server  |                                 response                                 | milliseconds  |
| http_upstream_server_response_header_time          |  http upstream server  |                                  header                                  | milliseconds  |
| http_upstream_server_traffic_rate                  |  http upstream server  |                              received, sent                              |    bytes/s    |
| http_upstream_server_state                         |  http upstream server  |             up, down, draining, unavail, checking, unhealthy             |     state     |
| http_upstream_server_connections_count             |  http upstream server  |                                  active                                  |  connections  |
| http_upstream_server_downtime                      |  http upstream server  |                                 downtime                                 |    seconds    |
| http_cache_state                                   |       http cache       |                                warm, cold                                |     state     |
| http_cache_iops                                    |       http cache       |                         served, written, bypass                          |  responses/s  |
| http_cache_io                                      |       http cache       |                         served, written, bypass                          |    bytes/s    |
| http_cache_size                                    |       http cache       |                                   size                                   |     bytes     |
| stream_server_zone_connections_rate                |   stream server zone   |                                 accepted                                 | connections/s |
| stream_server_zone_sessions_per_code_class_rate    |   stream server zone   |                              2xx, 4xx, 5xx                               |  sessions/s   |
| stream_server_zone_traffic_rate                    |   stream server zone   |                              received, sent                              |    bytes/s    |
| stream_server_zone_connections_processing_count    |   stream server zone   |                                processing                                |  connections  |
| stream_server_zone_connections_discarded_rate      |   stream server zone   |                                discarded                                 | connections/s |
| stream_upstream_peers_count                        |    stream upstream     |                                  peers                                   |     peers     |
| stream_upstream_zombies_count                      |    stream upstream     |                                  zombie                                  |    servers    |
| stream_upstream_server_connections_rate            | stream upstream server |                                forwarded                                 | connections/s |
| stream_upstream_server_traffic_rate                | stream upstream server |                              received, sent                              |    bytes/s    |
| stream_upstream_server_state                       | stream upstream server |                  up, down, unavail, checking, unhealthy                  |     state     |
| stream_upstream_server_downtime                    | stream upstream server |                                 downtime                                 |    seconds    |
| stream_upstream_server_connections_count           | stream upstream server |                                  active                                  |  connections  |
| resolver_zone_requests_rate                        |     resolver zone      |                             name, srv, addr                              |  requests/s   |
| resolver_zone_responses_rate                       |     resolver zone      | noerror, formerr, servfail, nxdomain, notimp, refused, timedout, unknown |  responses/s  |

## Configuration

Edit the `go.d/nginxplus.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/nginxplus.conf
```

Needs only server's `url`. Here is an example for local and remote servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1

  - name: remote
    url: http://203.0.113.10
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/nginx.conf).

## Troubleshooting

To troubleshoot issues with the `nginxplus` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m nginxplus
  ```
