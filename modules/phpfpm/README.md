# PHP-FPM collector

## Overview

[PHP-FPM](https://php-fpm.org/) is an alternative PHP FastCGI implementation with some additional features useful for
sites of any size, especially busier sites.

This collector monitors one or more PHP-FPM instances, depending on your configuration..

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                  |             Dimensions              |     Unit     |
|-------------------------|:-----------------------------------:|:------------:|
| phpfpm.connections      |      active, max_active, idle       | connections  |
| phpfpm.requests         |              requests               |  requests/s  |
| phpfpm.performance      | max_children_reached, slow_requests |    status    |
| phpfpm.request_duration |            min, max, avg            | milliseconds |
| phpfpm.request_cpu      |            min, max, avg            |  percentage  |
| phpfpm.request_mem      |            min, max, avg            |      KB      |

## Setup

### Prerequisites

#### Enable status page

Uncomment the `pm.status_path = /status` variable in the `php-fpm` config file.

### Configuration

#### File

The configuration file name is `go.d/phpfpm.conf`.

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
sudo ./edit-config go.d/phpfpm.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |              Default              | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:---------------------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |                 1                 |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |                 0                 |          |
|         url          | Server URL.                                                                                               | http://127.0.0.1/status?full&json |   yes    |
|        socket        | Server Unix socket.                                                                                       |                                   |          |
|       address        | Server address in IP:PORT format.                                                                         |                                   |          |
|      fcgi_path       | Status path.                                                                                              |              /status              |          |
|       timeout        | HTTP request timeout.                                                                                     |                 1                 |          |
|       username       | Username for basic HTTP authentication.                                                                   |                                   |          |
|       password       | Password for basic HTTP authentication.                                                                   |                                   |          |
|      proxy_url       | Proxy URL.                                                                                                |                                   |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                                   |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                                   |          |
|        method        | HTTP request method.                                                                                      |                GET                |          |
|         body         | HTTP request body.                                                                                        |                                   |          |
|       headers        | HTTP request headers.                                                                                     |                                   |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |                no                 |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |                no                 |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                                   |          |
|       tls_cert       | Client TLS certificate.                                                                                   |                                   |          |
|       tls_key        | Client TLS key.                                                                                           |                                   |          |

</details>

#### Examples

##### HTTP

Collecting data from a local instance over HTTP.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://localhost/status?full&json
```

</details>

##### Unix socket

Collecting data from a local instance over Unix socket.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    socket: '/tmp/php-fpm.sock'
```

</details>

##### TCP socket

Collecting data from a local instance over TCP socket.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:9000
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
    url: http://localhost/status?full&json

  - name: remote
    url: http://203.0.113.10/status?full&json
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `phpfpm` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m phpfpm
  ```
