# dnsmasq_dhcp

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
