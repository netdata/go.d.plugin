# PHP-FPM monitoring with Netdata

[`PHP-FPM`](https://php-fpm.org/) is an alternative PHP FastCGI implementation with some additional features useful for sites of any size, especially busier sites.

This module will monitor one or more `php-fpm` instances, depending on your configuration.

## Requirements

-   `php-fpm` with enabled `status` page
-   access to `status` page via web server

## Charts

It produces following charts:

-   Active Connections in `connections`
-   Requests in `requests/s`
-   Performance in `status`
-   Requests Duration Among All Idle Processes in `milliseconds`
-   Last Request CPU Usage Among All Idle Processes in `percentage`
-   Last Request Memory Usage Among All Idle Processes in `KB`

## Configuration

Edit the `go.d/phpfpm.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/phpfpm.conf
```

Needs only `url` to server's `status`. Here is an example for local server an remote servers:

```yaml
jobs:
  - name: local
    url: http://localhost/status?full&json

  - name: local
    url: http://[::1]/status?full&json

  - name: remote
    url: http://203.0.113.10/status?full&json
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/phpfpm.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m phpfpm
