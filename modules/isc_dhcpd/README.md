<!--
title: "Monitoring ISC dhcp lease files with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/isc_dhcpd/README.md
sidebar_label: "ISC dhcp lease files"
-->

# ISC dhcpd lease

This module monitors lease files.

Pool utilization:

-   Pool Utilization
-   Total leases
-   Active leases

## Charts

This plugin shows three charts

### Pools utilization

-   Aggregate chart for all pools. Utilization in `percentage`

### Total Leases

-   leases(Overall number of leases for all pools)

### Leases total

-   leases(number of leases for all pools).  

## Configuration

Edit the `go.d/isc_dhcpd.conf` configuration file using `edit-config` from the Agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/isc_dhcpd.conf
```

### Requirement

-   Dhcpd leases file MUST BE readable by Netdata.
-   pools MUST BE in CIDR format

### Example

Here is an example:

```yaml
jobs:
 - name: ipv4_leases_example
   leases_path: '/path/to/ipv4_lease_file'
   pools:
      - office:          '192.168.2.0/24'
      - wifi:            '192.168.3.10-192.168.3.20'
      - 192.168.4.0/24:  '192.168.4.0/24'                            
      - wifi-guest:      '192.168.5.0/24 192.168.6.10-192.168.6.20'  

 - name: ipv6_leases_example
   leases_path: '/path/to/ipv6_lease_file'
   pools:
      - office:          '192.168.2.0/24'
```

For all available options, see the ISC dhcpd collector's [configuration
file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/isc_dhcpd.conf).

## Troubleshooting

To troubleshoot issues with the ISC dhcpd collector, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` orchestrator to debug the collector:

```bash
./go.d.plugin -d -m isc_dhcpd
```
