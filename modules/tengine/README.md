# Tengine monitoring with Netdata

[`Tengine`](https://tengine.taobao.org/) is a web server originated by Taobao, the largest e-commerce website in Asia. It is based on the Nginx HTTP server and has many advanced features.

This module monitors one or more `Tengine` instances, depending on your configuration.

## Requirements

-   `tengine` with configured [`ngx_http_reqstat_module`](http://tengine.taobao.org/document/http_reqstat.html).
-   collector expects [default line format](http://tengine.taobao.org/document/http_reqstat.html).

## Charts

It produces the following summary charts:

-   Bandwidth in `B/s`
-   Connections in `connections/s`
-   Requests in `requests/s`
-   Requests Per Response Code Family in `requests/s`
-   Requests Per Response Code Detailed in `requests/s`
-   Number Of Requests Calling For Upstream in `requests/s`
-   Number Of Times Calling For Upstream in `calls/s`
-   Requests Per Response Code Family in `requests/s`

## Configuration

Edit the `go.d/tengine.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/tengine.conf
```

Needs only `url` to server's `/us`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/us
      
  - name: remote
    url: http://203.0.113.10/us
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/tengine.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m tengine
