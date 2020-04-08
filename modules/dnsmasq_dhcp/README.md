# Dnsmasq DHCP monitoring with Netdata

[`Dnsmasq`](http://www.thekelleys.org.uk/dnsmasq/doc.html) is a lightweight, easy to configure, DNS forwarder and DHCP server.

This module monitors `Dnsmasq DHCP` leases database.

## Charts

It produces the following set of charts for every dhcp-range:

-   DHCP Range Allocated Leases in `leases`
-   DHCP Range Utilization in `percentage`

## Auto-detection

Module automatically detects all configured dhcp-ranges reading `dnsmasq` configuration files.

By default it uses:

-   `/var/lib/misc/dnsmasq.leases` to read leases.
-   `/etc/dnsmasq.conf` to detect dhcp-ranges.
-   `/etc/dnsmasq.d` to find additional configurations.

## Configuration 

Edit the `go.d/dnsmasq_dhcp.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/dnsmasq_dhcp.conf
```

Here is an example:

```yaml
jobs:
  - name         : dnsmasq_dhcp
    leases_path  : /var/lib/misc/dnsmasq.leases
    conf_path    : /etc/dnsmasq.conf
    conf_dir     : /etc/dnsmasq.d
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dnsmasq_dhcp.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m dnsmasq_dhcp
