# phpfpm

[`PHP-FPM`](https://php-fpm.org/) is an alternative PHP FastCGI implementation with some additional features useful for sites of any size, especially busier sites.

This module will monitor one or more `php-fpm` instances depending on configuration.

## Requirements

-   `php-fpm` with enabled `status` page
-   access to `status` page via web server

## Charts

It produces following charts:

-   Active Connections in `connections`
-   Requests in `requests/s`
-   Performance in `status`
-   Request Duration in `milliseconds`
-   Request CPU in `percentage`
-   Request Memory in `KB`

## Configuration

Needs only `url` to server's `status`. Here is an example for local server:

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
