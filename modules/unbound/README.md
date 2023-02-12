<!--
title: "Unbound monitoring with Netdata"
description: "Monitor the health and performance of Unbound DNS resolvers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/unbound/README.md"
sidebar_label: "Unbound"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Networking"
-->

# Unbound monitoring with Netdata

[`Unbound`](https://nlnetlabs.nl/projects/unbound/about/) is a validating, recursive, and caching DNS resolver product
from NLnet Labs.

This module monitors one or more `Unbound` servers, depending on your configuration.

## Requirements

- `Unbound` with enabled `remote-control` interface (
  see [unbound.conf](https://nlnetlabs.nl/documentation/unbound/unbound.conf))

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

## Metrics

All metrics have "vcsa." prefix.

| Metric                            | Scope  |                       Dimensions                       |    Units     |
|-----------------------------------|:------:|:------------------------------------------------------:|:------------:|
| queries                           | global |                        queries                         |   queries    |
| queries_ip_ratelimited            | global |                      ratelimited                       |   queries    |
| dnscrypt_queries                  | global |          crypted, cert, cleartext, malformed           |   queries    |
| cache                             | global |                       hits, miss                       |    events    |
| cache_percentage                  | global |                       hits, miss                       |  percentage  |
| prefetch                          | global |                       prefetches                       |  prefetches  |
| expired                           | global |                        expired                         |   replies    |
| zero_ttl_replies                  | global |                        zero_ttl                        |   replies    |
| recursive_replies                 | global |                       recursive                        |   replies    |
| recursion_time                    | global |                      avg, median                       | milliseconds |
| request_list_usage                | global |                        avg, max                        |   queries    |
| current_request_list_usage        | global |                       all, users                       |   queries    |
| request_list_jostle_list          | global |                  overwritten, dropped                  |   queries    |
| tcpusage                          | global |                         usage                          |   buffers    |
| uptime                            | global |                          time                          |   seconds    |
| thread_cache                      | thread |                       hits, miss                       |    events    |
| thread_cache_percentage           | thread |                       hits, miss                       |  percentage  |
| thread_prefetch                   | thread |                       prefetches                       |  prefetches  |
| thread_expired                    | thread |                        expired                         |   replies    |
| thread_zero_ttl_replies           | thread |                        zero_ttl                        |   replies    |
| thread_recursive_replies          | thread |                       recursive                        |   replies    |
| thread_recursion_time             | thread |                      avg, median                       | milliseconds |
| thread_request_list_usage         | thread |                        avg, max                        |   queries    |
| thread_current_request_list_usage | thread |                       all, users                       |   queries    |
| thread_request_list_jostle_list   | thread |                  overwritten, dropped                  |   queries    |
| thread_tcpusage                   | thread |                         usage                          |   buffers    |
| cache_memory                      | global | message, rrset, dnscrypt_nonce, dnscrypt_shared_secret |      KB      |
| mod_memory                        | global |       iterator, respip, validator, subnet, ipsec       |      KB      |
| mem_streamwait                    | global |                       streamwait                       |      KB      |
| cache_count                       | global | infra, key, msg, rrset, dnscrypt_nonce, shared_secret  |    items     |
| type_queries                      | global |           <i>a dimension per query type</i>            |   queries    |
| class_queries                     | global |           <i>a dimension per query class</i>           |   queries    |
| opcode_queries                    | global |          <i>a dimension per query opcode</i>           |   queries    |
| flag_queries                      | global |             qr, aa, tc, rd, ra, z, ad, cd              |   queries    |
| rcode_answers                     | global |           <i>a dimension per reply rcode</i>           |   replies    |
| thread_queries                    | global |                        queries                         |   queries    |
| thread_queries_ip_ratelimited     | global |                      ratelimited                       |   queries    |
| thread_dnscrypt_queries           | global |          crypted, cert, cleartext, malformed           |   queries    |

## Configuration

Edit the `go.d/unbound.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/unbound.conf
```

This Unbound collector only needs the `address` to a server's `remote-control` interface if TLS is disabled or `address`
of unix socket. Otherwise, you need to set path to the `control-key-file` and `control-cert-file` files.

The module tries to auto-detect following parameters reading `unbound.conf`:

- address
- cumulative_stats
- use_tls
- tls_cert
- tls_key

Module supports both cumulative and non-cumulative modes. Default is non-cumulative. If your server has enabled
`statistics-cumulative`, but the module fails to auto-detect it (`unbound.conf` is not readable, or it is a remote
server), you need to set it manually in the configuration file.

Here is an example for several servers:

```yaml
jobs:
  - name: local
    address: 127.0.0.1:8953
    use_tls: yes
    tls_skip_verify: yes
    tls_cert: /etc/unbound/unbound_control.pem
    tls_key: /etc/unbound/unbound_control.key

  - name: remote
    address: 203.0.113.10:8953
    use_tls: no

  - name: remote_cumulative
    address: 203.0.113.11:8953
    use_tls: no
    cumulative_stats: yes

  - name: socket
    address: /var/run/unbound.sock
```

For all available options, please see the
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/unbound.conf).

## Troubleshooting

Ensure that the control protocol is actually configured correctly. Run following command as `root` user:

```bash
unbound-control stats_noreset
```

It should print out a bunch of info about the internal statistics of the server. If this returns an error, you don't
have the control protocol set up correctly.

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
