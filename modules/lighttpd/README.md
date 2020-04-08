# Lighttpd monitoring with Netdata

[`Lighttpd`](https://www.lighttpd.net/) is an open-source web server optimized for speed-critical environments while remaining standards-compliant, secure and flexible

This module will monitor one or more `Lighttpd` servers, depending on your configuration.

## Requirements

-   `lighttpd` with enabled [`mod_status`](https://redmine.lighttpd.net/projects/lighttpd/wiki/Docs_ModStatus).

## Charts

It produces the following charts:

-   Requests in `requests/s`
-   Bandwidth in `kilobytes/s`
-   Servers in `servers`
-   Scoreboard in `connections`
-   Uptime in `seconds`

## Configuration

Edit the `go.d/lighttpd.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/lighttpd.conf
```

Needs only `url` to server's `server-status?auto`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1/server-status?auto
      
  - name: remote
    url : http://203.0.113.10/server-status?auto
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/lighttpd.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m lighttpd
