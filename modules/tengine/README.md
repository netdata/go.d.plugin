# Tengine collector

## Overview

[Tengine](https://tengine.taobao.org/) is a web server originated by Taobao and is based on
the [NGINX](https://nginx.org/en/).

This collector monitors one or more Tengine instances, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                                   |                               Dimensions                               |     Unit      |
|----------------------------------------------------------|:----------------------------------------------------------------------:|:-------------:|
| tengine.bandwidth_total                                  |                                in, out                                 |      B/s      |
| tengine.connections_total                                |                                accepted                                | connections/s |
| tengine.requests_total                                   |                               processed                                |  requests/s   |
| tengine.requests_per_response_code_family_total          |                       2xx, 3xx, 4xx, 5xx, other                        |  requests/s   |
| tengine.requests_per_response_code_detailed_total        | 200, 206, 302, 304, 403, 404, 419, 499, 500, 502, 503, 504, 508, other |  requests/s   |
| tengine.requests_upstream_total                          |                                requests                                |  requests/s   |
| tengine.tries_upstream_total                             |                                 calls                                  |    calls/s    |
| tengine.requests_upstream_per_response_code_family_total |                                4xx, 5xx                                |  requests/s   |

## Setup

### Prerequisites

#### Enable ngx_http_reqstat_module module.

See [ngx_http_reqstat_module](https://tengine.taobao.org/document/http_reqstat.html) documentation.
The default line format is the only supported format.

### Configuration

#### File

The configuration file name is `go.d/tengine.conf`.

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
sudo ./edit-config go.d/tengine.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |       Default       | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:-------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |          1          |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |          0          |          |
|         url          | Server URL.                                                                                               | http://127.0.0.1/us |   yes    |
|       timeout        | HTTP request timeout.                                                                                     |          2          |          |
|       username       | Username for basic HTTP authentication.                                                                   |                     |          |
|       password       | Password for basic HTTP authentication.                                                                   |                     |          |
|      proxy_url       | Proxy URL.                                                                                                |                     |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                     |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                     |          |
|        method        | HTTP request method.                                                                                      |         GET         |          |
|         body         | HTTP request body.                                                                                        |                     |          |
|       headers        | HTTP request headers.                                                                                     |                     |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |         no          |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |         no          |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                     |          |
|       tls_cert       | Client TLS certificate.                                                                                   |                     |          |
|       tls_key        | Client TLS key.                                                                                           |                     |          |

</details>

#### Examples

##### Basic

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/us
```

</details>

##### HTTP authentication

Local server with basic HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/us
    username: foo
    password: bar
```

</details>

##### HTTPS with self-signed certificate

Tengine with enabled HTTPS and self-signed certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: https://127.0.0.1/us
    tls_skip_verify: yes
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
    url: http://127.0.0.1/us

  - name: remote
    url: http://203.0.113.10/us
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `tengine` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m tengine
  ```
