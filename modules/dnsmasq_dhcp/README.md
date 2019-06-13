# dnsmasq_dhcp

This module will monitor [`Dnsmasq DHCP`](http://www.thekelleys.org.uk/dnsmasq/doc.html) leases database.

It produces the following set of charts for both ipv4 and ipv6 dhcp-ranges:

1. **DHCP Range Active Leases** in active leases
  * per dhcp-range


2. **DHCP Range Utilization** in percentage per
  * per dhcp-range 


### configuration

Module automatically detects all configured dhcp-ranges reading `dnsmasq` configuration files.

By default it uses:
 - `/var/lib/misc/dnsmasq.dnsmasq.leases` to read leases.
 - `/etc/dnsmasq.conf` to detect dhcp-ranges.
 - `/etc/dnsmasq.d` to find additional configurations.

Here is an example:

```yaml
jobs:
  - name         : dnsmasq_dhcp
    leases_path  : /var/lib/misc/dnsmasq.dnsmasq.leases
    conf_path    : /etc/dnsmasq.conf
    conf_dir     : /etc/dnsmasq.d
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dnsmasq_dhcp.conf).

---
