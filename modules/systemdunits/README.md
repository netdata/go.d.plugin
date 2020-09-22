<!--
title: "Systemd units monitoring with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/systemdunits/README.md
sidebar_label: "systemdunits"
-->

# Systemd states monitoring with Netdata

[`Systemd`](https://www.freedesktop.org/wiki/Software/systemd/) is a suite of basic building blocks for a Linux system.

This module will monitor `Systemd` unit states.



## Charts

It produces the following charts:

-   Services in `service_states`
-   Sockets in `socket_states`
-   Targets in `target_states`
-   Paths in `path_states`
-   Devices in `device_states`
-   Mounts in `mount_states`
-   Autimounts in `automount_states`
-   Swaps in `swap_states`
-   Timers in `timer_states`
-   Scopes in `scope_states`

## Configuration

Edit the `go.d/systemdunits.conf` configuration file using `edit-config` from the your agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/systemdunits.conf
```


```yaml
jobs:
  - name: systemd-service-units
    selector:
      includes:
         - '* *.service'

  - name: systemd-socket-units
    selector:
      includes:
         - '* *.socket'
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/systemdunits.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m systemdunits
