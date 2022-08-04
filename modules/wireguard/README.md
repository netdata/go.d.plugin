<!--
title: "WireGuard monitoring with Netdata"
description: "Monitor WireGuard VPN network interfaces and peers traffic."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/wireguard/README.md
sidebar_label: "WireGuard"
-->

# WireGuard monitoring with Netdata

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

First, navigate to your plugins' directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m wireguard
```
