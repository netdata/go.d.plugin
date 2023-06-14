# WireGuard collector

## Overview

[WireGuard](https://www.wireguard.com/) is an extremely simple yet fast and modern VPN that utilizes state-of-the-art
cryptography.

This collector monitors WireGuard VPN network interfaces and peers traffic.
CAP_NET_ADMIN capability is required. It is set automatically during installation.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### device

These metrics refer to the VPN network interface.

Labels:

| Label  | Description           |
|--------|-----------------------|
| device | VPN network interface |

Metrics:

| Metric                      |    Dimensions     | Unit  |
|-----------------------------|:-----------------:|:-----:|
| wireguard.device_network_io | receive, transmit |  B/s  |
| wireguard.device_peers      |       peers       | peers |

### peer

These metrics refer to the VPN peer.

Labels:

| Label      | Description           |
|------------|-----------------------|
| device     | VPN network interface |
| public_key | Public key of a peer  |

Metrics:

| Metric                              |    Dimensions     |  Unit   |
|-------------------------------------|:-----------------:|:-------:|
| wireguard.peer_network_io           | receive, transmit |   B/s   |
| wireguard.peer_latest_handshake_ago |       time        | seconds |

## Setup

### Prerequisites

No action required.

### Configuration

No configuration required.

## Troubleshooting

### Debug mode

To troubleshoot issues with the `wireguard` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m wireguard
  ```
