# dnsmasq_dhcp

This module will monitor [`Dnsmasq DHCP`](http://www.thekelleys.org.uk/dnsmasq/doc.html) leases database.

It produces the following set of charts for every dhcp-range:

1. **DHCP Range Allocated Leases** in leases

2. **DHCP Range Utilization** in percentage

### configuration

Module automatically detects all configured dhcp-ranges reading `dnsmasq` configuration files.

By default it uses:
 - `/var/lib/misc/dnsmasq.leases` to read leases.
 - `/etc/dnsmasq.conf` to detect dhcp-ranges.
 - `/etc/dnsmasq.d` to find additional configurations.

Here is an example:

```yaml
jobs:
  - name         : dnsmasq_dhcp
    leases_path  : /var/lib/misc/dnsmasq.leases
    conf_path    : /etc/dnsmasq.conf
    conf_dir     : /etc/dnsmasq.d
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dnsmasq_dhcp.conf).

---
