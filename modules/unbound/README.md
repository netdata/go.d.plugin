# Unbound collector

## Overview

[Unbound](https://nlnetlabs.nl/projects/unbound/about/) is a validating, recursive, and caching DNS resolver product
from NLnet Labs.

This collector monitors one or more Unbound servers, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                             |                       Dimensions                       |     Unit     |
|------------------------------------|:------------------------------------------------------:|:------------:|
| unbound.queries                    |                        queries                         |   queries    |
| unbound.queries_ip_ratelimited     |                      ratelimited                       |   queries    |
| unbound.dnscrypt_queries           |          crypted, cert, cleartext, malformed           |   queries    |
| unbound.cache                      |                       hits, miss                       |    events    |
| unbound.cache_percentage           |                       hits, miss                       |  percentage  |
| unbound.prefetch                   |                       prefetches                       |  prefetches  |
| unbound.expired                    |                        expired                         |   replies    |
| unbound.zero_ttl_replies           |                        zero_ttl                        |   replies    |
| unbound.recursive_replies          |                       recursive                        |   replies    |
| unbound.recursion_time             |                      avg, median                       | milliseconds |
| unbound.request_list_usage         |                        avg, max                        |   queries    |
| unbound.current_request_list_usage |                       all, users                       |   queries    |
| unbound.request_list_jostle_list   |                  overwritten, dropped                  |   queries    |
| unbound.tcpusage                   |                         usage                          |   buffers    |
| unbound.uptime                     |                          time                          |   seconds    |
| unbound.cache_memory               | message, rrset, dnscrypt_nonce, dnscrypt_shared_secret |      KB      |
| unbound.mod_memory                 |       iterator, respip, validator, subnet, ipsec       |      KB      |
| unbound.mem_streamwait             |                       streamwait                       |      KB      |
| unbound.cache_count                | infra, key, msg, rrset, dnscrypt_nonce, shared_secret  |    items     |
| unbound.type_queries               |               a dimension per query type               |   queries    |
| unbound.class_queries              |              a dimension per query class               |   queries    |
| unbound.opcode_queries             |              a dimension per query opcode              |   queries    |
| unbound.flag_queries               |             qr, aa, tc, rd, ra, z, ad, cd              |   queries    |
| unbound.rcode_answers              |              a dimension per reply rcode               |   replies    |

### thread

These metrics refer to threads.

This scope has no labels.

Metrics:

| Metric                                    |             Dimensions              |     Unit     |
|-------------------------------------------|:-----------------------------------:|:------------:|
| unbound.thread_queries                    |               queries               |   queries    |
| unbound.thread_queries_ip_ratelimited     |             ratelimited             |   queries    |
| unbound.thread_dnscrypt_queries           | crypted, cert, cleartext, malformed |   queries    |
| unbound.thread_cache                      |             hits, miss              |    events    |
| unbound.thread_cache_percentage           |             hits, miss              |  percentage  |
| unbound.thread_prefetch                   |             prefetches              |  prefetches  |
| unbound.thread_expired                    |               expired               |   replies    |
| unbound.thread_zero_ttl_replies           |              zero_ttl               |   replies    |
| unbound.thread_recursive_replies          |              recursive              |   replies    |
| unbound.thread_recursion_time             |             avg, median             | milliseconds |
| unbound.thread_request_list_usage         |              avg, max               |   queries    |
| unbound.thread_current_request_list_usage |             all, users              |   queries    |
| unbound.thread_request_list_jostle_list   |        overwritten, dropped         |   queries    |
| unbound.thread_tcpusage                   |                usage                |   buffers    |

## Setup

### Prerequisites

#### Enable remote control interface

Set `control-enable` to yes in [unbound.conf](https://nlnetlabs.nl/documentation/unbound/unbound.conf).

#### Check permissions and adjust if necessary

If using unix socket:

- socket should be readable and writeable by `netdata` user

If using ip socket and TLS is disabled:

- socket should be accessible via network

If TLS is enabled, in addition:

- `control-key-file` should be readable by `netdata` user
- `control-cert-file` should be readable by `netdata` user

For auto-detection parameters from `unbound.conf`:

- `unbound.conf` should be readable by `netdata` user
- if you have several configuration files (include feature) all of them should be readable by `netdata` user

### Configuration

#### File

The configuration file name is `go.d/unbound.conf`.

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
sudo ./edit-config go.d/unbound.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                                                        |             Default              | Required |
|:-------------------:|------------------------------------------------------------------------------------------------------------------------------------|:--------------------------------:|:--------:|
|    update_every     | Data collection frequency.                                                                                                         |                5                 |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                                                                 |                0                 |          |
|       address       | Server address in IP:PORT format.                                                                                                  |          127.0.0.1:8953          |   yes    |
|       timeout       | Connection/read/write/ssl handshake timeout.                                                                                       |                1                 |          |
|      conf_path      | Absolute path to the unbound configuration file.                                                                                   |    /etc/unbound/unbound.conf     |          |
|  cumulative_stats   | Statistics collection mode. Should have the same value as the `statistics-cumulative` parameter in the unbound configuration file. |    /etc/unbound/unbound.conf     |          |
|       use_tls       | Whether to use TLS or not.                                                                                                         |               yes                |          |
|   tls_skip_verify   | Server certificate chain and hostname validation policy. Controls whether the client performs this check.                          |               yes                |          |
|       tls_ca        | Certificate authority that client use when verifying server certificates.                                                          |                                  |          |
|      tls_cert       | Client tls certificate.                                                                                                            | /etc/unbound/unbound_control.pem |          |
|       tls_key       | Client tls key.                                                                                                                    | /etc/unbound/unbound_control.key |          |

</details>

#### Examples

##### Basic

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:8953
```

</details>

##### Unix socket

Connecting through Unix socket.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: socket
    address: /var/run/unbound.sock
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
    address: 127.0.0.1:8953

  - name: remote
    address: 203.0.113.11:8953
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `unbound` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m unbound
  ```
