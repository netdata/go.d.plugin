# OpenVPN monitoring with Netdata

[`OpenVPN`](https://openvpn.net/) is an open-source commercial software that implements virtual private network techniques to create secure point-to-point or site-to-site connections in routed or bridged configurations and remote access facilities.

This module will monitor one or more `OpenVPN` instances via Management Interface.

## Requirements

-   `OpenVPN` with enabled [`management-interface`](https://openvpn.net/community-resources/management-interface/).

## Charts

It produces the following charts:

-   Total Number Of Active Clients in `clients`
-   Total Traffic in `kilobits/s`

Per user charts (disabled by default, see `per_user_stats` in the module config file):

-   User Traffic in `kilobits/s`
-   User Connection Time in `seconds`
 
## Configuration

This collector is disabled by default. Should be explicitly enabled in [go.d.conf](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).

Reason:
 >  Currently,the OpenVPN daemon can at most support a single management client any one time.

We disabled it to not break other tools which uses `Management Interface`.

Edit the `go.d/openvpn.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/openvpn.conf
```

Needs only `address` of OpenVPN `Management Interface`. Here is an example for 2 `OpenVPN` instances:

```yaml
jobs:
  - name: local
    address : /dev/openvpn
      
  - name: remote
    address : 203.0.113.10:7505
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/openvpn.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m openvpn
