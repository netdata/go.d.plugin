# Lighttpd2 monitoring with Netdata

[`Lighttpd2`](https://redmine.lighttpd.net/projects/lighttpd2) is a work in progress version of open-source web server.

This module will monitor one or more `Lighttpd2` servers, depending on your configuration.

## Requirements

-   `lighttpd2` with enabled [`mod_status`](https://doc.lighttpd.net/lighttpd2/mod_status.html)

## Charts

It produces the following charts:

-   Requests in `requests/s`
-   Status Codes in `requests/s`
-   Traffic in `kilobits/s`
-   Connections in `connections`
-   Connection States in  `connection`
-   Memory Usage in `KiB`
-   Uptime in `seconds`

## Configuration

Edit the `go.d/lighttpd2.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/lighttpd2.conf
```

Needs only `url` to server's `server-status?format=plain`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1/server-status?format=plain
      
  - name: remote
    url : http://203.0.113.10/server-status?format=plain
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/lighttpd2.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m lighttpd2

