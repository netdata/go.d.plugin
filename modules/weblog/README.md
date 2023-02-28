<!--
title: "Web server log (Apache, NGINX) monitoring with Netdata"
description: "Monitor the health and performance of Apache or Nginx logs with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/weblog/README.md"
sidebar_label: "Web server logs (Apache, NGINX, Microsoft IIS)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Quantify logs to metrics"
-->

# Web server log (Apache, NGINX) monitoring with Netdata

This module parses [`Apache`](https://httpd.apache.org/), [`NGINX`](https://nginx.org/en/) and [Microsoft IIS](https://www.iis.net/) web servers logs.

## Metrics

All metrics have "web_log." prefix.

| Metric                              |       Scope       |                 Dimensions                  |    Units     |
|-------------------------------------|:-----------------:|:-------------------------------------------:|:------------:|
| requests                            |      global       |                  requests                   |  requests/s  |
| excluded_requests                   |      global       |                  unmatched                  |  requests/s  |
| type_requests                       |      global       |        success, bad, redirect, error        |  requests/s  |
| status_code_class_responses         |      global       |           1xx, 2xx, 3xx, 4xx, 5xx           | responses/s  |
| status_code_class_1xx_responses     |      global       |       <i>a dimension per 1xx code</i>       | responses/s  |
| status_code_class_2xx_responses     |      global       |       <i>a dimension per 2xx code</i>       | responses/s  |
| status_code_class_3xx_responses     |      global       |       <i>a dimension per 3xx code</i>       | responses/s  |
| status_code_class_4xx_responses     |      global       |       <i>a dimension per 4xx code</i>       | responses/s  |
| status_code_class_5xx_responses     |      global       |       <i>a dimension per 5xx code</i>       | responses/s  |
| bandwidth                           |      global       |               received, sent                |  kilobits/s  |
| request_processing_time             |      global       |                min, max, avg                | milliseconds |
| requests_processing_time_histogram  |      global       |        <i>a dimension per bucket</i>        |  requests/s  |
| upstream_response_time              |      global       |                min, max, avg                | milliseconds |
| upstream_responses_time_histogram   |      global       |        <i>a dimension per bucket</i>        |  requests/s  |
| current_poll_uniq_clients           |      global       |                 ipv4, ipv6                  |   clients    |
| vhost_requests                      |      global       |        <i>a dimension per vhost</i>         |  requests/s  |
| port_requests                       |      global       |         <i>a dimension per port</i>         |  requests/s  |
| scheme_requests                     |      global       |                 http, https                 |  requests/s  |
| http_method_requests                |      global       |     <i>a dimension per HTTP method</i>      |  requests/s  |
| http_version_requests               |      global       |     <i>a dimension per HTTP version</i>     |  requests/s  |
| ip_proto_requests                   |      global       |                 ipv4, ipv6                  |  requests/s  |
| ssl_proto_requests                  |      global       |     <i>a dimension per SSL protocol</i>     |  requests/s  |
| ssl_cipher_suite_requests           |      global       |   <i>a dimension per SSL cipher suite</i>   |  requests/s  |
| url_pattern_requests                |      global       |     <i>a dimension per URL pattern</i>      |  requests/s  |
| custom_field_pattern_requests       |      global       | <i>a dimension per custom field pattern</i> |  requests/s  |
| custom_time_field_summary           | custom time field |                min, max, avg                | milliseconds |
| custom_time_field_histogram         | custom time field |        <i>a dimension per bucket</i>        | observations |
| url_pattern_status_code_responses   |    URL pattern    |       <i>a dimension per pattern</i>        | responses/s  |
| url_pattern_http_method_requests    |    URL pattern    |     <i>a dimension per HTTP method</i>      |  requests/s  |
| url_pattern_bandwidth               |    URL pattern    |               received, sent                |  kilobits/s  |
| url_pattern_request_processing_time |    URL pattern    |                min, max, avg                | milliseconds |

## Log Parsers

Weblog supports 4 different log parsers:

- `CSV`
- [`JSON`](https://www.json.org/json-en.html)
- [`LTSV`](http://ltsv.org/)
- `RegExp`

Try to avoid using `RegExp` because it's much slower than the other parsers. Prefer to use `LTSV` or `CSV` parser.

There is an example job for every log parser.

```yaml
jobs:
  - name: csv_parser_example
    path: /path/to/file.log
    log_type: csv
    csv_config:
      format: 'FORMAT'
      fields_per_record: -1
      delimiter: ' '
      trim_leading_space: no

  - name: json_parser_example
    path: /path/to/file.log
    log_type: json
    json_config:
      mapping:
        label1: field1
        label2: field2

  - name: ltsv_parser_example
    path: /path/to/file.log
    log_type: ltsv
    ltsv_config:
      field_delimiter: ' '
      value_delimiter: ':'
      mapping:
        label1: field1
        label2: field2

  - name: regexp_parser_example
    path: /path/to/file.log
    log_type: regexp
    regexp_config:
      pattern: 'PATTERN'
```

## Log Parser Auto-Detection

If `log_type` parameter set to `auto` (which is default), weblog will try to auto-detect appropriate log parser and log
format using the last line of the log file.

- checks if format is `CSV` (using regexp).
- checks if format is `JSON` (using regexp).
- assumes format is `CSV` and tries to find appropriate `CSV` log format using predefind list of formats. It tries to
  parse the line using each of them in the following order:

```sh
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time $upstream_response_time
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time $upstream_response_time
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time $upstream_response_time
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time $upstream_response_time
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent
```

The first one matches is used later. If you use default Apache/NGINX log format auto-detect will do for you. If it
doesn't work you need [to set format manually](#custom-log-format).

## Known Fields

These are [NGINX](http://nginx.org/en/docs/varindex.html)
and [Apache](http://httpd.apache.org/docs/current/mod/mod_log_config.html) log format variables.

Weblog is aware how to parse and interpret the fields:

| nginx                   | apache   | description                                                                              |
|-------------------------|----------|------------------------------------------------------------------------------------------|
| $host ($http_host)      | %v       | Name of the server which accepted a request.                                             |
| $server_port            | %p       | Port of the server which accepted a request.                                             |
| $scheme                 | -        | Request scheme. "http" or "https".                                                       |
| $remote_addr            | %a (%h)  | Client address.                                                                          |
| $request                | %r       | Full original request line. The line is "$request_method $request_uri $server_protocol". |
| $request_method         | %m       | Request method. Usually "GET" or "POST".                                                 |
| $request_uri            | %U       | Full original request URI.                                                               |
| $server_protocol        | %H       | Request protocol. Usually "HTTP/1.0", "HTTP/1.1", or "HTTP/2.0".                         |
| $status                 | %s (%>s) | Response status code.                                                                    |
| $request_length         | %I       | Bytes received from a client, including request and headers.                             |
| $bytes_sent             | %O       | Bytes sent to a client, including request and headers.                                   |
| $body_bytes_sent        | %B (%b)  | Bytes sent to a client, not counting the response header.                                |
| $request_time           | %D       | Request processing time.                                                                 |
| $upstream_response_time | -        | Time spent on receiving the response from the upstream server.                           |
| $ssl_protocol           | -        | Protocol of an established SSL connection.                                               |
| $ssl_cipher             | -        | String of ciphers used for an established SSL connection.                                |

In addition to that weblog understands [user defined fields](#custom-fields-feature).

Notes:

- Apache `%h` logs the IP address if [HostnameLookups](https://httpd.apache.org/docs/2.4/mod/core.html#hostnamelookups)
  is Off. The web log collector counts hostnames as IPv4 addresses. We recommend either to disable HostnameLookups or
  use `%a` instead of `%h`.
- Since httpd 2.0, unlike 1.3, the `%b` and `%B` format strings do not represent the number of bytes sent to the client,
  but simply the size in bytes of the HTTP response. It will differ, for instance, if the connection is aborted, or
  if SSL is used. The `%O` format provided by [`mod_logio`](https://httpd.apache.org/docs/2.4/mod/mod_logio.html)
  will log the actual number of bytes sent over the network.
- To get `%I` and `%O` working you need to enable `mod_logio` on Apache.
- NGINX logs URI with query parameters, Apache doesnt.
- `$request` is parsed into `$request_method`, `$request_uri` and `$server_protocol`. If you have `$request` in your log
  format, there is no sense to have others.
- Don't use both `$bytes_sent` and `$body_bytes_sent` (`%O` and `%B` or `%b`). The module does not distinguish between
  these parameters.

## Custom Log Format

Custom log format is easy. Use [known fields](#known-fields) to construct your log format.

- If using `CSV` parser

Since weblog understands NGINX and Apache variables all you need is to copy your log format and... that is it!
If there is a field that is not known by the weblog it's not a problem. It will skip it during parsing. We suggest
replace all unknown fields with `-` for optimization purposes.

Let's take as an example some non default format.

```bash
# apache
LogFormat "\"%{Referer}i\" \"%{User-agent}i\" %h %l %u %t \"%r\" %>s %b" custom

# nginx
log_format custom '"$http_referer" "$http_user_agent" '
                  '$remote_addr - $remote_user [$time_local] '
                  '"$request" $status $body_bytes_sent'
```

To get it working we need to copy the format without any changes (make it a line for nginx). Replacing unknown fields is
optional but recommended.

Special case:

Both `%t` and `$time_local` fields represent time
in [Common Log Format](https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format). It is a special case
because it's in fact 2 fields after csv parse (ex.: `[22/Mar/2009:09:30:31 +0100]`). Weblog understands it, and you
don't
need to replace it with `-` (if we want to do it we need to make it `- -`).

```yaml
jobs:
  - name: apache_csv_custom_format_example
    path: /path/to/file.log
    log_type: csv
    csv_config:
      format: '- - %h - - %t \"%r\" %>s %b'

  - name: nginx_csv_custom_format_example
    path: /path/to/file.log
    log_type: csv
    csv_config:
      format: '- - $remote_addr - - [$time_local] "$request" $status $body_bytes_sent'
```

- If using `JSON` parser

Provide fields [mapping](#known-fields) if needed. Don't use `$` and `%` prefixes for mapped field names. They are only
needed in `CSV` format.

- If using `LTSV` parser

Provide fields [mapping](#known-fields) if needed. Don't use `$` and `%` prefixes for mapped field names. They are only
needed in `CSV` format.

- If using `RegExp` parser

Use pattern with subexpressions names. These names should be known by weblog.

## Custom Fields Feature

Weblog is able to extract user defined fields and count patterns matches against these fields.

This feature needs:

- custom log format with user defined fields
- list of patterns to match against appropriate fields

Pattern syntax: [matcher](https://github.com/netdata/go.d.plugin/tree/master/pkg/matcher#supported-format).

There is an example with 2 custom fields - `$http_referer` and `$http_user_agent`. Weblog is unaware of these fields,
but we still can get some info from them.

```yaml
  - name: nginx_csv_custom_fields_example
    path: /path/to/file.log
    log_type: csv
    csv_config:
      format: '- - $remote_addr - - [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"'
    custom_fields:
      - name: http_referer     # same name as in 'format' without $
        patterns:
          - name: cacti
            match: '~ cacti'
          - name: observium
            match: '~ observium'
      - name: http_user_agent  # same name as in 'format' without $
        patterns:
          - name: android
            match: '~ Android'
          - name: iphone
            match: '~ iPhone'
          - name: other
            match: '* *'
```

## Custom time fields feature

The web log collector is also able to extract user defined time fields and could count min/avg/max + histogram against
these fields.

This feature needs:

- A custom log format with user-defined time fields.
- A histogram to show response time in seconds, which is optional.

As an example, Apache [`mod_logio`](https://httpd.apache.org/docs/2.4/mod/mod_logio.html) adds a `^FB` logging
directive. This value shows a delay in microseconds between when the request arrived, and the first byte of the response
headers are written.

As with the custom fields feature, Netdata's web log collector is unaware of these fields, but we can still get some
info from them.

```yaml
  - name: apache_csv_custom_fields_example
    path: /path/to/file.log
    log_type: csv
    csv_config:
      format: '%v %a %p %m %H \"%U\" %t %>s %O %I %D %^FB \"%{Referer}i\" \"%{User-Agent}i\" \"%r\"'
    custom_time_fields:
      - name: '^FB'
        histogram: [.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10] # optional field
```

## Configuration

Edit the `go.d/web_log.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/web_log.conf
```

This module needs only `path` to log file. If it fails to auto-detect your log format you
need [to set it manually](#custom-log-format).

```yaml
jobs:
  - name: nginx
    path: /var/log/nginx/access.log

  - name: apache
    path: /var/log/apache2/access.log
    log_type: csv
    csv_config:
      format: '- - %h - - %t \"%r\" %>s %b'

  - name: iis
    path: /mnt/c/inetpub/logs/LogFiles/W3SVC1/u_ex*.log
    log_type: csv
    csv_config:
      format: '- - $host $request_method $request_uri - $server_port - $remote_addr - - $status - - $request_time'
```

For all available options, please see the
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/web_log.conf).

## Troubleshooting

To troubleshoot issues with the `web_log` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m web_log
  ```
