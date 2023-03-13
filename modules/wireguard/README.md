<!--
title: "WireGuard monitoring with Netdata"
description: "Monitor WireGuard VPN network interfaces and peers traffic."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/wireguard/README.md"
sidebar_label: "WireGuard"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Networking"
-->

# WireGuard collector

[WireGuard](https://www.wireguard.com/) is an extremely simple yet fast and modern VPN that utilizes state-of-the-art
cryptography.

This module monitors WireGuard VPN network interfaces and peers traffic.

## Requirements

- Grant `CAP_NET_ADMIN` capability to `go.d.plugin`.

  ```bash
  sudo setcap CAP_NET_ADMIN+epi <INSTALL_PREFIX>/usr/libexec/netdata/plugins.d/go.d.plugin
  ```

## Metrics

All metrics have "wireguard." prefix.

Labels per scope:

- device: device.
- peer: device, public_key.

| Metric                    | Scope  |    Dimensions     |  Units  |
|---------------------------|:------:|:-----------------:|:-------:|
| device_peers              | device |       peers       |  peers  |
| device_network_io         | device | receive, transmit |   B/s   |
| peer_network_io           |  peer  | receive, transmit |   B/s   |
| peer_latest_handshake_ago |  peer  |       time        | seconds |

## Configuration

No configuration needed.

## Troubleshooting

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
