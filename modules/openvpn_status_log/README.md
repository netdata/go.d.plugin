<!--
title: "OpenVPN monitoring with Netdata(based on status log)"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/openvpn_status/README.md
sidebar_label: "OpenVPN(StatusLog)"
-->

# OpenVPN monitoring with Netdata

Parses server log files and provides summary (client, traffic) metrics. Please *note* that this collector is similar to
another OpenVPN collector. However, this collector requires status logs from OpenVPN in contrast to another collector
that requires Management Interface enabled.

## Requirements

- make sure `netdata` user can read `openvpn-status.log`.

- `update_every` interval must match the interval on which OpenVPN writes the operational status to the log file.

## Charts

It produces the following charts:

- Active Clients in `clients`
- Traffic in `kilobits/s`

Per user charts (disabled by default, see `per_user_stats` in the module config file):

- User Traffic in `kilobits/s`
- User Connection Time in `seconds`

## Configuration

Edit the `go.d/openvpn_status_log.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata   # Replace this path with your Netdata config directory, if different
sudo ./edit-config go.d/openvpn_status_log.conf
```

Configuration example:

```yaml
jobs:
  - name: local
    log_path: '/var/log/openvpn/status.log'
    per_user_stats:
      includes:
        - "* *"
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/openvpn_status_log.conf).

---

## Troubleshooting

To troubleshoot issues with the `openvpn_status_log` collector, run the `go.d.plugin` with the debug option enabled. The
output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m openvpn_status_log
```
