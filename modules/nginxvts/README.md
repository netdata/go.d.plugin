<!--
title: "NGINX Vts monitoring"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/nginxvts/README.md
sidebar_label: "NGINX VTS"
-->

# NGINX VTS monitoring with Netdata

[`NGINX`](https://www.nginx.com/) is a web server which can also be used as a reverse proxy, load balancer, mail proxy and HTTP cache. 

This module will monitor one or more `NGINX` servers, depending on your configuration.

## Requirements

 -   `NGINX` with configured [`nginx-module-vts`](https://github.com/vozlt/nginx-module-vts).

## Charts

Nginx VTS will produce following charts:
- Main charts:
  - Nginx running time (`milliseconds`): 
    - Starting time
    - Up time
  - Nginx connections (`number`):
    - active,	reading, writing, waiting, accepted, handled,	requests

- SharedZones
  - Shared memory size (`bytes`)
    - Maximum size of shared memory
    - Current size of shared memory
  - Number of node using in shared memory (`number`)

- ServerZones charts
  - Total number of client requests (`number`)
  - Response code (`number`)
    - 1xx, 2xx, 3xx, 4xx, 5xx
  - IO (`bytes`)
    - The total number of bytes received from clients
    - The total number of bytes sent to clients
  - cache (`number`)
    - miss, bypass, expired, stale, updating, revalidated, hit, scarce

- UpstreamZones charts
  - Total number of client connections forwarded to this server (`number`)
  - Response code (`number`)
    - 1xx, 2xx, 3xx, 4xx, 5xx
  - IO (`bytes`)
    - The total number of bytes received from clients
    - The total number of bytes sent to clients

- FilterZones charts
  - Total number of client requests (`number`)
  - Response code (`number`)
    - 1xx, 2xx, 3xx, 4xx, 5xx
  - IO (`bytes`)
    - The total number of bytes received from clients
    - The total number of bytes sent to clients
  - cache (`number`)
    - miss, bypass, expired, stale, updating, revalidated, hit, scarce
  
- CacheZones charts
  - Cache size (`bytes`)
  - Shared memory size (`bytes`)
    - Maximum size of cache
    - Current size of cache
  - IO (`bytes`)
    - The total number of bytes received from cache
    - The total number of bytes sent to cache
  - Cache responses (`number`)
    - miss, bypass, expired, stale, updating, revalidated, hit, scarce

Refer [`nginx-module-vts`](https://github.com/vozlt/nginx-module-vts#json) for more information.
  

## Configuration

Edit the `go.d/nginxvts.conf` configuration file using `edit-config` from the your agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/nginxvts.conf
```

Needs only `url` to server's `stub_status`. Here is an example for local and remote servers:

```yaml
jobs:
  - name: local
    url: http://192.168.66.6/status/format/json
  - name: remote
    url: http://8.8.8.8/status/format/json
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/nginxvts.conf).


## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m nginxvts
