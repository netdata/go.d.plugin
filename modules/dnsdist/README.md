# DNSdist collector

## Overview

[DNSdist](https://dnsdist.org/) is a highly DNS-, DoS- and abuse-aware loadbalancer.

This collector monitors load-balancer performance and health metrics.

It collects metrics from [the internal webserver](https://dnsdist.org/guides/webserver.html).

Used endpoints:

- [/jsonstat?command=stats](https://dnsdist.org/guides/webserver.html#get--jsonstat).

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                     |                     Dimensions                     |     Unit     |
|----------------------------|:--------------------------------------------------:|:------------:|
| dnsdist.queries            |               all, recursive, empty                |  queries/s   |
| dnsdist.queries_dropped    | rule_drop, dynamic_blocked, no_policy, non_queries |  queries/s   |
| dnsdist.packets_dropped    |                        acl                         |  packets/s   |
| dnsdist.answers            |  self_answered, nxdomain, refused, trunc_failures  |  answers/s   |
| dnsdist.backend_responses  |                     responses                      | responses/s  |
| dnsdist.backend_commerrors |                    send_errors                     |   errors/s   |
| dnsdist.backend_errors     |         timeouts, servfail, non_compliant          | responses/s  |
| dnsdist.cache              |                    hits, misses                    |  answers/s   |
| dnsdist.servercpu          |              system_state, user_state              |     ms/s     |
| dnsdist.servermem          |                    memory_usage                    |     MiB      |
| dnsdist.query_latency      |         1ms, 10ms, 50ms, 100ms, 1sec, slow         |  queries/s   |
| dnsdist.query_latency_avg  |                100, 1k, 10k, 1000k                 | microseconds |

## Setup

### Prerequisites

#### Enable DNSdist's built-in Webserver

For collecting metrics via HTTP, you need to [enable the built-in webserver](https://dnsdist.org/guides/webserver.html).

### Configuration

#### File

The configuration file name is `go.d/dnsdist.conf`.

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
sudo ./edit-config go.d/dnsdist.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |        Default        | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:---------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |           1           |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |           0           |          |
|         url          | Server URL.                                                                                               | http://127.0.0.1:8083 |   yes    |
|       username       | Username for basic HTTP authentication.                                                                   |                       |          |
|       password       | Password for basic HTTP authentication.                                                                   |                       |          |
|      proxy_url       | Proxy URL.                                                                                                |                       |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                       |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                       |          |
|       timeout        | HTTP request timeout.                                                                                     |           1           |          |
|        method        | HTTP request method.                                                                                      |          GET          |          |
|         body         | HTTP request body.                                                                                        |           -           |          |
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
    url: http://127.0.0.1:8083
    headers:
      X-API-Key: your-api-key # static pre-shared authentication key for access to the REST API (api-key).
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
    url: http://127.0.0.1:8083
    headers:
      X-API-Key: 'your-api-key' # static pre-shared authentication key for access to the REST API (api-key).

  - name: remote
    url: http://203.0.113.0:8083
    headers:
      X-API-Key: 'your-api-key'
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `dnsdist` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m dnsdist
  ```
