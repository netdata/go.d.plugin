# Wireguard

[`Wireguard`](https://www.wireguard.com/) is an modern open-source VPN that utilizes **state-of-the-art** cryptograph.

This module will monitor one or more vpn network interface and wireguard peers, depending on your configuration.

## Requirements

- `Wireguard` enabled
- Set network capability to golang collector (default file is: `/usr/libexec/netdata/plugins.d/go.d.plugin`)

E.g:
```bash
$ sudo setcap CAP_NET_ADMIN+ep /usr/libexec/netdata/plugins.d/go.d.plugin
```

## Charts

It produces the following charts:

#### Total data
- total received data (KB) by peers on wireguard interface in `received`
- total sent data (KB) by peers on wireguard interface in `sent`

#### Peer's bandwitch
- badwidth received data (Kb/s) in `received`
- total sent data (Kb/s) in `sent`


## Configuration

Edit the `go.d/wireguard.conf` configuration file using `edit-config` from the your agent's [config directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/wireguard.conf
```

Define your Wireguard interfaces:
```yaml
jobs:
  - interface: wg0

  - interface: wg1
```

If you do not define any interface, wireguard's collector will define `wg0` as **default**.

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m wireguard
