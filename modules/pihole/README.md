# Pi-hole monitoring with Netdata

[`Pi-hole`](https://pi-hole.net) is a Linux network-level advertisement and Internet tracker blocking application which acts as a DNS sinkhole, intended for use on a private network.

This module will monitor one or more `Pi-hole` instances using [PHP API](https://github.com/pi-hole/AdminLTE).

The API exposed data time frame is `for the last 24 hr`. All collected values are for that time time frame, not for the module collection interval.

## Charts 

It produces the following set of charts:

-   DNS Queries Total (Cached, Blocked and Forwarded) in `queries`
-   DNS Queries in `queries`
-   DNS Queries Percentage in `percentage`  
-   Unique Clients in `clients`
-   Domains On Blocklist in `domains`
-   Blocklist Last Update in `seconds`
-   Unwanted Domains Blocking Status in `boolean`
 
If the web password is set and valid following charts will be added:

-   DNS Queries Per Type in `percentage`
-   DNS Queries Per Destination in `percentage`
-   Top Clients in `requests`
-   Top Permitted Domains in `hits`
-   Top Blocked Domains in `hits`

## Configuration

Edit the `go.d/pihole.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/pihole.conf
```

Module automatically detects `Pihole` web password reading `setupVars.conf` file. It expects to find the file in the `/etc/pihole/` directory.

If you want to monitor remote instance you need to set the password in the module configuration file. 

Here is an example for local and remote instances:

```yaml
jobs:
  - name: local
    top_clients_entries: 10
    top_items_entries: 10  # top permitted and top blocked domains charts
    
  - name: remote
    url: http://203.0.113.10
    password: 1ebd33f882f9aa5fac26a7cb74704742f91100228eb322e41b7bd6e6aeb8f74b
    
  - name: remote_https
    url: https://203.0.113.11
    password: 1ebd33f882f9aa5fac26a7cb74704742f91100228eb322e41b7bd6e6aeb8f74b
    tls_skip_verify: yes  # self signed certificate verification skip
    
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/pihole.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m pihole
