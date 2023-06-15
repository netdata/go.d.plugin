# Dnsmasq collector

## Overview

[Dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) is a lightweight, easy to configure DNS forwarder, designed to provide DNS (and optionally DHCP and TFTP) services to a small-scale network.

This collector monitors one or more Dnsmasq DNS Forwarder instances, depending on your configuration.

It collects DNS cache statistics by [reading the response on the following query](https://manpages.debian.org/stretch/dnsmasq-base/dnsmasq.8.en.html#NOTES):

```cmd
;; opcode: QUERY, status: NOERROR, id: 37862
;; flags: rd; QUERY: 7, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;cachesize.bind.   CH  TXT
;insertions.bind.  CH  TXT
;evictions.bind.   CH  TXT
;hits.bind.        CH  TXT
;misses.bind.      CH  TXT
;servers.bind.     CH  TXT
```

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

The metrics apply to the entire monitored application.

This scope has no labels.

Metrics:

|          Metric           |      Dimensions       |     Unit     |
| ------------------------- | :-------------------: | :----------: |
| dnsmasq.servers_queries   |    success, failed    |  queries/s   |
| dnsmasq.cache_performance |     hist, misses      |   events/s   |
| dnsmasq.cache_operations  | insertions, evictions | operations/s |
| dnsmasq.cache_size        |         size          |   entries    |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/dnsmasq.conf`.

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
sudo ./edit-config go.d/dnsmasq.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>All options</summary>

|        Name         |                            Description                             |    Default     | Required |
| :-----------------: | ------------------------------------------------------------------ | :------------: | :------: |
|    update_every     | Data collection frequency.                                         |       1        |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check. |       0        |          |
|       address       | Server's address. Format is 'ip_address:port'.                     | `127.0.0.1:53` |   yes    |
|      protocol       | DNS query transport protocol. Valid options: udp, tcp, tcp-tls.    |       -        |          |
|       timeout       | DNS query timeout (dial, write and read) in seconds.               |       1        |          |

</details>

#### Examples

##### Basic

An example configuration.

```yaml
jobs:
  - name: local
    address: '127.0.0.1:53'
```

##### Basic example with `protocol` option

Local server with defined DNS query transport protocol.

```yaml
jobs:
  - name: local
    address: '127.0.0.1:53'
    protocol: udp
```

##### Multi-instance

When you are defining more than one jobs, you must be careful to use different job names, to not override each other.

```yaml
jobs:
  - name: local
    address: '127.0.0.1:53'

  - name: remote
    address: '203.0.113.0:53'
```

## Troubleshooting

### Debug mode

To troubleshoot issues with the `dnsmasq` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m dnsmasq
  ```
