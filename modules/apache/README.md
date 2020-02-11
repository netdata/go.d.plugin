# Apache monitoring with Netdata

[`Apache`](https://httpd.apache.org/) is an open-source HTTP server for modern operating systems including UNIX and Windows.

This module will monitor one or more `Apache` servers, depending on your configuration.

## Requirements

-   `Apache` with enabled [`mod_status`](https://httpd.apache.org/docs/2.4/mod/mod_status.html)

## Charts

It produces the following charts:

-   Requests in `requests/s`
-   Connections in `connections`
-   Async Connections in `connections`
-   Scoreboard in `connections`
-   Bandwidth in `kilobits/s`
-   Workers in `workers`
-   Lifetime Average Number Of Requests Per Second in `requests/s`
-   Lifetime Average Number Of Bytes Served Per Second in `KiB/s`
-   Lifetime Average Response Size in `KiB`

## Configuration

Edit the `go.d/apache.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/apache.conf
```

Needs only `url` to server's `server-status?auto`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/server-status?auto
      
  - name: remote
    url: http://203.0.113.10/server-status?auto
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/apache.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m apache
