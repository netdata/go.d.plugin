# PowerDNS Authoritative Server collector

## Overview

[PowerDNS Authoritative Server](https://doc.powerdns.com/authoritative/) is a versatile nameserver which supports a
large number of backends.

This collector monitors one or more PowerDNS Authoritative Server instances, depending on your configuration.

It collects metrics from [the internal webserver](https://doc.powerdns.com/authoritative/http-api/index.html#webserver).

Used endpoints:

- [`/api/v1/servers/localhost/statistics`](https://doc.powerdns.com/authoritative/http-api/statistics.html)

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                 |                              Dimensions                              |     Unit     |
|------------------------|:--------------------------------------------------------------------:|:------------:|
| powerdns.questions_in  |                               udp, tcp                               | questions/s  |
| powerdns.questions_out |                               udp, tcp                               | questions/s  |
| powerdns.cache_usage   | query-cache-hit, query-cache-miss, packetcache-hit, packetcache-miss |   events/s   |
| powerdns.cache_size    |           query-cache, packet-cache, key-cache, meta-cache           |   entries    |
| powerdns.latency       |                               latency                                | microseconds |

## Setup

### Prerequisites

#### Enable webserver

Follow [webserver](https://doc.powerdns.com/authoritative/http-api/index.html#webserver) documentation.

#### Enable HTTP API

Follow [HTTP API](https://doc.powerdns.com/authoritative/http-api/index.html#enabling-the-api) documentation.

### Configuration

#### File

The configuration file name is `go.d/powerdns.conf`.

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
sudo ./edit-config go.d/powerdns.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |        Default        | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:---------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |           1           |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |           0           |          |
|         url          | Server URL.                                                                                               | http://127.0.0.1:8081 |   yes    |
|       timeout        | HTTP request timeout.                                                                                     |           1           |          |
|       username       | Username for basic HTTP authentication.                                                                   |                       |          |
|       password       | Password for basic HTTP authentication.                                                                   |                       |          |
|      proxy_url       | Proxy URL.                                                                                                |                       |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                       |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                       |          |
|        method        | HTTP request method.                                                                                      |          GET          |          |
|         body         | HTTP request body.                                                                                        |                       |          |
|       headers        | HTTP request headers.                                                                                     |                       |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |          no           |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |          no           |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                       |          |
|       tls_cert       | Client TLS certificate.                                                                                   |                       |          |
|       tls_key        | Client TLS key.                                                                                           |                       |          |

</details>

#### Examples

##### Basic

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8081
```

</details>

##### HTTP authentication

Basic HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8081
    username: admin
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
    url: http://127.0.0.1:8081

  - name: remote
    url: http://203.0.113.0:8081
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `powerdns` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m powerdns
  ```
