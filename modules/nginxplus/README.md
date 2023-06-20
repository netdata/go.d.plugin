# NGINX Plus collector

## Overview

[NGINX Plus](https://www.nginx.com/products/nginx/) is a software load balancer, API gateway, and reverse proxy built on
top of NGINX.

This collector will monitor one or more NGINX Plus servers, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                 |                            Dimensions                             |     Unit      |
|----------------------------------------|:-----------------------------------------------------------------:|:-------------:|
| nginxplus.client_connections_rate      |                         accepted, dropped                         | connections/s |
| nginxplus.client_connections_count     |                           active, idle                            |  connections  |
| nginxplus.ssl_handshakes_rate          |                        successful, failed                         | handshakes/s  |
| nginxplus.ssl_handshakes_failures_rate | no_common_protocol, no_common_cipher, timeout, peer_rejected_cert |  failures/s   |
| nginxplus.ssl_verification_errors_rate |   no_cert, expired_cert, revoked_cert, hostname_mismatch, other   |   errors/s    |
| nginxplus.ssl_session_reuses_rate      |                            ssl_session                            |   reuses/s    |
| nginxplus.http_requests_rate           |                             requests                              |  requests/s   |
| nginxplus.http_requests_count          |                             requests                              |   requests    |
| nginxplus.uptime                       |                              uptime                               |    seconds    |

### http server zone

These metrics refer to the HTTP server zone.

Labels:

| Label            | Description           |
|------------------|-----------------------|
| http_server_zone | HTTP server zone name |

Metrics:

| Metric                                                   |       Dimensions        |    Unit     |
|----------------------------------------------------------|:-----------------------:|:-----------:|
| nginxplus.http_server_zone_requests_rate                 |        requests         | requests/s  |
| nginxplus.http_server_zone_responses_per_code_class_rate | 1xx, 2xx, 3xx, 4xx, 5xx | responses/s |
| nginxplus.http_server_zone_traffic_rate                  |     received, sent      |   bytes/s   |
| nginxplus.http_server_zone_requests_processing_count     |       processing        |  requests   |
| nginxplus.http_server_zone_requests_discarded_rate       |        discarded        | requests/s  |

### http location zone

These metrics refer to the HTTP location zone.

Labels:

| Label              | Description             |
|--------------------|-------------------------|
| http_location_zone | HTTP location zone name |

Metrics:

| Metric                                                     |       Dimensions        |    Unit     |
|------------------------------------------------------------|:-----------------------:|:-----------:|
| nginxplus.http_location_zone_requests_rate                 |        requests         | requests/s  |
| nginxplus.http_location_zone_responses_per_code_class_rate | 1xx, 2xx, 3xx, 4xx, 5xx | responses/s |
| nginxplus.http_location_zone_traffic_rate                  |     received, sent      |   bytes/s   |
| nginxplus.http_location_zone_requests_discarded_rate       |        discarded        | requests/s  |

### http upstream

These metrics refer to the HTTP upstream.

Labels:

| Label              | Description             |
|--------------------|-------------------------|
| http_upstream_name | HTTP upstream name      |
| http_upstream_zone | HTTP upstream zone name |

Metrics:

| Metric                                  | Dimensions |    Unit     |
|-----------------------------------------|:----------:|:-----------:|
| nginxplus.http_upstream_peers_count     |   peers    |    peers    |
| nginxplus.http_upstream_zombies_count   |   zombie   |   servers   |
| nginxplus.http_upstream_keepalive_count | keepalive  | connections |

### http upstream server

These metrics refer to the HTTP upstream server.

Labels:

| Label                        | Description                                      |
|------------------------------|--------------------------------------------------|
| http_upstream_name           | HTTP upstream name                               |
| http_upstream_zone           | HTTP upstream zone name                          |
| http_upstream_server_address | HTTP upstream server address (e.g. 127.0.0.1:81) |
| http_upstream_server_name    | HTTP upstream server name                        |

Metrics:

| Metric                                                       |                    Dimensions                    |     Unit     |
|--------------------------------------------------------------|:------------------------------------------------:|:------------:|
| nginxplus.http_upstream_server_requests_rate                 |                     requests                     |  requests/s  |
| nginxplus.http_upstream_server_responses_per_code_class_rate |             1xx, 2xx, 3xx, 4xx, 5xx              | responses/s  |
| nginxplus.http_upstream_server_response_time                 |                     response                     | milliseconds |
| nginxplus.http_upstream_server_response_header_time          |                      header                      | milliseconds |
| nginxplus.http_upstream_server_traffic_rate                  |                  received, sent                  |   bytes/s    |
| nginxplus.http_upstream_server_state                         | up, down, draining, unavail, checking, unhealthy |    state     |
| nginxplus.http_upstream_server_connections_count             |                      active                      | connections  |
| nginxplus.http_upstream_server_downtime                      |                     downtime                     |   seconds    |

### http cache

These metrics refer to the HTTP cache.

Labels:

| Label      | Description     |
|------------|-----------------|
| http_cache | HTTP cache name |

Metrics:

| Metric                     |       Dimensions        |    Unit     |
|----------------------------|:-----------------------:|:-----------:|
| nginxplus.http_cache_state |       warm, cold        |    state    |
| nginxplus.http_cache_iops  | served, written, bypass | responses/s |
| nginxplus.http_cache_io    | served, written, bypass |   bytes/s   |
| nginxplus.http_cache_size  |          size           |    bytes    |

### stream server zone

These metrics refer to the Stream server zone.

Labels:

| Label              | Description             |
|--------------------|-------------------------|
| stream_server_zone | Stream server zone name |

Metrics:

| Metric                                                    |   Dimensions   |     Unit      |
|-----------------------------------------------------------|:--------------:|:-------------:|
| nginxplus.stream_server_zone_connections_rate             |    accepted    | connections/s |
| nginxplus.stream_server_zone_sessions_per_code_class_rate | 2xx, 4xx, 5xx  |  sessions/s   |
| nginxplus.stream_server_zone_traffic_rate                 | received, sent |    bytes/s    |
| nginxplus.stream_server_zone_connections_processing_count |   processing   |  connections  |
| nginxplus.stream_server_zone_connections_discarded_rate   |   discarded    | connections/s |

### stream upstream

These metrics refer to the Stream upstream.

Labels:

| Label                | Description               |
|----------------------|---------------------------|
| stream_upstream_name | Stream upstream name      |
| stream_upstream_zone | Stream upstream zone name |

Metrics:

| Metric                                  | Dimensions |  Unit   |
|-----------------------------------------|:----------:|:-------:|
| nginxplus.stream_upstream_peers_count   |   peers    |  peers  |
| nginxplus.stream_upstream_zombies_count |   zombie   | servers |

### stream upstream server

These metrics refer to the Stream upstream server.

Labels:

| Label                          | Description                                           |
|--------------------------------|-------------------------------------------------------|
| stream_upstream_name           | Stream upstream name                                  |
| stream_upstream_zone           | Stream upstream zone name                             |
| stream_upstream_server_address | Stream upstream server address (e.g. 127.0.0.1:12346) |
| stream_upstream_server_name    | Stream upstream server name                           |

Metrics:

| Metric                                             |               Dimensions               |     Unit      |
|----------------------------------------------------|:--------------------------------------:|:-------------:|
| nginxplus.stream_upstream_server_connections_rate  |               forwarded                | connections/s |
| nginxplus.stream_upstream_server_traffic_rate      |             received, sent             |    bytes/s    |
| nginxplus.stream_upstream_server_state             | up, down, unavail, checking, unhealthy |     state     |
| nginxplus.stream_upstream_server_downtime          |                downtime                |    seconds    |
| nginxplus.stream_upstream_server_connections_count |                 active                 |  connections  |

### resolver zone

These metrics refer to the resolver zone.

Labels:

| Label         | Description        |
|---------------|--------------------|
| resolver_zone | resolver zone name |

Metrics:

| Metric                                 |                                Dimensions                                |    Unit     |
|----------------------------------------|:------------------------------------------------------------------------:|:-----------:|
| nginxplus.resolver_zone_requests_rate  |                             name, srv, addr                              | requests/s  |
| nginxplus.resolver_zone_responses_rate | noerror, formerr, servfail, nxdomain, notimp, refused, timedout, unknown | responses/s |

## Setup

### Prerequisites

#### Enable API

See [Configure the API](https://docs.nginx.com/nginx/admin-guide/monitoring/live-activity-monitoring/#configuring-the-api).

### Configuration

#### File

The configuration file name is `go.d/nginxplus.conf`.

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
sudo ./edit-config go.d/nginxplus.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |     Default      | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:----------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |        1         |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |        0         |          |
|         url          | Server URL.                                                                                               | http://127.0.0.1 |   yes    |
|       timeout        | HTTP request timeout.                                                                                     |        1         |          |
|       username       | Username for basic HTTP authentication.                                                                   |                  |          |
|       password       | Password for basic HTTP authentication.                                                                   |                  |          |
|      proxy_url       | Proxy URL.                                                                                                |                  |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                  |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                  |          |
|        method        | HTTP request method.                                                                                      |       GET        |          |
|         body         | HTTP request body.                                                                                        |                  |          |
|       headers        | HTTP request headers.                                                                                     |                  |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |        no        |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |        no        |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                  |          |
|       tls_cert       | Client TLS certificate.                                                                                   |                  |          |
|       tls_key        | Client TLS key.                                                                                           |                  |          |

</details>

#### Examples

##### Basic

A basic example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1
```

</details>

##### HTTP authentication

Basic HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1
    username: username
    password: password
```

</details>

##### HTTPS with self-signed certificate

NGINX Plus with enabled HTTPS and self-signed certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: https://127.0.0.1
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
    url: http://127.0.0.1

  - name: remote
    url: http://192.0.2.1
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `nginxplus` collector, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

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
