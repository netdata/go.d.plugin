<!--
title: "NGINX VTS monitoring"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/nginxvts/README.md
sidebar_label: "NGINX VTS"
-->

# NGINX VTS monitoring with Netdata

`nginxvts` can monitor statistics of NGINX which configured
with [`nginx-module-vts`](https://github.com/vozlt/nginx-module-vts), including:

- Nginx uptime (`seconds`):
    - Uptime
- Nginx connections (`requests/s`):
    - active, reading, writing, waiting, accepted, handled, total

- Shared memory size (`bytes`)
    - Maximum size of shared memory
    - Current size of shared memory
- Number of node using in shared memory (`nodes`)

- Total number of client requests (`requests/s`)
- Total Response code (`responses/s`)
    - 1xx, 2xx, 3xx, 4xx, 5xx
- Total server traffic (`bytes/s`)
    - The total number of bytes received from clients
    - The total number of bytes sent to clients
- Total server cache (`responses/s`)
    - miss, bypass, expired, stale, updating, revalidated, hit, scarce

Refer [`nginx-module-vts`](https://github.com/vozlt/nginx-module-vts#json) for more information.

`(Other statistics like UpsteamZones, FilterZones will be added later)`

## Configuration

Edit the `go.d/nginxvts.conf` configuration file using `edit-config` from the your
agent's [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/nginxvts.conf
```

Needs only `url` to server's `stub_status`. Here is an example for local and remote servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/status/format/json
  - name: remote
    url: http://203.0.113.0/status/format/json
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/nginxvts.conf).

## Troubleshooting

To troubleshoot issues with the `nginxvts` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` orchestrator to debug the collector:

```bash
./go.d.plugin -d -m nginxvts
```
